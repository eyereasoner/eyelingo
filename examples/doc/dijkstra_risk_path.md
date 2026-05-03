# Dijkstra Risk Path

`dijkstra_risk_path` is a Go translation/adaptation of Eyeling's `dijkstra.n3`.

## Background

Path-finding in a weighted graph chooses a route by minimizing a cost assigned to each edge. Dijkstra's algorithm is the standard method when all edge costs are non-negative. In a risk-aware variant, the weights may combine distance, hazard, exposure, or operational risk, so the selected path is the lowest-risk route rather than simply the shortest route.

## What it demonstrates

**Category:** Engineering. Risk-adjusted path selection using weighted network edges.

## How the Go implementation works

The Go code builds a weighted graph where each edge has both cost and risk. It applies Dijkstra's algorithm to the risk-adjusted score, carrying the current path with each label so the chosen route can be printed directly.

The check section verifies the selected path, total raw cost, accumulated risk, and final score.

## Files

Input JSON: [../input/dijkstra_risk_path.json](../input/dijkstra_risk_path.json)

Go implementation: [../dijkstra_risk_path.go](../dijkstra_risk_path.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/dijkstra_risk_path.md](../output/dijkstra_risk_path.md)
