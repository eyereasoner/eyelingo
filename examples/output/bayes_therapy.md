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
C1 OK - all priors are valid probabilities and the model has nonzero prior mass  
C2 OK - all evidence symptoms are available for every disease and absent evidence uses complement factors  
C3 OK - the Bayesian normalizing constant is recomputed independently  
C4 OK - the reported posterior table contains one row for each disease  
C5 OK - each reported unnormalized disease likelihood matches the Python recomputation  
C6 OK - each reported posterior equals likelihood divided by the evidence total  
C7 OK - therapy success vectors align with diseases and all therapy probabilities are in [0, 1]  
C8 OK - each expected therapy success is recomputed as posterior-weighted disease success  
C9 OK - each therapy utility applies benefitWeight × expectedSuccess − harmWeight × adverse  
C10 OK - the reported recommendation is the maximum-utility therapy  
