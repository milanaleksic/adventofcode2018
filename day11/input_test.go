package main

import "testing"

func TestInput1(t *testing.T) {
	values := []struct {
		x, y, serial, expected int
	}{
		{3, 5, 8, 4},
		{122, 79, 57, -5},
		{217, 196, 39, 0},
		{101, 153, 71, 4},
	}
	for _, v := range values {
		calc := calculate(v.x, v.y, v.serial)
		if calc != v.expected {
			t.Errorf("Failed to match result %d, expected %d", calc, v.expected)
		} else {
			t.Log("Match!")
		}
	}
}

func TestSolution1(t *testing.T) {
	values := []struct {
		x, y, serial, block, expectedX, expectedY int
	}{
		{300, 300, 18, 3, 33, 45},
		{300, 300, 42, 3, 21, 61},
	}
	for _, v := range values {
		_, solutionX, solutionY := part1(v.x, v.y, v.serial, v.block)
		if solutionX != v.expectedX || solutionY != v.expectedY {
			t.Errorf("Failed to match result %d,%d; expected %d,%d", solutionX, solutionY, v.expectedX, v.expectedY)
		} else {
			t.Log("Match!")
		}
	}
}
