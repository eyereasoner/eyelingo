# GPS — Goal driven route planning

`gps` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on goal-driven route planning over a small road network. Its input fixture is organized around `Traveller`, `Question`, `Routes`, `Edges`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Take the direct route via Brugge. Recommended route: Gent → Brugge → Oostende

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - the direct Gent → Brugge → Oostende route was derived. C2 OK - the alternative Gent → Kortrijk → Brugge → Oostende route was derived. C3 OK - the recommended route is faster than the alternative.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/gps.json](../input/gps.json)

Go translation: [../gps.go](../gps.go)

Expected Markdown output: [../output/gps.md](../output/gps.md)
