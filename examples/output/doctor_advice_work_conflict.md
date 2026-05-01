# Doctor Advice Work Conflict  

## Answer  
overall decision for Jos : RemoteWorkOnly  
Request_Jos_Prog_Home : raw=Deny+Permit status=BothPermitDeny effective=Permit  
Request_Jos_Prog_Office : raw=Deny+Permit status=BothPermitDeny effective=Deny  

## Reason  
Jos has Flu, so the health context classifies the agent as Sick.  
The doctor's statement permits ProgrammingWork, but the general sick-work policy also denies work, so the raw closure deliberately keeps the conflict visible.  
The conflict resolver permits sick ProgrammingWork at Home but denies Office work because of colleague-infection risk.  
Since Home is permitted and Office is denied, the combined recommendation is RemoteWorkOnly.  

## Check  
C1 OK - Flu classifies Jos as sick for the policy conflict  
C2 OK - ProgrammingWork is closed upward to Work  
C3 OK - doctor advice contributes Permit for every ProgrammingWork request  
C4 OK - sick-work default contributes Deny for both requests  
C5 OK - the home request keeps the raw Permit+Deny conflict before resolution  
C6 OK - the office request keeps the raw Permit+Deny conflict before resolution  
C7 OK - conflict resolution permits sick programming work at Home  
C8 OK - conflict resolution denies Office work  
C9 OK - the combined recommendation recomputes to RemoteWorkOnly  
