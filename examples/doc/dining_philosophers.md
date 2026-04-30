# Dining Philosophers

`dining_philosophers` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on chandy-Misra dining-philosophers trace with concurrency conflict checks. Its input fixture is a list with 9 entries.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The Chandy-Misra dining-philosophers trace completes without conflict. philosophers : 5 forks : 5

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - the translated run follows the nine start-of-round configs from the N3 source C2 OK - 26 dirty-fork requests transfer and the other 19 fork-round pairs are kept C3 OK - rounds 1/4/7 feed P1,P3; rounds 2/5/8 feed P2,P4; rounds 3/6/9 feed P5

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/dining_philosophers.json](../input/dining_philosophers.json)

Go translation: [../dining_philosophers.go](../dining_philosophers.go)

Expected Markdown output: [../output/dining_philosophers.md](../output/dining_philosophers.md)
