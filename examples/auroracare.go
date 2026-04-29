// auroracare.go
//
// A self-contained Go translation of the Eyeling AuroraCare purpose-based
// medical-data exchange example.
//
// The source N3 file models a policy decision point (PDP). The same health data
// can be allowed or denied depending on purpose, requester role, care-team
// relationship, environment, privacy safeguards, and patient consent. This Go
// version keeps those concrete policy checks visible as ordinary structs and
// functions.
//
// This is intentionally not a generic RDF, ODRL, DPV, EHDS, or N3 reasoner. It
// translates the source facts and rules for this one example into explicit Go
// data and deterministic checks, then emits a compact report.
//
// Run:
//
//	go run auroracare.go
//
// The program has no third-party dependencies.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

const eyelingoExampleName = "auroracare"

const sourceFile = "auroracare.n3"

const (
	PurposePrimaryCare      = "PrimaryCareManagement"
	PurposeRemoteMonitoring = "PatientRemoteMonitoring"
	PurposeQualitySafety    = "EnsureQualitySafetyHealthcare"
	PurposeResearch         = "HealthcareScientificResearch"
	PurposeInsurance        = "InsuranceManagement"
	PurposeAITraining       = "TrainTestAndEvaluateAISystemsAlgorithms"

	CategoryPatientSummary = "PATIENT_SUMMARY"
	CategoryLabResults     = "LAB_RESULTS"
	CategoryImagingReport  = "IMAGING_REPORT"

	TOMAnonymisation = "Anonymisation"
)

type Dataset struct {
	CaseName   string
	Question   string
	Policies   []Policy
	Subjects   map[string]Subject
	Requesters map[string]Requester
	Scenarios  []Scenario
}

type Subject struct {
	ID           string
	Label        string
	ConsentAllow []string
	ConsentDeny  []string
}

type Requester struct {
	ID          string
	Label       string
	LinkedTo    string
	DefaultRole string
}

type Policy struct {
	ID                   string
	UID                  string
	Name                 string
	Kind                 string
	AllowedPurposes      []string
	ProhibitedPurposes   []string
	RequiredRole         string
	RequiredEnvironment  string
	RequiredTOM          string
	AllowAnyCategories   []string
	RequireAllCategories []string
	Duties               []string
}

type Scenario struct {
	Key         string
	Label       string
	Description string
	RequesterID string
	Role        string
	SubjectID   string
	Purpose     string
	Environment string
	TOM         string
	Categories  []string
	Expected    string
}

type ScenarioResult struct {
	Scenario          Scenario
	Decision          string
	Reason            string
	MatchedPolicyUID  string
	MatchedPolicyName string
	MatchedDuties     []string
	Trace             []string
	CareTeamLinked    bool
	SubjectOptIn      bool
	SubjectOptOut     bool
	CategoryMode      string
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

type Analysis struct {
	Results     []ScenarioResult
	Checks      []Check
	PermitCount int
	DenyCount   int
	PolicyUIDs  []string
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, fixture())
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(analysis)
	printChecks(analysis.Checks)
	printAudit(ds, analysis)
	if !allChecksOK(analysis.Checks) {
		os.Exit(1)
	}
}

