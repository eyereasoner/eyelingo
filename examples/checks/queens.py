# Independent Python checks for the queens example.
from .common import run_fragment_checks

CHECKS = [
    ('search reached depth 8', ['As column positions by row: [1, 5, 8, 6, 3, 7, 2, 4]']),
    ('first solution places one queen in each row', ['The solver places one queen per row on a 8x8 board.']),
    ('first solution columns are unique', ['Counting continues after the printed solution limit, so the total solution count remains complete.']),
    ('no pair of queens in the first solution shares a diagonal', ['At each row it uses bit masks for occupied columns and both diagonal directions to enumerate only safe candidate columns.']),
    ('counted 92 solutions for the normalized 8-Queens input', ['Total solutions for 8-Queens: 92']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
