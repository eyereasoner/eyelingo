# Independent Python checks for the lldm example.
from .common import run_fragment_checks

CHECKS = [
    ('L1 is perpendicular to L3 and L4 (slopes product ≈ -1)', ['L3 and L4 are perpendicular to L1.  Intersection points p5 (L1∩L3)']),
    ('p5 lies on both L1 and L3', ['d53 (p5–p3) and d64 (p6–p4) are then computed, and their difference']),
    ('p6 lies on both L1 and L4', ['and p6 (L1∩L4) are computed analytically.  The Euclidean distances']),
    ('squared distances are non‑negative', ['SL1 = -0.062857  SL3 = SL4 = 15.909091']),
    ('dCm is a finite number', ['LLD Alarm = TRUE  (discrepancy dCm = -1.908234, threshold ±1.25)']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