func fixture() Dataset {
	return Dataset{
		CaseName: "auroracare",
		Question: "For each AuroraCare scenario, should the PDP permit or deny the requested use of health data, and why?",
		Policies: []Policy{
			{
				ID:                 "policy_primary",
				UID:                "urn:policy:primary-care-001",
				Name:               "primary care",
				Kind:               "permission",
				AllowedPurposes:    []string{PurposePrimaryCare, PurposeRemoteMonitoring},
				RequiredRole:       "clinician",
				AllowAnyCategories: []string{CategoryPatientSummary, CategoryLabResults},
			},
			{
				ID:                   "policy_qi",
				UID:                  "urn:policy:qi-2025-aurora",
				Name:                 "quality improvement",
				Kind:                 "permission",
				AllowedPurposes:      []string{PurposeQualitySafety},
				RequiredEnvironment:  "secure_env",
				RequireAllCategories: []string{CategoryLabResults, CategoryPatientSummary},
				Duties:               []string{"requireConsent", "noExfiltration"},
			},
			{
				ID:                  "policy_research",
				UID:                 "urn:policy:research-aurora-diabetes",
				Name:                "diabetes research",
				Kind:                "permission",
				AllowedPurposes:     []string{PurposeResearch},
				RequiredEnvironment: "secure_env",
				RequiredTOM:         TOMAnonymisation,
				AllowAnyCategories:  []string{CategoryLabResults, CategoryPatientSummary, CategoryImagingReport},
				Duties:              []string{"annualOutcomeReport", "noReidentification", "noExfiltration"},
			},
			{
				ID:                 "policy_deny_insurance",
				UID:                "urn:policy:deny-insurance",
				Name:               "insurance prohibition",
				Kind:               "prohibition",
				ProhibitedPurposes: []string{PurposeInsurance},
			},
		},
		Subjects: map[string]Subject{
			"ruben": {
				ID:           "ruben",
				Label:        "Ruben",
				ConsentAllow: []string{PurposeResearch},
				ConsentDeny:  []string{PurposeAITraining},
			},
		},
		Requesters: map[string]Requester{
			"clinician_alba": {
				ID:          "clinician_alba",
				Label:       "Clinician Alba",
				LinkedTo:    "ruben",
				DefaultRole: "clinician",
			},
			"gp_ruben": {
				ID:          "gp_ruben",
				Label:       "Ruben's GP",
				LinkedTo:    "ruben",
				DefaultRole: "clinician",
			},
			"qi_analyst": {
				ID:          "qi_analyst",
				Label:       "QI analyst",
				DefaultRole: "data_user",
			},
			"researcher_aurora": {
				ID:          "researcher_aurora",
				Label:       "AuroraCare researcher",
				DefaultRole: "data_user",
			},
			"insurer_bot": {
				ID:          "insurer_bot",
				Label:       "Insurance bot",
				DefaultRole: "data_user",
			},
			"ml_ops": {
				ID:          "ml_ops",
				Label:       "ML operations",
				DefaultRole: "data_user",
			},
		},
		Scenarios: []Scenario{
			{
				Key:         "A",
				Label:       "Primary care visit",
				Description: "Clinician in the patient's care team accessing the patient summary for primary care management.",
				RequesterID: "clinician_alba",
				Role:        "clinician",
				SubjectID:   "ruben",
				Purpose:     PurposePrimaryCare,
				Environment: "api_gateway",
				Categories:  []string{CategoryPatientSummary},
				Expected:    "PERMIT",
			},
			{
				Key:         "B",
				Label:       "Quality improvement (in scope)",
				Description: "QI analyst using lab results + summary in a secure environment.",
				RequesterID: "qi_analyst",
				Role:        "data_user",
				SubjectID:   "ruben",
				Purpose:     PurposeQualitySafety,
				Environment: "secure_env",
				Categories:  []string{CategoryLabResults, CategoryPatientSummary},
				Expected:    "PERMIT",
			},
			{
				Key:         "C",
				Label:       "Quality improvement (out of scope)",
				Description: "QI analyst with only lab results; policy expects labs + summary.",
				RequesterID: "qi_analyst",
				Role:        "data_user",
				SubjectID:   "ruben",
				Purpose:     PurposeQualitySafety,
				Environment: "secure_env",
				Categories:  []string{CategoryLabResults},
				Expected:    "DENY",
			},
			{
				Key:         "D",
				Label:       "Insurance management",
				Description: "Insurance bot attempting to use health data for insurance management.",
				RequesterID: "insurer_bot",
				Role:        "data_user",
				SubjectID:   "ruben",
				Purpose:     PurposeInsurance,
				Environment: "secure_env",
				Categories:  []string{CategoryPatientSummary},
				Expected:    "DENY",
			},
			{
				Key:         "E",
				Label:       "GP checks labs",
				Description: "GP for the same patient checking lab results via the API gateway.",
				RequesterID: "gp_ruben",
				Role:        "clinician",
				SubjectID:   "ruben",
				Purpose:     PurposePrimaryCare,
				Environment: "api_gateway",
				Categories:  []string{CategoryLabResults},
				Expected:    "PERMIT",
			},
			{
				Key:         "F",
				Label:       "Research on anonymised dataset",
				Description: "Researcher using anonymised labs + summary in a secure environment, with opt-in.",
				RequesterID: "researcher_aurora",
				Role:        "data_user",
				SubjectID:   "ruben",
				Purpose:     PurposeResearch,
				Environment: "secure_env",
				TOM:         TOMAnonymisation,
				Categories:  []string{CategoryPatientSummary, CategoryLabResults},
				Expected:    "PERMIT",
			},
			{
				Key:         "G",
				Label:       "AI training (opt-out)",
				Description: "Data user wants to train AI, but the subject opted out of AI training.",
				RequesterID: "ml_ops",
				Role:        "data_user",
				SubjectID:   "ruben",
				Purpose:     PurposeAITraining,
				Environment: "secure_env",
				Categories:  []string{CategoryPatientSummary, CategoryLabResults},
				Expected:    "DENY",
			},
		},
	}
}

