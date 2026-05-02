# Digital Product Passport

`digital_product_passport` is a Go translation/adaptation of Eyeling's `digital-product-passport.n3`.

## Background

A digital product passport summarizes lifecycle information about a product, such as material composition, recycled content, carbon footprint, critical materials, and repairability. Many of those values are component-level facts that must be rolled up into product-level indicators. The example focuses on the aggregation and threshold checks needed to turn component data into a concise passport assessment.

## What it demonstrates

**Category:** Science. Component roll-up for recycled content, carbon footprint, repairability, and critical materials.

## How the Go implementation works

The implementation loads component, material, lifecycle, repair, and footprint data, then rolls those facts up into product-level indicators. It computes total mass, recycled-content share, carbon footprint, critical-material exposure, and repairability hints from the component table.

Checks compare the roll-up totals and threshold classifications with the expected passport summary.

## Files

Input JSON: [../input/digital_product_passport.json](../input/digital_product_passport.json)

Go implementation: [../digital_product_passport.go](../digital_product_passport.go)

Python check: [../checks/digital_product_passport.py](../checks/digital_product_passport.py)

Expected Markdown output: [../output/digital_product_passport.md](../output/digital_product_passport.md)
