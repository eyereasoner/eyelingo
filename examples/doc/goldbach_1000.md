# Goldbach 1000

`goldbach_1000` is a Go translation/adaptation of Eyeling's `goldbach-1000.n3`.

The example checks the strong even Goldbach condition for every even integer from 4 through 1000. It caches primes and records representative witnesses such as 1000 = 3 + 997.

## What it demonstrates

This is mainly a **Mathematics** example. This is not a proof of the full conjecture; it is a bounded verification example with explicit arithmetic checks and audit counters.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, counters, thresholds, derived facts, or platform details.

## Files

Input JSON: [../input/goldbach_1000.json](../input/goldbach_1000.json)

Go translation: [../goldbach_1000.go](../goldbach_1000.go)

Expected Markdown output: [../output/goldbach_1000.md](../output/goldbach_1000.md)
