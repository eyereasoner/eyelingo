package main

import (
	"math"
)

func checkEbike(ctx *Context) []Check {
	d := ctx.M()
	lower := num(d["InitialMotorC"]) - num(d["AmbientC"])
	upper := lower
	rows := []struct {
		assist string
		lo, up float64
	}{}
	heatmap := asMap(d["HeatingEnvelopeByAssist"])
	for _, a := range sarr(d["AssistPlan"]) {
		h := asMap(heatmap[a])
		lower = lower*num(d["CoolingLower"]) + num(h["Lower"])
		upper = upper*num(d["CoolingUpper"]) + num(h["Upper"])
		rows = append(rows, struct {
			assist string
			lo, up float64
		}{a, num(d["AmbientC"]) + lower, num(d["AmbientC"]) + upper})
	}
	maxUp := -1.0
	first := 0
	rec := 0
	for i, r := range rows {
		if r.up > maxUp {
			maxUp = r.up
		}
		if first == 0 && r.up > num(d["WarningLimitC"]) {
			first = i + 1
		}
		if first > 0 && rec == 0 && i+1 >= first && r.up <= num(d["WarningLimitC"]) {
			rec = i + 1
		}
	}
	hardSafe := true
	for _, r := range rows {
		hardSafe = hardSafe && r.up < num(d["HardLimitC"])
	}
	exact := math.Exp(-num(d["SamplePeriodSec"]) / num(d["ThermalTimeConstantSec"]))
	repMax := fieldFloat(ctx.Answer, "maximum upper motor temperature")
	repRec := grabInt(ctx.Answer, `warning recovery : step (\d+) at`)
	repHard := fieldFloat(ctx.Answer, "hard limit")
	envelopesOK := true
	for _, v := range heatmap {
		h := asMap(v)
		envelopesOK = envelopesOK && num(h["Lower"]) <= num(h["Upper"]) && num(h["Lower"]) >= 0
	}
	cool := true
	for i := 8; i < 12; i++ {
		cool = cool && rows[i].up < rows[i-1].up
	}
	return []Check{{"the cooling certificate brackets exp(-sample/tau)", num(d["CoolingLower"]) <= exact && exact <= num(d["CoolingUpper"])}, {"the assist plan has twelve sampled thermal updates", len(rows) == 12}, {"Turbo, Tour, Eco, and Coast heating envelopes are nonnegative intervals", envelopesOK}, {"interval propagation recomputes the maximum upper motor temperature", close(repMax, maxUp, 5e-4)}, {"the upper envelope first exceeds the warning limit during Turbo", first == 1 && rows[0].assist == "Turbo"}, {"the reported warning-recovery step matches the independently propagated envelope", repRec == rec && rec == integer(asMap(d["Expected"])["WarningRecoveryStep"])}, {"all upper temperatures remain below the hard thermal limit", hardSafe && repHard == num(d["HardLimitC"])}, {"the final Coast samples cool monotonically", cool}, {"the final decision matches the safety envelope", contains(ctx.Answer, str(asMap(d["Expected"])["Decision"])) && hardSafe}}
}
