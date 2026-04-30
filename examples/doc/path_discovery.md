# Path Discovery

`path_discovery` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on airport path discovery with stopover and routing constraints. Its input fixture is organized around `Question`, `SourceID`, `DestinationID`, `Labels`, `Edges`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Technology** example. It demonstrates data representation, interoperability, policies, and computational artifacts in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The path discovery query finds 3 air routes with at most 2 stopovers. from : Ostend-Bruges International Airport to : Václav Havel Airport Prague

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - source and destination airport labels are known. C2 OK - Ostend-Bruges has one outbound route in the full N3 graph, to Liège Airport. C3 OK - the discovered route set matches the N3 query answer.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/path_discovery.json](../input/path_discovery.json)

Go translation: [../path_discovery.go](../path_discovery.go)

Expected Markdown output: [../output/path_discovery.md](../output/path_discovery.md)
