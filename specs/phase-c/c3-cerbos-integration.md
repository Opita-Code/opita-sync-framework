# C.3 — Cerbos integration

## Objetivo

Definir la integración implementable entre el kernel de Opyta Sync y **Cerbos** para que el kernel actúe como **Policy Enforcement Point (PEP)**, Cerbos actúe como **Policy Decision Point (PDP)**, y toda decisión contextual de autorización/governance sensible se evalúe sobre un input canonizado, versionado, auditable y consistente con A.2, A.4, B.2, C.1 y C.2.

Este bloque NO implementa policies ni código productivo. Cierra el diseño operativo necesario para que el runtime, approvals y clasificación consulten policy sin inventar reglas ad hoc dentro del kernel.

## Principios de implementación de la integración

- El **kernel es el PEP**. Cerbos es el **PDP**.
- Cerbos **no decide el lifecycle del runtime**; decide autorización/policy contextual sobre un input canonizado.
- El kernel **no envía payloads arbitrarios** a Cerbos; envía shapes canonizados, estables y versionados.
- La policy debe evaluarse sobre objetos del dominio ya compilados o proyectados; no sobre texto libre o estado implícito.
- El kernel debe **traducir el output de Cerbos a transiciones conocidas** del runtime/approval/classification; no puede reinterpretarlo libremente por cada caller.
- El enforcement debe ser **tenant-scoped**, correlacionable y auditable.
- Toda decisión material de policy debe dejar evidencia durable mediante `policy_decision_record` o equivalente embebido/referenciado.
- La integración debe ser compatible con `compiled_contract`, `execution_record`, `approval_request`, clasificación y resolution de capabilities sin mover los boundaries ya fijados.
- El failure mode por defecto es **fail closed** para mutaciones y operaciones sensibles.
- La integración debe poder crecer en complejidad sin transformar al kernel en una mini-plataforma de policy propia.

## Boundary exacto PEP / PDP

### Qué hace el PEP del kernel

El PEP:

1. determina si el lifecycle actual exige consulta de policy;
2. toma el objeto de dominio relevante (`compiled_contract`, `execution_record`, `approval_request`, `capability`, `policy_artifact`);
3. construye un `policy_input_v1` canonizado y versionado;
4. invoca a Cerbos vía adapter estable;
5. recibe respuesta del PDP;
6. valida que la respuesta sea mapeable;
7. la traduce a consecuencias conocidas del kernel (`allow`, `deny_block`, `require_approval`, `require_escalation`, `restricted_view`);
8. persiste evidencia y correlación;
9. entrega una decisión ya interpretada al runtime o subsistema consumidor.

### Qué hace el PDP en Cerbos

Cerbos:

1. evalúa principal + resource + action + context + atributos canonizados;
2. aplica policy versionada tenant-scoped/global según corresponda;
3. devuelve decisión de autorización y metadata útil para explicación/mapping;
4. no decide estados de workflow, tiempos de approval, retries, compensación ni lifecycle operativo.

### Regla de boundary

- **PEP = orquesta enforcement y traduce decisión al dominio.**
- **PDP = evalúa policy y emite decisión contextual.**
- Si una regla requiere conocer transición de workflow, estado durable o invariantes de A.5/C.2, esa decisión queda en el kernel.
- Si una regla requiere evaluar permiso contextual, clearance, tenant scope, clasificación, delegación, riesgo o governance mode, la evaluación queda en Cerbos.

## Responsabilidades explícitas del PEP del kernel

- Resolver el `resource_kind` correcto para cada punto del lifecycle.
- Resolver la `action` exacta y no ambigua que se está consultando.
- Garantizar que `tenant_id` y correlación mínima estén presentes antes de consultar policy.
- Canonizar `principal`, `resource` y `context` en `policy_input_v1`.
- Adjuntar `policy_input_schema_version` y `policy_reference_version` / `policy_version` relevante.
- Invocar a Cerbos mediante `cerbos_client_adapter` y manejar timeouts/fallas.
- Validar semánticamente la respuesta del PDP.
- Aplicar `policy_decision_mapper` para traducir salida PDP a decisión del kernel.
- Persistir `policy_decision_record` con evidencia suficiente.
- Aplicar `policy_cache_guard` cuando exista caché permitida.
- Rechazar respuestas no mapeables o inconsistentes con el estado local.
- Re-evaluar policy cuando cambie un input material como clasificación, delegación, approval mode efectivo o policy version.

