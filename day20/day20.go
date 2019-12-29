package main

import (
	"bufio"
	"fmt"
	"github.com/albertorestifo/dijkstra"
	//"gonum.org/v1/gonum/graph"
	//"gonum.org/v1/gonum/graph/traverse"

	//"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"os"
	"unicode"
)

type Portal struct {
	name rune
	nodeId int64
	depth int
	outer bool
}

type Coord struct {
	x int
	y int
}

type Grid [][]rune
type Items map[rune]Coord

type KeysHeld map[rune]bool
func (k KeysHeld) allTrue() bool {
	for _, v := range k {
		if !v {
			return false
		}
	}
	return true
}

var nodeIdToCoord = make(map[int64]Coord)

func partOne(startingPos Coord, grid Grid, keys Items, doors Items) {
	visited := make(Grid, len(grid))
	for i := range grid {
		visited[i] = make([]rune, len(grid[i]))
		copy(visited[i], grid[i])
	}
}

func getNeighbours(pos Coord, grid Grid) (neighbours []Coord) {
	if pos.x > 0 && isHallway(grid[pos.y][pos.x-1]) {
		neighbours = append(neighbours, Coord{pos.x-1,pos.y})
	}
	if pos.x < len(grid[pos.y])-1 && isHallway(grid[pos.y][pos.x+1]) {
		neighbours = append(neighbours, Coord{pos.x+1,pos.y})
	}
	if pos.y > 0 && isHallway(grid[pos.y-1][pos.x]) {
		neighbours = append(neighbours, Coord{pos.x,pos.y-1})
	}
	if pos.y < len(grid)-1 && isHallway(grid[pos.y+1][pos.x]) {
		neighbours = append(neighbours, Coord{pos.x,pos.y+1})
	}
	return neighbours
}

func isPortal(r rune) bool {
	return unicode.IsNumber(r) || unicode.IsLetter(r)
}

func isHallway(r rune) bool {
	return r == '.' || isPortal(r)
}

func findNode(name rune, depth int, outer bool, portals []Portal) Portal {
	for _, v := range portals {
		if v.name == name && v.depth == depth && v.outer == outer {
			return v
		}
	}
	return Portal{}
}

func buildRecursiveGraph2(g dijkstra.Graph, grid Grid, portals []Portal, depth int, maxDepth int) (dijkstra.Graph, []Portal) {

	runeToId := make(map[rune]int64)
	nodeCoordToId := make(map[Coord]int64)
	portalToNodeId := make(map[rune][]int64)

	if depth > maxDepth {
		return g, portals
	}

	// create nodes, flag portals
	nodeCount := int64(depth*100000)
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			tmpCoord := Coord{x, y}
			if isHallway(grid[y][x]) {
				nodeKey :=fmt.Sprintf("%d-%d-%d",nodeCount,x,y)
				g[nodeKey] = make(map[string]int, 0)
				if _, ok := nodeCoordToId[tmpCoord]; !ok {
					nodeCoordToId[tmpCoord] = nodeCount
					nodeIdToCoord[nodeCount] = tmpCoord
				}
			}
			if isPortal(grid[y][x]) {
				isOuter := outer(tmpCoord, grid)
				portals = append(portals, Portal{grid[y][x], nodeCount, depth , isOuter})
				if _, ok := portalToNodeId[grid[y][x]]; !ok {
					portalToNodeId[grid[y][x]] = []int64{nodeCount}
				} else {
					portalToNodeId[grid[y][x]] = append(portalToNodeId[grid[y][x]], nodeCount)
				}
				runeToId[grid[y][x]] = nodeCount
			}
			nodeCount += 1
		}
	}

	// add edges between neighbours and portals
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			gVal := grid[y][x]
			tmpCoord := Coord{x, y}
			if nodeId, ok := nodeCoordToId[tmpCoord]; ok {
				neighs := getNeighbours(Coord{x, y}, grid)
				for _, neigh := range neighs {
					neighNodeId := nodeCoordToId[neigh]
					fromKey := fmt.Sprintf("%d-%d-%d",nodeId,x,y)
					toKey := fmt.Sprintf("%d-%d-%d",neighNodeId,neigh.x,neigh.y)
					g[fromKey][toKey] = 1
					g[toKey][fromKey] = 1
				}
				// join portals
				if isPortal(gVal) {
					// link to portal of higher depth
					if outer(tmpCoord, grid) && depth > 1 {
						fromPortal := findNode(gVal, depth, true, portals)
						toPortal := findNode(gVal, depth-1, false, portals)
						if toPortal.nodeId > 0 {
							fromKey := fmt.Sprintf("%d-%d-%d",fromPortal.nodeId,nodeIdToCoord[fromPortal.nodeId].x, nodeIdToCoord[fromPortal.nodeId].y)
							toKey := fmt.Sprintf("%d-%d-%d",toPortal.nodeId,nodeIdToCoord[toPortal.nodeId].x, nodeIdToCoord[toPortal.nodeId].y)
							g[fromKey][toKey] = 1
							g[toKey][fromKey] = 1
						}
					}
				}
			}
		}
	}

	return buildRecursiveGraph2(g, grid, portals, depth+1, maxDepth)
}

