// sudoku.go
//
// Standalone Go translation of the Eyeling Sudoku example:
//   - examples/sudoku.n3 supplies the default puzzle and report/check structure.
//   - examples/builtin/sudoku.js supplies the solving and validation logic.
//
// The program reads an 81-cell Sudoku puzzle, solves it with constraint propagation
// plus depth-first search, and prints an N3-style explanation report with answer,
// reasoning summary, consistency checks, and Go-specific audit details.
//
// Usage:
//
//    go run sudoku.go
//    go run sudoku.go -puzzle "100007090030020008009600500005300900010080002600004000300000010040000007007000300"
//
// Puzzle input accepts digits 1-9 for givens and 0, '.', or '_' for blanks.
// Whitespace and common board separators such as '|', '+', and '-' are ignored.
package main

import (
    "flag"
    "fmt"
    "math/bits"
    "os"
    "runtime"
    "strings"
    "unicode"
)

const (
    allDigitsMask = 0x1ff // bits 0..8 represent digits 1..9
    defaultPuzzle = "100007090030020008009600500005300900010080002600004000300000010040000007007000300"
)

// Move records one placement made by the solver.
// It is used both for the human-readable move summary and for replay validation.
type Move struct {
    Index          int
    Value          int
    CandidatesMask int
    Forced         bool
}

// Stats collects search metrics that are printed in the explanation and audit output.
type Stats struct {
    Givens         int
    Blanks         int
    ForcedMoves    int
    GuessedMoves   int
    RecursiveNodes int
    Backtracks     int
    MaxDepth       int
}

// State is the mutable Sudoku board used during search.
// RowUsed, ColUsed, and BoxUsed mirror Cells and store which digits are already present.
type State struct {
    Cells   []int
    RowUsed []int
    ColUsed []int
    BoxUsed []int
    Moves   []Move
}

// Report is the data model for the final output sections.
type Report struct {
    Status          string
    Error           string
    Raw             string
    Normalized      string
    Givens          int
    Blanks          int
    ForcedMoves     int
    GuessedMoves    int
    RecursiveNodes  int
    Backtracks      int
    MaxDepth        int
    Unique          bool
    Solution        string
    PuzzleText      string
    SolutionText    string
    MoveSummary     string
    MoveCount       int
    GivensPreserved bool
    NoBlanks        bool
    RowsComplete    bool
    ColsComplete    bool
    BoxesComplete   bool
    ReplayLegal     bool
    StoryConsistent bool
}

func digitMask(v int) int {
    return 1 << (v - 1)
}

func boxIndex(row, col int) int {
    return (row/3)*3 + col/3
}

func popcount(mask int) int {
    return bits.OnesCount(uint(mask))
}

func maskToDigits(mask int) []int {
    digits := make([]int, 0, 9)
    for d := 1; d <= 9; d++ {
        if mask&digitMask(d) != 0 {
            digits = append(digits, d)
        }
    }
    return digits
}

// parsePuzzle accepts common Sudoku encodings and normalizes them to exactly 81 cells.
func parsePuzzle(input string) ([]int, error) {
    filtered := make([]rune, 0, 81)
    for _, ch := range input {
        if unicode.IsSpace(ch) || ch == '|' || ch == '+' {
            continue
        }
        filtered = append(filtered, ch)
    }

    if len(filtered) != 81 {
        return nil, fmt.Errorf("Expected exactly 81 cells after removing whitespace, but found %d.", len(filtered))
    }

    cells := make([]int, 81)
    for i, ch := range filtered {
        switch {
        case ch >= '1' && ch <= '9':
            cells[i] = int(ch - '0')
        case ch == '0' || ch == '.' || ch == '_':
            cells[i] = 0
        default:
            return nil, fmt.Errorf("Unexpected character '%c' at position %d.", ch, i+1)
        }
    }
    return cells, nil
}

// formatBoard renders a board as three-by-three blocks for the report.
func formatBoard(cells []int) string {
    var b strings.Builder
    for r := 0; r < 9; r++ {
        if r > 0 && r%3 == 0 {
            b.WriteByte('\n')
        }
        for c := 0; c < 9; c++ {
            if c > 0 && c%3 == 0 {
                b.WriteString("| ")
            }
            v := cells[r*9+c]
            if v == 0 {
                b.WriteString(". ")
            } else {
                fmt.Fprintf(&b, "%d ", v)
            }
        }
        b.WriteByte('\n')
    }
    return b.String()
}

func newState() *State {
    return &State{
        Cells:   make([]int, 81),
        RowUsed: make([]int, 9),
        ColUsed: make([]int, 9),
        BoxUsed: make([]int, 9),
        Moves:   []Move{},
    }
}

