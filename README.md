# SEE

**Structured Evidence Explanation**

SEE is an approach for turning input facts into gated insights with human-readable explanations.

```text
input facts -> rules -> trust gate -> insight + explanation
```

Each example is a small Node.js program:

- `examples/input/<name>.json` contains the input facts.
- `examples/<name>.js` derives the insight and checks the trust gate.
- `examples/output/<name>.md` is the expected Markdown output.
- `examples/doc/<name>.md` explains the example.

The trust gate is executable verification. If a required fact is missing, the program fails instead of emitting an unsupported insight.

SEE was largely inspired by Prof. Ruben Verborgh's essay [Inside the Insight Economy](https://ruben.verborgh.org/blog/2025/08/12/inside-the-insight-economy/).

## Run

Run all examples:

```sh
npm test
```

Run one example:

```sh
node examples/delfour.js
```

Compare one output:

```sh
node examples/bmi.js > /tmp/bmi.md
diff -u examples/output/bmi.md /tmp/bmi.md
```

## Files

```text
examples/
  <name>.js             executable example
  _see.js               shared helper
  input/<name>.json     input facts
  output/<name>.md      expected insight and explanation
  doc/<name>.md         guide
```

## Current examples

| Example | Idea | Files |
|---|---|---|
| Bayes diagnosis | posterior probability from prior, likelihood, and evidence | [doc](examples/doc/bayes_diagnosis.md), [js](examples/bayes_diagnosis.js), [input](examples/input/bayes_diagnosis.json), [output](examples/output/bayes_diagnosis.md) |
| Bayes therapy | expected-utility ranking from Bayesian posteriors | [doc](examples/doc/bayes_therapy.md), [js](examples/bayes_therapy.js), [input](examples/input/bayes_therapy.json), [output](examples/output/bayes_therapy.md) |
| BMI | unit normalization, BMI category, and healthy-weight range | [doc](examples/doc/bmi.md), [js](examples/bmi.js), [input](examples/input/bmi.json), [output](examples/output/bmi.md) |
| Complex matrix stability | spectral-radius check for discrete-time stability | [doc](examples/doc/complex_matrix_stability.md), [js](examples/complex_matrix_stability.js), [input](examples/input/complex_matrix_stability.json), [output](examples/output/complex_matrix_stability.md) |
| Delfour | authorized, minimized shopping insight | [doc](examples/doc/delfour.md), [js](examples/delfour.js), [input](examples/input/delfour.json), [output](examples/output/delfour.md) |
| Digital product passport | component and footprint facts folded into a public passport decision | [doc](examples/doc/digital_product_passport.md), [js](examples/digital_product_passport.js), [input](examples/input/digital_product_passport.json), [output](examples/output/digital_product_passport.md) |
| Dijkstra risk path | shortest path with a risk penalty | [doc](examples/doc/dijkstra_risk_path.md), [js](examples/dijkstra_risk_path.js), [input](examples/input/dijkstra_risk_path.json), [output](examples/output/dijkstra_risk_path.md) |
| Eco route insight | local route insight without exporting raw logistics data | [doc](examples/doc/eco_route_insight.md), [js](examples/eco_route_insight.js), [input](examples/input/eco_route_insight.json), [output](examples/output/eco_route_insight.md) |
| Euler identity certificate | finite numerical certificate for exp(iπ)+1 | [doc](examples/doc/euler_identity_certificate.md), [js](examples/euler_identity_certificate.js), [input](examples/input/euler_identity_certificate.json), [output](examples/output/euler_identity_certificate.md) |
| EV roadtrip planner | fastest route that satisfies cost, time, and reliability limits | [doc](examples/doc/ev_roundtrip_planner.md), [js](examples/ev_roundtrip_planner.js), [input](examples/input/ev_roundtrip_planner.json), [output](examples/output/ev_roundtrip_planner.md) |
| FFT8 numeric | dominant frequency bins in eight samples | [doc](examples/doc/fft8_numeric.md), [js](examples/fft8_numeric.js), [input](examples/input/fft8_numeric.json), [output](examples/output/fft8_numeric.md) |
| Fibonacci | exact big-integer Fibonacci computation | [doc](examples/doc/fibonacci.md), [js](examples/fibonacci.js), [input](examples/input/fibonacci.json), [output](examples/output/fibonacci.md) |
| Fundamental theorem arithmetic | prime factors checked by reconstruction | [doc](examples/doc/fundamental_theorem_arithmetic.md), [js](examples/fundamental_theorem_arithmetic.js), [input](examples/input/fundamental_theorem_arithmetic.json), [output](examples/output/fundamental_theorem_arithmetic.md) |
| Genetic knapsack selection | deterministic one-bit mutation search for a feasible knapsack | [doc](examples/doc/genetic_knapsack_selection.md), [js](examples/genetic_knapsack_selection.js), [input](examples/input/genetic_knapsack_selection.json), [output](examples/output/genetic_knapsack_selection.md) |
| Goldbach 1000 | prime-sum witnesses for even numbers up to 1000 | [doc](examples/doc/goldbach_1000.md), [js](examples/goldbach_1000.js), [input](examples/input/goldbach_1000.json), [output](examples/output/goldbach_1000.md) |
| GPS route planning | route choice from duration, cost, belief, and comfort | [doc](examples/doc/gps.md), [js](examples/gps.js), [input](examples/input/gps.json), [output](examples/output/gps.md) |
| Gray code counter | 4-bit cycle with one-bit transitions | [doc](examples/doc/gray_code_counter.md), [js](examples/gray_code_counter.js), [input](examples/input/gray_code_counter.json), [output](examples/output/gray_code_counter.md) |
| Kaprekar 6174 | four-digit chains checked against the 6174 attractor | [doc](examples/doc/kaprekar_6174.md), [js](examples/kaprekar_6174.js), [input](examples/input/kaprekar_6174.json), [output](examples/output/kaprekar_6174.md) |
| ODRL + DPV risk ranking | policy clauses turned into ranked risks and mitigations | [doc](examples/doc/odrl_dpv_risk_ranked.md), [js](examples/odrl_dpv_risk_ranked.js), [input](examples/input/odrl_dpv_risk_ranked.json), [output](examples/output/odrl_dpv_risk_ranked.md) |
| Path discovery | bounded airport routes with at most two stopovers | [doc](examples/doc/path_discovery.md), [js](examples/path_discovery.js), [input](examples/input/path_discovery.json), [output](examples/output/path_discovery.md) |
| 8-Queens | one board plus the count of all 92 solutions | [doc](examples/doc/queens.md), [js](examples/queens.js), [input](examples/input/queens.json), [output](examples/output/queens.md) |
| RC discharge envelope | decay interval checked against a voltage tolerance | [doc](examples/doc/rc_discharge_envelope.md), [js](examples/rc_discharge_envelope.js), [input](examples/input/rc_discharge_envelope.json), [output](examples/output/rc_discharge_envelope.md) |
| School placement audit | assignment checked against route barriers and preferences | [doc](examples/doc/school_placement_audit.md), [js](examples/school_placement_audit.js), [input](examples/input/school_placement_audit.json), [output](examples/output/school_placement_audit.md) |
| Sudoku | grid emitted after row, column, box, and clue checks | [doc](examples/doc/sudoku.md), [js](examples/sudoku.js), [input](examples/input/sudoku.json), [output](examples/output/sudoku.md) |
| Wind turbine envelope | wind-speed intervals classified into energy output | [doc](examples/doc/wind_turbine.md), [js](examples/wind_turbine.js), [input](examples/input/wind_turbine.json), [output](examples/output/wind_turbine.md) |

## Add an example

Create these files:

```text
examples/<name>.js
examples/input/<name>.json
examples/output/<name>.md
examples/doc/<name>.md
```

The JavaScript file should load the JSON, derive the insight, check the trust gate, and emit only:

```md
## Insight

...

## Explanation

...
```

Then run:

```sh
npm test
```
