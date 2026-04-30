# E-Bike Motor Thermal Envelope

`ebike_motor_thermal_envelope` is a Go translation/adaptation of Eyeling's `decimal-ebike-motor-thermal-envelope.n3`.

The context is embedded mobility safety. A planned e-bike assist profile is checked against a decimal motor-temperature envelope, making the safety margin visible.

## What it demonstrates

This example is mainly in the **Science** category. Certified e-bike motor-temperature envelope for an assist plan.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/ebike_motor_thermal_envelope.json](../input/ebike_motor_thermal_envelope.json)

Go translation: [../ebike_motor_thermal_envelope.go](../ebike_motor_thermal_envelope.go)

Expected Markdown output: [../output/ebike_motor_thermal_envelope.md](../output/ebike_motor_thermal_envelope.md)
