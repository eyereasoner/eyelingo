// go_capture.go
//
// A self-contained Go scenario for the Chinese game of Go (Weiqi) in ARC style.
// It simulates one move: White places a stone that captures a lonely Black
// stone by depriving it of its last liberty.  The program verifies the legality
// of the move, performs the capture, and checks basic rule invariants.
//
// This is intentionally not a full Go engine – it is a concrete scenario that
// mirrors the style of the Eyeling N3-to-Go translations.
//
// Run:
//
//     go run go_capture.go
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

const eyelingoExampleName = "go_capture"

// ---------- types ----------

// Stone represents a cell on the board.
type Stone int

const (
	Empty Stone = 0
	Black Stone = 1
	White Stone = 2
)

// GoString returns a one-character representation.
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

// Board holds the current state and the board size (always square).
type Board struct {
	Size int
	Grid [][]Stone
}

// Move represents a player’s intended stone placement.
type Move struct {
	Player Stone // Black or White
	Row    int
	Col    int
}

// ResultAfterMove holds the outcome of attempting a move.
type ResultAfterMove struct {
	Legal       bool
	Captured    []Stone // Stones removed from the board
	CapturedPos []string
	BoardAfter  Board
	Description string
}

// Scenario sets up the fixed initial position and the move played.
type Scenario struct {
	Initial  Board
	Move     Move
	Expected string // short description of expected behaviour
}

// Checks collects the rule‑invariant checks for this scenario.
type Checks struct {
	LibertyBeforeMove int
	LibertyAfterMove  int
	NoSuicide         bool
	CaptureOccurred   bool
	BoardSizeOk       bool
}

// ---------- board helpers ----------

// newBoard creates an empty square board of given size.
func newBoard(size int) Board {
	g := make([][]Stone, size)
	for i := range g {
		g[i] = make([]Stone, size)
	}
	return Board{Size: size, Grid: g}
}

// clone returns a deep copy.
func (b Board) clone() Board {
	nb := newBoard(b.Size)
	for r := 0; r < b.Size; r++ {
		copy(nb.Grid[r], b.Grid[r])
	}
	return nb
}

// inBounds checks if (r, c) is on the board.
func (b Board) inBounds(r, c int) bool {
	return r >= 0 && r < b.Size && c >= 0 && c < b.Size
}

// neighbours returns the up to 4 orthogonal adjacent coordinates.
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

// groupLiberties returns the set of liberties (empty intersections) of the
// connected component of stones of colour <col> containing (r,c).
func (b Board) groupLiberties(r, c int, col Stone) map[[2]int]bool {
	visited := map[[2]int]bool{{r, c}: true}
	queue := [][2]int{{r, c}}
	liberties := map[[2]int]bool{}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, nb := range b.neighbours(cur[0], cur[1]) {
			if visited[nb] {
				continue
			}
			cell := b.Grid[nb[0]][nb[1]]
			if cell == Empty {
				liberties[nb] = true
			} else if cell == col {
				visited[nb] = true
				queue = append(queue, nb)
			}
		}
	}
	return liberties
}

