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

type cell struct {
	underlying  cellType
	fallThrough bool
}

type cellType byte

func (c *cellType) String() string {
	return fmt.Sprintf("%s", string(*c))
}

const (
	CLAY          cellType = '#'
	SAND          cellType = '.'
	SPRING        cellType = '+'
	WATER_FALLING cellType = '|'
	WATER_STABLE  cellType = '~'
)

func main() {
	//file, err := os.Open("day17/input2.txt")
	file, err := os.Open("day17/input.txt")
	//file, err := os.Open("day17/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)
	minX, maxX, minY, maxY, track := fromInput(ls)
	fmt.Printf("minX=%v, maxX=%v, minY=%v, maxY=%v\n", minX, maxX, minY, maxY)
	count, countStable := part1And2(minX, maxX, minY, maxY, track)
	fmt.Printf("Solution to part 1 is: %v\n", count)
	fmt.Printf("Solution to part 2 is: %v\n", countStable)
}

func part1And2(minX int, maxX int, minY, maxY int, ground []*cell) (count, countStable int) {
	maxIter := 35000
	//maxIter := 38489
	ground[linear(500, 0, maxX, maxY)].underlying = SPRING
	iter := 1
	defer func() {
		state := recover()
		fmt.Printf("failed with state %v, in iteration %v\n", state, iter)
		printState(minX, maxY, maxX, ground)
	}()
	for iter = 1; iter <= maxIter; iter++ {
		//fmt.Printf("########### ITER:%v\n", iter)
		dropX, dropY := 500, 1
		newDropX, newDropY, _ := findNextRestingPlace(dropX, dropY, minX, maxX, minY, maxY, ground, 0)
		if newDropX == -1 && newDropY == -1 || (newDropX == dropX && newDropY == dropY) {
			ground[linear(500, 1, maxX, maxY)].underlying = WATER_FALLING
			printState(minX, maxY, maxX, ground)
			return countDrops(minX, maxX, minY, maxY, ground)
		}
		//printState(minX, maxY, maxX, ground)
	}
	return countDrops(minX, maxX, minY, maxY, ground)
}

func countDrops(minX int, maxX int, minY, maxY int, cells []*cell) (count, countStable int) {
	count = 0
	countStable = 0
	for x := minX; x <= maxX; x++ {
		for y := 0; y <= maxY; y++ {
			cell := cells[linear(x, y, maxX, maxY)]
			if cell.underlying == WATER_FALLING {
				if y >= minY {
					count++
				}
			}
			if cell.underlying == WATER_STABLE {
				if y >= minY {
					count++
					countStable++
				}
			}
		}
	}
	return count, countStable
}

