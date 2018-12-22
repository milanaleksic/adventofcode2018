package main

import (
	"testing"
)

func TestInput1(t *testing.T) {
	values := []struct {
		initialState      string
		expectedMaxLength int
	}{
		{`^NESW$`, 4},
		{"^WSS(S|NE)$", 3}, // my examples
		{`^N(EE|N)N$`, 4},
		//
		{`^EEE(WWEE|)$`, 3},
		{`^WNE$`, 3}, //their examples
		{`^ENWWW(NEEE|SSE(EE|N))$`, 10},
		{`^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$`, 18},
		{`^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$`, 23},
		{"^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$", 31},
	}
	for _, v := range values {
		t.Run(v.initialState, func(t *testing.T) {
			actualMaxLength := part1(v.initialState)
			if v.expectedMaxLength != actualMaxLength {
				t.Errorf("Failed to match result actualMaxLength=%d (expected %d)", actualMaxLength, v.expectedMaxLength)
			}
		})
	}
}

func TestExplode(t *testing.T) {
	logger := func(output map[string]bool) func(path string) {
		return func(path string) {
			output[path] = true
		}
	}
	values := []struct {
		initialState string
		expectedMap  map[string]bool
	}{
		{`^NESW$`, map[string]bool{
			"NESW": true,
		}},
		{"^WSS(S|NE)$", map[string]bool{
			"WSSS":  true,
			"WSSNE": true,
		}},
		{`^N(EE|N)N$`, map[string]bool{
			"NEEN": true,
			"NNN":  true,
		}},
		{`^ENWWW(SSE(EE|N))$`, map[string]bool{
			"ENWWWSSEEE": true,
			"ENWWWSSEN":  true,
		}},
		{`^EEE(WWEE|)$`, map[string]bool{
			"EEE":     true,
			"EEEWWEE": true,
		}},
		{`^EEE(|WWEE)$`, map[string]bool{
			"EEE":     true,
			"EEEWWEE": true,
		}},
		{`^ENWWW(NEEE|SSE(EE|N))$`, map[string]bool{
			"ENWWWNEEE":  true,
			"ENWWWSSEEE": true,
			"ENWWWSSEN":  true,
		}},
		{`^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$`, map[string]bool{
			"ENNWSWWSSSEENEENNN":             true,
			"ENNWSWWSSSEENEESWENNNN":         true,
			"ENNWSWWSSSEENWNSEEENNN":         true,
			"ENNWSWWSSSEENWNSEEESWENNNN":     true,
			"ENNWSWWNEWSSSSEENEENNN":         true,
			"ENNWSWWNEWSSSSEENEESWENNNN":     true,
			"ENNWSWWNEWSSSSEENWNSEEENNN":     true,
			"ENNWSWWNEWSSSSEENWNSEEESWENNNN": true,
		}},
		{`^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$`, map[string]bool{
			"ESSWWNE":                 true,
			"ESSWWNNNENNEESSSSS":      true,
			"ESSWWNNNENNEESSWNSESSS":  true,
			"ESSWWNNNENNWWWSSSSESW":   true,
			"ESSWWNNNENNWWWSSSSENNNE": true,
		}},
		{`^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$`, map[string]bool{
			`WSSEESWWWNWS`:                    true,
			`WSSEESWWWNWNENNEEEENNESSSSWNWSW`: true,
			`WSSEESWWWNWNENNEEEENNESSSSWSSEN`: true,
			`WSSEESWWWNWNENNEEEENNWSWWNE`:     true,
			`WSSEESWWWNWNENNEEEENNWSWWNWWSE`:  true,
			`WSSEESWWWNWNENNEEEENNWSWWNWWSSS`: true,
		}},
	}
	for _, v := range values {
		t.Run(v.initialState, func(t *testing.T) {
			actualMap := make(map[string]bool)
			explodeRegex(v.initialState, logger(actualMap))
			if len(actualMap) != len(v.expectedMap) {
				t.Errorf("Failed to match size of the maps actual=%d (expected %d), received map=%v", len(actualMap), len(v.expectedMap), actualMap)
			}
			for key := range v.expectedMap {
				_, ok := actualMap[key]
				if !ok {
					t.Errorf("Failed to match key %v, received map=%v", key, actualMap)
				}
			}
		})
	}
}

//func TestCoord(t *testing.T) {
//	values := []struct {
//		path             string
//		expectedDistance float64
//	}{
//		// my examples
//		{"NWSE", 0},
//		{"NN", 2},
//		{"NW", math.Sqrt2},
//	}
//	for _, v := range values {
//		t.Run(v.path, func(t *testing.T) {
//			dists := make(map[coord]*path)
//			howFar(v.path, dists)
//			fmt.Printf("dist: %v\n", dists)
//			//if math.Abs(actualDistance-v.expectedDistance) > 0.1 {
//			//	t.Errorf("Failed to match result for path %v, expected %v but got %v", v.path, v.expectedDistance, actualDistance)
//			//}
//		})
//	}
//}

//func TestSimplify(t *testing.T) {
//	nodes := []*node{
//		{
//			value:    "X",
//			operator: And,
//		},
//		{
//			value:    "C",
//			operator: And,
//		},
//		{
//			children: []*node{
//				{
//					value:    "E",
//					operator: Or,
//				},
//				{
//					value:    "F",
//					operator: Or,
//				},
//			},
//		},
//	}
//	result, ok := simplify(nodes)
//	if !ok {
//		t.Fatalf("Expected to simplify!")
//	}
//	if len(result) != 2 {
//		t.Fatalf("Length is not 2 but %v", len(result))
//	}
//	if result[0][2].value != "E" {
//		t.Errorf("E is not matched")
//	}
//	if result[1][2].value != "F" {
//		t.Errorf("F is not matched")
//	}
//}