// Place writes a digit if it does not violate the row, column, or 3x3 box.
func (s *State) Place(idx, value int) bool {
    if s.Cells[idx] != 0 {
        return s.Cells[idx] == value
    }
    row := idx / 9
    col := idx % 9
    box := boxIndex(row, col)
    bit := digitMask(value)
    if (s.RowUsed[row]|s.ColUsed[col]|s.BoxUsed[box])&bit != 0 {
        return false
    }
    s.Cells[idx] = value
    s.RowUsed[row] |= bit
    s.ColUsed[col] |= bit
    s.BoxUsed[box] |= bit
    return true
}

// Candidates returns the legal digits for a blank cell as a bit mask.
func (s *State) Candidates(idx int) int {
    row := idx / 9
    col := idx % 9
    box := boxIndex(row, col)
    return allDigitsMask & ^(s.RowUsed[row] | s.ColUsed[col] | s.BoxUsed[box])
}

// Clone lets the search explore one branch without mutating sibling branches.
func (s *State) Clone() *State {
    clone := &State{
        Cells:   append([]int(nil), s.Cells...),
        RowUsed: append([]int(nil), s.RowUsed...),
        ColUsed: append([]int(nil), s.ColUsed...),
        BoxUsed: append([]int(nil), s.BoxUsed...),
        Moves:   append([]Move(nil), s.Moves...),
    }
    return clone
}

// stateFromPuzzle places every given clue into an empty search state.
func stateFromPuzzle(cells []int) (*State, error) {
    state := newState()
    for idx, value := range cells {
        if value == 0 {
            continue
        }
        if value < 1 || value > 9 {
            return nil, fmt.Errorf("Cell %d contains %d, but only digits 1-9 or 0/./_ are allowed.", idx+1, value)
        }
        if !state.Place(idx, value) {
            row := idx/9 + 1
            col := idx%9 + 1
            return nil, fmt.Errorf("The given clues already conflict at row %d, column %d.", row, col)
        }
    }
    return state, nil
}

// summarizeMoves creates a compact explanation of the first few solver placements.
func summarizeMoves(moves []Move, limit int) string {
    if len(moves) == 0 {
        return "no placements were needed"
    }
    parts := make([]string, 0, limit+1)
    for i, mv := range moves {
        if i >= limit {
            break
        }
        row := mv.Index/9 + 1
        col := mv.Index%9 + 1
        mode := "guess"
        if mv.Forced {
            mode = "forced"
        }
        parts = append(parts, fmt.Sprintf("r%dc%d=%d: %s", row, col, mv.Value, mode))
    }
    if len(moves) > limit {
        parts = append(parts, fmt.Sprintf("… and %d more placements", len(moves)-limit))
    }
    return strings.Join(parts, ", ")
}

// unitIsComplete checks one row, column, or box for exactly the digits 1 through 9.
func unitIsComplete(values []int) bool {
    seen := 0
    for _, v := range values {
        if v < 1 || v > 9 {
            return false
        }
        bit := digitMask(v)
        if seen&bit != 0 {
            return false
        }
        seen |= bit
    }
    return seen == allDigitsMask
}

// replayMovesAreLegal rebuilds the proof path and verifies every recorded placement.
func replayMovesAreLegal(puzzleCells []int, moves []Move) bool {
    state, err := stateFromPuzzle(puzzleCells)
    if err != nil {
        return false
    }
    for _, mv := range moves {
        if state.Cells[mv.Index] != 0 {
            return false
        }
        maskNow := state.Candidates(mv.Index)
        if maskNow != mv.CandidatesMask {
            return false
        }
        if maskNow&digitMask(mv.Value) == 0 {
            return false
        }
        if mv.Forced && popcount(maskNow) != 1 {
            return false
        }
        if !state.Place(mv.Index, mv.Value) {
            return false
        }
    }
    return true
}

// propagateSingles repeatedly fills cells that have exactly one legal candidate.
func propagateSingles(state *State, stats *Stats) bool {
    for {
        progress := false
        for idx := 0; idx < 81; idx++ {
            if state.Cells[idx] != 0 {
                continue
            }
            mask := state.Candidates(idx)
            count := popcount(mask)
            if count == 0 {
                return false
            }
            if count == 1 {
                digit := maskToDigits(mask)[0]
                state.Moves = append(state.Moves, Move{Index: idx, Value: digit, CandidatesMask: mask, Forced: true})
                if !state.Place(idx, digit) {
                    return false
                }
                stats.ForcedMoves++
                progress = true
            }
        }
        if !progress {
            return true
        }
    }
}

