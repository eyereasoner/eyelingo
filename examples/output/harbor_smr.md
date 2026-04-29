# HarborSMR Insight Dispatch

## Answer
PERMIT - North Quay Hydrogen Hub may use https://example.org/insight/harborsmr to run PEM_electrolyzer_train_2 at 16 MW from 2026-06-18T14:00:00Z to 2026-06-18T18:00:00Z.

## Reason why
The SMR operator exposes a bounded 18 MW flexible-export insight for day_ahead_balancing, not raw reactor telemetry.
The requested 16 MW electrolysis dispatch fits inside that window, safety margins clear the thresholds, no outage is planned, and the policy permits use only for electrolysis_dispatch while forbidding market resale distribution.
The approved dispatch is 64 MWh over the four-hour window, scoped to port-hydrogen-hub and PEM_electrolyzer_train_2.

## Check
C1 OK - reserve margin 24 MW exceeds threshold 19 MW
C2 OK - cooling margin 18% exceeds threshold 14%
C3 OK - no planned outage blocks the balancing window
C4 OK - requested 16 MW fits inside the 18 MW flexible-export insight
C5 OK - serialized insight omits sensitive telemetry terms
C6 OK - aggregate flags confirm raw reactor telemetry stays local
C7 OK - policy permits use for electrolysis dispatch before the insight expires
C8 OK - policy prohibits redistribution for market resale
C9 OK - scope is explicit: device, event, start, and expiry
C10 OK - dispatch plan matches the requested load, power, and insight window

## Go audit details
platform : go1.26.2 linux/amd64
case : harborsmr
question : May the port hydrogen hub use a temporary SMR flexible-export insight to schedule electrolysis for the next four-hour balancing window?
hub : portHydrogenHub (North Quay Hydrogen Hub, BE)
aggregate : available=18MW reserve=24MW cooling=18% outage=false
insight : id=https://example.org/insight/harborsmr metric=temporary_flexible_export_window export=18MW target=PEM_electrolyzer_train_2
scope : device=port-hydrogen-hub event=day_ahead_balancing start=2026-06-18T14:00:00Z expires=2026-06-18T18:00:00Z
request : action=odrl:use purpose=electrolysis_dispatch requested=16MW target=PEM_electrolyzer_train_2 authAt=2026-06-18T13:10:00Z
dispatch : 16MW 2026-06-18T14:00:00Z to 2026-06-18T18:00:00Z load=PEM_electrolyzer_train_2 energy=64MWh
policy : profile=HarborSMR-Insight-Policy permission=(odrl:use target=https://example.org/insight/harborsmr purpose=electrolysis_dispatch) prohibition=(odrl:distribute target=https://example.org/insight/harborsmr purpose=market_resale)
privacy flags : coreTemperature=false rodPosition=false neutronFlux=false operatorBadgeIDs=false
serialized sensitive term hits : none
checks passed : 10/10
decision : PERMIT
