# Cell Marker Panel

## Answer
The exact marker panel uses 8 genes and separates all 28 cell-type pairs.
- case : pbmc-plus-epithelium-marker-panel
- positive anchors : 8/8 cell populations
- assay cost : 86

- Selected panel:
 - CD3D (T-cell receptor complex) anchors=:T pairSeparations=7
 - MS4A1 (B-cell membrane marker) anchors=:B pairSeparations=7
 - NKG7 (cytotoxic granule marker) anchors=:NK pairSeparations=7
 - CD14 (monocyte co-receptor) anchors=:Mono pairSeparations=7
 - FCER1A (dendritic-cell receptor) anchors=:DC pairSeparations=7
 - EPCAM (epithelial adhesion) anchors=:Epi pairSeparations=7
 - PECAM1 (endothelial junction) anchors=:Endo pairSeparations=7
 - COL1A1 (fibroblast matrix) anchors=:Fib pairSeparations=7

## Reason why
- Candidate genes are first filtered by assay policy: housekeeping, ribosomal, mitochondrial-stress, and ambient-RNA features cannot be used as lineage markers. Each remaining gene then derives two facts: which cell-type pairs it separates by the expression-gap threshold, and which cell type, if any, it positively anchors with low off-target signal. Dominated backup markers are removed before an exact branch-and-bound search solves the paired set-cover goal.
- cell populations : 8
- cell-type pairs : 28
- candidate genes : 16
- usable after QC : 12
- dominated markers pruned : 2
- retained markers searched : 10
- search states visited : 71
- branch-and-bound prunes : impossible=28 size=0 lowerBound=0
- minimality subsets checked : 968
- panel pair margin sum : 466
- panel anchor margin sum : 57
- Derived positive anchors:
 - CD3D anchors :T: on=9 maxOff=2 margin=7
 - MS4A1 anchors :B: on=9 maxOff=1 margin=8
 - NKG7 anchors :NK: on=9 maxOff=3 margin=6
 - CD14 anchors :Mono: on=8 maxOff=2 margin=6
 - FCER1A anchors :DC: on=9 maxOff=2 margin=7
 - EPCAM anchors :Epi: on=9 maxOff=1 margin=8
 - PECAM1 anchors :Endo: on=9 maxOff=1 margin=8
 - COL1A1 anchors :Fib: on=9 maxOff=2 margin=7

## Check
C1 OK - every pair of cell populations is separated by at least one selected gene.
C2 OK - every cell population has a positive high-confidence anchor.
C3 OK - excluded housekeeping, ribosomal, mitochondrial, and ambient markers are not selected.
C4 OK - dominated backup markers are not selected.
C5 OK - selected genes are stable and assayable.
C6 OK - selected expression signatures are unique for all populations.
C7 OK - all reported separations meet the configured expression-gap threshold.
C8 OK - no smaller panel satisfies both pair separation and positive-anchor constraints.
C9 OK - reported margin and cost totals match the selected genes.

## Go audit details
- platform : go1.26.2 linux/amd64
- question : Which compact marker panel separates the reference cell populations?
- cell populations : 8
- pairwise separation obligations : 28
- separation threshold : 5
- anchor rule : on >= 7 and every off-target <= 3
- genes in fixture : 16
- usable genes : 12
- excluded genes : 4
- dominated genes pruned : 2
 - KRT18 dominated by EPCAM pairSeparations=7 anchors=1
 - VWF dominated by PECAM1 pairSeparations=7 anchors=1
- retained search genes : 10
- selected genes : CD3D, MS4A1, NKG7, CD14, FCER1A, EPCAM, PECAM1, COL1A1
- selected gene count : 8
- selected assay cost : 86
- covered pair mask : 1111111111111111111111111111
- covered anchor mask : 11111111
- pair margin sum : 466
- anchor margin sum : 57
- cell signatures:
 - :T T cell           CD3D=9 MS4A1=1 NKG7=3 CD14=0 FCER1A=0 EPCAM=0 PECAM1=0 COL1A1=0
 - :B B cell           CD3D=1 MS4A1=9 NKG7=1 CD14=0 FCER1A=1 EPCAM=0 PECAM1=0 COL1A1=0
 - :NK NK cell          CD3D=2 MS4A1=1 NKG7=9 CD14=2 FCER1A=0 EPCAM=0 PECAM1=0 COL1A1=0
 - :Mono Monocyte         CD3D=1 MS4A1=1 NKG7=2 CD14=8 FCER1A=2 EPCAM=0 PECAM1=0 COL1A1=0
 - :DC Dendritic cell   CD3D=1 MS4A1=1 NKG7=2 CD14=2 FCER1A=9 EPCAM=0 PECAM1=0 COL1A1=0
 - :Epi Epithelial cell  CD3D=0 MS4A1=0 NKG7=0 CD14=0 FCER1A=0 EPCAM=9 PECAM1=1 COL1A1=1
 - :Endo Endothelial cell CD3D=0 MS4A1=0 NKG7=0 CD14=0 FCER1A=0 EPCAM=1 PECAM1=9 COL1A1=2
 - :Fib Fibroblast       CD3D=0 MS4A1=0 NKG7=0 CD14=0 FCER1A=0 EPCAM=1 PECAM1=1 COL1A1=9
- retained gene coverage:
 - PTPRC pairSeparations=15 anchors= cost=10
 - CD3D pairSeparations=7 anchors=:T cost=10
 - MS4A1 pairSeparations=7 anchors=:B cost=10
 - NKG7 pairSeparations=7 anchors=:NK cost=11
 - CD14 pairSeparations=7 anchors=:Mono cost=11
 - FCER1A pairSeparations=7 anchors=:DC cost=12
 - EPCAM pairSeparations=7 anchors=:Epi cost=10
 - PECAM1 pairSeparations=7 anchors=:Endo cost=10
 - COL1A1 pairSeparations=7 anchors=:Fib cost=12
 - LYZ pairSeparations=13 anchors= cost=11
- excluded gene reasons:
 - ACTB : housekeeping gene is not lineage-specific
 - RPLP0 : ribosomal gene fails marker policy
 - MT-CO1 : mitochondrial stress signal is batch-sensitive
 - MALAT1 : ambient RNA risk in dissociation protocol
- search states visited : 71
- leaf states reached : 8
- include branches : 35
- exclude branches : 35
- impossible prunes : 28
- size prunes : 0
- lower-bound prunes : 0
- complete solutions found : 4
- checks passed : 9/9
- all checks pass : yes
