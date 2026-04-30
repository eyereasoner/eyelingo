# GPS — Goal driven route planning

`gps` is a Go translation/adaptation of Eyeling's `gps.n3`.

The context is goal-driven route planning. Road-network facts and goals are turned into a concrete route with derived cost and step information.

## What it demonstrates

This example is mainly in the **Engineering** category. Goal-driven route planning over a small road network.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/gps.json](../input/gps.json)

Go translation: [../gps.go](../gps.go)

Expected Markdown output: [../output/gps.md](../output/gps.md)
