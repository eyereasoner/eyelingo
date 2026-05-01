// barley_seed_lineage.go
//
// Inspired by Eyeling's `examples/act-barley-seed-lineage.n3`.
//
// The example keeps the N3 CAN/CAN'T shape: one viable barley lineage and four
// contrast lineages that fail for a single missing ingredient.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"strings"
)

const eyelingoExampleName = "barley_seed_lineage"

type Dataset struct {
	CaseName string `json:"caseName"`
	Question string `json:"question"`
	World    struct {
		NoDesignLaws     bool     `json:"noDesignLaws"`
		Greenhouse       []string `json:"greenhouse"`
		SelectionFavours string   `json:"selectionFavours"`
	} `json:"world"`
	Lineages []Lineage `json:"lineages"`
	Expected struct {
		Evolvable string   `json:"evolvable"`
		Blocked   []string `json:"blocked"`
	} `json:"expected"`
}

type Lineage struct {
	Name               string `json:"name"`
	DigitalHeredity    bool   `json:"digitalHeredity"`
	Repair             bool   `json:"repair"`
	DormancyProtection bool   `json:"dormancyProtection"`
	HeritableVariation bool   `json:"heritableVariation"`
	Variant            string `json:"variant"`
}

type Result struct {
	Name      string
	Copy      bool
	Accurate  bool
	Closure   bool
	Adaptive  bool
	Evolvable bool
	Blockers  []string
}

type Check struct {
	ID   string
	OK   bool
	Text string
}
type Analysis struct {
	Results []Result
	Checks  []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printReport(ds, analysis)
	if !allOK(analysis.Checks) {
		os.Exit(1)
	}
}

func derive(ds Dataset) Analysis {
	results := []Result{}
	for _, l := range ds.Lineages {
		copyOK := ds.World.NoDesignLaws && l.DigitalHeredity
		accurate := copyOK && l.Repair
		closure := l.DormancyProtection && has(ds.World.Greenhouse, "warmth") && has(ds.World.Greenhouse, "moisture") && has(ds.World.Greenhouse, "light")
		adaptive := accurate && l.HeritableVariation && l.Variant == ds.World.SelectionFavours
		evolvable := closure && adaptive
		blockers := []string{}
		if !l.DigitalHeredity {
			blockers = append(blockers, "non-digital heredity blocks accurate genome copying")
		}
		if !l.Repair {
			blockers = append(blockers, "missing repair blocks reliable damage correction")
		}
		if !l.DormancyProtection {
			blockers = append(blockers, "missing dormancy protection blocks seed-stage closure")
		}
		if !l.HeritableVariation {
			blockers = append(blockers, "missing heritable variation blocks adaptive evolution")
		}
		results = append(results, Result{l.Name, copyOK, accurate, closure, adaptive, evolvable, blockers})
	}
	checks := []Check{
		{"C1", ds.World.NoDesignLaws, "no-design laws are loaded"},
		{"C2", byName(results, "mainLine").Copy, "mainLine can copy its digitally instantiated genome"},
		{"C3", byName(results, "mainLine").Accurate && byName(results, "mainLine").Closure, "mainLine has repair, protected dormancy, and greenhouse support"},
		{"C4", byName(results, "mainLine").Evolvable, "mainLine is evolvable under the saline selection environment"},
		{"C5", containsBlocker(byName(results, "analogLine"), "non-digital"), "analogLine is blocked by non-digital heredity"},
		{"C6", containsBlocker(byName(results, "fragileLine"), "repair"), "fragileLine is blocked by missing repair"},
		{"C7", containsBlocker(byName(results, "coatlessLine"), "dormancy"), "coatlessLine is blocked by missing dormancy protection"},
		{"C8", containsBlocker(byName(results, "staticLine"), "variation"), "staticLine is blocked by missing heritable variation"},
	}
	return Analysis{results, checks}
}

func has(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
func byName(rs []Result, name string) Result {
	for _, r := range rs {
		if r.Name == name {
			return r
		}
	}
	return Result{}
}
func containsBlocker(r Result, needle string) bool {
	return strings.Contains(strings.Join(r.Blockers, "|"), needle)
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
func evolvableNames(rs []Result) []string {
	out := []string{}
	for _, r := range rs {
		if r.Evolvable {
			out = append(out, r.Name)
		}
	}
	return out
}
func blockedNames(rs []Result) []string {
	out := []string{}
	for _, r := range rs {
		if !r.Evolvable {
			out = append(out, r.Name)
		}
	}
	return out
}

func printReport(ds Dataset, analysis Analysis) {
	fmt.Println("# Barley Seed Lineage")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("evolvable lineage : %s\n", strings.Join(evolvableNames(analysis.Results), ", "))
	fmt.Printf("blocked contrast lineages : %s\n", strings.Join(blockedNames(analysis.Results), ", "))
	fmt.Println("mainLine CAN : genome-copy, protected-dormancy, germination, propagule-production, accurate-self-reproduction, lineage-closure, adaptive-persistence")
	fmt.Println("mainLine CAN'T : none of the modeled blockers apply")
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The main lineage satisfies the constructor-theory style CAN side: digital heredity under no-design laws, repair support, a protected dormant seed stage, germination resources, propagule production, and heritable variation.")
	fmt.Println("The contrast lineages are deliberately near misses so the CAN'T side is explicit.")
	fmt.Println("analogLine lacks a digital hereditary medium, fragileLine lacks repair, coatlessLine lacks the protected dormant compartment, and staticLine lacks heritable variation.")
	fmt.Println("Only mainLine closes the life cycle and adaptively persists in the saline selection environment.")
	fmt.Println()
	return
}
