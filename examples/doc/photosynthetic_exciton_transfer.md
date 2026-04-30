# Photosynthetic Exciton Transfer

`photosynthetic_exciton_transfer` is a Go translation/adaptation of Eyeling's `act-photosynthetic-exciton-transfer.n3`.

## Background

Photosynthesis begins when absorbed light creates an exciton, an electronic excitation that must reach a reaction center efficiently. Coupling, energy gaps, vibronic effects, coherence, and dephasing all influence whether transfer is possible. This example compares a tuned complex that can deliver the excitation with a detuned contrast case that cannot.

## What it demonstrates

**Category:** Science. CAN/CAN'T reasoning for tuned versus detuned exciton delivery to a reaction center.

## How the Go implementation works

The Go code evaluates each photosynthetic complex against coupling, delocalization, vibronic bridge, energy landscape, coherence, dephasing, and connection-to-center fields. It derives CAN entries for the tuned complex and CAN'T entries for the detuned contrast.

Checks make sure the successful transfer has all required conditions and that the failed case is blocked by the expected missing properties.

## Files

Input JSON: [../input/photosynthetic_exciton_transfer.json](../input/photosynthetic_exciton_transfer.json)

Go implementation: [../photosynthetic_exciton_transfer.go](../photosynthetic_exciton_transfer.go)

Expected Markdown output: [../output/photosynthetic_exciton_transfer.md](../output/photosynthetic_exciton_transfer.md)
