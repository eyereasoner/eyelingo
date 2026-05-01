# BMI — ARC-style Body Mass Index example  

## Answer  
BMI = 22.7  
Category = Normal  
At height 178 cm, a healthy-weight range is about 58.6–78.9 kg (BMI 18.5–24.9).  

## Reason why  
BMI is defined as weight in kilograms divided by height in meters squared.  
The normalized weight and height were used to compute BMI, then the result was mapped to the WHO adult category table.  
The input was metric, so Inputs were already metric, so kilograms stay kilograms and centimeters are divided by 100 to obtain meters.  

## Check  
C1 OK - the input was normalized into positive SI values  
C2 OK - height squared was reconstructed from the normalized height  
C3 OK - the BMI value matches the BMI = kg / m² formula  
C4 OK - a BMI of 18.49 stays below the normal-weight threshold  
C5 OK - the lower boundary is half-open: BMI 18.5 is classified as Normal  
C6 OK - BMI 25.0 starts the Overweight category  
C7 OK - BMI 30.0 starts the Obesity I category  
C8 OK - classification behavior is monotonic across representative BMI values  
C9 OK - the healthy-weight band was reconstructed from BMI 18.5 to 24.9 at the same height  
