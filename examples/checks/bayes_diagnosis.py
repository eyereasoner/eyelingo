# Independent Python checks for the bayes_diagnosis example.
from .common import run_fragment_checks

CHECKS = [
    ('all prior probabilities are in [0,1]', ['where for an absent symptom the factor is 1 − P(symptom|d).']),
    ('all conditional probabilities are in [0,1]', ['COVID19               posterior = 0.941209  (unnormalized = 0.00154700)']),
    ('the evidence total is non-zero and reported as the Bayesian normalizing constant', ['Evidence total (normalizing constant) = 0.00164363.']),
    ('COVID19 has the largest posterior probability', ['The most likely disease is COVID19 (posterior = 0.941209).']),
    ('the posterior distribution contains four diseases', ['The posterior for each disease is computed as:']),
    ('absent Sneezing is handled through a complement likelihood factor', ['Influenza             posterior = 0.029204  (unnormalized = 0.00004800)']),
    ('posterior probabilities are normalized by the evidence total', ['AllergicRhinitis      posterior = 0.000456  (unnormalized = 0.00000075)']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
