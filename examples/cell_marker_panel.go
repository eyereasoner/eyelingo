// cell_marker_panel.go
//
// A self-contained Go translation of a small Eyeling-style single-cell biology
// example.
//
// The scenario models a tiny reference atlas with eight cell populations and a
// candidate pool of gene markers. The goal is to design the smallest qPCR / flow
// cytometry style marker panel that gives every population a positive anchor and
// separates every pair of populations by at least one robust expression gap.
//
// This is intentionally not a generic RDF/N3 reasoner. The biological facts and
// rules are represented as Go structs and ordinary functions so the inference
// mechanics stay visible and directly runnable.
//
// Run:
//
//	go run cell_marker_panel.go
//
// The program has no third-party dependencies.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"math/bits"
	"os"
	"runtime"
	"sort"
	"strings"
)

const eyelingoExampleName = "cell_marker_panel"

const (
	separationThreshold = 5
	anchorHigh          = 7
	anchorMaxOffTarget  = 3
)

type Dataset struct {
	CaseName  string
	Question  string
	CellTypes []CellType
	Genes     []Gene
}

type CellType struct {
	ID      string
	Label   string
	Lineage string
}

type Gene struct {
	ID             string
	Label          string
	Expression     []int
	AssayCost      int
	Stable         bool
	Assayable      bool
	ExcludedReason string
}

type Pair struct {
	Index int
	A     int
	B     int
}

type Marker struct {
	GeneIndex       int
	PairMask        uint64
	AnchorMask      uint16
	Margins         []int
	PairMarginSum   int
	AnchorMarginSum int
}

type DominatedMarker struct {
	GeneIndex      int
	DominatedBy    int
	PairCoverage   int
	AnchorCoverage int
}

type Panel struct {
	MarkerIndexes   []int
	PairMask        uint64
	AnchorMask      uint16
	PairMarginSum   int
	AnchorMarginSum int
	Cost            int
}

type SearchStats struct {
	StatesVisited       int
	LeafStates          int
	IncludeBranches     int
	ExcludeBranches     int
	ImpossiblePrunes    int
	SizePrunes          int
	LowerBoundPrunes    int
	SolutionsFound      int
	MinimalProofSubsets int
	RetainedMarkers     int
	DominatedPruned     int
}

type Checks struct {
	AllPairsSeparated        bool
	AllCellTypesAnchored     bool
	NoExcludedGenesSelected  bool
	NoDominatedGenesSelected bool
	AssayRulesSatisfied      bool
	SignaturesUnique         bool
	ThresholdsRespected      bool
	ExactMinimumProved       bool
	ScoreConsistent          bool
}

