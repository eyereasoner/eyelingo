#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'dijkstra_risk_path';

function edgeScore(edge, riskWeight) {
  return Number(edge.cost) + riskWeight * Number(edge.risk);
}

function buildGraph(data) {
  const graph = new Map();
  for (const edge of data.edges) {
    if (!graph.has(edge.from)) graph.set(edge.from, []);
    graph.get(edge.from).push(edge);
  }
  return graph;
}

function shortestPath(data) {
  const graph = buildGraph(data);
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

function enumerateSimplePaths(data) {
  const graph = buildGraph(data);
  const riskWeight = Number(data.riskWeight);
  const paths = [];

  function walk(node, path, rawCost, riskSum, score) {
    if (node === data.goal) {
      paths.push({ path: [...path], rawCost, riskSum, score });
      return;
    }
    for (const edge of graph.get(node) ?? []) {
      if (path.includes(edge.to)) continue;
      walk(
        edge.to,
        [...path, edge.to],
        rawCost + Number(edge.cost),
        riskSum + Number(edge.risk),
        score + edgeScore(edge, riskWeight),
      );
    }
  }

  walk(data.start, [data.start], 0.0, 0.0, 0.0);
  return paths;
}

function samePath(a, b) {
  return JSON.stringify(a) === JSON.stringify(b);
}

function trustedDerivation(data) {
  const result = shortestPath(data);
  const allPaths = enumerateSimplePaths(data);
  const bestScore = Math.min(...allPaths.map((path) => path.score));
  fail('Dijkstra risk path derivation failed', {
    'edge costs are non-negative': data.edges.every((edge) => Number(edge.cost) >= 0),
    'edge risks are non-negative': data.edges.every((edge) => Number(edge.risk) >= 0),
    'at least one path reaches the goal': allPaths.length > 0,
    'selected path starts at source': result.path[0] === data.start,
    'selected path ends at goal': result.path[result.path.length - 1] === data.goal,
    'selected path appears in enumerated paths': allPaths.some((path) => samePath(path.path, result.path)),
    'selected score is minimal among simple paths': Math.abs(result.score - bestScore) < 1e-9,
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
  emit('Each edge contributes its delivery cost plus the configured risk penalty.');
  emit("Dijkstra's queue expands the lowest accumulated score first, so the first time HubZ is popped the selected route is optimal for the weighted graph.");
  emit('The trust gate also enumerates the simple paths in this small graph and checks that no other path has a lower risk-adjusted score.');
  emit('The selected route balances cost and risk through DepotB and LabD.');
}

if (require.main === module) main();
