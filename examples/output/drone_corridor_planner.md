# Drone Corridor Planner  

## Answer  
selected plan : fly_gent_brugge -> public_coastline_brugge_oostende  
duration : 2800 s  
cost : 0.012  
belief : 0.960300  
comfort : 0.950400  
end state : location=Oostende battery=low permit=none fuelLeft=5  
surviving plans : 17  

## Reason why  
The planner treats each corridor description as a state transition over location, battery, and permit state.  
Search is fuel-bounded to 7 steps, which keeps cycles finite while still allowing charging and permit actions.  
For composed plans, duration and cost are summed while belief and comfort are multiplied along the path.  
Plans are retained only when belief is greater than 0.94 and cost is less than 0.03.  
The selected plan is the lowest-cost surviving plan; the next cheapest starts with fly_gent_brugge and costs 0.014.  

## Check  
C1 OK - 14 corridor actions were loaded from JSON  
C2 OK - bounded search found 17 plans meeting belief and cost thresholds  
C3 OK - lowest-cost selected plan starts with fly_gent_brugge  
C4 OK - selected plan reaches Oostende  
C5 OK - selected belief 0.960300 is above 0.94  
C6 OK - selected cost 0.012 is below 0.03  
