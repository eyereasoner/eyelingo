# AuroraCare  

## Answer  
For each AuroraCare scenario, should the PDP permit or deny the requested use of health data, and why?  
permit count : 4  
deny count : 3  

Scenario decisions:  
  A – Primary care visit : PERMIT (urn:policy:primary-care-001)  
  B – Quality improvement (in scope) : PERMIT (urn:policy:qi-2025-aurora)  
  C – Quality improvement (out of scope) : DENY (no policy matched)  
  D – Insurance management : DENY (urn:policy:deny-insurance)  
  E – GP checks labs : PERMIT (urn:policy:primary-care-001)  
  F – Research on anonymised dataset : PERMIT (urn:policy:research-aurora-diabetes)  
  G – AI training (opt-out) : DENY (no policy matched)  

## Reason why  
A – Primary care visit  
  request : role=clinician purpose=PrimaryCareManagement environment=api_gateway categories=PATIENT_SUMMARY  
  decision : PERMIT  
  reason : Permitted: a clinician in the patient's care team matched the primary-care policy.  
  care-team linked : yes  
  subject opt-in : no  
  subject opt-out : no  
  trace : care_team_linked,urn:policy:primary-care-001:permit:permission_matched  

B – Quality improvement (in scope)  
  request : role=data_user purpose=EnsureQualitySafetyHealthcare environment=secure_env categories=LAB_RESULTS,PATIENT_SUMMARY  
  decision : PERMIT  
  reason : Permitted: the quality-improvement policy matched the secure environment and required data categories.  
  care-team linked : no  
  subject opt-in : no  
  subject opt-out : no  
  duties : requireConsent,noExfiltration  
  trace : urn:policy:qi-2025-aurora:permit:permission_matched  

C – Quality improvement (out of scope)  
  request : role=data_user purpose=EnsureQualitySafetyHealthcare environment=secure_env categories=LAB_RESULTS  
  decision : DENY  
  reason : Denied: no permission matched the purpose, environment, safeguards, role, or requested categories.  
  care-team linked : no  
  subject opt-in : no  
  subject opt-out : no  
  trace : deny:no_permission_matched  

D – Insurance management  
  request : role=data_user purpose=InsuranceManagement environment=secure_env categories=PATIENT_SUMMARY  
  decision : DENY  
  reason : Denied: the requested purpose is prohibited by policy.  
  care-team linked : no  
  subject opt-in : no  
  subject opt-out : no  
  trace : urn:policy:deny-insurance:deny:prohibition_matched  

E – GP checks labs  
  request : role=clinician purpose=PrimaryCareManagement environment=api_gateway categories=LAB_RESULTS  
  decision : PERMIT  
  reason : Permitted: a clinician in the patient's care team matched the primary-care policy.  
  care-team linked : yes  
  subject opt-in : no  
  subject opt-out : no  
  trace : care_team_linked,urn:policy:primary-care-001:permit:permission_matched  

F – Research on anonymised dataset  
  request : role=data_user purpose=HealthcareScientificResearch environment=secure_env categories=PATIENT_SUMMARY,LAB_RESULTS  
  safeguard : tom=Anonymisation  
  decision : PERMIT  
  reason : Permitted: the subject opted in, the dataset is anonymised, and the research policy matched.  
  care-team linked : no  
  subject opt-in : yes  
  subject opt-out : no  
  duties : annualOutcomeReport,noReidentification,noExfiltration  
  trace : subject_opted_in,urn:policy:research-aurora-diabetes:permit:permission_matched  

G – AI training (opt-out)  
  request : role=data_user purpose=TrainTestAndEvaluateAISystemsAlgorithms environment=secure_env categories=PATIENT_SUMMARY,LAB_RESULTS  
  decision : DENY  
  reason : Denied: the subject opted out of data use for AI training.  
  care-team linked : no  
  subject opt-in : no  
  subject opt-out : yes  
  trace : subject_opted_out,deny:subject_opted_out_ai_training  

## Check  
C1 OK - all seven scenarios match the PERMIT/DENY outcomes encoded in the N3 example  
C2 OK - primary-care access requires a clinician role and a care-team link  
C3 OK - quality improvement is allowed only when both lab results and patient summary are requested in the secure environment  
C4 OK - insurance management is denied by a matching prohibition  
C5 OK - research is allowed only with patient opt-in and anonymisation in the secure environment  
C6 OK - AI training is denied because the subject opted out  
C7 OK - four scenarios are permitted and three are denied  
