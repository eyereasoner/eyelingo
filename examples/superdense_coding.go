// superdense_coding.go
//
// A self-contained Go translation of superdense-coding.n3 from the Eyeling
// examples.
//
// The N3 example models superdense coding with modal bits, or "mobits".
// Think of a mobit as the small teaching-model analogue of a quantum bit. Alice
// and Bob share the entangled state |R). Alice encodes one of four two-bit
// messages by applying one relation to her half. Bob then applies a joint test
// to decode it. Because the model uses GF(2), meaning arithmetic modulo 2,
// duplicate derivations cancel: an answer is kept only when it appears an odd
// number of times.
//
// This Go version keeps that rule structure visible: relation composition,
// superdense candidate generation, parity cancellation, and audit checks are all
// explicit.
//
// Run:
//
//    go run superdense_coding.go
//
// The program has no third-party dependencies.
package main

import (
    "fmt"
    "os"
    "runtime"
    "sort"
    "strings"
)

type Bit bool

type Pair struct {
    A Bit
    B Bit
}

type Relation map[Pair]bool

type Candidate struct {
    Message       int
    Decoded       int
    SharedAlice   Bit
    AliceOutbound Bit
    BobBit        Bit
    Reason        string
}

type Result struct {
    States          map[string]Relation
    Primitive       map[string]Relation
    Composed        map[string]Relation
    AliceOps        map[int]string
    BobTests        map[int]string
    Candidates      []Candidate
    Counts          map[[2]int]int
    Survivors       map[int]int
    EncodedSupports map[int]Relation
    Checks          []Check
}

type Check struct {
    Label string
    OK    bool
    Text  string
}

func main() {
    result := derive()
    printAnswer(result)
    printReason(result)
    printChecks(result.Checks)
    printAudit(result)
    if !allChecksOK(result.Checks) {
        os.Exit(1)
    }
}

func derive() Result {
    states := map[string]Relation{
        "R": relation(pair(false, false), pair(true, true)),
        "S": relation(pair(false, true), pair(true, false)),
        "U": relation(pair(false, false), pair(true, false), pair(true, true)),
        "V": relation(pair(false, false), pair(false, true), pair(true, false)),
    }

    primitive := map[string]Relation{
        "id": relation(pair(false, false), pair(true, true)),
        "g":  relation(pair(false, true), pair(true, false)),
        "k":  relation(pair(false, false), pair(true, false), pair(true, true)),
    }
    composed := map[string]Relation{
        "kg": compose(primitive["g"], primitive["k"]), // {?X g ?Z. ?Z k ?Y}
        "gk": compose(primitive["k"], primitive["g"]), // {?X k ?Z. ?Z g ?Y}
    }

    aliceOps := map[int]string{0: "id", 1: "g", 2: "k", 3: "kg"}
    bobTests := map[int]string{0: "gk", 1: "k", 2: "g", 3: "id"}

    candidates := make([]Candidate, 0)
    counts := make(map[[2]int]int)
    encoded := make(map[int]Relation)

    for msg := 0; msg < 4; msg++ {
        encoded[msg] = make(Relation)
        aliceRel := lookup(aliceOps[msg], primitive, composed)
        for _, shared := range sortedPairs(states["R"]) {
            for _, move := range sortedPairs(aliceRel) {
                if move.A != shared.A {
                    continue
                }
                aliceOutbound := move.B
                encoded[msg][pair(aliceOutbound, shared.B)] = true

                for decoded := 0; decoded < 4; decoded++ {
                    bobRel := lookup(bobTests[decoded], primitive, composed)
                    if !bobRel[pair(aliceOutbound, shared.B)] {
                        continue
                    }
                    counts[[2]int{msg, decoded}]++
                    candidates = append(candidates, Candidate{
                        Message:       msg,
                        Decoded:       decoded,
                        SharedAlice:   shared.A,
                        AliceOutbound: aliceOutbound,
                        BobBit:        shared.B,
                        Reason: fmt.Sprintf("R has %s; Alice %d uses %s to emit %s; Bob test %d/%s accepts (%s,%s)",
                            formatPair(shared), msg, aliceOps[msg], bit(aliceOutbound), decoded, bobTests[decoded], bit(aliceOutbound), bit(shared.B)),
                    })
                }
            }
        }
    }

    survivors := make(map[int]int)
    for msg := 0; msg < 4; msg++ {
        for decoded := 0; decoded < 4; decoded++ {
            if counts[[2]int{msg, decoded}]%2 == 1 {
                survivors[msg] = decoded
            }
        }
    }

    result := Result{
        States:          states,
        Primitive:       primitive,
        Composed:        composed,
        AliceOps:        aliceOps,
        BobTests:        bobTests,
        Candidates:      candidates,
        Counts:          counts,
        Survivors:       survivors,
        EncodedSupports: encoded,
    }
    result.Checks = buildChecks(result)
    return result
}

