# HarborSMR Insight Dispatch  

## Answer  
PERMIT - North Quay Hydrogen Hub may use https://example.org/insight/harborsmr to run PEM_electrolyzer_train_2 at 16 MW from 2026-06-18T14:00:00Z to 2026-06-18T18:00:00Z.  

## Reason  
The SMR operator exposes a bounded 18 MW flexible-export insight for day_ahead_balancing, not raw reactor telemetry.  
The requested 16 MW electrolysis dispatch fits inside that window, safety margins clear the thresholds, no outage is planned, and the policy permits use only for electrolysis_dispatch while forbidding market resale distribution.  
The approved dispatch is 64 MWh over the four-hour window, scoped to port-hydrogen-hub and PEM_electrolyzer_train_2.  

## Check  
C1 OK - reserve margin exceeds the configured threshold  
C2 OK - cooling margin exceeds the configured threshold  
C3 OK - no planned outage blocks the balancing window  
C4 OK - requested dispatch fits inside the flexible-export insight  
C5 OK - serialized insight omits sensitive reactor telemetry terms  
C6 OK - aggregate flags keep raw reactor telemetry local  
C7 OK - permission policy authorizes electrolysis dispatch before expiry  
C8 OK - market-resale redistribution is separately prohibited  
C9 OK - scope is explicit for device, event, start, and expiry  
C10 OK - dispatch energy recomputes to 64 MWh over the four-hour window  
C11 OK - the reported load, power, and window match the request and dispatch  
