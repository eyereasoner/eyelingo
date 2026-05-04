# Euler Identity Certificate  

## Insight  
expression : exp(iπ) + 1  
terms used : 28  
computed real part of exp(iπ) : -1.000000000000000  
computed imaginary part of exp(iπ) : 0.000000000000000  
residual magnitude : 3.767e-16  
within tolerance : true  

## Explanation  
The example approximates exp(iπ) by a finite Taylor series over complex numbers.  
The resulting residual is not claimed to be mathematically exact zero; it is checked against the explicit tolerance from JSON.  
The computed real part is effectively -1 and the imaginary part is near 0 at the chosen precision.  
That gives a reproducible finite certificate for the familiar Euler-identity witness.  
