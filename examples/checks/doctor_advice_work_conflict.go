package main

import (
	"strings"
)

func checkDoctor(ctx *Context) []Check {
	d := ctx.M()
	computed := doctorEval(d)
	reported := parseDoctorRows(ctx.Answer)
	home, office := "", ""
	for _, r := range maps(d["Requests"]) {
		if str(r["Location"]) == "Home" {
			home = str(r["ID"])
		}
		if str(r["Location"]) == "Office" {
			office = str(r["ID"])
		}
	}
	overall := "Other"
	if computed[home]["effective"] == "Permit" && computed[office]["effective"] == "Deny" {
		overall = "RemoteWorkOnly"
	}
	allPerm, allDeny := true, true
	for _, r := range maps(d["Requests"]) {
		allPerm = allPerm && strings.Contains(computed[str(r["ID"])]["raw"], "Permit")
		allDeny = allDeny && strings.Contains(computed[str(r["ID"])]["raw"], "Deny")
	}
	repHome := reported[home]
	repOffice := reported[office]
	return []Check{{"Flu classifies Jos as sick for the policy conflict", str(asMap(d["Person"])["Condition"]) == "Flu" && contains(ctx.Reason, "Jos has Flu")}, {"ProgrammingWork is closed upward to Work", doctorClosure("ProgrammingWork", asMap(d["SubclassOf"]))["Work"]}, {"doctor advice contributes Permit for every ProgrammingWork request", allPerm}, {"sick-work default contributes Deny for both requests", allDeny}, {"the home request keeps the raw Permit+Deny conflict before resolution", repHome["raw"] == computed[home]["raw"] && repHome["status"] == computed[home]["status"] && repHome["effective"] == computed[home]["effective"] && computed[home]["status"] == "BothPermitDeny"}, {"the office request keeps the raw Permit+Deny conflict before resolution", repOffice["raw"] == computed[office]["raw"] && repOffice["status"] == computed[office]["status"] && repOffice["effective"] == computed[office]["effective"] && computed[office]["status"] == "BothPermitDeny"}, {"conflict resolution permits sick programming work at Home", computed[home]["effective"] == "Permit"}, {"conflict resolution denies Office work", computed[office]["effective"] == "Deny"}, {"the combined recommendation recomputes to RemoteWorkOnly", overall == str(asMap(d["Expected"])["OverallDecision"]) && contains(ctx.Answer, "overall decision for "+str(asMap(d["Person"])["Name"])+" : "+overall)}}
}
