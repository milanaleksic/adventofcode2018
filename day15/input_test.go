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
		//{"#####\n#G..#\n#E#E#\n#G.##\n#####", 37, 36334},

		{"#######\n#.G...#\n#...EG#\n#.#.#G#\n#..G#E#\n#.....#\n#######", 47, 27730},
		{"#######\n#G..#E#\n#E#E.E#\n#G.##.#\n#...#E#\n#...E.#\n#######", 37, 36334},
		{"#######\n#E..EG#\n#.#G.E#\n#E.##E#\n#G..#.#\n#..E#.#\n#######", 46, 39514},
		{"#######\n#E.G#.#\n#.#G..#\n#G.#.G#\n#G..#.#\n#...E.#\n#######", 35, 27755},
		{"#######\n#.E...#\n#.#..G#\n#.###.#\n#E#G#G#\n#...#G#\n#######", 54, 28944},
		{"#########\n#G......#\n#.E.#...#\n#..##..G#\n#...##..#\n#...#...#\n#.G...G.#\n#.....G.#\n#########", 20, 18740},
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
