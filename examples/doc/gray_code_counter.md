# Gray Code Counter

`gray_code_counter` is a Go translation/adaptation of Eyeling's `gray-code-counter.n3`.

The context is digital encoding. The generated Gray-code sequence is checked so adjacent codes differ by exactly one bit, a core property for robust counters.

## How it works

Inspired by Eyeling's `examples/gray-code-counter.n3`.

The example generates a cyclic reflected binary Gray counter and audits its
one-bit transition property.

## What it demonstrates

This example is mainly in the **Technology** category. n-bit Gray-code sequence with one-bit transition checks.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/gray_code_counter.json](../input/gray_code_counter.json)

Go translation: [../gray_code_counter.go](../gray_code_counter.go)

Expected Markdown output: [../output/gray_code_counter.md](../output/gray_code_counter.md)
