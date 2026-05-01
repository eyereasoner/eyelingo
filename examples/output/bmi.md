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
C1 OK - input units are normalized to positive SI kg and m  
C2 OK - height squared is recomputed independently from the normalized height  
C3 OK - reported BMI matches independent kg/m² computation rounded to one decimal  
C4 OK - the unrounded BMI independently rounds to the expected two-decimal value  
C5 OK - reported category is the independent WHO category for the computed BMI  
C6 OK - WHO category boundaries are half-open at 18.5, 25, 30, 35, and 40  
C7 OK - reported healthy-weight lower bound equals BMI 18.5 at the same height  
C8 OK - reported healthy-weight upper bound equals BMI 24.9 at the same height  
C9 OK - the explanation mentions the same formula that the Python check recomputed  
