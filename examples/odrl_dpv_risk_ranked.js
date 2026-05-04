#!/usr/bin/env node
const { emit, fail, loadInput } = require('./_see');

const NAME = 'odrl_dpv_risk_ranked';

function byId(items) {
  return new Map(items.map((item) => [item.ID, item]));
}

function rulesByAction(data) {
  const map = new Map();
  for (const rule of data.Agreement.Rules) {
    if (!map.has(rule.Action)) map.set(rule.Action, []);
    map.get(rule.Action).push(rule);
  }
  return map;
}

function hasDuty(rule, action) {
  return rule.Duties.some((duty) => duty.Action === action);
}

function constraintValue(rule, leftOperand) {
  const constraint = rule.Constraints.find((c) => c.LeftOperand === leftOperand);
  return constraint ? constraint.RightOperand : undefined;
}

function dutyConstraintValue(rule, dutyAction, leftOperand) {
  for (const duty of rule.Duties) {
    if (duty.Action !== dutyAction) continue;
    const constraint = duty.Constraints.find((c) => c.LeftOperand === leftOperand);
    if (constraint) return constraint.RightOperand;
  }
  return undefined;
}

function consentIsExplicit(rule) {
  return rule.Constraints.some((constraint) => (
    constraint.LeftOperand === 'consent'
    && constraint.Operator === 'eq'
    && constraint.RightOperand === true
  ));
}

function scoreRisk(base, importance, scoring) {
  return Math.min(scoring.MaximumScore, base + importance);
}

function classify(score, scoring) {
  if (score >= scoring.HighAtLeast) {
    return { Level: 'HighRisk', Severity: 'HighSeverity' };
  }
  if (score >= scoring.ModerateAtLeast) {
    return { Level: 'ModerateRisk', Severity: 'ModerateSeverity' };
  }
  return { Level: 'LowRisk', Severity: 'LowSeverity' };
}

function makeRisk(data, rule, need, base, kind, description, consequences, impacts, mitigations) {
  const score = scoreRisk(base, need.Importance, data.Scoring);
  const classification = classify(score, data.Scoring);
  const clause = byId(data.Agreement.Clauses).get(rule.Clause);
  return {
    Rule: rule.ID,
    Clause: clause.ID,
    ClauseId: clause.ClauseId,
    ClauseText: clause.Text,
    Need: need.ID,
    Kind: kind,
    Score: score,
    Level: classification.Level,
    Severity: classification.Severity,
    Consequences: consequences.slice(),
    Impacts: impacts.slice(),
    Description: description(clause),
    Mitigations: mitigations.slice(),
  };
}

function deriveRisks(data) {
  const needs = byId(data.ConsumerProfile.Needs);
  const actionRules = rulesByAction(data);
  const risks = [];

  const deleteRules = actionRules.get('removeAccount') || [];
  for (const rule of deleteRules) {
    const missingNotice = constraintValue(rule, 'noticeDays') === undefined;
    const missingInform = !hasDuty(rule, 'inform');
    if (rule.Type === 'Permission' && missingNotice && missingInform) {
      risks.push(makeRisk(
        data,
        rule,
        needs.get('Need_DataCannotBeRemoved'),
        90,
        'UnwantedDataDeletion',
        (clause) => `account/data removal is permitted without notice safeguards. Clause ${clause.ClauseId}: ${clause.Text}`,
        ['DataLoss', 'DataUnavailable', 'CustomerConfidenceLoss'],
        ['FinancialLoss', 'NonMaterialDamage'],
        [
          'Add a notice constraint (minimum noticeDays) before account removal.',
          'Add a duty to inform the consumer prior to account removal.',
        ],
      ));
    }
  }

  const changeRules = actionRules.get('changeTerms') || [];
  for (const rule of changeRules) {
    const need = needs.get('Need_ChangeOnlyWithPriorNotice');
    const days = dutyConstraintValue(rule, 'inform', 'noticeDays');
    if (rule.Type === 'Permission' && typeof days === 'number' && days < need.MinNoticeDays) {
      risks.push(makeRisk(
        data,
        rule,
        need,
        70,
        'PolicyRisk',
        (clause) => `terms may change with notice (${days} days) below consumer requirement (${need.MinNoticeDays} days). Clause ${clause.ClauseId}: ${clause.Text}`,
        ['CustomerConfidenceLoss'],
        ['NonMaterialDamage'],
        ['Increase minimum noticeDays in the inform duty to meet the consumer requirement.'],
      ));
    }
  }

  const shareRules = actionRules.get('shareData') || [];
  for (const rule of shareRules) {
    if (rule.Type === 'Permission' && !consentIsExplicit(rule)) {
      risks.push(makeRisk(
        data,
        rule,
        needs.get('Need_NoSharingWithoutConsent'),
        85,
        'UnwantedDisclosureData',
        (clause) => `user data sharing is permitted without an explicit consent constraint. Clause ${clause.ClauseId}: ${clause.Text}`,
        ['CustomerConfidenceLoss'],
        ['NonMaterialDamage', 'FinancialLoss'],
        ['Add an explicit consent constraint before data sharing.'],
      ));
    }
  }

  const exportRules = actionRules.get('exportData') || [];
  for (const rule of exportRules) {
    if (rule.Type === 'Prohibition') {
      risks.push(makeRisk(
        data,
        rule,
        needs.get('Need_DataPortability'),
        60,
        'PolicyRisk',
        (clause) => `portability is restricted because exporting user data is prohibited. Clause ${clause.ClauseId}: ${clause.Text}`,
        ['CustomerConfidenceLoss'],
        ['NonMaterialDamage'],
        ['Add a permission allowing data export, or remove the prohibition, to support portability.'],
      ));
    }
  }

  return risks.sort((a, b) => b.Score - a.Score || a.ClauseId.localeCompare(b.ClauseId));
}

