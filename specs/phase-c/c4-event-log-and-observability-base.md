# C.4 — Event log and observability base

## Objetivo

Definir la construcción operativa v1 del plano de evidencia del kernel para que **PostgreSQL** conserve la **verdad operativa durable** mediante un event log canónico mínimo y append-only, mientras **OpenTelemetry + Grafana LGTM** reciban únicamente **telemetría derivada** para operación, troubleshooting y análisis, sin competir con el source of truth del core ni contradecir A.5, B.3, C.2 o C.3.

## Principios de implementación del event log y observabilidad

1. **PostgreSQL es la verdad operativa.** El estado canónico vive en `compiled_contract`, `execution_record`, `policy_decision_record`, `approval_request` y en el event log mínimo append-only.
2. **OTel/LGTM es telemetría derivada.** Sirve para explicar, correlacionar, monitorear y depurar; NO decide estado canónico ni reemplaza auditoría durable.
3. **Temporal history no reemplaza el event log del kernel.** Temporal puede ayudar a debugging runtime, pero el kernel necesita su propio log canónico tenant-scoped y gobernable.
4. **Append-only obligatorio.** Los hechos operativos se agregan; no se mutan destructivamente.
5. **Separación de planos.** Estado operativo, evidencia auditable y telemetría analítica no se mezclan en un solo storage ni en un solo contrato semántico.
6. **Correlación mínima estable.** Toda emisión material debe enlazarse como mínimo con `tenant_id`, `environment`, `trace_id` y los IDs causales aplicables.
7. **Clasificación y redacción antes de salir.** Ningún payload sensible se persiste/exporta sin pasar por `classification_redaction_guard`.
8. **Observabilidad no puede romper el kernel.** Si OTel/LGTM falla, el kernel conserva estado, transiciones y evidencia canónica.
9. **Idempotencia explícita.** Retry técnico, replay, export retry y duplicado lógico no pueden duplicar hechos canónicos.
10. **Primero hitos del corredor principal.** En v1 se materializan los eventos que permiten reconstruir compilación, policy, approvals, ejecución, aplicación y compensación.

## Boundary exacto entre verdad operativa y telemetría derivada

### Verdad operativa

Pertenece al plano operativo todo hecho que el kernel necesita para:

- reconstruir causalidad auditable,
- demostrar qué decisión tomó y por qué,
- responder consultas operativas confiables,
- reanudar/reintentar/compensar sin ambigüedad,
- soportar deduplicación e idempotencia,
- sostener cumplimiento tenant-scoped.

Eso se persiste en PostgreSQL y se gobierna como evidencia canónica.

### Telemetría derivada

Pertenece al plano OTel/LGTM toda señal usada para:

- trazabilidad distribuida,
- logs diagnósticos estructurados,
- métricas agregadas,
- dashboards, alertas y troubleshooting,
- vistas de performance, saturación, retries y errores.

La telemetría derivada puede resumir o referenciar evidencia operativa, pero **nunca** sustituirla ni ser la única fuente de reconstrucción operativa.

## Responsabilidades explícitas del event log operativo

El event log operativo del kernel DEBE:

1. registrar hechos canónicos mínimos append-only;
2. correlacionar `compiled_contract`, `execution_record`, `policy_decision_record`, approvals y results;
3. permitir reconstrucción del corredor causal por `tenant_id`, `execution_id`, `contract_id`, `trace_id` y IDs vecinos;
4. dejar evidencia durable de hitos del runtime aunque la exportación OTel falle;
5. soportar deduplicación por `event_key`/`causation_key` material;
6. persistir clasificación efectiva y resultado de redacción aplicado al evento;
7. conservar referencias a evidence refs, reason codes y fingerprints materiales;
8. habilitar queries operativas y auditoría sin depender de Temporal history ni de dashboards;
9. mantener monotonicidad temporal/lógica del corredor principal;
10. servir como base de proyección hacia observabilidad derivada.

## No-responsabilidades explícitas del event log operativo

El event log operativo NO debe:

1. reemplazar `execution_record` como snapshot/estado vivo del runtime;
2. almacenar spans nativos, histogramas ni series temporales como si fuera backend de observabilidad;
3. guardar payloads raw sensibles “por las dudas”;
4. convertirse en data lake analítico o motor de dashboards;
5. depender de Temporal search/history para completar semántica faltante;
6. absorber lógica de policy, approvals o clasificación que pertenece a otros seams;
7. permitir updates destructivos, compactaciones semánticas opacas o reescrituras del pasado;
8. deducir retrospectivamente hechos no emitidos.

