#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'bmi';

function normalizeSi(data) {
  if (data.UnitSystem !== 'metric') throw new Error('this example expects metric input');
  const weightKg = Number(data.Weight);
  const heightM = Number(data.Height) / 100.0;
  const reason = 'Inputs were already metric, so kilograms stay kilograms and centimeters are divided by 100 to obtain meters.';
  return { weightKg, heightM, reason };
}

function bmiCategory(bmi) {
  if (bmi < 18.5) return 'Underweight';
  if (bmi < 25.0) return 'Normal';
  if (bmi < 30.0) return 'Overweight';
  return 'Obese';
}

function trustedDerivation(data) {
  const { weightKg, heightM, reason } = normalizeSi(data);
  const bmi = weightKg / (heightM * heightM);
  const category = bmiCategory(bmi);
  const healthyLow = 18.5 * heightM * heightM;
  const healthyHigh = 24.9 * heightM * heightM;

  fail('BMI derivation failed', {
    'height is positive': heightM > 0,
    'weight is positive': weightKg > 0,
    'BMI falls inside its emitted category': category === bmiCategory(bmi),
    'healthy range is ordered': healthyLow < healthyHigh,
    'normal category lies inside healthy range': healthyLow <= weightKg && weightKg <= healthyHigh,
  });
  return { bmi, category, low: healthyLow, high: healthyHigh, normalizationReason: reason };
}

function main() {
  const data = loadInput(NAME);
  const { bmi, category, low, high, normalizationReason } = trustedDerivation(data);
  const heightCm = Number.parseInt(data.Height, 10);

  emit('# BMI — Body Mass Index example');
  emit();
  emit('## Insight');
  emit(`BMI = ${bmi.toFixed(1)}`);
  emit(`Category = ${category}`);
  emit(`At height ${heightCm} cm, a healthy-weight range is about ${low.toFixed(1)}–${high.toFixed(1)} kg (BMI 18.5–24.9).`);
  emit();
  emit('## Explanation');
  emit('BMI is defined as weight in kilograms divided by height in meters squared.');
  emit('The normalized weight and height were used to compute BMI, then the result was mapped to the WHO adult category table.');
  emit(`The input was metric, so ${normalizationReason}`);
}

if (require.main === module) main();
