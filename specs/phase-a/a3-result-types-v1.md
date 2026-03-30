# A.3 — Tipos de resultado v1

## Principios base

- El resultado no es un string ni un JSON libre; es un artefacto tipado, clasificado, auditable y gobernable.
- El tipo de resultado determina el input contract requerido, el output contract esperado, la evidencia mínima, el nivel de approval necesario y la telemetría que se captura.
- Un resultado puede ser completo, parcial o redactado. Si la clasificación no permite devolver el resultado completo, el sistema devuelve una versión degradada cuando esa versión parcial esté permitida.
- El éxito de un resultado es formal y medible, no subjetivo. Se define por criterios explícitos antes de ejecutar.
- Todo resultado deja evidencia auditable. No existe resultado sin traza.
- El tipo de resultado no lo infiere el motor libremente: sale del `tipo_de_resultado_esperado` del intent contract, compilado por el sistema según la capability invocada.

---

## Nota sobre la taxonomía

El source-of-truth lista 8 tipos nativos de resultado. Tras el análisis, se propone un ajuste:

**`aprobación/rechazo`** se reemplaza por **`governance_decision`**.

**Razón**: A.4 modela approvals como un subsistema completo con sus propios objetos, ciclo de vida y eventos. Un `approval_request` no es un resultado del motor — es un artefacto del subsistema de governance que habilita o bloquea otros resultados. Sin embargo, el motor sí puede emitir decisiones de governance como output explícito de una capability (evaluación de compliance, resolución de delegación, evaluación de policy eligibility). `governance_decision` captura ese concepto con precisión sin colisionar con el ciclo de approvals de A.4.

Los 7 tipos restantes se conservan sin modificación.

---

## Taxonomía de tipos de resultado

### 8 tipos canónicos

| ID | Nombre canónico | Familia |
|----|-----------------|---------|
| `plan` | Plan | read-only |
| `inspection` | Inspección/análisis | read-only |
| `query` | Consulta | read-only |
| `report` | Reporte/evidencia | read-only |
| `change_proposal` | Propuesta de cambio | mutation |
| `execution` | Ejecución | mutation |
| `system_update` | Actualización en sistemas | mutation |
| `governance_decision` | Decisión de governance | governance |

### Familias

#### Familia read-only
Resultados que no modifican estado en sistemas externos ni en la plataforma.
- Tipos: `plan`, `inspection`, `query`, `report`
- Efecto externo: `none`
- Approval mode mínimo: `auto` (salvo clasificación o riesgo eleven)

#### Familia mutation
Resultados que proponen o aplican cambios en sistemas externos o en la plataforma.
- Tipos: `change_proposal`, `execution`, `system_update`
- Efecto externo: `reversible` o `irreversible` según la acción
- Approval mode mínimo: `pre_application` para cambios reales tenant-scoped en producción

#### Familia governance
Resultados que emiten una decisión de governance sobre una entidad del sistema.
- Tipos: `governance_decision`
- Efecto externo: depende del tipo de decisión (`none` o `reversible`)
- Approval mode mínimo: `pre_execution` por defecto; `double` para decisiones que afecten policy, connector o scope global

---

## Objeto `result_type` canónico

Objeto de configuración que define el comportamiento de un tipo de resultado dentro del sistema. Vive en el catálogo; no se genera por request.

Campos obligatorios:

- `id` — identificador canónico del tipo
- `family` — `read_only` | `mutation` | `governance`
- `name` — nombre legible
- `description` — qué produce este tipo de resultado
- `produces_external_effect` — `none` | `reversible` | `irreversible`
- `min_approval_mode` — floor mínimo de approval para este tipo
- `requires_plan_snapshot` — bool; si la ejecución debe tener snapshot del plan aprobado
- `requires_diff_preview` — bool; aplica para mutation antes de aplicar
- `input_requirements` — lista de campos obligatorios que el intent contract debe tener
- `output_contract` — definición del contrato de salida
- `evidence_requirements` — evidencia mínima que el resultado debe dejar
- `classification_behavior` — cómo se propaga la clasificación al resultado
- `audit_behavior` — qué capas de telemetría se activan
- `allowed_outcome_levels` — lista de `result_outcome_level` válidos para este tipo
- `can_be_partial` — bool; si el tipo soporta resultado parcial/redactado
- `can_be_degraded` — bool; si el tipo soporta nivel `degraded`
- `supports_simulation` — bool; si el tipo puede ejecutarse en modo simulación
- `supports_replay` — bool; si el tipo puede reproducirse para auditoría
- `version` — versión del contrato
- `status` — `active` | `disabled` | `retired`

