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

## Go audit details  
platform : go1.26.2 linux/amd64  
case : ebike_motor_thermal_envelope  
question : Does the sampled e-bike climb remain thermally safe when exp(-1/4) is represented by a certified decimal interval?  
sample period : 5 s  
thermal time constant : 20 s  
assist plan length : 12  
trace step 00 : mode=Initial tempLower=33.0000C tempUpper=33.0000C  
trace step 01 : mode=Turbo tempLower=36.0304C tempUpper=36.0304C  
trace step 02 : mode=Turbo tempLower=38.3905C tempUpper=38.3905C  
trace step 03 : mode=Turbo tempLower=40.2285C tempUpper=40.2285C  
trace step 04 : mode=Tour tempLower=39.2600C tempUpper=39.2600C  
trace step 05 : mode=Tour tempLower=38.5057C tempUpper=38.5057C  
trace step 06 : mode=Eco tempLower=36.7182C tempUpper=36.7182C  
trace step 07 : mode=Eco tempLower=35.3262C tempUpper=35.3262C  
trace step 08 : mode=Eco tempLower=34.2420C tempUpper=34.2420C  
trace step 09 : mode=Coast tempLower=32.1977C tempUpper=32.1977C  
trace step 10 : mode=Coast tempLower=30.6056C tempUpper=30.6056C  
trace step 11 : mode=Coast tempLower=29.3656C tempUpper=29.3656C  
trace step 12 : mode=Coast tempLower=28.4000C tempUpper=28.4000C  
decreasing steps : [3 4 5 6 7 8 9 10 11]  
checks passed : 6/6  
