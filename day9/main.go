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

// test
const input = "day9/test.txt"

//production
//const input = "day9/input.txt"

func main() {
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	regex, err := regexp.Compile("(\\d+) players; last marble is worth (\\d+) points")
	if err != nil {
		log.Fatal(err)
	}

	ls := list.New()
	readAll(file, ls)
	for line := ls.Front(); line != nil; line = line.Next() {
		matches := regex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		players, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		lastMarble, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found players=%d, lastMarble=%d\n", players, lastMarble)
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
