# Genetic Knapsack Selection

`genetic_knapsack_selection` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on deterministic genetic selection for a bounded knapsack. Its input fixture is organized around `CaseName`, `Question`, `Capacity`, `MaxGenerations`, `StartGenome`, `Items`, `Expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: final genome : 101000000101 selected items : item01, item03, item10, item12 weight : 50 / 50

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - 12 items align with a 12-bit genome C2 OK - final weight 50 is within capacity 50 C3 OK - final genome is 101000000101

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/genetic_knapsack_selection.json](../input/genetic_knapsack_selection.json)

Go translation: [../genetic_knapsack_selection.go](../genetic_knapsack_selection.go)

Expected Markdown output: [../output/genetic_knapsack_selection.md](../output/genetic_knapsack_selection.md)
