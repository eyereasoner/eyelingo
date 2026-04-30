# Drone Corridor Planner

`drone_corridor_planner` is a Go translation/adaptation of Eyeling's `drone-corridor-planner.n3`.

## Background

Drone route planning is a constrained graph-search problem. A feasible route must connect the origin to the destination while respecting corridor permissions, no-fly zones, battery or distance limits, and operational risk. The planner therefore evaluates paths not only by reachability but also by whether each segment satisfies the safety envelope.

## What it demonstrates

**Category:** Engineering. Constrained drone route planning through corridors and restricted zones.

## How the Go implementation works

The Go planner composes corridor actions into bounded route plans. For each partial plan it accumulates duration and cost, multiplies belief and comfort scores, and prunes candidates that violate restrictions or query thresholds.

The accepted plan is then checked for route continuity, budget limits, risk constraints, and expected endpoint.

## Files

Input JSON: [../input/drone_corridor_planner.json](../input/drone_corridor_planner.json)

Go implementation: [../drone_corridor_planner.go](../drone_corridor_planner.go)

Expected Markdown output: [../output/drone_corridor_planner.md](../output/drone_corridor_planner.md)
