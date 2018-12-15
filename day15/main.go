package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"

	"github.com/RyanCarrier/dijkstra"
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
	maxTick := 4
	for tick := 0; tick < maxTick; tick++ {
		fmt.Println("######################### TICK:", tick)
		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				cur := track[linear(x, y, maxX, maxY)]
				if cur.player == nil || cur.player.lastTick >= tick {
					continue
				}
				if newX, newY, ok := deduceNextCell(x, y, maxX, maxY, track); ok {
					fmt.Printf("%v moves from %d,%d (%v) to %d,%d (%v)\n", cur.player.plType, x, y, linear(x, y, maxX, maxY), newX, newY, linear(newX, newY, maxX, maxY))
					newLoc := track[linear(newX, newY, maxX, maxY)]
					newLoc.player = cur.player
					newLoc.player.lastTick = tick
					cur.player = nil
				}
			}
		}
		printState(tick, maxY, maxX, track)
	}
	return 0, 0
}

func deduceNextCell(oldX int, oldY int, maxX int, maxY int, track []*cell) (newX, newY int, found bool) {
	idsOfEnemies := make([]int, 0)
	oldId := linear(oldX, oldY, maxX, maxY)
	old := track[oldId]
	graph := dijkstra.NewGraph()
	graph.AddVertex(oldId)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			id := linear(x, y, maxX, maxY)
			iter := track[id]
			if iter.underlying == WALL {
				continue
			}
			if old.player != nil && iter.player != nil && old.player.plType != iter.player.plType {
				idsOfEnemies = append(idsOfEnemies, id)
			}
			if iter.player == nil || (iter.player != nil && old.player.plType != iter.player.plType) {
				//fmt.Printf("Adding vertex %v,%v (%v)\n", x, y, id)
				graph.AddVertex(id)
			}
		}
	}
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			id := linear(x, y, maxX, maxY)
			iter := track[id]
			if iter.underlying == WALL {
				continue
			}
			maybeReport(x, y-1, maxX, maxY, track, graph, id, oldId, 0)
			maybeReport(x-1, y, maxX, maxY, track, graph, id, oldId, 1)
			maybeReport(x+1, y, maxX, maxY, track, graph, id, oldId, 10)
			maybeReport(x, y+1, maxX, maxY, track, graph, id, oldId, 100)
		}
	}
	var bestPath *dijkstra.BestPath = nil
	for _, idOfEnemy := range idsOfEnemies {
		fmt.Printf("shortest between %v->%v\n", oldId, idOfEnemy)
		path, err := graph.Shortest(oldId, idOfEnemy)
		if err != nil {
			return -1, -1, false
		}
		if bestPath == nil {
			bestPath = &path
		} else if path.Distance < bestPath.Distance {
			bestPath = &path
		} else if path.Distance == bestPath.Distance {
			if bestPath.Path[0] > path.Path[0] {
				bestPath = &path
			}
		}
	}
	if bestPath == nil {
		fmt.Printf("No path found from %d:%d\n", oldX, oldY)
		return -1, -1, false
	}
	if len(bestPath.Path) == 2 {
		fmt.Println("Ready to attack!")
		return -1, -1, false
	}
	firstStep := bestPath.Path[1]
	return firstStep % maxX, firstStep / maxX, true
}

func maybeReport(x int, y int, maxX int, maxY int, track []*cell, graph *dijkstra.Graph, id, oldId int, price int64) {
	curr := track[id]
	source := track[oldId]
	if neighborId, ok := maybeGet(x, y, maxX, maxY); ok {
		neighbor := track[neighborId]
		if neighbor.underlying == EMPTY && curr.underlying == EMPTY {
			if neighbor.player != nil && neighbor.player.plType == source.player.plType {
				return
			}
			//fmt.Printf("Adding arc %v->%v (price %d)\n", id, neighborId, price)
			err := graph.AddArc(id, neighborId, price)
			if err != nil {
				if err.Error() == "Source/Destination not found" {
					fmt.Printf("Source(%v)/Destination(%v) not found\n", id, neighborId)
					return
				} else {
					log.Fatalf("Could not add arc between %v and %v: %v", neighborId, id, err)
				}
			}
		}
	}
}

func maybeGet(x int, y int, maxX int, maxY int) (int, bool) {
	if x < 0 || y < 0 || x >= maxX || y >= maxY {
		return -1, false
	}
	return linear(x, y, maxX, maxY), true
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
