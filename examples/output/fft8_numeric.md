# FFT8 Numeric  

## Answer  
sample vector : 0.000000, 0.707107, 1.000000, 0.707107, 0.000000, -0.707107, -1.000000, -0.707107  
dominant bins : k=1 magnitude=4.000000 phase=-1.570796; k=7 magnitude=4.000000 phase=1.570796  
time-domain energy : 4.000000  
frequency-domain energy / 8 : 4.000000  

## Reason  
The input samples describe one sine cycle over eight equally spaced samples.  
The DFT projects the signal onto eight complex roots of unity.  
A real sine wave has equal magnitude at the positive and negative frequency bins.  
All non-dominant bins cancel to zero within the configured numerical tolerance.  

## Check  
C1 OK - the input contains exactly eight samples  
C2 OK - dominant bins are recomputed as k=1 and k=7  
C3 OK - the DC component cancels to zero within tolerance  
C4 OK - the two sine bins have magnitude four  
C5 OK - reported dominant magnitudes and phases match recomputation  
C6 OK - Parseval energy is preserved under the unnormalized DFT convention  
C7 OK - reported time and frequency-domain energies match recomputation  
C8 OK - real samples produce conjugate-symmetric bins  
