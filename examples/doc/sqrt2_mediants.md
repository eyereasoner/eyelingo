# Integer-First Sqrt2 Mediants

`sqrt2_mediants` is a Go translation/adaptation of Eyeling's `integer-first-sqrt2-mediants.n3`.

The context is exact rational approximation. Integer square comparisons certify lower and upper mediants for sqrt(2), avoiding floating-point dependence.

## What it demonstrates

This example is mainly in the **Mathematics** category. Rational lower/upper bounds for sqrt(2) certified by integer square comparisons.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/sqrt2_mediants.json](../input/sqrt2_mediants.json)

Go translation: [../sqrt2_mediants.go](../sqrt2_mediants.go)

Expected Markdown output: [../output/sqrt2_mediants.md](../output/sqrt2_mediants.md)
