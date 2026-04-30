# Gradient Descent Step

`gradient_descent_step` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on certified single gradient-descent step for a quadratic objective. Its input fixture is organized around `caseName`, `question`, `function`, `start`, `stepSize`, `maxStepNorm`, `expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: start point : (8.000, 5.000) gradient : (10.000, 24.000) step size : 0.100

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - gradient was derived from f(x,y) = (x-3)^2 + 2(y+1)^2 C2 OK - step size is positive and below the conservative bound C3 OK - new point is (7.000, 2.600)

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/gradient_descent_step.json](../input/gradient_descent_step.json)

Go translation: [../gradient_descent_step.go](../gradient_descent_step.go)

Expected Markdown output: [../output/gradient_descent_step.md](../output/gradient_descent_step.md)
