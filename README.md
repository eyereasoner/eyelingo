# eyelingo

Small Go translations of selected EyeReasoner/Eyeling N3 examples.

## Layout

```text
go.mod                          local module so examples can share input loading
internal/exampleinput/          shared JSON input loader
examples/                       Go examples
examples/input/                 example-specific JSON data and parameters
examples/output/                expected Markdown output for each example
examples/checks/                independent Python check implementations
examples/doc/                   short explanatory notes for each example
tools/run_check.py              run one Python Check implementation
tools/build_output.py           append Python Check output to a Go report prefix
test                            run examples and compare with expected Markdown output
```

## Example structure

Each example now has five pieces:

```text
examples/example_xyz.go
examples/input/example_xyz.json
examples/checks/example_xyz.py
examples/output/example_xyz.md
examples/doc/example_xyz.md
```

The Go file contains the example logic and prints the ARC-style title, `Answer`, and `Reason why` report prefix as Markdown. The corresponding Python file implements the independent `Check` section. The matching input JSON file contains the example-specific facts, data, or parameters that are feasible to externalize. The companion doc file gives a short plain-language explanation of what the example demonstrates, how to read the output, and where the main files live.

Output is delivered in structured Markdown. Snapshotss use plain lines rather than Markdown list markers, and add two trailing spaces to every non-empty line so rendered Markdown keeps the same line breaks as `stdout`.

Most examples load their domain fixture directly from `examples/input/<name>.json` through `exampleinput.Load`. A few examples still keep complex internal relation structures in Go, but they also have a matching JSON input file documenting the corresponding data or parameters.

## Rationale

Eyelingo is a small translation laboratory: it takes selected EyeReasoner/Eyeling N3 examples and rewrites them as compact, runnable Go programs. The goal is not to replace N3, but to make the reasoning patterns easy to inspect in a mainstream systems language: facts become typed input data, rules become explicit Go functions, and derived conclusions become reproducible reports.

The examples keep an ARC-style shape: an `Answer` gives the computed result, `Reason why` explains the rule chain or decision path, and `Check` records the invariants that should still hold. The `Check` section is intentionally produced by Python, not Go, so it cannot call Go helper functions or reuse Go intermediate state from the answer path. The regression test builds each snapshot as `Go report prefix + Python Check`.

The visible output no longer includes Go audit details. Implementation diagnostics are kept out of the report so the Markdown focuses on the domain answer, explanation, and independent verification.

## Python Check flow

The Go examples only produce the human-readable report prefix: the title, `## Answer`, and `## Reason why`. They do not print the `## Check` section and they do not call the Python check modules directly.

During `./test`, the runner first captures that Go output prefix. It then calls `tools/build_output.py`, which imports the matching module from `examples/checks/<example>.py`. That Python module independently reconstructs the relevant facts from the JSON fixture and/or from the captured prefix, runs its own assertions, and returns the lines for the Markdown `## Check` section.

The final snapshot compared by the test suite is therefore:

```text
Go title + Answer + Reason why
+ Python-generated Check
```

This keeps the visible checks in a separate implementation language and prevents them from sharing Go helper functions or intermediate Go state.

STEM is the core of the collection. The examples are chosen to cover scientific measurement, technical interoperability, engineered systems, and mathematical reasoning. Together they show that rule-based examples can remain readable while still exercising realistic concerns: exact arithmetic, graph search, certificates, constraints, policy checks, safety envelopes, Bayesian reasoning, scheduling, routing, and optimization.

## ARC-style STEM examples

The examples are grouped by their main emphasis. Each row links to the example-specific JSON input, the Go translation, the independent Python checks, the expected Markdown output, and the companion documentation.

### Science

