# Bayes Therapy Decision Support

`bayes_therapy` is a Go translation/adaptation of Eyeling's `bayes-therapy-decision-support.n3`.

## Background

Therapy choice under uncertainty usually combines two steps: estimate what condition is likely, then evaluate how good each treatment would be under those possibilities. Expected utility does this by weighting benefits, harms, and penalties by posterior probabilities. A treatment can therefore rank highly because it performs well for likely diagnoses or because it avoids severe downside in less likely ones.

## What it demonstrates

**Category:** Science. Posterior-weighted utility selection of the best therapy.

## How the Go implementation works

The implementation first computes disease posteriors from the evidence, then evaluates therapy candidates with an expected-utility calculation. Each therapy's utility is the posterior-weighted combination of outcomes, benefits, and penalties encoded in the fixture.

The program ranks the therapies, selects the highest utility option, and checks both the Bayesian normalization and the decision ordering.

## Files

Input JSON: [../input/bayes_therapy.json](../input/bayes_therapy.json)

Go implementation: [../bayes_therapy.go](../bayes_therapy.go)

Expected Markdown output: [../output/bayes_therapy.md](../output/bayes_therapy.md)
