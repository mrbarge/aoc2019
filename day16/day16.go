package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

var basePattern = []int {0, 1, 0, -1}

func makePattern(base []int, repeat int, patternLength int) []int {
	first := true
	lengthReached := false
	r := make([]int, 0)

	for !lengthReached {
		for i := 0; i < len(base) && !lengthReached; i++ {
			for j := 0; j < repeat && !lengthReached; j++ {
				if first {
					// skip the first
					first = false
					continue
				}

				r = append(r, base[i])
				if len(r) == patternLength {
					lengthReached = true
					break
				}
			}
		}
	}
	return r
}

func firstDigit(i int) int {
	tmp := int(math.Abs(float64(i)))
	for tmp > 9 {
		tmp %= 10
	}
	return tmp
}

func applySum(input []int, pattern []int) int {
	result := 0
	for i := 0; i < len(input); i++ {
		result += input[i] * pattern[i]
	}
	return firstDigit(result)
}

func cycle(input []int) []int {
	r := make([]int, 0)
	inputLen := len(input)

	for i := 0; i < len(input); i++ {
		pattern := makePattern(basePattern, i+1, inputLen)
		//fmt.Printf("Input: %+v Pattern: %+v\n", input, pattern)
		r = append(r, applySum(input, pattern))
	}
	//fmt.Printf("Returning: %+v\n", r)
	return r
}

func cyclep2(input []int, offset int) []int {
	for i := len(input)-2; i > offset-7; i-- {
		n := input[i+1] + input[i]
		input[i] = firstDigit(n)
	}

	return input
}

func main() {
	bd, err := ioutil.ReadFile("input.txt")
	//bd, err := ioutil.ReadFile("test.txt")
	if err != nil {
		os.Exit(1)
	}
	inputLine := string(bd)
	input := make([]int, 0)
	for _, v := range inputLine {
		num, _ := strconv.Atoi(string(v))
		input = append(input, num)
	}

	fmt.Println(len(input))

	p1Input := make([]int, len(input))
	copy(p1Input, input)
	phases := 0
	for phases < 100 {
		p1Input = cycle(p1Input)
		phases += 1
	}

	// part 2
	codeLoc := 5977341

	p2Input := make([]int, len(input))
	for i := 0; i < 9999; i++ {
		p2Input = append(p2Input, input...)
	}
	phases = 0
	for phases < 100 {
		p2Input = cyclep2(p2Input, codeLoc)
		phases += 1
	}
	for i := codeLoc; i < codeLoc+8; i++ {
		fmt.Print(p2Input[i])
	}

}
