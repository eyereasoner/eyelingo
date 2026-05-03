package main

import (
	"fmt"
	"math/cmplx"
	"strings"
)

func checkFFT32(ctx *Context) []Check {
	d := ctx.M()
	n := integer(d["length"])
	tol := num(d["tolerance"])
	rows := parseFFT32Rows(ctx.Answer)
	spectra := map[string][]complex128{}
	samples := map[string][]float64{}
	domOK, magOK, energyOK, conjOK, reportOK := true, true, true, true, true
	for _, sp := range maps(d["waveforms"]) {
		name := str(sp["name"])
		samples[name] = fftWave(sp, n)
		spectra[name] = dftReal(samples[name])
		mags := []float64{}
		for _, z := range spectra[name] {
			mags = append(mags, cmplx.Abs(z))
		}
		dom := dominantBins(spectra[name], tol)
		expDom := []int{}
		for _, v := range asSlice(sp["expectedDominantBins"]) {
			expDom = append(expDom, integer(v))
		}
		domOK = domOK && intsEq(dom, expDom)
		for _, k := range dom {
			magOK = magOK && close(mags[k], num(sp["expectedDominantMagnitude"]), 1e-8)
		}
		te := 0.0
		for _, x := range samples[name] {
			te += x * x
		}
		fe := 0.0
		for _, z := range spectra[name] {
			fe += cmplx.Abs(z) * cmplx.Abs(z)
		}
		fe /= float64(n)
		energyOK = energyOK && close(te, fe, 1e-8)
		for k := 1; k < n; k++ {
			conjOK = conjOK && cmplx.Abs(spectra[name][k]-cmplx.Conj(spectra[name][n-k])) <= 1e-8
		}
		row, ok := rows[name]
		if boolean(sp["expectedFlatSpectrum"]) {
			reportOK = reportOK && ok && strings.Contains(row, "all 32 bins magnitude")
		} else {
			for _, k := range expDom {
				reportOK = reportOK && ok && strings.Contains(row, fmt.Sprintf("k=%d", k))
			}
		}
	}
	impulseOK := true
	for _, z := range spectra["impulse"] {
		impulseOK = impulseOK && close(cmplx.Abs(z), 1.0, 1e-10)
	}
	lenOK := true
	for _, s := range samples {
		lenOK = lenOK && len(s) == n
	}
	return []Check{{"each generated waveform has exactly 32 samples", lenOK}, {"dominant bins match the fixture expectations", domOK}, {"dominant magnitudes match the configured certificates", magOK}, {"Parseval energy is preserved for every waveform", energyOK}, {"real-valued waveforms produce conjugate-symmetric spectra", conjOK}, {"the impulse waveform has a flat unit-magnitude spectrum", impulseOK}, {"the answer reports every expected dominant bin by waveform name", reportOK}, {"the configured operation certificate matches six full 32-bin direct DFTs", integer(d["expectedOperations"]) == len(maps(d["waveforms"]))*n*n}}
}
