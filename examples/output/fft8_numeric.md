# FFT8 Numeric  

## Answer  
sample vector : 0.000000, 0.707107, 1.000000, 0.707107, 0.000000, -0.707107, -1.000000, -0.707107  
dominant bins : k=1 magnitude=4.000000 phase=-1.570796; k=7 magnitude=4.000000 phase=1.570796  
time-domain energy : 4.000000  
frequency-domain energy / 8 : 4.000000  

## Reason why  
The input samples describe one sine cycle over eight equally spaced samples.  
The DFT projects the signal onto eight complex roots of unity.  
A real sine wave has equal magnitude at the positive and negative frequency bins.  
All non-dominant bins cancel to zero within the configured numerical tolerance.  

## Check  
C1 OK - the input contains exactly 8 time-domain samples  
C2 OK - the dominant bins are k=1 and k=7  
C3 OK - the DC component is zero  
C4 OK - the two conjugate sine bins have magnitude 4  
C5 OK - Parseval energy is preserved by the unnormalized DFT convention  
C6 OK - real-valued samples produce conjugate-symmetric bins  

## Go audit details  
platform : go1.26.2 linux/amd64  
case : fft8-numeric  
question : Compute the dominant DFT bins for an 8-sample sine wave.  
sample count : 8  
tolerance : 1.0e-09  
bins evaluated : 8  
dominant bin count : 2  
checks passed : 6/6  
