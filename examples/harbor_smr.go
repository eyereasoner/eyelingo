// harbor_smr.go
//
// A self-contained Go translation inspired by Eyeling's `examples/harborsmr.n3`.
//
// The scenario models a port hydrogen hub that asks to use a narrow,
// permissioned SMR flexible-export insight for one four-hour electrolyzer
// dispatch window. The key point is data minimization: the hub receives only a
// bounded decision object, while raw reactor telemetry stays local to the SMR
// operator.
//
// This is intentionally not a generic RDF/N3/ODRL reasoner. The concrete facts
// and rules are represented as Go structs and ordinary functions so the policy,
// safety, and minimization checks stay visible and directly runnable.
//
// Run:
//
//	go run harbor_smr.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

const eyelingoExampleName = "harbor_smr"

// Dataset is the complete concrete fixture: case context, the narrow insight,
// the aggregate safety facts, the policy, and the dispatch request.
type Dataset struct {
	CaseName       string
	Question       string
	RequestPurpose string
	RequestAction  string
	HubAuthAt      string
	Hub            Hub
	Aggregate      SecureAggregate
	Insight        Insight
	Request        ElectrolyzerRequest
	Dispatch       DispatchPlan
	Policy         Policy
	Thresholds     Thresholds
}

type Hub struct {
	ID      string
	Name    string
	Country string
}

type SecureAggregate struct {
	ID                        string
	AvailableFlexibleExportMW int
	ReserveMarginMW           int
	CoolingMarginPct          int
	PlannedOutage             bool
	ContainsCoreTemperature   bool
	ContainsRodPosition       bool
	ContainsNeutronFlux       bool
	ContainsOperatorBadgeIDs  bool
}

type Insight struct {
	ID                  string
	Metric              string
	TargetLoad          string
	ExportMW            int
	WindowStart         string
	ExpiresAt           string
	ScopeDevice         string
	ScopeEvent          string
	SerializedLowercase string
}

type ElectrolyzerRequest struct {
	Requester   string
	RequestedMW int
	Purpose     string
	TargetLoad  string
}

type DispatchPlan struct {
	DispatchMW  int
	WindowStart string
	WindowEnd   string
	ForLoad     string
}

type Policy struct {
	Profile     string
	Permission  PolicyRule
	Prohibition PolicyRule
}

type PolicyRule struct {
	Action  string
	Target  string
	Purpose string
}

type Thresholds struct {
	MinReserveMarginMW  int
	MinCoolingMarginPct int
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Decision           string
	Checks             []Check
	SensitiveTerms     []string
	ParsedHubAuthAt    time.Time
	ParsedInsightStart time.Time
	ParsedInsightEnd   time.Time
	ParsedDispatchFrom time.Time
	ParsedDispatchTo   time.Time
	EnergyMWh          int
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
	printChecks(analysis)
	printAudit(ds, analysis)
	if !allChecksOK(analysis.Checks) {
		os.Exit(1)
	}
}

