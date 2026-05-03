package main

import (
	"fmt"
	"strings"
)

func checkHarbor(ctx *Context) []Check {
	d := ctx.M()
	agg := asMap(d["Aggregate"])
	insight := asMap(d["Insight"])
	req := asMap(d["Request"])
	dispatch := asMap(d["Dispatch"])
	policy := asMap(d["Policy"])
	th := asMap(d["Thresholds"])
	duration := float64(parseRFC3339ish(str(dispatch["WindowEnd"]))-parseRFC3339ish(str(dispatch["WindowStart"]))) / 3600.0
	energy := num(dispatch["DispatchMW"]) * duration
	serialized := strings.ToLower(str(insight["SerializedLowercase"]))
	sensitive := []string{"coretemperature", "rodposition", "neutronflux", "operatorbadge"}
	omit := true
	for _, t := range sensitive {
		omit = omit && !contains(serialized, t)
	}
	rawLocal := true
	for _, k := range []string{"ContainsCoreTemperature", "ContainsRodPosition", "ContainsNeutronFlux", "ContainsOperatorBadgeIDs"} {
		rawLocal = rawLocal && !boolean(agg[k])
	}
	perm := asMap(policy["Permission"])
	permit := str(d["RequestPurpose"]) == str(req["Purpose"]) && str(req["Purpose"]) == str(perm["Purpose"]) && str(d["RequestAction"]) == str(perm["Action"]) && str(perm["Target"]) == str(insight["ID"]) && str(req["TargetLoad"]) == str(insight["TargetLoad"]) && str(insight["TargetLoad"]) == str(dispatch["ForLoad"]) && num(req["RequestedMW"]) == num(dispatch["DispatchMW"]) && num(req["RequestedMW"]) <= num(insight["ExportMW"]) && num(insight["ExportMW"]) <= num(agg["AvailableFlexibleExportMW"]) && num(agg["ReserveMarginMW"]) >= num(th["MinReserveMarginMW"]) && num(agg["CoolingMarginPct"]) >= num(th["MinCoolingMarginPct"]) && !boolean(agg["PlannedOutage"]) && parseRFC3339ish(str(d["HubAuthAt"])) < parseRFC3339ish(str(insight["ExpiresAt"]))
	pro := asMap(policy["Prohibition"])
	scope := true
	for _, k := range []string{"ScopeDevice", "ScopeEvent", "WindowStart", "ExpiresAt"} {
		scope = scope && str(insight[k]) != ""
	}
	return []Check{{"reserve margin exceeds the configured threshold", num(agg["ReserveMarginMW"]) >= num(th["MinReserveMarginMW"])}, {"cooling margin exceeds the configured threshold", num(agg["CoolingMarginPct"]) >= num(th["MinCoolingMarginPct"])}, {"no planned outage blocks the balancing window", !boolean(agg["PlannedOutage"])}, {"requested dispatch fits inside the flexible-export insight", num(req["RequestedMW"]) <= num(insight["ExportMW"]) && num(insight["ExportMW"]) <= num(agg["AvailableFlexibleExportMW"])}, {"serialized insight omits sensitive reactor telemetry terms", omit}, {"aggregate flags keep raw reactor telemetry local", rawLocal}, {"permission policy authorizes electrolysis dispatch before expiry", permit && contains(ctx.Answer, "PERMIT")}, {"market-resale redistribution is separately prohibited", str(pro["Action"]) == "odrl:distribute" && str(pro["Purpose"]) == "market_resale"}, {"scope is explicit for device, event, start, and expiry", scope}, {"dispatch energy recomputes to 64 MWh over the four-hour window", energy == 64 && contains(ctx.Reason, "64 MWh")}, {"the reported load, power, and window match the request and dispatch", contains(ctx.Answer, str(req["TargetLoad"])) && contains(ctx.Answer, fmt.Sprintf("at %.0f MW", num(req["RequestedMW"]))) && contains(ctx.Answer, str(dispatch["WindowStart"])) && contains(ctx.Answer, str(dispatch["WindowEnd"]))}}
}
