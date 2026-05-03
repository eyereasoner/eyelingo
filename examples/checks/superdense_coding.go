package main

import (
	"fmt"
	"sort"
	"strings"
)

func checkSuperdense(ctx *Context) []Check {
	d := ctx.M()
	states := map[string]map[mobPair]bool{}
	for name, rows := range asMap(d["States"]) {
		states[name] = rel(asSlice(rows))
	}
	prim := map[string]map[mobPair]bool{}
	for name, rows := range asMap(d["Primitive"]) {
		prim[name] = rel(asSlice(rows))
	}
	comp := map[string]map[mobPair]bool{"kg": composeRel(prim["g"], prim["k"]), "gk": composeRel(prim["k"], prim["g"])}
	lookup := func(name string) map[mobPair]bool {
		if r, ok := prim[name]; ok {
			return r
		}
		return comp[name]
	}
	alice := map[int]string{}
	for k, v := range asMap(d["AliceOps"]) {
		alice[parseInt(k)] = str(v)
	}
	bob := map[int]string{0: "gk", 1: "k", 2: "g", 3: "id"}
	counts := map[[2]int]int{}
	candidates := 0
	encoded := map[int]map[mobPair]bool{}
	for msg := 0; msg < 4; msg++ {
		encoded[msg] = map[mobPair]bool{}
		for sh := range states["R"] {
			for mv := range lookup(alice[msg]) {
				if mv.a != sh.a {
					continue
				}
				encoded[msg][mobPair{mv.b, sh.b}] = true
				for dec := 0; dec < 4; dec++ {
					if lookup(bob[dec])[mobPair{mv.b, sh.b}] {
						counts[[2]int{msg, dec}]++
						candidates++
					}
				}
			}
		}
	}
	survivors := map[int]int{}
	diag := []int{}
	offEven := true
	for i := 0; i < 4; i++ {
		diag = append(diag, counts[[2]int{i, i}])
		for j := 0; j < 4; j++ {
			if i != j && counts[[2]int{i, j}]%2 != 0 {
				offEven = false
			}
			if counts[[2]int{i, j}]%2 == 1 {
				survivors[i] = j
			}
		}
	}
	reported := map[int]int{}
	for _, m := range reAll(ctx.Answer, `(?m)^\s*(\d+) dqc:superdense-coding (\d+)`) {
		reported[parseInt(m[1])] = parseInt(m[2])
	}
	encodedSupports := map[string]bool{}
	for _, v := range encoded {
		parts := []string{}
		for p := range v {
			parts = append(parts, fmt.Sprintf("%v/%v", p.a, p.b))
		}
		sort.Strings(parts)
		encodedSupports[strings.Join(parts, ",")] = true
	}
	return []Check{{"shared entanglement R contains exactly |0,0) and |1,1)", relEq(states["R"], map[mobPair]bool{{false, false}: true, {true, true}: true})}, {"composition KG is recomputed by composing G then K", relEq(comp["kg"], map[mobPair]bool{{false, false}: true, {false, true}: true, {true, false}: true})}, {"composition GK is recomputed by composing K then G", relEq(comp["gk"], map[mobPair]bool{{false, true}: true, {true, false}: true, {true, true}: true})}, {"the raw superdense rule creates 24 candidate derivations", candidates == 24}, {"GF(2) cancellation leaves odd diagonal and even off-diagonal counts", intsEq(diag, []int{1, 1, 1, 1}) && offEven}, {"reported surviving decoded messages match the parity survivors", mapIntEq(reported, survivors) && mapIntEq(survivors, map[int]int{0: 0, 1: 1, 2: 2, 3: 3})}, {"the four Alice operations produce distinct encoded supports", len(encodedSupports) == 4}, {"the JSON relation facts match the primitive teaching-model relations", relEq(prim["id"], states["R"]) && relEq(prim["g"], states["S"]) && relEq(prim["k"], states["U"])}}
}