func buildChecks(r Result) []Check {
    checks := []Check{
        {
            Label: "shared entanglement",
            OK:    sameRelation(r.States["R"], relation(pair(false, false), pair(true, true))),
            Text:  "|R) contains exactly |0,0) and |1,1)",
        },
        {
            Label: "composition KG",
            OK:    sameRelation(r.Composed["kg"], relation(pair(false, false), pair(false, true), pair(true, false))),
            Text:  "KG is obtained by composing G then K, exactly as in the N3 rule",
        },
        {
            Label: "composition GK",
            OK:    sameRelation(r.Composed["gk"], relation(pair(false, true), pair(true, false), pair(true, true))),
            Text:  "GK is obtained by composing K then G, exactly as in the N3 rule",
        },
        {
            Label: "candidate multiplicity",
            OK:    len(r.Candidates) == 24,
            Text:  "the raw superdense rule creates 24 candidate derivations before parity cancellation",
        },
        {
            Label: "GF(2) cancellation",
            OK:    allDiagonalOnly(r.Counts),
            Text:  "off-diagonal answers occur zero or two times and cancel; diagonal answers occur once",
        },
        {
            Label: "decoded messages",
            OK:    r.Survivors[0] == 0 && r.Survivors[1] == 1 && r.Survivors[2] == 2 && r.Survivors[3] == 3 && len(r.Survivors) == 4,
            Text:  "after odd-parity filtering, Bob recovers Alice's original two-bit message",
        },
        {
            Label: "encoded supports distinct",
            OK:    distinctSupports(r.EncodedSupports),
            Text:  "the four Alice operations produce four different supports over the two mobits",
        },
    }
    return checks
}

func compose(first, second Relation) Relation {
    out := make(Relation)
    for p := range first {
        for q := range second {
            if p.B == q.A {
                out[pair(p.A, q.B)] = true
            }
        }
    }
    return out
}

func lookup(name string, primitive, composed map[string]Relation) Relation {
    if rel, ok := primitive[name]; ok {
        return rel
    }
    return composed[name]
}

func printAnswer(r Result) {
    fmt.Println("=== Answer ===")
    fmt.Println("Superdense-coding facts that survive GF(2) parity cancellation:")
    for msg := 0; msg < 4; msg++ {
        fmt.Printf("  %d dqc:superdense-coding %d\n", msg, r.Survivors[msg])
    }
    fmt.Println()
}

