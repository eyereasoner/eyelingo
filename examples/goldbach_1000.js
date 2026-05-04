#!/usr/bin/env node
const { emit, fail, loadInput, range } = require('./_see');

const NAME = 'goldbach_1000';

function isPrime(n) {
  if (n < 2) return false;
  if (n === 2) return true;
  if (n % 2 === 0) return false;
  let d = 3;
  while (d * d <= n) {
    if (n % d === 0) return false;
    d += 2;
  }
  return true;
}

function witness(even) {
  for (let p = 2; p <= Math.trunc(even / 2); p += 1) {
    if (isPrime(p) && isPrime(even - p)) return [p, even - p];
  }
  throw new Error(`no Goldbach witness for ${even}`);
}

function trustedDerivation(data) {
  const maxEven = Number.parseInt(data.maxEven, 10);
  const evens = range(4, maxEven + 1, 2);
  const witnesses = new Map(evens.map((even) => [even, witness(even)]));
  const sampleWitnesses = new Map(data.sampleEvens.map((even) => [even, witnesses.get(even)]));
  fail('Goldbach derivation failed', {
    '499 even integers are covered': evens.length === 499,
    'every witness sums to its even number': [...witnesses.entries()].every(([e, [a, b]]) => a + b === e),
    'every witness consists of primes': [...witnesses.values()].every(([a, b]) => isPrime(a) && isPrime(b)),
    'sample witness for 1000 is stable': JSON.stringify(sampleWitnesses.get(1000)) === JSON.stringify([3, 997]),
  });
  return { evens, witnesses };
}

function main() {
  const data = loadInput(NAME);
  const { evens, witnesses } = trustedDerivation(data);
  const samples = data.sampleEvens.map((e) => `${e}=${witnesses.get(e)[0]}+${witnesses.get(e)[1]}`).join('; ');

  emit('# Goldbach 1000');
  emit();
  emit('## Insight');
  emit(`All ${evens.length} even integers from 4 through ${data.maxEven} have a Goldbach witness.`);
  emit(`sample witnesses : ${samples}`);
  emit();
  emit('## Explanation');
  emit('The checker caches primes up to the configured bound and then searches each even number E for a prime P not greater than E/2 where E-P is also prime.');
  emit('No counterexample is found in the bounded range, so the bounded Goldbach condition succeeds for this dataset.');
}

if (require.main === module) main();
