# GPS — Goal driven route planning

`gps` is a Go translation/adaptation of Eyeling's `gps.n3`.

## Background

Goal-driven route planning starts with a desired destination and searches for paths that satisfy the goal. Each candidate route can carry metrics such as distance, time, or number of legs, and the planner can choose the route that best matches the scoring rule. This example uses a small Belgian map so the route derivation and recommendation are easy to audit.

## What it demonstrates

**Category:** Engineering. Goal-driven route planning over a small road network.

## How the Go implementation works

The Go implementation turns road-network facts into explicit route states. It expands candidate paths from the start city to the goal, computes route metrics, ranks the alternatives, and formats the selected route with step-by-step evidence.

Checks compare the recommendation with the computed distance, cost, and route consistency.

## Files

Input JSON: [../input/gps.json](../input/gps.json)

Go implementation: [../gps.go](../gps.go)

Python check: [../checks/gps.py](../checks/gps.py)

Expected Markdown output: [../output/gps.md](../output/gps.md)
