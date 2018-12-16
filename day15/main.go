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
	lastTick    int
	plType      cellType
	hitPoints   int
	attackPower int
}

type cell struct {
	x          int
	y          int
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
	file, err := os.Open("day15/input.txt")
	//file, err := os.Open("day15/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	maxX, maxY, track := fromInput(ls)
	endTick, solution1 := part1(maxX, maxY, track)
	fmt.Printf("Solution part 1 is: %v (on tick %v)", solution1, endTick)
}

func fromInput(ls *list.List) (int, int, []*cell) {
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
			track[linear(x, y, maxX, maxY)] = deduceCell(x, y, line[x])
		}
	}
	return maxX, maxY, track
}

func part1(maxX, maxY int, track []*cell) (maxReachedTick int, result int) {
	maxTick := 500
	for tick := 0; tick < maxTick; tick++ {
		fmt.Println("######################### TICK:", tick)
		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				cur := track[linear(x, y, maxX, maxY)]
				if cur.player == nil || cur.player.lastTick >= tick {
					continue
				}
				if newX, newY, ok := deduceNextCell(x, y, maxX, maxY, track); ok {
					//fmt.Printf("%v moves from %d,%d (%v) to %d,%d (%v)\n", cur.player.plType, x, y, linear(x, y, maxX, maxY), newX, newY, linear(newX, newY, maxX, maxY))
					newLoc := track[linear(newX, newY, maxX, maxY)]
					newLoc.player = cur.player
					newLoc.player.lastTick = tick
					cur.player = nil
					attackPhase(newX, newY, maxX, maxY, track)
				} else {
					kindsOfEnemies := 0
					for y := 0; y < maxY; y++ {
						for x := 0; x < maxX; x++ {
							iter := track[linear(x, y, maxX, maxY)]
							if iter.player != nil {
								if iter.player.plType == ELF {
									kindsOfEnemies |= 1
								} else if iter.player.plType == GOBLIN {
									kindsOfEnemies |= 2
								}
							}
						}
					}
					if kindsOfEnemies != 3 {
						fmt.Println("Not enough enemy types left on the field!")
						return tick - 1, (tick - 1) * sumState(maxX, maxY, track)
					}
					attackPhase(x, y, maxX, maxY, track)
				}
			}
		}
		printState(maxY, maxX, track)
	}
	return 0, 0
}

func sumState(maxX, maxY int, track []*cell) int {
	sum := 0
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			iter := track[linear(x, y, maxX, maxY)]
			if iter.player != nil {
				fmt.Printf("Player on %x, %v has this many points: %v\n", x, y, iter.player.hitPoints)
				sum += iter.player.hitPoints
			}
		}
	}
	fmt.Printf("Sum of all hit points is: %v\n", sum)
	return sum
}

func hasEnemyAround(x int, y int, maxX int, maxY int, track []*cell) bool {
	source := track[linear(x, y, maxX, maxY)]
	if _, ok := isEnemy(x, y-1, maxX, maxY, track, source); ok {
		return true
	}
	if _, ok := isEnemy(x-1, y, maxX, maxY, track, source); ok {
		return true
	}
	if _, ok := isEnemy(x+1, y, maxX, maxY, track, source); ok {
		return true
	}
	if _, ok := isEnemy(x, y+1, maxX, maxY, track, source); ok {
		return true
	}
	return false
}

func attackPhase(x int, y int, maxX int, maxY int, track []*cell) bool {
	var targets = make([]*cell, 0)
	source := track[linear(x, y, maxX, maxY)]
	if source.player == nil {
		return false
	}
	if neighbor, ok := isEnemy(x, y-1, maxX, maxY, track, source); ok {
		targets = append(targets, neighbor)
	}
	if neighbor, ok := isEnemy(x-1, y, maxX, maxY, track, source); ok {
		targets = append(targets, neighbor)
	}
	if neighbor, ok := isEnemy(x+1, y, maxX, maxY, track, source); ok {
		targets = append(targets, neighbor)
	}
	if neighbor, ok := isEnemy(x, y+1, maxX, maxY, track, source); ok {
		targets = append(targets, neighbor)
	}
	if len(targets) == 0 {
		return false
	}
	minTargetPoints := -1
	var chosenTarget *cell
	for _, target := range targets {
		if target.player.hitPoints < minTargetPoints || minTargetPoints == -1 {
			minTargetPoints = target.player.hitPoints
			chosenTarget = target
		}
	}
	chosenTarget.player.hitPoints -= source.player.attackPower
	if chosenTarget.player.hitPoints <= 0 {
		fmt.Printf("player died on %v,%v\n", chosenTarget.x, chosenTarget.y)
		chosenTarget.player = nil
		return true
	}
	return false
}

func isEnemy(x int, y int, maxX int, maxY int, track []*cell, source *cell) (*cell, bool) {
	if neighborId, ok := maybeGet(x, y, maxX, maxY); !ok {
		return nil, false
	} else {
		neighbor := track[neighborId]
		if neighbor.player == nil {
			return nil, false
		}
		return neighbor, neighbor.player.plType != source.player.plType
	}
}

func deduceNextCell(oldX int, oldY int, maxX int, maxY int, track []*cell) (newX, newY int, found bool) {
	if hasEnemyAround(oldX, oldY, maxX, maxY, track) {
		return -1, -1, false
	}
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
			maybeReport(x, y-1, maxX, maxY, track, graph, id, oldId, 1)
			maybeReport(x-1, y, maxX, maxY, track, graph, id, oldId, 2)
			maybeReport(x+1, y, maxX, maxY, track, graph, id, oldId, 3)
			maybeReport(x, y+1, maxX, maxY, track, graph, id, oldId, 4)
		}
	}
	var bestPath *dijkstra.BestPath = nil
	for _, idOfEnemy := range idsOfEnemies {
		//fmt.Printf("shortest between %v->%v\n", oldId, idOfEnemy)
		path, err := graph.Shortest(oldId, idOfEnemy)
		if path.Path == nil {
			continue
		}
		if err != nil {
			continue
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
		//fmt.Printf("No path found from %d:%d\n", oldX, oldY)
		return -1, -1, false
	}
	if len(bestPath.Path) == 2 {
		//fmt.Println("Ready to attack!")
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
					//fmt.Printf("Source(%v)/Destination(%v) not found\n", id, neighborId)
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

func deduceCell(x, y int, cellContents byte) (result *cell) {
	result = &cell{
		x: x,
		y: y,
	}
	c := cellType(cellContents)
	switch c {
	case GOBLIN:
		result.underlying = EMPTY
		result.player = &player{
			plType:      GOBLIN,
			attackPower: 3,
			hitPoints:   200,
		}
	case ELF:
		result.underlying = EMPTY
		result.player = &player{
			plType:      ELF,
			attackPower: 3,
			hitPoints:   200,
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

func printState(maxY int, maxX int, track []*cell) {
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
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			iter := track[linear(x, y, maxX, maxY)]
			if iter.player != nil {
				fmt.Printf("(%v on %x, %v)=%v ", iter.player.plType, x, y, iter.player.hitPoints)
			}
		}
	}
	fmt.Printf("\n")
}
