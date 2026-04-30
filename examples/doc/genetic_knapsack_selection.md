# Genetic Knapsack Selection

`genetic_knapsack_selection` is a Go translation/adaptation of Eyeling's `genetic-knapsack-selection.n3`.

The context is optimization with deterministic genetic steps. Candidate knapsack selections are scored, mutated, and filtered while preserving capacity and value checks.

## What it demonstrates

This example is mainly in the **Mathematics** category. Deterministic genetic selection for a bounded knapsack.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/genetic_knapsack_selection.json](../input/genetic_knapsack_selection.json)

Go translation: [../genetic_knapsack_selection.go](../genetic_knapsack_selection.go)

Expected Markdown output: [../output/genetic_knapsack_selection.md](../output/genetic_knapsack_selection.md)
