# FFT8 Numeric

`fft8_numeric` is a Go translation/adaptation of Eyeling's `fft8-numeric.n3`.

## Background

A discrete Fourier transform converts a short sampled signal into frequency-domain coefficients. An 8-point transform is small enough to inspect manually but still shows key properties such as dominant bins, conjugate symmetry for real signals, and energy preservation between time and frequency domains.

## What it demonstrates

**Category:** Technology. Eight-point Fourier transform over a sampled sine wave with conjugate-bin and energy checks.

## How the Go implementation works

The implementation builds an 8-sample sine-wave fixture and computes its discrete Fourier transform directly. It records each complex bin, identifies the dominant frequency pair, and checks that the DC component and conjugate bins behave as expected.

Energy preservation is checked by comparing time-domain and frequency-domain totals within tolerance.

## Files

Input JSON: [../input/fft8_numeric.json](../input/fft8_numeric.json)

Go implementation: [../fft8_numeric.go](../fft8_numeric.go)

Expected Markdown output: [../output/fft8_numeric.md](../output/fft8_numeric.md)
