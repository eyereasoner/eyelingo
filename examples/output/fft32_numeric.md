# FFT32 Numeric  

## Answer  
waveforms evaluated : 6  
alternating : k=16 magnitude=32.000000 phase=0.000000; energy=32.000000  
constant_ones : k=0 magnitude=32.000000 phase=0.000000; energy=32.000000  
cosine_bin3 : k=3 magnitude=16.000000 phase=0.000000; k=29 magnitude=16.000000 phase=0.000000; energy=16.000000  
impulse : all 32 bins magnitude=1.000000; energy=1.000000  
ramp_0_31 : k=0 magnitude=496.000000 phase=0.000000; energy=10416.000000  
sine_bin5 : k=5 magnitude=16.000000 phase=-1.570796; k=27 magnitude=16.000000 phase=1.570796; energy=16.000000  

## Reason why  
The upstream FFT32 fixture defines several 32-sample waveforms and asks for the whole spectrum of each waveform.  
The Go translation evaluates every frequency bin by summing samples against the corresponding complex root of unity.  
Constant, alternating, cosine, sine, impulse, and ramp fixtures exercise different spectral shapes.  
The checks verify dominant bins, magnitudes, flat impulse spectrum, conjugate symmetry, and energy preservation.  

## Check  
C1 OK - each generated waveform has exactly 32 samples  
C2 OK - dominant bins match the fixture expectations  
C3 OK - dominant magnitudes match the configured certificates  
C4 OK - Parseval energy is preserved for every waveform  
C5 OK - real-valued waveforms produce conjugate-symmetric spectra  
C6 OK - the impulse waveform has a flat unit-magnitude spectrum  
C7 OK - the answer reports every expected dominant bin by waveform name  
C8 OK - the configured operation certificate matches six full 32-bin direct DFTs  
