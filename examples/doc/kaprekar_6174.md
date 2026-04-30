# Kaprekar 6174

`kaprekar_6174` is a Go translation/adaptation of Eyeling's `kaprekar-6174.n3`.

The context is finite iterative arithmetic. Four-digit Kaprekar chains are followed to 6174 with checks for bounded convergence and trace consistency.

## How it works

A self-contained Go translation of kaprekar-6174.n3 from the Eyeling
examples.

Kaprekar's routine starts with any four digits, including leading zeroes.
At each step it sorts those digits from high to low to make one number,
sorts them from low to high to make another number, and subtracts the smaller
from the larger. Most starts eventually reach 6174, Kaprekar's constant.
Starts with four equal digits, such as 0000 or 2222, fall into 0000 instead
and are intentionally omitted from the final :kaprekar facts.

The N3 example unrolls the proof to a fixed seven-step bound, because no
four-digit start that reaches 6174 needs more than seven Kaprekar steps. This
Go version keeps the same bounded, explicit approach while using ordinary Go
data structures instead of a general RDF/N3 reasoner.

## What it demonstrates

This example is mainly in the **Mathematics** category. Kaprekar chains and basin facts ending at 6174.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/kaprekar_6174.json](../input/kaprekar_6174.json)

Go translation: [../kaprekar_6174.go](../kaprekar_6174.go)

Expected Markdown output: [../output/kaprekar_6174.md](../output/kaprekar_6174.md)
