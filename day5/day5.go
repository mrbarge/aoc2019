package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

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

func cycle(program []int64, pos int, input int64) (p int, finished bool) {

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
		program[posDest] = input
		pos = pos + 2

	case OUTPUT:
		posDest := getPositionOrImmediate(program, modes[0], program[pos+1])
		fmt.Printf("Output: %d\n",posDest)
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

	return pos, finished
}

func runUntilHalt(program []int64, input int64) int64 {
	done := false
	pos := 0
	for !done {
		pos, done = cycle(program, pos, input)
	}
	return program[0]
}

func main() {

	//bd, err := ioutil.ReadFile("test.txt")
	bd, err := ioutil.ReadFile("input.txt")
	if err != nil {
		os.Exit(1)
	}
	sProg := strings.Split(string(bd), ",")
	originalProg := make([]int64, len(sProg))
	for i, v := range sProg {
		originalProg[i], _ = strconv.ParseInt(v, 10, 64)
	}

	candidateProg := make([]int64, len(originalProg))
	copy(candidateProg, originalProg)

	// part 1
	runUntilHalt(candidateProg, 1)
	// part 2
	copy(candidateProg, originalProg)
	runUntilHalt(candidateProg, 5)
}

