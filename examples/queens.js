#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'queens';

function solve(n) {
  const solutions = [];
  let count = 0;
  function search(row, cols, diag1, diag2, positions) {
    if (row === n) {
      count += 1;
      solutions.push([...positions]);
      return;
    }
    let available = ((1 << n) - 1) & ~(cols | diag1 | diag2);
    while (available) {
      const bit = available & -available;
      available -= bit;
      const col = Math.floor(Math.log2(bit));
      search(row + 1, cols | bit, (diag1 | bit) << 1, (diag2 | bit) >> 1, [...positions, col]);
    }
  }
  search(0, 0, 0, 0, []);
  return { solutions, count };
}

function trustedDerivation(data) {
  const n = Number.parseInt(data.N, 10);
  const { solutions, count } = solve(n);
  const first = solutions[0];
  fail('8-Queens derivation failed', {
    'first solution is stable': JSON.stringify(first.map((c) => c + 1)) === JSON.stringify([1, 5, 8, 6, 3, 7, 2, 4]),
    'total count is complete': count === 92,
    'one queen per row': first.length === n,
    'columns are unique': new Set(first).size === n,
  });
  return { first, total: count };
}

function rowText(n, col) {
  return [...Array(n).keys()].map((i) => i === col ? 'Q' : '.').join(' ');
}

function main() {
  const data = loadInput(NAME);
  const n = Number.parseInt(data.N, 10);
  const maxPrint = Number.parseInt(data.MaxPrint, 10);
  const { first, total } = trustedDerivation(data);

  emit('# 8-Queens');
  emit();
  emit('## Insight');
  emit(`Solving ${n}-Queens...`);
  emit(`Printing at most ${maxPrint} solution(s).`);
  emit();
  emit('Solution 1:');
  for (const col of first) emit(rowText(n, col));
  emit(`As column positions by row: [${first.map((c) => c + 1).join(', ')}]`);
  emit();
  emit(`Total solutions for ${n}-Queens: ${total}`);
  emit();
  emit('## Explanation');
  emit(`The solver places one queen per row on a ${n}x${n} board.`);
  emit('At each row it uses bit masks for occupied columns and both diagonal directions to enumerate only safe candidate columns.');
  emit('Counting continues after the printed solution limit, so the total solution count remains complete.');
}

if (require.main === module) main();
