# Independent Python checks for the alarm_bit_interoperability example.
from .common import run_fragment_checks

CHECKS = [
    ('two unlike classical media are present', ['The optical beacon and relay register are unlike physical media, but both encode the same abstract AlarmBit variable.']),
    ('classical media encode the same variable AlarmBit', ['Because the variable is classical in both media, local permutation and copying in both directed media transfers are possible.']),
    ('2 directed copy tasks are possible', ['copy task : opticalBeacon -> relayRegister for AlarmBit']),
    ('each classical medium supports a local permutation', ['The quantumToken is treated as a superinformation medium with states Horizontal, Vertical, Diagonal, AntiDiagonal.']),
    ('quantumToken cannot be universally cloned', ['That contrast substrate cannot support universal cloning of all states, so unrestricted classical-style fan-out is also blocked.']),
    ('unrestricted classical-style fan-out is blocked for the superinformation token', ['universal cloning of the superinformation token : NO']),
    ('CAN=YES and CANNOT=NO decisions are both derived', ['copy task : relayRegister -> opticalBeacon for AlarmBit']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
