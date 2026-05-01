# Independent Python checks for the complex_matrix_stability example.
from .common import run_fragment_checks

CHECKS = [
    ('A_unstable has spectral radius 2, so it is unstable', ['A_unstable : spectral radius 2 -> unstable']),
    ('A_stable has spectral radius 1, so it is marginally stable', ['A_stable : spectral radius 1 -> marginally stable']),
    ('A_damped has spectral radius 0, so every mode decays', ['A_damped : spectral radius 0 -> damped']),
    ('squared modulus of z*w equals the product of squared moduli', ['Because the matrices are diagonal, the eigenvalues are the diagonal entries; the largest modulus gives the spectral radius and therefore the stability class.']),
    ('scaling A_unstable by 2 multiplies spectral-radius-squared by 4', ['For a discrete-time linear system x_{k+1} = A x_k, the eigenvalues of A govern the modal behaviour.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
