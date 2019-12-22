package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func buildDeck(size int) []int {
	deck := make([]int, size)
	for i := 0; i < size; i++ {
		deck[i] = i
	}
	return deck
}

func dealNewStack(cards []int) []int {
	r := make([]int, len(cards))
	for i, j := len(cards)-1, 0; i >= 0; i-- {
		r[i] = cards[j]
		j++
	}
	return r
}

func cutCards(n int, cards []int) []int {
	r := make([]int, 0)
	if n < 0 {
		r = append(r, cards[len(cards)+n:]...)
		r = append(r, cards[0:(len(cards)+n)]...)
	} else {
		r = append(r, cards[n:]...)
		r = append(r, cards[0:n]...)
	}
	return r
}

func dealIncrement(n int, cards []int) []int {
	numCards := len(cards)
	r := make([]int, numCards)
	count := 0
	insertPos := 0
	for count < numCards {
		r[insertPos] = cards[count]
		insertPos = (insertPos + n) % numCards
		count++
	}
	return r
}

func main() {

	cards := buildDeck(10007)
	//cards := buildDeck(10)

	//file, _ := os.Open("test.txt")
	file, _ := os.Open("input.txt")
	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "cut ") {
			sCutField := strings.Split(line, " ")[1]
			cutField, _ := strconv.Atoi(sCutField)
			cards = cutCards(cutField, cards)
		} else if strings.HasPrefix(line, "deal with increment") {
			sDealField := strings.Split(line, " ")[3]
			dealField, _ := strconv.Atoi(sDealField)
			cards = dealIncrement(dealField, cards)
		} else if strings.HasPrefix(line, "deal into new stack") {
			cards = dealNewStack(cards)
		} else {
			fmt.Println("Unknown line ", line)
		}
	}

	for i := range cards {
		if cards[i] == 2019 {
			fmt.Println(i)
		}
	}
}
