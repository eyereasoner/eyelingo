# Bayes Therapy Decision Support

## What this example is about

This example extends the Bayesian diagnosis idea into a simple treatment-ranking exercise. It is a toy decision-support example, not medical advice and not a recommendation for real treatment.

The input describes possible diseases, symptoms, and therapies. Each therapy has an estimated success rate for each disease and an adverse-effect penalty.

## How it works, in plain language

The first step is the same as the Bayes diagnosis example: symptoms are used to compute posterior probabilities for each disease.

The second step asks a practical question: given that uncertainty, which therapy has the best expected trade-off? For each therapy, the program multiplies each disease posterior by the therapy success rate for that disease. It adds those weighted values to get expected success. Then it subtracts a harm penalty for adverse effects.

The therapy with the highest utility score wins.

## What to notice in the output

The output prints both the disease posteriors and the therapy scores. That lets a reader see why the chosen therapy is not simply attached to one diagnosis. It is selected because it performs best after averaging over all still-plausible diagnoses and applying the configured benefit/harm weights.

## What the trust gate checks

The trust gate verifies that priors look like probabilities, the Bayesian normalizer is positive, posteriors sum to one, every disease has evidence likelihoods for every symptom, therapy vectors align with the disease list, and the emitted therapy has the maximum computed utility.

## Run it

From the repository root:

```sh
node examples/bayes_therapy.js
```

## Files

- [JavaScript example](../bayes_therapy.js)
- [Input data](../input/bayes_therapy.json)
- [Reference output](../output/bayes_therapy.md)
