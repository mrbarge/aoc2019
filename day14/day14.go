package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Formula struct {
	produces string
	quantity int
	requires map[string]int
}

func parseEntry(s string) (amount int, name string) {
	elems := strings.Split(s, " ")
	if len(elems) < 2 {
		panic("Invalid entry " + s)
	}
	amount, _ = strconv.Atoi(elems[0])
	name = elems[1]
	return amount, name
}

func reachedGoal(need map[string]int, goal string) bool {
	_, ok := need[goal]
	// if the goal isn't in there, we aren't done
	if !ok {
		return false
	}
	// if there are any non-zero non-goals we aren't done
	for k, v := range need {
		if k != goal && v > 0 {
			return false
		}
	}
	return true
}

func partOne(f map[string]Formula, fuelGoal int) int {

	// keep track of quantities needed
	needs := make(map[string]int, 0)

	// part 1 wants 1 fuel
	needs["FUEL"] = fuelGoal

	// we want to find how much ore is needed
	goal := "ORE"

	// keep track of excess product
	have := make(map[string]int, 0)

	// cycle until we only have ore left in the needs
	for !reachedGoal(needs, goal) {
		// cycle through each needs
		for need, _ := range needs {
			// don't attempt to find constituents for our goal..
			if need == goal {
				continue
			}

			// can we use some excess before creating more?
			if _, ok := have[need]; ok {
				needs[need] -= have[need]
				delete(have, need)
			}

			// still got to make some?
			if needs[need] > 0 {

				// how much do we need to produce to get to our goal?
				product := 1
				if f[need].quantity < needs[need] {
					product = needs[need] / f[need].quantity
				}
				// get our list of things needed to satisfy this requirement
				constituents := f[need].requires

				for moreNeeds, moreNeedVal := range constituents {
					needs[moreNeeds] += (moreNeedVal * product)
				}

				needs[need] -= (f[need].quantity * product)

				// done with this need
				if needs[need] <= 0 {
					// record excess
					have[need] = 0 - needs[need]
					// remove need
					delete(needs, need)
				}
			}
		}
	}

	return needs["ORE"]
}

func partTwo(f map[string]Formula) {
	fmt.Println(sort.Search(1000000000000, func(n int) bool {
		return partOne(f, n) > 1000000000000
	}) - 1)
}

func main() {

	file, _ := os.Open("input.txt")
	s := bufio.NewScanner(file)

	r_log, _ := regexp.Compile(`^(.+) => (.+)$`)

	formulas := make(map[string]Formula, 0)

	for s.Scan() {
		line := s.Text()
		res_log := r_log.FindStringSubmatch(line)
		if res_log != nil {
			rFuel, rName := parseEntry(res_log[2])

			formulas[rName] = Formula{rName, rFuel, make(map[string]int, 0)}

			lElems := strings.Split(res_log[1], ", ")
			for _, lElem := range lElems {
				lFuel, lName := parseEntry(lElem)
				formulas[rName].requires[lName] = lFuel
			}
		}
	}

	fmt.Println(partOne(formulas, 1))
	partTwo(formulas)
}