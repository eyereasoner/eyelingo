# GPS — Goal driven route planning  

## Answer  
Take the direct route via Brugge.  
Recommended route: Gent → Brugge → Oostende  

## Reason  
From Gent to Oostende, the planner found two routes in this small map.  
The direct route (Gent → Brugge → Oostende) takes 2400.0 seconds at cost 0.01, with belief 0.9408 and comfort 0.99. The alternative (Gent → Kortrijk → Brugge → Oostende) takes 4100.0 seconds at cost 0.018, with belief 0.903168 and comfort 0.9801.  
So the direct route is faster, cheaper, more reliable, and slightly more comfortable.  

## Check  
C1 OK - the direct Gent-Brugge-Oostende route is derived from edges  
C2 OK - the Kortrijk alternative route is derived from edges  
C3 OK - exactly two simple routes connect the traveller to the destination  
C4 OK - route duration and cost are additive over edges  
C5 OK - route belief and comfort are multiplicative over edges  
C6 OK - the recommended route is faster than the alternative  
C7 OK - the recommended route is cheaper than the alternative  
C8 OK - the recommended route has higher belief and comfort scores  
C9 OK - the answer names the independently chosen direct route  
