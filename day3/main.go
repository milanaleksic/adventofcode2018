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

type Field struct {
	x int
	y int
}

type Count struct {
	id    int
	count int
}

func main() {
	file, err := os.Open("day3/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	list := list.New()
	readAll(file, list)

	// #1 @ 1,3: 4x4
	regex, err := regexp.Compile("#(\\d+) @ (\\d+),(\\d+): (\\d+)x(\\d+)")
	if err != nil {
		log.Fatal(err)
	}

	fabric := make(map[Field]Count)
	knownClaims := make(map[int]bool)

	// part 1
	for line := list.Front(); line != nil; line = line.Next() {
		matches := regex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		x, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(matches[3])
		if err != nil {
			log.Fatal(err)
		}
		w, err := strconv.Atoi(matches[4])
		if err != nil {
			log.Fatal(err)
		}
		h, err := strconv.Atoi(matches[5])
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Found id=%d x=%d, y=%d, w=%d, h=%d\n", id, x, y, w, h)
		conflict := false
		for i := x; i < x+w; i++ {
			for j := y; j < y+h; j++ {
				f := Field{i, j}
				existing, ok := fabric[f]
				if !ok {
					fabric[f] = Count{id: id, count: 1}
				} else {
					delete(knownClaims, existing.id)
					conflict = true
					fabric[f] = Count{id: -1, count: existing.count + 1}
				}
			}
		}
		if !conflict {
			knownClaims[id] = true
		}
	}
	count := 0
	for _, c := range fabric {
		if c.count >= 2 {
			count++
		}
	}
	fmt.Printf("Count of non-unique is: %d\n", count)
	fmt.Printf("Remaining non-conflicted claims: %d\n", knownClaims)
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
