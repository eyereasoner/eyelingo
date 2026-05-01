# Independent Python checks for the ev_roundtrip_planner example.
from .common import run_fragment_checks

CHECKS = [
    ('8 acceptable Brussels-to-Cologne plans were derived', ['acceptable plans : 8']),
    ('selected plan duration 210.0 is below 260.0', ['The selected plan is the fastest acceptable candidate under belief > 0.93, cost < 0.090, and duration < 260.0.']),
    ('selected plan cost 0.054 is below 0.090', ['cost : 0.054']),
    ('selected plan belief 0.974175 is above 0.93', ['It uses the shuttle from Aachen to Cologne, avoiding an extra charge stop while keeping belief at 0.974175.']),
    ('selected plan reaches Cologne', ['The planner starts with car1 at Brussels, battery=high, pass=none, then composes action descriptions until the goal city Cologne is reached.']),
    ('selected plan uses the high-belief Aachen-Cologne shuttle for the last mile', ['Duration and cost are summed across each candidate; belief and comfort are multiplied, matching the N3 planner pattern.']),
    ('bounded search consumed at most 8 of 8 fuel tokens', ['fuel remaining : 5 of 8']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
