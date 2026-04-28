// complex_numbers.go
//
// A self-contained Go translation of complex.n3 from the Eyeling examples.
//
// The original N3 file defines rules for complex polar conversion,
// quadrant-sensitive angle selection, complex exponentiation, and inverse sine /
// cosine over complex numbers. The test query asks for six derived complex
// values:
//
//    sqrt(-1), e^(i*pi), i^i, e^(-pi/2), asin(2), and acos(2)
//
// This Go version keeps those rules explicit instead of using Go's cmplx.Pow,
// cmplx.Asin, or cmplx.Acos helpers. That makes the same mathematical proof
// steps visible and auditable.
//
// Run:
//
//    go run complex_numbers.go
//
// The program has no third-party dependencies.
package main

import (
    "fmt"
    "math"
    "os"
    "runtime"
    "strings"
)

const (
    epsilon = 1e-9
)

type Complex struct {
    Re float64
    Im float64
}

type Polar struct {
    R         float64
    Theta     float64
    Quadrant  string
    DialRule  string
    BaseAngle float64
}

type ExponentCase struct {
    ID             string
    Expression     string
    Base           Complex
    Power          Complex
    Polar          Polar
    MagnitudeScale float64
    AngleMix       float64
    Result         Complex
    Expected       Complex
    Reason         string
}

type InverseCase struct {
    ID          string
    Operation   string
    Expression  string
    Input       Complex
    LeftRadius  float64
    RightRadius float64
    E           float64
    F           float64
    LogTerm     float64
    Result      Complex
    Expected    Complex
    Reason      string
}

type Analysis struct {
    Question         string
    ExponentCases    []ExponentCase
    InverseCases     []InverseCase
    Checks           []Check
    PrimitiveFacts   int
    RuleApplications int
    PolarDerivations int
    FiniteValues     int
}

type Check struct {
    Label string
    OK    bool
    Text  string
}

func main() {
    analysis := derive()
    printAnswer(analysis)
    printReason(analysis)
    printChecks(analysis)
    printAudit(analysis)
    if !allChecksOK(analysis.Checks) {
        os.Exit(1)
    }
}

func derive() Analysis {
    question := "Evaluate the complex-number facts queried by complex.n3."

    expInputs := []struct {
        id       string
        expr     string
        base     Complex
        power    Complex
        expected Complex
        reason   string
    }{
        {
            id:       "C1",
            expr:     "((-1 0) (0.5 0)) complex:exponentiation ?C1",
            base:     Complex{-1, 0},
            power:    Complex{0.5, 0},
            expected: Complex{0, 1},
            reason:   "the polar angle of -1+0i is pi, so the half power rotates to pi/2",
        },
        {
            id:       "C2",
            expr:     "((e 0) (0 pi)) complex:exponentiation ?C2",
            base:     Complex{math.E, 0},
            power:    Complex{0, math.Pi},
            expected: Complex{-1, 0},
            reason:   "Euler's identity falls out of the exponent rule because ln(e)=1",
        },
        {
            id:       "C3",
            expr:     "((0 1) (0 1)) complex:exponentiation ?C3",
            base:     Complex{0, 1},
            power:    Complex{0, 1},
            expected: Complex{math.Exp(-math.Pi / 2), 0},
            reason:   "i has polar angle pi/2, so i^i becomes exp(-pi/2)",
        },
        {
            id:       "C4",
            expr:     "((e 0) (-pi/2 0)) complex:exponentiation ?C4",
            base:     Complex{math.E, 0},
            power:    Complex{-math.Pi / 2, 0},
            expected: Complex{math.Exp(-math.Pi / 2), 0},
            reason:   "a real exponent of e gives the same exp(-pi/2) value as i^i",
        },
    }

    exps := make([]ExponentCase, 0, len(expInputs))
    for _, in := range expInputs {
        polar := toPolarN3(in.base)
        result, magnitudeScale, angleMix := exponentiateN3(in.base, in.power, polar)
        exps = append(exps, ExponentCase{
            ID:             in.id,
            Expression:     in.expr,
            Base:           in.base,
            Power:          in.power,
            Polar:          polar,
            MagnitudeScale: magnitudeScale,
            AngleMix:       angleMix,
            Result:         result,
            Expected:       in.expected,
            Reason:         in.reason,
        })
    }

    inverseInputs := []struct {
        id       string
        op       string
        expr     string
        input    Complex
        expected Complex
        reason   string
    }{
        {
            id:       "C5",
            op:       "asin",
            expr:     "(2 0) complex:asin ?C5",
            input:    Complex{2, 0},
            expected: Complex{math.Pi / 2, math.Log(2 + math.Sqrt(3))},
            reason:   "the real input is outside [-1,1], so the inverse sine has a positive imaginary part",
        },
        {
            id:       "C6",
            op:       "acos",
            expr:     "(2 0) complex:acos ?C6",
            input:    Complex{2, 0},
            expected: Complex{0, -math.Log(2 + math.Sqrt(3))},
            reason:   "the companion inverse cosine carries the opposite imaginary part",
        },
    }

    inverses := make([]InverseCase, 0, len(inverseInputs))
    for _, in := range inverseInputs {
        inv := inverseTrigN3(in.id, in.op, in.expr, in.input, in.expected, in.reason)
        inverses = append(inverses, inv)
    }

    analysis := Analysis{
        Question:         question,
        ExponentCases:    exps,
        InverseCases:     inverses,
        PrimitiveFacts:   len(expInputs) + len(inverseInputs),
        PolarDerivations: len(expInputs),
        // Rule applications are counted as the named N3-style derivations that this
        // translation performs: one polar plus one exponentiation per exponent case,
        // and one inverse rule per inverse case.
        RuleApplications: len(expInputs)*2 + len(inverseInputs),
    }
    analysis.FiniteValues = countFiniteValues(analysis)
    analysis.Checks = buildChecks(analysis)
    return analysis
}

