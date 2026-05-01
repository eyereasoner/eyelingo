// calidor.go
//
// A self-contained Go translation of the Eyeling Calidor heat-response
// insight-economy example.
//
// The source N3 file models an urgent city service without making the city
// collect raw household sensor, health, or prepaid-energy data. A local gateway
// turns those private signals into a narrow insight: this household needs
// priority cooling support during a heat-alert window. The city may use that
// signed and expiring envelope for heatwave response, but not for unrelated
// purposes such as tenant screening.
//
// This program is intentionally not a generic RDF, ODRL, crypto, or N3 reasoner.
// It translates the concrete facts and rules for this one example into ordinary
// local, while a minimal decision object is shared for a specific public-service
// purpose.
//
// Run:
//
//	go run calidor.go
//
// The program has no third-party dependencies.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"sort"
	"strings"
	"time"
)

const eyelingoExampleName = "calidor"

const sourceFile = "calidor.n3"

const (
	PurposeHeatwaveResponse = "heatwave_response"
	PurposeTenantScreening  = "tenant_screening"

	ActionUse        = "odrl:use"
	ActionDistribute = "odrl:distribute"
	ActionDelete     = "odrl:delete"

	ExpectedPayloadDigest = "3780df1071b0f2eec8a881ffd48425c3a1a60738d11cc2ba7debdf1cea992d63"
	ExpectedHMAC          = "e635c7c1991742a5c36992fc0da32a7abc80b32aa5777a1142adaab55183681c"
)

type Dataset struct {
	CaseName                       string
	Question                       string
	RequestPurpose                 string
	RequestAction                  string
	GatewayCreatedAt               string
	GatewayExpiresAt               string
	CityAuthAt                     string
	CityDutyAt                     string
	CurrentAlertLevel              int
	AlertLevelAtLeast              int
	CurrentIndoorTempC             float64
	IndoorTempCAtLeast             float64
	HoursAtOrAboveThreshold        int
	HoursAtOrAboveThresholdAtLeast int
	RemainingPrepaidCreditEur      float64
	EnergyCreditEurAtMost          float64
	MinimumActiveNeedCount         int
	MaxPackageCostEur              int
	DispatchesLogged               int
	VulnerabilityFlags             []string
	Packages                       []SupportPackage
	Insight                        Insight
	Policy                         Policy
	Signature                      Signature
	ReasonText                     string
}

type SupportPackage struct {
	ID           string
	PackageID    string
	Name         string
	CostEur      int
	Capabilities []string
}

type Insight struct {
	ID                  string
	Metric              string
	ThresholdCount      int
	ThresholdDisplay    string
	SupportPolicy       string
	ScopeDevice         string
	ScopeEvent          string
	Municipality        string
	CreatedAt           string
	ExpiresAt           string
	SerializedLowercase string
}

type Policy struct {
	Profile            string
	PermissionAction   string
	PermissionTarget   string
	PermissionPurpose  string
	ProhibitionAction  string
	ProhibitionTarget  string
	ProhibitionPurpose string
	DutyAction         string
	DutyAt             string
}

type Signature struct {
	Algorithm            string
	KeyID                string
	Created              string
	PayloadHashSHA256    string
	SignatureHMAC        string
	HMACVerificationMode string
}

type Need struct {
	Name         string
	Active       bool
	Why          string
	Capabilities []string
}

type Candidate struct {
	Package      SupportPackage
	Eligible     bool
	CoversActive int
	Reason       string
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

type Analysis struct {
	Needs                        []Need
	ActiveNeedCount              int
	RequiredCapabilities         []string
	Candidates                   []Candidate
	RecommendedPackage           SupportPackage
	RecommendationFound          bool
	PayloadHashMatches           bool
	SignatureVerifies            bool
	Minimized                    bool
	ScopeComplete                bool
	AuthorizationAllowed         bool
	DutyTimingConsistent         bool
	TenantScreeningProhibited    bool
	PriorityCoolingSupportNeeded bool
	Checks                       []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, fixture())
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
}

