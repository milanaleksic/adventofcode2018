package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// test
//const input = "day9/test.txt"

//production
const input = "day9/input.txt"

func main() {
	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	regex, err := regexp.Compile("(\\d+) players; last marble is worth (\\d+) points")
	if err != nil {
		log.Fatal(err)
	}

	ls := list.New()
	readAll(file, ls)
	for line := ls.Front(); line != nil; line = line.Next() {
		matches := regex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		players, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		lastMarble, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found players=%d, lastMarble=%d\n", players, lastMarble)
		part1(players, lastMarble)
		part1(players, lastMarble*100)
	}
}

type player struct {
	points int
}

func part1(noPlayers int, lastMarbleWorth int) int {
	players := make([]*player, noPlayers)
	for i := 0; i < noPlayers; i++ {
		players[i] = &player{}
	}
	playingField := []int{0}
	marble := 0
	iter := 0
	loc := -1
	for {
		iter = loc
		marble++
		if marble%100000 == 0 {
			fmt.Println("Reached marble", marble)
		}
		player := players[marble%noPlayers]
		if marble%23 == 0 {
			//fmt.Printf("------------\nPlayer %d should get some points\n", player)
			iter -= 7
			if iter >= 0 {
				loc = iter % len(playingField)
			} else {
				loc = len(playingField) + iter
			}
			//fmt.Printf("len=%v, loc=%v\n", len(playingField), loc)
			player.points += marble
			player.points += playingField[loc]
			prev := len(playingField)
			playingField = append(playingField[:loc], playingField[loc+1:]...)
			if len(playingField) != prev-1 {
				log.Fatalf("Wrong size!")
				break
			}
		} else {
			iter += 2
			loc = iter % len(playingField)
			//fmt.Printf("------------\nmarble=%d iter=%d loc=%d len(playingField)=%v before=%v\n", marble, iter, loc, len(playingField)+1, playingField)
			if loc == 0 {
				//fmt.Println("Adding to the end only!")
				playingField = append(playingField, marble)
				loc = len(playingField) - 1
			} else {
				//fmt.Printf("combining like this: ...:%d and %d:..., len=%v\n", loc, loc, len(playingField))
				newPlayingField := make([]int, len(playingField)+1)
				for i := 0; i < loc; i++ {
					newPlayingField[i] = playingField[i]
				}
				newPlayingField[loc] = marble
				for i := loc + 1; i <= len(playingField); i++ {
					newPlayingField[i] = playingField[i-1]
				}
				playingField = newPlayingField
			}
			//fmt.Printf("after: %v of size %d\n", playingField, len(playingField))
		}
		if marble >= lastMarbleWorth {
			//fmt.Printf("Final: %v\n", playingField)
			break
		}
	}
	max := 0
	for _, player := range players {
		if player.points > max {
			max = player.points
		}
	}
	fmt.Printf("Solution is: %v\n", max)
	return max
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
