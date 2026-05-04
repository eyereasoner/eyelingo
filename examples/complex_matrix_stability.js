#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'complex_matrix_stability';

function modulus(z) {
  return Math.hypot(Number(z.re), Number(z.im));
}

function fmtRadius(x) {
  return Math.abs(x - Math.round(x)) < 1e-12 ? String(Math.trunc(x)) : Number(x).toPrecision(6).replace(/\.?(0+)$/, '');
}

function classify(radius) {
  if (radius < 1 - 1e-12) return 'damped';
  if (Math.abs(radius - 1) <= 1e-12) return 'marginally stable';
  return 'unstable';
}

function trustedDerivation(data) {
  const rows = [];
  for (const m of data.matrices) {
    const radius = Math.max(...m.diagonal.map(modulus));
    rows.push([m.name, radius, classify(radius)]);
  }
  const z = data.sampleProduct.z;
  const w = data.sampleProduct.w;
  const productSq = (z.re * w.re - z.im * w.im) ** 2 + (z.re * w.im + z.im * w.re) ** 2;
  fail('Complex matrix stability derivation failed', {
    'three matrices are classified': rows.length === 3,
    'unstable case has radius two': JSON.stringify(rows[0]) === JSON.stringify(['A_unstable', 2.0, 'unstable']),
    'stable case has radius one': JSON.stringify(rows[1]) === JSON.stringify(['A_stable', 1.0, 'marginally stable']),
    'damped case has radius zero': JSON.stringify(rows[2]) === JSON.stringify(['A_damped', 0.0, 'damped']),
    'complex product modulus squares multiply': Math.abs(productSq - (modulus(z) ** 2) * (modulus(w) ** 2)) < 1e-12,
  });
  return rows;
}

function main() {
  const data = loadInput(NAME);
  const rows = trustedDerivation(data);

  emit('# Complex Matrix Stability');
  emit();
  emit('## Insight');
  for (const [name, radius, cls] of rows) {
    emit(`${name} : spectral radius ${fmtRadius(radius)} -> ${cls}`);
  }
  emit();
  emit('## Explanation');
  emit('For a discrete-time linear system x_{k+1} = A x_k, the eigenvalues of A govern the modal behaviour.');
  emit('Because the matrices are diagonal, the eigenvalues are the diagonal entries; the largest modulus gives the spectral radius and therefore the stability class.');
}

if (require.main === module) main();
