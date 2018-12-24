package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	//input = "day24/input.txt"
	input = "day24/test.txt"
)

type attackType int

func (a attackType) String() string {
	switch a {
	case fire:
		return "fire"
	case cold:
		return "cold"
	case bludgeoning:
		return "bludgeoning"
	case radiation:
		return "radiation"
	case slashing:
		return "slashing"
	}
	log.Fatalf("Unknown attack type: %f", a)
	return ""
}

const (
	fire attackType = iota
	cold
	bludgeoning
	radiation
	slashing
)

type group struct {
	id               int
	units            int
	hitPointsPerUnit int
	weakTo           []attackType
	immuneTo         []attackType
	might            int
	attack           attackType
	initiative       int
}

func (g *group) String() string {
	return fmt.Sprintf("Group %v with %v units with %v hit points (weak to: %+v, immune to: %+v), with might %v and initiative %v", g.id, g.units, g.hitPointsPerUnit, g.weakTo, g.immuneTo, g.might, g.initiative)
}

type army []*group

func main() {
	file, err := os.Open(input)
	check(err)
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	// 801 units each with 4706 hit points (weak to bludgeoning, fire; immune to cold) with an attack that does 116 bludgeoning damage at initiative 1
	// 801 units each with 4706 hit points (immune to cold) with an attack that does 116 bludgeoning damage at initiative 1
	// 801 units each with 4706 hit points (weak to bludgeoning) with an attack that does 116 fire damage at initiative 1
	groupDefinition := regexp.MustCompile(`(\d+) units each with (\d+) hit points \(?((?:(?:weak)|(?:immune)) to (?:(?:(?:cold)|(?:fire)|(?:radiation)|(?:slashing)|(?:bludgeoning))[, ]*)+)?[; ]*((?:(?:weak)|(?:immune)) to (?:(?:(?:cold)|(?:fire)|(?:radiation)|(?:slashing)|(?:bludgeoning))[, ]*)+)?[) ]*with an attack that does (\d+) ((?:cold)|(?:fire)|(?:radiation)|(?:slashing)|(?:bludgeoning)) damage at initiative (\d+)`)

	attackTypeDefinition := regexp.MustCompile("((?:cold)|(?:fire)|(?:radiation)|(?:slashing)|(?:bludgeoning))")

	immune := make(army, 0)
	infection := make(army, 0)
	isImmuneUnit := true
	idIter := 1

	for line := ls.Front(); line != nil; line = line.Next() {
		lineContents := line.Value.(string)
		if lineContents == "Immune System:" {
			isImmuneUnit = true
			idIter = 1
			continue
		} else if lineContents == "Infection:" {
			isImmuneUnit = false
			idIter = 1
			continue
		} else if lineContents == "" {
			continue
		}
		//fmt.Printf("Parsing %v\n", lineContents)
		matches := groupDefinition.FindAllStringSubmatch(lineContents, -1)[0]

		groupIter := 1
		unitCount := mustInt(matches[groupIter])
		groupIter++
		hitPointsPerUnit := mustInt(matches[groupIter])
		groupIter++

		immuneTo := make([]attackType, 0)
		weakTo := make([]attackType, 0)

		groupToRead := matches[groupIter]
		if len(groupToRead) != 0 {
			if strings.Contains(groupToRead, "weak to") {
				weakTo = append(weakTo, getActions(attackTypeDefinition, matches[groupIter])...)
			} else if strings.Contains(groupToRead, "immune to") {
				immuneTo = append(immuneTo, getActions(attackTypeDefinition, groupToRead)...)
			}
		}
		groupIter++
		groupToRead = matches[groupIter]
		if len(groupToRead) != 0 {
			if strings.Contains(groupToRead, "weak to") {
				weakTo = append(weakTo, getActions(attackTypeDefinition, matches[groupIter])...)
			} else if strings.Contains(groupToRead, "immune to") {
				immuneTo = append(immuneTo, getActions(attackTypeDefinition, groupToRead)...)
			}
		}
		groupIter++

		attackMight := mustInt(matches[groupIter])
		groupIter++
		attack := attackFromString(matches[groupIter])
		groupIter++
		groupInitiative := mustInt(matches[groupIter])

		g := &group{
			id:               idIter,
			units:            unitCount,
			hitPointsPerUnit: hitPointsPerUnit,
			might:            attackMight,
			attack:           attack,
			initiative:       groupInitiative,
			weakTo:           weakTo,
			immuneTo:         immuneTo,
		}
		if isImmuneUnit {
			immune = append(immune, g)
		} else {
			infection = append(infection, g)
		}
		idIter++
	}
	fmt.Printf("Immune army:\n%v\n", immune)
	fmt.Printf("Infection army:\n%v\n", infection)
}

func getActions(attackTypeDefinition *regexp.Regexp, attackTypesAsCSV string) []attackType {
	result := make([]attackType, 0)
	for _, match := range attackTypeDefinition.FindAllString(attackTypesAsCSV, -1) {
		result = append(result, attackFromString(match))
	}
	return result
}

func attackFromString(attack string) attackType {
	switch attack {
	case "fire":
		return fire
	case "cold":
		return cold
	case "bludgeoning":
		return bludgeoning
	case "radiation":
		return radiation
	case "slashing":
		return slashing
	default:
		panic("Unknown attack: " + attack)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustInt(str string) int {
	result, err := strconv.Atoi(str)
	check(err)
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
