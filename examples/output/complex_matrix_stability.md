# Complex Matrix Stability  

## Answer  
A_unstable : spectral radius 2 -> unstable  
A_stable : spectral radius 1 -> marginally stable  
A_damped : spectral radius 0 -> damped  

## Reason why  
For a discrete-time linear system x_{k+1} = A x_k, the eigenvalues of A govern the modal behaviour.  
Because the matrices are diagonal, the eigenvalues are the diagonal entries; the largest modulus gives the spectral radius and therefore the stability class.  

## Check  
C1 OK - A_unstable has spectral radius 2, so it is unstable  
C2 OK - A_stable has spectral radius 1, so it is marginally stable  
C3 OK - A_damped has spectral radius 0, so every mode decays  
C4 OK - squared modulus of z*w equals the product of squared moduli  
C5 OK - scaling A_unstable by 2 multiplies spectral-radius-squared by 4  

## Go audit details  
platform : go1.26.2 linux/amd64  
question : Classify three diagonal 2x2 complex matrices for discrete-time stability.  
translated source : complex-matrix-stability.n3  
matrices checked : 3  
scale factor : 2  
scaled unstable radius squared : 16  
checks passed : 5/5  
