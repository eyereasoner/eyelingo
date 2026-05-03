# FFT32 Numeric

`fft32_numeric` is a Go translation/adaptation of Eyeling's `fft32-numeric.n3`.

## Background

A discrete Fourier transform decomposes a finite sequence of samples into frequency bins. For real-valued signals, positive and negative frequency bins appear as conjugate pairs, and Parseval's theorem relates time-domain energy to frequency-domain energy. A 32-point transform gives enough bins to test several waveform shapes while keeping every bin visible in the output.

## What it demonstrates

**Category:** Technology. Thirty-two-point Fourier transform over several sampled waveforms with dominant-bin, flat-spectrum, conjugate-symmetry, and energy checks.

## How the Go implementation works

The Go code creates 32-sample waveform fixtures, computes a full discrete Fourier spectrum for each waveform, and stores magnitude and phase per bin. It then identifies dominant bins and checks special cases such as impulse flatness.

The audit also verifies conjugate symmetry and Parseval energy preservation within the configured tolerance.

## Files

Input JSON: [../input/fft32_numeric.json](../input/fft32_numeric.json)

Go implementation: [../fft32_numeric.go](../fft32_numeric.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/fft32_numeric.md](../output/fft32_numeric.md)