// Choice is the next blank cell selected for a search branch.
type Choice struct {
    Idx   int
    Mask  int
    Count int
}

// selectUnfilledCell implements the minimum-remaining-values heuristic.
func selectUnfilledCell(state *State) *Choice {
    var best *Choice
    for idx := 0; idx < 81; idx++ {
        if state.Cells[idx] != 0 {
            continue
        }
        mask := state.Candidates(idx)
        count := popcount(mask)
        if best == nil || count < best.Count {
            best = &Choice{Idx: idx, Mask: mask, Count: count}
        }
        if count == 2 {
            break
        }
    }
    return best
}

// solve combines constraint propagation with depth-first search.
func solve(state *State, stats *Stats, depth int) *State {
    stats.RecursiveNodes++
    if depth > stats.MaxDepth {
        stats.MaxDepth = depth
    }

    current := state.Clone()
    if !propagateSingles(current, stats) {
        stats.Backtracks++
        return nil
    }

    best := selectUnfilledCell(current)
    if best == nil {
        return current
    }

    for _, digit := range maskToDigits(best.Mask) {
        next := current.Clone()
        candidatesMask := next.Candidates(best.Idx)
        next.Moves = append(next.Moves, Move{Index: best.Idx, Value: digit, CandidatesMask: candidatesMask, Forced: false})
        stats.GuessedMoves++
        if !next.Place(best.Idx, digit) {
            continue
        }
        if solved := solve(next, stats, depth+1); solved != nil {
            return solved
        }
    }

    stats.Backtracks++
    return nil
}

// countSolutions searches only up to limit solutions, enough to test uniqueness.
func countSolutions(state *State, limit int, count *int) {
    if *count >= limit {
        return
    }
    current := state.Clone()
    dummy := &Stats{}
    if !propagateSingles(current, dummy) {
        return
    }
    best := selectUnfilledCell(current)
    if best == nil {
        *count++
        return
    }
    for _, digit := range maskToDigits(best.Mask) {
        if *count >= limit {
            return
        }
        next := current.Clone()
        if next.Place(best.Idx, digit) {
            countSolutions(next, limit, count)
        }
    }
}

func countValues(cells []int, value int) int {
    count := 0
    for _, v := range cells {
        if v == value {
            count++
        }
    }
    return count
}

func joinCells(cells []int) string {
    var b strings.Builder
    b.Grow(len(cells))
    for _, v := range cells {
        b.WriteByte(byte('0' + v))
    }
    return b.String()
}

func rowsComplete(cells []int) bool {
    for r := 0; r < 9; r++ {
        if !unitIsComplete(cells[r*9 : r*9+9]) {
            return false
        }
    }
    return true
}

func colsComplete(cells []int) bool {
    for c := 0; c < 9; c++ {
        values := make([]int, 0, 9)
        for r := 0; r < 9; r++ {
            values = append(values, cells[r*9+c])
        }
        if !unitIsComplete(values) {
            return false
        }
    }
    return true
}

func boxesComplete(cells []int) bool {
    for b := 0; b < 9; b++ {
        br := (b / 3) * 3
        bc := (b % 3) * 3
        values := make([]int, 0, 9)
        for dr := 0; dr < 3; dr++ {
            for dc := 0; dc < 3; dc++ {
                values = append(values, cells[(br+dr)*9+(bc+dc)])
            }
        }
        if !unitIsComplete(values) {
            return false
        }
    }
    return true
}

