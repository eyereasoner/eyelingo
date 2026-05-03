// flandor.go
//
// A self-contained Go translation of the Eyeling Flandor insight-economy
// example.
//
// The source N3 file models a regional policy decision. Exporters, training
// actors, and grid operators keep their sensitive details local. The region only
// receives a narrow macro-economic insight: enough aggregated pressure is active
// to justify a temporary retooling response for Flanders.
//
// This program is intentionally not a generic RDF, ODRL, crypto, or N3 reasoner.
// It translates the concrete facts and rules for this example into ordinary Go
// become a minimal, signed, expiring decision object for a public policy board.
//
// Run:
//
//	go run flandor.go
//
// The program has no third-party dependencies.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"sort"
	"strings"
	"time"
)

const seeExampleName = "flandor"

const sourceFile = "flandor.n3"

const (
	PurposeRegionalStabilization = "regional_stabilization"
	PurposeFirmSurveillance      = "firm_surveillance"

	ActionUse        = "odrl:use"
	ActionDistribute = "odrl:distribute"
	ActionDelete     = "odrl:delete"

	ExpectedPayloadDigest = "10a85e861075bef2a96c01c7f3238735f82b8f368deb62eafcedd1c4b7f7c707"
	DisplayPayloadDigest  = "718f5b17d07ab6a95503bc04a1000ddb132409f600659c03d21def81914b780b"
	ExpectedHMAC          = "955968ca99a191783bc00cba068128ccb9ff40a5e6114fda13a52c74ee27329e"
)

type Dataset struct {
	CaseName             string
	Question             string
	ExpectedFilesWritten int
	RequestPurpose       string
	RequestAction        string
	HubCreatedAt         string
	HubExpiresAt         string
	BoardAuthAt          string
	BoardDutyAt          string
	FilesWritten         int
	AuditEntries         int
	Region               Region
	Signals              SecureAggregate
	Clusters             []IndustrialCluster
	Labour               LabourMarket
	Grid                 GridSignal
	Budget               BudgetWindow
	Packages             []PolicyPackage
	Insight              Insight
	Policy               Policy
	Signature            Signature
}

type Region struct {
	ID   string
	Name string
}

type SecureAggregate struct {
	ObservedFirms       int
	AggregationLevel    string
	ContainsFirmNames   bool
	ContainsPayrollRows bool
}

type IndustrialCluster struct {
	ID                string
	Name              string
	ExportOrdersIndex int
	EnergyIntensity   int
}

type LabourMarket struct {
	TechVacancyRateTenths int
	TechVacancyRate       float64
}

type GridSignal struct {
	CongestionHours        int
	RenewableCurtailmentMW int
}

type BudgetWindow struct {
	MaxMEUR    int
	WindowName string
}

type PolicyPackage struct {
	ID                   string
	PackageID            string
	Name                 string
	CostMEUR             int
	WorkerCoverage       int
	GridReliefMW         int
	CoversExportWeakness bool
	CoversSkillsStrain   bool
	CoversGridStress     bool
}

type Insight struct {
	ID                  string
	Metric              string
	ThresholdScore      int
	ThresholdDisplay    string
	SuggestionPolicy    string
	ScopeDevice         string
	ScopeEvent          string
	Region              string
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
	DisplayPayloadSHA256 string
	SignatureHMAC        string
	HMACVerificationMode string
}

type Need struct {
	Name   string
	Active bool
	Why    string
}

