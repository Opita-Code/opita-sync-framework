# A.5 — Tenant schema y onboarding v1

## Principios base

- El `tenant` es la frontera canónica de aislamiento administrativo, de datos, policy, approvals, memoria y operación del sistema.
- Existir en base de datos no alcanza: un `tenant` solo puede recibir tráfico real cuando su estado explícito es `operable`.
- `operable` requiere cumplimiento completo del mínimo duro; no admite interpretación laxa ni “bootstrap implícito”.
- Lo que varía por tenant se define al crearlo o en enriquecimientos posteriores gobernados; no se resuelve ad hoc por ejecución.
- El producto debe soportar `single_user` tenant sin exigir organigrama ni subadmins, pero sin relajar seguridad, approvals ni clasificación.
- El onboarding del tenant debe ser determinista, auditable y orientado a evidencia: cada paso material deja estado, snapshot y evento.
- Plataforma/superadmin gobierna el catálogo global, los conectores globales y los límites máximos; el tenant solo configura dentro de esos límites.
- El authoring humano de configuración declarativa usa YAML y el runtime opera sobre JSON canónico determinístico, consistente con A.1.
- Toda operación runtime que referencie tenancy debe seguir el contrato de A.2 con `tenant_id` y `environment`.
- Todo efecto del onboarding que condicione ejecución, approvals o exposición de catálogo debe ser explicable con artifacts y eventos.

---

## Definición de `tenant` como objeto canónico y rol en el sistema

`tenant` es el objeto canónico first-class que representa una instancia operativa aislada de Opyta Sync para una organización, unidad o usuario individual.

### Rol sistémico

- delimita aislamiento lógico y de governance
- ancla `subject`, `approval_profile`, `policy_artifact` tenant-scoped y `memory_record`
- consume `connector` y `capability` globales con habilitación efectiva por tenant
- define el contexto inicial para compilación de contratos, approvals y clasificación
- determina si el runtime puede o no aceptar tráfico real para ese scope

### Consecuencia normativa

Todo `intent_contract`, `approval_request`, `approval_decision`, `execution_record`, `result_record`, `memory_record` y evento canónico debe poder responder sin ambigüedad a qué `tenant_id` pertenece y en qué `environment` opera.

---

## Schema final de `tenant`

El `tenant` declarativo compila al envelope canónico de A.1:

```json
{
  "api_version": "v1",
  "kind": "tenant",
  "metadata": {},
  "spec": {},
  "status": {}
}
```

### 1. `metadata`

Campos canónicos recomendados:

- `id` — identificador estable en kebab-case del tenant
- `display_name` — nombre legible para humanos
- `owner_scope` — `platform`
- `tenant_mode` — `single_user` | `organization`
- `created_at`
- `updated_at`
- `labels`
- `tags`

`metadata` no sustituye la semántica operativa. Toda regla material vive en `spec` o en `status` observado.

### 2. `spec.identidad`

- `tenant_id` — igual a `metadata.id` en el canon compilado o derivable de forma determinista
- `display_name`
- `legal_entity_name` — opcional
- `tenant_slug` — opcional, estable para routing amigable
- `tenant_mode` — `single_user` | `organization`
- `single_user` — boolean explícito
- `primary_admin_subject_ref` — referencia al admin inicial obligatorio
- `initial_subject_refs` — lista inicial de subjects a crear o bindear
- `org_profile` — bloque descriptivo opcional de empresa/unidad
- `org_chart_present` — boolean observacional/declarativo

### 3. `spec.configuracion_operativa`

- `environment_bindings` — ambientes habilitados: `dev` | `staging` | `prod`
- `default_locale` — opcional
- `default_timezone` — opcional
- `default_currency` — opcional
- `allowed_regions` — opcional
- `runtime_defaults`
  - `default_classification_level`
  - `default_approval_mode_floor`
  - `default_result_retention_profile`
  - `default_memory_retention_profile`
- `traffic_enablement`
  - `accepts_runtime_traffic` — derivado de `operable`, no libre
  - `accepts_contract_creation` — derivado de estado
  - `accepts_execution` — derivado de estado

### 4. `spec.governance`

