# A.5 — Runtime, estados de ejecución y eventos generales v1

## Principios base del runtime

- El runtime general existe para **orquestar** contrato, approvals, ejecución técnica, aplicación de efectos y cierre auditable; no reemplaza los lifecycles específicos de A.2, A.3 ni A.4.
- El runtime sigue el flujo base del core: **contrato -> acción de negocio -> tools/sistemas base**. Nunca ejecuta si el contrato no quedó consistente.
- Ante intención incompleta, el runtime respeta la secuencia fundacional: **inspeccionar -> memoria -> asumir/proponer -> preguntar solo si persiste ambigüedad crítica -> ejecutar únicamente cuando el contrato quede consistente**.
- El runtime distingue explícitamente entre **habilitación para correr** y **habilitación para aplicar efectos**. Esa separación es obligatoria para evitar falsos positivos de seguridad y auditoría.
- El runtime general debe poder representar tres clases de detención sin ambigüedad: **bloqueo por governance**, **falla técnica** y **compensación posterior**.
- Los estados del runtime no deben duplicar semántica de contrato, approvals o resultado; deben actuar como **capa orquestadora superior**.
- Todo cambio de estado del runtime debe ser trazable mediante `execution_record` + `telemetry_event` + correlación estable de IDs.
- No existe cierre válido de ejecución sin trazabilidad completa de: contrato efectivo, approvals relevantes, resultado(s), evidencia y eventos terminales.

---

## Definición de `execution_record` como objeto canónico y su rol

`execution_record` es el objeto canónico runtime-owned que representa una **instancia concreta de ejecución orquestada** bajo un `intent_contract` determinado.

Su rol no es describir el contrato ni materializar el resultado final, sino registrar la **verdad operacional** de una corrida específica:

- qué contrato intentó ejecutarse,
- bajo qué condiciones de approval y governance,
- cuándo quedó elegible,
- cuándo fue liberada la ejecución,
- cuándo se ejecutó técnicamente,
- cuándo se liberó y aplicó el efecto externo si correspondía,
- qué resultado(s) produjo,
- si falló, quedó bloqueada, fue compensada o cerró.

`execution_record` es first-class porque:

1. Tiene lifecycle propio y distinto del contrato.
2. Puede existir aunque la ejecución nunca llegue a correr, falle técnicamente o quede bloqueada.
3. Necesita soportar idempotencia, correlación, auditoría y evidencia append-only.
4. Es la frontera formal entre **capacidad ejecutable** y **operación realmente intentada/realizada**.

Normativamente:

- Un `intent_contract` puede originar **cero, una o varias** ejecuciones a lo largo del tiempo, pero cada `execution_record` referencia exactamente un `intent_contract` efectivo.
- Un `result_record` referencia la ejecución que lo produjo, no reemplaza a la ejecución.
- Una `approval_decision` puede habilitar o bloquear una ejecución, pero no absorbe el lifecycle del runtime.

---

## Schema sugerido de `execution_record`

Schema sugerido mínimo v1 a nivel lógico:

```yaml
api_version: v1
kind: execution_record
metadata:
  id: exe-2026-0001
  tenant_id: ten-...
  environment: prod
  object_version: "1.0"
  schema_version: v1
spec:
  execution_id: exe-2026-0001
  trace_id: trc-...
  parent_execution_id: null
  intent_contract_id: ctr-...
  contract_fingerprint: sha256:...
  capability_id: cap-...
  result_type: execution
  approval_mode_effective: pre_application
  execution_state: created
  release_model:
    execution_authorized: false
    execution_released: false
    application_required: true
    application_authorized: false
    application_released: false
  governance_status:
    is_blocked: false
    block_scope: none
    block_reason_code: null
  failure_status:
    has_failed: false
    failure_stage: null
    failure_reason_code: null
  compensation_status:
    is_compensation_required: false
    compensation_state: not_required
  linked_ids:
    active_approval_request_id: null
    approval_decision_ids: []
    result_ids: []
    telemetry_event_ids: []
  timing:
    created_at: 2026-03-29T10:00:00Z
    eligibility_checked_at: null
    execution_released_at: null
    execution_started_at: null
    execution_completed_at: null
    application_released_at: null
    application_started_at: null
    application_completed_at: null
    compensated_at: null
    failed_at: null
    closed_at: null
  actors:
    initiated_by_subject_id: sub-...
    released_by_subject_id: null
    applied_by_subject_id: null
    closed_by_subject_id: null
  evidence:
    policy_snapshot_ref: polsnap-...
    approval_snapshot_ref: null
    plan_snapshot_ref: plansnap-...
    execution_evidence_ref: null
    application_evidence_ref: null
    compensation_evidence_ref: null
status:
  terminal: false
  last_event_at: 2026-03-29T10:00:00Z
  last_event_type: execution.created
```

