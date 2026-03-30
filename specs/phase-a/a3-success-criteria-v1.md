# A.3 — Criterios formales de éxito v1

## Principios base

- El éxito no es un bool. Es un objeto con nivel, evidencia y razón.
- Los criterios de éxito se definen antes de ejecutar, no después. Si no se pueden definir, el contrato no está listo para ejecutar.
- Un resultado parcial o degradado es un outcome legítimo si está documentado, clasificado y auditado. No es un fallo silencioso.
- El fallo siempre tiene razón. No existe fallo sin reason code y evidencia de lo que se intentó.
- Los criterios de éxito son parte del intent contract (`criterios_de_exito`) y del `result_type` (criterios formales por tipo). Ambos deben ser consistentes.

---

## Objeto `result_outcome`

Representa el resultado de éxito o fallo de una ejecución. Es parte del output de todo resultado.

Campos:

- `outcome_level` — `success` | `partial_success` | `degraded` | `failed` | `blocked`
- `outcome_reason_code` — reason code normalizado (ver sección de reason codes)
- `outcome_reason` — descripción legible
- `criteria_evaluated` — lista de criterios evaluados con su resultado individual
- `criteria_met` — cuántos criterios se cumplieron
- `criteria_total` — cuántos criterios se evaluaron
- `evidence_refs` — referencias a la evidencia que respalda el outcome
- `partial_output_available` — bool; si hay versión parcial disponible
- `degraded_reason` — razón por la que el resultado quedó degradado si aplica
- `blocked_by` — qué bloqueó el resultado si `outcome_level = blocked`
- `recoverable` — bool; si el fallo es recuperable con reintento o corrección
- `recovery_hint` — sugerencia de recuperación si aplica

---

## Niveles de outcome

### `success`
- Todos los criterios formales cumplidos.
- Output completo, no redactado.
- Evidencia mínima presente.
- Clasificación no bloqueó nada.

### `partial_success`
- Al menos un criterio formal cumplido pero no todos.
- Output incompleto por razones documentadas (clasificación, scope parcial, datos insuficientes).
- La parte producida es válida y auditable.
- Aceptable cuando el intent contract explícitamente permite resultado parcial.

### `degraded`
- Output producido pero con calidad reducida respecto a lo esperado.
- Ejemplo: reporte con secciones redactadas, inspección con baja confianza, plan con asunciones sin validar.
- Requiere documentar la razón del degradado y qué se sacrificó.
- No es un fallo pero no cumple el estándar de `success`.

### `failed`
- No se pudo producir el resultado esperado.
- Razón documentada con reason code.
- Puede ser recuperable o no.

### `blocked`
- El resultado no se ejecutó porque governance lo impidió.
- Motivos: sin aprobación, clasificación bloqueante, policy violation, SoD violation.
- No es un error de ejecución: es el sistema funcionando correctamente.
- Siempre deja evento de bloqueo y evidencia.

---

## Criterios formales de éxito por tipo

### `plan`
Criterios mínimos para `success`:
1. El plan tiene al menos un paso ejecutable definido.
2. Cada paso tiene tool o capability asignada.
3. Las asunciones están documentadas.
4. El riesgo estimado fue calculado.
5. El approval mode requerido para ejecutar fue derivado.

Degradación aceptable (`degraded`):
- Plan con pasos sin tool asignada si se documenta como limitación de contexto.
- Plan con asunciones sin validar si se marcan explícitamente.

Fallo (`failed`):
- No se pudo producir ningún paso ejecutable.
- El objetivo es contradictorio con las restricciones.
- Contexto insuficiente y el motor no pudo asumir ni proponer.

### `inspection`
Criterios mínimos para `success`:
1. Se inspeccionaron todas las entidades en scope.
2. Los hallazgos cubren el objetivo declarado.
3. Se documentaron los datos accedidos.
4. La confianza del análisis es `medium` o `high`.
5. La clasificación de los datos accedidos quedó en snapshot.

Degradación aceptable (`degraded`):
- Inspección con confianza `low` si se documenta la razón.
- Hallazgos parciales por acceso limitado a datos restringidos.

Fallo (`failed`):
- No se pudo acceder a ninguna entidad en scope.
- Permisos insuficientes para todo el alcance definido.
- Datos insuficientes para producir hallazgos mínimos.

### `query`
Criterios mínimos para `success`:
1. La consulta devolvió resultados de al menos una fuente.
2. El query_snapshot está registrado.
3. La clasificación de los resultados está documentada.

Degradación aceptable (`degraded`):
- Resultados redactados por clasificación si se documenta qué se redactó.
- Resultados de fuentes parciales si el resto no estaba accesible.

Fallo (`failed`):
- La consulta no devolvió resultados y el scope era válido.
- Acceso bloqueado a todas las fuentes del alcance.
- Query inválido o no compilable.

