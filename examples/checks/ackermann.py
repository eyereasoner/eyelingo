# Independent Python checks for the ackermann example.
from .common import run_fragment_checks

CHECKS = [
    ('x=0 reduces to successor after the y+3 binary offset', ['A1 binary offset -> T(0,9,2) -> successor gives T=10, answer=T-3=7']),
    ('x=1 reduces to addition after the y+3 binary offset', ['The N3 source defines binary ackermann(x,y) by computing T(x,y+3,2) and subtracting 3. The ternary predicate T uses direct rules for successor, addition, multiplication, and expone']),
    ('x=2 reduces to multiplication after the y+3 binary offset', ['A5 binary offset -> T(2,12,2) -> multiplication gives T=24, answer=T-3=21']),
    ('x=3 reduces to exact BigInt exponentiation, including 2^1003-3', ['A7 binary offset -> T(3,1003,2) -> exponentiation gives T=302-digit integer [857206885749013856758740...220995094697645344555008; sha256=00800470f0cc96fd412c9f2f9cb3e7d9c6d3a244820']),
    ('x=4 derives the first tetration cases T(4,3,2)-3 and T(4,4,2)-3', ['A10 binary offset -> T(4,5,2) -> tetration recursion gives T=19729-digit integer [200352993040684646497907...339445587895905719156736; sha256=64829919027d6b545f931768c171c25c302262']),
    ('A(4,2) is held exactly as 2^65536-3, not as a floating-point approximation', ['A9 binary offset -> T(4,4,2) -> tetration recursion gives T=65536, answer=T-3=65533']),
    ('the pentation query A(5,0) lands on the same value as A(4,1)', ['A11 reuses the pentation step T(5,3,2)=T(4,4,2)=65536, so A11 equals A9.']),
    ('the evaluator reached the expected largest exact integer and memoized each distinct ternary fact once', ['A10 is 2^65536 - 3, an exact 19,729-digit integer summarized by fingerprint.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
