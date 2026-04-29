// odrl_dpv_risk_ranked.go
//
// A self-contained Go translation of `examples/odrl-dpv-risk-ranked.n3`.
//
// The original N3 example models a Terms-of-Service style agreement using
// policy and privacy-risk vocabulary terms. It links machine-readable
// permissions and prohibitions to human-readable clauses, derives risks from
// missing or weak safeguards, scores those risks, and emits a ranked report.
//
// This program is intentionally not a generic RDF, ODRL, DPV, or N3 reasoner.
// Instead, it translates the concrete facts and rules from the fixture into
// typed Go data structures and explicit inference functions. That keeps the
// rule mechanics visible while preserving the deterministic ranked output style
// of the N3 `log:outputString` section.
//
// Run:
//
//	go run odrl_dpv_risk_ranked.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
)

const eyelingoExampleName = "odrl_dpv_risk_ranked"

const (
	// Action names are kept as vocabulary-like strings to mirror the N3 source.
	ActionRemoveAccount = "tosl:removeAccount"
	ActionChangeTerms   = "tosl:changeTerms"
	ActionShareData     = "tosl:shareData"
	ActionExportData    = "tosl:exportData"
	ActionInform        = "odrl:inform"

	LeftNoticeDays = "tosl:noticeDays"
	LeftConsent    = "tosl:consent"

	OperatorGTE = "odrl:gteq"
	OperatorEQ  = "odrl:eq"
)

// Dataset is the complete translated fixture: consumer needs, agreement facts,
// ODRL policy rules, linked clause text, and the process context label.
type Dataset struct {
	Consumer  ConsumerProfile
	Agreement Agreement
	Process   ProcessContext
}

// ConsumerProfile corresponds to :ConsumerExample and its linked :Need values.
type ConsumerProfile struct {
	ID    string
	Title string
	Needs map[string]Need
}

// Need holds the consumer preference used as an input weight for risk scoring.
// MinNoticeDays is only populated for the prior-notice need.
type Need struct {
	ID            string
	Importance    int
	MinNoticeDays int
	Description   string
}

// Agreement contains the human-readable agreement label plus its ODRL policy.
type Agreement struct {
	ID      string
	Title   string
	Policy  Policy
	Clauses map[string]Clause
}

// Policy mirrors :Policy1 with permissions and prohibitions keyed by resource ID.
type Policy struct {
	ID           string
	Permissions  map[string]Permission
	Prohibitions map[string]Prohibition
}

// Permission is the subset of ODRL permission structure used by this fixture.
type Permission struct {
	ID          string
	Assigner    string
	Assignee    string
	Action      string
	Target      string
	ClauseID    string
	Duties      []Duty
	Constraints []Constraint
}

// Prohibition is the subset of ODRL prohibition structure used by this fixture.
type Prohibition struct {
	ID       string
	Assigner string
	Assignee string
	Action   string
	Target   string
	ClauseID string
}

// Duty and Constraint model nested ODRL duty/constraint blank nodes.
type Duty struct {
	Action      string
	Constraints []Constraint
}

type Constraint struct {
	LeftOperand  string
	Operator     string
	RightOperand Value
}

// Value keeps constraint operands simple while still distinguishing int and bool.
type Value struct {
	Int  int
	Bool bool
	Kind string
}

// Clause stores the text outside the policy graph. The N3 rules only match ODRL
// structure inside the graph and then join to clause text for explanation.
type Clause struct {
	ID   string
	Text string
}

// ProcessContext corresponds to :ProcessContext1, which receives dpv:hasRisk
// links in the N3 consequences.
type ProcessContext struct {
	ID     string
	Title  string
	Source string
}

// Risk is the materialized DPV/DPV-RISK result of one translated N3 rule.
type Risk struct {
	ID           string
	SourceRuleID string
	Clause       Clause
	ViolatesNeed string
	Description  string
	ScoreRaw     int
	Score        int
	Severity     string
	RiskLevel    string
	Types        []string
	Consequences []string
	Impacts      []string
	Mitigations  []Mitigation
}

// Mitigation corresponds to a dpv:RiskMitigationMeasure in the N3 output.
type Mitigation struct {
	ID          string
	Description string
	SuggestAdd  string
}

