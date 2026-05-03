package main

import (
	"fmt"
	"strings"
)

func checkGenetic(ctx *Context) []Check {
	d := ctx.M()
	items := maps(d["Items"])
	cap := integer(d["Capacity"])
	genome := str(d["StartGenome"])
	history := 0
	var final map[string]any
	for gen := 0; gen <= integer(d["MaxGenerations"]); gen++ {
		history++
		cand := append([]string{genome}, mutants(genome)...)
		best := knapBest(cand, items, cap)
		if str(best["genome"]) == genome {
			final = knapEval(genome, items, cap)
			break
		}
		genome = str(best["genome"])
	}
	parsedGenome := ""
	if m := reFind(ctx.Answer, `final genome\s*:\s*([01]+)`); m != nil {
		parsedGenome = m[0]
	}
	parsedItems := []string{}
	if m := reFind(ctx.Answer, `selected items\s*:\s*(.+)`); m != nil {
		for _, p := range strings.Split(m[0], ",") {
			parsedItems = append(parsedItems, strings.TrimSpace(p))
		}
	}
	bestNeighbor := knapBest(mutants(str(final["genome"])), items, cap)
	feasible := 0
	for v := 0; v < (1 << len(items)); v++ {
		g := fmt.Sprintf("%0*b", len(items), v)
		if integer(knapEval(g, items, cap)["weight"]) <= cap {
			feasible++
		}
	}
	exp := asMap(d["Expected"])
	startOK := len(items) == len(str(d["StartGenome"])) && len(strings.Trim(str(d["StartGenome"]), "01")) == 0
	return []Check{{"the input has one item per genome bit and a valid binary start genome", startOK}, {"the checker independently simulates the deterministic single-bit local search to the reported final genome", parsedGenome == str(final["genome"])}, {"reported selected items are exactly the one-bits in the final genome", sliceEq(parsedItems, final["items"].([]string))}, {"reported weight, value, and fitness match independent genome evaluation", fieldInt(ctx.Answer, "weight") == integer(final["weight"]) && fieldInt(ctx.Answer, "value") == integer(final["value"]) && fieldInt(ctx.Answer, "fitness") == integer(final["fitness"])}, {"the final candidate is feasible and matches the expected fixture totals", integer(final["weight"]) <= cap && str(final["genome"]) == str(exp["FinalGenome"]) && integer(final["weight"]) == integer(exp["FinalWeight"]) && integer(final["value"]) == integer(exp["FinalValue"])}, {"no one-bit neighbor has a lower fitness than the final candidate", integer(bestNeighbor["fitness"]) >= integer(final["fitness"])}, {"the reported generation count matches the independent simulation history length", fieldInt(ctx.Answer, "generations evaluated") == history}, {"the fixture has many feasible genomes, so the check validates the local-search rule rather than a text fragment", feasible > 1000}}
}
