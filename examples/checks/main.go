package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type anymap = map[string]any

type Context struct {
	Root   string
	Name   string
	Prefix string
	Answer string
	Reason string
	data   any
	loaded bool
}

type Check struct {
	Desc string
	OK   bool
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: see-checks EXAMPLE PREFIX_FILE")
		os.Exit(2)
	}
	root := os.Getenv("SEE_ROOT")
	if root == "" {
		wd, _ := os.Getwd()
		if filepath.Base(wd) == "checks" && filepath.Base(filepath.Dir(wd)) == "examples" {
			root = filepath.Dir(filepath.Dir(wd))
		} else if filepath.Base(wd) == "checks" {
			root = filepath.Dir(wd)
		} else {
			root = wd
		}
	}
	prefixBytes, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	prefix := normalizePrefix(string(prefixBytes))
	ctx := &Context{Root: root, Name: os.Args[1], Prefix: prefix, Answer: section(prefix, "Answer"), Reason: section(prefix, "Reason")}
	checks, err := run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("## Check")
	all := true
	for i, c := range checks {
		status := "FAIL"
		if c.OK {
			status = "OK"
		}
		fmt.Printf("C%d %s - %s\n", i+1, status, c.Desc)
		all = all && c.OK
	}
	if !all {
		os.Exit(1)
	}
}

func normalizePrefix(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\r\n \t"), "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t\r")
	}
	return strings.Join(lines, "\n") + "\n"
}

func section(prefix, heading string) string {
	marker := "## " + heading
	idx := strings.Index(prefix, marker)
	if idx < 0 {
		return ""
	}
	tail := prefix[idx+len(marker):]
	if j := strings.Index(tail, "\n## "); j >= 0 {
		tail = tail[:j]
	}
	return strings.TrimSpace(tail)
}

func (c *Context) Load() any {
	if c.loaded {
		return c.data
	}
	path := filepath.Join(c.Root, "examples", "input", c.Name+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if err := dec.Decode(&c.data); err != nil {
		panic(err)
	}
	c.loaded = true
	return c.data
}

func (c *Context) M() anymap { return asMap(c.Load()) }
func (c *Context) A() []any  { return asSlice(c.Load()) }

func run(ctx *Context) ([]Check, error) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	switch ctx.Name {
	case "ackermann":
		return checkAckermann(ctx), nil
	case "alarm_bit_interoperability":
		return checkAlarmBit(ctx), nil
	case "allen_interval_calculus":
		return checkAllen(ctx), nil
	case "auroracare":
		return checkAuroraCare(ctx), nil
	case "barley_seed_lineage":
		return checkBarley(ctx), nil
	case "bayes_diagnosis":
		return checkBayesDiagnosis(ctx), nil
	case "bayes_therapy":
		return checkBayesTherapy(ctx), nil
	case "bmi":
		return checkBMI(ctx), nil
	case "calidor":
		return checkCalidor(ctx), nil
	case "complex_matrix_stability":
		return checkComplexMatrix(ctx), nil
	case "complex_numbers":
		return checkComplexNumbers(ctx), nil
	case "control_system":
		return checkControlSystem(ctx), nil
	case "deep_taxonomy_100000":
		return checkDeepTaxonomy(ctx), nil
	case "delfour":
		return checkDelfour(ctx), nil
	case "digital_product_passport":
		return checkDPP(ctx), nil
	case "dijkstra_risk_path":
		return checkDijkstra(ctx), nil
	case "dining_philosophers":
		return checkDining(ctx), nil
	case "docking_abort_token":
		return checkDocking(ctx), nil
	case "doctor_advice_work_conflict":
		return checkDoctor(ctx), nil
	case "drone_corridor_planner":
		return checkDrone(ctx), nil
	case "ebike_motor_thermal_envelope":
		return checkEbike(ctx), nil
	case "eco_route_insight":
		return checkEcoRoute(ctx), nil
	case "euler_identity_certificate":
		return checkEuler(ctx), nil
	case "ev_roundtrip_planner":
		return checkEVRoundtrip(ctx), nil
	case "fft32_numeric":
		return checkFFT32(ctx), nil
	case "fft8_numeric":
		return checkFFT8(ctx), nil
	case "fibonacci":
		return checkFibonacci(ctx), nil
	case "flandor":
		return checkFlandor(ctx), nil
	case "fundamental_theorem_arithmetic":
		return checkFTA(ctx), nil
	case "genetic_knapsack_selection":
		return checkGenetic(ctx), nil
	case "goldbach_1000":
		return checkGoldbach(ctx), nil
	case "gps":
		return checkGPS(ctx), nil
	case "gray_code_counter":
		return checkGray(ctx), nil
	case "harbor_smr":
		return checkHarbor(ctx), nil
	case "high_trust_bloom_envelope":
		return checkBloom(ctx), nil
	case "isolation_breach_token":
		return checkIsolation(ctx), nil
	case "kaprekar_6174":
		return checkKaprekar(ctx), nil
	case "lldm":
		return checkLLDM(ctx), nil
	case "odrl_dpv_risk_ranked":
		return checkODRL(ctx), nil
	case "parcellocker":
		return checkParcel(ctx), nil
	case "path_discovery":
		return checkPathDiscovery(ctx), nil
	case "photosynthetic_exciton_transfer":
		return checkPhoto(ctx), nil
	case "queens":
		return checkQueens(ctx), nil
	case "rc_discharge_envelope":
		return checkRC(ctx), nil
	case "school_placement_audit":
		return checkSchool(ctx), nil
	case "sudoku":
		return checkSudoku(ctx), nil
	case "superdense_coding":
		return checkSuperdense(ctx), nil
	case "wind_turbine":
		return checkWind(ctx), nil
	default:
		return nil, fmt.Errorf("no Go check registered for %s", ctx.Name)
	}
}

