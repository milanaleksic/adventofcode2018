package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

type typeNode int

type coord struct {
	x, y int
}

type node struct {
	value    string
	children []*node
	operator typeNode
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	file, err := os.Open("day20/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)
	inputLine := ls.Front()
	part1, part2 := solve(inputLine.Value.(string))
	fmt.Printf("Solution to part 1 is: %v\n", part1)
	fmt.Printf("Solution to part 2 is: %v\n", part2)
}

func solve(input string) (part1, part2 int) {
	locations := list.New()
	newCoord := coord{0, 0}
	iterCoord := coord{0, 0}
	dists := make(map[coord]int)
	for iterInput := 0; iterInput < len(input); iterInput++ {
		direction := string(input[iterInput])
		switch direction {
		case "N":
			newCoord = coord{iterCoord.x, iterCoord.y + 1}
			if _, ok := dists[newCoord]; !ok || dists[newCoord] > dists[iterCoord]+1 {
				dists[newCoord] = dists[iterCoord] + 1
			}
		case "W":
			newCoord = coord{iterCoord.x - 1, iterCoord.y}
			if _, ok := dists[newCoord]; !ok || dists[newCoord] > dists[iterCoord]+1 {
				dists[newCoord] = dists[iterCoord] + 1
			}
		case "E":
			newCoord = coord{iterCoord.x + 1, iterCoord.y}
			if _, ok := dists[newCoord]; !ok || dists[newCoord] > dists[iterCoord]+1 {
				dists[newCoord] = dists[iterCoord] + 1
			}
		case "S":
			newCoord = coord{iterCoord.x, iterCoord.y - 1}
			if _, ok := dists[newCoord]; !ok || dists[newCoord] > dists[iterCoord]+1 {
				dists[newCoord] = dists[iterCoord] + 1
			}
		case "(":
			locations.PushFront(iterCoord)
		case "|":
			front := locations.Front()
			newCoord = front.Value.(coord)
		case ")":
			front := locations.Front()
			locations.Remove(front)
			newCoord = front.Value.(coord)
		}
		iterCoord = newCoord
	}
	fmt.Printf("%v\n", dists)
	maxD := 0
	for _, d := range dists {
		if d > maxD {
			maxD = d
		}
	}
	countMoreThan1K := 0
	for _, d := range dists {
		if d >= 1000 {
			countMoreThan1K++
		}
	}
	return maxD, countMoreThan1K
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
