# see

See is an experiment in a simpler style of checkable reasoning example:

```text
plain data + small rules + a trust gate -> Answer + Reason
```

Each example is written in ordinary JavaScript for Node.js. The input is plain JSON. The code is intentionally shaped like an executable specification: facts are loaded from `examples/input/`, rules are small functions, and a `trustedDerivation(...)` gate must succeed before the program emits the visible explanation.

The goal is not to claim that one JavaScript program is a perfect oracle. The goal is to make the assumptions, calculations, and emitted explanation inspectable in one widely familiar language.

## Why this project exists

Many reasoning examples are easy to print but harder to trust. See keeps the example small enough to read while making the obligations explicit in code.

The compact shape is:

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

The trust gate is executable scrutiny, not magic truth. It makes assumptions and failure points visible before the explanation is emitted. Human review still matters, but the example gives the reviewer something concrete to inspect.

## Run

Run all examples:

```sh
npm test
```

Run one example:

```sh
node examples/delfour.js
```

Compare manually:

```sh
node examples/bmi.js > /tmp/bmi.md
diff -u examples/output/bmi.md /tmp/bmi.md
```

The test script prints green `OK`, red `FAIL`, light gray per-example timings, and a final total with the number of examples run and elapsed time.

## Repository layout

```text
examples/
  <name>.js             executable example
  _see.js               tiny shared helper for loading JSON and emitting Markdown
  input/<name>.json     structured input data
  output/<name>.md      expected Answer + Reason snapshot
  doc/<name>.md         short per-example note
```

## Current examples

See contains the following runnable examples:

| Example | What it demonstrates | Files |
|---|---|---|
| [Bayes Diagnosis](examples/doc/bayes_diagnosis.md) | Posterior probabilities are recomputed from priors, likelihoods, and present/absent evidence. | [js](examples/bayes_diagnosis.js), [input](examples/input/bayes_diagnosis.json), [output](examples/output/bayes_diagnosis.md) |
| [Bayes Therapy Decision Support](examples/doc/bayes_therapy.md) | Bayesian posteriors feed a small expected-utility therapy ranking. | [js](examples/bayes_therapy.js), [input](examples/input/bayes_therapy.json), [output](examples/output/bayes_therapy.md) |
| [BMI — Body Mass Index example](examples/doc/bmi.md) | Unit normalization, BMI calculation, category boundaries, and healthy-weight range. | [js](examples/bmi.js), [input](examples/input/bmi.json), [output](examples/output/bmi.md) |
| [Complex Matrix Stability](examples/doc/complex_matrix_stability.md) | Diagonal complex matrices are classified by spectral radius for discrete-time stability. | [js](examples/complex_matrix_stability.js), [input](examples/input/complex_matrix_stability.json), [output](examples/output/complex_matrix_stability.md) |
| [Delfour](examples/doc/delfour.md) | A shopping insight is emitted only if authorization, minimization, payload hash, signature metadata, duty timing, and product choice all hold. | [js](examples/delfour.js), [input](examples/input/delfour.json), [output](examples/output/delfour.md) |
| [Digital Product Passport](examples/doc/digital_product_passport.md) | Component, material, document, lifecycle, and footprint facts are folded into a public passport decision. | [js](examples/digital_product_passport.js), [input](examples/input/digital_product_passport.json), [output](examples/output/digital_product_passport.md) |
| [Dijkstra Risk Path](examples/doc/dijkstra_risk_path.md) | A weighted shortest path balances route cost and risk penalty. | [js](examples/dijkstra_risk_path.js), [input](examples/input/dijkstra_risk_path.json), [output](examples/output/dijkstra_risk_path.md) |
| [Eco Route Insight](examples/doc/eco_route_insight.md) | A local route insight signs a minimal envelope instead of exporting raw logistics data. | [js](examples/eco_route_insight.js), [input](examples/input/eco_route_insight.json), [output](examples/output/eco_route_insight.md) |
| [Euler Identity Certificate](examples/doc/euler_identity_certificate.md) | A finite Taylor-series calculation gives a numerical certificate for exp(iπ)+1. | [js](examples/euler_identity_certificate.js), [input](examples/input/euler_identity_certificate.json), [output](examples/output/euler_identity_certificate.md) |
| [EV Roadtrip Planner](examples/doc/ev_roundtrip_planner.md) | A bounded EV planner composes route actions and selects the fastest candidate that satisfies reliability, cost, and duration thresholds. | [js](examples/ev_roundtrip_planner.js), [input](examples/input/ev_roundtrip_planner.json), [output](examples/output/ev_roundtrip_planner.md) |
| [FFT8 Numeric](examples/doc/fft8_numeric.md) | A discrete Fourier transform identifies the dominant frequency bins of an eight-sample sine wave. | [js](examples/fft8_numeric.js), [input](examples/input/fft8_numeric.json), [output](examples/output/fft8_numeric.md) |
| [Fibonacci Example (Big)](examples/doc/fibonacci.md) | Exact arbitrary-precision Fibonacci computation for index 10000. | [js](examples/fibonacci.js), [input](examples/input/fibonacci.json), [output](examples/output/fibonacci.md) |
| [Fundamental Theorem Arithmetic](examples/doc/fundamental_theorem_arithmetic.md) | Trial division factors several integers and checks multiplicative reconstruction and prime powers. | [js](examples/fundamental_theorem_arithmetic.js), [input](examples/input/fundamental_theorem_arithmetic.json), [output](examples/output/fundamental_theorem_arithmetic.md) |
| [Genetic Knapsack Selection](examples/doc/genetic_knapsack_selection.md) | A deterministic one-bit mutation search reaches a feasible local optimum for a knapsack genome. | [js](examples/genetic_knapsack_selection.js), [input](examples/input/genetic_knapsack_selection.json), [output](examples/output/genetic_knapsack_selection.md) |
| [Goldbach 1000](examples/doc/goldbach_1000.md) | Every even integer from 4 through 1000 is given a prime-sum witness. | [js](examples/goldbach_1000.js), [input](examples/input/goldbach_1000.json), [output](examples/output/goldbach_1000.md) |
| [GPS — Goal driven route planning](examples/doc/gps.md) | A small route planner compares duration, cost, belief, and comfort for two Gent-to-Oostende routes. | [js](examples/gps.js), [input](examples/input/gps.json), [output](examples/output/gps.md) |
| [Gray Code Counter](examples/doc/gray_code_counter.md) | A cyclic 4-bit Gray counter visits all states while flipping one bit per transition. | [js](examples/gray_code_counter.js), [input](examples/input/gray_code_counter.json), [output](examples/output/gray_code_counter.md) |
| [Kaprekar 6174](examples/doc/kaprekar_6174.md) | Four-digit Kaprekar chains are generated, counted, and checked against the 6174 attractor. | [js](examples/kaprekar_6174.js), [input](examples/input/kaprekar_6174.json), [output](examples/output/kaprekar_6174.md) |
| [Path Discovery](examples/doc/path_discovery.md) | A bounded airport-graph query finds all simple routes with at most two stopovers. | [js](examples/path_discovery.js), [input](examples/input/path_discovery.json), [output](examples/output/path_discovery.md) |
| [8-Queens](examples/doc/queens.md) | Bit-mask search prints one board while still counting all 92 solutions. | [js](examples/queens.js), [input](examples/input/queens.json), [output](examples/output/queens.md) |
| [RC Discharge Envelope](examples/doc/rc_discharge_envelope.md) | A finite decay interval certifies when an RC capacitor envelope falls below tolerance. | [js](examples/rc_discharge_envelope.js), [input](examples/input/rc_discharge_envelope.json), [output](examples/output/rc_discharge_envelope.md) |
| [School Placement Route Audit](examples/doc/school_placement_audit.md) | A straight-line assignment rule is audited against walking routes, barriers, and preferences. | [js](examples/school_placement_audit.js), [input](examples/input/school_placement_audit.json), [output](examples/output/school_placement_audit.md) |
| [Sudoku](examples/doc/sudoku.md) | A completed Sudoku grid is emitted only after clue preservation and row, column, and box constraints hold. | [js](examples/sudoku.js), [input](examples/input/sudoku.json), [output](examples/output/sudoku.md) |
| [Wind Turbine Envelope](examples/doc/wind_turbine.md) | Wind-speed intervals are classified against turbine thresholds and accumulated into energy. | [js](examples/wind_turbine.js), [input](examples/input/wind_turbine.json), [output](examples/output/wind_turbine.md) |

## Example design pattern

Each example should follow this shape:

1. Load only its JSON input.
2. Compute the answer with small, named functions.
3. Define a `trustedDerivation(...)` function that checks the obligations that make the explanation safe to emit.
4. Print only `## Insight` and `## Explanation`.
5. Keep expected output in `examples/output/<name>.md`.

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