func derive(ds Dataset) Analysis {
	authAt := mustParseTime("hubAuthAt", ds.HubAuthAt)
	insightStart := mustParseTime("insight.windowStart", ds.Insight.WindowStart)
	insightEnd := mustParseTime("insight.expiresAt", ds.Insight.ExpiresAt)
	dispatchFrom := mustParseTime("dispatch.windowStart", ds.Dispatch.WindowStart)
	dispatchTo := mustParseTime("dispatch.windowEnd", ds.Dispatch.WindowEnd)
	sensitiveTerms := findSensitiveTerms(ds.Insight.SerializedLowercase)
	energyMWh := ds.Dispatch.DispatchMW * int(dispatchTo.Sub(dispatchFrom).Hours())

	checks := []Check{
		{
			ID:   "C1",
			OK:   ds.Aggregate.ReserveMarginMW > ds.Thresholds.MinReserveMarginMW,
			Text: fmt.Sprintf("reserve margin %d MW exceeds threshold %d MW", ds.Aggregate.ReserveMarginMW, ds.Thresholds.MinReserveMarginMW),
		},
		{
			ID:   "C2",
			OK:   ds.Aggregate.CoolingMarginPct > ds.Thresholds.MinCoolingMarginPct,
			Text: fmt.Sprintf("cooling margin %d%% exceeds threshold %d%%", ds.Aggregate.CoolingMarginPct, ds.Thresholds.MinCoolingMarginPct),
		},
		{
			ID:   "C3",
			OK:   !ds.Aggregate.PlannedOutage,
			Text: "no planned outage blocks the balancing window",
		},
		{
			ID:   "C4",
			OK:   ds.Request.RequestedMW <= ds.Insight.ExportMW,
			Text: fmt.Sprintf("requested %d MW fits inside the %d MW flexible-export insight", ds.Request.RequestedMW, ds.Insight.ExportMW),
		},
		{
			ID:   "C5",
			OK:   len(sensitiveTerms) == 0,
			Text: "serialized insight omits sensitive telemetry terms",
		},
		{
			ID: "C6",
			OK: !ds.Aggregate.ContainsCoreTemperature && !ds.Aggregate.ContainsRodPosition &&
				!ds.Aggregate.ContainsNeutronFlux && !ds.Aggregate.ContainsOperatorBadgeIDs,
			Text: "aggregate flags confirm raw reactor telemetry stays local",
		},
		{
			ID: "C7",
			OK: ds.Policy.Permission.Action == ds.RequestAction &&
				ds.Policy.Permission.Target == ds.Insight.ID &&
				ds.Policy.Permission.Purpose == ds.RequestPurpose &&
				!authAt.After(insightEnd),
			Text: "policy permits use for electrolysis dispatch before the insight expires",
		},
		{
			ID: "C8",
			OK: ds.Policy.Prohibition.Action == "odrl:distribute" &&
				ds.Policy.Prohibition.Target == ds.Insight.ID &&
				ds.Policy.Prohibition.Purpose == "market_resale",
			Text: "policy prohibits redistribution for market resale",
		},
		{
			ID: "C9",
			OK: ds.Insight.ScopeDevice != "" && ds.Insight.ScopeEvent != "" &&
				!insightStart.IsZero() && !insightEnd.IsZero(),
			Text: "scope is explicit: device, event, start, and expiry",
		},
		{
			ID: "C10",
			OK: ds.Dispatch.ForLoad == ds.Request.TargetLoad &&
				ds.Request.RequestedMW <= ds.Dispatch.DispatchMW &&
				!dispatchFrom.Before(insightStart) && !dispatchTo.After(insightEnd),
			Text: "dispatch plan matches the requested load, power, and insight window",
		},
	}

	decision := "DENY"
	if allChecksOK(checks) {
		decision = "PERMIT"
	}

	return Analysis{
		Decision:           decision,
		Checks:             checks,
		SensitiveTerms:     sensitiveTerms,
		ParsedHubAuthAt:    authAt,
		ParsedInsightStart: insightStart,
		ParsedInsightEnd:   insightEnd,
		ParsedDispatchFrom: dispatchFrom,
		ParsedDispatchTo:   dispatchTo,
		EnergyMWh:          energyMWh,
	}
}

func mustParseTime(label, value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse %s=%q as RFC3339: %v\n", label, value, err)
		os.Exit(1)
	}
	return parsed
}

