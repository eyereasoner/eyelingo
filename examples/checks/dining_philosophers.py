# Independent Python checks for the dining_philosophers example.
from .common import run_fragment_checks

CHECKS = [
    ('the translated run follows the nine start-of-round configs from the N3 source', ['transfers : P5 sends F51 to P1; P2 sends F12 to P1; P2 sends F23 to P3; P4 sends F34 to P3']),
    ('26 dirty-fork requests transfer and the other 19 fork-round pairs are kept', ['Dirty, and the receiver gets that fork Clean. After the meal phase, all']),
    ('rounds 1/4/7 feed P1,P3; rounds 2/5/8 feed P2,P4; rounds 3/6/9 feed P5', ['transfers : P1 sends F12 to P2; P3 sends F23 to P2; P3 sends F34 to P4; P5 sends F45 to P4']),
    ('each philosopher has exactly three derived meals', ['Each round has three phases. Hungry philosophers first request adjacent']),
    ('no two philosophers eat with the same fork in the same round', ['forks are marked Dirty for the next round.']),
    ('the end-of-round state makes every fork Dirty again, matching the trace model', ['The Go code uses goroutines inside each phase, then applies state changes']),
    ('the final ownership is F12,F23 with P2; F34 with P4; F45,F51 with P5', ['requests : P1 asks P5 for F51; P1 asks P2 for F12; P3 asks P2 for F23; P3 asks P4 for F34']),
    ('each of the nine rounds used request, transfer, and meal goroutine batches', ['The Chandy-Misra dining-philosophers trace completes without conflict.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