### Campos normativos mínimos

- `execution_id`
- `trace_id`
- `tenant_id`
- `environment`
- `intent_contract_id`
- `contract_fingerprint`
- `capability_id`
- `approval_mode_effective`
- `execution_state`
- `linked_ids.result_ids`
- timestamps de transición relevantes
- referencias de evidencia mínima

### Reglas de modelado

- `execution_state` refleja únicamente el lifecycle general del runtime.
- Los snapshots de policy, plan, approval y clasificación viven referenciados o embebidos según corresponda, pero el `execution_record` conserva las referencias canónicas de evidencia.
- El record es **append-only**: cambios materiales se representan por nuevos eventos, nuevos snapshots, nuevos timestamps y nuevas transiciones válidas; nunca por overwrite destructivo.

---

## Estados generales de ejecución

Estados propuestos y adoptados para el runtime general v1:

- `created` — la ejecución fue creada como record operativo, pero todavía no se evaluó si puede correr.
- `eligibility_check` — el runtime está validando si la ejecución es elegible bajo contrato, approvals, policy, fingerprint y precondiciones técnicas.
- `blocked` — la ejecución no puede avanzar por una restricción de governance, autorización, vigencia, fingerprint o compliance. No representa falla técnica.
- `awaiting_approval` — la ejecución está detenida porque requiere resolución de approval o release y todavía no existe autorización vigente suficiente.
- `execution_released` — la ejecución quedó liberada para correr técnicamente; esto NO implica autorización para aplicar efectos externos.
- `executing` — la parte técnica de la ejecución está corriendo.
- `execution_completed` — la corrida técnica terminó y dejó evidencia/resultados, pero todavía no implica aplicación efectiva.
- `application_released` — la aplicación de efectos quedó liberada; solo aplica cuando hay mutación o external effect gobernado.
- `applying` — el runtime está aplicando efectos externos o persistiendo cambios materiales gobernados.
- `application_completed` — la aplicación material terminó satisfactoriamente.
- `partially_compensated` — se aplicaron medidas de compensación parciales, pero queda remanente pendiente o no reversible por completo.
- `compensated` — la compensación definida por policy/procedimiento terminó y quedó evidenciada.
- `failed` — la ejecución sufrió una falla técnica u operacional y no pudo completar el paso actual esperado.
- `closed` — la ejecución quedó cerrada terminalmente con trazabilidad completa.

### Intención semántica de los estados más sensibles

- **Elegible para ejecutar** no es un estado final; es una condición evaluada durante `eligibility_check`.
- **Liberado para ejecutar** se materializa en `execution_released`.
- **Ejecutando** se materializa en `executing`.
- **Ejecutado** se materializa en `execution_completed`.
- **Aplicado** se materializa en `application_completed`, nunca en `execution_completed`.
- **Compensado** se materializa en `compensated`.
- **Cerrado** se materializa en `closed` y exige trazabilidad terminal completa.

---

## Reglas de relación entre execution state vs contract state vs approval state vs result state

### 1. Execution state vs contract state