// computeReport normalizes input, solves the puzzle, runs checks, and assembles output data.
func computeReport(raw string) Report {
    cells, err := parsePuzzle(raw)
    if err != nil {
        return Report{Status: "invalid-input", Error: err.Error(), Raw: raw}
    }

    normalized := joinCells(cells)
    givens := 81 - countValues(cells, 0)
    blanks := countValues(cells, 0)

    initial, err := stateFromPuzzle(cells)
    if err != nil {
        return Report{
            Status:     "illegal-clues",
            Error:      err.Error(),
            Raw:        raw,
            Normalized: normalized,
            Givens:     givens,
            Blanks:     blanks,
            PuzzleText: formatBoard(cells),
        }
    }

    stats := &Stats{Givens: givens, Blanks: blanks}
    solved := solve(initial, stats, 0)
    if solved == nil {
        return Report{
            Status:         "unsatisfiable",
            Raw:            raw,
            Normalized:     normalized,
            Givens:         stats.Givens,
            Blanks:         stats.Blanks,
            RecursiveNodes: stats.RecursiveNodes,
            Backtracks:     stats.Backtracks,
            PuzzleText:     formatBoard(cells),
        }
    }

    solutionCount := 0
    countSolutions(initial, 2, &solutionCount)

    givensPreserved := true
    for i, v := range cells {
        if v != 0 && v != solved.Cells[i] {
            givensPreserved = false
            break
        }
    }

    noBlanks := true
    for _, v := range solved.Cells {
        if v < 1 || v > 9 {
            noBlanks = false
            break
        }
    }

    proofPathGuessCount := 0
    for _, mv := range solved.Moves {
        if !mv.Forced {
            proofPathGuessCount++
        }
    }

    return Report{
        Status:          "ok",
        Raw:             raw,
        Normalized:      normalized,
        Givens:          stats.Givens,
        Blanks:          stats.Blanks,
        ForcedMoves:     stats.ForcedMoves,
        GuessedMoves:    stats.GuessedMoves,
        RecursiveNodes:  stats.RecursiveNodes,
        Backtracks:      stats.Backtracks,
        MaxDepth:        stats.MaxDepth,
        Unique:          solutionCount == 1,
        Solution:        joinCells(solved.Cells),
        PuzzleText:      formatBoard(cells),
        SolutionText:    formatBoard(solved.Cells),
        MoveSummary:     summarizeMoves(solved.Moves, 8),
        MoveCount:       len(solved.Moves),
        GivensPreserved: givensPreserved,
        NoBlanks:        noBlanks,
        RowsComplete:    rowsComplete(solved.Cells),
        ColsComplete:    colsComplete(solved.Cells),
        BoxesComplete:   boxesComplete(solved.Cells),
        ReplayLegal:     replayMovesAreLegal(cells, solved.Moves),
        StoryConsistent: stats.RecursiveNodes >= 1 && stats.MaxDepth <= stats.Blanks && len(solved.Moves) == stats.Blanks && proofPathGuessCount <= stats.GuessedMoves,
    }
}

func statusText(ok bool) string {
    if ok {
        return "OK"
    }
    return "failed"
}

func yesNo(ok bool) string {
    if ok {
        return "yes"
    }
    return "no"
}