func asMap(v any) anymap {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	panic(fmt.Sprintf("not map: %T", v))
}
func asSlice(v any) []any {
	if a, ok := v.([]any); ok {
		return a
	}
	panic(fmt.Sprintf("not slice: %T", v))
}
func mget(m anymap, key string) any {
	v, ok := m[key]
	if !ok {
		return nil
	}
	return v
}
func str(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case json.Number:
		return x.String()
	case nil:
		return ""
	default:
		return fmt.Sprint(x)
	}
}
func num(v any) float64 {
	switch x := v.(type) {
	case json.Number:
		f, _ := x.Float64()
		return f
	case float64:
		return x
	case int:
		return float64(x)
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	case bool:
		if x {
			return 1
		}
		return 0
	default:
		return math.NaN()
	}
}
func integer(v any) int { return int(math.Round(num(v))) }
func boolean(v any) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return str(v) == "true"
}
func sliceAny(v any) []any {
	if v == nil {
		return nil
	}
	return asSlice(v)
}
func sarr(v any) []string {
	if v == nil {
		return nil
	}
	a := asSlice(v)
	out := make([]string, len(a))
	for i, x := range a {
		out[i] = str(x)
	}
	return out
}
func farr(v any) []float64 {
	a := asSlice(v)
	out := make([]float64, len(a))
	for i, x := range a {
		out[i] = num(x)
	}
	return out
}
func maps(v any) []anymap {
	if v == nil {
		return nil
	}
	a := asSlice(v)
	out := make([]anymap, len(a))
	for i, x := range a {
		out[i] = asMap(x)
	}
	return out
}
func contains(s, sub string) bool { return strings.Contains(s, sub) }
func reFind(text, pat string) []string {
	r := regexp.MustCompile(pat)
	m := r.FindStringSubmatch(text)
	if m == nil {
		return nil
	}
	return m[1:]
}
func reAll(text, pat string) [][]string {
	r := regexp.MustCompile(pat)
	return r.FindAllStringSubmatch(text, -1)
}
func parseFloat(s string) float64  { f, _ := strconv.ParseFloat(s, 64); return f }
func parseInt(s string) int        { i, _ := strconv.Atoi(s); return i }
func close(a, b, tol float64) bool { return math.Abs(a-b) <= tol }
func setOf(xs []string) map[string]bool {
	m := map[string]bool{}
	for _, x := range xs {
		m[x] = true
	}
	return m
}
func setContainsAll(m map[string]bool, xs []string) bool {
	for _, x := range xs {
		if !m[x] {
			return false
		}
	}
	return true
}
func setEq(a, b []string) bool {
	ma, mb := setOf(a), setOf(b)
	if len(ma) != len(mb) {
		return false
	}
	for k := range ma {
		if !mb[k] {
			return false
		}
	}
	return true
}
func sortedStrings(xs []string) []string {
	ys := append([]string(nil), xs...)
	sort.Strings(ys)
	return ys
}
func sumFloat(xs []float64) float64 {
	s := 0.0
	for _, x := range xs {
		s += x
	}
	return s
}
func sha256Hex(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func hmacSHA256Hex(secret, msg string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}
func parseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05", value)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04", value)
		}
	}
	if err != nil {
		panic(err)
	}
	return t
}
func jsonStable(v any) string { b, _ := json.Marshal(v); return string(b) }
func answerField(answer, label string) string {
	m := reFind(answer, regexp.QuoteMeta(label)+`\s*:\s*([^\n]+)`)
	if m == nil {
		return ""
	}
	return strings.TrimSpace(m[0])
}
func parseIntField(answer, label string) (int, bool) {
	m := reFind(answer, regexp.QuoteMeta(label)+`\s*:\s*(\d+)`)
	if m == nil {
		return 0, false
	}
	return parseInt(m[0]), true
}
func parseFloatField(answer, label string) (float64, bool) {
	m := reFind(answer, regexp.QuoteMeta(label)+`\s*:\s*(-?\d+(?:\.\d+)?)`)
	if m == nil {
		return 0, false
	}
	return parseFloat(m[0]), true
}

func checkFileExists(root, rel string) bool {
	_, err := os.Stat(filepath.Join(root, rel))
	return err == nil || !errors.Is(err, fs.ErrNotExist)
}
