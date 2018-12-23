package main

import (
	"bufio"
	"container/heap"
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

// An Item is something we manage in a priority queue.
type Item struct {
	x         int
	y         int
	z         int
	radius    float64
	precision int
	priority  int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value *point, priority int) {
	item.x = value.x
	item.y = value.y
	item.z = value.z
	item.priority = priority
	heap.Fix(pq, item.index)
}

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

	//countP1 := 0
	//p1 := &point{x: 11382527, y: 29059459, z: 39808804}
	//for _, p2 := range points {
	//	d1 := dist(p1, p2)
	//	if d1 <= p2.radius {
	//		countP1++
	//	}
	//}
	//fmt.Printf("Number of bots that match: %v, distance: %v", countP1, int(dist(&point{}, p1)))

	fmt.Printf("Solution to part 2 is: %v\n", int(part2(points)))
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

	var maxResolution = (maxX - minX) / 100
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	var maxRP *point = nil
	var howManyMax = 0
	found := func(x, y, z, countP1, found int) {
		fmt.Printf("Floating countP1=%v\n", countP1)
		p1 := &point{x: x, y: y, z: z}
		if countP1 > howManyMax {
			fmt.Printf("New best sample found at %v, %v, %v, count=%v (distance=%v)\n", x, y, z, countP1, int(found))
			howManyMax = countP1
			maxRP = p1
		} else if countP1 == howManyMax {
			if maxRP == nil || dist(zero, p1) < dist(zero, maxRP) {
				fmt.Printf("New best sample found at %v, %v, %v, count=%v (distance=%v)\n", x, y, z, countP1, int(found))
				maxRP = p1
			}
		}
	}

	mapIntoPQ(int(minX), int(maxX), int(minY), int(maxY), int(minZ), int(maxZ), maxResolution, points, &pq, found)
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		mapIntoPQ(item.x-item.precision, item.x+item.precision, item.y-item.precision, item.y+item.precision, item.z-item.precision, item.z+item.precision, float64(item.precision/2), points, &pq, found)
	}

	//fmt.Printf("Best point is: %v since it is covered by %d bots\n", totalMaxRP, totalHowManyMax)
	//return dist(totalMaxRP, zero)
	return -1
}

var zero = &point{}

func mapIntoPQ(minX int, maxX int, minY int, maxY int, minZ int, maxZ int, resolution float64, points []*point, pq *PriorityQueue, foundHandler func(x, y, z, countP1, found int)) {
	if resolution < 1 {
		resolution = 1
	}
	var p = &point{}
	for x := minX; x <= maxX; x += int(resolution) {
		//fmt.Printf("Processed %v%% so far (resolution: %f)\n", 100*(x-minX)/(maxX-minX), resolution)
		p.x = int(x)
		for y := minY; y <= maxY; y += int(resolution) {
			p.y = int(y)
			for z := minZ; z <= maxZ; z += int(resolution) {
				p.z = int(z)
				countP1 := 0
				for _, p2 := range points {
					d1 := dist(p, p2)
					var radius float64
					if resolution == 1 {
						radius = p2.radius
					} else {
						radius = p2.radius + resolution*3
					}
					if d1 <= radius {
						countP1++
					}
				}
				if countP1 == 0 {
					continue
				} else if resolution == 1 {
					foundHandler(x, y, z, countP1, int(dist(zero, p)))
				} else {
					fmt.Printf("Resampling at %v, %v, %v (count:%v) with precision %v\n", x, y, z, countP1, resolution)
					item := &Item{
						x:         int(x),
						y:         int(y),
						z:         int(z),
						priority:  countP1,
						precision: int(resolution),
					}
					heap.Push(pq, item)
				}
			}
		}
	}
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
