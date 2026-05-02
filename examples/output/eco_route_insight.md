# Eco Route Insight  

## Answer  
insight status : issue  
show eco banner : yes  
audience : Depot X  
allowed use : ui.eco.banner  
suggested route : alt-low-fuel  
current fuel index : 120.75  
suggested fuel index : 99.75  
estimated saving : 21.00  
expires at : 2025-01-01T11:00:00Z  
raw data exported : no  
signature algorithm : HMAC-SHA256  
payload digest : 00e19becd91e81d6881749655d23d43002d9ea714bba61e855eafbc8ef9a5135  
signature key : local-demo-key  
signature : 7fFGBN8fyI7xrmRz5VreeAUSf3LC_ywbj32NGk2ovUs  

## Reason  
The current route uses fuel index = distanceKm × (payloadKg / 1000) × gradientFactor.  
For shipment-1, Current urban route gives 42.00 × 2.50 × 1.15 = 120.75.  
The policy threshold is 120.00, so a local eco banner is justified.  
The selected alternative alt-low-fuel gives 38.00 × 2.50 × 1.05 = 99.75, saving 21.00 while staying within the ETA delay limit.  
The signed envelope exposes audience, use, expiry, route suggestion, and compact fuel indices, but not raw payload, GPS trace, driver behavior, or raw telemetry.  
This follows the insight pattern: ship the decision, not the data.  

## Check  
C1 OK - the current fuel index is recomputed from distance, payload tonnes, and gradient  
C2 OK - the policy threshold triggers a local eco banner  
C3 OK - the selected alternative is the best eligible lower-fuel route  
C4 OK - the alternative fuel index and saving are independently recomputed  
C5 OK - the selected route stays within the allowed ETA delay  
C6 OK - the envelope audience, allowed use, expiry, and raw-data flag match the policy  
C7 OK - the canonical envelope omits forbidden raw logistics terms  
C8 OK - the payload digest is SHA-256 over the independently rebuilt envelope  
C9 OK - the signature is the expected base64url HMAC-SHA256 value  
C10 OK - the Reason text reports the same arithmetic and the insight pattern  