func toPolarN3(z Complex) Polar {
    r := math.Sqrt(z.Re*z.Re + z.Im*z.Im)
    if r == 0 {
        return Polar{R: 0, Theta: 0, Quadrant: "origin", DialRule: "degenerate zero magnitude", BaseAngle: 0}
    }
    baseAngle := math.Acos(math.Abs(z.Re) / r)
    p := Polar{R: r, BaseAngle: baseAngle}
    switch {
    case z.Re >= 0 && z.Im >= 0:
        p.Theta = baseAngle
        p.Quadrant = "I or +axis"
        p.DialRule = "x>=0 and y>=0 => theta=T"
    case z.Re < 0 && z.Im >= 0:
        p.Theta = math.Pi - baseAngle
        p.Quadrant = "II or -real axis"
        p.DialRule = "x<0 and y>=0 => theta=pi-T"
    case z.Re < 0 && z.Im < 0:
        p.Theta = math.Pi + baseAngle
        p.Quadrant = "III"
        p.DialRule = "x<0 and y<0 => theta=pi+T"
    default:
        p.Theta = 2*math.Pi - baseAngle
        p.Quadrant = "IV"
        p.DialRule = "x>=0 and y<0 => theta=2*pi-T"
    }
    return p
}

func exponentiateN3(base, power Complex, polar Polar) (Complex, float64, float64) {
    if polar.R == 0 {
        return Complex{0, 0}, 0, 0
    }
    // N3 variables, renamed:
    // Z1 = R^C
    // Z4 = e^(-D*T)
    // Z5 = ln(R), recovered in N3 by solving e^Z5 = R
    // Z8 = D*ln(R) + C*T
    z1 := math.Pow(polar.R, power.Re)
    z4 := math.Exp(-power.Im * polar.Theta)
    z8 := power.Im*math.Log(polar.R) + power.Re*polar.Theta
    magnitudeScale := z1 * z4
    return Complex{
        Re: magnitudeScale * math.Cos(z8),
        Im: magnitudeScale * math.Sin(z8),
    }, magnitudeScale, z8
}

func inverseTrigN3(id, op, expr string, input, expected Complex, reason string) InverseCase {
    // Shared N3 core for both asin and acos:
    // Z5 = sqrt((1+A)^2 + B^2)
    // Z9 = sqrt((1-A)^2 + B^2)
    // E  = (Z5-Z9)/2
    // F  = (Z5+Z9)/2
    leftRadius := math.Sqrt(math.Pow(1+input.Re, 2) + input.Im*input.Im)
    rightRadius := math.Sqrt(math.Pow(1-input.Re, 2) + input.Im*input.Im)
    e := (leftRadius - rightRadius) / 2
    f := (leftRadius + rightRadius) / 2
    logTerm := math.Log(f + math.Sqrt(f*f-1))
    var result Complex
    if op == "asin" {
        result = Complex{Re: math.Asin(e), Im: logTerm}
    } else {
        result = Complex{Re: math.Acos(e), Im: -logTerm}
    }
    return InverseCase{
        ID:          id,
        Operation:   op,
        Expression:  expr,
        Input:       input,
        LeftRadius:  leftRadius,
        RightRadius: rightRadius,
        E:           e,
        F:           f,
        LogTerm:     logTerm,
        Result:      result,
        Expected:    expected,
        Reason:      reason,
    }
}

