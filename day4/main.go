package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Time struct {
	bitRepresentations *list.List
	totalTime          int
}

func main() {
	file, err := os.Open("day4/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	listInput := list.New()
	readAll(file, listInput)

	// [1518-11-01 00:00] Guard #10 begins shift
	// [1518-11-01 00:00] wakes up
	// [1518-11-01 00:00] falls asleep
	regex3, err := regexp.Compile("\\[(\\d{4})-(\\d{2})-(\\d{2}) (\\d{2}):(\\d{2})\\] (.*)")
	if err != nil {
		log.Fatal(err)
	}
	regex, err := regexp.Compile("\\[\\d{4}-\\d{2}-\\d{2} (\\d{2}):(\\d{2})\\] (.*)")
	if err != nil {
		log.Fatal(err)
	}
	regex2, err := regexp.Compile("Guard #(\\d+) begins shift")
	if err != nil {
		log.Fatal(err)
	}

	currentGuard := -1
	currentSleepStart := -1
	guardTime := make(map[int]*Time)
	work := make([]string, listInput.Len())
	iter := 0
	for line := listInput.Front(); line != nil; line = line.Next() {
		work[iter] = line.Value.(string)
		iter++
	}
	sort.Slice(work, func(i1, i2 int) bool {
		// \[(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2})\] (.*)
		w1 := work[i1]
		w2 := work[i2]
		matches1 := regex3.FindAllStringSubmatch(w1, -1)[0]
		matches2 := regex3.FindAllStringSubmatch(w2, -1)[0]
		year1, err := strconv.Atoi(matches1[1])
		if err != nil {
			log.Fatal(err)
		}
		year2, err := strconv.Atoi(matches2[1])
		if err != nil {
			log.Fatal(err)
		}
		if year1 < year2 {
			return true
		} else if year1 > year2 {
			return false
		}

		month1, err := strconv.Atoi(matches1[2])
		if err != nil {
			log.Fatal(err)
		}
		month2, err := strconv.Atoi(matches2[2])
		if err != nil {
			log.Fatal(err)
		}
		if month1 < month2 {
			return true
		} else if month1 > month2 {
			return false
		}
		day1, err := strconv.Atoi(matches1[3])
		if err != nil {
			log.Fatal(err)
		}
		day2, err := strconv.Atoi(matches2[3])
		if err != nil {
			log.Fatal(err)
		}
		if day1 < day2 {
			return true
		} else if day1 > day2 {
			return false
		}
		hour1, err := strconv.Atoi(matches1[4])
		if err != nil {
			log.Fatal(err)
		}
		hour2, err := strconv.Atoi(matches2[4])
		if err != nil {
			log.Fatal(err)
		}
		if hour1 < hour2 {
			return true
		} else if hour1 > hour2 {
			return false
		}
		min1, err := strconv.Atoi(matches1[5])
		if err != nil {
			log.Fatal(err)
		}
		min2, err := strconv.Atoi(matches2[5])
		if err != nil {
			log.Fatal(err)
		}
		if min1 < min2 {
			return true
		} else if hour1 > hour2 {
			return false
		}
		return false
	})

	for _, line := range work {
		matches := regex.FindAllStringSubmatch(line, -1)[0]
		//fmt.Println(line)
		hour, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		min, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		if hour == 23 {
			hour = 0
			min = 0
		}
		//fmt.Printf("started at min: %d\n", min)
		action := matches[3]
		var bitRepresentation uint64
		if action == "wakes up" {
			//fmt.Printf("Wakes up recognized!\n")
			soFar, ok := guardTime[currentGuard]
			amount := min - currentSleepStart
			for i := currentSleepStart; i < min; i++ {
				bitRepresentation |= uint64(1) << uint(i)
			}
			//fmt.Printf("%d - %064b\n", currentGuard, bitRepresentation)
			if ok {
				soFar.bitRepresentations.PushBack(bitRepresentation)
				soFar.totalTime = soFar.totalTime + amount
			} else {
				bitRepresentations := list.New()
				bitRepresentations.PushBack(bitRepresentation)
				guardTime[currentGuard] = &Time{totalTime: amount, bitRepresentations: bitRepresentations}
			}
			currentSleepStart = -1
		} else if action == "falls asleep" {
			//fmt.Printf("Sleep recognized!\n")
			currentSleepStart = min
		} else {
			guardId, err := strconv.Atoi(regex2.FindAllStringSubmatch(action, -1)[0][1])
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("guard id: %d\n", guardId)
			currentGuard = guardId
			currentSleepStart = -1
		}
	}

	// part 1
	part1(guardTime)

	// part 2
	part2(guardTime)
}

func part1(guardTime map[int]*Time) map[int]int {
	maxGuardId := 0
	maxGuardTime := 0
	for guardId, amountSlept := range guardTime {
		if amountSlept.totalTime > maxGuardTime {
			maxGuardTime = amountSlept.totalTime
			maxGuardId = guardId
		}
	}
	fmt.Printf("Guard %d slept for %d minutes\n", maxGuardId, maxGuardTime)
	bitRepresentations := guardTime[maxGuardId].bitRepresentations
	countsPerMin := make(map[int]int)
	for bitRepresentation := bitRepresentations.Front(); bitRepresentation != nil; bitRepresentation = bitRepresentation.Next() {
		bits := bitRepresentation.Value.(uint64)
		for b := 1; b <= 60; b++ {
			if uint64(bits)&(uint64(1)<<uint64(b)) != 0 {
				soFar, ok := countsPerMin[b]
				if ok {
					countsPerMin[b] = soFar + 1
				} else {
					countsPerMin[b] = 1
				}
			}
		}
	}
	maxCount := 0
	maxMin := 0
	for min, count := range countsPerMin {
		if count > maxCount {
			maxCount = count
			maxMin = min
		}
	}
	//fmt.Printf("%v", countsPerMin)
	fmt.Printf("Minute with most overlap is: %d, Guard id is: %d; thus Solution is: %d\n", maxMin, maxGuardId, maxMin*maxGuardId)
	return countsPerMin
}

func part2(guardTime map[int]*Time) {
	maxGuardCount := -1
	maxGuardMin := -1
	maxGuardId := -1
	for guardId, amountSlept := range guardTime {
		countsPerMin := make(map[int]int)
		for bitRepresentation := amountSlept.bitRepresentations.Front(); bitRepresentation != nil; bitRepresentation = bitRepresentation.Next() {
			bits := bitRepresentation.Value.(uint64)
			for b := 1; b <= 60; b++ {
				if uint64(bits)&(uint64(1)<<uint64(b)) != 0 {
					soFar, ok := countsPerMin[b]
					if ok {
						countsPerMin[b] = soFar + 1
					} else {
						countsPerMin[b] = 1
					}
				}
			}
		}
		maxCount := 0
		maxMin := 0
		for min, count := range countsPerMin {
			if count > maxCount {
				maxCount = count
				maxMin = min
			}
		}
		if maxCount > maxGuardCount {
			maxGuardCount = maxCount
			maxGuardId = guardId
			maxGuardMin = maxMin
		}
	}
	fmt.Printf("Minute with most overlap is: %d, Guard id is: %d; thus Solution is: %d\n", maxGuardMin, maxGuardId, maxGuardMin*maxGuardId)
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
