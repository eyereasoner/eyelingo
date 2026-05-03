package main

import (
	"sort"
	"strings"
)

func checkSchool(ctx *Context) []Check {
	d := ctx.M()
	straight := map[string]schoolAssign{}
	audited := map[string]schoolAssign{}
	changes := []struct {
		student           anymap
		straight, audited schoolAssign
	}{}
	maxWalk := integer(asMap(d["policy"])["maxWalkingMeters"])
	for _, st := range maps(d["students"]) {
		fl := chooseSchool(d, st, false)
		be := chooseSchool(d, st, true)
		straight[str(st["id"])] = fl
		audited[str(st["id"])] = be
		if fl.school != be.school || integer(fl.distance["walkingMeters"]) > maxWalk {
			changes = append(changes, struct {
				student           anymap
				straight, audited schoolAssign
			}{st, fl, be})
		}
	}
	largest := changes[0]
	for _, ch := range changes {
		if hiddenDetour(ch.straight) > hiddenDetour(largest.straight) {
			largest = ch
		}
	}
	schools := map[string]bool{}
	for _, s := range maps(d["schools"]) {
		schools[str(s["name"])] = true
	}
	prefsOK := true
	for _, st := range maps(d["students"]) {
		prefsOK = prefsOK && len(sarr(st["preferences"])) == len(schools) && setEq(sarr(st["preferences"]), keysBool(schools))
	}
	parsedAssign := map[string]string{}
	if m := reFind(ctx.Answer, `recommended assignments\s*:\s*(.+)`); m != nil {
		for _, part := range strings.Split(m[0], ";") {
			if strings.Contains(part, "->") {
				p := strings.SplitN(part, "->", 2)
				parsedAssign[strings.TrimSpace(p[0])] = strings.TrimSpace(p[1])
			}
		}
	}
	affected := []string{}
	if m := reFind(ctx.Answer, `children affected by straight-line rule\s*:\s*(.+)`); m != nil {
		for _, p := range strings.Split(m[0], ",") {
			if strings.TrimSpace(p) != "" {
				affected = append(affected, strings.TrimSpace(p))
			}
		}
	}
	expectedAssign := map[string]string{}
	provisional := map[string]string{}
	auditedNames := map[string]string{}
	affectedNames := []string{}
	over := 0
	changed := false
	for _, st := range maps(d["students"]) {
		name := str(st["name"])
		expectedAssign[name] = audited[str(st["id"])].school
		provisional[name] = straight[str(st["id"])].school
		auditedNames[name] = audited[str(st["id"])].school
	}
	for _, ch := range changes {
		affectedNames = append(affectedNames, str(ch.student["name"]))
		if integer(ch.straight.distance["walkingMeters"]) > maxWalk {
			over++
		}
		if ch.straight.school != ch.audited.school {
			changed = true
		}
	}
	sort.Strings(affectedNames)
	detourName := ""
	detourMeters := 0
	if m := reFind(ctx.Answer, `largest hidden detour\s*:\s*([^,]+),\s*(\d+)\s*m`); m != nil {
		detourName = strings.TrimSpace(m[0])
		detourMeters = parseInt(m[1])
	}
	return []Check{{"the fixture has four students, four schools, and a complete 4 × 4 distance matrix", len(maps(d["students"])) == 4 && len(maps(d["schools"])) == 4 && len(maps(d["distances"])) == 16}, {"every student preference list covers every school exactly once", prefsOK}, {"the independent straight-line rule assigns Ada and Björn to Centrum", provisional["Ada"] == "Centrum" && provisional["Björn"] == "Centrum"}, {"the independent route-aware rule computes walking distance plus preference penalty", mapStrEq(auditedNames, map[string]string{"Ada": "Lindholmen", "Björn": "Backa", "Clara": "Haga", "Davi": "Haga"})}, {"the reported recommended assignments match the Go audit", mapStrEq(parsedAssign, expectedAssign)}, {"the reported affected children are exactly those whose provisional placement is flagged", sliceEq(affected, affectedNames) && sliceEq(affectedNames, []string{"Ada", "Björn", "Davi"})}, {"the reported largest hidden detour matches the flawed straight-line placement", detourName == str(largest.student["name"]) && detourMeters == hiddenDetour(largest.straight) && detourMeters == 3000}, {"the failure result follows from at least one over-limit walking route and changed assignment", answerField(ctx.Answer, "audit result") == "fail" && over >= 2 && changed}, {"the Reason text names the support-tool rule, walking-route recomputation, and inspectability requirement", contains(strings.ToLower(ctx.Reason), "support-tool") && contains(strings.ToLower(ctx.Reason), "walking-route") && contains(strings.ToLower(ctx.Reason), "inspectable")}, {"the report explicitly requests an explanation for the affected placement decisions", answerField(ctx.Answer, "explanation requested") == "yes"}}
}
