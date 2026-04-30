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
C2 OK - dominant bins match the FFT32 fixture expectations  
C3 OK - dominant magnitudes match the configured certificates  
C4 OK - Parseval energy is preserved for every waveform under the unnormalized DFT convention  
C5 OK - real-valued waveforms produce conjugate-symmetric spectra  
C6 OK - the impulse waveform has a flat unit-magnitude spectrum  
C7 OK - the full spectrum is computed once per waveform  

## Go audit details  
platform : go1.23.2 linux/amd64  
case : fft32-numeric  
source example : fft32-numeric.n3  
question : Compute whole 32-point FFT spectra for several waveform fixtures and certify their dominant bins.  
sample count per waveform : 32  
waveform count : 6  
complex bin sums : 6144  
tolerance : 1.0e-08  
checks passed : 7/7  
