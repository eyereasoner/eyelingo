# Independent Python checks for the wind_turbine example.
from .common import run_fragment_checks

CHECKS = [
    ('cut-in, rated, and cut-out thresholds are ordered', ['operating thresholds : cut-in 3.5 m/s, rated 12.0 m/s, cut-out 25.0 m/s']),
    ('usable intervals are exactly the samples inside the operating envelope', ['usable intervals : 4']),
    ('two intervals reach rated power', ['Wind between cut-in and rated speed follows a cubic power curve normalized to the rated point.']),
    ('two intervals are stopped by low or cut-out wind', ['Wind below cut-in and at or above cut-out is stopped for production and safety.']),
    ('a below-rated usable wind speed follows the cubic power curve', ['Wind between rated speed and cut-out is capped at rated power.']),
    ('the six ten-minute samples yield about 1.57 MWh', ['Energy is accumulated by multiplying each interval power by the ten-minute interval duration.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
