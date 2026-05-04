# see

**see** is a small collection of executable reasoning examples built around one repeatable shape:

```text
plain data + small rules + a trust gate -> insight + explanation
```

Each example is ordinary Node.js JavaScript. It loads plain JSON from `examples/input/`, computes with small named functions, verifies the assumptions that make the result safe to emit, and then prints only the user-facing Markdown sections `## Insight` and `## Explanation`.

The point is not to turn JavaScript into an oracle. The point is to make reasoning examples inspectable: the data, calculations, trust checks, and emitted explanation all live in files that can be read, run, and diffed.

## Why this project exists

Many reasoning demos are persuasive after they print an answer, but vague about what was checked before the answer appeared. **see** keeps every example compact enough to review while making the obligations explicit in code.

The project pattern is:

```text
JSON input
   ↓
JavaScript rules
   ↓
trustedDerivation(...)
   ↓
## Insight
## Explanation
```

The trust gate is executable verification, not magic truth. It collects the facts that must hold before an insight and explanation is emitted: input shape, unit conversions, tolerance bounds, authorization checks, route constraints, conservation checks, snapshot expectations, or domain-specific invariants. Human review still matters, but reviewers get concrete code and fixtures instead of a black-box answer.

## Quick start

Run every example and compare each generated Markdown result with its snapshot:

```sh
npm test
```

Run one example directly:

```sh
node examples/delfour.js
```

Manually compare one generated output:

```sh
node examples/bmi.js > /tmp/bmi.md
diff -u examples/output/bmi.md /tmp/bmi.md
```

The test script prints green `OK`, red `FAIL`, light gray per-example timings, and a final total with the number of examples run and elapsed time.

## Repository layout

```text
examples/
  <name>.js             executable example
  _see.js               shared helper for loading JSON, failing checks, and emitting Markdown
  input/<name>.json     structured input data
  output/<name>.md      expected insight + explanation snapshot
  doc/<name>.md         short per-example note and file index
```

## Current examples

The examples below are runnable files in `examples/`. Each row links to the short note, executable JavaScript, input fixture, and expected Markdown output.

