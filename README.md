# eyelingo

Small Go translations of selected EyeReasoner/Eyeling N3 examples.

## Layout

```text
go.mod                          local module so examples can share input loading
internal/exampleinput/          shared JSON input loader
examples/                       Go examples
examples/input/                 example-specific JSON data and parameters
examples/output/                expected Markdown output for each example
test                            run examples and compare with expected Markdown output
```

## Example structure

Each example now has three pieces:

```text
examples/example_xyz.go
examples/input/example_xyz.json
examples/output/example_xyz.md
```

The Go file contains the example logic and prints the original ARC-style report as Markdown. The matching input JSON file contains the example-specific facts, data, or parameters that are feasible to externalize.

The output is Markdown, with structured sections such as:

Snapshot files use plain lines rather than Markdown list markers, and add two trailing spaces to every non-empty line so rendered Markdown keeps the same line breaks as stdout.

```md
# Example Title

## Answer
...

## Reason why
...

## Check
...

## Go audit details
...
```

Most examples load their domain fixture directly from `examples/input/<name>.json` through `exampleinput.Load`. A few examples still keep complex internal relation structures in Go, but they also have a matching JSON input file documenting the corresponding data or parameters.

## Rationale

Eyelingo is a small translation laboratory: it takes selected EyeReasoner/Eyeling N3 examples and rewrites them as compact, runnable Go programs. The goal is not to replace N3, but to make the reasoning patterns easy to inspect in a mainstream systems language: facts become typed input data, rules become explicit Go functions, and derived conclusions become reproducible reports.

The examples keep an ARC-style shape: an `Answer` gives the computed result, `Reason why` explains the rule chain or decision path, and `Check` records the invariants that should still hold. This makes every example useful both as a demonstration and as a regression fixture. A reader should be able to scan the output and see not only what was concluded, but also why it was concluded and what was verified.

The `Go audit details` section is intentional. It documents translation-level evidence such as source scenario, input counts, selected thresholds, rule counters, search bounds, precision choices, or platform details. These audit lines help distinguish domain conclusions from implementation evidence, and they make it easier to review changes when a Go translation evolves away from the original N3 file.

STEM is the core of the collection. The examples are chosen to cover scientific measurement, technical interoperability, engineered systems, and mathematical reasoning. Together they show that rule-based examples can remain readable while still exercising realistic concerns: exact arithmetic, graph search, certificates, constraints, policy checks, safety envelopes, Bayesian reasoning, scheduling, routing, and optimization.

## STEM examples

The examples are grouped by their main emphasis. Each row links to the example-specific JSON input, the Go translation, and the expected Markdown output.

### Science

| Example | Description | Input | Go | Output |
|---|---|---|---|---|
| AuroraCare | Health-data permit/deny scenarios across care, quality improvement, and research. | [json](examples/input/auroracare.json) | [go](examples/auroracare.go) | [md](examples/output/auroracare.md) |
| Barley Seed Lineage | Seed-lineage CAN/CAN'T reasoning for reproduction, dormancy, variation, and persistence. | [json](examples/input/barley_seed_lineage.json) | [go](examples/barley_seed_lineage.go) | [md](examples/output/barley_seed_lineage.md) |
| Bayes Diagnosis | Bayesian posterior ranking of possible diseases from symptoms and test evidence. | [json](examples/input/bayes_diagnosis.json) | [go](examples/bayes_diagnosis.go) | [md](examples/output/bayes_diagnosis.md) |
| Bayes Therapy Decision Support | Posterior-weighted utility selection of the best therapy. | [json](examples/input/bayes_therapy.json) | [go](examples/bayes_therapy.go) | [md](examples/output/bayes_therapy.md) |
| BMI — ARC-style Body Mass Index example | Adult BMI calculation, category assignment, and healthy-weight interval. | [json](examples/input/bmi.json) | [go](examples/bmi.go) | [md](examples/output/bmi.md) |
| Cell Marker Panel | Gene-marker panel selection that separates modeled cell populations. | [json](examples/input/cell_marker_panel.json) | [go](examples/cell_marker_panel.go) | [md](examples/output/cell_marker_panel.md) |
| Cold Chain Release | Biologic lot release and allocation under cold-chain and transit constraints. | [json](examples/input/cold_chain_release.json) | [go](examples/cold_chain_release.go) | [md](examples/output/cold_chain_release.md) |
| Digital Product Passport | Component roll-up for recycled content, carbon footprint, repairability, and critical materials. | [json](examples/input/digital_product_passport.json) | [go](examples/digital_product_passport.go) | [md](examples/output/digital_product_passport.md) |
| E-Bike Motor Thermal Envelope | Certified e-bike motor-temperature envelope for an assist plan. | [json](examples/input/ebike_motor_thermal_envelope.json) | [go](examples/ebike_motor_thermal_envelope.go) | [md](examples/output/ebike_motor_thermal_envelope.md) |
| LLD — Leg Length Discrepancy Measurement | Leg-length discrepancy measurement and alarm thresholding. | [json](examples/input/lldm.json) | [go](examples/lldm.go) | [md](examples/output/lldm.md) |
| RC Discharge Envelope | Certified exponential decay envelope for an RC discharge. | [json](examples/input/rc_discharge_envelope.json) | [go](examples/rc_discharge_envelope.go) | [md](examples/output/rc_discharge_envelope.md) |
| Superdense Coding | Quantum-information parity facts for superdense coding. | [json](examples/input/superdense_coding.json) | [go](examples/superdense_coding.go) | [md](examples/output/superdense_coding.md) |