type Candidate struct {
	Package      PolicyPackage
	Eligible     bool
	Reason       string
	CoversActive int
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

type Analysis struct {
	Needs                  []Need
	ActiveNeedCount        int
	ExportWeakClusters     []IndustrialCluster
	Candidates             []Candidate
	RecommendedPackage     PolicyPackage
	RecommendationFound    bool
	AuthorizationAllowed   bool
	DutyTimingConsistent   bool
	Minimized              bool
	ScopeComplete          bool
	PayloadHashMatches     bool
	SignatureVerifies      bool
	HMACMatches            bool
	SurveillanceProhibited bool
	FilesWrittenExpected   bool
	Checks                 []Check
}

func main() {
	ds := exampleinput.Load(seeExampleName, fixture())
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
}

func fixture() Dataset {
	insightID := "https://example.org/insight/flandor"
	return Dataset{
		CaseName:             "flandor",
		Question:             "Is the Flemish Economic Resilience Board allowed to use a neutral macro-economic insight for regional stabilization, and if so which package should it activate for Flanders?",
		ExpectedFilesWritten: 6,
		RequestPurpose:       PurposeRegionalStabilization,
		RequestAction:        ActionUse,
		HubCreatedAt:         "2026-04-08T07:00:00+00:00",
		HubExpiresAt:         "2026-04-08T19:00:00+00:00",
		BoardAuthAt:          "2026-04-08T09:15:00+00:00",
		BoardDutyAt:          "2026-04-08T18:30:00+00:00",
		FilesWritten:         6,
		AuditEntries:         1,
		Region:               Region{ID: "flanders", Name: "Flanders"},
		Signals: SecureAggregate{
			ObservedFirms:       217,
			AggregationLevel:    "regional_cluster",
			ContainsFirmNames:   false,
			ContainsPayrollRows: false,
		},
		Clusters: []IndustrialCluster{
			{ID: "cluster_antwerp", Name: "Antwerp chemicals", ExportOrdersIndex: 84, EnergyIntensity: 92},
			{ID: "cluster_ghent", Name: "Ghent manufacturing", ExportOrdersIndex: 87, EnergyIntensity: 76},
		},
		Labour: LabourMarket{TechVacancyRateTenths: 46, TechVacancyRate: 4.6},
		Grid:   GridSignal{CongestionHours: 19, RenewableCurtailmentMW: 240},
		Budget: BudgetWindow{MaxMEUR: 140, WindowName: "Q2 resilience window"},
		Packages: []PolicyPackage{
			{ID: "pkg_training_only", PackageID: "pkg:TRAIN_070", Name: "Flanders Skills Sprint", CostMEUR: 70, WorkerCoverage: 900, GridReliefMW: 0, CoversSkillsStrain: true},
			{ID: "pkg_logistics_only", PackageID: "pkg:PORT_095", Name: "Schelde Trade Buffer", CostMEUR: 95, WorkerCoverage: 300, GridReliefMW: 10, CoversExportWeakness: true},
			{ID: "pkg_flandor", PackageID: "pkg:RET_FLEX_120", Name: "Flandor Retooling Pulse", CostMEUR: 120, WorkerCoverage: 1200, GridReliefMW: 85, CoversExportWeakness: true, CoversSkillsStrain: true, CoversGridStress: true},
			{ID: "pkg_full_corridor", PackageID: "pkg:CORRIDOR_165", Name: "Full Corridor Shock Shield", CostMEUR: 165, WorkerCoverage: 1600, GridReliefMW: 110, CoversExportWeakness: true, CoversSkillsStrain: true, CoversGridStress: true},
		},
		Insight: Insight{
			ID:                  insightID,
			Metric:              "regional_retooling_priority",
			ThresholdScore:      3,
			ThresholdDisplay:    "3 active needs",
			SuggestionPolicy:    "lowest_cost_package_covering_all_active_needs",
			ScopeDevice:         "economic-resilience-board",
			ScopeEvent:          "budget-prep-window",
			Region:              "Flanders",
			CreatedAt:           "2026-04-08T07:00:00+00:00",
			ExpiresAt:           "2026-04-08T19:00:00+00:00",
			SerializedLowercase: "{\"createdat\":\"2026-04-08t07:00:00+00:00\",\"expiresat\":\"2026-04-08t19:00:00+00:00\",\"id\":\"https://example.org/insight/flandor\",\"metric\":\"regional_retooling_priority\",\"region\":\"flanders\",\"scopedevice\":\"economic-resilience-board\",\"scopeevent\":\"budget-prep-window\",\"suggestionpolicy\":\"lowest_cost_package_covering_all_active_needs\",\"threshold\":3,\"type\":\"ins:insight\"}",
		},
		Policy: Policy{
			Profile:            "Flandor-Insight-Policy",
			PermissionAction:   ActionUse,
			PermissionTarget:   insightID,
			PermissionPurpose:  PurposeRegionalStabilization,
			ProhibitionAction:  ActionDistribute,
			ProhibitionTarget:  insightID,
			ProhibitionPurpose: PurposeFirmSurveillance,
			DutyAction:         ActionDelete,
			DutyAt:             "2026-04-08T19:00:00+00:00",
		},
		Signature: Signature{
			Algorithm:            "HMAC-SHA256",
			KeyID:                "demo-shared-secret",
			Created:              "2026-04-08T07:00:00+00:00",
			PayloadHashSHA256:    ExpectedPayloadDigest,
			DisplayPayloadSHA256: DisplayPayloadDigest,
			SignatureHMAC:        ExpectedHMAC,
			HMACVerificationMode: "trustedPrecomputedInput",
		},
	}
}

func derive(ds Dataset) Analysis {
	exportClusters := exportWeakClusters(ds.Clusters)
	needs := []Need{
		{
			Name:   "export weakness",
			Active: len(exportClusters) > 0,
			Why:    "at least one cluster has exportOrdersIndex < 90",
		},
		{
			Name:   "skills strain",
			Active: ds.Labour.TechVacancyRateTenths > 39,
			Why:    "technical vacancy rate is above 3.9%",
		},
		{
			Name:   "grid stress",
			Active: ds.Grid.CongestionHours > 11,
			Why:    "grid congestion hours are above 11",
		},
	}

	activeNeedCount := countActiveNeeds(needs)
	candidates, recommended, found := choosePackage(ds.Packages, ds.Budget, needs)

	analysis := Analysis{
		Needs:                  needs,
		ActiveNeedCount:        activeNeedCount,
		ExportWeakClusters:     exportClusters,
		Candidates:             candidates,
		RecommendedPackage:     recommended,
		RecommendationFound:    found,
		AuthorizationAllowed:   authorizationAllowed(ds),
		DutyTimingConsistent:   notAfter(ds.BoardDutyAt, ds.Insight.ExpiresAt),
		Minimized:              minimized(ds.Insight.SerializedLowercase),
		ScopeComplete:          ds.Insight.ScopeDevice != "" && ds.Insight.ScopeEvent != "" && ds.Insight.ExpiresAt != "",
		PayloadHashMatches:     ds.Signature.PayloadHashSHA256 == ExpectedPayloadDigest,
		SignatureVerifies:      ds.Signature.HMACVerificationMode == "trustedPrecomputedInput",
		HMACMatches:            ds.Signature.SignatureHMAC == ExpectedHMAC && ds.Signature.HMACVerificationMode == "trustedPrecomputedInput",
		SurveillanceProhibited: ds.Policy.ProhibitionAction == ActionDistribute && ds.Policy.ProhibitionPurpose == PurposeFirmSurveillance,
		FilesWrittenExpected:   ds.FilesWritten == ds.ExpectedFilesWritten,
	}
	analysis.Checks = buildChecks(ds, analysis)
	return analysis
}

func exportWeakClusters(clusters []IndustrialCluster) []IndustrialCluster {
	weak := []IndustrialCluster{}
	for _, cluster := range clusters {
		if cluster.ExportOrdersIndex < 90 {
			weak = append(weak, cluster)
		}
	}
	return weak
}

func choosePackage(packages []PolicyPackage, budget BudgetWindow, needs []Need) ([]Candidate, PolicyPackage, bool) {
	candidates := make([]Candidate, 0, len(packages))
	for _, pkg := range packages {
		covers := coversActiveNeeds(pkg, needs)
		withinBudget := pkg.CostMEUR <= budget.MaxMEUR
		eligible := covers && withinBudget
		reason := packageReason(pkg, budget, needs, covers, withinBudget)
		candidates = append(candidates, Candidate{
			Package:      pkg,
			Eligible:     eligible,
			Reason:       reason,
			CoversActive: countCoveredActiveNeeds(pkg, needs),
		})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Package.CostMEUR == candidates[j].Package.CostMEUR {
			return candidates[i].Package.PackageID < candidates[j].Package.PackageID
		}
		return candidates[i].Package.CostMEUR < candidates[j].Package.CostMEUR
	})

	for _, candidate := range candidates {
		if candidate.Eligible {
			return candidates, candidate.Package, true
		}
	}
	return candidates, PolicyPackage{}, false
}

