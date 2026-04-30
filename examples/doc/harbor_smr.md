# HarborSMR Insight Dispatch

`harbor_smr` is a Go translation/adaptation of Eyeling's `harborsmr.n3`.

## Background

A port hydrogen hub may need flexible electricity or heat insight from a small modular reactor operator, but raw reactor telemetry is highly sensitive. A safer workflow shares only a bounded dispatch insight for a specific electrolyzer window and keeps safety-critical source data local. The example combines energy-dispatch reasoning with policy, minimization, and safety-margin checks.

## What it demonstrates

**Category:** Engineering. Port electrolysis dispatch decision with safety margin and policy checks.

## How the Go implementation works

The program models a permissioned dispatch decision for a port hydrogen hub. It checks the insight envelope for scope, expiry, privacy minimization, safety margin, and threshold conditions before approving the electrolyzer window.

The report shows the dispatch conclusion together with the policy and safety checks that made it valid.

## Files

Input JSON: [../input/harbor_smr.json](../input/harbor_smr.json)

Go implementation: [../harbor_smr.go](../harbor_smr.go)

Expected Markdown output: [../output/harbor_smr.md](../output/harbor_smr.md)