## No-responsabilidades explícitas del PEP del kernel

- No escribir lógica de autorización ad hoc fuera del mapping formal.
- No enviar blobs arbitrarios, prompts, payloads libres o estructuras no versionadas a Cerbos.
- No pedirle a Cerbos que decida estados como `executing`, `applying`, `failed`, `unknown_outcome` o `compensated`.
- No permitir que Cerbos sustituya las validaciones estructurales del contrato o del runtime.
- No decidir autoridad humana final de approvals por fuera del subsystem de A.4.
- No degradar silenciosamente a `allow` cuando Cerbos no responde.
- No cachear decisiones ignorando `policy_version`, clasificación, delegación o fingerprints relevantes.
- No reabrir la comparativa B.2 ni cambiar el baseline Cerbos.

## Qué decisiones siguen fuera de Cerbos

Las siguientes decisiones permanecen en el kernel o en otros subsistemas, aunque puedan usar policy como input:

- validez estructural de `compiled_contract` según A.2/C.1;
- creación idempotente de `execution_record` y arranque de `execution_workflow` según C.2;
- transición exacta entre `execution_completed` y `application_completed`;
- retries técnicos, timers, `unknown_outcome`, compensación y cierre manual;
- construcción y lifecycle interno de `approval_request`/`approval_decision`;
- clasificación compilada del contrato como objeto canónico;
- detección de cambio material que invalida approvals;
- persistencia canónica en PostgreSQL y proyección de workflow;
- validaciones de ejecutabilidad del contrato y presence de `plan_snapshot`/approval decision;
- resolución final de cómo una decisión `require_approval` se convierte en `approval_request` concreto y quiénes integran el set efectivo de aprobadores.

## Puntos exactos del lifecycle donde el kernel debe consultar policy

Las consultas de policy obligatorias en v1 son:

1. **Durante compilación del contrato (`compile`)**
   - objetivo: validar autorización contextual para compilar/intentar materializar el contrato;
   - resource principal: `intent_contract`.

2. **Antes de liberar ejecución técnica (`execute`)**
   - ocurre en `preconditions_validating -> policy_evaluating` del `execution_workflow`;
   - resource principal: `execution_record` o `intent_contract` referenciado;
   - resultado posible: permitir, bloquear o derivar a approval/escalation.

3. **Antes de `release_execution`**
   - cuando llega señal externa para sacar la ejecución de `awaiting_approval`;
   - resource principal: `approval_request` + `execution_record` correlacionados;
   - verifica que la liberación siga autorizada bajo policy vigente.

4. **Antes de `release_application`**
   - después de `execution_completed` y antes de habilitar `applying`;
   - resource principal: `execution_record`;
   - obligatorio para mutation y casos con impacto externo.

5. **Al crear o visualizar un `approval_request` sensible**
   - acciones: `execute`, `release_execution`, `release_application`, `view_restricted` según el punto;
   - resource principal: `approval_request`.

6. **Al consultar vistas restringidas o evidencia clasificada (`view_restricted`)**
   - resource principal: `execution_record`, `approval_request`, `intent_contract` o `policy_artifact` según la vista.

7. **Cuando cambia un input material después de una decisión previa**
   - cambio de `policy_version`;
   - cambio de `classification_level`;
   - cambio de delegación;
   - cambio de `approval_mode_effective`;
   - cambio de scope o external effect.

## Input canonizado exacto a Cerbos

El kernel debe construir un envelope estable llamado `policy_input_v1`.

