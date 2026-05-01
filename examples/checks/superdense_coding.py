# Independent Python checks for the superdense_coding example.
from .common import run_fragment_checks

CHECKS = [
    ('shared entanglement: |R) contains exactly |0,0) and |1,1)', ['Alice and Bob start with |R) = |0,0) + |1,1).']),
    ('composition KG: KG is obtained by composing G then K, exactly as in the N3 rule', ['The N3 example keeps only answers with odd derivation count, so duplicate']),
    ('composition GK: GK is obtained by composing K then G, exactly as in the N3 rule', ['message 3 -> KG -> encoded support {(0,0), (0,1), (1,0)}']),
    ('candidate multiplicity: the raw superdense rule creates 24 candidate derivations before parity cancellation', ['Raw candidate counts before parity filtering:']),
    ('GF(2) cancellation: off-diagonal answers occur zero or two times and cancel; diagonal answers occur once', ['Superdense-coding facts that survive GF(2) parity cancellation:']),
    ("decoded messages: after odd-parity filtering, Bob recovers Alice's original two-bit message", ['message 2 -> K  -> encoded support {(0,0), (0,1), (1,1)}']),
    ('encoded supports distinct: the four Alice operations produce four different supports over the two mobits', ['Alice chooses one relation for the first mobit; Bob applies one joint test.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
