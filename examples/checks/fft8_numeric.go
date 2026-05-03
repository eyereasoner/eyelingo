package main

import (
	"math"
	"math/cmplx"
)

func checkFFT8(ctx *Context) []Check {
	d := ctx.M()
	samples := farr(d["samples"])
	spectrum := dftReal(samples)
	mags := make([]float64, len(spectrum))
	maxMag := 0.0
	for i, z := range spectrum {
		mags[i] = cmplx.Abs(z)
		if mags[i] > maxMag {
			maxMag = mags[i]
		}
	}
	tol := num(asMap(d["expected"])["tolerance"])
	dom := []int{}
	for k, mag := range mags {
		if math.Abs(mag-maxMag) <= tol {
			dom = append(dom, k)
		}
	}
	reported := map[int][2]float64{}
	for _, m := range reAll(ctx.Answer, `k=([0-9]+) magnitude=([0-9.]+) phase=([-0-9.]+)`) {
		reported[parseInt(m[1])] = [2]float64{parseFloat(m[2]), parseFloat(m[3])}
	}
	timeEnergy := 0.0
	for _, x := range samples {
		timeEnergy += x * x
	}
	freqEnergy := 0.0
	for _, z := range spectrum {
		freqEnergy += cmplx.Abs(z) * cmplx.Abs(z)
	}
	freqEnergy /= float64(len(samples))
	expDom := []int{}
	for _, v := range asSlice(asMap(d["expected"])["dominantBins"]) {
		expDom = append(expDom, integer(v))
	}
	repOK := len(reported) == len(dom)
	for _, k := range dom {
		r, ok := reported[k]
		repOK = repOK && ok && close(r[0], mags[k], 1e-6) && close(r[1], cmplx.Phase(spectrum[k]), 1e-6)
	}
	conj := true
	for k := 1; k < len(samples); k++ {
		conj = conj && cmplx.Abs(spectrum[k]-cmplx.Conj(spectrum[len(samples)-k])) <= 1e-9
	}
	return []Check{{"the input contains exactly eight samples", len(samples) == 8}, {"dominant bins are recomputed as k=1 and k=7", intsEq(dom, expDom)}, {"the DC component cancels to zero within tolerance", cmplx.Abs(spectrum[0]) <= tol}, {"the two sine bins have magnitude four", allDominantMag(mags, dom, 4.0, 1e-9)}, {"reported dominant magnitudes and phases match recomputation", repOK}, {"Parseval energy is preserved under the unnormalized DFT convention", close(timeEnergy, freqEnergy, 1e-9)}, {"reported time and frequency-domain energies match recomputation", close(fieldFloat(ctx.Answer, "time-domain energy"), timeEnergy, 1e-6) && close(fieldFloat(ctx.Answer, "frequency-domain energy / 8"), freqEnergy, 1e-6)}, {"real samples produce conjugate-symmetric bins", conj}}
}
