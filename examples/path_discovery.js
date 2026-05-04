#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'path_discovery';
const SOURCE_GRAPH_AIRPORT_LABELS = 7698;
const SOURCE_GRAPH_OUTBOUND_FACTS = 37505;
const DEFAULT_MAX_STOPOVERS = 2;

function label(data, term) { return data.Labels[term] ?? term; }

function buildAdjacency(data) {
  const adj = new Map();
  for (const edge of data.Edges) {
    if (!adj.has(edge.From)) adj.set(edge.From, []);
    adj.get(edge.From).push(edge.To);
  }
  for (const [node, outbound] of adj) {
    outbound.sort((a, b) => label(data, a).localeCompare(label(data, b)) || a.localeCompare(b));
    adj.set(node, outbound);
  }
  return adj;
}

function routeLabel(data, route) {
  return route.map((node) => label(data, node)).join(' -> ');
}

function dfs(data, adj, sourceID, destinationID, maxStopovers) {
  const maxHops = maxStopovers + 1;
  const routes = [];
  const stats = { recursive_calls: 0, edge_tests: 0, depth_limit_leaves: 0, max_depth: 0 };

  function walk(current, path) {
    const depth = path.length - 1;
    stats.recursive_calls += 1;
    stats.max_depth = Math.max(stats.max_depth, depth);
    if (path.length > 1 && current === destinationID) {
      routes.push([...path]);
      return;
    }
    if (depth >= maxHops) {
      stats.depth_limit_leaves += 1;
      return;
    }
    const outbound = adj.get(current) ?? [];
    stats.edge_tests += outbound.length;
    for (const nxt of outbound) {
      if (path.includes(nxt)) continue;
      walk(nxt, [...path, nxt]);
    }
  }

  walk(sourceID, [sourceID]);
  routes.sort((a, b) => routeLabel(data, a).localeCompare(routeLabel(data, b)));
  return { routes, stats };
}

function stopovers(route) { return Math.max(0, route.length - 2); }

function routeDistribution(routes) {
  let direct = 0;
  let oneStop = 0;
  let twoStop = 0;
  for (const route of routes) {
    const s = stopovers(route);
    if (s === 0) direct += 1;
    else if (s === 1) oneStop += 1;
    else if (s === 2) twoStop += 1;
  }
  return [direct, oneStop, twoStop];
}

function expandedAirports(data, adj, sourceID, maxDepth = 2) {
  const seen = new Set();
  const ordered = [];
  function walk(current, path) {
    const depth = path.length - 1;
    if (depth > maxDepth) return;
    if (!seen.has(current)) {
      seen.add(current);
      ordered.push(current);
    }
    if (depth === maxDepth) return;
    for (const nxt of adj.get(current) ?? []) {
      if (!path.includes(nxt)) walk(nxt, [...path, nxt]);
    }
  }
  walk(sourceID, [sourceID]);
  return ordered;
}

function array2SetPairs(edges) {
  return new Set(edges.map((e) => `${e.From}\u0000${e.To}`));
}

function trustRouteDiscovery(data, routes, stats, sourceID, destinationID, maxStopovers) {
  const edgeSet = array2SetPairs(data.Edges);
  const dist = routeDistribution(routes);
  fail('Path discovery derivation failed', {
    'source and destination labels are known': label(data, sourceID) !== sourceID && label(data, destinationID) !== destinationID,
    'route distribution matches route count': dist.reduce((acc, x) => acc + x, 0) === routes.length,
    'all routes start at source': routes.every((route) => route[0] === sourceID),
    'all routes end at destination': routes.every((route) => route[route.length - 1] === destinationID),
    'all routes respect stopover bound': routes.every((route) => stopovers(route) <= maxStopovers),
    'all hops are graph facts': routes.every((route) => route.slice(0, -1).every((_, i) => edgeSet.has(`${route[i]}\u0000${route[i + 1]}`))),
    'routes do not revisit airports': routes.every((route) => route.length === new Set(route).size),
    'full source graph is loaded': Object.keys(data.Labels).length === SOURCE_GRAPH_AIRPORT_LABELS && data.Edges.length === SOURCE_GRAPH_OUTBOUND_FACTS,
    'search touched graph facts': stats.edge_tests > 0,
    'search depth respects bound': stats.max_depth <= maxStopovers + 1,
  });
}

