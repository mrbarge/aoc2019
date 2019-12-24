package main

import (
	"bufio"
	"fmt"
	"os"
)

type Grid [][]bool

type Coord struct {
	x int
	y int
}

func (c Coord) getNeighbours(maxX int, maxY int) []Coord {
	r := make([]Coord, 0)
	if c.x > 0 {
		r = append(r, Coord{c.x-1, c.y})
	}
	if c.x < maxX {
		r = append(r, Coord{c.x+1, c.y})
	}
	if c.y < maxY {
		r = append(r, Coord{c.x, c.y+1})
	}
	if c.y > 0 {
		r = append(r, Coord{c.x, c.y-1})
	}
	return r
}

func cycle(g Grid) Grid {
	ret := make(Grid, len(g))
	for i := range g {
		ret[i] = make([]bool, len(g[i]))
	}

	for y := range g {
		for x := range g[y] {
			c := Coord{x,y}
			bugs := g.countNeighbourBugs(c)
			ret[y][x] = g[y][x]
			if g[y][x] && bugs != 1 {
				ret[y][x] = false
			} else if !g[y][x] && bugs > 0 && bugs < 3 {
				ret[y][x] = true
			}
		}
	}
	return ret
}

func (g Grid) countNeighbourBugs(c Coord) int {
	bugs := 0
	neighs := c.getNeighbours(len(g[0])-1, len(g)-1)
	for _, n := range neighs {
		if g[n.y][n.x] {
			bugs++
		}
	}
	return bugs
}

func (g Grid) printGrid() {
	for y := range g {
		for x := range g[y] {
			if g[y][x] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (g Grid) biodiversity() int {
	r := 0
	rating := [][]int{{1,2,4,8,16},{32,64,128,256,512},{1024,2048,4096,8192,16384},{32768,65536,131072,262144,524288},{1048576,2097152,4194304,8388608,16777216}}
	for y := range g {
		for x := range g[y] {
			if g[y][x] {
				r += rating[y][x]
			}
		}
	}
	return r
}

func main() {

	file, _ := os.Open("input.txt")
	s := bufio.NewScanner(file)
	grid := make(Grid, 0)

	for s.Scan() {
		line := s.Text()
		row := make([]bool, len(line))
		for i, c := range line {
			if c == '#' {
				row[i] = true
			} else {
				row[i] = false
			}
		}
		grid = append(grid, row)
	}

	done := false
	seen := make(map[int]bool, 0)
	for !done {
		bd := grid.biodiversity()
		if _, ok := seen[bd]; ok {
			done = true
		} else {
			seen[bd] = true
			grid = cycle(grid)
		}
	}
	fmt.Println(grid.biodiversity())
}