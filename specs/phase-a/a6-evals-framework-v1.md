# A.6 — Evals framework v1

## Principios base

- Un eval en Opyta Sync NO es un unit test tradicional. Es una validación del comportamiento esperado del motor frente a un caso controlado, con input, contexto, oracle, evidencia y criterio de aceptación explícitos.
- El eval valida comportamiento observable del sistema, no implementación interna. Si el comportamiento correcto se conserva, el eval debe seguir pasando aunque cambie la arquitectura interna.
- Los evals deben respetar las verdades previas del sistema: A.3 define éxito formal y `result_outcome`; A.4 define lifecycle, SLA, fallback y bordes de approvals; A.5 define runtime general, idempotencia y compensación; el source-of-truth fija las capas mínimas de captura y métricas de calidad/costo.
- Todo eval debe dejar evidencia auditable, trazable y correlacionable con artifacts del sistema. No existe eval “verde” sin referencias verificables a lo evaluado.
- El outcome de un eval NO se reduce a pass/fail bool. Debe poder medir outcome estructurado, score, severidad, fallas duras, fallas blandas y warnings.
- Debe existir `eval_case_id` estable, durable y reutilizable entre corridas, fixtures, regresiones y releases.
- Los evals deben poder correrse en modo advisory o gating sin redefinir el caso. Cambia el modo de enforcement; no cambia la verdad esperada.
- Los evals deben ser reproducibles sobre fixtures controladas y también auditables cuando se corren contra artifacts reales snapshotteados.
- El framework debe distinguir con precisión entre defecto de producto, limitación esperada, bloqueo correcto por governance y ruido observacional.

---

## Qué es un eval en Opyta Sync y qué NO es

### Qué es

Un eval es un caso formalizado que verifica si el motor, dado un contrato/approval/runtime/result controlado o reconstruido, produce el comportamiento esperado y deja la evidencia correcta.

Ese comportamiento esperado puede abarcar, según la familia:

- outcome de resultado por tipo,
- bloqueo correcto por policy o clasificación,
- coherencia entre approvals y releases,
- consistencia entre contrato compilado, ejecución y resultado,
- resiliencia ante retries, deduplicación, unknown outcome y compensación.

### Qué NO es

- No es un test unitario de una función aislada.
- No es un benchmark libre sin oracle formal.
- No es una revisión humana ad hoc sin fixture, expectativa y evidencia.
- No es solo un smoke test de “no crasheó”.
- No es observabilidad pasiva; requiere juicio explícito contra una expectativa definida.

---

## Taxonomía general de evals

Familias mínimas obligatorias v1:

- `result_eval` — valida calidad formal y outcome de resultados tipados según A.3.
- `policy_eval` — valida decisiones de policy, floors, overrides, restricciones y reason codes.
- `approval_eval` — valida lifecycle, releases, invalidación, SLA, fallback y consistencia de approvals según A.4.
- `classification_eval` — valida clasificación compilada, redacción parcial, bloqueo total y consistencia entre clasificación y output entregado.
- `runtime_eval` — valida consistencia orquestadora entre `intent_contract`, `approval_request/decision`, `execution_record`, `result_record` y eventos.
- `resilience_eval` — valida idempotencia, deduplicación, retries, unknown outcome y compensación según A.5.

Subtipos recomendados:

- `golden_case` — fixture canónica que no debería cambiar salvo decisión explícita.
- `edge_case` — borde diseñado para detectar regresiones conceptuales.
- `policy_regression` — caso histórico de rule enforcement.
- `cross_artifact_consistency` — caso centrado en correlación y trazabilidad.
- `operational_resilience` — caso centrado en runtime degradado o incierto.

---

## `eval_case` como objeto lógico (no necesariamente first-class top-level) con schema sugerido

`eval_case` es la unidad lógica del framework. Puede existir embebido en catálogos, suites, manifests o tooling, pero conceptualmente debe tener identidad estable.

Schema lógico sugerido v1:

```yaml
eval_case:
  id: eval.result.plan.success.baseline.v1
  family: result_eval
  title: plan básico exitoso con evidencia mínima
  intent: validar que un plan válido produzca outcome formal consistente con A.3
  description: caso controlado con contrato compilado, sin bloqueo de policy y con output completo
  fixture_refs:
    - fixture://contracts/plan-basic-v1
    - fixture://policy/default-low-risk-v1
    - fixture://runtime/read-only-baseline-v1
  artifact_refs:
    contract_id: ctr-...
    contract_version: "1.0"
    approval_request_id: null
    approval_decision_id: null
    execution_id: exe-...
    result_id: res-...
  preconditions:
    - contract compilado y executable
    - clasificación entregable
    - evidence mínima disponible
  expected_outcome:
    status: pass
    result_outcome_level: success
    reason_codes_allowed:
      - success.*
    required_assertions:
      - plan_steps >= 1
      - approval_mode_required presente
      - evidence_refs no vacía
  oracle_type: deterministic
  oracle_ref: rules://a3/plan/success-minimum
  severity: high
  scoring:
    max_score: 100
    pass_threshold: 95
  gating_mode: required_for_merge
  tags:
    - phase-a
    - plan
    - regression
    - core
  owner: spec
  version: v1
```

Campos mínimos normativos:

- `id` (`eval_case_id`) — estable, único, semanticamente durable.
- `family`
- `title`
- `intent`
- `fixture_refs`
- `expected_outcome`
- `severity` — `critical` | `high` | `medium` | `low`.
- `oracle_type` — al menos `deterministic` o `rule_assisted`.
- `gating_mode` — `advisory` | `required_for_merge` | `required_for_release`.
- `tags`

Reglas sobre `eval_case_id`:

- Debe ser estable entre corridas equivalentes.
- Debe cambiar solo ante cambio material del caso, no por refresh cosmético.
- SHOULD codificar familia, dominio y variante (`eval.<family>.<domain>.<scenario>.vN`).
- No debe reutilizarse para dos oráculos distintos.

---

## Entradas mínimas de un eval

Todo eval debe poder resolver, de manera directa o por referencias, estas entradas mínimas:

- `eval_case_id`
- `family`
- `fixture_refs` o `artifact_refs` suficientes para reconstruir el caso
- `oracle_type`
- `expected_outcome`
- `severity`
- `gating_mode`
- contexto de ejecución del eval:
  - `tenant_id` si aplica
  - `environment`
  - `policy_version` o snapshot equivalente
  - `classification_snapshot` si aplica
  - `risk_snapshot` si aplica
  - `trace_id`/correlation root si se evalúa corrida real

Artifacts mínimos referenciables según familia:

- `contract_id` / `intent_contract_id`
- `contract_fingerprint`
- `approval_request_id`
- `approval_decision_id`
- `execution_id`
- `result_id`
- `event_ids` o `telemetry_refs`
- `evidence_refs`

Principio: si el eval no puede reconstruir qué caso se evaluó, bajo qué verdad esperada y contra qué artifacts, el eval está mal especificado.

---

## Salidas mínimas de un eval

Output lógico mínimo sugerido v1:

```yaml
eval_run:
  eval_run_id: evalrun-...
  eval_case_id: eval.result.plan.success.baseline.v1
  status: pass
  score: 98
  hard_failures: []
  soft_failures: []
  warnings: []
  observed_outcome:
    result_outcome_level: success
    reason_codes: [success.criteria_met]
  evidence_refs:
    - evr-...
    - event://execution.completed/...
  artifact_refs:
    contract_id: ctr-...
    execution_id: exe-...
    result_id: res-...
  executed_at: 2026-03-29T10:00:00Z
```

Campos mínimos obligatorios:

- `eval_run_id`
- `eval_case_id`
- `status` — `pass` | `hard_fail` | `soft_fail` | `warning`
- `score` — numérico normalizado, recomendado 0-100
- `hard_failures`
- `soft_failures`
- `warnings`
- `evidence_refs`
- `artifact_refs`
- `executed_at`

Semántica mínima de estado:

- `pass` — el caso cumple su oracle y supera el threshold.
- `hard_fail` — incumplimiento no tolerable para esa familia/severidad; invalida gating cuando aplica.
- `soft_fail` — incumplimiento relevante pero no necesariamente bloqueante según threshold/modo.
- `warning` — desvío informativo, deuda de calidad o señal de riesgo, sin invalidar el caso.

Norma: el framework debe poder exponer simultáneamente un `status` general y el detalle interno de `hard_failures`, `soft_failures` y `warnings`.

---

## Evals por tipo de resultado

Todo `result_eval` debe medir, como mínimo:

- consistencia del `result_type` con el contrato compilado,
- `result_outcome` estructurado según A.3,
- evidence mínima requerida por tipo,
- reason code correcto,
- clasificación y deliverability,
- completitud/calidad específica del tipo,
- consistencia entre output producido y restricciones del contrato.

### Matriz mínima obligatoria por tipo de resultado

| Tipo | Se evalúa sí o sí |
|---|---|
| `plan` | existencia de pasos ejecutables, tools/capabilities asignadas, asunciones documentadas, riesgo estimado, `approval_mode_required`, evidencia mínima |
| `inspection` | cobertura de entidades en scope, hallazgos alineados al objetivo, `data_accessed` documentado, `confidence_level`, `classification_snapshot`, outcome formal |
| `query` | `query_snapshot`, fuentes consultadas, `result_count`/resultado consistente, clasificación del output, tratamiento de redacción, evidencia mínima |
| `report` | cobertura del período, secciones requeridas, `completeness_level`, fuentes documentadas, secciones faltantes o redactadas justificadas, outcome formal |
| `change_proposal` | `diff_preview`, entidades afectadas, riesgo estimado, `approval_mode_required`, reversibilidad/rollback outline, consistencia con restricciones |
| `execution` | `plan_executed_snapshot` fingerprint, pasos completados/fallidos, outputs producidos, `approval_decision_ref`, evidencia mínima, consistencia con `criterios_de_exito` |
| `system_update` | todo lo de `execution` + sistemas modificados, cambios aplicados, `external_effect_confirmed`, verificabilidad post-aplicación, tratamiento de reversibilidad |
| `governance_decision` | `decision_outcome`, `policy_refs` versionadas, `effective_scope`, condiciones si `conditional`, `effective_until`, reason code y evidencia |

---

## Evals de policy

`policy_eval` verifica que el motor aplique correctamente reglas compiladas, floors, overrides y prohibiciones.

Debe cubrir, como mínimo:

- derivación de `approval_mode` efectiva,
- floors duros por clasificación, irreversible, `manage_policy`, `manage_connector`, cross-tenant y prod,
- resolución de riesgo y consistencia con policy snapshot,
- emisión correcta de `reason_code` de bloqueo o rechazo,
- bloqueo correcto de caminos prohibidos,
- consistencia entre policy compilada y artifacts emitidos.

Casos obligatorios de policy v1:

- `restricted` nunca por debajo de `pre_execution`.
- `external_effect = irreversible` fuerza `double`.
- `manage_policy` y `manage_connector` nunca por debajo de `pre_application`.
- `broad_scope + delegated` endurece el modo.
- `prod` endurece configure/publish/apply.
- change/global/cross-tenant nunca auto.

---

## Evals de approval

`approval_eval` valida el subsistema de approvals como comportamiento observable, no como tabla aislada de estados.

Debe cubrir:

- lifecycle válido y transiciones prohibidas,
- invalidación por cambio material,
- separación `execution_released` vs `application_released`,
- SoD y elegibilidad de aprobadores,
- expiración, revocación y supersession,
- SLA, reminders, timeout y fallback,
- coherencia entre `approval_request`, `approval_decision` y `execution_record`.

Checks mínimos obligatorios:

- sin approval vigente no hay ejecución cuando A.4 exige release,
- approval expirada/revocada/superseded bloquea correctamente,
- `pre_application` puede permitir ejecución técnica pero no aplicación,
- `double` respeta SoD y autoridad,
- todo cambio terminal deja evento y evidencia.

