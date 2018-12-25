package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x float64
	y float64
	z float64
	t float64
}

func (p *point) String() string {
	return fmt.Sprintf("(%v,%v,%v,%v)", p.x, p.y, p.z, p.t)
}

func (p *point) dist(p2 *point) float64 {
	result := math.Abs(p.x-p2.x) + math.Abs(p.y-p2.y) + math.Abs(p.z-p2.z) + math.Abs(p.t-p2.t)
	//fmt.Printf("Distance between %v and %v is %v\n", p, p2, result)
	return result
}

func (p *point) inConstellationWith(p2 *point) bool {
	return p.dist(p2) <= 3
}

func main() {
	file, err := os.Open("day25/input.txt")
	check(err)
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	fmt.Printf("Solution to part 1 is: %v\n", part1(fromInput(ls)))
}

func part1(points []*point) int {
	constellations := make([][]*point, 0)
	for _, p := range points {
		constellationIds := make([]int, 0)
		for i, constellation := range constellations {
			for _, member := range constellation {
				if member.inConstellationWith(p) {
					constellationIds = append(constellationIds, i)
				}
			}
		}
		if len(constellationIds) == 0 {
			constellations = append(constellations, []*point{p})
		} else if len(constellationIds) == 1 {
			constellations[constellationIds[0]] = append(constellations[constellationIds[0]], p)
		} else {
			newConstellation := make([]*point, 0)
			for _, oldConstellationId := range constellationIds {
				newConstellation = append(newConstellation, constellations[oldConstellationId]...)
				constellations[oldConstellationId] = []*point{}
			}
			newConstellation = append(newConstellation, p)
			constellations = append(constellations, newConstellation)
		}
	}
	count := 0
	for _, constellation := range constellations {
		if len(constellation) > 0 {
			count++
		}
	}
	return count
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

func fromInput(ls *list.List) []*point {
	points := make([]*point, 0)
	for line := ls.Front(); line != nil; line = line.Next() {
		lineContents := line.Value.(string)
		coords := strings.Split(lineContents, ",")
		x := mustInt(coords[0])
		y := mustInt(coords[1])
		z := mustInt(coords[2])
		t := mustInt(coords[3])
		points = append(points, &point{x: float64(x), y: float64(y), z: float64(z), t: float64(t)})
	}
	return points
}

func mustInt(str string) int {
	result, err := strconv.Atoi(strings.TrimSpace(str))
	check(err)
	return result
}
