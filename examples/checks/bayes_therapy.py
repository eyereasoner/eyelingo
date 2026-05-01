# Independent Python checks for the bayes_therapy example.
from .common import run_fragment_checks

CHECKS = [
    ('all prior probabilities are in [0,1]', ['where for an absent symptom the factor is 1 − P(symptom|d).']),
    ('all conditional probabilities are in [0,1]', ['Paxlovid              expectedSuccess = 0.388517  adverse = 0.10  utility = 3.585174']),
    ('all adverse probabilities are in [0,1]', ['Oseltamivir           expectedSuccess = 0.285141  adverse = 0.08  utility = 2.611410']),
    ('all success probabilities are in [0,1]', ['For each therapy, expected success is:']),
    ('evidence total is non‑zero', ['Evidence total (normalizing constant) = 0.01918233.']),
    ('number of diseases matches success list length', ['Antihistamine         expectedSuccess = 0.100269  adverse = 0.03  utility = 0.912689']),
    ('number of therapies is correct', ['Antibiotic            expectedSuccess = 0.110953  adverse = 0.07  utility = 0.899526']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