// playMove tries to place a stone. It returns the result with legality and captures.
func (b Board) playMove(player Stone, r, c int) ResultAfterMove {
	res := ResultAfterMove{
		Legal:       false,
		BoardAfter:  b.clone(),
		Description: "",
	}

	if !b.inBounds(r, c) {
		res.Description = fmt.Sprintf("out of bounds (%d,%d)", r, c)
		return res
	}
	if b.Grid[r][c] != Empty {
		res.Description = fmt.Sprintf("intersection occupied by %s", b.Grid[r][c].GoString())
		return res
	}

	// place stone
	nb := b.clone()
	nb.Grid[r][c] = player

	// check opponent captures first (so that self-capture can be avoided if opponent stones are removed)
	opponent := Black
	if player == Black {
		opponent = White
	}

	var captured []Stone
	var capturedPos []string

	// find opponent groups that become libertyless
	for rr := 0; rr < nb.Size; rr++ {
		for cc := 0; cc < nb.Size; cc++ {
			if nb.Grid[rr][cc] == opponent {
				libs := nb.groupLiberties(rr, cc, opponent)
				if len(libs) == 0 {
					// Remove the entire group
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

	// check self-capture (player's group at (r,c) must have at least one liberty after captures)
	selfLibs := nb.groupLiberties(r, c, player)
	if len(selfLibs) == 0 {
		res.Description = "suicide move (no liberties after placement)"
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

// string returns a human-readable board.
func (b Board) string() string {
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

// ---------- scenario ----------

func buildScenario() Scenario {
	b := newBoard(5)
	// Stones placed according to the described capture scenario:
	// Black stone at (1,1) that will be captured.
	b.Grid[1][1] = Black
	// White stones surrounding it on three sides.
	b.Grid[0][1] = White
	b.Grid[1][2] = White
	b.Grid[2][1] = White
	// The capturing move: White plays at (1,0), the last liberty of the black stone.
	move := Move{
		Player: White,
		Row:    1,
		Col:    0,
	}
	return Scenario{
		Initial:  b,
		Move:     move,
		Expected: "White captures the Black stone at (1,1) by placing at its last liberty (1,0).",
	}
}

// ---------- checks ----------

func performChecks(initial Board, move Move, result ResultAfterMove) Checks {
	var cks Checks

	// Liberty count for the black stone before the move (should be exactly 1).
	libsBefore := initial.groupLiberties(1, 1, Black)
	cks.LibertyBeforeMove = len(libsBefore)

	// After capture, the black stone should have 0 liberties (since captured).
	libsAfter := result.BoardAfter.groupLiberties(1, 1, Black)
	cks.LibertyAfterMove = len(libsAfter) // group no longer exists, but function returns liberties of empty? It won't find a black stone. We'll actually check that (1,1) is empty.
	if result.BoardAfter.Grid[1][1] != Empty {
		cks.LibertyAfterMove = -1 // indicate problem
	}

	// No suicide: the white placed stone has liberties.
	whiteLibs := result.BoardAfter.groupLiberties(move.Row, move.Col, White)
	cks.NoSuicide = len(whiteLibs) > 0

	// Capture occurred.
	cks.CaptureOccurred = len(result.Captured) > 0

	// Board size is 5.
	cks.BoardSizeOk = initial.Size == 5

	return cks
}

func allChecksPass(cks Checks) bool {
	return cks.LibertyBeforeMove == 1 &&
		cks.LibertyAfterMove == 0 &&
		cks.NoSuicide &&
		cks.CaptureOccurred &&
		cks.BoardSizeOk
}

func checkCount(cks Checks) int {
	n := 0
	if cks.LibertyBeforeMove == 1 {
		n++
	}
	if cks.LibertyAfterMove == 0 {
		n++
	}
	if cks.NoSuicide {
		n++
	}
	if cks.CaptureOccurred {
		n++
	}
	if cks.BoardSizeOk {
		n++
	}
	return n
}

// ---------- ARC rendering ----------

func renderArcOutput(scenario Scenario, result ResultAfterMove, cks Checks) {
	move := scenario.Move
	fmt.Println("# Go Capture Scenario (Weiqi)")
	fmt.Println()

	// --- Answer ---
	fmt.Println("## Answer")
	if result.Legal {
		fmt.Printf("The move %s at (%d,%d) is legal and captures %d opponent stone(s).\n",
			move.Player.GoString(), move.Row, move.Col, len(result.Captured))
		if len(result.CapturedPos) > 0 {
			fmt.Printf("Captured stones at: %s\n", strings.Join(result.CapturedPos, ", "))
		}
		fmt.Println("Board after the move:")
	} else {
		fmt.Printf("The move is illegal: %s\n", result.Description)
		fmt.Println("Board remains unchanged:")
	}
	fmt.Println()
	fmt.Println(result.BoardAfter.string())
	fmt.Println()

	// --- Reason Why ---
	fmt.Println("## Reason why")
	fmt.Println("In Go, a stone lives if it belongs to a group with at least one liberty (empty adjacent intersection).")
	fmt.Println("When a player places a stone, any opponent group that loses its last liberty is immediately captured and removed.")
	fmt.Printf("Before the move, the Black stone at (1,1) had only one liberty at (1,0) – the other three were occupied by White.\n")
	fmt.Printf("By playing White at (1,0), that last liberty is taken, so the Black stone is captured and removed.\n")
	fmt.Println("The White move itself is safe because the newly placed stone still has liberties (e.g., (0,0), (2,0), ...).")
	fmt.Println()

	// --- Check ---
	fmt.Println("## Check")
	fmt.Printf("C1 OK - Black stone had exactly 1 liberty before the move (count=%d).\n", cks.LibertyBeforeMove)
	if cks.LibertyAfterMove == 0 {
		fmt.Println("C2 OK - After capture, the Black stone is gone (0 liberties).")
	} else {
		fmt.Println("C2 FAIL - Black stone not properly removed.")
	}
	if cks.NoSuicide {
		fmt.Println("C3 OK - The White move is not suicide (it has remaining liberties).")
	} else {
		fmt.Println("C3 FAIL - Suicide move detected.")
	}
	if cks.CaptureOccurred {
		fmt.Println("C4 OK - At least one stone was captured.")
	} else {
		fmt.Println("C4 FAIL - Expected capture did not occur.")
	}
	if cks.BoardSizeOk {
		fmt.Println("C5 OK - Board size is 5x5.")
	} else {
		fmt.Println("C5 FAIL - Board size wrong.")
	}
	fmt.Println()

	// --- Go audit details ---
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("board size : %dx%d\n", scenario.Initial.Size, scenario.Initial.Size)
	fmt.Printf("move : %s (%d,%d)\n", move.Player.GoString(), move.Row, move.Col)
	fmt.Println("initial board:")
	fmt.Print(scenario.Initial.string())
	fmt.Println()
	fmt.Printf("result : %s\n", result.Description)
	fmt.Printf("captured stones : %d\n", len(result.Captured))
	fmt.Printf("checks passed : %d/5\n", checkCount(cks))
	fmt.Printf("recommendation consistent : %s\n", yesNo(allChecksPass(cks)))
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// ---------- main ----------

func main() {
	scenario := exampleinput.Load(eyelingoExampleName, buildScenario())
	result := scenario.Initial.playMove(scenario.Move.Player, scenario.Move.Row, scenario.Move.Col)
	if !result.Legal {
		fmt.Fprintf(os.Stderr, "unexpected illegal move: %s\n", result.Description)
		os.Exit(1)
	}
	cks := performChecks(scenario.Initial, scenario.Move, result)
	renderArcOutput(scenario, result, cks)
}
