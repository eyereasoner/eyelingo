#!/usr/bin/env node
const { emit, fail, loadInput, minBy, sum } = require('./_see');

const NAME = 'genetic_knapsack_selection';

function evaluate(genome, data) {
  const pairs = [...genome].map((bit, i) => [bit, data.Items[i]]);
  const weight = sum(pairs.filter(([bit]) => bit === '1').map(([, item]) => item.Weight));
  const value = sum(pairs.filter(([bit]) => bit === '1').map(([, item]) => item.Value));
  const fitness = weight <= data.Capacity ? 1_000_000 - value : 2_000_000 + (weight - data.Capacity);
  return [weight, value, fitness];
}

function flip(genome, i) {
  return genome.slice(0, i) + (genome[i] === '1' ? '0' : '1') + genome.slice(i + 1);
}

function search(data) {
  let genome = data.StartGenome;
  const history = [genome];
  for (let generation = 0; generation < data.MaxGenerations; generation += 1) {
    const candidates = [genome, ...[...Array(genome.length).keys()].map((i) => flip(genome, i))];
    const best = minBy(candidates, (g) => [evaluate(g, data)[2], g]);
    if (best === genome) return { genome, generations: history.length, history };
    genome = best;
    history.push(genome);
  }
  return { genome, generations: history.length, history };
}

function trustedDerivation(data) {
  const { genome, generations, history } = search(data);
  const [weight, value, fitness] = evaluate(genome, data);
  const expected = data.Expected;
  const neighbors = [...Array(genome.length).keys()].map((i) => flip(genome, i));
  fail('Genetic knapsack derivation failed', {
    'genome has one bit per item': genome.length === data.Items.length,
    'final genome matches fixture': genome === expected.FinalGenome,
    'weight matches fixture': weight === expected.FinalWeight,
    'value matches fixture': value === expected.FinalValue,
    'final candidate is feasible': weight <= data.Capacity,
    'no one-bit neighbor improves fitness': neighbors.every((n) => evaluate(genome, data)[2] <= evaluate(n, data)[2]),
  });
  return { genome, weight, value, fitness, generations, history };
}

function main() {
  const data = loadInput(NAME);
  const { genome, weight, value, fitness, generations } = trustedDerivation(data);
  const selected = [...genome].map((bit, i) => bit === '1' ? data.Items[i].Name : null).filter(Boolean);

  emit('# Genetic Knapsack Selection');
  emit();
  emit('## Insight');
  emit(`final genome : ${genome}`);
  emit(`selected items : ${selected.join(', ')}`);
  emit(`weight : ${weight} / ${data.Capacity}`);
  emit(`value : ${value}`);
  emit(`fitness : ${fitness}`);
  emit(`generations evaluated : ${generations}`);
  emit();
  emit('## Explanation');
  emit('Each genome bit says whether the corresponding item is selected for the knapsack.');
  emit('Feasible candidates get fitness 1000000 minus value, so higher value means lower fitness; overweight candidates are penalized above every feasible candidate.');
  emit('At every generation, all single-bit mutants of the parent are compared with the parent, and the lowest-fitness candidate is selected deterministically.');
  emit(`The run stops at ${genome} because every one-bit neighbor is no better under the capacity ${data.Capacity} rule.`);
}

if (require.main === module) main();
