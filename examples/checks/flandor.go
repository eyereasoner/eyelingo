package main

import (
	"strings"
)

func checkFlandor(ctx *Context) []Check {
	d := ctx.M()
	needs := map[string]bool{"export": false, "skills": num(asMap(d["Labour"])["TechVacancyRate"]) > 3.9, "grid": num(asMap(d["Grid"])["CongestionHours"]) > 11}
	for _, c := range maps(d["Clusters"]) {
		if num(c["ExportOrdersIndex"]) < 90 {
			needs["export"] = true
		}
	}
	active := 0
	for _, v := range needs {
		if v {
			active++
		}
	}
	covers := func(p anymap) bool {
		return (!needs["export"] || boolean(p["CoversExportWeakness"])) && (!needs["skills"] || boolean(p["CoversSkillsStrain"])) && (!needs["grid"] || boolean(p["CoversGridStress"]))
	}
	var selected anymap
	for _, p := range maps(d["Packages"]) {
		if covers(p) && num(p["CostMEUR"]) <= num(asMap(d["Budget"])["MaxMEUR"]) {
			if selected == nil || num(p["CostMEUR"]) < num(selected["CostMEUR"]) || (num(p["CostMEUR"]) == num(selected["CostMEUR"]) && str(p["PackageID"]) < str(selected["PackageID"])) {
				selected = p
			}
		}
	}
	repActive := grabInt(ctx.Answer, `Active need count: (\d+)/`)
	repPkg := ""
	if m := reFind(ctx.Answer, `Recommended package: ([^\n]+)`); m != nil {
		repPkg = m[0]
	}
	repCost := grabInt(ctx.Answer, `Package cost: €(\d+)M`)
	repHash := ""
	if m := reFind(ctx.Answer, `Payload SHA-256: ([0-9a-f]+)`); m != nil {
		repHash = m[0]
	}
	repHMAC := ""
	if m := reFind(ctx.Answer, `Envelope HMAC-SHA-256: ([0-9a-f]+)`); m != nil {
		repHMAC = m[0]
	}
	insight := asMap(d["Insight"])
	policy := asMap(d["Policy"])
	sig := asMap(d["Signature"])
	allowed := str(d["RequestAction"]) == str(policy["PermissionAction"]) && str(d["RequestPurpose"]) == str(policy["PermissionPurpose"]) && str(policy["PermissionTarget"]) == str(insight["ID"]) && str(d["BoardAuthAt"]) < str(insight["ExpiresAt"]) && active >= integer(insight["ThresholdScore"]) && selected != nil
	cheaperReject := true
	for _, p := range maps(d["Packages"]) {
		if selected != nil && num(p["CostMEUR"]) < num(selected["CostMEUR"]) {
			cheaperReject = cheaperReject && !covers(p)
		}
	}
	full := maps(d["Packages"])[len(maps(d["Packages"]))-1]
	serialized := strings.ToLower(str(insight["SerializedLowercase"]))
	exportWeak := []string{}
	for _, c := range maps(d["Clusters"]) {
		if num(c["ExportOrdersIndex"]) < 90 {
			exportWeak = append(exportWeak, str(c["Name"]))
		}
	}
	return []Check{{"export weakness, skills strain, and grid stress are all active", needs["export"] && needs["skills"] && needs["grid"] && active == 3}, {"active need count meets the insight threshold", repActive == active && active == integer(insight["ThresholdScore"])}, {"the lowest-cost package covering all active needs is selected", str(selected["Name"]) == "Flandor Retooling Pulse" && repPkg == str(selected["Name"])}, {"the selected package fits inside the €140M budget", integer(selected["CostMEUR"]) == repCost && repCost == 120 && num(selected["CostMEUR"]) <= num(asMap(d["Budget"])["MaxMEUR"])}, {"cheaper packages are rejected because each covers only one active need", cheaperReject}, {"the full corridor package covers all needs but is over budget", covers(full) && num(full["CostMEUR"]) > num(asMap(d["Budget"])["MaxMEUR"])}, {"policy permission authorizes regional-stabilization use before expiry", allowed && contains(ctx.Answer, "decision : ALLOWED")}, {"firm-surveillance redistribution is prohibited", str(policy["ProhibitionAction"]) == "odrl:distribute" && str(policy["ProhibitionPurpose"]) == "firm_surveillance"}, {"deletion duty is scheduled before envelope expiry", str(d["BoardDutyAt"]) < str(insight["ExpiresAt"])}, {"shared insight omits firm names and payroll rows", !boolean(asMap(d["Signals"])["ContainsFirmNames"]) && !boolean(asMap(d["Signals"])["ContainsPayrollRows"]) && !contains(serialized, "firm") && !contains(serialized, "payroll")}, {"reported signature metadata matches the trusted precomputed input", repHash == str(sig["DisplayPayloadSHA256"]) && repHMAC == str(sig["SignatureHMAC"]) && str(sig["Algorithm"]) == "HMAC-SHA256"}, {"the expected six files and one audit entry are recorded", integer(d["FilesWritten"]) == integer(d["ExpectedFilesWritten"]) && integer(d["FilesWritten"]) == 6 && integer(d["AuditEntries"]) == 1}, {"export-weak cluster names are independently identified", sliceEq(exportWeak, []string{"Antwerp chemicals", "Ghent manufacturing"})}}
}
