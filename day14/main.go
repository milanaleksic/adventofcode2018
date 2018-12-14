package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
)

// test
const input = "day14/test.txt"

//production
//const input = "day14/input.txt"

func main() {
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)
	recipeCountString := ls.Front().Value.(string)
	recipeCount, err := strconv.Atoi(recipeCountString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Input is recipeCount=%d\n", recipeCount)
	fmt.Printf("Solution is: %s", part1(recipeCount))
}

type linkedList struct {
	head, tail *cell
	length     int
}

func (l *linkedList) appendCell(value int) *cell {
	newTail := &cell{
		value: value,
	}
	newTail.prev = l.tail
	if l.tail != nil {
		l.tail.next = newTail
	}
	l.tail = newTail
	if l.head == nil {
		l.head = newTail
	}
	l.length++
	return newTail
}

func (l *linkedList) printState() {
	iter := l.head
	fmt.Printf("State: ")
	for iter != nil {
		fmt.Printf("%d ", iter.value)
		iter = iter.next
	}
	fmt.Printf("\n")
}

type cell struct {
	next  *cell
	prev  *cell
	value int
}

func part1(recipeCount int) string {
	linkedList := &linkedList{}
	linkedList.appendCell(3)
	linkedList.appendCell(7)
	maxIter := 10
	iter := 0
	for iter <= maxIter {
		linkedList.printState()
		iter++
	}
	return "0124515891"
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