- `policy_bindings`
  - `base_policy_set_ref` — set base obligatorio
  - `tenant_policy_artifact_refs` — overrides permitidos dentro de límites
  - `platform_constraints_snapshot_ref`
- `approval_bootstrap`
  - `base_approval_profile_refs`
  - `default_approval_rules_ref`
  - `sod_profile_ref` — opcional si la base ya lo embebe
- `classification_bootstrap`
  - `classification_policy_ref`
  - `default_data_classification` — obligatorio
  - `output_control_policy_ref`
- `delegation_policy_ref` — opcional
- `governance_notes` — opcional y no material salvo decisión explícita

### 5. `spec.catalogo_conectores`

- `catalog_visibility`
  - `visible_capability_refs` — catálogo visible inicial
  - `hidden_by_default_capability_refs` — opcional
  - `catalog_policy_ref`
- `connector_enablement`
  - `enabled_connector_refs` — al menos uno en mínimo duro
  - `connector_limit_profile_ref`
  - `connector_credentials_binding_status` — snapshot/resumen, no secreto
  - `platform_managed_connector_set_ref`
- `capability_constraints`
  - `tenant_allowed_capability_refs`
  - `tenant_denied_capability_refs` — opcional dentro de límites permitidos

### 6. `spec.memoria_contexto`

- `bootstrap_context`
  - `organization_summary` — contexto inicial obligatorio
  - `business_domain` — opcional pero recomendado
  - `key_entities` — opcional
  - `known_systems` — opcional
  - `operational_constraints` — opcional
  - `classification_hints` — opcional
- `initial_memory_record_refs` — referencias a memoria inicial obligatoria o derivada
- `memory_policy_ref`
- `memory_bootstrap_status` — declarativo esperado

### 7. `spec.lifecycle_tecnico`

- `bootstrap_profile` — bloque canónico que define el onboarding inicial esperado
- `provisioning_requirements`
  - `requires_subject_seed`
  - `requires_policy_seed`
  - `requires_connector_enablement`
  - `requires_catalog_seed`
  - `requires_approval_seed`
  - `requires_classification_seed`
  - `requires_memory_seed`
- `operable_requirements_version` — `v1`
- `allowed_state_transitions_version` — `v1`
- `archive_policy_ref` — opcional

### 8. `status`

Campos observados mínimos:

- `lifecycle_state`
- `operable` — boolean observado; verdadero solo si `lifecycle_state = operable`
- `hard_minimum_complete`
- `validation_summary`
- `bootstrap_progress`
- `last_bootstrap_event_at`
- `last_validation_at`
- `suspension_reason_code` — si aplica
- `archive_reason_code` — si aplica
- `active_environment_refs`

---

## Campos obligatorios del tenant operable

Un tenant solo puede considerarse `operable` si, como mínimo, el canon compilado contiene y valida estos campos materiales:

- `metadata.id`
- `spec.identidad.display_name`
- `spec.identidad.single_user`
- `spec.identidad.primary_admin_subject_ref`
- `spec.configuracion_operativa.environment_bindings`
- `spec.governance.policy_bindings.base_policy_set_ref`
- `spec.governance.approval_bootstrap.base_approval_profile_refs`
- `spec.governance.classification_bootstrap.classification_policy_ref`
- `spec.governance.classification_bootstrap.default_data_classification`
- `spec.catalogo_conectores.catalog_visibility.visible_capability_refs`
- `spec.catalogo_conectores.connector_enablement.enabled_connector_refs`
- `spec.memoria_contexto.bootstrap_context.organization_summary` o `initial_memory_record_refs` equivalente resoluble
- `spec.memoria_contexto.memory_policy_ref`
- `spec.lifecycle_tecnico.bootstrap_profile`
- `status.lifecycle_state`
- `status.hard_minimum_complete`
- `status.operable`

Regla dura: si cualquiera de estos campos falta, está vacío cuando debe ser no vacío, o referencia un artefacto no resoluble, el tenant no puede entrar en `operable`.

---

## Campos opcionales / enriquecibles posterior al mínimo duro

Los siguientes campos pueden completarse después del mínimo duro sin impedir creación del tenant, siempre que no contradigan límites de plataforma:

