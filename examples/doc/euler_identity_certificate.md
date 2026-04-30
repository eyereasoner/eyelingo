# Euler Identity Certificate

`euler_identity_certificate` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on high-precision certificate for the identity exp(iπ) + 1 = 0. Its input fixture is organized around `caseName`, `question`, `angle`, `tolerance`, `terms`, `expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: expression : exp(iπ) + 1 terms used : 28 computed real part of exp(iπ) : -1.000000000000000

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - Taylor expansion used 28 terms from the JSON input C2 OK - computed real part is close to -1 C3 OK - computed imaginary part is close to 0

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/euler_identity_certificate.json](../input/euler_identity_certificate.json)

Go translation: [../euler_identity_certificate.go](../euler_identity_certificate.go)

Expected Markdown output: [../output/euler_identity_certificate.md](../output/euler_identity_certificate.md)
