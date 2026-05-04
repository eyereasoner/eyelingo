# BMI — Body Mass Index Example

## What this example is about

This is a small unit-normalization and category-boundary example. It takes a person's height and weight in metric units, computes Body Mass Index, and reports the adult BMI category.

It is intended to demonstrate transparent calculation and boundary checks. It is not a health assessment.

## How it works, in plain language

The input gives weight in kilograms and height in centimeters. The program converts centimeters to meters, then applies the standard BMI calculation: weight divided by height squared.

After computing the value, it maps the number to a category: underweight, normal, overweight, or obese. It also calculates the weight range that would correspond to a normal BMI range at the same height.

## What to notice in the output

The output does more than print a BMI number. It also shows the category and a height-specific healthy-weight range. That makes the answer easier to interpret for readers who do not remember the category thresholds.

## What the trust gate checks

The trust gate verifies that height is positive, the computed category is normal for this input, the rounded BMI is stable, and the lower and upper weight-range bounds are stable. This catches mistakes in unit conversion, rounding, or threshold handling.

## Run it

From the repository root:

```sh
node examples/bmi.js
```

## Files

- [JavaScript example](../bmi.js)
- [Input data](../input/bmi.json)
- [Expected output](../output/bmi.md)
