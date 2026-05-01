# Alarm Bit Interoperability  

## Answer  
classical alarm-bit interoperability : YES  
universal cloning of the superinformation token : NO  
copy task : opticalBeacon -> relayRegister for AlarmBit  
copy task : relayRegister -> opticalBeacon for AlarmBit  
blocked tasks : UniversalClone, UnrestrictedStateFanOut  

## Reason why  
The optical beacon and relay register are unlike physical media, but both encode the same abstract AlarmBit variable.  
Because the variable is classical in both media, local permutation and copying in both directed media transfers are possible.  
The quantumToken is treated as a superinformation medium with states Horizontal, Vertical, Diagonal, AntiDiagonal.  
That contrast substrate cannot support universal cloning of all states, so unrestricted classical-style fan-out is also blocked.  

## Check  
C1 OK - all classical media encode the same abstract variable  
C2 OK - the directed copy-task count is recomputed from the media graph  
C3 OK - the report lists exactly the expected directed copy tasks  
C4 OK - each classical medium has distinguishable zero and one states  
C5 OK - the superinformation contrast has more than two named states  
C6 OK - all expected impossible tasks are reported as blocked  
C7 OK - the reported CAN/CANNOT decisions match the expected polarities  
