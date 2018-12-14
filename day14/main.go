package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// test
//const input1 = "day14/test1.txt"
//const input2 = "day14/test2.txt"

//production
const input1 = "day14/input1.txt"
const input2 = "day14/input2.txt"

func main() {
	file, err := os.Open(input1)
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
	fmt.Printf("Solution 1 is: %s\n", part1(recipeCount))

	file2, err := os.Open(input2)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	ls2 := list.New()
	readAll(file2, ls2)
	expectedScores := ls2.Front().Value.(string)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Input is recipeCount=%s\n", expectedScores)
	fmt.Printf("Solution 2 is: %d\n", part2(expectedScores))
}

type linkedList struct {
	head, tail *cell
	length     int
}

func (l *linkedList) appendCell(value int) *cell {
	newTail := &cell{
		value: value,
		owner: l,
	}
	newTail.prev = l.tail
	if l.tail != nil {
		l.tail.next = newTail
	}
	l.tail = newTail
	if l.head == nil {
		l.head = newTail
	}
	newTail.index = l.length
	l.length++
	return newTail
}

func (l *linkedList) printState(player1, player2 *cell) {
	iter := l.head
	fmt.Printf("State: ")
	for iter != nil {
		if iter == player1 {
			fmt.Print("(")
		} else if iter == player2 {
			fmt.Print("[")
		}
		fmt.Printf("%d", iter.value)
		if iter == player1 {
			fmt.Print(")")
		} else if iter == player2 {
			fmt.Print("]")
		} else {
			fmt.Print(" ")
		}
		iter = iter.next
	}
	fmt.Printf("\n")
}

type cell struct {
	next  *cell
	prev  *cell
	value int
	index int
	owner *linkedList
}

func (c *cell) String() string {
	return fmt.Sprintf("value=%d,index=%d", c.value, c.index)
}

func (cell *cell) move(spaces int) *cell {
	iter := cell
	for i := 1; i <= spaces; i++ {
		iter = iter.next
		if iter == nil {
			iter = cell.owner.head
		}
	}
	return iter
}

func part1(recipeCount int) string {
	l := &linkedList{}
	l.appendCell(3)
	l.appendCell(7)
	iter := 0
	player1 := l.head
	player2 := l.tail
	for {
		//l.printState(player1, player2)
		iter++

		sum := player1.value + player2.value
		//fmt.Printf("Sum=%d\n", sum)
		sumAsString := strconv.Itoa(sum)
		for _, c := range sumAsString {
			i, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatal(err)
			}
			l.appendCell(i)
		}
		player1 = player1.move(player1.value + 1)
		player2 = player2.move(player2.value + 1)
		if l.length >= recipeCount+10 {
			iter := l.head
			for i := 0; i < recipeCount; i++ {
				//fmt.Printf("Skipping %d\n", iter.value)
				iter = iter.next
			}
			solution := make([]string, 0)
			for i := 0; i < 10; i++ {
				solution = append(solution, strconv.Itoa(iter.value))
				//fmt.Printf("Including %d\n", iter.value)
				iter = iter.next
			}
			return strings.Join(solution, "")
		}
	}
}

func part2(scores string) int {
	l := &linkedList{}
	l.appendCell(3)
	l.appendCell(7)
	//iter := 0
	//maxIter := 12
	player1 := l.head
	player2 := l.tail
	//for iter < maxIter {
	for {
		//l.printState(player1, player2)
		//iter++

		sum := player1.value + player2.value
		//fmt.Printf("Sum=%d\n", sum)
		sumAsString := strconv.Itoa(sum)
		for _, c := range sumAsString {
			i, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatal(err)
			}
			l.appendCell(i)
		}
		player1 = player1.move(player1.value + 1)
		player2 = player2.move(player2.value + 1)

		if l.length < len(scores)*2+1 {
			continue
		}
		for i := l.length - len(scores)*2; i < l.length-len(scores); i++ {
			c := l.tail
			for j := l.length - 1; j > i; j-- {
				c = c.prev
			}
			for j := 0; j < len(scores); j++ {
				if strconv.Itoa(c.value)[0] != scores[j] {
					break
				}
				//if j > 3 {
				//	fmt.Printf("Matched %d\n", j+1)
				//}
				if j == len(scores)-1 {
					return i
				}
				c = c.next
			}
		}
		//if l.length%1000 == 0 {
		//	fmt.Println("Reached size: ", l.length)
		//}
	}
	return 0
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
