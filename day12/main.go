package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	//file, err := os.Open("day12/input.txt")
	file, err := os.Open("day12/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// initial initialStateString: #..#.#..##......###...###
	regexInitial, err := regexp.Compile("initial state: ([\\.#]+)")
	if err != nil {
		log.Fatal(err)
	}
	regexOperation, err := regexp.Compile("([\\.#]+) => ([\\.#])")
	if err != nil {
		log.Fatal(err)
	}

	line := ls.Front()
	matchesInitial := regexInitial.FindAllStringSubmatch(line.Value.(string), -1)[0]
	initialStateString := matchesInitial[1]
	//fmt.Println(initialStateString)
	state := make([]byte, len(initialStateString))
	for i, s := range initialStateString {
		state[i] = toByte(s)
	}
	fmt.Printf("Initially: %+v\n", state)

	line = line.Next().Next()
	ops := make(map[string]byte)
	for ; line != nil; line = line.Next() {
		lineValue := line.Value.(string)
		fmt.Printf("Analyzing line %s\n", lineValue)
		matches := regexOperation.FindAllStringSubmatch(lineValue, -1)[0]
		op := toByte(int32(matches[2][0]))
		stateMatch := matches[1]
		ops[stateMatch] = op
		fmt.Printf("Found stateTransition: %+v -> %d\n", stateMatch, op)
	}

	//maxValue, solutionX, solutionY := part1(maxX, maxY, serialNumber, 3)
	//fmt.Printf("Solution is: %d at coordinate %d,%d", maxValue, solutionX, solutionY)
	//maxValue2, solutionX2, solutionY2, blockSize2 := part2(maxX, maxY, serialNumber)
	//fmt.Printf("Solution 2 is: %d at coordinate %d,%d,%d", maxValue2, solutionX2, solutionY2, blockSize2)
}

func toByte(s int32) byte {
	if s == int32('#') {
		return 1
	} else {
		return 0
	}
}

func readAll(file *os.File, list *list.List) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := scanner.Text()
		list.PushBack(val)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
