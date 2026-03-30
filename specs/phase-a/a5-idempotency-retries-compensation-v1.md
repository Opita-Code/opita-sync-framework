# A.5 — Idempotencia, reintentos y compensación v1

## Principios base

- Esta especificación define las reglas runtime-generales de idempotencia, deduplicación, reintento y compensación para `execution_record`.
- NO reemplaza reglas específicas de subsistemas. En particular, A.4 ya define idempotencia y reintentos del subsistema de approvals; este documento define la capa superior del runtime general.
- La unidad canónica para estas reglas es la **ejecución concreta** (`execution_record`), no el contrato abstracto, no el approval request y no el resultado aislado.
- La semántica base del core se conserva: `intent_contract -> acción de negocio -> tools/sistemas base`.
- A.5 debe respetar la separación ya establecida entre `execution_completed` y `application_completed`. Reintentar o compensar una corrida técnica NO equivale a re-aplicar un efecto externo.
- A.3 sigue vigente: `mutation` y `read-only` tienen obligaciones distintas; `system_update` irreversible NO soporta tratamiento laxo de `partial_success` como si fuera seguro volver a aplicar.
- `rollback` físico y `compensación lógica` NO son sinónimos. El primero intenta revertir el efecto material; la segunda restaura consistencia operacional o de negocio aunque el efecto físico no pueda deshacerse exactamente.
- Toda decisión de deduplicación, retry, compensación o escalación debe dejar trazabilidad auditable vía `execution_record` + `telemetry_event`.
- Debe existir manejo explícito de **unknown outcome** cuando un sistema externo no confirma si aplicó o no aplicó un efecto después de timeout, corte de red o pérdida de confirmación.
- No existe “replay para auditoría” con efectos reales. Replay de auditoría es reconstrucción/verificación de evidencia, no una nueva aplicación.

---

## Definición formal de idempotencia en Opyta Sync

En Opyta Sync, una operación es **idempotente** si múltiples intentos materialmente equivalentes sobre el mismo scope producen un único efecto operacional válido y una única interpretación canónica del outcome, sin multiplicar efectos externos no deseados.

Formalmente, dadas dos solicitudes `R1` y `R2`, se consideran equivalentes a efectos de idempotencia runtime si coinciden al menos en:

- `tenant_id`
- `capability_id`
- `contract_fingerprint`
- `target_scope`
- `operation_phase`

Y además:

- operan dentro de la misma ventana de validez de deduplicación,
- no existe cambio material del contrato, policy o approvals,
- no cambia el modo efectivo de aprobación ni la clasificación material,
- no cambia el tipo de efecto externo esperado.

Bajo esas condiciones, el runtime DEBE garantizar una de estas salidas canónicas:

1. **Return existing outcome**: devolver referencia a la ejecución ya existente o a su resultado terminal compatible.
2. **Attach as duplicate attempt**: registrar el intento como duplicado sin volver a aplicar el efecto.
3. **Resume safely**: continuar una ejecución interrumpida si el diseño de la operación lo permite y el estado es reanudable.
4. **Reject retry/replay**: rechazar el intento si la naturaleza de la operación hace inseguro repetirla.

Idempotencia en este motor NO significa “misma respuesta byte a byte”; significa **misma interpretación canónica del efecto permitido** bajo las mismas condiciones materiales.

---

## Qué operaciones deben ser idempotentes por diseño

Deben ser idempotentes por diseño runtime o por contrato operativo:

1. **Creación de `execution_record`** para una misma intención material y misma fase operativa.
2. **Registro de eventos canónicos** cuando el productor reemite por retry técnico, usando `event_id`/correlación estable para deduplicar ingestión.
3. **Operaciones read-only** (`plan`, `inspection`, `query`, `report`) cuando se reejecutan sobre el mismo contrato material.
4. **Consultas de estado** a sistemas externos.
5. **Fetch de evidencia**, polling de confirmación y verificación post-aplicación.
6. **Release checks** y validaciones previas sin efecto externo.
7. **Reanudación de pasos técnicos internos** que escriben solo en estado append-only controlado.
8. **Compensaciones lógicas declarativas** que marcan remediación, invalidación, supersession o cierre controlado sin duplicar side effects.
9. **Replay de auditoría** sobre trazas, snapshots y evidencia, siempre que se ejecute en modo no-efectivo.

