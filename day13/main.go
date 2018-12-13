package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

type cell struct {
	driver     cellType
	underlying cellType
	lastTick   int
}

type cellType byte

func (c cellType) String() string {
	return fmt.Sprintf("%s", string(c))
}

const (
	NS          cellType = '|'
	WE          cellType = '-'
	NE_SW       cellType = '\\'
	SE_NW       cellType = '/'
	NS_WE       cellType = '+'
	DRIVER_NS   cellType = 'v'
	DRIVER_EW   cellType = '<'
	DRIVER_WE   cellType = '>'
	DRIVER_SN   cellType = '^'
	DRIVER_NONE cellType = '?'
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
	track := make([]*cell, maxX*maxY)
	fmt.Printf("Size of the field: %vx%v\n", maxX, maxY)
	for y := 0; y < maxY; y++ {
		line := []byte(lines[y])
		for x := 0; x < maxX; x++ {
			track[linear(x, y, maxX, maxY)] = deduceCell(line[x])
		}
	}
	x, y := part1(maxX, maxY, track)
	fmt.Printf("Collision on: %v,%v", x, y)
}

func part1(maxX, maxY int, track []*cell) (int, int) {
	maxTick := 1000
	for tick := 0; tick < maxTick; tick++ {
		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				cur := track[linear(x, y, maxX, maxY)]
				var target *cell
				if cur.lastTick >= tick {
					continue
				}
				switch cur.driver {
				case DRIVER_NS:
					target = track[linear(x, y+1, maxX, maxY)]
					switch target.underlying {
					case NE_SW:
						cur.driver = DRIVER_WE
					case SE_NW:
						cur.driver = DRIVER_EW
					}
				case DRIVER_SN:
					target = track[linear(x, y-1, maxX, maxY)]
					switch target.underlying {
					case NE_SW:
						cur.driver = DRIVER_EW
					case SE_NW:
						cur.driver = DRIVER_WE
					}
				case DRIVER_WE:
					target = track[linear(x+1, y, maxX, maxY)]
					switch target.underlying {
					case NE_SW:
						cur.driver = DRIVER_NS
					case SE_NW:
						cur.driver = DRIVER_SN
					}
				case DRIVER_EW:
					target = track[linear(x-1, y, maxX, maxY)]
					switch target.underlying {
					case NE_SW:
						cur.driver = DRIVER_SN
					case SE_NW:
						cur.driver = DRIVER_NS
					}
				default:
					continue
				}
				if target.driver != DRIVER_NONE {
					return x, y
				}
				target.driver = cur.driver
				cur.driver = DRIVER_NONE
				target.lastTick = tick
			}
		}
		printState(tick, maxY, maxX, track)
	}
	return 0, 0
}

func deduceCell(cellContents byte) (result *cell) {
	result = &cell{
		lastTick: 0,
	}
	c := cellType(cellContents)
	switch c {
	case DRIVER_NS:
		result.underlying = NS
		result.driver = c
	case DRIVER_SN:
		result.underlying = NS
		result.driver = c
	case DRIVER_WE:
		result.underlying = WE
		result.driver = c
	case DRIVER_EW:
		result.underlying = WE
		result.driver = c
	default:
		result.underlying = c
		result.driver = DRIVER_NONE
	}
	return
}

func printState(tick int, maxY int, maxX int, track []*cell) {
	fmt.Printf("Tick: %d\n", tick)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			c := track[linear(x, y, maxX, maxY)]
			if c.driver != DRIVER_NONE {
				fmt.Printf("%v", c.driver)
			} else {
				fmt.Printf("%v", c.underlying)
			}
		}
		fmt.Println()
	}
}

func linear(x int, y int, maxX int, maxY int) int {
	scaled := x + y*maxX
	if x >= maxX {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	if y >= maxY {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	//fmt.Printf("Scaled %d:%d to %d\n", x, y, scaled)
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