func printReason(r Result) {
    fmt.Println("=== Reason Why ===")
    fmt.Println("Alice and Bob start with |R) = |0,0) + |1,1).")
    fmt.Println("Alice chooses one relation for the first mobit; Bob applies one joint test.")
    fmt.Println("The N3 example keeps only answers with odd derivation count, so duplicate")
    fmt.Println("modal paths cancel just like addition over GF(2).")
    fmt.Println()
    fmt.Println("Alice operations:")
    for msg := 0; msg < 4; msg++ {
        fmt.Printf("  message %d -> %-2s -> encoded support %s\n", msg, strings.ToUpper(r.AliceOps[msg]), formatRelation(r.EncodedSupports[msg]))
    }
    fmt.Println()
    fmt.Println("Raw candidate counts before parity filtering:")
    for msg := 0; msg < 4; msg++ {
        pieces := make([]string, 0, 4)
        for decoded := 0; decoded < 4; decoded++ {
            pieces = append(pieces, fmt.Sprintf("%d:%d", decoded, r.Counts[[2]int{msg, decoded}]))
        }
        fmt.Printf("  encoded %d -> decoded counts {%s}\n", msg, strings.Join(pieces, ", "))
    }
    fmt.Println()
    fmt.Println("Surviving explanations:")
    for msg := 0; msg < 4; msg++ {
        fmt.Printf("  %d -> %d because count=%d is odd; all competing counts are even.\n", msg, r.Survivors[msg], r.Counts[[2]int{msg, r.Survivors[msg]}])
    }
    fmt.Println()
}

func printChecks(checks []Check) {
    fmt.Println("=== Check ===")
    for _, check := range checks {
        status := "FAIL"
        if check.OK {
            status = "OK"
        }
        fmt.Printf("  %s %-24s %s\n", status, check.Label+":", check.Text)
    }
    fmt.Println()
}

func printAudit(r Result) {
    fmt.Println("=== Go audit details ===")
    fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
    fmt.Printf("mobit values: %s, %s\n", bit(false), bit(true))
    fmt.Printf("states      : R=%s S=%s U=%s V=%s\n", formatRelation(r.States["R"]), formatRelation(r.States["S"]), formatRelation(r.States["U"]), formatRelation(r.States["V"]))
    fmt.Printf("relations   : id=%s G=%s K=%s KG=%s GK=%s\n", formatRelation(r.Primitive["id"]), formatRelation(r.Primitive["g"]), formatRelation(r.Primitive["k"]), formatRelation(r.Composed["kg"]), formatRelation(r.Composed["gk"]))
    fmt.Printf("raw candidates before cancellation: %d\n", len(r.Candidates))
    fmt.Printf("surviving facts after odd-parity filter: %d\n", len(r.Survivors))
}

func relation(pairs ...Pair) Relation {
    r := make(Relation)
    for _, p := range pairs {
        r[p] = true
    }
    return r
}

func pair(a, b Bit) Pair { return Pair{A: a, B: b} }

func sortedPairs(r Relation) []Pair {
    out := make([]Pair, 0, len(r))
    for p := range r {
        out = append(out, p)
    }
    sort.Slice(out, func(i, j int) bool {
        if out[i].A != out[j].A {
            return !bool(out[i].A) && bool(out[j].A)
        }
        return !bool(out[i].B) && bool(out[j].B)
    })
    return out
}

func formatRelation(r Relation) string {
    pairs := sortedPairs(r)
    parts := make([]string, 0, len(pairs))
    for _, p := range pairs {
        parts = append(parts, formatPair(p))
    }
    return "{" + strings.Join(parts, ", ") + "}"
}

func formatPair(p Pair) string {
    return fmt.Sprintf("(%s,%s)", bit(p.A), bit(p.B))
}

func bit(b Bit) string {
    if b {
        return "1"
    }
    return "0"
}

func sameRelation(a, b Relation) bool {
    if len(a) != len(b) {
        return false
    }
    for p := range a {
        if !b[p] {
            return false
        }
    }
    return true
}

func allDiagonalOnly(counts map[[2]int]int) bool {
    for msg := 0; msg < 4; msg++ {
        for decoded := 0; decoded < 4; decoded++ {
            count := counts[[2]int{msg, decoded}]
            if msg == decoded {
                if count%2 != 1 {
                    return false
                }
            } else if count%2 != 0 {
                return false
            }
        }
    }
    return true
}

func distinctSupports(supports map[int]Relation) bool {
    seen := make(map[string]bool)
    for msg := 0; msg < 4; msg++ {
        key := formatRelation(supports[msg])
        if seen[key] {
            return false
        }
        seen[key] = true
    }
    return true
}

func allChecksOK(checks []Check) bool {
    for _, check := range checks {
        if !check.OK {
            return false
        }
    }
    return true
}