## Responsabilidades explícitas del plano OTel/LGTM

El plano OTel/LGTM DEBE:

1. recibir señales derivadas desde el kernel;
2. modelar traces para compilación, evaluation gate, execution phase y application phase;
3. recibir logs estructurados con `reason_code`, `event_type`, severidad y correlation ids;
4. exponer métricas mínimas operables para throughput, errores, bloqueos, retries y `unknown_outcome`;
5. tolerar redacción/sanitización previa y operar con payload resumido o referenciado;
6. permitir alerting y troubleshooting sin introducir semántica nueva de negocio;
7. aceptar replay/proyección diferida de eventos canónicos ya persistidos.

## Event log canónico mínimo del kernel

Se adopta un objeto lógico canónico llamado **`telemetry_event`**, persistido en una estructura append-only física recomendada llamada **`event_record`**. En v1 ambos nombres refieren al mismo hecho canónico en distinto nivel:

- **`telemetry_event`**: envelope lógico del hecho emitido por el kernel.
- **`event_record`**: representación persistida append-only en PostgreSQL.

Cada `event_record` representa un único hecho material, con metadata suficiente para:

- identificar causalidad,
- deduplicar emisión,
- clasificar sensibilidad,
- proyectar a observabilidad,
- enlazar con records canónicos vecinos.

## Schema mínimo recomendado de `telemetry_event` / `event_record`

Schema lógico mínimo recomendado:

```yaml
api_version: v1
kind: telemetry_event
metadata:
  event_id: evt-...
  event_type: execution.started
  tenant_id: ten-...
  environment: prod
  occurred_at: 2026-03-29T10:00:00Z
  recorded_at: 2026-03-29T10:00:01Z
  correlation_version: v1
  classification_level: internal|confidential|restricted
  redaction_status: full|summary_only|redacted|blocked_for_export
spec:
  trace_id: trc-...
  span_id: spn-...            # nullable
  contract_id: ctr-...
  contract_fingerprint: sha256:...
  execution_id: exe-...
  approval_request_id: apr-... # nullable
  result_id: res-...           # nullable
  policy_decision_id: pol-...  # nullable
  causation_key: exe-...:execution.started:1
  idempotency_key: idem-...    # nullable para eventos no execution-scoped
  producer: compiler|runtime|policy_pep|approval_flow
  phase: compilation|policy|execution|application|compensation|approval
  state_from: created          # nullable
  state_to: executing          # nullable
  reason_code: runtime.started # nullable
  summary:
    message: Execution started
    attributes:
      capability_id: cap-...
      is_read_only: false
      has_external_effect: true
  references:
    compiled_contract_ref: cc-...
    execution_record_ref: exe-...
    policy_decision_record_ref: pol-...
    approval_request_ref: apr-...
    result_record_ref: res-...
  payload_classified:
    allowed_for_operational_storage: {}
    allowed_for_export_summary: {}
    restricted_fields_present: true
status:
  export_projection:
    export_state: pending|exported|export_failed|not_exportable
    last_export_attempt_at: null
    export_attempt_count: 0
    export_batch_id: null
```

Campos normativos obligatorios en v1:

- `event_id`
- `event_type`
- `tenant_id`
- `environment`
- `occurred_at`
- `trace_id`
- `contract_id` cuando el hecho nace de compilación/ejecución/policy
- `contract_fingerprint` cuando el hecho refiere contrato ejecutable o policy material
- `execution_id` cuando el hecho nace del corredor runtime
- `approval_request_id` cuando aplica
- `result_id` cuando aplica
- `policy_decision_id` cuando aplica
- `classification_level`
- `redaction_status`
- `causation_key`
- `producer`
- `phase`

Reglas:

1. `event_id` es único e inmutable.
2. `event_record` es append-only; correcciones se expresan como nuevos eventos de compensación o anotación, nunca update destructivo.
3. `occurred_at` refleja el hecho de dominio; `recorded_at` refleja persistencia física.
4. `payload_classified` separa explícitamente lo operativo persistible de lo exportable.
5. `export_state` NO define verdad del hecho; sólo estado de proyección.