func findNextRestingPlace(x, y int, minX int, maxX int, minY int, maxY int, ground []*cell, direction int) (newX, newY int, fallth bool) {
	//fmt.Printf("visiting x=%v, y=%v\n", x, y)
	cur := linear(x, y, maxX, maxY)
	if ground[cur].underlying != SAND {
		//fmt.Println("returning false 2")
		return -1, -1, ground[cur].fallThrough
	}
	if x < minX {
		ground[cur].fallThrough = true
		return minX, y, true
	}
	if x >= maxX {
		ground[cur].fallThrough = true
		return maxX, y, true
	}
	if y >= maxY {
		ground[cur].fallThrough = true
		return x, maxY, true
	}
	below := linear(x, y+1, maxX, maxY)
	if ground[below].underlying == SAND {
		return findNextRestingPlace(x, y+1, minX, maxX, minY, maxY, ground, 0)
	}
	if ground[below].underlying == CLAY || ground[below].underlying == WATER_STABLE || ground[below].underlying == WATER_FALLING {
		if ground[below].fallThrough {
			ground[linear(x, y, maxX, maxY)].fallThrough = true
			return x, y, true
		}
		if direction == 0 {
			leftLimitFound, rightLimitFound := false, false
			if newX, newY, fallth = findNextRestingPlace(x-1, y, minX, maxX, minY, maxY, ground, -1); newX != -1 && newY != -1 {
				target := linear(newX, newY, maxX, maxY)
				if ground[target].underlying == SAND {
					//fmt.Printf("Setting to falling 4: %v, %v\n", newX, newY)
					ground[target].underlying = WATER_FALLING
				}
				if ground[target].fallThrough {
					ground[linear(x, y, maxX, maxY)].fallThrough = true
				}
				//if !ground[target].fallThrough {
				//	fmt.Println("returning false 3")
				//}
				return newX, newY, ground[target].fallThrough
			} else {
				leftLimitFound = !fallth
			}
			if newX, newY, fallth = findNextRestingPlace(x+1, y, minX, maxX, minY, maxY, ground, +1); newX != -1 && newY != -1 {
				target := linear(newX, newY, maxX, maxY)
				if ground[target].underlying == SAND {
					//fmt.Printf("Setting to falling 3: %v, %v\n", newX, newY)
					ground[target].underlying = WATER_FALLING
				}
				if ground[target].fallThrough {
					ground[linear(x, y, maxX, maxY)].fallThrough = true
				}
				//if !ground[target].fallThrough {
				//	fmt.Println("returning false 4")
				//}
				return newX, newY, ground[target].fallThrough
			} else {
				rightLimitFound = !fallth
			}
			if leftLimitFound || rightLimitFound {
				cascadedDeduceForLine(x, y, minX, maxX, maxY, ground)
				//if !fallth {
				//	fmt.Println("returning false 5")
				//}
				return x, y, !fallth
			} else {
				//fmt.Printf("Setting to falling 1: %v, %v\n", x, y)
				ground[linear(x, y, maxX, maxY)].underlying = WATER_FALLING
			}
		} else {
			if newX, newY, fallth = findNextRestingPlace(x+direction, y, minX, maxX, minY, maxY, ground, direction); newX != -1 && newY != -1 {
				return newX, newY, fallth
			} else {
				//if fallth {
				//	return -1, -1, false
				//}
				//if ground[linear(x+1, y+1, maxX, maxY)].fallThrough && ground[linear(x-1, y+1, maxX, maxY)].underlying == WATER_STABLE {
				//	return -1, -1
				//}
				//fmt.Printf("Setting to falling 2: %v, %v (fallth=%v)\n", x, y, fallth)
				ground[linear(x, y, maxX, maxY)].underlying = WATER_FALLING
				if ground[linear(x+direction, y, maxX, maxY)].fallThrough {
					ground[linear(x, y, maxX, maxY)].fallThrough = true
					return x, y, true
				}
			}
		}
	}
	return x, y, false
}

func cascadedDeduceForLine(x int, y int, minX int, maxX int, maxY int, ground []*cell) {
	oneSideIsFallthrough := false
	for ix := x; ix >= minX; ix-- {
		iter := ground[linear(ix, y, maxX, maxY)]
		if iter.fallThrough {
			//fmt.Printf("Setting to fall through: %v, %v\n", ix, y)
			iter.underlying = WATER_FALLING
			iter.fallThrough = true
			oneSideIsFallthrough = true
		} else if iter.underlying == WATER_FALLING || ix == x {
			//fmt.Printf("Setting to stable: %v, %v\n", ix, y)
			iter.underlying = WATER_STABLE
		} else if iter.underlying == SAND || iter.underlying == CLAY {
			break
		}
	}
	for ix := x; ix <= maxX; ix++ {
		iter := ground[linear(ix, y, maxX, maxY)]
		if iter.fallThrough {
			//fmt.Printf("Setting to fall through: %v, %v\n", ix, y)
			iter.underlying = WATER_FALLING
			iter.fallThrough = true
			oneSideIsFallthrough = true
		} else if iter.underlying == WATER_FALLING || ix == x {
			//fmt.Printf("Setting to stable: %v, %v\n", ix, y)
			iter.underlying = WATER_STABLE
		} else if iter.underlying == SAND || iter.underlying == CLAY {
			break
		}
	}
	if oneSideIsFallthrough {
		for ix := x; ix >= minX; ix-- {
			iter := ground[linear(ix, y, maxX, maxY)]
			if iter.underlying == WATER_STABLE {
				iter.underlying = WATER_FALLING
			} else if iter.underlying == SAND || iter.underlying == CLAY {
				break
			}
			iter.fallThrough = true
		}
		for ix := x; ix <= maxX; ix++ {
			iter := ground[linear(ix, y, maxX, maxY)]
			if iter.underlying == WATER_STABLE {
				iter.underlying = WATER_FALLING
			} else if iter.underlying == SAND || iter.underlying == CLAY {
				break
			}
			iter.fallThrough = true
		}
	}
}