Norma: si una operación puede diseñarse como idempotente y el documento de capability no la declara así, se considera un déficit de diseño.

---

## Qué operaciones no pueden ser idempotentes pero deben ser detectadas/controladas

Hay operaciones cuyo efecto material puede no admitir repetición segura. El runtime NO puede fingir idempotencia donde no existe.

Casos típicos:

1. `system_update` con `external_effect = irreversible`.
2. Mutaciones que disparan side effects externos no deduplicables por el sistema destino.
3. Operaciones con semántica “append irreversible” en terceros: envío de email definitivo, transferencia, publicación externa, rotación destructiva, disparo de workflow irreversible.
4. Escrituras sobre APIs sin soporte de idempotency token ni lectura confiable de estado final.
5. Acciones cuyo segundo intento crea un recurso adicional indistinguible del primero.
6. Pasos donde el timeout deja outcome incierto y el sistema externo no provee confirmación transaccional.

Para estas operaciones el runtime DEBE:

- detectarlas antes de ejecutar,
- marcarlas como `non_idempotent_risk` a nivel de ejecución/paso,
- exigir policy de retry ultra conservadora,
- bloquear auto-retry si hay riesgo real de doble aplicación,
- exigir verificación externa antes de cualquier nuevo intento,
- escalar a compensación lógica o intervención manual si el outcome queda incierto.

---

## Objeto o envelope recomendado para `idempotency_key` / deduplication metadata dentro de `execution_record`

`execution_record` DEBE tener `idempotency_key` a nivel ejecución y metadata de deduplicación explícita.

Envelope recomendado:

```yaml
idempotency:
  idempotency_key: idem-...
  dedup_fingerprint: sha256:...
  dedup_window:
    starts_at: 2026-03-29T10:00:00Z
    expires_at: 2026-03-29T16:00:00Z
  dedup_scope:
    tenant_id: ten-...
    environment: prod
    capability_id: cap-...
    contract_fingerprint: sha256:...
    target_scope: dst:crm/account:123|fields:status,tier
    operation_phase: execution|application|compensation|verification|audit_replay
  dedup_policy:
    mode: strict
    return_existing_execution: true
    allow_resume_in_place: false
    allow_new_execution_if_terminal: conditional
  causal_links:
    deduped_against_execution_id: null
    parent_execution_id: null
    supersedes_execution_id: null
    replay_of_execution_id: null
  attempt_counters:
    technical_step_attempt_count: 0
    execution_attempt_count: 1
    audit_replay_count: 0
  risk_flags:
    non_idempotent_risk: false
    unknown_outcome_risk: false
    external_confirmation_required: false
```

### Reglas del envelope

- `idempotency_key` identifica la intención de deduplicación de la corrida.
- `dedup_fingerprint` puede derivarse de una serialización determinística de los campos materiales de deduplicación.
- `target_scope` DEBE representar el alcance material sobre el que el efecto operaría; no puede ser texto libre ambiguo.
- `operation_phase` es obligatorio para distinguir reintentos de ejecución técnica, aplicación material, compensación y replay.
- `attempt_counters.technical_step_attempt_count` NO reemplaza `execution_attempt_count`.
- `replay_of_execution_id` solo aplica a auditoría y nunca habilita efectos reales.

---

## Reglas de deduplicación

### Clave mínima obligatoria de deduplicación

La deduplicación runtime DEBE considerar como mínimo:

- `tenant_id`
- `capability_id`
- `contract_fingerprint`
- `target_scope`
- `operation_phase`

Y SHOULD considerar además cuando aplique:

- `environment`
- `approval_mode_effective`
- `result_type`
- `external_effect`
- `destination_snapshot` fingerprint
- `subject/delegation` cuando afecten autoridad material

### Reglas normativas

