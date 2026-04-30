# Brussels Beijing

`brussels_beijing` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on cost-optimized flight routing from Brussels to Beijing while avoiding a carrier. Its input fixture is organized around `Flights`, `Labels`, `StartCity`, `EndCity`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Cheapest route from Brussels to Beijing (avoiding Turkish Airlines) costs €670: Brussels → Frankfurt (Lufthansa, €150), Frankfurt → Beijing (Lufthansa, €520)

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - a path exists between Brussels and Beijing. C2 OK - the total cost matches the computed optimal cost. C3 OK - the chosen path is cheaper than the direct flight (€800).

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/brussels_beijing.json](../input/brussels_beijing.json)

Go translation: [../brussels_beijing.go](../brussels_beijing.go)

Expected Markdown output: [../output/brussels_beijing.md](../output/brussels_beijing.md)
