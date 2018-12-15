package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

type player struct {
	lastTick int
	plType   cellType
}

type cell struct {
	player     *player
	underlying cellType
}

type cellType byte

func (c cellType) String() string {
	return fmt.Sprintf("%s", string(c))
}

const (
	WALL   cellType = '#'
	GOBLIN cellType = 'G'
	ELF    cellType = 'E'
	EMPTY  cellType = '.'
)

func main() {
	//file, err := os.Open("day15/input.txt")
	file, err := os.Open("day15/test.txt")
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
	track := make([]*cell, maxX*maxY)
	fmt.Printf("Size of the field: %vx%v\n", maxX, maxY)
	for y := 0; y < maxY; y++ {
		line := []byte(lines[y])
		for x := 0; x < maxX; x++ {
			track[linear(x, y, maxX, maxY)] = deduceCell(line[x])
		}
	}
	x, y := part1(maxX, maxY, track)
	fmt.Printf("Solution part 1 is: %v,%v", x, y)
}

func part1(maxX, maxY int, track []*cell) (int, int) {
	maxTick := 3
	for tick := 0; tick < maxTick; tick++ {
		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				//cur := track[linear(x, y, maxX, maxY)]
			}
		}
		printState(tick, maxY, maxX, track)
	}
	return 0, 0
}

func deduceCell(cellContents byte) (result *cell) {
	result = &cell{}
	c := cellType(cellContents)
	switch c {
	case GOBLIN:
		result.underlying = EMPTY
		result.player = &player{
			plType: GOBLIN,
		}
	case ELF:
		result.underlying = EMPTY
		result.player = &player{
			plType: ELF,
		}
	default:
		result.underlying = c
		result.player = nil
	}
	return
}

func linear(x int, y int, maxX int, maxY int) int {
	scaled := x + y*maxX
	if x >= maxX {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	if y >= maxY {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
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

func printState(tick int, maxY int, maxX int, track []*cell) {
	fmt.Printf("Tick: %d\n", tick)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			c := track[linear(x, y, maxX, maxY)]
			if c.player != nil {
				fmt.Printf(string(c.player.plType))
			} else {
				fmt.Printf("%v", c.underlying)
			}
		}
		fmt.Println()
	}
}