1. Si cambia cualquiera de los campos mínimos, NO es el mismo grupo de deduplicación.
2. Si el `contract_fingerprint` cambia, toda aprobación previa puede quedar invalidada y la deduplicación debe reiniciarse.
3. Dos operaciones en distinta `operation_phase` NUNCA se deduplican entre sí aunque compartan el resto de los campos.
4. `execution` y `application` deben poder deduplicarse por separado.
5. Una corrida en `audit_replay` NO puede colisionar con una corrida efectiva.
6. Si una ejecución previa está terminalmente consistente y el nuevo intento es materialmente equivalente, el runtime debe preferir **return existing outcome** sobre crear una nueva aplicación.
7. Si la ejecución previa quedó en estado no terminal pero reanudable y la policy lo permite, el runtime puede **resume safely**; de lo contrario debe crear nueva ejecución vinculada causalmente.
8. Si la ejecución previa quedó con `unknown_outcome`, la deduplicación debe entrar en modo defensivo: no se autoriza re-aplicación hasta verificar estado externo o escalar.
9. Duplicados detectados deben registrarse como evidencia; no se descartan silenciosamente.
10. La deduplicación es por tenant y ambiente; nunca debe colapsar ejecuciones cross-tenant.

### Ventana de deduplicación

- Para read-only: puede ser amplia y configurable.
- Para mutation reversible: debe cubrir al menos la ventana razonable en la que un timeout tardío todavía puede materializar el efecto.
- Para mutation irreversible: debe ser conservadora y preferentemente asociada a retención duradera de evidencia, no solo TTL corto.
- Para replay de auditoría: ventana propia, separada de la operativa.

---

## Estrategia general de reintentos

El runtime distingue TRES conceptos que NO deben mezclarse:

### 1. Retry del paso técnico

Es la repetición de una sub-operación técnica interna o de una llamada puntual a tool/sistema base dentro de la misma ejecución lógica.

Ejemplos:

- reintentar un GET a una API,
- reintentar polling de confirmación,
- reintentar persistencia de telemetría,
- reintentar una llamada idempotente al conector.

### 2. Retry de la ejecución completa

Es una nueva corrida operacional respecto de una ejecución previa fallida, bloqueada o inconclusa. Debe quedar reflejada con **nuevo `execution_id`** y vínculo causal al intento anterior, en línea con A.5 runtime states.

### 3. Replay de auditoría

Es la reproducción o reconstrucción de evidencia, trazas y decisiones para auditoría, debugging o compliance. NO debe ejecutar side effects reales, NO debe usar canales efectivos de aplicación y NO debe confundirse con retry.

### Principios de retry

- Retry no es derecho automático; es una decisión gobernada por tipo de operación, riesgo, fase, evidencia y estado conocido del sistema destino.
- Primero verificar, después reintentar. Especialmente en mutaciones.
- Cuanto más irreversible sea el efecto, más fuerte debe ser la preferencia por **verify-before-retry**.
- Si existe posibilidad de doble aplicación, el default es **no auto-retry**.
- Unknown outcome post-timeout fuerza verificación externa, compensación lógica o intervención manual; nunca retry ciego.

---

## Política de reintentos por tipo de operación

### `read-only`

- Política: **agresiva**.
- Se permite auto-retry de pasos técnicos y, si corresponde, de la ejecución completa.
- Se prioriza disponibilidad y completitud por encima de costo incremental moderado.
- Si el resultado cambia entre intentos por volatilidad normal del sistema fuente, eso NO rompe idempotencia mientras no haya side effects.
- El runtime puede devolver último resultado consistente o regenerarlo según freshness policy.

### `mutation reversible`

- Política: **cautelosa**.
- Se permite retry técnico solo si el paso es idempotente o si la plataforma puede verificar que el efecto no fue aplicado todavía.
- Retry de ejecución completa requiere nuevo `execution_id` y referencia causal.
- Debe existir plan de rollback o compensación definido antes de aplicar.
- Si el outcome es incierto, detener auto-retry y pasar a verificación/compensación.

### `mutation irreversible`

- Política: **ultra conservadora**.
- `system_update` irreversible NO debe auto-reintentarse si existe riesgo de doble aplicación.
- Solo se permite retry técnico en pasos previos no-efectivos: validación, simulación, lectura, polling, telemetría, preparación.
- Una vez emitido el side effect irreversible o si no puede descartarse que haya sido emitido, el runtime debe bloquear auto-retry de aplicación.
- Requiere verificación externa fuerte; si no existe, escalar a intervención manual.

