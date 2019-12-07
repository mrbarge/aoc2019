package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Amplifier struct {
	program []int64
	pos int
	input []int64
	inputPos int
	output int64
	feedback bool
	next *Amplifier
	state AmplifierState
	id int
}

type AmplifierState int
const (
	RUN = iota
	EMIT
	BLOCK
	FINISH
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

func cycle(amp *Amplifier) (output int64, finished bool) {

	op := parseOpCode(amp.program[amp.pos])
	modes := parseModes(amp.program[amp.pos])

	// set amp as running
	amp.state = RUN
	switch op {

	case ADDITION:
		value1 := getPositionOrImmediate(amp.program, modes[0], amp.program[amp.pos+1])
		value2 := getPositionOrImmediate(amp.program, modes[1], amp.program[amp.pos+2])
		posDest := amp.program[amp.pos+3]
		amp.program[posDest] = value1 + value2
		amp.pos = amp.pos + 4

	case MULTIPLY:
		value1 := getPositionOrImmediate(amp.program, modes[0], amp.program[amp.pos+1])
		value2 := getPositionOrImmediate(amp.program, modes[1], amp.program[amp.pos+2])
		posDest := amp.program[amp.pos+3]
		amp.program[posDest] = value1 * value2
		amp.pos = amp.pos + 4

	case STORE:
		posDest := amp.program[amp.pos+1]
		if amp.inputPos > len(amp.input)-1 {
			// waiting for new input
			amp.state = BLOCK
		} else {
			amp.program[posDest] = amp.input[amp.inputPos]
			amp.inputPos += 1
			amp.pos = amp.pos + 2
		}

	case OUTPUT:
		posDest := getPositionOrImmediate(amp.program, modes[0], amp.program[amp.pos+1])
		amp.output = posDest
		amp.pos = amp.pos + 2
		amp.state = EMIT

	case JIT:
		value1 := getPositionOrImmediate(amp.program, modes[0], amp.program[amp.pos+1])
		value2 := getPositionOrImmediate(amp.program, modes[1], amp.program[amp.pos+2])
		if value1 != 0 {
			amp.pos = int(value2)
		} else {
			amp.pos += 3
		}

	case JIF:
		value1 := getPositionOrImmediate(amp.program, modes[0], amp.program[amp.pos+1])
		value2 := getPositionOrImmediate(amp.program, modes[1], amp.program[amp.pos+2])
		if value1 == 0 {
			amp.pos = int(value2)
		} else {
			amp.pos += 3
		}

	case LT:
		value1 := getPositionOrImmediate(amp.program, modes[0], amp.program[amp.pos+1])
		value2 := getPositionOrImmediate(amp.program, modes[1], amp.program[amp.pos+2])
		posDest := amp.program[amp.pos+3]
		if value1 < value2 {
			amp.program[posDest] = 1
		} else {
			amp.program[posDest] = 0
		}
		amp.pos += 4

	case EQ:
		value1 := getPositionOrImmediate(amp.program, modes[0], amp.program[amp.pos+1])
		value2 := getPositionOrImmediate(amp.program, modes[1], amp.program[amp.pos+2])
		posDest := amp.program[amp.pos+3]
		if value1 == value2 {
			amp.program[posDest] = 1
		} else {
			amp.program[posDest] = 0
		}
		amp.pos += 4

	case QUIT:
		amp.state = FINISH
		finished = true

	default:
		fmt.Printf("Unknown op code: %d\n", amp.program[amp.pos])
	}

	return amp.output, finished
}

func runUntilInterrupt(amp *Amplifier) (signal int64, done bool) {
	done = false
	for !done && amp.state == RUN {
		signal, done = cycle(amp)
	}
	return signal, done
}

func simulation(program []int64, phaseSequence []int, feedback bool) int64 {

	// create our amplifiers
	amps := make([]*Amplifier, 0)
	for i, phase := range phaseSequence {
		// clone the program
		candidateProg := make([]int64, len(program))
		copy(candidateProg, program)

		// create the amplifier
		amp := Amplifier{candidateProg, 0, []int64{int64(phase)}, 0,  0,false, nil,RUN,i}
		amps = append(amps, &amp)
	}

	// chain them together in a loop
	for i := 0; i < len(amps); i++ {
		if i == len(amps)-1 {
			if feedback {
				amps[i].next = amps[0]
			}
		} else {
			amps[i].next = amps[i+1]
		}
	}

	currentAmp := amps[0]
	// 0 signal for first amp
	currentAmp.input = append(currentAmp.input, 0)

	finished := false
	signal := int64(0)
	for ; currentAmp != nil && !finished; {
		ampquit := false
		signal, ampquit = runUntilInterrupt(currentAmp)
		// reset amp for next time
		currentAmp.state = RUN

		if ampquit && currentAmp.id == 4 {
			finished = true
		}

		// move on to next amp
		currentAmp = currentAmp.next
		if currentAmp != nil {
			currentAmp.input = append(currentAmp.input, signal)
		}


	}

	return signal
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

	inputPerms := permutations([]int{0,1,2,3,4})

	thrusterSignals := make([]int64, 0)
	for _, inputPerm := range inputPerms {
		thrust := simulation(originalProg, inputPerm, false)
		thrusterSignals = append(thrusterSignals, thrust)
	}

	sort.Slice(thrusterSignals, func(i, j int) bool { return thrusterSignals[i] < thrusterSignals[j] })
	fmt.Println(thrusterSignals[len(thrusterSignals)-1])

	// Feedback loop mode begins
	feedbackPerms := permutations([]int{5,6,7,8,9})
	thrusterSignals = make([]int64, 0)
	for _, feedbackPerm := range feedbackPerms {
		thrust := simulation(originalProg, feedbackPerm, true)
		thrusterSignals = append(thrusterSignals, thrust)
	}

	sort.Slice(thrusterSignals, func(i, j int) bool { return thrusterSignals[i] < thrusterSignals[j] })
	fmt.Println(thrusterSignals[len(thrusterSignals)-1])

}

