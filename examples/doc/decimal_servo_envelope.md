# Decimal Servo Envelope

`decimal_servo_envelope` is a Go translation/adaptation of Eyeling's `decimal-transcendental-servo-envelope.n3`.

The context is certified servo behavior. Decimal intervals bound a pole and settling envelope, which lets the example express engineering guarantees without relying only on approximate simulation.

## What it demonstrates

This example is mainly in the **Engineering** category. Certified servo pole interval and settling-step envelope.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/decimal_servo_envelope.json](../input/decimal_servo_envelope.json)

Go translation: [../decimal_servo_envelope.go](../decimal_servo_envelope.go)

Expected Markdown output: [../output/decimal_servo_envelope.md](../output/decimal_servo_envelope.md)