func fixture() Dataset {
	return Dataset{
		CaseName: "pbmc-plus-epithelium-marker-panel",
		Question: "Which compact marker panel separates the reference cell populations?",
		CellTypes: []CellType{
			{ID: ":T", Label: "T cell", Lineage: "immune"},
			{ID: ":B", Label: "B cell", Lineage: "immune"},
			{ID: ":NK", Label: "NK cell", Lineage: "immune"},
			{ID: ":Mono", Label: "Monocyte", Lineage: "immune"},
			{ID: ":DC", Label: "Dendritic cell", Lineage: "immune"},
			{ID: ":Epi", Label: "Epithelial cell", Lineage: "barrier"},
			{ID: ":Endo", Label: "Endothelial cell", Lineage: "vascular"},
			{ID: ":Fib", Label: "Fibroblast", Lineage: "stromal"},
		},
		Genes: []Gene{
			{ID: "PTPRC", Label: "pan-leukocyte gate", Expression: []int{8, 8, 8, 8, 8, 1, 1, 1}, AssayCost: 10, Stable: true, Assayable: true},
			{ID: "CD3D", Label: "T-cell receptor complex", Expression: []int{9, 1, 2, 1, 1, 0, 0, 0}, AssayCost: 10, Stable: true, Assayable: true},
			{ID: "MS4A1", Label: "B-cell membrane marker", Expression: []int{1, 9, 1, 1, 1, 0, 0, 0}, AssayCost: 10, Stable: true, Assayable: true},
			{ID: "NKG7", Label: "cytotoxic granule marker", Expression: []int{3, 1, 9, 2, 2, 0, 0, 0}, AssayCost: 11, Stable: true, Assayable: true},
			{ID: "CD14", Label: "monocyte co-receptor", Expression: []int{0, 0, 2, 8, 2, 0, 0, 0}, AssayCost: 11, Stable: true, Assayable: true},
			{ID: "FCER1A", Label: "dendritic-cell receptor", Expression: []int{0, 1, 0, 2, 9, 0, 0, 0}, AssayCost: 12, Stable: true, Assayable: true},
			{ID: "EPCAM", Label: "epithelial adhesion", Expression: []int{0, 0, 0, 0, 0, 9, 1, 1}, AssayCost: 10, Stable: true, Assayable: true},
			{ID: "PECAM1", Label: "endothelial junction", Expression: []int{0, 0, 0, 0, 0, 1, 9, 1}, AssayCost: 10, Stable: true, Assayable: true},
			{ID: "COL1A1", Label: "fibroblast matrix", Expression: []int{0, 0, 0, 0, 0, 1, 2, 9}, AssayCost: 12, Stable: true, Assayable: true},
			{ID: "LYZ", Label: "shared myeloid lysozyme", Expression: []int{1, 1, 5, 9, 7, 0, 0, 0}, AssayCost: 11, Stable: true, Assayable: true},
			{ID: "KRT18", Label: "epithelial keratin backup", Expression: []int{0, 0, 0, 0, 0, 8, 1, 1}, AssayCost: 11, Stable: true, Assayable: true},
			{ID: "VWF", Label: "endothelial secretion backup", Expression: []int{0, 0, 0, 0, 0, 1, 8, 1}, AssayCost: 11, Stable: true, Assayable: true},
			{ID: "ACTB", Label: "housekeeping control", Expression: []int{6, 6, 6, 6, 6, 6, 6, 6}, AssayCost: 6, Stable: true, Assayable: true, ExcludedReason: "housekeeping gene is not lineage-specific"},
			{ID: "RPLP0", Label: "ribosomal control", Expression: []int{5, 5, 5, 5, 5, 5, 5, 5}, AssayCost: 6, Stable: true, Assayable: true, ExcludedReason: "ribosomal gene fails marker policy"},
			{ID: "MT-CO1", Label: "mitochondrial stress signal", Expression: []int{4, 4, 5, 6, 6, 7, 5, 4}, AssayCost: 8, Stable: false, Assayable: true, ExcludedReason: "mitochondrial stress signal is batch-sensitive"},
			{ID: "MALAT1", Label: "ambient RNA sentinel", Expression: []int{5, 5, 5, 4, 5, 5, 5, 5}, AssayCost: 7, Stable: false, Assayable: true, ExcludedReason: "ambient RNA risk in dissociation protocol"},
		},
	}
}

func buildPairs(cellCount int) []Pair {
	pairs := []Pair{}
	for i := 0; i < cellCount; i++ {
		for j := i + 1; j < cellCount; j++ {
			pairs = append(pairs, Pair{Index: len(pairs), A: i, B: j})
		}
	}
	return pairs
}

func buildMarkers(data Dataset, pairs []Pair) ([]Marker, int) {
	markers := []Marker{}
	excluded := 0
	for gi, gene := range data.Genes {
		if gene.ExcludedReason != "" || !gene.Stable || !gene.Assayable {
			excluded++
			continue
		}
		m := Marker{
			GeneIndex: gi,
			Margins:   make([]int, len(pairs)),
		}
		for _, pair := range pairs {
			margin := abs(gene.Expression[pair.A] - gene.Expression[pair.B])
			m.Margins[pair.Index] = margin
			if margin >= separationThreshold {
				m.PairMask |= 1 << pair.Index
				m.PairMarginSum += margin
			}
		}
		for ci := range data.CellTypes {
			if isAnchorFor(gene, ci) {
				m.AnchorMask |= 1 << ci
				m.AnchorMarginSum += anchorMargin(gene, ci)
			}
		}
		markers = append(markers, m)
	}
	return markers, excluded
}