## Event types mínimos a materializar primero

El kernel debe materializar, como mínimo, estos tipos:

### Compilación

- `contract.compilation_started`
- `contract.compilation_completed`
- `contract.compilation_failed`

### Policy

- `policy.decision_requested`
- `policy.decision_recorded`

### Approvals

- `approval.awaiting`
- `approval.released`

### Lifecycle de ejecución

- `execution.created`
- `execution.released`
- `execution.started`
- `execution.completed`
- `execution.blocked`
- `execution.failed`
- `execution.unknown_outcome`

### Application phase

- `application.released`
- `application.started`
- `application.completed`

### Compensación

- `compensation.requested`
- `compensation.completed`

Estos son los mínimos. Se pueden sumar eventos auxiliares, pero NO reemplazar estos hitos base.

## Correlación exacta de IDs y reglas de propagación

La correlación mínima obligatoria del sistema queda fijada en:

- `tenant_id`
- `environment`
- `trace_id`
- `contract_id`
- `contract_fingerprint`
- `execution_id`
- `approval_request_id`
- `result_id`
- `policy_decision_id`

### Reglas exactas

1. **Toda emisión material requiere `tenant_id` y `environment`.** Si faltan, el kernel rechaza emisión canónica.
2. **`trace_id` es obligatorio** para compilación, policy, ejecución, aplicación y activación. Si falta, el kernel debe crear/propagar uno antes de emitir, salvo que la operación se rechace por invariantes más duros.
3. **`contract_id` + `contract_fingerprint`** son obligatorios en eventos que dependan del contrato compilado o de una evaluación de policy material.
4. **`execution_id`** es obligatorio en todo evento del corredor runtime, approvals runtime-bound, application y compensación.
5. **`approval_request_id`** es obligatorio en `approval.awaiting`, `approval.released` y en cualquier policy event cuya salida causal sea approval.
6. **`result_id`** es obligatorio cuando el hecho depende de evidencia/result outcome ya materializado.
7. **`policy_decision_id`** es obligatorio en `policy.decision_recorded` y SHOULD estar presente en eventos runtime causados directamente por esa decisión.
8. **Propagación ascendente:** `compiled_contract` → `execution_record` → `policy_decision_record` / `approval_request` / `result_record` → `event_record`.
9. **Propagación transversal:** `correlation_context_propagator` debe inyectar el mismo set de IDs en spans, logs estructurados y métricas etiquetadas permitidas.
10. **No mezcla de corredores:** un `trace_id` puede contener varios spans, pero no debe mezclar ejecuciones independientes sin causal link explícito (`parent_execution_id`, `supersedes_execution_id` o equivalente).
11. **Versionado de correlación:** cambios futuros del envelope se expresan por `correlation_version`, no por reinterpretación silenciosa.

## Qué persiste en PostgreSQL

En PostgreSQL deben persistirse completos o con referencia operativa durable:

1. `compiled_contract` y `compilation_report` definidos en C.1.
2. `execution_record` y snapshots/timestamps canónicos definidos en C.2.
3. `policy_decision_record` o equivalente durable definido en C.3.
4. `approval_request` y decisiones/release relevantes según A.4.
5. `event_record` append-only como event log mínimo del kernel.
6. metadata de exportación/proyección (`export_state`, attempts, batch refs) asociada al evento o en tabla vecina.
7. references a evidence, reason codes, hashes y fingerprints materiales.

Regla v1: PostgreSQL conserva **hecho operativo canónico** y, cuando el payload bruto sea sensible, conserva sólo la versión clasificada permitida o una referencia gobernada al artefacto sensible.

## Qué se proyecta a OTel/LGTM

Hacia OTel/LGTM se proyecta solamente señal derivada:

1. spans de compilación, evaluation gate, execution phase y application phase;
2. logs estructurados con `event_type`, `reason_code`, severidad, `trace_id` y correlation ids permitidos;
3. métricas derivadas de conteo, latencia, error, retry, bloqueo y `unknown_outcome`;
4. atributos resumidos o referenciados, nunca payload operativo completo sensible;
5. eventos resumidos que permitan troubleshooting sin reconstituir secretos ni contenido restricted.

La proyección a observabilidad puede ser:

- **full summary** para eventos de baja sensibilidad,
- **summary + refs** para eventos moderados,
- **redacted summary only** para eventos sensibles,
- **blocked_for_export** para eventos restricted que no deban salir del plano operativo.

