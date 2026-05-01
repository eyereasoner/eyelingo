// doctor_advice_work_conflict.go
//
// A compact Go translation inspired by Eyeling's
// `examples/doctor-advice-work-conflict.n3`.
//
// The example keeps conflicting permit/deny conclusions visible and then applies
// deterministic conflict-resolution rules for home versus office work.
//
// Run:
//
//	go run examples/doctor_advice_work_conflict.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"sort"
	"strings"
)

const eyelingoExampleName = "doctor_advice_work_conflict"

type Dataset struct {
	CaseName         string
	Question         string
	Person           Person
	DoctorCanDoJobs  []string
	SubclassOf       map[string]string
	Requests         []Request
	ConflictPolicies ConflictPolicies
	Expected         Expected
}

type Person struct {
	Name      string
	Condition string
}

type Request struct {
	ID       string
	Job      string
	Location string
}

type ConflictPolicies struct {
	SickWorkDefault       string
	SickOfficeWorkDefault string
	HomeProgrammingWork   string
}

type Expected struct {
	OverallDecision string
}

type RequestDecision struct {
	Request           Request
	RawDecisions      []string
	DecisionStatus    string
	EffectiveDecision string
	Reasons           []string
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	HealthStatus string
	Decisions    []RequestDecision
	Overall      string
	Checks       []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
}

func derive(ds Dataset) Analysis {
	health := "Well"
	if ds.Person.Condition != "" {
		health = "Sick"
	}
	decisions := make([]RequestDecision, 0, len(ds.Requests))
	for _, request := range ds.Requests {
		raw := make([]string, 0)
		reasons := make([]string, 0)
		if contains(ds.DoctorCanDoJobs, request.Job) || canDoSuper(ds.DoctorCanDoJobs, request.Job, ds.SubclassOf) {
			raw = append(raw, "Permit")
			reasons = append(reasons, "doctor statement permits the requested job")
		}
		if health == "Sick" && isKindOf(request.Job, "Work", ds.SubclassOf) {
			raw = append(raw, ds.ConflictPolicies.SickWorkDefault)
			reasons = append(reasons, "general sick-work policy denies work")
		}
		if health == "Sick" && request.Location == "Office" && isKindOf(request.Job, "Work", ds.SubclassOf) {
			raw = append(raw, ds.ConflictPolicies.SickOfficeWorkDefault)
			reasons = append(reasons, "office policy denies work to avoid infecting colleagues")
		}
		raw = uniqueSorted(raw)
		status := decisionStatus(raw)
		effective := resolve(request, status, ds.ConflictPolicies)
		decisions = append(decisions, RequestDecision{Request: request, RawDecisions: raw, DecisionStatus: status, EffectiveDecision: effective, Reasons: reasons})
	}
	overall := "NotSpecified"
	if hasEffective(decisions, "Home", "Permit") && hasEffective(decisions, "Office", "Deny") {
		overall = "RemoteWorkOnly"
	}
	checks := []Check{
		{ID: "C1", OK: health == "Sick", Text: fmt.Sprintf("%s is classified as Sick from condition %s", ds.Person.Name, ds.Person.Condition)},
		{ID: "C2", OK: hasStatus(decisions, "Request_Jos_Prog_Home", "BothPermitDeny"), Text: "home programming request keeps both Permit and Deny before resolution"},
		{ID: "C3", OK: hasEffectiveByID(decisions, "Request_Jos_Prog_Home", "Permit"), Text: "home programming conflict resolves to Permit"},
		{ID: "C4", OK: hasStatus(decisions, "Request_Jos_Prog_Office", "BothPermitDeny"), Text: "office programming request keeps both Permit and Deny before resolution"},
		{ID: "C5", OK: hasEffectiveByID(decisions, "Request_Jos_Prog_Office", "Deny"), Text: "office programming conflict resolves to Deny"},
		{ID: "C6", OK: overall == ds.Expected.OverallDecision, Text: fmt.Sprintf("overall work decision is %s", overall)},
	}
	return Analysis{HealthStatus: health, Decisions: decisions, Overall: overall, Checks: checks}
}

func canDoSuper(canDo []string, job string, subclasses map[string]string) bool {
	for _, allowed := range canDo {
		if isKindOf(job, allowed, subclasses) {
			return true
		}
	}
	return false
}

func isKindOf(child string, parent string, subclasses map[string]string) bool {
	if child == parent {
		return true
	}
	for current, ok := child, true; ok; current, ok = subclasses[current] {
		if current == parent {
			return true
		}
		if next, exists := subclasses[current]; exists && next == parent {
			return true
		}
	}
	return false
}

func decisionStatus(raw []string) string {
	hasPermit := contains(raw, "Permit")
	hasDeny := contains(raw, "Deny")
	switch {
	case hasPermit && hasDeny:
		return "BothPermitDeny"
	case hasPermit:
		return "PermitOnly"
	case hasDeny:
		return "DenyOnly"
	default:
		return "Neither"
	}
}

func resolve(request Request, status string, policies ConflictPolicies) string {
	switch status {
	case "PermitOnly":
		return "Permit"
	case "DenyOnly":
		return "Deny"
	case "Neither":
		return "NotSpecified"
	case "BothPermitDeny":
		if request.Location == "Home" && request.Job == "ProgrammingWork" {
			return policies.HomeProgrammingWork
		}
		if request.Location == "Office" {
			return "Deny"
		}
		return "Undecided"
	default:
		return "Undecided"
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	for _, value := range values {
		seen[value] = true
	}
	out := make([]string, 0, len(seen))
	for value := range seen {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func hasStatus(decisions []RequestDecision, id string, status string) bool {
	for _, decision := range decisions {
		if decision.Request.ID == id && decision.DecisionStatus == status {
			return true
		}
	}
	return false
}

func hasEffectiveByID(decisions []RequestDecision, id string, effective string) bool {
	for _, decision := range decisions {
		if decision.Request.ID == id && decision.EffectiveDecision == effective {
			return true
		}
	}
	return false
}

func hasEffective(decisions []RequestDecision, location string, effective string) bool {
	for _, decision := range decisions {
		if decision.Request.Location == location && decision.EffectiveDecision == effective {
			return true
		}
	}
	return false
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
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

func printAnswer(ds Dataset, analysis Analysis) {
	fmt.Println("# Doctor Advice Work Conflict")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("overall decision for %s : %s\n", ds.Person.Name, analysis.Overall)
	for _, decision := range analysis.Decisions {
		fmt.Printf("%s : raw=%s status=%s effective=%s\n", decision.Request.ID, strings.Join(decision.RawDecisions, "+"), decision.DecisionStatus, decision.EffectiveDecision)
	}
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason")
	fmt.Printf("%s has %s, so the health context classifies the agent as %s.\n", ds.Person.Name, ds.Person.Condition, analysis.HealthStatus)
	fmt.Println("The doctor's statement permits ProgrammingWork, but the general sick-work policy also denies work, so the raw closure deliberately keeps the conflict visible.")
	fmt.Println("The conflict resolver permits sick ProgrammingWork at Home but denies Office work because of colleague-infection risk.")
	fmt.Println("Since Home is permitted and Office is denied, the combined recommendation is RemoteWorkOnly.")
	fmt.Println()
}
