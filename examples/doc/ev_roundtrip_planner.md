# EV Roadtrip Planner

`ev_roundtrip_planner` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on eV route planning with battery, duration, cost, and comfort constraints. Its input fixture is organized around `CaseName`, `Question`, `Vehicle`, `Goal`, `FuelSteps`, `Thresholds`, `Actions`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Select plan : drive_bru_liege -> drive_liege_aachen -> shuttle_aachen_cologne. route result : Cologne battery=low pass=none duration : 210.0 minutes

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - 8 acceptable Brussels-to-Cologne plans were derived C2 OK - selected plan duration 210.0 is below 260.0 C3 OK - selected plan cost 0.054 is below 0.090

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/ev_roundtrip_planner.json](../input/ev_roundtrip_planner.json)

Go translation: [../ev_roundtrip_planner.go](../ev_roundtrip_planner.go)

Expected Markdown output: [../output/ev_roundtrip_planner.md](../output/ev_roundtrip_planner.md)
