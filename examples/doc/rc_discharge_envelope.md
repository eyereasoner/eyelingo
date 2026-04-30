# RC Discharge Envelope

`rc_discharge_envelope` is a Go translation/adaptation of Eyeling's `rc-discharge-envelope.n3`.

The context is circuit behavior. An RC discharge is bounded by an exponential envelope, and the checks verify that the decay stays within the certified range.

## How it works

Inspired by Eyeling's `examples/floating-point-first-rc-discharge.n3`.

The example propagates a certified floating-point decay interval for a
sampled RC discharge and finds the first safe voltage sample.

## What it demonstrates

This example is mainly in the **Science** category. Certified exponential decay envelope for an RC discharge.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/rc_discharge_envelope.json](../input/rc_discharge_envelope.json)

Go translation: [../rc_discharge_envelope.go](../rc_discharge_envelope.go)

Expected Markdown output: [../output/rc_discharge_envelope.md](../output/rc_discharge_envelope.md)
