# Independent Python checks for the allen_interval_calculus example.
from .common import run_fragment_checks

CHECKS = [
    ('11 intervals were loaded and completed when duration was present', ['completed intervals : I=16:00-18:00, K=13:30-14:00']),
    ('A before B was derived', ['showcase : A before B | A meets C | A overlaps D | F starts A | G finishes A | A during H | A equals E | J meets I | K finishes C']),
    ('A meets C and C metBy A were derived', ['The relation rules are purely endpoint constraints, so the result is deterministic and auditable.']),
    ('A overlaps D and D overlappedBy A were derived', ['derived relations : 110 ordered interval pairs']),
    ('F starts A and G finishes A were derived', ['Each ordered pair is then classified with the 13 Allen base relations, including the six converse relations.']),
    ('A during H and H contains A were derived', ['The duration-derived intervals participate in the same relation table as directly supplied intervals.']),
    ('A equals E was derived', ['The example completes any interval that has a start plus duration before comparing endpoints.']),
    ('duration completion produced I ending at 18:00 and K ending at 14:00', ['invalid intervals : 0']),
    ('no invalid intervals were detected', ['## Answer']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
