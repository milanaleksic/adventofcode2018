package main

import (
	"bufio"
	"container/list"
	"log"
	"os"
)

func main() {
	file, err := os.Open("day2/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	count2 := 0
	count3 := 0

	// part 1
	for line := ls.Front(); line != nil; line = line.Next() {
		counts := make(map[rune]int)
		for _, c := range line.Value.(string) {
			count, ok := counts[c]
			if !ok {
				counts[c] = 1
			} else {
				counts[c] = 1 + count
			}
		}
		localCount2 := false
		localCount3 := false
		for _, c := range counts {
			if c == 2 {
				localCount2 = true
			} else if c == 3 {
				localCount3 = true
			}
		}
		if localCount2 {
			count2 += 1
		}
		if localCount3 {
			count3 += 1
		}
	}
	println(count2 * count3)

	// part 2
	for line1 := ls.Front(); line1 != nil; line1 = line1.Next() {
		for line2 := ls.Front(); line2 != nil; line2 = line2.Next() {
			if line1 == line2 {
				continue
			}
			wrong := -1
			i := 0
			line1 := line1.Value.(string)
			line2 := line2.Value.(string)
			for i = range line1 {
				c1 := line1[i]
				c2 := line2[i]
				if c1 != c2 {
					if wrong == -1 {
						wrong = i
					} else {
						wrong = -1
						break
					}
				}
			}
			if i == len(line1)-1 && wrong != -1 {
				println("Found!", line1[0:wrong]+line1[wrong+1:])
				//return
			}
		}
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
