# A.2 — Versionado, eventos y tests del contrato v1

## Principios de versionado

- El contrato es versionado por cambio material, no por cada edición cosmética.
- Todo contrato tiene una versión inicial `v1` al momento de crearse en `draft`.
- Todo cambio material incrementa la versión del contrato y recalcula el `fingerprint`.
- Los cambios no materiales pueden registrarse como revisiones menores, pero no invalidan aprobación ni generan nueva versión material.
- Un contrato `superseded` no se reutiliza: si el usuario quiere continuar, se crea un nuevo contrato derivado con `parent_contract_id` apuntando al anterior.

---

## Versionado del contrato

### `contract_version`

Formato recomendado v1:
- `major.minor`
- `major` incrementa con cada cambio material
- `minor` incrementa con cada cambio no material documentado

Ejemplos:
- `1.0` — contrato inicial en `draft`
- `1.1` — se agregaron `notas_del_usuario` sin impacto material
- `2.0` — cambió `alcance` y eso obligó a recompilar fingerprint

### Cuándo incrementa `major`

Cambios materiales que incrementan `major`:
- `objetivo`
- `alcance`
- `tipo_de_resultado_esperado`
- `sistemas_confirmados`
- `datos_permitidos`
- `herramientas_permitidas`
- `capability_id`
- `classification_level`
- `nivel_de_riesgo`
- `approval_mode_efectivo`
- `policy_snapshot` (versión)
- `plan_snapshot`
- `destination_snapshot`
- `criterios_de_exito` para `execution` y `system_update`

### Cuándo incrementa `minor`

Cambios no materiales que incrementan `minor`:
- `notas_del_usuario`
- `notas_de_contexto`
- `contexto_relevante`
- correcciones de formato o normalización textual que no alteren semántica
- metadata operativa no usada en fingerprint

### Política de snapshots por versión

Cada versión `major` del contrato debe capturar y conservar:
- `fingerprint`
- `policy_snapshot`
- `classification_snapshot`
- `risk_snapshot`
- `permission_snapshot`
- `delegation_snapshot` si aplica
- `context_snapshot`
- `plan_snapshot` si existe en esa versión
- `destination_snapshot` si aplica

---

## Diffs del contrato

### Tipos de diff

#### `semantic_diff`
Compara los campos materiales y determina si cambió el significado operativo del contrato.

Campos incluidos:
- `objetivo`
- `alcance`
- `tipo_de_resultado_esperado`
- `sistemas_confirmados`
- `datos_permitidos`
- `herramientas_permitidas`
- `capability_id`
- `classification_level`
- `nivel_de_riesgo`
- `approval_mode_efectivo`
- `plan_snapshot`
- `destination_snapshot`
- `criterios_de_exito` (si material para el tipo)

Salida recomendada:
- `changed_fields`
- `material_change_detected` — bool
- `old_fingerprint`
- `new_fingerprint`
- `requires_reapproval` — bool
- `reason_codes` — lista de razones del diff

#### `presentation_diff`
Compara campos no materiales y cambios cosméticos.

Campos incluidos:
- `notas_del_usuario`
- `notas_de_contexto`
- `contexto_relevante`
- orden o formateo de listas sin cambio semántico

Salida recomendada:
- `changed_fields`
- `material_change_detected = false`
- `requires_reapproval = false`

### Reason codes sugeridos para diff

- `diff.material.objective_changed`
- `diff.material.scope_changed`
- `diff.material.result_type_changed`
- `diff.material.systems_changed`
- `diff.material.data_scope_changed`
- `diff.material.tools_changed`
- `diff.material.capability_changed`
- `diff.material.classification_changed`
- `diff.material.risk_changed`
- `diff.material.approval_mode_changed`
- `diff.material.policy_version_changed`
- `diff.material.plan_changed`
- `diff.material.destination_changed`
- `diff.material.success_criteria_changed`
- `diff.non_material.notes_changed`
- `diff.non_material.context_notes_changed`
- `diff.non_material.memory_context_changed`

