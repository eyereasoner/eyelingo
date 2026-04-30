# 8-Queens

`queens` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on 8-Queens constraint satisfaction with a valid board solution. Its input fixture is organized around `N`, `MaxPrint`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Solving 8-Queens... Printing at most 1 solution(s). Solution 1:

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - search reached depth 8. C2 OK - first solution places one queen in each row. C3 OK - first solution columns are unique.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/queens.json](../input/queens.json)

Go translation: [../queens.go](../queens.go)

Expected Markdown output: [../output/queens.md](../output/queens.md)
