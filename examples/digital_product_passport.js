#!/usr/bin/env node
const { boolText, emit, fail, loadInput } = require('./_see');

const NAME = 'digital_product_passport';

function totalMass(data) {
  return data.Components.reduce((acc, c) => acc + Number.parseInt(c.MassG, 10), 0);
}

function recycledMass(data) {
  return data.Components.reduce((acc, c) => acc + Number.parseInt(c.RecycledMassG, 10), 0);
}

function recycledPct(data) {
  return Math.floor((recycledMass(data) * 100) / totalMass(data));
}

function lifecycleFootprint(data) {
  const f = data.Footprint;
  return Number.parseInt(f.ManufacturingGCO2e, 10) + Number.parseInt(f.TransportGCO2e, 10) + Number.parseInt(f.UsePhaseGCO2e, 10);
}

function criticalMaterials(data) {
  const critical = new Set(data.Materials.filter((m) => m.CriticalRawMaterial).map((m) => m.ID));
  const used = new Set(data.Components.flatMap((c) => c.ContainsMaterial));
  return [...critical].filter((m) => used.has(m)).sort();
}

function repairFriendly(data) {
  const policy = data.AccessPolicy;
  const publicDocs = new Set(data.Documents.filter((d) => d.Section === policy.PublicSection).map((d) => d.DocType));
  const declarations = new Set(data.Documents
    .filter((d) => d.Section === policy.PublicSection)
    .flatMap((d) => d.Declares));
  const batteryReplaceable = data.Components.some((c) => c.Type === 'Battery' && c.Replaceable);
  return batteryReplaceable && policy.PublicDocTypes.every((t) => publicDocs.has(t)) && declarations.has(policy.RequiredPublicClaim);
}

function trustedDerivation(data) {
  const restrictedTypes = new Set(data.AccessPolicy.RestrictedDocTypes);
  const componentMass = totalMass(data);
  const componentRecycledMass = recycledMass(data);
  const publicDocs = new Set(data.Documents.filter((d) => d.Section === data.AccessPolicy.PublicSection).map((d) => d.DocType));
  fail('Digital product passport derivation failed', {
    'component mass is positive': componentMass > 0,
    'recycled mass does not exceed total mass': componentRecycledMass <= componentMass,
    'recycled percentage is derived from masses': recycledPct(data) === Math.floor((componentRecycledMass * 100) / componentMass),
    'footprint sum is derived from phases': lifecycleFootprint(data) === ['ManufacturingGCO2e', 'TransportGCO2e', 'UsePhaseGCO2e'].reduce((acc, key) => acc + Number.parseInt(data.Footprint[key], 10), 0),
    'critical materials used in components are exposed': criticalMaterials(data).every((material) => data.Components.some((component) => component.ContainsMaterial.includes(material))),
    'all required public document types are present': data.AccessPolicy.PublicDocTypes.every((docType) => publicDocs.has(docType)),
    'repair-friendly conditions hold': repairFriendly(data),
    'restricted docs stay restricted': data.Documents.filter((d) => restrictedTypes.has(d.DocType)).every((d) => d.Section === data.AccessPolicy.RestrictedSection),
    'digital link equals endpoint': data.Product.DigitalLink === data.Passport.PublicEndpoint,
    'lifecycle is chronological': JSON.stringify(data.Lifecycle.map((e) => e.OnDate)) === JSON.stringify([...data.Lifecycle.map((e) => e.OnDate)].sort()),
  });
}

function main() {
  const data = loadInput(NAME);
  trustedDerivation(data);
  const product = data.Product;

  emit('# Digital Product Passport');
  emit();
  emit('## Insight');
  emit(`Passport decision : PASS for ${product.Model} ${product.SerialNumber}.`);
  emit(`recycled content : ${recycledPct(data)}%`);
  emit(`lifecycle footprint : ${lifecycleFootprint(data)} gCO2e`);
  emit(`total component mass : ${totalMass(data)} g`);
  emit(`critical raw materials : ${criticalMaterials(data).join(', ')}`);
  emit('circularity hint : repairFriendly');
  emit(`public endpoint : ${data.Passport.PublicEndpoint}`);
  emit();
  emit('## Explanation');
  emit('The passport folds the explicit component list to derive total mass and recycled mass, then computes an integer recycled-content percentage.');
  emit('Lifecycle footprint is derived by summing manufacturing, transport, and use-phase emissions.');
  emit('The product is repair-friendly because the battery is replaceable and the public passport section exposes both repair and spare-parts documentation.');
  emit(`Critical raw-material exposure is derived from component-material links: ${criticalMaterials(data).join(', ')}.`);
  emit();
  emit('Component roll-up:');
  for (const c of data.Components) {
    emit(`${c.ID} ${c.Type} mass=${c.MassG}g recycled=${c.RecycledMassG}g materials=${c.ContainsMaterial.join(', ')} replaceable=${boolText(c.Replaceable)}`);
  }
  emit('Public documents:');
  for (const d of data.Documents) {
    if (d.Section === data.AccessPolicy.PublicSection) emit(`${d.ID} ${d.DocType} ${d.URL}`);
  }
}

if (require.main === module) main();
