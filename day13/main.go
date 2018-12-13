package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

type cellType byte

func (c cellType) String() string {
	return fmt.Sprintf("%s", string(c))
}

const (
	NS        cellType = '|'
	WE        cellType = '-'
	NE_SW     cellType = '\\'
	SE_NW     cellType = '/'
	NS_WE     cellType = '+'
	DRIVER_NS cellType = 'v'
	DRIVER_EW cellType = '<'
	DRIVER_WE cellType = '>'
	DRIVER_SN cellType = '^'
)

func main() {
	//file, err := os.Open("day13/input.txt")
	file, err := os.Open("day13/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	lines := make([]string, 0)
	for line := ls.Front(); line != nil; line = line.Next() {
		trackLine := line.Value.(string)
		if trackLine != "" {
			lines = append(lines, trackLine)
		}
	}
	maxX := len(lines[0])
	maxY := len(lines)
	track := make([]cellType, maxX*maxY)
	fmt.Printf("Size of the field: %vx%v\n", maxX, maxY)
	for y := 0; y < maxY; y++ {
		line := []byte(lines[y])
		for x := 0; x < maxX; x++ {
			track[linear(x, y, maxX)] = cellType(line[x])
		}
	}
	maxTick := 10
	for tick := 0; tick < maxTick; tick++ {
		fmt.Printf("Tick: %d\n", tick)
		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				fmt.Printf("%v", track[linear(x, y, maxX)])
			}
			fmt.Println()
		}
	}
}

func linear(x int, y int, maxX int) int {
	scaled := x + y*maxX
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
