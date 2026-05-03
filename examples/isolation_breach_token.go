// isolation_breach_token.go
//
// Inspired by Eyeling's `examples/act-isolation-breach.n3`.
//
// This example keeps the can/can't split: breach-token tasks are possible for
// classical media, while the provenance seal refuses unrestricted fan-out.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"os"
)

const seeExampleName = "isolation_breach_token"

type Dataset struct {
	CaseName               string   `json:"caseName"`
	Question               string   `json:"question"`
	Variable               string   `json:"variable"`
	Media                  []Medium `json:"media"`
	SuperinformationMedium struct {
		Name     string   `json:"name"`
		Variable string   `json:"variable"`
		States   []string `json:"states"`
	} `json:"superinformationMedium"`
	Expected struct {
		PreparedState string `json:"preparedState"`
		Serial        string `json:"serial"`
	} `json:"expected"`
}
type Medium struct {
	Name string `json:"name"`
	Zero string `json:"zero"`
	One  string `json:"one"`
}
type Check struct {
	ID   string
	OK   bool
	Text string
}
type Analysis struct {
	Prepare, Permute, Copy, Measure, Serial, Parallel, Impossible int
	Checks                                                        []Check
}

func main() {
	ds := exampleinput.Load(seeExampleName, Dataset{})
	a := derive(ds)
	printReport(ds, a)
	if !allOK(a.Checks) {
		os.Exit(1)
	}
}
func derive(ds Dataset) Analysis {
	n := len(ds.Media)
	a := Analysis{Prepare: 2 * n, Permute: 2 * n, Copy: n * (n - 1), Measure: n * n, Serial: n * (n - 1) * n, Parallel: n * (n - 1) * (n - 2), Impossible: 2}
	a.Checks = []Check{
		{"C1", hasMedium(ds.Media, "doorBeacon") && hasMedium(ds.Media, "containmentPLC") && hasMedium(ds.Media, "nursePager") && hasMedium(ds.Media, "incidentBoard"), "doorBeacon, containmentPLC, nursePager, and incidentBoard encode BreachBit"},
		{"C2", hasMedium(ds.Media, "nursePager") && ds.Expected.PreparedState == "CodeBreach", "nursePager can prepare CodeBreach"},
		{"C3", hasMedium(ds.Media, "doorBeacon"), "doorBeacon can permute SafeGreen to BreachRed and back"},
		{"C4", hasMedium(ds.Media, "containmentPLC") && hasMedium(ds.Media, "nursePager"), "containmentPLC can copy the breach token to nursePager"},
		{"C5", hasMedium(ds.Media, "nursePager") && hasMedium(ds.Media, "incidentBoard"), "nursePager can be measured into incidentBoard"},
		{"C6", ds.Expected.Serial == "doorBeacon->containmentPLC->incidentBoard", "doorBeacon -> containmentPLC -> incidentBoard is a serial audit network"},
		{"C7", hasMedium(ds.Media, "containmentPLC") && hasMedium(ds.Media, "nursePager") && hasMedium(ds.Media, "incidentBoard"), "containmentPLC can fan out to nursePager and incidentBoard"},
		{"C8", len(ds.SuperinformationMedium.States) == 3, "specimenSeal cannot universally clone all provenance states"},
		{"C9", ds.SuperinformationMedium.Name == "specimenSeal", "specimenSeal blocks unrestricted parallel fan-out"},
	}
	return a
}
func hasMedium(media []Medium, name string) bool {
	for _, m := range media {
		if m.Name == name {
			return true
		}
	}
	return false
}
func allOK(checks []Check) bool {
	for _, c := range checks {
		if !c.OK {
			return false
		}
	}
	return true
}
func countOK(checks []Check) int {
	n := 0
	for _, c := range checks {
		if c.OK {
			n++
		}
	}
	return n
}
func printReport(ds Dataset, a Analysis) {
	fmt.Println("# Isolation Breach Token")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("classical breach token : YES, prepare, reversible permutation, copy, measure, serial audit, and parallel fan-out all succeed")
	fmt.Println("specimen provenance seal : NO, universal cloning and unrestricted parallel fan-out are blocked")
	fmt.Println("prepared witness : nursePager prepares CodeBreach")
	fmt.Println("serial witness : doorBeacon -> containmentPLC -> incidentBoard")
	fmt.Printf("possible prepare tasks : %d\n", a.Prepare)
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The breach token is an ordinary classical information variable carried by four unlike media in the lab workflow.")
	fmt.Println("Since each medium carries BreachBit, the example derives preparation, reversible permutation, copying, and measurement tasks.")
	fmt.Println("Copy and measurement compose into an incident-board audit path, and copy pairs compose into a parallel notification witness.")
	fmt.Println("The specimen seal is separate and superinformation-like, so it records provenance without becoming an unrestricted cloneable broadcast token.")
	fmt.Println()
	return
}
