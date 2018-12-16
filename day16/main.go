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

type state struct {
	reg1, reg2, reg3, reg4 int
}

type behavior struct {
	before      *state
	after       *state
	instruction *instruction
}

func (b *behavior) String() string {
	return fmt.Sprintf("%+v -> %+v (via %v)", b.before, b.after, b.instruction)
}

type instruction struct {
	instructionType int
	input1, input2  int
	output          int
}

func main() {
	file, err := os.Open("day16/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// Before: [1, 1, 2, 2]
	beforeRegex, err := regexp.Compile(`Before:\s*\[(\d+), (\d+), (\d+), (\d+)\]`)
	if err != nil {
		log.Fatal(err)
	}
	// 9 0 1 0
	instructionRegex, err := regexp.Compile(`(\d+) (\d+) (\d+) (\d+)`)
	if err != nil {
		log.Fatal(err)
	}
	// After: [1, 1, 2, 2]
	afterRegex, err := regexp.Compile(`After:\s*\[(\d+), (\d+), (\d+), (\d+)\]`)
	if err != nil {
		log.Fatal(err)
	}

	behaviors := make([]*behavior, 0)
	instructions := make([]*instruction, 0)
	instructionReading := false
	for line := ls.Front(); line != nil; line = line.Next() {
		if line.Value.(string) == "" {
			continue
		}
		if !instructionReading {
			fullMatchesBefore := beforeRegex.FindAllStringSubmatch(line.Value.(string), -1)
			if len(fullMatchesBefore) == 0 {
				instructionReading = true
			}
		}
		if instructionReading {
			// part 2!
			matchesInstruction := instructionRegex.FindAllStringSubmatch(line.Value.(string), -1)[0]
			instrType := forceInt(matchesInstruction[1])
			input1 := forceInt(matchesInstruction[2])
			input2 := forceInt(matchesInstruction[3])
			output := forceInt(matchesInstruction[4])
			instruction := &instruction{
				instructionType: instrType,
				input1:          input1,
				input2:          input2,
				output:          output,
			}
			fmt.Printf("Found instruction x=%v\n", instruction)
			instructions = append(instructions, instruction)
		} else {
			fullMatchesBefore := beforeRegex.FindAllStringSubmatch(line.Value.(string), -1)
			matchesBefore := fullMatchesBefore[0]
			beforeState := &state{
				reg1: forceInt(matchesBefore[1]),
				reg2: forceInt(matchesBefore[2]),
				reg3: forceInt(matchesBefore[3]),
				reg4: forceInt(matchesBefore[4]),
			}
			line = line.Next()

			matchesInstruction := instructionRegex.FindAllStringSubmatch(line.Value.(string), -1)[0]
			instrType := forceInt(matchesInstruction[1])
			input1 := forceInt(matchesInstruction[2])
			input2 := forceInt(matchesInstruction[3])
			output := forceInt(matchesInstruction[4])
			line = line.Next()

			fullMatchesAfter := afterRegex.FindAllStringSubmatch(line.Value.(string), -1)
			matchesAfter := fullMatchesAfter[0]
			afterState := &state{
				reg1: forceInt(matchesAfter[1]),
				reg2: forceInt(matchesAfter[2]),
				reg3: forceInt(matchesAfter[3]),
				reg4: forceInt(matchesAfter[4]),
			}
			line = line.Next()

			b := behavior{
				before: beforeState,
				after:  afterState,
				instruction: &instruction{
					instructionType: instrType,
					input1:          input1,
					input2:          input2,
					output:          output,
				},
			}
			fmt.Printf("Found behavior x=%v\n", &b)
			behaviors = append(behaviors, &b)
		}
	}

}

func forceInt(x string) int {
	result, err := strconv.Atoi(x)
	if err != nil {
		log.Fatalf("Could not convert %x to integer", x)
	}
	return result
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
