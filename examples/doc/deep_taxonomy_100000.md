# Deep Taxonomy 100000

`deep_taxonomy_100000` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on large taxonomy materialization benchmark through a very deep class chain. Its input fixture is organized around `TaxonomyDepth`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Technology** example. It demonstrates data representation, interoperability, policies, and computational artifacts in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The deep taxonomy test succeeds. Starting fact : :ind a :N0 Reached class : :ind a :N100000

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - the starting classification N0 is present. C2 OK - the first expansion produced N1 together with side labels I1 and J1. C3 OK - the chain reaches the midpoint N50000 and still carries both side-label branches.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/deep_taxonomy_100000.json](../input/deep_taxonomy_100000.json)

Go translation: [../deep_taxonomy_100000.go](../deep_taxonomy_100000.go)

Expected Markdown output: [../output/deep_taxonomy_100000.md](../output/deep_taxonomy_100000.md)