func fixture() Dataset {
	insightID := "https://example.org/insight/calidor"
	return Dataset{
		CaseName:                       "calidor",
		Question:                       "Is the Calidor heat-response system allowed to use a narrow household support insight for heatwave response, and if so which support package should it recommend?",
		RequestPurpose:                 PurposeHeatwaveResponse,
		RequestAction:                  ActionUse,
		GatewayCreatedAt:               "2026-07-18T09:00:00+00:00",
		GatewayExpiresAt:               "2026-07-18T21:00:00+00:00",
		CityAuthAt:                     "2026-07-18T09:05:00+00:00",
		CityDutyAt:                     "2026-07-18T20:30:00+00:00",
		CurrentAlertLevel:              4,
		AlertLevelAtLeast:              3,
		CurrentIndoorTempC:             31.4,
		IndoorTempCAtLeast:             30.0,
		HoursAtOrAboveThreshold:        9,
		HoursAtOrAboveThresholdAtLeast: 6,
		RemainingPrepaidCreditEur:      3.2,
		EnergyCreditEurAtMost:          5.0,
		MinimumActiveNeedCount:         3,
		MaxPackageCostEur:              20,
		DispatchesLogged:               1,
		VulnerabilityFlags:             []string{"heat_sensitive_condition", "mobility_limitation"},
		Packages: []SupportPackage{
			{ID: "pkg_CHECK", PackageID: "pkg:CHECK", Name: "Cooling Check Call", CostEur: 8, Capabilities: []string{"welfare_check"}},
			{ID: "pkg_VOUCHER", PackageID: "pkg:VOUCHER", Name: "Cooling Center Transport Voucher", CostEur: 12, Capabilities: []string{"transport", "welfare_check"}},
			{ID: "pkg_BUNDLE", PackageID: "pkg:BUNDLE", Name: "Calidor Priority Cooling Bundle", CostEur: 18, Capabilities: []string{"bill_credit", "cooling_kit", "transport", "welfare_check"}},
			{ID: "pkg_DELUXE", PackageID: "pkg:DELUXE", Name: "Extended Resilience Package", CostEur: 28, Capabilities: []string{"bill_credit", "cooling_kit", "transport", "welfare_check", "followup_visit"}},
		},
		Insight: Insight{
			ID:                  insightID,
			Metric:              "active_need_count",
			ThresholdCount:      3,
			ThresholdDisplay:    "3.0",
			SupportPolicy:       "lowest_cost_covering_package",
			ScopeDevice:         "household-gateway",
			ScopeEvent:          "heat-alert-window",
			Municipality:        "Calidor",
			CreatedAt:           "2026-07-18T09:00:00+00:00",
			ExpiresAt:           "2026-07-18T21:00:00+00:00",
			SerializedLowercase: "{\"createdat\":\"2026-07-18t09:00:00+00:00\",\"expiresat\":\"2026-07-18t21:00:00+00:00\",\"id\":\"https://example.org/insight/calidor\",\"metric\":\"active_need_count\",\"municipality\":\"calidor\",\"scopedevice\":\"household-gateway\",\"scopeevent\":\"heat-alert-window\",\"supportpolicy\":\"lowest_cost_covering_package\",\"threshold\":3,\"type\":\"ins:insight\"}",
		},
		Policy: Policy{
			Profile:            "Calidor-Heatwave-Policy",
			PermissionAction:   ActionUse,
			PermissionTarget:   insightID,
			PermissionPurpose:  PurposeHeatwaveResponse,
			ProhibitionAction:  ActionDistribute,
			ProhibitionTarget:  insightID,
			ProhibitionPurpose: PurposeTenantScreening,
			DutyAction:         ActionDelete,
			DutyAt:             "2026-07-18T21:00:00+00:00",
		},
		Signature: Signature{
			Algorithm:            "HMAC-SHA256",
			KeyID:                "calidor-demo-shared-secret",
			Created:              "2026-07-18T09:00:00+00:00",
			PayloadHashSHA256:    ExpectedPayloadDigest,
			SignatureHMAC:        ExpectedHMAC,
			HMACVerificationMode: "trustedPrecomputedInput",
		},
		ReasonText: "The gateway keeps raw indoor heat, vulnerability, and prepaid-energy data local, derives a priority-support signal, and shares only a scoped heatwave-response envelope with expiry.",
	}
}

