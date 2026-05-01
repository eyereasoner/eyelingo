# Doctor Advice Work Conflict  

## Answer  
overall decision for Jos : RemoteWorkOnly  
Request_Jos_Prog_Home : raw=Deny+Permit status=BothPermitDeny effective=Permit  
Request_Jos_Prog_Office : raw=Deny+Permit status=BothPermitDeny effective=Deny  

## Reason why  
Jos has Flu, so the health context classifies the agent as Sick.  
The doctor's statement permits ProgrammingWork, but the general sick-work policy also denies work, so the raw closure deliberately keeps the conflict visible.  
The conflict resolver permits sick ProgrammingWork at Home but denies Office work because of colleague-infection risk.  
Since Home is permitted and Office is denied, the combined recommendation is RemoteWorkOnly.  

## Check  
C1 OK - Jos is classified as Sick from condition Flu  
C2 OK - home programming request keeps both Permit and Deny before resolution  
C3 OK - home programming conflict resolves to Permit  
C4 OK - office programming request keeps both Permit and Deny before resolution  
C5 OK - office programming conflict resolves to Deny  
C6 OK - overall work decision is RemoteWorkOnly  