- El runtime no puede crear una ejecución real si el contrato no existe como unidad válida runtime-owned.
- `execution_record.created` puede existir cuando el contrato ya está en una fase posterior de preparación, pero **no puede avanzar a `execution_released`** si el contrato no está en `executable`.
- `executing` del runtime exige que el contrato ya haya alcanzado `executable` y que la transición a ejecución no haya sido invalidada por cambio material.
- `execution_completed` del runtime no obliga al contrato a pasar automáticamente a `closed`; el contrato puede requerir cierre posterior una vez completado el resultado/aplicación.

### 2. Execution state vs approval state

- `awaiting_approval` existe cuando el runtime requiere una approval o release vigente y todavía no la tiene.
- `execution_released` del runtime solo es válido si existe una condición equivalente en approvals: `approved` con modo `auto`, `execution_released`, o decisión vigente que normativamente habilite correr.
- `application_released` del runtime solo es válido si el approval lifecycle alcanzó una condición efectiva equivalente a `application_released` o si el modo `auto` permite aplicar sin segundo release.
- Si la approval vence, es revocada o queda superseded antes de correr, la ejecución debe volver a `blocked` o permanecer en `awaiting_approval`; nunca puede seguir como liberada.

### 3. Execution state vs result state

- El runtime puede llegar a `execution_completed` cuando la corrida técnica terminó, aunque el `result_record` todavía esté transitando `produced -> classifying -> delivered`.
- `execution_completed` no implica `application_completed` ni `result_record.closed`.
- Para capacidades read-only, `execution_completed` suele coexistir con un resultado en lifecycle A.3; `application_released` y `applying` pueden no aplicar.
- Para mutaciones, la separación `execution_completed -> application_released -> applying -> application_completed` es obligatoria aunque el resultado se produzca antes o durante la aplicación.

### 4. Dominancia de estados

- **Bloqueo por governance** domina sobre liberación pendiente: si una restricción dura invalida la corrida, el runtime cae en `blocked`.
- **Falla técnica** domina sobre progreso técnico: si falla la ejecución o la aplicación, el runtime cae en `failed`.
- **Compensación** no borra la falla o el efecto previo; la compensa. Por eso `compensated` no significa “nunca falló”, sino “se restauró o mitigó bajo procedimiento”.
- `closed` solo es válido cuando los demás lifecycles relevantes ya tienen estado terminal consistente o evidencia de por qué no aplican.

---

## Transiciones válidas del execution lifecycle

```text
created -> eligibility_check
created -> blocked

eligibility_check -> awaiting_approval
eligibility_check -> execution_released
eligibility_check -> blocked
eligibility_check -> failed

awaiting_approval -> execution_released
awaiting_approval -> blocked
awaiting_approval -> failed

execution_released -> executing
execution_released -> blocked
execution_released -> failed

executing -> execution_completed
executing -> failed

execution_completed -> application_released
execution_completed -> closed
execution_completed -> failed

application_released -> applying
application_released -> blocked
application_released -> failed

applying -> application_completed
applying -> partially_compensated
applying -> compensated
applying -> failed

application_completed -> closed
application_completed -> partially_compensated
application_completed -> compensated

failed -> partially_compensated
failed -> compensated
failed -> closed

partially_compensated -> compensated
partially_compensated -> failed
partially_compensated -> closed

compensated -> closed
blocked -> closed
```

### Reglas adicionales de transición

- `execution_completed -> closed` solo es válido para ejecuciones read-only o sin fase de aplicación requerida.
- `application_completed -> compensated` es válido si la policy exige rollback/compensación posterior por verificación negativa o revocación tardía.
- `failed -> closed` solo es válido si ya existe evidencia suficiente para cierre sin compensación adicional requerida.

---

## Transiciones prohibidas importantes

- `created -> executing`.
- `created -> execution_completed`.
- `eligibility_check -> applying`.
- `awaiting_approval -> executing` sin release vigente.
- `blocked -> executing` sin crear nueva evaluación/liberación válida.
- `execution_released -> application_completed`.
- `executing -> application_completed` sin pasar por `execution_completed`.
- `execution_completed -> applying` sin `application_released` en mutaciones.
- `execution_completed -> application_completed` directo en mutaciones.
- `application_released -> execution_released` (no se retrocede semánticamente).
- `application_completed -> executing`.
- `compensated -> applying`.
- `closed -> cualquier_otro_estado`.

