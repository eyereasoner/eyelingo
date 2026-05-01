# French Cities  

## Answer  
Four cities in this small network can reach Nantes: Paris, Chartres, Le Mans, Angers.  

## Reason  
The original example says that every :oneway link is also a :path, and that :path is transitive. So once Angers can reach Nantes directly, longer routes can be built by chaining earlier links. Angers reaches Nantes directly. Le Mans reaches Nantes through Angers. Chartres reaches Nantes through Le Mans and Angers. Paris reaches Nantes through Chartres, Le Mans, and Angers.  

## Check  
C1 OK - Angers has a direct one-way edge to Nantes  
C2 OK - Le Mans reaches Nantes through Angers  
C3 OK - Chartres reaches Nantes through Le Mans and Angers  
C4 OK - Paris reaches Nantes through Chartres, Le Mans, and Angers  
C5 OK - exactly four non-destination cities can reach Nantes  
C6 OK - the reported city set matches the transitive closure  
C7 OK - cities without any chain to Nantes are rejected  
