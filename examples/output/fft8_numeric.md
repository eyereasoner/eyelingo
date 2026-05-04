# FFT8 Numeric  

## Insight  
sample vector : 0.000000, 0.707107, 1.000000, 0.707107, 0.000000, -0.707107, -1.000000, -0.707107  
dominant bins : k=1 magnitude=4.000000 phase=-1.570796; k=7 magnitude=4.000000 phase=1.570796  
time-domain energy : 4.000000  
frequency-domain energy / 8 : 4.000000  

## Explanation  
The input samples describe one sine cycle over eight equally spaced samples.  
The DFT projects the signal onto eight complex roots of unity.  
A real sine wave has equal magnitude at the positive and negative frequency bins.  
All non-dominant bins cancel to zero within the configured numerical tolerance.  
