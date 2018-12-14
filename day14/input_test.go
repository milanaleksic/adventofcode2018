package main

import "testing"

func TestInput1(t *testing.T) {
	values := []struct {
		recipeCount int
		expected    string
	}{
		{9, "5158916779"},
		{5, "0124515891"},
		{18, "9251071085"},
		{2018, "5941429882"},
	}
	for _, v := range values {
		actual := part1(v.recipeCount)
		if v.expected != actual {
			t.Errorf("Failed to match result %d; expected %d", actual, v.expected)
		} else {
			t.Log("Match!")
		}
	}
}

func TestInput2(t *testing.T) {
	values := []struct {
		recipeScores string
		expected     int
	}{
		{"51589", 9},
		{"01245", 5},
		{"92510", 18},
		{"59414", 2018},
	}
	for _, v := range values {
		actual := part2(v.recipeScores)
		if v.expected != actual {
			t.Errorf("Failed to match result %d; expected %d", actual, v.expected)
		} else {
			t.Log("Match!")
		}
	}
}
