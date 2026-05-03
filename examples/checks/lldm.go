package main

import (
	"math"
)

func checkLLDM(ctx *Context) []Check {
	d := ctx.M()
	p1 := []float64{num(d["P1x"]), num(d["P1y"])}
	p2 := []float64{num(d["P2x"]), num(d["P2y"])}
	p3 := []float64{num(d["P3x"]), num(d["P3y"])}
	p4 := []float64{num(d["P4x"]), num(d["P4y"])}
	sl1 := (p2[1] - p1[1]) / (p2[0] - p1[0])
	b1 := p1[1] - sl1*p1[0]
	sl3 := -1 / sl1
	b3 := p3[1] - sl3*p3[0]
	b4 := p4[1] - sl3*p4[0]
	inter := func(m1, b1, m2, b2 float64) (float64, float64) { x := (b2 - b1) / (m1 - m2); return x, m1*x + b1 }
	p5x, p5y := inter(sl1, b1, sl3, b3)
	p6x, p6y := inter(sl1, b1, sl3, b4)
	dist := func(a, b []float64) float64 { return math.Hypot(a[0]-b[0], a[1]-b[1]) }
	d53 := dist([]float64{p5x, p5y}, p3)
	d64 := dist([]float64{p6x, p6y}, p4)
	dcm := d53 - d64
	val := func(pat string) float64 {
		if m := reFind(ctx.Answer, pat); m != nil {
			return parseFloat(m[0])
		}
		return math.NaN()
	}
	return []Check{{"L1 slope is recomputed from p1 and p2", close(val(`SL1 = (-?[0-9.]+)`), sl1, 5e-6)}, {"L3/L4 slopes are perpendicular to L1", close(val(`SL3 = SL4 = (-?[0-9.]+)`), sl3, 5e-6) && math.Abs(sl1*sl3+1) < 1e-12}, {"p5 is the analytic intersection of L1 and the perpendicular through p3", close(val(`p5\s+= \((-?[0-9.]+),`), p5x, 5e-4) && close(val(`p5\s+= \(-?[0-9.]+, (-?[0-9.]+)\)`), p5y, 5e-4)}, {"p6 is the analytic intersection of L1 and the perpendicular through p4", close(val(`p6\s+= \((-?[0-9.]+),`), p6x, 5e-4) && close(val(`p6\s+= \(-?[0-9.]+, (-?[0-9.]+)\)`), p6y, 5e-4)}, {"d53 recomputes as the Euclidean distance from p5 to p3", close(val(`d53 = ([0-9.]+) cm`), d53, 5e-6)}, {"d64 recomputes as the Euclidean distance from p6 to p4", close(val(`d64 = ([0-9.]+) cm`), d64, 5e-6)}, {"dCm recomputes as d53 minus d64", close(val(`dCm = (-?[0-9.]+) cm`), dcm, 5e-6)}, {"the discrepancy is finite and negative for this geometry", !math.IsInf(dcm, 0) && !math.IsNaN(dcm) && dcm < 0}, {"the alarm condition follows from |dCm| > 1.25 cm", math.Abs(dcm) > 1.25 && contains(ctx.Answer, "LLD Alarm = TRUE")}}
}
