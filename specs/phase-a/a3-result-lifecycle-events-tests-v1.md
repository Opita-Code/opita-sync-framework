# A.3 — Lifecycle, eventos, tests v1

## Estados del ciclo de resultado

- `initializing` — el resultado está siendo preparado; el input contract está siendo validado
- `input_validated` — el input contract pasó validación; el motor puede proceder
- `governance_check` — se está evaluando si se puede producir este tipo de resultado (permisos, approval, policy)
- `governance_blocked` — governance bloqueó la producción; el resultado no puede avanzar
- `producing` — el motor está generando el resultado activamente
- `produced` — el resultado fue generado; pendiente de clasificación y control de salida
- `classifying` — se está aplicando la clasificación y control de salida al resultado
- `classified` — la clasificación fue aplicada; el resultado tiene su nivel de clasificación
- `redacting` — se está aplicando redacción porque la clasificación excede el permiso del solicitante
- `redacted` — la redacción fue aplicada; el resultado es parcial
- `delivering` — el resultado está siendo entregado al solicitante
- `delivered` — el resultado fue entregado exitosamente
- `failed` — el resultado falló en algún punto del ciclo
- `closed` — el resultado está cerrado (terminal); aplica a todos los outcomes

## Reglas operativas clave

- Todo resultado empieza en `initializing`.
- `governance_blocked` es un estado terminal válido: el sistema funcionó correctamente.
- `produced` no implica que el resultado sea entregable; debe pasar por `classifying` primero.
- `redacting` solo se activa si `classifying` determinó que el resultado completo no puede entregarse pero existe versión parcial permitida.
- Si no existe versión parcial permitida, el estado va directo a `governance_blocked` (clasificación bloqueante total).
- Todo estado terminal debe terminar en `closed`.
- `failed` siempre debe tener `outcome_reason_code` y evidencia mínima antes de ir a `closed`.

## Transiciones válidas

```
initializing -> input_validated
initializing -> failed (input inválido)
input_validated -> governance_check
governance_check -> producing (governance ok)
governance_check -> governance_blocked (governance denegó)
producing -> produced
producing -> failed (error de ejecución)
produced -> classifying
classifying -> classified (output completo permitido)
classifying -> redacting (output parcial disponible)
classifying -> governance_blocked (output completamente bloqueado por clasificación)
redacting -> redacted
classified -> delivering
redacted -> delivering
delivering -> delivered
delivering -> failed (error de entrega)
delivered -> closed
failed -> closed
governance_blocked -> closed
```

## Transiciones prohibidas

- `initializing -> delivered`
- `governance_check -> delivered`
- `produced -> delivered` (debe pasar por classifying)
- `produced -> closed` (sin clasificar)
- `classified -> closed` (sin entregar)
- `failed -> delivered`
- `governance_blocked -> producing`
- Reabrir cualquier estado terminal (`closed`, `delivered` final, `governance_blocked`)

---

## Eventos canónicos de resultado

### Eventos del ciclo principal

- `result.initializing_started`
- `result.input_validated`
- `result.input_validation_failed`
- `result.governance_check_started`
- `result.governance_check_passed`
- `result.governance_blocked`
- `result.producing_started`
- `result.produced`
- `result.production_failed`
- `result.classification_started`
- `result.classified`
- `result.redaction_started`
- `result.redacted`
- `result.delivering_started`
- `result.delivered`
- `result.delivery_failed`
- `result.closed`

### Eventos por causa específica

- `result.approval_required_check`
- `result.approval_missing_blocked`
- `result.approval_valid_confirmed`
- `result.plan_fingerprint_checked`
- `result.plan_fingerprint_mismatch`
- `result.classification_level_compiled`
- `result.redaction_applied`
- `result.partial_output_generated`
- `result.evidence_captured`
- `result.criteria_evaluated`
- `result.outcome_level_set`
- `result.rollback_triggered`
- `result.external_effect_confirmed`
- `result.external_effect_unconfirmed`

### Payload mínimo por evento

- `event_id`
- `event_type`
- `tenant_id`
- `environment`
- `result_id`
- `result_type`
- `result_family`
- `execution_id`
- `intent_contract_id`
- `capability_id`
- `from_state`
- `to_state`
- `outcome_level` (cuando ya está determinado)
- `outcome_reason_code` (cuando aplica)
- `classification_level` (cuando ya está determinado)
- `is_redacted` (cuando aplica)
- `triggered_by_subject_id`
- `occurred_at`

