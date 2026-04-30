# Wind Turbine Envelope  

## Answer  
operating thresholds : cut-in 3.5 m/s, rated 12.0 m/s, cut-out 25.0 m/s  
rated power : 3.2 MW  
interval classifications : t1 3.0 m/s stopped 0.000 MW; t2 6.5 m/s partial 0.440 MW; t3 11.2 m/s partial 2.586 MW; t4 15.0 m/s rated 3.200 MW; t5 24.5 m/s rated 3.200 MW; t6 27.0 m/s stopped 0.000 MW  
usable intervals : 4  
total energy : 1.571 MWh  

## Reason why  
Wind below cut-in and at or above cut-out is stopped for production and safety.  
Wind between cut-in and rated speed follows a cubic power curve normalized to the rated point.  
Wind between rated speed and cut-out is capped at rated power.  
Energy is accumulated by multiplying each interval power by the ten-minute interval duration.  

## Check  
C1 OK - cut-in, rated, and cut-out thresholds are ordered  
C2 OK - usable intervals are exactly the samples inside the operating envelope  
C3 OK - two intervals reach rated power  
C4 OK - two intervals are stopped by low or cut-out wind  
C5 OK - a below-rated usable wind speed follows the cubic power curve  
C6 OK - the six ten-minute samples yield about 1.57 MWh  

## Go audit details  
platform : go1.26.2 linux/amd64  
case : wind-turbine  
question : Classify wind-speed samples against a turbine envelope and compute interval energy.  
samples evaluated : 6  
interval minutes : 10.0  
partial intervals : 2  
rated intervals : 2  
stopped intervals : 2  
checks passed : 6/6  
