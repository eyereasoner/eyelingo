# French Cities

## Answer
Four cities in this small network can reach Nantes: Paris, Chartres, Le Mans, Angers.

## Reason why
The original example says that every :oneway link is also a :path, and that :path is transitive. So once Angers can reach Nantes directly, longer routes can be built by chaining earlier links. Angers reaches Nantes directly. Le Mans reaches Nantes through Angers. Chartres reaches Nantes through Le Mans and Angers. Paris reaches Nantes through Chartres, Le Mans, and Angers.

## Check
C1 OK - Angers has a direct one-way connection to Nantes.
C2 OK - Le Mans reaches Nantes by chaining Le Mans → Angers → Nantes.
C3 OK - Chartres reaches Nantes by chaining Chartres → Le Mans → Angers → Nantes.
C4 OK - Paris reaches Nantes by chaining Paris → Chartres → Le Mans → Angers → Nantes.
C5 OK - cities without a chain to Nantes are rejected by fail-loud fuse rules.

## Go audit details
platform : go1.26.2 linux/amd64
total edges (one-way roads) : 10
one-way connections:
  Paris → Orléans
  Paris → Chartres
  Paris → Amiens
  Orléans → Blois
  Orléans → Bourges
  Blois → Tours
  Chartres → Le Mans
  Le Mans → Angers
  Le Mans → Tours
  Angers → Nantes

reachable from Nantes:
  Paris : yes
  Chartres : yes
  Le Mans : yes
  Angers : yes
  Orléans : no
  Amiens : no
  Bourges : no
  Blois : no
  Tours : no
checks passed : 5/5
recommendation consistent : yes
