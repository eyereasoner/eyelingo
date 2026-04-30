# Calidor

`calidor` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on municipal cooling intervention bundle chosen from active needs and budget constraints. Its input fixture is organized around `CaseName`, `Question`, `RequestPurpose`, `RequestAction`, `GatewayCreatedAt`, `GatewayExpiresAt`, `CityAuthAt`, `CityDutyAt`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Name: Calidor Municipality: Calidor Metric: active_need_count

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - the trusted precomputed HMAC signature verifies C2 OK - payload hash matches the source envelope digest C3 OK - the shared insight strips raw heat, vulnerability, credit, and meter-trace terms

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/calidor.json](../input/calidor.json)

Go translation: [../calidor.go](../calidor.go)

Expected Markdown output: [../output/calidor.md](../output/calidor.md)
