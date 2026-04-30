# 8-Queens

`queens` is a Go translation/adaptation of Eyeling's `queens.n3`.

The context is constraint satisfaction. The eight queens are placed so that rows, columns, and diagonals remain conflict-free, giving a compact combinatorial proof object.

## What it demonstrates

This example is mainly in the **Mathematics** category. 8-Queens constraint satisfaction with a valid board solution.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/queens.json](../input/queens.json)

Go translation: [../queens.go](../queens.go)

Expected Markdown output: [../output/queens.md](../output/queens.md)