### Prohibición estructural clave

El runtime NO puede reutilizar un `execution_record` cerrado, fallido terminal o compensado como si fuera una nueva corrida. Si hace falta reintentar, corresponde un **nuevo `execution_id`** con referencia causal al anterior si aplica.

---

## Release model claro

### 1. `authorized`

Una ejecución está **autorizada** cuando existe base normativa vigente para permitir el siguiente paso gobernado. La autorización puede provenir de:

- `approval_mode = auto` con policy válida,
- approval humana/política aprobada,
- release específico de ejecución o aplicación.

Autorización no equivale todavía a transición de runtime. Es una condición de governance.

### 2. `eligible to execute`

Una ejecución es **elegible para ejecutar** cuando simultáneamente:

1. el contrato está en `executable`,
2. el fingerprint material coincide con el aprobado,
3. la approval vigente cubre el paso requerido,
4. no existe bloqueo duro por policy, clasificación, delegación o expiración,
5. las precondiciones técnicas mínimas están disponibles.

La elegibilidad se determina durante `eligibility_check`.

### 3. `execution released`

Una ejecución está **liberada para correr** cuando, además de ser elegible, el runtime registra una liberación efectiva y auditable para iniciar la corrida técnica. Se materializa en el estado `execution_released`.

### 4. `executing`

La ejecución está **ejecutando** cuando la capability ya está corriendo sobre tools/sistemas base y consumiendo recursos operativos.

### 5. `execution completed`

La ejecución está **ejecutada** cuando la corrida técnica finalizó y dejó evidencia suficiente de outcome técnico. Esto puede producir:

- resultados read-only listos para clasificación/entrega,
- una mutación todavía no aplicada,
- una preparación técnica previa a aplicación,
- evidencia de que el intento terminó sin poder aplicar.

Norma dura: `execution_completed` **NO implica** `application_completed`.

### 6. `application released`

La ejecución queda **liberada para aplicar** cuando existe autorización vigente para materializar efectos externos o persistir cambios gobernados. En mutaciones esta separación es obligatoria.

### 7. `application completed`

La ejecución queda **aplicada** cuando el efecto externo comprometido fue realizado y verificado conforme a policy/capability/result_type. No basta con “llamé a la API”; debe existir evidencia mínima de aplicación o confirmación equivalente.

### 8. `compensated`

La ejecución queda **compensada** cuando, tras una falla o reversión requerida, se aplicó un procedimiento de rollback, mitigación o corrección definido y evidenciado. Puede ser:

- `partially_compensated` si el rollback es incompleto,
- `compensated` si el remanente quedó resuelto dentro del alcance permitido.

### 9. `completed`

En este documento, “completada” NO se usa como alias ambiguo. Debe distinguirse siempre entre:

- `execution_completed`,
- `application_completed`,
- `closed`.

### 10. `closed`

La ejecución queda **cerrada** cuando ya no admite nuevas transiciones y el expediente auditable está completo. `closed` es un estado administrativo terminal, no un outcome semántico por sí solo.

---

## Eventos canónicos generales del runtime

Estos eventos conviven con los de contrato, approvals y resultado. No los reemplazan.

### Eventos principales del lifecycle general

- `execution.created`
- `execution.eligibility_check_started`
- `execution.eligibility_confirmed`
- `execution.eligibility_rejected`
- `execution.awaiting_approval`
- `execution.execution_released`
- `execution.started`
- `execution.completed`
- `execution.application_released`
- `execution.application_started`
- `execution.application_completed`
- `execution.partially_compensated`
- `execution.compensated`
- `execution.failed`
- `execution.closed`

### Eventos generales por causa/orquestación

