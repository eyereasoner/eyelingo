# Bayes Therapy Decision Support  

## Answer  
The recommended therapy is Paxlovid (utility = 3.585174).  

Full posterior distribution:  
  COVID19               posterior = 0.483883  (unnormalized = 0.00928200)  
  Influenza             posterior = 0.427894  (unnormalized = 0.00820800)  
  AllergicRhinitis      posterior = 0.006686  (unnormalized = 0.00012825)  
  BacterialPneumonia    posterior = 0.081538  (unnormalized = 0.00156408)  

Therapy utility scores:  
  Paxlovid              expectedSuccess = 0.388517  adverse = 0.10  utility = 3.585174  
  Oseltamivir           expectedSuccess = 0.285141  adverse = 0.08  utility = 2.611410  
  Antihistamine         expectedSuccess = 0.100269  adverse = 0.03  utility = 0.912689  
  Antibiotic            expectedSuccess = 0.110953  adverse = 0.07  utility = 0.899526  
  SupportiveCare        expectedSuccess = 0.291512  adverse = 0.01  utility = 2.885120  

## Reason  
Evidence: Fever=present, DryCough=present, LossOfSmell=absent, Sneezing=absent, ShortBreath=absent.  
Evidence total (normalizing constant) = 0.01918233.  

The posterior for each disease is computed as:  
  posterior(d) = prior(d) × ∏ P(symptom|d) / evidenceTotal  
where for an absent symptom the factor is 1 − P(symptom|d).  

For each therapy, expected success is:  
  expectedSuccess(t) = Σ_i posterior(i) × successByDisease(i)  
and utility = benefitWeight × expectedSuccess − harmWeight × adverse.  
The recommended therapy is the one with the highest utility.  

## Check  
C1 OK - all disease priors are probabilities and the model is deliberately incomplete  
C2 OK - every conditional probability is in [0, 1]  
C3 OK - the evidence normalizing constant is recomputed independently  
C4 OK - reported posteriors and unnormalized likelihoods match recomputation  
C5 OK - reported posteriors sum to one after rounding  
C6 OK - all therapy rows are reported  
C7 OK - expected success and adverse-risk inputs match independent utility computation  
C8 OK - each reported utility is benefit-weighted success minus harm-weighted adverse risk  
C9 OK - the recommended therapy is independently selected as maximum utility  
