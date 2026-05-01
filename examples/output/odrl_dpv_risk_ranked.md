# Ranked DPV Risk Report  

## Answer  
Agreement: Example Agreement  
Profile: Example consumer profile  

score=100 (risk:HighRisk, risk:HighSeverity) clause C1  
 Risk: account/data removal is permitted without notice safeguards (no notice constraint and no duty to inform). Clause C1: Provider may remove the user account (and associated data) at its discretion.  

mitigation for clause C1: Add a notice constraint (minimum noticeDays) before account removal.  
mitigation for clause C1: Add a duty to inform the consumer prior to account removal.  
score=97 (risk:HighRisk, risk:HighSeverity) clause C3  
 Risk: user data sharing is permitted without an explicit consent constraint. Clause C3: Provider may share user data with partners for business purposes.  

mitigation for clause C3: Add an explicit consent constraint before data sharing.  
score=85 (risk:HighRisk, risk:HighSeverity) clause C2  
 Risk: terms may change with notice (3 days) below consumer requirement (14 days). Clause C2: Provider may change terms by informing users at least 3 days in advance.  

mitigation for clause C2: Increase minimum noticeDays in the inform duty to meet the consumer requirement.  
score=70 (risk:ModerateRisk, risk:ModerateSeverity) clause C4  
 Risk: portability is restricted because exporting user data is prohibited. Clause C4: Users are not permitted to export their data.  

mitigation for clause C4: Add a permission allowing data export (or remove the prohibition) to support portability.  

## Reason why  
The agreement policy is scanned for permissions and prohibitions that conflict with the consumer profile needs.  
Each triggered rule derives a risk row with a normalized score, a source clause, and one or more mitigation measures.  
Rows are sorted by descending score so the highest-risk clauses are reviewed first.  

## Check  
C1 OK - 4 risk rows were derived  
C2 OK - ranked output is in descending score order  
C3 OK - score range is 70 to 100  
C4 OK - high=3 moderate=1 low=0 risk levels were derived  
C5 OK - 5 mitigation measures were generated  
