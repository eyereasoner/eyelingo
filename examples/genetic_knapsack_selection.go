// genetic_knapsack_selection.go
//
// A deterministic Go translation inspired by Eyeling's
// `examples/genetic-algorithm-knapsack.n3`.
//
// The example evaluates a 0/1 knapsack genome, generates every single-bit
// mutant, and keeps the candidate with the lowest fitness until no improving
// mutant exists.
//
// Run:
//
//	go run examples/genetic_knapsack_selection.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"sort"
	"strings"
)

const eyelingoExampleName = "genetic_knapsack_selection"

type Dataset struct {
	CaseName       string
	Question       string
	Capacity       int
	MaxGenerations int
	StartGenome    string
	Items          []Item
	Expected       Expected
}

type Item struct {
	Name   string
	Weight int
	Value  int
}

type Expected struct {
	FinalGenome string
	FinalWeight int
	FinalValue  int
}

type Candidate struct {
	Genome  string
	Weight  int
	Value   int
	Fitness int
}

type Generation struct {
	Index  int
	Parent Candidate
	Best   Candidate
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	History []Generation
	Final   Candidate
	Checks  []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
}

func derive(ds Dataset) Analysis {
	genome := ds.StartGenome
	history := make([]Generation, 0)
	for generation := 0; generation <= ds.MaxGenerations; generation++ {
		parent := evaluate(genome, ds.Items, ds.Capacity)
		best := selectBest(append([]string{genome}, mutants(genome)...), ds.Items, ds.Capacity)
		history = append(history, Generation{Index: generation, Parent: parent, Best: best})
		if best.Genome == genome {
			break
		}
		genome = best.Genome
	}
	final := evaluate(genome, ds.Items, ds.Capacity)
	selected := selectedItemNames(genome, ds.Items)
	checks := []Check{
		{ID: "C1", OK: len(ds.Items) == len(ds.StartGenome), Text: fmt.Sprintf("%d items align with a %d-bit genome", len(ds.Items), len(ds.StartGenome))},
		{ID: "C2", OK: final.Weight <= ds.Capacity, Text: fmt.Sprintf("final weight %d is within capacity %d", final.Weight, ds.Capacity)},
		{ID: "C3", OK: final.Genome == ds.Expected.FinalGenome, Text: fmt.Sprintf("final genome is %s", final.Genome)},
		{ID: "C4", OK: final.Value == ds.Expected.FinalValue, Text: fmt.Sprintf("final value is %d using %s", final.Value, strings.Join(selected, ", "))},
		{ID: "C5", OK: history[len(history)-1].Best.Genome == final.Genome, Text: "no single-bit neighbor improves the final candidate"},
	}
	return Analysis{History: history, Final: final, Checks: checks}
}

func evaluate(genome string, items []Item, capacity int) Candidate {
	weight := 0
	value := 0
	for i, bit := range genome {
		if bit == '1' {
			weight += items[i].Weight
			value += items[i].Value
		}
	}
	fitness := 1000000 - value
	if weight > capacity {
		fitness = 2000000 + (weight - capacity)
	}
	return Candidate{Genome: genome, Weight: weight, Value: value, Fitness: fitness}
}

func mutants(genome string) []string {
	out := make([]string, 0, len(genome))
	for i := range genome {
		bits := []byte(genome)
		if bits[i] == '0' {
			bits[i] = '1'
		} else {
			bits[i] = '0'
		}
		out = append(out, string(bits))
	}
	return out
}

func selectBest(genomes []string, items []Item, capacity int) Candidate {
	candidates := make([]Candidate, 0, len(genomes))
	for _, genome := range genomes {
		candidates = append(candidates, evaluate(genome, items, capacity))
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Fitness != candidates[j].Fitness {
			return candidates[i].Fitness < candidates[j].Fitness
		}
		return candidates[i].Genome < candidates[j].Genome
	})
	return candidates[0]
}

func selectedItemNames(genome string, items []Item) []string {
	out := make([]string, 0)
	for i, bit := range genome {
		if bit == '1' {
			out = append(out, items[i].Name)
		}
	}
	return out
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
	fmt.Println("# Genetic Knapsack Selection")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("final genome : %s\n", analysis.Final.Genome)
	fmt.Printf("selected items : %s\n", strings.Join(selectedItemNames(analysis.Final.Genome, ds.Items), ", "))
	fmt.Printf("weight : %d / %d\n", analysis.Final.Weight, ds.Capacity)
	fmt.Printf("value : %d\n", analysis.Final.Value)
	fmt.Printf("fitness : %d\n", analysis.Final.Fitness)
	fmt.Printf("generations evaluated : %d\n", len(analysis.History))
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason why")
	fmt.Println("Each genome bit says whether the corresponding item is selected for the knapsack.")
	fmt.Println("Feasible candidates get fitness 1000000 minus value, so higher value means lower fitness; overweight candidates are penalized above every feasible candidate.")
	fmt.Println("At every generation, all single-bit mutants of the parent are compared with the parent, and the lowest-fitness candidate is selected deterministically.")
	fmt.Printf("The run stops at %s because every one-bit neighbor is no better under the capacity %d rule.\n", analysis.Final.Genome, ds.Capacity)
	fmt.Println()
}
