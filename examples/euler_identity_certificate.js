#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'euler_identity_certificate';

function cadd(a, b) { return { re: a.re + b.re, im: a.im + b.im }; }
function cmul(a, b) { return { re: a.re * b.re - a.im * b.im, im: a.re * b.im + a.im * b.re }; }
function cdivReal(a, x) { return { re: a.re / x, im: a.im / x }; }
function cabs(a) { return Math.hypot(a.re, a.im); }

function expI(angle, terms) {
  let total = { re: 0.0, im: 0.0 };
  let term = { re: 1.0, im: 0.0 };
  const z = { re: 0.0, im: angle };
  for (let n = 0; n < terms; n += 1) {
    if (n === 0) term = { re: 1.0, im: 0.0 };
    else term = cmul(term, cdivReal(z, n));
    total = cadd(total, term);
  }
  return total;
}

function trustedDerivation(data) {
  const value = expI(Number(data.angle), Number.parseInt(data.terms, 10));
  const residual = cabs(cadd(value, { re: 1.0, im: 0.0 }));
  fail('Euler identity derivation failed', {
    'terms count is positive': data.terms > 0,
    'residual satisfies expectation': (residual < data.tolerance) === data.expected.residualBelowTolerance,
    'real part is near minus one': Math.abs(value.re + 1) < data.tolerance,
    'imaginary part is near zero': Math.abs(value.im) < data.tolerance,
  });
  return { value, residual };
}

function main() {
  const data = loadInput(NAME);
  const { value, residual } = trustedDerivation(data);

  emit('# Euler Identity Certificate');
  emit();
  emit('## Insight');
  emit('expression : exp(iπ) + 1');
  emit(`terms used : ${data.terms}`);
  emit(`computed real part of exp(iπ) : ${value.re.toFixed(15)}`);
  emit(`computed imaginary part of exp(iπ) : ${value.im.toFixed(15)}`);
  emit(`residual magnitude : ${residual.toExponential(3)}`);
  emit(`within tolerance : ${String(residual < data.tolerance)}`);
  emit();
  emit('## Explanation');
  emit('The example approximates exp(iπ) by a finite Taylor series over complex numbers.');
  emit('The resulting residual is not claimed to be mathematically exact zero; it is checked against the explicit tolerance from JSON.');
  emit('The computed real part is effectively -1 and the imaginary part is near 0 at the chosen precision.');
  emit('That gives a reproducible finite certificate for the familiar Euler-identity witness.');
}

if (require.main === module) main();