func derive(ds Dataset) Analysis {
	results := make([]ScenarioResult, 0, len(ds.Scenarios))
	permitCount := 0
	denyCount := 0

	for _, scenario := range ds.Scenarios {
		result := evaluateScenario(ds, scenario)
		results = append(results, result)
		switch result.Decision {
		case "PERMIT":
			permitCount++
		case "DENY":
			denyCount++
		}
	}

	policyUIDs := make([]string, 0, len(ds.Policies))
	for _, policy := range ds.Policies {
		policyUIDs = append(policyUIDs, policy.UID)
	}
	sort.Strings(policyUIDs)

	analysis := Analysis{
		Results:     results,
		PermitCount: permitCount,
		DenyCount:   denyCount,
		PolicyUIDs:  policyUIDs,
	}
	analysis.Checks = buildChecks(ds, analysis)
	return analysis
}

func evaluateScenario(ds Dataset, scenario Scenario) ScenarioResult {
	subject := ds.Subjects[scenario.SubjectID]
	requester := ds.Requesters[scenario.RequesterID]
	result := ScenarioResult{
		Scenario:       scenario,
		CareTeamLinked: requester.LinkedTo == scenario.SubjectID,
		SubjectOptIn:   contains(subject.ConsentAllow, scenario.Purpose),
		SubjectOptOut:  contains(subject.ConsentDeny, scenario.Purpose),
	}

	if result.CareTeamLinked {
		result.Trace = append(result.Trace, "care_team_linked")
	}
	if result.SubjectOptIn {
		result.Trace = append(result.Trace, "subject_opted_in")
	}
	if result.SubjectOptOut {
		result.Trace = append(result.Trace, "subject_opted_out")
	}

	if policy, ok := matchingProhibition(ds.Policies, scenario.Purpose); ok {
		result.Decision = "DENY"
		result.Reason = "Denied: the requested purpose is prohibited by policy."
		result.MatchedPolicyUID = policy.UID
		result.MatchedPolicyName = policy.Name
		result.Trace = append(result.Trace, policy.UID+":deny:prohibition_matched")
		return result
	}

	if scenario.Purpose == PurposeAITraining && result.SubjectOptOut {
		result.Decision = "DENY"
		result.Reason = "Denied: the subject opted out of data use for AI training."
		result.Trace = append(result.Trace, "deny:subject_opted_out_ai_training")
		return result
	}

	for _, policy := range ds.Policies {
		if policy.Kind != "permission" {
			continue
		}
		matched, mode := permissionMatches(policy, scenario, result)
		if matched {
			result.Decision = "PERMIT"
			result.MatchedPolicyUID = policy.UID
			result.MatchedPolicyName = policy.Name
			result.MatchedDuties = append([]string(nil), policy.Duties...)
			result.CategoryMode = mode
			result.Reason = permitReason(policy)
			result.Trace = append(result.Trace, policy.UID+":permit:permission_matched")
			return result
		}
	}

	result.Decision = "DENY"
	result.Reason = "Denied: no permission matched the purpose, environment, safeguards, role, or requested categories."
	result.Trace = append(result.Trace, "deny:no_permission_matched")
	return result
}

