# French Cities

`french_cities` is a Go translation/adaptation of Eyeling's `french-cities.n3`.

The context is small graph reachability. A route network between French cities is explored to derive which destinations are reachable under the supplied links.

## How it works

A self-contained Go translation of examples/french-cities.n3 from the Eyeling
example suite, in ARC style.

The original N3 program encodes a small graph of French cities connected by
one‑way roads.  It uses RDFS/OWL rules to derive longer paths from shorter
ones and answers the question: which cities can reach Nantes?

This is intentionally not a full N3 reasoner – it is a concrete scenario that
mirrors the structure of the original N3 example.

## What it demonstrates

This example is mainly in the **Technology** category. Reachability over a small French city route graph.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/french_cities.json](../input/french_cities.json)

Go translation: [../french_cities.go](../french_cities.go)

Expected Markdown output: [../output/french_cities.md](../output/french_cities.md)
