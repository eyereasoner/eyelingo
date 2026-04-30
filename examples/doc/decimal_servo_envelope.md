# Decimal Servo Envelope

`decimal_servo_envelope` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on certified servo pole interval and settling-step envelope. Its input fixture is organized around `caseName`, `question`, `samplePeriod`, `timeConstant`, `exactPoleSymbol`, `poleLower`, `poleUpper`, `initialAbsError`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: exact pole symbol : exp(-1/3) certified pole interval : [0.7165313105, 0.7165313106] first settled step : 10

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - pole certificate is nonempty, positive, and below 1 C2 OK - upper envelope strictly decreases at every sampled step C3 OK - step 9 is not yet below tolerance

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/decimal_servo_envelope.json](../input/decimal_servo_envelope.json)

Go translation: [../decimal_servo_envelope.go](../decimal_servo_envelope.go)

Expected Markdown output: [../output/decimal_servo_envelope.md](../output/decimal_servo_envelope.md)
