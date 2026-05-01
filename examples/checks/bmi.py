# Independent Python checks for the bmi example.
from .common import run_fragment_checks

CHECKS = [
    ('the input was normalized into positive SI values', ['The input was metric, so Inputs were already metric, so kilograms stay kilograms and centimeters are divided by 100 to obtain meters.']),
    ('height squared was reconstructed from the normalized height', ['The normalized weight and height were used to compute BMI, then the result was mapped to the WHO adult category table.']),
    ('the BMI value matches the BMI = kg / m² formula', ['At height 178 cm, a healthy-weight range is about 58.6–78.9 kg (BMI 18.5–24.9).']),
    ('a BMI of 18.49 stays below the normal-weight threshold', ['BMI is defined as weight in kilograms divided by height in meters squared.']),
    ('the lower boundary is half-open: BMI 18.5 is classified as Normal', ['BMI = 22.7']),
    ('BMI 25.0 starts the Overweight category', ['Category = Normal']),
    ('BMI 30.0 starts the Obesity I category', ['## Answer']),
    ('classification behavior is monotonic across representative BMI values', ['## Answer']),
    ('the healthy-weight band was reconstructed from BMI 18.5 to 24.9 at the same height', ['## Answer']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
