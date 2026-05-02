# eyelingo

Eyelingo is a collection of small, runnable Go translations of selected EyeReasoner/Eyeling N3 examples. Each example keeps the ARC-style report shape:

```text
Answer
Reason
Check
```

## Quick start

Run the full regression suite:

```sh
./test
```

Run one Go example directly:

```sh
go run examples/bmi.go
```

That command prints only the report prefix: title, `## Answer`, and `## Reason`. The `## Check` section is appended by the test/update pipeline.

Regenerate expected Markdown outputs after an intentional change:

```sh
./test --update
```

## What is in this repository

This project is a translation laboratory. Facts from the original examples become typed input data, rules become explicit Go functions, and derived conclusions become reproducible Markdown reports.

The main inspiration is Prof. Ruben Verborgh's [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/). In that spirit, many examples are small insight derivations: structured inputs are transformed into answers, decisions, rankings, or certificates, then checked independently rather than trusted because the report reads well.

The examples focus on STEM reasoning: scientific measurement, technical interoperability, engineered systems, and mathematics. They cover exact arithmetic, graph search, certificates, constraints, policy checks, safety envelopes, Bayesian reasoning, scheduling, routing, and optimization.

## Repository guide

The repository is organized around one stem per example. For an example named `<name>`, the main files are:

```text
examples/<name>.go          Go translation that computes Answer and Reason
examples/input/<name>.json  example-specific facts, data, or parameters
examples/checks/<name>.py   independent Python implementation of Check
examples/output/<name>.md   expected combined Markdown report
examples/doc/<name>.md      short explanatory note
```

Supporting code and tools live alongside those example files:

```text
go.mod                         local module so examples can share input loading
internal/exampleinput/         shared JSON input loader
tools/run_check.py             run one Python Check implementation
tools/build_output.py          append Python Check output to a Go report prefix
test                           run examples and compare with expected Markdown output
```

Most Go examples load their domain fixture from `examples/input/<name>.json` through `internal/exampleinput`. A few examples still keep complex relation structures directly in Go; those still have matching JSON input files that document the corresponding data or parameters.

Expected Markdown outputs use plain lines rather than Markdown list markers. Non-empty output lines end with two spaces so rendered Markdown preserves the same line breaks as stdout.

## Independent Python checks

The checks deliberately live in a different programming language from the answer implementation.

During `./test`:

1. The test runner executes `go run examples/<name>.go` and captures the Go report prefix.
2. `tools/build_output.py` calls `tools/run_check.py` for the same example.
3. `tools/run_check.py` imports `examples/checks/<name>.py`.
4. The Python check module reconstructs the relevant facts from JSON and/or parses the captured report prefix.
5. Python emits the visible `## Check` section.
6. The combined report is compared with `examples/output/<name>.md`.

The snapshot under test is therefore:

```text
Go title + Answer + Reason
+ Python-generated Check
```

Following a design suggestion by Prof. Ruben Verborgh, the `Check` section is deliberately not produced by Go. Go computes and explains the answer; Python independently verifies it and appends the visible `## Check` section during testing.

Eyelingo makes reasoning auditable: not by trusting the explanation, but by checking it independently. When Go and Python disagree, the test fails. Eyelingo does not assume either side is automatically correct; the disagreement points to a bug in the Go implementation, the Python verifier, the stored snapshot, or the example specification. The independent check is not a second source of truth by itself; it is a deliberately separate witness that makes disagreements visible and reviewable.

This separation prevents the checks from calling Go helper functions or reusing Go intermediate state from the answer path. Shared Python helper code lives in `examples/checks/common.py`; substantive checks are implemented in the per-example modules.

The visible output no longer includes Go audit details. Implementation diagnostics stay out of the report so the Markdown focuses on the domain answer, explanation, and independent verification.

## Example catalog

Each row links to the example-specific JSON input, Go translation, independent Python checks, expected Markdown output, and companion documentation.

### Science

