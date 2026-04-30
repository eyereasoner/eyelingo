# Photosynthetic Exciton Transfer

`photosynthetic_exciton_transfer` is a Go translation/adaptation of Eyeling's `act-photosynthetic-exciton-transfer.n3`.

The example asks whether a tuned photosynthetic antenna can deliver an excitation to a reaction center while a detuned contrast antenna cannot. It is careful about the claim: short-lived quantum-assisted transfer is enough in the tuned regime; long-lived coherence is not assumed.

## What it demonstrates

This is mainly a **Science** example. The tuned complex has strong excitonic coupling, delocalization, a tuned vibronic bridge, moderate dephasing, and a downhill route. The contrast complex has weak coupling, no useful bridge, strong dephasing, and a trapping mismatch.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, counters, thresholds, derived facts, or platform details.

## Files

Input JSON: [../input/photosynthetic_exciton_transfer.json](../input/photosynthetic_exciton_transfer.json)

Go translation: [../photosynthetic_exciton_transfer.go](../photosynthetic_exciton_transfer.go)

Expected Markdown output: [../output/photosynthetic_exciton_transfer.md](../output/photosynthetic_exciton_transfer.md)
