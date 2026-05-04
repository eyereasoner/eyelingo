# Bayes Diagnosis  

## Insight  
The most likely disease is COVID19 (posterior = 0.941209).  

Full posterior distribution:  
  COVID19               posterior = 0.941209  (unnormalized = 0.00154700)  
  Influenza             posterior = 0.029204  (unnormalized = 0.00004800)  
  AllergicRhinitis      posterior = 0.000456  (unnormalized = 0.00000075)  
  BacterialPneumonia    posterior = 0.029131  (unnormalized = 0.00004788)  

## Explanation  
Evidence: Fever=present, DryCough=present, LossOfSmell=present, Sneezing=absent, ShortBreath=present.  
Evidence total (normalizing constant) = 0.00164363.  
The posterior for each disease is computed as:  
  posterior(d) = prior(d) × ∏ P(symptom|d) / evidenceTotal  
where for an absent symptom the factor is 1 − P(symptom|d).  
