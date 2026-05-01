# Independent Python checks for the fft8_numeric example.
from .common import run_fragment_checks

CHECKS = [
    ('the input contains exactly 8 time-domain samples', ['The input samples describe one sine cycle over eight equally spaced samples.']),
    ('the dominant bins are k=1 and k=7', ['dominant bins : k=1 magnitude=4.000000 phase=-1.570796; k=7 magnitude=4.000000 phase=1.570796']),
    ('the DC component is zero', ['All non-dominant bins cancel to zero within the configured numerical tolerance.']),
    ('the two conjugate sine bins have magnitude 4', ['A real sine wave has equal magnitude at the positive and negative frequency bins.']),
    ('Parseval energy is preserved by the unnormalized DFT convention', ['The DFT projects the signal onto eight complex roots of unity.']),
    ('real-valued samples produce conjugate-symmetric bins', ['sample vector : 0.000000, 0.707107, 1.000000, 0.707107, 0.000000, -0.707107, -1.000000, -0.707107']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
