# Alarm Bit Interoperability

`alarm_bit_interoperability` is a Go translation/adaptation of Eyeling's `act-alarm-bit-interoperability.n3`.

The context is interoperability between classical alarm systems. The example separates ordinary copying and relay of a classical alarm bit from quantum-like no-cloning constraints, which makes it a small but concrete technology-facing rule example.

## How it works

A Go translation inspired by Eyeling's
`examples/act-alarm-bit-interoperability.n3`.

The example distinguishes what can be done with classical alarm-bit media
from what cannot be done with a superinformation-like token.

## What it demonstrates

This example is mainly in the **Technology** category. Classical alarm-bit copy and relay tasks contrasted with forbidden universal cloning.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/alarm_bit_interoperability.json](../input/alarm_bit_interoperability.json)

Go translation: [../alarm_bit_interoperability.go](../alarm_bit_interoperability.go)

Expected Markdown output: [../output/alarm_bit_interoperability.md](../output/alarm_bit_interoperability.md)
