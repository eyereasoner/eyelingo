#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'rc_discharge_envelope';

function upperVoltage(data, step) {
  return Number(data.initialVoltage) * (Number(data.decayUpper) ** step);
}

function firstSettled(data) {
  for (let step = 0; step <= Number.parseInt(data.maxStep, 10); step += 1) {
    if (upperVoltage(data, step) < Number(data.tolerance)) return step;
  }
  throw new Error('not settled');
}

function trustedDerivation(data) {
  const step = firstSettled(data);
  fail('RC discharge derivation failed', {
    'interval is positive and contracting': 0 < data.decayLower && data.decayLower <= data.decayUpper && data.decayUpper < 1,
    'previous step is not settled': upperVoltage(data, step - 1) >= data.tolerance,
    'settled step is below tolerance': upperVoltage(data, step) < data.tolerance,
    'witness occurs before max step': step <= data.maxStep,
  });
  return step;
}

function main() {
  const data = loadInput(NAME);
  const step = trustedDerivation(data);
  const time = step * Number(data.samplePeriod);
  const voltage = upperVoltage(data, step);

  emit('# RC Discharge Envelope');
  emit();
  emit('## Insight');
  emit(`exact decay symbol : ${data.exactDecaySymbol}`);
  emit(`certified decay interval : [${Number(data.decayLower).toFixed(10)}, ${Number(data.decayUpper).toFixed(10)}]`);
  emit(`first below tolerance step : ${step}`);
  emit(`first below tolerance time : ${time.toFixed(3)} s`);
  emit(`upper voltage at step ${step} : ${voltage.toFixed(6)} V`);
  emit();
  emit('## Explanation');
  emit('The physical decay factor is exp(-1/4), but the example uses a finite double interval as the certificate.');
  emit('Because the interval lies strictly between 0 and 1, the capacitor voltage envelope contracts each sample.');
  emit('The upper envelope is the safety-relevant bound: once it falls below 1.0 V, every compatible exact trajectory is below tolerance.');
  emit('The first such witness occurs before the configured maximum step.');
}

if (require.main === module) main();
