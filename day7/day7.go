package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Input struct {
	value int64
	next *Input
}

type ParameterMode int
const (
	POSITION = iota
	IMMEDIATE
)

type OpCode int
const (
	ADDITION = iota + 1
	MULTIPLY		// 2
	STORE			// 3
	OUTPUT			// 4
	JIT				// 5
	JIF				// 6
	LT				// 7
	EQ				// 8
	QUIT = 99
)

func parseOpCode(i int64) OpCode {
	s := strconv.Itoa(int(i))
	if len(s) > 1 {
		r, _ := strconv.Atoi(s[len(s)-2:])
		return OpCode(r)
	} else {
		r, _ := strconv.Atoi(s)
		return OpCode(r)
	}
}

// Taken from
// https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func permutations(arr []int)[][]int{
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int){
		if n == 1{
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++{
				helper(arr, n - 1)
				if n % 2 == 1{
					tmp := arr[i]
					arr[i] = arr[n - 1]
					arr[n - 1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n - 1]
					arr[n - 1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func parseModes(i int64) []ParameterMode {
	// fill modes with default position mode - generic to cater for numbers > 4 digits
	s := strconv.Itoa(int(i))
	rl := 3
	if len(s) - 3 > 3 {
		rl = len(s) - 3
	}
	r := make([]ParameterMode, rl)
	for i, _ := range r {
		r[i] = POSITION
	}

	ptr := 0
	for i := len(s)-3; i >= 0; i-- {
		if s[i] == '0' {
			r[ptr] = POSITION
			ptr += 1
		} else if s[i] == '1' {
			r[ptr] = IMMEDIATE
			ptr += 1
		} else {
			fmt.Printf("Error parsing mode: %d\n", i)
		}
	}
	return r
}

func getPositionOrImmediate(program []int64, mode ParameterMode, value int64) int64 {
	if mode == POSITION {
		return program[value]
	} else if mode == IMMEDIATE {
		return value
	} else {
		fmt.Printf("Error with mode argument: %d\n", mode)
		return -1
	}
}

func cycle(program []int64, pos int, input *Input) (p int, output int64, finished bool) {

	op := parseOpCode(program[pos])
	modes := parseModes(program[pos])

	switch op {

	case ADDITION:
		value1 := getPositionOrImmediate(program, modes[0], program[pos+1])
		value2 := getPositionOrImmediate(program, modes[1], program[pos+2])
		posDest := program[pos+3]
		program[posDest] = value1 + value2
		pos = pos + 4

	case MULTIPLY:
		value1 := getPositionOrImmediate(program, modes[0], program[pos+1])
		value2 := getPositionOrImmediate(program, modes[1], program[pos+2])
		posDest := program[pos+3]
		program[posDest] = value1 * value2
		pos = pos + 4

	case STORE:
		posDest := program[pos+1]
		program[posDest] = input.value
		if input.next != nil {
			input.next = input.next.next
		}
		pos = pos + 2

	case OUTPUT:
		posDest := getPositionOrImmediate(program, modes[0], program[pos+1])
		output = posDest
		pos = pos + 2

	case JIT:
		value1 := getPositionOrImmediate(program, modes[0], program[pos+1])
		value2 := getPositionOrImmediate(program, modes[1], program[pos+2])
		if value1 != 0 {
			pos = int(value2)
		} else {
			pos += 3
		}

	case JIF:
		value1 := getPositionOrImmediate(program, modes[0], program[pos+1])
		value2 := getPositionOrImmediate(program, modes[1], program[pos+2])
		if value1 == 0 {
			pos = int(value2)
		} else {
			pos += 3
		}

	case LT:
		value1 := getPositionOrImmediate(program, modes[0], program[pos+1])
		value2 := getPositionOrImmediate(program, modes[1], program[pos+2])
		posDest := program[pos+3]
		if value1 < value2 {
			program[posDest] = 1
		} else {
			program[posDest] = 0
		}
		pos += 4

	case EQ:
		value1 := getPositionOrImmediate(program, modes[0], program[pos+1])
		value2 := getPositionOrImmediate(program, modes[1], program[pos+2])
		posDest := program[pos+3]
		if value1 == value2 {
			program[posDest] = 1
		} else {
			program[posDest] = 0
		}
		pos += 4

	case QUIT:
		finished = true
	default:
		fmt.Printf("Unknown op code: %d\n", program[pos])
	}

	return pos, output, finished
}

func runUntilHalt(program []int64, input *Input) (output int64) {
	done := false
	pos := 0
	for !done {
		pos, output, done = cycle(program, pos, input)
	}
	return output
}

func thrusterSignal(program []int64, ampSequence []int) int64 {

	// first amplifier gets 0 input
	chainInput := int64(0)

	for _, amplifier := range ampSequence {
		// clone the program
		candidateProg := make([]int64, len(program))
		copy(candidateProg, program)

		// craft our input chain
		inputs := Input{int64(amplifier), &Input{chainInput, nil}}

		fmt.Printf("Running amplifier %d with input %d\n", amplifier, chainInput)
		chainInput = runUntilHalt(candidateProg, &inputs)
	}

	return chainInput
}

func main() {

	bd, err := ioutil.ReadFile("test.txt")
	//bd, err := ioutil.ReadFile("input.txt")
	if err != nil {
		os.Exit(1)
	}
	sProg := strings.Split(string(bd), ",")
	originalProg := make([]int64, len(sProg))
	for i, v := range sProg {
		originalProg[i], _ = strconv.ParseInt(v, 10, 64)
	}

	inputPerms := permutations([]int{0,1,2,3,4})

	thrusterSignals := make([]int64, 0)
	for _, inputPerm := range inputPerms {
		thrust := thrusterSignal(originalProg, inputPerm)
		thrusterSignals = append(thrusterSignals, thrust)
	}

	sort.Slice(thrusterSignals, func(i, j int) bool { return thrusterSignals[i] < thrusterSignals[j] })
	fmt.Println(thrusterSignals)
}

