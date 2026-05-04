#!/usr/bin/env node
const { emit, fail, loadInput, sum } = require('./_see');

const NAME = 'wind_turbine';

function classify(speed, cutIn, rated, cutOut, ratedPower) {
  if (speed < cutIn || speed >= cutOut) return ['stopped', 0.0];
  if (speed < rated) {
    const power = ratedPower * ((speed ** 3 - cutIn ** 3) / (rated ** 3 - cutIn ** 3));
    return ['partial', power];
  }
  return ['rated', ratedPower];
}

function trustedDerivation(data) {
  const rows = data.windSpeedsMS.map((speed, idx) => {
    const [status, power] = classify(speed, data.cutInMS, data.ratedMS, data.cutOutMS, data.ratedPowerMW);
    return [idx + 1, speed, status, power];
  });
  const usable = rows.filter((row) => row[2] !== 'stopped').length;
  const ratedCount = rows.filter((row) => row[2] === 'rated').length;
  const stopped = rows.filter((row) => row[2] === 'stopped').length;
  const energy = sum(rows.map((row) => row[3])) * data.intervalMinutes / 60.0;
  fail('Wind turbine derivation failed', {
    'one classification per wind sample': rows.length === data.windSpeedsMS.length,
    'intervals partition into stopped partial and rated': usable + stopped === rows.length && ratedCount <= usable,
    'all generated power is non-negative': rows.every((row) => row[3] >= 0),
    'all generated power is capped at rated power': rows.every((row) => row[3] <= data.ratedPowerMW),
  });
  return rows;
}

function main() {
  const data = loadInput(NAME);
  const rows = trustedDerivation(data);
  const energy = sum(rows.map((row) => row[3])) * data.intervalMinutes / 60.0;
  const usable = rows.filter((row) => row[2] !== 'stopped').length;
  const classifications = rows.map(([i, speed, status, power]) => `t${i} ${speed.toFixed(1)} m/s ${status} ${power.toFixed(3)} MW`).join('; ');

  emit('# Wind Turbine Envelope');
  emit();
  emit('## Insight');
  emit(`operating thresholds : cut-in ${data.cutInMS.toFixed(1)} m/s, rated ${data.ratedMS.toFixed(1)} m/s, cut-out ${data.cutOutMS.toFixed(1)} m/s`);
  emit(`rated power : ${data.ratedPowerMW.toFixed(1)} MW`);
  emit(`interval classifications : ${classifications}`);
  emit(`usable intervals : ${usable}`);
  emit(`total energy : ${energy.toFixed(3)} MWh`);
  emit();
  emit('## Explanation');
  emit('Wind below cut-in and at or above cut-out is stopped for production and safety.');
  emit('Wind between cut-in and rated speed follows a cubic power curve normalized to the rated point.');
  emit('Wind between rated speed and cut-out is capped at rated power.');
  emit('Energy is accumulated by multiplying each interval power by the ten-minute interval duration.');
}

if (require.main === module) main();
