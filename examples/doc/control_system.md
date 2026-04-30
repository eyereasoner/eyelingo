# Control System — ARC-style control-system example

`control_system` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on translated measurement and control rules for actuators, inputs, and disturbances. Its input fixture is organized around `input1`, `input2`, `input3`, `disturbance1`, `disturbance2`, `state1`, `state2`, `state3`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: actuator1 control1 = 39.273462 actuator2 control1 = 26.080000 input1 measurement10 = 2.236068 (lessThan branch)

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - input1 measurement10 follows the lessThan backward rule. C2 OK - disturbance2 measurement10 follows the notLessThan backward rule. C3 OK - input2 boolean guard is true for the feedforward rule.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/control_system.json](../input/control_system.json)

Go translation: [../control_system.go](../control_system.go)

Expected Markdown output: [../output/control_system.md](../output/control_system.md)
