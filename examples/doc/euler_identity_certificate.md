# Euler Identity Certificate

`euler_identity_certificate` is a Go translation/adaptation of Eyeling's `euler-identity-certificate.n3`.

## Background

Euler's identity connects five fundamental constants through complex exponentiation: e, i, π, 1, and 0. Numerically, certifying the identity means evaluating exp(iπ) + 1 and showing that the remaining real and imaginary parts are within a chosen tolerance. The challenge is not the formula itself but presenting enough precision evidence that the near-zero result is credible.

## What it demonstrates

**Category:** Mathematics. High-precision certificate for the identity exp(iπ) + 1 = 0.

## How the Go implementation works

The code evaluates the finite residual for `exp(i*pi)+1` using the configured precision and tolerance. It separates the computed complex residual, its magnitude, and the pass/fail certificate so the numerical claim is auditable.

Checks verify that the residual is within tolerance and that the reported certificate matches the computed value.

## Files

Input JSON: [../input/euler_identity_certificate.json](../input/euler_identity_certificate.json)

Go implementation: [../euler_identity_certificate.go](../euler_identity_certificate.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/euler_identity_certificate.md](../output/euler_identity_certificate.md)
