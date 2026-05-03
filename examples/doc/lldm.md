# LLD — Leg Length Discrepancy Measurement

`lldm` is a Go translation/adaptation of Eyeling's `lldm.n3`.

## Background

Leg length discrepancy can be estimated from medical-image landmarks by constructing geometric lines and distances. The calculation depends on coordinate differences, slopes, intersections, and Euclidean distance. An alarm threshold then turns the measured discrepancy into a clinical flag while preserving the numeric steps that led to it.

## What it demonstrates

**Category:** Science. Leg-length discrepancy measurement and alarm thresholding.

## How the Go implementation works

The program processes medical-image measurement points through a small geometry pipeline. It computes coordinate differences, line slopes, intersections, and Euclidean distances, then compares the discrepancy with the alarm threshold.

Checks cover each geometry stage so the final clinical flag can be traced back to the measured points.

## Files

Input JSON: [../input/lldm.json](../input/lldm.json)

Go implementation: [../lldm.go](../lldm.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/lldm.md](../output/lldm.md)
