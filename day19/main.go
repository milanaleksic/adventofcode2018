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
	registers []int
}

type instruction struct {
	instructionType string
	input1, input2  int
	output          int
}

type instrLogic func(input1, input2, output int, inputState, outputState *state)

func main() {
	file, err := os.Open("day19/input.txt")
	//file, err := os.Open("day19/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// #ip 0
	firstInstruction, err := regexp.Compile(`#ip (\d+)`)
	if err != nil {
		log.Fatal(err)
	}
	// seti 5 0 1
	instructionRegex, err := regexp.Compile(`([a-z]+) (\d+) (\d+) (\d+)`)
	if err != nil {
		log.Fatal(err)
	}

	instructions := make([]*instruction, 0)
	line := ls.Front()
	matches := firstInstruction.FindAllStringSubmatch(line.Value.(string), -1)[0]
	ipRegister := forceInt(matches[1])
	line = line.Next()

	for ; line != nil; line = line.Next() {
		if line.Value.(string) == "" {
			continue
		}
		matchesInstruction := instructionRegex.FindAllStringSubmatch(line.Value.(string), -1)[0]
		instrType := matchesInstruction[1]
		input1 := forceInt(matchesInstruction[2])
		input2 := forceInt(matchesInstruction[3])
		output := forceInt(matchesInstruction[4])
		instructions = append(instructions, &instruction{
			instructionType: instrType,
			input1:          input1,
			input2:          input2,
			output:          output,
		})
	}
	//for ind, instr := range instructions {
	//	fmt.Printf("Added instruction %v %v %v %v %v\n", ind, instr.instructionType, instr.input1, instr.input2, instr.output)
	//}

	architecture := createKnownArchitecture()

	// PART 1
	cpu := &state{
		registers: []int{0, 0, 0, 0, 0, 0},
	}
	iter := 0
	maxIter := 1000000000000
	activeInstruction := cpu.registers[ipRegister]
	//fmt.Printf("CPU: %v\n", cpu.registers)
	last0 := -1
	for {
		//fmt.Printf("iter=%v ip=%v %v", iter, activeInstruction, cpu.registers)
		iter++
		i := instructions[activeInstruction]
		architecture[i.instructionType](i.input1, i.input2, i.output, cpu, cpu)
		activeInstruction = cpu.registers[ipRegister]
		//fmt.Printf(" %v ", i)
		activeInstruction++
		if activeInstruction >= len(instructions) || iter >= maxIter {
			fmt.Printf("(FINAL PART 1) %v\n", cpu.registers)
			break
		}
		cpu.registers[ipRegister] = activeInstruction
		//fmt.Printf("%v\n", cpu.registers)
		if cpu.registers[0] != last0 {
			if cpu.registers[0] == 0 {
				continue
			}
			if last0 == -1 {
				fmt.Println("New value added: ", cpu.registers[0])
			} else {
				fmt.Println("New value added: ", cpu.registers[0]-last0)
			}
			last0 = cpu.registers[0]
		}
	}

	//PART 2
	fmt.Printf("Optimized calcultion for part 1 is: %v\n", sumFactors(1008))
	fmt.Printf("Optimized calcultion for part 2 is: %v\n", sumFactors(10551408))
}

func sumFactors(number int) (result int) {
	for i := 1; i <= number/2; i++ {
		if number%i == 0 {
			fmt.Println("Adding factor", i)
			result += i
		}
	}
	fmt.Println("Adding factor", number)
	return result + number
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