func makeRecursiveGraph2(grid Grid)  (dijkstra.Graph, []Portal) {
	g := dijkstra.Graph{}
	p := make([]Portal, 0)
	return buildRecursiveGraph2(g, grid, p, 1, 100)
}

func makeGraph(grid Grid) (*simple.WeightedUndirectedGraph, map[Coord]int64, map[int64]Coord, map[rune]int64) {

	g := simple.NewWeightedUndirectedGraph(0, 0)

	runeToId := make(map[rune]int64)
	nodeCoordToId := make(map[Coord]int64)
	//nodeIdToCoord := make(map[int64]Coord)
	portalToNodeId := make(map[rune][]int64)

	// create nodes, flag portals
	nodeCount := int64(0)
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			tmpCoord := Coord{x, y}
			if isHallway(grid[y][x]) {
				g.AddNode(simple.Node(nodeCount))
				if _, ok := nodeCoordToId[tmpCoord]; !ok {
					nodeCoordToId[tmpCoord] = nodeCount
					nodeIdToCoord[nodeCount] = tmpCoord
				}
			}
			if isPortal(grid[y][x]) {
				if _, ok := portalToNodeId[grid[y][x]]; !ok {
					portalToNodeId[grid[y][x]] = []int64{nodeCount}
				} else {
					portalToNodeId[grid[y][x]] = append(portalToNodeId[grid[y][x]], nodeCount)
				}
				runeToId[grid[y][x]] = nodeCount
			}
			nodeCount += 1
		}
	}

	// add edges between neighbours and portals
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			gVal := grid[y][x]
			tmpCoord := Coord{x, y}
			if nodeId, ok := nodeCoordToId[tmpCoord]; ok {
				neighs := getNeighbours(Coord{x, y}, grid)
				for _, neigh := range neighs {
					neighNodeId := nodeCoordToId[neigh]
					g.SetWeightedEdge(simple.WeightedEdge{
						F: simple.Node(nodeId),
						T: simple.Node(neighNodeId),
						W: 1,
					})
				}
				// join portals
				if isPortal(gVal) && len(portalToNodeId[gVal]) > 1 {
					g.SetWeightedEdge(simple.WeightedEdge {
						F: simple.Node(portalToNodeId[gVal][0]),
						T: simple.Node(portalToNodeId[gVal][1]),
						W: 1,
					})
				}
			}
		}
	}

	return g, nodeCoordToId, nodeIdToCoord, runeToId
}

func outer(c Coord, grid Grid) bool {
	return (c.y == 0 || c.x == 0 || c.y == len(grid)-1 || c.x == len(grid[0])-1)
}

func inner(c Coord, grid Grid) bool {
	return !outer(c, grid)
}

func (g Grid) isDoor(c Coord) bool {
	return unicode.IsUpper(g[c.y][c.x])
}

func (g Grid) isKey(c Coord) bool {
	return unicode.IsLower(g[c.y][c.x])
}

func main() {

	file, _ := os.Open("input.txt")
	//file, _ := os.Open("test.txt")
	s := bufio.NewScanner(file)

	grid := make(Grid, 0)

	yPos := 0
	for s.Scan() {
		line := s.Text()
		grid = append(grid, make([]rune, len(line)))

		for xPos, v := range line {
			grid[yPos][xPos] = v
		}
		yPos++
	}

	g, portals := makeRecursiveGraph2(grid)

	fromId2 := findNode('5', 1, true, portals)
	toId2 := findNode('W', 1, true, portals)
	fmt.Printf("Finding between %d and %d\n", fromId2.nodeId, toId2.nodeId)

	fromKey := fmt.Sprintf("%d-%d-%d", fromId2.nodeId, nodeIdToCoord[fromId2.nodeId].x, nodeIdToCoord[fromId2.nodeId].y)
	toKey := fmt.Sprintf("%d-%d-%d", toId2.nodeId, nodeIdToCoord[toId2.nodeId].x, nodeIdToCoord[toId2.nodeId].y)
	path, _, _ := g.Path(fromKey,toKey)
	fmt.Println(len(path)-1)

}
