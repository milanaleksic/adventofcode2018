package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
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
	file, err := os.Open("day5/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)
	for line := ls.Front(); line != nil; line = line.Next() {
		l := []byte(line.Value.(string))
		part1(l)
		part2(l)
	}

}

func part1(ls []byte) int {
	fmt.Println(string(ls))
	l := make([]byte, len(ls))
	copy(l, ls)
	iter := 1
	for {
		if iter >= len(l) {
			break
		}
		if iter == 0 {
			iter++
		} else {
			//fmt.Printf("len=%d, iter=%d\n", len(l), iter)
			curr := float64(l[iter])
			prev := float64(l[iter-1])
			diff := math.Abs(curr - prev)
			if diff == 32 {
				l = append(l[:iter-1], l[iter+1:]...)
				if iter > 0 {
					iter--
				}
				//fmt.Println(string(l))
			} else {
				iter++
			}
		}
	}
	fmt.Printf("After elimination, remainder is: %d\n", len(l))
	return len(l)
}

func part2(ls []byte) {
	l := make([]byte, len(ls))
	copy(l, ls)
	knownTypes := make(map[byte]int)
	for _, b := range l {
		if b < 'A' || b > 'z' {
			log.Fatalf("Wrong byte found: %d", b)
		}
		_, ok := knownTypes[b]
		if !ok {
			knownTypes[b] = 1 // low is there
		} else {
			knownTypes[b] |= 0x01 // low is there
		}
		if b <= 'Z' {
			_, ok = knownTypes[b+32]
			if !ok {
				knownTypes[b] = 2 // high is there
			} else {
				knownTypes[b] |= 0x02 // high is there
			}
		}
	}
	min := -1
	for knownType, occurrence := range knownTypes {
		if occurrence == 0x03 {
			candidate := make([]byte, len(l))
			iter := 0
			for _, b := range l {
				if b != knownType && b != (knownType+32) {
					candidate[iter] = b
					iter++
				}
			}
			localMin := part1(candidate[:iter])
			if min > localMin || min == -1 {
				min = localMin
			}
		}
	}
	fmt.Printf("Global minimum is: %d", min)
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