| Example | Description | Input | Go | Checks | Output | Doc |
|---|---|---|---|---|---|---|
| AuroraCare | Health-data permit/deny scenarios across care, quality improvement, and research, echoing [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/). | [json](examples/input/auroracare.json) | [go](examples/auroracare.go) | [py](examples/checks/auroracare.py) | [md](examples/output/auroracare.md) | [md](examples/doc/auroracare.md) |
| Barley Seed Lineage | Seed-lineage CAN/CAN'T reasoning for reproduction, dormancy, variation, and persistence. | [json](examples/input/barley_seed_lineage.json) | [go](examples/barley_seed_lineage.go) | [py](examples/checks/barley_seed_lineage.py) | [md](examples/output/barley_seed_lineage.md) | [md](examples/doc/barley_seed_lineage.md) |
| Bayes Diagnosis | Bayesian posterior ranking of possible diseases from symptoms and test evidence. | [json](examples/input/bayes_diagnosis.json) | [go](examples/bayes_diagnosis.go) | [py](examples/checks/bayes_diagnosis.py) | [md](examples/output/bayes_diagnosis.md) | [md](examples/doc/bayes_diagnosis.md) |
| Bayes Therapy Decision Support | Posterior-weighted utility selection of the best therapy. | [json](examples/input/bayes_therapy.json) | [go](examples/bayes_therapy.go) | [py](examples/checks/bayes_therapy.py) | [md](examples/output/bayes_therapy.md) | [md](examples/doc/bayes_therapy.md) |
| BMI — Body Mass Index | Adult BMI calculation, category assignment, and healthy-weight interval. | [json](examples/input/bmi.json) | [go](examples/bmi.go) | [py](examples/checks/bmi.py) | [md](examples/output/bmi.md) | [md](examples/doc/bmi.md) |
| Digital Product Passport | Component roll-up for recycled content, carbon footprint, repairability, and critical materials. | [json](examples/input/digital_product_passport.json) | [go](examples/digital_product_passport.go) | [py](examples/checks/digital_product_passport.py) | [md](examples/output/digital_product_passport.md) | [md](examples/doc/digital_product_passport.md) |
| E-Bike Motor Thermal Envelope | Certified e-bike motor-temperature envelope for an assist plan. | [json](examples/input/ebike_motor_thermal_envelope.json) | [go](examples/ebike_motor_thermal_envelope.go) | [py](examples/checks/ebike_motor_thermal_envelope.py) | [md](examples/output/ebike_motor_thermal_envelope.md) | [md](examples/doc/ebike_motor_thermal_envelope.md) |
| Gravity Mediator Witness | Mediator-only entanglement witness contrasting non-classical and purely classical gravitational mediators. | [json](examples/input/gravity_mediator_witness.json) | [go](examples/gravity_mediator_witness.go) | [py](examples/checks/gravity_mediator_witness.py) | [md](examples/output/gravity_mediator_witness.md) | [md](examples/doc/gravity_mediator_witness.md) |
| LLD — Leg Length Discrepancy Measurement | Leg-length discrepancy measurement and alarm thresholding. | [json](examples/input/lldm.json) | [go](examples/lldm.go) | [py](examples/checks/lldm.py) | [md](examples/output/lldm.md) | [md](examples/doc/lldm.md) |
| Photosynthetic Exciton Transfer | CAN/CAN'T reasoning for tuned versus detuned exciton delivery to a reaction center. | [json](examples/input/photosynthetic_exciton_transfer.json) | [go](examples/photosynthetic_exciton_transfer.go) | [py](examples/checks/photosynthetic_exciton_transfer.py) | [md](examples/output/photosynthetic_exciton_transfer.md) | [md](examples/doc/photosynthetic_exciton_transfer.md) |
| RC Discharge Envelope | Certified exponential decay envelope for an RC discharge. | [json](examples/input/rc_discharge_envelope.json) | [go](examples/rc_discharge_envelope.go) | [py](examples/checks/rc_discharge_envelope.py) | [md](examples/output/rc_discharge_envelope.md) | [md](examples/doc/rc_discharge_envelope.md) |
| Superdense Coding | Quantum-information parity facts for superdense coding. | [json](examples/input/superdense_coding.json) | [go](examples/superdense_coding.go) | [py](examples/checks/superdense_coding.py) | [md](examples/output/superdense_coding.md) | [md](examples/doc/superdense_coding.md) |

### Technology