func findSensitiveTerms(s string) []string {
	terms := []string{"coretemp", "rodposition", "neutronflux", "operatorbadge"}
	lower := strings.ToLower(s)
	var hits []string
	for _, term := range terms {
		if strings.Contains(lower, term) {
			hits = append(hits, term)
		}
	}
	sort.Strings(hits)
	return hits
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func printAnswer(ds Dataset, analysis Analysis) {
	fmt.Println("# HarborSMR Insight Dispatch")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("%s - %s may use %s to run %s at %d MW from %s to %s.\n",
		analysis.Decision,
		ds.Hub.Name,
		ds.Insight.ID,
		ds.Request.TargetLoad,
		ds.Dispatch.DispatchMW,
		ds.Dispatch.WindowStart,
		ds.Dispatch.WindowEnd,
	)
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason why")
	fmt.Printf("The SMR operator exposes a bounded %d MW flexible-export insight for %s, not raw reactor telemetry.\n",
		ds.Insight.ExportMW,
		ds.Insight.ScopeEvent,
	)
	fmt.Printf("The requested %d MW electrolysis dispatch fits inside that window, safety margins clear the thresholds, no outage is planned, and the policy permits use only for %s while forbidding market resale distribution.\n",
		ds.Request.RequestedMW,
		ds.RequestPurpose,
	)
	fmt.Printf("The approved dispatch is %d MWh over the four-hour window, scoped to %s and %s.\n",
		analysis.EnergyMWh,
		ds.Insight.ScopeDevice,
		ds.Insight.TargetLoad,
	)
	fmt.Println()
}

func printChecks(analysis Analysis) {
	fmt.Println("## Check")
	for _, check := range analysis.Checks {
		status := "FAIL"
		if check.OK {
			status = "OK"
		}
		fmt.Printf("%s %s - %s\n", check.ID, status, check.Text)
	}
	fmt.Println()
}

func printAudit(ds Dataset, analysis Analysis) {
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("hub : %s (%s, %s)\n", ds.Hub.ID, ds.Hub.Name, ds.Hub.Country)
	fmt.Printf("aggregate : available=%dMW reserve=%dMW cooling=%d%% outage=%v\n",
		ds.Aggregate.AvailableFlexibleExportMW,
		ds.Aggregate.ReserveMarginMW,
		ds.Aggregate.CoolingMarginPct,
		ds.Aggregate.PlannedOutage,
	)
	fmt.Printf("insight : id=%s metric=%s export=%dMW target=%s\n",
		ds.Insight.ID,
		ds.Insight.Metric,
		ds.Insight.ExportMW,
		ds.Insight.TargetLoad,
	)
	fmt.Printf("scope : device=%s event=%s start=%s expires=%s\n",
		ds.Insight.ScopeDevice,
		ds.Insight.ScopeEvent,
		ds.Insight.WindowStart,
		ds.Insight.ExpiresAt,
	)
	fmt.Printf("request : action=%s purpose=%s requested=%dMW target=%s authAt=%s\n",
		ds.RequestAction,
		ds.RequestPurpose,
		ds.Request.RequestedMW,
		ds.Request.TargetLoad,
		ds.HubAuthAt,
	)
	fmt.Printf("dispatch : %dMW %s to %s load=%s energy=%dMWh\n",
		ds.Dispatch.DispatchMW,
		ds.Dispatch.WindowStart,
		ds.Dispatch.WindowEnd,
		ds.Dispatch.ForLoad,
		analysis.EnergyMWh,
	)
	fmt.Printf("policy : profile=%s permission=(%s target=%s purpose=%s) prohibition=(%s target=%s purpose=%s)\n",
		ds.Policy.Profile,
		ds.Policy.Permission.Action,
		ds.Policy.Permission.Target,
		ds.Policy.Permission.Purpose,
		ds.Policy.Prohibition.Action,
		ds.Policy.Prohibition.Target,
		ds.Policy.Prohibition.Purpose,
	)
	fmt.Printf("privacy flags : coreTemperature=%v rodPosition=%v neutronFlux=%v operatorBadgeIDs=%v\n",
		ds.Aggregate.ContainsCoreTemperature,
		ds.Aggregate.ContainsRodPosition,
		ds.Aggregate.ContainsNeutronFlux,
		ds.Aggregate.ContainsOperatorBadgeIDs,
	)
	if len(analysis.SensitiveTerms) == 0 {
		fmt.Println("serialized sensitive term hits : none")
	} else {
		fmt.Printf("serialized sensitive term hits : %s\n", strings.Join(analysis.SensitiveTerms, ", "))
	}
	fmt.Printf("checks passed : %d/%d\n", countChecksOK(analysis.Checks), len(analysis.Checks))
	fmt.Printf("decision : %s\n", analysis.Decision)
}

func countChecksOK(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}
