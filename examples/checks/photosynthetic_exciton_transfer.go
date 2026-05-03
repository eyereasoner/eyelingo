package main

func checkPhoto(ctx *Context) []Check {
	d := ctx.M()
	comps := map[string]anymap{}
	for _, c := range maps(d["complexes"]) {
		comps[str(c["role"])] = c
	}
	tuned := comps["positive"]
	detuned := comps["contrast"]
	rc := str(d["reactionCenter"])
	tunedOK := photoEff(tuned, rc)
	detunedOK := photoEff(detuned, rc)
	return []Check{{"the tuned complex has strong excitonic coupling", str(tuned["excitonCoupling"]) == "Strong"}, {"the tuned complex has delocalized states", str(tuned["stateDelocalization"]) == "Present"}, {"the tuned complex has a tuned vibronic bridge", str(tuned["vibronicBridge"]) == "Tuned"}, {"the tuned complex has moderate dephasing and short-lived coherence", str(tuned["dephasing"]) == "Moderate" && str(tuned["electronicCoherenceLifetime"]) == "Short"}, {"the tuned complex is connected downhill to the reaction center", str(tuned["energyLandscape"]) == "DownhillToReactionCenter" && str(tuned["connectedTo"]) == rc}, {"efficient transfer is independently derived for the tuned complex", tunedOK && contains(ctx.Answer, "YES for the tuned antenna complex")}, {"the detuned complex lacks strong coupling, delocalization, and a vibronic bridge", str(detuned["excitonCoupling"]) == "Weak" && str(detuned["stateDelocalization"]) == "Absent" && str(detuned["vibronicBridge"]) == "Absent"}, {"the detuned complex has strong dephasing and a trapping mismatch", str(detuned["dephasing"]) == "Strong" && str(detuned["energyLandscape"]) == "TrappingMismatch"}, {"efficient delivery is blocked for the detuned contrast complex", !detunedOK && contains(ctx.Answer, "NO for the detuned")}, {"the reaction-center connection alone is insufficient without the other conditions", str(detuned["connectedTo"]) == rc && !detunedOK}}
}
