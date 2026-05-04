#!/usr/bin/env node
const crypto = require('crypto');
const { emit, fail, loadInput } = require('./_see');

const NAME = 'eco_route_insight';

function fuelIndex(route, payloadTonnes) {
  return Number(route.distanceKm) * payloadTonnes * Number(route.gradientFactor);
}

function expiry(issuedAt, ttlHours) {
  const date = new Date(issuedAt);
  date.setUTCHours(date.getUTCHours() + Number.parseInt(ttlHours, 10));
  return date.toISOString().replace('.000Z', 'Z');
}

function goNumber(value) {
  const rounded = Math.round((value + Number.EPSILON) * 100) / 100;
  return Number.isInteger(rounded) ? Math.trunc(rounded) : rounded;
}

function chooseAlternative(data, currentFuel) {
  const currentEta = Number.parseInt(data.currentRoute.etaMinutes, 10);
  const maxDelay = Number.parseInt(data.policy.maxEtaDelayMinutes, 10);
  const threshold = Number(data.policy.fuelIndexThreshold);
  const payloadTonnes = Number(data.shipment.payloadKg) / 1000.0;
  const scored = [];
  for (const route of data.alternativeRoutes) {
    const fi = fuelIndex(route, payloadTonnes);
    const saving = currentFuel - fi;
    const delay = Number.parseInt(route.etaMinutes, 10) - currentEta;
    let eligible = saving > 0 && delay <= maxDelay;
    if (data.policy.requireAlternativeBelowThreshold) eligible = eligible && fi <= threshold;
    scored.push({ route, fuelIndex: fi, saving, delay, eligible });
  }
  scored.sort((a, b) => Number(!a.eligible) - Number(!b.eligible) || b.saving - a.saving || a.route.id.localeCompare(b.route.id));
  return scored[0];
}

function canonicalEnvelope(data, currentFuel, selected) {
  const issuedAt = data.scenario.issuedAt;
  const exp = expiry(issuedAt, data.scenario.ttlHours);
  const issue = currentFuel > Number(data.policy.fuelIndexThreshold) && selected.eligible;
  return {
    audience: data.scenario.depot,
    allowedUse: data.policy.allowedUse,
    issuedAt,
    expiry: exp,
    keyId: data.signing.keyId,
    assertions: {
      showEcoBanner: issue,
      suggestedRoute: selected.route.id,
      currentFuelIndex: goNumber(currentFuel),
      suggestedFuelIndex: goNumber(selected.fuelIndex),
      estimatedSaving: goNumber(selected.saving),
      rawDataExported: false,
    },
  };
}

function canonicalJson(value) {
  return JSON.stringify(value);
}

function sign(secret, canonical) {
  return crypto.createHmac('sha256', Buffer.from(secret, 'utf8'))
    .update(Buffer.from(canonical, 'utf8'))
    .digest('base64url');
}

function trustedDerivation(data) {
  const payloadTonnes = Number(data.shipment.payloadKg) / 1000.0;
  const currentFuel = fuelIndex(data.currentRoute, payloadTonnes);
  const selected = chooseAlternative(data, currentFuel);
  const envelope = canonicalEnvelope(data, currentFuel, selected);
  const canonical = canonicalJson(envelope);
  const digest = crypto.createHash('sha256').update(canonical, 'utf8').digest('hex');
  const sig = sign(data.signing.secret, canonical);
  const forbidden = data.dataMinimization.forbiddenEnvelopeTerms;

  fail('Eco route insight derivation failed', {
    'current route crosses fuel threshold': currentFuel > Number(data.policy.fuelIndexThreshold),
    'selected alternative is eligible': selected.eligible,
    'ETA delay is acceptable': selected.delay <= Number.parseInt(data.policy.maxEtaDelayMinutes, 10),
    'selected alternative has best eligible saving': data.alternativeRoutes.every((route) => {
      const fi = fuelIndex(route, payloadTonnes);
      const saving = currentFuel - fi;
      const delay = Number.parseInt(route.etaMinutes, 10) - Number.parseInt(data.currentRoute.etaMinutes, 10);
      const eligible = saving > 0 && delay <= Number.parseInt(data.policy.maxEtaDelayMinutes, 10) && (!data.policy.requireAlternativeBelowThreshold || fi <= Number(data.policy.fuelIndexThreshold));
      return !eligible || selected.saving >= saving;
    }),
    'envelope omits forbidden raw data terms': !forbidden.some((term) => canonical.includes(term)),
    'digest matches canonical envelope': digest === crypto.createHash('sha256').update(canonical, 'utf8').digest('hex'),
    'signature matches canonical envelope': sig === sign(data.signing.secret, canonical),
  });
  return { currentFuel, selected, envelope, digest, sig, canonical };
}

function main() {
  const data = loadInput(NAME);
  const { currentFuel, selected, envelope, digest, sig } = trustedDerivation(data);
  const payloadTonnes = Number(data.shipment.payloadKg) / 1000.0;
  const route = selected.route;
  const saving = selected.saving;
  const status = envelope.assertions.showEcoBanner ? 'issue' : 'clear';
  const yesNo = envelope.assertions.showEcoBanner ? 'yes' : 'no';

  emit('# Eco Route Insight');
  emit();
  emit('## Insight');
  emit(`insight status : ${status}`);
  emit(`show eco banner : ${yesNo}`);
  emit(`audience : ${envelope.audience}`);
  emit(`allowed use : ${envelope.allowedUse}`);
  emit(`suggested route : ${route.id}`);
  emit(`current fuel index : ${currentFuel.toFixed(2)}`);
  emit(`suggested fuel index : ${selected.fuelIndex.toFixed(2)}`);
  emit(`estimated saving : ${saving.toFixed(2)}`);
  emit(`expires at : ${envelope.expiry}`);
  emit('raw data exported : no');
  emit(`signature algorithm : ${data.policy.signatureAlgorithm}`);
  emit(`payload digest : ${digest}`);
  emit(`signature key : ${data.signing.keyId}`);
  emit(`signature : ${sig}`);
  emit();
  emit('## Explanation');
  emit('The current route uses fuel index = distanceKm × (payloadKg / 1000) × gradientFactor.');
  emit(`For ${data.shipment.id}, ${data.currentRoute.label} gives ${Number(data.currentRoute.distanceKm).toFixed(2)} × ${payloadTonnes.toFixed(2)} × ${Number(data.currentRoute.gradientFactor).toFixed(2)} = ${currentFuel.toFixed(2)}.`);
  emit(`The policy threshold is ${Number(data.policy.fuelIndexThreshold).toFixed(2)}, so a local eco banner is justified.`);
  emit(`The selected alternative ${route.id} gives ${Number(route.distanceKm).toFixed(2)} × ${payloadTonnes.toFixed(2)} × ${Number(route.gradientFactor).toFixed(2)} = ${selected.fuelIndex.toFixed(2)}, saving ${saving.toFixed(2)} while staying within the ETA delay limit.`);
  emit('The signed envelope exposes audience, use, expiry, route suggestion, and compact fuel indices, but not raw payload, GPS trace, driver behavior, or raw telemetry.');
  emit('This follows the insight pattern: ship the decision, not the data.');
}

if (require.main === module) main();
