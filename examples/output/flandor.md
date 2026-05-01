# Flandor  

## Answer  
Name: Flandor  
Region: Flanders  
Metric: regional_retooling_priority  
Active need count: 3/3  
Recommended package: Flandor Retooling Pulse  
Budget cap: €140M  
Package cost: €120M  
Worker coverage: 1200  
Grid relief: 85 MW  
Payload SHA-256: 718f5b17d07ab6a95503bc04a1000ddb132409f600659c03d21def81914b780b  
Envelope HMAC-SHA-256: 955968ca99a191783bc00cba068128ccb9ff40a5e6114fda13a52c74ee27329e  
decision : ALLOWED  

## Reason why  
question : Is the Flemish Economic Resilience Board allowed to use a neutral macro-economic insight for regional stabilization, and if so which package should it activate for Flanders?  
aggregate : observedFirms=217 level=regional_cluster containsFirmNames=no containsPayrollRows=no  

export weakness : active - at least one cluster has exportOrdersIndex < 90  
skills strain : active - technical vacancy rate is above 3.9%  
grid stress : active - grid congestion hours are above 11  
export weak clusters : Antwerp chemicals=84, Ghent manufacturing=87  
technical vacancy rate : 4.6% (threshold > 3.9%)  
grid congestion : 19 hours (threshold > 11)  

recommendation policy : lowest_cost_package_covering_all_active_needs  
candidate packages:  
  pkg:TRAIN_070   : reject, cost=€70M, covers 1/3 active needs; within budget  
  pkg:PORT_095    : reject, cost=€95M, covers 1/3 active needs; within budget  
  pkg:RET_FLEX_120 : selected, cost=€120M, covers all active needs; within budget  
  pkg:CORRIDOR_165 : reject, cost=€165M, covers all active needs; over budget by €25M  

Selected package "Flandor Retooling Pulse" covers export=yes, skills=yes, grid=yes, cost=€120M.  
Usage is permitted only for purpose "regional_stabilization" and the envelope expires at 2026-04-08T19:00:00+00:00.  
Surveillance reuse is blocked by a prohibition on odrl:distribute for purpose "firm_surveillance".  
Deletion duty time : 2026-04-08T18:30:00+00:00  

## Check  
C1 OK - export weakness, skills strain, and grid stress are all active  
C2 OK - active need count meets the insight threshold  
C3 OK - the lowest-cost package covering all active needs is selected  
C4 OK - the selected package fits inside the €140M budget  
C5 OK - cheaper packages are rejected because each covers only one active need  
C6 OK - the full corridor package covers all needs but is over budget  
C7 OK - policy permission authorizes regional-stabilization use before expiry  
C8 OK - firm-surveillance redistribution is prohibited  
C9 OK - deletion duty is scheduled before envelope expiry  
C10 OK - shared insight omits firm names and payroll rows  
C11 OK - reported signature metadata matches the trusted precomputed input  
C12 OK - the expected six files and one audit entry are recorded  
C13 OK - export-weak cluster names are independently identified  
