# Complex Matrix Stability

`complex_matrix_stability` is a Go translation/adaptation of Eyeling's `complex-matrix-stability.n3`.

## Background

A discrete-time linear system repeatedly applies a matrix to the current state. For a diagonal complex matrix, the size of each diagonal entry determines whether that component shrinks, stays bounded, or grows. The spectral radius—the largest magnitude of any eigenvalue—is the key stability value: values below the chosen threshold indicate bounded decay, while values above it indicate growth.

## What it demonstrates

**Category:** Engineering. Discrete-time stability classification using spectral radii of diagonal complex matrices.

## How the Go implementation works

The program reads diagonal complex entries, computes their squared magnitudes, and classifies the discrete-time system from the resulting spectral radius. The stability decision is made by comparing the largest magnitude against the fixture threshold.

Checks verify the dominant entry, the radius calculation, and the final bounded/unbounded classification.

## Files

Input JSON: [../input/complex_matrix_stability.json](../input/complex_matrix_stability.json)

Go implementation: [../complex_matrix_stability.go](../complex_matrix_stability.go)

Expected Markdown output: [../output/complex_matrix_stability.md](../output/complex_matrix_stability.md)
