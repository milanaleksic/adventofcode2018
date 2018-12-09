package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("day1/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// part 1
	iter := 0
	for e := ls.Front(); e != nil; e = e.Next() {
		iter += e.Value.(int)
	}
	fmt.Println("Result: (part 1)", iter)

	// part 2
	cache := make(map[int]bool)
	iter2 := 0
	for {
		for e := ls.Front(); e != nil; e = e.Next() {
			iter2 += e.Value.(int)
			_, ok := cache[iter2]
			if ok {
				fmt.Println("Another time I see frequency: ", iter2)
				return
			} else {
				cache[iter2] = true
				//fmt.Println(iter2)
			}
		}
	}
}

func readAll(file *os.File, list *list.List) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		list.PushBack(val)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