## Qué nunca debe persistirse/exportarse en crudo

Nunca debe persistirse/exportarse en crudo en el event log ni en observabilidad derivada:

- secretos, tokens, API keys, credentials;
- payloads completos de herramientas con datos sensibles innecesarios;
- contenido restricted/raw de prompts, adjuntos, snapshots o evidence sensible;
- PII/PHI/financiero identificable cuando la política exija referencia o hash;
- documentos completos de approvals o policy inputs cuando basta snapshot resumido/ref;
- cuerpos completos de errores externos que arrastren datos del tenant.

Regla dura: si un campo está clasificado como `restricted` y no existe sanitizer aprobado, **NO sale** a OTel/LGTM y tampoco entra completo al `event_record`; se reemplaza por resumen, hash, ref o marcador de bloqueo.

## Redacción y clasificación previa a persistencia/exportación

Antes de persistir/exportar, toda emisión debe pasar por `classification_redaction_guard`.

### Pipeline mínimo

1. resolver `classification_level` efectivo del hecho;
2. identificar campos `public`, `internal`, `confidential`, `restricted`;
3. construir `payload_operational` permitido para PostgreSQL;
4. construir `payload_export_summary` permitido para OTel/LGTM;
5. bloquear/exportar según policy y sanitize rules;
6. anotar `redaction_status`, versión de reglas y restricted flags aplicados.

### Regla de clasificación

- **PostgreSQL operativo** puede conservar más detalle que OTel/LGTM, pero sólo dentro del límite aprobado por clasificación.
- **OTel/LGTM** siempre recibe una vista igual o más reducida que la operativa.
- cambios de clasificación posteriores NO reescriben el hecho canónico; generan nueva proyección/redacción y, si aplica, evento adicional de reclasificación.

## Integración con `compiled_contract`

1. `contract.compilation_started` se emite al iniciar el corredor de compilación con `tenant_id`, `environment`, `trace_id` y correlación tentativa de contrato.
2. `contract.compilation_completed` sólo se emite cuando `compiled_contract` y `compilation_report` ya quedaron persistidos en PostgreSQL.
3. `contract.compilation_failed` debe persistir evidencia suficiente aunque el contrato final no exista todavía.
4. `contract_id` y `contract_fingerprint` usados por C.1 deben quedar reflejados en el `event_record` una vez materializados.
5. El event log NO recompila ni interpreta el contrato; sólo evidencia el hecho de compilación y su resultado.

## Integración con `execution_workflow`

1. `execution_workflow` sigue siendo el runtime durable de C.2, pero no reemplaza el event log canónico.
2. Cada transición material del corredor principal debe producir un `event_record` compatible con el estado proyectado a `execution_record`.
3. Orden recomendado del commit lógico:
   - transición/snapshot canónico en PostgreSQL,
   - append al `event_record`,
   - proyección asíncrona a OTel/LGTM.
4. Si falla la exportación OTel, la transición runtime sigue siendo válida si `execution_record` + `event_record` quedaron persistidos.
5. `execution.completed` no puede emitirse como terminal equivalente a `application.completed` en mutaciones; se respeta A.5.
6. `application.completed` no aplica a read-only; para esos casos el corredor cierra desde `execution.completed` según A.5/C.2.
7. `execution.unknown_outcome` debe existir como hecho canónico separado cuando hay duda material sobre aplicación.

## Integración con Cerbos / policy decisions

1. `policy.decision_requested` se emite cuando el PEP prepara una evaluación material.
2. `policy.decision_recorded` se emite sólo cuando la decisión quedó persistida como `policy_decision_record` o evidencia equivalente durable.
3. El event log debe conservar `policy_decision_id`, `policy_version` y `request_hash`/equivalente resumido cuando aplique.
4. Un deny/bloqueo de policy puede causar `execution.blocked`, pero la evidencia de policy sigue siendo evento separado.
5. Si Cerbos respondió pero falló la persistencia canónica de la decisión, NO debe considerarse éxito operativo final para mutación.

## Integración con approvals

### Approvals

1. `approval.awaiting` se emite cuando el corredor entra formalmente en espera de approval.
2. `approval.released` se emite sólo cuando existe release válido y persistido.
3. approvals no reemplazan eventos de ejecución; son hechos correlacionados pero distintos.
4. `approval_request_id` debe estar presente y ser consistente con `execution_id`/`contract_fingerprint` vigentes.

