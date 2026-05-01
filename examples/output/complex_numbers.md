# Complex Numbers  

## Answer  
The complex.n3 test query derives 6 complex-number facts.  
Computed values:  
C1 sqrt(-1+0i) = 0 + 1i  
C2 e^(i*pi) = -1 + 0i  
C3 i^i = 0.207879576351 + 0i  
C4 e^(-pi/2) = 0.207879576351 + 0i  
C5 asin(2+0i) = 1.570796326795 + 1.316957896925i  
C6 acos(2+0i) = 0 - 1.316957896925i  
Key equivalences:  
i^i = e^(-pi/2) = 0.207879576351  
asin(2) + acos(2) = 1.570796326795 + 0i  

## Reason why  
The N3 source first converts each complex base (x,y) to polar form using r=sqrt(x^2+y^2) and quadrant-sensitive dial rules. It then applies (a+bi)^(c+di)=r^c*e^(-d*theta)*(cos(d*ln(r)+c*theta), sin(d*ln(r)+c*theta)). The asin/acos rules use the same pair of distances from 1+a and 1-a, then recover the imaginary part with ln(F+sqrt(F^2-1)).  
primitive test facts : 6  
polar derivations : 4  
rule applications : 10  
polar bases:  
C1 base=-1 + 0i r=1 theta=3.14159265359 quadrant=II or -real axis  
C2 base=2.718281828459 + 0i r=2.718281828459 theta=0 quadrant=I or +axis  
C3 base=0 + 1i r=1 theta=1.570796326795 quadrant=I or +axis  
C4 base=2.718281828459 + 0i r=2.718281828459 theta=0 quadrant=I or +axis  
exponentiation traces:  
C1 scale=1 angleMix=1.570796326795 result=0 + 1i; the polar angle of -1+0i is pi, so the half power rotates to pi/2  
C2 scale=1 angleMix=3.14159265359 result=-1 + 0i; Euler's identity falls out of the exponent rule because ln(e)=1  
C3 scale=0.207879576351 angleMix=0 result=0.207879576351 + 0i; i has polar angle pi/2, so i^i becomes exp(-pi/2)  
C4 scale=0.207879576351 angleMix=0 result=0.207879576351 + 0i; a real exponent of e gives the same exp(-pi/2) value as i^i  
inverse-trig traces:  
C5 asin: E=1 F=2 lnTerm=1.316957896925 result=1.570796326795 + 1.316957896925i; the real input is outside [-1,1], so the inverse sine has a positive imaginary part  
C6 acos: E=1 F=2 lnTerm=1.316957896925 result=0 - 1.316957896925i; the companion inverse cosine carries the opposite imaginary part  

## Check  
C1 OK - principal polar angles for -1, e, and i match the N3 dial cases  
C2 OK - all four complex exponentiation answers match independent complex arithmetic  
C3 OK - i^i and e^(-pi/2) recompute to the same real value  
C4 OK - asin(2) and acos(2) match independent inverse-trig recomputation  
C5 OK - sin(asin(2)) and cos(acos(2)) round-trip to 2+0i  
C6 OK - asin(2) + acos(2) equals pi/2 with cancelling imaginary parts  
C7 OK - all six reported complex outputs match recomputation to displayed precision  
C8 OK - the report contains four exponentiation and two inverse-trig queries  
C9 OK - all recomputed outputs are finite real/imaginary pairs  
