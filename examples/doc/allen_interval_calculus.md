# Allen Interval Calculus

`allen_interval_calculus` is a Go translation/adaptation of Eyeling's `allen-interval-calculus.n3`.

The context is temporal reasoning. Intervals, endpoints, and Allen relations are derived and cross-checked so that ordering facts can be used as reliable inputs to planning or scheduling logic.

## How it works

A compact Go translation inspired by Eyeling's
`examples/allen-interval-calculus.n3`.

The example completes intervals with duration fields and classifies every
ordered interval pair using Allen's 13 base relations.

## What it demonstrates

This example is mainly in the **Mathematics** category. Allen temporal interval relation closure over completed and explicit intervals.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/allen_interval_calculus.json](../input/allen_interval_calculus.json)

Go translation: [../allen_interval_calculus.go](../allen_interval_calculus.go)

Expected Markdown output: [../output/allen_interval_calculus.md](../output/allen_interval_calculus.md)
