# Bayes Diagnosis  

## Answer  
The most likely disease is COVID19 (posterior = 0.941209).  

Full posterior distribution:  
  COVID19               posterior = 0.941209  (unnormalized = 0.00154700)  
  Influenza             posterior = 0.029204  (unnormalized = 0.00004800)  
  AllergicRhinitis      posterior = 0.000456  (unnormalized = 0.00000075)  
  BacterialPneumonia    posterior = 0.029131  (unnormalized = 0.00004788)  

## Reason why  
Evidence: Fever=present, DryCough=present, LossOfSmell=present, Sneezing=absent, ShortBreath=present.  
Evidence total (normalizing constant) = 0.00164363.  
The posterior for each disease is computed as:  
  posterior(d) = prior(d) × ∏ P(symptom|d) / evidenceTotal  
where for an absent symptom the factor is 1 − P(symptom|d).  

## Check  
C1 OK - all prior probabilities are in [0,1]  
C2 OK - all conditional probabilities are in [0,1]  
C3 OK - the evidence total is non-zero and reported as the Bayesian normalizing constant  
C4 OK - COVID19 has the largest posterior probability  
C5 OK - the posterior distribution contains four diseases  
C6 OK - absent Sneezing is handled through a complement likelihood factor  
C7 OK - posterior probabilities are normalized by the evidence total  
