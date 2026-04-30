# Drone Corridor Planner

`drone_corridor_planner` is a Go translation/adaptation of Eyeling's `drone-corridor-planner.n3`.

The context is constrained aerial navigation. Candidate drone corridors are filtered by restrictions, risk, and route feasibility before a path is accepted.

## How it works

A compact Go translation inspired by Eyeling's
`examples/drone-corridor-planner.n3`.

The example composes corridor actions into bounded plans from Gent to
Oostende. Duration and cost are summed; belief and comfort are multiplied.
Plans are kept only when they satisfy the same pruning style as the N3 file.

## What it demonstrates

This example is mainly in the **Engineering** category. Constrained drone route planning through corridors and restricted zones.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/drone_corridor_planner.json](../input/drone_corridor_planner.json)

Go translation: [../drone_corridor_planner.go](../drone_corridor_planner.go)

Expected Markdown output: [../output/drone_corridor_planner.md](../output/drone_corridor_planner.md)