### Technology

| Example | Description | Input | Go | Output |
|---|---|---|---|---|
| Alarm Bit Interoperability | Classical alarm-bit copy and relay tasks contrasted with forbidden universal cloning. | [json](examples/input/alarm_bit_interoperability.json) | [go](examples/alarm_bit_interoperability.go) | [md](examples/output/alarm_bit_interoperability.md) |
| Deep Taxonomy 100000 | Large taxonomy materialization benchmark through a very deep class chain. | [json](examples/input/deep_taxonomy_100000.json) | [go](examples/deep_taxonomy_100000.go) | [md](examples/output/deep_taxonomy_100000.md) |
| Delfour | Privacy-preserving retail insight and recommendation policy. | [json](examples/input/delfour.json) | [go](examples/delfour.go) | [md](examples/output/delfour.md) |
| Doctor Advice Work Conflict | Policy conflict resolution for remote-work and office-work advice. | [json](examples/input/doctor_advice_work_conflict.json) | [go](examples/doctor_advice_work_conflict.go) | [md](examples/output/doctor_advice_work_conflict.md) |
| French Cities | Reachability over a small French city route graph. | [json](examples/input/french_cities.json) | [go](examples/french_cities.go) | [md](examples/output/french_cities.md) |
| Go Capture Scenario (Weiqi) | Weiqi move legality, capture detection, and board update. | [json](examples/input/go_capture.json) | [go](examples/go_capture.go) | [md](examples/output/go_capture.md) |
| Go Eye Capture Scenario (Weiqi) | Weiqi single-eye capture reasoning for a surrounded group. | [json](examples/input/go_eye_capture.json) | [go](examples/go_eye_capture.go) | [md](examples/output/go_eye_capture.md) |
| Go Ko Rule Scenario (Weiqi) | Weiqi ko-rule prevention through board-state comparison. | [json](examples/input/go_ko.json) | [go](examples/go_ko.go) | [md](examples/output/go_ko.md) |
| Gray Code Counter | n-bit Gray-code sequence with one-bit transition checks. | [json](examples/input/gray_code_counter.json) | [go](examples/gray_code_counter.go) | [md](examples/output/gray_code_counter.md) |
| High Trust RDF Bloom Envelope | Bloom-envelope acceptance using canonical graph, index, and false-positive checks. | [json](examples/input/high_trust_bloom_envelope.json) | [go](examples/high_trust_bloom_envelope.go) | [md](examples/output/high_trust_bloom_envelope.md) |
| Ranked DPV Risk Report | ODRL/DPV clause risk ranking by severity and risk class. | [json](examples/input/odrl_dpv_risk_ranked.json) | [go](examples/odrl_dpv_risk_ranked.go) | [md](examples/output/odrl_dpv_risk_ranked.md) |
| Parcel Locker | Delegated parcel pickup-token authorization. | [json](examples/input/parcellocker.json) | [go](examples/parcellocker.go) | [md](examples/output/parcellocker.md) |
| Path Discovery | Airport path discovery with stopover and routing constraints. | [json](examples/input/path_discovery.json) | [go](examples/path_discovery.go) | [md](examples/output/path_discovery.md) |

### Engineering

