# ODRL + DPV Risk Ranking

## What this example is about

This example translates an N3 policy-reasoning example into the repository's JavaScript example format. It looks at a small terms-of-service style agreement and compares it with a consumer profile.

The agreement says the provider may remove an account, change terms with only three days of notice, share user data with partners, and prohibit data export. The consumer profile says these areas matter: account/data removal, prior notice for changes, consent for sharing, and data portability.

The purpose is to show how policy clauses can be turned into a ranked risk report. Instead of saying only "this agreement has problems," the example lists which clauses matter most, why they matter, and what mitigation would make them safer.

## How it works, in plain language

The input uses a simple JSON version of the original policy graph. Each rule has a type, such as `Permission` or `Prohibition`, an action, a target, a clause link, and optional duties or constraints.

The program checks four policy patterns:

1. Account/data removal is permitted without a notice constraint and without a duty to inform.
2. Terms can change with fewer notice days than the consumer requires.
3. Data sharing is permitted without an explicit consent constraint.
4. Data export is prohibited, which conflicts with the consumer's portability need.

When one of those patterns is found, the program creates a risk record. The score is the rule's base risk plus the importance of the consumer need, capped at 100. Scores of 80 or above are classified as high risk; scores from 50 to 79 are moderate risk.

## What to notice in the output

The output starts with the agreement, profile, number of risks found, and the highest-risk clause. It then prints a ranked report. The top item is clause C1 because the account/data removal risk reaches the maximum score of 100.

Each risk line includes the score, risk level, severity, clause ID, plain-language reason, consequences, impacts, and mitigation advice. This makes the report useful to both technical readers and people who mainly want to know which contract language should be reviewed first.

## What the trust gate checks

The trust gate verifies that at least one risk is derived, every risk keeps its link back to a human-readable clause, every risk links to a consumer need, consequences and impacts are present, mitigation advice is present, scores are capped, levels and severities follow the configured thresholds, and the final ranking is in descending score order.

## Run it

From the repository root:

```sh
node examples/odrl_dpv_risk_ranked.js
```

## Files

- [JavaScript example](../odrl_dpv_risk_ranked.js)
- [Input data](../input/odrl_dpv_risk_ranked.json)
- [Reference output](../output/odrl_dpv_risk_ranked.md)
