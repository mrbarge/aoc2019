package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func fuel(mass int, fuelNeedsFuel bool) int {
	fuelRequired := int(math.Floor(float64(mass) / 3) - 2)
	if fuelRequired < 0 {
		return 0
	} else if fuelNeedsFuel {
		return fuelRequired + fuel(fuelRequired, fuelNeedsFuel)
	} else {
		return fuelRequired
	}
}

func main() {

	file, _ := os.Open("input.txt")
	s := bufio.NewScanner(file)

	partOneSum := 0
	partTwoSum := 0
	for s.Scan() {
		line := s.Text()
		mass, err := strconv.Atoi(line)
		if err == nil {
			partOneSum += fuel(mass, false)
			partTwoSum += fuel(mass, true)
		}
	}
	fmt.Printf("%d\n", partOneSum)
	fmt.Printf("%d\n", partTwoSum)

}