| Example | What it demonstrates | Files |
|---|---|---|
| Bayes Diagnosis | Posterior probabilities are recomputed from priors, likelihoods, and present/absent evidence. | [doc](examples/doc/bayes_diagnosis.md), [js](examples/bayes_diagnosis.js), [input](examples/input/bayes_diagnosis.json), [output](examples/output/bayes_diagnosis.md) |
| Bayes Therapy Decision Support | Bayesian posteriors feed a small expected-utility therapy ranking. | [doc](examples/doc/bayes_therapy.md), [js](examples/bayes_therapy.js), [input](examples/input/bayes_therapy.json), [output](examples/output/bayes_therapy.md) |
| BMI — Body Mass Index example | Unit normalization, BMI calculation, category boundaries, and healthy-weight range. | [doc](examples/doc/bmi.md), [js](examples/bmi.js), [input](examples/input/bmi.json), [output](examples/output/bmi.md) |
| Complex Matrix Stability | Diagonal complex matrices are classified by spectral radius for discrete-time stability. | [doc](examples/doc/complex_matrix_stability.md), [js](examples/complex_matrix_stability.js), [input](examples/input/complex_matrix_stability.json), [output](examples/output/complex_matrix_stability.md) |
| Delfour | A shopping insight is emitted only if authorization, minimization, payload hash, signature metadata, duty timing, and product choice all hold. | [doc](examples/doc/delfour.md), [js](examples/delfour.js), [input](examples/input/delfour.json), [output](examples/output/delfour.md) |
| Digital Product Passport | Component, material, document, lifecycle, and footprint facts are folded into a public passport decision. | [doc](examples/doc/digital_product_passport.md), [js](examples/digital_product_passport.js), [input](examples/input/digital_product_passport.json), [output](examples/output/digital_product_passport.md) |
| Dijkstra Risk Path | A weighted shortest path balances route cost and risk penalty. | [doc](examples/doc/dijkstra_risk_path.md), [js](examples/dijkstra_risk_path.js), [input](examples/input/dijkstra_risk_path.json), [output](examples/output/dijkstra_risk_path.md) |
| Eco Route Insight | A local route insight signs a minimal envelope instead of exporting raw logistics data. | [doc](examples/doc/eco_route_insight.md), [js](examples/eco_route_insight.js), [input](examples/input/eco_route_insight.json), [output](examples/output/eco_route_insight.md) |
| Euler Identity Certificate | A finite Taylor-series calculation gives a numerical certificate for exp(iπ)+1. | [doc](examples/doc/euler_identity_certificate.md), [js](examples/euler_identity_certificate.js), [input](examples/input/euler_identity_certificate.json), [output](examples/output/euler_identity_certificate.md) |
| EV Roadtrip Planner | A bounded EV planner composes route actions and selects the fastest candidate that satisfies reliability, cost, and duration thresholds. | [doc](examples/doc/ev_roundtrip_planner.md), [js](examples/ev_roundtrip_planner.js), [input](examples/input/ev_roundtrip_planner.json), [output](examples/output/ev_roundtrip_planner.md) |
| FFT8 Numeric | A discrete Fourier transform identifies the dominant frequency bins of an eight-sample sine wave. | [doc](examples/doc/fft8_numeric.md), [js](examples/fft8_numeric.js), [input](examples/input/fft8_numeric.json), [output](examples/output/fft8_numeric.md) |
| Fibonacci Example (Big) | Exact arbitrary-precision Fibonacci computation for index 10000. | [doc](examples/doc/fibonacci.md), [js](examples/fibonacci.js), [input](examples/input/fibonacci.json), [output](examples/output/fibonacci.md) |
| Fundamental Theorem Arithmetic | Trial division factors several integers and checks multiplicative reconstruction and prime powers. | [doc](examples/doc/fundamental_theorem_arithmetic.md), [js](examples/fundamental_theorem_arithmetic.js), [input](examples/input/fundamental_theorem_arithmetic.json), [output](examples/output/fundamental_theorem_arithmetic.md) |
| Genetic Knapsack Selection | A deterministic one-bit mutation search reaches a feasible local optimum for a knapsack genome. | [doc](examples/doc/genetic_knapsack_selection.md), [js](examples/genetic_knapsack_selection.js), [input](examples/input/genetic_knapsack_selection.json), [output](examples/output/genetic_knapsack_selection.md) |
| Goldbach 1000 | Every even integer from 4 through 1000 is given a prime-sum witness. | [doc](examples/doc/goldbach_1000.md), [js](examples/goldbach_1000.js), [input](examples/input/goldbach_1000.json), [output](examples/output/goldbach_1000.md) |
| GPS — Goal driven route planning | A small route planner compares duration, cost, belief, and comfort for two Gent-to-Oostende routes. | [doc](examples/doc/gps.md), [js](examples/gps.js), [input](examples/input/gps.json), [output](examples/output/gps.md) |
| Gray Code Counter | A cyclic 4-bit Gray counter visits all states while flipping one bit per transition. | [doc](examples/doc/gray_code_counter.md), [js](examples/gray_code_counter.js), [input](examples/input/gray_code_counter.json), [output](examples/output/gray_code_counter.md) |
| Kaprekar 6174 | Four-digit Kaprekar chains are generated, counted, and checked against the 6174 attractor. | [doc](examples/doc/kaprekar_6174.md), [js](examples/kaprekar_6174.js), [input](examples/input/kaprekar_6174.json), [output](examples/output/kaprekar_6174.md) |
| ODRL + DPV Risk Ranking | ODRL-style policy clauses are checked against consumer needs and emitted as a ranked DPV risk report with mitigation advice. | [doc](examples/doc/odrl_dpv_risk_ranked.md), [js](examples/odrl_dpv_risk_ranked.js), [input](examples/input/odrl_dpv_risk_ranked.json), [output](examples/output/odrl_dpv_risk_ranked.md) |
| Path Discovery | A bounded airport-graph query finds all simple routes with at most two stopovers. | [doc](examples/doc/path_discovery.md), [js](examples/path_discovery.js), [input](examples/input/path_discovery.json), [output](examples/output/path_discovery.md) |
| 8-Queens | Bit-mask search prints one board while still counting all 92 solutions. | [doc](examples/doc/queens.md), [js](examples/queens.js), [input](examples/input/queens.json), [output](examples/output/queens.md) |
| RC Discharge Envelope | A finite decay interval certifies when an RC capacitor envelope falls below tolerance. | [doc](examples/doc/rc_discharge_envelope.md), [js](examples/rc_discharge_envelope.js), [input](examples/input/rc_discharge_envelope.json), [output](examples/output/rc_discharge_envelope.md) |
| School Placement Route Audit | A straight-line assignment rule is audited against walking routes, barriers, and preferences. | [doc](examples/doc/school_placement_audit.md), [js](examples/school_placement_audit.js), [input](examples/input/school_placement_audit.json), [output](examples/output/school_placement_audit.md) |
| Sudoku | A completed Sudoku grid is emitted only after clue preservation and row, column, and box constraints hold. | [doc](examples/doc/sudoku.md), [js](examples/sudoku.js), [input](examples/input/sudoku.json), [output](examples/output/sudoku.md) |
| Wind Turbine Envelope | Wind-speed intervals are classified against turbine thresholds and accumulated into energy. | [doc](examples/doc/wind_turbine.md), [js](examples/wind_turbine.js), [input](examples/input/wind_turbine.json), [output](examples/output/wind_turbine.md) |

## Example design pattern

Each example should follow this shape:

1. Load only its JSON input.
2. Compute the insight with small, named functions.
3. Define a `trustedDerivation(...)` function that checks the obligations that make the explanation safe to emit.
4. Print only `## Insight` and `## Explanation`.
5. Keep expected output in `examples/output/<name>.md`.
6. Add a short note in `examples/doc/<name>.md`.

A typical trust gate looks like this:

```javascript
const obligations = {
  "normalizer is positive": total > 0,
  "posteriors sum to one": Math.abs(sum(posteriors) - 1.0) < 1e-12,
  "winner is stable": winner === "COVID19",
};

const failed = Object.entries(obligations)
  .filter(([, ok]) => !ok)
  .map(([name]) => name);
if (failed.length) {
  throw new Error("derivation failed: " + failed.join(", "));
}
```

## Adding an example

To add a new example, create the matching files as a set:

```text
examples/<name>.js
examples/input/<name>.json
examples/output/<name>.md
examples/doc/<name>.md
```

Then run `npm test`. The test runner discovers every top-level `examples/*.js` file except `_see.js`, so no separate registration step is needed.