// dataset hard-codes the concrete facts from the N3 file.
func dataset() Dataset {
	return Dataset{
		Consumer: ConsumerProfile{
			ID:    ":ConsumerExample",
			Title: "Example consumer profile",
			Needs: map[string]Need{
				":Need_DataCannotBeRemoved": {
					ID:          ":Need_DataCannotBeRemoved",
					Importance:  20,
					Description: "Provider must not remove the consumer account/data.",
				},
				":Need_ChangeOnlyWithPriorNotice": {
					ID:            ":Need_ChangeOnlyWithPriorNotice",
					Importance:    15,
					MinNoticeDays: 14,
					Description:   "Agreement may change only with prior notice (>= 14 days).",
				},
				":Need_NoSharingWithoutConsent": {
					ID:          ":Need_NoSharingWithoutConsent",
					Importance:  12,
					Description: "No data sharing without explicit consent.",
				},
				":Need_DataPortability": {
					ID:          ":Need_DataPortability",
					Importance:  10,
					Description: "Consumer must be able to export their data.",
				},
			},
		},

		Agreement: Agreement{
			ID:    ":Agreement1",
			Title: "Example Agreement",
			Policy: Policy{
				ID: ":Policy1",
				Permissions: map[string]Permission{
					":PermDeleteAccount": {
						ID:       ":PermDeleteAccount",
						Assigner: ":Provider",
						Assignee: ":ConsumerExample",
						Action:   ActionRemoveAccount,
						Target:   ":UserAccount",
						ClauseID: ":ClauseC1",
					},
					":PermChangeTerms": {
						ID:       ":PermChangeTerms",
						Assigner: ":Provider",
						Assignee: ":ConsumerExample",
						Action:   ActionChangeTerms,
						Target:   ":AgreementText",
						ClauseID: ":ClauseC2",
						Duties: []Duty{
							{
								Action: ActionInform,
								Constraints: []Constraint{
									{
										LeftOperand:  LeftNoticeDays,
										Operator:     OperatorGTE,
										RightOperand: intValue(3),
									},
								},
							},
						},
					},
					":PermShareData": {
						ID:       ":PermShareData",
						Assigner: ":Provider",
						Assignee: ":ConsumerExample",
						Action:   ActionShareData,
						Target:   ":UserData",
						ClauseID: ":ClauseC3",
					},
				},
				Prohibitions: map[string]Prohibition{
					":ProhibitExportData": {
						ID:       ":ProhibitExportData",
						Assigner: ":Provider",
						Assignee: ":ConsumerExample",
						Action:   ActionExportData,
						Target:   ":UserData",
						ClauseID: ":ClauseC4",
					},
				},
			},
			Clauses: map[string]Clause{
				":ClauseC1": {
					ID:   "C1",
					Text: "Provider may remove the user account (and associated data) at its discretion.",
				},
				":ClauseC2": {
					ID:   "C2",
					Text: "Provider may change terms by informing users at least 3 days in advance.",
				},
				":ClauseC3": {
					ID:   "C3",
					Text: "Provider may share user data with partners for business purposes.",
				},
				":ClauseC4": {
					ID:   "C4",
					Text: "Users are not permitted to export their data.",
				},
			},
		},

		Process: ProcessContext{
			ID:     ":ProcessContext1",
			Title:  "Service operation under Agreement1",
			Source: ":Agreement1",
		},
	}
}

func intValue(value int) Value {
	return Value{Int: value, Kind: "int"}
}

func boolValue(value bool) Value {
	return Value{Bool: value, Kind: "bool"}
}