```yaml
policy_input_v1:
  policy_input_schema_version: v1
  request:
    decision_point: compile | execute | release_execution | release_application | view_restricted
    requested_at: timestamp
    trace_id: string
    idempotency_key: string?
  principal:
    subject_id: string
    subject_type: user | service | worker | system | approver
    roles: [string]
    tenant_id: string
    acting_for_subject_id: string?
    delegation_id: string?
    clearance_level: public | internal | confidential | restricted
    attributes:
      authority_level: string?
      approver_types: [string]?
      authn_assurance_level: string?
      actor_mode: direct | delegated
      delegation_valid: bool
  resource:
    resource_kind: intent_contract | execution_record | approval_request | capability | policy_artifact
    resource_id: string
    tenant_id: string
    resource_version: string?
    policy_version: string
    contract_id: string?
    execution_id: string?
    approval_request_id: string?
    capability_id: string?
    attributes: {}
  action:
    name: compile | execute | release_execution | release_application | view_restricted
  context:
    environment: dev | staging | prod
    risk_level: low | medium | high | critical
    classification_level: public | internal | confidential | restricted
    approval_mode_effective: auto | pre_execution | pre_application | double
    external_effect: none | reversible | irreversible
    scope_size: single_resource | bounded_set | broad_scope
    flags: {}
  correlation:
    tenant_id: string
    contract_id: string?
    execution_id: string?
    approval_request_id: string?
    trace_id: string
```

### Reglas obligatorias del input canonizado

- `policy_input_schema_version` es obligatorio.
- `tenant_id` debe existir en `principal`, `resource` y `correlation` cuando aplique.
- `resource_kind` y `action.name` deben pertenecer al catálogo cerrado v1.
- Los campos opcionales ausentes deben omitirse o enviarse nulos de forma consistente; no mezclar formas.
- El kernel debe rechazar inputs incompletos antes de llamar al PDP.
- El mismo estado material debe producir el mismo `policy_input_v1`.

## Resource model recomendado para Cerbos en Opyta Sync

### 1. `intent_contract`

Representa el contrato compilado o en punto de compilación.

Campos mínimos en `resource.attributes`:

- `contract_state`
- `contract_version`
- `contract_fingerprint`
- `capability_id`
- `result_type`
- `approval_mode_effective`
- `classification_level`
- `risk_level`
- `policy_snapshot_version`
- `plan_snapshot_present`
- `delegation_present`

### 2. `execution_record`

Representa la ejecución durable proyectada por C.2.

Campos mínimos:

- `execution_status`
- `execution_phase_status`
- `application_phase_status`
- `contract_id`
- `contract_fingerprint`
- `capability_id`
- `approval_status`
- `classification_level`
- `risk_level`
- `external_effect`
- `unknown_outcome`
- `release_stage`: execution | application

### 3. `approval_request`

Representa la solicitud de approval gobernada por A.4.

Campos mínimos:

- `request_status`
- `approval_mode`
- `required_approvals_count`
- `received_approvals_count`
- `material_change_fingerprint`
- `classification_level`
- `risk_level`
- `policy_version_binding`
- `expires_at`
- `human_approval_required`

### 4. `capability`

Representa una capability registrable/resoluble visible al kernel.

Campos mínimos:

- `capability_id`
- `capability_version`
- `capability_type`
- `sensitivity`
- `supports_external_effect`
- `supports_irreversible_effect`
- `allowed_environments`

### 5. `policy_artifact` (aplica en v1 para vistas y auditoría)

No es el foco principal, pero permite gobernar acceso a evidencia/policies.

Campos mínimos:

- `artifact_type`: policy_bundle | decision_record | audit_view
- `artifact_version`
- `classification_level`
- `tenant_scope`
- `contains_restricted_reasons`

## Action model recomendado para Cerbos en Opyta Sync

El catálogo mínimo cerrado v1 es:

- `compile`
- `execute`
- `release_execution`
- `release_application`
- `view_restricted`

### Semántica v1 de cada action

- `compile`: autoriza materializar o recompilar `intent_contract`.
- `execute`: autoriza liberar la fase de ejecución técnica sobre `execution_record`/`intent_contract`.
- `release_execution`: autoriza sacar una ejecución de un gate de governance previo a execution.
- `release_application`: autoriza liberar la fase de aplicación luego de `execution_completed`.
- `view_restricted`: autoriza acceso a vistas, evidencia o artifacts clasificados/restringidos.

No se agregan acciones libres en runtime. Nuevas acciones exigen cambio de catálogo y nueva versión de input canonizado.

