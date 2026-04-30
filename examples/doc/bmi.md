# BMI — ARC-style Body Mass Index example

`bmi` is a Go translation/adaptation of Eyeling's `bmi.n3`.

The context is a familiar clinical measurement. Height and weight facts are converted into a BMI value, category, and healthy-weight interval, making this a simple entry point for numeric rules with checks.

## How it works

A self-contained Go translation of examples/bmi.n3 from the Eyeling
example suite, in ARC style.

The original N3 program encodes a small Body Mass Index calculator that
normalizes metric or US inputs, computes BMI, assigns a WHO adult category,
derives a healthy-range weight band, and runs nine independent checks.

This is intentionally not a generic N3 reasoner.  The concrete N3 facts and
rules are represented as ordinary Go data and functions so the
probabilistic inference is easy to read and directly runnable.

## What it demonstrates

This example is mainly in the **Science** category. Adult BMI calculation, category assignment, and healthy-weight interval.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/bmi.json](../input/bmi.json)

Go translation: [../bmi.go](../bmi.go)

Expected Markdown output: [../output/bmi.md](../output/bmi.md)
