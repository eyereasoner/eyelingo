package main

func checkAuroraCare(ctx *Context) []Check {
	d := ctx.M()
	computed := auroraEvaluate(d)
	reported := auroraParseRows(ctx.Answer)
	permit, deny := 0, 0
	for _, r := range computed {
		if r.decision == "PERMIT" {
			permit++
		} else if r.decision == "DENY" {
			deny++
		}
	}
	match := len(computed) == len(reported)
	for k, r := range computed {
		match = match && reported[k].decision == r.decision
	}
	return []Check{{"all seven scenario decisions are recomputed", len(computed) == 7 && len(reported) == 7}, {"reported decisions match the independent PDP evaluation", match}, {"primary-care access requires clinician role and care-team link", computed["A"].decision == "PERMIT" && computed["A"].uid == "urn:policy:primary-care-001" && computed["E"].decision == "PERMIT" && computed["E"].uid == "urn:policy:primary-care-001"}, {"quality improvement is allowed only with both required categories in the secure environment", computed["B"].decision == "PERMIT" && computed["B"].uid == "urn:policy:qi-2025-aurora" && computed["C"].decision == "DENY"}, {"insurance management is denied by the matching prohibition", computed["D"].decision == "DENY" && computed["D"].uid == "urn:policy:deny-insurance"}, {"research requires opt-in, anonymisation, and the secure environment", computed["F"].decision == "PERMIT" && computed["F"].uid == "urn:policy:research-aurora-diabetes"}, {"AI training is denied because the subject opted out", computed["G"].decision == "DENY" && computed["G"].uid == "" && computed["G"].why == "subject_opted_out"}, {"permit and deny counts match the report", permit == 4 && deny == 3 && contains(ctx.Answer, "permit count : 4") && contains(ctx.Answer, "deny count : 3")}}
}