function parseCli(data) {
  const [sourceArg, destinationArg, maxStopoversArg] = process.argv.slice(2);
  const sourceID = sourceArg ?? data.SourceID;
  const destinationID = destinationArg ?? data.DestinationID;
  const maxStopovers = maxStopoversArg === undefined ? DEFAULT_MAX_STOPOVERS : Number.parseInt(maxStopoversArg, 10);
  fail('Invalid path_discovery arguments', {
    'source airport exists in labels': Object.hasOwn(data.Labels, sourceID),
    'destination airport exists in labels': Object.hasOwn(data.Labels, destinationID),
    'max stopovers is a non-negative integer': Number.isInteger(maxStopovers) && maxStopovers >= 0,
  });
  return { sourceID, destinationID, maxStopovers };
}

function main() {
  const data = loadInput(NAME);
  const { sourceID, destinationID, maxStopovers } = parseCli(data);
  const adj = buildAdjacency(data);
  const { routes, stats } = dfs(data, adj, sourceID, destinationID, maxStopovers);
  trustRouteDiscovery(data, routes, stats, sourceID, destinationID, maxStopovers);
  const secondHopBase = routes[0]?.[1];
  const firstHopOut = secondHopBase ? [...(adj.get(secondHopBase) ?? [])].sort((a, b) => label(data, a).localeCompare(label(data, b)) || a.localeCompare(b)) : [];
  const [direct, oneStop, twoStop] = routeDistribution(routes);
  const endpointAirports = new Set(data.Edges.flatMap((edge) => [edge.From, edge.To])).size;
  const expanded = expandedAirports(data, adj, sourceID);

  emit('# Path Discovery');
  emit();
  emit('## Insight');
  emit(`The path discovery query finds ${routes.length} air routes with at most ${maxStopovers} stopovers.`);
  emit(`from : ${label(data, sourceID)}`);
  emit(`to : ${label(data, destinationID)}`);
  emit(`max stopovers : ${maxStopovers}`);
  emit();
  emit('Discovered routes:');
  routes.forEach((route, i) => emit(`route ${i + 1} (${stopovers(route)} stopovers): ${routeLabel(data, route)}`));
  emit();
  emit('## Explanation');
  emit('The N3 source defines a recursive :route relation over nepo:hasOutboundRouteTo facts. A route can use a direct edge when the current length is within the maximum, or extend through a non-visited intermediate airport and recurse with length+1. The final log:collectAllIn query collects the labels of each airport in every route from the source to the destination.');
  emit(`source N3 airport labels : ${SOURCE_GRAPH_AIRPORT_LABELS}`);
  emit(`source N3 outbound-route facts : ${SOURCE_GRAPH_OUTBOUND_FACTS}`);
  emit(`translated full airport labels : ${Object.keys(data.Labels).length}`);
  emit(`translated full outbound-route facts : ${data.Edges.length}`);
  emit(`airport terms appearing in outbound facts : ${endpointAirports}`);
  emit(`frontier airports expanded : ${expanded.length}`);
  emit(`bounded search outbound facts touched : ${stats.edge_tests}`);
  emit(`source outbound candidates : ${(adj.get(sourceID) ?? []).length}`);
  if (secondHopBase) emit(`${label(data, secondHopBase)} outbound candidates : ${firstHopOut.length}`);
  emit(`direct routes : ${direct}`);
  emit(`one-stop routes : ${oneStop}`);
  emit(`two-stopover routes : ${twoStop}`);
  emit(`search recursive calls : ${stats.recursive_calls}`);
  emit(`search edge tests : ${stats.edge_tests}`);
  emit(`search depth-limit leaves : ${stats.depth_limit_leaves}`);
  if (secondHopBase) {
    emit(`Second-hop candidates from ${label(data, secondHopBase)}:`);
    for (const airport of firstHopOut) emit(`${label(data, airport)} (${airport})`);
  }
}

if (require.main === module) main();
