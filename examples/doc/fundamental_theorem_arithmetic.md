# Fundamental Theorem Arithmetic

`fundamental_theorem_arithmetic` is a Go translation/adaptation of Eyeling's `fundamental-theorem-arithmetic.n3`.

The context is integer factorization. The example checks that a number decomposes into prime factors and that the grouped prime-power representation is consistent.

## What it demonstrates

This example is mainly in the **Mathematics** category. Prime factorization and prime-power representation.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/fundamental_theorem_arithmetic.json](../input/fundamental_theorem_arithmetic.json)

Go translation: [../fundamental_theorem_arithmetic.go](../fundamental_theorem_arithmetic.go)

Expected Markdown output: [../output/fundamental_theorem_arithmetic.md](../output/fundamental_theorem_arithmetic.md)
