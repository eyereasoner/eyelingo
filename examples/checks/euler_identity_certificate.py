# Independent Python checks for the euler_identity_certificate example.
from .common import run_fragment_checks

CHECKS = [
    ('Taylor expansion used 28 terms from the JSON input', ['terms used : 28']),
    ('computed real part is close to -1', ['The computed real part is effectively -1 and the imaginary part is near 0 at the chosen precision.']),
    ('computed imaginary part is close to 0', ['computed imaginary part of exp(iπ) : -0.000000000000000']),
    ('|exp(iπ)+1| is below the configured tolerance', ['expression : exp(iπ) + 1']),
    ('the audit records the finite residual rather than asserting exact real arithmetic', ['The resulting residual is not claimed to be mathematically exact zero; it is checked against the explicit tolerance from JSON.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