function trustedDerivation(data, risks) {
  const expected = data.Expected;
  const scoreByClause = Object.fromEntries(risks.map((risk) => [risk.ClauseId, risk.Score]));
  const levelByClause = Object.fromEntries(risks.map((risk) => [risk.ClauseId, risk.Level]));
  const mitigationCountByClause = Object.fromEntries(risks.map((risk) => [risk.ClauseId, risk.Mitigations.length]));

  fail('ODRL DPV risk ranking derivation failed', {
    'all expected clauses produced a risk': JSON.stringify(risks.map((risk) => risk.ClauseId)) === JSON.stringify(expected.RankedClauses),
    'scores match expected capped values': JSON.stringify(scoreByClause) === JSON.stringify(expected.ScoresByClause),
    'risk levels match expected thresholds': JSON.stringify(levelByClause) === JSON.stringify(expected.LevelsByClause),
    'mitigation counts match expected': JSON.stringify(mitigationCountByClause) === JSON.stringify(expected.MitigationCountsByClause),
    'highest risk is first': risks[0].ClauseId === expected.TopClause,
    'every risk keeps a clause link': risks.every((risk) => risk.Clause && risk.ClauseText.length > 0),
    'every risk has consequences and impacts': risks.every((risk) => risk.Consequences.length > 0 && risk.Impacts.length > 0),
    'rank is descending by score': risks.every((risk, i) => i === 0 || risks[i - 1].Score >= risk.Score),
  });
}

function main() {
  const data = loadInput(NAME);
  const risks = deriveRisks(data);
  trustedDerivation(data, risks);

  emit('# ODRL + DPV Risk Ranking');
  emit();
  emit('## Insight');
  emit(`agreement : ${data.Agreement.Title}`);
  emit(`profile : ${data.ConsumerProfile.Title}`);
  emit(`risks found : ${risks.length}`);
  emit(`highest risk : clause ${risks[0].ClauseId} score=${risks[0].Score} ${risks[0].Level}`);
  emit();
  emit('## Explanation');
  emit('The example treats the agreement as a small policy graph: permissions describe what the provider may do, prohibitions describe what users may not do, duties describe required follow-up actions, and constraints describe safeguards such as notice days or consent.');
  emit('Each rule is compared with the consumer profile. When a required safeguard is missing or too weak, the program creates a DPV-style risk record with a score, severity, risk level, consequences, impacts, and mitigation advice.');
  emit('Scores are base risk plus consumer-need importance, capped at 100, then ranked from highest to lowest so the most urgent clause appears first.');
  emit();
  emit('Ranked risk report:');
  for (const risk of risks) {
    emit(`score=${risk.Score} (${risk.Level}, ${risk.Severity}) clause ${risk.ClauseId}`);
    emit(` ${risk.Description}`);
    emit(` consequences=${risk.Consequences.join(', ')} impacts=${risk.Impacts.join(', ')}`);
    for (const mitigation of risk.Mitigations) {
      emit(` - mitigation: ${mitigation}`);
    }
  }
}

if (require.main === module) main();
