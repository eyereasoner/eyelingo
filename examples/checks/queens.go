package main

import (
	"math"
	"sort"
	"strings"
)

func checkQueens(ctx *Context) []Check {
	d := ctx.M()
	n := integer(d["N"])
	maxp := integer(d["MaxPrint"])
	cols := []int{}
	if m := reFind(ctx.Answer, `As column positions by row: \[([^\]]+)\]`); m != nil {
		for _, p := range strings.Split(m[0], ",") {
			cols = append(cols, parseInt(strings.TrimSpace(p)))
		}
	}
	total := grabInt(ctx.Answer, `Total solutions for \d+-Queens: (\d+)`)
	zero := []int{}
	for _, c := range cols {
		zero = append(zero, c-1)
	}
	diagOK := true
	for r1 := 0; r1 < len(zero); r1++ {
		for r2 := r1 + 1; r2 < len(zero); r2++ {
			if int(math.Abs(float64(zero[r1]-zero[r2]))) == r2-r1 {
				diagOK = false
			}
		}
	}
	rows := 0
	queens := 0
	for _, line := range strings.Split(ctx.Answer, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			ok := true
			q := 0
			for _, f := range fields {
				if f != "Q" && f != "." {
					ok = false
				}
				if f == "Q" {
					q++
				}
			}
			if ok && q > 0 {
				rows++
				queens += q
			}
		}
	}
	indep := countQueens(n)
	sorted := append([]int(nil), cols...)
	sort.Ints(sorted)
	want := []int{}
	for i := 1; i <= n; i++ {
		want = append(want, i)
	}
	inRange := true
	for _, c := range cols {
		inRange = inRange && c >= 1 && c <= n
	}
	return []Check{{"the checker loaded the normalized 8-Queens input", n == 8 && maxp == 1}, {"the printed solution gives one column for each row", len(cols) == n}, {"all printed columns are within the board", inRange}, {"the printed solution uses each column exactly once", intsEq(sorted, want)}, {"no pair of printed queens shares a diagonal", diagOK}, {"the rendered board contains exactly eight rows and eight queens", rows == n && queens == n}, {"an independent bit-mask search counts 92 total solutions", indep == 92}, {"the reported total matches the independent solution count", total == indep}}
}