---

## Input contract por tipo

Campos del intent contract que deben estar presentes y resueltos antes de que el motor produzca ese tipo de resultado.

### `plan`
- `objetivo` — resuelto y no ambiguo
- `alcance` — definido
- `tipo_de_resultado_esperado` = `plan`
- `sistemas_posibles` — al menos uno identificado
- `autonomia_solicitada` — declarada
- Condición: no requiere `criterios_de_exito` formales previos, pero el motor los debe inferir o solicitar si la capability lo exige.

### `inspection`
- `objetivo` — qué inspeccionar, sobre qué entidad o sistema
- `alcance` — qué datos o contexto son relevantes
- `tipo_de_resultado_esperado` = `inspection`
- `datos_permitidos` — resuelto por el sistema según permisos y clasificación
- `sistemas_posibles` — al menos uno
- Condición: el acceso a datos restringidos requiere validación de permisos antes de inspeccionar. La salida puede ser parcial si la clasificación bloquea parte del contexto.

### `query`
- `objetivo` — pregunta o consulta específica
- `tipo_de_resultado_esperado` = `query`
- `datos_permitidos` — resuelto
- `alcance` — acotado
- Condición: más acotado que `inspection`. `query` devuelve datos puntuales; `inspection` produce análisis y hallazgos.

### `report`
- `objetivo` — qué evidencia o reporte producir
- `tipo_de_resultado_esperado` = `report`
- `alcance` — qué período, entidades o eventos cubrir
- `datos_permitidos` — resuelto
- `criterios_de_exito` — al menos un criterio de completitud
- Condición: la clasificación puede requerir redacción de secciones completas.

### `change_proposal`
- `objetivo` — qué cambio proponer
- `alcance` — qué sistemas y entidades afecta
- `tipo_de_resultado_esperado` = `change_proposal`
- `sistemas_posibles` — identificados
- `restricciones` — declaradas
- `datos_permitidos` — resuelto
- `aprobacion_requerida` — el sistema debe poder calcular el approval mode antes de proponer
- Condición: el motor debe generar un `diff_preview` como parte del output. El cambio propuesto no se aplica.

### `execution`
- `objetivo` — qué ejecutar
- `alcance` — exactamente definido, no ambiguo
- `tipo_de_resultado_esperado` = `execution`
- `sistemas_posibles` — confirmados, no solo candidatos
- `restricciones` — resueltas
- `datos_permitidos` — resuelto
- `autonomia_solicitada` — resuelta
- `aprobacion_requerida` — resuelta y consistente con el approval mode calculado
- `criterios_de_exito` — definidos antes de ejecutar
- `plan_snapshot` — plan aprobado que se va a ejecutar debe estar snapshotteado
- Condición: el `plan_snapshot` que se ejecuta debe coincidir fingerprint con el aprobado.

### `system_update`
Todos los campos de `execution`, más:
- `external_effect` — declarado explícitamente (`reversible` | `irreversible`)
- `rollback_plan` — debe existir si el efecto es reversible; si es `irreversible`, requiere confirmación explícita en la aprobación
- `destination_snapshot` — sistema destino snapshotteado al momento de la aprobación
- Condición: `irreversible` siempre fuerza `double` approval y risk mínimo `high`. Hereda floors duros de A.4. Diferencia con `execution`: `system_update` aplica cambios persistentes en sistemas externos; `execution` puede ser una ejecución sin persistencia externa.

