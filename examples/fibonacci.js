#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'fibonacci';

function fib(n) {
  let a = 0n;
  let b = 1n;
  for (let i = 0; i < n; i += 1) {
    [a, b] = [b, a + b];
  }
  return a;
}

function trustedDerivation(data) {
  const target = Number.parseInt(data.Target, 10);
  const sampleIndices = data.SampleChecks.map((n) => Number.parseInt(n, 10));
  const value = fib(target);
  const obligations = {
    'target is a non-negative integer': Number.isInteger(target) && target >= 0,
    'sample indices are valid': sampleIndices.every((n) => Number.isInteger(n) && n >= 0 && n <= target),
    'base case F(0) is zero': fib(0) === 0n,
    'base case F(1) is one': fib(1) === 1n,
    'sample recurrence checks hold': sampleIndices.every((n) => n < 2 || fib(n) === fib(n - 1) + fib(n - 2)),
    'target value is an integer string': /^\d+$/.test(value.toString()),
  };
  fail('Fibonacci derivation failed', obligations);
  return { target, value };
}

function main() {
  const data = loadInput(NAME);
  const { target, value } = trustedDerivation(data);

  emit('# Fibonacci Example (Big)');
  emit();
  emit('## Insight');
  emit(`The Fibonacci number for index ${target} is:`);
  emit(value.toString());
  emit();
  emit('## Explanation');
  emit('The Fibonacci sequence is defined by F(0)=0, F(1)=1,');
  emit('and F(n)=F(n-1)+F(n-2) for n>=2.');
  emit('Arbitrary‑precision arithmetic (BigInt) is used to');
  emit('compute the exact value without overflow, even for');
  emit(`indices as large as ${target}.`);
}

if (require.main === module) main();