## Principal model recomendado para Cerbos en Opyta Sync

El principal mínimo recomendado es:

- `subject_id`
- `subject_type`
- `roles`
- `tenant_id`
- `acting_for_subject_id`
- `delegation_id`
- `clearance_level`
- atributos contextuales relevantes

### Shape recomendado

```yaml
principal:
  subject_id: string
  subject_type: user | service | worker | system | approver
  roles: [string]
  tenant_id: string
  acting_for_subject_id: string?
  delegation_id: string?
  clearance_level: public | internal | confidential | restricted
  attributes:
    actor_mode: direct | delegated
    delegation_valid: bool
    delegation_expires_at: timestamp?
    authority_level: string?
    approver_types: [string]?
    employment_status: active | suspended | disabled?
    break_glass: bool?
    managed_service_identity: bool?
```

### Reglas v1 del principal

- `tenant_id` es obligatorio excepto en actores estrictamente platform/global, que igualmente deben explicitar scope efectivo.
- Si existe `delegation_id`, debe existir `acting_for_subject_id` o evidencia equivalente.
- `clearance_level` nunca puede omitirse para `view_restricted`.
- `roles` deben llegar ya resueltos por el kernel; Cerbos no resuelve identidad primaria.

## Context model mínimo recomendado para Cerbos

El context mínimo recomendado es:

- `environment`
- `risk_level`
- `classification_level`
- `approval_mode_effective`
- `external_effect`
- `scope_size`
- flags relevantes de execution/runtime state

### Shape recomendado

```yaml
context:
  environment: dev | staging | prod
  risk_level: low | medium | high | critical
  classification_level: public | internal | confidential | restricted
  approval_mode_effective: auto | pre_execution | pre_application | double
  external_effect: none | reversible | irreversible
  scope_size: single_resource | bounded_set | broad_scope
  flags:
    contract_executable: bool?
    approval_present: bool?
    approval_valid: bool?
    execution_completed: bool?
    application_release_requested: bool?
    policy_recheck_required: bool?
    cross_tenant: bool?
    runtime_mutation: bool?
    read_only: bool?
    unknown_outcome: bool?
```

### Regla de minimalismo

El contexto debe incluir sólo atributos necesarios para policy contextual repetible. Datos enormes, evidence blobs o snapshots completos deben persistirse fuera del request y referenciarse por IDs/versiones.

## Output esperado de Cerbos

Cerbos debe responder con una decisión que el kernel pueda reducir a una de estas consecuencias canónicas:

- `allow`
- `deny_block`
- `require_approval`
- `require_escalation`
- `restricted_view`

### Forma lógica esperada del adapter

El `cerbos_client_adapter` debe producir una respuesta interna normalizada similar a:

```yaml
policy_decision_response_v1:
  cerbos_decision_id: string?
  policy_version: string
  effect: allow | deny
  matched_policy_refs: [string]
  outputs:
    governance_hint: none | require_approval | require_escalation | restricted_view
    reason_codes: [string]
    obligations: [string]
  raw_decision_ref: string?
```

El kernel nunca consume la respuesta raw directo desde runtime; siempre pasa por normalización + mapping.

## Mapping exacto de decisiones Cerbos → runtime / approvals / clasificación

### Tabla de mapping v1

| Cerbos normalizado | Condición adicional local | Decisión del kernel | Efecto en runtime/approval/classification |
|---|---|---|---|
| `effect=allow`, `governance_hint=none` | ninguna contradicción local | `allow` | permite continuar al siguiente estado válido |
| `effect=deny` | cualquier caso mapeable | `deny_block` | transición a `blocked` o rechazo de operación con evidencia |
| `effect=allow`, `governance_hint=require_approval` | operation compatible con approval | `require_approval` | crea/reutiliza `approval_request` y entra o permanece en `awaiting_approval` |
| `effect=deny`, `governance_hint=require_escalation` o deny por razón sensible | escalación requerida | `require_escalation` | bloquea avance automático y deriva a escalación/manual governance |
| `effect=allow`, `governance_hint=restricted_view` | consulta de lectura sensible | `restricted_view` | permite sólo vista redacted/restringida, no full payload |

### Reglas duras de mapping

