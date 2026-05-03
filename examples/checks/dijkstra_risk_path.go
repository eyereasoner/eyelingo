package main

import (
	"math"
)

func checkDijkstra(ctx *Context) []Check {
	d := ctx.M()
	path, raw, risk, score, edgesCount := parseRiskPath(ctx.Answer)
	w := num(d["riskWeight"])
	edges := maps(d["edges"])
	tr, tk, ts, valid := riskPathTotals(edges, path, w)
	all := riskAllPaths(d)
	bestScore := math.Inf(1)
	bestPath := []string{}
	depotBest := math.Inf(1)
	allNotLower := true
	for _, p := range all {
		_, _, s, _ := riskPathTotals(edges, p, w)
		if s < bestScore {
			bestScore = s
			bestPath = p
		}
		if setOf(p)["DepotC"] && s < depotBest {
			depotBest = s
		}
	}
	for _, p := range all {
		_, _, s, _ := riskPathTotals(edges, p, w)
		allNotLower = allNotLower && s >= bestScore-1e-12
	}
	edgeWeights := true
	for _, e := range edges {
		edgeWeights = edgeWeights && edgeRiskScore(e, w) > num(e["cost"])
	}
	return []Check{{"the graph fixture has the requested start, goal, risk weight, and eight directed edges", str(d["start"]) == "ClinicA" && str(d["goal"]) == "HubZ" && w == 2.0 && len(edges) == 8}, {"every edge weight is independently computed as cost + riskWeight × risk", edgeWeights}, {"the reported path is made only of directed edges present in the input graph", valid && path[0] == str(d["start"]) && path[len(path)-1] == str(d["goal"])}, {"reported raw cost, risk sum, score, and edge count match the parsed path", valid && close(raw, tr, 5e-9) && close(risk, tk, 5e-9) && close(score, ts, 5e-9) && edgesCount == len(path)-1}, {"an independent Dijkstra-style search selects ClinicA -> DepotB -> LabD -> HubZ", sliceEq(bestPath, []string{"ClinicA", "DepotB", "LabD", "HubZ"})}, {"the independent shortest-path score matches the fixture expectation", close(bestScore, num(asMap(d["expected"])["score"]), 1e-9)}, {"exhaustive simple-path enumeration finds no lower risk-adjusted score", allNotLower}, {"the best route through DepotC is independently more expensive than the selected route", depotBest > bestScore && !setOf(path)["DepotC"]}}
}
