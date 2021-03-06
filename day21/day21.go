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

var WALK = []rune{'W','A','L','K','\n'}
var RUN = []rune{'R','U','N','\n'}

type Computer struct {
	program []int64
	inputs []int64
	relativeBase int64
	pos int64
	finished bool
	outputs []int64
	inputBlocked bool
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
		if len(c.inputs) == 0 {
			c.inputBlocked = true
		} else {
			c.inputBlocked = false
			posDest := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], false)
			c.program[posDest] = c.inputs[0]
			c.inputs = c.inputs[1:]
			c.pos = c.pos + 2
		}

	case OUTPUT:
		posDest := getPositionOrImmediate(c, modes[0], c.program[c.pos+1], true)
		c.outputs = append(c.outputs, posDest)
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

func addInput(c *Computer, input int64) {
	c.inputs = append(c.inputs, input)
}

func popOutput(c *Computer) int64 {
	if len(c.outputs) > 0 {
		r := c.outputs[0]
		c.outputs = c.outputs[1:]
		return r
	} else {
		panic("Computer has no outputs to provide")
	}
}

func toInputArr(r []rune) []int64 {
	m := make([]int64,len(r))
	for i := range r {
		m[i] = int64(r[i])
	}
	return m
}

func notToInput(arg1 rune, arg2 rune) []int64 {
	return toInputArr([]rune{'N','O','T',' ',arg1,' ',arg2,'\n'})
}

func orToInput(arg1 rune, arg2 rune) []int64 {
	return toInputArr([]rune{'O','R',' ',arg1,' ',arg2,'\n'})
}
func andToInput(arg1 rune, arg2 rune) []int64 {
	return toInputArr([]rune{'A','N','D',' ',arg1,' ',arg2,'\n'})
}


func partOne(program []int64)  {
	// state of the world

	inputs := make([]int64, 0)

	// p1
	//inputs = append(inputs, orToInput('A','J')...)
	//inputs = append(inputs, andToInput('B','J')...)
	//inputs = append(inputs, andToInput('B','J')...)
	//inputs = append(inputs, andToInput('C','J')...)
	//inputs = append(inputs, notToInput('J','J')...)
	//inputs = append(inputs, andToInput('D','J')...)

	inputs = append(inputs, orToInput('A','J')...)
	inputs = append(inputs, andToInput('B','J')...)
	inputs = append(inputs, andToInput('B','J')...)
	inputs = append(inputs, andToInput('C','J')...)
	inputs = append(inputs, notToInput('J','J')...)
	inputs = append(inputs, andToInput('D','J')...)
	inputs = append(inputs, orToInput('E','T')...)
	inputs = append(inputs, orToInput('H','T')...)
	inputs = append(inputs, andToInput('T','J')...)

	inputs = append(inputs, toInputArr(RUN)...)

	c := Computer{program, inputs,0, 0, false, []int64{}, false}

	for !c.finished {

		// run program, get colour and direction
		for !c.finished && !c.inputBlocked && len(c.outputs) == 0 {
			cycle(&c)
		}

		if len(c.outputs) > 0 && !c.finished {
			fmt.Println(popOutput(&c))
		}
	}
}

func main() {

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

	bufferSpace := make([]int64, 10000)
	candidateProg = append(candidateProg, bufferSpace...)

	partOne(candidateProg)
}

