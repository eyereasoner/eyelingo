# Complex Matrix Stability  

## Answer  
A_unstable : spectral radius 2 -> unstable  
A_stable : spectral radius 1 -> marginally stable  
A_damped : spectral radius 0 -> damped  

## Reason  
For a discrete-time linear system x_{k+1} = A x_k, the eigenvalues of A govern the modal behaviour.  
Because the matrices are diagonal, the eigenvalues are the diagonal entries; the largest modulus gives the spectral radius and therefore the stability class.  

## Check  
C1 OK - diagonal entries are used as the eigenvalues  
C2 OK - A_unstable has independently recomputed spectral radius 2  
C3 OK - A_stable has spectral radius exactly 1 and is marginal  
C4 OK - A_damped has spectral radius 0 and is damped  
C5 OK - reported matrix classes and radii match recomputation  
C6 OK - squared modulus of z*w equals product of squared moduli  
C7 OK - scaling a matrix by 2 multiplies spectral-radius-squared by 4  
