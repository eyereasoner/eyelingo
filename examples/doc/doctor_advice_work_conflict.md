# Doctor Advice Work Conflict

`doctor_advice_work_conflict` is a Go translation/adaptation of Eyeling's `doctor-advice-work-conflict.n3`.

The context is conflicting recommendations. Remote-work and office-work advice are evaluated through policy precedence and conflict-resolution rules, producing a clear final decision.

## What it demonstrates

This example is mainly in the **Technology** category. Policy conflict resolution for remote-work and office-work advice.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/doctor_advice_work_conflict.json](../input/doctor_advice_work_conflict.json)

Go translation: [../doctor_advice_work_conflict.go](../doctor_advice_work_conflict.go)

Expected Markdown output: [../output/doctor_advice_work_conflict.md](../output/doctor_advice_work_conflict.md)
