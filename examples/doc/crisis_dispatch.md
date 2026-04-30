# Crisis Dispatch

`crisis_dispatch` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on storm-incident resource assignment optimized for triage score and travel time. Its input fixture is organized around `CaseName`, `Question`, `Locations`, `Roads`, `Responders`, `Incidents`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The exact dispatch plan serves all 6 storm incidents with triage score 475. case : storm-response finish minute : 30

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - every incident is served by exactly one selected route. C2 OK - selected responders have the required capabilities for every stop. C3 OK - all service finishes meet incident deadlines and responder shifts.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/crisis_dispatch.json](../input/crisis_dispatch.json)

Go translation: [../crisis_dispatch.go](../crisis_dispatch.go)

Expected Markdown output: [../output/crisis_dispatch.md](../output/crisis_dispatch.md)
