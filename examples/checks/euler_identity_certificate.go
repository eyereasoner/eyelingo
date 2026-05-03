package main

import (
	"fmt"
	"math"
)

func checkEuler(ctx *Context) []Check {
	d := ctx.M()
	angle := num(d["angle"])
	terms := integer(d["terms"])
	tol := num(d["tolerance"])
	z := complex(0, angle)
	term := complex(1, 0)
	total := complex(1, 0)
	for k := 1; k < terms; k++ {
		term *= z / complex(float64(k), 0)
		total += term
	}
	residual := math.Hypot(real(total)+1, imag(total))
	rr := fieldFloat(ctx.Answer, "computed real part of exp(iπ)")
	ri := fieldFloat(ctx.Answer, "computed imaginary part of exp(iπ)")
	rres := fieldFloat(ctx.Answer, "residual magnitude")
	return []Check{{"the Taylor expansion uses the terms count from JSON", terms == 28 && contains(ctx.Answer, fmt.Sprintf("terms used : %d", terms))}, {"the input angle is pi to floating precision", math.Abs(angle-math.Pi) <= 1e-15}, {"the finite Taylor real part is close to -1", math.Abs(real(total)+1) <= tol}, {"the finite Taylor imaginary part is close to zero", math.Abs(imag(total)) <= tol}, {"the residual |exp(iπ)+1| is below the configured tolerance", residual < tol && boolean(asMap(d["expected"])["residualBelowTolerance"])}, {"reported real, imaginary, and residual values match recomputation", close(rr, real(total), 1e-12) && close(ri, imag(total), 1e-12) && close(rres, residual, 5e-16)}, {"the explanation explicitly treats the result as a finite certificate", contains(ctx.Reason, "finite") && contains(ctx.Reason, "not claimed") && contains(ctx.Reason, "tolerance")}}
}
