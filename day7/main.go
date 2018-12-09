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

// test
//const numWorkers = 2
//const input = "day7/test.txt"
//const cost = 0

//production
const numWorkers = 5
const input = "day7/input.txt"
const cost = 60

type depends struct {
	dependsOn map[string]bool
	goesTo    map[string]bool
}

type workerData struct {
	busyUntil int
	workingOn string
}

func (d *depends) String() string {
	return fmt.Sprintf("dependsOn=%+v, goesTo=%+v", d.dependsOn, d.goesTo)
}

func main() {
	file, err := os.Open(input)
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

	part1(graph)
	part2(graph)
}

func part2(graphInput map[string]*depends) {
	graph := copyGraph(graphInput)
	var workers = make([]*workerData, numWorkers)
	for i := range workers {
		workers[i] = &workerData{
			workingOn: "",
			busyUntil: -1,
		}
	}
	var visited = make(map[string]bool)
	var available = make(map[string]bool)
	var ongoing = make(map[string]bool)
	var time = 0
	for point, dep := range graph {
		if len(dep.dependsOn) == 0 {
			available[point] = true
		}
	}
	fmt.Printf("available: %+v\n", available)
	for {
		time++
		fmt.Printf("t=%v", time)
		for _, worker := range workers {
			if worker.busyUntil != -1 && worker.busyUntil >= time {
				continue
			}
			//fmt.Printf("Worker %v is available!\n", workerId)

			if worker.workingOn != "" {
				delete(ongoing, worker.workingOn)
				visited[worker.workingOn] = true
				for followingFromIter := range graph[worker.workingOn].goesTo {
					delete(graph[followingFromIter].dependsOn, worker.workingOn)
				}
				for point, dep := range graph {
					if len(dep.dependsOn) == 0 {
						_, ok := visited[point]
						if !ok {
							_, ok := ongoing[point]
							if !ok {
								_, ok := available[point]
								if !ok {
									fmt.Printf(" (became available since %v is done: %+v)", worker.workingOn, point)
									available[point] = true
								}
							}
						}
					}
				}
			}
		}
		for _, worker := range workers {
			if worker.busyUntil != -1 && worker.busyUntil >= time {
				continue
			}
			iter := min(available)
			if iter == "" {
				// no work available at this time
				worker.busyUntil = -1
				worker.workingOn = ""
				continue
			}
			delete(available, iter)
			worker.busyUntil = time + cost + int([]byte(iter)[0]-'A')
			worker.workingOn = iter
			ongoing[iter] = true
		}
		ongoing := false
		for workerId, worker := range workers {
			if worker.workingOn != "" {
				if worker.busyUntil > time {
					ongoing = true
				}
				fmt.Printf(" %d:%s", workerId, worker.workingOn)
			}
		}
		fmt.Printf(" available: %+v", available)
		fmt.Println()
		if len(graph) == len(visited) && !ongoing {
			fmt.Printf("Ending on moment: %d\n", time-1)
			break
		}
	}
}

func part1(graphInput map[string]*depends) {
	graph := copyGraph(graphInput)
	var visited = make(map[string]bool)
	var active = make(map[string]bool)
	var result = make([]string, 0)
	for point, dep := range graph {
		if len(dep.dependsOn) == 0 {
			active[point] = true
		}
	}
	for {
		iter := min(active)
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

func copyGraph(graphInput map[string]*depends) map[string]*depends {
	var graph = make(map[string]*depends)
	for taskName, dep := range graphInput {
		depNew := &depends{goesTo: make(map[string]bool), dependsOn: make(map[string]bool)}
		graph[taskName] = depNew
		for taskName := range dep.dependsOn {
			depNew.dependsOn[taskName] = true
		}
		for taskName := range dep.goesTo {
			depNew.goesTo[taskName] = true
		}
	}
	return graph
}

func min(candidates map[string]bool) string {
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
