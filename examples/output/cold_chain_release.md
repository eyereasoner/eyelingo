# Cold Chain Release

## Answer
Release decision : 3 of 6 candidate lots can ship.
Safe release lots : BIO-17A, BIO-20A, BIO-21A
Quarantined lots : BIO-17B, BIO-18A, BIO-19A
Best allocation value : 7860 priority-dose points
Doses delivered : 860
Priority allocation:
 - BIO-17A -> North Clinic : 380/420 doses, transit=110m, value=3420
 - BIO-20A -> Dialysis Unit : 300/300 doses, transit=90m, value=3000
 - BIO-21A -> Field Station : 180/180 doses, transit=70m, value=1440
Reserved/unmet:
 - BIO-17A remainder held under ledger root sha256:e57b3009e0aa42a9…: 40 doses
 - Field Station unmet demand: 70 doses
 - Rural Mobile Unit unmet demand: 180 doses

## Reason why
The dataset is treated as linked facts: parent/child lot relationships create a lineage closure; recalled ancestors contaminate descendants; custody rows are verified as a SHA-256 hash chain; temperature segments are checked against the 2-8 °C policy; only releasable lots enter an exact memoized allocation search.
lineage parent facts : 8
lineage closure facts : 12 ancestor links
custody event facts : 18
temperature observations : 18
release candidates : 6
recalled ancestors : Seed-Beta
releasable lots : 3
quarantine reasons:
 - BIO-17B : temperature excursion 21m exceeds single-segment limit 15m
 - BIO-18A : descendant of recalled ancestor Seed-Beta
 - BIO-19A : custody ledger break at step carrier handoff (expected sha256:60c4879ecc7b7e78…, got sha256:0000000000000000…)
allocation candidates : 7 feasible lot→clinic assignments + 3 hold options
optimizer states : 30
optimizer transitions : 36
memo hits : 3

## Check
C1 OK - every final lot has a computed ancestry closure.
C2 OK - recalled Seed-Beta propagates to BIO-18A and not to unrelated safe lots.
C3 OK - custody ledgers are hash-linked or explicitly rejected.
C4 OK - temperature policy rejects BIO-17B but accepts the three released lots.
C5 OK - no quarantined lot is offered to allocation.
C6 OK - every selected shipment respects clinic transit limits.
C7 OK - clinic demand caps are never exceeded.
C8 OK - the exact optimizer proves no higher-value allocation exists.
C9 OK - delivered dose/value totals equal the selected shipments.
C10 OK - reserve doses remain bound to the original ledger root.

## Go audit details
platform : go1.26.2 linux/amd64
case : cold-chain-release
question : Which cold-chain lots can be released, and where should scarce doses ship first?
product : Stable mRNA pediatric booster
lots total : 11
release candidate lots : 6
lineage parent facts : 8
lineage ancestor closure facts : 12
recalled source lots : Seed-Beta
custody event facts : 18
temperature observations : 18
ledger-valid release candidates : 5/6
temperature-valid release candidates : 5/6
releasable candidates : 3/6
quarantined candidates : 3/6
ledger roots:
 - BIO-17A root=sha256:e57b3009e0aa42a9… status=valid
 - BIO-17B root=sha256:73dea83c53094fcb… status=valid
 - BIO-18A root=sha256:c4fdb3d693abb849… status=valid
 - BIO-19A root=sha256:5a0ab02c71efbb91… status=invalid
 - BIO-20A root=sha256:1ffb664e71e3c68a… status=valid
 - BIO-21A root=sha256:67bdcbb64f4012e6… status=valid
release decisions:
 - BIO-17A release doses=420 facility=CityHub tempExcursion=6m maxSegment=6m ancestors=Bulk-A, Seed-Alpha
 - BIO-17B quarantine reason=temperature excursion 21m exceeds single-segment limit 15m
 - BIO-18A quarantine reason=descendant of recalled ancestor Seed-Beta
 - BIO-19A quarantine reason=custody ledger break at step carrier handoff (expected sha256:60c4879ecc7b7e78…, got sha256:0000000000000000…)
 - BIO-20A release doses=300 facility=NorthDepot tempExcursion=12m maxSegment=12m ancestors=Bulk-A, Seed-Alpha
 - BIO-21A release doses=180 facility=AirportHub tempExcursion=9m maxSegment=9m ancestors=Seed-Gamma
clinics:
 - Dialysis Unit demand=300 priority=10 maxTransit=100m need=immunocompromised pediatric list
 - North Clinic demand=380 priority=9 maxTransit=130m need=school outbreak ring
 - Field Station demand=250 priority=8 maxTransit=90m need=temporary shelter
 - Rural Mobile Unit demand=180 priority=6 maxTransit=140m need=mobile catch-up route
transit candidates:
 - BIO-17A -> Dialysis Unit 70m value=3000
 - BIO-17A -> North Clinic 110m value=3420
 - BIO-20A -> Dialysis Unit 90m value=3000
 - BIO-20A -> North Clinic 40m value=2700
 - BIO-21A -> Field Station 70m value=1440
 - BIO-21A -> North Clinic 130m value=1620
 - BIO-21A -> Rural Mobile Unit 90m value=1080
selected allocation:
 - BIO-17A -> North Clinic delivered=380 reserve=40 transit=110m value=3420
 - BIO-20A -> Dialysis Unit delivered=300 reserve=0 transit=90m value=3000
 - BIO-21A -> Field Station delivered=180 reserve=0 transit=70m value=1440
delivered doses : 860
reserved doses : 40
unmet demand : 250 (Field Station 70 + Rural Mobile Unit 180)
best value : 7860
total selected transit : 270m
optimizer states explored : 30
optimizer memo hits : 3
optimizer transitions : 36
checks passed : 10/10
all checks pass : yes
