#!/usr/bin/env node
const { emit, fail, loadInput, sum } = require('./_see');

const NAME = 'fft8_numeric';

function cadd(a, b) { return { re: a.re + b.re, im: a.im + b.im }; }
function cmul(a, b) { return { re: a.re * b.re - a.im * b.im, im: a.re * b.im + a.im * b.re }; }
function cabs(a) { return Math.hypot(a.re, a.im); }
function cphase(a) { return Math.atan2(a.im, a.re); }

function dft(samples) {
  const n = samples.length;
  const bins = [];
  for (let k = 0; k < n; k += 1) {
    let total = { re: 0.0, im: 0.0 };
    for (let j = 0; j < samples.length; j += 1) {
      const x = samples[j];
      const angle = -2 * Math.PI * k * j / n;
      total = cadd(total, cmul({ re: x, im: 0.0 }, { re: Math.cos(angle), im: Math.sin(angle) }));
    }
    bins.push(total);
  }
  return bins;
}

function trustedDerivation(data) {
  const samples = data.samples.map(Number);
  const bins = dft(samples);
  const magnitudes = bins.map(cabs);
  const maxMag = Math.max(...magnitudes);
  const tol = Number(data.expected.tolerance);
  const dominant = magnitudes.map((mag, k) => Math.abs(mag - maxMag) <= tol ? k : null).filter((x) => x !== null);
  const energyTime = sum(samples.map((x) => x * x));
  const energyFreq = sum(bins.map((value) => cabs(value) ** 2)) / samples.length;
  fail('FFT8 derivation failed', {
    'input has eight samples': samples.length === 8,
    'dominant bins are expected': JSON.stringify(dominant) === JSON.stringify(data.expected.dominantBins),
    'DC component cancels': cabs(bins[0]) <= tol,
    'sine bins have magnitude four': Math.abs(cabs(bins[1]) - 4.0) < 1e-9 && Math.abs(cabs(bins[7]) - 4.0) < 1e-9,
    'Parseval energy is preserved': Math.abs(energyTime - energyFreq) < 1e-9,
    'real signal gives conjugate symmetry': Math.abs(bins[1].re - bins[7].re) < 1e-9 && Math.abs(bins[1].im + bins[7].im) < 1e-9,
  });
  return { bins, dominant, energyTime, energyFreq };
}

function main() {
  const data = loadInput(NAME);
  const samples = data.samples.map(Number);
  const { bins, dominant, energyTime, energyFreq } = trustedDerivation(data);

  emit('# FFT8 Numeric');
  emit();
  emit('## Insight');
  emit(`sample vector : ${samples.map((x) => x.toFixed(6)).join(', ')}`);
  emit(`dominant bins : ${dominant.map((k) => `k=${k} magnitude=${cabs(bins[k]).toFixed(6)} phase=${cphase(bins[k]).toFixed(6)}`).join('; ')}`);
  emit(`time-domain energy : ${energyTime.toFixed(6)}`);
  emit(`frequency-domain energy / 8 : ${energyFreq.toFixed(6)}`);
  emit();
  emit('## Explanation');
  emit('The input samples describe one sine cycle over eight equally spaced samples.');
  emit('The DFT projects the signal onto eight complex roots of unity.');
  emit('A real sine wave has equal magnitude at the positive and negative frequency bins.');
  emit('All non-dominant bins cancel to zero within the configured numerical tolerance.');
}

if (require.main === module) main();
