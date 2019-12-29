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

func (c Coord) getNeighboursP2() []Coord {
	r := make([]Coord, 0)
	r = append(r, Coord{c.x-1, c.y})
	r = append(r, Coord{c.x+1, c.y})
	r = append(r, Coord{c.x, c.y+1})
	r = append(r, Coord{c.x, c.y-1})
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

func countNeighbourBugsP2(grids map[int]Grid, depth int, c Coord) int {
	bugs := 0
	neighs := c.getNeighboursP2()

	for _, neigh := range neighs {
		if neigh.x == 2 && neigh.y == 2 {
			// middle square - look inner
			if c.x == 1 {
				// test all ys on inner right
				for ny := 0; ny < 5; ny++ {
					if grids[depth+1][ny][4] {
						bugs++
					}
				}
			} else if c.x == 3 {
				// test all ys on inner left
				for ny := 0; ny < 5; ny++ {
					if grids[depth+1][ny][0] {
						bugs++
					}
				}
			} else if c.y == 1 {
				// test all xs on inner bottom
				for nx := 0; nx < 5; nx++ {
					if grids[depth+1][0][nx] {
						bugs++
					}
				}
			} else if c.y == 3 {
				// test all xs on inner top
				for nx := 0; nx < 5; nx++ {
					if grids[depth+1][4][nx] {
						bugs++
					}
				}
			}
			// end for middle square
		} else if neigh.x == -1 || neigh.x == 5 || neigh.y == -1 || neigh.y == 5 {
			// outer edges.. look outer mid-left
			if neigh.x == -1 && grids[depth-1][2][1] {
				bugs++
			}
			if neigh.x == 5 && grids[depth-1][2][3] {
				bugs++
			}
			if neigh.y == -1 && grids[depth-1][1][2] {
				bugs++
			}
			if neigh.y == 5 && grids[depth-1][3][2] {
				bugs++
			}
			// end for outer edges
		} else {
			if grids[depth][neigh.y][neigh.x] {
				bugs++
			}
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

func partOne(grid Grid) {
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

func emptyGrid() Grid {
	m := make(Grid, 0)
	for y := 0; y < 5; y++ {
		m = append(m, make([]bool, 5))
	}
	return m
}

func countBugs(grid Grid) int {
	nb := 0
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] {
				nb++
			}
		}
	}
	return nb
}

func partTwo(grid Grid) {

	depthGrids := make(map[int]Grid, 0)
	for i := 0; i < 400; i++ {
		depthGrids[i] = emptyGrid()
	}
	depthGrids[200] = grid

	minDepth, maxDepth := 199, 201

	for min := 0; min < 200; min++ {
		tmpDepthGrids := make(map[int]Grid, 0)
		for i := 0; i < 400; i++ {
			tmpDepthGrids[i] = emptyGrid()
		}

		for depth := minDepth; depth <= maxDepth; depth++ {

			tmpGrid := emptyGrid()

			for y := 0; y < 5; y++ {
				for x := 0; x < 5; x++ {
					if x == 2 && y == 2 {
						continue
					}
					c := Coord{x,y }

					numBugs := countNeighbourBugsP2(depthGrids, depth, c)
					tmpGrid[y][x] = depthGrids[depth][y][x]
					if depthGrids[depth][y][x] && numBugs != 1 {
						tmpGrid[y][x] = false
					} else if !depthGrids[depth][y][x] && numBugs > 0 && numBugs < 3 {
						tmpGrid[y][x] = true
					}

					tmpDepthGrids[depth] = tmpGrid
				}
			}
		}
		depthGrids = tmpDepthGrids

		for i := minDepth-2; i < maxDepth+2; i++ {
			if countBugs(depthGrids[i]) > 0 {
				if i < minDepth {
					minDepth = i
				} else if i > maxDepth {
					maxDepth = i
				}
			}
		}
	}

	total := 0
	for _, v := range depthGrids {
		total += countBugs(v)
	}
	fmt.Println(total)
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

	p1Grid := make(Grid, len(grid))
	p2Grid := make(Grid, len(grid))
	for i := range grid {
		p1Grid[i] = make([]bool, len(grid[i]))
		copy(p1Grid[i], grid[i])
		p2Grid[i] = make([]bool, len(grid[i]))
		copy(p2Grid[i], grid[i])
	}

	//partOne(p1Grid)
	partTwo(p2Grid)
}