func matchingProhibition(policies []Policy, purpose string) (Policy, bool) {
	for _, policy := range policies {
		if policy.Kind == "prohibition" && contains(policy.ProhibitedPurposes, purpose) {
			return policy, true
		}
	}
	return Policy{}, false
}

func permissionMatches(policy Policy, scenario Scenario, result ScenarioResult) (bool, string) {
	if !contains(policy.AllowedPurposes, scenario.Purpose) {
		return false, ""
	}
	if policy.RequiredRole != "" && policy.RequiredRole != scenario.Role {
		return false, ""
	}
	if policy.RequiredRole == "clinician" && !result.CareTeamLinked {
		return false, ""
	}
	if policy.RequiredEnvironment != "" && policy.RequiredEnvironment != scenario.Environment {
		return false, ""
	}
	if policy.RequiredTOM != "" && policy.RequiredTOM != scenario.TOM {
		return false, ""
	}
	if policy.RequiredTOM != "" && !result.SubjectOptIn {
		return false, ""
	}
	if len(policy.RequireAllCategories) > 0 && !containsAll(scenario.Categories, policy.RequireAllCategories) {
		return false, ""
	}
	if len(policy.AllowAnyCategories) > 0 && !intersects(scenario.Categories, policy.AllowAnyCategories) {
		return false, ""
	}
	if len(policy.RequireAllCategories) > 0 {
		return true, "all-of"
	}
	if len(policy.AllowAnyCategories) > 0 {
		return true, "any-of"
	}
	return true, "none"
}

func permitReason(policy Policy) string {
	switch policy.ID {
	case "policy_primary":
		return "Permitted: a clinician in the patient's care team matched the primary-care policy."
	case "policy_qi":
		return "Permitted: the quality-improvement policy matched the secure environment and required data categories."
	case "policy_research":
		return "Permitted: the subject opted in, the dataset is anonymised, and the research policy matched."
	default:
		return "Permitted: a policy permission matched."
	}
}

func buildChecks(ds Dataset, analysis Analysis) []Check {
	byKey := resultsByKey(analysis.Results)
	expectedOK := true
	for _, scenario := range ds.Scenarios {
		if byKey[scenario.Key].Decision != scenario.Expected {
			expectedOK = false
			break
		}
	}

	primaryOK := byKey["A"].Decision == "PERMIT" && byKey["A"].CareTeamLinked && byKey["E"].Decision == "PERMIT" && byKey["E"].CareTeamLinked
	qiOK := byKey["B"].Decision == "PERMIT" && byKey["C"].Decision == "DENY"
	insuranceOK := byKey["D"].Decision == "DENY" && strings.Contains(byKey["D"].MatchedPolicyUID, "deny-insurance")
	researchOK := byKey["F"].Decision == "PERMIT" && byKey["F"].SubjectOptIn && byKey["F"].Scenario.TOM == TOMAnonymisation
	aiOK := byKey["G"].Decision == "DENY" && byKey["G"].SubjectOptOut
	totalsOK := analysis.PermitCount == 4 && analysis.DenyCount == 3

	return []Check{
		{
			Label: "expected decisions",
			OK:    expectedOK,
			Text:  "all seven scenarios match the PERMIT/DENY outcomes encoded in the N3 example",
		},
		{
			Label: "primary care",
			OK:    primaryOK,
			Text:  "primary-care access requires a clinician role and a care-team link",
		},
		{
			Label: "quality improvement",
			OK:    qiOK,
			Text:  "quality improvement is allowed only when both lab results and patient summary are requested in the secure environment",
		},
		{
			Label: "insurance prohibition",
			OK:    insuranceOK,
			Text:  "insurance management is denied by a matching prohibition",
		},
		{
			Label: "research safeguards",
			OK:    researchOK,
			Text:  "research is allowed only with patient opt-in and anonymisation in the secure environment",
		},
		{
			Label: "AI opt-out",
			OK:    aiOK,
			Text:  "AI training is denied because the subject opted out",
		},
		{
			Label: "decision totals",
			OK:    totalsOK,
			Text:  "four scenarios are permitted and three are denied",
		},
	}
}

