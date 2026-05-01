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
C1 OK - Taylor expansion used 28 terms from the JSON input  
C2 OK - computed real part is close to -1  
C3 OK - computed imaginary part is close to 0  
C4 OK - |exp(iπ)+1| is below the configured tolerance  
C5 OK - the audit records the finite residual rather than asserting exact real arithmetic  
