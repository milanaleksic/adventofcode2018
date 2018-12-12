package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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
	for g := 1; g <= 20; g++ {
		nextGen := make([]byte, len(state)+4) // adding 2 left and right
		for i := 0; i < len(state); i++ {
			for stateMatch, op := range ops {
				stateBytes := make([]byte, 0)
				for j := i - 2; j <= i+2; j++ {
					var b byte
					if j < 0 || j >= len(state) {
						b = 0
					} else {
						b = state[j]
					}
					stateBytes = append(stateBytes, b)
				}
				match := true
				for j := 0; j < 5; j++ {
					stateMatchByte := toByte(int32(stateMatch[j]))
					thisStateByte := stateBytes[j]
					if stateMatchByte != thisStateByte {
						match = false
						break
					}
				}
				if match {
					//fmt.Printf("Applying op %d on index %d since match was ok: %s\n", op, i, stateMatch)
					nextGen[i+2] = op
				}
			}
		}
		state = nextGen
		fmt.Printf("Generation %02d: ", g)
		for i := 0 + (g-1)*2; i < len(state); i++ {
			if state[i] == 1 {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	solution := 0
	for i := 0; i < len(state); i++ {
		if state[i] == 1 {
			solution += i - 40
		}
	}
	fmt.Printf("Solution for 1 is: %d", solution)

	//maxValue, solutionX, solutionY := part1(maxX, maxY, serialNumber, 3)
	//fmt.Printf("Solution is: %d at coordinate %d,%d", maxValue, solutionX, solutionY)
	//maxValue2, solutionX2, solutionY2, blockSize2 := part2(maxX, maxY, serialNumber)
	//fmt.Printf("Solution 2 is: %d at coordinate %d,%d,%d", maxValue2, solutionX2, solutionY2, blockSize2)
}

func toByte(s int32) byte {
	if s == int32('#') {
		return 1
	} else if s == int32('.') {
		return 0
	} else {
		panic("unexpected input!: " + strconv.Itoa(int(s)))
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