- `execution.contract_linked`
- `execution.fingerprint_verified`
- `execution.fingerprint_mismatch_detected`
- `execution.governance_blocked`
- `execution.approval_linked`
- `execution.approval_missing`
- `execution.approval_expired_detected`
- `execution.release_revoked`
- `execution.result_linked`
- `execution.technical_failure_detected`
- `execution.compensation_required`
- `execution.closure_validated`
- `execution.closure_rejected`

### Alineación con eventos base del source-of-truth

- `intent_received` e `inspection_contract_created` siguen viviendo en la capa contractual/orquestación temprana.
- `approval_requested` y `approval_resolved` siguen siendo eventos base cross-domain; aquí se profesionalizan en correlación con `execution.awaiting_approval`, `execution.approval_linked` y releases específicos.
- `tool_executed` sigue siendo granular de capa de ejecución; no reemplaza `execution.started` ni `execution.completed`.
- `trace_closed` sigue existiendo, pero `execution.closed` expresa el cierre específico del runtime general.

---

## Taxonomía de eventos por dominio/capa

### Dominio 1 — Contract

- `contract.*`
- Responsable de lifecycle contractual A.2.

### Dominio 2 — Approval

- `approval.*`
- Responsable de lifecycle de approvals A.4.

### Dominio 3 — Execution runtime general

- `execution.*`
- Responsable de la orquestación superior definida en A.5.

### Dominio 4 — Result

- `result.*`
- Responsable del lifecycle de resultados A.3.

### Dominio 5 — Tool/step/span técnico

- `tool.*`, `span.*`, `trace.*` o equivalente interno.
- Responsable de granularidad de ejecución operativa y observabilidad de bajo nivel.

### Regla de capa

- Un evento `execution.*` nunca debe duplicar payload semántico completo de `approval.*` o `result.*`; debe **referenciarlos y orquestarlos**.
- Los dashboards o auditorías pueden proyectar una vista unificada, pero el canon mantiene los dominios separados.

---

## Payload mínimo sugerido por evento general

Todo evento `execution.*` debe incluir como mínimo:

- `event_id`
- `event_type`
- `tenant_id`
- `environment`
- `trace_id`
- `execution_id`
- `intent_contract_id`
- `capability_id`
- `approval_request_id` (nullable cuando no aplica)
- `result_id` (nullable cuando no aplica)
- `from_state`
- `to_state`
- `release_scope` (`none`, `execution`, `application`)
- `reason_code` (cuando aplica)
- `triggered_by_subject_id`
- `occurred_at`

### Payload recomendado ampliado

- `contract_fingerprint`
- `approval_mode_effective`
- `result_type`
- `has_external_effect`
- `is_read_only`
- `is_compensation_required`
- `failure_stage`
- `evidence_ref`
- `correlation_version`

---

## Reglas de correlación e IDs

### IDs canónicos mínimos

- `execution_id` — identidad de la corrida runtime.
- `trace_id` — correlación transversal de observabilidad.
- `contract_id` / `intent_contract_id` — identidad del contrato efectivo.
- `approval_request_id` — identidad del request de approval activo o causalmente relevante.
- `result_id` — identidad del resultado producido por la ejecución.

### Reglas normativas

1. Todo `execution_record` tiene `execution_id` único dentro del tenant y ambiente.
2. Todo `telemetry_event` asociado a una ejecución debe incluir `execution_id` y `trace_id`.
3. `trace_id` puede agrupar múltiples eventos y subspans, pero no puede mezclar ejecuciones independientes sin referencia causal explícita.
4. Un `result_record` debe apuntar a un único `execution_id` originante.
5. Si una ejecución requiere más de un approval request a lo largo de su vida, el `execution_record` conserva el activo y el historial referenciado.
6. Si existe reintento, replay o compensación separada, debe emitirse nuevo `execution_id`; `parent_execution_id` o vínculo causal debe capturarse explícitamente.
7. `contract_fingerprint` fijado en la ejecución debe corresponder al fingerprint validado al liberar ejecución o aplicación.

