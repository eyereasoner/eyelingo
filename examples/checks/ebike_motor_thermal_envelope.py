# Independent Python checks for the ebike_motor_thermal_envelope example.
from .common import run_fragment_checks

CHECKS = [
    ('cooling interval 0.7788007830..0.7788007831 is positive, ordered, and contractive', ['cooling certificate : exp(-1/4) in 0.7788007830 .. 0.7788007831']),
    ('temperature trace has 13 samples including the initial state', ['The model keeps an interval for motor temperature excess over ambient instead of pretending to know the transcendental cooling factor exactly.']),
    ('maximum upper temperature 40.2285 C stays below hard limit 45.0 C', ['The upper envelope returns below the 35.0 C warning limit at step 8 and remains below the 45.0 C hard limit throughout.']),
    ('warning recovery occurs at step 8 after 40 seconds', ['warning recovery : step 8 at 40 s']),
    ('upper envelope decreases from the first post-Turbo sample onward', ['Turbo pushes the upper envelope to 40.2285 C, then Tour, Eco, and Coast allow the envelope to decrease.']),
    ('ride decision is ThermallySafeForThisAssistPlan', ['decision : ThermallySafeForThisAssistPlan']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
