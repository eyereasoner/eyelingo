#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'fibonacci';
const TARGET = 10000;

function fib(n) {
  let a = 0n;
  let b = 1n;
  for (let i = 0; i < n; i += 1) {
    [a, b] = [b, a + b];
  }
  return a;
}

function trustedDerivation(data) {
  const value = fib(TARGET);
  const sampleIndices = [0, 1, 10, 100, 1000, 10000];
  const obligations = {};
  for (const n of sampleIndices) obligations[`F(${n}) matches input reference`] = fib(n).toString() === data[String(n)];
  obligations['target is configured'] = String(TARGET) in data;
  obligations['target value has expected digit count'] = value.toString().length === data[String(TARGET)].length;
  obligations['target value matches reference'] = value.toString() === data[String(TARGET)];
  fail('Fibonacci derivation failed', obligations);
  return value;
}

function main() {
  const data = loadInput(NAME);
  const value = trustedDerivation(data);

  emit('# Fibonacci Example (Big)');
  emit();
  emit('## Insight');
  emit(`The Fibonacci number for index ${TARGET} is:`);
  emit(value.toString());
  emit();
  emit('## Explanation');
  emit('The Fibonacci sequence is defined by F(0)=0, F(1)=1,');
  emit('and F(n)=F(n-1)+F(n-2) for n>=2.');
  emit('Arbitrary‑precision arithmetic (math/big) is used to');
  emit('compute the exact value without overflow, even for');
  emit('indices as large as 10000.');
}

if (require.main === module) main();