### `governance_decision`
- `objetivo` — qué decisión de governance emitir
- `tipo_de_resultado_esperado` = `governance_decision`
- `alcance` — sobre qué entidad, tenant, capability o policy aplica la decisión
- `datos_permitidos` — resuelto
- `aprobacion_requerida` — siempre resuelta; nunca `auto` para decisiones sobre policy o connector
- Condición: no reemplaza el subsistema de approvals de A.4; lo complementa como resultado explícito de capabilities específicas.

---

## Output contract por tipo

### Campos comunes a todos los tipos

- `result_id` — identificador único del resultado
- `result_type` — tipo canónico
- `result_family` — familia
- `tenant_id`
- `environment`
- `execution_id` — ejecución que produjo este resultado
- `intent_contract_id` — contrato de intención que originó la ejecución
- `capability_id` — capability que produjo el resultado
- `outcome_level` — `success` | `partial_success` | `degraded` | `failed` | `blocked`
- `outcome_reason_code` — reason code normalizado (ver a3-success-criteria-v1)
- `outcome_reason` — descripción legible del outcome
- `classification_level` — clasificación del resultado
- `is_redacted` — bool; si el resultado fue redactado
- `redaction_reason` — qué se redactó y por qué
- `produced_at` — timestamp de producción
- `expires_at` — si el resultado tiene TTL
- `evidence` — evidencia mínima del resultado
- `audit_refs` — referencias a los eventos de telemetría asociados

### Campos adicionales: `plan`
- `plan_steps` — pasos propuestos (puede estar clasificado/redactado)
- `tools_proposed` — tools que se usarían para ejecutar
- `assumptions` — asunciones tomadas por el motor
- `open_questions` — preguntas que quedan abiertas si persiste ambigüedad
- `estimated_risk_snapshot` — snapshot del riesgo estimado al producir el plan
- `approval_mode_required` — qué approval mode requeriría ejecutar este plan

### Campos adicionales: `inspection`
- `findings` — hallazgos del análisis (puede estar parcialmente redactado)
- `entities_inspected` — entidades analizadas
- `data_accessed` — datos accedidos (referencia, no contenido bruto)
- `context_snapshot` — snapshot del contexto al momento de la inspección
- `anomalies_detected` — anomalías encontradas si las hay
- `confidence_level` — `high` | `medium` | `low`

### Campos adicionales: `query`
- `query_result` — resultado de la consulta (puede estar redactado)
- `data_sources` — fuentes consultadas
- `result_count` — cantidad de registros si aplica
- `query_snapshot` — query o parámetros de búsqueda normalizados

### Campos adicionales: `report`
- `report_sections` — secciones del reporte (cada una puede estar clasificada individualmente)
- `coverage_period` — período cubierto
- `data_sources` — fuentes usadas
- `completeness_level` — `full` | `partial` | `redacted`
- `missing_sections` — secciones que no se pudieron incluir y razón
- `summary` — resumen ejecutivo

### Campos adicionales: `change_proposal`
- `proposal_summary` — qué se propone en lenguaje claro
- `diff_preview` — preview del cambio propuesto
- `affected_entities` — entidades que cambiarían
- `estimated_risk_snapshot` — riesgo calculado para ejecutar esta propuesta
- `approval_mode_required` — qué approval mode requeriría aplicar este cambio
- `reversibility` — si el cambio sería reversible y cómo
- `rollback_plan_outline` — esquema del rollback si aplica
- `simulation_result` — resultado de simulación si se corrió previamente

### Campos adicionales: `execution`
- `execution_summary` — descripción de qué se ejecutó
- `tools_used` — tools efectivamente invocadas
- `plan_executed_snapshot` — plan que se ejecutó (fingerprint debe coincidir con el aprobado)
- `steps_completed` — pasos completados exitosamente
- `steps_failed` — pasos fallidos con reason code
- `outputs_produced` — outputs de cada paso relevante
- `approval_decision_ref` — referencia a la `approval_decision` que habilitó la ejecución

### Campos adicionales: `system_update`
Todos los campos de `execution`, más:
- `systems_modified` — sistemas que recibieron cambios
- `changes_applied` — detalle de cambios persistidos
- `rollback_available` — bool
- `rollback_procedure_ref` — referencia al procedimiento de rollback si existe
- `external_effect_confirmed` — confirmación del efecto externo real observado