| Example | Description | Input | Go | Checks | Output | Doc |
|---|---|---|---|---|---|---|
| Alarm Bit Interoperability | Classical alarm-bit copy and relay tasks contrasted with forbidden universal cloning. | [json](examples/input/alarm_bit_interoperability.json) | [go](examples/alarm_bit_interoperability.go) | [py](examples/checks/alarm_bit_interoperability.py) | [md](examples/output/alarm_bit_interoperability.md) | [md](examples/doc/alarm_bit_interoperability.md) |
| Deep Taxonomy 100000 | Large taxonomy materialization benchmark through a very deep class chain. | [json](examples/input/deep_taxonomy_100000.json) | [go](examples/deep_taxonomy_100000.go) | [py](examples/checks/deep_taxonomy_100000.py) | [md](examples/output/deep_taxonomy_100000.md) | [md](examples/doc/deep_taxonomy_100000.md) |
| Delfour | Privacy-preserving retail insight and recommendation policy, echoing [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/). | [json](examples/input/delfour.json) | [go](examples/delfour.go) | [py](examples/checks/delfour.py) | [md](examples/output/delfour.md) | [md](examples/doc/delfour.md) |
| Doctor Advice Work Conflict | Policy conflict resolution for remote-work and office-work advice. | [json](examples/input/doctor_advice_work_conflict.json) | [go](examples/doctor_advice_work_conflict.go) | [py](examples/checks/doctor_advice_work_conflict.py) | [md](examples/output/doctor_advice_work_conflict.md) | [md](examples/doc/doctor_advice_work_conflict.md) |
| FFT8 Numeric | Eight-point Fourier transform over a sampled sine wave with conjugate-bin and energy checks. | [json](examples/input/fft8_numeric.json) | [go](examples/fft8_numeric.go) | [py](examples/checks/fft8_numeric.py) | [md](examples/output/fft8_numeric.md) | [md](examples/doc/fft8_numeric.md) |
| FFT32 Numeric | Thirty-two-point Fourier transform over several sampled waveforms with spectral invariant checks. | [json](examples/input/fft32_numeric.json) | [go](examples/fft32_numeric.go) | [py](examples/checks/fft32_numeric.py) | [md](examples/output/fft32_numeric.md) | [md](examples/doc/fft32_numeric.md) |
| French Cities | Reachability over a small French city route graph. | [json](examples/input/french_cities.json) | [go](examples/french_cities.go) | [py](examples/checks/french_cities.py) | [md](examples/output/french_cities.md) | [md](examples/doc/french_cities.md) |
| Gray Code Counter | n-bit Gray-code sequence with one-bit transition checks. | [json](examples/input/gray_code_counter.json) | [go](examples/gray_code_counter.go) | [py](examples/checks/gray_code_counter.py) | [md](examples/output/gray_code_counter.md) | [md](examples/doc/gray_code_counter.md) |
| High Trust RDF Bloom Envelope | Bloom-envelope acceptance using canonical graph, index, and false-positive checks, echoing [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/). | [json](examples/input/high_trust_bloom_envelope.json) | [go](examples/high_trust_bloom_envelope.go) | [py](examples/checks/high_trust_bloom_envelope.py) | [md](examples/output/high_trust_bloom_envelope.md) | [md](examples/doc/high_trust_bloom_envelope.md) |
| Parcel Locker | Delegated parcel pickup-token authorization, echoing [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/). | [json](examples/input/parcellocker.json) | [go](examples/parcellocker.go) | [py](examples/checks/parcellocker.py) | [md](examples/output/parcellocker.md) | [md](examples/doc/parcellocker.md) |
| Path Discovery | Airport path discovery with stopover and routing constraints. | [json](examples/input/path_discovery.json) | [go](examples/path_discovery.go) | [py](examples/checks/path_discovery.py) | [md](examples/output/path_discovery.md) | [md](examples/doc/path_discovery.md) |
| Ranked DPV Risk Report | ODRL/DPV clause risk ranking by severity and risk class. | [json](examples/input/odrl_dpv_risk_ranked.json) | [go](examples/odrl_dpv_risk_ranked.go) | [py](examples/checks/odrl_dpv_risk_ranked.py) | [md](examples/output/odrl_dpv_risk_ranked.md) | [md](examples/doc/odrl_dpv_risk_ranked.md) |

### Engineering