func derive(ds Dataset) Analysis {
	needs := []Need{
		{
			Name:         "heat alert",
			Active:       ds.CurrentAlertLevel >= ds.AlertLevelAtLeast,
			Why:          fmt.Sprintf("alert level %d is at least %d", ds.CurrentAlertLevel, ds.AlertLevelAtLeast),
			Capabilities: []string{},
		},
		{
			Name:         "unsafe indoor heat",
			Active:       ds.CurrentIndoorTempC >= ds.IndoorTempCAtLeast && ds.HoursAtOrAboveThreshold >= ds.HoursAtOrAboveThresholdAtLeast,
			Why:          fmt.Sprintf("%.1f°C for %d hours reaches the %.1f°C/%d hour threshold", ds.CurrentIndoorTempC, ds.HoursAtOrAboveThreshold, ds.IndoorTempCAtLeast, ds.HoursAtOrAboveThresholdAtLeast),
			Capabilities: []string{"cooling_kit"},
		},
		{
			Name:         "vulnerability present",
			Active:       len(ds.VulnerabilityFlags) > 0,
			Why:          "the local gateway sees heat-sensitive and mobility flags",
			Capabilities: []string{"welfare_check", "transport"},
		},
		{
			Name:         "energy constraint",
			Active:       ds.RemainingPrepaidCreditEur <= ds.EnergyCreditEurAtMost,
			Why:          fmt.Sprintf("€%.1f prepaid credit is at or below the €%.1f limit", ds.RemainingPrepaidCreditEur, ds.EnergyCreditEurAtMost),
			Capabilities: []string{"bill_credit"},
		},
	}

	activeNeedCount := countActiveNeeds(needs)
	required := requiredCapabilities(needs)
	candidates, recommended, found := choosePackage(ds.Packages, ds.MaxPackageCostEur, required)

	analysis := Analysis{
		Needs:                        needs,
		ActiveNeedCount:              activeNeedCount,
		RequiredCapabilities:         required,
		Candidates:                   candidates,
		RecommendedPackage:           recommended,
		RecommendationFound:          found,
		PayloadHashMatches:           ds.Signature.PayloadHashSHA256 == ExpectedPayloadDigest,
		SignatureVerifies:            ds.Signature.SignatureHMAC == ExpectedHMAC && ds.Signature.HMACVerificationMode == "trustedPrecomputedInput",
		Minimized:                    minimized(ds.Insight.SerializedLowercase),
		ScopeComplete:                ds.Insight.ScopeDevice != "" && ds.Insight.ScopeEvent != "" && ds.Insight.ExpiresAt != "",
		AuthorizationAllowed:         authorizationAllowed(ds),
		DutyTimingConsistent:         notAfter(ds.CityDutyAt, ds.Insight.ExpiresAt),
		TenantScreeningProhibited:    ds.Policy.ProhibitionAction == ActionDistribute && ds.Policy.ProhibitionPurpose == PurposeTenantScreening,
		PriorityCoolingSupportNeeded: activeNeedCount >= ds.MinimumActiveNeedCount,
	}
	analysis.Checks = buildChecks(ds, analysis)
	return analysis
}

func requiredCapabilities(needs []Need) []string {
	seen := map[string]bool{}
	caps := []string{}
	for _, need := range needs {
		if !need.Active {
			continue
		}
		for _, cap := range need.Capabilities {
			if !seen[cap] {
				seen[cap] = true
				caps = append(caps, cap)
			}
		}
	}
	sort.Strings(caps)
	return caps
}

func choosePackage(packages []SupportPackage, budget int, required []string) ([]Candidate, SupportPackage, bool) {
	candidates := make([]Candidate, 0, len(packages))
	for _, pkg := range packages {
		covers := coversAll(pkg, required)
		withinBudget := pkg.CostEur <= budget
		eligible := covers && withinBudget
		candidates = append(candidates, Candidate{
			Package:      pkg,
			Eligible:     eligible,
			CoversActive: countCoveredCapabilities(pkg, required),
			Reason:       packageReason(pkg, budget, required, covers, withinBudget),
		})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Package.CostEur == candidates[j].Package.CostEur {
			return candidates[i].Package.PackageID < candidates[j].Package.PackageID
		}
		return candidates[i].Package.CostEur < candidates[j].Package.CostEur
	})

	for _, candidate := range candidates {
		if candidate.Eligible {
			return candidates, candidate.Package, true
		}
	}
	return candidates, SupportPackage{}, false
}

