# Wind Turbine Envelope  

## Answer  
operating thresholds : cut-in 3.5 m/s, rated 12.0 m/s, cut-out 25.0 m/s  
rated power : 3.2 MW  
interval classifications : t1 3.0 m/s stopped 0.000 MW; t2 6.5 m/s partial 0.440 MW; t3 11.2 m/s partial 2.586 MW; t4 15.0 m/s rated 3.200 MW; t5 24.5 m/s rated 3.200 MW; t6 27.0 m/s stopped 0.000 MW  
usable intervals : 4  
total energy : 1.571 MWh  

## Reason  
Wind below cut-in and at or above cut-out is stopped for production and safety.  
Wind between cut-in and rated speed follows a cubic power curve normalized to the rated point.  
Wind between rated speed and cut-out is capped at rated power.  
Energy is accumulated by multiplying each interval power by the ten-minute interval duration.  

## Check  
C1 OK - cut-in, rated, and cut-out thresholds are strictly ordered  
C2 OK - usable intervals are exactly the samples inside the operating envelope  
C3 OK - rated intervals are speeds at or above rated and below cut-out  
C4 OK - stopped intervals are below cut-in or at/above cut-out  
C5 OK - below-rated usable speeds follow the cubic normalized power curve  
C6 OK - total interval energy is recomputed in MWh  
C7 OK - reported usable count and total energy match recomputation  
C8 OK - the answer reports every sample classification and power  