| Example | Description | Input | Go | Checks | Output | Doc |
|---|---|---|---|---|---|---|
| Calidor | Municipal cooling intervention bundle chosen from active needs and budget constraints, echoing [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/). | [json](examples/input/calidor.json) | [go](examples/calidor.go) | [py](examples/checks/calidor.py) | [md](examples/output/calidor.md) | [md](examples/doc/calidor.md) |
| Complex Matrix Stability | Discrete-time stability classification using spectral radii of diagonal complex matrices. | [json](examples/input/complex_matrix_stability.json) | [go](examples/complex_matrix_stability.go) | [py](examples/checks/complex_matrix_stability.py) | [md](examples/output/complex_matrix_stability.md) | [md](examples/doc/complex_matrix_stability.md) |
| Control System | Translated measurement and control rules for actuators, inputs, and disturbances. | [json](examples/input/control_system.json) | [go](examples/control_system.go) | [py](examples/checks/control_system.py) | [md](examples/output/control_system.md) | [md](examples/doc/control_system.md) |
| Dijkstra Risk Path | Risk-adjusted path selection using weighted network edges. | [json](examples/input/dijkstra_risk_path.json) | [go](examples/dijkstra_risk_path.go) | [py](examples/checks/dijkstra_risk_path.py) | [md](examples/output/dijkstra_risk_path.md) | [md](examples/doc/dijkstra_risk_path.md) |
| Docking Abort Token | Docking abort audit-token flow and safety-system copy restrictions. | [json](examples/input/docking_abort_token.json) | [go](examples/docking_abort_token.go) | [py](examples/checks/docking_abort_token.py) | [md](examples/output/docking_abort_token.md) | [md](examples/doc/docking_abort_token.md) |
| Drone Corridor Planner | Constrained drone route planning through corridors and restricted zones. | [json](examples/input/drone_corridor_planner.json) | [go](examples/drone_corridor_planner.go) | [py](examples/checks/drone_corridor_planner.py) | [md](examples/output/drone_corridor_planner.md) | [md](examples/doc/drone_corridor_planner.md) |
| EV Roadtrip Planner | EV route planning with battery, duration, cost, and comfort constraints. | [json](examples/input/ev_roundtrip_planner.json) | [go](examples/ev_roundtrip_planner.go) | [py](examples/checks/ev_roundtrip_planner.py) | [md](examples/output/ev_roundtrip_planner.md) | [md](examples/doc/ev_roundtrip_planner.md) |
| Flandor | Regional retooling priority calculation for a Flanders scenario, echoing [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/). | [json](examples/input/flandor.json) | [go](examples/flandor.go) | [py](examples/checks/flandor.py) | [md](examples/output/flandor.md) | [md](examples/doc/flandor.md) |
| GPS — Goal driven route planning | Goal-driven route planning over a small road network. | [json](examples/input/gps.json) | [go](examples/gps.go) | [py](examples/checks/gps.py) | [md](examples/output/gps.md) | [md](examples/doc/gps.md) |
| HarborSMR Insight Dispatch | Port electrolysis dispatch decision with safety margin and policy checks, echoing [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/). | [json](examples/input/harbor_smr.json) | [go](examples/harbor_smr.go) | [py](examples/checks/harbor_smr.py) | [md](examples/output/harbor_smr.md) | [md](examples/doc/harbor_smr.md) |
| Isolation Breach Token | Isolation-breach audit-token flow with cloning and fan-out restrictions. | [json](examples/input/isolation_breach_token.json) | [go](examples/isolation_breach_token.go) | [py](examples/checks/isolation_breach_token.py) | [md](examples/output/isolation_breach_token.md) | [md](examples/doc/isolation_breach_token.md) |
| Wind Turbine Envelope | Wind-speed envelope classification with cubic power curve and interval energy audit. | [json](examples/input/wind_turbine.json) | [go](examples/wind_turbine.go) | [py](examples/checks/wind_turbine.py) | [md](examples/output/wind_turbine.md) | [md](examples/doc/wind_turbine.md) |

### Mathematics

