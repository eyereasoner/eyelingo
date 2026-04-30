# Photosynthetic Exciton Transfer

`photosynthetic_exciton_transfer` is a Go translation/adaptation of Eyeling's `act-photosynthetic-exciton-transfer.n3`.

The context is quantum biology. Tuned excitonic coupling, vibronic bridging, dephasing, and downhill transfer are contrasted with a detuned complex that fails to deliver the excitation.

## What it demonstrates

This example is mainly in the **Science** category. CAN/CAN'T reasoning for tuned versus detuned exciton delivery to a reaction center.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/photosynthetic_exciton_transfer.json](../input/photosynthetic_exciton_transfer.json)

Go translation: [../photosynthetic_exciton_transfer.go](../photosynthetic_exciton_transfer.go)

Expected Markdown output: [../output/photosynthetic_exciton_transfer.md](../output/photosynthetic_exciton_transfer.md)