func isAnchorFor(gene Gene, cellIndex int) bool {
	if gene.Expression[cellIndex] < anchorHigh {
		return false
	}
	for i, expr := range gene.Expression {
		if i == cellIndex {
			continue
		}
		if expr > anchorMaxOffTarget {
			return false
		}
	}
	return true
}

func anchorMargin(gene Gene, cellIndex int) int {
	maxOther := 0
	for i, expr := range gene.Expression {
		if i != cellIndex && expr > maxOther {
			maxOther = expr
		}
	}
	return gene.Expression[cellIndex] - maxOther
}

func pruneDominated(data Dataset, markers []Marker) ([]Marker, []DominatedMarker) {
	dominated := map[int]DominatedMarker{}
	for bi, b := range markers {
		for ai, a := range markers {
			if ai == bi {
				continue
			}
			if dominates(data, a, b) {
				dominated[bi] = DominatedMarker{
					GeneIndex:      b.GeneIndex,
					DominatedBy:    a.GeneIndex,
					PairCoverage:   bits.OnesCount64(b.PairMask),
					AnchorCoverage: bits.OnesCount16(b.AnchorMask),
				}
				break
			}
		}
	}
	kept := []Marker{}
	for i, marker := range markers {
		if _, ok := dominated[i]; !ok {
			kept = append(kept, marker)
		}
	}
	domList := []DominatedMarker{}
	for _, dom := range dominated {
		domList = append(domList, dom)
	}
	sort.Slice(domList, func(i, j int) bool {
		return data.Genes[domList[i].GeneIndex].ID < data.Genes[domList[j].GeneIndex].ID
	})
	return kept, domList
}

func dominates(data Dataset, a, b Marker) bool {
	if data.Genes[a.GeneIndex].AssayCost > data.Genes[b.GeneIndex].AssayCost {
		return false
	}
	if (a.PairMask | b.PairMask) != a.PairMask {
		return false
	}
	if (a.AnchorMask | b.AnchorMask) != a.AnchorMask {
		return false
	}
	for idx, margin := range b.Margins {
		if margin >= separationThreshold && a.Margins[idx] < margin {
			return false
		}
	}
	if b.AnchorMask != 0 {
		for ci := range data.CellTypes {
			mask := uint16(1 << ci)
			if b.AnchorMask&mask != 0 && anchorMargin(data.Genes[a.GeneIndex], ci) < anchorMargin(data.Genes[b.GeneIndex], ci) {
				return false
			}
		}
	}
	return true
}

