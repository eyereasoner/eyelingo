# AuroraCare

`auroracare` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on health-data permit/deny scenarios across care, quality improvement, and research. Its input fixture is organized around `CaseName`, `Question`, `Policies`, `Subjects`, `Requesters`, `Scenarios`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: For each AuroraCare scenario, should the PDP permit or deny the requested use of health data, and why? permit count : 4 deny count : 3

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - all seven scenarios match the PERMIT/DENY outcomes encoded in the N3 example C2 OK - primary-care access requires a clinician role and a care-team link C3 OK - quality improvement is allowed only when both lab results and patient summary are requested in the secure environment

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/auroracare.json](../input/auroracare.json)

Go translation: [../auroracare.go](../auroracare.go)

Expected Markdown output: [../output/auroracare.md](../output/auroracare.md)
