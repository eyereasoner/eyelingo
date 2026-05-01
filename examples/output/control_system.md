# Control System — ARC-style control-system example  

## Answer  
actuator1 control1 = 39.273462  
actuator2 control1 = 26.080000  
input1 measurement10 = 2.236068 (lessThan branch)  
disturbance2 measurement10 = 45.000000 (notLessThan branch)  

## Reason why  
The backward measurement10 rules first normalize measurement1 pairs into scalar measurement10 values.  
For input1, 6 < 11, so sqrt(11 - 6) = 2.236068.  
The first forward rule computes feedforward control from input1 measurement10, the input2 boolean guard, and disturbance1 compensation.  
The second forward rule computes PND feedback control from target/measurement error and the state/output differential error.  

## Check  
C1 OK - input1 measurement10 follows the lessThan backward rule  
C2 OK - disturbance2 measurement10 follows the notLessThan backward rule  
C3 OK - input2 boolean guard is true for the feedforward rule  
C4 OK - feedforward control equals product minus log10 compensation  
C5 OK - feedback error is target minus measurement  
C6 OK - differential error is state observation minus output measurement  
C7 OK - nonlinear differential part equals (7.3 / error) times differential error  
C8 OK - actuator2 control is proportional part plus nonlinear differential part  
