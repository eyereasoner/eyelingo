# Kaprekar 6174

`kaprekar_6174` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on kaprekar chains and basin facts ending at 6174. Its input fixture is organized around `StartCount`, `TargetConstant`, `ZeroBasin`, `MaxKaprekarSteps`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Kaprekar chains that end at 6174 are emitted as :kaprekar facts. total emitted : 9990 omitted 0000 basin : 10

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - all digit patterns from 0000 through 9999 were considered C2 OK - the identity-based step matches direct digit sorting for every start C3 OK - 3524 follows the classic 3087 -> 8352 -> 6174 chain

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/kaprekar_6174.json](../input/kaprekar_6174.json)

Go translation: [../kaprekar_6174.go](../kaprekar_6174.go)

Expected Markdown output: [../output/kaprekar_6174.md](../output/kaprekar_6174.md)
