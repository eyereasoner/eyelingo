# Independent Python checks for the gray_code_counter example.
from .common import run_fragment_checks

CHECKS = [
    ('16 states were generated for a 4-bit counter', ['For 4 bits, the first 16 integers cover the full state space without duplicates.']),
    ('all generated states are unique', ['unique states : 16']),
    ('each adjacent transition flips exactly one bit', ['A valid cyclic Gray counter therefore changes exactly one bit at every step.']),
    ('the final state wraps to the first with one bit flip', ['The Hamming-distance check compares each state with the next state, including the final wraparound transition.']),
    ('first four states match the reflected binary Gray-code prefix', ['The counter maps each integer n to n xor (n >> 1), which is the reflected binary Gray-code construction.']),
    ('the numeric generator is n xor (n >> 1)', ['maximum adjacent Hamming distance : 1']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
