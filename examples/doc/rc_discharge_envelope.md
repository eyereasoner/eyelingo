# RC Discharge Envelope

`rc_discharge_envelope` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on certified exponential decay envelope for an RC discharge. Its input fixture is organized around `caseName`, `question`, `samplePeriod`, `timeConstant`, `exactDecaySymbol`, `decayLower`, `decayUpper`, `initialVoltage`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: exact decay symbol : exp(-1/4) certified decay interval : [0.7788007830, 0.7788007831] first below tolerance step : 13

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - decay certificate is nonempty, positive, and below 1 C2 OK - voltage upper envelope decreases at every sample C3 OK - step 12 remains above the voltage tolerance

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/rc_discharge_envelope.json](../input/rc_discharge_envelope.json)

Go translation: [../rc_discharge_envelope.go](../rc_discharge_envelope.go)

Expected Markdown output: [../output/rc_discharge_envelope.md](../output/rc_discharge_envelope.md)
