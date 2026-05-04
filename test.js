#!/usr/bin/env node
const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

const GREEN = '\x1b[32m';
const RED = '\x1b[31m';
const GRAY = '\x1b[90m';
const RESET = '\x1b[0m';

const ROOT = __dirname;
const EXAMPLES = path.join(ROOT, 'examples');
const OUTPUT = path.join(EXAMPLES, 'output');

let fail = false;
let count = 0;
const totalStart = process.hrtime.bigint();

function secondsSince(start) {
  return Number(process.hrtime.bigint() - start) / 1e9;
}

function timing(seconds) {
  return `${GRAY}${seconds.toFixed(3)}s${RESET}`;
}

function unifiedDiff(expected, actual, expectedName, actualName) {
  // Small, dependency-free fallback: print full expected/actual when snapshots differ.
  return [
    `--- ${expectedName}`,
    `+++ ${actualName}`,
    '@@ snapshot differs @@',
    '--- expected ---',
    expected,
    '--- actual ---',
    actual,
  ].join('\n');
}

const files = fs.readdirSync(EXAMPLES)
  .filter((file) => file.endsWith('.js') && !file.startsWith('_'))
  .sort();

for (const file of files) {
  count += 1;
  const name = path.basename(file, '.js');
  const expectedFile = path.join(OUTPUT, `${name}.md`);
  const start = process.hrtime.bigint();
  const result = spawnSync(process.execPath, [path.join(EXAMPLES, file)], {
    cwd: ROOT,
    encoding: 'utf8',
    maxBuffer: 64 * 1024 * 1024,
  });
  const elapsed = secondsSince(start);

  if (result.status !== 0) {
    const message = (result.stderr || result.stdout || '').trim();
    console.log(`${RED}FAIL${RESET} ${name} ${timing(elapsed)} execution failed: ${message}`);
    fail = true;
    continue;
  }

  const actual = result.stdout;
  const expected = fs.readFileSync(expectedFile, 'utf8');

  if (actual !== expected) {
    console.log(`${RED}FAIL${RESET} ${name} ${timing(elapsed)} output differs`);
    console.log(unifiedDiff(expected, actual, expectedFile, `${name} actual`));
    fail = true;
  } else {
    console.log(`${GREEN}OK${RESET} ${name} ${timing(elapsed)}`);
  }
}

const totalElapsed = secondsSince(totalStart);
const status = fail ? `${RED}FAIL${RESET}` : `${GREEN}OK${RESET}`;
console.log(`${status} ${count} examples run ${timing(totalElapsed)} total`);
process.exit(fail ? 1 : 0);
