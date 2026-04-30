# Goldbach 1000

`goldbach_1000` is a Go translation/adaptation of Eyeling's `goldbach-1000.n3`.

The context is bounded number-theory checking. Every even integer in the selected range is tested for a Goldbach decomposition, giving a compact certificate for the finite bound.

## How it works

Inspired by Eyeling's `examples/goldbach-1000.n3`.

## What it demonstrates

This example is mainly in the **Mathematics** category. Bounded strong-Goldbach checker for every even integer from 4 through 1000.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/goldbach_1000.json](../input/goldbach_1000.json)

Go translation: [../goldbach_1000.go](../goldbach_1000.go)

Expected Markdown output: [../output/goldbach_1000.md](../output/goldbach_1000.md)
