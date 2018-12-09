package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// test
//const input = "day8/test.txt"

//production
const input = "day8/input.txt"

type node struct {
	children []*node
	metadata []int
}

func (n *node) String() string {
	return fmt.Sprintf("{children=%+v, metadata=%+v}", n.children, n.metadata)
}

func main() {
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)
	for line := ls.Front(); line != nil; line = line.Next() {
		data := strings.Split(line.Value.(string), " ")
		node, _ := readNode(data, 0)
		part1(node)
		part2(node)
	}
}

func part1(root *node) {
	//fmt.Printf("Tree: %+v\n", node)
	fmt.Printf("Solution to part 1 is: %v\n", sum(root))
}

func sum(root *node) (result int) {
	for _, metadata := range root.metadata {
		result += metadata
	}
	for _, child := range root.children {
		result += sum(child)
	}
	return result
}

func part2(root *node) {
	//fmt.Printf("Tree: %+v\n", node)
	fmt.Printf("Solution to part 2 is: %v\n", value(root))
}

func value(node *node) (result int) {
	//fmt.Printf("inside node starting on %v, len=%v, children=%v\n", node, len(node.children), node.children)
	if len(node.children) == 0 {
		for _, metadata := range node.metadata {
			result += metadata
		}
	} else {
		for _, metadata := range node.metadata {
			if metadata >= 1 && metadata <= len(node.children) {
				result += value(node.children[metadata-1])
			}
		}
	}
	return result
}

func readNode(data []string, start int) (result *node, end int) {
	result = &node{children: make([]*node, 0), metadata: make([]int, 0)}
	//fmt.Printf("reading data, starting from %v\n", start)
	noChildren, err := strconv.Atoi(data[start])
	if err != nil {
		log.Fatalf("Could not deduce the number of children: %v, %v", data[start], err)
	}
	noMetadata, err := strconv.Atoi(data[start+1])
	if err != nil {
		log.Fatalf("Could not deduce the number of metadata: %v, %v", data[start], err)
	}
	iter := start + 2
	for i := 1; i <= noChildren; i++ {
		child, end := readNode(data, iter)
		iter = end
		result.children = append(result.children, child)
	}
	for i := iter; i < iter+noMetadata; i++ {
		//fmt.Printf("Accessing %v, metadata count=%v\n", i, noMetadata)
		metadata, err := strconv.Atoi(data[i])
		if err != nil {
			log.Fatalf("Could not deduce metadata: %v, %v", data[start], err)
		}
		result.metadata = append(result.metadata, metadata)
	}
	return result, iter + noMetadata
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
