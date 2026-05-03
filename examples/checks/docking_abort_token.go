package main

import (
	"strings"
)

func checkDocking(ctx *Context) []Check {
	d := ctx.M()
	media := maps(d["media"])
	names := []string{}
	for _, m := range media {
		names = append(names, str(m["name"]))
	}
	copyCount := len(names) * (len(names) - 1)
	reported := grabInt(ctx.Answer, `possible copy tasks : (\d+)`)
	serial := strings.Split(str(asMap(d["expected"])["serial"]), "->")
	parallel := str(asMap(d["expected"])["parallelSource"])
	nameSet := setOf(names)
	serialKnown := true
	for _, x := range serial {
		serialKnown = serialKnown && nameSet[x]
	}
	distinct := true
	for _, m := range media {
		distinct = distinct && str(m["zero"]) != str(m["one"])
	}
	super := asMap(d["superinformationMedium"])
	return []Check{{"all four media carry the same AbortBit variable", len(media) == 4 && str(d["variable"]) == "AbortBit"}, {"each medium has distinct zero and one states", distinct}, {"the directed copy-task count is recomputed as n*(n-1)", reported == copyCount && copyCount == 12}, {"the expected serial witness uses known media in order", serialKnown && contains(ctx.Answer, serial[0]+" -> "+serial[1]+" -> "+serial[2])}, {"the serial witness is backed by two legal copy/measure edges", len(serial) >= 3 && serial[0] != serial[1] && serial[1] != serial[2]}, {"the expected parallel source can fan out to two other media", nameSet[parallel] && copyCount/len(names) >= 3 && contains(ctx.Answer, "parallel witness : flightPLC -> radioFrame and auditDisplay")}, {"the quantum seal is separate from the classical AbortBit variable", str(super["variable"]) != str(d["variable"])}, {"the answer blocks universal cloning and unrestricted audit fan-out for the seal", contains(ctx.Answer, "cannot be universally cloned") && contains(ctx.Answer, "unrestricted audit fan-out")}}
}
