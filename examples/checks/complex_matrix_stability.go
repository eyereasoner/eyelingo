package main

import (
	"math"
	"math/cmplx"
	"strings"
)

func checkComplexMatrix(ctx *Context) []Check {
	d := ctx.M()
	reported := map[string]struct {
		radius float64
		class  string
	}{}
	for _, m := range reAll(ctx.Answer, `(?m)^([A-Za-z0-9_]+) : spectral radius ([0-9.]+) -> ([^\n]+)$`) {
		reported[m[1]] = struct {
			radius float64
			class  string
		}{parseFloat(m[2]), strings.TrimSpace(m[3])}
	}
	classify := func(radius float64) string {
		if radius < 1.0 {
			return "damped"
		}
		if math.Abs(radius-1.0) <= 1e-12 {
			return "marginally stable"
		}
		return "unstable"
	}
	radii := map[string]float64{}
	for _, mat := range maps(d["matrices"]) {
		name := str(mat["name"])
		mx := 0.0
		for _, z := range maps(mat["diagonal"]) {
			r := math.Hypot(num(z["re"]), num(z["im"]))
			if r > mx {
				mx = r
			}
		}
		radii[name] = mx
	}
	rowsOK := len(reported) == len(radii)
	for name, radius := range radii {
		row, ok := reported[name]
		rowsOK = rowsOK && ok && math.Abs(row.radius-radius) <= 1e-12 && row.class == classify(radius)
	}
	zm := asMap(asMap(d["sampleProduct"])["z"])
	wm := asMap(asMap(d["sampleProduct"])["w"])
	z := complex(num(zm["re"]), num(zm["im"]))
	w := complex(num(wm["re"]), num(wm["im"]))
	scale := num(d["scale"])
	return []Check{{"diagonal entries are used as the eigenvalues", len(radii) == len(maps(d["matrices"]))}, {"A_unstable has independently recomputed spectral radius 2", math.Abs(radii["A_unstable"]-2.0) <= 1e-12 && classify(radii["A_unstable"]) == "unstable"}, {"A_stable has spectral radius exactly 1 and is marginal", math.Abs(radii["A_stable"]-1.0) <= 1e-12 && classify(radii["A_stable"]) == "marginally stable"}, {"A_damped has spectral radius 0 and is damped", radii["A_damped"] == 0.0 && classify(radii["A_damped"]) == "damped"}, {"reported matrix classes and radii match recomputation", rowsOK}, {"squared modulus of z*w equals product of squared moduli", math.Abs(math.Pow(cmplx.Abs(z*w), 2)-math.Pow(cmplx.Abs(z), 2)*math.Pow(cmplx.Abs(w), 2)) <= 1e-12}, {"scaling a matrix by 2 multiplies spectral-radius-squared by 4", math.Abs(math.Pow(scale*radii["A_unstable"], 2)-math.Pow(scale, 2)*math.Pow(radii["A_unstable"], 2)) <= 1e-12}}
}
