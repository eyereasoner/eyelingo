package main

import (
	"fmt"
	"strings"
)

func checkCalidor(ctx *Context) []Check {
	d := ctx.M()
	needs := calidorActiveNeeds(d)
	active := 0
	for _, v := range needs {
		if v {
			active++
		}
	}
	required := []string{"bill_credit", "cooling_kit", "transport", "welfare_check"}
	selected := chooseCapabilityPackage(maps(d["Packages"]), num(d["MaxPackageCostEur"]), required)
	reported := parseCalidorAnswer(ctx.Answer)
	insight := asMap(d["Insight"])
	policy := asMap(d["Policy"])
	sig := asMap(d["Signature"])
	allowed := str(d["RequestAction"]) == str(policy["PermissionAction"]) &&
		str(d["RequestPurpose"]) == str(policy["PermissionPurpose"]) &&
		str(policy["PermissionTarget"]) == str(insight["ID"]) &&
		parseTime(str(d["CityAuthAt"])).Before(parseTime(str(d["GatewayExpiresAt"]))) &&
		active >= integer(d["MinimumActiveNeedCount"]) &&
		selected != nil
	deleteBefore := parseTime(str(d["CityDutyAt"])).Before(parseTime(str(d["GatewayExpiresAt"])))
	serialized := strings.ToLower(str(insight["SerializedLowercase"]))
	localTerms := []string{"heat_sensitive_condition", "mobility_limitation", "prepaid"}
	omitted := true
	for _, t := range localTerms {
		omitted = omitted && !contains(serialized, t)
	}
	cheaperReject := true
	deluxeOverBudget := false
	for _, pkg := range maps(d["Packages"]) {
		if selected != nil && integer(pkg["CostEur"]) < integer(selected["CostEur"]) {
			cheaperReject = cheaperReject && !setContainsAll(setOf(sarr(pkg["Capabilities"])), required)
		}
		if str(pkg["PackageID"]) == "pkg:DELUXE" {
			deluxeOverBudget = integer(pkg["CostEur"]) > integer(d["MaxPackageCostEur"])
		}
	}
	scopeOK := true
	for _, k := range []string{"ScopeDevice", "ScopeEvent", "Municipality", "CreatedAt", "ExpiresAt"} {
		scopeOK = scopeOK && str(insight[k]) != ""
	}
	selectedOK := selected != nil && str(selected["Name"]) == "Calidor Priority Cooling Bundle"
	return []Check{
		{"four active heat-response needs are recomputed from the local signals", active == 4 && reported["active"] == "4"},
		{"the insight threshold of three active needs is met", active >= integer(d["MinimumActiveNeedCount"]) && integer(d["MinimumActiveNeedCount"]) == integer(insight["ThresholdCount"])},
		{"the lowest-cost eligible package covering all required capabilities is selected", selectedOK && reported["pkg"] == str(selected["Name"])},
		{"the selected package fits inside the €20 budget", selected != nil && reported["cost"] == fmt.Sprint(integer(selected["CostEur"])) && integer(selected["CostEur"]) == 18 && integer(selected["CostEur"]) <= integer(d["MaxPackageCostEur"])},
		{"cheaper packages are rejected because they do not cover all capabilities", cheaperReject},
		{"the deluxe package is rejected because it is over budget", deluxeOverBudget},
		{"policy permission authorizes heatwave-response use before expiry", allowed && contains(ctx.Answer, "decision : ALLOWED")},
		{"tenant-screening reuse is prohibited by the policy", str(policy["ProhibitionAction"]) == "odrl:distribute" && str(policy["ProhibitionPurpose"]) == "tenant_screening"},
		{"deletion duty is scheduled before envelope expiry", deleteBefore},
		{"vulnerability flags and local raw stress signals are omitted from the serialized insight", omitted},
		{"reported signature metadata matches the trusted precomputed input", reported["hash"] == str(sig["PayloadHashSHA256"]) && reported["hmac"] == str(sig["SignatureHMAC"]) && str(sig["Algorithm"]) == "HMAC-SHA256"},
		{"scope metadata is explicit for device, event, municipality, creation, and expiry", scopeOK},
	}
}
