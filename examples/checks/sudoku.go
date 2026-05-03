package main

func checkSudoku(ctx *Context) []Check {
	d := ctx.M()
	puzzle := str(d["Puzzle"])
	grid := parseSudokuGrid(ctx.Answer)
	givens := [][3]int{}
	for i, ch := range puzzle {
		if ch != '0' {
			givens = append(givens, [3]int{i / 9, i % 9, int(ch - '0')})
		}
	}
	preserved := len(grid) == 9
	full := len(grid) == 9
	for _, g := range givens {
		preserved = preserved && grid[g[0]][g[1]] == g[2]
	}
	for _, row := range grid {
		for _, c := range row {
			full = full && c >= 1 && c <= 9
		}
	}
	rowsOK := len(grid) == 9
	colsOK := len(grid) == 9
	boxesOK := len(grid) == 9
	for i := 0; i < 9 && len(grid) == 9; i++ {
		rowsOK = rowsOK && digitSet(grid[i])
		col := []int{}
		for r := 0; r < 9; r++ {
			col = append(col, grid[r][i])
		}
		colsOK = colsOK && digitSet(col)
	}
	if len(grid) == 9 {
		for br := 0; br < 9; br += 3 {
			for bc := 0; bc < 9; bc += 3 {
				vals := []int{}
				for r := br; r < br+3; r++ {
					for c := bc; c < bc+3; c++ {
						vals = append(vals, grid[r][c])
					}
				}
				boxesOK = boxesOK && digitSet(vals)
			}
		}
	}
	legal := len(grid) == 9
	for r := 0; r < 9 && len(grid) == 9; r++ {
		for c := 0; c < 9; c++ {
			legal = legal && sudokuLegal(grid, r, c, grid[r][c])
		}
	}
	expected := [][]int{{1, 6, 2, 8, 5, 7, 4, 9, 3}, {5, 3, 4, 1, 2, 9, 6, 7, 8}, {7, 8, 9, 6, 4, 3, 5, 2, 1}, {4, 7, 5, 3, 1, 2, 9, 8, 6}, {9, 1, 3, 5, 8, 6, 7, 4, 2}, {6, 2, 8, 7, 9, 4, 1, 3, 5}, {3, 5, 6, 4, 7, 8, 2, 1, 9}, {2, 4, 1, 9, 3, 5, 8, 6, 7}, {8, 9, 7, 2, 6, 1, 3, 5, 4}}
	eq := gridEq(grid, expected)
	return []Check{{"the input puzzle has 81 cells and exactly 23 given clues", len(puzzle) == 81 && len(givens) == 23}, {"the completed grid is parsed as nine rows of nine digits", len(grid) == 9 && allRowsLen(grid, 9)}, {"every original clue is preserved at the same row and column", preserved}, {"the final grid contains only digits 1 through 9", full}, {"each completed row is a permutation of 1 through 9", rowsOK}, {"each completed column is a permutation of 1 through 9", colsOK}, {"each completed 3×3 box is a permutation of 1 through 9", boxesOK}, {"every filled cell is legal against its row, column, and box peers", legal}, {"the completed grid matches a separately embedded expected solution fixture", eq}}
}
