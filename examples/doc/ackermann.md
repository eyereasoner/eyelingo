# Ackermann

`ackermann` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on exact Ackermann and hyperoperation facts, including very large integer results. Its input fixture is a list with 12 entries.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The ackermann.n3 test query derives 12 Ackermann facts. Computed values: A0 ackermann(0,0) = 1

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - x=0 reduces to successor after the y+3 binary offset. C2 OK - x=1 reduces to addition after the y+3 binary offset. C3 OK - x=2 reduces to multiplication after the y+3 binary offset.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/ackermann.json](../input/ackermann.json)

Go translation: [../ackermann.go](../ackermann.go)

Expected Markdown output: [../output/ackermann.md](../output/ackermann.md)
