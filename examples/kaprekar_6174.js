#!/usr/bin/env node
const { emit, fail, loadInput, range } = require('./_see');

const NAME = 'kaprekar_6174';

function pad4(n) { return String(n).padStart(4, '0'); }

function kaprekarStep(n) {
  const digits = [...pad4(n)].sort();
  const asc = Number.parseInt(digits.join(''), 10);
  const desc = Number.parseInt([...digits].reverse().join(''), 10);
  return desc - asc;
}

function chain(start, target, zero, maxSteps) {
  const values = [start];
  let current = start;
  for (let i = 0; i < maxSteps; i += 1) {
    if (current === target || current === zero) return values;
    current = kaprekarStep(current);
    values.push(current);
  }
  return values;
}

function counterInc(map, key) { map.set(key, (map.get(key) ?? 0) + 1); }

function trustedDerivation(data) {
  const target = data.TargetConstant;
  const zero = data.ZeroBasin;
  const maxSteps = data.MaxKaprekarSteps;
  const emitted = new Map();
  const omitted = new Map();
  const distribution = new Map();

  for (const start of range(data.StartCount)) {
    const ch = chain(start, target, zero, maxSteps);
    if (ch[ch.length - 1] === target) {
      emitted.set(start, ch);
      counterInc(distribution, ch.length - 1);
    } else if (ch[ch.length - 1] === zero) {
      omitted.set(start, ch);
    }
  }

  fail('Kaprekar derivation failed', {
    'all starts are classified': emitted.size + omitted.size === data.StartCount,
    'maximum step count stays within bound': Math.max(...distribution.keys()) <= maxSteps,
    'target is fixed when started from target': JSON.stringify(chain(target, target, zero, maxSteps)) === JSON.stringify([target]),
    'emitted chains end at target': [...emitted.values()].every((ch) => ch[ch.length - 1] === target),
    'omitted chains end at zero basin': [...omitted.values()].every((ch) => ch[ch.length - 1] === zero),
  });
  return { emitted, omitted, distribution };
}

function chainText(values) {
  const shown = values.length === 1 ? values : values.slice(1);
  return `(${shown.map(pad4).join(' ')})`;
}

function omittedText(values) {
  const shown = values.length === 1 ? values : values.slice(1);
  return `(${shown.map(pad4).join(' ')})`;
}

function main() {
  const data = loadInput(NAME);
  const { emitted, omitted, distribution } = trustedDerivation(data);

  emit('# Kaprekar 6174');
  emit();
  emit('## Insight');
  emit('Kaprekar chains that end at 6174 are emitted as :kaprekar facts.');
  emit(`total emitted : ${emitted.size}`);
  emit(`omitted 0000 basin : ${omitted.size}`);
  emit(`maximum steps to 6174 : ${Math.max(...distribution.keys())}`);
  emit();
  emit('Selected facts, shown with four-digit padding for readability:');
  for (const start of [1, 3524, 6174, 9831, 9998]) emit(`  ${pad4(start)} :kaprekar ${chainText(emitted.get(start))}`);
  emit();
  emit('## Explanation');
  emit('Each start is read as four digits, so 1 is treated as 0001.');
  emit('The digits are sorted once, then the optimized identity computes the');
  emit('same result as descending-number minus ascending-number.');
  emit('The search is bounded to seven steps, matching the N3 source: any');
  emit('four-digit start that reaches 6174 does so within that bound.');
  emit();
  emit('Step-count distribution for emitted starts:');
  for (const steps of range(0, 8)) emit(`  ${steps} step(s) : ${distribution.get(steps) ?? 0} start(s)`);
  emit();
  emit('Examples omitted because they fall to 0000:');
  for (const start of [0, 1111, 2222, 9999]) emit(`  ${pad4(start)} -> ${omittedText(omitted.get(start))}`);
}

if (require.main === module) main();