func packageReason(pkg SupportPackage, budget int, required []string, covers bool, withinBudget bool) string {
	parts := []string{}
	if covers {
		parts = append(parts, "covers all required capabilities")
	} else {
		parts = append(parts, fmt.Sprintf("covers %d/%d required capabilities", countCoveredCapabilities(pkg, required), len(required)))
	}
	if withinBudget {
		parts = append(parts, "within budget")
	} else {
		parts = append(parts, fmt.Sprintf("over budget by €%d", pkg.CostEur-budget))
	}
	return strings.Join(parts, "; ")
}

func countActiveNeeds(needs []Need) int {
	count := 0
	for _, need := range needs {
		if need.Active {
			count++
		}
	}
	return count
}

func countCoveredCapabilities(pkg SupportPackage, required []string) int {
	count := 0
	for _, cap := range required {
		if hasCapability(pkg, cap) {
			count++
		}
	}
	return count
}

func coversAll(pkg SupportPackage, required []string) bool {
	for _, cap := range required {
		if !hasCapability(pkg, cap) {
			return false
		}
	}
	return true
}

func hasCapability(pkg SupportPackage, cap string) bool {
	for _, own := range pkg.Capabilities {
		if own == cap {
			return true
		}
	}
	return false
}

func authorizationAllowed(ds Dataset) bool {
	return ds.Policy.PermissionAction == ds.RequestAction &&
		ds.Policy.PermissionPurpose == ds.RequestPurpose &&
		notAfter(ds.CityAuthAt, ds.Insight.ExpiresAt)
}

func notAfter(left string, right string) bool {
	leftTime, leftErr := time.Parse(time.RFC3339, left)
	rightTime, rightErr := time.Parse(time.RFC3339, right)
	if leftErr != nil || rightErr != nil {
		return false
	}
	return !leftTime.After(rightTime)
}

func minimized(serialized string) bool {
	forbidden := []string{"heat_sensitive_condition", "mobility_limitation", "credit", "meter_trace"}
	lower := strings.ToLower(serialized)
	for _, term := range forbidden {
		if strings.Contains(lower, term) {
			return false
		}
	}
	return true
}

func buildChecks(ds Dataset, analysis Analysis) []Check {
	recommendedEligible := analysis.RecommendationFound && analysis.RecommendedPackage.CostEur <= ds.MaxPackageCostEur && coversAll(analysis.RecommendedPackage, analysis.RequiredCapabilities)
	lowestEligible := analysis.RecommendationFound && isLowestEligible(analysis.RecommendedPackage, analysis.Candidates)
	allCore := analysis.SignatureVerifies && analysis.PayloadHashMatches && analysis.Minimized && analysis.ScopeComplete && analysis.AuthorizationAllowed && analysis.PriorityCoolingSupportNeeded && recommendedEligible && analysis.DutyTimingConsistent && analysis.TenantScreeningProhibited

	return []Check{
		{Label: "C1", OK: analysis.SignatureVerifies, Text: "the trusted precomputed HMAC signature verifies"},
		{Label: "C2", OK: analysis.PayloadHashMatches, Text: "payload hash matches the source envelope digest"},
		{Label: "C3", OK: analysis.Minimized, Text: "the shared insight strips raw heat, vulnerability, credit, and meter-trace terms"},
		{Label: "C4", OK: analysis.ScopeComplete, Text: "the insight has a device scope, event scope, and expiry"},
		{Label: "C5", OK: analysis.AuthorizationAllowed, Text: "heatwave-response use is authorized before the insight expires"},
		{Label: "C6", OK: analysis.Needs[0].Active, Text: "the heat alert is active"},
		{Label: "C7", OK: analysis.Needs[1].Active, Text: "indoor heat is unsafe for long enough to trigger support"},
		{Label: "C8", OK: analysis.PriorityCoolingSupportNeeded, Text: "the active-need count reaches the priority threshold"},
		{Label: "C9", OK: recommendedEligible, Text: "the recommended package is within budget and covers every required capability"},
		{Label: "C10", OK: lowestEligible, Text: "the lowest-cost eligible package is chosen"},
		{Label: "C11", OK: analysis.DutyTimingConsistent, Text: "the deletion duty is scheduled before expiry"},
		{Label: "C12", OK: analysis.TenantScreeningProhibited, Text: "reuse for tenant screening is prohibited"},
		{Label: "C13", OK: allCore, Text: "all source-level policy, scope, signature, and package checks pass"},
	}
}

