# Control System — ARC-style control-system example  

## Answer  
actuator1 control1 = 39.273462  
actuator2 control1 = 26.080000  
input1 measurement10 = 2.236068 (lessThan branch)  
disturbance2 measurement10 = 45.000000 (notLessThan branch)  

## Reason  
The backward measurement10 rules first normalize measurement1 pairs into scalar measurement10 values.  
For input1, 6 < 11, so sqrt(11 - 6) = 2.236068.  
The first forward rule computes feedforward control from input1 measurement10, the input2 boolean guard, and disturbance1 compensation.  
The second forward rule computes PND feedback control from target/measurement error and the state/output differential error.  

## Check  
C1 OK - input1 measurement10 recomputes through the lessThan square-root branch  
C2 OK - disturbance2 measurement10 recomputes through the notLessThan branch  
C3 OK - feedforward guard is true before actuator1 arithmetic is applied  
C4 OK - actuator1 control recomputes product minus log10 compensation  
C5 OK - target-minus-measurement error is recomputed as 5  
C6 OK - state/output differential error is recomputed as -2  
C7 OK - nonlinear feedback term uses 7.3/error times the differential  
C8 OK - actuator2 control recomputes proportional plus nonlinear feedback  
