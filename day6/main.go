package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type coord struct {
	x int
	y int
}

type dist struct {
	minDist   int
	totalDist float64
}

func main() {
	file, err := os.Open("day6/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// #1 @ 1,3: 4x4
	regex, err := regexp.Compile("(\\d+), (\\d+)")
	if err != nil {
		log.Fatal(err)
	}

	// part 1
	maxX := 0
	maxY := 0
	field := make(map[coord]int)
	iter := 0
	for line := ls.Front(); line != nil; line = line.Next() {
		matches := regex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		x, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Found x=%d, y=%d\n", x, y)
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		field[coord{x: x, y: y}] = iter
		iter++
	}
	//fmt.Printf("Size of the branch: %d:%d\n", maxX+1, maxY+1)
	//fmt.Printf("Field: %+v", field)

	fieldMapping := make(map[coord]dist)
	for y := 0; y < maxY+2; y++ {
		for x := 0; x < maxX+2; x++ {
			dist := minDist(field, x, y)
			fieldMapping[coord{x: x, y: y}] = dist
			//fmt.Printf("% 2d", dist)
		}
		//fmt.Print("\n")
	}
	knownWithUnlimited := make(map[int]bool)
	for x := 0; x < maxX+2; x++ {
		knownWithUnlimited[fieldMapping[coord{x: x, y: 0}].minDist] = true
		knownWithUnlimited[fieldMapping[coord{x: x, y: maxY + 1}].minDist] = true
	}
	for y := 0; y < maxY+2; y++ {
		knownWithUnlimited[fieldMapping[coord{x: 0, y: y}].minDist] = true
		knownWithUnlimited[fieldMapping[coord{x: maxX + 1, y: y}].minDist] = true
	}
	mappingSizes := make(map[int]int)
	for y := 0; y < maxY+2; y++ {
		for x := 0; x < maxX+2; x++ {
			id := fieldMapping[coord{x: x, y: y}].minDist
			if !knownWithUnlimited[id] {
				soFar, ok := mappingSizes[id]
				if ok {
					mappingSizes[id] = soFar + 1
				} else {
					mappingSizes[id] = 1
				}
			}
		}
	}
	//fmt.Printf("%v\n", mappingSizes)
	maxSize := -1
	for _, size := range mappingSizes {
		if size > maxSize {
			maxSize = size
		}
	}
	fmt.Printf("Solution 1 is: %v\n", maxSize)
	// part 2
	result := 0
	// testing
	//sizeLimit := float64(32)
	// real input
	sizeLimit := float64(10000)
	for y := 0; y < maxY+2; y++ {
		for x := 0; x < maxX+2; x++ {
			totalDist := fieldMapping[coord{x: x, y: y}].totalDist
			if totalDist < sizeLimit {
				result++
			}
		}
	}
	fmt.Printf("Solution 2 is: %v", result)
}

func minDist(field map[coord]int, x, y int) dist {
	minDist := math.MaxFloat64
	minDistID := -1
	conflict := false
	totalDist := float64(0)
	for coord, id := range field {
		dist := math.Abs(float64(x-coord.x)) + math.Abs(float64(y-coord.y))
		if dist < minDist {
			minDist = dist
			minDistID = id
			conflict = false
		} else if dist == minDist {
			conflict = true
		}
		totalDist += dist
	}
	if conflict {
		return dist{minDist: -1, totalDist: totalDist}
	}
	return dist{minDist: minDistID, totalDist: totalDist}
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