### `governance`

- Política: **determinista y conservadora**.
- Reintentos de evaluación policy/read pueden ser automáticos.
- Reintentos que impliquen reemitir una decisión de governance efectiva deben respetar snapshot material y no duplicar efectos jurídicos/operativos.
- Si una decisión ya fue emitida y registrada, un intento equivalente debe deduplicarse o supersederse explícitamente; nunca duplicarse de forma ambigua.

---

## Errores retryable vs non-retryable

## Taxonomía mínima `retryable_errors`

`retryable_errors` son errores donde repetir controladamente puede resolver la falla sin introducir riesgo desproporcionado de doble efecto.

Clases mínimas:

1. `transient_network_error`
2. `connection_reset`
3. `dns_resolution_transient`
4. `gateway_timeout`
5. `upstream_timeout_before_commit_confirmation`
6. `rate_limited_retry_after`
7. `temporary_service_unavailable`
8. `optimistic_lock_retryable`
9. `ephemeral_connector_session_expired`
10. `telemetry_ingest_transient`
11. `polling_inconclusive_but_safe`
12. `platform_internal_transient`

Condiciones para clasificarlos como retryable:

- el paso es read-only o idempotente,
- o existe prueba suficiente de no-aplicación,
- o el retry ocurre antes de la fase efectiva de aplicación,
- o el sistema destino provee token/idempotency semantics verificables.

## Taxonomía mínima `non_retryable_errors`

`non_retryable_errors` son errores donde reintentar automáticamente podría repetir un efecto, violar governance o persistir en una condición estructuralmente inválida.

Clases mínimas:

1. `contract_fingerprint_mismatch`
2. `approval_missing_or_invalid`
3. `approval_superseded`
4. `policy_blocked`
5. `scope_violation`
6. `classification_violation`
7. `sod_violation`
8. `validation_error_non_transient`
9. `unsupported_operation`
10. `irreversible_effect_may_have_been_applied`
11. `unknown_outcome_requires_verification`
12. `duplicate_non_idempotent_application_risk`
13. `manual_hold_required`
14. `compensation_precondition_missing`

Norma: un timeout en mutación NO es automáticamente retryable. Si existe duda sobre aplicación parcial o completa, debe reclasificarse como `unknown_outcome_requires_verification`.

---

## Backoff y límites recomendados

### Reglas generales

- Usar backoff exponencial con jitter.
- Respetar `Retry-After` si el upstream lo provee.
- Limitar por tipo de operación, fase y riesgo.
- El presupuesto de retry debe ser observable y auditable.

### Recomendaciones v1

#### Read-only
- intentos técnicos automáticos: 3 a 6
- backoff inicial: 250ms a 1s
- máximo entre intentos: 15s
- jitter: 10% a 30%

#### Mutation reversible
- intentos técnicos automáticos sobre pasos seguros: 1 a 3
- backoff inicial: 1s a 3s
- máximo entre intentos: 30s
- requerir verificación entre intentos si el paso toca aplicación

#### Mutation irreversible
- intentos técnicos automáticos en pasos pre-application: 1 a 2
- aplicación efectiva: 0 auto-retries si hay riesgo de doble aplicación
- verificación post-timeout: polling acotado y luego escalación

#### Governance
- policy checks/read calls: 2 a 4
- emisión efectiva de decisión: 0 a 1 según idempotencia garantizada del subsistema

### Límites duros

1. No encadenar retries indefinidos.
2. No reintentar ejecución completa dentro del mismo `execution_id`.
3. No auto-reintentar aplicación irreversible ante unknown outcome.
4. No usar backoff para “esperar magia” en errores de validación/policy.

---

## Reglas de compensación vs rollback

### Rollback físico

Rollback físico es la acción que busca revertir el cambio material en el sistema afectado, restaurando el estado previo o equivalente cercano.

Ejemplos:

- revertir un cambio de configuración,
- restaurar un valor previo,
- eliminar un recurso recién creado si la plataforma destino lo soporta.

