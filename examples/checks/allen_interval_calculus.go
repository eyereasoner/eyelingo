package main

import (
	"strings"
	"time"
)

func checkAllen(ctx *Context) []Check {
	d := ctx.M()
	intervalRows := maps(d["intervals"])
	type iv struct{ start, end time.Time }
	intervals := map[string]iv{}
	completed := map[string]iv{}
	for _, it := range intervalRows {
		start := parseTime(strings.ReplaceAll(str(it["start"]), "Z", "+00:00"))
		var end time.Time
		if _, ok := it["end"]; ok {
			end = parseTime(strings.ReplaceAll(str(it["end"]), "Z", "+00:00"))
		} else {
			end = start.Add(time.Duration(integer(it["durationMinutes"])) * time.Minute)
			completed[str(it["name"])] = iv{start, end}
		}
		intervals[str(it["name"])] = iv{start, end}
	}
	rel := func(a, b iv) string {
		a0, a1, b0, b1 := a.start, a.end, b.start, b.end
		switch {
		case a1.Before(b0):
			return "before"
		case a0.After(b1):
			return "after"
		case a1.Equal(b0):
			return "meets"
		case a0.Equal(b1):
			return "metBy"
		case a0.Before(b0) && a1.After(b0) && a1.Before(b1):
			return "overlaps"
		case a0.After(b0) && a0.Before(b1) && a1.After(b1):
			return "overlappedBy"
		case a0.Equal(b0) && a1.Before(b1):
			return "starts"
		case a0.Equal(b0) && a1.After(b1):
			return "startedBy"
		case a0.After(b0) && a1.Before(b1):
			return "during"
		case a0.Before(b0) && a1.After(b1):
			return "contains"
		case a0.After(b0) && a1.Equal(b1):
			return "finishes"
		case a0.Before(b0) && a1.Equal(b1):
			return "finishedBy"
		case a0.Equal(b0) && a1.Equal(b1):
			return "equals"
		}
		return "unknown"
	}
	ordered := map[[2]string]string{}
	relset := map[string]bool{}
	invalid := false
	for a, iva := range intervals {
		if !iva.end.After(iva.start) {
			invalid = true
		}
		for b, ivb := range intervals {
			if a != b {
				r := rel(iva, ivb)
				ordered[[2]string{a, b}] = r
				relset[r] = true
			}
		}
	}
	req := asMap(asMap(d["expected"])["requiredRelations"])
	requiredOK := true
	for k, v := range req {
		parts := strings.Split(k, "|")
		requiredOK = requiredOK && ordered[[2]string{parts[0], parts[1]}] == str(v)
	}
	reportedCount := -1
	if m := reFind(ctx.Answer, `derived relations : (\d+) ordered interval pairs`); m != nil {
		reportedCount = parseInt(m[0])
	}
	showcase := true
	pairs := map[[2]string]bool{{"A", "B"}: true, {"A", "C"}: true, {"A", "D"}: true, {"F", "A"}: true, {"G", "A"}: true, {"A", "H"}: true, {"A", "E"}: true, {"J", "I"}: true, {"K", "C"}: true}
	for k, v := range req {
		parts := strings.Split(k, "|")
		if pairs[[2]string{parts[0], parts[1]}] {
			showcase = showcase && contains(ctx.Answer, parts[0]+" "+str(v)+" "+parts[1])
		}
	}
	return []Check{{"duration-based intervals are completed from start plus minutes", completed["I"].end.Hour() == 18 && completed["K"].end.Hour() == 14 && completed["K"].end.Minute() == 0}, {"all completed intervals have strictly positive duration", !invalid}, {"all ordered non-self interval pairs are classified", reportedCount == len(intervalRows)*(len(intervalRows)-1) && len(ordered) == reportedCount}, {"every required Allen relation is recomputed from endpoints", requiredOK}, {"A/B, A/C, and A/D demonstrate before, meets, and overlaps", ordered[[2]string{"A", "B"}] == "before" && ordered[[2]string{"A", "C"}] == "meets" && ordered[[2]string{"A", "D"}] == "overlaps"}, {"converse relations are independently recovered", ordered[[2]string{"B", "A"}] == "after" && ordered[[2]string{"C", "A"}] == "metBy" && ordered[[2]string{"D", "A"}] == "overlappedBy"}, {"start, finish, during, contains, and equals cases all occur", relset["starts"] && relset["finishes"] && relset["during"] && relset["contains"] && relset["equals"]}, {"the showcase text includes each required forward example", showcase}}
}
