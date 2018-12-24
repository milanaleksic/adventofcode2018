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
	"strings"
)

const (
	input = "day24/input.txt"
	//input = "day24/test.txt"
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
	log.Fatal("Unknown attack type: ", int(a))
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
	designator       string
}

func (g *group) String() string {
	return fmt.Sprintf("Group %v with %v units with %v hit points (weak to: %+v, immune to: %+v), with might %v and initiative %v; effective=%v", g.id, g.units, g.hitPointsPerUnit, g.weakTo, g.immuneTo, g.might, g.initiative, g.effectivePower())
}

func (g *group) effectivePower() int {
	return g.units * g.might
}

func (g *group) inBattle() bool {
	return g.units > 0
}

func (g *group) isImmuneTo(a attackType) bool {
	for _, immunity := range g.immuneTo {
		if immunity == a {
			return true
		}
	}
	return false
}

func (g *group) isWeakTo(a attackType) bool {
	for _, weakness := range g.weakTo {
		if weakness == a {
			return true
		}
	}
	return false
}

type army []*group

// 801 units each with 4706 hit points (weak to bludgeoning, fire; immune to cold) with an attack that does 116 bludgeoning damage at initiative 1
// 801 units each with 4706 hit points (immune to cold) with an attack that does 116 bludgeoning damage at initiative 1
// 801 units each with 4706 hit points (weak to bludgeoning) with an attack that does 116 fire damage at initiative 1
var groupDefinition = regexp.MustCompile(`(\d+) units each with (\d+) hit points \(?((?:(?:weak)|(?:immune)) to (?:(?:(?:cold)|(?:fire)|(?:radiation)|(?:slashing)|(?:bludgeoning))[, ]*)+)?[; ]*((?:(?:weak)|(?:immune)) to (?:(?:(?:cold)|(?:fire)|(?:radiation)|(?:slashing)|(?:bludgeoning))[, ]*)+)?[) ]*with an attack that does (\d+) ((?:cold)|(?:fire)|(?:radiation)|(?:slashing)|(?:bludgeoning)) damage at initiative (\d+)`)
var attackTypeDefinition = regexp.MustCompile("((?:cold)|(?:fire)|(?:radiation)|(?:slashing)|(?:bludgeoning))")

func main() {
	file, err := os.Open(input)
	check(err)
	defer file.Close()

	ls := list.New()
	readAll(file, ls)

	immune, infection := readInput(ls)
	fmt.Printf("Immune army:\n%v\n", immune)
	fmt.Printf("Infection army:\n%v\n", infection)
	solution, winner := part1(immune, infection)
	fmt.Printf("Solution to part 1 is: %v (winner is %v)\n", solution, winner)
	boost, solution := part2(ls)
	fmt.Printf("Solution to part 2 is: %v (with boost %v)\n", solution, boost)
}

func readInput(ls *list.List) (immune army, infection army) {
	immune = make(army, 0)
	infection = make(army, 0)
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

		designator := "infection"
		if isImmuneUnit {
			designator = "immune"
		}
		g := &group{
			id:               idIter,
			units:            unitCount,
			hitPointsPerUnit: hitPointsPerUnit,
			might:            attackMight,
			attack:           attack,
			initiative:       groupInitiative,
			weakTo:           weakTo,
			immuneTo:         immuneTo,
			designator:       designator,
		}
		if isImmuneUnit {
			immune = append(immune, g)
		} else {
			infection = append(infection, g)
		}
		idIter++
	}
	return
}

func part2(input *list.List) (boost int, solution int) {
	maxBoost := 1571
	for boost := 0; boost < maxBoost; boost++ {
		immune, infection := readInput(input)
		for _, i := range immune {
			i.might += boost
		}
		solution, winner := part1(immune, infection)
		if winner == "immune" {
			return boost, solution
		}
	}
	return 0, 0
}

func comparatorGroupsOrderOfSelection(target army) func(i, j int) bool {
	return func(i, j int) bool {
		powerI := target[i].effectivePower()
		powerJ := target[j].effectivePower()
		if powerI > powerJ {
			return true
		} else if powerI < powerJ {
			return false
		}
		return target[i].initiative > target[j].initiative
	}
}

