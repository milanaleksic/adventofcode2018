package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

type state struct {
	registers []int
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

type instrLogic func(input1, input2, output int, inputState, outputState *state)

func main() {
	file, err := os.Open("day16/input.txt")
	//file, err := os.Open("day16/test.txt")
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
				registers: []int{
					forceInt(matchesBefore[1]),
					forceInt(matchesBefore[2]),
					forceInt(matchesBefore[3]),
					forceInt(matchesBefore[4]),
				},
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
				registers: []int{
					forceInt(matchesAfter[1]),
					forceInt(matchesAfter[2]),
					forceInt(matchesAfter[3]),
					forceInt(matchesAfter[4]),
				},
			}

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
	architecture := createKnownArchitecture()
	numberOfMoreThanThreeMatches := part1(architecture, behaviors)
	fmt.Printf("Number of samples that have behavior like >=3 instructions: %v", numberOfMoreThanThreeMatches)
}

func part1(architecture map[string]instrLogic, behaviors []*behavior) (numberOfMoreThanThreeMatches int) {
	for _, b := range behaviors {
		numberOfMatches := 0
		for opName, op := range architecture {
			sandbox := &state{
				registers: make([]int, 4),
			}
			for i := 0; i < len(b.before.registers); i++ {
				sandbox.registers[i] = b.before.registers[i]
			}
			op(b.instruction.input1, b.instruction.input2, b.instruction.output, sandbox, sandbox)
			if reflect.DeepEqual(sandbox, b.after) {
				fmt.Println("Behaves like: ", opName)
				numberOfMatches++
			}
		}
		if numberOfMatches >= 3 {
			numberOfMoreThanThreeMatches++
		}
	}
	return
}

func createKnownArchitecture() map[string]instrLogic {
	knownOperations := make(map[string]instrLogic, 0)
	knownOperations["addr"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1] + inputData.registers[input2]
	}
	knownOperations["addi"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1] + input2
	}
	knownOperations["mulr"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1] * inputData.registers[input2]
	}
	knownOperations["muli"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1] * input2
	}
	knownOperations["banr"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1] & inputData.registers[input2]
	}
	knownOperations["bani"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1] & input2
	}
	knownOperations["borr"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1] | inputData.registers[input2]
	}
	knownOperations["bori"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1] | input2
	}
	knownOperations["setr"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = inputData.registers[input1]
	}
	knownOperations["seti"] = func(input1, input2, output int, inputData *state, outputData *state) {
		outputData.registers[output] = input1
	}
	knownOperations["gtir"] = func(input1, input2, output int, inputData *state, outputData *state) {
		if input1 > inputData.registers[input2] {
			outputData.registers[output] = 1
		} else {
			outputData.registers[output] = 0
		}
	}
	knownOperations["gtri"] = func(input1, input2, output int, inputData *state, outputData *state) {
		if inputData.registers[input1] > input2 {
			outputData.registers[output] = 1
		} else {
			outputData.registers[output] = 0
		}
	}
	knownOperations["gtrr"] = func(input1, input2, output int, inputData *state, outputData *state) {
		if inputData.registers[input1] > inputData.registers[input2] {
			outputData.registers[output] = 1
		} else {
			outputData.registers[output] = 0
		}
	}
	knownOperations["eqir"] = func(input1, input2, output int, inputData *state, outputData *state) {
		if input1 == inputData.registers[input2] {
			outputData.registers[output] = 1
		} else {
			outputData.registers[output] = 0
		}
	}
	knownOperations["eqri"] = func(input1, input2, output int, inputData *state, outputData *state) {
		if inputData.registers[input1] == input2 {
			outputData.registers[output] = 1
		} else {
			outputData.registers[output] = 0
		}
	}
	knownOperations["eqrr"] = func(input1, input2, output int, inputData *state, outputData *state) {
		if inputData.registers[input1] == inputData.registers[input2] {
			outputData.registers[output] = 1
		} else {
			outputData.registers[output] = 0
		}
	}
	return knownOperations
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
