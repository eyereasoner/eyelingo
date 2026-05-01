# Drone Corridor Planner  

## Answer  
selected plan : fly_gent_brugge -> public_coastline_brugge_oostende  
duration : 2800 s  
cost : 0.012  
belief : 0.960300  
comfort : 0.950400  
end state : location=Oostende battery=low permit=none fuelLeft=5  
surviving plans : 17  

## Reason  
The planner treats each corridor description as a state transition over location, battery, and permit state.  
Search is fuel-bounded to 7 steps, which keeps cycles finite while still allowing charging and permit actions.  
For composed plans, duration and cost are summed while belief and comfort are multiplied along the path.  
Plans are retained only when belief is greater than 0.94 and cost is less than 0.03.  
The selected plan is the lowest-cost surviving plan; the next cheapest starts with fly_gent_brugge and costs 0.014.  

## Check  
C1 OK - 14 corridor actions are loaded from JSON  
C2 OK - bounded search independently finds 17 surviving plans  
C3 OK - the selected path is the lowest-cost survivor  
C4 OK - the selected path starts with the expected first action  
C5 OK - duration, cost, belief, and comfort are recomputed along the selected path  
C6 OK - the selected end state reaches Oostende with low battery and no permit  
C7 OK - the selected belief and cost satisfy the thresholds  
C8 OK - state-cycle pruning keeps every selected-path state unique  
C9 OK - the next cheapest survivor costs 0.014 as stated  
