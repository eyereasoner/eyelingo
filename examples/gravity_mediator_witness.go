// gravity_mediator_witness.go
//
// Inspired by Eyeling's `examples/act-gravity-mediator-witness.n3`.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
)

const eyelingoExampleName = "gravity_mediator_witness"

type Run struct {
	Name                 string   `json:"name"`
	Role                 string   `json:"role"`
	Mediator             string   `json:"mediator"`
	MediatorModel        string   `json:"mediatorModel"`
	CouplingMode         string   `json:"couplingMode"`
	Assumes              []string `json:"assumes"`
	DirectCouplingStatus string   `json:"directCouplingStatus"`
	Observed             string   `json:"observed"`
	ProbeStatus          string   `json:"probeStatus"`
	ControlStatus        string   `json:"controlStatus"`
}

type Dataset struct {
	CaseName string `json:"caseName"`
	Source   string `json:"source"`
	Question string `json:"question"`
	Runs     []Run  `json:"runs"`
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Result struct {
	Positive           Run
	Contrast           Run
	PositiveConclusion bool
	ContrastBlock      bool
	Checks             []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	r := derive(ds)
	printReport(ds, r)
	if !allOK(r.Checks) {
		os.Exit(1)
	}
}

func derive(ds Dataset) Result {
	var r Result
	for _, run := range ds.Runs {
		if run.Role == "positive" {
			r.Positive = run
		} else {
			r.Contrast = run
		}
	}

	posMediatorOnly := has(r.Positive.Assumes, "Locality") && r.Positive.DirectCouplingStatus == "NoDirectCoupling"
	posInterfaces := has(r.Positive.Assumes, "Interoperability") &&
		r.Positive.ControlStatus == "CopyLikeControlPresent" &&
		r.Positive.ProbeStatus == "LocalProbeReadoutPresent"
	r.PositiveConclusion = posMediatorOnly && posInterfaces &&
		r.Positive.Observed == "EntanglementWitnessPassed" &&
		r.Positive.CouplingMode == "Gravitational"

	conMediatorOnly := has(r.Contrast.Assumes, "Locality") && r.Contrast.DirectCouplingStatus == "NoDirectCoupling"
	r.ContrastBlock = conMediatorOnly && has(r.Contrast.Assumes, "Interoperability") &&
		r.Contrast.MediatorModel == "PurelyClassical"

	r.Checks = []Check{
		{"C1", has(r.Positive.Assumes, "Locality"), "locality is assumed in the positive run"},
		{"C2", has(r.Positive.Assumes, "Interoperability"), "interoperability is assumed in the positive run"},
		{"C3", r.Positive.DirectCouplingStatus == "NoDirectCoupling", "direct coupling between the two quantum systems is excluded"},
		{"C4", posMediatorOnly, "the positive run has a mediator-only interaction path"},
		{"C5", r.Positive.Observed == "EntanglementWitnessPassed", "an entanglement witness is observed in the positive run"},
		{"C6", posInterfaces, "the positive run has both information-transfer and local-readout interfaces"},
		{"C7", r.PositiveConclusion, "the gravitational mediator is derived to be non-classical"},
		{"C8", r.PositiveConclusion, "a purely classical mediator model is ruled out by the positive run"},
		{"C9", conMediatorOnly, "the contrast run is also mediator-only"},
		{"C10", r.ContrastBlock, "the purely classical contrast mediator cannot support the witness"},
	}
	return r
}

func printReport(ds Dataset, r Result) {
	fmt.Println("# Gravity Mediator Witness")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("YES for the mediator-only witness run.")
	fmt.Println("NO for a purely classical mediator model under the same mediator-only conditions.")
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The positive run assumes locality and interoperability, excludes direct coupling, and observes entanglement after interaction through the gravitational mediator alone.")
	fmt.Println("Under those conditions the mediator-only witness supports a non-classical-mediator conclusion, while the purely classical contrast model cannot support the same witness.")
	fmt.Println()
	return
}

func has(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func allOK(cs []Check) bool {
	for _, c := range cs {
		if !c.OK {
			return false
		}
	}
	return true
}

func countOK(cs []Check) int {
	n := 0
	for _, c := range cs {
		if c.OK {
			n++
		}
	}
	return n
}