// infer runs the translated N3 rules in a deterministic order.
func infer(data Dataset) ([]Risk, error) {
	var risks []Risk

	if risk, ok, err := inferAccountRemovalWithoutSafeguards(data); err != nil {
		return nil, err
	} else if ok {
		risks = append(risks, risk)
	}

	if risk, ok, err := inferChangeTermsNoticeTooShort(data); err != nil {
		return nil, err
	} else if ok {
		risks = append(risks, risk)
	}

	if risk, ok, err := inferShareDataWithoutConsent(data); err != nil {
		return nil, err
	} else if ok {
		risks = append(risks, risk)
	}

	if risk, ok, err := inferExportProhibited(data); err != nil {
		return nil, err
	} else if ok {
		risks = append(risks, risk)
	}

	for i := range risks {
		classifyRisk(&risks[i])
	}

	// The N3 report sorts by inverse score (1000 - score), making the highest
	// score print first. Clause ID provides deterministic ordering after that.
	sort.SliceStable(risks, func(i, j int) bool {
		if risks[i].Score != risks[j].Score {
			return risks[i].Score > risks[j].Score
		}
		return risks[i].Clause.ID < risks[j].Clause.ID
	})

	return risks, nil
}

// R1: remove account/data WITHOUT a notice constraint AND WITHOUT an inform duty.
func inferAccountRemovalWithoutSafeguards(data Dataset) (Risk, bool, error) {
	need, ok := data.Consumer.Needs[":Need_DataCannotBeRemoved"]
	if !ok {
		return Risk{}, false, nil
	}

	perm, ok := data.Agreement.Policy.Permissions[":PermDeleteAccount"]
	if !ok || perm.Action != ActionRemoveAccount {
		return Risk{}, false, nil
	}

	// These checks correspond to the N3 log:notIncludes patterns.
	if hasConstraint(perm.Constraints, LeftNoticeDays, "") || hasDutyAction(perm.Duties, ActionInform) {
		return Risk{}, false, nil
	}

	clause, err := clauseFor(data, perm.ClauseID)
	if err != nil {
		return Risk{}, false, err
	}

	raw := 90 + need.Importance
	return Risk{
		ID:           "_:risk1",
		SourceRuleID: perm.ID,
		Clause:       clause,
		ViolatesNeed: need.ID,
		Description: fmt.Sprintf(
			"Risk: account/data removal is permitted without notice safeguards (no notice constraint and no duty to inform). Clause %s: %s",
			clause.ID,
			clause.Text,
		),
		ScoreRaw: raw,
		Types: []string{
			"dpv:Risk",
			"risk:UnwantedDataDeletion",
			"risk:DataUnavailable",
			"risk:DataErasureError",
			"risk:DataLoss",
		},
		Consequences: []string{
			"risk:DataLoss",
			"risk:DataUnavailable",
			"risk:CustomerConfidenceLoss",
		},
		Impacts: []string{
			"risk:FinancialLoss",
			"risk:NonMaterialDamage",
		},
		Mitigations: []Mitigation{
			{
				ID:          "_:m11",
				Description: "Add a notice constraint (minimum noticeDays) before account removal.",
				SuggestAdd:  ":PermDeleteAccount odrl:constraint [ a odrl:Constraint ; odrl:leftOperand tosl:noticeDays ; odrl:operator odrl:gteq ; odrl:rightOperand 14 ] .",
			},
			{
				ID:          "_:m21",
				Description: "Add a duty to inform the consumer prior to account removal.",
				SuggestAdd:  ":PermDeleteAccount odrl:duty [ a odrl:Duty ; odrl:action odrl:inform ] .",
			},
		},
	}, true, nil
}

// R2: terms may change with a notice period below the consumer's requirement.
func inferChangeTermsNoticeTooShort(data Dataset) (Risk, bool, error) {
	need, ok := data.Consumer.Needs[":Need_ChangeOnlyWithPriorNotice"]
	if !ok {
		return Risk{}, false, nil
	}

	perm, ok := data.Agreement.Policy.Permissions[":PermChangeTerms"]
	if !ok || perm.Action != ActionChangeTerms {
		return Risk{}, false, nil
	}

	days, ok := noticeDaysFromInformDuty(perm.Duties)
	if !ok || days >= need.MinNoticeDays {
		return Risk{}, false, nil
	}

	clause, err := clauseFor(data, perm.ClauseID)
	if err != nil {
		return Risk{}, false, err
	}

	raw := 70 + need.Importance
	return Risk{
		ID:           "_:risk2",
		SourceRuleID: perm.ID,
		Clause:       clause,
		ViolatesNeed: need.ID,
		Description: fmt.Sprintf(
			"Risk: terms may change with notice (%d days) below consumer requirement (%d days). Clause %s: %s",
			days,
			need.MinNoticeDays,
			clause.ID,
			clause.Text,
		),
		ScoreRaw: raw,
		Types: []string{
			"dpv:Risk",
			"risk:PolicyRisk",
			"risk:CustomerConfidenceLoss",
		},
		Consequences: []string{"risk:CustomerConfidenceLoss"},
		Impacts:      []string{"risk:NonMaterialDamage"},
		Mitigations: []Mitigation{
			{
				ID:          "_:m12",
				Description: "Increase minimum noticeDays in the inform duty to meet the consumer requirement.",
				SuggestAdd:  ":PermChangeTerms odrl:duty [ a odrl:Duty ; odrl:action odrl:inform ; odrl:constraint [ a odrl:Constraint ; odrl:leftOperand tosl:noticeDays ; odrl:operator odrl:gteq ; odrl:rightOperand 14 ] ] .",
			},
		},
	}, true, nil
}

