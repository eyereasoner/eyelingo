# Independent Python checks for the rc_discharge_envelope example.
from .common import run_fragment_checks

CHECKS = [
    ('decay certificate is nonempty, positive, and below 1', ['The physical decay factor is exp(-1/4), but the example uses a finite double interval as the certificate.']),
    ('voltage upper envelope decreases at every sample', ['The upper envelope is the safety-relevant bound: once it falls below 1.0 V, every compatible exact trajectory is below tolerance.']),
    ('step 12 remains above the voltage tolerance', ['upper voltage at step 13 : 0.930581 V']),
    ('step 13 is the first certified discharge step', ['first below tolerance step : 13']),
    ('settling time is 0.325 s', ['first below tolerance time : 0.325 s']),
    ('the certificate uses the JSON double interval rather than an exact transcendental value', ['Because the interval lies strictly between 0 and 1, the capacitor voltage envelope contracts each sample.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
