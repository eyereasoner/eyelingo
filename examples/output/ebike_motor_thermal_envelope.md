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
C1 OK - cooling interval 0.7788007830..0.7788007831 is positive, ordered, and contractive  
C2 OK - temperature trace has 13 samples including the initial state  
C3 OK - maximum upper temperature 40.2285 C stays below hard limit 45.0 C  
C4 OK - warning recovery occurs at step 8 after 40 seconds  
C5 OK - upper envelope decreases from the first post-Turbo sample onward  
C6 OK - ride decision is ThermallySafeForThisAssistPlan  
