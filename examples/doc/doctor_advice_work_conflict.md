# Doctor Advice Work Conflict

`doctor_advice_work_conflict` is a Go translation/adaptation of Eyeling's `doctor-advice-work-conflict.n3`.

## Background

Policy engines often receive advice or rules from multiple authorities, and those instructions can conflict. A work-location example may involve medical advice, employer requirements, remote-work permissions, and office-attendance rules. The useful output is not just a final answer but a conflict analysis that records which rule has priority and why the losing rule is overridden.

## What it demonstrates

**Category:** Technology. Policy conflict resolution for remote-work and office-work advice.

## How the Go implementation works

The implementation derives permit and deny conclusions for home-work and office-work recommendations, then runs a deterministic conflict-resolution pass. Precedence rules choose the final outcome when both supporting and opposing advice are present.

The report keeps the conflicting intermediate conclusions visible before showing the resolved decision.

## Files

Input JSON: [../input/doctor_advice_work_conflict.json](../input/doctor_advice_work_conflict.json)

Go implementation: [../doctor_advice_work_conflict.go](../doctor_advice_work_conflict.go)

Expected Markdown output: [../output/doctor_advice_work_conflict.md](../output/doctor_advice_work_conflict.md)