### Campos adicionales: `governance_decision`
- `decision_type` — tipo de decisión emitida (compliance_evaluation, delegation_resolution, policy_eligibility, etc.)
- `decision_subject` — entidad sobre la que aplica la decisión
- `decision_outcome` — `approved` | `rejected` | `conditional` | `deferred`
- `decision_conditions` — condiciones si el outcome es `conditional`
- `policy_refs` — políticas evaluadas con su versión
- `effective_scope` — alcance efectivo de la decisión
- `effective_until` — vigencia de la decisión si aplica

---

## Evidencia mínima por tipo

| Tipo | Evidencia mínima obligatoria |
|------|------------------------------|
| `plan` | intent contract snapshot, assumptions, estimated risk snapshot, approval mode requerido para ejecutar |
| `inspection` | intent contract snapshot, entities inspected, data accessed refs, context snapshot, classification snapshot |
| `query` | intent contract snapshot, query snapshot, data sources, classification snapshot |
| `report` | intent contract snapshot, data sources, coverage period, completeness level, classification snapshot |
| `change_proposal` | intent contract snapshot, diff preview, estimated risk snapshot, approval mode requerido, simulation result si corrió |
| `execution` | intent contract snapshot, plan executed snapshot con fingerprint, approval decision ref, tools used, outputs produced, risk snapshot al momento de ejecución |
| `system_update` | todo lo de `execution` + systems modified, changes applied, rollback availability, external effect confirmed |
| `governance_decision` | intent contract snapshot, policy refs con versión, effective scope, decision conditions, decision subject snapshot |

---

## Relación resultado vs clasificación

- La clasificación del resultado se compila a partir de: clasificación del intent contract + clasificación de los datos accedidos + clasificación de los sistemas involucrados. Se toma el nivel más restrictivo.
- Si la clasificación del resultado excede el nivel permitido para el solicitante, el motor produce versión redactada cuando esté definida; de lo contrario, bloquea completamente.
- `is_redacted = true` implica que `outcome_level` puede bajar a `partial_success` o `degraded`.
- La clasificación del resultado queda en snapshot dentro del resultado y en el evento de cierre.
- Un resultado `governance_decision` sobre datos `restricted` nunca puede tener approval mode `auto` ni clasificación `public`.

---

## Relación resultado vs approvals

- El tipo de resultado determina el `min_approval_mode` como floor.
- El `min_approval_mode` del tipo se combina con el riesgo calculado (A.4) para derivar el approval mode efectivo: siempre se toma el más restrictivo.
- Familia `read_only`: floor `auto`, puede elevarse por clasificación o riesgo.
- Familia `mutation`: floor mínimo `pre_application` en producción con efecto externo real.
- `system_update` con `irreversible`: siempre `double` (hereda el override duro de A.4).
- Familia `governance`: floor mínimo `pre_execution`; `double` para decisiones sobre policy, connector, scope global o cross-tenant.
- Un resultado de familia `mutation` o `governance` no puede quedar en estado `applied` o `closed` sin `approval_decision_ref` válido cuando el approval mode no sea `auto`.

---

## Relación resultado vs auditoría

| Tipo | Capas de telemetría activadas |
|------|-------------------------------|
| `plan` | Capa 1 (negocio), Capa 2 (orquestación) |
| `inspection` | Capa 1, Capa 2, Capa 3 — énfasis en tool calls y datos accedidos |
| `query` | Capa 1, Capa 3 |
| `report` | Capa 1, Capa 2, Capa 3, Capa 4 (calidad — completeness) |
| `change_proposal` | Capa 1, Capa 2, Capa 3 — incluyendo diff y simulation |
| `execution` | Capas 1–4 completas |
| `system_update` | Capas 1–5 completas |
| `governance_decision` | Capas 1, 2, 3 — énfasis en policy refs y decision subject |

Capa 5 (costo/rendimiento) se activa siempre para todos los tipos.