func part1(immunity army, infection army) (solution int, winner string) {
	battleground := make([]*group, 0)
	comparatorGroupsOrderOfAttack := func(i, j int) bool {
		return battleground[i].initiative > battleground[j].initiative
	}
	maxIter := 10000
	for iter := 1; iter <= maxIter; iter++ {
		//fmt.Println("################### ITER", iter)
		attackedByInfection := make(map[int]*group)
		sort.Slice(infection, comparatorGroupsOrderOfSelection(infection))
		for _, attacker := range infection {
			findOptimalTarget(attacker, immunity, attackedByInfection, "Infection")
		}
		attackedByImmune := make(map[int]*group)
		sort.Slice(immunity, comparatorGroupsOrderOfSelection(immunity))
		for _, attacker := range immunity {
			findOptimalTarget(attacker, infection, attackedByImmune, "Immune")
		}

		battleground = make([]*group, 0)
		for _, group := range immunity {
			if group.inBattle() {
				battleground = append(battleground, group)
			}
		}
		sizeBG := len(battleground)
		if sizeBG == 0 {
			//fmt.Println("Battle finished, no immunity groups any more")
			return sumUnits(infection), "infection"
		}
		for _, group := range infection {
			if group.inBattle() {
				battleground = append(battleground, group)
			}
		}
		if len(battleground) == sizeBG {
			//fmt.Println("Battle finished, no infection groups any more")
			return sumUnits(immunity), "immune"
		}
		sort.Slice(battleground, comparatorGroupsOrderOfAttack)
		for _, attacker := range battleground {
			var target *group
			switch attacker.designator {
			case "infection":
				target = findTarget(attackedByInfection, immunity, attacker)
			case "immune":
				target = findTarget(attackedByImmune, infection, attacker)
			}
			if target == nil || !target.inBattle() {
				continue
			}
			//fmt.Printf("Attack by %v group %v", attacker.designator, attacker.id)
			attackMight := calculateAttack(attacker, target)
			unitsToKill := attackMight / target.hitPointsPerUnit
			if unitsToKill > target.units {
				unitsToKill = target.units
			}
			target.units -= unitsToKill
			//fmt.Printf(" removing %v units down to %v\n", unitsToKill, target.units)
		}
	}
	return -1, ""
}

func sumUnits(groups army) int {
	sum := 0
	for _, g := range groups {
		sum += g.units
	}
	return sum
}

func findTarget(planning map[int]*group, targetArmy army, attacker *group) *group {
	for targetId, g := range planning {
		if g.id == attacker.id {
			return find(targetArmy, targetId)
		}
	}
	//log.Fatalf("Could not find designated target!")
	return nil
}

func findOptimalTarget(attacker *group, targetArmy army, planning map[int]*group, designator string) {
	if !attacker.inBattle() {
		return
	}
	//fmt.Printf("%v group %v has effective power %v (across %v units) and initiative %v", designator, attacker.id, attacker.effectivePower(), attacker.units, attacker.initiative)
	maxAttack := 0
	maxAttackTargetId := -1
	for _, groupTargetCandidate := range targetArmy {
		if !groupTargetCandidate.inBattle() {
			continue
		}
		if _, ok := planning[groupTargetCandidate.id]; ok {
			// already marked by higher prio attacker
			continue
		}
		attackMight := calculateAttack(attacker, groupTargetCandidate)
		if attackMight > maxAttack {
			maxAttack = attackMight
			maxAttackTargetId = groupTargetCandidate.id
		} else if attackMight == maxAttack && attackMight != 0 {
			candidate := find(targetArmy, maxAttackTargetId)
			if groupTargetCandidate.effectivePower() > candidate.effectivePower() {
				maxAttackTargetId = groupTargetCandidate.id
			} else if groupTargetCandidate.effectivePower() == candidate.effectivePower() {
				if groupTargetCandidate.initiative > candidate.initiative {
					maxAttackTargetId = groupTargetCandidate.id
				}
			}
		}
	}
	if maxAttack != 0 {
		planning[maxAttackTargetId] = attacker
		//fmt.Printf("; it will attack enemy group %v with might %v\n", maxAttackTargetId, maxAttack)
	} else {
		//fmt.Printf("; it will not attack\n")
	}
}

func find(groups army, id int) *group {
	for _, group := range groups {
		if group.id == id {
			return group
		}
	}
	log.Fatalf("Could not find group with id %v in army %v", id, groups)
	return nil
}

func calculateAttack(attacker *group, targetCandidate *group) int {
	if targetCandidate.isImmuneTo(attacker.attack) {
		return 0
	}
	if targetCandidate.isWeakTo(attacker.attack) {
		return attacker.effectivePower() * 2
	}
	return attacker.effectivePower()
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
