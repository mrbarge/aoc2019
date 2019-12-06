package main

import (
	"bufio"
	"fmt"
	"os"
)

type Planet struct {
	name string
	orbits *Planet
}

func countOrbits(p *Planet) int {

	if p.orbits == nil {
		return 0
	} else {
		return 1 + countOrbits(p.orbits)
	}

}

func getRouteBack(p *Planet) []string {
	if p.orbits == nil {
		return []string{}
	} else {
		return append([]string{p.orbits.name}, getRouteBack(p.orbits)...)
	}
}

func main() {

	file, _ := os.Open("input.txt")
	s := bufio.NewScanner(file)

	planets := make(map[string]*Planet)

	for s.Scan() {
		line := s.Text()

		p1 := line[0:3]
		p2 := line[4:]

		if _, ok := planets[p1]; !ok {
			planet1 := Planet{p1, nil}
			planets[p1] = &planet1
		}

		if _, ok := planets[p2]; !ok {
			planet2 := Planet{p2, nil}
			planets[p2] = &planet2
		}

		planets[p2].orbits = planets[p1]

	}

	total := 0
	for _, v := range planets {
		total += countOrbits(v)
	}
	fmt.Printf("%d\n", total)

	youRoutes := getRouteBack(planets["YOU"])
	sanRoutes := getRouteBack(planets["SAN"])

	matchFound := false
	for i, x := range youRoutes {
		// does san have this as a parent planet?
		for j, y := range sanRoutes {
			if x == y {
				fmt.Printf("Match found: %d", i+j)
				matchFound = true
				break
			}
		}
		if matchFound {
			break
		}
	}
}
