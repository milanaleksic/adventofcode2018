package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type point struct {
	x      int
	y      int
	z      int
	radius float64
}

const (
	input = "day23/input.txt"
	//input = "day23/test.txt"
	//input = "day23/test2.txt"
)

func main() {
	file, err := os.Open(input)
	check(err)
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// pos=<0,0,0>, r=4
	regex, err := regexp.Compile(`pos=<(-?\d+),(-?\d+),(-?\d+)>,\s*r=(\d+)`)
	check(err)

	points := make([]*point, 0)
	for line := ls.Front(); line != nil; line = line.Next() {
		matches := regex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		x, err := strconv.Atoi(matches[1])
		check(err)
		y, err := strconv.Atoi(matches[2])
		check(err)
		z, err := strconv.Atoi(matches[3])
		check(err)
		radius, err := strconv.Atoi(matches[4])
		check(err)
		points = append(points, &point{x: x, y: y, z: z, radius: float64(radius)})
	}

	fmt.Printf("Solution to part 1 is: %v\n", part1(points))

	countP1 := 0
	p1 := &point{x: 11382618, y: 29059365, z: 39808800}
	for _, p2 := range points {
		d1 := dist(p1, p2)
		if d1 <= p2.radius {
			countP1++
		}
	}
	fmt.Printf("Number of bots that match: %v, distance: %v", countP1, int(dist(&point{}, p1)))

	//fmt.Printf("Solution to part 2 is: %v\n", int(part2(points)))
}

func part1(points []*point) int {
	var maxRP *point = nil
	for _, p := range points {
		if maxRP == nil || p.radius > maxRP.radius {
			maxRP = p
		}
	}
	//fmt.Printf("Point with max radius is: %v\n", maxRP)
	countInRadius := 0
	for _, p := range points {
		if dist(p, maxRP) <= maxRP.radius {
			//fmt.Printf("Point is in radius: %v\n", p)
			countInRadius++
		}
	}
	return countInRadius
}

func part2(points []*point) float64 {
	var totalMaxRP *point = nil
	var totalHowManyMax = 0
	var zero = &point{}
	var minX, minY, minZ, maxX, maxY, maxZ = math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1)
	for _, p := range points {
		if float64(p.x) < minX || math.IsInf(minX, -1) {
			minX = float64(p.x)
		}
		if float64(p.y) < minY || math.IsInf(minY, -1) {
			minY = float64(p.y)
		}
		if float64(p.z) < minZ || math.IsInf(minZ, -1) {
			minZ = float64(p.z)
		}
		if float64(p.x) > maxX || math.IsInf(maxX, -1) {
			maxX = float64(p.x)
		}
		if float64(p.y) > maxY || math.IsInf(maxY, -1) {
			maxY = float64(p.y)
		}
		if float64(p.z) > maxZ || math.IsInf(maxZ, -1) {
			maxZ = float64(p.z)
		}
	}

	var resolution float64
	var maxResolution = (maxX - minX) / 100
	var minResolution float64 = 1
	for resolution = maxResolution; ; resolution /= 10 {
		var maxRP *point = nil
		var howManyMax = 0
		if resolution < minResolution {
			resolution = minResolution
		}
		var p1 = &point{}
		fmt.Printf("Going from (%f,%f,%f) to (%f,%f,%f), resolution=%v\n", minX, minY, minZ, maxX, maxY, maxZ, resolution)
		for x := minX; x <= maxX; x += resolution {
			//fmt.Printf("Processed %v%% so far (resolution: %f)\n", 100*(x-minX)/(maxX-minX), resolution)
			p1.x = int(x)
			for y := minY; y <= maxY; y += resolution {
				p1.y = int(y)
				for z := minZ; z <= maxZ; z += resolution {
					p1.z = int(z)
					countP1 := 0
					for _, p2 := range points {
						d1 := dist(p1, p2)
						if d1 <= p2.radius+resolution {
							countP1++
						}
					}
					if countP1 > howManyMax {
						howManyMax = countP1
						maxRP = &point{x: p1.x, y: p1.y, z: p1.z}
					} else if countP1 == howManyMax {
						if maxRP == nil || dist(zero, p1) < dist(zero, maxRP) {
							maxRP = p1
						}
					}
				}
			}
		}
		if resolution > minResolution {
			minX = float64(maxRP.x) - float64(resolution/2)
			maxX = float64(maxRP.x) + float64(resolution/2)
			minY = float64(maxRP.y) - float64(resolution/2)
			maxY = float64(maxRP.y) + float64(resolution/2)
			minZ = float64(maxRP.z) - float64(resolution/2)
			maxZ = float64(maxRP.z) + float64(resolution/2)
		} else if resolution == minResolution {
			totalMaxRP = maxRP
			totalHowManyMax = howManyMax
			break
		}
	}
	fmt.Printf("Best point is: %v since it is covered by %d bots\n", totalMaxRP, totalHowManyMax)
	return dist(totalMaxRP, zero)
}

func dist(p1, p2 *point) float64 {
	return math.Abs(float64(p1.x)-float64(p2.x)) + math.Abs(float64(p1.y)-float64(p2.y)) + math.Abs(float64(p1.z)-float64(p2.z))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
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
