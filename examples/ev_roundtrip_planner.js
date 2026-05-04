#!/usr/bin/env node
const { emit, fail, loadInput, compareKeys } = require('./_see');

const NAME = 'ev_roundtrip_planner';

class State {
  constructor(At, Battery, Pass) {
    this.At = At;
    this.Battery = Battery;
    this.Pass = Pass;
  }
  key() { return `${this.At}\u0000${this.Battery}\u0000${this.Pass}`; }
}

class Plan {
  constructor(actions, final, duration, cost, belief, comfort, fuelRemaining) {
    this.actions = actions;
    this.final = final;
    this.duration = duration;
    this.cost = cost;
    this.belief = belief;
    this.comfort = comfort;
    this.fuel_remaining = fuelRemaining;
  }
}

function wildcard(pattern, value) {
  return pattern === '*' || pattern === value;
}

function matchesState(pattern, state) {
  return wildcard(pattern.At, state.At) && wildcard(pattern.Battery, state.Battery) && wildcard(pattern.Pass, state.Pass);
}

function matchesGoal(state, goal) {
  return wildcard(goal.At, state.At) && wildcard(goal.Battery, state.Battery) && wildcard(goal.Pass, state.Pass);
}

function applyState(pattern, current) {
  return new State(
    pattern.At === '*' ? current.At : pattern.At,
    pattern.Battery === '*' ? current.Battery : pattern.Battery,
    pattern.Pass === '*' ? current.Pass : pattern.Pass,
  );
}

function search(data) {
  const start = new State(data.Vehicle.At, data.Vehicle.Battery, data.Vehicle.Pass);
  const goal = data.Goal;
  const thresholds = data.Thresholds;
  const fuel = Number.parseInt(data.FuelSteps, 10);
  const plans = [];

  function walk(state, path, duration, cost, belief, comfort, fuelLeft, seen) {
    if (matchesGoal(state, goal)) {
      if (belief > thresholds.MinBelief && cost < thresholds.MaxCost && duration < thresholds.MaxDuration) {
        plans.push(new Plan([...path], state, duration, cost, belief, comfort, fuelLeft));
      }
      return;
    }
    if (fuelLeft === 0) return;
    for (const action of data.Actions) {
      if (!matchesState(action.From, state)) continue;
      const nxt = applyState(action.To, state);
      if (seen.has(nxt.key()) && nxt.key() !== state.key()) continue;
      const nextSeen = new Set(seen);
      nextSeen.add(nxt.key());
      walk(
        nxt,
        [...path, action.Name],
        duration + Number(action.Duration),
        cost + Number(action.Cost),
        belief * Number(action.Belief),
        comfort * Number(action.Comfort),
        fuelLeft - 1,
        nextSeen,
      );
    }
  }

  walk(start, [], 0.0, 0.0, 1.0, 1.0, fuel, new Set([start.key()]));
  return plans.sort((a, b) => compareKeys([a.duration, a.cost, a.actions.join(' / ')], [b.duration, b.cost, b.actions.join(' / ')]));
}

function trustedDerivation(data) {
  const plans = search(data);
  const best = plans[0];
  fail('EV roadtrip derivation failed', {
    'bounded search finds eight acceptable plans': plans.length === 8,
    'best plan is fastest acceptable candidate': JSON.stringify(best.actions) === JSON.stringify(['drive_bru_liege', 'drive_liege_aachen', 'shuttle_aachen_cologne']),
    'best reaches Cologne': best.final.At === 'Cologne',
    'best cost is recomputed': Math.abs(best.cost - 0.054) < 1e-12,
    'best belief is recomputed': Math.abs(best.belief - 0.974175) < 1e-6,
    'best fuel remaining is stable': best.fuel_remaining === 5,
    'best satisfies thresholds': best.belief > data.Thresholds.MinBelief && best.cost < data.Thresholds.MaxCost && best.duration < data.Thresholds.MaxDuration,
  });
  return plans;
}

function main() {
  const data = loadInput(NAME);
  const plans = trustedDerivation(data);
  const best = plans[0];
  const start = new State(data.Vehicle.At, data.Vehicle.Battery, data.Vehicle.Pass);
  const thresholds = data.Thresholds;

  emit('# EV Roadtrip Planner');
  emit();
  emit('## Insight');
  emit(`Select plan : ${best.actions.join(' -> ')}.`);
  emit(`route result : ${best.final.At} battery=${best.final.Battery} pass=${best.final.Pass}`);
  emit(`duration : ${best.duration.toFixed(1)} minutes`);
  emit(`cost : ${best.cost.toFixed(3)}`);
  emit(`belief : ${best.belief.toFixed(6)}`);
  emit(`comfort : ${best.comfort.toFixed(6)}`);
  emit(`acceptable plans : ${plans.length}`);
  emit(`fuel remaining : ${best.fuel_remaining} of ${data.FuelSteps}`);
  emit();
  emit('## Explanation');
  emit(`The planner starts with ${data.Vehicle.ID} at ${start.At}, battery=${start.Battery}, pass=${start.Pass}, then composes action descriptions until the goal city ${data.Goal.At} is reached.`);
  emit('Duration and cost are summed across each candidate; belief and comfort are multiplied, matching the N3 planner pattern.');
  emit(`The selected plan is the fastest acceptable candidate under belief > ${thresholds.MinBelief.toFixed(2)}, cost < ${thresholds.MaxCost.toFixed(3)}, and duration < ${thresholds.MaxDuration.toFixed(1)}.`);
  emit(`It uses the shuttle from Aachen to Cologne, avoiding an extra charge stop while keeping belief at ${best.belief.toFixed(6)}.`);
  emit();
  emit('Top acceptable plans:');
  plans.slice(0, 5).forEach((plan, i) => {
    emit(`${i + 1}. ${plan.actions.join(' -> ')} | duration=${plan.duration.toFixed(1)} cost=${plan.cost.toFixed(3)} belief=${plan.belief.toFixed(6)} comfort=${plan.comfort.toFixed(6)} final=${plan.final.At}/${plan.final.Battery}/${plan.final.Pass}`);
  });
}

if (require.main === module) main();
