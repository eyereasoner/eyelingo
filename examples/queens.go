// queens.go
//
// A small, standalone Go implementation of the N-Queens puzzle.
//
// The program counts every way to place N queens on an N×N chessboard so that
// no two queens attack each other. It uses a compact bit-mask backtracking
// search: each recursive level represents one board row, and three masks track
// columns plus the two diagonal directions that are already under attack.
//
// Run with the default 8-Queens puzzle:
//
//	go run queens.go
//
// Or choose the board size and number of example boards to print:
//
//	go run queens.go 14 1
//
// the normalized command-line settings and search counters so the translation is
// easy to inspect, compare, and regression-test.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"math/bits"
	"os"
	"strconv"
)

const seeExampleName = "queens"

// SolveStats records implementation-level counters from the backtracking run.
// These values are not needed to solve the puzzle, but they make the Go output
type SolveStats struct {
	Nodes               uint64
	CandidatePlacements uint64
	DeadEnds            uint64
	SolutionsPrinted    uint64
	MaxDepth            int
	FirstSolution       []int
}

// solveNQueens validates the board size, prepares the initial bit mask, and
// counters collected along the way.
func solveNQueens(n int, maxPrint int) (uint64, SolveStats) {
	if n <= 0 {
		return 0, SolveStats{}
	}

	// The search stores column availability in a uint64. N=64 would require a
	// special full-mask case, but the shifted diagonal masks would overflow past
	// the board edge, so this implementation intentionally caps N at 63.
	if n > 63 {
		fmt.Fprintln(os.Stderr, "This implementation supports N up to 63.")
		os.Exit(1)
	}

	// allColumns has the low N bits set to 1. A 1 bit means the corresponding
	// column is part of the board and may still be considered during search.
	allColumns := (uint64(1) << n) - 1

	board := make([]int, n)
	var count uint64
	stats := SolveStats{}

	search(n, 0, 0, 0, 0, allColumns, board, &count, maxPrint, &stats)

	return count, stats
}

// search places queens one row at a time.
//
// columns marks occupied columns. diagLeft and diagRight mark squares attacked
// by queens on the two diagonal directions for the current row. At each row,
// available is the set of safe columns. The expression position := available &
// -available extracts the lowest set bit, which is the next candidate column.
func search(
	n int,
	row int,
	columns uint64,
	diagLeft uint64,
	diagRight uint64,
	allColumns uint64,
	board []int,
	count *uint64,
	maxPrint int,
	stats *SolveStats,
) {
	stats.Nodes++
	if row > stats.MaxDepth {
		stats.MaxDepth = row
	}

	// Reaching row N means every row has a queen and all attack constraints were
	// respected, so this branch is a complete solution.
	if row == n {
		*count++

		if len(stats.FirstSolution) == 0 {
			stats.FirstSolution = append([]int(nil), board...)
		}

		if *count <= uint64(maxPrint) {
			fmt.Printf("Solution %d:\n", *count)
			printBoard(board)
			fmt.Println()
			stats.SolutionsPrinted++
		}

		return
	}

	// Combine all attacked columns, invert them, and keep only bits inside the
	// board. The remaining 1 bits are safe columns for this row.
	available := allColumns & ^(columns | diagLeft | diagRight)
	if available == 0 {
		stats.DeadEnds++
	}

	for available != 0 {
		// Pick one available column and remove it from this row's candidate mask.
		position := available & -available
		available ^= position
		stats.CandidatePlacements++

		column := bits.TrailingZeros64(position)
		board[row] = column

		// Moving to the next row shifts diagonal attacks outward by one column.
		search(
			n,
			row+1,
			columns|position,
			(diagLeft|position)<<1,
			(diagRight|position)>>1,
			allColumns,
			board,
			count,
			maxPrint,
			stats,
		)
	}
}

// printBoard renders a found solution using Q for queens and . for empty
// squares, then prints the same solution as one-based column positions by row.
func printBoard(board []int) {
	n := len(board)

	for _, queenCol := range board {
		for col := 0; col < n; col++ {
			if col == queenCol {
				fmt.Print("Q ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}

	fmt.Print("As column positions by row: [")
	for i, col := range board {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(col + 1)
	}
	fmt.Println("]")
}

func renderReason(n int) {
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Printf("The solver places one queen per row on a %dx%d board.\n", n, n)
	fmt.Println("At each row it uses bit masks for occupied columns and both diagonal directions to enumerate only safe candidate columns.")
	fmt.Println("Counting continues after the printed solution limit, so the total solution count remains complete.")
}

func columnsUnique(board []int) bool {
	seen := map[int]bool{}
	for _, col := range board {
		if seen[col] {
			return false
		}
		seen[col] = true
	}
	return true
}

func diagonalsSafe(board []int) bool {
	for r1, c1 := range board {
		for r2 := r1 + 1; r2 < len(board); r2++ {
			c2 := board[r2]
			if abs(r1-r2) == abs(c1-c2) {
				return false
			}
		}
	}
	return true
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func statusWord(ok bool) string {
	if ok {
		return "OK"
	}
	return "FAIL"
}

// what was solved and how the search behaved. These lines are diagnostic rather
// than part of the mathematical N-Queens answer.

// formatOneBasedColumns converts the internal zero-based board representation
// into a compact, human-readable one-based column list.
func formatOneBasedColumns(board []int) string {
	result := "["
	for i, col := range board {
		if i > 0 {
			result += ", "
		}
		result += strconv.Itoa(col + 1)
	}
	return result + "]"
}

// translated examples.
func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func main() {
	cfg := exampleinput.Load(seeExampleName, struct {
		N        int
		MaxPrint int
	}{N: 8, MaxPrint: 1})
	n := cfg.N
	maxPrint := cfg.MaxPrint

	// Optional argument 1 selects N, the board width/height and queen count.
	if len(os.Args) >= 2 {
		parsed, err := strconv.Atoi(os.Args[1])
		if err != nil || parsed < 0 {
			fmt.Fprintln(os.Stderr, "Usage: go run queens.go [N] [MAX_PRINT]")
			os.Exit(1)
		}
		n = parsed
	}

	// Optional argument 2 limits how many complete boards are printed. Counting
	// still continues after this limit; only rendering is capped.
	if len(os.Args) >= 3 {
		parsed, err := strconv.Atoi(os.Args[2])
		if err != nil || parsed < 0 {
			fmt.Fprintln(os.Stderr, "MAX_PRINT must be a non-negative integer.")
			os.Exit(1)
		}
		maxPrint = parsed
	}

	fmt.Println("# 8-Queens")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("Solving %d-Queens...\n", n)
	fmt.Printf("Printing at most %d solution(s).\n", maxPrint)
	fmt.Println()

	count, _ := solveNQueens(n, maxPrint)

	fmt.Printf("Total solutions for %d-Queens: %d\n", n, count)
	renderReason(n)
}
