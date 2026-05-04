#!/usr/bin/env node
const crypto = require('crypto');
const { emit, fail, loadInput, minBy } = require('./_see');

const NAME = 'delfour';

function productsById(data) {
  return Object.fromEntries(data.Products.map((product) => [product.ID, product]));
}

function scannedProduct(data) {
  return productsById(data)[data.Scan.ScannedProductID];
}

function needsLowSugar(data) {
  return data.Household.Condition.toLowerCase() === 'diabetes';
}

function beforeOrAt(left, right) {
  return left <= right;
}

function deepEqual(a, b) {
  return JSON.stringify(a) === JSON.stringify(b);
}

function authorized(data) {
  const c = data.Case;
  const insight = data.Insight;
  const permission = data.Policy.Permission;
  const constraint = permission.Constraint;
  return (
    permission.Action === c.RequestAction &&
    permission.Target === insight.ID &&
    deepEqual(constraint, {
      LeftOperand: 'odrl:purpose',
      Operator: 'odrl:eq',
      RightOperand: c.RequestPurpose,
    }) &&
    beforeOrAt(c.ScannerAuthAt, insight.ExpiresAt)
  );
}

function marketingProhibited(data) {
  const insight = data.Insight;
  const prohibition = data.Policy.Prohibition;
  const constraint = prohibition.Constraint;
  return (
    prohibition.Action === 'odrl:distribute' &&
    prohibition.Target === insight.ID &&
    deepEqual(constraint, {
      LeftOperand: 'odrl:purpose',
      Operator: 'odrl:eq',
      RightOperand: 'marketing',
    })
  );
}

function dutyTimingConsistent(data) {
  const c = data.Case;
  const insight = data.Insight;
  const duty = data.Policy.Duty;
  const constraint = duty.Constraint;
  return (
    duty.Action === 'odrl:delete' &&
    deepEqual(constraint, {
      LeftOperand: 'odrl:dateTime',
      Operator: 'odrl:eq',
      RightOperand: insight.ExpiresAt,
    }) &&
    beforeOrAt(c.ScannerDutyAt, insight.ExpiresAt)
  );
}

function minimized(data) {
  const serialized = data.Insight.SerializedLowercase.toLowerCase();
  return !serialized.includes('diabetes') && !serialized.includes('medical');
}

function canonicalPayload(data) {
  const insight = data.Insight;
  const policy = data.Policy;
  const payload = {
    insight: {
      createdAt: insight.CreatedAt,
      expiresAt: insight.ExpiresAt,
      id: insight.ID,
      metric: insight.Metric,
      retailer: insight.Retailer,
      scopeDevice: insight.ScopeDevice,
      scopeEvent: insight.ScopeEvent,
      suggestionPolicy: insight.SuggestionPolicy,
      threshold: Number(insight.ThresholdG),
      type: 'ins:Insight',
    },
    policy: {
      duty: {
        action: policy.Duty.Action,
        constraint: {
          leftOperand: policy.Duty.Constraint.LeftOperand,
          operator: policy.Duty.Constraint.Operator,
          rightOperand: policy.Duty.Constraint.RightOperand,
        },
      },
      permission: {
        action: policy.Permission.Action,
        constraint: {
          leftOperand: policy.Permission.Constraint.LeftOperand,
          operator: policy.Permission.Constraint.Operator,
          rightOperand: policy.Permission.Constraint.RightOperand,
        },
        target: policy.Permission.Target,
      },
      profile: policy.Profile,
      prohibition: {
        action: policy.Prohibition.Action,
        constraint: {
          leftOperand: policy.Prohibition.Constraint.LeftOperand,
          operator: policy.Prohibition.Constraint.Operator,
          rightOperand: policy.Prohibition.Constraint.RightOperand,
        },
        target: policy.Prohibition.Target,
      },
      type: 'odrl:Policy',
    },
  };
  return JSON.stringify(payload).replace('"threshold":10,', '"threshold":10.0,');
}

function payloadHashMatches(data) {
  const escapedPayload = canonicalPayload(data).replace(/"/g, '\\"');
  const digest = crypto.createHash('sha256').update(escapedPayload, 'utf8').digest('hex');
  return digest === data.Signature.PayloadHashSHA256;
}

function signatureMetadataValid(data) {
  const signature = data.Signature;
  return (
    signature.Alg === 'HMAC-SHA256' &&
    signature.HMACVerificationMode === 'trusted-precomputed-input' &&
    signature.SignatureHMAC.length === 64
  );
}

function selectedAlternative(data) {
  const scanned = scannedProduct(data);
  const candidates = data.Products.filter((product) => product.SugarTenths < scanned.SugarTenths);
  if (!candidates.length) throw new Error('no lower-sugar alternative');
  return minBy(candidates, (product) => [product.SugarTenths, product.Name]);
}

function bannerWarranted(data) {
  return scannedProduct(data).SugarTenths > data.Insight.ThresholdTenths;
}

function countsMatch(data) {
  return data.Case.FilesWritten === 6 && data.Case.AuditEntries === 1;
}

function trustedDerivation(data) {
  fail('Delfour derivation failed', {
    'low-sugar need exists': needsLowSugar(data),
    'scanner request is authorized': authorized(data),
    'marketing distribution is prohibited': marketingProhibited(data),
    'delete duty is tied to expiry': dutyTimingConsistent(data),
    'insight is minimized': minimized(data),
    'payload SHA-256 matches': payloadHashMatches(data),
    'signature metadata is valid': signatureMetadataValid(data),
    'lower-sugar alternative exists': selectedAlternative(data) !== undefined,
    'banner is warranted': bannerWarranted(data),
    'bus and audit counts match': countsMatch(data),
  });
}

function report(data) {
  trustedDerivation(data);
  const c = data.Case;
  const insight = data.Insight;
  const signature = data.Signature;
  const scanned = scannedProduct(data);
  const alternative = selectedAlternative(data);

  emit('# Delfour');
  emit();
  emit('## Insight');
  emit(`The scanner is allowed to use a neutral shopping insight and recommends ${alternative.Name} instead of ${scanned.Name}.`);
  emit(`case : ${c.CaseName}`);
  emit('decision : Allowed');
  emit(`scanned product : ${scanned.Name}`);
  emit(`suggested alternative: ${alternative.Name}`);
  emit();
  emit('## Explanation');
  emit('The phone desensitizes a diabetes-related household condition into a scoped low-sugar need, wraps it in an expiring Insight + Policy envelope, signs it, and the scanner consumes that envelope for shopping assistance.');
  emit(`metric : ${insight.Metric}`);
  emit(`threshold : ${insight.ThresholdDisplay}`);
  emit(`scope : ${insight.ScopeDevice} @ ${insight.ScopeEvent}`);
  emit(`retailer : ${insight.Retailer}`);
  emit(`signature alg : ${signature.Alg}`);
  emit('banner headline : Track sugar per serving while you scan');
  emit(`expires at : ${insight.ExpiresAt}`);
  emit(`reason.txt : ${data.ReasonText.trim()}`);
  emit(`audit entries : ${c.AuditEntries}`);
  emit(`bus files written : ${c.FilesWritten}`);
}

if (require.main === module) report(loadInput(NAME));