func solvePanel(data Dataset, markers []Marker, pairCount int, cellCount int) (Panel, SearchStats) {
	sort.SliceStable(markers, func(i, j int) bool {
		ci := bits.OnesCount64(markers[i].PairMask) + 3*bits.OnesCount16(markers[i].AnchorMask)
		cj := bits.OnesCount64(markers[j].PairMask) + 3*bits.OnesCount16(markers[j].AnchorMask)
		if ci != cj {
			return ci > cj
		}
		if markers[i].PairMarginSum != markers[j].PairMarginSum {
			return markers[i].PairMarginSum > markers[j].PairMarginSum
		}
		return data.Genes[markers[i].GeneIndex].ID < data.Genes[markers[j].GeneIndex].ID
	})

	allPairs := uint64(1<<pairCount) - 1
	allAnchors := uint16(1<<cellCount) - 1
	suffixPair := make([]uint64, len(markers)+1)
	suffixAnchor := make([]uint16, len(markers)+1)
	for i := len(markers) - 1; i >= 0; i-- {
		suffixPair[i] = suffixPair[i+1] | markers[i].PairMask
		suffixAnchor[i] = suffixAnchor[i+1] | markers[i].AnchorMask
	}

	best := Panel{MarkerIndexes: nil, Cost: math.MaxInt}
	stats := SearchStats{RetainedMarkers: len(markers)}
	selected := []int{}

	var dfs func(index int, pairMask uint64, anchorMask uint16, pairMargin int, anchorMarginTotal int, cost int)
	dfs = func(index int, pairMask uint64, anchorMask uint16, pairMargin int, anchorMarginTotal int, cost int) {
		stats.StatesVisited++

		if pairMask == allPairs && anchorMask == allAnchors {
			stats.LeafStates++
			candidate := Panel{
				MarkerIndexes:   append([]int(nil), selected...),
				PairMask:        pairMask,
				AnchorMask:      anchorMask,
				PairMarginSum:   pairMargin,
				AnchorMarginSum: anchorMarginTotal,
				Cost:            cost,
			}
			if betterPanel(data, markers, candidate, best) {
				best = candidate
			}
			stats.SolutionsFound++
			return
		}

		if index == len(markers) {
			stats.LeafStates++
			return
		}

		if best.MarkerIndexes != nil && len(selected) >= len(best.MarkerIndexes) {
			stats.SizePrunes++
			return
		}

		if ((pairMask | suffixPair[index]) != allPairs) || ((anchorMask | suffixAnchor[index]) != allAnchors) {
			stats.ImpossiblePrunes++
			return
		}

		if best.MarkerIndexes != nil {
			remainingPairBits := bits.OnesCount64(allPairs &^ pairMask)
			remainingAnchorBits := bits.OnesCount16(allAnchors &^ anchorMask)
			maxNew := 0
			for i := index; i < len(markers); i++ {
				gain := bits.OnesCount64(markers[i].PairMask&^pairMask) + 3*bits.OnesCount16(markers[i].AnchorMask&^anchorMask)
				if gain > maxNew {
					maxNew = gain
				}
			}
			need := remainingPairBits + 3*remainingAnchorBits
			if maxNew == 0 || len(selected)+ceilDiv(need, maxNew) > len(best.MarkerIndexes) {
				stats.LowerBoundPrunes++
				return
			}
		}

		marker := markers[index]
		selected = append(selected, index)
		stats.IncludeBranches++
		dfs(
			index+1,
			pairMask|marker.PairMask,
			anchorMask|marker.AnchorMask,
			pairMargin+marker.PairMarginSum,
			anchorMarginTotal+marker.AnchorMarginSum,
			cost+data.Genes[marker.GeneIndex].AssayCost,
		)
		selected = selected[:len(selected)-1]

		stats.ExcludeBranches++
		dfs(index+1, pairMask, anchorMask, pairMargin, anchorMarginTotal, cost)
	}

	dfs(0, 0, 0, 0, 0, 0)
	return best, stats
}

func betterPanel(data Dataset, markers []Marker, candidate Panel, current Panel) bool {
	if current.MarkerIndexes == nil {
		return true
	}
	if len(candidate.MarkerIndexes) != len(current.MarkerIndexes) {
		return len(candidate.MarkerIndexes) < len(current.MarkerIndexes)
	}
	if candidate.AnchorMarginSum != current.AnchorMarginSum {
		return candidate.AnchorMarginSum > current.AnchorMarginSum
	}
	if candidate.PairMarginSum != current.PairMarginSum {
		return candidate.PairMarginSum > current.PairMarginSum
	}
	if candidate.Cost != current.Cost {
		return candidate.Cost < current.Cost
	}
	return panelGeneList(data, markers, candidate) < panelGeneList(data, markers, current)
}

func proveMinimal(data Dataset, markers []Marker, best Panel, pairCount int, cellCount int) (bool, int) {
	allPairs := uint64(1<<pairCount) - 1
	allAnchors := uint16(1<<cellCount) - 1
	checked := 0
	for size := 0; size < len(best.MarkerIndexes); size++ {
		comb := make([]int, 0, size)
		found := false
		var rec func(start int, pairMask uint64, anchorMask uint16)
		rec = func(start int, pairMask uint64, anchorMask uint16) {
			if found {
				return
			}
			if len(comb) == size {
				checked++
				if pairMask == allPairs && anchorMask == allAnchors {
					found = true
				}
				return
			}
			remainingSlots := size - len(comb)
			for i := start; i <= len(markers)-remainingSlots; i++ {
				comb = append(comb, i)
				rec(i+1, pairMask|markers[i].PairMask, anchorMask|markers[i].AnchorMask)
				comb = comb[:len(comb)-1]
			}
		}
		rec(0, 0, 0)
		if found {
			return false, checked
		}
	}
	return true, checked
}

