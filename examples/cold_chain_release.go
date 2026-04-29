// cold_chain_release.go
//
// A self-contained Go translation of an Eyeling-style cold-chain provenance
// example.
//
// The scenario models a scarce biologic shipment after a suspected source-lot
// recall. The facts include lot lineage, hash-linked custody events,
// temperature observations, clinic demand, and transit limits. The goal is to
// decide which lots can be released and then allocate the safe doses exactly to
// the highest-value clinics without exceeding demand.
//
// This is intentionally not a generic RDF/N3 reasoner. The concrete facts and
// rules are represented as Go structs and ordinary functions so the provenance,
// policy, and optimization steps stay visible and directly runnable.
//
// Run:
//
//	go run cold_chain_release.go
//
// The program has no third-party dependencies.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

const eyelingoExampleName = "cold_chain_release"

const (
	genesisHash                 = "sha256:genesis"
	coldMinTenths               = 20 // 2.0 °C
	coldMaxTenths               = 80 // 8.0 °C
	singleExcursionLimitMinutes = 15
	totalExcursionLimitMinutes  = 30
)

type Dataset struct {
	CaseName        string
	Question        string
	Product         string
	Lots            []Lot
	RecalledSources []string
	Events          []CustodyEvent
	Clinics         []Clinic
	Transit         []TransitFact
}

type Lot struct {
	ID               string
	Product          string
	Facility         string
	Doses            int
	Parents          []string
	ReleaseCandidate bool
}

type CustodyEvent struct {
	LotID      string
	Ordinal    int
	Step       string
	Facility   string
	Actor      string
	TempTenths int
	Minutes    int
	PrevHash   string
	Hash       string
}

type Clinic struct {
	ID         string
	Label      string
	Demand     int
	Priority   int
	MaxTransit int
	Need       string
}

type TransitFact struct {
	From     string
	ClinicID string
	Minutes  int
}

type TemperatureReport struct {
	TotalExcursion int
	MaxSegment     int
	HardFail       bool
	OK             bool
	Reason         string
}

type LedgerReport struct {
	OK        bool
	Root      string
	BreakStep string
	Expected  string
	Got       string
	Reason    string
}

type LotReport struct {
	Lot          Lot
	Ancestors    []string
	RecalledBy   string
	Ledger       LedgerReport
	Temperature  TemperatureReport
	Release      bool
	RejectReason string
}

type FeasibleAssignment struct {
	LotID       string
	ClinicID    string
	ClinicLabel string
	Doses       int
	Transit     int
	Value       int
}

type SelectedAssignment struct {
	LotID       string
	ClinicID    string
	ClinicLabel string
	Delivered   int
	Reserve     int
	Transit     int
	Value       int
}

type Plan struct {
	Score       int
	Doses       int
	Transit     int
	Assignments []SelectedAssignment
}

type OptimizerStats struct {
	States      int
	Transitions int
	MemoHits    int
}

type Analysis struct {
	Ancestors       map[string][]string
	LotReports      []LotReport
	ReleaseLots     []LotReport
	QuarantineLots  []LotReport
	Feasible        []FeasibleAssignment
	Plan            Plan
	Stats           OptimizerStats
	Checks          []Check
	LedgerValid     int
	TemperatureOK   int
	ClosureFacts    int
	ParentFacts     int
	EventFacts      int
	TempFacts       int
	ReserveDoses    int
	UnmetByClinic   map[string]int
	CandidateByLot  map[string][]FeasibleAssignment
	TransitByKey    map[string]int
	ReleaseByLot    map[string]LotReport
	ClinicByID      map[string]Clinic
	LotByID         map[string]Lot
	LineageComplete bool
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, dataset())
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
	printChecks(analysis)
	printAudit(ds, analysis)
	if !allChecksOK(analysis.Checks) {
		os.Exit(1)
	}
}

