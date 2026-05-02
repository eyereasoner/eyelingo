# Parcel Locker

`parcellocker` is a Go translation/adaptation of Eyeling's `parcellocker.n3`.

## Background

Delegated pickup systems let one person authorize another to collect a parcel, but the token must be limited by parcel, locker, identity, time window, and purpose. One-time use and anti-redirection checks prevent the same authorization from being reused or applied to unrelated data. This example models those constraints as an authorization decision.

## What it demonstrates

**Category:** Technology. Delegated parcel pickup-token authorization.

## How the Go implementation works

The Go code models a one-time delegated pickup token and evaluates it against identity, parcel, locker, purpose, time window, and consumption state. It rejects attempts to reuse the token, redirect the parcel, use a different locker, or expose unrelated billing data.

The final authorization is reported with the specific constraints that passed or failed.

## Files

Input JSON: [../input/parcellocker.json](../input/parcellocker.json)

Go implementation: [../parcellocker.go](../parcellocker.go)

Python check: [../checks/parcellocker.py](../checks/parcellocker.py)

Expected Markdown output: [../output/parcellocker.md](../output/parcellocker.md)
