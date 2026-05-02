# Wind Turbine Envelope

`wind_turbine` is a Go translation/adaptation of Eyeling's `wind-turbine.n3`.

## Background

Wind-turbine output depends strongly on wind speed. Below cut-in speed the turbine produces no useful power, between cut-in and rated speed power is often modeled with a cubic curve, and above cut-out speed the turbine shuts down for safety. This example classifies sampled wind speeds and integrates the accepted power contributions.

## What it demonstrates

**Category:** Engineering. Wind-speed envelope classification with cubic power curve and interval energy audit.

## How the Go implementation works

The implementation classifies wind-speed samples against cut-in, rated, and cut-out regions, then applies the cubic power curve for usable intervals. It integrates the interval contributions into an energy estimate and tracks operating-envelope status.

Checks verify the speed classifications, power estimates, and total certified energy.

## Files

Input JSON: [../input/wind_turbine.json](../input/wind_turbine.json)

Go implementation: [../wind_turbine.go](../wind_turbine.go)

Python check: [../checks/wind_turbine.py](../checks/wind_turbine.py)

Expected Markdown output: [../output/wind_turbine.md](../output/wind_turbine.md)
