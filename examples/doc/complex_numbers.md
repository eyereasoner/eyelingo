# Complex Numbers

`complex_numbers` is a Go translation/adaptation of Eyeling's `complex.n3`.

The context is complex arithmetic. The example translates polar-form and transcendental identities into Go calculations and checks quadrant-sensitive behavior in a way that is easy to audit.

## How it works

A self-contained Go translation of complex.n3 from the Eyeling examples.

The original N3 file defines rules for complex polar conversion,
quadrant-sensitive angle selection, complex exponentiation, and inverse sine /
cosine over complex numbers. The test query asks for six derived complex
values:

	sqrt(-1), e^(i*pi), i^i, e^(-pi/2), asin(2), and acos(2)

This Go version keeps those rules explicit instead of using Go's cmplx.Pow,
cmplx.Asin, or cmplx.Acos helpers. That makes the same mathematical proof
steps visible and auditable.

## What it demonstrates

This example is mainly in the **Mathematics** category. Complex arithmetic and transcendental identity checks.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/complex_numbers.json](../input/complex_numbers.json)

Go translation: [../complex_numbers.go](../complex_numbers.go)

Expected Markdown output: [../output/complex_numbers.md](../output/complex_numbers.md)
