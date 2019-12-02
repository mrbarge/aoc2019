package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type OpCode int
const (
	ADDITION = iota + 1
	MULTIPLY
	QUIT = 99
)

func cycle(program []int, pos int) (int, bool) {
	finished := false
	switch program[pos] {
	case ADDITION:
		pos1 := program[pos+1]
		pos2 := program[pos+2]
		posDest := program[pos+3]
		program[posDest] = program[pos1] + program[pos2]
		pos = pos + 4
	case MULTIPLY:
		pos1 := program[pos+1]
		pos2 := program[pos+2]
		posDest := program[pos+3]
		program[posDest] = program[pos1] * program[pos2]
		pos = pos + 4
	case QUIT:
		finished = true
	default:
		fmt.Printf("Unknown op code: %d\n", program[pos])
	}

	return pos, finished
}

func runUntilHalt(program []int) int {
	done := false
	pos := 0
	for !done {
		pos, done = cycle(program, pos)
	}
	return program[0]
}

func main() {

	partOne := false

	bd, err := ioutil.ReadFile("input.txt")
	if err != nil {
		os.Exit(1)
	}
	sProg := strings.Split(string(bd), ",")
	originalProg := make([]int, len(sProg))
	for i, v := range sProg {
		originalProg[i], _ = strconv.Atoi(v)
	}

	if partOne {
		originalProg[1] = 12
		originalProg[2] = 2
		output := runUntilHalt(originalProg)
		fmt.Println(output)
	} else {
		done := false
		noun, verb := 0, 0
		for noun <= 99 && !done {
			for verb <= 99 && !done {
				candidateProg := make([]int, len(originalProg))
				copy(candidateProg, originalProg)
				candidateProg[1] = noun
				candidateProg[2] = verb
				result := runUntilHalt(candidateProg)
				if result == 19690720 {
					done = true
				} else {
					verb += 1
				}
			}
			if !done {
				noun += 1
				verb = 0
			}
		}
		fmt.Printf("Noun: %d, Verb: %d\n", noun, verb)
	}

}
