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
	inputs []int64
	relativeBase int64
	pos int64
	finished bool
	outputs []int64
	inputBlocked bool
}

type Game struct {
	ball Coord
	paddle Coord
	grid [][]Tile
}

type Tile int
const (
	EMPTY = iota
	WALL
	BLOCK
	PADDLE
	BALL
)

type Coord struct {
	x int
	y int
}

type JoystickState int
const (
	NEUTRAL = iota
	LEFT = -1
	RIGHT = 1
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
		fmt.Println("Store time")
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
		fmt.Println("Output time")
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

func partOne(program []int64) {
	c := Computer{program, []int64{},0, 0, false, []int64{}, false}

	grid := make([][]Tile, 0)
	for i := 0; i < 1000; i++ {
		r := make([]Tile, 1000)
		grid = append(grid, r)
	}

	for !c.finished {
		// run program, get colour and direction
		for !c.finished && len(c.outputs) != 3 {
			cycle(&c)
		}

		if c.finished {
			break
		}

		xpos := c.outputs[0]
		ypos := c.outputs[1]
		tile := Tile(c.outputs[2])


		c.outputs = []int64{}
		grid[xpos][ypos] = tile
	}

	count := 0
	for _, r := range grid {
		for _, v := range r {
			if v == BLOCK {
				count += 1
			}
		}
	}
	fmt.Println(count)
}

func partTwo(program []int64) (score int64) {
	c := Computer{program, []int64{},0, 0, false, []int64{}, false }

	grid := make([][]Tile, 0)
	for i := 0; i < 1000; i++ {
		r := make([]Tile, 1000)
		grid = append(grid, r)
	}

	game := Game{Coord{0,0}, Coord{0, 0}, grid }

	for !c.finished {
		// run program, get colour and direction
		for !c.finished && !c.inputBlocked && len(c.outputs) != 3 {
			cycle(&c)
		}

		// are we done?
		if c.finished {
			break
		}

		// are we input blocked?
		if c.inputBlocked {
			// feed it a smart input based on where the ball is
			if game.ball.x < game.paddle.x {
				// need to move left
				c.inputs = append(c.inputs, LEFT)
			} else if game.ball.x > game.paddle.x {
				// need to move right
				c.inputs = append(c.inputs, RIGHT)
			} else {
				// keep it cool
				c.inputs = append(c.inputs, NEUTRAL)
			}
			c.inputBlocked = false
			fmt.Println("Giving input")
			continue
		}

		// otherwise we must have outputs
		xpos := c.outputs[0]
		ypos := c.outputs[1]
		tile := Tile(c.outputs[2])

		if xpos == -1 && ypos == 0 {
			score = c.outputs[2]
		} else {
			c.outputs = []int64{}
			grid[xpos][ypos] = tile

			if tile == BALL {
				game.ball.x = int(xpos)
				game.ball.y = int(ypos)
			} else if tile == PADDLE {
				game.paddle.x = int(xpos)
				game.paddle.y = int(ypos)
			}
		}
	}

	return score
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

	// part 2
	bufferSpace = make([]int64, 10000)
	copy(candidateProg, originalProg)
	candidateProg = append(candidateProg, bufferSpace...)
	candidateProg[0] = 2
	partTwo(candidateProg)

}

