package main

import (
	"bufio"
	"container/list"
	"log"
	"strings"
	"testing"
)

func TestInput1(t *testing.T) {
	values := []struct {
		initialState     string
		expectedEndTick  int
		expectedSolution int
	}{
		{`#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`, 47, 27730},
		{`#######
#G..#E#
#E#E.E#
#G.##.#
#...#E#
#...E.#
#######`, 37, 36334},
		{`#######
#E..EG#
#.#G.E#
#E.##E#
#G..#.#
#..E#.#
#######`, 46, 39514},
		{`#######
#E.G#.#
#.#G..#
#G.#.G#
#G..#.#
#...E.#
#######`, 35, 27755},
		{`#######
#.E...#
#.#..G#
#.###.#
#E#G#G#
#...#G#
#######`, 54, 28944},
		{`#########
#G......#
#.E.#...#
#..##..G#
#...##..#
#...#...#
#.G...G.#
#.....G.#
#########`, 20, 18740},
	}
	for _, v := range values {
		actualEndTick, actualSolution := part1(fromInput(makeList(v.initialState)))
		if v.expectedEndTick != actualEndTick || actualSolution != v.expectedSolution {
			t.Errorf("Failed to match result tick=%d (expected %d), solution=%d (expected %d)", actualEndTick, v.expectedEndTick, actualSolution, v.expectedSolution)
		} else {
			t.Log("Match!")
		}
	}
}

func makeList(input string) *list.List {
	ls := list.New()
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		val := scanner.Text()
		ls.PushBack(val)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ls
}
