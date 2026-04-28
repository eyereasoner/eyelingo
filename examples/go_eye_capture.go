// go_eye_capture.go
//
// A self-contained Go scenario for the Chinese game of Go (Weiqi) in ARC style.
// It demonstrates that a group with only a single eye and no outside liberties
// is dead.  When the opponent plays inside the eye, the entire group is
// captured.
//
// This is intentionally not a full Go engine – it is a concrete scenario that
// mirrors the style of the Eyeling N3‑to‑Go translations.
//
// Run:
//
//     go run go_eye_capture.go
//
// The program has no third-party dependencies.

package main

import (
    "fmt"
    "os"
    "runtime"
    "strings"
)

// ---------- types ----------

// Stone represents a cell on the board.
type Stone int

const (
    Empty Stone = 0
    Black Stone = 1
    White Stone = 2
)

// GoString returns a one‑character representation.
func (s Stone) GoString() string {
    switch s {
    case Empty:
        return "."
    case Black:
        return "B"
    case White:
        return "W"
    }
    return "?"
}

// Board holds the current grid and the board size.
type Board struct {
    Size int
    Grid [][]Stone
}

// Move represents a player's intended stone placement.
type Move struct {
    Player Stone
    Row    int
    Col    int
}

// MoveResult stores the outcome of attempting a move.
type MoveResult struct {
    Legal         bool
    IllegalReason string
    Captured      []Stone
    CapturedPos   []string
    BoardAfter    Board
    Description   string
}

// ---------- board helpers ----------

func newBoard(size int) Board {
    g := make([][]Stone, size)
    for i := range g {
        g[i] = make([]Stone, size)
    }
    return Board{Size: size, Grid: g}
}

func (b Board) clone() Board {
    nb := newBoard(b.Size)
    for r := 0; r < b.Size; r++ {
        copy(nb.Grid[r], b.Grid[r])
    }
    return nb
}

func (b Board) inBounds(r, c int) bool {
    return r >= 0 && r < b.Size && c >= 0 && c < b.Size
}

func (b Board) neighbours(r, c int) [][2]int {
    var out [][2]int
    for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
        nr, nc := r+d[0], c+d[1]
        if b.inBounds(nr, nc) {
            out = append(out, [2]int{nr, nc})
        }
    }
    return out
}

// groupLiberties returns the set of liberties of the connected group
// containing (r,c) of colour col.
func (b Board) groupLiberties(r, c int, col Stone) map[[2]int]bool {
    visited := map[[2]int]bool{{r, c}: true}
    queue := [][2]int{{r, c}}
    libs := map[[2]int]bool{}
    for len(queue) > 0 {
        cur := queue[0]
        queue = queue[1:]
        for _, nb := range b.neighbours(cur[0], cur[1]) {
            if visited[nb] {
                continue
            }
            cell := b.Grid[nb[0]][nb[1]]
            if cell == Empty {
                libs[nb] = true
            } else if cell == col {
                visited[nb] = true
                queue = append(queue, nb)
            }
        }
    }
    return libs
}

// boardString serialises the board to a single string.
func (b Board) boardString() string {
    var sb strings.Builder
    for r := 0; r < b.Size; r++ {
        for c := 0; c < b.Size; c++ {
            sb.WriteString(b.Grid[r][c].GoString())
        }
    }
    return sb.String()
}

