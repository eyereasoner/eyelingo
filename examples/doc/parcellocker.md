# Parcel Locker

`parcellocker` is a Go translation/adaptation of Eyeling's `parcellocker.n3`.

The context is delegated authorization. A one-time parcel pickup request is permitted only when identity, token, time, locker, and delegation constraints are all satisfied.

## How it works

A self-contained Go translation of the Eyeling ParcelLocker delegation
example.

The source N3 file models a narrow, one-time permission: Maya lets Noah pick
up one specific parcel from one specific locker for the pickup purpose only.
The same token must not reveal billing details, redirect the parcel, work for
another person, work at another locker, or be used after it has been consumed.

This program is intentionally not a generic RDF or N3 reasoner. It translates
the concrete facts and rules for this example into ordinary Go structs and
checks so the authorization boundary is easy to read and run.

## What it demonstrates

This example is mainly in the **Technology** category. Delegated parcel pickup-token authorization.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/parcellocker.json](../input/parcellocker.json)

Go translation: [../parcellocker.go](../parcellocker.go)

Expected Markdown output: [../output/parcellocker.md](../output/parcellocker.md)
