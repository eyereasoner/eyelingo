package main

import (
	"strings"
)

func checkAlarmBit(ctx *Context) []Check {
	d := ctx.M()
	media := maps(d["ClassicalMedia"])
	variable := ""
	if len(media) > 0 {
		variable = str(media[0]["Variable"])
	}
	directed := [][2]string{}
	for _, a := range media {
		for _, b := range media {
			if str(a["Name"]) != str(b["Name"]) && str(a["Variable"]) == str(b["Variable"]) {
				directed = append(directed, [2]string{str(a["Name"]), str(b["Name"])})
			}
		}
	}
	reported := map[[2]string]bool{}
	for _, m := range reAll(ctx.Answer, `copy task : ([A-Za-z0-9_]+) -> ([A-Za-z0-9_]+) for ([A-Za-z0-9_]+)`) {
		reported[[2]string{m[1], m[2]}] = true
	}
	blocked := []string{}
	if m := reFind(ctx.Answer, `blocked tasks : ([^\n]+)`); m != nil {
		for _, x := range strings.Split(m[0], ",") {
			blocked = append(blocked, strings.TrimSpace(x))
		}
	}
	same := len(media) >= 2
	distinguish := true
	for _, m := range media {
		same = same && str(m["Variable"]) == variable
		distinguish = distinguish && str(m["ZeroState"]) != str(m["OneState"])
	}
	tasksOK := len(reported) == len(directed)
	for _, p := range directed {
		tasksOK = tasksOK && reported[p]
	}
	super := asMap(d["Superinformation"])
	answerCan := contains(ctx.Answer, "classical alarm-bit interoperability : YES")
	answerCannot := contains(ctx.Answer, "universal cloning of the superinformation token : NO")
	return []Check{{"all classical media encode the same abstract variable", same}, {"the directed copy-task count is recomputed from the media graph", len(directed) == integer(d["ExpectedCopyTasks"])}, {"the report lists exactly the expected directed copy tasks", tasksOK}, {"each classical medium has distinguishable zero and one states", distinguish}, {"the superinformation contrast has more than two named states", len(asSlice(super["States"])) > 2}, {"all expected impossible tasks are reported as blocked", setEq(blocked, sarr(d["ExpectedImpossible"]))}, {"the reported CAN/CANNOT decisions match the expected polarities", answerCan == (str(d["ExpectedCanDecision"]) == "YES") && answerCannot == (str(d["ExpectedCantDecision"]) == "NO")}}
}
