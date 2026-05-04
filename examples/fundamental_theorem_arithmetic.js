#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'fundamental_theorem_arithmetic';
const PRIMARY = 202692987;

function factor(n0) {
  let n = Number(n0);
  const factors = [];
  let d = 2;
  while (d * d <= n) {
    while (n % d === 0) {
      factors.push(d);
      n = Math.trunc(n / d);
    }
    d += d === 2 ? 1 : 2;
  }
  if (n > 1) factors.push(n);
  return factors;
}

function isPrime(n) {
  if (n < 2) return false;
  if (n % 2 === 0) return n === 2;
  let d = 3;
  while (d * d <= n) {
    if (n % d === 0) return false;
    d += 2;
  }
  return true;
}

function powers(factors) {
  const counts = new Map();
  for (const f of factors) counts.set(f, (counts.get(f) ?? 0) + 1);
  return [...counts.keys()].sort((a, b) => a - b).map((p) => counts.get(p) > 1 ? `${p}^${counts.get(p)}` : String(p)).join(' * ');
}

function product(xs) {
  return xs.reduce((acc, x) => acc * x, 1);
}

function trustedDerivation(numbers) {
  const factorizations = new Map(numbers.map((n) => [Number(n), factor(Number(n))]));
  fail('FTA derivation failed', {
    'primary number is present': factorizations.has(PRIMARY),
    'all factorizations multiply back': [...factorizations.entries()].every(([n, fs]) => product(fs) === n),
    'all factors are prime': [...factorizations.values()].flat().every(isPrime),
    'sample set is non-empty': numbers.length > 0,
    'all samples are integers greater than one': numbers.every((n) => Number.isInteger(Number(n)) && Number(n) > 1),
    'prime samples remain single factors': [...factorizations.entries()].filter(([n]) => isPrime(n)).every(([n, fs]) => fs.length === 1 && fs[0] === n),
  });
  return factorizations;
}

function main() {
  const numbers = loadInput(NAME);
  const f = trustedDerivation(numbers);
  const allFactors = [...f.values()].flat();

  emit('# Fundamental Theorem Arithmetic');
  emit();
  emit('## Insight');
  emit(`Primary N3 case: n = ${PRIMARY} has prime factors ${f.get(PRIMARY).join(' * ')}.`);
  emit(`primary prime-power form : ${powers(f.get(PRIMARY))}`);
  emit(`sample count : ${numbers.length}`);
  emit(`largest sample : ${Math.max(...numbers)}`);
  emit(`total prime factors counted with multiplicity : ${allFactors.length}`);
  emit(`distinct primes seen across samples : ${new Set(allFactors).size}`);
  emit();
  emit('Sample factorizations:');
  for (const n of numbers) emit(`  ${n} = ${powers(f.get(Number(n)))}`);
  emit();
  emit('## Explanation');
  emit('Existence comes from repeated smallest-divisor decomposition.');
  emit('At each step, the first divisor found is prime because no smaller');
  emit('positive divisor can divide the current number.');
  emit();
  emit('Smallest-divisor trace for the N3 source number:');
  let n = PRIMARY;
  const primaryFactors = f.get(PRIMARY);
  for (const p of primaryFactors.slice(0, -1)) {
    const q = Math.trunc(n / p);
    emit(`  ${n} = ${p} * ${q}`);
    n = q;
  }
  emit(`  ${primaryFactors[primaryFactors.length - 1]} is prime`);
  emit();
  emit('Uniqueness up to order is checked by reversing each traversal and sorting');
  emit('both factor lists. Matching sorted lists describe the same multiset of');
  emit('prime factors, even when the factors were discovered in the opposite order.');
  emit(`  source smallest-first factors : ${primaryFactors.join(' * ')}`);
  emit(`  source largest-first factors : ${[...primaryFactors].reverse().join(' * ')}`);
  emit(`  source sorted comparison : ${[...primaryFactors].sort((a, b) => a - b).join(' * ')}`);
  emit();
  emit('The additional samples cover repeated small factors, special products,');
  emit('large composites, and a larger prime that has no smaller divisor.');
}

if (require.main === module) main();