---

## Reglas de cierre de una ejecución

Una ejecución solo puede pasar a `closed` si se cumplen TODAS las condiciones aplicables:

1. existe `execution_record` persistido e íntegro,
2. el estado actual es terminalmente consistente (`blocked`, `application_completed`, `compensated`, `failed` o `execution_completed` read-only),
3. todas las referencias relevantes (`intent_contract_id`, `execution_id`, `trace_id`) están presentes,
4. los approvals relevantes están terminales, vigentes históricamente o explícitamente marcados como no aplicables,
5. los `result_record` vinculados están creados o la ausencia de resultado está justificada por `reason_code`,
6. la evidencia mínima de ejecución/aplicación/compensación/falla está registrada,
7. se emitió `execution.closed`,
8. el `trace_closed` técnico puede correlacionarse con la ejecución cerrada.

### Reglas específicas

- Si la ejecución falló y requiere compensación por policy, no puede cerrar antes de resolver `partially_compensated` o justificar por qué se cierra con remanente aceptado.
- Si hubo mutación, no puede cerrar solo con `execution_completed`; necesita `application_completed`, `compensated` o `failed` con evidencia suficiente del no-aplicado / aplicado parcial.
- Si falta correlación de IDs o evidencia terminal, el cierre debe rechazarse y emitirse `execution.closure_rejected`.

---

## Reglas de observabilidad mínima obligatoria

### Obligatorio por ejecución

- un evento `execution.created`
- un evento de inicio de evaluación (`execution.eligibility_check_started`)
- al menos un evento de resolución de elegibilidad o bloqueo
- un evento de release cuando corresponda (`execution.execution_released`, `execution.application_released`)
- un evento terminal específico (`execution.completed`, `execution.application_completed`, `execution.failed`, `execution.compensated` o `execution.governance_blocked`)
- un evento `execution.closed`

### Obligatorio por correlación

- `trace_id` estable en todos los eventos de la ejecución
- `execution_id` presente en runtime, resultado y telemetría relacionada
- timestamps monotónicos por transición válida
- reason code normalizado en todo `blocked`, `failed`, `partially_compensated` o `closure_rejected`

### Obligatorio por evidencia

- referencia a contract fingerprint validado
- referencia a approval vigente o razón de no aplicación
- referencia a evidencia de ejecución técnica
- referencia a evidencia de aplicación o de compensación cuando corresponda

### Regla de separación memoria/telemetría

- La telemetría del runtime vive en `telemetry_event`.
- La memoria operativa no reemplaza la evidencia de cierre ni los eventos append-only.

---

## Tests borde mínimos

### T1 — Contrato no executable
Dado un contrato en `compiled` pero no en `executable`.
Esperado: `created -> eligibility_check -> blocked`; reason `blocked.contract.not_executable`.

### T2 — Approval requerida ausente
Dado un modo `pre_execution` sin approval vigente.
Esperado: `eligibility_check -> awaiting_approval`.

### T3 — Approval expirada antes de liberar ejecución
Dado un approval request previamente aprobado pero vencido.
Esperado: `eligibility_check -> blocked`; reason `blocked.approval.expired`.

### T4 — Fingerprint mismatch antes de correr
Dado un contrato con material change después de aprobar.
Esperado: `eligibility_check -> blocked`; `execution.fingerprint_mismatch_detected`.

### T5 — Read-only autoaprobado
Dado un caso read-only con `approval_mode = auto` y contrato executable.
Esperado: `eligibility_check -> execution_released -> executing -> execution_completed -> closed`.

### T6 — Mutación con separación obligatoria
Dado un `system_update` mutable.
Esperado: nunca salta de `execution_completed` a `application_completed` directo.

### T7 — Intento de aplicar sin release
Dado un caso mutable sin `application_released`.
Esperado: transición `execution_completed -> applying` rechazada.

### T8 — Bloqueo por governance después de ejecución técnica
Dado un modo `pre_application` donde la segunda approval no llega.
Esperado: `execution_completed -> awaiting_approval` o `blocked`, pero no `application_released`.

