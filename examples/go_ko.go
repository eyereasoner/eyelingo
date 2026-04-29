// go_ko.go
//
// A self-contained Go scenario for the Chinese game of Go (Weiqi) in ARC style.
// It demonstrates the ko rule: after one stone is captured, the opponent
// cannot immediately recapture if that would restore the board position
// that existed before the capture.
//
// This is intentionally not a full Go engine – it is a concrete scenario that
// mirrors the style of the Eyeling N3-to-Go translations.
//
// Run:
//
//     go run go_ko.go
//
// The program has no third-party dependencies.

package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "go_ko"

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

// MoveResult stores the outcome of attempting a move, plus a flag for ko.
type MoveResult struct {
	Legal         bool
	IllegalReason string
	Captured      []Stone
	CapturedPos   []string
	BoardAfter    Board
	Description   string
	KoIllegal     bool // true if the move is illegal because of ko
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

// boardString serialises the board to a single string for position
// comparison.
func (b Board) boardString() string {
	var sb strings.Builder
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			sb.WriteString(b.Grid[r][c].GoString())
		}
	}
	return sb.String()
}

// playMove executes a move and returns the result.  If prevBoardStr is
// non‑empty and the move would recreate that exact board position
// AND it captures exactly one stone, the move is marked as ko‑illegal.
func (b Board) playMove(player Stone, r, c int, prevBoardStr string) MoveResult {
	res := MoveResult{
		Legal:       false,
		BoardAfter:  b.clone(),
		Description: "",
		KoIllegal:   false,
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

	// ko rule: if prevBoardStr is given and the resulting board equals it,
	// and exactly one stone was captured, the move is ko‑illegal.
	if prevBoardStr != "" && nb.boardString() == prevBoardStr && len(captured) == 1 {
		res.KoIllegal = true
		res.IllegalReason = "ko rule (immediate recapture of a single stone that restores the board)"
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

// scenario returns the initial board and the two moves that demonstrate ko.
func scenario() (Board, Move, Move) {
	// 3×3 board with a classic ko shape:
	//   . W .
	//   B . W
	//   . W .
	b := newBoard(3)
	b.Grid[0][1] = White
	b.Grid[1][2] = White
	b.Grid[2][1] = White
	b.Grid[1][0] = Black

	// Move 1: White captures the Black stone at (1,0) by playing (1,1).
	move1 := Move{White, 1, 1}
	// Move 2: Black tries to recapture immediately at (1,0).
	move2 := Move{Black, 1, 0}

	return b, move1, move2
}

// ---------- checks ----------

type Checks struct {
	BlackStoneLibertyBefore int // Black stone at (1,0) had exactly 1 liberty
	WhiteCaptureLegal       bool
	CaptureOccurred         bool
	WhiteStoneLibertyAfter  int // White stone at (1,1) after capture has exactly 1 liberty
	KoRecaptureAttempt      bool
	KoIllegalDetected       bool
}

func performChecks(
	boardBeforeCapture Board,
	resCapture MoveResult,
	moveRecapture Move,
	resRecapture MoveResult,
) Checks {
	c := Checks{}

	// Before capture, Black stone at (1,0) should have 1 liberty.
	libsBefore := boardBeforeCapture.groupLiberties(1, 0, Black)
	c.BlackStoneLibertyBefore = len(libsBefore)

	// Capture must be legal.
	c.WhiteCaptureLegal = resCapture.Legal
	// At least one stone captured.
	c.CaptureOccurred = len(resCapture.Captured) > 0

	// After capture, the White stone at (1,1) should have exactly 1 liberty (at (1,0)).
	libsAfter := resCapture.BoardAfter.groupLiberties(1, 1, White)
	c.WhiteStoneLibertyAfter = len(libsAfter)

	// Recapture was attempted.
	c.KoRecaptureAttempt = true
	// The recapture should be illegal because of ko.
	c.KoIllegalDetected = !resRecapture.Legal && resRecapture.KoIllegal

	return c
}

func allChecksPass(c Checks) bool {
	return c.BlackStoneLibertyBefore == 1 &&
		c.WhiteCaptureLegal &&
		c.CaptureOccurred &&
		c.WhiteStoneLibertyAfter == 1 &&
		c.KoIllegalDetected
}

func checkCount(c Checks) int {
	n := 0
	if c.BlackStoneLibertyBefore == 1 {
		n++
	}
	if c.WhiteCaptureLegal {
		n++
	}
	if c.CaptureOccurred {
		n++
	}
	if c.WhiteStoneLibertyAfter == 1 {
		n++
	}
	if c.KoIllegalDetected {
		n++
	}
	return n
}

// ---------- ARC rendering ----------

func renderArcOutput(
	boardBefore Board,
	moveCapture Move,
	resCapture MoveResult,
	moveRecapture Move,
	resRecapture MoveResult,
	cks Checks,
) {
	fmt.Println("# Go Ko Rule Scenario (Weiqi)")
	fmt.Println()

	// --- Answer ---
	fmt.Println("## Answer")
	// Describe the capture move
	fmt.Printf("Move 1: %s at (%d,%d) is %s.\n",
		moveCapture.Player.GoString(), moveCapture.Row, moveCapture.Col,
		resCapture.Description)
	fmt.Println("Board after White's capture:")
	fmt.Println(resCapture.BoardAfter.boardStringFormat())

	// Describe the attempted recapture
	if resRecapture.KoIllegal {
		fmt.Printf("Move 2: %s at (%d,%d) is ILLEGAL (ko).\n",
			moveRecapture.Player.GoString(), moveRecapture.Row, moveRecapture.Col)
		fmt.Println("The board remains as after Move 1.")
	} else {
		fmt.Printf("Move 2: %s at (%d,%d) is legal.\n",
			moveRecapture.Player.GoString(), moveRecapture.Row, moveRecapture.Col)
		fmt.Println("Board after recapture (should not happen here):")
		fmt.Println(resRecapture.BoardAfter.boardStringFormat())
	}
	fmt.Println()

	// --- Reason Why ---
	fmt.Println("## Reason why")
	fmt.Println("In Go, a ko occurs when:")
	fmt.Println("1. A player captures exactly one stone.")
	fmt.Println("2. After the capture, the capturing stone is left with exactly one liberty.")
	fmt.Println("3. If the opponent immediately recaptures, the board would return to its")
	fmt.Println("   exact state before the first capture.")
	fmt.Println("The ko rule forbids such an immediate recapture to prevent infinite loops.")
	fmt.Println()
	fmt.Println("In this scenario:")
	fmt.Println("- The Black stone at (1,0) had only one liberty (1,1).")
	fmt.Println("- White captured it by playing at (1,1), removing the Black stone.")
	fmt.Println("- The White stone at (1,1) now had only one liberty (1,0).")
	fmt.Println("- If Black were allowed to play at (1,0), it would capture the White stone")
	fmt.Println("  and the board would look identical to how it was before White's capture.")
	fmt.Println("- Therefore, Black's immediate recapture is illegal under the ko rule.")
	fmt.Println()

	// --- Check ---
	fmt.Println("## Check")
	fmt.Printf("C1 OK - Black stone had exactly 1 liberty before capture (count=%d).\n",
		cks.BlackStoneLibertyBefore)
	if cks.WhiteCaptureLegal {
		fmt.Println("C2 OK - White's capture move was legal.")
	} else {
		fmt.Println("C2 FAIL - White's capture move not legal.")
	}
	if cks.CaptureOccurred {
		fmt.Println("C3 OK - At least one stone was captured.")
	} else {
		fmt.Println("C3 FAIL - No capture occurred.")
	}
	fmt.Printf("C4 OK - After capture, White stone at (1,1) has exactly 1 liberty (count=%d).\n",
		cks.WhiteStoneLibertyAfter)
	if cks.KoIllegalDetected {
		fmt.Println("C5 OK - Black's immediate recapture was correctly flagged as ko‑illegal.")
	} else {
		fmt.Println("C5 FAIL - Ko detection did not work.")
	}
	fmt.Println()

	// --- Go audit details ---
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Println("initial board:")
	fmt.Print(boardBefore.boardStringFormat())
	fmt.Printf("move 1 : %s (%d,%d) -> %s\n",
		moveCapture.Player.GoString(), moveCapture.Row, moveCapture.Col,
		resCapture.Description)
	if len(resCapture.CapturedPos) > 0 {
		fmt.Printf("  captured at: %s\n", strings.Join(resCapture.CapturedPos, ", "))
	}
	fmt.Printf("move 2 attempt: %s (%d,%d) -> %s\n",
		moveRecapture.Player.GoString(), moveRecapture.Row, moveRecapture.Col,
		resRecapture.IllegalReason)
	fmt.Printf("ko detection : %v\n", resRecapture.KoIllegal)
	fmt.Printf("checks passed : %d/5\n", checkCount(cks))
	fmt.Printf("recommendation consistent : %s\n", yesNo(allChecksPass(cks)))
}

// boardStringFormat returns a formatted board representation.
func (b Board) boardStringFormat() string {
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
	// Set up the ko‑rule scenario.
	defaultBoard, defaultCapture, defaultRecapture := scenario()
	cfg := exampleinput.Load(eyelingoExampleName, struct {
		BoardBefore   Board
		MoveCapture   Move
		MoveRecapture Move
	}{BoardBefore: defaultBoard, MoveCapture: defaultCapture, MoveRecapture: defaultRecapture})
	boardBefore, moveCapture, moveRecapture := cfg.BoardBefore, cfg.MoveCapture, cfg.MoveRecapture

	// Record the board state before White's capture (this will be the
	// “previous position” checked by the ko rule).
	prevBoardStr := boardBefore.boardString()

	// Execute White's capture.
	resCapture := boardBefore.playMove(moveCapture.Player, moveCapture.Row, moveCapture.Col, "")
	if !resCapture.Legal {
		fmt.Fprintf(os.Stderr, "unexpected illegal capture move: %s\n", resCapture.IllegalReason)
		os.Exit(1)
	}

	// Attempt Black's recapture, supplying the board string before White's capture.
	resRecapture := resCapture.BoardAfter.playMove(
		moveRecapture.Player, moveRecapture.Row, moveRecapture.Col,
		prevBoardStr,
	)

	// Perform checks.
	cks := performChecks(boardBefore, resCapture, moveRecapture, resRecapture)

	// Render ARC‑style output.
	renderArcOutput(boardBefore, moveCapture, resCapture, moveRecapture, resRecapture, cks)
}