## Señales mínimas OTel (traces / logs / metrics)

### Traces

Deben existir spans mínimos para:

- compilación (`contract.compile`),
- evaluation gate (`policy.evaluate`),
- execution phase (`execution.run`),
- application phase (`application.run`).

Cada span debe llevar, cuando aplique:

- `tenant_id`
- `environment`
- `trace_id`
- `contract_id`
- `contract_fingerprint`
- `execution_id`
- `approval_request_id`
- `result_id`
- `policy_decision_id`

### Logs estructurados

Deben incluir al menos:

- `timestamp`
- `severity`
- `event_type`
- `reason_code`
- `tenant_id`
- `environment`
- `trace_id`
- IDs correlacionables aplicables
- `redaction_status`
- `export_state` cuando aplique

### Metrics mínimas

- `contract_compilations_total`
- `policy_decisions_total{decision=allow|deny|approval|escalation}`
- `execution_started_total`
- `execution_completed_total`
- `execution_failed_total`
- `execution_blocked_total`
- `application_completed_total`
- `compensation_count`
- `retry_count`
- `unknown_outcome_count`

SHOULD agregarse latencias por fase, pero las métricas mínimas v1 son los conteos anteriores.

## Estrategia de emisión y proyección

Se definen estos componentes mínimos:

1. **`event_log_writer`**
   - persiste `event_record` append-only en PostgreSQL;
   - valida envelope mínimo, correlación e idempotencia.

2. **`event_projector`**
   - consume `event_record` canónico;
   - proyecta logs/spans/metrics derivados;
   - no altera verdad operativa.

3. **`telemetry_exporter`**
   - encapsula export a OTel/LGTM;
   - soporta batch, retry técnico y reporte de falla.

4. **`correlation_context_propagator`**
   - construye/inyecta correlation ids consistentes en runtime, compiler, policy y activation.

5. **`classification_redaction_guard`**
   - calcula clasificación efectiva y vistas permitidas antes de persistencia/exportación.

6. **`observability_failure_guard`**
   - evita que fallas del plano OTel/LGTM rompan commit canónico, transición runtime o respuesta operativa.

### Estrategia recomendada v1

1. componente de dominio confirma la transición/decisión material;
2. `event_log_writer` valida y persiste `event_record` en PostgreSQL;
3. `event_projector` marca evento como pendiente de export;
4. `telemetry_exporter` proyecta a OTel/LGTM de manera asíncrona o diferida;
5. `export_state` se actualiza sin reescribir el hecho del evento.

Regla dura: la exportación es **side effect derivado**, no parte del commit semántico del kernel.

## Falla del emisor / backpressure / degradación segura

1. si OTel collector, LGTM o red fallan, el kernel conserva `event_record` y estado operativo en PostgreSQL;
2. `telemetry_exporter` puede reintentar por batch sin duplicar el evento canónico;
3. bajo backpressure, se prioriza append canónico y se degrada la exportación a modo diferido;
4. si el backlog crece, el sistema puede exportar sólo summaries mínimos antes que payloads enriquecidos;
5. si el guard detecta saturación extrema, puede marcar `export_state=export_failed` o `pending_retry` sin bloquear ejecución;
6. dashboards/alertas pueden degradarse; la verdad operativa NO;
7. si falla la escritura canónica, el export exitoso debe considerarse irrelevante para verdad del kernel y debe tratarse como señal huérfana/reconciliable.

## Idempotencia y deduplicación del event log

1. `event_record` debe tener unicidad al menos por `event_id`.
2. SHOULD existir `causation_key` material para detectar retries del mismo hecho lógico.
3. retries técnicos de exportación NO crean nuevos `event_record`.
4. retries técnicos del runtime pueden incrementar métricas, pero no duplicar hechos canónicos si el hito lógico ya fue persistido.
5. si un mismo evento se reemite por retry antes de confirmación, `event_log_writer` debe:
   - insertar una sola vez, o
   - registrar deduplicación explícita sin duplicar hecho canónico.
6. duplicados detectados deben dejar evidencia operativa auditable.
7. el event log hereda la semántica de A.5: retry técnico, retry de ejecución completa y replay de auditoría son conceptos distintos.

