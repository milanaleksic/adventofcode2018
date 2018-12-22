package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/RyanCarrier/dijkstra"
)

type cell struct {
	underlying   cellType
	erosionLevel int
}

type cellType byte

func (c cellType) String() string {
	return fmt.Sprintf("%s", string(c))
}

type toolType byte

const (
	Wet    cellType = '='
	Rocky  cellType = '.'
	Narrow cellType = '|'

	Torch toolType = iota
	Climb
	Neither

	input = "day22/input.txt"
	//input = "day22/test.txt"
)

func main() {
	file, err := os.Open(input)
	check(err)
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	line1 := ls.Front()
	depth, err := strconv.Atoi(strings.Split(line1.Value.(string), " ")[1])
	check(err)
	line2coords := strings.Split(line1.Next().Value.(string), " ")
	coords := strings.Split(line2coords[1], ",")
	targetX, err := strconv.Atoi(coords[0])
	check(err)
	targetY, err := strconv.Atoi(coords[1])
	check(err)
	fmt.Printf("Depth is %v, target is @ %v,%v\n", depth, targetX, targetY)

	maxX := targetX + 50
	maxY := targetY + 50
	field := make([]*cell, maxX*maxY)
	calculateField(field, depth, targetX, targetY, maxX, maxY)
	//printState(field, targetX, targetY, maxX, maxY)

	//fmt.Printf("Solution to part 1 (risk level) is: %v\n", riskLevel(0, 0, targetX, targetY, field, maxX, maxY))
	createGraph(field, targetX, targetY, maxX, maxY)

	//solution := part1(field)
	//fmt.Printf("Solution for part 1 is: %v\n", solution)
	//fmt.Printf("Solution for part 2 is: %v\n", part2())
}

func createGraph(cells []*cell, targetX, targetY, maxX int, maxY int) {
	graph := dijkstra.NewGraph()
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			types := allowedTypes[cells[linear(x, y, maxX, maxY)].underlying]
			for _, t := range types {
				graph.AddVertex(linear3d(x, y, t, maxX, maxY))
			}
			t0 := linear3d(x, y, types[0], maxX, maxY)
			t1 := linear3d(x, y, types[1], maxX, maxY)
			check(graph.AddArc(t0, t1, 7))
			check(graph.AddArc(t1, t0, 7))
		}
	}
	for y := 0; y < maxY-1; y++ {
		for x := 0; x < maxX-1; x++ {
			for _, t := range allowedTypes[cells[linear(x, y, maxX, maxY)].underlying] {
				vertexCur := linear3d(x, y, t, maxX, maxY)
				for _, tr := range allowedTypes[cells[linear(x+1, y, maxX, maxY)].underlying] {
					vertexRight := linear3d(x+1, y, tr, maxX, maxY)
					if t == tr {
						check(graph.AddArc(vertexCur, vertexRight, 1))
						check(graph.AddArc(vertexRight, vertexCur, 1))
					}
				}
				for _, tb := range allowedTypes[cells[linear(x, y+1, maxX, maxY)].underlying] {
					vertexDown := linear3d(x, y+1, tb, maxX, maxY)
					if t == tb {
						check(graph.AddArc(vertexCur, vertexDown, 1))
						check(graph.AddArc(vertexDown, vertexCur, 1))
					}
				}
			}
		}
	}
	path, err := graph.Shortest(linear3d(0, 0, Torch, maxX, maxY), linear3d(targetX, targetY, Torch, maxX, maxY))
	check(err)
	fmt.Printf("Path price: %v\n", path.Distance)
	//last := -1
	//for _, p := range path.Path {
	//	x, y := breakLinear3D(p, maxX, maxY)
	//	var price int64 = 0
	//	var ok bool
	//	if last != -1 {
	//		price, ok = graph.Verticies[last].GetArc(p)
	//		if !ok {
	//			log.Fatalf("Programming error, expected arc to exist")
	//		}
	//	}
	//	fmt.Printf("Path: %v, %v (price=%v)\n", x, y, price)
	//	last = p
	//}
}

var allowedTypes = map[cellType][2]toolType{
	Wet:    {Climb, Neither},
	Rocky:  {Climb, Torch},
	Narrow: {Torch, Neither},
}

func breakLinear3D(id int, maxX, maxY int) (x int, y int) {
	return id % maxX, id % (maxX * maxY) / maxX
}

func riskLevel(x1 int, y1 int, x2 int, y2 int, cells []*cell, maxX, maxY int) int {
	totalRisk := 0
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			switch cells[linear(x, y, maxX, maxY)].underlying {
			case Wet:
				totalRisk += 1
			case Narrow:
				totalRisk += 2
			}
		}
	}
	return totalRisk
}

func calculateField(cells []*cell, depth, targetX, targetY int, maxX int, maxY int) {
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			var geologicIndex int
			if (x == 0 && y == 0) || (x == targetX && y == targetY) {
				geologicIndex = 0
			} else if y == 0 {
				geologicIndex = 16807 * x
			} else if x == 0 {
				geologicIndex = 48271 * y
			} else {
				erosionLevel1 := cells[linear(x-1, y, maxX, maxY)].erosionLevel
				erosionLevel2 := cells[linear(x, y-1, maxX, maxY)].erosionLevel
				geologicIndex = erosionLevel1 * erosionLevel2
			}
			erosionLevel := (geologicIndex + depth) % 20183
			switch erosionLevel % 3 {
			case 0:
				cells[linear(x, y, maxX, maxY)] = &cell{underlying: Rocky, erosionLevel: erosionLevel}
			case 1:
				cells[linear(x, y, maxX, maxY)] = &cell{underlying: Wet, erosionLevel: erosionLevel}
			case 2:
				cells[linear(x, y, maxX, maxY)] = &cell{underlying: Narrow, erosionLevel: erosionLevel}
			}
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
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

func linear(x int, y int, maxX, maxY int) int {
	scaled := x + y*maxX
	if x >= maxX {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	if y >= maxY {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	return scaled
}

func linear3d(x int, y int, tool toolType, maxX, maxY int) int {
	scaled := x + y*maxX + maxX*maxY*int(tool)
	if x >= maxX {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	if y >= maxY {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	return scaled
}

func printState(track []*cell, targetX, targetY int, maxX, maxY int) {
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			if x == 0 && y == 0 {
				fmt.Print("M")
			} else if x == targetX && y == targetY {
				fmt.Print("T")
			} else {
				c := track[linear(x, y, maxX, maxY)]
				fmt.Printf("%v", c.underlying)
			}
		}
		fmt.Println()
	}
	fmt.Printf("\n")
}
