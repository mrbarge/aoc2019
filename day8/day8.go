package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
)

const width int = 25
const height int = 6

type Color int
const (
	NONE = -1
	BLACK = 0
	WHITE = 1
	TRANSPARENT = 2
)

func printLayer(l [width][height]int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if l[x][y] == BLACK {
				fmt.Print(" ")
			} else if l[x][y] == WHITE {
				fmt.Print("#")
			} else {
				fmt.Print("_")
			}
		}
		fmt.Println()
	}
}

func initialiseLayer(l *[width][height]int) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			l[x][y] = NONE
		}
	}
}

func countDigit(l [width][height]int, n int) (numDigits int) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if l[x][y] == n {
				numDigits += 1
			}
		}
	}
	return numDigits
}

func main() {

	bd, _ := ioutil.ReadFile("input.txt")
	s := string(bd)

	layers := make([][width][height]int, 0)
	makeLayer := true
	x, y := 0, 0

	for _, c := range s {

		ir, err := strconv.Atoi(string(c))
		if err != nil {
			log.Fatal(err)
		}

		// make a new layer if needed
		if makeLayer {
			layers = append(layers, [width][height]int{})
			makeLayer = false
		}

		layers[len(layers)-1][x][y] = ir
		x += 1; y += 1
		if x == width && y == height {
			makeLayer = true
			x, y = 0, 0
		} else {
			x = x % width
			y = y % height
		}

	}

	sort.Slice(layers, func(i, j int) bool { return countDigit(layers[i], 0) < countDigit(layers[j], 0) })
	fmt.Println(countDigit(layers[0], 1) * countDigit(layers[0], 2))

	// part 2 - just read it backwards
	messageLayer := [width][height]int{}
	initialiseLayer(&messageLayer)

	x, y = width-1, height-1
	for i := len(s)-1; i >= 0; i-- {
		ir, _ := strconv.Atoi(string(s[i]))

		if ir == BLACK || ir == WHITE {
			messageLayer[x][y] = ir
		}
		x -= 1
		if x < 0 {
			x = width-1
			y -= 1
		}
		if y < 0 {
			y = height-1
		}
	}
	printLayer(messageLayer)
}
