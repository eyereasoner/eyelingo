# Independent Python checks for the control_system example.
from .common import run_fragment_checks

CHECKS = [
    ('input1 measurement10 follows the lessThan backward rule', ['The first forward rule computes feedforward control from input1 measurement10, the input2 boolean guard, and disturbance1 compensation.']),
    ('disturbance2 measurement10 follows the notLessThan backward rule', ['disturbance2 measurement10 = 45.000000 (notLessThan branch)']),
    ('input2 boolean guard is true for the feedforward rule', ['The second forward rule computes PND feedback control from target/measurement error and the state/output differential error.']),
    ('feedforward control equals product minus log10 compensation', ['input1 measurement10 = 2.236068 (lessThan branch)']),
    ('feedback error is target minus measurement', ['For input1, 6 < 11, so sqrt(11 - 6) = 2.236068.']),
    ('differential error is state observation minus output measurement', ['The backward measurement10 rules first normalize measurement1 pairs into scalar measurement10 values.']),
    ('nonlinear differential part equals (7.3 / error) times differential error', ['actuator1 control1 = 39.273462']),
    ('actuator2 control is proportional part plus nonlinear differential part', ['actuator2 control1 = 26.080000']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
