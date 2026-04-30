# Kaprekar 6174

`kaprekar_6174` is a Go translation/adaptation of Eyeling's `kaprekar-6174.n3`.

The context is finite iterative arithmetic. Four-digit Kaprekar chains are followed to 6174 with checks for bounded convergence and trace consistency.

## What it demonstrates

This example is mainly in the **Mathematics** category. Kaprekar chains and basin facts ending at 6174.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/kaprekar_6174.json](../input/kaprekar_6174.json)

Go translation: [../kaprekar_6174.go](../kaprekar_6174.go)

Expected Markdown output: [../output/kaprekar_6174.md](../output/kaprekar_6174.md)
