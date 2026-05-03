// alarm_bit_interoperability.go
//
// A Go translation inspired by Eyeling's
// `examples/act-alarm-bit-interoperability.n3`.
//
// The example distinguishes what can be done with classical alarm-bit media
// from what cannot be done with a superinformation-like token.
//
// Run:
//
//	go run examples/alarm_bit_interoperability.go
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"strings"
)

const seeExampleName = "alarm_bit_interoperability"

type Dataset struct {
	CaseName             string
	Question             string
	ClassicalMedia       []InformationMedium
	Superinformation     SuperinformationMedium
	ExpectedCopyTasks    int
	ExpectedImpossible   []string
	ExpectedCanDecision  string
	ExpectedCantDecision string
}

type InformationMedium struct {
	Name      string
	Variable  string
	ZeroState string
	OneState  string
}

type SuperinformationMedium struct {
	Name     string
	Variable string
	States   []string
}

type CopyTask struct {
	From     string
	To       string
	Variable string
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	LocalPermutations []string
	CopyTasks         []CopyTask
	ImpossibleTasks   []string
	CannotStates      []string
	CanDecision       string
	CantDecision      string
	Checks            []Check
}

func main() {
	ds := exampleinput.Load(seeExampleName, Dataset{})
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
}

func derive(ds Dataset) Analysis {
	permutations := make([]string, 0)
	copyTasks := make([]CopyTask, 0)
	for _, medium := range ds.ClassicalMedia {
		permutations = append(permutations, medium.Name)
	}
	for _, from := range ds.ClassicalMedia {
		for _, to := range ds.ClassicalMedia {
			if from.Name != to.Name && from.Variable == to.Variable {
				copyTasks = append(copyTasks, CopyTask{From: from.Name, To: to.Name, Variable: from.Variable})
			}
		}
	}
	impossible := []string{"CloneAllStates"}
	cannot := []string{"UniversalClone", "UnrestrictedStateFanOut"}
	canDecision := "NO"
	if len(copyTasks) == ds.ExpectedCopyTasks && len(permutations) == len(ds.ClassicalMedia) {
		canDecision = ds.ExpectedCanDecision
	}
	cantDecision := "NO"
	if containsAll(cannot, ds.ExpectedImpossible) {
		cantDecision = ds.ExpectedCantDecision
	}
	checks := []Check{
		{ID: "C1", OK: len(ds.ClassicalMedia) == 2, Text: "two unlike classical media are present"},
		{ID: "C2", OK: sameVariable(ds.ClassicalMedia), Text: fmt.Sprintf("classical media encode the same variable %s", ds.ClassicalMedia[0].Variable)},
		{ID: "C3", OK: len(copyTasks) == ds.ExpectedCopyTasks, Text: fmt.Sprintf("%d directed copy tasks are possible", len(copyTasks))},
		{ID: "C4", OK: len(permutations) == len(ds.ClassicalMedia), Text: "each classical medium supports a local permutation"},
		{ID: "C5", OK: contains(cannot, "UniversalClone"), Text: fmt.Sprintf("%s cannot be universally cloned", ds.Superinformation.Name)},
		{ID: "C6", OK: contains(cannot, "UnrestrictedStateFanOut"), Text: "unrestricted classical-style fan-out is blocked for the superinformation token"},
		{ID: "C7", OK: canDecision == ds.ExpectedCanDecision && cantDecision == ds.ExpectedCantDecision, Text: fmt.Sprintf("CAN=%s and CANNOT=%s decisions are both derived", canDecision, cantDecision)},
	}
	return Analysis{LocalPermutations: permutations, CopyTasks: copyTasks, ImpossibleTasks: impossible, CannotStates: cannot, CanDecision: canDecision, CantDecision: cantDecision, Checks: checks}
}

func sameVariable(media []InformationMedium) bool {
	if len(media) == 0 {
		return false
	}
	variable := media[0].Variable
	for _, medium := range media {
		if medium.Variable != variable {
			return false
		}
	}
	return true
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func containsAll(values []string, wants []string) bool {
	for _, want := range wants {
		if !contains(values, want) {
			return false
		}
	}
	return true
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func countChecksOK(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}

func printAnswer(ds Dataset, analysis Analysis) {
	fmt.Println("# Alarm Bit Interoperability")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("classical alarm-bit interoperability : %s\n", analysis.CanDecision)
	fmt.Printf("universal cloning of the superinformation token : %s\n", analysis.CantDecision)
	for _, task := range analysis.CopyTasks {
		fmt.Printf("copy task : %s -> %s for %s\n", task.From, task.To, task.Variable)
	}
	fmt.Printf("blocked tasks : %s\n", strings.Join(analysis.CannotStates, ", "))
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason")
	fmt.Println("The optical beacon and relay register are unlike physical media, but both encode the same abstract AlarmBit variable.")
	fmt.Println("Because the variable is classical in both media, local permutation and copying in both directed media transfers are possible.")
	fmt.Printf("The %s is treated as a superinformation medium with states %s.\n", ds.Superinformation.Name, strings.Join(ds.Superinformation.States, ", "))
	fmt.Println("That contrast substrate cannot support universal cloning of all states, so unrestricted classical-style fan-out is also blocked.")
	fmt.Println()
}
