#!/usr/bin/env node
const { emit, fail, loadInput, maxBy, minBy } = require('./_see');

const NAME = 'school_placement_audit';

function schoolNameById(data) {
  return Object.fromEntries(data.schools.map((school) => [school.id, school.name]));
}

function distancesByStudent(data) {
  const rows = {};
  for (const row of data.distances) {
    if (!(row.student in rows)) rows[row.student] = [];
    rows[row.student].push(row);
  }
  return rows;
}

function preferenceRank(student, schoolName) {
  return student.preferences.indexOf(schoolName);
}

function straightAssignment(student, distances) {
  return minBy(distances, (row) => [row.straightMeters, preferenceRank(student, row.school)]);
}

function auditedAssignment(student, distances, penalty) {
  return minBy(distances, (row) => [row.walkingMeters + penalty * preferenceRank(student, row.school), row.walkingMeters]);
}

function analyze(data) {
  const penalty = Number.parseInt(data.policy.preferencePenaltyMeters, 10);
  const maxWalk = Number.parseInt(data.policy.maxWalkingMeters, 10);
  const byStudent = distancesByStudent(data);
  const rows = [];
  for (const student of data.students) {
    const sid = student.id;
    const straight = straightAssignment(student, byStudent[sid]);
    const audited = auditedAssignment(student, byStudent[sid], penalty);
    const hiddenDetour = Number.parseInt(straight.walkingMeters, 10) - Number.parseInt(straight.straightMeters, 10);
    const affected = straight.school !== audited.school || Number.parseInt(straight.walkingMeters, 10) > maxWalk;
    rows.push({ student, straight, audited, hidden_detour: hiddenDetour, affected });
  }
  return { rows, max_walk: maxWalk, penalty };
}

function trustedDerivation(data) {
  const result = analyze(data);
  const rows = result.rows;
  const affected = rows.filter((row) => row.affected).map((row) => row.student.name);
  const recommended = Object.fromEntries(rows.map((row) => [row.student.name, row.audited.school]));
  const largest = maxBy(rows, (row) => row.hidden_detour);
  fail('School placement audit derivation failed', {
    'complete 4x4 distance matrix': data.students.length === 4 && data.schools.length === 4 && data.distances.length === 16,
    'Ada and Björn are affected': affected.includes('Ada') && affected.includes('Björn'),
    'Davi is affected by preference-aware audit': affected.includes('Davi'),
    'Clara remains at Haga': recommended.Clara === 'Haga',
    'Ada has largest hidden detour': largest.student.name === 'Ada' && largest.hidden_detour === 3000,
  });
  return result;
}

function main() {
  const data = loadInput(NAME);
  const result = trustedDerivation(data);
  const rows = result.rows;
  const affected = rows.filter((row) => row.affected).map((row) => row.student.name);
  const largest = maxBy(rows, (row) => row.hidden_detour);
  const assignments = rows.map((row) => `${row.student.name} -> ${row.audited.school}`).join('; ');

  emit('# School Placement Route Audit');
  emit();
  emit('## Insight');
  emit('audit result : fail');
  emit(`children affected by straight-line rule : ${affected.join(', ')}`);
  emit(`largest hidden detour : ${largest.student.name}, ${largest.hidden_detour} m`);
  emit(`recommended assignments : ${assignments}`);
  emit('explanation requested : yes');
  emit();
  emit('## Explanation');
  emit('The support-tool rule chooses the school with the smallest straight-line distance, using preference rank only as a tie-breaker.');
  emit('The independent audit recomputes each candidate with walking-route distance plus 600 m per preference step.');
  emit('Any provisional assignment that is not the audited best, or that requires more than 2500 m of walking, is flagged.');
  emit('Ada and Björn look close to Centrum on a map, but their walking routes cross barriers and exceed the walking limit; Davi is also better served by the first-preference Haga route.');
  emit('This illustrates why a decision-support label is not enough: route geometry, preferences, and audit records must be inspectable.');
}

if (require.main === module) main();
