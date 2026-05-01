# Digital Product Passport  

## Answer  
Passport decision : PASS for ACME X1000 SN123.  
recycled content : 13%  
lifecycle footprint : 52500 gCO2e  
total component mass : 105 g  
critical raw materials : Cobalt, Lithium  
circularity hint : repairFriendly  
public endpoint : https://example.org/dpp/ACME-X1000-SN123  

## Reason why  
The passport folds the explicit component list to derive total mass and recycled mass, then computes an integer recycled-content percentage.  
Lifecycle footprint is derived by summing manufacturing, transport, and use-phase emissions.  
The product is repair-friendly because the battery is replaceable and the public passport section exposes both repair and spare-parts documentation.  
Critical raw-material exposure is derived from component-material links: Cobalt, Lithium.  

Component roll-up:  
BatteryPack-01 Battery mass=48g recycled=0g materials=Lithium, Cobalt, Nickel replaceable=true  
Chassis-01 Housing mass=32g recycled=12g materials=Aluminium replaceable=false  
Mainboard-01 Electronics mass=25g recycled=2g materials=Copper, GoldTrace replaceable=false  
Public documents:  
Doc-UserManual UserManual https://example.org/manuals/acme-x1000  
Doc-RepairGuide RepairGuide https://example.org/repair/acme-x1000  
Doc-SpareParts SparePartsCatalog https://example.org/spares/acme-x1000  

## Check  
C1 OK - component masses are folded from the JSON component list  
C2 OK - recycled mass and integer recycled-content percentage are recomputed independently  
C3 OK - manufacturing, transport, and use-phase footprint values sum to the reported lifecycle footprint  
C4 OK - critical raw-material exposure is derived by joining component materials to material declarations  
C5 OK - repairFriendly is derived from a replaceable battery plus public repair and spare-parts documents  
C6 OK - all required public document types are present only in the public document section  
C7 OK - all restricted declaration document types stay in the restricted section  
C8 OK - lifecycle events are chronological and follow manufacturing before sale before repair  
C9 OK - the passport endpoint equals the product digital link and is reported publicly  
C10 OK - PASS is reported only because every independent passport check succeeds  
