#!/usr/bin/env node
const { emit, fail, loadInput, compareKeys } = require('./_see');

const NAME = 'gps';

class RouteScore {
  constructor(path, duration, cost, belief, comfort) {
    this.path = path;
    this.duration = duration;
    this.cost = cost;
    this.belief = belief;
    this.comfort = comfort;
  }
}

function enumerateRoutes(data, start, goal) {
  const graph = new Map();
  for (const edge of data.Edges) {
    if (!graph.has(edge.From)) graph.set(edge.From, []);
    graph.get(edge.From).push(edge);
  }
  const routes = [];
  function walk(node, path, duration, cost, belief, comfort) {
    if (node === goal) {
      routes.push(new RouteScore([...path], duration, cost, belief, comfort));
      return;
    }
    for (const edge of graph.get(node) ?? []) {
      const nxt = edge.To;
      if (path.includes(nxt)) continue;
      walk(nxt, [...path, nxt], duration + Number(edge.Duration), cost + Number(edge.Cost), belief * Number(edge.Belief), comfort * Number(edge.Comfort));
    }
  }
  walk(start, [start], 0.0, 0.0, 1.0, 1.0);
  return routes.sort((a, b) => compareKeys([a.duration, a.cost, a.path.join('\u0000')], [b.duration, b.cost, b.path.join('\u0000')]));
}

function routeLabel(data, route) {
  const pathLabel = route.path.join(' → ');
  const match = Object.values(data.Routes).find((candidate) => candidate.Label === pathLabel);
  return match ? match.Label : pathLabel;
}

function trustedDerivation(data) {
  const start = data.Traveller.Location;
  const goal = data.Goal;
  const routes = enumerateRoutes(data, start, goal);
  const best = routes[0];
  const alternative = routes[1];
  fail('GPS derivation failed', {
    'goal is configured': typeof goal === 'string' && goal.length > 0,
    'at least one simple route connects start to goal': routes.length > 0,
    'every route starts at traveller location': routes.every((route) => route.path[0] === start),
    'every route ends at goal': routes.every((route) => route.path[route.path.length - 1] === goal),
    'routes are simple': routes.every((route) => route.path.length === new Set(route.path).size),
    'edge metrics are non-negative probabilities or costs': data.Edges.every((edge) => edge.Duration >= 0 && edge.Cost >= 0 && edge.Belief >= 0 && edge.Belief <= 1 && edge.Comfort >= 0 && edge.Comfort <= 1),
    'best route is fastest after sorting': routes.every((route) => best.duration <= route.duration),
  });
  return { best, alternative, routes };
}

function main() {
  const data = loadInput(NAME);
  const { best, alternative, routes } = trustedDerivation(data);
  const bestLabel = routeLabel(data, best);

  emit('# GPS — Goal driven route planning');
  emit();
  emit('## Insight');
  emit('Take the fastest route found.');
  emit(`Recommended route: ${bestLabel}`);
  emit();
  emit('## Explanation');
  emit(`From ${data.Traveller.Location} to ${data.Goal}, the planner found ${routes.length} route(s) in this small map.`);
  if (alternative) {
    const alternativeLabel = routeLabel(data, alternative);
    emit(`The recommended route (${bestLabel}) takes ${best.duration.toFixed(1)} seconds at cost ${best.cost.toFixed(2)}, with belief ${best.belief.toFixed(4)} and comfort ${best.comfort.toFixed(2)}. The alternative (${alternativeLabel}) takes ${alternative.duration.toFixed(1)} seconds at cost ${alternative.cost.toFixed(3)}, with belief ${alternative.belief.toFixed(6)} and comfort ${alternative.comfort.toFixed(4)}.`);
    emit('So the recommended route is faster for this input, and the trust gate checks that no enumerated route is faster.');
  } else {
    emit(`The recommended route takes ${best.duration.toFixed(1)} seconds at cost ${best.cost.toFixed(2)}, with belief ${best.belief.toFixed(4)} and comfort ${best.comfort.toFixed(2)}.`);
  }
}

if (require.main === module) main();
