package main

import (
	"fmt"
	"math"
)

func checkRC(ctx *Context) []Check {
	d := ctx.M()
	lower, upper := num(d["decayLower"]), num(d["decayUpper"])
	initial := num(d["initialVoltage"])
	tol := num(d["tolerance"])
	period := num(d["samplePeriod"])
	maxs := integer(d["maxStep"])
	up := make([]float64, maxs+1)
	lo := make([]float64, maxs+1)
	first := -1
	for s := 0; s <= maxs; s++ {
		up[s] = initial * math.Pow(upper, float64(s))
		lo[s] = initial * math.Pow(lower, float64(s))
		if first < 0 && up[s] < tol {
			first = s
		}
	}
	dec := true
	bracket := true
	for s := 0; s < maxs; s++ {
		dec = dec && up[s+1] < up[s]
	}
	for s := 0; s <= maxs; s++ {
		bracket = bracket && lo[s] <= up[s]
	}
	rv := fieldFloat(ctx.Answer, fmt.Sprintf("upper voltage at step %d", first))
	return []Check{{"the decay interval is nonempty, positive, and below one", 0 < lower && lower <= upper && upper < 1}, {"the certified upper voltage envelope decreases at every sample", dec}, {"the lower and upper envelopes bracket every compatible decay", bracket}, {"step 12 remains above the voltage tolerance", up[first-1] >= tol}, {"step 13 is the first certified below-tolerance sample", first == integer(asMap(d["expected"])["firstSettledStep"])}, {"reported first step, time, and upper voltage match recomputation", fieldInt(ctx.Answer, "first below tolerance step") == first && close(fieldFloat(ctx.Answer, "first below tolerance time"), float64(first)*period, 5e-4) && close(rv, up[first], 5e-4)}, {"the report uses the JSON double interval rather than only the exact symbol", contains(ctx.Answer, fmt.Sprintf("[%0.10f, %0.10f]", lower, upper)) && contains(ctx.Answer, str(d["exactDecaySymbol"]))}}
}
