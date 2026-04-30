# Control System — ARC-style control-system example

`control_system` is a Go translation/adaptation of Eyeling's `control-system.n3`, with the rule structure kept close to EYE's `reasoning/control-system/rules-001.n3`.

The context is rule-based control. Measurements, actuator states, disturbances, and input conditions are used to derive safe control conclusions while keeping the original forward/backward rule shape visible.

## What it demonstrates

This example is mainly in the **Engineering** category. Translated measurement and control rules for actuators, inputs, and disturbances.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/control_system.json](../input/control_system.json)

Go translation: [../control_system.go](../control_system.go)

Expected Markdown output: [../output/control_system.md](../output/control_system.md)
