# Flandor

`flandor` is a Go translation/adaptation of Eyeling's `flandor.n3`.

The context is regional industrial policy. Exporters, training capacity, intervention packages, and payload-signature checks are combined to select a retooling priority.

## What it demonstrates

This example is mainly in the **Engineering** category. Regional retooling priority calculation for a Flanders scenario.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/flandor.json](../input/flandor.json)

Go translation: [../flandor.go](../flandor.go)

Expected Markdown output: [../output/flandor.md](../output/flandor.md)