1. El kernel debe mapear la salida de Cerbos a este catálogo cerrado; no inventa un sexto resultado por caller.
2. Si Cerbos devuelve `allow` pero el runtime detecta prerequisito estructural faltante o approval requerida por A.2/A.4, el kernel NO fuerza avance; prevalecen invariantes locales y se deriva a `require_approval` o bloqueo estructural según corresponda.
3. `restricted_view` sólo aplica a lecturas/vistas; nunca habilita mutación.
4. `require_escalation` no equivale a approval humana estándar; representa bloqueo con intervención reforzada.
5. Respuesta no mapeable = falla cerrada.

## Audit trail mínimo de policy decision

Cada consulta material de policy debe dejar como mínimo:

- `policy_decision_record_id`
- `decision_point`
- `resource_kind`
- `resource_id`
- `action`
- `tenant_id`
- `contract_id` si aplica
- `execution_id` si aplica
- `approval_request_id` si aplica
- `trace_id`
- `principal_subject_id`
- `principal_roles_snapshot`
- `input_schema_version`
- `policy_version`
- `cerbos_decision_id` si existe
- `decision_effect_raw`
- `decision_mapped`
- `reason_codes[]`
- `matched_policy_refs[]`
- `cache_status`: miss | hit | bypass
- `requested_at`
- `decided_at`
- `persisted_at`

El audit trail puede guardar input completo redacted, hash del input o referencia al payload persistido según clasificación, pero no puede perder correlación mínima.

## Persistencia y correlación operativa

Debe existir `policy_decision_record` como objeto persistible o equivalente embebido/referenciado. En v1 se recomienda objeto explícito.

### Campos mínimos

- `policy_decision_record_id`
- `tenant_id`
- `contract_id`
- `execution_id`
- `approval_request_id`
- `trace_id`
- `policy_version`
- `policy_input_schema_version`
- `resource_kind`
- `resource_id`
- `action`
- `principal_subject_id`
- `principal_tenant_id`
- `decision_mapped`
- `decision_effect_raw`
- `reason_codes[]`
- `correlation_hash`
- `request_hash`
- `created_at`

### Regla de correlación mínima obligatoria

Toda decisión material debe poder correlacionarse como mínimo con:

- `tenant_id`
- `contract_id`
- `execution_id`
- `approval_request_id`
- `trace_id`
- `policy_version`

Cuando alguno no aplique, debe persistirse explícitamente como `null` y no omitirse sin criterio.

## Failure mode y fallback seguro

### Regla general

El failure mode por defecto es **fail closed** para mutación y operaciones sensibles.

### Casos v1

1. **Cerbos no responde en mutation o acción sensible**
   - resultado: `deny_block` o `require_escalation` según el punto;
   - nunca avanzar a `executing` o `applying`.

2. **Cerbos no responde en read-only no sensible**
   - degradado futuro posible, pero en v1 sólo puede permitirse si existe regla explícita del kernel para diagnóstico/read-only low risk;
   - si se habilita, debe ser auditado como degraded decision. Por defecto: bloqueo seguro.

3. **Input canonizado inválido**
   - no llamar al PDP;
   - tratar como falla cerrada del PEP.

4. **Output de Cerbos no mapeable o inconsistente**
   - tratar como falla cerrada;
   - registrar evidencia de incompatibilidad.

5. **Persistencia de `policy_decision_record` falla**
   - para mutaciones/sensibles, la decisión no puede considerarse confirmada operativamente;
   - el kernel debe bloquear avance o reintentar persistencia antes de liberar el runtime.

## Caché / performance / invalidación (solo nivel diseño, no tuning fino)

La caché es permitida sólo como optimización secundaria mediante `policy_cache_guard`.

### Reglas v1

- La caché nunca reemplaza el registro auditable.
- La key de caché debe incluir, como mínimo:
  - `policy_input_schema_version`
  - `policy_version`
  - `tenant_id`
  - `resource_kind`
  - `resource_id` o fingerprint material equivalente
  - `action`
  - `principal_subject_id`
  - `roles hash`
  - `classification_level`
  - `risk_level`
  - `approval_mode_effective`
  - `delegation_id` o estado equivalente
