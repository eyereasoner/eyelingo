# School Placement Route Audit

`school_placement_audit` is a small algorithmic-governance example inspired by school-placement disputes where automated or semi-automated tools can hide important local geography.

## Background

A placement rule that uses straight-line distance can look objective while ignoring the real route a child must walk. Rivers, highways, bridges, and rail corridors can turn a short map distance into a long or unsuitable journey. The example keeps the data fictionalized and small, but it mirrors the audit question: what changes when the claimed support tool is checked against route distance and family preferences?

## What it demonstrates

**Category:** Technology. Route-aware audit of a straight-line school-placement support tool.

## How the Go implementation works

The Go implementation loads four students, four schools, ranked preferences, and a complete distance matrix. It derives two assignments:

- a provisional support-tool assignment that chooses the smallest straight-line distance, using preference rank only as a tie-breaker;
- an audited assignment that scores each candidate as walking-route distance plus a fixed penalty for each preference step.

It then flags students whose provisional assignment differs from the audited assignment or exceeds the walking-distance limit.

The isolated Go check independently recomputes both assignment models, parses the reported affected children and recommended assignments, verifies the largest hidden detour, and checks that the report requests an inspectable explanation.

## Files

Input JSON: [../input/school_placement_audit.json](../input/school_placement_audit.json)

Go implementation: [../school_placement_audit.go](../school_placement_audit.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/school_placement_audit.md](../output/school_placement_audit.md)
