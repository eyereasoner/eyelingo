# Bayes Diagnosis

`bayes_diagnosis` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on bayesian posterior ranking of possible diseases from symptoms and test evidence. Its input fixture is organized around `Diseases`, `ProbGiven`, `Evidence`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The most likely disease is COVID19 (posterior = 0.941209). Full posterior distribution: COVID19 posterior = 0.941209 (unnormalized = 0.00154700)

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - all prior probabilities are in [0,1]. C2 OK - all conditional probabilities are in [0,1].

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/bayes_diagnosis.json](../input/bayes_diagnosis.json)

Go translation: [../bayes_diagnosis.go](../bayes_diagnosis.go)

Expected Markdown output: [../output/bayes_diagnosis.md](../output/bayes_diagnosis.md)
