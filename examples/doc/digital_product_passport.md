# Digital Product Passport

## What this example is about

This example builds a public product-passport summary for a smartphone. It combines component data, material data, public documents, restricted documents, lifecycle events, and carbon-footprint fields into one pass/fail passport decision.

The purpose is to show how a public explanation can reveal useful circularity information while keeping restricted documents scoped to the right audiences.

## How it works, in plain language

The program totals the mass of the listed components and separately totals the recycled mass. It then computes a recycled-content percentage. It also sums manufacturing, transport, and use-phase emissions to produce a lifecycle footprint.

For circularity, it checks whether the battery is replaceable and whether the public document set contains the needed repair and spare-parts information. For material exposure, it follows component-to-material links and reports which used materials are marked as critical raw materials.

## What to notice in the output

The output gives a passport `PASS`, recycled content, lifecycle footprint, total component mass, critical raw materials, a repair-friendly hint, and the public endpoint. It also prints a component roll-up and the public document list so the result can be inspected.

## What the trust gate checks

The trust gate verifies mass totals, recycled-content percentage, footprint sum, critical-material exposure, repair-friendly conditions, restricted-document scoping, that the product digital link equals the passport endpoint, and that lifecycle dates are chronological.

## Run it

From the repository root:

```sh
node examples/digital_product_passport.js
```

## Files

- [JavaScript example](../digital_product_passport.js)
- [Input data](../input/digital_product_passport.json)
- [Expected output](../output/digital_product_passport.md)
