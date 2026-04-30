# Bayes Therapy Decision Support

`bayes_therapy` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on posterior-weighted utility selection of the best therapy. Its input fixture is organized around `Diseases`, `Therapies`, `ProbGiven`, `Evidence`, `BenefitWeight`, `HarmWeight`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The recommended therapy is Paxlovid (utility = 3.585174). Full posterior distribution: COVID19 posterior = 0.483883 (unnormalized = 0.00928200)

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - all prior probabilities are in [0,1]. C2 OK - all conditional probabilities are in [0,1]. C3 OK - all adverse probabilities are in [0,1].

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/bayes_therapy.json](../input/bayes_therapy.json)

Go translation: [../bayes_therapy.go](../bayes_therapy.go)

Expected Markdown output: [../output/bayes_therapy.md](../output/bayes_therapy.md)
