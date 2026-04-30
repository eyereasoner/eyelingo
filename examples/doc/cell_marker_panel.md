# Cell Marker Panel

`cell_marker_panel` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on gene-marker panel selection that separates modeled cell populations. Its input fixture is organized around `CaseName`, `Question`, `CellTypes`, `Genes`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The exact marker panel uses 8 genes and separates all 28 cell-type pairs. case : pbmc-plus-epithelium-marker-panel positive anchors : 8/8 cell populations

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - every pair of cell populations is separated by at least one selected gene. C2 OK - every cell population has a positive high-confidence anchor. C3 OK - excluded housekeeping, ribosomal, mitochondrial, and ambient markers are not selected.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/cell_marker_panel.json](../input/cell_marker_panel.json)

Go translation: [../cell_marker_panel.go](../cell_marker_panel.go)

Expected Markdown output: [../output/cell_marker_panel.md](../output/cell_marker_panel.md)
