package main

import (
	"errors"
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

type Tile int
const (
	WALL = iota
	EMPTY
	UNEXPLORED
	OXYGEN
)

type Direction int
const (
	NORTH = iota + 1
	SOUTH
	WEST
	EAST
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

func updatePos(c Coord, d Direction) Coord {
	switch d {
	case NORTH:
		return Coord{c.x, c.y - 1}
	case SOUTH:
		return Coord{c.x, c.y + 1}
	case EAST:
		return Coord{c.x + 1, c.y }
	case WEST:
		return Coord{c.x - 1, c.y }
	default:
		return c
	}
}

func selectDirection(c Coord, world map[Coord]Tile) (Direction, error) {
	// north?
	if _, ok := world[Coord{c.x, c.y-1}]; !ok {
		return NORTH, nil
	}
	// east?
    if _, ok := world[Coord{c.x+1, c.y}]; !ok {
		return EAST, nil
	}
	// south?
	if _, ok := world[Coord{c.x, c.y+1}]; !ok {
		return SOUTH, nil
	}
	// west?
	if _, ok := world[Coord{c.x-1, c.y}]; !ok {
		return WEST, nil
	}
	// can't move
	//fmt.Println("Can't move")
	return -1, errors.New("Can't move anywhere")
}

func backtrack(d Direction) Direction {
	if d == NORTH { return SOUTH } else if d == EAST { return WEST } else if d == WEST { return EAST } else { return NORTH }
}

func partOne(program []int64) {
	c := Computer{program, []int64{},0, 0, false, []int64{}, false}

	// state of the world
	grid := make(map[Coord]Tile, 0)
	// keep track of where we've been
	seen := []Direction{}
	// keep track of where we are
	pos := Coord{0,0}
	grid[pos] = EMPTY
	// go north first
	lastDirection := Direction(NORTH)
	addInput(&c, int64(lastDirection))

	steps := 0
	stepsUntilOxygen := 0
	finished := false

	for !finished {
		// run program, get colour and direction
		for !c.finished && !c.inputBlocked && len(c.outputs) == 0 {
			cycle(&c)
		}

		if c.finished {
			break
		}

		moveResult := popOutput(&c)
		switch moveResult {
		case 0:
			// couldn't move, flag this as a wall
			nextPos := updatePos(pos, lastDirection)
			grid[nextPos] = WALL

		case 1:
			// moved successfully
			pos = updatePos(pos, lastDirection)
			// flag we've visited here, if not a backtrack
			if _, ok := grid[pos]; !ok {
				seen = append(seen, lastDirection)
				steps += 1
			}
			grid[pos] = EMPTY

		case 2:
			// we're good!
			pos = updatePos(pos, lastDirection)
			if _, ok := grid[pos]; !ok {
				grid[pos] = OXYGEN
				steps += 1
				stepsUntilOxygen = steps
			}
		}

		// choose a new destination
		nextDirection, err := selectDirection(pos, grid)
		if err == nil {
			lastDirection = nextDirection
		} else {
			// move backwards

			// can we? if we can't, we're done boyyy
			if len(seen) == 0 {
				finished = true
			} else {
				lastDirection = backtrack(seen[len(seen)-1])
				seen = seen[:len(seen)-1]
				steps -= 1
			}
		}
		addInput(&c, int64(lastDirection))
	}

	// part 1 answer
	fmt.Println(stepsUntilOxygen)

	minutes := 0
	stillSpreading := true

	for stillSpreading {
		stillSpreading = false
		// don't spread until we know every possible spread space
		spreads := make([]Coord, 0)
		for spos, space := range grid {
			// if this is oxygen
			if space == OXYGEN {
				// find what empty neighbours it could spread to
				neighs := getNeighbours(spos, grid)
				for _, npos := range neighs {
					// they now have oxygen, flag that we're still spreading
					spreads = append(spreads, npos)
					stillSpreading = true
				}
			}
		}
		for _, spos := range spreads {
			grid[spos] = OXYGEN
		}
		// if no oxygen spread, we must be full
		minutes += 1
	}

	fmt.Println(minutes-1)
}

func getUnexploredNeighbours(c Coord, grid map[Coord]Tile) []Coord {
	possibleCoords := []Coord{updatePos(c, NORTH), updatePos(c, SOUTH), updatePos(c, WEST), updatePos(c, EAST)}
	retCoords := []Coord{}
	for _, p := range possibleCoords {
		if _, ok := grid[p]; !ok {
			retCoords = append(retCoords, p)
		}
	}
	return retCoords
}

func getNeighbours(c Coord, grid map[Coord]Tile) []Coord {
	possibleCoords := []Coord{updatePos(c, NORTH), updatePos(c, SOUTH), updatePos(c, WEST), updatePos(c, EAST)}
	retCoords := []Coord{}
	for _, p := range possibleCoords {
		if _, ok := grid[p]; (ok && grid[p] == EMPTY) {
			retCoords = append(retCoords, p)
		}
	}
	return retCoords
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

	// part 1
	bufferSpace := make([]int64, 10000)
	candidateProg = append(candidateProg, bufferSpace...)
	partOne(candidateProg)

}

