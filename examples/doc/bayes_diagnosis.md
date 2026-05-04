# Bayes Diagnosis

## What this example is about

This example shows how a small Bayesian classifier turns symptoms into a ranked list of possible diseases. It is a toy reasoning demo, not medical advice and not a clinical diagnostic tool.

The input contains four candidate diseases, a prior probability for each one, and symptom likelihoods such as `P(Fever | COVID19)`. The evidence says which symptoms are present or absent.

## How it works, in plain language

The program starts with the prior probability for each disease. It then adjusts that score by multiplying in the likelihood of every observed symptom. For symptoms that are absent, it uses the opposite likelihood: `1 - P(symptom | disease)`.

Those scores are still unnormalized, so the program divides each score by the total score across all diseases. The result is a posterior distribution: probabilities that add up to 1.

## What to notice in the output

The output shows COVID19 with the largest posterior probability. It also prints the full posterior distribution, so the reader can see that other possibilities were not ignored. This is important because a good explanation should show the competing alternatives, not only the final winner.

## What the trust gate checks

Before printing the explanation, the example verifies that the normalizing total is positive, the posterior probabilities sum to one, the emitted winner really has the maximum posterior probability, and all posteriors are valid probabilities. These checks support the result without relying on a stored answer.

## Run it

From the repository root:

```sh
node examples/bayes_diagnosis.js
```

## Files

- [JavaScript example](../bayes_diagnosis.js)
- [Input data](../input/bayes_diagnosis.json)
- [Reference output](../output/bayes_diagnosis.md)
