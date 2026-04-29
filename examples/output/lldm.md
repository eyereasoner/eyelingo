# LLD — Leg Length Discrepancy Measurement

## Answer
LLD Alarm = TRUE  (discrepancy dCm = -1.908234, threshold ±1.25)

- Key computed values:
  SL1 = -0.062857  SL3 = SL4 = 15.909091
  p5  = (2.2482, 8.2935)
  p6  = (53.2285, 5.0891)
  d53 = 21.548900 cm
  d64 = 23.457134 cm
  dCm = -1.908234 cm

## Reason why
The measurement points p1-p4 define two parallel lines L1 and L3/L4.
L1 is defined by p1 and p2; L3 passes through p3; L4 passes through p4.
L3 and L4 are perpendicular to L1.  Intersection points p5 (L1∩L3)
and p6 (L1∩L4) are computed analytically.  The Euclidean distances
d53 (p5–p3) and d64 (p6–p4) are then computed, and their difference
dCm = d53 − d64 is the leg-length discrepancy.  An alarm fires when
|dCm| > 1.25 cm.

## Check
- C1 OK - L1 is perpendicular to L3 and L4 (slopes product ≈ -1).
- C2 OK - p5 lies on both L1 and L3.
- C3 OK - p6 lies on both L1 and L4.
- C4 OK - squared distances are non‑negative.
- C5 OK - dCm is a finite number.

## Go audit details
- platform : go1.26.2 linux/amd64
- input points : p1(10.1, 7.8) p2(45.1, 5.6) p3(3.6, 29.8) p4(54.7, 28.5)
- cL1 (slope L1)          : -0.06285714
- cL3 (slope L3, L4)      : 15.90909091
- p5 (x, y)               : 2.24816562, 8.29354388
- p6 (x, y)               : 53.22845573, 5.08906850
- d53                     : 21.54890046 cm
- d64                     : 23.45713445 cm
- dCm (LL discrepancy)    : -1.90823398 cm
- alarm threshold         : ±1.25 cm
- LLD alarm               : true
- checks passed           : 5/5
- recommendation consistent : yes