| Example | Description | Input | Go | Checks | Output | Doc |
|---|---|---|---|---|---|---|
| 8-Queens | 8-Queens constraint satisfaction with a valid board solution. | [json](examples/input/queens.json) | [go](examples/queens.go) | [py](examples/checks/queens.py) | [md](examples/output/queens.md) | [md](examples/doc/queens.md) |
| Ackermann | Exact Ackermann and hyperoperation facts, including very large integer results. | [json](examples/input/ackermann.json) | [go](examples/ackermann.go) | [py](examples/checks/ackermann.py) | [md](examples/output/ackermann.md) | [md](examples/doc/ackermann.md) |
| Allen Interval Calculus | Allen temporal interval relation closure over completed and explicit intervals. | [json](examples/input/allen_interval_calculus.json) | [go](examples/allen_interval_calculus.go) | [py](examples/checks/allen_interval_calculus.py) | [md](examples/output/allen_interval_calculus.md) | [md](examples/doc/allen_interval_calculus.md) |
| Complex Numbers | Complex arithmetic and transcendental identity checks. | [json](examples/input/complex_numbers.json) | [go](examples/complex_numbers.go) | [py](examples/checks/complex_numbers.py) | [md](examples/output/complex_numbers.md) | [md](examples/doc/complex_numbers.md) |
| Dining Philosophers | Chandy-Misra dining-philosophers trace with concurrency conflict checks. | [json](examples/input/dining_philosophers.json) | [go](examples/dining_philosophers.go) | [py](examples/checks/dining_philosophers.py) | [md](examples/output/dining_philosophers.md) | [md](examples/doc/dining_philosophers.md) |
| Euler Identity Certificate | High-precision certificate for the identity exp(iπ) + 1 = 0. | [json](examples/input/euler_identity_certificate.json) | [go](examples/euler_identity_certificate.go) | [py](examples/checks/euler_identity_certificate.py) | [md](examples/output/euler_identity_certificate.md) | [md](examples/doc/euler_identity_certificate.md) |
| Fibonacci Example (Big) | Exact computation of a large Fibonacci number. | [json](examples/input/fibonacci.json) | [go](examples/fibonacci.go) | [py](examples/checks/fibonacci.py) | [md](examples/output/fibonacci.md) | [md](examples/doc/fibonacci.md) |
| Fundamental Theorem Arithmetic | Prime factorization and prime-power representation. | [json](examples/input/fundamental_theorem_arithmetic.json) | [go](examples/fundamental_theorem_arithmetic.go) | [py](examples/checks/fundamental_theorem_arithmetic.py) | [md](examples/output/fundamental_theorem_arithmetic.md) | [md](examples/doc/fundamental_theorem_arithmetic.md) |
| Genetic Knapsack Selection | Deterministic genetic selection for a bounded knapsack. | [json](examples/input/genetic_knapsack_selection.json) | [go](examples/genetic_knapsack_selection.go) | [py](examples/checks/genetic_knapsack_selection.py) | [md](examples/output/genetic_knapsack_selection.md) | [md](examples/doc/genetic_knapsack_selection.md) |
| Goldbach 1000 | Bounded strong-Goldbach checker for every even integer from 4 through 1000. | [json](examples/input/goldbach_1000.json) | [go](examples/goldbach_1000.go) | [py](examples/checks/goldbach_1000.py) | [md](examples/output/goldbach_1000.md) | [md](examples/doc/goldbach_1000.md) |
| Kaprekar 6174 | Kaprekar chains and basin facts ending at 6174. | [json](examples/input/kaprekar_6174.json) | [go](examples/kaprekar_6174.go) | [py](examples/checks/kaprekar_6174.py) | [md](examples/output/kaprekar_6174.md) | [md](examples/doc/kaprekar_6174.md) |
| Sudoku | Sudoku constraint solving with a unique completed grid. | [json](examples/input/sudoku.json) | [go](examples/sudoku.go) | [py](examples/checks/sudoku.py) | [md](examples/output/sudoku.md) | [md](examples/doc/sudoku.md) |

## Run

Run one example from the repository root:

```sh
go run examples/bmi.go
```

The program writes the Markdown report prefix to stdout. During `./test`, Python appends the `## Check` section before the combined report is compared with `examples/output/*.md`.

Run the full regression test:

```sh
./test
```

The test prints `OK` or `FAIL` for each example, per-example timing, and total time. It compares against `examples/output/*.md` after appending the corresponding independent Python `Check` implementation.

Regenerate expected outputs after intentional changes:

```sh
./test --update
```
