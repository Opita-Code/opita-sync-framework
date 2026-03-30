# 03 — A.4 Approvals v1

## Principios base
- No modelar approvals como un flag, sino como un subsystem.
- Separar formalmente `authorization`, `approval` y `execution`.
- Toda aprobación es determinista, auditable, explicable, revocable y re-ejecutable.
- Toda aprobación queda invalidada si cambia algo material.
- Lo configurable se resuelve por policy/perfil/contexto, no libremente por ejecución.

## Invalidación por cambio material
Toda aprobación queda inválida si cambia alguno de estos campos materiales:
- plan
- datos
- conectores
- destino
- acción final
- clasificación
- risk score
- policy version

## Modos exactos de aprobación
- `auto`
  - no requiere aprobación humana
  - si policy permite, ejecuta y aplica
- `pre_execution`
  - requiere aprobación antes de correr
  - sin aprobación no ejecuta nada
- `pre_application`
  - puede analizar, planear o simular
  - no puede aplicar cambios hasta aprobación
  - es el modo por defecto para cambios reales tenant-scoped
- `double`
  - requiere dos aprobaciones según profile/policy/constraints
  - se usa para irreversible, global, cross-tenant, policy y connector sensitive

## SoD y autoridad de aprobación
### Regla de diseño
SoD y autoridad no quedan hardcoded; quedan configurables pero gobernados por:
- policy
- tipo de acción
- tenant
- contexto

### Variantes soportadas
- mismo aprobador permitido o prohibido
- aprobadores distintos obligatorios
- solicitante puede o no puede aprobar
- segundo aprobador de rol distinto
- segundo aprobador de función distinta
- doble aprobación obligatoria solo en ciertos casos

## approval_profile
Objeto reusable que define el comportamiento gobernado de approvals.

Campos sugeridos v1:
- `id`
- `name`
- `enabled`
- `scope`: global | tenant | capability | context
- `applies_when`
- `approval_mode`: auto | pre_execution | pre_application | double
- `required_approvals_count`
- `allowed_approver_types`
- `approver_conditions`
- `sod_constraints`
- `request_evidence_requirements`
- `decision_evidence_requirements`
- `timeout_policy`
- `fallback_policy`
- `revocation_rules`
- `invalidation_rules`
- `material_change_fields`
- `classification_constraints`
- `delegation_constraints`
- `execution_constraints`
- `policy_version_binding`
- `risk_policy_ref`
- `classification_policy_ref`
- `audit_policy_ref`
- `simulation_enabled`
- `replay_enabled`
- `version`
- `status`: draft | review | approved | published | retired

### allowed_approver_types sugeridos v1
- owner
- admin
- subadmin
- compliance
- designated_approver
- approver_group
- role_based
- external_governed

### sod_constraints sugeridos v1
- `requester_can_approve`
- `same_approver_can_approve_twice`
- `require_distinct_approvers`
- `require_distinct_roles`
- `require_distinct_functions`
- `forbid_requester_as_second_approver`

## approval_request
Captura una solicitud concreta con snapshots inmutables.

Campos sugeridos v1:
- `id`
- `tenant_id`
- `environment`
- `capability_id`
- `workflow_id`
- `execution_id`
- `intent_contract_id`
- `approval_profile_id`
- `approval_profile_version`
- `approval_mode`
- `request_status`: draft | pending_policy | pending_human | approved | rejected | expired | revoked | superseded | executed | applied
- `requested_by_subject_id`
- `requested_by_subject_type`
- `acting_for_subject_id`
- `delegation_id`
- `policy_decision_snapshot`
- `classification_snapshot`
- `risk_snapshot`
- `plan_snapshot`
- `data_scope_snapshot`
- `connectors_snapshot`
- `destination_snapshot`
- `action_snapshot`
- `material_change_fingerprint`
- `request_evidence`
- `request_reason`
- `human_approval_required`
- `required_approvals_count`
- `received_approvals_count`
- `approver_selection_policy`
- `candidate_approvers`
- `expires_at`
- `timeout_policy_snapshot`
- `fallback_policy_snapshot`
- `revocation_rules_snapshot`
- `invalidation_rules_snapshot`
- `created_at`
- `updated_at`