func printAnswer(ds Dataset, analysis Analysis) {
	fmt.Println("# AuroraCare")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println(ds.Question)
	fmt.Printf("- permit count : %d\n", analysis.PermitCount)
	fmt.Printf("- deny count : %d\n", analysis.DenyCount)
	fmt.Println()
	fmt.Println("- Scenario decisions:")
	for _, result := range analysis.Results {
		policy := "no policy matched"
		if result.MatchedPolicyUID != "" {
			policy = result.MatchedPolicyUID
		}
		fmt.Printf("  %s – %s : %s (%s)\n", result.Scenario.Key, result.Scenario.Label, result.Decision, policy)
	}
	fmt.Println()
}

func printReason(analysis Analysis) {
	fmt.Println("## Reason why")
	for _, result := range analysis.Results {
		fmt.Printf("%s – %s\n", result.Scenario.Key, result.Scenario.Label)
		fmt.Printf("  request : role=%s purpose=%s environment=%s categories=%s\n", result.Scenario.Role, result.Scenario.Purpose, result.Scenario.Environment, joinStrings(result.Scenario.Categories, ","))
		if result.Scenario.TOM != "" {
			fmt.Printf("  safeguard : tom=%s\n", result.Scenario.TOM)
		}
		fmt.Printf("  decision : %s\n", result.Decision)
		fmt.Printf("  reason : %s\n", result.Reason)
		fmt.Printf("  care-team linked : %s\n", yesNo(result.CareTeamLinked))
		fmt.Printf("  subject opt-in : %s\n", yesNo(result.SubjectOptIn))
		fmt.Printf("  subject opt-out : %s\n", yesNo(result.SubjectOptOut))
		if len(result.MatchedDuties) > 0 {
			fmt.Printf("  duties : %s\n", joinStrings(result.MatchedDuties, ","))
		}
		fmt.Printf("  trace : %s\n", joinStrings(result.Trace, ","))
		fmt.Println()
	}
}

func printChecks(checks []Check) {
	fmt.Println("## Check")
	for i, check := range checks {
		status := "FAIL"
		if check.OK {
			status = "OK"
		}
		fmt.Printf("C%d %s - %s\n", i+1, status, check.Text)
	}
	fmt.Println()
}

func printAudit(ds Dataset, analysis Analysis) {
	passed := countPassed(analysis.Checks)

	fmt.Println("## Go audit details")
	fmt.Printf("- platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("- source file : %s\n", sourceFile)
	fmt.Printf("- case : %s\n", ds.CaseName)
	fmt.Printf("- scenario count : %d\n", len(analysis.Results))
	fmt.Printf("- policy count : %d\n", len(ds.Policies))
	fmt.Printf("- policy uids : %s\n", joinStrings(analysis.PolicyUIDs, ","))
	fmt.Printf("- permit count : %d\n", analysis.PermitCount)
	fmt.Printf("- deny count : %d\n", analysis.DenyCount)
	fmt.Printf("- checks passed : %d/%d\n", passed, len(analysis.Checks))
	fmt.Printf("- all checks pass : %s\n", yesNo(passed == len(analysis.Checks)))
}

func resultsByKey(results []ScenarioResult) map[string]ScenarioResult {
	out := map[string]ScenarioResult{}
	for _, result := range results {
		out[result.Scenario.Key] = result
	}
	return out
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func containsAll(values []string, required []string) bool {
	for _, item := range required {
		if !contains(values, item) {
			return false
		}
	}
	return true
}

func intersects(a []string, b []string) bool {
	for _, x := range a {
		if contains(b, x) {
			return true
		}
	}
	return false
}

func joinStrings(values []string, separator string) string {
	if len(values) == 0 {
		return "none"
	}
	return strings.Join(values, separator)
}

func allChecksOK(checks []Check) bool {
	return countPassed(checks) == len(checks)
}

func countPassed(checks []Check) int {
	passed := 0
	for _, check := range checks {
		if check.OK {
			passed++
		}
	}
	return passed
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}