func buildChecks(data Dataset, markers []Marker, dominated []DominatedMarker, panel Panel, pairs []Pair, minimal bool) Checks {
	selectedGene := map[int]bool{}
	for _, markerIndex := range panel.MarkerIndexes {
		selectedGene[markers[markerIndex].GeneIndex] = true
	}
	dominatedGene := map[int]bool{}
	for _, dom := range dominated {
		dominatedGene[dom.GeneIndex] = true
	}
	allPairs := uint64(1<<len(pairs)) - 1
	allAnchors := uint16(1<<len(data.CellTypes)) - 1
	noExcluded := true
	noDominated := true
	assayOK := true
	for gi := range selectedGene {
		gene := data.Genes[gi]
		if gene.ExcludedReason != "" {
			noExcluded = false
		}
		if dominatedGene[gi] {
			noDominated = false
		}
		if !gene.Stable || !gene.Assayable {
			assayOK = false
		}
	}
	return Checks{
		AllPairsSeparated:        panel.PairMask == allPairs,
		AllCellTypesAnchored:     panel.AnchorMask == allAnchors,
		NoExcludedGenesSelected:  noExcluded,
		NoDominatedGenesSelected: noDominated,
		AssayRulesSatisfied:      assayOK,
		SignaturesUnique:         signaturesUnique(data, markers, panel),
		ThresholdsRespected:      thresholdsRespected(data, markers, panel, pairs),
		ExactMinimumProved:       minimal,
		ScoreConsistent:          panel.PairMarginSum == recomputePairMargin(data, markers, panel) && panel.AnchorMarginSum == recomputeAnchorMargin(data, markers, panel),
	}
}

func signaturesUnique(data Dataset, markers []Marker, panel Panel) bool {
	seen := map[string]bool{}
	for ci := range data.CellTypes {
		parts := []string{}
		for _, markerIndex := range sortedPanelMarkerIndexes(data, markers, panel) {
			gene := data.Genes[markers[markerIndex].GeneIndex]
			parts = append(parts, fmt.Sprintf("%s=%d", gene.ID, gene.Expression[ci]))
		}
		sig := strings.Join(parts, ";")
		if seen[sig] {
			return false
		}
		seen[sig] = true
	}
	return true
}