- Si cambia `policy_version`, la decisión cacheada es inválida.
- Si cambia `classification_level`, `approval_mode_effective`, delegación o fingerprint material, la decisión cacheada es inválida.
- Para `release_application` y vistas `restricted`, se recomienda caché muy acotada o bypass.
- Si la caché devuelve una decisión con `policy_version` distinta de la esperada, el PEP debe ignorarla y consultar al PDP.

## Integración con `compiled_contract`

- `compiled_contract` debe contener o referenciar `policy_evaluation_input` derivado en C.1.
- C.3 formaliza que ese input no es libre: debe proyectarse a `policy_input_v1`.
- `policy_snapshot` del contrato es insumo obligatorio para `policy_version`/binding.
- Cambios materiales del contrato que afecten `policy_snapshot`, `classification_level`, `approval_mode_efectivo`, `capability_id`, `scope`, `datos_permitidos` o `herramientas_permitidas` invalidan decisiones previas relevantes.
- El compilador prepara el input; el PEP decide cuándo usarlo y cuándo re-canonizar con estado más fresco.

## Integración con `execution_workflow`

- `execution_workflow` nunca llama a Cerbos directo; usa el seam `policy_gate_adapter` / `policy_enforcement_point`.
- El workflow sólo consume decisiones ya mapeadas por el PEP.
- Estados relevantes de C.2:
  - `preconditions_validating -> policy_evaluating`
  - `policy_evaluating -> executing`
  - `policy_evaluating -> awaiting_approval`
  - `policy_evaluating -> blocked`
- `release_execution` y `release_application` deben provocar nueva evaluación de policy cuando el release sea material.
- `execution_completed` no implica permiso automático para `release_application`; requiere consulta separada cuando corresponda.

## Integración con approvals

- `require_approval` desde policy no reemplaza A.4; dispara o mantiene el corredor formal de approvals.
- El `approval_request` debe capturar `policy_decision_snapshot` o referencia al `policy_decision_record`.
- `policy_version`, `classification_snapshot`, `risk_snapshot` y `material_change_fingerprint` deben quedar alineados entre contrato, approval y policy record.
- Si Cerbos requiere approval pero la operación es estructuralmente incompatible con approval flow definido, el kernel debe escalar en vez de improvisar.
- La liberación de approvals (`release_execution`, `release_application`) debe revalidar policy vigente y no confiar ciegamente en una decisión vieja.

## Integración con clasificación/redacción

- Cerbos no redefine `classification_level`; consume clasificación canonizada del kernel.
- La policy puede determinar `restricted_view` y obligaciones de redacción/acceso parcial.
- Si la clasificación sube después de compilación o durante ejecución, la policy debe reevaluarse.
- `view_restricted` requiere clearance suficiente y debe poder devolver vista redacted aunque niegue full view.
- La clasificación aplicada al contrato hereda a evidencia y payloads de `policy_decision_record`.

## Idempotencia y deduplicación de decisiones

- El PEP debe ser idempotente para el mismo `policy_input_v1` material y misma `policy_version`.
- Repetir una consulta igual puede reutilizar resultado cacheado o persistido, pero siempre preservando evidencia/correlación.
- Debe existir `request_hash` o equivalente sobre el input canonizado redacted/estable.
- Si el mismo `trace_id` y `request_hash` ya existen con decisión persistida válida, el PEP puede deduplicar.
- Si el runtime recibió respuesta pero no se persistió evidencia, no debe asumirse éxito operativo final para mutation.

## Componentes mínimos esperados de la integración

1. `policy_enforcement_point`
   - coordina consulta, validación y resultado consumible por el kernel.

2. `policy_input_canonicalizer`
   - transforma objetos del dominio a `policy_input_v1`.

3. `cerbos_client_adapter`
   - encapsula API de Cerbos, timeouts y respuesta raw.

4. `policy_decision_mapper`
   - traduce respuesta normalizada a `allow`, `deny_block`, `require_approval`, `require_escalation`, `restricted_view`.

5. `policy_decision_recorder`
   - persiste `policy_decision_record` y correlación mínima.

6. `policy_cache_guard`
   - decide hit/miss/bypass e invalidación por `policy_version` y materialidad.

