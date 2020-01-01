package main

import (
	"bufio"
	"fmt"
	"github.com/albertorestifo/dijkstra"
	"gonum.org/v1/gonum/graph/simple"
	"os"
	"sort"
	"unicode"
)

type ExploreState struct {
	steps int
	state SeenState
}

type SeenState struct {
	at Coord
	keys []rune
}

func (s SeenState) hasKey(k rune) bool {
	for _, v := range s.keys {
		if v == unicode.ToLower(k) {
			return true
		}
	}
	return false
}

func (s SeenState) hash() string {
	sort.Slice(s.keys, func(i, j int) bool {
		return s.keys[i] < s.keys[j]
	})
	return fmt.Sprintf("%d:%d:%s", s.at.x, s.at.y, string(s.keys))
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

func getNeighbours(pos Coord, grid Grid) (neighbours []Coord) {
	if grid[pos.y][pos.x-1] != '#'  {
		neighbours = append(neighbours, Coord{pos.x-1,pos.y})
	}
	if grid[pos.y][pos.x+1] != '#'  {
		neighbours = append(neighbours, Coord{pos.x+1,pos.y})
	}
	if grid[pos.y-1][pos.x] != '#' {
		neighbours = append(neighbours, Coord{pos.x,pos.y-1})
	}
	if grid[pos.y+1][pos.x] != '#' {
		neighbours = append(neighbours, Coord{pos.x,pos.y+1})
	}
	return neighbours
}

func makeGraph(grid Grid) (*simple.UndirectedGraph, map[Coord]int64, map[int64]Coord) {

	g := simple.NewUndirectedGraph()

	nodeCoordToId := make(map[Coord]int64)
	nodeIdToCoord := make(map[int64]Coord)

	nodeCount := int64(0)
	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[y])-1; x++ {
			tmpCoord := Coord{x, y}
			g.AddNode(simple.Node(nodeCount))
			if _, ok := nodeCoordToId[tmpCoord]; !ok {
				nodeCoordToId[tmpCoord] = nodeCount
				nodeIdToCoord[nodeCount] = tmpCoord
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

	return g, nodeCoordToId, nodeIdToCoord
}

func makeGraph2(grid Grid) (dijkstra.Graph, map[Coord]string, map[string]Coord) {

	g := dijkstra.Graph{}

	nodeCoordToId := make(map[Coord]string)
	nodeIdToCoord := make(map[string]Coord)

	nodeCount := int64(0)
	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[y])-1; x++ {
			tmpCoord := Coord{x, y}
			nodeId := fmt.Sprintf("%d-%d",x,y)
			if _, ok := nodeCoordToId[tmpCoord]; !ok {
				g[nodeId] = make(map[string]int, 0)
				nodeCoordToId[tmpCoord] = nodeId
				nodeIdToCoord[nodeId] = tmpCoord
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
				g[nodeId][neighNodeId] = 1
				g[neighNodeId][nodeId] = 1
			}
		}
	}

	return g, nodeCoordToId, nodeIdToCoord
}

func (g Grid) isDoor(c Coord) bool {
	return unicode.IsUpper(g[c.y][c.x])
}

func (g Grid) isKey(c Coord) bool {
	return unicode.IsLower(g[c.y][c.x])
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
				doors[v] = Coord{xPos, yPos }
			} else if unicode.IsLower(v) {
				// found a key
				keys[v] = Coord{xPos, yPos }
			} else if v == '@' {
				startingPos = Coord{xPos, yPos }
			}
		}
		yPos++
	}

	//g, nodeCoordToId, nodeIdToCoord := makeGraph2(grid)
	//startingNode := nodeCoordToId[startingPos]
	//
	partOne(grid, startingPos, doors, keys)

	//partOne(startingPos, grid, keys, doors)

	//
	//b := traverse.BreadthFirst {
	//	Traverse: func(e graph.Edge) bool {
	//		// is the 'to' edge a door?
	//		nodeTo := e.To()
	//		toCoord := nodeIdToCoord[nodeTo.ID()]
	//
	//		if grid.isDoor(toCoord) {
	//			// do we have a key
	//			keyName := unicode.ToLower(grid[toCoord.y][toCoord.x])
	//			keyCoord := keys[keyName]
	//			keyNodeId := nodeCoordToId[keyCoord]
	//			fmt.Println(keyNodeId)
	//		}
	//		return true
	//	},
	//	Visit: func(n graph.Node) {
	//	},
	//}
}

func partOne(g Grid, start Coord, doors Items, keys Items) {

	seen := make(map[string]bool, 0)
	queue := make([]ExploreState, 0)
	queueRec := make(map[string]bool, 0)

	// add starting pos
	queue = append(queue, ExploreState{0, SeenState{start, []rune{}}})

	done := false
	for !done {
		// pop the queue

		if len(queue) == 0 {
			fmt.Println("Bad")
			break
		}

		current := queue[0]
		queue = queue[1:]
		neighs := getNeighbours(current.state.at, g)

		// have we gotten all keys yet?
		if len(current.state.keys) == len(keys) {
			fmt.Println(current.steps)
			break
		}

		//fmt.Println("Adding to hash: ", current.state.hash())
		seen[current.state.hash()] = true

		for _, neigh := range neighs {
			newkeys := make([]rune,len(current.state.keys))
			copy(newkeys, current.state.keys)
			seenIt := ExploreState{current.steps+1, SeenState{neigh, newkeys}}
			seenItHash := seenIt.state.hash()

			//fmt.Println("Seen it hash: ", seenIt.state.hash())
			// have we already been at this node with the amount of keys we have?
			if _, ok := seen[seenItHash]; ok {
				// have already been in this state, ignore it
				//fmt.Println("Ignoring pos because weve been here\n")
				continue
			}
			if _, ok := queueRec[seenItHash]; ok {
				continue
			}

			tile := g[neigh.y][neigh.x]

			if _, ok := doors[tile]; ok {
				// we are at a door
				if !seenIt.state.hasKey(tile) {
					// we don't have a key, so ignore this path
					//fmt.Printf("Canot move beyond door %s\n",string(tile))
					continue
				}
			}

			if _, ok := keys[tile]; ok {
				// we are at a key
				//fmt.Printf("Before key %s (%s)\n",string(tile),seenIt.state.hash())
				if ! seenIt.state.hasKey(tile) {
					seenIt.state.keys = append(seenIt.state.keys, tile)
				}
				//fmt.Printf("Found key %s (%s %d)\n",string(tile),seenItHash,len(seenIt.state.keys))
			}

			// add this to the queue of explored spaces
			if _, ok := queueRec[seenItHash]; !ok {
				queue = append(queue, seenIt)
				queueRec[seenItHash] = true
			}
		}
	}

}

func isInQueue(s string, queue []ExploreState) bool {
	for _, v := range queue {
		if v.state.hash() == s {
			return true
		}
	}
	return false
}
func printQueue(q map[string]bool) {
	for k, _ := range q {
		fmt.Print(k)
		fmt.Print(" ")
	}
	fmt.Println()
}
//func partOne(g dijkstra.Graph, beginNodeId string, doors Items, keys Items, nodeCoordToId map[Coord]string, nodeIdToCoord map[string]Coord) {
//
//}