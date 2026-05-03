package main

import (
	"fmt"
	"sort"
	"strings"
)

func checkDPP(ctx *Context) []Check {
	d := ctx.M()
	total, recycled := 0, 0
	for _, c := range maps(d["Components"]) {
		total += integer(c["MassG"])
		recycled += integer(c["RecycledMassG"])
	}
	pct := 0
	if total > 0 {
		pct = recycled * 100 / total
	}
	fp := asMap(d["Footprint"])
	lifecycle := integer(fp["ManufacturingGCO2e"]) + integer(fp["TransportGCO2e"]) + integer(fp["UsePhaseGCO2e"])
	critMat := map[string]bool{}
	for _, m := range maps(d["Materials"]) {
		critMat[str(m["ID"])] = boolean(m["CriticalRawMaterial"])
	}
	crit := []string{}
	seen := map[string]bool{}
	battery := false
	for _, c := range maps(d["Components"]) {
		if strings.ToLower(str(c["Type"])) == "battery" && boolean(c["Replaceable"]) {
			battery = true
		}
		for _, mat := range sarr(c["ContainsMaterial"]) {
			if critMat[mat] && !seen[mat] {
				seen[mat] = true
				crit = append(crit, mat)
			}
		}
	}
	sort.Strings(crit)
	ap := asMap(d["AccessPolicy"])
	pubSec, restSec := str(ap["PublicSection"]), str(ap["RestrictedSection"])
	pubDocs, restDocs := map[string]bool{}, map[string]bool{}
	declares := func(section, claim string) bool {
		for _, doc := range maps(d["Documents"]) {
			if str(doc["Section"]) == section {
				for _, cl := range sarr(doc["Declares"]) {
					if cl == claim {
						return true
					}
				}
			}
		}
		return false
	}
	for _, doc := range maps(d["Documents"]) {
		if str(doc["Section"]) == pubSec {
			pubDocs[str(doc["DocType"])] = true
		}
		if str(doc["Section"]) == restSec {
			restDocs[str(doc["DocType"])] = true
		}
	}
	repair := battery && pubDocs["RepairGuide"] && pubDocs["SparePartsCatalog"] && declares(pubSec, str(ap["RequiredPublicClaim"]))
	pubTypesOK, restTypesOK := true, true
	for _, t := range sarr(ap["PublicDocTypes"]) {
		pubTypesOK = pubTypesOK && pubDocs[t] && !restDocs[t]
	}
	for _, t := range sarr(ap["RestrictedDocTypes"]) {
		restTypesOK = restTypesOK && restDocs[t] && !pubDocs[t]
	}
	orderOK := true
	lastDay := ""
	lastOrder := -1
	ord := map[string]int{"ManufacturingEvent": 0, "SaleEvent": 1, "RepairEvent": 2}
	for _, e := range maps(d["Lifecycle"]) {
		day := str(e["OnDate"])
		eo := ord[str(e["Type"])]
		if lastDay != "" && day < lastDay {
			orderOK = false
		}
		if eo < lastOrder {
			orderOK = false
		}
		lastDay = day
		lastOrder = eo
	}
	ansCrit := answerField(ctx.Answer, "critical raw materials")
	decision := answerField(ctx.Answer, "Passport decision")
	exp := asMap(d["Expected"])
	return []Check{{"component masses are folded from the JSON component list", fieldInt(ctx.Answer, "total component mass") == total && total == integer(exp["TotalMassG"])}, {"recycled mass and integer recycled-content percentage are recomputed independently", recycled == integer(exp["RecycledMassG"]) && fieldInt(ctx.Answer, "recycled content") == pct && pct == integer(exp["RecycledContentPct"])}, {"manufacturing, transport, and use-phase footprint values sum to the reported lifecycle footprint", fieldInt(ctx.Answer, "lifecycle footprint") == lifecycle && lifecycle == integer(exp["LifecycleGCO2e"])}, {"critical raw-material exposure is derived by joining component materials to material declarations", ansCrit == strings.Join(crit, ", ") && sliceEq(crit, []string{"Cobalt", "Lithium"})}, {"repairFriendly is derived from a replaceable battery plus public repair and spare-parts documents", repair && answerField(ctx.Answer, "circularity hint") == str(exp["CircularityHint"])}, {"all required public document types are present only in the public document section", pubTypesOK}, {"all restricted declaration document types stay in the restricted section", restTypesOK}, {"lifecycle events are chronological and follow manufacturing before sale before repair", orderOK}, {"the passport endpoint equals the product digital link and is reported publicly", str(asMap(d["Passport"])["PublicEndpoint"]) == str(asMap(d["Product"])["DigitalLink"]) && str(asMap(d["Product"])["DigitalLink"]) == answerField(ctx.Answer, "public endpoint")}, {"PASS is reported only because every independent passport check succeeds", decision == fmt.Sprintf("PASS for %s %s.", str(asMap(d["Product"])["Model"]), str(asMap(d["Product"])["SerialNumber"]))}}
}