func thresholdsRespected(data Dataset, markers []Marker, panel Panel, pairs []Pair) bool {
	for _, pair := range pairs {
		ok := false
		for _, markerIndex := range panel.MarkerIndexes {
			gene := data.Genes[markers[markerIndex].GeneIndex]
			if abs(gene.Expression[pair.A]-gene.Expression[pair.B]) >= separationThreshold {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	for ci := range data.CellTypes {
		ok := false
		for _, markerIndex := range panel.MarkerIndexes {
			gene := data.Genes[markers[markerIndex].GeneIndex]
			if isAnchorFor(gene, ci) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func recomputePairMargin(data Dataset, markers []Marker, panel Panel) int {
	total := 0
	for _, markerIndex := range panel.MarkerIndexes {
		marker := markers[markerIndex]
		total += marker.PairMarginSum
	}
	return total
}

func recomputeAnchorMargin(data Dataset, markers []Marker, panel Panel) int {
	total := 0
	for _, markerIndex := range panel.MarkerIndexes {
		marker := markers[markerIndex]
		total += marker.AnchorMarginSum
	}
	return total
}

func renderAnswer(data Dataset, markers []Marker, panel Panel, pairs []Pair) {
	fmt.Println("=== Answer ===")
	fmt.Printf("The exact marker panel uses %d genes and separates all %d cell-type pairs.\n", len(panel.MarkerIndexes), len(pairs))
	fmt.Printf("case : %s\n", data.CaseName)
	fmt.Printf("positive anchors : %d/%d cell populations\n", bits.OnesCount16(panel.AnchorMask), len(data.CellTypes))
	fmt.Printf("assay cost : %d\n\n", panel.Cost)

	fmt.Println("Selected panel:")
	for _, markerIndex := range sortedPanelMarkerIndexes(data, markers, panel) {
		marker := markers[markerIndex]
		gene := data.Genes[marker.GeneIndex]
		anchors := anchorLabels(data, marker.AnchorMask)
		if anchors == "" {
			anchors = "none"
		}
		fmt.Printf(" - %s (%s) anchors=%s pairSeparations=%d\n", gene.ID, gene.Label, anchors, bits.OnesCount64(marker.PairMask))
	}
}

func renderReason(data Dataset, markers []Marker, dominated []DominatedMarker, panel Panel, pairs []Pair, stats SearchStats, minimalSubsets int) {
	fmt.Println()
	fmt.Println("=== Reason Why ===")
	fmt.Println("Candidate genes are first filtered by assay policy: housekeeping, ribosomal, mitochondrial-stress, and ambient-RNA features cannot be used as lineage markers. Each remaining gene then derives two facts: which cell-type pairs it separates by the expression-gap threshold, and which cell type, if any, it positively anchors with low off-target signal. Dominated backup markers are removed before an exact branch-and-bound search solves the paired set-cover goal.")
	fmt.Printf("cell populations : %d\n", len(data.CellTypes))
	fmt.Printf("cell-type pairs : %d\n", len(pairs))
	fmt.Printf("candidate genes : %d\n", len(data.Genes))
	fmt.Printf("usable after QC : %d\n", countUsable(data.Genes))
	fmt.Printf("dominated markers pruned : %d\n", len(dominated))
	fmt.Printf("retained markers searched : %d\n", stats.RetainedMarkers)
	fmt.Printf("search states visited : %d\n", stats.StatesVisited)
	fmt.Printf("branch-and-bound prunes : impossible=%d size=%d lowerBound=%d\n", stats.ImpossiblePrunes, stats.SizePrunes, stats.LowerBoundPrunes)
	fmt.Printf("minimality subsets checked : %d\n", minimalSubsets)
	fmt.Printf("panel pair margin sum : %d\n", panel.PairMarginSum)
	fmt.Printf("panel anchor margin sum : %d\n", panel.AnchorMarginSum)

	fmt.Println("Derived positive anchors:")
	for _, markerIndex := range sortedPanelMarkerIndexes(data, markers, panel) {
		marker := markers[markerIndex]
		gene := data.Genes[marker.GeneIndex]
		if marker.AnchorMask == 0 {
			continue
		}
		for ci := range data.CellTypes {
			if marker.AnchorMask&(1<<ci) != 0 {
				fmt.Printf(" - %s anchors %s: on=%d maxOff=%d margin=%d\n", gene.ID, data.CellTypes[ci].ID, gene.Expression[ci], maxOffTarget(gene, ci), anchorMargin(gene, ci))
			}
		}
	}
}

func renderChecks(checks Checks) {
	fmt.Println()
	fmt.Println("=== Check ===")
	items := []struct {
		OK   bool
		Text string
	}{
		{checks.AllPairsSeparated, "every pair of cell populations is separated by at least one selected gene."},
		{checks.AllCellTypesAnchored, "every cell population has a positive high-confidence anchor."},
		{checks.NoExcludedGenesSelected, "excluded housekeeping, ribosomal, mitochondrial, and ambient markers are not selected."},
		{checks.NoDominatedGenesSelected, "dominated backup markers are not selected."},
		{checks.AssayRulesSatisfied, "selected genes are stable and assayable."},
		{checks.SignaturesUnique, "selected expression signatures are unique for all populations."},
		{checks.ThresholdsRespected, "all reported separations meet the configured expression-gap threshold."},
		{checks.ExactMinimumProved, "no smaller panel satisfies both pair separation and positive-anchor constraints."},
		{checks.ScoreConsistent, "reported margin and cost totals match the selected genes."},
	}
	for i, item := range items {
		status := "FAIL"
		if item.OK {
			status = "OK"
		}
		fmt.Printf("C%d %s - %s\n", i+1, status, item.Text)
	}
}

func renderAudit(data Dataset, markers []Marker, dominated []DominatedMarker, panel Panel, pairs []Pair, stats SearchStats, checks Checks) {
	fmt.Println()
	fmt.Println("=== Go audit details ===")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("question : %s\n", data.Question)
	fmt.Printf("cell populations : %d\n", len(data.CellTypes))
	fmt.Printf("pairwise separation obligations : %d\n", len(pairs))
	fmt.Printf("separation threshold : %d\n", separationThreshold)
	fmt.Printf("anchor rule : on >= %d and every off-target <= %d\n", anchorHigh, anchorMaxOffTarget)
	fmt.Printf("genes in fixture : %d\n", len(data.Genes))
	fmt.Printf("usable genes : %d\n", countUsable(data.Genes))
	fmt.Printf("excluded genes : %d\n", len(data.Genes)-countUsable(data.Genes))
	fmt.Printf("dominated genes pruned : %d\n", len(dominated))
	for _, dom := range dominated {
		fmt.Printf(" - %s dominated by %s pairSeparations=%d anchors=%d\n", data.Genes[dom.GeneIndex].ID, data.Genes[dom.DominatedBy].ID, dom.PairCoverage, dom.AnchorCoverage)
	}
	fmt.Printf("retained search genes : %d\n", stats.RetainedMarkers)
	fmt.Printf("selected genes : %s\n", panelGeneList(data, markers, panel))
	fmt.Printf("selected gene count : %d\n", len(panel.MarkerIndexes))
	fmt.Printf("selected assay cost : %d\n", panel.Cost)
	fmt.Printf("covered pair mask : %0*b\n", len(pairs), panel.PairMask)
	fmt.Printf("covered anchor mask : %0*b\n", len(data.CellTypes), panel.AnchorMask)
	fmt.Printf("pair margin sum : %d\n", panel.PairMarginSum)
	fmt.Printf("anchor margin sum : %d\n", panel.AnchorMarginSum)
	fmt.Println("cell signatures:")
	for ci, cell := range data.CellTypes {
		fmt.Printf(" - %s %-16s %s\n", cell.ID, cell.Label, expressionSignature(data, markers, panel, ci))
	}
	fmt.Println("retained gene coverage:")
	for _, marker := range sortedMarkersByFixtureOrder(data, markers) {
		gene := data.Genes[marker.GeneIndex]
		fmt.Printf(" - %s pairSeparations=%d anchors=%s cost=%d\n", gene.ID, bits.OnesCount64(marker.PairMask), anchorLabels(data, marker.AnchorMask), gene.AssayCost)
	}
	fmt.Println("excluded gene reasons:")
	for _, gene := range data.Genes {
		if gene.ExcludedReason != "" || !gene.Stable || !gene.Assayable {
			reason := gene.ExcludedReason
			if reason == "" {
				reason = "failed stability or assayability rule"
			}
			fmt.Printf(" - %s : %s\n", gene.ID, reason)
		}
	}
	fmt.Printf("search states visited : %d\n", stats.StatesVisited)
	fmt.Printf("leaf states reached : %d\n", stats.LeafStates)
	fmt.Printf("include branches : %d\n", stats.IncludeBranches)
	fmt.Printf("exclude branches : %d\n", stats.ExcludeBranches)
	fmt.Printf("impossible prunes : %d\n", stats.ImpossiblePrunes)
	fmt.Printf("size prunes : %d\n", stats.SizePrunes)
	fmt.Printf("lower-bound prunes : %d\n", stats.LowerBoundPrunes)
	fmt.Printf("complete solutions found : %d\n", stats.SolutionsFound)
	fmt.Printf("checks passed : %d/9\n", countChecks(checks))
	fmt.Printf("all checks pass : %s\n", yesNo(countChecks(checks) == 9))
}

func sortedPanelMarkerIndexes(data Dataset, markers []Marker, panel Panel) []int {
	indexes := append([]int(nil), panel.MarkerIndexes...)
	sort.Slice(indexes, func(i, j int) bool {
		return markers[indexes[i]].GeneIndex < markers[indexes[j]].GeneIndex
	})
	return indexes
}

func sortedMarkersByFixtureOrder(data Dataset, markers []Marker) []Marker {
	copyMarkers := append([]Marker(nil), markers...)
	sort.Slice(copyMarkers, func(i, j int) bool {
		return copyMarkers[i].GeneIndex < copyMarkers[j].GeneIndex
	})
	return copyMarkers
}

func panelGeneList(data Dataset, markers []Marker, panel Panel) string {
	ids := []string{}
	for _, markerIndex := range sortedPanelMarkerIndexes(data, markers, panel) {
		ids = append(ids, data.Genes[markers[markerIndex].GeneIndex].ID)
	}
	return strings.Join(ids, ", ")
}

func anchorLabels(data Dataset, mask uint16) string {
	labels := []string{}
	for ci, cell := range data.CellTypes {
		if mask&(1<<ci) != 0 {
			labels = append(labels, cell.ID)
		}
	}
	return strings.Join(labels, ",")
}

func expressionSignature(data Dataset, markers []Marker, panel Panel, cellIndex int) string {
	parts := []string{}
	for _, markerIndex := range sortedPanelMarkerIndexes(data, markers, panel) {
		gene := data.Genes[markers[markerIndex].GeneIndex]
		parts = append(parts, fmt.Sprintf("%s=%d", gene.ID, gene.Expression[cellIndex]))
	}
	return strings.Join(parts, " ")
}

func countUsable(genes []Gene) int {
	count := 0
	for _, gene := range genes {
		if gene.ExcludedReason == "" && gene.Stable && gene.Assayable {
			count++
		}
	}
	return count
}

func countChecks(checks Checks) int {
	values := []bool{
		checks.AllPairsSeparated,
		checks.AllCellTypesAnchored,
		checks.NoExcludedGenesSelected,
		checks.NoDominatedGenesSelected,
		checks.AssayRulesSatisfied,
		checks.SignaturesUnique,
		checks.ThresholdsRespected,
		checks.ExactMinimumProved,
		checks.ScoreConsistent,
	}
	count := 0
	for _, value := range values {
		if value {
			count++
		}
	}
	return count
}

func maxOffTarget(gene Gene, cellIndex int) int {
	maxOther := 0
	for i, expr := range gene.Expression {
		if i != cellIndex && expr > maxOther {
			maxOther = expr
		}
	}
	return maxOther
}

func ceilDiv(a int, b int) int {
	if b <= 0 {
		return math.MaxInt / 4
	}
	return (a + b - 1) / b
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func main() {
	data := exampleinput.Load(eyelingoExampleName, fixture())
	pairs := buildPairs(len(data.CellTypes))
	markers, _ := buildMarkers(data, pairs)
	markers, dominated := pruneDominated(data, markers)
	panel, stats := solvePanel(data, markers, len(pairs), len(data.CellTypes))
	if panel.MarkerIndexes == nil {
		fmt.Fprintln(os.Stderr, "No feasible marker panel found.")
		os.Exit(1)
	}
	stats.DominatedPruned = len(dominated)
	minimal, minimalSubsets := proveMinimal(data, markers, panel, len(pairs), len(data.CellTypes))
	stats.MinimalProofSubsets = minimalSubsets
	checks := buildChecks(data, markers, dominated, panel, pairs, minimal)

	renderAnswer(data, markers, panel, pairs)
	renderReason(data, markers, dominated, panel, pairs, stats, minimalSubsets)
	renderChecks(checks)
	renderAudit(data, markers, dominated, panel, pairs, stats, checks)
}
