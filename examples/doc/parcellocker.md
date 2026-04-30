# Parcel Locker

`parcellocker` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on delegated parcel pickup-token authorization. Its input fixture is organized around `CaseName`, `Question`, `People`, `Parcel`, `Locker`, `Authorization`, `Requests`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Technology** example. It demonstrates data representation, interoperability, policies, and computational artifacts in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: May Noah use Maya's one-time pickup token to collect parcel123 from locker B17? decision : PERMIT release : Noah may collect parcel123 for Maya from locker B17 at Station West

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - the source pickup request satisfies all ten authorization conditions C2 OK - the same token is denied for billing, redirect, wrong person, wrong locker, and reuse C3 OK - every request matches its expected PERMIT or DENY result

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/parcellocker.json](../input/parcellocker.json)

Go translation: [../parcellocker.go](../parcellocker.go)

Expected Markdown output: [../output/parcellocker.md](../output/parcellocker.md)
