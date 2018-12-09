package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type depends struct {
	dependsOn map[string]bool
	goesTo    map[string]bool
}

func (d *depends) String() string {
	return fmt.Sprintf("dependsOn=%+v, goesTo=%+v", d.dependsOn, d.goesTo)
}

func main() {
	file, err := os.Open("day7/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// Step X must be finished before step U can begin.
	regex, err := regexp.Compile("Step (.*) must be finished before step (.*) can begin.")
	if err != nil {
		log.Fatal(err)
	}

	// part 1
	graph := make(map[string]*depends)
	for line := ls.Front(); line != nil; line = line.Next() {
		matches := regex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		x := matches[1]
		y := matches[2]
		//fmt.Printf("Found x=%s, y=%s\n", x, y)
		_, ok := graph[x]
		if !ok {
			graph[x] = &depends{goesTo: make(map[string]bool), dependsOn: make(map[string]bool)}
		}
		_, ok = graph[y]
		if !ok {
			graph[y] = &depends{goesTo: make(map[string]bool), dependsOn: make(map[string]bool)}
		}

		graph[x].goesTo[y] = true
		graph[y].dependsOn[x] = true
	}
	//fmt.Printf("%+v\n", graph)

	var visited = make(map[string]bool)
	var active = make(map[string]bool)
	for point, dep := range graph {
		if len(dep.dependsOn) == 0 {
			active[point] = true
		}
	}
	result := make([]string, 0)
	iter := min(active)
	for {
		//fmt.Printf("Visiting: %s\n", iter)
		visited[iter] = true
		delete(active, iter)
		result = append(result, iter)

		for followingFromIter := range graph[iter].goesTo {
			delete(graph[followingFromIter].dependsOn, iter)
		}
		for point, dep := range graph {
			if len(dep.dependsOn) == 0 {
				_, ok := visited[point]
				if !ok {
					active[point] = true
				}
			}
		}
		if len(graph) == len(visited) {
			break
		}
		iter = min(active)
	}
	fmt.Println("Solution is: ", strings.Join(result, ""))
}

func min(candidates map[string]bool) string {
	if len(candidates) == 0 {
		log.Fatal("No more candidates")
	}
	//fmt.Printf("Finding minimum inside: %+v\n", candidates)
	result := ""
	for point := range candidates {
		if before(point, result) {
			result = point
		}
	}
	return result
}

func before(a, b string) bool {
	return b == "" || strings.Compare(a, b) < 0
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