- `legal_entity_name`
- `tenant_slug`
- `org_profile`
- `org_chart_present = true` con estructura formal posterior
- `subadmins` o subjects administrativos adicionales
- `delegation_policy_ref`
- `allowed_regions`
- `default_locale`, `default_timezone`, `default_currency`
- `hidden_by_default_capability_refs`
- `tenant_denied_capability_refs` permitidos por plataforma
- `business_domain`, `key_entities`, `known_systems`, `operational_constraints`
- perfiles avanzados de retención, analítica u observabilidad
- bindings de más conectores además del mínimo de uno
- enriquecimiento de memoria inicial con más `memory_record`

Regla: “opcional” no significa “arbitrario”; todo enriquecimiento posterior sigue validación, límites de policy y eventos canónicos.

---

## Regla formal del mínimo duro de creación

Se define el predicado canónico:

`tenant_hard_minimum_complete(tenant) = true` si y solo si se cumplen simultáneamente las condiciones siguientes:

1. existe exactamente un `primary_admin_subject_ref` resoluble y habilitado
2. existe un set de `policies` base resueltas y bindadas al tenant
3. existe al menos un `connector` global habilitado efectivamente para el tenant
4. existe un catálogo visible inicial no vacío y consistente con el set de conectores/policies
5. existen reglas de aprobación base resolubles y publicadas para el tenant
6. existe clasificación base resuelta y publicable para el tenant
7. existe memoria/contexto inicial resoluble suficiente para explicar a qué organización o usuario representa el tenant
8. todas las referencias anteriores respetan límites platform-governed y environment válido

### Corolarios normativos

- `tenant.created` no implica `tenant.operable`.
- persistencia física en DB no implica habilitación para tráfico real.
- `hard_minimum_complete = false` obliga a estado no operable.
- un tenant con admin pero sin conector, o con conector pero sin approvals base, sigue siendo no operable.

---

## Objeto `tenant_bootstrap_profile`

Se recomienda introducir `tenant_bootstrap_profile` como bloque declarativo embebido dentro de `tenant.spec.lifecycle_tecnico.bootstrap_profile`, no como objeto first-class separado en A.5.

### Justificación

- el objeto canónico first-class sigue siendo `tenant`
- el bootstrap inicial es material para operabilidad, pero no requiere identidad top-level separada en esta fase
- embebido evita ambigüedad entre “tenant deseado” y “plan de bootstrap efectivo”
- puede originarse desde una plantilla platform-governed, pero el canon final queda embebido y snapshotteado dentro del tenant

### Campos sugeridos de `bootstrap_profile`

- `profile_id` — identificador del perfil usado para bootstrap
- `profile_version`
- `profile_source` — `platform_template` | `custom_within_limits`
- `single_user_mode`
- `required_admin_count` — para v1 debe resolver a `1`
- `required_connector_min_count` — para v1 debe resolver a `1`
- `required_visible_capability_min_count` — para v1 debe ser `>= 1`
- `required_base_policy_refs`
- `required_base_approval_profile_refs`
- `required_base_classification_policy_refs`
- `required_initial_memory_rules`
- `org_chart_required` — para v1 debe poder ser `false`
- `subadmins_required` — para v1 debe ser `false`
- `completion_rule` — expresión declarativa equivalente al mínimo duro

### Regla de compilación

Si el bootstrap proviene de plantilla platform-governed, el runtime debe snapshotear el perfil resuelto dentro del JSON canónico del tenant para que la evidencia de onboarding no dependa de lookup mutable externo.

---

## Proceso exacto de onboarding del tenant

### Paso 1 — Recepción de solicitud de tenant

Se recibe una intención de creación con identidad mínima, `tenant_mode`, `single_user`, admin inicial, ambientes objetivo y perfil de bootstrap deseado.

Resultado esperado:

- validación sintáctica y de schema
- rechazo inmediato si faltan campos mínimos de creación
- emisión de `tenant.creation_requested`

### Paso 2 — Resolución de límites platform-governed

La plataforma resuelve:

- catálogo global publicable para ese tipo de tenant
- conectores globales elegibles
- set base de policies
- set base de approvals
- clasificación base permitida

