# Complex Numbers

`complex_numbers` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on complex arithmetic and transcendental identity checks. Its input fixture is organized around `Question`, `Exponents`, `Inverses`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The complex.n3 test query derives 6 complex-number facts. Computed values: C1 sqrt(-1+0i) = 0 + 1i

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - N3 dial rules assign the expected polar angles for -1, e, and i. C2 OK - all four complex exponentiation answers match the complex.n3 test facts. C3 OK - i^i and e^(-pi/2) derive the same real value.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/complex_numbers.json](../input/complex_numbers.json)

Go translation: [../complex_numbers.go](../complex_numbers.go)

Expected Markdown output: [../output/complex_numbers.md](../output/complex_numbers.md)
