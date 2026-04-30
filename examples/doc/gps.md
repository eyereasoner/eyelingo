# GPS — Goal driven route planning

`gps` is a Go translation/adaptation of Eyeling's `gps.n3`.

The context is goal-driven route planning. Road-network facts and goals are turned into a concrete route with derived cost and step information.

## How it works

A self-contained Go translation of examples/gps.n3 from the Eyeling example
suite.

The original N3 example is an ARC-style, goal-driven route planner for a tiny
western-Belgium map. It starts from Gent, derives possible paths to Oostende,
compares their computed metrics, recommends the better route, renders a small
explanation, and checks that the recommendation is consistent with the route
metrics.

This is intentionally not a generic RDF/N3 reasoner. The concrete N3 facts
and rules are represented as Go structs and ordinary functions so the path
derivation and decision logic are easy to read and directly runnable.

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