// R3: data sharing is permitted WITHOUT an explicit consent constraint.
func inferShareDataWithoutConsent(data Dataset) (Risk, bool, error) {
	need, ok := data.Consumer.Needs[":Need_NoSharingWithoutConsent"]
	if !ok {
		return Risk{}, false, nil
	}

	perm, ok := data.Agreement.Policy.Permissions[":PermShareData"]
	if !ok || perm.Action != ActionShareData {
		return Risk{}, false, nil
	}

	// This is the N3 log:notIncludes pattern for odrl:constraint consent == true.
	if hasBoolConstraint(perm.Constraints, LeftConsent, OperatorEQ, true) {
		return Risk{}, false, nil
	}

	clause, err := clauseFor(data, perm.ClauseID)
	if err != nil {
		return Risk{}, false, err
	}

	raw := 85 + need.Importance
	return Risk{
		ID:           "_:risk3",
		SourceRuleID: perm.ID,
		Clause:       clause,
		ViolatesNeed: need.ID,
		Description: fmt.Sprintf(
			"Risk: user data sharing is permitted without an explicit consent constraint. Clause %s: %s",
			clause.ID,
			clause.Text,
		),
		ScoreRaw: raw,
		Types: []string{
			"dpv:Risk",
			"risk:UnwantedDisclosureData",
			"risk:CustomerConfidenceLoss",
		},
		Consequences: []string{"risk:CustomerConfidenceLoss"},
		Impacts: []string{
			"risk:NonMaterialDamage",
			"risk:FinancialLoss",
		},
		Mitigations: []Mitigation{
			{
				ID:          "_:m13",
				Description: "Add an explicit consent constraint before data sharing.",
				SuggestAdd:  ":PermShareData odrl:constraint [ a odrl:Constraint ; odrl:leftOperand tosl:consent ; odrl:operator odrl:eq ; odrl:rightOperand true ] .",
			},
		},
	}, true, nil
}

// R4: data export is prohibited, which conflicts with the portability need.
func inferExportProhibited(data Dataset) (Risk, bool, error) {
	need, ok := data.Consumer.Needs[":Need_DataPortability"]
	if !ok {
		return Risk{}, false, nil
	}

	prohibition, ok := data.Agreement.Policy.Prohibitions[":ProhibitExportData"]
	if !ok || prohibition.Action != ActionExportData {
		return Risk{}, false, nil
	}

	clause, err := clauseFor(data, prohibition.ClauseID)
	if err != nil {
		return Risk{}, false, err
	}

	raw := 60 + need.Importance
	return Risk{
		ID:           "_:risk4",
		SourceRuleID: prohibition.ID,
		Clause:       clause,
		ViolatesNeed: need.ID,
		Description: fmt.Sprintf(
			"Risk: portability is restricted because exporting user data is prohibited. Clause %s: %s",
			clause.ID,
			clause.Text,
		),
		ScoreRaw: raw,
		Types: []string{
			"dpv:Risk",
			"risk:PolicyRisk",
			"risk:CustomerConfidenceLoss",
		},
		Consequences: []string{"risk:CustomerConfidenceLoss"},
		Impacts:      []string{"risk:NonMaterialDamage"},
		Mitigations: []Mitigation{
			{
				ID:          "_:m14",
				Description: "Add a permission allowing data export (or remove the prohibition) to support portability.",
				SuggestAdd:  ":Policy1 odrl:permission [ a odrl:Permission ; odrl:assigner :Provider ; odrl:assignee :ConsumerExample ; odrl:action tosl:exportData ; odrl:target :UserData ] .",
			},
		},
	}, true, nil
}