### T9 — Falla técnica durante ejecución
Dado error de tool crítica en corrida técnica.
Esperado: `executing -> failed`; no `blocked`.

### T10 — Falla técnica durante aplicación
Dado error al persistir efecto externo.
Esperado: `applying -> failed`; `failure_stage = application`.

### T11 — Compensación parcial
Dado rollback posible solo sobre parte del alcance.
Esperado: `failed -> partially_compensated`.

### T12 — Compensación completa
Dado rollback exitoso completo.
Esperado: `failed -> compensated -> closed`.

### T13 — Cierre sin trace_id
Dado un execution record terminal sin `trace_id`.
Esperado: rechazo de cierre; `execution.closure_rejected`.

### T14 — Cierre de mutación sin evidencia de aplicación
Dado `application_completed` sin evidence ref.
Esperado: no puede pasar a `closed`.

### T15 — Reintento sobre execution_id cerrado
Dado intento de reutilizar `execution_id` en `closed`.
Esperado: inválido; crear nuevo `execution_id`.

### T16 — Resultado producido pero ejecución no completada
Dado un `result_record.produced` sin `execution_completed` previo o concurrente coherente.
Esperado: inconsistencia de correlación; rechazar cierre.

### T17 — Read-only con `application_released`
Dado un resultado read-only puro.
Esperado: `application_released` y `applying` marcados como no aplicables; no deben emitirse salvo caso excepcional explícito.

### T18 — Approval revocada tras `execution_released` y antes de `executing`
Dado release ya emitido pero luego revocado antes de comenzar.
Esperado: `execution_released -> blocked`; reason `blocked.approval.revoked`.

### T19 — Approval revocada tras `execution_completed` en mutación pre-application
Dado `execution_completed` y revocación antes de aplicar.
Esperado: no puede pasar a `application_released`; queda `blocked` o `awaiting_approval` según policy.

### T20 — Cierre sin evento terminal
Dado un record con timestamps completos pero sin `execution.closed`.
Esperado: cierre inválido por observabilidad incompleta.

---

## Criterios de aceptación del runtime general

- Existe `execution_record` como objeto canónico first-class separado de contrato, approval y resultado.
- Los estados del runtime son semánticamente distintos y no colisionan con los estados de A.2, A.3 y A.4.
- Queda explícita la diferencia entre `eligible to execute`, `execution_released`, `executing`, `execution_completed`, `application_released`, `application_completed`, `compensated` y `closed`.
- `execution_completed` y `application_completed` quedan definidos como hechos distintos y no intercambiables.
- El runtime soporta al menos un estado de bloqueo por governance (`blocked`) y uno de falla técnica (`failed`) con reason codes diferenciables.
- Para resultados read-only, la fase de aplicación puede marcarse como no aplicable sin romper el lifecycle general.
- Para mutaciones, la separación entre ejecución y aplicación queda como regla obligatoria.
- Los eventos `execution.*` conviven con `contract.*`, `approval.*`, `result.*` y eventos base de observabilidad sin reemplazarlos.
- Toda ejecución puede correlacionarse de punta a punta con `execution_id`, `trace_id`, `intent_contract_id`, approvals y resultados.
- Ninguna ejecución puede cerrarse sin trazabilidad completa, evidencia mínima y evento terminal de cierre.
- Las transiciones válidas/prohibidas permiten auditar y validar mecánicamente el runtime general v1.

---

## Decisión final de A.5 v1

- `execution_record` queda fijado como la unidad canónica de orquestación runtime.
- El runtime general gobierna la corrida completa sin absorber los lifecycles especializados de contrato, approval ni resultado.
- La separación entre **autorizado**, **liberado para ejecutar**, **ejecutado**, **liberado para aplicar**, **aplicado**, **compensado** y **cerrado** queda normativamente obligatoria.
- La ejecución solo puede cerrarse con trazabilidad integral y observabilidad mínima verificable.
