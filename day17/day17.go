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

type Tile rune
const (
	WALL = '#'
	EMPTY =	'.'
	ROBOT_NORTH = '^'
	ROBOT_SOUTH = 'v'
	ROBOT_WEST = '<'
	ROBOT_EAST = '>'
	NEWLINE = '\n'
)

type Computer struct {
	program []int64
	inputs []int64
	relativeBase int64
	pos int64
	finished bool
	outputs []int64
	inputBlocked bool
}

type Coord struct {
	x int
	y int
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

func partOne(program []int64) int {
	c := Computer{program, []int64{},0, 0, false, []int64{}, false}

	// state of the world
	grid := make([][]Tile, 0)
	finished := false

	//xPos := 0
	yPos := 0

	grid = append(grid, make([]Tile, 0))
	for !finished {
		// run program, get colour and direction
		for !c.finished && !c.inputBlocked && len(c.outputs) == 0 {
			cycle(&c)
		}

		if c.finished {
			break
		}

		mapResult := popOutput(&c)
		switch mapResult {
		case NEWLINE:
			if len(grid[yPos]) > 0 {
				yPos += 1
				grid = append(grid, make([]Tile, 0))
			}
		default:
			grid[yPos] = append(grid[yPos], Tile(mapResult))
		}
	}

	sum := findIntersectionSum(grid)

	printMap(grid)
	return sum
}

func partTwo(program []int64) int64 {
	// part 2 i did the route by hand
	// L4 L6 L8 L12 L8 R12 L12 L8 R12 L12 L4 L6 L8 L12 L8 R12 L12 R12 L6 L6 L8 L4 L6 L8 L12 R12 L6 L6 L8 L8 R12 L12 R12 L6 L6 L8
	// which simplified down to
	// L4 L6 L8 L12 A
	// L8 R12 L12 B
	// R12 L6 L6 L8 C
	// A B B A B C A C B C
	// which turns into input
	// 65,44,66,44,66,44,65,44,66,44,67,44,65,44,67,44,66,44,67,10  (program)
	// 76,44,52,44,76,44,54,44,76,44,56,44,76,44,49,50,10   (A)
	// 76,44,56,44,82,44,49,50,44,76,44,49,50,10 (B)
	// 82,44,49,50,44,76,44,54,44,76,44,54,44,76,44,56,10  (C)

	c := Computer{program, []int64{},0, 0, false, []int64{}, false}

	// state of the world
	finished := false

	// wake up robot
	c.program[0] = 2

	// feed inputs
	c.inputs = []int64{65,44,66,44,66,44,65,44,66,44,67,44,65,44,67,44,66,44,67,10,
		76,44,52,44,76,44,54,44,76,44,56,44,76,44,49,50,10,
		76,44,56,44,82,44,49,50,44,76,44,49,50,10,
		82,44,49,50,44,76,44,54,44,76,44,54,44,76,44,56,10,
		110,10,
	}

	dust := int64(0)

	for !finished {
		// run program, get colour and direction
		for !c.finished && !c.inputBlocked && len(c.outputs) == 0 {
			cycle(&c)
		}

		if c.finished {
			break
		}

		if len(c.outputs) > 0 {
			fmt.Println(len(c.inputs))
			dust = popOutput(&c)
		}
	}

	return dust
}

func printMap(grid [][]Tile) {
	for _, i := range grid {
		for _, j := range i {
			fmt.Print(string(j))
		}
		fmt.Println()
	}
}

func getNeighbours(x int, y int, grid [][]Tile) []Coord {
	neighbours := []Coord{}
	if x < len(grid[0])-1 {
		neighbours = append(neighbours, Coord{x+1, y})
	}
	if x > 0 {
		neighbours = append(neighbours, Coord{x-1, y})
	}
	if y < len(grid)-1 {
		neighbours = append(neighbours, Coord{x, y+1})
	}
	if y > 0 {
		neighbours = append(neighbours, Coord{x, y-1})
	}
	return neighbours
}

func findIntersectionSum(grid [][]Tile) (sum int) {

	for y, i := range grid {
		for x, _ := range i {
			if grid[y][x] == EMPTY {
				continue
			}
			n := getNeighbours(x, y, grid)
			intersection := true
			for _, coord := range n {
				if len(grid[coord.y]) > 0 && grid[coord.y][coord.x] != WALL {
					intersection = false
				}
			}
			if intersection {
				sum += (x*y)
			}
		}
	}
	return sum
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

	fmt.Println(partOne(candidateProg))

	candidateProg = make([]int64, len(originalProg))
	copy(candidateProg, originalProg)
	bufferSpace = make([]int64, 10000)
	candidateProg = append(candidateProg, bufferSpace...)

	fmt.Println(partTwo(candidateProg))
}