### Compensación lógica

Compensación lógica es una acción correctiva posterior que restablece consistencia operacional, contractual, de negocio o de auditoría cuando el rollback físico total no existe, no es seguro o no alcanza.

Ejemplos:

- marcar una operación como revertida lógicamente aunque el sistema externo conserve huella histórica,
- emitir una acción compensatoria inversa en otro sistema,
- invalidar resultados posteriores dependientes,
- generar tarea obligatoria de remediación manual,
- congelar nuevas ejecuciones sobre el scope afectado.

### Reglas normativas

1. Rollback físico y compensación lógica pueden coexistir.
2. Si no hay rollback físico posible, DEBE existir compensación lógica o escalación/manual intervention explícita.
3. Compensar NO borra la historia ni la falla original.
4. La compensación debe quedar evidenciada y correlacionada causalmente con la ejecución afectada.
5. `failed` no implica automáticamente `compensated`; debe cumplirse el procedimiento y registrarse evidencia.
6. `partially_compensated` es obligatorio cuando parte del daño/remanente sigue abierto.

---

## Qué significa compensación lógica en este motor

En Opyta Sync, compensación lógica significa que el motor declara y evidencia una secuencia de acciones destinada a restaurar **consistencia suficiente** aunque el mundo externo no vuelva exactamente al estado previo.

La compensación lógica puede incluir una o más de estas dimensiones:

- **consistencia de negocio**: restaurar el outcome esperado para el tenant o el usuario,
- **consistencia de workflow**: impedir que flujos posteriores asuman un éxito inválido,
- **consistencia de governance**: registrar excepción, revocación, bloqueo o supersession,
- **consistencia de auditoría**: dejar claro qué pasó, qué quedó incierto y qué remanente sigue abierto,
- **consistencia operativa**: abrir remediación manual, alerta, quarantine o hold sobre el scope impactado.

Una compensación lógica es válida cuando:

1. existe reason code de por qué compensar,
2. existe procedimiento definido o manual intervention formalizada,
3. existe evidencia del efecto compensatorio logrado o del remanente aceptado,
4. el estado terminal resultante deja clara la deuda residual.

---

## Reglas por tipo de resultado (`plan`, `inspection`, `query`, `report`, `change_proposal`, `execution`, `system_update`, `governance_decision`)

### `plan`

- read-only puro
- retry agresivo permitido
- deduplicación preferente por contrato + scope
- compensación normalmente `not_required`

### `inspection`

- read-only
- retry agresivo permitido
- si la evidencia quedó incompleta, se puede regenerar
- compensación no aplica salvo saneamiento de clasificación/auditoría

### `query`

- read-only
- retry agresivo permitido
- si la consulta es eventual-consistent, documentar freshness
- compensación no requerida

### `report`

- read-only con evidencia
- retry agresivo permitido
- replay de auditoría especialmente relevante
- compensación solo si hubo entrega/clasificación errónea que exige invalidación o redistribución corregida

### `change_proposal`

- mutation semántica pero sin aplicación externa real
- retry generalmente seguro
- deduplicación por proposal scope + fingerprint
- compensación suele ser lógica: supersede, invalidate, withdraw

### `execution`

- mutation general
- retry depende de si hubo side effects reales
- si solo hubo ejecución técnica sin aplicación, puede reintentarse con más libertad
- si hubo unknown outcome sobre efectos, bloquear retry ciego
- compensación puede ser rollback físico o lógica según alcance

### `system_update`

- mutation de mayor sensibilidad
- reversible: retry con extrema cautela y verificación
- irreversible: NO auto-retry si hay riesgo de doble aplicación
- `partial_success` de A.3 no habilita relanzar aplicación irreversible como si nada hubiera pasado
- si no puede saberse si aplicó, pasar a `unknown outcome` + verificación + compensación/manual

### `governance_decision`

- efecto gobernado
- reintentar evaluación puede ser seguro
- reemitir decisión efectiva exige deduplicación estricta y snapshots materiales iguales
- compensación típica: revoke, supersede, expire, exception record, hold operacional

---

## Estados o flags de compensación en `execution_record`