Resultado esperado:

- se fijan límites máximos
- se emite `tenant.platform_constraints_resolved`

### Paso 3 — Compilación del `bootstrap_profile`

Se compila el perfil efectivo de onboarding, incluyendo reglas para single-user u organization.

Resultado esperado:

- `bootstrap_profile` queda materializado en el tenant
- se emite `tenant.bootstrap_profile_compiled`

### Paso 4 — Creación del registro canónico del tenant

Se crea el objeto `tenant` en estado `draft` o `provisioning`, pero todavía no habilitado para tráfico real.

Resultado esperado:

- existe `tenant_id`
- `status.operable = false`
- se emite `tenant.created`

### Paso 5 — Seed de identidad administrativa

Se crea o vincula el `primary_admin_subject_ref` y cualquier subject inicial requerido.

Resultado esperado:

- admin inicial resoluble
- se emite `tenant.admin_seeded`

### Paso 6 — Seed de governance base

Se bindan policies base, approvals base y clasificación base del tenant.

Resultado esperado:

- policies y profiles resolubles
- clasificación default fijada
- se emiten `tenant.base_policies_bound`, `tenant.base_approvals_bound`, `tenant.base_classification_bound`

### Paso 7 — Habilitación inicial de conectores

Se habilita al menos un conector global permitido para el tenant, sin exponer secretos en el canon del tenant.

Resultado esperado:

- `enabled_connector_refs` no vacío
- se emite `tenant.connectors_enabled`

### Paso 8 — Seed del catálogo visible inicial

Se publica al tenant un catálogo visible inicial compatible con sus conectores y policies.

Resultado esperado:

- `visible_capability_refs` no vacío
- ninguna capability visible contradice policies o connector enablement
- se emite `tenant.initial_catalog_published`

### Paso 9 — Seed de memoria/contexto inicial

Se crea o vincula la memoria/contexto inicial mínima de la empresa o del usuario único.

Resultado esperado:

- existe contexto inicial resoluble
- se emite `tenant.initial_memory_seeded`

### Paso 10 — Validación del mínimo duro

El sistema ejecuta validación formal sobre admin, policies, conectores, catálogo, approvals, clasificación y memoria.

Resultado esperado:

- si falla algo: `hard_minimum_complete = false`, se emite `tenant.operability_validation_failed`
- si todo pasa: `hard_minimum_complete = true`, se emite `tenant.hard_minimum_completed`

### Paso 11 — Promoción a `operable`

Solo si la validación anterior es satisfactoria, el tenant se promueve a `operable`.

Resultado esperado:

- `status.lifecycle_state = operable`
- `status.operable = true`
- `traffic_enablement` pasa a habilitado por derivación
- se emite `tenant.operable`

### Paso 12 — Inicio de operación auditada

Desde este momento el tenant puede recibir tráfico real de contratos y ejecución, sujeto a su environment y governance.

Resultado esperado:

- contratos A.2 con `tenant_id` y `environment` válidos pueden aceptarse
- antes de este punto, deben rechazarse

---

## Estados del lifecycle del tenant

### Estados canónicos

- `draft` — intención recibida o tenant compilado parcialmente; todavía sin provisión suficiente
- `provisioning` — se están creando bindings y artifacts base
- `bootstrap_pending` — el tenant existe, pero aún no completó el mínimo duro
- `operable` — tenant habilitado para tráfico real
- `suspended` — tenant temporalmente bloqueado para tráfico real por governance, riesgo o operación
- `archiving` — cierre técnico en progreso
- `archived` — tenant fuera de operación activa y no reactivable directamente sin proceso explícito de restore/recreate
- `failed` — provisioning/bootstrap falló de forma consistente y requiere intervención o reintento controlado

### Reglas operativas clave

- solo `operable` acepta tráfico real
- `draft`, `provisioning`, `bootstrap_pending` y `failed` no aceptan contratos operativos
- `suspended` conserva identidad y evidencia, pero bloquea tráfico real
- `archived` es terminal para operación normal
- `failed` no habilita uso aunque algunos bindings parciales existan

### Transiciones válidas

