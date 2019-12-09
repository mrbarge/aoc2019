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
	RELATIVE
)

type Computer struct {
	program []int64
	input int64
	relativeBase int64
	pos int64
	finished bool
}

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
	RBO				// 9
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
		} else if s[i] == '2' {
			r[ptr] = RELATIVE
			ptr += 1
		} else {
			fmt.Printf("Error parsing mode: %d\n", i)
		}
	}
	return r
}

func getPositionOrImmediate(c *Computer, mode ParameterMode, value int64, read bool) int64 {
	if mode == POSITION {
		if read {
			return c.program[value]
		} else {
			return value
		}
	} else if mode == IMMEDIATE {
		return value
	} else if mode == RELATIVE {
		if read {
			return c.program[c.relativeBase+value]
		} else {
			return c.relativeBase+value
		}
	} else {
		fmt.Printf("Error with mode argument: %d\n", mode)
		return -1
	}
}

func cycle(c *Computer) {

	op := parseOpCode(c.program[c.pos])
	modes := parseModes(c.program[c.pos])

	switch op {

	case ADDITION:
		value1 := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		value2 := getPositionOrImmediate(c, modes[1], c.program[c.pos+2], true)
		posDest := getPositionOrImmediate(c, modes[2], c.program[c.pos+3], false)
		c.program[posDest] = value1 + value2
		c.pos = c.pos + 4

	case MULTIPLY:
		value1 := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		value2 := getPositionOrImmediate(c, modes[1], c.program[c.pos+2], true)
		posDest := getPositionOrImmediate(c, modes[2], c.program[c.pos+3], false)
		c.program[posDest] = value1 * value2
		c.pos = c.pos + 4

	case STORE:
		posDest := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], false)
		c.program[posDest] = c.input
		c.pos = c.pos + 2

	case OUTPUT:
		posDest := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		fmt.Printf("Output: %d\n",posDest)
		c.pos = c.pos + 2

	case JIT:
		value1 := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		value2 := getPositionOrImmediate(c, modes[1], c.program[c.pos+2], true)
		if value1 != 0 {
			c.pos = int64(value2)
		} else {
			c.pos += 3
		}

	case JIF:
		value1 := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		value2 := getPositionOrImmediate(c, modes[1], c.program[c.pos+2], true)
		if value1 == 0 {
			c.pos = int64(value2)
		} else {
			c.pos += 3
		}

	case LT:
		value1 := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		value2 := getPositionOrImmediate(c, modes[1], c.program[c.pos+2], true)
		posDest := getPositionOrImmediate(c, modes[2], c.program[c.pos+3], false)
		if value1 < value2 {
			c.program[posDest] = 1
		} else {
			c.program[posDest] = 0
		}
		c.pos += 4

	case EQ:
		value1 := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		value2 := getPositionOrImmediate(c, modes[1], c.program[c.pos+2], true)
		posDest := getPositionOrImmediate(c, modes[2], c.program[c.pos+3], false)
		if value1 == value2 {
			c.program[posDest] = 1
		} else {
			c.program[posDest] = 0
		}
		c.pos += 4

	case RBO:
		adjustment := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		c.relativeBase += adjustment
		c.pos += 2

	case QUIT:
		c.finished = true

	default:
		fmt.Printf("Unknown op code: %d\n", c.program[c.pos])
	}
}

func runUntilHalt(program []int64, input int64) int64 {
	c := Computer{program, input,0, 0, false }
	for !c.finished {
		cycle(&c)
	}
	return c.program[0]
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
	bufferSpace := make([]int64, 10000)
	candidateProg = append(candidateProg, bufferSpace...)
	runUntilHalt(candidateProg, 1)

	// part 2
	candidateProg = make([]int64, len(originalProg))
	copy(candidateProg, originalProg)
	candidateProg = append(candidateProg, bufferSpace...)
	runUntilHalt(candidateProg, 2)

}

