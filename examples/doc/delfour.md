# Delfour

## What this example is about

This example imagines a supermarket self-scanner that can show a low-sugar shopping suggestion without exposing sensitive household health information.

The scenario begins with a household need related to diabetes. Instead of sending that medical fact to the scanner, the system creates a neutral, scoped insight: track sugar per serving while scanning at a specific retailer, for a limited time, under a specific policy.

## How it works, in plain language

The scanner sees the current product and a signed insight envelope. The envelope says what the scanner may do, what it must not do, when the insight expires, and what metric matters. In this fixture, the scanned biscuits exceed the sugar threshold, so the program selects a lower-sugar alternative.

The important design idea is data minimization. The scanner receives just enough information to help the shopper. It does not receive the original medical condition as raw data.

## What to notice in the output

The output says the scanner is allowed to show a neutral shopping insight, identifies the scanned product, and suggests a lower-sugar alternative. It also prints policy and audit details: signature algorithm, expiry, reason text, audit count, and bus files written.

## What the trust gate checks

Before emitting the explanation, the trust gate verifies authorization, marketing prohibition, delete-duty timing, payload minimization, payload hash, signature metadata, existence of a lower-sugar alternative, that a banner is warranted, and that audit/bus counts match the fixture.

## Run it

From the repository root:

```sh
node examples/delfour.js
```

## Files

- [JavaScript example](../delfour.js)
- [Input data](../input/delfour.json)
- [Expected output](../output/delfour.md)
