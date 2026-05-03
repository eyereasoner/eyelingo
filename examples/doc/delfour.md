# Delfour

`delfour` is a Go translation/adaptation of Eyeling's `delfour.n3`.

## Background

Retail personalization can be useful but risky when it depends on sensitive household information. A privacy-preserving design keeps raw personal signals on a trusted device and shares only a narrow, expiring insight that a store can verify for a specific purpose. This example combines signed evidence, policy checks, and product recommendation without exposing the underlying private condition.

## What it demonstrates

**Category:** Technology. Privacy-preserving retail insight and recommendation policy.

## How the Go implementation works

The program models the privacy-preserving retail flow with structs for household insight, scanner context, policy requirements, signatures, and recommendation candidates. It checks the envelope before allowing the product suggestion: purpose, expiry, minimization, hash, signature, and policy scope all have to pass.

Only after the authorization path succeeds does the code choose and report the lower-sugar recommendation.

## Files

Input JSON: [../input/delfour.json](../input/delfour.json)

Go implementation: [../delfour.go](../delfour.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/delfour.md](../output/delfour.md)
