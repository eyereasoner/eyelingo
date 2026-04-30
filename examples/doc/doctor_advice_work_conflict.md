# Doctor Advice Work Conflict

`doctor_advice_work_conflict` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on policy conflict resolution for remote-work and office-work advice. Its input fixture is organized around `CaseName`, `Question`, `Person`, `DoctorCanDoJobs`, `SubclassOf`, `Requests`, `ConflictPolicies`, `Expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Technology** example. It demonstrates data representation, interoperability, policies, and computational artifacts in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: overall decision for Jos : RemoteWorkOnly Request_Jos_Prog_Home : raw=Deny+Permit status=BothPermitDeny effective=Permit Request_Jos_Prog_Office : raw=Deny+Permit status=BothPermitDeny effective=Deny

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - Jos is classified as Sick from condition Flu C2 OK - home programming request keeps both Permit and Deny before resolution C3 OK - home programming conflict resolves to Permit

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/doctor_advice_work_conflict.json](../input/doctor_advice_work_conflict.json)

Go translation: [../doctor_advice_work_conflict.go](../doctor_advice_work_conflict.go)

Expected Markdown output: [../output/doctor_advice_work_conflict.md](../output/doctor_advice_work_conflict.md)
