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

function trustedDerivation(data) {
  const start = data.Traveller.Location;
  const goal = 'Oostende';
  const routes = enumerateRoutes(data, start, goal);
  const direct = routes.find((route) => JSON.stringify(route.path) === JSON.stringify(['Gent', 'Brugge', 'Oostende']));
  const alternative = routes.find((route) => JSON.stringify(route.path) === JSON.stringify(['Gent', 'Kortrijk', 'Brugge', 'Oostende']));
  fail('GPS derivation failed', {
    'exactly two simple routes connect Gent to Oostende': routes.length === 2,
    'direct route duration is additive': direct.duration === 2400.0,
    'alternative route duration is additive': alternative.duration === 4100.0,
    'direct route is cheaper': direct.cost < alternative.cost,
    'direct route is more reliable': direct.belief > alternative.belief,
    'direct route is more comfortable': direct.comfort > alternative.comfort,
  });
  return { direct, alternative };
}

function main() {
  const data = loadInput(NAME);
  const { direct, alternative } = trustedDerivation(data);
  const directLabel = data.Routes.routeDirect.Label;
  const altLabel = data.Routes.routeViaKortrijk.Label;

  emit('# GPS — Goal driven route planning');
  emit();
  emit('## Insight');
  emit('Take the direct route via Brugge.');
  emit(`Recommended route: ${directLabel}`);
  emit();
  emit('## Explanation');
  emit('From Gent to Oostende, the planner found two routes in this small map.');
  emit(`The direct route (${directLabel}) takes ${direct.duration.toFixed(1)} seconds at cost ${direct.cost.toFixed(2)}, with belief ${direct.belief.toFixed(4)} and comfort ${direct.comfort.toFixed(2)}. The alternative (${altLabel}) takes ${alternative.duration.toFixed(1)} seconds at cost ${alternative.cost.toFixed(3)}, with belief ${alternative.belief.toFixed(6)} and comfort ${alternative.comfort.toFixed(4)}.`);
  emit('So the direct route is faster, cheaper, more reliable, and slightly more comfortable.');
}

if (require.main === module) main();
