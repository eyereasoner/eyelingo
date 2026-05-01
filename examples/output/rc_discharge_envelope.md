# RC Discharge Envelope  

## Answer  
exact decay symbol : exp(-1/4)  
certified decay interval : [0.7788007830, 0.7788007831]  
first below tolerance step : 13  
first below tolerance time : 0.325 s  
upper voltage at step 13 : 0.930581 V  

## Reason  
The physical decay factor is exp(-1/4), but the example uses a finite double interval as the certificate.  
Because the interval lies strictly between 0 and 1, the capacitor voltage envelope contracts each sample.  
The upper envelope is the safety-relevant bound: once it falls below 1.0 V, every compatible exact trajectory is below tolerance.  
The first such witness occurs before the configured maximum step.  

## Check  
C1 OK - the decay interval is nonempty, positive, and below one  
C2 OK - the certified upper voltage envelope decreases at every sample  
C3 OK - the lower and upper envelopes bracket every compatible decay  
C4 OK - step 12 remains above the voltage tolerance  
C5 OK - step 13 is the first certified below-tolerance sample  
C6 OK - reported first step, time, and upper voltage match recomputation  
C7 OK - the report uses the JSON double interval rather than only the exact symbol  
