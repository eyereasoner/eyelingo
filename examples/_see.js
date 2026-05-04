const fs = require('fs');
const path = require('path');

const EXAMPLES_DIR = __dirname;
const INPUT_DIR = path.join(EXAMPLES_DIR, 'input');

function loadInput(name) {
  return JSON.parse(fs.readFileSync(path.join(INPUT_DIR, `${name}.json`), 'utf8'));
}

function emit(line = '') {
  if (line) process.stdout.write(`${line}  \n`);
  else process.stdout.write('\n');
}

function emitLines(lines) {
  for (const line of lines) emit(line);
}

function fail(prefix, obligations) {
  const failed = Object.entries(obligations)
    .filter(([, ok]) => !ok)
    .map(([name]) => name);
  if (failed.length) {
    throw new Error(`${prefix}: ${failed.join(', ')}`);
  }
}

function sum(values) {
  let out = 0;
  for (const value of values) out += value;
  return out;
}

function minBy(values, keyFn) {
  if (!values.length) throw new Error('minBy on empty array');
  let best = values[0];
  let bestKey = keyFn(best);
  for (let i = 1; i < values.length; i += 1) {
    const key = keyFn(values[i]);
    if (compareKeys(key, bestKey) < 0) {
      best = values[i];
      bestKey = key;
    }
  }
  return best;
}

function maxBy(values, keyFn) {
  if (!values.length) throw new Error('maxBy on empty array');
  let best = values[0];
  let bestKey = keyFn(best);
  for (let i = 1; i < values.length; i += 1) {
    const key = keyFn(values[i]);
    if (compareKeys(key, bestKey) > 0) {
      best = values[i];
      bestKey = key;
    }
  }
  return best;
}

function compareKeys(a, b) {
  const aa = Array.isArray(a) ? a : [a];
  const bb = Array.isArray(b) ? b : [b];
  const n = Math.min(aa.length, bb.length);
  for (let i = 0; i < n; i += 1) {
    if (aa[i] < bb[i]) return -1;
    if (aa[i] > bb[i]) return 1;
  }
  return aa.length - bb.length;
}

function range(start, stop = undefined, step = 1) {
  if (stop === undefined) {
    stop = start;
    start = 0;
  }
  const out = [];
  for (let i = start; step > 0 ? i < stop : i > stop; i += step) out.push(i);
  return out;
}

function roundTo(value, digits = 0) {
  const factor = 10 ** digits;
  return Math.round((value + Number.EPSILON) * factor) / factor;
}

function boolText(value) {
  return value ? 'true' : 'false';
}

module.exports = {
  loadInput,
  emit,
  emitLines,
  fail,
  sum,
  minBy,
  maxBy,
  compareKeys,
  range,
  roundTo,
  boolText,
};
