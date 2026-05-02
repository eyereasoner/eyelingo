# Bayes Diagnosis

`bayes_diagnosis` is a Go translation/adaptation of Eyeling's `bayes-diagnosis.n3`.

## Background

Bayesian diagnosis starts with prior probabilities for possible diseases and updates them with evidence such as symptoms or test results. Each hypothesis receives a likelihood score, and the scores are normalized so they sum to one as posterior probabilities. The result is not a certain diagnosis; it is a ranked set of explanations under the assumptions encoded in the data.

## What it demonstrates

**Category:** Science. Bayesian posterior ranking of possible diseases from symptoms and test evidence.

## How the Go implementation works

The Go code keeps diseases, priors, symptoms, and evidence as ordinary data. For each disease it multiplies the prior by the appropriate likelihood factors, normalizes the scores into posterior probabilities, and sorts the diseases by posterior rank.

Guard checks validate the probability table before inference, while output checks verify normalization, ranking, and the selected diagnosis.

## Files

Input JSON: [../input/bayes_diagnosis.json](../input/bayes_diagnosis.json)

Go implementation: [../bayes_diagnosis.go](../bayes_diagnosis.go)

Python check: [../checks/bayes_diagnosis.py](../checks/bayes_diagnosis.py)

Expected Markdown output: [../output/bayes_diagnosis.md](../output/bayes_diagnosis.md)
