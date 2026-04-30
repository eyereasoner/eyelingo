# Fibonacci Example (Big)

`fibonacci` is a Go translation/adaptation of Eyeling's `fibonacci.n3`.

The context is exact recurrence computation. A large Fibonacci value is produced with integer arithmetic, making it useful for checking deterministic dynamic-programming behavior.

## What it demonstrates

This example is mainly in the **Mathematics** category. Exact computation of a large Fibonacci number.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/fibonacci.json](../input/fibonacci.json)

Go translation: [../fibonacci.go](../fibonacci.go)

Expected Markdown output: [../output/fibonacci.md](../output/fibonacci.md)