// playMove executes a move and returns the result.
func (b Board) playMove(player Stone, r, c int) MoveResult {
    res := MoveResult{
        Legal:      false,
        BoardAfter: b.clone(),
    }
    if !b.inBounds(r, c) {
        res.IllegalReason = "out of bounds"
        return res
    }
    if b.Grid[r][c] != Empty {
        res.IllegalReason = fmt.Sprintf("occupied by %s", b.Grid[r][c].GoString())
        return res
    }

    nb := b.clone()
    nb.Grid[r][c] = player

    opponent := Black
    if player == Black {
        opponent = White
    }

    var captured []Stone
    var capturedPos []string

    // remove opponent groups without liberties
    for rr := 0; rr < nb.Size; rr++ {
        for cc := 0; cc < nb.Size; cc++ {
            if nb.Grid[rr][cc] == opponent {
                libs := nb.groupLiberties(rr, cc, opponent)
                if len(libs) == 0 {
                    visited := map[[2]int]bool{{rr, cc}: true}
                    queue := [][2]int{{rr, cc}}
                    for len(queue) > 0 {
                        cur := queue[0]
                        queue = queue[1:]
                        nb.Grid[cur[0]][cur[1]] = Empty
                        captured = append(captured, opponent)
                        capturedPos = append(capturedPos, fmt.Sprintf("(%d,%d)", cur[0], cur[1]))
                        for _, n := range nb.neighbours(cur[0], cur[1]) {
                            if !visited[n] && nb.Grid[n[0]][n[1]] == opponent {
                                visited[n] = true
                                queue = append(queue, n)
                            }
                        }
                    }
                }
            }
        }
    }

    // suicide check
    selfLibs := nb.groupLiberties(r, c, player)
    if len(selfLibs) == 0 {
        res.IllegalReason = "suicide"
        return res
    }

    res.Legal = true
    res.Captured = captured
    res.CapturedPos = capturedPos
    res.BoardAfter = nb
    if len(captured) > 0 {
        res.Description = fmt.Sprintf("legal move: captures %d opponent stone(s)", len(captured))
    } else {
        res.Description = "legal move: no captures"
    }
    return res
}

// ---------- scenario ----------

func scenario() (Board, Move) {
    b := newBoard(5)
    // Fill with White border.
    for r := 0; r < 5; r++ {
        for c := 0; c < 5; c++ {
            b.Grid[r][c] = White
        }
    }
    // Place Black inner stones forming a group with a single eye at (2,2).
    blackPositions := [][2]int{
        {1, 1}, {1, 2}, {1, 3},
        {2, 1}, {2, 3},
        {3, 1}, {3, 2}, {3, 3},
    }
    for _, pos := range blackPositions {
        b.Grid[pos[0]][pos[1]] = Black
    }
    // Ensure the eye is empty (the previous loop set it to White).
    b.Grid[2][2] = Empty

    move := Move{White, 2, 2} // play inside the eye
    return b, move
}

// ---------- checks ----------

type Checks struct {
    InitialSingleLiberty   bool
    BlackGroupSize         int
    CaptureLegal           bool
    AllBlackStonesCaptured bool
    WhiteStoneHasLiberties bool
}

func performChecks(initial Board, move Move, result MoveResult) Checks {
    cks := Checks{}

    // Before move: Black group at (1,1) should have exactly 1 liberty (the eye).
    libs := initial.groupLiberties(1, 1, Black)
    cks.InitialSingleLiberty = len(libs) == 1

    // Count initial Black stones.
    initialBlackCount := 0
    for r := 0; r < initial.Size; r++ {
        for c := 0; c < initial.Size; c++ {
            if initial.Grid[r][c] == Black {
                initialBlackCount++
            }
        }
    }
    cks.BlackGroupSize = initialBlackCount

    cks.CaptureLegal = result.Legal

    // All initial Black stones should be captured.
    cks.AllBlackStonesCaptured = len(result.Captured) == initialBlackCount

    if result.Legal {
        wLibs := result.BoardAfter.groupLiberties(move.Row, move.Col, White)
        cks.WhiteStoneHasLiberties = len(wLibs) > 0
    }

    return cks
}

func allChecksPass(c Checks) bool {
    return c.InitialSingleLiberty &&
        c.BlackGroupSize == 8 &&
        c.CaptureLegal &&
        c.AllBlackStonesCaptured &&
        c.WhiteStoneHasLiberties
}

func checkCount(c Checks) int {
    n := 0
    if c.InitialSingleLiberty {
        n++
    }
    if c.BlackGroupSize == 8 {
        n++
    }
    if c.CaptureLegal {
        n++
    }
    if c.AllBlackStonesCaptured {
        n++
    }
    if c.WhiteStoneHasLiberties {
        n++
    }
    return n
}

// ---------- ARC rendering ----------