| Example | Description | Input | Go | Checks | Output | Doc |
|---|---|---|---|---|---|---|
| AuroraCare | Health-data permit/deny scenarios across care, quality improvement, and research. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/auroracare.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/auroracare.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/auroracare.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/auroracare.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/auroracare.md) |
| Barley Seed Lineage | Seed-lineage CAN/CAN'T reasoning for reproduction, dormancy, variation, and persistence. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/barley_seed_lineage.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/barley_seed_lineage.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/barley_seed_lineage.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/barley_seed_lineage.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/barley_seed_lineage.md) |
| Bayes Diagnosis | Bayesian posterior ranking of possible diseases from symptoms and test evidence. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/bayes_diagnosis.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/bayes_diagnosis.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/bayes_diagnosis.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/bayes_diagnosis.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/bayes_diagnosis.md) |
| Bayes Therapy Decision Support | Posterior-weighted utility selection of the best therapy. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/bayes_therapy.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/bayes_therapy.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/bayes_therapy.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/bayes_therapy.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/bayes_therapy.md) |
| BMI — Body Mass Index | Adult BMI calculation, category assignment, and healthy-weight interval. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/bmi.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/bmi.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/bmi.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/bmi.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/bmi.md) |
| Digital Product Passport | Component roll-up for recycled content, carbon footprint, repairability, and critical materials. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/digital_product_passport.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/digital_product_passport.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/digital_product_passport.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/digital_product_passport.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/digital_product_passport.md) |
| E-Bike Motor Thermal Envelope | Certified e-bike motor-temperature envelope for an assist plan. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/ebike_motor_thermal_envelope.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/ebike_motor_thermal_envelope.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/ebike_motor_thermal_envelope.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/ebike_motor_thermal_envelope.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/ebike_motor_thermal_envelope.md) |
| Gravity Mediator Witness | Mediator-only entanglement witness contrasting non-classical and purely classical gravitational mediators. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/gravity_mediator_witness.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/gravity_mediator_witness.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/gravity_mediator_witness.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/gravity_mediator_witness.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/gravity_mediator_witness.md) |
| LLD — Leg Length Discrepancy Measurement | Leg-length discrepancy measurement and alarm thresholding. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/lldm.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/lldm.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/lldm.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/lldm.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/lldm.md) |
| Photosynthetic Exciton Transfer | CAN/CAN'T reasoning for tuned versus detuned exciton delivery to a reaction center. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/photosynthetic_exciton_transfer.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/photosynthetic_exciton_transfer.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/photosynthetic_exciton_transfer.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/photosynthetic_exciton_transfer.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/photosynthetic_exciton_transfer.md) |
| RC Discharge Envelope | Certified exponential decay envelope for an RC discharge. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/rc_discharge_envelope.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/rc_discharge_envelope.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/rc_discharge_envelope.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/rc_discharge_envelope.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/rc_discharge_envelope.md) |
| Superdense Coding | Quantum-information parity facts for superdense coding. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/superdense_coding.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/superdense_coding.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/superdense_coding.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/superdense_coding.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/superdense_coding.md) |

### Technology