---

## Resultados parciales y redactados

### Cuándo se produce redacción

La redacción se activa cuando:
1. El resultado fue generado exitosamente (`produced`).
2. La clasificación compilada del resultado excede el nivel de acceso del solicitante.
3. Existe una versión parcial del resultado que sí puede ser entregada.

Si no existe versión parcial entregable → `governance_blocked` (no `redacted`).

### Qué se puede redactar por tipo

| Tipo | Qué se puede redactar |
|------|----------------------|
| `plan` | pasos con datos restringidos, tools con clasificación alta |
| `inspection` | findings con datos restringidos, secciones de entidades no accesibles |
| `query` | registros individuales, campos restringidos dentro de registros |
| `report` | secciones completas, campos específicos dentro de secciones |
| `change_proposal` | diff con datos restringidos, estimaciones de impacto si son sensibles |
| `execution` | outputs con datos restringidos, pasos que accedieron a datos clasificados |
| `system_update` | changes_applied si contienen datos restringidos, sistemas destino si son sensibles |
| `governance_decision` | decision_conditions si son clasificadas, policy_refs internas |

### Reglas de redacción

- Cada elemento redactado debe reemplazarse con un marcador de redacción que indique que fue omitido y por qué clase de restricción (no el contenido).
- `is_redacted = true` siempre que al menos un campo fue redactado.
- `redaction_reason` debe indicar la clasificación que aplicó, no el contenido redactado.
- Un resultado con redacción puede tener `outcome_level = success` si el objetivo principal se cumplió con la versión parcial entregada. Si no, `partial_success` o `degraded`.

---

## Tests borde mínimos

### T1 — Plan con objetivo ambiguo
Dado un intent contract con `objetivo` no resuelto.
Esperado: `initializing -> failed` con `failed.input.objective_ambiguous`.

### T2 — Execution sin plan_snapshot
Dado un execution request sin `plan_snapshot` en el contrato.
Esperado: `initializing -> failed` con `failed.input.missing_required_field`.

### T3 — Execution con fingerprint mismatch
Dado un execution request donde el plan a ejecutar difiere del aprobado.
Esperado: `governance_check -> governance_blocked` con `blocked.governance.approval_superseded` (o `failed.input.plan_fingerprint_mismatch`).

### T4 — system_update irreversible sin double approval
Dado un `system_update` con `external_effect = irreversible` y approval mode `pre_application` (no `double`).
Esperado: `governance_check -> governance_blocked` con `blocked.governance.no_approval` (floor duro de A.4 no satisfecho).

### T5 — governance_decision con policy_refs inaccesibles
Dado un `governance_decision` donde las policies requeridas no existen en el tenant.
Esperado: `producing -> failed` con `failed.execution.connector_unavailable` o similar; `outcome_level = failed`.

### T6 — Resultado con clasificación que excede permiso del solicitante (con versión parcial)
Dado un `inspection` cuyos findings incluyen datos `restricted` y el solicitante tiene acceso `confidential`.
Esperado: `classifying -> redacting -> redacted -> delivering -> delivered`; `is_redacted = true`; `outcome_level = partial_success` o `degraded`.

### T7 — Resultado con clasificación completamente bloqueante (sin versión parcial)
Dado un `query` donde todo el resultado es `restricted` y el solicitante no tiene acceso.
Esperado: `classifying -> governance_blocked -> closed`; `blocked.governance.classification_restricted`.

### T8 — system_update con partial_success irreversible
Dado un `system_update` con `external_effect = irreversible` que se aplicó en 2 de 3 sistemas.
Esperado: `outcome_level = failed` (no `partial_success`); `failed.execution.partial_apply_irreversible`. Regla dura.

### T9 — Approval expirada entre producción y entrega
Dado un resultado que fue aprobado, luego la aprobación expiró antes de que se completara la entrega.
Esperado: si el resultado ya fue `produced` y `classified`, la entrega puede proceder (la approval ya habilitó la ejecución). Si el approval se necesitaba para `delivering` (caso application_released de A.4), `governance_blocked`.