// classifyRisk implements the N3 score normalization and severity/risk-level rules.
func classifyRisk(risk *Risk) {
	if risk.ScoreRaw > 100 {
		risk.Score = 100
	} else {
		risk.Score = risk.ScoreRaw
	}

	switch {
	case risk.Score > 79:
		risk.Severity = "risk:HighSeverity"
		risk.RiskLevel = "risk:HighRisk"
	case risk.Score < 80 && risk.Score > 49:
		risk.Severity = "risk:ModerateSeverity"
		risk.RiskLevel = "risk:ModerateRisk"
	default:
		risk.Severity = "risk:LowSeverity"
		risk.RiskLevel = "risk:LowRisk"
	}
}

func clauseFor(data Dataset, clauseResourceID string) (Clause, error) {
	clause, ok := data.Agreement.Clauses[clauseResourceID]
	if !ok {
		return Clause{}, fmt.Errorf("clause not found: %s", clauseResourceID)
	}
	return clause, nil
}

func hasConstraint(constraints []Constraint, leftOperand, operator string) bool {
	for _, constraint := range constraints {
		if constraint.LeftOperand != leftOperand {
			continue
		}
		if operator == "" || constraint.Operator == operator {
			return true
		}
	}
	return false
}

func hasBoolConstraint(constraints []Constraint, leftOperand, operator string, value bool) bool {
	for _, constraint := range constraints {
		if constraint.LeftOperand == leftOperand &&
			constraint.Operator == operator &&
			constraint.RightOperand.Kind == "bool" &&
			constraint.RightOperand.Bool == value {
			return true
		}
	}
	return false
}

func hasDutyAction(duties []Duty, action string) bool {
	for _, duty := range duties {
		if duty.Action == action {
			return true
		}
	}
	return false
}

func noticeDaysFromInformDuty(duties []Duty) (int, bool) {
	for _, duty := range duties {
		if duty.Action != ActionInform {
			continue
		}
		for _, constraint := range duty.Constraints {
			if constraint.LeftOperand == LeftNoticeDays && constraint.RightOperand.Kind == "int" {
				return constraint.RightOperand.Int, true
			}
		}
	}
	return 0, false
}

// renderRankedReport mirrors the N3 `log:outputString` formatting.
func renderRankedReport(data Dataset, risks []Risk) {
	fmt.Println("# Ranked DPV Risk Report")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("Agreement: %s\n", data.Agreement.Title)
	fmt.Printf("Profile: %s\n\n", data.Consumer.Title)

	for _, risk := range risks {
		fmt.Printf(
			"score=%d (%s, %s) clause %s\n %s\n\n",
			risk.Score,
			risk.RiskLevel,
			risk.Severity,
			risk.Clause.ID,
			risk.Description,
		)

		for _, mitigation := range risk.Mitigations {
			fmt.Printf(" - mitigation for clause %s: %s\n", risk.Clause.ID, mitigation.Description)
		}
	}
}

func renderReason(data Dataset, risks []Risk) {
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("The agreement policy is scanned for permissions and prohibitions that conflict with the consumer profile needs.")
	fmt.Println("Each triggered rule derives a risk row with a normalized score, a source clause, and one or more mitigation measures.")
	fmt.Println("Rows are sorted by descending score so the highest-risk clauses are reviewed first.")
}

func renderChecks(data Dataset, risks []Risk) {
	levelCounts := riskLevelCounts(risks)
	minScore, maxScore := scoreRange(risks)
	fmt.Println()
	fmt.Println("## Check")
	fmt.Printf("C1 %s - %d risk rows were derived.\n", checkStatus(len(risks) == 4), len(risks))
	fmt.Printf("C2 %s - ranked output is in descending score order.\n", checkStatus(rankedDescending(risks)))
	fmt.Printf("C3 %s - score range is %d to %d.\n", checkStatus(minScore == 70 && maxScore == 100), minScore, maxScore)
	fmt.Printf("C4 %s - high=%d moderate=%d low=%d risk levels were derived.\n", checkStatus(levelCounts["risk:HighRisk"] == 3 && levelCounts["risk:ModerateRisk"] == 1), levelCounts["risk:HighRisk"], levelCounts["risk:ModerateRisk"], levelCounts["risk:LowRisk"])
	fmt.Printf("C5 %s - %d mitigation measures were generated.\n", checkStatus(countMitigations(risks) == 5), countMitigations(risks))
}

