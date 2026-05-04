#!/usr/bin/env node
const { emit, fail, loadInput, maxBy, roundTo, sum } = require('./_see');

const NAME = 'bayes_diagnosis';

function evidenceText(data) {
  return data.Evidence.map((ev) => `${ev.Symptom}=${ev.Present ? 'present' : 'absent'}`).join(', ');
}

function scoreDisease(data, disease) {
  const name = disease.Name;
  let score = Number(disease.Prior);
  for (const ev of data.Evidence) {
    const likelihood = Number(data.ProbGiven[name][ev.Symptom]);
    score *= ev.Present ? likelihood : 1.0 - likelihood;
  }
  return score;
}

function trustedDerivation(data) {
  let rows = data.Diseases.map((disease) => [disease.Name, scoreDisease(data, disease), 0.0]);
  const total = sum(rows.map((row) => row[1]));
  rows = rows.map(([name, unnormalized]) => [name, unnormalized, unnormalized / total]);
  const winner = maxBy(rows, (row) => row[2])[0];
  const posteriorByName = Object.fromEntries(rows.map(([name, , posterior]) => [name, posterior]));

  fail('Bayes diagnosis derivation failed', {
    'normalizer is positive': total > 0,
    'posteriors sum to one': Math.abs(sum(rows.map((row) => row[2])) - 1.0) < 1e-12,
    'COVID19 has maximum posterior': winner === 'COVID19',
    'COVID19 posterior is stable': roundTo(posteriorByName.COVID19, 6) === 0.941209,
  });
  return { rows, total, winner };
}

function main() {
  const data = loadInput(NAME);
  const { rows, total, winner } = trustedDerivation(data);
  const winnerPosterior = Object.fromEntries(rows.map(([name, , posterior]) => [name, posterior]))[winner];

  emit('# Bayes Diagnosis');
  emit();
  emit('## Insight');
  emit(`The most likely disease is ${winner} (posterior = ${winnerPosterior.toFixed(6)}).`);
  emit();
  emit('Full posterior distribution:');
  for (const [name, unnormalized, posterior] of rows) {
    emit(`  ${name.padEnd(21)} posterior = ${posterior.toFixed(6)}  (unnormalized = ${unnormalized.toFixed(8)})`);
  }
  emit();
  emit('## Explanation');
  emit(`Evidence: ${evidenceText(data)}.`);
  emit(`Evidence total (normalizing constant) = ${total.toFixed(8)}.`);
  emit('The posterior for each disease is computed as:');
  emit('  posterior(d) = prior(d) × ∏ P(symptom|d) / evidenceTotal');
  emit('where for an absent symptom the factor is 1 − P(symptom|d).');
}

if (require.main === module) main();
