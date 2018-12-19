package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

type cell struct {
	underlying cellType
}

type cellType byte

func (c cellType) String() string {
	return fmt.Sprintf("%s", string(c))
}

const (
	LUMBERYARD cellType = '#'
	OPEN       cellType = '.'
	TREES      cellType = '|'

	input = "day18/input.txt"
	maxX  = 50
	maxY  = 50

	//input = "day18/test.txt"
	//maxX  = 10
	//maxY  = 10
)

func main() {
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	field := make([]*cell, maxX*maxY)
	y := 0
	for line := ls.Front(); line != nil; line = line.Next() {
		fieldLine := line.Value.(string)
		for x, c := range fieldLine {
			c := deduceCell(x, y, byte(c))
			field[linear(x, y)] = c
		}
		y++
	}
	printState(field)
	//solution := part1(field)
	//fmt.Printf("Solution for part 1 is: %v\n", solution)
	fmt.Printf("Solution for part 2 is: %v\n", part2())
}

func part2() int {
	j := -1
	detectedValues := []int{201335, 190896, 195952, 198744, 199134, 208208, 203750}
	for i := 30000; i <= 1000000000; i += 10000 {
		j++
		fmt.Printf("For value %v: %v\n", i, detectedValues[j%len(detectedValues)])
	}
	return detectedValues[j%len(detectedValues)]
}

func part1(cells []*cell) int {
	convergedNumber := 0
	maxIter := 1000000000
	for iter := 1; iter <= maxIter; iter++ {
		if iter%10000 == 0 {
			fmt.Printf("******ITER %v", iter)
		}
		newCells := copyField(cells)
		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				cell := cells[linear(x, y)]
				trees, lumberyards := counts(x, y, cells)
				if cell.underlying == OPEN && trees >= 3 {
					newCells[linear(x, y)].underlying = TREES
				} else if cell.underlying == TREES && lumberyards >= 3 {
					newCells[linear(x, y)].underlying = LUMBERYARD
				} else if cell.underlying == LUMBERYARD {
					if lumberyards > 0 && trees > 0 {
						newCells[linear(x, y)].underlying = LUMBERYARD
					} else {
						newCells[linear(x, y)].underlying = OPEN
					}
				}
			}
		}
		cells = newCells
		//printState(cells)
		countLumberyards := 0
		countTrees := 0
		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				cell := cells[linear(x, y)]
				if cell.underlying == TREES {
					countTrees++
				} else if cell.underlying == LUMBERYARD {
					countLumberyards++
				}
			}
		}
		if countLumberyards*countTrees == convergedNumber {
			return countTrees * countLumberyards
		} else {
			convergedNumber = countTrees * countLumberyards
		}
		if iter%10000 == 0 {
			fmt.Printf(", converged to: %v\n", convergedNumber)
		}
	}
	return -1
}

func counts(x int, y int, cells []*cell) (trees int, lumberyards int) {
	if y > 0 {
		count(cells, x, y-1, &lumberyards, &trees)
	}
	if y < maxY-1 {
		count(cells, x, y+1, &lumberyards, &trees)
	}
	if x > 0 {
		if y > 0 {
			count(cells, x-1, y-1, &lumberyards, &trees)
		}
		count(cells, x-1, y, &lumberyards, &trees)
		if y < maxY-1 {
			count(cells, x-1, y+1, &lumberyards, &trees)
		}
	}
	if x < maxX-1 {
		if y > 0 {
			count(cells, x+1, y-1, &lumberyards, &trees)
		}
		count(cells, x+1, y, &lumberyards, &trees)
		if y < maxY-1 {
			count(cells, x+1, y+1, &lumberyards, &trees)
		}
	}
	return trees, lumberyards
}

func count(cells []*cell, x int, y int, lumberyards *int, trees *int) {
	if cells[linear(x, y)].underlying == LUMBERYARD {
		*lumberyards++
	}
	if cells[linear(x, y)].underlying == TREES {
		*trees++
	}
}

func copyField(cells []*cell) []*cell {
	result := make([]*cell, len(cells))
	for i, c := range cells {
		result[i] = &cell{
			underlying: c.underlying,
		}
	}
	return result
}

func deduceCell(x, y int, cellContents byte) (result *cell) {
	result = &cell{
		underlying: cellType(cellContents),
	}
	return
}

func linear(x int, y int) int {
	scaled := x + y*maxX
	if x >= maxX {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	if y >= maxY {
		log.Fatalf("Tried to access %v,%v", x, y)
	}
	return scaled
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

func printState(track []*cell) {
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			c := track[linear(x, y)]
			fmt.Printf("%v", c.underlying)
		}
		fmt.Println()
	}
	fmt.Printf("\n")
}