## Tests borde mínimos (al menos 20)

1. **mismo evento emitido dos veces por retry** → se deduplica sin duplicar `event_record`.
2. **falta `tenant_id` en emisión** → rechazo del evento canónico.
3. **falta `trace_id` en ejecución** → se crea/propaga o se rechaza según invariant definido; nunca se persiste evento runtime sin correlación.
4. **contract compilation falla antes de persistir contrato** → existe `contract.compilation_failed` con evidencia suficiente y sin `compiled_contract` falso.
5. **policy decision se registra pero OTel export falla** → queda persistido `policy_decision_record` + `event_record`; export en retry.
6. **OTel export exitoso pero falla escritura canónica** → el evento NO se considera válido operativamente; se detecta como huérfano/reconciliable.
7. **redacción obligatoria en log sensible** → el payload exportado sale resumido/redacted.
8. **campo restricted intenta salir a observabilidad sin sanitize** → `classification_redaction_guard` bloquea export.
9. **`execution.completed` emitido sin `execution.started`** → inconsistencia rechazada o marcada como violación auditable.
10. **`application.completed` para read-only** → transición/evento inválido.
11. **compensación sin `execution.failed` previo** → sólo permitido si viene de `unknown_outcome` o corredor compatible; si no, rechazo.
12. **approval release sin `approval.awaiting`** → inconsistencia causal rechazada.
14. **policy event con `policy_version` faltante** → `policy.decision_recorded` inválido.
15. **reintento técnico incrementa métricas sin duplicar event log canónico** → métricas aumentan; `event_record` no se duplica.
16. **history de Temporal diverge del event log canónico** → prevalece PostgreSQL/event log; se dispara reconciliación.
17. **export batch retrasado por backpressure** → backlog crece pero el commit canónico no se pierde.
18. **clasificación cambia y evento proyectado requiere nueva redacción** → nueva proyección resumida sin mutar el hecho original.
19. **event_record append-only intenta mutación destructiva** → operación rechazada.
20. **consulta operativa requiere reconstrucción de corredor por correlación** → puede reconstruirse usando `tenant_id`, `trace_id`, `contract_id`, `execution_id` y IDs vecinos.
21. **`execution.blocked` causado por deny de policy** → existen ambos hechos correlacionados: policy + runtime.
22. **`approval.awaiting` emitido sin `approval_request_id`** → rechazo.
23. **retry del exporter reenvía el mismo batch** → LGTM puede ver duplicado técnico controlado, pero el kernel no duplica `event_record`.
24. **`execution.unknown_outcome` seguido de retry de aplicación ciego** → rechazo por A.5; debe requerir verificación/escalación.
25. **evento con `contract_fingerprint` que no coincide con `execution_record`** → rechazo o invalidación por conflicto material.

## Criterios de aceptación de C.4

1. Queda explícito que **PostgreSQL conserva verdad operativa y event log canónico mínimo**.
2. Queda explícito que **OTel/LGTM conserva telemetría derivada y no source of truth**.
3. Queda explícito que **Temporal history no reemplaza el event log canónico del kernel**.
4. Existe `telemetry_event` / `event_record` canónico mínimo con schema implementable.
5. Queda fijada la correlación mínima con `tenant_id`, `environment`, `trace_id`, `contract_id`, `contract_fingerprint`, `execution_id`, `approval_request_id`, `result_id` y `policy_decision_id`.
6. Existe un catálogo mínimo de event types para compilación, policy, approvals, ejecución, application phase y compensación.
7. Queda definido qué persiste completo/referenciado en PostgreSQL y qué se proyecta resumido/referenciado a OTel/LGTM.
8. Queda explícito qué nunca debe persistirse/exportarse en crudo.
9. Queda definida clasificación/redacción previa a persistencia/exportación.
10. Quedan definidos `event_log_writer`, `event_projector`, `telemetry_exporter`, `correlation_context_propagator`, `classification_redaction_guard` y `observability_failure_guard`.
11. Queda definida estrategia de emisión/proyección con degradación segura y backpressure.
12. Quedan fijadas append-only, idempotencia y deduplicación del event log.
13. La integración con `compiled_contract`, `execution_workflow`, Cerbos y approvals queda descrita sin contradecir A.5, B.3, C.2 ni C.3.
14. Existe batería mínima de tests borde suficiente para construcción operativa v1.
