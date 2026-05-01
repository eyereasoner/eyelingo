# Independent Python checks for the fft32_numeric example.
from .common import run_fragment_checks

CHECKS = [
    ('each generated waveform has exactly 32 samples', ['The upstream FFT32 fixture defines several 32-sample waveforms and asks for the whole spectrum of each waveform.']),
    ('dominant bins match the FFT32 fixture expectations', ['The checks verify dominant bins, magnitudes, flat impulse spectrum, conjugate symmetry, and energy preservation.']),
    ('dominant magnitudes match the configured certificates', ['The Go translation evaluates every frequency bin by summing samples against the corresponding complex root of unity.']),
    ('Parseval energy is preserved for every waveform under the unnormalized DFT convention', ['cosine_bin3 : k=3 magnitude=16.000000 phase=0.000000; k=29 magnitude=16.000000 phase=0.000000; energy=16.000000']),
    ('real-valued waveforms produce conjugate-symmetric spectra', ['waveforms evaluated : 6']),
    ('the impulse waveform has a flat unit-magnitude spectrum', ['impulse : all 32 bins magnitude=1.000000; energy=1.000000']),
    ('the full spectrum is computed once per waveform', ['sine_bin5 : k=5 magnitude=16.000000 phase=-1.570796; k=27 magnitude=16.000000 phase=1.570796; energy=16.000000']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
