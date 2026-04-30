# LLD — Leg Length Discrepancy Measurement

`lldm` is a Go translation/adaptation of Eyeling's `lldm.n3`.

The context is clinical measurement. Leg-length values are compared to derive discrepancy size and whether an alarm threshold has been crossed.

## How it works

A self-contained Go translation of examples/lldm.n3 from the Eyeling
example suite, in ARC style.

The original N3 program encodes a geometric model for Leg Length
Discrepancy Measurement.  Four measurement points on a medical image are
processed through a pipeline of coordinate differences, line slopes,
intersection points, and Euclidean distances to decide whether an alarm
should be raised.

This is intentionally not a generic N3 reasoner.  The concrete N3 facts
and derivation rules are represented as ordinary Go data and functions so
the decision logic is easy to read and directly runnable.

## What it demonstrates

This example is mainly in the **Science** category. Leg-length discrepancy measurement and alarm thresholding.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/lldm.json](../input/lldm.json)

Go translation: [../lldm.go](../lldm.go)

Expected Markdown output: [../output/lldm.md](../output/lldm.md)
