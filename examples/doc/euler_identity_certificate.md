# Euler Identity Certificate

`euler_identity_certificate` is a Go translation/adaptation of Eyeling's `euler-identity-certificate.n3`.

The context is numerical certification of a classic identity. High-precision arithmetic is used to show that exp(iπ)+1 is close enough to zero under an explicit tolerance.

## What it demonstrates

This example is mainly in the **Mathematics** category. High-precision certificate for the identity exp(iπ) + 1 = 0.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/euler_identity_certificate.json](../input/euler_identity_certificate.json)

Go translation: [../euler_identity_certificate.go](../euler_identity_certificate.go)

Expected Markdown output: [../output/euler_identity_certificate.md](../output/euler_identity_certificate.md)
