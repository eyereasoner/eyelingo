# Superdense Coding

`superdense_coding` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on quantum-information parity facts for superdense coding. Its input fixture is organized around `States`, `Primitive`, `AliceOps`, `BobTests`, `Note`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Superdense-coding facts that survive GF(2) parity cancellation: 0 dqc:superdense-coding 0 1 dqc:superdense-coding 1

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: OK shared entanglement: |R) contains exactly |0,0) and |1,1) OK composition KG: KG is obtained by composing G then K, exactly as in the N3 rule OK composition GK: GK is obtained by composing K then G, exactly as in the N3 rule

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/superdense_coding.json](../input/superdense_coding.json)

Go translation: [../superdense_coding.go](../superdense_coding.go)

Expected Markdown output: [../output/superdense_coding.md](../output/superdense_coding.md)
