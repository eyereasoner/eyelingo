# Barley Seed Lineage

`barley_seed_lineage` is a Go translation/adaptation of Eyeling's `act-barley-seed-lineage.n3`.

The context is biological lineage and constructor-theory style CAN/CAN'T reasoning. The scenario distinguishes viable reproduction, dormant preservation, and heritable variation from transformations that the supplied resources or laws do not permit.

## What it demonstrates

This example is mainly in the **Science** category. Seed-lineage CAN/CAN'T reasoning for reproduction, dormancy, variation, and persistence.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/barley_seed_lineage.json](../input/barley_seed_lineage.json)

Go translation: [../barley_seed_lineage.go](../barley_seed_lineage.go)

Expected Markdown output: [../output/barley_seed_lineage.md](../output/barley_seed_lineage.md)
