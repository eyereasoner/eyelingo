# E-Bike Motor Thermal Envelope

`ebike_motor_thermal_envelope` is a Go translation/adaptation of Eyeling's `decimal-ebike-motor-thermal-envelope.n3`.

## Background

Electric motors heat up when they deliver torque, and excessive temperature can damage components or reduce lifetime. A thermal envelope describes acceptable temperature ranges over time, often with margins around warning and cutoff thresholds. This example samples an assist plan and checks whether the motor remains inside the certified envelope.

## What it demonstrates

**Category:** Science. Certified e-bike motor-temperature envelope for an assist plan.

## How the Go implementation works

The implementation propagates a certified decimal interval through a sampled motor-temperature model. Each assist segment updates the temperature envelope, tracks the safety margin, and records the worst-case sample.

Checks confirm that the plan stays below the thermal limit and that the interval arithmetic remains inside the declared tolerance.

## Files

Input JSON: [../input/ebike_motor_thermal_envelope.json](../input/ebike_motor_thermal_envelope.json)

Go implementation: [../ebike_motor_thermal_envelope.go](../ebike_motor_thermal_envelope.go)

Python check: [../checks/ebike_motor_thermal_envelope.py](../checks/ebike_motor_thermal_envelope.py)

Expected Markdown output: [../output/ebike_motor_thermal_envelope.md](../output/ebike_motor_thermal_envelope.md)