---

## Eventos canónicos del contrato

### Eventos del ciclo principal

- `contract.created`
- `contract.inspection_started`
- `contract.inspection_completed`
- `contract.incomplete_detected`
- `contract.compilation_started`
- `contract.compilation_completed`
- `contract.compilation_failed`
- `contract.plan_generation_started`
- `contract.plan_proposed`
- `contract.plan_approved`
- `contract.executable_ready`
- `contract.execution_started`
- `contract.execution_completed`
- `contract.superseded`
- `contract.cancelled`
- `contract.closed`

### Eventos por cambio o causa específica

- `contract.material_change_detected`
- `contract.non_material_change_recorded`
- `contract.fingerprint_recomputed`
- `contract.version_incremented`
- `contract.policy_snapshot_captured`
- `contract.classification_snapshot_captured`
- `contract.risk_snapshot_captured`
- `contract.permission_snapshot_captured`
- `contract.delegation_snapshot_captured`
- `contract.context_snapshot_captured`
- `contract.approval_mode_derived`
- `contract.validation_failed`
- `contract.validation_passed`
- `contract.missing_user_input_detected`
- `contract.capability_selected`
- `contract.plan_fingerprint_bound`

### Payload mínimo sugerido por evento

- `event_id`
- `event_type`
- `tenant_id`
- `environment`
- `contract_id`
- `contract_version`
- `parent_contract_id`
- `fingerprint`
- `contract_state`
- `from_state`
- `to_state`
- `changed_fields` (si aplica)
- `material_change_detected` (si aplica)
- `approval_mode_efectivo` (si aplica)
- `classification_level` (si aplica)
- `nivel_de_riesgo` (si aplica)
- `triggered_by_subject_id`
- `occurred_at`

---

## Tests borde mínimos

### T1 — Contrato draft sin objetivo
Dado un contrato creado sin `objetivo`.
Esperado: `draft -> incomplete` o rechazo de recepción; `contract.validation_failed`; reason code: input objetivo faltante.

### T2 — Usuario cambia objetivo después de plan aprobado
Dado un contrato en `plan_approved` donde el usuario cambia `objetivo`.
Esperado: `material_change_detected = true`; incremento `major`; `fingerprint` cambia; contrato pasa a `superseded`; aprobación previa invalidada.

### T3 — Usuario cambia notas del usuario después de plan aprobado
Dado un contrato en `plan_approved` donde el usuario solo cambia `notas_del_usuario`.
Esperado: incremento `minor`; `material_change_detected = false`; aprobación previa sigue válida.

### T4 — Policy version cambia entre plan aprobado y executable
Dado un contrato en `plan_approved` y el tenant publica una nueva policy relevante.
Esperado: `policy_snapshot` cambia; `material_change_detected = true`; contrato pasa a `superseded`.

### T5 — Delegación expira durante compilación
Dado un contrato en `compiling` con `delegation_id` activa al inicio pero expirada antes de cerrar compilación.
Esperado: `delegation_snapshot` marca expiración; si el scope dependía de esa delegación, compilación falla o el contrato pasa a `incomplete`/`cancelled` según el caso.

### T6 — Classification sube después de compilar
Dado un contrato compilado con `classification_level = confidential`, pero durante ejecución se detecta dato `restricted`.
Esperado: `classification_snapshot` cambia; `material_change_detected = true`; contrato pasa a `superseded` antes de continuar.

### T7 — Riesgo insuficiente para calcular
Dado un contrato con variables de riesgo incompletas.
Esperado: el motor asume el peor caso (`critical`), registra nota contextual, y compila con `nivel_de_riesgo = critical`; no deja el riesgo sin resolver.

