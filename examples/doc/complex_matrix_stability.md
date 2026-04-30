# Complex Matrix Stability

`complex_matrix_stability` is a Go translation/adaptation of Eyeling's `complex-matrix-stability.n3`.

The example classifies three diagonal complex matrices by spectral radius for discrete-time dynamics. It also audits complex multiplication and the expected scaling behaviour of the spectral radius squared.

## What it demonstrates

This is mainly a **Engineering** example. The key idea is that diagonal entries are eigenvalues, so the largest eigenvalue modulus determines whether modes grow, remain marginal, or decay.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, counters, thresholds, derived facts, or platform details.

## Files

Input JSON: [../input/complex_matrix_stability.json](../input/complex_matrix_stability.json)

Go translation: [../complex_matrix_stability.go](../complex_matrix_stability.go)

Expected Markdown output: [../output/complex_matrix_stability.md](../output/complex_matrix_stability.md)
