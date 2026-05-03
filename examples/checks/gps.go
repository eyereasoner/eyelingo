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
	dm, am := metrics[direct], metrics[alt]
	best := ""
	for lab := range metrics {
		if best == "" || metrics[lab][0] < metrics[best][0] {
			best = lab
		}
	}
	return []Check{
		{"the direct Gent-Brugge-Oostende route is derived from edges", metrics[direct][0] != 0},
		{"the Kortrijk alternative route is derived from edges", metrics[alt][0] != 0},
		{"exactly two simple routes connect the traveller to the destination", len(metrics) == 2},
		{"route duration and cost are additive over edges", dm[0] == 2400.0 && close(dm[1], 0.010, 1e-12) && am[0] == 4100.0 && close(am[1], 0.018, 1e-12)},
		{"route belief and comfort are multiplicative over edges", close(dm[2], 0.9408, 1e-12) && close(dm[3], 0.99, 1e-12) && close(am[2], 0.903168, 1e-12) && close(am[3], 0.9801, 1e-12)},
		{"the recommended route is faster than the alternative", dm[0] < am[0] && best == direct},
		{"the recommended route is cheaper than the alternative", dm[1] < am[1]},
		{"the recommended route has higher belief and comfort scores", dm[2] > am[2] && dm[3] > am[3]},
		{"the answer names the independently chosen direct route", contains(ctx.Answer, direct) && contains(ctx.Answer, "Take the direct route")},
	}
}