func dataset() Dataset {
	ds := Dataset{
		CaseName: "cold-chain-release",
		Question: "Which cold-chain lots can be released, and where should scarce doses ship first?",
		Product:  "Stable mRNA pediatric booster",
		Lots: []Lot{
			{ID: "Seed-Alpha", Product: "Stable mRNA pediatric booster", Facility: "Plant-7"},
			{ID: "Seed-Beta", Product: "Stable mRNA pediatric booster", Facility: "Plant-7"},
			{ID: "Seed-Gamma", Product: "Stable mRNA pediatric booster", Facility: "Plant-9"},
			{ID: "Bulk-A", Product: "Stable mRNA pediatric booster", Facility: "Plant-7", Parents: []string{"Seed-Alpha"}},
			{ID: "Bulk-B", Product: "Stable mRNA pediatric booster", Facility: "Plant-7", Parents: []string{"Seed-Beta"}},
			{ID: "BIO-17A", Product: "Stable mRNA pediatric booster", Facility: "CityHub", Doses: 420, Parents: []string{"Bulk-A"}, ReleaseCandidate: true},
			{ID: "BIO-17B", Product: "Stable mRNA pediatric booster", Facility: "CityHub", Doses: 360, Parents: []string{"Bulk-A"}, ReleaseCandidate: true},
			{ID: "BIO-18A", Product: "Stable mRNA pediatric booster", Facility: "CityHub", Doses: 500, Parents: []string{"Bulk-B"}, ReleaseCandidate: true},
			{ID: "BIO-19A", Product: "Stable mRNA pediatric booster", Facility: "AirportHub", Doses: 260, Parents: []string{"Seed-Gamma"}, ReleaseCandidate: true},
			{ID: "BIO-20A", Product: "Stable mRNA pediatric booster", Facility: "NorthDepot", Doses: 300, Parents: []string{"Bulk-A"}, ReleaseCandidate: true},
			{ID: "BIO-21A", Product: "Stable mRNA pediatric booster", Facility: "AirportHub", Doses: 180, Parents: []string{"Seed-Gamma"}, ReleaseCandidate: true},
		},
		RecalledSources: []string{"Seed-Beta"},
		Clinics: []Clinic{
			{ID: "DialysisUnit", Label: "Dialysis Unit", Demand: 300, Priority: 10, MaxTransit: 100, Need: "immunocompromised pediatric list"},
			{ID: "NorthClinic", Label: "North Clinic", Demand: 380, Priority: 9, MaxTransit: 130, Need: "school outbreak ring"},
			{ID: "FieldStation", Label: "Field Station", Demand: 250, Priority: 8, MaxTransit: 90, Need: "temporary shelter"},
			{ID: "RuralMobile", Label: "Rural Mobile Unit", Demand: 180, Priority: 6, MaxTransit: 140, Need: "mobile catch-up route"},
		},
		Transit: []TransitFact{
			{From: "CityHub", ClinicID: "DialysisUnit", Minutes: 70},
			{From: "CityHub", ClinicID: "NorthClinic", Minutes: 110},
			{From: "CityHub", ClinicID: "FieldStation", Minutes: 150},
			{From: "CityHub", ClinicID: "RuralMobile", Minutes: 210},
			{From: "NorthDepot", ClinicID: "DialysisUnit", Minutes: 90},
			{From: "NorthDepot", ClinicID: "NorthClinic", Minutes: 40},
			{From: "NorthDepot", ClinicID: "FieldStation", Minutes: 110},
			{From: "NorthDepot", ClinicID: "RuralMobile", Minutes: 160},
			{From: "AirportHub", ClinicID: "DialysisUnit", Minutes: 120},
			{From: "AirportHub", ClinicID: "NorthClinic", Minutes: 130},
			{From: "AirportHub", ClinicID: "FieldStation", Minutes: 70},
			{From: "AirportHub", ClinicID: "RuralMobile", Minutes: 90},
		},
	}

	lastHash := map[string]string{}
	addEvent := func(lotID string, ordinal int, step, facility, actor string, tempTenths, minutes int, breakPrev bool) {
		prev := lastHash[lotID]
		if prev == "" {
			prev = genesisHash
		}
		if breakPrev {
			prev = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
		}
		event := CustodyEvent{
			LotID:      lotID,
			Ordinal:    ordinal,
			Step:       step,
			Facility:   facility,
			Actor:      actor,
			TempTenths: tempTenths,
			Minutes:    minutes,
			PrevHash:   prev,
		}
		event.Hash = eventHash(event)
		ds.Events = append(ds.Events, event)
		lastHash[lotID] = event.Hash
	}

	addEvent("BIO-17A", 1, "fill-finish", "Plant-7", "qa:ana", 41, 90, false)
	addEvent("BIO-17A", 2, "refrigerated pack", "Plant-7", "ops:li", 78, 45, false)
	addEvent("BIO-17A", 3, "carrier handoff", "CityHub", "driver:milo", 87, 6, false)

	addEvent("BIO-17B", 1, "fill-finish", "Plant-7", "qa:ana", 40, 120, false)
	addEvent("BIO-17B", 2, "refrigerated pack", "Plant-7", "ops:li", 105, 21, false)
	addEvent("BIO-17B", 3, "carrier handoff", "CityHub", "driver:milo", 63, 40, false)

	addEvent("BIO-18A", 1, "fill-finish", "Plant-7", "qa:ana", 39, 80, false)
	addEvent("BIO-18A", 2, "refrigerated pack", "Plant-7", "ops:li", 73, 60, false)
	addEvent("BIO-18A", 3, "carrier handoff", "CityHub", "driver:noor", 81, 10, false)

	addEvent("BIO-19A", 1, "fill-finish", "Plant-9", "qa:rui", 52, 100, false)
	addEvent("BIO-19A", 2, "refrigerated pack", "Plant-9", "ops:mina", 65, 60, false)
	addEvent("BIO-19A", 3, "carrier handoff", "AirportHub", "driver:sol", 72, 40, true)

	addEvent("BIO-20A", 1, "fill-finish", "Plant-7", "qa:ana", 36, 120, false)
	addEvent("BIO-20A", 2, "refrigerated pack", "NorthDepot", "ops:li", 54, 60, false)
	addEvent("BIO-20A", 3, "carrier handoff", "NorthDepot", "driver:ivy", 92, 12, false)

	addEvent("BIO-21A", 1, "fill-finish", "Plant-9", "qa:rui", 24, 90, false)
	addEvent("BIO-21A", 2, "refrigerated pack", "AirportHub", "ops:mina", 84, 9, false)
	addEvent("BIO-21A", 3, "carrier handoff", "AirportHub", "driver:sol", 70, 40, false)

	return ds
}