func buildChecks(a Analysis) []Check {
    checks := []Check{}
    checks = append(checks, Check{
        Label: "C1",
        OK:    close(a.ExponentCases[0].Polar.Theta, math.Pi) && close(a.ExponentCases[1].Polar.Theta, 0) && close(a.ExponentCases[2].Polar.Theta, math.Pi/2),
        Text:  "N3 dial rules assign the expected polar angles for -1, e, and i.",
    })
    checks = append(checks, Check{
        Label: "C2",
        OK:    allExponentResultsExpected(a.ExponentCases),
        Text:  "all four complex exponentiation answers match the complex.n3 test facts.",
    })
    checks = append(checks, Check{
        Label: "C3",
        OK:    closeComplex(a.ExponentCases[2].Result, a.ExponentCases[3].Result),
        Text:  "i^i and e^(-pi/2) derive the same real value.",
    })
    checks = append(checks, Check{
        Label: "C4",
        OK:    allInverseResultsExpected(a.InverseCases),
        Text:  "asin(2) and acos(2) match the N3 inverse-trig derivations.",
    })
    checks = append(checks, Check{
        Label: "C5",
        OK:    closeComplex(complexSin(a.InverseCases[0].Result), a.InverseCases[0].Input) && closeComplex(complexCos(a.InverseCases[1].Result), a.InverseCases[1].Input),
        Text:  "sin(asin(2)) and cos(acos(2)) round-trip back to 2+0i.",
    })
    checks = append(checks, Check{
        Label: "C6",
        OK:    closeComplex(add(a.InverseCases[0].Result, a.InverseCases[1].Result), Complex{math.Pi / 2, 0}),
        Text:  "asin(2) + acos(2) equals pi/2 with cancelling imaginary parts.",
    })
    checks = append(checks, Check{
        Label: "C7",
        OK:    a.FiniteValues == 12,
        Text:  "all six complex outputs are finite real/imaginary pairs.",
    })
    checks = append(checks, Check{
        Label: "C8",
        OK:    a.PrimitiveFacts == 6 && a.PolarDerivations == 4 && a.RuleApplications == 10,
        Text:  "the translated rule count matches four exponentiation and two inverse-trig queries.",
    })
    return checks
}

func printAnswer(a Analysis) {
    fmt.Println("=== Answer ===")
    fmt.Println("The complex.n3 test query derives 6 complex-number facts.")
    fmt.Println("Computed values:")
    for _, c := range a.ExponentCases {
        fmt.Printf(" - %s %s = %s\n", c.ID, shortName(c.ID), formatComplex(c.Result))
    }
    for _, c := range a.InverseCases {
        fmt.Printf(" - %s %s = %s\n", c.ID, shortName(c.ID), formatComplex(c.Result))
    }
    fmt.Println("Key equivalences:")
    fmt.Printf(" - i^i = e^(-pi/2) = %s\n", formatFloat(a.ExponentCases[2].Result.Re))
    fmt.Printf(" - asin(2) + acos(2) = %s\n", formatComplex(add(a.InverseCases[0].Result, a.InverseCases[1].Result)))
    fmt.Println()
}

func printReason(a Analysis) {
    fmt.Println("=== Reason Why ===")
    fmt.Println("The N3 source first converts each complex base (x,y) to polar form using r=sqrt(x^2+y^2) and quadrant-sensitive dial rules. It then applies (a+bi)^(c+di)=r^c*e^(-d*theta)*(cos(d*ln(r)+c*theta), sin(d*ln(r)+c*theta)). The asin/acos rules use the same pair of distances from 1+a and 1-a, then recover the imaginary part with ln(F+sqrt(F^2-1)).")
    fmt.Printf("primitive test facts : %d\n", a.PrimitiveFacts)
    fmt.Printf("polar derivations : %d\n", a.PolarDerivations)
    fmt.Printf("rule applications : %d\n", a.RuleApplications)
    fmt.Println("polar bases:")
    for _, c := range a.ExponentCases {
        fmt.Printf(" - %s base=%s r=%s theta=%s quadrant=%s\n", c.ID, formatComplex(c.Base), formatFloat(c.Polar.R), formatFloat(c.Polar.Theta), c.Polar.Quadrant)
    }
    fmt.Println("exponentiation traces:")
    for _, c := range a.ExponentCases {
        fmt.Printf(" - %s scale=%s angleMix=%s result=%s; %s\n", c.ID, formatFloat(c.MagnitudeScale), formatFloat(c.AngleMix), formatComplex(c.Result), c.Reason)
    }
    fmt.Println("inverse-trig traces:")
    for _, c := range a.InverseCases {
        fmt.Printf(" - %s %s: E=%s F=%s lnTerm=%s result=%s; %s\n", c.ID, c.Operation, formatFloat(c.E), formatFloat(c.F), formatFloat(c.LogTerm), formatComplex(c.Result), c.Reason)
    }
    fmt.Println()
}

func printChecks(a Analysis) {
    fmt.Println("=== Check ===")
    for _, check := range a.Checks {
        status := "FAIL"
        if check.OK {
            status = "OK"
        }
        fmt.Printf("%s %s - %s\n", check.Label, status, check.Text)
    }
    fmt.Println()
}

