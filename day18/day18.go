package main

import (
	"bufio"
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
	"os"
	"unicode"
)

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

	g, nodeCoordToId, nodeIdToCoord := makeGraph(grid)
	startingNode := nodeCoordToId[startingPos]

	b := traverse.BreadthFirst {
		Traverse: func(e graph.Edge) bool {
			// is the 'to' edge a door?
			nodeTo := e.To()
			toCoord := nodeIdToCoord[nodeTo.ID()]

			if grid.isDoor(toCoord) {
				// do we have a key
				keyName := unicode.ToLower(grid[toCoord.y][toCoord.x])
				keyCoord := keys[keyName]
				keyNodeId := nodeCoordToId[keyCoord]
				fmt.Println(keyNodeId)
			}
			return true
		},
		Visit: func(n graph.Node) {
		},
	}

	fmt.Println(startingNode)

	//partOne(startingPos, grid, keys, doors)
}