func isLowestEligible(recommended SupportPackage, candidates []Candidate) bool {
	for _, candidate := range candidates {
		if candidate.Eligible && candidate.Package.CostEur < recommended.CostEur {
			return false
		}
	}
	return true
}

func printAnswer(ds Dataset, analysis Analysis) {
	fmt.Println("# Calidor")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("Name: Calidor")
	fmt.Printf("Municipality: %s\n", ds.Insight.Municipality)
	fmt.Printf("Metric: %s\n", ds.Insight.Metric)
	fmt.Printf("Active need count: %d/%d\n", analysis.ActiveNeedCount, ds.MinimumActiveNeedCount)
	if analysis.RecommendationFound {
		fmt.Printf("Recommended package: %s\n", analysis.RecommendedPackage.Name)
		fmt.Printf("Package cost: €%d\n", analysis.RecommendedPackage.CostEur)
	} else {
		fmt.Println("Recommended package: none")
	}
	fmt.Printf("Budget cap: €%d\n", ds.MaxPackageCostEur)
	fmt.Printf("Payload SHA-256: %s\n", ds.Signature.PayloadHashSHA256)
	fmt.Printf("Envelope HMAC-SHA-256: %s\n", ds.Signature.SignatureHMAC)
	fmt.Printf("decision : %s\n", allowedText(analysis.AuthorizationAllowed))
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason")
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Println("The gateway desensitizes local heat, vulnerability, and prepaid-energy stress into an expiring municipal support insight.")
	fmt.Printf("metric : %s\n", ds.Insight.Metric)
	fmt.Printf("threshold : %s\n", ds.Insight.ThresholdDisplay)
	fmt.Printf("scope : %s @ %s\n", ds.Insight.ScopeDevice, ds.Insight.ScopeEvent)
	fmt.Printf("required capabilities: %s\n", strings.Join(analysis.RequiredCapabilities, ", "))
	fmt.Println()

	for _, need := range analysis.Needs {
		fmt.Printf("%s : %s - %s\n", need.Name, activeText(need.Active), need.Why)
	}
	fmt.Printf("vulnerability flags kept local : %d\n", len(ds.VulnerabilityFlags))
	fmt.Printf("expires at : %s\n", ds.Insight.ExpiresAt)
	fmt.Println()

	fmt.Printf("support policy : %s\n", ds.Insight.SupportPolicy)
	fmt.Println("candidate packages:")
	for _, candidate := range analysis.Candidates {
		marker := "reject"
		if candidate.Eligible {
			marker = "eligible"
		}
		if analysis.RecommendationFound && candidate.Package.ID == analysis.RecommendedPackage.ID {
			marker = "selected"
		}
		fmt.Printf("  %-12s : %s, cost=€%d, %s\n", candidate.Package.PackageID, marker, candidate.Package.CostEur, candidate.Reason)
	}
	fmt.Println()

	if analysis.RecommendationFound {
		fmt.Printf("Selected package %q covers %s.\n", analysis.RecommendedPackage.Name, strings.Join(analysis.RecommendedPackage.Capabilities, ", "))
	}
	fmt.Printf("Usage is permitted only for purpose %q and the envelope expires at %s.\n", ds.Policy.PermissionPurpose, ds.Insight.ExpiresAt)
	fmt.Printf("Tenant-screening reuse is blocked by a prohibition on %s for purpose %q.\n", ds.Policy.ProhibitionAction, ds.Policy.ProhibitionPurpose)
	fmt.Printf("reason.txt : %s\n", ds.ReasonText)
	fmt.Printf("dispatches logged : %d\n", ds.DispatchesLogged)
	fmt.Println()
}

func activeText(active bool) string {
	if active {
		return "active"
	}
	return "inactive"
}

func allowedText(allowed bool) string {
	if allowed {
		return "ALLOWED"
	}
	return "DENIED"
}

func okText(ok bool) string {
	if ok {
		return "OK"
	}
	return "FAIL"
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func passedCount(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}
