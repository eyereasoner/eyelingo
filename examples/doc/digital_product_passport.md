# Digital Product Passport

`digital_product_passport` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on component roll-up for recycled content, carbon footprint, repairability, and critical materials. Its input fixture is organized around `CaseName`, `Question`, `Organization`, `Site`, `Repairer`, `Passport`, `AccessPolicy`, `Product`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Passport decision : PASS for ACME X1000 SN123. recycled content : 13% lifecycle footprint : 52500 gCO2e

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - component masses fold to 105 g C2 OK - recycled component masses fold to 14 g C3 OK - integer recycled-content percentage is 13%

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/digital_product_passport.json](../input/digital_product_passport.json)

Go translation: [../digital_product_passport.go](../digital_product_passport.go)

Expected Markdown output: [../output/digital_product_passport.md](../output/digital_product_passport.md)