## Tests borde mínimos

Como mínimo deben existir estos casos de borde para validar la integración:

1. `resource_kind` desconocido → el PEP rechaza antes de llamar a Cerbos.
2. `action` desconocida → falla cerrada con razón auditable.
3. `tenant_id` faltante en input → no se consulta policy.
4. principal sin tenant → rechazo salvo actor global explícitamente soportado.
5. delegación expirada → deny o escalation según policy/local invariant.
6. clearance insuficiente para `view_restricted` → `restricted_view` parcial o denegación según policy definida.
7. Cerbos devuelve `allow` pero runtime detecta approval requerida → prevalece invariant local y no avanza directo.
8. Cerbos devuelve `require_approval` para operación read-only low risk → el kernel la acepta sólo si policy lo exige explícitamente y deja evidencia; no la degrada silenciosamente.
9. Cerbos no responde en mutation → fail closed.
10. Cerbos no responde en read-only no sensible → comportamiento explícito, auditable y por defecto bloqueado.
11. `policy_version` cambia entre compilación y ejecución → reevaluación obligatoria.
12. input canonizado con field faltante obligatorio → rechazo en canonicalizer.
13. output de Cerbos no mapeable → fail closed.
14. caché entrega decisión vieja con `policy_version` distinta → bypass de caché y reevaluación.
15. execution liberada sin `release_execution` permitido → bloqueo del release.
16. `release_application` denegado después de `execution_completed` → no puede entrar a `applying`.
17. `approval_request` visible pero no releaseable → `view_restricted` permitido sin `release_execution`/`release_application`.
18. clasificación cambia y policy debe reevaluarse → decisión previa queda inválida.
19. record de decisión no persiste pero el runtime ya recibió respuesta → la liberación no se confirma para mutación.
20. `principal.roles` vacío en acción sensible → policy debe evaluar con menor privilegio, no asumir defaults permisivos.
21. `resource.tenant_id` distinto de `principal.tenant_id` sin bandera cross-tenant → deny.
22. `approval_request` expirado intenta `release_execution` → deny_block.
23. `release_application` llega antes de `execution_completed` → el runtime rechaza aunque policy diga allow.
24. `restricted_view` permitido sobre evidencia clasificada pero payload full bloqueado → redacción correcta.
25. mismatch entre `contract_fingerprint` del `execution_record` y el del input de policy → invalidación/recheck.
26. cambio de delegación entre compile y execute → reevaluación obligatoria.
27. respuesta cacheada de sujeto A reutilizada para sujeto B por bug de key → debe detectarse por key incompleta en test negativo.

## Criterios de aceptación de C.3

1. Queda explícito y sin ambigüedad que el kernel es **PEP** y Cerbos es **PDP**.
2. Queda explícito que Cerbos decide autorización/policy contextual y **no** lifecycle del runtime.
3. Existe `policy_input_v1` canonizado, versionado y no arbitrario.
4. El catálogo mínimo de `resource_kind` incluye `intent_contract`, `execution_record`, `approval_request`, `capability` y `policy_artifact`.
5. El catálogo mínimo de `action` incluye `compile`, `execute`, `release_execution`, `release_application` y `view_restricted`.
6. El principal model, resource model y context model mínimos quedan definidos a nivel implementable.
7. El output de Cerbos queda mapeado como mínimo a `allow`, `deny_block`, `require_approval`, `require_escalation` y `restricted_view`.
8. El kernel traduce la salida de Cerbos a transiciones conocidas y no la reinterpretará libremente.
9. Existe `policy_decision_record` o equivalente con evidencia auditable mínima.
10. La correlación mínima obligatoria incluye `tenant_id`, `contract_id`, `execution_id`, `approval_request_id`, `trace_id` y `policy_version`.
11. El failure mode por defecto queda cerrado como **fail closed** para mutaciones y operaciones sensibles.
12. La integración con `compiled_contract`, `execution_workflow`, approvals y clasificación queda definida sin contradecir A.2, A.4, B.2, C.1 ni C.2.
13. Quedan definidos tests borde mínimos suficientes para validar el seam de policy en Fase C.
