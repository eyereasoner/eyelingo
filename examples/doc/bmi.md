# BMI — Body Mass Index

`bmi` is a Go translation/adaptation of Eyeling's `bmi.n3`.

## Background

Body Mass Index is defined as body mass in kilograms divided by height in meters squared. Adult BMI categories are threshold intervals, and a healthy-weight interval for a fixed height can be obtained by reversing the same formula. This makes BMI a compact example of unit normalization, arithmetic derivation, categorical classification, and boundary checks.

## What it demonstrates

**Category:** Science. Adult BMI calculation, category assignment, and healthy-weight interval.

## How the Go implementation works

The code normalizes the input height and weight units, computes BMI, assigns the adult category, and derives the healthy-weight interval for the same height. The calculation is intentionally small and direct, with the intermediate values kept available for the report.

Independent checks cover the numeric result, category threshold, unit handling, and interval bounds.

## Files

Input JSON: [../input/bmi.json](../input/bmi.json)

Go implementation: [../bmi.go](../bmi.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/bmi.md](../output/bmi.md)
