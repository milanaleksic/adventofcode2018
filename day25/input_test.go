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
		input                  string
		expectedConstellations int
	}{
		{` 0,0,0,0
		3,0,0,0
		0,3,0,0
		0,0,3,0
		0,0,0,3
		0,0,0,6
		9,0,0,0
		12,0,0,0`, 2},
		{
			`-1,2,2,0
				0,0,2,-2
				0,0,0,-2
				-1,2,0,0
				-2,-2,-2,2
				3,0,2,-1
				-1,3,2,2
				-1,0,-1,0
				0,2,1,-2
				3,0,0,0`, 4},
		{`1,-1,0,1
		2,0,-1,0
		3,2,-1,0
		0,0,3,1
		0,0,-1,-1
		2,3,-2,0
		-2,2,0,0
		2,-2,0,-1
		1,-1,0,-1
		3,2,0,2`, 3},
		{`1,-1,-1,-2
		-2,-2,0,1
		0,2,1,3
		-2,3,-2,1
		0,2,3,-2
		-1,-1,1,-2
		0,-2,-1,0
		-2,2,3,-1
		1,2,2,0
		-1,-2,0,-2`, 8},
	}
	for _, v := range values {
		t.Run(v.input, func(t *testing.T) {
			actualConstellations := part1(fromInput(makeList(v.input)))
			if v.expectedConstellations != actualConstellations {
				t.Errorf("Failed to match result actualConstellations=%d (expected %d)", actualConstellations, v.expectedConstellations)
			}
		})
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
