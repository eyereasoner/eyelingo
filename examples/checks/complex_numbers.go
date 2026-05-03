package main

import (
	"math"
	"math/cmplx"
)

func checkComplexNumbers(ctx *Context) []Check {
	d := ctx.M()
	expected := map[string]complex128{}
	recomputed := map[string]complex128{}
	angles := map[string]float64{}
	for _, it := range maps(d["Exponents"]) {
		id := str(it["ID"])
		base := toComplex(asMap(it["Base"]))
		pow := toComplex(asMap(it["Power"]))
		expected[id] = toComplex(asMap(it["Expected"]))
		recomputed[id] = cmplx.Pow(base, pow)
		angles[id] = cmplx.Phase(base)
	}
	for _, it := range maps(d["Inverses"]) {
		id := str(it["ID"])
		z := toComplex(asMap(it["Input"]))
		expected[id] = toComplex(asMap(it["Expected"]))
		if str(it["Operation"]) == "asin" {
			recomputed[id] = cmplx.Asin(z)
		} else {
			recomputed[id] = cmplx.Acos(z)
		}
	}
	reported := parseComplexAnswer(ctx.Answer)
	asinv := recomputed["C5"]
	acosv := recomputed["C6"]
	expOK := true
	for _, it := range maps(d["Exponents"]) {
		id := str(it["ID"])
		expOK = expOK && complexClose(recomputed[id], expected[id], 5e-12)
	}
	reportOK := len(reported) == len(recomputed)
	for k, v := range reported {
		reportOK = reportOK && complexClose(v, recomputed[k], 5e-12)
	}
	finite := true
	for _, z := range recomputed {
		finite = finite && !math.IsNaN(real(z)) && !math.IsNaN(imag(z)) && !math.IsInf(real(z), 0) && !math.IsInf(imag(z), 0)
	}
	return []Check{{"principal polar angles for -1, e, and i match the N3 dial cases", math.Abs(angles["C1"]-math.Pi) <= 1e-12 && math.Abs(angles["C2"]) <= 1e-12 && math.Abs(angles["C3"]-math.Pi/2) <= 1e-12}, {"all four complex exponentiation answers match independent complex arithmetic", expOK}, {"i^i and e^(-pi/2) recompute to the same real value", cmplx.Abs(recomputed["C3"]-recomputed["C4"]) <= 1e-12 && math.Abs(imag(recomputed["C3"])) <= 1e-12}, {"asin(2) and acos(2) match independent inverse-trig recomputation", complexClose(recomputed["C5"], expected["C5"], 5e-12) && complexClose(recomputed["C6"], expected["C6"], 5e-12)}, {"sin(asin(2)) and cos(acos(2)) round-trip to 2+0i", cmplx.Abs(cmplx.Sin(asinv)-2) <= 1e-12 && cmplx.Abs(cmplx.Cos(acosv)-2) <= 1e-12}, {"asin(2) + acos(2) equals pi/2 with cancelling imaginary parts", math.Abs(real(asinv+acosv)-math.Pi/2) <= 1e-12 && math.Abs(imag(asinv+acosv)) <= 1e-12}, {"all six reported complex outputs match recomputation to displayed precision", reportOK}, {"the report contains four exponentiation and two inverse-trig queries", len(maps(d["Exponents"])) == 4 && len(maps(d["Inverses"])) == 2 && contains(ctx.Reason, "primitive test facts : 6")}, {"all recomputed outputs are finite real/imaginary pairs", finite}}
}
