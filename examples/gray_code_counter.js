#!/usr/bin/env node
const { emit, fail, loadInput, range } = require('./_see');

const NAME = 'gray_code_counter';

function gray(n) { return n ^ (n >> 1); }
function bits(value, width) { return value.toString(2).padStart(width, '0'); }
function hamming(a, b) {
  let v = a ^ b;
  let count = 0;
  while (v) { count += v & 1; v >>= 1; }
  return count;
}

function trustedDerivation(data) {
  const width = Number.parseInt(data.bits, 10);
  const steps = Number.parseInt(data.steps, 10);
  const sequence = range(steps).map(gray);
  const distances = range(steps).map((i) => hamming(sequence[i], sequence[(i + 1) % steps]));
  fail('Gray-code derivation failed', {
    'sequence length matches configured steps': sequence.length === steps,
    'sequence covers every configured state once': new Set(sequence).size === steps && steps === 2 ** width,
    'wrap distance is one': distances[distances.length - 1] === 1,
    'maximum adjacent Hamming distance is one': Math.max(...distances) === 1,
  });
  return sequence;
}

function main() {
  const data = loadInput(NAME);
  const sequence = trustedDerivation(data);
  const width = data.bits;
  const prefix = sequence.slice(0, 8).map((x) => bits(x, width)).join(', ');
  const wrap = `${bits(sequence[sequence.length - 1], width)} -> ${bits(sequence[0], width)}`;
  const maxDistance = Math.max(...range(sequence.length).map((i) => hamming(sequence[i], sequence[(i + 1) % sequence.length])));

  emit('# Gray Code Counter');
  emit();
  emit('## Insight');
  emit(`bits : ${width}`);
  emit(`states visited : ${sequence.length}`);
  emit(`unique states : ${new Set(sequence).size}`);
  emit(`sequence prefix : ${prefix}`);
  emit(`wrap transition : ${wrap}`);
  emit(`maximum adjacent Hamming distance : ${maxDistance}`);
  emit();
  emit('## Explanation');
  emit('The counter maps each integer n to n xor (n >> 1), which is the reflected binary Gray-code construction.');
  emit('For 4 bits, the first 16 integers cover the full state space without duplicates.');
  emit('The Hamming-distance check compares each state with the next state, including the final wraparound transition.');
  emit('A valid cyclic Gray counter therefore changes exactly one bit at every step.');
}

if (require.main === module) main();
