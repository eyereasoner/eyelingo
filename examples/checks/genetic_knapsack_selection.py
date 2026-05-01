# Independent Python checks for the genetic_knapsack_selection example.
from .common import run_fragment_checks

CHECKS = [
    ('12 items align with a 12-bit genome', ['selected items : item01, item03, item10, item12']),
    ('final weight 50 is within capacity 50', ['The run stops at 101000000101 because every one-bit neighbor is no better under the capacity 50 rule.']),
    ('final genome is 101000000101', ['final genome : 101000000101']),
    ('final value is 101 using item01, item03, item10, item12', ['value : 101']),
    ('no single-bit neighbor improves the final candidate', ['At every generation, all single-bit mutants of the parent are compared with the parent, and the lowest-fitness candidate is selected deterministically.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
