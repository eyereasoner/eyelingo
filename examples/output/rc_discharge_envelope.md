# RC Discharge Envelope  

## Answer  
exact decay symbol : exp(-1/4)  
certified decay interval : [0.7788007830, 0.7788007831]  
first below tolerance step : 13  
first below tolerance time : 0.325 s  
upper voltage at step 13 : 0.930581 V  

## Reason why  
The physical decay factor is exp(-1/4), but the example uses a finite double interval as the certificate.  
Because the interval lies strictly between 0 and 1, the capacitor voltage envelope contracts each sample.  
The upper envelope is the safety-relevant bound: once it falls below 1.0 V, every compatible exact trajectory is below tolerance.  
The first such witness occurs before the configured maximum step.  

## Check  
C1 OK - decay certificate is nonempty, positive, and below 1  
C2 OK - voltage upper envelope decreases at every sample  
C3 OK - step 12 remains above the voltage tolerance  
C4 OK - step 13 is the first certified discharge step  
C5 OK - settling time is 0.325 s  
C6 OK - the certificate uses the JSON double interval rather than an exact transcendental value  

## Go audit details  
platform : go1.26.2 linux/amd64  
case : rc_discharge_envelope  
question : When is a sampled RC capacitor guaranteed below the voltage tolerance using a double decay certificate?  
sample period : 0.025 s  
time constant : 0.100 s  
max step : 18  
initial voltage : 24.000 V  
tolerance : 1.000 V  
envelope rows : 19  
checks passed : 6/6  
