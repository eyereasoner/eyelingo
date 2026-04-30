# Digital Product Passport

`digital_product_passport` is a Go translation/adaptation of Eyeling's `digital-product-passport.n3`.

The context is product sustainability reporting. Component facts are rolled up into recycled-content, footprint, repairability, and critical-material conclusions similar to a digital product passport.

## What it demonstrates

This example is mainly in the **Science** category. Component roll-up for recycled content, carbon footprint, repairability, and critical materials.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/digital_product_passport.json](../input/digital_product_passport.json)

Go translation: [../digital_product_passport.go](../digital_product_passport.go)

Expected Markdown output: [../output/digital_product_passport.md](../output/digital_product_passport.md)