```text
draft -> provisioning
draft -> failed
provisioning -> bootstrap_pending
provisioning -> failed
bootstrap_pending -> operable
bootstrap_pending -> failed
operable -> suspended
operable -> archiving
suspended -> operable
suspended -> archiving
failed -> provisioning
archiving -> archived
```

### Transiciones prohibidas

- `draft -> operable` sin validación formal del mínimo duro
- `provisioning -> operable` salteando `bootstrap_pending` y su validación
- `bootstrap_pending -> suspended` como sustituto de validación fallida
- `failed -> operable` sin reprovisioning y revalidación completa
- `archived -> operable` sin proceso explícito externo de restauración o recreación
- `suspended -> draft`
- cualquier transición que marque `status.operable = true` fuera de `lifecycle_state = operable`

---

## Compatibilidad con single-user tenant

### Regla estructural

`single_user = true` debe existir como campo explícito en `spec.identidad`.

### Efectos permitidos de `single_user = true`

- permite que `tenant_mode` represente una sola persona o un micro-scope sin organigrama
- permite ausencia de `org_chart_present`
- permite ausencia de subadmins
- permite que el `primary_admin_subject_ref` sea también el único subject inicial

### Efectos prohibidos de `single_user = true`

- no reduce floors de approval
- no reduce clasificación base
- no elimina policies base
- no permite operar sin conector habilitado
- no autoriza bypass de evidencia auditable

### Regla de seguridad

Un tenant single-user puede tener menos estructura administrativa, pero NO menos governance. Esa distinción es central.

---

## Reglas sobre conectores, catálogo visible inicial y policies base en onboarding

### Conectores

- los conectores son globales y platform-governed
- el tenant solo recibe habilitación efectiva dentro de límites de plataforma
- el onboarding debe dejar al menos un `enabled_connector_ref`
- un conector visible en tenant debe estar tanto globalmente activo como tenant-enabled
- credenciales o secretos no viven en el schema canónico del tenant; solo referencias/snapshots de binding no sensibles

### Catálogo visible inicial

- el catálogo visible inicial es tenant-effective, pero su universo posible es platform-governed
- debe existir al menos una capability visible
- ninguna capability visible puede depender exclusivamente de conectores no habilitados
- plataforma puede imponer capabilities obligatorias, prohibidas o invisibles por posture/riesgo
- el tenant puede restringir adicionalmente dentro de límites, nunca ampliar fuera del set permitido

### Policies base

- deben existir antes de `operable`
- incluyen como mínimo policy base general, approvals base, clasificación base y memory policy base
- platform/superadmin define los mínimos y floors no relajables
- el tenant puede introducir configuración adicional solo dentro de límites declarados
- cualquier override tenant-configurable debe quedar snapshotteado y auditable

---

## Reglas de memoria/contexto inicial del tenant

La memoria/contexto inicial no es decorativa; forma parte del mínimo duro.

### Requisitos mínimos

Debe existir al menos una de estas dos formas, y ambas deben ser resolubles:

1. `bootstrap_context.organization_summary` suficientemente descriptivo, o
2. uno o más `initial_memory_record_refs` que representen contexto inicial equivalente

### Contenido mínimo esperado

- quién es la empresa, unidad o usuario representado
- qué sistemas o dominios principales existen, si son conocidos
- restricciones operativas relevantes para compilación o explicación
- hints básicos de clasificación o sensibilidad si ya se conocen

### Reglas duras

- no se permite `operable` con memoria/contexto inicial vacío
- la memoria inicial debe quedar gobernada por `memory_policy_ref`
- memoria inicial y telemetría no se confunden: la primera es contexto operativo, la segunda observabilidad append-only
- si el onboarding usa texto libre para contexto inicial, debe normalizarse a representación canónica resoluble

---

## Validaciones duras para considerar a un tenant `operable`

Un tenant es `operable` si y solo si todas estas validaciones pasan simultáneamente:

1. `lifecycle_state = operable`
2. `hard_minimum_complete = true`
3. `primary_admin_subject_ref` existe, pertenece al tenant y está habilitado
4. `base_policy_set_ref` resuelve a policies vigentes y compatibles con el environment
5. `base_approval_profile_refs` resuelven a approvals publicadas y aplicables
6. `classification_policy_ref` y `default_data_classification` son válidos
7. existe al menos un `enabled_connector_ref` globalmente permitido y tenant-enabled
8. existe al menos una `visible_capability_ref` compatible con policies y conectores
9. existe memoria/contexto inicial resoluble y gobernada
10. `environment_bindings` contiene al menos un environment válido entre `dev`, `staging`, `prod`
11. no existe suspensión activa ni reason code bloqueante
12. el bootstrap_profile compilado coincide con el estado efectivo observado
13. todos los checks dejan evidencia y evento de validación

Si una sola validación falla, `operable = false`.

---

## Eventos canónicos del lifecycle del tenant

### Eventos del ciclo principal

- `tenant.creation_requested`
- `tenant.platform_constraints_resolved`
- `tenant.bootstrap_profile_compiled`
- `tenant.created`
- `tenant.provisioning_started`
- `tenant.admin_seeded`
- `tenant.base_policies_bound`
- `tenant.base_approvals_bound`
- `tenant.base_classification_bound`
- `tenant.connectors_enabled`
- `tenant.initial_catalog_published`
- `tenant.initial_memory_seeded`
- `tenant.operability_validation_started`
- `tenant.operability_validation_failed`
- `tenant.hard_minimum_completed`
- `tenant.operable`
- `tenant.suspended`
- `tenant.reactivated`
- `tenant.archiving_started`
- `tenant.archived`
- `tenant.provisioning_failed`

### Eventos por causa específica

- `tenant.single_user_mode_confirmed`
- `tenant.org_chart_skipped`
- `tenant.subadmin_seed_skipped`
- `tenant.connector_enablement_failed`
- `tenant.catalog_validation_failed`
- `tenant.policy_binding_failed`
- `tenant.classification_binding_failed`
- `tenant.memory_seed_failed`
- `tenant.transition_rejected`
- `tenant.traffic_blocked`

---

## Payload mínimo sugerido por evento

- `event_id`
- `event_type`
- `tenant_id`
- `environment`
- `tenant_state`
- `from_state`
- `to_state`
- `operable`
- `hard_minimum_complete`
- `bootstrap_profile_id`
- `bootstrap_profile_version`
- `single_user`
- `primary_admin_subject_id`
- `enabled_connector_refs`
- `visible_capability_refs`
- `base_policy_set_ref`
- `base_approval_profile_refs`
- `classification_policy_ref`
- `memory_record_refs`
- `reason_code` — cuando aplique
- `validation_errors` — cuando aplique
- `triggered_by_subject_id`
- `occurred_at`

Regla: si el evento no aplica a algunos campos, pueden omitirse o enviarse nulos, pero `tenant_id`, `event_type`, `tenant_state` y `occurred_at` deben existir siempre.

---

## Tests borde mínimos

### T1 — Tenant creado sin admin

Dado un tenant persistido sin `primary_admin_subject_ref`.
Esperado: `hard_minimum_complete = false`; no entra en `operable`; `tenant.operability_validation_failed`.

### T2 — Tenant con admin pero sin policies base

Dado un tenant con admin inicial pero sin `base_policy_set_ref` resoluble.
Esperado: onboarding bloqueado; no acepta tráfico real.

### T3 — Tenant con policies base pero sin approval base

Dado un tenant sin `base_approval_profile_refs`.
Esperado: no puede considerarse operable aunque el resto exista.

### T4 — Tenant con approvals base pero sin clasificación base

Dado un tenant sin `classification_policy_ref` o `default_data_classification`.
Esperado: validación falla; no puede compilar contratos reales.

### T5 — Tenant con todo menos conectores

Dado un tenant con catálogo visible pero `enabled_connector_refs = []`.
Esperado: `tenant.connector_enablement_failed` o validación fallida; no `operable`.

### T6 — Tenant con conector habilitado pero catálogo inicial vacío

Dado un tenant con al menos un conector pero sin `visible_capability_refs`.
Esperado: onboarding incompleto; `hard_minimum_complete = false`.

### T7 — Capability visible incompatible con conectores habilitados

