# Delfour

`delfour` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on privacy-preserving retail insight and recommendation policy. Its input fixture is organized around `Case`, `Products`, `Household`, `Scan`, `Insight`, `Policy`, `Signature`, `ReasonText`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Technology** example. It demonstrates data representation, interoperability, policies, and computational artifacts in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The scanner is allowed to use a neutral shopping insight and recommends Low-Sugar Tea Biscuits instead of Classic Tea Biscuits. case : delfour decision : Allowed

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: signature verifies : yes payload hash matches : yes minimization strips sensitive terms: yes

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/delfour.json](../input/delfour.json)

Go translation: [../delfour.go](../delfour.go)

Expected Markdown output: [../output/delfour.md](../output/delfour.md)
