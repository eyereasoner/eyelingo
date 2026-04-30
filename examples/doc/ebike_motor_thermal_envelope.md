# E-Bike Motor Thermal Envelope

`ebike_motor_thermal_envelope` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on certified e-bike motor-temperature envelope for an assist plan. Its input fixture is organized around `CaseName`, `Question`, `SamplePeriodSec`, `ThermalTimeConstantSec`, `ExactCoolingSymbol`, `CoolingLower`, `CoolingUpper`, `AmbientC`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: decision : ThermallySafeForThisAssistPlan cooling certificate : exp(-1/4) in 0.7788007830 .. 0.7788007831 maximum upper motor temperature : 40.2285 C

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - cooling interval 0.7788007830..0.7788007831 is positive, ordered, and contractive C2 OK - temperature trace has 13 samples including the initial state C3 OK - maximum upper temperature 40.2285 C stays below hard limit 45.0 C

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/ebike_motor_thermal_envelope.json](../input/ebike_motor_thermal_envelope.json)

Go translation: [../ebike_motor_thermal_envelope.go](../ebike_motor_thermal_envelope.go)

Expected Markdown output: [../output/ebike_motor_thermal_envelope.md](../output/ebike_motor_thermal_envelope.md)
