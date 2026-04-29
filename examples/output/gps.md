# GPS — Goal driven route planning

## Answer
Take the direct route via Brugge.
- Recommended route: Gent → Brugge → Oostende

## Reason why
From Gent to Oostende, the planner found two routes in this small map.
The direct route (Gent → Brugge → Oostende) takes 2400.0 seconds at cost 0.01, with belief 0.9408 and comfort 0.99. The alternative (Gent → Kortrijk → Brugge → Oostende) takes 4100.0 seconds at cost 0.018, with belief 0.903168 and comfort 0.9801.
So the direct route is faster, cheaper, more reliable, and slightly more comfortable.

## Check
- C1 OK - the direct Gent → Brugge → Oostende route was derived.
- C2 OK - the alternative Gent → Kortrijk → Brugge → Oostende route was derived.
- C3 OK - the recommended route is faster than the alternative.
- C4 OK - the recommended route is cheaper than the alternative.
- C5 OK - the recommended route has higher belief and comfort scores.

## Go audit details
- platform : go1.26.2 linux/amd64
- question : Which route should we take from Gent to Oostende?
- traveller : i1
- start : Gent
- goal : Oostende
- map edges : 4
- named routes : 2
- paths derived : 2
derived path 1 : drive_gent_brugge -> drive_brugge_oostende | duration=2400.0 cost=0.01 belief=0.9408 comfort=0.99
derived path 2 : drive_gent_kortrijk -> drive_kortrijk_brugge -> drive_brugge_oostende | duration=4100.0 cost=0.018 belief=0.903168 comfort=0.9801
- direct actions : drive_gent_brugge -> drive_brugge_oostende
- alternative actions : drive_gent_kortrijk -> drive_kortrijk_brugge -> drive_brugge_oostende
- duration advantage seconds : 1700.0
- cost advantage : 0.008
- belief advantage : 0.037632
- comfort advantage : 0.0099
- search recursive calls : 4
- search edge tests : 16
- search edges extended : 5
- search revisit prunes : 0
- search max depth : 2
- checks passed : 5/5
- recommendation consistent : yes