Además del `execution_state`, el `execution_record` debe mantener flags de compensación explícitos.

Envelope recomendado:

```yaml
compensation_status:
  is_compensation_required: false
  compensation_state: not_required
  compensation_reason_code: null
  compensation_strategy: none
  rollback_available: false
  rollback_attempted: false
  rollback_succeeded: false
  logical_compensation_available: false
  logical_compensation_applied: false
  manual_intervention_required: false
  residual_risk_level: none
  residual_effect_summary: null
  unknown_outcome: false
```

### Valores mínimos de `compensation_state`

- `not_required`
- `required`
- `in_progress`
- `partially_compensated`
- `compensated`
- `manual_intervention_pending`
- `not_possible`

### Regla de transición de estados terminales

Una ejecución pasa a:

- `compensated` cuando la remediación definida quedó completada y el remanente aceptable es nulo o explícitamente resuelto.
- `partially_compensated` cuando se mitigó una parte relevante pero persiste remanente, deuda o efecto no reversible completo.
- `failed` cuando la ejecución no completó el paso requerido y todavía no existe compensación suficiente o el remanente no fue aceptado/resuelto.

Norma adicional:

- `failed` puede coexistir históricamente con evidencia de compensación posterior, pero el estado runtime terminal visible debe seguir la transición formal definida por A.5 runtime states.

---

## Eventos canónicos de reintento/compensación

Eventos mínimos canónicos:

- `execution.idempotency_key_assigned`
- `execution.duplicate_detected`
- `execution.dedup_returned_existing`
- `execution.retry_scheduled`
- `execution.retry_started`
- `execution.retry_exhausted`
- `execution.step_retry_scheduled`
- `execution.step_retry_started`
- `execution.step_retry_succeeded`
- `execution.unknown_outcome_detected`
- `execution.external_confirmation_started`
- `execution.external_confirmation_succeeded`
- `execution.external_confirmation_inconclusive`
- `execution.compensation_required`
- `execution.rollback_started`
- `execution.rollback_succeeded`
- `execution.rollback_failed`
- `execution.logical_compensation_started`
- `execution.logical_compensation_applied`
- `execution.manual_intervention_required`
- `execution.partially_compensated`
- `execution.compensated`

Regla: estos eventos complementan, no sustituyen, los eventos generales ya definidos en A.5 runtime states.

---

## Payload mínimo por evento

Todo evento canónico de esta sección DEBE incluir como mínimo:

- `event_id`
- `event_type`
- `occurred_at`
- `tenant_id`
- `environment`
- `trace_id`
- `execution_id`
- `intent_contract_id`
- `capability_id`
- `contract_fingerprint`
- `operation_phase`
- `attempt_number`
- `technical_step_attempt_number` cuando aplique
- `idempotency_key`
- `dedup_fingerprint`
- `target_scope`
- `reason_code`
- `reason_detail`
- `retryable` bool cuando aplique
- `unknown_outcome` bool
- `compensation_required` bool
- `manual_intervention_required` bool
- `evidence_ref` o `external_reference` cuando exista

Payload adicional recomendado:

- `parent_execution_id`
- `replay_of_execution_id`
- `deduped_against_execution_id`
- `tool_call_id`
- `connector_id`
- `upstream_status_code`
- `retry_after_ms`
- `backoff_ms`
- `residual_risk_level`
- `rollback_plan_ref`
- `compensation_strategy`

---

## Taxonomía mínima adicional

### `compensation_required`

Debe activarse como verdadero al menos en estos casos:

1. hubo aplicación parcial de una mutación reversible,
2. hubo side effect confirmado pero el outcome global quedó fallido,
3. el resultado aplicado violó criterio de éxito material,
4. hubo rollback físico incompleto,
5. hubo unknown outcome con riesgo material que obliga contención,
6. una decisión de governance posterior exige deshacer, invalidar o mitigar lo ya aplicado.

### `manual_intervention_required`

Debe activarse como verdadero al menos en estos casos:

1. no hay confirmación confiable del estado externo,
2. no existe rollback físico posible y la compensación lógica automática no alcanza,
3. hay side effect irreversible potencialmente duplicado,
4. falta evidencia mínima para cerrar con seguridad,
5. el sistema destino exige acción humana fuera del motor,
6. hay conflicto entre fuentes de verdad externas que el motor no puede resolver determinísticamente.

---

## Reglas específicas para `unknown outcome` después de timeout en sistemas externos

`unknown outcome` ocurre cuando el runtime pierde certeza material sobre si el sistema externo aplicó, no aplicó o aplicó parcialmente el efecto.

### Causas típicas

- timeout esperando confirmación final,
- corte de red luego de enviar request,
- respuesta upstream perdida o corrupta,
- conector cae después de emitir el comando pero antes de persistir evidencia,
- API externa retorna estado ambiguo o asíncrono sin confirmación final.

### Reglas normativas

1. Timeout en mutación NO debe mapearse automáticamente a `failed` retryable.
2. El runtime DEBE marcar `unknown_outcome = true` en `execution_record`.
3. Debe emitirse `execution.unknown_outcome_detected`.
4. Debe iniciarse fase de `external_confirmation` si existe canal de verificación confiable.
5. Mientras el outcome sea desconocido, la aplicación efectiva queda en hold y NO se debe auto-reintentar si hay riesgo de doble aplicación.
6. Si la verificación confirma “no aplicado”, el runtime puede habilitar retry según policy.
7. Si la verificación confirma “aplicado”, el runtime debe avanzar según corresponda a `application_completed`, `failed`, o compensación, pero NO re-aplicar.
8. Si la verificación confirma “aplicado parcialmente”, debe marcar `compensation_required = true`.
9. Si la verificación sigue inconclusa tras el presupuesto de polling, debe activarse `manual_intervention_required = true` y dejarse remanente explícito.
10. En `system_update` irreversible, unknown outcome bloquea auto-retry de aplicación por defecto.

---

## Tests borde mínimos (al menos 18)

### T1 — Duplicado exacto read-only devuelve ejecución existente
Dado mismo `tenant_id`, `capability_id`, `contract_fingerprint`, `target_scope` y `operation_phase = execution` para `query`.
Esperado: no crea nueva aplicación; emite `execution.duplicate_detected` y referencia ejecución previa.

### T2 — Mismo contrato pero distinta `operation_phase`
Dado misma clave material salvo `operation_phase = audit_replay`.
Esperado: NO deduplica contra ejecución efectiva.

### T3 — Cambio de `contract_fingerprint`
Dado mismo scope pero fingerprint distinto.
Esperado: nuevo grupo de deduplicación; aprobaciones previas pueden quedar superseded.

### T4 — Retry técnico read-only por `gateway_timeout`
Dado `inspection` con timeout transiente antes de resultado.
Esperado: retry automático con backoff y jitter.

### T5 — Retry técnico reversible antes de aplicación
Dado `execution` reversible fallando en validación de conector antes de aplicar.
Esperado: retry técnico permitido dentro del mismo `execution_id`.

### T6 — Retry de ejecución completa crea nuevo `execution_id`
Dada ejecución previa fallida terminalmente.
Esperado: nuevo `execution_id` con `parent_execution_id`/causal link.

### T7 — Replay de auditoría no ejecuta side effects
Dado `report` previo reproducido para auditoría.
Esperado: solo reconstrucción/evidencia; ninguna tool efectiva de aplicación.

### T8 — `system_update` irreversible con timeout post-envío
Dado timeout luego de enviar request al sistema externo.
Esperado: `unknown_outcome = true`, NO auto-retry de aplicación.

### T9 — Verificación externa confirma no aplicado
Dado unknown outcome y canal confiable de lectura.
Esperado: puede habilitarse retry posterior según policy.

### T10 — Verificación externa confirma aplicado
Dado unknown outcome y lectura externa confirma efecto.
Esperado: no re-aplica; actualiza estado según outcome real.

### T11 — Verificación externa sigue inconclusa
Dado polling agotado sin certeza.
Esperado: `manual_intervention_required = true`.

### T12 — Aplicación parcial reversible
Dado cambio reversible aplicado sobre parte del scope.
Esperado: `compensation_required = true` y eventual `partially_compensated` si queda remanente.

