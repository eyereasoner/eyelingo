# Gradient Descent Step

`gradient_descent_step` is a Go translation/adaptation of Eyeling's `gd-step-certified.n3`.

## Background

Gradient descent minimizes a function by moving from the current point in the negative direction of the gradient. For a quadratic objective, the gradient is linear and a single step can be checked exactly against the chosen learning rate. The example focuses on certifying one update rather than running a long optimization loop.

## What it demonstrates

**Category:** Mathematics. Certified single gradient-descent step for a quadratic objective.

## How the Go implementation works

The program evaluates one gradient-descent update for a quadratic objective. It computes the gradient at the starting point, applies the configured step size, and then compares objective values before and after the move.

The certificate checks the expected decrease condition, step bound, and final point against the fixture tolerance.

## Files

Input JSON: [../input/gradient_descent_step.json](../input/gradient_descent_step.json)

Go implementation: [../gradient_descent_step.go](../gradient_descent_step.go)

Expected Markdown output: [../output/gradient_descent_step.md](../output/gradient_descent_step.md)