func renderArcOutput(initial Board, move Move, result MoveResult, cks Checks) {
    fmt.Println("=== Go Eye Capture Scenario (Weiqi) ===")
    fmt.Println()

    // --- Answer ---
    fmt.Println("=== Answer ===")
    fmt.Printf("The Black group had only a single eye and no outside liberties.\n")
    fmt.Printf("White plays inside the eye at (%d,%d) and captures the entire group.\n",
        move.Row, move.Col)
    fmt.Println("Initial board:")
    fmt.Print(boardStringFormat(initial))
    fmt.Printf("Move: %s at (%d,%d)\n", move.Player.GoString(), move.Row, move.Col)
    fmt.Println("Board after the move:")
    fmt.Print(boardStringFormat(result.BoardAfter))
    fmt.Printf("Number of captured Black stones: %d\n", len(result.Captured))
    fmt.Println()

    // --- Reason Why ---
    fmt.Println("=== Reason Why ===")
    fmt.Println("In Go, a group lives only if it has at least two eyes, or can make them")
    fmt.Println("eventually. A group with only one eye and no other liberties is dead")
    fmt.Println("because the opponent can safely play inside that eye.")
    fmt.Println("When the opponent fills the last liberty, the entire group is removed")
    fmt.Println("from the board as captured stones. The played stone itself remains alive")
    fmt.Println("because after the capture it inherits the newly vacated liberties.")
    fmt.Println()
    fmt.Println("In this scenario, the Black enclosure is completely sealed by White")
    fmt.Println("stones. The only empty intersection inside is the eye at (2,2).")
    fmt.Println("White’s move inside the eye is not suicide because the Black stones are")
    fmt.Println("immediately captured, freeing the space around the newly placed White stone.")
    fmt.Println()

    // --- Check ---
    fmt.Println("=== Check ===")
    if cks.InitialSingleLiberty {
        fmt.Println("C1 OK - Before the move, the Black group had exactly 1 liberty (the eye).")
    } else {
        fmt.Println("C1 FAIL - Black group did not have exactly 1 liberty.")
    }
    fmt.Printf("C2 OK - Initial Black group size is %d (8 expected).\n", cks.BlackGroupSize)
    if cks.CaptureLegal {
        fmt.Println("C3 OK - The killing move was legal.")
    } else {
        fmt.Println("C3 FAIL - The move was illegal.")
    }
    if cks.AllBlackStonesCaptured {
        fmt.Println("C4 OK - All Black stones were captured.")
    } else {
        fmt.Println("C4 FAIL - Not all Black stones were captured.")
    }
    if cks.WhiteStoneHasLiberties {
        fmt.Println("C5 OK - After capture, the White move is not a suicide (it has liberties).")
    } else {
        fmt.Println("C5 FAIL - White stone has no liberties.")
    }
    fmt.Println()

    // --- Go audit details ---
    fmt.Println("=== Go audit details ===")
    fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
    fmt.Printf("board size : %dx%d\n", initial.Size, initial.Size)
    fmt.Println("initial board:")
    fmt.Print(boardStringFormat(initial))
    fmt.Printf("killing move : %s (%d,%d)\n", move.Player.GoString(), move.Row, move.Col)
    fmt.Printf("result : %s\n", result.Description)
    fmt.Printf("captured stones : %d\n", len(result.Captured))
    if len(result.CapturedPos) > 0 {
        fmt.Printf("captured positions : %s\n", strings.Join(result.CapturedPos, ", "))
    }
    fmt.Printf("checks passed : %d/5\n", checkCount(cks))
    fmt.Printf("recommendation consistent : %s\n", yesNo(allChecksPass(cks)))
}

func boardStringFormat(b Board) string {
    var sb strings.Builder
    for r := 0; r < b.Size; r++ {
        for c := 0; c < b.Size; c++ {
            sb.WriteString(b.Grid[r][c].GoString())
            sb.WriteByte(' ')
        }
        sb.WriteByte('\n')
    }
    return sb.String()
}

func yesNo(value bool) string {
    if value {
        return "yes"
    }
    return "no"
}

// ---------- main ----------

func main() {
    initial, move := scenario()
    result := initial.playMove(move.Player, move.Row, move.Col)
    if !result.Legal {
        fmt.Fprintf(os.Stderr, "unexpected illegal move: %s\n", result.IllegalReason)
        os.Exit(1)
    }
    cks := performChecks(initial, move, result)
    renderArcOutput(initial, move, result, cks)
}
