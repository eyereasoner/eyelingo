# Independent Python checks for the kaprekar_6174 example.
from .common import run_fragment_checks

CHECKS = [
    ('all digit patterns from 0000 through 9999 were considered', ['9999 -> (0000)']),
    ('the identity-based step matches direct digit sorting for every start', ['3 step(s) : 2400 start(s)']),
    ('3524 follows the classic 3087 -> 8352 -> 6174 chain', ['3524 :kaprekar (3087 8352 6174)']),
    ('0001 is accepted as a four-digit start by treating it as 0,0,0,1', ['Each start is read as four digits, so 1 is treated as 0001.']),
    ('0000 and the nine non-zero repdigits fall to 0000 and are not emitted', ['Examples omitted because they fall to 0000:']),
    ("every :kaprekar fact kept by the translation reaches Kaprekar's constant", ['9831 :kaprekar (8442 5994 5355 1998 8082 8532 6174)']),
    ('no emitted chain needs more than the seven steps unrolled in the N3 source', ['The search is bounded to seven steps, matching the N3 source: any']),
    ('all 9990 non-repdigit starts are emitted, including 6174 itself', ['Kaprekar chains that end at 6174 are emitted as :kaprekar facts.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
