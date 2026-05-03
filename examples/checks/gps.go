package main

import (
	"strings"
)

func checkGPS(ctx *Context) []Check {
	d := ctx.M()
	start := str(asMap(d["Traveller"])["Location"])
	goal := "Oostende"
	if m := reFind(str(d["Question"]), `to ([A-Za-zÀ-ÿ]+)`); m != nil {
		goal = strings.TrimRight(m[0], "?")
	}
	paths := gpsPaths(maps(d["Edges"]), start, goal)
	metrics := map[string][4]float64{}
	for _, p := range paths {
		a, b, c, dd := gpsMetrics(p)
		metrics[gpsLabel(p)] = [4]float64{a, b, c, dd}
	}
	direct := str(asMap(asMap(d["Routes"])["routeDirect"])["Label"])
	alt := str(asMap(asMap(d["Routes"])["routeViaKortrijk"])["Label"])
	dm, directOK := metrics[direct]
	am, altOK := metrics[alt]
	chosen := gpsDominatingRoute(metrics)
	return []Check{
		{"the named direct route is derived from map edges", directOK},
		{"the named Kortrijk alternative route is derived from map edges", altOK},
		{"all simple routes connect the traveller to the destination", len(metrics) == len(paths) && len(metrics) > 0},
		{"route duration and cost are recomputed additively from edges", directOK && altOK && dm[0] > 0 && am[0] > dm[0] && dm[1] > 0 && am[1] > dm[1]},
		{"route belief and comfort are recomputed multiplicatively from edges", directOK && altOK && dm[2] > am[2] && dm[3] > am[3]},
		{"the independently chosen route is faster than its comparator", chosen == direct && dm[0] < am[0]},
		{"the independently chosen route is cheaper than its comparator", chosen == direct && dm[1] < am[1]},
		{"the independently chosen route has higher belief and comfort scores", chosen == direct && dm[2] > am[2] && dm[3] > am[3]},
		{"the answer names the independently chosen route", chosen != "" && contains(ctx.Answer, chosen)},
	}
}

func gpsDominatingRoute(metrics map[string][4]float64) string {
	best := ""
	for candidate, cm := range metrics {
		dominatesAll := true
		for other, om := range metrics {
			if candidate == other {
				continue
			}
			if !(cm[0] < om[0] && cm[1] < om[1] && cm[2] > om[2] && cm[3] > om[3]) {
				dominatesAll = false
				break
			}
		}
		if dominatesAll {
			if best != "" {
				return ""
			}
			best = candidate
		}
	}
	return best
}
