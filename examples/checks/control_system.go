package main

import (
	"math"
)

func checkControlSystem(ctx *Context) []Check {
	d := ctx.M()
	input1, branch1 := measurement10(farr(asMap(d["input1"])["measurement1"]))
	dist2, branch2 := measurement10(farr(asMap(d["disturbance2"])["measurement1"]))
	product := input1 * 19.6
	compensation := math.Log10(num(asMap(d["disturbance1"])["measurement3"]))
	actuator1 := product - compensation
	errorVal := num(asMap(d["output2"])["target2"]) - num(asMap(d["output2"])["measurement4"])
	diff := num(asMap(d["state3"])["observation3"]) - num(asMap(d["output2"])["measurement4"])
	nonlinear := 7.3 / errorVal * diff
	actuator2 := 5.8*errorVal + nonlinear
	rep := parseControlAnswer(ctx.Answer)
	cr := func(k string, exp float64) bool { v, ok := rep[k]; return ok && math.Abs(v-exp) <= 0.5e-6 }
	return []Check{{"input1 measurement10 recomputes through the lessThan square-root branch", branch1 == "lessThan" && cr("input1_m10", input1)}, {"disturbance2 measurement10 recomputes through the notLessThan branch", branch2 == "notLessThan" && cr("disturbance2_m10", dist2)}, {"feedforward guard is true before actuator1 arithmetic is applied", boolean(asMap(d["input2"])["measurement2"])}, {"actuator1 control recomputes product minus log10 compensation", cr("actuator1", actuator1)}, {"target-minus-measurement error is recomputed as 5", errorVal == 5}, {"state/output differential error is recomputed as -2", diff == -2}, {"nonlinear feedback term uses 7.3/error times the differential", math.Abs(nonlinear-(-2.92)) < 1e-12}, {"actuator2 control recomputes proportional plus nonlinear feedback", cr("actuator2", actuator2)}}
}
