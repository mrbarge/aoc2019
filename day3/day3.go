package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Direction rune
const (
	UP = 'U'
	DOWN = 'D'
	LEFT = 'L'
	RIGHT = 'R'
)

type Path struct {
	dir Direction
	distance int
}

type Coord struct {
	x int
	y int
}

func parsePath(p string) Path {
	dist, err := strconv.Atoi(p[1:])
	if err != nil {
		fmt.Printf("Invalid distance: %s\n", p)
	}
	return Path{dir: Direction(p[0]), distance: dist}
}

func followPaths(paths []Path, seen map[Coord]bool) []Coord {

	intersections := make([]Coord, 0)
	ownCoords := make(map[Coord]bool, 0)	// to detect self-crosses

	x, y := 0, 0
	for _, p := range paths {

		for i := 0; i < p.distance; i++ {
			c := Coord{x, y}
			_, alreadyVisited := ownCoords[c]
			if !alreadyVisited && seen[c] && (x != 0 && y != 0) {
				intersections = append(intersections, c)
			}
			seen[c] = true
			ownCoords[c] = true

			switch p.dir {
			case UP:
				y -= 1
			case DOWN:
				y += 1
			case LEFT:
				x -= 1
			case RIGHT:
				x += 1
			}
		}
	}
	return intersections
}

func followPathsSteps(paths []Path, seen map[Coord]int) map[Coord]int {

	intersections := make(map[Coord]int, 0)
	ownCoords := make(map[Coord]bool, 0)	// to detect self-crosses

	steps, x, y := 0, 0, 0
	for _, p := range paths {

		for i := 0; i < p.distance; i++ {
			c := Coord{x, y}

			_, alreadyVisited := ownCoords[c]

			if !alreadyVisited && seen[c] > 0 && (x != 0 && y != 0) {
				intersections[c] = seen[c] + steps
			}
			if seen[c] == 0 {
				seen[c] = steps
			}
			ownCoords[c] = true

			switch p.dir {
			case UP:
				y -= 1
			case DOWN:
				y += 1
			case LEFT:
				x -= 1
			case RIGHT:
				x += 1
			}
			steps += 1
		}
	}
	return intersections
}


func manhattan(c Coord) int {
	return int(math.Abs(float64(c.x)) + math.Abs(float64(c.y)))
}

func partOne(wire1Paths []Path, wire2Paths []Path) {
	// find the intersections
	seenCoordinates := make(map[Coord]bool, 0)
	followPaths(wire1Paths, seenCoordinates)

	intersections := followPaths(wire2Paths, seenCoordinates)

	// find the manhattan distance of the closest
	smallest := math.MaxInt64
	for _, c := range intersections {
		md := manhattan(c)
		if md < smallest {
			smallest = md
		}
	}

	// answer
	fmt.Println(smallest)
}

func partTwo(wire1Paths []Path, wire2Paths []Path) {

	// find the intersections
	seenCoordinates := make(map[Coord]int, 0)
	followPathsSteps(wire1Paths, seenCoordinates)

	intersections := followPathsSteps(wire2Paths, seenCoordinates)

	smallest := math.MaxInt64
	for _, v := range intersections {
		if v < smallest {
			smallest = v
		}
	}

	// answer
	fmt.Println(smallest)
}

func main() {

	// input data handling

	file, _ := os.Open("input.txt")
	s := bufio.NewScanner(file)
	s.Scan()
	wire1 := s.Text()
	s.Scan()
	wire2 := s.Text()

	wire1Paths := make([]Path, 0)
	wire2Paths := make([]Path, 0)
	for _, v := range strings.Split(wire1, ",") {
		p := parsePath(v)
		wire1Paths = append(wire1Paths, p)
	}
	for _, v := range strings.Split(wire2, ",") {
		p := parsePath(v)
		wire2Paths = append(wire2Paths, p)
	}

	partOne(wire1Paths, wire2Paths)
	partTwo(wire1Paths, wire2Paths)
}