func fromInput(ls *list.List) (minX, maxX, minY, maxY int, track []*cell) {
	// x=301, y=218..246
	regexLine := regexp.MustCompile(`([xy])=(\d+)\.?\.?(\d+)?`)
	type metDatum struct {
		xStart, xEnd, yStart, yEnd int
	}
	metData := make([]*metDatum, 0)
	minX = 100000
	minY = 100000
	for line := ls.Front(); line != nil; line = line.Next() {
		trackLine := line.Value.(string)
		matches := regexLine.FindAllStringSubmatch(trackLine, -1)
		if len(matches) == 0 {
			log.Fatalf("No matches found: in %v", trackLine)
		}
		d := &metDatum{}
		for _, match := range matches {
			start, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatalf("Could not convert string start %v, match=%v", match[2], match)
			}
			end := start
			if len(match) == 4 && match[3] != "" {
				end, err = strconv.Atoi(match[3])
				if err != nil {
					log.Fatalf("Could not convert string end %v, match=%v", match[3], match)
				}
			}
			switch match[1] {
			case "x":
				d.xStart = start
				d.xEnd = end
				if start < minX {
					minX = start
				}
				if end > maxX {
					maxX = end
				}
			case "y":
				d.yStart = start
				d.yEnd = end
				if start < minY {
					minY = start
				}
				if end > maxY {
					maxY = end
				}
			}
			metData = append(metData, d)
		}
	}
	maxX += 2
	minX -= 2
	track = make([]*cell, (maxX+2)*(maxY+2))
	for _, d := range metData {
		for x := d.xStart; x <= d.xEnd; x++ {
			for y := d.yStart; y <= d.yEnd; y++ {
				track[linear(x, y, maxX, maxY)] = &cell{underlying: CLAY}
				//fmt.Printf("%v,%v is CLAY\n", x, y)
			}
		}
	}
	for x := 0; x <= maxX; x++ {
		for y := 0; y <= maxY; y++ {
			if track[linear(x, y, maxX, maxY)] == nil {
				track[linear(x, y, maxX, maxY)] = &cell{underlying: SAND}
			}
		}
	}
	return minX, maxX, minY, maxY, track
}

func linear(x int, y int, maxX int, maxY int) int {
	scaled := x + y*maxX
	if x > maxX+1 {
		panic(fmt.Sprintf("Tried to access %v,%v (but maxX=%v)", x, y, maxX))
	}
	if y > maxY+1 {
		panic(fmt.Sprintf("Tried to access %v,%v (but maxY=%v)", x, y, maxY))
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

func printState(minX int, maxY int, maxX int, track []*cell) {
	for y := 0; y <= maxY+1; y++ {
		//if y < 449 || y >= 470 {
		//	continue
		//}
		for x := minX - 1; x <= maxX+1; x++ {
			if x == maxX+1 || y > maxY {
				fmt.Print(string(SAND))
				//} else if track[linear(x, y, maxX, maxY)].fallThrough {
				//	fmt.Printf("F")
			} else {
				fmt.Printf(string(track[linear(x, y, maxX, maxY)].underlying))
			}
		}
		fmt.Printf(" line=%v", y)
		fmt.Println()
	}
	fmt.Printf("\n")
}
