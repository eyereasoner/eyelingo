# Euler Identity Certificate  

## Answer  
expression : exp(iπ) + 1  
terms used : 28  
computed real part of exp(iπ) : -1.000000000000000  
computed imaginary part of exp(iπ) : -0.000000000000000  
residual magnitude : 2.628e-16  
within tolerance : true  

## Reason why  
The example approximates exp(iπ) by a finite Taylor series over complex numbers.  
The resulting residual is not claimed to be mathematically exact zero; it is checked against the explicit tolerance from JSON.  
The computed real part is effectively -1 and the imaginary part is near 0 at the chosen precision.  
That gives a reproducible finite certificate for the familiar Euler-identity witness.  

## Check  
C1 OK - the Taylor expansion uses the terms count from JSON  
C2 OK - the input angle is pi to floating precision  
C3 OK - the finite Taylor real part is close to -1  
C4 OK - the finite Taylor imaginary part is close to zero  
C5 OK - the residual |exp(iπ)+1| is below the configured tolerance  
C6 OK - reported real, imaginary, and residual values match recomputation  
C7 OK - the explanation explicitly treats the result as a finite certificate  
