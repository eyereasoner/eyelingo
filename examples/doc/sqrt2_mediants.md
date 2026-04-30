# Integer-First Sqrt2 Mediants

`sqrt2_mediants` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on rational lower/upper bounds for sqrt(2) certified by integer square comparisons. Its input fixture is organized around `caseName`, `question`, `maxDenominator`, `expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: lower bound : 1393/985 = 1.414213197970 upper bound : 577/408 = 1.414215686275 certified interval width : 0.000002488305

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - nine convergents stay within the denominator limit C2 OK - the best lower bound is 1393/985 C3 OK - the best upper bound is 577/408

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/sqrt2_mediants.json](../input/sqrt2_mediants.json)

Go translation: [../sqrt2_mediants.go](../sqrt2_mediants.go)

Expected Markdown output: [../output/sqrt2_mediants.md](../output/sqrt2_mediants.md)
