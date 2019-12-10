package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Coord struct {
	x int
	y int
}

type CoordDist struct {
	c Coord
	distance int
}

func getTangentAngle(c1 Coord, c2 Coord) float64 {
	atan := math.Atan2(float64(c2.x-c1.x), float64(c1.y-c2.y))
	ab :=  atan * 180 / math.Pi
	if ab < 0 {
		return ab + 360
	} else {
		return ab
	}
}

func manhattan(c1 Coord, c2 Coord) int {
	return int(math.Abs(float64(c1.x)-float64(c2.x)) + math.Abs(float64(c1.y)-float64(c2.y)))
}

func main() {

	file, _ := os.Open("input.txt")
	s := bufio.NewScanner(file)

	asteroids := make(map[Coord]bool, 0)
	visible := make(map[Coord]int, 0)

	y := 0
	for s.Scan() {
		line := s.Text()

		for x, v := range line {
			if v == '#' {
				asteroids[Coord{x, y}] = true
			}
		}
		y += 1
	}

	for asteroid, _ := range asteroids {
		targetTangents := make(map[float64][]Coord)
		for cmp, _ := range asteroids {
			if asteroid == cmp {
				continue
			}
			tangle := getTangentAngle(asteroid, cmp)
			if _, ok := targetTangents[tangle]; !ok {
				// first time we're seeing an asteroid at this angle
				targetTangents[tangle] = []Coord{cmp}
				// flag that it is see-able
				visible[asteroid] += 1
			}
		}
	}

	max := 0
	var maxCoord Coord
	for k, v := range visible {
		if v > max {
			max = v
			maxCoord = k
		}
	}

	// part 1
	fmt.Println(max)
	fmt.Println(maxCoord)

	targetDistances := make(map[float64][]CoordDist)
	for asteroid, _ := range asteroids {
		if maxCoord == asteroid {
			continue
		}
		tangle := getTangentAngle(maxCoord, asteroid)
		manhattanDist := manhattan(maxCoord, asteroid)
		if _, ok := targetDistances[tangle]; !ok {
			targetDistances[tangle] = []CoordDist{CoordDist{asteroid, manhattanDist}}
		} else {
			targetDistances[tangle] = append(targetDistances[tangle], CoordDist{asteroid, manhattanDist})
		}
	}
	// sort distances based on closest to furthest
	for k, _ := range targetDistances {
		sort.Slice(targetDistances[k][:], func(i, j int) bool { return targetDistances[k][i].distance < targetDistances[k][j].distance })
	}

	// get a list of angles from closest to furthest
	tangentAngles := make([]float64, 0)
	for k, _ := range targetDistances {
		tangentAngles = append(tangentAngles, k)
	}
	sort.Float64s(tangentAngles)

	// keep on spinning til there's nothing left
	done := false
	killCount := 0
	for !done {
		for _, angle := range tangentAngles {
			// is there an asteroid left at this angle
			if len(targetDistances[angle]) > 0 {
				// vaporized!!
				killCount += 1

				if killCount == 200 {
					fmt.Println(targetDistances[angle][0])
					done = true
				}

				targetDistances[angle] = targetDistances[angle][1:]
			}
		}
	}
}