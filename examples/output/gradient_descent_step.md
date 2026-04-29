# Gradient Descent Step  

## Answer  
start point : (8.000, 5.000)  
gradient : (10.000, 24.000)  
step size : 0.100  
next point : (7.000, 2.600)  
objective before : 97.000  
objective after : 41.920  
decrease : 55.080  

## Reason why  
The quadratic is convex and its gradient is computed symbolically from the JSON coefficients.  
The update uses x_next = x - alpha × gradient, with alpha fixed by the input.  
The new point stays within the declared step-norm bound and produces a strictly smaller objective value.  
That gives a small finite certificate for one safe descent step rather than an open-ended optimizer.  

## Check  
C1 OK - gradient was derived from f(x,y) = (x-3)^2 + 2(y+1)^2  
C2 OK - step size is positive and below the conservative bound  
C3 OK - new point is (7.000, 2.600)  
C4 OK - objective value decreases from 97.000 to 41.920  
C5 OK - step norm stays below 3.000  
C6 OK - the JSON expected decrease flag is satisfied  

## Go audit details  
platform : go1.26.2 linux/amd64  
case : gradient_descent_step  
question : Does the certified gradient step reduce a convex quadratic while staying inside the step bound?  
center : (3.000, -1.000)  
weightY : 2.000  
step norm : 2.600  
max step norm : 3.000  
checks passed : 6/6  