### T10 — Report con completeness_level partial pero criterio exige full
Dado un `report` cuyo intent contract define `criterios_de_exito` con completeness requerido = full, y la clasificación obligó a redactar secciones.
Esperado: `outcome_level = failed` o `degraded` según si el criterio es hard; `outcome_reason_code = failed.quality.completeness_insufficient` o `degraded.redaction_applied`.

### T11 — Execution donde un paso opcional falla
Dado un `execution` con 5 pasos donde 1 es opcional y falla.
Esperado: `outcome_level = partial_success`; `steps_failed` contiene el paso opcional con su reason code; objetivo principal cumplido.

### T12 — Execution donde un paso crítico falla
Dado un `execution` con 5 pasos donde 1 es crítico y falla.
Esperado: `outcome_level = failed`; `failed.execution.tool_error` o similar; el motor no continúa pasos dependientes del paso fallido.

### T13 — change_proposal con restricciones que contradicen el cambio
Dado un `change_proposal` donde el cambio propuesto viola una restricción dura del tenant.
Esperado: `producing -> failed`; `failed.input.contradictory_constraints`; el diff_preview puede mostrar el conflicto pero no puede proponer el cambio.

### T14 — governance_decision sobre scope global con approval mode auto
Dado un `governance_decision` con `alcance = global` y `aprobacion_requerida` resuelto como `auto`.
Esperado: el motor debe elevar el approval mode a al menos `pre_execution` (override de floor). Si el contrato fue compilado con `auto` para scope global, `governance_check -> governance_blocked` con `blocked.governance.policy_violation`.

### T15 — inspection con datos restringidos y sin permisos
Dado un `inspection` donde todos los datos en scope son `restricted` y el actor no tiene acceso.
Esperado: `governance_check -> governance_blocked`; `blocked.governance.insufficient_permissions`. No llega a `producing`.

### T16 — Resultado que intenta ir de produced a closed sin clasificar
Dado un resultado en estado `produced` que el sistema intenta cerrar sin pasar por `classifying`.
Esperado: transición rechazada; error de ciclo de vida; el resultado permanece en `produced`.

### T17 — system_update con external_effect_confirmed = false post-aplicación
Dado un `system_update` que completó la ejecución pero no pudo verificar el efecto externo.
Esperado: `outcome_level = failed`; `failed.execution.external_effect_unconfirmed`; el resultado queda en estado de alerta para revisión manual.

### T18 — plan donde todos los pasos requieren datos sin acceso
Dado un `plan` donde todos los pasos propuestos dependen de datos `restricted` inaccesibles para el actor.
Esperado: `outcome_level = degraded` con `degraded.assumptions_unvalidated`; el plan puede generarse con asunciones marcadas como no validables, pero no puede proceder a `execution` sin resolver el acceso.

### T19 — Reintento de un resultado en estado blocked
Dado un resultado en `governance_blocked` que el motor intenta reactivar sin cambios.
Esperado: el resultado en `governance_blocked` no puede reactivarse; debe crearse un nuevo result_id con el contrato corregido.

### T20 — report donde data_sources no están disponibles
Dado un `report` donde las fuentes de datos declaradas no están disponibles al momento de producir.
Esperado: `producing -> failed`; `failed.execution.connector_unavailable`; si alguna fuente secundaria sí está disponible y el contrato lo permite, `partial_success` con `missing_sections` documentadas.

---

## Criterios de aceptación de A.3

- Todo intent contract puede derivar el tipo de resultado sin ambigüedad a partir de `tipo_de_resultado_esperado` y la capability invocada.
- Todo tipo de resultado tiene input contract definido; el motor puede validar que el contrato está listo antes de producir.
- Todo resultado producido tiene `outcome_level`, `outcome_reason_code` y evidencia mínima.
- Todo resultado pasa por `classifying` antes de ser entregado; no existe entrega sin clasificación aplicada.
- Todo resultado bloqueado por governance emite evento canónico y queda auditado.
- Todo `system_update` con `irreversible` usa `double` approval y no soporta `partial_success`.
- Todo resultado redactado tiene `is_redacted = true` y `redaction_reason` documentado.
- Todo fallo tiene reason code de la taxonomía normalizada; no existe fallo sin reason code.
- El motor separa correctamente `blocked` (governance funcionó) de `failed` (error de ejecución).
- Todo estado terminal termina en `closed` con evidencia.
- No existe resultado en `applied` o equivalente sin `approval_decision_ref` cuando el approval mode no es `auto`.
