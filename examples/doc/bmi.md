# BMI — ARC-style Body Mass Index example

`bmi` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on adult BMI calculation, category assignment, and healthy-weight interval. Its input fixture is organized around `UnitSystem`, `Weight`, `Height`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: BMI = 22.7 Category = Normal At height 178 cm, a healthy-weight range is about 58.6–78.9 kg (BMI 18.5–24.9).

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - the input was normalized into positive SI values. C2 OK - height squared was reconstructed from the normalized height. C3 OK - the BMI value matches the BMI = kg / m² formula.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/bmi.json](../input/bmi.json)

Go translation: [../bmi.go](../bmi.go)

Expected Markdown output: [../output/bmi.md](../output/bmi.md)
