# Independent Python checks for the fundamental_theorem_arithmetic example.
from .common import run_fragment_checks

CHECKS = [
    ('the source example factors 202692987 as 3,3,7,829,3881', ['Primary N3 case: n = 202692987 has prime factors 3 * 3 * 7 * 829 * 3881.']),
    ('the source example groups repeated factors as 3^2 * 7 * 829 * 3881', ['source smallest-first factors : 3 * 3 * 7 * 829 * 3881']),
    ('multiplying each computed factor list reconstructs its original number', ['Uniqueness up to order is checked by reversing each traversal and sorting']),
    ('every distinct factor in every decomposition is prime by trial division', ['prime factors, even when the factors were discovered in the opposite order.']),
    ('smallest-first and largest-first traversals sort to the same multisets', ['both factor lists. Matching sorted lists describe the same multiset of']),
    ('the extended sample includes six cases and includes the ten-digit prime 9999999967', ['9999999967 = 9999999967']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
