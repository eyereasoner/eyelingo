# Decimal Servo Envelope  

## Answer  
exact pole symbol : exp(-1/3)  
certified pole interval : [0.7165313105, 0.7165313106]  
first settled step : 10  
first settled time : 0.200 s  
upper envelope at step 10 : 0.445925  

## Reason why  
The exact pole exp(-1/3) is not represented as an exact finite decimal, so the input provides a certified decimal interval.  
The upper bound of the pole interval is below 1, which makes the error envelope contractive.  
The certificate propagates the lower and upper absolute-error envelopes sample by sample.  
The first sample whose upper envelope is below the tolerance is the first guaranteed settling witness.  

## Check  
C1 OK - pole certificate is nonempty, positive, and below 1  
C2 OK - upper envelope strictly decreases at every sampled step  
C3 OK - step 9 is not yet below tolerance  
C4 OK - step 10 is the first certified settling step  
C5 OK - settling time is 0.200 s  
C6 OK - all values are derived from the JSON certificate parameters  

## Go audit details  
platform : go1.26.2 linux/amd64  
case : decimal_servo_envelope  
question : When is a sampled servo guaranteed inside the error tolerance using a decimal pole certificate?  
sample period : 0.020 s  
time constant : 0.060 s  
max step : 12  
initial absolute error : 12.500  
tolerance : 0.500  
envelope rows : 13  
checks passed : 6/6  
