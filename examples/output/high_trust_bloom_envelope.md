# High Trust RDF Bloom Envelope

## Answer
- Deployment decision : AcceptForHighTrustUse for artifact.
- lambda : 0.5126953125
- false-positive envelope : 0.001670806 .. 0.001670806
- expected extra exact lookups upper : 83.540 per 50000 negative lookups
- maybe-positive policy : ConfirmAgainstCanonicalGraph
- definite-negative policy : ReturnAbsent

## Reason why
The canonical graph and the SPO index have the same triple count, so exact membership remains grounded in the graph snapshot.
The Bloom prefilter has n=1200 triples, m=16384 bits, and k=7 hash functions, giving lambda 0.5126953125.
Instead of asking the engine to know exp(-k*n/m) exactly, the input carries a decimal interval certificate for exp(-lambda).
That certificate bounds (1-exp(-lambda))^k below the 0.002 false-positive budget and keeps extra exact confirmations below 100.0.
- Correctness never depends on the Bloom filter alone: maybe-positive answers must be confirmed against the canonical graph.

## Check
- C1 OK - numeric Bloom and workload parameters are positive
- C2 OK - canonical graph and SPO index agree on 1200 triples
- C3 OK - derived lambda 0.5126953125 matches the certified lambda
- C4 OK - decimal interval 0.5988792348..0.5988792349 is a valid exp(-lambda) certificate
- C5 OK - false-positive upper bound 0.001670806 is below 0.002
- C6 OK - expected extra exact lookups 83.540 stay below 100.0
- C7 OK - maybe-positive Bloom hits are confirmed against the canonical graph
- C8 OK - definite Bloom negatives may be returned absent without exact lookup
- C9 OK - deployment decision is AcceptForHighTrustUse

## Go audit details
- platform : go1.26.2 linux/amd64
- case : high_trust_bloom_envelope
- question : Can a high-trust RDF graph artifact use a Bloom prefilter while keeping correctness grounded in exact canonical graph checks?
- artifact : artifact
- triple counts : canonical=1200 spoIndex=1200 agreement=true
- bloom parameters : bits=16384 hashFunctions=7 lambda=0.5126953125
- certificate : symbol=exp(-k*n/m) certifiedLambda=0.5126953125 expLower=0.5988792348 expUpper=0.5988792349 certified=true
- fp envelope : lower=0.001670806 upper=0.001670806 budget=0.002 within=true
- lookup budget : negativeLookups=50000 extraUpper=83.540 budget=100.0 within=true
- policies : maybePositive=ConfirmAgainstCanonicalGraph definiteNegative=ReturnAbsent
- checks passed : 9/9
- decision : AcceptForHighTrustUse
