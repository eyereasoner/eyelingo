# EV Roadtrip Planner

`ev_roundtrip_planner` is a Go translation/adaptation of Eyeling's `ev-roundtrip-planner.n3`.

## Background

Electric-vehicle trip planning depends on state of charge, consumption, charger availability, dwell time, cost, and comfort constraints. A feasible round trip must keep the battery above reserve limits, choose charging stops that fit the route, and balance total duration against convenience. This example turns those constraints into an auditable route recommendation.

## What it demonstrates

**Category:** Engineering. EV route planning with battery, duration, cost, and comfort constraints.

## How the Go implementation works

The implementation treats route legs and charging stops as state transitions. It composes bounded plans from Brussels to Cologne, updating battery state, duration, cost, belief, and comfort at each step.

Candidate plans are filtered against the query thresholds, and the selected itinerary is checked for feasibility, final battery margin, and total trip metrics.

## Files

Input JSON: [../input/ev_roundtrip_planner.json](../input/ev_roundtrip_planner.json)

Go implementation: [../ev_roundtrip_planner.go](../ev_roundtrip_planner.go)

Expected Markdown output: [../output/ev_roundtrip_planner.md](../output/ev_roundtrip_planner.md)