func coversActiveNeeds(pkg PolicyPackage, needs []Need) bool {
	for _, need := range needs {
		if !need.Active {
			continue
		}
		if need.Name == "export weakness" && !pkg.CoversExportWeakness {
			return false
		}
		if need.Name == "skills strain" && !pkg.CoversSkillsStrain {
			return false
		}
		if need.Name == "grid stress" && !pkg.CoversGridStress {
			return false
		}
	}
	return true
}

func countCoveredActiveNeeds(pkg PolicyPackage, needs []Need) int {
	count := 0
	for _, need := range needs {
		if !need.Active {
			continue
		}
		if need.Name == "export weakness" && pkg.CoversExportWeakness {
			count++
		}
		if need.Name == "skills strain" && pkg.CoversSkillsStrain {
			count++
		}
		if need.Name == "grid stress" && pkg.CoversGridStress {
			count++
		}
	}
	return count
}

func packageReason(pkg PolicyPackage, budget BudgetWindow, needs []Need, covers bool, withinBudget bool) string {
	parts := []string{}
	if covers {
		parts = append(parts, "covers all active needs")
	} else {
		parts = append(parts, fmt.Sprintf("covers %d/%d active needs", countCoveredActiveNeeds(pkg, needs), countActiveNeeds(needs)))
	}
	if withinBudget {
		parts = append(parts, "within budget")
	} else {
		parts = append(parts, fmt.Sprintf("over budget by €%dM", pkg.CostMEUR-budget.MaxMEUR))
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

func authorizationAllowed(ds Dataset) bool {
	return ds.Policy.PermissionAction == ds.RequestAction &&
		ds.Policy.PermissionPurpose == ds.RequestPurpose &&
		notAfter(ds.BoardAuthAt, ds.Insight.ExpiresAt)
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
	forbidden := []string{"salary", "payroll", "invoice", "medical", "firmname"}
	lower := strings.ToLower(serialized)
	for _, term := range forbidden {
		if strings.Contains(lower, term) {
			return false
		}
	}
	return true
}

func buildChecks(ds Dataset, analysis Analysis) []Check {
	thresholdReached := analysis.ActiveNeedCount >= ds.Insight.ThresholdScore
	withinBudget := analysis.RecommendationFound && analysis.RecommendedPackage.CostMEUR <= ds.Budget.MaxMEUR
	coversAll := analysis.RecommendationFound && coversActiveNeeds(analysis.RecommendedPackage, analysis.Needs)
	lowest := analysis.RecommendationFound && isLowestEligible(analysis.RecommendedPackage, analysis.Candidates)
	allCore := analysis.SignatureVerifies && analysis.PayloadHashMatches && analysis.Minimized && analysis.ScopeComplete && analysis.AuthorizationAllowed && withinBudget && coversAll && analysis.DutyTimingConsistent && analysis.SurveillanceProhibited && analysis.FilesWrittenExpected

	return []Check{
		{Label: "C1", OK: analysis.PayloadHashMatches, Text: "payload hash matches the source envelope digest"},
		{Label: "C2", OK: analysis.HMACMatches, Text: "HMAC value matches the trusted precomputed signature"},
		{Label: "C3", OK: thresholdReached, Text: "export weakness, skills strain, and grid stress reach the three-need threshold"},
		{Label: "C4", OK: analysis.ScopeComplete, Text: "the insight has a device scope, event scope, and expiry"},
		{Label: "C5", OK: analysis.Minimized, Text: "the shared insight omits firm names, payroll rows, and other sensitive terms"},
		{Label: "C6", OK: analysis.AuthorizationAllowed, Text: "regional stabilization use is authorized before the insight expires"},
		{Label: "C7", OK: analysis.DutyTimingConsistent, Text: "the deletion duty is scheduled before expiry"},
		{Label: "C8", OK: analysis.SurveillanceProhibited, Text: "reuse for firm surveillance is prohibited"},
		{Label: "C9", OK: withinBudget, Text: "the recommended package fits inside the €140M budget"},
		{Label: "C10", OK: coversAll, Text: "the recommended package covers all active needs"},
		{Label: "C11", OK: lowest, Text: "the lowest-cost eligible package is chosen"},
		{Label: "C12", OK: analysis.FilesWrittenExpected, Text: "the expected six files are recorded as written"},
		{Label: "C13", OK: allCore, Text: "all source-level policy, scope, signature, and package checks pass"},
	}
}

func isLowestEligible(recommended PolicyPackage, candidates []Candidate) bool {
	for _, candidate := range candidates {
		if candidate.Eligible && candidate.Package.CostMEUR < recommended.CostMEUR {
			return false
		}
	}
	return true
}

func printAnswer(ds Dataset, analysis Analysis) {
	fmt.Println("# Flandor")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("Name: Flandor")
	fmt.Printf("Region: %s\n", ds.Region.Name)
	fmt.Printf("Metric: %s\n", ds.Insight.Metric)
	fmt.Printf("Active need count: %d/%d\n", analysis.ActiveNeedCount, ds.Insight.ThresholdScore)
	if analysis.RecommendationFound {
		fmt.Printf("Recommended package: %s\n", analysis.RecommendedPackage.Name)
	} else {
		fmt.Println("Recommended package: none")
	}
	fmt.Printf("Budget cap: €%dM\n", ds.Budget.MaxMEUR)
	if analysis.RecommendationFound {
		fmt.Printf("Package cost: €%dM\n", analysis.RecommendedPackage.CostMEUR)
		fmt.Printf("Worker coverage: %d\n", analysis.RecommendedPackage.WorkerCoverage)
		fmt.Printf("Grid relief: %d MW\n", analysis.RecommendedPackage.GridReliefMW)
	}
	fmt.Printf("Payload SHA-256: %s\n", ds.Signature.DisplayPayloadSHA256)
	fmt.Printf("Envelope HMAC-SHA-256: %s\n", ds.Signature.SignatureHMAC)
	fmt.Printf("decision : %s\n", allowedText(analysis.AuthorizationAllowed))
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason")
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("aggregate : observedFirms=%d level=%s containsFirmNames=%s containsPayrollRows=%s\n", ds.Signals.ObservedFirms, ds.Signals.AggregationLevel, yesNo(ds.Signals.ContainsFirmNames), yesNo(ds.Signals.ContainsPayrollRows))
	fmt.Println()

	for _, need := range analysis.Needs {
		fmt.Printf("%s : %s - %s\n", need.Name, activeText(need.Active), need.Why)
	}
	fmt.Printf("export weak clusters : %s\n", clusterSummary(analysis.ExportWeakClusters))
	fmt.Printf("technical vacancy rate : %.1f%% (threshold > 3.9%%)\n", ds.Labour.TechVacancyRate)
	fmt.Printf("grid congestion : %d hours (threshold > 11)\n", ds.Grid.CongestionHours)
	fmt.Println()

	fmt.Printf("recommendation policy : %s\n", ds.Insight.SuggestionPolicy)
	fmt.Println("candidate packages:")
	for _, candidate := range analysis.Candidates {
		marker := "reject"
		if candidate.Eligible {
			marker = "eligible"
		}
		if analysis.RecommendationFound && candidate.Package.ID == analysis.RecommendedPackage.ID {
			marker = "selected"
		}
		fmt.Printf("  %-15s : %s, cost=€%dM, %s\n", candidate.Package.PackageID, marker, candidate.Package.CostMEUR, candidate.Reason)
	}
	fmt.Println()

	if analysis.RecommendationFound {
		fmt.Printf("Selected package \"%s\" covers export=%s, skills=%s, grid=%s, cost=€%dM.\n", analysis.RecommendedPackage.Name, yesNo(analysis.RecommendedPackage.CoversExportWeakness), yesNo(analysis.RecommendedPackage.CoversSkillsStrain), yesNo(analysis.RecommendedPackage.CoversGridStress), analysis.RecommendedPackage.CostMEUR)
	}
	fmt.Printf("Usage is permitted only for purpose \"%s\" and the envelope expires at %s.\n", ds.Policy.PermissionPurpose, ds.Insight.ExpiresAt)
	fmt.Printf("Surveillance reuse is blocked by a prohibition on %s for purpose \"%s\".\n", ds.Policy.ProhibitionAction, ds.Policy.ProhibitionPurpose)
	fmt.Printf("Deletion duty time : %s\n", ds.BoardDutyAt)
	fmt.Println()
}

func clusterSummary(clusters []IndustrialCluster) string {
	if len(clusters) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(clusters))
	for _, cluster := range clusters {
		parts = append(parts, fmt.Sprintf("%s=%d", cluster.Name, cluster.ExportOrdersIndex))
	}
	return strings.Join(parts, ", ")
}

func allowedText(ok bool) string {
	if ok {
		return "ALLOWED"
	}
	return "DENIED"
}

func activeText(ok bool) string {
	if ok {
		return "active"
	}
	return "inactive"
}

func passedChecks(checks []Check) int {
	n := 0
	for _, check := range checks {
		if check.OK {
			n++
		}
	}
	return n
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}
