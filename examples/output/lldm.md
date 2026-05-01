# LLD — Leg Length Discrepancy Measurement  

## Answer  
LLD Alarm = TRUE  (discrepancy dCm = -1.908234, threshold ±1.25)  

Key computed values:  
  SL1 = -0.062857  SL3 = SL4 = 15.909091  
  p5  = (2.2482, 8.2935)  
  p6  = (53.2285, 5.0891)  
  d53 = 21.548900 cm  
  d64 = 23.457134 cm  
  dCm = -1.908234 cm  

## Reason  
The measurement points p1-p4 define two parallel lines L1 and L3/L4.  
L1 is defined by p1 and p2; L3 passes through p3; L4 passes through p4.  
L3 and L4 are perpendicular to L1.  Intersection points p5 (L1∩L3)  
and p6 (L1∩L4) are computed analytically.  The Euclidean distances  
d53 (p5–p3) and d64 (p6–p4) are then computed, and their difference  
dCm = d53 − d64 is the leg-length discrepancy.  An alarm fires when  
|dCm| > 1.25 cm.  

## Check  
C1 OK - L1 slope is recomputed from p1 and p2  
C2 OK - L3/L4 slopes are perpendicular to L1  
C3 OK - p5 is the analytic intersection of L1 and the perpendicular through p3  
C4 OK - p6 is the analytic intersection of L1 and the perpendicular through p4  
C5 OK - d53 recomputes as the Euclidean distance from p5 to p3  
C6 OK - d64 recomputes as the Euclidean distance from p6 to p4  
C7 OK - dCm recomputes as d53 minus d64  
C8 OK - the discrepancy is finite and negative for this geometry  
C9 OK - the alarm condition follows from |dCm| > 1.25 cm  