| Example | Description | Input | Go | Output |
|---|---|---|---|---|
| Brussels Beijing | Cost-optimized flight routing from Brussels to Beijing while avoiding a carrier. | [json](examples/input/brussels_beijing.json) | [go](examples/brussels_beijing.go) | [md](examples/output/brussels_beijing.md) |
| Calidor | Municipal cooling intervention bundle chosen from active needs and budget constraints. | [json](examples/input/calidor.json) | [go](examples/calidor.go) | [md](examples/output/calidor.md) |
| Control System — ARC-style control-system example | Translated measurement and control rules for actuators, inputs, and disturbances. | [json](examples/input/control_system.json) | [go](examples/control_system.go) | [md](examples/output/control_system.md) |
| Crisis Dispatch | Storm-incident resource assignment optimized for triage score and travel time. | [json](examples/input/crisis_dispatch.json) | [go](examples/crisis_dispatch.go) | [md](examples/output/crisis_dispatch.md) |
| Decimal Servo Envelope | Certified servo pole interval and settling-step envelope. | [json](examples/input/decimal_servo_envelope.json) | [go](examples/decimal_servo_envelope.go) | [md](examples/output/decimal_servo_envelope.md) |
| Dijkstra Risk Path | Risk-adjusted path selection using weighted network edges. | [json](examples/input/dijkstra_risk_path.json) | [go](examples/dijkstra_risk_path.go) | [md](examples/output/dijkstra_risk_path.md) |
| Docking Abort Token | Docking abort audit-token flow and safety-system copy restrictions. | [json](examples/input/docking_abort_token.json) | [go](examples/docking_abort_token.go) | [md](examples/output/docking_abort_token.md) |
| Drone Corridor Planner | Constrained drone route planning through corridors and restricted zones. | [json](examples/input/drone_corridor_planner.json) | [go](examples/drone_corridor_planner.go) | [md](examples/output/drone_corridor_planner.md) |
| EV Roadtrip Planner | EV route planning with battery, duration, cost, and comfort constraints. | [json](examples/input/ev_roundtrip_planner.json) | [go](examples/ev_roundtrip_planner.go) | [md](examples/output/ev_roundtrip_planner.md) |
| Flandor | Regional retooling priority calculation for a Flanders scenario. | [json](examples/input/flandor.json) | [go](examples/flandor.go) | [md](examples/output/flandor.md) |
| GPS — Goal driven route planning | Goal-driven route planning over a small road network. | [json](examples/input/gps.json) | [go](examples/gps.go) | [md](examples/output/gps.md) |
| HarborSMR Insight Dispatch | Port electrolysis dispatch decision with safety margin and policy checks. | [json](examples/input/harbor_smr.json) | [go](examples/harbor_smr.go) | [md](examples/output/harbor_smr.md) |
| Isolation Breach Token | Isolation-breach audit-token flow with cloning and fan-out restrictions. | [json](examples/input/isolation_breach_token.json) | [go](examples/isolation_breach_token.go) | [md](examples/output/isolation_breach_token.md) |

### Mathematics

| Example | Description | Input | Go | Output |
|---|---|---|---|---|
| Ackermann | Exact Ackermann and hyperoperation facts, including very large integer results. | [json](examples/input/ackermann.json) | [go](examples/ackermann.go) | [md](examples/output/ackermann.md) |
| Allen Interval Calculus | Allen temporal interval relation closure over completed and explicit intervals. | [json](examples/input/allen_interval_calculus.json) | [go](examples/allen_interval_calculus.go) | [md](examples/output/allen_interval_calculus.md) |
| Complex Numbers | Complex arithmetic and transcendental identity checks. | [json](examples/input/complex_numbers.json) | [go](examples/complex_numbers.go) | [md](examples/output/complex_numbers.md) |
| Dining Philosophers | Chandy-Misra dining-philosophers trace with concurrency conflict checks. | [json](examples/input/dining_philosophers.json) | [go](examples/dining_philosophers.go) | [md](examples/output/dining_philosophers.md) |
| Euler Identity Certificate | High-precision certificate for the identity exp(iπ) + 1 = 0. | [json](examples/input/euler_identity_certificate.json) | [go](examples/euler_identity_certificate.go) | [md](examples/output/euler_identity_certificate.md) |
| Fibonacci Example (Big) | Exact computation of a large Fibonacci number. | [json](examples/input/fibonacci.json) | [go](examples/fibonacci.go) | [md](examples/output/fibonacci.md) |
| Fundamental Theorem Arithmetic | Prime factorization and prime-power representation. | [json](examples/input/fundamental_theorem_arithmetic.json) | [go](examples/fundamental_theorem_arithmetic.go) | [md](examples/output/fundamental_theorem_arithmetic.md) |
| Genetic Knapsack Selection | Deterministic genetic selection for a bounded knapsack. | [json](examples/input/genetic_knapsack_selection.json) | [go](examples/genetic_knapsack_selection.go) | [md](examples/output/genetic_knapsack_selection.md) |
| Gradient Descent Step | Certified single gradient-descent step for a quadratic objective. | [json](examples/input/gradient_descent_step.json) | [go](examples/gradient_descent_step.go) | [md](examples/output/gradient_descent_step.md) |
| Kaprekar 6174 | Kaprekar chains and basin facts ending at 6174. | [json](examples/input/kaprekar_6174.json) | [go](examples/kaprekar_6174.go) | [md](examples/output/kaprekar_6174.md) |
| 8-Queens | 8-Queens constraint satisfaction with a valid board solution. | [json](examples/input/queens.json) | [go](examples/queens.go) | [md](examples/output/queens.md) |
| Sudoku | Sudoku constraint solving with a unique completed grid. | [json](examples/input/sudoku.json) | [go](examples/sudoku.go) | [md](examples/output/sudoku.md) |


## Run

Run one example from the repository root:

```sh
go run examples/bmi.go
```

The program writes Markdown output to stdout.

Run the full regression test:

```sh
./test
```

The test prints `OK` or `FAIL` for each example, per-example timing, and total time. It compares against `examples/output/*.md` while normalizing the Go platform audit value, since that varies by machine.

Regenerate expected outputs after intentional changes:

```sh
./test --update
```
