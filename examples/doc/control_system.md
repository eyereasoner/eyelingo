# Control System

`control_system` is a Go translation/adaptation of Eyeling's `control-system.n3`.

## Background

A control system links measurements, inferred state, and actuator decisions. Forward rules derive consequences from known facts, while goal-directed or backward-style logic asks which supporting measurements would justify a conclusion. This example keeps that structure explicit for a small rule base involving measured quantities, disturbances, actuator states, and outputs.

## What it demonstrates

**Category:** Engineering. Measurement and control rules for actuators, inputs, and disturbances.

## How the Go implementation works

The implementation loads measurement, observation, actuator, input, disturbance, state, and output facts into typed slices. Forward rules derive actuator and output conclusions, while backward-style helper logic derives compound measurements from paired measurements.

The report separates derived facts from checks, so the control decision and the rule evidence can be audited independently.

## Files

Input JSON: [../input/control_system.json](../input/control_system.json)

Go implementation: [../control_system.go](../control_system.go)

Expected Markdown output: [../output/control_system.md](../output/control_system.md)
