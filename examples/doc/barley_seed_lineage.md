# Barley Seed Lineage

`barley_seed_lineage` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on seed-lineage CAN/CAN'T reasoning for reproduction, dormancy, variation, and persistence. Its input fixture is organized around `caseName`, `question`, `world`, `lineages`, `expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: evolvable lineage : mainLine blocked contrast lineages : analogLine, fragileLine, coatlessLine, staticLine mainLine CAN : genome-copy, protected-dormancy, germination, propagule-production, accurate-self-reproduction, lineage-closure, adaptive-persistence

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - no-design laws are loaded C2 OK - mainLine can copy its digitally instantiated genome C3 OK - mainLine has repair, protected dormancy, and greenhouse support

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/barley_seed_lineage.json](../input/barley_seed_lineage.json)

Go translation: [../barley_seed_lineage.go](../barley_seed_lineage.go)

Expected Markdown output: [../output/barley_seed_lineage.md](../output/barley_seed_lineage.md)
