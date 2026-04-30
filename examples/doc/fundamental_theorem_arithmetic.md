# Fundamental Theorem Arithmetic

`fundamental_theorem_arithmetic` is a Go translation/adaptation of Eyeling's `fundamental-theorem-arithmetic.n3`.

The context is integer factorization. The example checks that a number decomposes into prime factors and that the grouped prime-power representation is consistent.

## How it works

A self-contained Go translation of fundamental-theorem-arithmetic.n3 from
the Eyeling examples.

The Fundamental Theorem of Arithmetic says that every integer greater than 1
has a prime factorization, and that the factorization is unique except for
the order of the factors. The source N3 file demonstrates this with:

	202692987 = 3^2 * 7 * 829 * 3881

This Go version keeps that source example as the primary case and adds a
wider set of numbers, including larger composites, repeated factors, and a
large prime. Each case is factored by repeatedly taking the smallest divisor.
The program then checks the product, checks that the factors are prime, and
compares smallest-first and largest-first traversals to confirm uniqueness up
to order.

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
