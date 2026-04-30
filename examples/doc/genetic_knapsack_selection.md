# Genetic Knapsack Selection

`genetic_knapsack_selection` is a Go translation/adaptation of Eyeling's `genetic-knapsack-selection.n3`.

## Background

The knapsack problem asks which items to choose when each item has value and weight and the total weight is bounded. Genetic algorithms search such spaces by keeping a population of candidate selections, scoring fitness, and applying selection, crossover, or mutation. This example uses a deterministic version so the optimization trace is reproducible.

## What it demonstrates

**Category:** Mathematics. Deterministic genetic selection for a bounded knapsack.

## How the Go implementation works

The implementation represents each knapsack candidate as a bit genome. It scores the current genome, generates all single-bit mutants, filters out capacity violations, and moves to the improving candidate with the best fitness until no mutant improves the score.

The checks verify capacity, selected value, mutation count, and convergence of the deterministic search.

## Files

Input JSON: [../input/genetic_knapsack_selection.json](../input/genetic_knapsack_selection.json)

Go implementation: [../genetic_knapsack_selection.go](../genetic_knapsack_selection.go)

Expected Markdown output: [../output/genetic_knapsack_selection.md](../output/genetic_knapsack_selection.md)
