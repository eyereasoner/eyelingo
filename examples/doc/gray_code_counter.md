# Gray Code Counter

`gray_code_counter` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on n-bit Gray-code sequence with one-bit transition checks. Its input fixture is organized around `caseName`, `question`, `bits`, `steps`, `expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Technology** example. It demonstrates data representation, interoperability, policies, and computational artifacts in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: bits : 4 states visited : 16 unique states : 16

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - 16 states were generated for a 4-bit counter C2 OK - all generated states are unique C3 OK - each adjacent transition flips exactly one bit

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/gray_code_counter.json](../input/gray_code_counter.json)

Go translation: [../gray_code_counter.go](../gray_code_counter.go)

Expected Markdown output: [../output/gray_code_counter.md](../output/gray_code_counter.md)
