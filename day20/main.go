package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

type typeNode int

type coord struct {
	x, y int
}

type path struct {
	dist int
	path string
}

func (p *path) String() string {
	return fmt.Sprintf("dist=%v, path=%v", p.dist, p.path)
}

const (
	Normal typeNode = iota
	Splitter
)

type node struct {
	value    string
	children []*node
	operator typeNode
}

func main() {
	file, err := os.Open("day20/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)
	inputLine := ls.Front()
	fmt.Printf("Solution to part 1 is: %v", part1(inputLine.Value.(string)))
}

func part1(input string) int {
	var dists = make(map[coord]*path)
	iter := 0
	explodeRegex(input, func(ln int, output string) {
		if iter%5000 == 0 && iter > 0 {
			fmt.Printf("Simplification found, current size (working): %d, processed so far: %d, best proposal: %v\n", ln, iter, len(getBestProposal(dists).path))
		}
		iter++
		howFar(output, dists)
	})
	maxDist := getBestProposal(dists)
	return len(maxDist.path)
}

func getBestProposal(dists map[coord]*path) *path {
	var maxDist *path
	for _, p := range dists {
		if maxDist == nil {
			maxDist = p
		} else if p.dist > maxDist.dist {
			maxDist = p
		}
	}
	return maxDist
}

func howFar(input string, dists map[coord]*path) {
	//fmt.Printf("Received input=%v\n", input)
	locX, locY := 0, 0
	iter := 0
	for _, x := range input {
		iter++
		switch x {
		case 'N':
			locY += 1
		case 'S':
			locY -= 1
		case 'E':
			locX += 1
		case 'W':
			locX -= 1
		}
		c := coord{x: locX, y: locY}
		if k, ok := dists[c]; ok {
			if iter < k.dist {
				k.dist = iter
				k.path = input
			}
		} else {
			dists[c] = &path{
				path: input,
				dist: iter,
			}
		}
	}
	return
}

func explodeRegex(input string, outputHandler func(int, string)) {
	l, _ := makePaths(input, 0)
	queue := list.New()
	queue.PushFront(l)
	for {
		item := queue.Front()
		queue.Remove(item)
		nodes := item.Value.([]*node)
		explodedN, changed := simplify(nodes)
		if changed {
			for _, e := range explodedN {
				queue.PushFront(e)
			}
		} else {
			outputHandler(queue.Len(), toString(nodes))
		}
		if queue.Len() == 0 {
			break
		} else {

		}
	}
}

func toString(nodes []*node) string {
	buffer := ""
	for _, x := range nodes {
		buffer += x.value
	}
	return buffer
}

func simplify(nodes []*node) (result [][]*node, changed bool) {
	if len(nodes) == 0 {
		log.Fatalf("Empty node list not supported: %v", nodes)
	}
	if len(nodes) == 1 {
		return result, false
	}
	for i, n := range nodes {
		if len(n.children) == 0 {
			continue
		}
		groups := createGroups(n.children)
		options := make([][]*node, 0)
		for _, group := range groups {
			startIter := group[0]
			endIter := group[1]
			var suffix []*node
			if startIter == -1 || endIter == -1 {
				suffix = []*node{}
			} else {
				suffix = n.children[startIter:endIter]
			}
			options = append(options, suffix)
		}
		result = make([][]*node, len(options))
		for j := range options {
			result[j] = append(result[j], nodes[:i]...)
			result[j] = append(result[j], options[j]...)
			result[j] = append(result[j], nodes[i+1:]...)
		}
		return result, true
	}
	return result, false
}

func createGroups(nodes []*node) (groups [][]int) {
	if nodes == nil || len(nodes) == 0 {
		return
	}
	groupIter := 0
	for _, n := range nodes {
		if n.operator == Splitter {
			groups = append(groups, []int{-1, -1})
		}
	}
	groups = append(groups, []int{-1, -1})
	nodeIterStart := 0
	nodeIterEnd := 0
	for {
		node := nodes[nodeIterEnd]
		if node.operator == Splitter {
			groups[groupIter][0] = nodeIterStart
			groups[groupIter][1] = nodeIterEnd
			groupIter++
			nodeIterStart = nodeIterEnd + 1
			nodeIterEnd = nodeIterStart
		} else {
			nodeIterEnd++
		}
		if nodeIterEnd >= len(nodes) {
			if nodeIterEnd > nodeIterStart {
				groups[groupIter][0] = nodeIterStart
				groups[groupIter][1] = nodeIterEnd
			}
			break
		}
	}
	return
}

func makePaths(input string, loc int) (ns []*node, end int) {
	ns = make([]*node, 0)
	n := &node{
		value: "",
	}
	ns = append(ns, n)

	iter := loc
	for iter = loc; iter < len(input); iter++ {
		direction := string(input[iter])
		switch direction {
		case "N":
			if n == nil {
				n = &node{
					value: "",
				}
				ns = append(ns, n)
			}
			n.value += "N"
		case "W":
			if n == nil {
				n = &node{
					value: "",
				}
				ns = append(ns, n)
			}
			n.value += "W"
		case "E":
			if n == nil {
				n = &node{
					value: "",
				}
				ns = append(ns, n)
			}
			n.value += "E"
		case "S":
			if n == nil {
				n = &node{
					value: "",
				}
				ns = append(ns, n)
			}
			n.value += "S"
		case "|":
			n = &node{
				value:    "",
				operator: Splitter,
			}
			ns = append(ns, n)
			n = nil
		case "(":
			children, end := makePaths(input, iter+1)
			iter = end - 1
			ns = append(ns, &node{
				value:    "",
				children: children,
			})
			n = nil
		case ")":
			return ns, iter + 1
		}
	}
	return ns, iter
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