func checkStatus(ok bool) string {
	if ok {
		return "OK"
	}
	return "FAIL"
}

// renderAuditDetails is extra Go-side output that makes the translation easier
// to inspect without changing the primary ranked report above.
func renderAuditDetails(data Dataset, risks []Risk) {
	levelCounts := riskLevelCounts(risks)
	minScore, maxScore := scoreRange(risks)

	fmt.Println()
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("process context : %s (%s -> %s)\n", data.Process.Title, data.Process.ID, data.Process.Source)
	fmt.Printf("consumer profile : %s needs=%d\n", data.Consumer.ID, len(data.Consumer.Needs))
	fmt.Printf("agreement : %s policy=%s\n", data.Agreement.ID, data.Agreement.Policy.ID)
	fmt.Printf("policy graph size : permissions=%d prohibitions=%d clauses=%d duties=%d constraints=%d\n",
		len(data.Agreement.Policy.Permissions),
		len(data.Agreement.Policy.Prohibitions),
		len(data.Agreement.Clauses),
		countDuties(data),
		countConstraints(data),
	)
	fmt.Printf("risks derived : %d\n", len(risks))
	fmt.Printf("risk levels : high=%d moderate=%d low=%d\n", levelCounts["risk:HighRisk"], levelCounts["risk:ModerateRisk"], levelCounts["risk:LowRisk"])
	fmt.Printf("score range : min=%d max=%d\n", minScore, maxScore)
	fmt.Printf("ranked descending : %s\n", yesNo(rankedDescending(risks)))
	fmt.Printf("mitigation measures : %d\n", countMitigations(risks))
	fmt.Printf("R1 account removal without safeguards : %s\n", yesNo(ruleTriggered(risks, ":PermDeleteAccount")))
	fmt.Printf("R2 change terms notice too short : %s\n", yesNo(ruleTriggered(risks, ":PermChangeTerms")))
	fmt.Printf("R3 share data without consent : %s\n", yesNo(ruleTriggered(risks, ":PermShareData")))
	fmt.Printf("R4 export data prohibited : %s\n", yesNo(ruleTriggered(risks, ":ProhibitExportData")))

	fmt.Println("need weights:")
	for _, need := range sortedNeeds(data.Consumer.Needs) {
		if need.MinNoticeDays > 0 {
			fmt.Printf(" - %s importance=%d minNoticeDays=%d\n", need.ID, need.Importance, need.MinNoticeDays)
		} else {
			fmt.Printf(" - %s importance=%d\n", need.ID, need.Importance)
		}
	}

	fmt.Println("policy actions:")
	for _, line := range policyActionLines(data) {
		fmt.Printf(" - %s\n", line)
	}

	fmt.Println("derived risk rows:")
	for rank, risk := range risks {
		fmt.Printf(
			" #%d %s clause=%s raw=%d normalized=%d level=%s severity=%s violates=%s source=%s mitigations=%d\n",
			rank+1,
			risk.ID,
			risk.Clause.ID,
			risk.ScoreRaw,
			risk.Score,
			risk.RiskLevel,
			risk.Severity,
			risk.ViolatesNeed,
			risk.SourceRuleID,
			len(risk.Mitigations),
		)
	}
}

// countDuties counts nested ODRL duty blocks in all translated permissions.
func countDuties(data Dataset) int {
	count := 0
	for _, permission := range data.Agreement.Policy.Permissions {
		count += len(permission.Duties)
	}
	return count
}

// countConstraints counts direct permission constraints plus duty constraints.
func countConstraints(data Dataset) int {
	count := 0
	for _, permission := range data.Agreement.Policy.Permissions {
		count += len(permission.Constraints)
		for _, duty := range permission.Duties {
			count += len(duty.Constraints)
		}
	}
	return count
}

