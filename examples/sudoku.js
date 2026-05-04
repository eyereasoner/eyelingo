#!/usr/bin/env node
const { emit, fail, loadInput, range } = require('./_see');

const NAME = 'sudoku';
const DIGITS = range(1, 10);
const DIGIT_SET = new Set(DIGITS);

function setEquals(a, b) {
  return a.size === b.size && [...a].every((x) => b.has(x));
}

function parseGrid(puzzle) {
  const chars = [...puzzle];
  fail('Sudoku input failed', {
    'puzzle has 81 cells': chars.length === 81,
    'puzzle contains only digits, dots, or zeros': chars.every((ch) => /[1-9.0]/.test(ch)),
  });
  const cells = chars.map((ch) => '.0'.includes(ch) ? 0 : Number.parseInt(ch, 10));
  const rows = [];
  for (let i = 0; i < 81; i += 9) rows.push(cells.slice(i, i + 9));
  return rows;
}

function groupHasNoDuplicates(values) {
  const nonZero = values.filter(Boolean);
  return nonZero.length === new Set(nonZero).size;
}

function validPuzzle(grid) {
  const rows = grid.every(groupHasNoDuplicates);
  const cols = range(9).every((c) => groupHasNoDuplicates(range(9).map((r) => grid[r][c])));
  const boxes = [0, 3, 6].every((br) => [0, 3, 6].every((bc) => {
    const vals = [];
    for (let r = br; r < br + 3; r += 1) for (let c = bc; c < bc + 3; c += 1) vals.push(grid[r][c]);
    return groupHasNoDuplicates(vals);
  }));
  return rows && cols && boxes;
}

function validSolution(grid) {
  const rows = grid.every((row) => setEquals(new Set(row), DIGIT_SET));
  const cols = range(9).every((c) => setEquals(new Set(range(9).map((r) => grid[r][c])), DIGIT_SET));
  const boxes = [0, 3, 6].every((br) => [0, 3, 6].every((bc) => {
    const vals = [];
    for (let r = br; r < br + 3; r += 1) for (let c = bc; c < bc + 3; c += 1) vals.push(grid[r][c]);
    return setEquals(new Set(vals), DIGIT_SET);
  }));
  return rows && cols && boxes;
}

function candidates(grid, row, col) {
  const used = new Set();
  for (let i = 0; i < 9; i += 1) {
    used.add(grid[row][i]);
    used.add(grid[i][col]);
  }
  const br = Math.floor(row / 3) * 3;
  const bc = Math.floor(col / 3) * 3;
  for (let r = br; r < br + 3; r += 1) for (let c = bc; c < bc + 3; c += 1) used.add(grid[r][c]);
  return DIGITS.filter((digit) => !used.has(digit));
}

function chooseEmptyCell(grid) {
  let best = null;
  for (let r = 0; r < 9; r += 1) {
    for (let c = 0; c < 9; c += 1) {
      if (grid[r][c] !== 0) continue;
      const choices = candidates(grid, r, c);
      if (choices.length === 0) return { row: r, col: c, choices };
      if (!best || choices.length < best.choices.length) best = { row: r, col: c, choices };
    }
  }
  return best;
}

function solve(grid, limit = 2) {
  const work = grid.map((row) => row.slice());
  const solutions = [];
  let searchNodes = 0;

  function search() {
    if (solutions.length >= limit) return;
    searchNodes += 1;
    const cell = chooseEmptyCell(work);
    if (!cell) {
      solutions.push(work.map((row) => row.slice()));
      return;
    }
    if (cell.choices.length === 0) return;
    for (const digit of cell.choices) {
      work[cell.row][cell.col] = digit;
      search();
      work[cell.row][cell.col] = 0;
      if (solutions.length >= limit) return;
    }
  }

  search();
  return { solutions, searchNodes };
}

function fmt(row) {
  return `${row.slice(0, 3).join(' ')} | ${row.slice(3, 6).join(' ')} | ${row.slice(6).join(' ')}`;
}

function fmtPuzzle(row) {
  const vals = row.map((v) => v ? String(v) : '.');
  return `${vals.slice(0, 3).join(' ')} | ${vals.slice(3, 6).join(' ')} | ${vals.slice(6).join(' ')}`;
}

function trustedDerivation(data) {
  const puzzle = parseGrid(data.Puzzle);
  const clues = puzzle.flat().filter(Boolean).length;
  const { solutions, searchNodes } = solve(puzzle, 2);
  const solution = solutions[0];

  fail('Sudoku derivation failed', {
    'puzzle clues do not conflict': validPuzzle(puzzle),
    'puzzle has given clues': clues > 0,
    'solver found a solution': solutions.length >= 1,
    'solution is unique within search limit': solutions.length === 1,
    'computed grid is legal': validSolution(solution),
    'all original clues are preserved': range(9).every((r) => range(9).every((c) => !puzzle[r][c] || puzzle[r][c] === solution[r][c])),
  });
  return { puzzle, solution, searchNodes };
}

function main() {
  const data = loadInput(NAME);
  const { puzzle, solution, searchNodes } = trustedDerivation(data);
  const clues = puzzle.flat().filter(Boolean).length;
  const empties = 81 - clues;

  emit('# Sudoku');
  emit();
  emit('## Insight');
  emit('The puzzle is solved by search, and the completed grid is a valid unique Sudoku solution.');
  emit('case : sudoku');
  emit(`default puzzle : ${data.Name}`);
  emit();
  emit('Puzzle');
  puzzle.forEach((row, i) => {
    emit(fmtPuzzle(row));
    if ([2, 5].includes(i)) emit();
  });
  emit();
  emit('Completed grid');
  solution.forEach((row, i) => {
    emit(fmt(row));
    if ([2, 5].includes(i)) emit();
  });
  emit();
  emit('## Explanation');
  emit(`The input contains ${clues} given clues and ${empties} empty cells.`);
  emit(`The solver fills empty cells with backtracking search, choosing the open cell with the fewest legal candidates first. The search visited ${searchNodes} node(s).`);
  emit('The trust gate checks that the original clues do not conflict, exactly one solution is found, every clue is preserved, and every row, column, and 3×3 box contains digits 1 through 9.');
  emit('Only after those constraints hold does the example emit the completed grid.');
}

if (require.main === module) main();
