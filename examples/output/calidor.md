# Calidor  

## Answer  
Name: Calidor  
Municipality: Calidor  
Metric: active_need_count  
Active need count: 4/3  
Recommended package: Calidor Priority Cooling Bundle  
Package cost: €18  
Budget cap: €20  
Payload SHA-256: 3780df1071b0f2eec8a881ffd48425c3a1a60738d11cc2ba7debdf1cea992d63  
Envelope HMAC-SHA-256: e635c7c1991742a5c36992fc0da32a7abc80b32aa5777a1142adaab55183681c  
decision : ALLOWED  

## Reason why  
question : Is the Calidor heat-response system allowed to use a narrow household support insight for heatwave response, and if so which support package should it recommend?  
The gateway desensitizes local heat, vulnerability, and prepaid-energy stress into an expiring municipal support insight.  
metric : active_need_count  
threshold : 3.0  
scope : household-gateway @ heat-alert-window  
required capabilities: bill_credit, cooling_kit, transport, welfare_check  

heat alert : active - alert level 4 is at least 3  
unsafe indoor heat : active - 31.4°C for 9 hours reaches the 30.0°C/6 hour threshold  
vulnerability present : active - the local gateway sees heat-sensitive and mobility flags  
energy constraint : active - €3.2 prepaid credit is at or below the €5.0 limit  
vulnerability flags kept local : 2  
expires at : 2026-07-18T21:00:00+00:00  

support policy : lowest_cost_covering_package  
candidate packages:  
  pkg:CHECK    : reject, cost=€8, covers 1/4 required capabilities; within budget  
  pkg:VOUCHER  : reject, cost=€12, covers 2/4 required capabilities; within budget  
  pkg:BUNDLE   : selected, cost=€18, covers all required capabilities; within budget  
  pkg:DELUXE   : reject, cost=€28, covers all required capabilities; over budget by €8  

Selected package "Calidor Priority Cooling Bundle" covers bill_credit, cooling_kit, transport, welfare_check.  
Usage is permitted only for purpose "heatwave_response" and the envelope expires at 2026-07-18T21:00:00+00:00.  
Tenant-screening reuse is blocked by a prohibition on odrl:distribute for purpose "tenant_screening".  
reason.txt : The gateway keeps raw indoor heat, vulnerability, and prepaid-energy data local, derives a priority-support signal, and shares only a scoped heatwave-response envelope with expiry.  
dispatches logged : 1

## Check  
C1 OK - four active heat-response needs are recomputed from the local signals  
C2 OK - the insight threshold of three active needs is met  
C3 OK - the lowest-cost eligible package covering all required capabilities is selected  
C4 OK - the selected package fits inside the €20 budget  
C5 OK - cheaper packages are rejected because they do not cover all capabilities  
C6 OK - the deluxe package is rejected because it is over budget  
C7 OK - policy permission authorizes heatwave-response use before expiry  
C8 OK - tenant-screening reuse is prohibited by the policy  
C9 OK - deletion duty is scheduled before envelope expiry  
C10 OK - vulnerability flags and local raw stress signals are omitted from the serialized insight  
C11 OK - reported signature metadata matches the trusted precomputed input  
C12 OK - scope metadata is explicit for device, event, municipality, creation, and expiry
