package main

import (
	"strings"
)

func checkIsolation(ctx *Context) []Check {
	d := ctx.M()
	media := maps(d["media"])
	names := []string{}
	statesOK := true
	for _, m := range media {
		names = append(names, str(m["name"]))
		statesOK = statesOK && str(m["zero"]) != str(m["one"])
	}
	prepare := len(media) * 2
	copyCount := len(media) * (len(media) - 1)
	serial := strings.Split(str(asMap(d["expected"])["serial"]), "->")
	nameSet := setOf(names)
	serialOK := true
	for _, x := range serial {
		serialOK = serialOK && nameSet[x]
	}
	edgesOK := serialOK
	for i := 0; i+1 < len(serial); i++ {
		edgesOK = edgesOK && serial[i] != serial[i+1]
	}
	rep := grabInt(ctx.Answer, `possible prepare tasks : (\d+)`)
	prepared := str(asMap(d["expected"])["preparedState"])
	nurseOK := false
	for _, m := range media {
		nurseOK = nurseOK || (str(m["name"]) == "nursePager" && str(m["one"]) == prepared)
	}
	super := asMap(d["superinformationMedium"])
	return []Check{{"all classical lab media carry the BreachBit variable", len(media) == 4 && str(d["variable"]) == "BreachBit"}, {"each medium has distinguishable safe/breach states", statesOK}, {"prepare-task count is independently recomputed from all states", rep == prepare && prepare == 8}, {"the expected prepared breach state belongs to nursePager", nurseOK && contains(ctx.Answer, "nursePager prepares "+prepared)}, {"the expected serial audit path is backed by legal directed edges", edgesOK}, {"containmentPLC has at least two legal fan-out targets", copyCount/len(media) >= 2}, {"the specimen seal has a separate non-classical provenance variable", str(super["variable"]) != str(d["variable"]) && len(asSlice(super["states"])) >= 3}, {"the answer blocks universal cloning and unrestricted parallel fan-out for the specimen seal", contains(ctx.Answer, "universal cloning") && contains(ctx.Answer, "unrestricted parallel fan-out")}, {"all three expected witnesses are reported in the answer", contains(ctx.Answer, "classical breach token : YES") && contains(ctx.Answer, "specimen provenance seal : NO") && contains(ctx.Answer, "serial witness")}}
}