| Example | Description | Input | Go | Checks | Output | Doc |
|---|---|---|---|---|---|---|
| Alarm Bit Interoperability | Classical alarm-bit copy and relay tasks contrasted with forbidden universal cloning. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/alarm_bit_interoperability.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/alarm_bit_interoperability.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/alarm_bit_interoperability.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/alarm_bit_interoperability.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/alarm_bit_interoperability.md) |
| Deep Taxonomy 100000 | Large taxonomy materialization benchmark through a very deep class chain. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/deep_taxonomy_100000.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/deep_taxonomy_100000.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/deep_taxonomy_100000.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/deep_taxonomy_100000.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/deep_taxonomy_100000.md) |
| Delfour | Privacy-preserving retail insight and recommendation policy. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/delfour.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/delfour.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/delfour.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/delfour.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/delfour.md) |
| Doctor Advice Work Conflict | Policy conflict resolution for remote-work and office-work advice. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/doctor_advice_work_conflict.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/doctor_advice_work_conflict.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/doctor_advice_work_conflict.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/doctor_advice_work_conflict.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/doctor_advice_work_conflict.md) |
| FFT8 Numeric | Eight-point Fourier transform over a sampled sine wave with conjugate-bin and energy checks. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/fft8_numeric.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/fft8_numeric.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/fft8_numeric.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/fft8_numeric.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/fft8_numeric.md) |
| FFT32 Numeric | Thirty-two-point Fourier transform over several sampled waveforms with spectral invariant checks. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/fft32_numeric.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/fft32_numeric.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/fft32_numeric.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/fft32_numeric.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/fft32_numeric.md) |
| School Placement Route Audit | Route-aware audit of a straight-line school-placement support tool. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/school_placement_audit.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/school_placement_audit.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/school_placement_audit.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/school_placement_audit.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/school_placement_audit.md) |
| Gray Code Counter | n-bit Gray-code sequence with one-bit transition checks. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/gray_code_counter.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/gray_code_counter.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/gray_code_counter.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/gray_code_counter.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/gray_code_counter.md) |
| High Trust RDF Bloom Envelope | Bloom-envelope acceptance using canonical graph, index, and false-positive checks. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/high_trust_bloom_envelope.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/high_trust_bloom_envelope.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/high_trust_bloom_envelope.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/high_trust_bloom_envelope.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/high_trust_bloom_envelope.md) |
| Parcel Locker | Delegated parcel pickup-token authorization. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/parcellocker.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/parcellocker.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/parcellocker.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/parcellocker.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/parcellocker.md) |
| Path Discovery | Airport path discovery with stopover and routing constraints. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/path_discovery.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/path_discovery.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/path_discovery.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/path_discovery.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/path_discovery.md) |
| Ranked DPV Risk Report | ODRL/DPV clause risk ranking by severity and risk class. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/odrl_dpv_risk_ranked.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/odrl_dpv_risk_ranked.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/odrl_dpv_risk_ranked.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/odrl_dpv_risk_ranked.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/odrl_dpv_risk_ranked.md) |

### Engineering

| Example | Description | Input | Go | Checks | Output | Doc |
|---|---|---|---|---|---|---|
| Calidor | Municipal cooling intervention bundle chosen from active needs and budget constraints. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/calidor.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/calidor.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/calidor.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/calidor.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/calidor.md) |
| Complex Matrix Stability | Discrete-time stability classification using spectral radii of diagonal complex matrices. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/complex_matrix_stability.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/complex_matrix_stability.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/complex_matrix_stability.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/complex_matrix_stability.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/complex_matrix_stability.md) |
| Control System | Translated measurement and control rules for actuators, inputs, and disturbances. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/control_system.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/control_system.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/control_system.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/control_system.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/control_system.md) |
| Dijkstra Risk Path | Risk-adjusted path selection using weighted network edges. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/dijkstra_risk_path.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/dijkstra_risk_path.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/dijkstra_risk_path.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/dijkstra_risk_path.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/dijkstra_risk_path.md) |
| Docking Abort Token | Docking abort audit-token flow and safety-system copy restrictions. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/docking_abort_token.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/docking_abort_token.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/docking_abort_token.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/docking_abort_token.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/docking_abort_token.md) |
| Drone Corridor Planner | Constrained drone route planning through corridors and restricted zones. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/drone_corridor_planner.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/drone_corridor_planner.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/drone_corridor_planner.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/drone_corridor_planner.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/drone_corridor_planner.md) |
| EV Roadtrip Planner | EV route planning with battery, duration, cost, and comfort constraints. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/ev_roundtrip_planner.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/ev_roundtrip_planner.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/ev_roundtrip_planner.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/ev_roundtrip_planner.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/ev_roundtrip_planner.md) |
| Flandor | Regional retooling priority calculation for a Flanders scenario. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/flandor.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/flandor.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/flandor.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/flandor.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/flandor.md) |
| GPS — Goal driven route planning | Goal-driven route planning over a small road network. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/gps.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/gps.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/gps.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/gps.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/gps.md) |
| HarborSMR Insight Dispatch | Port electrolysis dispatch decision with safety margin and policy checks. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/harbor_smr.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/harbor_smr.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/harbor_smr.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/harbor_smr.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/harbor_smr.md) |
| Isolation Breach Token | Isolation-breach audit-token flow with cloning and fan-out restrictions. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/isolation_breach_token.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/isolation_breach_token.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/isolation_breach_token.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/isolation_breach_token.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/isolation_breach_token.md) |
| Wind Turbine Envelope | Wind-speed envelope classification with cubic power curve and interval energy audit. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/wind_turbine.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/wind_turbine.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/wind_turbine.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/wind_turbine.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/wind_turbine.md) |

