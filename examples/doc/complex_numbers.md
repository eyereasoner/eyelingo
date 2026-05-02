# Complex Numbers

`complex_numbers` is a Go translation/adaptation of Eyeling's `complex.n3`.

## Background

Complex numbers extend real numbers with an imaginary component, written a + bi where i² = -1. They can also be represented in polar form by magnitude and angle, which makes exponentials, logarithms, powers, and rotations easier to express. Correct handling of the quadrant of the angle is important because several complex values can share related magnitudes or trigonometric components.

## What it demonstrates

**Category:** Mathematics. Complex arithmetic and transcendental identity checks.

## How the Go implementation works

The Go code spells out the complex-number derivations instead of delegating everything to high-level helpers. It computes polar conversion, quadrant-sensitive arguments, exponentials, logarithms, powers, and inverse trigonometric forms so the report can show the mathematical path.

The checks compare the derived values with expected identities and tolerances, including cases where branch choice matters.

## Files

Input JSON: [../input/complex_numbers.json](../input/complex_numbers.json)

Go implementation: [../complex_numbers.go](../complex_numbers.go)

Python check: [../checks/complex_numbers.py](../checks/complex_numbers.py)

Expected Markdown output: [../output/complex_numbers.md](../output/complex_numbers.md)
