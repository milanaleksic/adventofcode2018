package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

type turn int

const (
	left turn = iota
	straight
	right
)

type driver struct {
	direction cellType
	lastTick  int
	lastTurn  turn
}

type cell struct {
	driver     *driver
	underlying cellType
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
	file, err := os.Open("day13/input.txt")
	//file, err := os.Open("day13/test.txt")
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
				if cur.driver.lastTick >= tick {
					continue
				}
				switch cur.driver.direction {
				case DRIVER_NS:
					target = track[linear(x, y+1, maxX, maxY)]
					if target.driver.direction != DRIVER_NONE {
						return x, y + 1
					}
					switch target.underlying {
					case NE_SW:
						cur.driver.direction = DRIVER_WE
					case SE_NW:
						cur.driver.direction = DRIVER_EW
					case NS_WE:
						nextTurn := nextTurn(cur.driver.lastTurn)
						switch nextTurn {
						case left:
							cur.driver.direction = DRIVER_WE
						case right:
							cur.driver.direction = DRIVER_EW
						}
						cur.driver.lastTurn = nextTurn
					}
				case DRIVER_SN:
					target = track[linear(x, y-1, maxX, maxY)]
					if target.driver.direction != DRIVER_NONE {
						return x, y - 1
					}
					switch target.underlying {
					case NE_SW:
						cur.driver.direction = DRIVER_EW
					case SE_NW:
						cur.driver.direction = DRIVER_WE
					case NS_WE:
						nextTurn := nextTurn(cur.driver.lastTurn)
						switch nextTurn {
						case left:
							cur.driver.direction = DRIVER_EW
						case right:
							cur.driver.direction = DRIVER_WE
						}
						cur.driver.lastTurn = nextTurn
					}
				case DRIVER_WE:
					target = track[linear(x+1, y, maxX, maxY)]
					if target.driver.direction != DRIVER_NONE {
						return x + 1, y
					}
					switch target.underlying {
					case NE_SW:
						cur.driver.direction = DRIVER_NS
					case SE_NW:
						cur.driver.direction = DRIVER_SN
					case NS_WE:
						nextTurn := nextTurn(cur.driver.lastTurn)
						switch nextTurn {
						case left:
							cur.driver.direction = DRIVER_SN
						case right:
							cur.driver.direction = DRIVER_NS
						}
						cur.driver.lastTurn = nextTurn
					}
				case DRIVER_EW:
					target = track[linear(x-1, y, maxX, maxY)]
					if target.driver.direction != DRIVER_NONE {
						return x - 1, y
					}
					switch target.underlying {
					case NE_SW:
						cur.driver.direction = DRIVER_SN
					case SE_NW:
						cur.driver.direction = DRIVER_NS
					case NS_WE:
						nextTurn := nextTurn(cur.driver.lastTurn)
						switch nextTurn {
						case left:
							cur.driver.direction = DRIVER_NS
						case right:
							cur.driver.direction = DRIVER_SN
						}
						cur.driver.lastTurn = nextTurn
					}
				default:
					continue
				}
				target.driver.direction = cur.driver.direction
				target.driver.lastTurn = cur.driver.lastTurn
				target.driver.lastTick = tick
				cur.driver.direction = DRIVER_NONE
			}
		}
		//printState(tick, maxY, maxX, track)
	}
	return 0, 0
}

func deduceCell(cellContents byte) (result *cell) {
	result = &cell{
		driver: &driver{
			lastTurn: right,
		},
	}
	c := cellType(cellContents)
	switch c {
	case DRIVER_NS:
		result.underlying = NS
		result.driver.direction = c
	case DRIVER_SN:
		result.underlying = NS
		result.driver.direction = c
	case DRIVER_WE:
		result.underlying = WE
		result.driver.direction = c
	case DRIVER_EW:
		result.underlying = WE
		result.driver.direction = c
	default:
		result.underlying = c
		result.driver.direction = DRIVER_NONE
	}
	return
}

func printState(tick int, maxY int, maxX int, track []*cell) {
	fmt.Printf("Tick: %d\n", tick)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			c := track[linear(x, y, maxX, maxY)]
			if c.driver.direction != DRIVER_NONE {
				fmt.Printf("%v", c.driver.direction)
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

func nextTurn(lastTurn turn) (result turn) {
	result = lastTurn + 1
	if result == right+1 {
		result = left
	}
	return
}
