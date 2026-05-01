// fibonacci.go
//
// A self-contained Go translation of examples/fibonacci.n3 from the Eyeling
// example suite, extended to handle very large indices (up to F(10000)).
// It uses arbitrary‑precision integers to compute exact Fibonacci numbers.
//
// This is intentionally not a full N3 reasoner – it is a concrete scenario
// that mirrors the structure of the original N3 example.
//
// Run:
//
//     go run fibonacci.go
//
// The program has no third‑party dependencies.

package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math/big"
	"os"
	"sort" // needed for stable output ordering
)

const eyelingoExampleName = "fibonacci"

// ---------- expected values ----------
// These are the exact decimal representations of F(n) for selected n,
// taken from the well‑known Fibonacci sequence.  They are used as

var expected = map[int]string{
	0:     "0",
	1:     "1",
	10:    "55",
	100:   "354224848179261915075",
	1000:  "43466557686937456435688527675040625802564660517371780402481729089536555417949051890403879840079255169295922593080322634775209689623239873322471161642996440906533187938298969649928516003704476137795166849228875",
	10000: "33644764876431783266621612005107543310302148460680063906564769974680081442166662368155595513633734025582065332680836159373734790483865268263040892463056431887354544369559827491606602099884183933864652731300088830269235673613135117579297437854413752130520504347701602264758318906527890855154366159582987279682987510631200575428783453215515103870818298969791613127856265033195487140214287532698187962046936097879900350962302291026368131493195275630227837628441540360584402572114334961180023091208287046088923962328835461505776583271252546093591128203925285393434620904245248929403901706233888991085841065183173360437470737908552631764325733993712871937587746897479926305837065742830161637408969178426378624212835258112820516370298089332099905707920064367426202389783111470054074998459250360633560933883831923386783056136435351892133279732908133732642652633989763922723407882928177953580570993691049175470808931841056146322338217465637321248226383092103297701648054726243842374862411453093812206564914032751086643394517512161526545361333111314042436854805106765843493523836959653428071768775328348234345557366719731392746273629108210679280784718035329131176778924659089938635459327894523777674406192240337638674004021330343297496902028328145933418826817683893072003634795623117103101291953169794607632737589253530772552375943788434504067715555779056450443016640119462580972216729758615026968443146952034614932291105970676243268515992834709891284706740862008587135016260312071903172086094081298321581077282076353186624611278245537208532365305775956430072517744315051539600905168603220349163222640885248852433158051534849622434848299380905070483482449327453732624567755879089187190803662058009594743150052402532709746995318770724376825907419939632265984147498193609285223945039707165443156421328157688908058783183404917434556270520223564846495196112460268313970975069382648706613264507665074611512677522748621598642530711298441182622661057163515069260029861704945425047491378115154139941550671256271197133252763631939606902895650288268608362241082050562430701794976171121233066073310059947366875",
}

// ---------- inference ----------

// computeFibonacci returns a slice of *big.Int for F(0) .. F(target).
func computeFibonacci(target int) []*big.Int {
	if target < 0 {
		return nil
	}
	seq := make([]*big.Int, target+1)

	// Base cases using big.Int
	zero := new(big.Int)
	one := new(big.Int).SetInt64(1)
	seq[0] = zero
	if target >= 1 {
		seq[1] = one
	}

	// Recurrence: F(n) = F(n-1) + F(n-2)
	for n := 2; n <= target; n++ {
		sum := new(big.Int).Add(seq[n-1], seq[n-2])
		seq[n] = sum
	}
	return seq
}

// ---------- helper: sorted keys ----------

func sortedKeys(m map[int]string) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

type Checks struct {
	BaseCasesCorrect   bool
	RecurrenceHolds    bool
	AllTargetsMatch    bool
	NonNegative        bool
	MonotonicFromStart bool
}

