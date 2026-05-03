# Barley Seed Lineage

`barley_seed_lineage` is a Go translation/adaptation of Eyeling's `act-barley-seed-lineage.n3`.

## Background

Seed-lineage reasoning combines biological heredity with constraints on what transformations are physically or biologically possible. A viable seed line needs reproduction, persistence across time, enough copying accuracy, and some heritable variation for selection. Dormancy matters because seeds can remain viable without actively growing, while missing resources or broken heredity block a proposed lineage.

## What it demonstrates

**Category:** Science. Seed-lineage CAN/CAN'T reasoning for reproduction, dormancy, variation, and persistence.

## How the Go implementation works

The implementation evaluates each candidate lineage against the required conditions for viable heredity and selection. It derives copy, accuracy, closure, and adaptive-variation flags, then records the single missing ingredient that blocks each contrast lineage.

The final result separates the evolvable lineage from the blocked ones and checks both lists against the expected fixture.

## Files

Input JSON: [../input/barley_seed_lineage.json](../input/barley_seed_lineage.json)

Go implementation: [../barley_seed_lineage.go](../barley_seed_lineage.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/barley_seed_lineage.md](../output/barley_seed_lineage.md)
