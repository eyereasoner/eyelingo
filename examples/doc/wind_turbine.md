# Wind Turbine Envelope

`wind_turbine` is a Go translation/adaptation of Eyeling's `wind-turbine.n3`.

The context is renewable-energy engineering. Wind-speed intervals are classified and passed through a cubic power curve to estimate energy and operating-envelope status.

## How it works

Inspired by Eyeling's `examples/wind-turbine.n3`.

The example classifies wind-speed samples for a turbine power curve and
computes the certified energy contribution of the usable intervals.

## What it demonstrates

This example is mainly in the **Engineering** category. Wind-speed envelope classification with cubic power curve and interval energy audit.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/wind_turbine.json](../input/wind_turbine.json)

Go translation: [../wind_turbine.go](../wind_turbine.go)

Expected Markdown output: [../output/wind_turbine.md](../output/wind_turbine.md)