func performChecks(seq []*big.Int, expected map[int]string) Checks {
	c := Checks{}

	// C1: base cases
	zero := big.NewInt(0)
	one := big.NewInt(1)
	c.BaseCasesCorrect = (seq[0].Cmp(zero) == 0) && (seq[1].Cmp(one) == 0)

	// C2: recurrence for n>=2
	c.RecurrenceHolds = true
	for n := 2; n < len(seq); n++ {
		want := new(big.Int).Add(seq[n-1], seq[n-2])
		if seq[n].Cmp(want) != 0 {
			c.RecurrenceHolds = false
			break
		}
	}

	// C3: all given test targets match
	c.AllTargetsMatch = true
	for idx, expStr := range expected {
		if idx >= len(seq) {
			c.AllTargetsMatch = false
			break
		}
		exp, ok := new(big.Int).SetString(expStr, 10)
		if !ok {
			c.AllTargetsMatch = false
			break
		}
		if seq[idx].Cmp(exp) != 0 {
			c.AllTargetsMatch = false
			break
		}
	}

	// C4: non‑negative
	c.NonNegative = true
	for _, v := range seq {
		if v.Sign() < 0 {
			c.NonNegative = false
			break
		}
	}

	// C5: monotonic from F(2) onward
	c.MonotonicFromStart = true
	for n := 2; n < len(seq); n++ {
		if seq[n].Cmp(seq[n-1]) < 0 {
			c.MonotonicFromStart = false
			break
		}
	}
	return c
}

func allChecksPass(c Checks) bool {
	return c.BaseCasesCorrect && c.RecurrenceHolds && c.AllTargetsMatch && c.NonNegative && c.MonotonicFromStart
}

func checkCount(c Checks) int {
	count := 0
	if c.BaseCasesCorrect {
		count++
	}
	if c.RecurrenceHolds {
		count++
	}
	if c.AllTargetsMatch {
		count++
	}
	if c.NonNegative {
		count++
	}
	if c.MonotonicFromStart {
		count++
	}
	return count
}

// ---------- rendering ----------

// abbreviated returns a short description of a huge number.
func abbreviated(num *big.Int) string {
	s := num.String()
	if len(s) <= 60 {
		return s
	}
	return s[:30] + " ... " + s[len(s)-30:]
}

func renderArcOutput(target int, seq []*big.Int, cks Checks) {
	fmt.Println("# Fibonacci Example (Big)")
	fmt.Println()

	// --- Answer ---
	fmt.Println("## Answer")
	fmt.Printf("The Fibonacci number for index %d is:\n", target)
	fmt.Println(seq[target].String())
	fmt.Println()

	// --- Reason Why ---
	fmt.Println("## Reason why")
	fmt.Println("The Fibonacci sequence is defined by F(0)=0, F(1)=1,")
	fmt.Println("and F(n)=F(n-1)+F(n-2) for n>=2.")
	fmt.Println("Arbitrary‑precision arithmetic (math/big) is used to")
	fmt.Println("compute the exact value without overflow, even for")
	fmt.Println("indices as large as 10000.")
	fmt.Println()

	return
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// ---------- main ----------

func main() {
	inputExpected := exampleinput.Load(eyelingoExampleName, expected)
	expected = inputExpected

	// Use the maximum index from the expected set to compute everything needed.
	maxIdx := 0
	for k := range expected {
		if k > maxIdx {
			maxIdx = k
		}
	}
	if maxIdx < 0 {
		fmt.Fprintln(os.Stderr, "no targets provided")
		os.Exit(1)
	}

	// Compute the whole sequence in one pass.
	seq := computeFibonacci(maxIdx)
	if len(seq) != maxIdx+1 {
		fmt.Fprintf(os.Stderr, "internal error: expected %d elements, got %d\n", maxIdx+1, len(seq))
		os.Exit(1)
	}

	cks := performChecks(seq, expected)
	renderArcOutput(maxIdx, seq, cks)
}
