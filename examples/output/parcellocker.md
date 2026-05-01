# Parcel Locker  

## Answer  
May Noah use Maya's one-time pickup token to collect parcel123 from locker B17?  
decision : PERMIT  
release : Noah may collect parcel123 for Maya from locker B17 at Station West  
guardrail denials : 5/5  

Request decisions:  
  pickup        : PERMIT (source request)  
  billing       : DENY  
  redirect      : DENY  
  wrong-person  : DENY  
  wrong-locker  : DENY  
  reuse         : DENY  

## Reason why  
token : delegate=noah parcel=parcel123 locker=lockerB17 action=CollectParcel purpose=PickupOnly state=Active reuse=SingleUse  
privacy : billingAccess=None redirectAllowed=No  
parcel : parcel123 status=ReadyForPickup  

Noah collects the parcel  
  request : requester=noah parcel=parcel123 locker=lockerB17 action=CollectParcel purpose=PickupOnly  
  decision : PERMIT  
  reason : Permit: the requester, parcel, locker, action, purpose, active state, single-use limit, parcel readiness, and privacy restrictions all match.  
  passed conditions : 10/10  

Noah views billing details  
  request : requester=noah parcel=parcel123 locker=lockerB17 action=ViewBillingDetails purpose=BillingAccess  
  decision : DENY  
  reason : Deny: C4 requested action must be parcel collection; C5 requested purpose must be pickup only; C9 billing details must stay hidden.  
  passed conditions : 7/10  

Noah redirects the parcel  
  request : requester=noah parcel=parcel123 locker=lockerB17 action=RedirectParcel purpose=RedirectDelivery  
  decision : DENY  
  reason : Deny: C4 requested action must be parcel collection; C5 requested purpose must be pickup only; C10 parcel redirection must stay blocked.  
  passed conditions : 7/10  

Maya uses Noah's pickup token  
  request : requester=maya parcel=parcel123 locker=lockerB17 action=CollectParcel purpose=PickupOnly  
  decision : DENY  
  reason : Deny: C1 requester must match the named delegate.  
  passed conditions : 9/10  

Noah tries another locker  
  request : requester=noah parcel=parcel123 locker=lockerA04 action=CollectParcel purpose=PickupOnly  
  decision : DENY  
  reason : Deny: C3 requested locker must match the authorized locker.  
  passed conditions : 9/10  

Noah reuses the token  
  request : requester=noah parcel=parcel123 locker=lockerB17 action=CollectParcel purpose=PickupOnly  
  state override : token has already been used once  
  decision : DENY  
  reason : Deny: C6 authorization must be active and not already consumed.  
  passed conditions : 9/10

## Check  
C1 OK - the source pickup request satisfies all ten authorization conditions  
C2 OK - all reported request decisions match independent policy evaluation  
C3 OK - billing access is denied by the privacy guardrail  
C4 OK - redirect is denied by the parcel-redirection guardrail  
C5 OK - wrong-person use is denied because requester must match the delegate  
C6 OK - wrong-locker use is denied because locker must match the token  
C7 OK - single-use reuse is denied after the token is already consumed  
C8 OK - guardrail denials recompute to five out of five  
C9 OK - the release text matches parcel owner, delegate, parcel, locker, and site
