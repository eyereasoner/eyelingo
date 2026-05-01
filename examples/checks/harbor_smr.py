# Independent Python checks for the harbor_smr example.
from __future__ import annotations

import re
from datetime import datetime

from .common import run_checks


def dt(value: str) -> datetime:
    return datetime.fromisoformat(value.replace("Z", "+00:00"))


def run(ctx):
    data = ctx.load_input()
    agg = data["Aggregate"]
    insight = data["Insight"]
    req = data["Request"]
    dispatch = data["Dispatch"]
    policy = data["Policy"]
    thresholds = data["Thresholds"]
    auth_time = dt(data["HubAuthAt"])
    expires = dt(insight["ExpiresAt"])
    window_start = dt(dispatch["WindowStart"])
    window_end = dt(dispatch["WindowEnd"])
    duration_h = (window_end - window_start).total_seconds() / 3600.0
    energy_mwh = dispatch["DispatchMW"] * duration_h
    sensitive_terms = ["coretemperature", "rodposition", "neutronflux", "operatorbadge"]
    serialized = insight["SerializedLowercase"].lower()

    permit = (
        data["RequestPurpose"] == req["Purpose"] == policy["Permission"]["Purpose"]
        and data["RequestAction"] == policy["Permission"]["Action"]
        and policy["Permission"]["Target"] == insight["ID"]
        and req["TargetLoad"] == insight["TargetLoad"] == dispatch["ForLoad"]
        and req["RequestedMW"] == dispatch["DispatchMW"]
        and req["RequestedMW"] <= insight["ExportMW"] <= agg["AvailableFlexibleExportMW"]
        and agg["ReserveMarginMW"] >= thresholds["MinReserveMarginMW"]
        and agg["CoolingMarginPct"] >= thresholds["MinCoolingMarginPct"]
        and not agg["PlannedOutage"]
        and auth_time < expires
    )

    checks = [
        ("reserve margin exceeds the configured threshold", agg["ReserveMarginMW"] >= thresholds["MinReserveMarginMW"]),
        ("cooling margin exceeds the configured threshold", agg["CoolingMarginPct"] >= thresholds["MinCoolingMarginPct"]),
        ("no planned outage blocks the balancing window", agg["PlannedOutage"] is False),
        ("requested dispatch fits inside the flexible-export insight", req["RequestedMW"] <= insight["ExportMW"] <= agg["AvailableFlexibleExportMW"]),
        ("serialized insight omits sensitive reactor telemetry terms", not any(term in serialized for term in sensitive_terms)),
        ("aggregate flags keep raw reactor telemetry local", not any(agg[k] for k in ["ContainsCoreTemperature", "ContainsRodPosition", "ContainsNeutronFlux", "ContainsOperatorBadgeIDs"])),
        ("permission policy authorizes electrolysis dispatch before expiry", permit and "PERMIT" in ctx.answer),
        ("market-resale redistribution is separately prohibited", policy["Prohibition"]["Action"] == "odrl:distribute" and policy["Prohibition"]["Purpose"] == "market_resale"),
        ("scope is explicit for device, event, start, and expiry", all(insight.get(k) for k in ["ScopeDevice", "ScopeEvent", "WindowStart", "ExpiresAt"])),
        ("dispatch energy recomputes to 64 MWh over the four-hour window", energy_mwh == 64 and "64 MWh" in ctx.reason),
        ("the reported load, power, and window match the request and dispatch", req["TargetLoad"] in ctx.answer and f"at {req['RequestedMW']} MW" in ctx.answer and dispatch["WindowStart"] in ctx.answer and dispatch["WindowEnd"] in ctx.answer),
    ]
    return run_checks(checks)
