package main

import (
	"strings"
)

func checkBarley(ctx *Context) []Check {
	d := ctx.M()
	world := asMap(d["world"])
	blockers := map[string][]string{}
	evolv := []string{}
	blocked := []string{}
	for _, line := range maps(d["lineages"]) {
		name := str(line["name"])
		b := barleyBlockers(world, line)
		blockers[name] = b
		if len(b) == 0 {
			evolv = append(evolv, name)
		} else {
			blocked = append(blocked, name)
		}
	}
	reported := []string{}
	if m := reFind(ctx.Answer, `blocked contrast lineages : ([^\n]+)`); m != nil {
		for _, x := range strings.Split(m[0], ",") {
			reported = append(reported, strings.TrimSpace(x))
		}
	}
	expected := asMap(d["expected"])
	var mainLine anymap
	for _, l := range maps(d["lineages"]) {
		if str(l["name"]) == "mainLine" {
			mainLine = l
		}
	}
	return []Check{{"no-design laws and greenhouse resources are available", boolean(world["noDesignLaws"]) && setEq(sarr(world["greenhouse"]), []string{"warmth", "moisture", "light"})}, {"mainLine satisfies every CAN condition", len(blockers["mainLine"]) == 0}, {"mainLine is the unique evolvable lineage", len(evolv) == 1 && evolv[0] == str(expected["evolvable"]) && contains(ctx.Answer, "evolvable lineage : mainLine")}, {"analogLine is blocked by missing digital heredity", setEq(blockers["analogLine"], []string{"digital-heredity"})}, {"fragileLine is blocked by missing repair", setEq(blockers["fragileLine"], []string{"repair"})}, {"coatlessLine is blocked by missing dormancy protection", setEq(blockers["coatlessLine"], []string{"protected-dormancy"})}, {"staticLine is blocked by missing heritable variation", setEq(blockers["staticLine"], []string{"heritable-variation"})}, {"reported blocked lineages match independent blocker analysis", sliceEq(reported, sarr(expected["blocked"])) && sliceEq(sarr(expected["blocked"]), blocked)}, {"adaptive persistence follows from the selected salt-tolerant variant", str(mainLine["variant"]) == str(world["selectionFavours"])}}
}
