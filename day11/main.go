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

type point struct {
	x         int
	y         int
	velocityX int
	velocityY int
}

func main() {
	file, err := os.Open("day11/input.txt")
	//file, err := os.Open("day11/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// x,y,serial
	regex, err := regexp.Compile("(\\d+),(\\d+),(\\d+)")
	if err != nil {
		log.Fatal(err)
	}

	var maxX = 0
	var maxY = 0
	var serialNumber = 0
	for line := ls.Front(); line != nil; line = line.Next() {
		matches := regex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		maxX, err = strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		maxY, err = strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		serialNumber, err = strconv.Atoi(matches[3])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found maxX=%d, maxY=%d, serialNumber=%d\n", maxX, maxY, serialNumber)
	}

	maxValue, solutionX, solutionY := part1(maxX, maxY, serialNumber, 3)
	fmt.Printf("Solution is: %d at coordinate %d,%d", maxValue, solutionX, solutionY)
}

func part1(maxX, maxY, serialNumber, blockSize int) (maxValue, solutionX, solutionY int) {
	field := make([]int, maxX*maxY)
	for y := 1; y <= maxY; y++ {
		for x := 1; x <= maxX; x++ {
			field[linear(x, y, maxX)] = calculate(x, y, serialNumber)
		}
	}

	//for y := 1; y <= maxY; y++ {
	//	for x := 1; x <= maxX; x++ {
	//		fmt.Printf("%3d", field[linear(x, y, maxX)])
	//	}
	//	fmt.Println()
	//}
	for y := blockSize/2 + 1; y < maxY-blockSize/2; y++ {
		for x := blockSize/2 + 1; x < maxX-blockSize/2; x++ {
			v := field[linear(x-1, y-1, maxX)] + field[linear(x, y-1, maxX)] + field[linear(x+1, y-1, maxX)]
			v += field[linear(x-1, y, maxX)] + field[linear(x, y, maxX)] + field[linear(x+1, y, maxX)]
			v += field[linear(x-1, y+1, maxX)] + field[linear(x, y+1, maxX)] + field[linear(x+1, y+1, maxX)]
			if v > maxValue {
				maxValue = v
				solutionX = x - 1
				solutionY = y - 1
				//fmt.Printf("New maximum found at: %d,%d (%d)\n", x-1, y-1, maxValue)
			}
		}
	}
	return
}

func calculate(x int, y int, serialNumber int) int {
	rackId := x + 10
	powerLevelStart := rackId * y
	powerLevel := powerLevelStart + serialNumber
	value := powerLevel * rackId
	valueAsStr := strconv.Itoa(value)
	var valueBeforeSubtraction uint8 = 0
	if len(valueAsStr) >= 3 {
		valueBeforeSubtraction = valueAsStr[len(valueAsStr)-3] - '0'
	}
	finalValue := int(valueBeforeSubtraction) - 5
	return finalValue
}

func linear(x int, y int, maxX int) int {
	scaled := (x - 1) + (y-1)*maxX
	//fmt.Printf("Scaled %d:%d to %d", x, y, scaled)
	return scaled
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