---

## Evals de clasificación y redacción

`classification_eval` cubre tres obligaciones distintas que NO deben mezclarse:

1. clasificación compilada correcta,
2. política de entrega correcta,
3. consistencia entre clasificación compilada y output realmente entregado.

Debe incluir obligatoriamente:

- redacción parcial cuando el tipo soporta entrega degradada,
- bloqueo total cuando no existe versión parcial entregable,
- consistency check entre `classification_snapshot`, `classification_level` compilado y delivered output,
- coherencia de `is_redacted`, `redaction_reason`, markers de redacción y `outcome_level`,
- validación de que el output entregado no filtre contenido por encima del clearance permitido,
- validación de que el resultado completo pueda quedar auditable aunque la versión entregada sea parcial.

Reglas mínimas:

- `classification_eval` debe fallar duro si la salida entregada contradice la clasificación compilada.
- Debe fallar duro si el sistema devuelve contenido no redactado donde correspondía bloqueo o redacción.
- Debe fallar blando o warning si la redacción es correcta pero la explicación/auditoría es incompleta.

---

## Evals de runtime general (consistencia entre contract/approval/result/execution)

`runtime_eval` valida la capa orquestadora superior definida en A.5.

Debe cubrir, como mínimo:

- que `execution_record` referencia un `intent_contract` válido y fingerprint correcto,
- que `approval_mode_effective` y releases sean consistentes con approvals vigentes,
- que `execution_state` no contradiga `contract_state`, `approval_state` ni `result_state`,
- que las transiciones del runtime sean válidas,
- que `result_record` apunte a la ejecución correcta,
- que el cierre terminal tenga trazabilidad completa.

Inconsistencias que deben evaluarse explícitamente:

- `execution_completed` sin resultado ni evidencia,
- `application_completed` sin release válido,
- `result_record` emitido para contrato superseded,
- fingerprint del plan ejecutado distinto al aprobado,
- `closed` sin eventos terminales correlacionables.

---

## Evals de resiliencia operacional (idempotencia/reintentos/compensación)

`resilience_eval` valida comportamiento defensivo y consistente del runtime bajo condiciones no ideales.

Debe cubrir:

- deduplicación correcta por `idempotency_key` y fingerprint material,
- separación entre retry técnico, retry de ejecución y replay de auditoría,
- tratamiento de `unknown_outcome`,
- bloqueo de auto-retry cuando existe riesgo de doble aplicación,
- uso correcto de verificación externa antes de reintentar,
- compensación lógica vs rollback físico,
- vínculo causal entre intentos y corridas.

Checks obligatorios:

- operaciones read-only equivalentes deben deduplicarse o devolver outcome existente cuando corresponda,
- mutaciones irreversibles no deben reintentarse ciegamente,
- `unknown_outcome` debe forzar verify-before-retry o escalación,
- compensación no debe borrar la falla original ni la traza causal,
- duplicados detectados deben quedar evidenciados; nunca descartados silenciosamente.

---

## Niveles de severidad y scoring de evals

Severidades obligatorias:

- `critical` — incumplimiento que rompe seguridad, governance dura, integridad de artifacts o riesgo de side effect indebido.
- `high` — incumplimiento que invalida outcome formal, release, policy floor o coherencia fuerte entre artifacts.
- `medium` — incumplimiento relevante de calidad, completitud, evidencia o explicabilidad, sin romper el control duro principal.
- `low` — deuda menor, warning de observabilidad o detalle accesorio.

Modelo recomendado de scoring:

- score base por caso: `100`.
- restar penalidades por hallazgo según severidad:
  - `critical`: -100
  - `high`: -40
  - `medium`: -15
  - `low`: -5
- el score mínimo es `0`.
- si existe al menos un hallazgo `critical`, el caso queda `hard_fail` independientemente del score.

Interpretación recomendada:

- `95-100`: pass fuerte
- `80-94`: pass con warning o deuda menor
- `60-79`: `soft_fail`
- `<60`: `hard_fail`

