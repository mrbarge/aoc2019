package main

import (
	"fmt"
	"strconv"
)

func tests(n int, partOne bool) bool {

	adjacentFound := false
	sNum := strconv.Itoa(n)
	countAdjacent := 0
	lastNum := -1

	for i := 0; i < int(len(sNum)); i++ {
		// test for increasing digits
		if i < len(sNum)-1 && sNum[i+1] < sNum[i] {
			return false
		}

		if lastNum != int(sNum[i]) {
			if !partOne && countAdjacent == 2 {
				adjacentFound = true
			} else if partOne && countAdjacent >= 2 {
				adjacentFound = true
			}
			countAdjacent = 1
			lastNum = int(sNum[i])
		} else {
			countAdjacent += 1
		}
	}

	// for the end of the range
	if (partOne && countAdjacent >= 2) || (!partOne && countAdjacent == 2) {
		adjacentFound = true
	}

	return adjacentFound
}

func run(lower int, upper int, partOne bool) (count int) {

	for i := lower; i <= upper; i++ {
		if tests(i, partOne) {
			count += 1
		}
	}
	return count
}

func main() {
	fmt.Println(run(234208, 765869, true))
	fmt.Println(run(234208, 765869, false))
}
