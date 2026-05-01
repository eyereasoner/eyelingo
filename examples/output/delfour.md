# Delfour  

## Answer  
The scanner is allowed to use a neutral shopping insight and recommends Low-Sugar Tea Biscuits instead of Classic Tea Biscuits.  
case : delfour  
decision : Allowed  
scanned product : Classic Tea Biscuits  
suggested alternative: Low-Sugar Tea Biscuits  

## Reason why  
The phone desensitizes a diabetes-related household condition into a scoped low-sugar need, wraps it in an expiring Insight + Policy envelope, signs it, and the scanner consumes that envelope for shopping assistance.  
metric : sugar_g_per_serving  
threshold : 10.0  
scope : self-scanner @ pick_up_scanner  
retailer : Delfour  
signature alg : HMAC-SHA256  
banner headline : Track sugar per serving while you scan  
expires at : 2025-10-05T22:33:48.907185+00:00  
reason.txt : Household requires low-sugar guidance (diabetes in POD). A neutral Insight is scoped to device 'self-scanner', event 'pick_up_scanner', retailer 'Delfour', and expires soon; the policy confines use to shopping assistance.  
audit entries : 1  
bus files written : 6  

## Check  
C1 OK - the scanner request satisfies the ODRL permission independently from Go  
C2 OK - the policy prohibition independently blocks marketing distribution  
C3 OK - the delete duty is tied to the insight expiry and is still before expiry  
C4 OK - the minimized serialized insight omits the sensitive medical condition  
C5 OK - the scanned product exceeds the sugar threshold, so a banner is warranted  
C6 OK - the selected alternative is the lowest-sugar lower-sugar catalog item  
C7 OK - the suggested alternative reduces sugar grams per serving by 9 g  
C8 OK - the payload SHA-256 is recomputed from the canonical escaped JSON  
C9 OK - the signature metadata is structurally valid for the trusted precomputed HMAC mode  
C10 OK - the reported bus and audit counts match the independent input fixture  
