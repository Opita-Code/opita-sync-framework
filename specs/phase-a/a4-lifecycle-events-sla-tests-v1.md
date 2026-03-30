# 05 — A.4 Lifecycle, events, SLA, fallback, tests

## Estados del ciclo de approvals
- `draft`
- `pending_policy`
- `pending_human`
- `partially_approved`
- `approved`
- `rejected`
- `expired`
- `revoked`
- `superseded`
- `execution_released`
- `application_released`
- `executed`
- `applied`
- `closed`

## Reglas operativas clave
- `approved` no significa siempre que ya puede aplicar; depende de `approval_mode`.
- `execution_released` y `application_released` deben quedar separados.
- `partially_approved` solo existe cuando hay más de una aprobación requerida.
- `superseded` gana sobre estados pendientes si ocurre cambio material.
- Todo estado terminal debe terminar en `closed`.

## Transiciones prohibidas importantes
- `draft -> applied`
- `pending_human -> applied`
- `rejected -> approved`
- `expired -> approved`
- `revoked -> approved`
- reutilizar `superseded` como request activo
- pasar a `application_released` si el approval efectivo ya no está vigente

## Eventos canónicos finales
- `approval.request_created`
- `approval.request_submitted`
- `approval.policy_evaluation_started`
- `approval.policy_evaluation_completed`
- `approval.human_review_requested`
- `approval.human_review_received`
- `approval.partially_approved`
- `approval.approved`
- `approval.rejected`
- `approval.expired`
- `approval.revoked`
- `approval.superseded`
- `approval.execution_released`
- `approval.application_released`
- `approval.execution_blocked`
- `approval.application_blocked`
- `approval.closed`

### Eventos por causa
- `approval.material_change_detected`
- `approval.approver_resolution_completed`
- `approval.fallback_triggered`
- `approval.timeout_warning_emitted`
- `approval.timeout_reached`
- `approval.sod_violation_detected`
- `approval.policy_floor_applied`
- `approval.effective_risk_computed`

### Payload mínimo sugerido por evento
- `event_id`
- `event_type`
- `tenant_id`
- `environment`
- `approval_request_id`
- `execution_id`
- `intent_contract_id`
- `approval_profile_id`
- `approval_mode`
- `effective_risk_level`
- `from_state`
- `to_state`
- `reason_code`
- `reason`
- `triggered_by_subject_id`
- `policy_version`
- `occurred_at`

## SLA y fallback operativos
### Principios
- El SLA no es libre por ejecución; sale de policy/perfil.
- El fallback no puede reducir seguridad ni saltarse floors duros.
- Todo timeout o escalación deja evento y evidencia.

### SLA por defecto recomendado v1
#### low
- `auto`: inmediato
- `pre_execution`: 4 horas
- `pre_application`: 8 horas
- `double`: 12 horas

#### medium
- `auto`: inmediato
- `pre_execution`: 2 horas
- `pre_application`: 4 horas
- `double`: 8 horas

#### high
- `auto`: no aplica normalmente
- `pre_execution`: 1 hora
- `pre_application`: 2 horas
- `double`: 4 horas

#### critical
- `auto`: prohibido
- `pre_execution`: 30 minutos
- `pre_application`: 1 hora
- `double`: 2 horas

### Reminder schedule recomendado
- primer reminder al 50% del SLA
- segundo reminder al 80%
- último warning al 95%

### Timeout policy recomendada
- `low/medium`: `expire` o `fallback_profile` según capability
- `high`: `expire` o `escalate`, nunca autoaprobar
- `critical`: `escalate` o `reject`, nunca autoaprobar

### Fallback policy recomendada
- si no responde aprobador individual -> escalar a grupo aprobador
- si no responde grupo aprobador -> escalar a autoridad superior permitida
- si policy no permite escalar -> `expired` o `rejected`
- nunca degradar a un aprobador con menor autoridad
- nunca saltar constraints SoD ni clasificación

### Defaults específicos recomendados para Opyta Sync
- cambios de policy, connectors, credenciales, publicación global o cross-tenant:
  - sin respuesta -> `escalate`
  - segundo timeout -> `reject`
- cambios reversibles tenant-scoped:
  - sin respuesta -> `fallback_profile` o `expire`
- acciones delegated + broad scope:
  - sin respuesta -> `expire`

## Tests borde mínimos
- cambio material después de aprobar -> `superseded`
- mismo request con dos decisiones incompatibles -> compilar solo válidas
- `double` con mismo aprobador cuando SoD lo prohíbe -> rechazo o invalidez
- aprobador sin clearance suficiente -> rechazo
- tenant correcto pero role binding vencido -> no elegible
- `restricted + external exposure` -> floor mínimo `high`
- `cross_tenant` -> riesgo `critical`
- timeout sin fallback permitido -> `expired`
- timeout con fallback permitido -> `fallback_triggered`
- aprobación vigente revocada antes de aplicar -> bloquear aplicación
- policy version cambia después de aprobar -> `superseded`
- delegated + broad scope -> endurecimiento de riesgo y autoridad
- `pre_application` permite ejecución parcial pero bloquea aplicación sin approval
- `auto` nunca para irreversible/global/cross-tenant
- `manage_policy` o `manage_connector` nunca por debajo de `pre_application`
- `critical` nunca autoaprueba por timeout

## Criterios de aceptación de A.4
- toda solicitud puede resolver `approval_mode` sin ambigüedad
- todo request puede calcular `business_risk_score`, `security_risk_score` y `effective_risk_level`
- todo request puede resolver aprobadores válidos por policy compilada
- toda aprobación queda invalidada por cambio material
- todo cambio de estado emite evento canónico
- el motor separa `execution_released` de `application_released`
- el motor soporta `auto`, `pre_execution`, `pre_application` y `double`
- el motor soporta SoD configurable gobernado
- el motor soporta fallback y timeout sin romper constraints
- toda decisión deja snapshot y evidencia auditable
- no existe camino no auditado hacia `applied`
