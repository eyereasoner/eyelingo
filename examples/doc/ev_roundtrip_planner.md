# EV Roadtrip Planner

`ev_roundtrip_planner` is a Go translation/adaptation of Eyeling's `ev-roundtrip-planner.n3`.

The context is electric-vehicle trip planning. Route choices combine battery state, charging, duration, cost, and comfort thresholds so the answer reflects several operational constraints at once.

## What it demonstrates

This example is mainly in the **Engineering** category. EV route planning with battery, duration, cost, and comfort constraints.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/ev_roundtrip_planner.json](../input/ev_roundtrip_planner.json)

Go translation: [../ev_roundtrip_planner.go](../ev_roundtrip_planner.go)

Expected Markdown output: [../output/ev_roundtrip_planner.md](../output/ev_roundtrip_planner.md)
