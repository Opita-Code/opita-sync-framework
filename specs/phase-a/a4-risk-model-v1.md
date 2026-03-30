# 04 — A.4 Risk model v1

## Principio
El `risk_level` no sale de un enum fijo ni de intuición humana. Se compila por request usando score dinámico + floors duros.

## Dos scores separados
El motor calcula primero:
- `business_risk_score`
- `security_risk_score`

Y luego deriva:
- `effective_risk_level`
- `effective_risk_reasoning`

## business_risk_score
### Variables exactas recomendadas v1
- `process_criticality`
- `external_effect`
- `scope_blast_radius`
- `financial_impact`
- `operational_impact`
- `reputational_impact`
- `tenant_impact`
- `rollback_feasibility`
- `execution_confidence`
- `change_novelty`

### Fórmula sugerida v1
```text
business_risk_score = clamp(0,100,
 process_criticality +
 external_effect +
 scope_blast_radius +
 financial_impact +
 operational_impact +
 reputational_impact +
 tenant_impact +
 rollback_feasibility +
 execution_confidence +
 change_novelty
)
```

### Buckets
- `0-24` -> low
- `25-49` -> medium
- `50-74` -> high
- `75-100` -> critical

### Floors duros recomendados
- `external_effect = irreversible` -> mínimo `high`
- `scope_blast_radius = cross_tenant` -> mínimo `critical`
- `tenant_impact = global_platform` -> mínimo `critical`
- `process_criticality = mission_critical` + `operational_impact >= high` -> mínimo `high`
- `financial_impact = severe` + `rollback_feasibility = none` -> mínimo `high`

## security_risk_score
### Variables exactas recomendadas v1
- `data_classification_max`
- `connector_sensitivity`
- `credential_permission_sensitivity`
- `delegation_risk`
- `auth_session_risk`
- `actor_posture_risk`
- `behavior_anomaly_risk`
- `policy_sensitivity`
- `external_exposure_risk`
- `lateral_abuse_escalation_risk`

### Fórmula sugerida v1
```text
security_risk_score = clamp(0,100,
 data_classification_max +
 connector_sensitivity +
 credential_permission_sensitivity +
 delegation_risk +
 auth_session_risk +
 actor_posture_risk +
 behavior_anomaly_risk +
 policy_sensitivity +
 external_exposure_risk +
 lateral_abuse_escalation_risk
)
```

### Buckets
- `0-24` -> low
- `25-49` -> medium
- `50-74` -> high
- `75-100` -> critical

### Floors duros recomendados
- `data_classification_max = restricted` + `external_exposure_risk >= medium` -> mínimo `high`
- `credential_permission_sensitivity = critical` -> mínimo `high`
- `policy_sensitivity = critical` -> mínimo `high`
- `auth_session_risk = critical` -> mínimo `critical`
- `actor_posture_risk = unknown_or_untrusted` + cambio real -> mínimo `high`
- `lateral_abuse_escalation_risk = critical` -> mínimo `critical`
- `connector_sensitivity = critical` + efecto externo -> mínimo `high`

## effective_risk_level
### Fórmula operativa recomendada
- calcular `business_bucket`
- calcular `security_bucket`
- tomar el mayor como base
- aplicar floors duros globales
- elevar a `critical` cuando ambos estén en `high` y exista efecto externo o alcance amplio

```text
effective_risk_level = max(business_bucket, security_bucket, hard_floors)
```

Regla extra:
```text
high + high + (external_effect != none or scope_blast_radius in [broad_scope, cross_tenant]) -> critical
```

### hard_floors globales sugeridos
- irreversible -> `critical`
- cross_tenant -> `critical`
- restricted + external exposure -> mínimo `high`
- manage_policy sensible -> mínimo `high`
- manage_connector sensible -> mínimo `high`
- auth/session critical -> `critical`

## Recomendación específica para Opyta Sync
- Pesar fuerte `external_effect`, `scope_blast_radius` y `tenant_impact`.
- Pesar fuerte `data_classification_max`, `credential_permission_sensitivity` y `lateral_abuse_escalation_risk`.
- No confiar solo en rol; considerar sesión, postura, anomalías y delegación.
- Mantener `capability_risk_base` en catálogo y recalcular el resto en tiempo real.
- Permitir override por tenant solo dentro de límites del superadmin.
- No permitir que un override reduzca floors duros.
