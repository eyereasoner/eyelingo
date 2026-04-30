# Dining Philosophers

`dining_philosophers` is a Go translation/adaptation of Eyeling's `dining-philosophers.n3`.

The context is concurrent resource use. A Chandy-Misra style trace is replayed so the output can check fairness, ownership, and conflict freedom between philosophers and forks.

## What it demonstrates

This example is mainly in the **Mathematics** category. Chandy-Misra dining-philosophers trace with concurrency conflict checks.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/dining_philosophers.json](../input/dining_philosophers.json)

Go translation: [../dining_philosophers.go](../dining_philosophers.go)

Expected Markdown output: [../output/dining_philosophers.md](../output/dining_philosophers.md)
