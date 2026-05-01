# High Trust RDF Bloom Envelope  

## Answer  
Deployment decision : AcceptForHighTrustUse for artifact.  
lambda : 0.5126953125  
false-positive envelope : 0.001670806 .. 0.001670806  
expected extra exact lookups upper : 83.540 per 50000 negative lookups  
maybe-positive policy : ConfirmAgainstCanonicalGraph  
definite-negative policy : ReturnAbsent  

## Reason why  
The canonical graph and the SPO index have the same triple count, so exact membership remains grounded in the graph snapshot.  
The Bloom prefilter has n=1200 triples, m=16384 bits, and k=7 hash functions, giving lambda 0.5126953125.  
Instead of asking the engine to know exp(-k*n/m) exactly, the input carries a decimal interval certificate for exp(-lambda).  
That certificate bounds (1-exp(-lambda))^k below the 0.002 false-positive budget and keeps extra exact confirmations below 100.0.  
Correctness never depends on the Bloom filter alone: maybe-positive answers must be confirmed against the canonical graph.

## Check  
C1 OK - numeric Bloom parameters are positive  
C2 OK - canonical graph and SPO index triple counts agree  
C3 OK - lambda recomputes as k*n/m  
C4 OK - the exp(-lambda) decimal certificate brackets the exact Python value  
C5 OK - false-positive envelope is recomputed from the certificate  
C6 OK - false-positive upper bound stays below the policy budget  
C7 OK - expected extra exact lookups stay below budget  
C8 OK - maybe-positive answers must be confirmed against the canonical graph  
C9 OK - definite negatives may return absent without exact lookup  
C10 OK - the deployment decision matches the recomputed envelope