### `report`
Criterios mínimos para `success`:
1. Todas las secciones del alcance fueron producidas.
2. El período de cobertura fue completamente cubierto.
3. Las fuentes de datos están documentadas.
4. El `completeness_level` es `full`.
5. Los criterios de completitud del intent contract se cumplieron.

Degradación aceptable (`degraded`):
- `completeness_level = partial` si las secciones faltantes están documentadas.
- Secciones redactadas por clasificación si están identificadas.

Fallo (`failed`):
- No se pudo producir ninguna sección.
- El período de cobertura no pudo ser cubierto en absoluto.
- Datos insuficientes para cumplir el criterio mínimo de completitud.

### `change_proposal`
Criterios mínimos para `success`:
1. El `diff_preview` fue generado y es legible.
2. Las entidades afectadas están identificadas.
3. El riesgo estimado fue calculado.
4. El approval mode requerido para aplicar fue derivado.
5. La reversibilidad del cambio está documentada.

Degradación aceptable (`degraded`):
- Diff con secciones estimadas (no exactas) si se documentan las limitaciones.
- Riesgo estimado con menor precisión por contexto insuficiente.

Fallo (`failed`):
- No se pudo generar el diff preview.
- No se pudo identificar ninguna entidad afectada.
- El cambio propuesto contradice restricciones duras del tenant.

### `execution`
Criterios mínimos para `success`:
1. Todos los pasos del plan ejecutado se completaron.
2. El fingerprint del plan ejecutado coincide con el aprobado.
3. Los outputs producidos cumplen los criterios_de_exito del intent contract.
4. La evidencia mínima está presente.
5. No hubo pasos fallidos.

Degradación aceptable (`degraded`):
- Algunos pasos opcionales fallaron pero el objetivo principal se cumplió.
- Outputs con menor completitud si el objetivo central está cubierto.

Fallo (`failed`):
- Al menos un paso crítico falló.
- El fingerprint del plan ejecutado no coincide con el aprobado.
- Los outputs no cumplen ninguno de los criterios_de_exito.

### `system_update`
Criterios mínimos para `success`:
1. Todos los criterios de `execution` cumplidos.
2. Los cambios fueron aplicados en todos los sistemas en scope.
3. El efecto externo real coincide con el declarado en el intent contract.
4. `external_effect_confirmed = true`.
5. El estado post-aplicación es verificable.

Degradación aceptable (`degraded`):
- Aplicación parcial si los sistemas restantes no estaban accesibles y se documenta.
- No aplica para `irreversible`: o se aplica completo o es fallo.

Fallo (`failed`):
- No se aplicó ningún cambio en sistemas externos.
- El efecto externo real difiere del declarado.
- El estado post-aplicación no es verificable.

### `governance_decision`
Criterios mínimos para `success`:
1. La decisión fue emitida con `decision_outcome` definido.
2. Las políticas evaluadas están referenciadas con versión.
3. El alcance efectivo está documentado.
4. Si el outcome es `conditional`, las condiciones están explicitadas.
5. La vigencia de la decisión (`effective_until`) está definida.

Degradación aceptable (`degraded`):
- Decisión `deferred` con razón documentada.
- Decisión `conditional` con condiciones parcialmente definidas si se documenta la limitación.

Fallo (`failed`):
- No se pudo emitir ninguna decisión.
- Las políticas requeridas no existen o no son accesibles.
- El alcance es inválido o no gobernable.

---

## Failure taxonomy

### Familia 1 — Fallos de input (el contrato no estaba listo)

| Reason code | Descripción |
|-------------|-------------|
| `failed.input.objective_ambiguous` | El objetivo no se pudo resolver sin ambigüedad |
| `failed.input.scope_undefined` | El alcance no está definido o es inconsistente |
| `failed.input.missing_required_field` | Falta un campo obligatorio del input contract |
| `failed.input.contradictory_constraints` | Las restricciones del contrato se contradicen entre sí |
| `failed.input.success_criteria_missing` | No se definieron criterios de éxito para tipos que los requieren |
| `failed.input.plan_fingerprint_mismatch` | El plan a ejecutar no coincide fingerprint con el aprobado |

### Familia 2 — Fallos de governance (el sistema bloqueó correctamente)

| Reason code | Descripción |
|-------------|-------------|
| `blocked.governance.no_approval` | No existe approval_decision válida para el modo requerido |
| `blocked.governance.approval_expired` | La aprobación existía pero expiró antes de ejecutar |
| `blocked.governance.approval_revoked` | La aprobación fue revocada antes de aplicar |
| `blocked.governance.approval_superseded` | La aprobación fue invalidada por cambio material |
| `blocked.governance.classification_restricted` | La clasificación del resultado bloquea la entrega completa |
| `blocked.governance.policy_violation` | La acción viola una policy activa del tenant |
| `blocked.governance.sod_violation` | La ejecución viola separación de responsabilidades |
| `blocked.governance.insufficient_permissions` | El actor no tiene permisos para producir este tipo de resultado |
| `blocked.governance.scope_exceeds_delegation` | El alcance supera los límites de la delegación activa |

