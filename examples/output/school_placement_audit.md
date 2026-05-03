# School Placement Route Audit  

## Answer  
audit result : fail  
children affected by straight-line rule : Ada, Björn, Davi  
largest hidden detour : Ada, 3000 m  
recommended assignments : Ada -> Lindholmen; Björn -> Backa; Clara -> Haga; Davi -> Haga  
explanation requested : yes  

## Reason  
The support-tool rule chooses the school with the smallest straight-line distance, using preference rank only as a tie-breaker.  
The independent audit recomputes each candidate with walking-route distance plus 600 m per preference step.  
Any provisional assignment that is not the audited best, or that requires more than 2500 m of walking, is flagged.  
Ada and Björn look close to Centrum on a map, but their walking routes cross barriers and exceed the walking limit; Davi is also better served by the first-preference Haga route.  
This illustrates why a decision-support label is not enough: route geometry, preferences, and audit records must be inspectable.  

## Check  
C1 OK - the fixture has four students, four schools, and a complete 4 × 4 distance matrix  
C2 OK - every student preference list covers every school exactly once  
C3 OK - the independent straight-line rule assigns Ada and Björn to Centrum  
C4 OK - the independent route-aware rule computes walking distance plus preference penalty  
C5 OK - the reported recommended assignments match the Go audit  
C6 OK - the reported affected children are exactly those whose provisional placement is flagged  
C7 OK - the reported largest hidden detour matches the flawed straight-line placement  
C8 OK - the failure result follows from at least one over-limit walking route and changed assignment  
C9 OK - the Reason text names the support-tool rule, walking-route recomputation, and inspectability requirement  
C10 OK - the report explicitly requests an explanation for the affected placement decisions  