El score NO reemplaza el juicio semántico. Un `critical` puede dejar score `0` aunque el resto del caso esté correcto.

---

## Criterios de pass/fail por familia

### Regla transversal

- `hard_fail` cuando el motor viola una regla dura del dominio, entrega output inseguro, contradice snapshots materiales o rompe trazabilidad esencial.
- `soft_fail` cuando el comportamiento principal es reconocible pero incompleto, inconsistente en calidad o insuficiente en evidencia.
- `warning` cuando el comportamiento es aceptable pero deja deuda menor no bloqueante.

### Por familia

- `result_eval`: `hard_fail` si `result_outcome`, evidence mínima o criteria formales contradicen A.3.
- `policy_eval`: `hard_fail` si un floor duro no se aplica o se permite una acción prohibida.
- `approval_eval`: `hard_fail` si existe camino no auditado hacia ejecución/aplicación o si approvals inválidas siguen liberando.
- `classification_eval`: `hard_fail` si hay fuga de contenido o inconsistencia material entre clasificación y output entregado.
- `runtime_eval`: `hard_fail` si artifacts correlacionados se contradicen materialmente o el runtime entra en estado prohibido.
- `resilience_eval`: `hard_fail` si el sistema duplica side effects, reintenta ciegamente una mutación insegura o compensa sin preservar causalidad.

---

## Dataset strategy / fixtures mínimas recomendadas

Estrategia recomendada v1:

- mantener fixtures canónicas por familia,
- separar fixtures sintéticas de fixtures reconstruidas desde corridas reales anonimizadas,
- versionar snapshots materiales usados por el oracle,
- cubrir tanto happy paths como bordes críticos.

Fixtures mínimas recomendadas:

1. contrato read-only low-risk entregable.
2. contrato `report` con secciones redactables.
3. contrato `query` con bloqueo total por clasificación.
4. `change_proposal` reversible con diff preview correcto.
5. `execution` con fingerprint aprobado consistente.
6. `system_update` irreversible con approval `double`.
7. caso delegated + broad scope.
8. caso `manage_policy` en prod.
9. approval expirada antes de aplicar.
10. policy snapshot cambiada después de aprobar.
11. duplicate attempt dentro de misma ventana de dedup.
12. unknown outcome post-timeout con verificación externa pendiente.
13. compensación parcial.
14. corrida cerrada con trazabilidad completa.
15. corrida con evidencia incompleta pero outcome principal correcto.

Regla: cada fixture SHOULD declarar qué capas de captura activa (negocio, orquestación, ejecución, calidad, costo/rendimiento) para alinear el framework con el source-of-truth de observabilidad.

---

## Trazabilidad entre evals y artifacts del sistema

Todo eval debe poder referenciar explícitamente artifacts del sistema. Campos mínimos recomendados:

- `contract_id`
- `contract_version`
- `contract_fingerprint`
- `approval_request_id`
- `approval_decision_id`
- `execution_id`
- `result_id`
- `event_ids` / `telemetry_event_ids`
- `trace_id`
- `evidence_refs`

Reglas de trazabilidad:

- un eval sin `artifact_refs` suficientes solo puede existir como eval puramente sintética si su fixture es auto-contenida y versionada;
- si el eval se corre sobre artifacts reales, las referencias deben quedar persistidas en el output del eval;
- la evidencia del eval debe distinguir claramente entre evidence del sistema evaluado y evidence de la corrida de evaluación;
- el framework debe permitir explicar qué artifact falló, en qué campo y contra qué oracle.

---

## Eventos mínimos de ejecución de evals

Eventos mínimos recomendados v1:

- `eval.case_registered`
- `eval.run_started`
- `eval.fixture_resolved`
- `eval.oracle_resolved`
- `eval.assertion_executed`
- `eval.warning_emitted`
- `eval.soft_failure_detected`
- `eval.hard_failure_detected`
- `eval.score_computed`
- `eval.run_completed`

Payload mínimo sugerido:

