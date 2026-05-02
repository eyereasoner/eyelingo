# RC Discharge Envelope

`rc_discharge_envelope` is a Go translation/adaptation of Eyeling's `rc-discharge-envelope.n3`.

## Background

An RC circuit discharges through a resistor with an exponential voltage decay. The time constant determines how quickly the voltage falls, and an envelope bounds the acceptable voltage range at sampled times. This example checks monotonic decay and identifies when the voltage has safely crossed below a threshold.

## What it demonstrates

**Category:** Science. Certified exponential decay envelope for an RC discharge.

## How the Go implementation works

The program propagates a floating-point decay interval through sampled RC-discharge times. It computes voltage bounds at each sample, finds the first safe sample, and records the margin against the threshold.

Checks verify monotonic decay, interval containment, and the selected safe point.

## Files

Input JSON: [../input/rc_discharge_envelope.json](../input/rc_discharge_envelope.json)

Go implementation: [../rc_discharge_envelope.go](../rc_discharge_envelope.go)

Python check: [../checks/rc_discharge_envelope.py](../checks/rc_discharge_envelope.py)

Expected Markdown output: [../output/rc_discharge_envelope.md](../output/rc_discharge_envelope.md)
