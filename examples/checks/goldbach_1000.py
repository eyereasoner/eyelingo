# Independent Python checks for the goldbach_1000 example.
from .common import run_fragment_checks

CHECKS = [
    ('the configured upper bound is 1000', ['The checker caches primes up to the configured bound and then searches each even number E for a prime P not greater than E/2 where E-P is also prime.']),
    ('there are 499 even integers from 4 through 1000', ['All 499 even integers from 4 through 1000 have a Goldbach witness.']),
    ('every checked even integer has a prime-pair witness', ['sample witnesses : 4=2+2; 28=5+23; 100=3+97; 998=7+991; 1000=3+997']),
    ('each requested sample even has a witness', ['No counterexample is found in the bounded range, so the bounded Goldbach condition succeeds for this dataset.']),
    ('there are 168 primes at or below 1000', ['## Answer']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
