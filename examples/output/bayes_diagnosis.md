# Bayes Diagnosis  

## Answer  
The most likely disease is COVID19 (posterior = 0.941209).  

Full posterior distribution:  
  COVID19               posterior = 0.941209  (unnormalized = 0.00154700)  
  Influenza             posterior = 0.029204  (unnormalized = 0.00004800)  
  AllergicRhinitis      posterior = 0.000456  (unnormalized = 0.00000075)  
  BacterialPneumonia    posterior = 0.029131  (unnormalized = 0.00004788)  

## Reason  
Evidence: Fever=present, DryCough=present, LossOfSmell=present, Sneezing=absent, ShortBreath=present.  
Evidence total (normalizing constant) = 0.00164363.  
The posterior for each disease is computed as:  
  posterior(d) = prior(d) × ∏ P(symptom|d) / evidenceTotal  
where for an absent symptom the factor is 1 − P(symptom|d).  

## Check  
C1 OK - all priors are probabilities and the prior mass is less than one  
C2 OK - every conditional probability is in [0, 1]  
C3 OK - all evidence symptoms are available for every disease  
C4 OK - the absent Sneezing evidence uses the complement likelihood  
C5 OK - the Bayesian normalizing constant is recomputed independently  
C6 OK - the reported distribution contains one posterior for each disease  
C7 OK - each reported unnormalized likelihood matches the Go recomputation  
C8 OK - each reported posterior matches likelihood divided by evidence total  
C9 OK - the reported posteriors sum to one after rounding  
C10 OK - COVID19 is independently selected as the maximum-posterior disease  