func printAudit(a Analysis) {
    fmt.Println("=== Go audit details ===")
    fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
    fmt.Printf("question : %s\n", a.Question)
    fmt.Println("translated source : complex.n3")
    fmt.Printf("primitive test facts : %d\n", a.PrimitiveFacts)
    fmt.Printf("polar derivations : %d\n", a.PolarDerivations)
    fmt.Printf("rule applications : %d\n", a.RuleApplications)
    fmt.Println("N3 expressions:")
    for _, c := range a.ExponentCases {
        fmt.Printf(" - %s %s\n", c.ID, c.Expression)
    }
    for _, c := range a.InverseCases {
        fmt.Printf(" - %s %s\n", c.ID, c.Expression)
    }
    fmt.Println("dial details:")
    for _, c := range a.ExponentCases {
        fmt.Printf(" - %s base=%s absAngle=%s theta=%s rule=%s\n", c.ID, formatComplex(c.Base), formatFloat(c.Polar.BaseAngle), formatFloat(c.Polar.Theta), c.Polar.DialRule)
    }
    fmt.Println("derived exponentiation facts:")
    for _, c := range a.ExponentCases {
        fmt.Printf(" - %s base=%s power=%s result=%s expected=%s\n", c.ID, formatComplex(c.Base), formatComplex(c.Power), formatComplex(c.Result), formatComplex(c.Expected))
    }
    fmt.Println("derived inverse-trig facts:")
    for _, c := range a.InverseCases {
        fmt.Printf(" - %s input=%s result=%s expected=%s leftRadius=%s rightRadius=%s\n", c.ID, formatComplex(c.Input), formatComplex(c.Result), formatComplex(c.Expected), formatFloat(c.LeftRadius), formatFloat(c.RightRadius))
    }
    fmt.Printf("finite output components : %d/12\n", a.FiniteValues)
    passed := 0
    for _, check := range a.Checks {
        if check.OK {
            passed++
        }
    }
    fmt.Printf("checks passed : %d/%d\n", passed, len(a.Checks))
    if passed == len(a.Checks) {
        fmt.Println("all checks pass : yes")
    } else {
        fmt.Println("all checks pass : no")
    }
}

func shortName(id string) string {
    switch id {
    case "C1":
        return "sqrt(-1+0i)"
    case "C2":
        return "e^(i*pi)"
    case "C3":
        return "i^i"
    case "C4":
        return "e^(-pi/2)"
    case "C5":
        return "asin(2+0i)"
    case "C6":
        return "acos(2+0i)"
    default:
        return id
    }
}

func allExponentResultsExpected(cases []ExponentCase) bool {
    for _, c := range cases {
        if !closeComplex(c.Result, c.Expected) {
            return false
        }
    }
    return true
}

func allInverseResultsExpected(cases []InverseCase) bool {
    for _, c := range cases {
        if !closeComplex(c.Result, c.Expected) {
            return false
        }
    }
    return true
}

func countFiniteValues(a Analysis) int {
    count := 0
    for _, c := range a.ExponentCases {
        if finite(c.Result.Re) {
            count++
        }
        if finite(c.Result.Im) {
            count++
        }
    }
    for _, c := range a.InverseCases {
        if finite(c.Result.Re) {
            count++
        }
        if finite(c.Result.Im) {
            count++
        }
    }
    return count
}

func complexSin(z Complex) Complex {
    return Complex{
        Re: math.Sin(z.Re) * math.Cosh(z.Im),
        Im: math.Cos(z.Re) * math.Sinh(z.Im),
    }
}

func complexCos(z Complex) Complex {
    return Complex{
        Re: math.Cos(z.Re) * math.Cosh(z.Im),
        Im: -math.Sin(z.Re) * math.Sinh(z.Im),
    }
}

func add(a, b Complex) Complex {
    return Complex{Re: a.Re + b.Re, Im: a.Im + b.Im}
}

func closeComplex(a, b Complex) bool {
    return close(a.Re, b.Re) && close(a.Im, b.Im)
}

func close(a, b float64) bool {
    return math.Abs(a-b) <= epsilon
}

func finite(v float64) bool {
    return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func formatComplex(z Complex) string {
    re := cleanZero(z.Re)
    im := cleanZero(z.Im)
    sign := "+"
    if im < 0 {
        sign = "-"
        im = -im
    }
    return fmt.Sprintf("%s %s %si", formatFloat(re), sign, formatFloat(im))
}

func formatFloat(v float64) string {
    v = cleanZero(v)
    return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.12f", v), "0"), ".")
}

func cleanZero(v float64) float64 {
    if math.Abs(v) < 5e-13 {
        return 0
    }
    return v
}

func allChecksOK(checks []Check) bool {
    for _, check := range checks {
        if !check.OK {
            return false
        }
    }
    return true
}