### Mathematics

| Example | Description | Input | Go | Checks | Output | Doc |
|---|---|---|---|---|---|---|
| 8-Queens | 8-Queens constraint satisfaction with a valid board solution. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/queens.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/queens.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/queens.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/queens.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/queens.md) |
| Ackermann | Exact Ackermann and hyperoperation facts, including very large integer results. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/ackermann.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/ackermann.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/ackermann.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/ackermann.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/ackermann.md) |
| Allen Interval Calculus | Allen temporal interval relation closure over completed and explicit intervals. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/allen_interval_calculus.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/allen_interval_calculus.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/allen_interval_calculus.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/allen_interval_calculus.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/allen_interval_calculus.md) |
| Complex Numbers | Complex arithmetic and transcendental identity checks. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/complex_numbers.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/complex_numbers.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/complex_numbers.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/complex_numbers.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/complex_numbers.md) |
| Dining Philosophers | Chandy-Misra dining-philosophers trace with concurrency conflict checks. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/dining_philosophers.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/dining_philosophers.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/dining_philosophers.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/dining_philosophers.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/dining_philosophers.md) |
| Euler Identity Certificate | High-precision certificate for the identity exp(iπ) + 1 = 0. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/euler_identity_certificate.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/euler_identity_certificate.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/euler_identity_certificate.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/euler_identity_certificate.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/euler_identity_certificate.md) |
| Fibonacci Example (Big) | Exact computation of a large Fibonacci number. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/fibonacci.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/fibonacci.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/fibonacci.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/fibonacci.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/fibonacci.md) |
| Fundamental Theorem Arithmetic | Prime factorization and prime-power representation. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/fundamental_theorem_arithmetic.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/fundamental_theorem_arithmetic.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/fundamental_theorem_arithmetic.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/fundamental_theorem_arithmetic.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/fundamental_theorem_arithmetic.md) |
| Genetic Knapsack Selection | Deterministic genetic selection for a bounded knapsack. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/genetic_knapsack_selection.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/genetic_knapsack_selection.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/genetic_knapsack_selection.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/genetic_knapsack_selection.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/genetic_knapsack_selection.md) |
| Goldbach 1000 | Bounded strong-Goldbach checker for every even integer from 4 through 1000. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/goldbach_1000.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/goldbach_1000.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/goldbach_1000.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/goldbach_1000.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/goldbach_1000.md) |
| Kaprekar 6174 | Kaprekar chains and basin facts ending at 6174. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/kaprekar_6174.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/kaprekar_6174.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/kaprekar_6174.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/kaprekar_6174.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/kaprekar_6174.md) |
| Sudoku | Sudoku constraint solving with a unique completed grid. | [json](https://github.com/eyereasoner/eyelingo/blob/main/examples/input/sudoku.json) | [go](https://github.com/eyereasoner/eyelingo/blob/main/examples/sudoku.go) | [py](https://github.com/eyereasoner/eyelingo/blob/main/examples/checks/sudoku.py) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/output/sudoku.md) | [md](https://github.com/eyereasoner/eyelingo/blob/main/examples/doc/sudoku.md) |

