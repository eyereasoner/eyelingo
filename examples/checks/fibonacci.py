# Independent Python checks for the fibonacci example.
from .common import run_fragment_checks

CHECKS = [
    ('base cases F(0)=0 and F(1)=1 hold', ['The Fibonacci sequence is defined by F(0)=0, F(1)=1,']),
    ('recurrence holds for all computed steps', ['and F(n)=F(n-1)+F(n-2) for n>=2.']),
    ('all requested Fibonacci numbers match expected values:', ['The Fibonacci number for index 10000 is:']),
    ('F(0) = 0', ['336447648764317832666216120051075433103021484606800639065647699746800814421666623681555955136337340255820653326808361593737347904838652682630408924630564318873545443695598274916066']),
    ('F(1) = 1', ['indices as large as 10000.']),
    ('F(10) = 55', ['Arbitrary‑precision arithmetic (math/big) is used to']),
    ('F(100) = 354224848179261915075', ['compute the exact value without overflow, even for']),
    ('F(1000) = 434665576869374564356885276750 ... 516003704476137795166849228875', ['## Answer']),
    ('F(10000) = 336447648764317832666216120051 ... 171121233066073310059947366875', ['## Answer']),
    ('all Fibonacci numbers are non‑negative', ['## Answer']),
    ('the sequence is non‑decreasing from F(2) onward', ['## Answer']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