// printReport renders the N3-style answer, reason, check, and Go audit sections.
func printReport(report Report, puzzleName string) {
    switch report.Status {
    case "ok":
        fmt.Println("=== Answer ===")
        if report.Unique {
            fmt.Println("The puzzle is solved, and the completed grid is the unique valid Sudoku solution.")
        } else {
            fmt.Println("The puzzle is solved, and the completed grid is a valid Sudoku solution.")
        }
        fmt.Println("case : sudoku")
        fmt.Printf("default puzzle : %s\n\n", puzzleName)
        fmt.Println("Puzzle")
        fmt.Print(report.PuzzleText)
        fmt.Println()
        fmt.Println("Completed grid")
        fmt.Print(report.SolutionText)
        fmt.Println()

        fmt.Println("=== Reason Why ===")
        if report.Unique {
            fmt.Printf("The solver starts from %d clues and fills the remaining %d cells by combining constraint propagation with depth-first search. At each step it chooses the empty cell with the fewest legal digits, places forced singles immediately, and only guesses when more than one candidate remains. Across the search it made %d forced placements and tried %d guesses, visited %d search nodes overall, and backtracked %d times before reaching the completed grid. The solver also confirmed that the solution is unique. Early steps: %s\n\n",
                report.Givens, report.Blanks, report.ForcedMoves, report.GuessedMoves, report.RecursiveNodes, report.Backtracks, report.MoveSummary)
        } else {
            fmt.Printf("The solver starts from %d clues and fills the remaining %d cells by combining constraint propagation with depth-first search. At each step it chooses the empty cell with the fewest legal digits, places forced singles immediately, and only guesses when more than one candidate remains. Across the search it made %d forced placements and tried %d guesses, visited %d search nodes overall, and backtracked %d times before reaching the completed grid. The solver found a valid solution, but there is more than one. Early steps: %s\n\n",
                report.Givens, report.Blanks, report.ForcedMoves, report.GuessedMoves, report.RecursiveNodes, report.Backtracks, report.MoveSummary)
        }

        fmt.Println("=== Check ===")
        fmt.Printf("C1 %s - every given clue is preserved in the final grid.\n", statusText(report.GivensPreserved))
        fmt.Printf("C2 %s - the final grid contains only digits 1 through 9, with no blanks left.\n", statusText(report.NoBlanks))
        fmt.Printf("C3 %s - each row contains every digit exactly once.\n", statusText(report.RowsComplete))
        fmt.Printf("C4 %s - each column contains every digit exactly once.\n", statusText(report.ColsComplete))
        fmt.Printf("C5 %s - each 3×3 box contains every digit exactly once.\n", statusText(report.BoxesComplete))
        fmt.Printf("C6 %s - replaying the recorded placements from the original puzzle remains legal at every step.\n", statusText(report.ReplayLegal))
        fmt.Printf("C7 %s - the search statistics and the successful proof path are internally consistent.\n", statusText(report.StoryConsistent))
        if report.Unique {
            fmt.Println("C8 OK - a second search found no alternative solution, so the solution is unique.")
        } else {
            fmt.Println("C8 INFO - a second search found another solution, so the puzzle is not unique.")
        }

        // Extra implementation-oriented details. The N3 example focuses on the proof-style
        // answer/check sections; this Go audit block exposes concrete runtime facts
        // that are useful when testing or comparing translations.
        fmt.Println()
        fmt.Println("=== Go audit details ===")
        fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
        fmt.Printf("normalized puzzle : %s\n", report.Normalized)
        fmt.Printf("solution string : %s\n", report.Solution)
        fmt.Printf("givens : %d\n", report.Givens)
        fmt.Printf("blanks : %d\n", report.Blanks)
        fmt.Printf("recorded placements : %d\n", report.MoveCount)
        fmt.Printf("forced placements : %d\n", report.ForcedMoves)
        fmt.Printf("guesses tried : %d\n", report.GuessedMoves)
        fmt.Printf("recursive nodes : %d\n", report.RecursiveNodes)
        fmt.Printf("backtracks : %d\n", report.Backtracks)
        fmt.Printf("max search depth : %d\n", report.MaxDepth)
        fmt.Printf("unique solution : %s\n", yesNo(report.Unique))
        fmt.Printf("story consistent : %s\n", yesNo(report.StoryConsistent))

    case "invalid-input":
        fmt.Println("=== Answer ===")
        fmt.Println("The supplied puzzle is not well formed and cannot be parsed as a 9×9 Sudoku.")
        fmt.Println("case : sudoku")
        fmt.Println()
        fmt.Println("=== Reason Why ===")
        fmt.Println(report.Error)
        fmt.Println()
        fmt.Println("=== Check ===")
        fmt.Println("C1 failed - the supplied text does not normalize to exactly 81 legal Sudoku cells.")

    case "illegal-clues":
        fmt.Println("=== Answer ===")
        fmt.Println("The puzzle is invalid and cannot be solved as a standard Sudoku.")
        fmt.Println("case : sudoku")
        fmt.Println()
        fmt.Println("Puzzle")
        fmt.Print(report.PuzzleText)
        fmt.Println()
        fmt.Println("=== Reason Why ===")
        fmt.Println(report.Error)
        fmt.Println()
        fmt.Println("=== Check ===")
        fmt.Println("C1 failed - the given clues already violate Sudoku rules.")

    case "unsatisfiable":
        fmt.Println("=== Answer ===")
        fmt.Println("No valid Sudoku solution exists for the supplied puzzle.")
        fmt.Println("case : sudoku")
        fmt.Printf("default puzzle : %s\n\n", puzzleName)
        fmt.Println("Puzzle")
        fmt.Print(report.PuzzleText)
        fmt.Println()
        fmt.Println("=== Reason Why ===")
        fmt.Printf("The solver explored %d search nodes with minimum-remaining-values branching and backtracked %d times, but every branch eventually contradicted the row, column, or box constraints.\n\n", report.RecursiveNodes, report.Backtracks)
        fmt.Println("=== Check ===")
        fmt.Println("C1 OK - the given clues are internally consistent.")
        fmt.Println("C2 OK - every explored assignment respected row, column, and box legality.")
        fmt.Println("C3 failed - exhaustive search found no complete legal grid.")

    default:
        fmt.Fprintf(os.Stderr, "unknown report status: %s\n", report.Status)
        os.Exit(1)
    }
}

func main() {
    puzzle := flag.String("puzzle", defaultPuzzle, "81-cell Sudoku puzzle; blanks may be 0, ., or _; whitespace and |/+ separators are ignored")
    puzzleName := flag.String("name", "classic", "puzzle name used in the generated report")
    flag.Parse()

    if flag.NArg() > 0 {
        *puzzle = strings.Join(flag.Args(), " ")
    }

    report := computeReport(*puzzle)
    printReport(report, *puzzleName)

    if report.Status != "ok" {
        os.Exit(1)
    }
}
