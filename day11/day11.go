package main

import (
	"fmt"
	"io/ioutil"
	"math"
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

type Robot struct {
	pos Coord
	dir Direction
}

type Computer struct {
	program []int64
	input int64
	relativeBase int64
	pos int64
	finished bool
	outputs []int64
}

type Direction int
const (
	LEFT = iota
	RIGHT
	UP
	DOWN
)

type Coord struct {
	x int
	y int
}

type Paint int
const (
	BLACK = iota
	WHITE
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

func turn(facing Direction, turning Direction) Direction {
	switch facing {
	case UP:
		if turning == LEFT { return LEFT } else { return RIGHT }
	case DOWN:
		if turning == LEFT { return RIGHT } else { return LEFT }
	case LEFT:
		if turning == LEFT { return DOWN } else { return UP }
	case RIGHT:
		if turning == LEFT { return UP } else { return DOWN }
	default:
		panic("Something has gone wrong here")
	}
}

func moveRobot(r *Robot, d Direction) {
	r.dir = turn(r.dir, d)
	switch r.dir {
	case UP:
			r.pos.y -= 1
	case DOWN:
			r.pos.y += 1
	case LEFT:
			r.pos.x -= 1
	case RIGHT:
			r.pos.x += 1
	}
}

func partOne(program []int64) (painted int) {
	c := Computer{program, 0,0, 0, false, []int64{} }
	r := Robot{ Coord{0, 0}, UP }

	visited := make(map[Coord]Paint, 0)

	for !c.finished {
		// set input to current panel
		if _, ok := visited[r.pos]; !ok {
			c.input = 0
		} else {
			c.input = int64(visited[r.pos])
		}

		// run program, get colour and direction
		for !c.finished && len(c.outputs) != 2 {
			cycle(&c)
		}

		if c.finished {
			break
		}

		colour := Paint(c.outputs[0])
		direction := Direction(c.outputs[1])
		c.outputs = []int64{}

		visited[r.pos] = colour
		moveRobot(&r, direction)
	}

	for k, _ := range visited {
		fmt.Println(k)
	}
	return len(visited)
}

func partTwo(program []int64) (painted int) {
	c := Computer{program, 1,0, 0, false, []int64{} }
	r := Robot{ Coord{0, 0}, UP }

	visited := make(map[Coord]Paint, 0)

	for !c.finished {
		// set input to current panel
		if _, ok := visited[r.pos]; !ok {
			c.input = 1
		} else {
			c.input = int64(visited[r.pos])
		}

		// run program, get colour and direction
		for !c.finished && len(c.outputs) != 2 {
			cycle(&c)
		}

		if c.finished {
			break
		}

		colour := Paint(c.outputs[0])
		direction := Direction(c.outputs[1])
		c.outputs = []int64{}

		visited[r.pos] = colour
		moveRobot(&r, direction)
	}

	minX := math.MaxInt64
	maxX := math.MinInt64
	minY := math.MaxInt64
	maxY := math.MinInt64

	for k, _ := range visited {
		if k.x < minX { minX = k.x }
		if k.x > maxX { maxX = k.x }
		if k.y < minY { minY = k.y }
		if k.y > maxY { maxY = k.y }
	}

	region := make([][]Paint, 0)
	for i := 0; i < (maxY+(-1*minY)+1); i++ {
		region = append(region, make([]Paint, maxX+(-1*minX)+1))
	}

	for k, v := range visited {
		cx := k.x + (-1*minX)
		cy := k.y + (-1*minY)
		region[cy][cx] = v
	}

	printit(region)
	return len(visited)
}

func printit(d [][]Paint) {
	for _, v := range d {
		for _, w := range v {
			if w == WHITE {
				fmt.Printf("X")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println("")
	}

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
	fmt.Println(partOne(candidateProg))

	// part 2
	copy(candidateProg, originalProg)
	bufferSpace = make([]int64, 10000)
	candidateProg = append(candidateProg, bufferSpace...)
	partTwo(candidateProg)
}

