#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'dijkstra_risk_path';

function edgeScore(edge, riskWeight) {
  return Number(edge.cost) + riskWeight * Number(edge.risk);
}

function shortestPath(data) {
  const graph = new Map();
  for (const edge of data.edges) {
    if (!graph.has(edge.from)) graph.set(edge.from, []);
    graph.get(edge.from).push(edge);
  }
  const start = data.start;
  const goal = data.goal;
  const riskWeight = Number(data.riskWeight);
  const queue = [[0.0, start, [start], 0.0, 0.0]];
  const best = new Map([[start, 0.0]]);

  while (queue.length) {
    queue.sort((a, b) => a[0] - b[0] || String(a[1]).localeCompare(String(b[1])));
    const [score, node, path, rawCost, riskSum] = queue.shift();
    if (node === goal) return { path, rawCost, riskSum, score };
    if (score > (best.get(node) ?? Infinity)) continue;
    for (const edge of graph.get(node) ?? []) {
      const nxt = edge.to;
      const nextCost = rawCost + Number(edge.cost);
      const nextRisk = riskSum + Number(edge.risk);
      const nextScore = score + edgeScore(edge, riskWeight);
      if (nextScore < (best.get(nxt) ?? Infinity)) {
        best.set(nxt, nextScore);
        queue.push([nextScore, nxt, [...path, nxt], nextCost, nextRisk]);
      }
    }
  }
  throw new Error('no path');
}

function trustedDerivation(data) {
  const result = shortestPath(data);
  fail('Dijkstra risk path derivation failed', {
    'expected path is selected': JSON.stringify(result.path) === JSON.stringify(data.expected.path),
    'expected score is selected': Math.abs(result.score - Number(data.expected.score)) < 1e-9,
    'raw cost is stable': Number(result.rawCost.toFixed(2)) === 10.00,
    'risk sum is stable': Number(result.riskSum.toFixed(2)) === 0.55,
  });
  return result;
}

function main() {
  const data = loadInput(NAME);
  const { path, rawCost, riskSum, score } = trustedDerivation(data);

  emit('# Dijkstra Risk Path');
  emit();
  emit('## Insight');
  emit(`selected path : ${path.join(' -> ')}`);
  emit(`raw cost : ${rawCost.toFixed(2)}`);
  emit(`risk sum : ${riskSum.toFixed(2)}`);
  emit(`risk-adjusted score : ${score.toFixed(2)}`);
  emit(`edges in selected path : ${path.length - 1}`);
  emit();
  emit('## Explanation');
  emit("Each edge contributes its delivery cost plus the configured risk penalty.");
  emit("Dijkstra's queue expands the lowest accumulated score first, so the first time HubZ is popped the selected route is optimal for the weighted graph.");
  emit('The DepotC shortcut has lower early cost but carries enough risk to lose under the configured risk weight.');
  emit('The selected route balances cost and risk through DepotB and LabD.');
}

if (require.main === module) main();
