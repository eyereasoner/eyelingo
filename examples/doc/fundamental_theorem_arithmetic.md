# Fundamental Theorem Arithmetic

`fundamental_theorem_arithmetic` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on prime factorization and prime-power representation. Its input fixture is a list with 6 entries.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Primary N3 case: n = 202692987 has prime factors 3 * 3 * 7 * 829 * 3881. primary prime-power form : 3^2 * 7 * 829 * 3881 sample count : 6

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - the source example factors 202692987 as 3,3,7,829,3881 C2 OK - the source example groups repeated factors as 3^2 * 7 * 829 * 3881 C3 OK - multiplying each computed factor list reconstructs its original number

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/fundamental_theorem_arithmetic.json](../input/fundamental_theorem_arithmetic.json)

Go translation: [../fundamental_theorem_arithmetic.go](../fundamental_theorem_arithmetic.go)

Expected Markdown output: [../output/fundamental_theorem_arithmetic.md](../output/fundamental_theorem_arithmetic.md)