- `event_id`
- `event_type`
- `eval_run_id`
- `eval_case_id`
- `family`
- `gating_mode`
- `severity`
- `tenant_id` si aplica
- `environment`
- `trace_id` si aplica
- `artifact_refs`
- `status`
- `score`
- `reason_codes`
- `occurred_at`

Correspondencia con capas de captura:

- Capa 1 negocio: caso, intent, tipo y outcome esperado.
- Capa 2 orquestación: suite, oracle, reglas cargadas, artifacts resueltos.
- Capa 3 ejecución: assertions, evidence y fallas detectadas.
- Capa 4 calidad: score, pass/fail por familia, cobertura.
- Capa 5 costo/rendimiento: duración del eval, costo estimado, tamaño del fixture, cantidad de artifacts consultados.

---

## Tests borde del framework de evals (mínimo 15)

1. Un `eval_case` sin `eval_case_id` estable debe ser rechazado.
2. Dos casos distintos no pueden compartir el mismo `eval_case_id` con oracle material distinto.
3. Un `result_eval` que solo devuelve bool sin `observed_outcome` estructurado debe fallar especificación.
4. Un `classification_eval` donde la salida entregada contiene contenido no redactado por encima del clearance permitido debe ser `hard_fail critical`.
5. Un `classification_eval` con redacción correcta pero sin `redaction_reason` debe ser al menos `soft_fail`.
6. Un `policy_eval` que permite `external_effect = irreversible` sin `double` debe ser `hard_fail`.
7. Un `approval_eval` que muestra `application_released` sin approval vigente debe ser `hard_fail`.
8. Un `approval_eval` donde cambio material no invalida aprobación previa debe ser `hard_fail`.
9. Un `runtime_eval` con `execution_completed` pero sin `result_id` ni evidencia mínima debe ser `hard_fail`.
10. Un `runtime_eval` con `plan_executed_snapshot` fingerprint distinto al aprobado debe ser `hard_fail`.
11. Un `resilience_eval` que reintenta automáticamente una mutación irreversible con `unknown_outcome` debe ser `hard_fail critical`.
12. Un `resilience_eval` que detecta duplicado y lo descarta sin evidencia debe ser `soft_fail` o `hard_fail` según severidad material.
13. Un eval advisory con hallazgos `medium` no debe bloquear release pero sí dejar score y warnings auditables.
14. Un eval `required_for_merge` con `high` y score bajo threshold debe bloquear merge aunque no sea `critical`.
15. Un eval `required_for_release` con cualquier `critical` debe bloquear release independientemente del score agregado.
16. Un eval sobre artifacts reales sin `trace_id` ni correlación suficiente debe fallar trazabilidad.
17. Un oracle `deterministic` no puede depender de juicio humano implícito no codificado.
18. Un oracle `rule_assisted` debe registrar qué reglas aplicó; si no, debe fallar auditabilidad.

---

## Criterios de aceptación del framework

- Existe una definición formal de eval alineada con A.3, A.4, A.5 y observabilidad source-of-truth.
- Existe `eval_case_id` estable y normativamente obligatorio.
- Existen al menos las familias `result_eval`, `policy_eval`, `approval_eval`, `classification_eval`, `runtime_eval` y `resilience_eval`.
- El framework soporta outcome estructurado, no solo bool.
- El framework define `severity`, `score`, `hard_fail`, `soft_fail` y `warning`.
- El framework soporta `oracle_type` al menos `deterministic` y `rule_assisted`.
- El framework soporta `gating_mode` al menos `advisory`, `required_for_merge` y `required_for_release`.
- Existe matriz mínima obligatoria por tipo de resultado para `plan`, `inspection`, `query`, `report`, `change_proposal`, `execution`, `system_update` y `governance_decision`.
- Los evals de clasificación cubren redacción parcial, bloqueo total y consistencia entre clasificación compilada y output entregado.
- Los evals dejan evidencia auditable y referencias a artifacts como `contract_id`, `execution_id`, `result_id` y relacionados.
- El framework define dataset strategy, fixtures mínimas, eventos de ejecución y tests borde.
- El framework queda apto para cerrar la primera mitad de A.6 como capa formal de evaluación del motor.
