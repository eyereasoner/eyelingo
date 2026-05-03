package main

func checkODRL(ctx *Context) []Check {
	d := ctx.M()
	rows := odrlRiskRows(d)
	reported := []riskRow{}
	for _, m := range reAll(ctx.Answer, `score=(\d+) \(risk:([A-Za-z]+), risk:[A-Za-z]+\) clause (C\d)`) {
		reported = append(reported, riskRow{m[3], parseInt(m[1]), m[2], 0})
	}
	reportOK := len(reported) == len(rows)
	for i := range reported {
		if i < len(rows) {
			reportOK = reportOK && reported[i].clause == rows[i].clause && reported[i].score == rows[i].score
		}
	}
	scoresDesc := true
	for i := 1; i < len(rows); i++ {
		scoresDesc = scoresDesc && rows[i-1].score >= rows[i].score
	}
	mitig := len(reAll(ctx.Answer, `(?m)^mitigation for clause`))
	high, mod := 0, 0
	totalMit := 0
	for _, r := range rows {
		if r.level == "HighRisk" {
			high++
		}
		if r.level == "ModerateRisk" {
			mod++
		}
		totalMit += r.mitigations
	}
	change := asMap(asMap(asMap(asMap(d["Agreement"])["Policy"])["Permissions"])[":PermChangeTerms"])
	got, _ := odrlNoticeDays(change)
	need := integer(asMap(asMap(asMap(d["Consumer"])["Needs"])[":Need_ChangeOnlyWithPriorNotice"])["MinNoticeDays"])
	return []Check{{"four risk rows are derived from the policy/profile conflict scan", len(rows) == 4 && len(reported) == 4}, {"reported rows match independently computed clauses and scores", reportOK}, {"ranked output is in descending score order", scoresDesc}, {"account/data removal without notice safeguards is the highest risk", rows[0].clause == "C1" && rows[0].score == 100}, {"user-data sharing without explicit consent is scored as high risk", anyRisk(rows, "C3", 97, "HighRisk")}, {"three-day terms-change notice is below the fourteen-day consumer requirement", got == 3 && need == 14}, {"data-export prohibition creates the portability risk row", anyRisk(rows, "C4", 70, "ModerateRisk")}, {"risk level counts recompute to high=3, moderate=1, low=0", high == 3 && mod == 1}, {"five mitigation measures are generated", mitig == totalMit && totalMit == 5}}
}
