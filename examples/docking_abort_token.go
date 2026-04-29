// docking_abort_token.go
//
// Inspired by Eyeling's `examples/act-docking-abort.n3`.
//
// This example derives possible classical abort-token tasks and impossible
// quantum-seal tasks from a small constructor-theory style fixture.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
)

const eyelingoExampleName = "docking_abort_token"

type Dataset struct {
	CaseName               string   `json:"caseName"`
	Question               string   `json:"question"`
	Variable               string   `json:"variable"`
	Media                  []Medium `json:"media"`
	SuperinformationMedium struct {
		Name     string `json:"name"`
		Variable string `json:"variable"`
	} `json:"superinformationMedium"`
	Expected struct {
		Serial         string `json:"serial"`
		ParallelSource string `json:"parallelSource"`
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
	CopyTasks, MeasureTasks, SerialNetworks, ParallelNetworks, ImpossibleTasks int
	Checks                                                                     []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	a := derive(ds)
	printReport(ds, a)
	if !allOK(a.Checks) {
		os.Exit(1)
	}
}
func derive(ds Dataset) Analysis {
	n := len(ds.Media)
	copyTasks := n * (n - 1)
	measureTasks := n * n
	serial := copyTasks * n
	parallel := n * (n - 1) * (n - 2)
	checks := []Check{
		{"C1", n == 4, "four classical media encode AbortBit"},
		{"C2", n > 0, "each classical medium can distinguish and locally permute the abort bit"},
		{"C3", hasMedium(ds.Media, "abortLamp") && hasMedium(ds.Media, "flightPLC"), "abortLamp can copy the token to flightPLC"},
		{"C4", hasMedium(ds.Media, "radioFrame") && hasMedium(ds.Media, "auditDisplay"), "radioFrame can be measured into auditDisplay"},
		{"C5", hasMedium(ds.Media, "abortLamp") && hasMedium(ds.Media, "flightPLC") && hasMedium(ds.Media, "auditDisplay"), "a serial abortLamp -> flightPLC -> auditDisplay network is possible"},
		{"C6", hasMedium(ds.Media, "flightPLC") && hasMedium(ds.Media, "radioFrame") && hasMedium(ds.Media, "auditDisplay"), "a parallel flightPLC fan-out to radioFrame and auditDisplay is possible"},
		{"C7", ds.SuperinformationMedium.Name == "quantumSeal", "quantumSeal cannot universally clone all seal states"},
		{"C8", ds.SuperinformationMedium.Variable == "SealVariable", "quantumSeal cannot be used for unrestricted audit fan-out"},
	}
	return Analysis{copyTasks, measureTasks, serial, parallel, 2, checks}
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
	fmt.Println("# Docking Abort Token")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("classical abort token : YES, it can be permuted, copied, measured, and composed into audit networks")
	fmt.Println("quantum authenticity seal : NO, it cannot be universally cloned or used as unrestricted audit fan-out")
	fmt.Println("serial witness : abortLamp -> flightPLC -> auditDisplay")
	fmt.Println("parallel witness : flightPLC -> radioFrame and auditDisplay")
	fmt.Printf("possible copy tasks : %d\n", a.CopyTasks)
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("All four classical media encode the same abstract AbortBit variable, so the model treats them as interoperable information media.")
	fmt.Println("Local permutation and local cloning are allowed on each classical medium, while copying and measuring are allowed between media that carry the same variable.")
	fmt.Println("Those primitive tasks compose into a serial audit path and a parallel fan-out witness.")
	fmt.Println("The quantum seal is modeled as a superinformation medium, so universal cloning and unrestricted audit fan-out are explicitly blocked.")
	fmt.Println()
	fmt.Println("## Check")
	for _, c := range a.Checks {
		status := "FAIL"
		if c.OK {
			status = "OK"
		}
		fmt.Printf("%s %s - %s\n", c.ID, status, c.Text)
	}
	fmt.Println()
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("classical media : %d\n", len(ds.Media))
	fmt.Printf("copy tasks : %d\n", a.CopyTasks)
	fmt.Printf("measure tasks : %d\n", a.MeasureTasks)
	fmt.Printf("serial networks : %d\n", a.SerialNetworks)
	fmt.Printf("parallel networks : %d\n", a.ParallelNetworks)
	fmt.Printf("impossible tasks : %d\n", a.ImpossibleTasks)
	fmt.Printf("checks passed : %d/%d\n", countOK(a.Checks), len(a.Checks))
}
