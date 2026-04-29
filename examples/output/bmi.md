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
C1 OK - the input was normalized into positive SI values.  
C2 OK - height squared was reconstructed from the normalized height.  
C3 OK - the BMI value matches the BMI = kg / m² formula.  
C4 OK - a BMI of 18.49 stays below the normal-weight threshold.  
C5 OK - the lower boundary is half-open: BMI 18.5 is classified as Normal.  
C6 OK - BMI 25.0 starts the Overweight category.  
C7 OK - BMI 30.0 starts the Obesity I category.  
C8 OK - classification behavior is monotonic across representative BMI values.  
C9 OK - the healthy-weight band was reconstructed from BMI 18.5 to 24.9 at the same height.  

## Go audit details  
platform : go1.26.2 linux/amd64  
unit system : metric  
input weight : 72.0  
input height : 178.0  
normalized weight (kg) : 72.000000  
normalized height (m) : 1.780000  
height squared (m²) : 3.168400  
bmi : 22.724403  
bmi rounded : 22.7  
category : Normal  
healthy min (kg) : 58.6  
healthy max (kg) : 78.9  
checks passed : 9/9  
recommendation consistent : yes  
