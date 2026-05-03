package main

import (
	"fmt"
)

func checkWind(ctx *Context) []Check {
	d := ctx.M()
	cut, rated, cutout := num(d["cutInMS"]), num(d["ratedMS"]), num(d["cutOutMS"])
	rp := num(d["ratedPowerMW"])
	hours := num(d["intervalMinutes"]) / 60
	speeds := farr(d["windSpeedsMS"])
	classes := []string{}
	powers := []float64{}
	for _, v := range speeds {
		kind, p := windPower(v, cut, rated, cutout, rp)
		classes = append(classes, kind)
		powers = append(powers, p)
	}
	usable, ratedCount, stopped := 0, 0, 0
	energy := 0.0
	for i, k := range classes {
		if k != "stopped" {
			usable++
		}
		if k == "rated" {
			ratedCount++
		}
		if k == "stopped" {
			stopped++
		}
		energy += powers[i] * hours
	}
	exp := asMap(d["expected"])
	linesOK := true
	for i, v := range speeds {
		linesOK = linesOK && contains(ctx.Answer, fmt.Sprintf("t%d %.1f m/s %s %.3f MW", i+1, v, classes[i], powers[i]))
	}
	return []Check{{"cut-in, rated, and cut-out thresholds are strictly ordered", cut < rated && rated < cutout}, {"usable intervals are exactly the samples inside the operating envelope", usable == integer(exp["usableIntervals"])}, {"rated intervals are speeds at or above rated and below cut-out", ratedCount == integer(exp["ratedIntervals"])}, {"stopped intervals are below cut-in or at/above cut-out", stopped == integer(exp["stoppedIntervals"])}, {"below-rated usable speeds follow the cubic normalized power curve", close(powers[1], 0.44008604702915216, 1e-9) && close(powers[2], 2.586496313329871, 1e-9)}, {"total interval energy is recomputed in MWh", close(energy, 1.5710970600598373, 1e-9)}, {"reported usable count and total energy match recomputation", fieldInt(ctx.Answer, "usable intervals") == usable && close(fieldFloat(ctx.Answer, "total energy"), energy, 5e-4)}, {"the answer reports every sample classification and power", linesOK}}
}