### Familia 3 — Fallos de ejecución (el motor lo intentó y falló)

| Reason code | Descripción |
|-------------|-------------|
| `failed.execution.tool_error` | Una tool invocada retornó error irrecuperable |
| `failed.execution.tool_timeout` | Una tool invocada no respondió dentro del timeout |
| `failed.execution.connector_unavailable` | El conector al sistema externo no está disponible |
| `failed.execution.external_system_error` | El sistema externo retornó error |
| `failed.execution.partial_apply_irreversible` | Aplicación parcial de un cambio marcado como irreversible |
| `failed.execution.rollback_triggered` | Se ejecutó rollback por fallo en medio de la ejecución |
| `failed.execution.external_effect_unconfirmed` | No se pudo verificar el efecto externo post-aplicación |
| `failed.execution.max_retries_exceeded` | Se superó el límite de reintentos |

### Familia 4 — Fallos de calidad (el resultado no cumple el estándar)

| Reason code | Descripción |
|-------------|-------------|
| `failed.quality.criteria_not_met` | Los outputs producidos no cumplen los criterios_de_exito |
| `failed.quality.confidence_too_low` | La confianza del análisis es insuficiente para el tipo de resultado |
| `failed.quality.completeness_insufficient` | El nivel de completitud no alcanza el mínimo requerido |
| `failed.quality.evidence_incomplete` | La evidencia mínima requerida no se pudo producir |
| `failed.quality.output_undeliverable` | El output fue producido pero no puede ser entregado al solicitante |

### Reason codes de éxito parcial y degradado

| Reason code | Descripción |
|-------------|-------------|
| `partial_success.scope_limited` | Se cumplió el objetivo principal pero no todo el alcance |
| `partial_success.data_access_restricted` | Parte de los datos no eran accesibles por clasificación |
| `partial_success.optional_steps_failed` | Los pasos opcionales fallaron pero el objetivo central se cumplió |
| `degraded.redaction_applied` | El output fue producido pero con redacciones por clasificación |
| `degraded.low_confidence` | El análisis fue producido pero con confianza reducida |
| `degraded.assumptions_unvalidated` | El plan incluye asunciones que no pudieron validarse |
| `degraded.partial_coverage` | El reporte o inspección cubre menos del alcance esperado |

### Reason codes de éxito completo

| Reason code | Descripción |
|-------------|-------------|
| `success.all_criteria_met` | Todos los criterios formales se cumplieron |
| `success.within_policy` | El resultado fue producido dentro de los límites de policy |
| `success.with_approved_exception` | El resultado fue producido con una excepción aprobada explícita |

---

## Métricas por tipo de resultado

Qué se mide para evaluar calidad (conecta con Capa 4 de telemetría).

| Tipo | Métricas clave |
|------|----------------|
| `plan` | tasa de planes ejecutables sin corrección, tasa de asunciones validadas, precisión del riesgo estimado vs real |
| `inspection` | cobertura del scope, confianza promedio, tasa de anomalías detectadas vs confirmadas |
| `query` | tasa de resultados completos vs redactados, latencia, tasa de cache hit si aplica |
| `report` | completeness_level promedio, tasa de secciones redactadas, cobertura temporal |
| `change_proposal` | tasa de propuestas aceptadas, precisión del diff_preview vs cambio real aplicado |
| `execution` | tasa de éxito de pasos, tasa de fingerprint match, tasa de reintentos |
| `system_update` | tasa de aplicación completa, tasa de rollback, tasa de external_effect_confirmed |
| `governance_decision` | tasa de decisiones emitidas vs bloqueadas, tasa de condicionales vs definitivas |

---

## Matriz de outcome permitido por tipo y familia

| Tipo | success | partial_success | degraded | failed | blocked |
|------|---------|-----------------|----------|--------|---------|
| `plan` | sí | sí | sí | sí | sí |
| `inspection` | sí | sí | sí | sí | sí |
| `query` | sí | sí | sí | sí | sí |
| `report` | sí | sí | sí | sí | sí |
| `change_proposal` | sí | sí | sí | sí | sí |
| `execution` | sí | sí (solo pasos opcionales) | no | sí | sí |
| `system_update` | sí | solo en reversible | no | sí | sí |
| `governance_decision` | sí | sí (deferred/conditional) | no | sí | sí |

`system_update` irreversible no soporta `partial_success`: o se aplica completo o es `failed`. Esta regla es dura.
