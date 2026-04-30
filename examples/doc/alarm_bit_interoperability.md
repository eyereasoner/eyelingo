# Alarm Bit Interoperability

`alarm_bit_interoperability` is a Go translation/adaptation of Eyeling's `act-alarm-bit-interoperability.n3`.

## Background

A classical bit can be copied, relayed, and fanned out without changing its value, which is why ordinary alarm signals can be distributed to multiple recipients. Quantum states are different: an unknown state cannot be cloned perfectly. This example uses that contrast to separate safe classical alarm-bit interoperability from protected-token cases where unrestricted copying would be invalid.

## What it demonstrates

**Category:** Technology. Classical alarm-bit copy and relay tasks contrasted with forbidden universal cloning.

## How the Go implementation works

The program loads the classical media and protected-token fixture into structs, then derives which copy and relay tasks are valid for ordinary alarm bits. It keeps the restricted token in a separate branch so prohibited cloning or fan-out cases are recorded as impossible tasks.

The checks compare the derived can/can't sets with the expected decisions, making the boundary between ordinary copying and protected evidence explicit.

## Files

Input JSON: [../input/alarm_bit_interoperability.json](../input/alarm_bit_interoperability.json)

Go implementation: [../alarm_bit_interoperability.go](../alarm_bit_interoperability.go)

Expected Markdown output: [../output/alarm_bit_interoperability.md](../output/alarm_bit_interoperability.md)
