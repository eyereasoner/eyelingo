# French Cities

`french_cities` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on reachability over a small French city route graph. Its input fixture is organized around `Edges`, `Labels`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Technology** example. It demonstrates data representation, interoperability, policies, and computational artifacts in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Four cities in this small network can reach Nantes: Paris, Chartres, Le Mans, Angers.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - Angers has a direct one-way connection to Nantes. C2 OK - Le Mans reaches Nantes by chaining Le Mans → Angers → Nantes. C3 OK - Chartres reaches Nantes by chaining Chartres → Le Mans → Angers → Nantes.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/french_cities.json](../input/french_cities.json)

Go translation: [../french_cities.go](../french_cities.go)

Expected Markdown output: [../output/french_cities.md](../output/french_cities.md)