### T13 — Rollback físico exitoso total
Dado failure post-application con rollback completo.
Esperado: `compensated` con evidencia de rollback.

### T14 — Rollback físico imposible pero compensación lógica posible
Dado efecto irreversible en tercero con medida correctiva lógica disponible.
Esperado: no rollback físico; sí compensación lógica o remediación formal.

### T15 — Compensación parcial con deuda residual
Dado remanente no reversible completo.
Esperado: `partially_compensated`, no `compensated`.

### T16 — Error de policy no retryable
Dado `policy_blocked`.
Esperado: sin auto-retry; clasificación `non_retryable_errors`.

### T17 — Error de approval superseded
Dada aprobación invalidada por cambio material.
Esperado: bloqueo, no retry automático hasta nueva aprobación.

### T18 — Duplicado en mutación irreversible detectado antes de aplicar
Dado segundo intento equivalente antes de side effect.
Esperado: dedup/return-existing; no doble envío.

### T19 — Duplicado detectado después de `application_completed`
Dado intento equivalente posterior sobre mutación ya aplicada.
Esperado: deduplicación defensiva y retorno de outcome existente o rechazo explícito.

### T20 — Error transiente de telemetría no repite aplicación
Dada aplicación ya confirmada y falla al registrar evento.
Esperado: retry solo de telemetría, nunca de side effect.

### T21 — Governance decision reemitida con mismo snapshot
Dada misma decisión efectiva ya emitida.
Esperado: deduplicación o supersession explícita; no duplicado ambiguo.

### T22 — Unknown outcome en reversible con lectura concluyente de no cambio
Dado timeout en aplicación reversible y verificación posterior negativa.
Esperado: retry puede habilitarse bajo nueva evidencia.

### T23 — Unknown outcome en irreversible con ausencia de canal de lectura
Dado sistema externo sin API de confirmación.
Esperado: `manual_intervention_required = true`, auto-retry bloqueado.

### T24 — Retry agotado en read-only
Dado múltiples fallas transientes consecutivas.
Esperado: `execution.retry_exhausted` con evidencia completa.

---

## Criterios de aceptación del bloque

1. Existe `idempotency_key` a nivel de `execution_record`.
2. La deduplicación considera al menos `tenant_id`, `capability_id`, `contract_fingerprint`, `target_scope` y `operation_phase`.
3. Queda normativamente diferenciada la semántica entre retry del paso técnico, retry de la ejecución completa y replay de auditoría.
4. Read-only queda definido como agresivamente reintentable.
5. Mutation reversible queda definida como reintentable solo bajo verificación y cautela.
6. Mutation irreversible queda definida con reglas ultra conservadoras.
7. `system_update` irreversible NO se auto-reintenta cuando hay riesgo de doble aplicación.
8. Se distingue formalmente rollback físico de compensación lógica.
9. Si no hay rollback físico posible, existe compensación lógica o `manual_intervention_required` explícito.
10. Queda definido cuándo una ejecución pasa a `compensated`, `partially_compensated` o `failed`.
11. Existe taxonomía mínima para `retryable_errors`, `non_retryable_errors`, `compensation_required` y `manual_intervention_required`.
12. `unknown outcome` post-timeout en sistemas externos queda normado con verificación, bloqueo de retry ciego y escalación cuando corresponda.
13. La observabilidad mínima incluye trazas, tool calls, errores, reintentos y tasa de rollback/compensación coherente con source-of-truth.
14. La especificación no contradice A.3, A.4, A.5 runtime states ni source-of-truth del core.

---

## Decisiones cerradas por esta versión

- Se adopta `idempotency_key` obligatorio a nivel ejecución.
- Se adopta deduplicación phase-aware; `execution`, `application`, `compensation` y `audit_replay` no colisionan entre sí.
- Se adopta tratamiento explícito de `unknown outcome` como estado de riesgo operacional, no como simple timeout retryable.
- Se adopta prioridad de `verify-before-retry` para mutaciones y prohibición práctica de auto-retry ciego para irreversibles.
- Se adopta `manual_intervention_required` como salida normativa válida cuando el motor no puede cerrar la verdad material del efecto.