### T8 — approval_mode efectivo menor al floor del tipo
Dado un contrato `system_update` que el usuario fuerza con `aprobacion_requerida = auto`.
Esperado: la compilación eleva a al menos `pre_application` o `double` según el caso; nunca acepta `auto` por debajo del floor.

### T9 — Plan snapshot no coincide con fingerprint actual
Dado un contrato en `executable` cuyo `plan_snapshot` fue aprobado con fingerprint viejo.
Esperado: `contract.plan_fingerprint_bound` falla; contrato pasa a `superseded`; no ejecuta.

### T10 — Reapertura de contrato closed
Dado un contrato en `closed` que el motor intenta pasar a `draft`.
Esperado: transición rechazada. Un contrato `closed` no se reabre; se crea uno nuevo derivado.

### T11 — Cambio de capability sin impacto en objetivo pero sí en toolchain
Dado un contrato compilado donde cambia `capability_id` por otra capability equivalente.
Esperado: sigue siendo material; `major` incrementa; `requires_reapproval = true`.

### T12 — Cambio de contexto_relevante solamente
Dado un contrato compilado donde la memoria recuperada cambia pero no altera campos materiales.
Esperado: `minor` incrementa; `material_change_detected = false`; solo `contract.context_snapshot_captured`.

### T13 — Falta tool compatible con el tipo de resultado
Dado un contrato `execution` donde `herramientas_permitidas` queda vacío.
Esperado: `compilation_failed` o `incomplete`; nunca pasa a `compiled`.

### T14 — Contradicción entre restricciones del usuario y del sistema
Dado un contrato donde el usuario exige usar un sistema que policy bloquea.
Esperado: compilación falla o el contrato queda `incomplete` con la contradicción explícita; nunca pasa a `compiled` ocultando el conflicto.

### T15 — system_update irreversible sin rollback_plan
Dado un contrato `system_update` con `external_effect = irreversible` y sin `rollback_plan` ni nota explícita.
Esperado: no pasa a `executable`; validación falla.

### T16 — Contract superseded genera nuevo draft derivado
Dado un contrato `superseded` y el usuario decide continuar con los cambios nuevos.
Esperado: se crea un nuevo contrato `draft` con `parent_contract_id` apuntando al anterior y `contract_version = 1.0` en el nuevo objeto.

### T17 — approval válida pero contract_state incorrecto
Dado un contrato con `approval_decision` válida pero `contract_state = compiled`.
Esperado: no ejecuta; debe pasar por `planning`/`plan_approved`/`executable` primero.

### T18 — Cambio de criterios_de_exito en report
Dado un contrato `report` donde cambian los criterios de completitud antes de ejecutar.
Esperado: si afecta la definición operativa del reporte, `material_change_detected = true`; si es solo aclaración no semántica, `minor`.

### T19 — Tenant inactivo durante recepción
Dado un contrato recibido con `tenant_id` de un tenant deshabilitado.
Esperado: rechazo de recepción; no se crea contrato válido.

### T20 — Environment inválido
Dado un contrato con `environment = qa`.
Esperado: validación falla; solo se aceptan `dev`, `staging`, `prod`.

---

## Criterios de aceptación de A.2

- Todo contrato tiene schema exacto con separación clara entre campos del usuario, campos del sistema, campos técnicos y snapshots.
- Todo contrato puede clasificarse como válido, compilable o ejecutable según reglas explícitas.
- Todo contrato tiene validaciones duras por fase (recepción, compilación, ejecución).
- Todo contrato tiene estados formales y transiciones válidas/prohibidas.
- Todo cambio material recalcula fingerprint, incrementa versión material y puede invalidar approvals previas.
- Todo cambio no material puede registrarse sin romper approvals.
- Todo contrato `superseded` genera un camino claro para derivar un nuevo contrato.
- Todo evento canónico del contrato tiene payload mínimo definido.
- El contrato es consistente con A.3 (tipos de resultado) y A.4 (approvals y risk model).
- No existe ejecución desde un contrato que no esté en estado `executable`.