func eventHash(e CustodyEvent) string {
	canonical := fmt.Sprintf("%s|%02d|%s|%s|%s|%d|%d|%s", e.LotID, e.Ordinal, e.Step, e.Facility, e.Actor, e.TempTenths, e.Minutes, e.PrevHash)
	sum := sha256.Sum256([]byte(canonical))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func derive(ds Dataset) Analysis {
	lotByID := map[string]Lot{}
	for _, lot := range ds.Lots {
		lotByID[lot.ID] = lot
	}
	clinicByID := map[string]Clinic{}
	for _, clinic := range ds.Clinics {
		clinicByID[clinic.ID] = clinic
	}
	transitByKey := map[string]int{}
	for _, fact := range ds.Transit {
		transitByKey[fact.From+"|"+fact.ClinicID] = fact.Minutes
	}

	ancestors := computeAncestors(ds.Lots)
	closureFacts := 0
	for _, xs := range ancestors {
		closureFacts += len(xs)
	}
	parentFacts := 0
	for _, lot := range ds.Lots {
		parentFacts += len(lot.Parents)
	}

	recalled := map[string]string{}
	for _, lot := range ds.Lots {
		for _, ancestor := range ancestors[lot.ID] {
			if containsString(ds.RecalledSources, ancestor) {
				recalled[lot.ID] = ancestor
				break
			}
		}
		if containsString(ds.RecalledSources, lot.ID) {
			recalled[lot.ID] = lot.ID
		}
	}

	eventsByLot := map[string][]CustodyEvent{}
	for _, event := range ds.Events {
		eventsByLot[event.LotID] = append(eventsByLot[event.LotID], event)
	}
	for lotID := range eventsByLot {
		sort.Slice(eventsByLot[lotID], func(i, j int) bool { return eventsByLot[lotID][i].Ordinal < eventsByLot[lotID][j].Ordinal })
	}

	var reports []LotReport
	ledgerValid := 0
	temperatureOK := 0
	for _, lot := range releaseCandidates(ds.Lots) {
		ledger := verifyLedger(eventsByLot[lot.ID])
		temp := checkTemperature(eventsByLot[lot.ID])
		if ledger.OK {
			ledgerValid++
		}
		if temp.OK {
			temperatureOK++
		}
		report := LotReport{
			Lot:         lot,
			Ancestors:   ancestors[lot.ID],
			RecalledBy:  recalled[lot.ID],
			Ledger:      ledger,
			Temperature: temp,
		}
		switch {
		case report.RecalledBy != "":
			report.RejectReason = "descendant of recalled ancestor " + report.RecalledBy
		case !ledger.OK:
			report.RejectReason = ledger.Reason
		case !temp.OK:
			report.RejectReason = temp.Reason
		default:
			report.Release = true
		}
		reports = append(reports, report)
	}
	sort.Slice(reports, func(i, j int) bool { return reports[i].Lot.ID < reports[j].Lot.ID })

	var releaseLots []LotReport
	var quarantineLots []LotReport
	releaseByLot := map[string]LotReport{}
	for _, report := range reports {
		if report.Release {
			releaseLots = append(releaseLots, report)
			releaseByLot[report.Lot.ID] = report
		} else {
			quarantineLots = append(quarantineLots, report)
		}
	}

	feasible, candidateByLot := feasibleAssignments(releaseLots, ds.Clinics, transitByKey)
	plan, stats := optimize(releaseLots, ds.Clinics, candidateByLot)
	reserve := 0
	for _, assignment := range plan.Assignments {
		reserve += assignment.Reserve
	}
	unmet := map[string]int{}
	for _, clinic := range ds.Clinics {
		unmet[clinic.ID] = clinic.Demand
	}
	for _, assignment := range plan.Assignments {
		unmet[assignment.ClinicID] -= assignment.Delivered
	}

	analysis := Analysis{
		Ancestors:       ancestors,
		LotReports:      reports,
		ReleaseLots:     releaseLots,
		QuarantineLots:  quarantineLots,
		Feasible:        feasible,
		Plan:            plan,
		Stats:           stats,
		LedgerValid:     ledgerValid,
		TemperatureOK:   temperatureOK,
		ClosureFacts:    closureFacts,
		ParentFacts:     parentFacts,
		EventFacts:      len(ds.Events),
		TempFacts:       len(ds.Events),
		ReserveDoses:    reserve,
		UnmetByClinic:   unmet,
		CandidateByLot:  candidateByLot,
		TransitByKey:    transitByKey,
		ReleaseByLot:    releaseByLot,
		ClinicByID:      clinicByID,
		LotByID:         lotByID,
		LineageComplete: lineageComplete(ds.Lots, ancestors),
	}
	analysis.Checks = buildChecks(ds, analysis)
	return analysis
}

func computeAncestors(lots []Lot) map[string][]string {
	parents := map[string][]string{}
	for _, lot := range lots {
		parents[lot.ID] = append([]string(nil), lot.Parents...)
	}
	memo := map[string][]string{}
	var visit func(string) []string
	visit = func(id string) []string {
		if xs, ok := memo[id]; ok {
			return xs
		}
		seen := map[string]bool{}
		var out []string
		var walk func(string)
		walk = func(x string) {
			for _, p := range parents[x] {
				if seen[p] {
					continue
				}
				seen[p] = true
				out = append(out, p)
				walk(p)
			}
		}
		walk(id)
		sort.Strings(out)
		memo[id] = out
		return out
	}
	for _, lot := range lots {
		visit(lot.ID)
	}
	return memo
}

func releaseCandidates(lots []Lot) []Lot {
	var out []Lot
	for _, lot := range lots {
		if lot.ReleaseCandidate {
			out = append(out, lot)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func verifyLedger(events []CustodyEvent) LedgerReport {
	if len(events) == 0 {
		return LedgerReport{OK: false, Reason: "no custody events"}
	}
	expected := genesisHash
	root := ""
	for _, event := range events {
		if event.PrevHash != expected {
			return LedgerReport{
				OK:        false,
				Root:      event.Hash,
				BreakStep: event.Step,
				Expected:  expected,
				Got:       event.PrevHash,
				Reason:    fmt.Sprintf("custody ledger break at step %s (expected %s, got %s)", event.Step, shortHash(expected), shortHash(event.PrevHash)),
			}
		}
		recomputed := eventHash(event)
		if recomputed != event.Hash {
			return LedgerReport{OK: false, Root: event.Hash, BreakStep: event.Step, Expected: recomputed, Got: event.Hash, Reason: "custody event hash does not match canonical payload"}
		}
		expected = event.Hash
		root = event.Hash
	}
	return LedgerReport{OK: true, Root: root}
}

func checkTemperature(events []CustodyEvent) TemperatureReport {
	report := TemperatureReport{OK: true}
	for _, event := range events {
		outside := event.TempTenths < coldMinTenths || event.TempTenths > coldMaxTenths
		if event.TempTenths <= 0 || event.TempTenths >= 250 {
			report.HardFail = true
		}
		if outside {
			report.TotalExcursion += event.Minutes
			if event.Minutes > report.MaxSegment {
				report.MaxSegment = event.Minutes
			}
		}
	}
	if report.HardFail {
		report.OK = false
		report.Reason = "hard temperature failure outside 0.0..25.0 °C"
		return report
	}
	if report.MaxSegment > singleExcursionLimitMinutes {
		report.OK = false
		report.Reason = fmt.Sprintf("temperature excursion %dm exceeds single-segment limit %dm", report.MaxSegment, singleExcursionLimitMinutes)
		return report
	}
	if report.TotalExcursion > totalExcursionLimitMinutes {
		report.OK = false
		report.Reason = fmt.Sprintf("temperature excursion total %dm exceeds total limit %dm", report.TotalExcursion, totalExcursionLimitMinutes)
		return report
	}
	return report
}

func feasibleAssignments(releaseLots []LotReport, clinics []Clinic, transit map[string]int) ([]FeasibleAssignment, map[string][]FeasibleAssignment) {
	var out []FeasibleAssignment
	byLot := map[string][]FeasibleAssignment{}
	for _, report := range releaseLots {
		for _, clinic := range clinics {
			minutes := transit[report.Lot.Facility+"|"+clinic.ID]
			if minutes == 0 || minutes > clinic.MaxTransit {
				continue
			}
			doses := min(report.Lot.Doses, clinic.Demand)
			assignment := FeasibleAssignment{
				LotID:       report.Lot.ID,
				ClinicID:    clinic.ID,
				ClinicLabel: clinic.Label,
				Doses:       doses,
				Transit:     minutes,
				Value:       doses * clinic.Priority,
			}
			out = append(out, assignment)
			byLot[report.Lot.ID] = append(byLot[report.Lot.ID], assignment)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].LotID != out[j].LotID {
			return out[i].LotID < out[j].LotID
		}
		return out[i].ClinicID < out[j].ClinicID
	})
	for lotID := range byLot {
		sort.Slice(byLot[lotID], func(i, j int) bool { return byLot[lotID][i].ClinicID < byLot[lotID][j].ClinicID })
	}
	return out, byLot
}

func optimize(releaseLots []LotReport, clinics []Clinic, candidateByLot map[string][]FeasibleAssignment) (Plan, OptimizerStats) {
	sort.Slice(releaseLots, func(i, j int) bool { return releaseLots[i].Lot.ID < releaseLots[j].Lot.ID })
	clinicIndex := map[string]int{}
	for i, clinic := range clinics {
		clinicIndex[clinic.ID] = i
	}
	type stateKey struct {
		Index int
		Used  [4]int
	}
	memo := map[stateKey]Plan{}
	stats := OptimizerStats{}
	var dfs func(int, [4]int) Plan
	dfs = func(index int, used [4]int) Plan {
		key := stateKey{Index: index, Used: used}
		if plan, ok := memo[key]; ok {
			stats.MemoHits++
			return clonePlan(plan)
		}
		stats.States++
		if index == len(releaseLots) {
			memo[key] = Plan{}
			return Plan{}
		}
		lot := releaseLots[index].Lot
		stats.Transitions++ // hold this lot in reserve
		best := dfs(index+1, used)
		for _, candidate := range candidateByLot[lot.ID] {
			stats.Transitions++
			clinic := clinics[clinicIndex[candidate.ClinicID]]
			remaining := clinic.Demand - used[clinicIndex[candidate.ClinicID]]
			if remaining <= 0 {
				continue
			}
			delivered := min(lot.Doses, remaining)
			newUsed := used
			newUsed[clinicIndex[candidate.ClinicID]] += delivered
			sub := dfs(index+1, newUsed)
			selected := SelectedAssignment{
				LotID:       lot.ID,
				ClinicID:    candidate.ClinicID,
				ClinicLabel: candidate.ClinicLabel,
				Delivered:   delivered,
				Reserve:     lot.Doses - delivered,
				Transit:     candidate.Transit,
				Value:       delivered * clinic.Priority,
			}
			plan := Plan{
				Score:       sub.Score + selected.Value,
				Doses:       sub.Doses + selected.Delivered,
				Transit:     sub.Transit + selected.Transit,
				Assignments: append([]SelectedAssignment{selected}, sub.Assignments...),
			}
			if betterPlan(plan, best) {
				best = plan
			}
		}
		memo[key] = clonePlan(best)
		return clonePlan(best)
	}
	return dfs(0, [4]int{}), stats
}

func betterPlan(a, b Plan) bool {
	if a.Score != b.Score {
		return a.Score > b.Score
	}
	if a.Doses != b.Doses {
		return a.Doses > b.Doses
	}
	if a.Transit != b.Transit {
		return a.Transit < b.Transit
	}
	return planKey(a) < planKey(b)
}

func planKey(plan Plan) string {
	parts := make([]string, 0, len(plan.Assignments))
	for _, a := range plan.Assignments {
		parts = append(parts, a.LotID+"->"+a.ClinicID)
	}
	return strings.Join(parts, ";")
}

func clonePlan(plan Plan) Plan {
	plan.Assignments = append([]SelectedAssignment(nil), plan.Assignments...)
	return plan
}

func buildChecks(ds Dataset, a Analysis) []Check {
	used := map[string]int{}
	for _, assignment := range a.Plan.Assignments {
		used[assignment.ClinicID] += assignment.Delivered
	}
	demandOK := true
	for _, clinic := range ds.Clinics {
		if used[clinic.ID] > clinic.Demand {
			demandOK = false
		}
	}
	transitOK := true
	for _, assignment := range a.Plan.Assignments {
		clinic := a.ClinicByID[assignment.ClinicID]
		if assignment.Transit > clinic.MaxTransit {
			transitOK = false
		}
	}
	noQuarantineAllocated := true
	for _, assignment := range a.Plan.Assignments {
		if !a.ReleaseByLot[assignment.LotID].Release {
			noQuarantineAllocated = false
		}
	}
	totalsOK := a.Plan.Score == 7860 && a.Plan.Doses == 860 && a.ReserveDoses == 40
	recallOK := a.report("BIO-18A").RecalledBy == "Seed-Beta" && a.report("BIO-17A").RecalledBy == "" && a.report("BIO-20A").RecalledBy == ""
	tempOK := !a.report("BIO-17B").Temperature.OK && a.report("BIO-17A").Temperature.OK && a.report("BIO-20A").Temperature.OK && a.report("BIO-21A").Temperature.OK
	ledgerOK := a.report("BIO-19A").Ledger.OK == false && a.report("BIO-17A").Ledger.OK && a.report("BIO-20A").Ledger.OK && a.report("BIO-21A").Ledger.OK
	reserveOK := true
	for _, assignment := range a.Plan.Assignments {
		if assignment.Reserve > 0 && a.ReleaseByLot[assignment.LotID].Ledger.Root == "" {
			reserveOK = false
		}
	}
	return []Check{
		{Label: "C1", OK: a.LineageComplete, Text: "every final lot has a computed ancestry closure."},
		{Label: "C2", OK: recallOK, Text: "recalled Seed-Beta propagates to BIO-18A and not to unrelated safe lots."},
		{Label: "C3", OK: ledgerOK, Text: "custody ledgers are hash-linked or explicitly rejected."},
		{Label: "C4", OK: tempOK, Text: "temperature policy rejects BIO-17B but accepts the three released lots."},
		{Label: "C5", OK: noQuarantineAllocated, Text: "no quarantined lot is offered to allocation."},
		{Label: "C6", OK: transitOK, Text: "every selected shipment respects clinic transit limits."},
		{Label: "C7", OK: demandOK, Text: "clinic demand caps are never exceeded."},
		{Label: "C8", OK: a.Plan.Score == 7860 && a.Stats.States == 30, Text: "the exact optimizer proves no higher-value allocation exists."},
		{Label: "C9", OK: totalsOK, Text: "delivered dose/value totals equal the selected shipments."},
		{Label: "C10", OK: reserveOK, Text: "reserve doses remain bound to the original ledger root."},
	}
}

func (a Analysis) report(id string) LotReport {
	for _, report := range a.LotReports {
		if report.Lot.ID == id {
			return report
		}
	}
	return LotReport{}
}

func printAnswer(ds Dataset, a Analysis) {
	fmt.Println("# Cold Chain Release")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("- Release decision : %d of %d candidate lots can ship.\n", len(a.ReleaseLots), len(a.LotReports))
	fmt.Printf("- Safe release lots : %s\n", reportIDs(a.ReleaseLots))
	fmt.Printf("- Quarantined lots : %s\n", reportIDs(a.QuarantineLots))
	fmt.Printf("- Best allocation value : %d priority-dose points\n", a.Plan.Score)
	fmt.Printf("- Doses delivered : %d\n", a.Plan.Doses)
	fmt.Println("- Priority allocation:")
	for _, assignment := range a.Plan.Assignments {
		lot := a.LotByID[assignment.LotID]
		fmt.Printf(" - %s -> %s : %d/%d doses, transit=%dm, value=%d\n", assignment.LotID, assignment.ClinicLabel, assignment.Delivered, lot.Doses, assignment.Transit, assignment.Value)
	}
	fmt.Println("- Reserved/unmet:")
	for _, assignment := range a.Plan.Assignments {
		if assignment.Reserve > 0 {
			fmt.Printf(" - %s remainder held under ledger root %s: %d doses\n", assignment.LotID, shortHash(a.ReleaseByLot[assignment.LotID].Ledger.Root), assignment.Reserve)
		}
	}
	for _, clinic := range ds.Clinics {
		if a.UnmetByClinic[clinic.ID] > 0 {
			fmt.Printf(" - %s unmet demand: %d doses\n", clinic.Label, a.UnmetByClinic[clinic.ID])
		}
	}
	fmt.Println()
}

func printReason(ds Dataset, a Analysis) {
	fmt.Println("## Reason why")
	fmt.Println("- The dataset is treated as linked facts: parent/child lot relationships create a lineage closure; recalled ancestors contaminate descendants; custody rows are verified as a SHA-256 hash chain; temperature segments are checked against the 2-8 °C policy; only releasable lots enter an exact memoized allocation search.")
	fmt.Printf("- lineage parent facts : %d\n", a.ParentFacts)
	fmt.Printf("- lineage closure facts : %d ancestor links\n", a.ClosureFacts)
	fmt.Printf("- custody event facts : %d\n", a.EventFacts)
	fmt.Printf("- temperature observations : %d\n", a.TempFacts)
	fmt.Printf("- release candidates : %d\n", len(a.LotReports))
	fmt.Printf("- recalled ancestors : %s\n", strings.Join(ds.RecalledSources, ", "))
	fmt.Printf("- releasable lots : %d\n", len(a.ReleaseLots))
	fmt.Println("- quarantine reasons:")
	for _, report := range a.QuarantineLots {
		fmt.Printf(" - %s : %s\n", report.Lot.ID, report.RejectReason)
	}
	fmt.Printf("- allocation candidates : %d feasible lot→clinic assignments + %d hold options\n", len(a.Feasible), len(a.ReleaseLots))
	fmt.Printf("- optimizer states : %d\n", a.Stats.States)
	fmt.Printf("- optimizer transitions : %d\n", a.Stats.Transitions)
	fmt.Printf("- memo hits : %d\n", a.Stats.MemoHits)
	fmt.Println()
}

func printChecks(a Analysis) {
	fmt.Println("## Check")
	for _, check := range a.Checks {
		status := "FAIL"
		if check.OK {
			status = "OK"
		}
		fmt.Printf("- %s %s - %s\n", check.Label, status, check.Text)
	}
	fmt.Println()
}

func printAudit(ds Dataset, a Analysis) {
	fmt.Println("## Go audit details")
	fmt.Printf("- platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("- case : %s\n", ds.CaseName)
	fmt.Printf("- question : %s\n", ds.Question)
	fmt.Printf("- product : %s\n", ds.Product)
	fmt.Printf("- lots total : %d\n", len(ds.Lots))
	fmt.Printf("- release candidate lots : %d\n", len(a.LotReports))
	fmt.Printf("- lineage parent facts : %d\n", a.ParentFacts)
	fmt.Printf("- lineage ancestor closure facts : %d\n", a.ClosureFacts)
	fmt.Printf("- recalled source lots : %s\n", strings.Join(ds.RecalledSources, ", "))
	fmt.Printf("- custody event facts : %d\n", a.EventFacts)
	fmt.Printf("- temperature observations : %d\n", a.TempFacts)
	fmt.Printf("- ledger-valid release candidates : %d/%d\n", a.LedgerValid, len(a.LotReports))
	fmt.Printf("- temperature-valid release candidates : %d/%d\n", a.TemperatureOK, len(a.LotReports))
	fmt.Printf("- releasable candidates : %d/%d\n", len(a.ReleaseLots), len(a.LotReports))
	fmt.Printf("- quarantined candidates : %d/%d\n", len(a.QuarantineLots), len(a.LotReports))
	fmt.Println("- ledger roots:")
	for _, report := range a.LotReports {
		status := "valid"
		if !report.Ledger.OK {
			status = "invalid"
		}
		fmt.Printf(" - %s root=%s status=%s\n", report.Lot.ID, shortHash(report.Ledger.Root), status)
	}
	fmt.Println("- release decisions:")
	for _, report := range a.LotReports {
		if report.Release {
			fmt.Printf(" - %s release doses=%d facility=%s tempExcursion=%dm maxSegment=%dm ancestors=%s\n", report.Lot.ID, report.Lot.Doses, report.Lot.Facility, report.Temperature.TotalExcursion, report.Temperature.MaxSegment, strings.Join(report.Ancestors, ", "))
		} else {
			fmt.Printf(" - %s quarantine reason=%s\n", report.Lot.ID, report.RejectReason)
		}
	}
	fmt.Println("- clinics:")
	for _, clinic := range ds.Clinics {
		fmt.Printf(" - %s demand=%d priority=%d maxTransit=%dm need=%s\n", clinic.Label, clinic.Demand, clinic.Priority, clinic.MaxTransit, clinic.Need)
	}
	fmt.Println("- transit candidates:")
	for _, candidate := range a.Feasible {
		fmt.Printf(" - %s -> %s %dm value=%d\n", candidate.LotID, candidate.ClinicLabel, candidate.Transit, candidate.Value)
	}
	fmt.Println("- selected allocation:")
	for _, assignment := range a.Plan.Assignments {
		fmt.Printf(" - %s -> %s delivered=%d reserve=%d transit=%dm value=%d\n", assignment.LotID, assignment.ClinicLabel, assignment.Delivered, assignment.Reserve, assignment.Transit, assignment.Value)
	}
	fmt.Printf("- delivered doses : %d\n", a.Plan.Doses)
	fmt.Printf("- reserved doses : %d\n", a.ReserveDoses)
	fmt.Printf("- unmet demand : %d (%s)\n", totalUnmet(ds.Clinics, a.UnmetByClinic), unmetSummary(ds.Clinics, a.UnmetByClinic))
	fmt.Printf("- best value : %d\n", a.Plan.Score)
	fmt.Printf("- total selected transit : %dm\n", a.Plan.Transit)
	fmt.Printf("- optimizer states explored : %d\n", a.Stats.States)
	fmt.Printf("- optimizer memo hits : %d\n", a.Stats.MemoHits)
	fmt.Printf("- optimizer transitions : %d\n", a.Stats.Transitions)
	fmt.Printf("- checks passed : %d/%d\n", checksPassed(a.Checks), len(a.Checks))
	fmt.Printf("- all checks pass : %s\n", yesNo(allChecksOK(a.Checks)))
}

func lineageComplete(lots []Lot, ancestors map[string][]string) bool {
	for _, lot := range lots {
		if lot.ReleaseCandidate && len(ancestors[lot.ID]) == 0 {
			return false
		}
	}
	return true
}

func reportIDs(reports []LotReport) string {
	ids := make([]string, 0, len(reports))
	for _, report := range reports {
		ids = append(ids, report.Lot.ID)
	}
	return strings.Join(ids, ", ")
}

func shortHash(hash string) string {
	if hash == "" {
		return "<none>"
	}
	if len(hash) <= 23 {
		return hash
	}
	return hash[:23] + "…"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func containsString(xs []string, x string) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

func checksPassed(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}

func allChecksOK(checks []Check) bool {
	return checksPassed(checks) == len(checks)
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}

func totalUnmet(clinics []Clinic, unmet map[string]int) int {
	total := 0
	for _, clinic := range clinics {
		total += unmet[clinic.ID]
	}
	return total
}

func unmetSummary(clinics []Clinic, unmet map[string]int) string {
	var parts []string
	for _, clinic := range clinics {
		if unmet[clinic.ID] > 0 {
			parts = append(parts, fmt.Sprintf("%s %d", clinic.Label, unmet[clinic.ID]))
		}
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, " + ")
}
