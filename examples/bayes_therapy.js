#!/usr/bin/env node
const { emit, fail, loadInput, maxBy, sum } = require('./_see');

const NAME = 'bayes_therapy';

function diseaseScores(data) {
  let rows = [];
  for (const disease of data.Diseases) {
    const name = disease.Name;
    let likelihood = Number(disease.Prior);
    for (const ev of data.Evidence) {
      const p = Number(data.ProbGiven[name][ev.Symptom]);
      likelihood *= ev.Present ? p : 1.0 - p;
    }
    rows.push([name, likelihood, 0.0]);
  }
  const total = sum(rows.map((row) => row[1]));
  return rows.map(([name, likelihood]) => [name, likelihood, likelihood / total]);
}

function therapyScores(data, posteriors) {
  const scores = [];
  for (const therapy of data.Therapies) {
    const expectedSuccess = sum(posteriors.map((row, i) => row[2] * Number(therapy.SuccessByDisease[i])));
    const adverse = Number(therapy.Adverse);
    const utility = Number(data.BenefitWeight) * expectedSuccess - Number(data.HarmWeight) * adverse;
    scores.push([therapy.Name, expectedSuccess, adverse, utility]);
  }
  return scores;
}

function trustedDerivation(data) {
  const posteriors = diseaseScores(data);
  const total = sum(posteriors.map((row) => row[1]));
  const scores = therapyScores(data, posteriors);
  const winner = maxBy(scores, (row) => row[3]);
  const symptoms = new Set(data.Evidence.map((ev) => ev.Symptom));

  fail('Bayes therapy derivation failed', {
    'priors are probabilities': data.Diseases.every((d) => 0 <= d.Prior && d.Prior <= 1),
    'normalizer is positive': total > 0,
    'posteriors sum to one': Math.abs(sum(posteriors.map((row) => row[2])) - 1.0) < 1e-12,
    'evidence available for all diseases': data.Diseases.every((d) => [...symptoms].every((s) => s in data.ProbGiven[d.Name])),
    'therapy vectors align': data.Therapies.every((t) => t.SuccessByDisease.length === data.Diseases.length),
    'winner is stable': winner[0] === 'Paxlovid',
  });
  return { posteriors, scores, winner };
}

function main() {
  const data = loadInput(NAME);
  const { posteriors, scores, winner } = trustedDerivation(data);
  const total = sum(posteriors.map((row) => row[1]));

  emit('# Bayes Therapy Decision Support');
  emit();
  emit('## Insight');
  emit(`The recommended therapy is ${winner[0]} (utility = ${winner[3].toFixed(6)}).`);
  emit();
  emit('Full posterior distribution:');
  for (const [name, likelihood, posterior] of posteriors) {
    emit(`  ${name.padEnd(21)} posterior = ${posterior.toFixed(6)}  (unnormalized = ${likelihood.toFixed(8)})`);
  }
  emit();
  emit('Therapy utility scores:');
  for (const [name, expectedSuccess, adverse, utility] of scores) {
    emit(`  ${name.padEnd(21)} expectedSuccess = ${expectedSuccess.toFixed(6)}  adverse = ${adverse.toFixed(2)}  utility = ${utility.toFixed(6)}`);
  }
  emit();
  emit('## Explanation');
  const evidence = data.Evidence.map((ev) => `${ev.Symptom}=${ev.Present ? 'present' : 'absent'}`).join(', ');
  emit(`Evidence: ${evidence}.`);
  emit(`Evidence total (normalizing constant) = ${total.toFixed(8)}.`);
  emit();
  emit('The posterior for each disease is computed as:');
  emit('  posterior(d) = prior(d) × ∏ P(symptom|d) / evidenceTotal');
  emit('where for an absent symptom the factor is 1 − P(symptom|d).');
  emit();
  emit('For each therapy, expected success is:');
  emit('  expectedSuccess(t) = Σ_i posterior(i) × successByDisease(i)');
  emit('and utility = benefitWeight × expectedSuccess − harmWeight × adverse.');
  emit('The recommended therapy is the one with the highest utility.');
}

if (require.main === module) main();
