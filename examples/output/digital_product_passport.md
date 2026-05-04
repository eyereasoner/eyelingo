# Digital Product Passport  

## Insight  
Passport decision : PASS for ACME X1000 SN123.  
recycled content : 13%  
lifecycle footprint : 52500 gCO2e  
total component mass : 105 g  
critical raw materials : Cobalt, Lithium  
circularity hint : repairFriendly  
public endpoint : https://example.org/dpp/ACME-X1000-SN123  

## Explanation  
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
