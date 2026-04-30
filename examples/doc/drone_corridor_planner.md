# Drone Corridor Planner

`drone_corridor_planner` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on constrained drone route planning through corridors and restricted zones. Its input fixture is organized around `CaseName`, `Question`, `Start`, `GoalLocation`, `Fuel`, `Thresholds`, `Actions`, `Expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: selected plan : fly_gent_brugge -> public_coastline_brugge_oostende duration : 2800 s cost : 0.012

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - 14 corridor actions were loaded from JSON C2 OK - bounded search found 17 plans meeting belief and cost thresholds C3 OK - lowest-cost selected plan starts with fly_gent_brugge

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/drone_corridor_planner.json](../input/drone_corridor_planner.json)

Go translation: [../drone_corridor_planner.go](../drone_corridor_planner.go)

Expected Markdown output: [../output/drone_corridor_planner.md](../output/drone_corridor_planner.md)
