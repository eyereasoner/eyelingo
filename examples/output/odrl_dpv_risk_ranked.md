# ODRL + DPV Risk Ranking  

## Insight  
agreement : Example Agreement  
profile : Example consumer profile  
risks found : 4  
highest risk : clause C1 score=100 HighRisk  

## Explanation  
The example treats the agreement as a small policy graph: permissions describe what the provider may do, prohibitions describe what users may not do, duties describe required follow-up actions, and constraints describe safeguards such as notice days or consent.  
Each rule is compared with the consumer profile. When a required safeguard is missing or too weak, the program creates a DPV-style risk record with a score, severity, risk level, consequences, impacts, and mitigation advice.  
Scores are base risk plus consumer-need importance, capped at 100, then ranked from highest to lowest so the most urgent clause appears first.  

Ranked risk report:  
score=100 (HighRisk, HighSeverity) clause C1  
 account/data removal is permitted without notice safeguards. Clause C1: Provider may remove the user account (and associated data) at its discretion.  
 consequences=DataLoss, DataUnavailable, CustomerConfidenceLoss impacts=FinancialLoss, NonMaterialDamage  
 - mitigation: Add a notice constraint (minimum noticeDays) before account removal.  
 - mitigation: Add a duty to inform the consumer prior to account removal.  
score=97 (HighRisk, HighSeverity) clause C3  
 user data sharing is permitted without an explicit consent constraint. Clause C3: Provider may share user data with partners for business purposes.  
 consequences=CustomerConfidenceLoss impacts=NonMaterialDamage, FinancialLoss  
 - mitigation: Add an explicit consent constraint before data sharing.  
score=85 (HighRisk, HighSeverity) clause C2  
 terms may change with notice (3 days) below consumer requirement (14 days). Clause C2: Provider may change terms by informing users at least 3 days in advance.  
 consequences=CustomerConfidenceLoss impacts=NonMaterialDamage  
 - mitigation: Increase minimum noticeDays in the inform duty to meet the consumer requirement.  
score=70 (ModerateRisk, ModerateSeverity) clause C4  
 portability is restricted because exporting user data is prohibited. Clause C4: Users are not permitted to export their data.  
 consequences=CustomerConfidenceLoss impacts=NonMaterialDamage  
 - mitigation: Add a permission allowing data export, or remove the prohibition, to support portability.  
