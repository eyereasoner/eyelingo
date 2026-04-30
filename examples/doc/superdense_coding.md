# Superdense Coding

`superdense_coding` is a Go translation/adaptation of Eyeling's `superdense-coding.n3`.

The context is quantum information. Bell-state and parity facts are used to demonstrate how two classical bits are represented in a superdense-coding scenario.

## How it works

A self-contained Go translation of superdense-coding.n3 from the Eyeling
examples.

The N3 example models superdense coding with modal bits, or "mobits".
Think of a mobit as the small teaching-model analogue of a quantum bit. Alice
and Bob share the entangled state |R). Alice encodes one of four two-bit
messages by applying one relation to her half. Bob then applies a joint test
to decode it. Because the model uses GF(2), meaning arithmetic modulo 2,
duplicate derivations cancel: an answer is kept only when it appears an odd
number of times.

This Go version keeps that rule structure visible: relation composition,
superdense candidate generation, parity cancellation, and audit checks are all
explicit.

## What it demonstrates

This example is mainly in the **Science** category. Quantum-information parity facts for superdense coding.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/superdense_coding.json](../input/superdense_coding.json)

Go translation: [../superdense_coding.go](../superdense_coding.go)

Expected Markdown output: [../output/superdense_coding.md](../output/superdense_coding.md)