Dado un tenant cuya única capability visible requiere un conector no habilitado.
Esperado: `tenant.catalog_validation_failed`; no `operable`.

### T8 — Tenant sin memoria/contexto inicial

Dado un tenant sin `organization_summary` ni `initial_memory_record_refs`.
Esperado: validación falla; no entra a operación.

### T9 — Single-user tenant sin organigrama

Dado un tenant con `single_user = true` y sin organigrama.
Esperado: válido si el resto del mínimo duro está completo.

### T10 — Single-user tenant sin approvals base

Dado un tenant single-user que intenta omitir approvals base por ser de un solo usuario.
Esperado: rechazo; single-user no reduce governance.

### T11 — Tenant existente en DB intenta recibir tráfico en `bootstrap_pending`

Dado un contrato A.2 con `tenant_id` válido pero tenant en `bootstrap_pending`.
Esperado: rechazo de recepción; `tenant.traffic_blocked`.

### T12 — Transition directa `draft -> operable`

Dado un intento de transición que salta provisioning y validación.
Esperado: `tenant.transition_rejected`.

### T13 — Environment inválido

Dado un tenant con `environment_bindings` que incluye `qa`.
Esperado: schema inválido o validación fallida; solo `dev`, `staging`, `prod` son válidos.

### T14 — Override tenant fuera de límites de plataforma

Dado un tenant que intenta exponer capability no permitida por plataforma.
Esperado: rechazo de compilación o publicación tenant-effective.

### T15 — Conector global retirado después de operable y sin alternativas

Dado un tenant operable cuyo único conector habilitado deja de estar permitido globalmente.
Esperado: pasa a `suspended` o deja `operable = false` tras revalidación.

### T16 — Policies base cambian materialmente después del bootstrap

Dado un tenant operable con cambio material en set base obligatorio.
Esperado: revalidación obligatoria; si el estado efectivo ya no cumple mínimos, suspensión.

### T17 — Admin inicial deshabilitado

Dado un tenant operable cuyo `primary_admin_subject_ref` queda inválido y no existe reemplazo válido.
Esperado: suspensión o invalidación de operabilidad tras revalidación.

### T18 — Tenant archivado intenta reabrirse directamente

Dado un tenant en `archived` que intenta volver a `operable` por cambio directo de estado.
Esperado: transición prohibida.

---

## Criterios de aceptación del bloque tenant/onboarding

Se considera cerrado este bloque de A.5 si y solo si:

1. `tenant` queda definido como objeto canónico con schema final separado por identidad, operación, governance, catálogo/conectores, memoria/contexto y lifecycle técnico.
2. queda explícita la distinción entre tenant creado y tenant habilitado para tráfico real.
3. existe estado explícito `operable` y no se deriva solo de existencia en DB.
4. `operable` requiere cumplimiento completo del mínimo duro.
5. `single_user` queda explícito o inequívocamente derivable, y solo relaja organigrama/subadmins.
6. approvals, clasificación y seguridad NO se relajan por single-user.
7. conectores, catálogo y policies base quedan claramente separados entre platform-governed y tenant-configurable dentro de límites.
8. onboarding queda definido en pasos exactos y auditables desde creación hasta `operable`.
9. lifecycle del tenant tiene estados formales y transiciones válidas/prohibidas.
10. memoria/contexto inicial queda incluida como requisito del mínimo duro.
11. existen eventos canónicos del lifecycle con payload mínimo sugerido.
12. existen validaciones duras para determinar operabilidad.
13. existe una suite mínima de tests borde que cubre creación, bootstrap, single-user, límites de plataforma y revalidación.
14. el documento es consistente con A.1, A.2, A.3, A.4 y los source-of-truth de tenant model y tenant minimum hard.

---

## Decisión final de A.5 para tenant schema y onboarding v1

- `tenant` es el anchor canónico de tenancy y operación.
- un tenant puede existir sin estar operativo.
- `operable` es un estado explícito y auditado.
- el mínimo duro no es recomendación: es condición necesaria de operabilidad.
- `single_user` cambia estructura administrativa mínima, no floors de governance.
- onboarding, validación y transición a tráfico real deben dejar evidencia material y eventos canónicos.
