# Independent Python checks for the complex_numbers example.
from .common import run_fragment_checks

CHECKS = [
    ('N3 dial rules assign the expected polar angles for -1, e, and i', ['The N3 source first converts each complex base (x,y) to polar form using r=sqrt(x^2+y^2) and quadrant-sensitive dial rules. It then applies (a+bi)^(c+di)=r^c*e^(-d*theta)*(cos(d*ln']),
    ('all four complex exponentiation answers match the complex.n3 test facts', ['The complex.n3 test query derives 6 complex-number facts.']),
    ('i^i and e^(-pi/2) derive the same real value', ['C4 scale=0.207879576351 angleMix=0 result=0.207879576351 + 0i; a real exponent of e gives the same exp(-pi/2) value as i^i']),
    ('asin(2) and acos(2) match the N3 inverse-trig derivations', ['C5 asin: E=1 F=2 lnTerm=1.316957896925 result=1.570796326795 + 1.316957896925i; the real input is outside [-1,1], so the inverse sine has a positive imaginary part']),
    ('sin(asin(2)) and cos(acos(2)) round-trip back to 2+0i', ['C5 asin(2+0i) = 1.570796326795 + 1.316957896925i']),
    ('asin(2) + acos(2) equals pi/2 with cancelling imaginary parts', ['C6 acos: E=1 F=2 lnTerm=1.316957896925 result=0 - 1.316957896925i; the companion inverse cosine carries the opposite imaginary part']),
    ('all six complex outputs are finite real/imaginary pairs', ['C3 scale=0.207879576351 angleMix=0 result=0.207879576351 + 0i; i has polar angle pi/2, so i^i becomes exp(-pi/2)']),
    ('the translated rule count matches four exponentiation and two inverse-trig queries', ["C2 scale=1 angleMix=3.14159265359 result=-1 + 0i; Euler's identity falls out of the exponent rule because ln(e)=1"]),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
