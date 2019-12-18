package main

import (
	"bufio"
	"fmt"
	"gonum.org/v1/gonum/graph/simple"
	"os"
	"unicode"
)

//type Node struct {
//	pos Coord
//	tile rune
//	next []*Node
//}
//
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

func partOne(startingPos Coord, grid Grid, keys Items, doors Items) {
	visited := make(Grid, len(grid))
	for i := range grid {
		visited[i] = make([]rune, len(grid[i]))
		copy(visited[i], grid[i])
	}
}

func getNeighbours(pos Coord, grid Grid) (neighbours []Coord) {
	if grid[pos.y][pos.x-1] == '.'  {
		neighbours = append(neighbours, Coord{pos.x-1,pos.y})
	}
	if grid[pos.y][pos.x+1] == '.'  {
		neighbours = append(neighbours, Coord{pos.x+1,pos.y})
	}
	if grid[pos.y-1][pos.x] == '.' {
		neighbours = append(neighbours, Coord{pos.x,pos.y-1})
	}
	if grid[pos.y+1][pos.x] == '.' {
		neighbours = append(neighbours, Coord{pos.x,pos.y+1})
	}
	return neighbours
}

func makeGraph(grid Grid) (*simple.UndirectedGraph, map[Coord]int64) {

	g := simple.NewUndirectedGraph()

	nodeCoordToId := make(map[Coord]int64)

	nodeCount := int64(0)
	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[y])-1; x++ {
			tmpCoord := Coord{x, y}
			g.AddNode(simple.Node(nodeCount))
			if _, ok := nodeCoordToId[tmpCoord]; !ok {
				nodeCoordToId[tmpCoord] = nodeCount
			}
			nodeCount += 1
		}
	}

	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[y])-1; x++ {
			tmpCoord := Coord{x, y}
			nodeId := nodeCoordToId[tmpCoord]
			neighs := getNeighbours(Coord{x, y}, grid)
			for _, neigh := range neighs {
				neighNodeId := nodeCoordToId[neigh]
				g.SetEdge(simple.Edge{
					F: g.Node(nodeId),
					T: g.Node(neighNodeId),
				})
			}
		}
	}

	return g, nodeCoordToId
}

func main() {

	file, _ := os.Open("input.txt")
	s := bufio.NewScanner(file)

	grid := make(Grid, 0)

	keys := make(Items, 0)
	doors := make(Items, 0)

	yPos := 0
	startingPos := Coord{}
	for s.Scan() {
		line := s.Text()
		grid = append(grid, make([]rune, len(line)))

		for xPos, v := range line {
			grid[yPos][xPos] = v

			if unicode.IsUpper(v) {
				// found a door
				doors[unicode.ToLower(v)] = Coord{xPos, yPos }
			} else if unicode.IsLower(v) {
				// found a key
				keys[v] = Coord{xPos, yPos }
			} else if v == '@' {
				startingPos = Coord{xPos, yPos }
			}
		}
		yPos++
	}

	_, nodemap := makeGraph(grid)
	startingNode := nodemap[startingPos]

	fmt.Println(startingNode)

	//partOne(startingPos, grid, keys, doors)
}
