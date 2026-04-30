# Fibonacci Example (Big)

`fibonacci` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on exact computation of a large Fibonacci number. Its input fixture is organized around `0`, `1`, `10`, `100`, `1000`, `10000`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The Fibonacci number for index 10000 is:...

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - base cases F(0)=0 and F(1)=1 hold. C2 OK - recurrence holds for all computed steps. C3 OK - all requested Fibonacci numbers match expected values:

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/fibonacci.json](../input/fibonacci.json)

Go translation: [../fibonacci.go](../fibonacci.go)

Expected Markdown output: [../output/fibonacci.md](../output/fibonacci.md)
