# Ackermann

`ackermann` is a Go translation/adaptation of Eyeling's `ackermann.n3`.

The context is exact integer reasoning over the Ackermann hierarchy. The Go version keeps the huge-number behavior explicit, so the example is useful for checking recursion, memoization, and exact arithmetic without hiding the result behind floating-point approximations.

## How it works

A self-contained Go translation of ackermann.n3 from the Eyeling examples.

The original N3 file defines Ackermann answers by first calling a more
general three-argument helper:

	ackermann(x,y) = ackermann3(x, y+3, 2) - 3

That helper covers a ladder of operations: successor, addition,
multiplication, exponentiation, and then larger repeated-power operations
such as tetration and pentation. The test query asks for twelve values,
including exact integers with hundreds and tens of thousands of digits.

This Go version keeps the reduction and recursive rules explicit. It uses
math/big for exact arithmetic and reports large answers by decimal digit
count plus SHA-256 fingerprint so the output stays readable and auditable.

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
