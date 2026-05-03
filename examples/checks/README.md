# Isolated Go checks

This directory contains the independent Go implementation of SEE `## Check` sections. It is a separate Go module located under `examples/checks/` so the checker can be built without importing the root `see` module or any example program.

The test runner captures an example report prefix (`# ...`, `## Answer`, and `## Reason`), then calls this checker with the example name and the captured prefix file. The checker reloads the project-root input at `examples/input/<example>.json`, recomputes the relevant facts, parses the captured report where useful, and emits the visible `## Check` section.

The important constraint is module separation: this checker module has no dependency on the root module. Checks receive only the example name, the captured report prefix, and the fixture JSON. That prevents the check path from importing example helper functions or reusing intermediate state from the answer path.

Each example has a registered domain-specific Go checker routine. If a new example is added without one, the checker binary fails closed with a clear `no Go check registered` error instead of falling back to a generic pass.

File layout:

- `main.go` contains the checker harness, dispatcher, and shared loader/context.
- `shared_helpers.go` contains common parsing and recomputation helpers.
- `<example>.go` files contain the example-specific check routines.
