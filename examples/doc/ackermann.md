# Ackermann

`ackermann` is a Go translation/adaptation of Eyeling's `ackermann.n3`.

The context is exact integer reasoning over the Ackermann hierarchy. The Go version keeps the huge-number behavior explicit, so the example is useful for checking recursion, memoization, and exact arithmetic without hiding the result behind floating-point approximations.

## What it demonstrates

This example is mainly in the **Mathematics** category. Exact Ackermann and hyperoperation facts, including very large integer results.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/ackermann.json](../input/ackermann.json)

Go translation: [../ackermann.go](../ackermann.go)

Expected Markdown output: [../output/ackermann.md](../output/ackermann.md)
