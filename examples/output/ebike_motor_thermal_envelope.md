# E-Bike Motor Thermal Envelope  

## Answer  
decision : ThermallySafeForThisAssistPlan  
cooling certificate : exp(-1/4) in 0.7788007830 .. 0.7788007831  
maximum upper motor temperature : 40.2285 C  
warning recovery : step 8 at 40 s  
hard limit : 45.0 C  

## Reason why  
The model keeps an interval for motor temperature excess over ambient instead of pretending to know the transcendental cooling factor exactly.  
For each 5 second sample, the lower and upper excess envelopes are propagated with the certified cooling interval and the mode-specific heat injection.  
Turbo pushes the upper envelope to 40.2285 C, then Tour, Eco, and Coast allow the envelope to decrease.  
The upper envelope returns below the 35.0 C warning limit at step 8 and remains below the 45.0 C hard limit throughout.  

## Check  
C1 OK - the cooling certificate brackets exp(-sample/tau)  
C2 OK - the assist plan has twelve sampled thermal updates  
C3 OK - Turbo, Tour, Eco, and Coast heating envelopes are nonnegative intervals  
C4 OK - interval propagation recomputes the maximum upper motor temperature  
C5 OK - the upper envelope first exceeds the warning limit during Turbo  
C6 OK - the reported warning-recovery step matches the independently propagated envelope  
C7 OK - all upper temperatures remain below the hard thermal limit  
C8 OK - the final Coast samples cool monotonically  
C9 OK - the final decision matches the safety envelope  