// countMitigations totals all mitigation suggestions across derived risks.
func countMitigations(risks []Risk) int {
	count := 0
	for _, risk := range risks {
		count += len(risk.Mitigations)
	}
	return count
}

// riskLevelCounts groups the classified risks by DPV-RISK level.
func riskLevelCounts(risks []Risk) map[string]int {
	counts := map[string]int{
		"risk:HighRisk":     0,
		"risk:ModerateRisk": 0,
		"risk:LowRisk":      0,
	}
	for _, risk := range risks {
		counts[risk.RiskLevel]++
	}
	return counts
}

// scoreRange returns the normalized score interval for the derived risk list.
func scoreRange(risks []Risk) (int, int) {
	if len(risks) == 0 {
		return 0, 0
	}
	minScore := risks[0].Score
	maxScore := risks[0].Score
	for _, risk := range risks[1:] {
		if risk.Score < minScore {
			minScore = risk.Score
		}
		if risk.Score > maxScore {
			maxScore = risk.Score
		}
	}
	return minScore, maxScore
}

// rankedDescending verifies the report order produced after the score sort.
func rankedDescending(risks []Risk) bool {
	for i := 1; i < len(risks); i++ {
		if risks[i].Score > risks[i-1].Score {
			return false
		}
		if risks[i].Score == risks[i-1].Score && risks[i].Clause.ID < risks[i-1].Clause.ID {
			return false
		}
	}
	return true
}

// ruleTriggered reports whether a risk was emitted from a specific permission or
// prohibition resource.
func ruleTriggered(risks []Risk, sourceRuleID string) bool {
	for _, risk := range risks {
		if risk.SourceRuleID == sourceRuleID {
			return true
		}
	}
	return false
}

// sortedNeeds returns consumer needs in stable ID order for deterministic audit output.
func sortedNeeds(needs map[string]Need) []Need {
	ids := make([]string, 0, len(needs))
	for id := range needs {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	out := make([]Need, 0, len(ids))
	for _, id := range ids {
		out = append(out, needs[id])
	}
	return out
}

// policyActionLines returns a deterministic inventory of ODRL permissions and prohibitions.
func policyActionLines(data Dataset) []string {
	var lines []string

	permissionIDs := make([]string, 0, len(data.Agreement.Policy.Permissions))
	for id := range data.Agreement.Policy.Permissions {
		permissionIDs = append(permissionIDs, id)
	}
	sort.Strings(permissionIDs)
	for _, id := range permissionIDs {
		permission := data.Agreement.Policy.Permissions[id]
		lines = append(lines, fmt.Sprintf("permission %s action=%s target=%s clause=%s duties=%d constraints=%d", permission.ID, permission.Action, permission.Target, permission.ClauseID, len(permission.Duties), len(permission.Constraints)))
	}

	prohibitionIDs := make([]string, 0, len(data.Agreement.Policy.Prohibitions))
	for id := range data.Agreement.Policy.Prohibitions {
		prohibitionIDs = append(prohibitionIDs, id)
	}
	sort.Strings(prohibitionIDs)
	for _, id := range prohibitionIDs {
		prohibition := data.Agreement.Policy.Prohibitions[id]
		lines = append(lines, fmt.Sprintf("prohibition %s action=%s target=%s clause=%s", prohibition.ID, prohibition.Action, prohibition.Target, prohibition.ClauseID))
	}

	return lines
}

// yesNo prints booleans in the same compact style used by the other examples.
func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func main() {
	data := exampleinput.Load(eyelingoExampleName, dataset())
	risks, err := infer(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ODRL/DPV risk inference failed: %v\n", err)
		os.Exit(1)
	}

	renderRankedReport(data, risks)
	renderReason(data, risks)
	renderChecks(data, risks)
	renderAuditDetails(data, risks)

	// Keep boolValue referenced because the fixture uses boolean operands in the
	// missing-consent check. A future dataset variant can add that safeguard by
	// appending Constraint{LeftOperand: LeftConsent, Operator: OperatorEQ,
	// RightOperand: boolValue(true)} to :PermShareData.
	_ = boolValue
}