### request_evidence sugerida v1
- `summary`
- `goal`
- `expected_outcome`
- `affected_resources`
- `data_classes_involved`
- `external_effects`
- `estimated_risk`
- `diff_preview`
- `simulation_result`
- `policy_explanation`
- `attachments_refs`

## approval_decision
Guarda cada decisión política o humana.

Campos sugeridos v1:
- `id`
- `approval_request_id`
- `tenant_id`
- `environment`
- `decision_type`: policy | human | hybrid
- `decision_status`: approved | rejected | expired | revoked | superseded
- `decided_by_subject_id`
- `decided_by_subject_type`
- `decided_by_role_snapshot`
- `acting_for_subject_id`
- `delegation_id`
- `decision_reason_code`
- `decision_reason`
- `decision_evidence`
- `policy_version_snapshot`
- `classification_snapshot`
- `risk_snapshot`
- `material_change_fingerprint`
- `sod_evaluation_snapshot`
- `constraints_evaluation_snapshot`
- `effective_scope_snapshot`
- `decision_order`
- `is_final_decision`
- `supersedes_decision_id`
- `expires_at`
- `revoked_at`
- `created_at`

### decision_reason_code sugeridos v1
- `approved_within_policy`
- `approved_with_exception`
- `rejected_by_policy`
- `rejected_by_risk`
- `rejected_by_classification`
- `rejected_by_sod`
- `rejected_by_scope`
- `expired_no_response`
- `revoked_by_requester`
- `revoked_by_admin`
- `superseded_by_material_change`

## Matriz v1 riesgo × acción × contexto → approval_mode
### Variables efectivas
- `risk_level`: low | medium | high | critical
- `action_type`: read | analyze | simulate | configure | publish | execute | apply_external_change | manage_policy | manage_delegation | manage_connector | manage_capability
- `data_classification_max`: public | internal | confidential | restricted
- `external_effect`: none | reversible | irreversible
- `scope_size`: single_resource | bounded_set | broad_scope
- `actor_mode`: direct | delegated
- `environment`: dev | staging | prod

### Reglas base
- `auto`
  - read, analyze o simulate
  - sin efecto externo
  - riesgo low
  - clasificación hasta internal
  - scope acotado
- `pre_execution`
  - ejecución con riesgo medium o high
  - lectura/análisis ampliado sobre datos confidential o restricted
  - acciones delegadas con impacto relevante
  - gestión operativa con posible efecto indirecto
- `pre_application`
  - cuando existe cambio externo reversible o potencialmente sensible
  - publish, configure o execute con impacto controlable antes de aplicar
- `double`
  - cambios irreversibles
  - riesgo critical
  - clasificación restricted con efecto externo
  - policy, connector, delegation o capability management sensible
  - publicación global o cross-tenant

### Overrides duros
- Si `data_classification_max = restricted`, nunca bajar de `pre_execution`
- Si `external_effect = irreversible`, forzar `double`
- Si `action_type = manage_policy`, nunca bajar de `pre_application`
- Si `action_type = manage_connector`, nunca bajar de `pre_application`
- Si `scope_size = broad_scope` y `actor_mode = delegated`, subir un nivel
- Si `environment = prod`, subir un nivel para configure/publish/apply
- Si el cambio es global o cross-tenant, nunca usar `auto`

## Authority resolution de aprobadores válidos
La autoridad efectiva de aprobación se resuelve así:
- global policy default
- tenant approval policy override
- capability approval profile
- contextual constraints
- compiled effective approver set

### Fuentes de autoridad permitidas v1
- `owner`
- `tenant_admin`
- `tenant_subadmin`
- `compliance_officer`
- `designated_approver`
- `approver_group`
- `role_based_subject`
- `function_based_subject`
- `external_governed_approver`

### Resolution pipeline
1. Resolver universo elegible
2. Aplicar constraints del approval_profile
3. Aplicar elevaciones por contexto
4. Compilar set efectivo

### Constraints configurables recomendados
- `requester_can_approve`
- `same_approver_can_approve_twice`
- `require_distinct_subjects`
- `require_distinct_roles`
- `require_distinct_functions`
- `minimum_authority_level`
- `required_clearance_level`
- `allowed_approver_types`
- `max_delegation_depth_for_approver`
- `allow_external_governed_approver`
