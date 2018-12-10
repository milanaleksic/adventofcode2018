package main

import "testing"

func TestInput1(t *testing.T) {
	res := part1(9, 25)
	if res != 32 {
		t.Fatalf("Expected 32, got %s", res)
	}
}

func TestInput2(t *testing.T) {
	res := part1(1, 48)
	if res != 95 {
		t.Fatalf("Expected 32, got %s", res)
	}
}

func TestInput3(t *testing.T) {
	res := part1(9, 48)
	if res != 63 {
		t.Fatalf("Expected 32, got %s", res)
	}
}

func TestInput4(t *testing.T) {
	res := part1(10, 1618)
	if res != 8317 {
		t.Fatalf("Expected 32, got %s", res)
	}
}

func TestInput5(t *testing.T) {
	res := part1(13, 7999)
	if res != 146373 {
		t.Fatalf("Expected 32, got %s", res)
	}
}

func TestInput6(t *testing.T) {
	res := part1(17, 1104)
	if res != 2764 {
		t.Fatalf("Expected 32, got %s", res)
	}
}

func TestInput7(t *testing.T) {
	res := part1(21, 6111)
	if res != 54718 {
		t.Fatalf("Expected 32, got %s", res)
	}
}

func TestInput8(t *testing.T) {
	res := part1(30, 5807)
	if res != 37305 {
		t.Fatalf("Expected 32, got %s", res)
	}
}
