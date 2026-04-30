# Calidor

`calidor` is a Go translation/adaptation of Eyeling's `calidor.n3`.

The context is municipal response planning during heat or cooling stress. Needs, packages, budgets, signatures, and policy constraints are combined to decide whether an intervention bundle is admissible.

## What it demonstrates

This example is mainly in the **Engineering** category. Municipal cooling intervention bundle chosen from active needs and budget constraints.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/calidor.json](../input/calidor.json)

Go translation: [../calidor.go](../calidor.go)

Expected Markdown output: [../output/calidor.md](../output/calidor.md)
