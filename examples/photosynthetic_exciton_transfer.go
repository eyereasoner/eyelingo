// photosynthetic_exciton_transfer.go
//
// Inspired by Eyeling's `examples/act-photosynthetic-exciton-transfer.n3`.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
)

const eyelingoExampleName = "photosynthetic_exciton_transfer"

type Complex struct {
	Name                        string `json:"name"`
	Role                        string `json:"role"`
	ExcitonCoupling             string `json:"excitonCoupling"`
	StateDelocalization         string `json:"stateDelocalization"`
	VibronicBridge              string `json:"vibronicBridge"`
	EnergyLandscape             string `json:"energyLandscape"`
	ElectronicCoherenceLifetime string `json:"electronicCoherenceLifetime"`
	Dephasing                   string `json:"dephasing"`
	ConnectedTo                 string `json:"connectedTo"`
}

type Dataset struct {
	CaseName       string    `json:"caseName"`
	Source         string    `json:"source"`
	Question       string    `json:"question"`
	ReactionCenter string    `json:"reactionCenter"`
	Complexes      []Complex `json:"complexes"`
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Result struct {
	Tuned   Complex
	Detuned Complex
	Can     map[string][]string
	Cannot  map[string][]string
	Checks  []Check
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
	r := Result{Can: map[string][]string{}, Cannot: map[string][]string{}}

	for _, c := range ds.Complexes {
		if c.Role == "positive" {
			r.Tuned = c
		} else {
			r.Detuned = c
		}

		if c.ExcitonCoupling == "Strong" && c.StateDelocalization == "Present" {
			r.Can[c.Name] = append(r.Can[c.Name], "CoherentPathwaySampling")
		}
		if c.VibronicBridge == "Tuned" && c.Dephasing == "Moderate" {
			r.Can[c.Name] = append(r.Can[c.Name], "VibronicallyAssistedTransfer")
		}
		if c.ElectronicCoherenceLifetime == "Short" && c.EnergyLandscape == "DownhillToReactionCenter" {
			r.Can[c.Name] = append(r.Can[c.Name], "ShortLivedQuantumAssistance")
		}
		if has(r.Can[c.Name], "CoherentPathwaySampling") &&
			has(r.Can[c.Name], "VibronicallyAssistedTransfer") &&
			has(r.Can[c.Name], "ShortLivedQuantumAssistance") {
			r.Can[c.Name] = append(r.Can[c.Name], "EfficientExcitonTransfer")
		}
		if has(r.Can[c.Name], "EfficientExcitonTransfer") && c.ConnectedTo == ds.ReactionCenter {
			r.Can[c.Name] = append(r.Can[c.Name], "DeliverExcitationToReactionCenter")
		}

		if c.ExcitonCoupling == "Weak" && c.StateDelocalization == "Absent" {
			r.Cannot[c.Name] = append(r.Cannot[c.Name], "CoherentPathwaySampling")
		}
		if c.VibronicBridge == "Absent" && c.Dephasing == "Strong" {
			r.Cannot[c.Name] = append(r.Cannot[c.Name], "VibronicallyAssistedTransfer")
		}
		if c.EnergyLandscape == "TrappingMismatch" {
			r.Cannot[c.Name] = append(r.Cannot[c.Name], "DirectedReactionCenterTransfer")
		}
		if has(r.Cannot[c.Name], "CoherentPathwaySampling") &&
			has(r.Cannot[c.Name], "VibronicallyAssistedTransfer") {
			r.Cannot[c.Name] = append(r.Cannot[c.Name], "EfficientExcitonTransfer")
		}
		if has(r.Cannot[c.Name], "DirectedReactionCenterTransfer") ||
			has(r.Cannot[c.Name], "EfficientExcitonTransfer") {
			r.Cannot[c.Name] = append(r.Cannot[c.Name], "DeliverExcitationToReactionCenter")
		}
	}

	r.Checks = []Check{
		{"C1", has(r.Can[r.Tuned.Name], "CoherentPathwaySampling"), "the tuned complex can sample exciton pathways coherently"},
		{"C2", has(r.Can[r.Tuned.Name], "VibronicallyAssistedTransfer"), "the tuned complex can use vibronically assisted transfer"},
		{"C3", has(r.Can[r.Tuned.Name], "ShortLivedQuantumAssistance"), "short-lived quantum assistance is enough in the tuned downhill regime"},
		{"C4", has(r.Can[r.Tuned.Name], "EfficientExcitonTransfer"), "efficient exciton transfer is possible in the tuned complex"},
		{"C5", has(r.Can[r.Tuned.Name], "DeliverExcitationToReactionCenter"), "the tuned complex can deliver excitation to the reaction center"},
		{"C6", has(r.Cannot[r.Detuned.Name], "CoherentPathwaySampling"), "the detuned complex cannot sample pathways coherently"},
		{"C7", has(r.Cannot[r.Detuned.Name], "VibronicallyAssistedTransfer"), "the detuned complex cannot use vibronically assisted transfer"},
		{"C8", has(r.Cannot[r.Detuned.Name], "DirectedReactionCenterTransfer"), "the detuned complex cannot achieve directed reaction-center transfer"},
		{"C9", has(r.Cannot[r.Detuned.Name], "EfficientExcitonTransfer"), "the detuned complex cannot achieve efficient exciton transfer"},
		{"C10", has(r.Cannot[r.Detuned.Name], "DeliverExcitationToReactionCenter"), "the detuned complex cannot deliver excitation efficiently to the reaction center"},
	}
	return r
}

func printReport(ds Dataset, r Result) {
	fmt.Println("# Photosynthetic Exciton Transfer")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("YES for the tuned antenna complex.")
	fmt.Println("NO for the detuned, strongly decohered contrast complex.")
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The tuned complex combines strong excitonic coupling, delocalization, a tuned vibronic bridge, moderate dephasing, short-lived coherence, and a downhill route to the reaction center.")
	fmt.Println("The detuned contrast complex has weak coupling, absent delocalization, no vibronic bridge, strong dephasing, and a trapping mismatch, so the same efficient delivery task is blocked.")
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
