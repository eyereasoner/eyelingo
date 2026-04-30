# Allen Interval Calculus

`allen_interval_calculus` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on allen temporal interval relation closure over completed and explicit intervals. Its input fixture is organized around `caseName`, `question`, `intervals`, `expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: derived relations : 110 ordered interval pairs showcase : A before B | A meets C | A overlaps D | F starts A | G finishes A | A during H | A equals E | J meets I | K finishes C completed intervals : I=16:00-18:00, K=13:30-14:00

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - 11 intervals were loaded and completed when duration was present C2 OK - A before B was derived C3 OK - A meets C and C metBy A were derived

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/allen_interval_calculus.json](../input/allen_interval_calculus.json)

Go translation: [../allen_interval_calculus.go](../allen_interval_calculus.go)

Expected Markdown output: [../output/allen_interval_calculus.md](../output/allen_interval_calculus.md)
