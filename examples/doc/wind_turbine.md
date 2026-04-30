# Wind Turbine Envelope

`wind_turbine` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on wind-speed envelope classification with cubic power curve and interval energy audit. Its input fixture is organized around `caseName`, `question`, `cutInMS`, `ratedMS`, `cutOutMS`, `ratedPowerMW`, `intervalMinutes`, `windSpeedsMS`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: operating thresholds : cut-in 3.5 m/s, rated 12.0 m/s, cut-out 25.0 m/s rated power : 3.2 MW interval classifications : t1 3.0 m/s stopped 0.000 MW; t2 6.5 m/s partial 0.440 MW; t3 11.2 m/s partial 2.586 MW; t4 15.0 m/s rated 3.200 MW; t5 24.5 m/s rated 3.200 MW; t6 27.0 m/s stopped 0.000 MW

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - cut-in, rated, and cut-out thresholds are ordered C2 OK - usable intervals are exactly the samples inside the operating envelope C3 OK - two intervals reach rated power

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/wind_turbine.json](../input/wind_turbine.json)

Go translation: [../wind_turbine.go](../wind_turbine.go)

Expected Markdown output: [../output/wind_turbine.md](../output/wind_turbine.md)
