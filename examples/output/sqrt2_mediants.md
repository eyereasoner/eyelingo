# Integer-First Sqrt2 Mediants  

## Answer  
lower bound : 1393/985 = 1.414213197970  
upper bound : 577/408 = 1.414215686275  
certified interval width : 0.000002488305  
convergents used : 1/1, 3/2, 7/5, 17/12, 41/29, 99/70, 239/169, 577/408, 1393/985  

## Reason why  
The continued fraction for sqrt(2) is [1; 2, 2, 2, ...], so each convergent is derived by an integer recurrence.  
Each rational p/q is classified without floating-point square roots by comparing p*p with 2*q*q.  
The latest lower and upper convergents within the denominator limit form a tight certified bracket.  
Their cross-difference is one, so no simpler Stern-Brocot rational lies strictly between them.  

## Check  
C1 OK - nine convergents stay within the denominator limit  
C2 OK - the best lower bound is 1393/985  
C3 OK - the best upper bound is 577/408  
C4 OK - the bracket is certified by integer square comparisons  
C5 OK - the lower rational is strictly below the upper rational  
C6 OK - the chosen bounds are adjacent Stern-Brocot neighbors  

## Go audit details  
platform : go1.26.2 linux/amd64  
case : integer-first-sqrt2-mediants  
question : Find an integer-certified rational bracket for sqrt(2) under a denominator limit.  
max denominator : 1000  
convergents generated : 9  
lower square comparison : 1393^2 < 2*985^2  
upper square comparison : 577^2 > 2*408^2  
checks passed : 6/6  
