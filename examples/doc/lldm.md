# LLD — Leg Length Discrepancy Measurement

`lldm` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on leg-length discrepancy measurement and alarm thresholding. Its input fixture is organized around `P1x`, `P1y`, `P2x`, `P2y`, `P3x`, `P3y`, `P4x`, `P4y`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: LLD Alarm = TRUE (discrepancy dCm = -1.908234, threshold ±1.25) Key computed values: SL1 = -0.062857 SL3 = SL4 = 15.909091

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - L1 is perpendicular to L3 and L4 (slopes product ≈ -1). C2 OK - p5 lies on both L1 and L3. C3 OK - p6 lies on both L1 and L4.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/lldm.json](../input/lldm.json)

Go translation: [../lldm.go](../lldm.go)

Expected Markdown output: [../output/lldm.md](../output/lldm.md)
