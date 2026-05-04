# School Placement Route Audit

## What this example is about

This example audits a school-placement rule. A support tool initially chooses schools by straight-line distance, but the audit asks whether walking-route distance and family preferences tell a different story.

The example is about explainability in public-service decision support: a map distance can look objective while hiding barriers and detours.

## How it works, in plain language

For each student and school, the input includes straight-line distance and walking distance. The original rule chooses the school with the smallest straight-line distance, using preference rank only as a tie-breaker.

The audit recomputes the choice using walking distance plus a preference penalty. It flags a student if the original assignment differs from the audited best choice or if the walking distance is above the configured maximum.

## What to notice in the output

The audit fails and names Ada, Björn, and Davi as affected. Ada has the largest hidden detour. The output also gives recommended assignments, so the reader can see the proposed correction rather than only the failure label.

## What the trust gate checks

The trust gate verifies that the distance matrix is complete, every student has a distance to every school, policy limits are positive, the affected flag follows the audit rule, recommended assignments are known schools, and the largest hidden detour is truly maximal.

## Run it

From the repository root:

```sh
node examples/school_placement_audit.js
```

## Files

- [JavaScript example](../school_placement_audit.js)
- [Input data](../input/school_placement_audit.json)
- [Reference output](../output/school_placement_audit.md)
