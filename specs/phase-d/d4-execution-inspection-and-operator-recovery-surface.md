# D.4 — Execution inspection and operator recovery surface

## Objetivo

Definir la surface v1 desde la cual operator, reviewer y developer puedan **inspeccionar ejecuciones reales, correlacionar artifacts/IDs/evidencia y preparar recovery operacional permitido** sobre el kernel ya cerrado en C.2, C.4, C.6 y la preview de D.3, sin crear una segunda verdad de runtime ni habilitar mutación directa del estado canónico.

D.4 opera sobre evidencia y seams ya cerrados del kernel: `execution_record`, `policy_decision_record`, `approval_request`, `result_record`, `event_record`, runtime durable sobre Temporal y resolution `capability -> bundle -> binding -> provider`. No recompila contrato, no reinterpreta policy y no redefine lifecycle del motor.

## Principios de implementación de la operator surface

- La surface de inspección **lee, correlaciona y presenta**; NO redefine estado canónico.
- El kernel sigue siendo el único responsable de ejecutar transiciones durables y persistirlas.
- `execution_inspection_view`, `recovery_action_candidate` y la acción real de recovery ejecutada por el kernel son **tres boundaries distintos**.
- La operator surface puede preparar, validar y solicitar recovery; nunca ejecutarlo por fuera del workflow/runtime canónico.
- La evidencia canónica vive en PostgreSQL y records del kernel; Temporal history y OTel/LGTM son soporte complementario, no source of truth.
- Si la surface no puede sostener una lectura con correlación y evidencia mínima, debe degradar a `inspection_incomplete` o bloquear el recovery candidateado, nunca inventar certeza.
- Recovery v1 es estrictamente **operacional y gobernado**: corrige pausas, retries técnicos y escalaciones permitidas; no reescribe resultados.
- Toda lectura y todo recovery debe ser tenant-scoped, auditado, idempotente y correlable.
- Los estados problemáticos (`blocked`, `failed`, `unknown_outcome`, `compensation_pending`, `compensated`) deben exponerse explícitamente; no pueden esconderse en un genérico “error”.

## Boundary exacto entre inspección, recovery permitido y acción prohibida

### 1. `execution_inspection_view`

Es la representación consultable y correlada de una ejecución. Agrega snapshot canónico, timeline, policy, approvals, evidence refs, resolution path y findings operativos.

`execution_inspection_view`:

- puede leer y correlacionar records ya persistidos;
- puede exponer findings y recomendar acciones;
- NO cambia estado canónico;
- NO equivale a `execution_record` ni reemplaza al runtime.

### 2. `recovery_action_candidate`

Es la evaluación gobernada de una acción de recovery potencial sobre una ejecución ya inspeccionada. Expresa qué recovery se pide, bajo qué precondiciones, por qué está permitido o bloqueado y si está listo para ejecución.

`recovery_action_candidate`:

- puede existir aunque todavía no se ejecute nada;
- puede quedar `ready_for_execution = false`;
- puede requerir aprobación adicional;
- NO muta runtime, events ni results;
- NO equivale a una signal, command o transición real del kernel.

### 3. Acción real de recovery ejecutada por el kernel

Es la solicitud formal que el kernel acepta y materializa mediante Temporal/workflow/signals/commands compatibles con C.2 y evidencia canónica compatible con C.4.

La acción real:

- sólo puede originarse desde un `recovery_action_candidate` válido o equivalente normativo;
- sólo el kernel la ejecuta y la persiste;
- debe dejar event log, outcome y evidence trail propios;
- puede ser aceptada, rechazada, deduplicada o quedar pendiente.

### Regla normativa fuerte

- `execution_inspection_view` != `recovery_action_candidate` != acción real de recovery.
- El operador puede inspeccionar y solicitar; el kernel ejecuta y persiste.
- Ninguna UI, API o shell operativa puede puentear este boundary.

## Responsabilidades explícitas de la surface de inspección/recovery

- Exponer un snapshot legible del estado actual de una ejecución real.
- Correlacionar `tenant_id`, `trace_id`, `contract_id`, `execution_id`, `policy_decision_id`, `approval_request_id`, `result_id`, `event_ids[]`, `capability_id`, `bundle_digest`, `binding_id` y `provider_ref`.
- Mostrar timeline operativo reconstruible desde evidencia canónica.
- Mostrar decisiones de policy, approvals y restricciones vigentes que afectan recovery.
- Hacer visible el resolution path usado por la ejecución fallida o bloqueada.
- Exponer findings sobre gaps de evidencia, correlación rota, staleness o drift relevante.
- Construir `recovery_action_candidate` para recoveries permitidos.
- Aplicar gates previos antes de habilitar solicitud de recovery.
- Emitir evidence trail del recovery: quién pidió, por qué, bajo qué estado y con qué outcome.
- Integrarse con runtime/Temporal, event log, observabilidad y policy/approvals sin mover sus seams.

## No-responsabilidades explícitas

- No redefinir ni sobrescribir `execution_record`.
- No cambiar estado canónico directamente.
- No forzar `success` ni “cerrar bien” una ejecución sin evidencia terminal del kernel.
- No sobrescribir policy ni reinterpretar un `deny` vigente como `allow`.
- No saltear approvals ni inventar approvals implícitas.
- No editar event log histórico append-only.
- No reescribir evidence refs, results históricos ni artifacts originales.
- No recompilar contrato ni cambiar `compiled_contract` desde recovery.
- No reemplazar Temporal history, event log o policy records con estado de UI.
- No usar observabilidad derivada como única base para acciones de recovery.

## Shape mínimo de `execution_inspection_view`

Debe cubrir como mínimo:

- `inspection_view_id`
- `execution_id`
- `tenant_id`
- `trace_id`
- `contract_id`
- `current_runtime_state`
- `current_outcome_state`
- `policy_decision_refs[]`
- `approval_request_refs[]`
- `result_refs[]`
- `event_refs[]`
- `resolved_capability_ref`
- `resolved_binding_ref`
- `resolved_provider_ref`
- `operator_summary`
- `operator_findings[]`

### Campos normativos adicionales recomendados

- `inspection_generated_at`
- `inspection_scope`
- `environment`
- `contract_fingerprint`
- `policy_decision_id`
- `approval_request_id`
- `result_id`
- `event_ids[]`
- `capability_id`
- `bundle_digest`
- `binding_id`
- `provider_ref`
- `runtime_state_snapshot_ref`
- `timeline_ref`
- `evidence_gaps[]`
- `classification_visibility`
- `recovery_candidates_refs[]`
- `inspection_status`

### Semántica mínima obligatoria

- `inspection_view_id` identifica la vista calculada, no la ejecución.
- `execution_id`, `tenant_id`, `trace_id` y `contract_id` son obligatorios; si falta alguno, la vista no está completa.
- `current_runtime_state` refleja estado canónico del runtime, no inferencia de observabilidad.
- `current_outcome_state` expresa outcome operativo visible (`none`, `success`, `blocked`, `failed`, `unknown_outcome`, `compensation_pending`, `compensated`, etc.) sin colapsar runtime state y outcome state en un solo campo.
- `policy_decision_refs[]`, `approval_request_refs[]`, `result_refs[]` y `event_refs[]` deben ser resolubles o marcarse explícitamente como faltantes.
- `resolved_capability_ref`, `resolved_binding_ref` y `resolved_provider_ref` deben representar la cadena efectivamente usada o el gap explícito.
- `operator_summary` debe ser legible por humano pero anclado a evidencia.
- `operator_findings[]` debe distinguir entre findings informativos, warnings, blockers y evidence gaps.

## Shape mínimo de `recovery_action_candidate`

Debe cubrir como mínimo:

- `recovery_action_candidate_id`
- `execution_id`
- `requested_action`
- `requested_by_subject_id`
- `current_runtime_state`
- `preconditions_refs[]`
- `blocking_constraints[]`
- `reason_codes[]`
- `ready_for_execution` bool
- `requires_additional_approval` bool

### Campos normativos adicionales recomendados

- `tenant_id`
- `trace_id`
- `inspection_view_id`
- `current_outcome_state`
- `requested_at`
- `approval_request_id`
- `policy_decision_id`
- `result_id`
- `event_ids[]`
- `candidate_deduplication_key`
- `proposed_kernel_command`
- `evidence_refs[]`
- `expected_outcome`
- `recovery_scope`
- `supersedes_candidate_id`

### Catálogo mínimo de `requested_action`

- `retry_technical_step`
- `resume_after_approval`
- `request_manual_compensation`
- `acknowledge_unknown_outcome`
- `escalate_for_human_review`

### Semántica mínima obligatoria

- `requested_action` pertenece al catálogo cerrado v1.
- `requested_by_subject_id` es obligatorio para auditoría; no se aceptan recoveries anónimos.
- `current_runtime_state` debe coincidir con la inspección usada para candidatear la acción o declararse stale.
- `preconditions_refs[]` es obligatorio aun cuando quede vacío sólo si la acción no requiere precondiciones adicionales; en v1, las cinco acciones mínimas requieren al menos una referencia resoluble.
- `blocking_constraints[]` explica por qué la acción no está lista o qué la sigue condicionando.
- `reason_codes[]` expresa tanto motivo positivo como bloqueo de recovery.
- `ready_for_execution = true` no ejecuta nada; sólo habilita solicitud al kernel.
- `requires_additional_approval = true` bloquea ejecución directa hasta que approval/policy cierren por camino gobernado.

## Vistas mínimas de inspección

### 1. `execution_summary_view`

Debe responder rápidamente:

- qué ejecución es;
- en qué estado canónico está;
- qué outcome operativo visible tiene;
- qué capability, binding y provider usó;
- cuál es el último reason code material;
- si existe recovery permitido candidateable.

### 2. `execution_timeline_view`

Debe reconstruir, como mínimo:

- arranque de ejecución;
- gate de policy;
- pauses/bloqueos/approvals;
- tramo execution;
- tramo application cuando aplique;
- falla, `unknown_outcome`, compensación o cierre;
- events canónicos y enlaces a observabilidad derivada cuando existan.

La timeline se reconstruye primariamente desde event log canónico; OTel/Temporal pueden enriquecer, nunca sustituir.

### 3. `policy_and_approval_view`

Debe exponer:

- `policy_decision_id` y decisión vigente relevante;
- `approval_request_id` y estado de approval cuando aplique;
- constraints de policy/approval que siguen bloqueando recovery;
- si una acción candidateada requiere nueva approval.

### 4. `evidence_and_event_view`

Debe exponer:

- `event_ids[]` relevantes y sus `event_type`;
- `result_refs[]`, evidence refs y gaps detectados;
- clasificación/redacción aplicada a la vista;
- qué parte de la lectura viene de verdad canónica y qué parte de soporte derivado.

### 5. `recovery_decision_view`

Debe exponer:

- recovery actions candidatas y su estado;
- precondiciones, blockers y approvals faltantes;
- reason codes de habilitación o bloqueo;
- outcome esperado y evidencia mínima requerida antes de pedir ejecución al kernel.

## Correlación visible entre IDs y artifacts

La surface debe volver visible y navegable, como mínimo, la correlación entre:

- `tenant_id`
- `trace_id`
- `contract_id`
- `execution_id`
- `policy_decision_id`
- `approval_request_id`
- `result_id`
- `event_ids[]`
- `capability_id`
- `bundle_digest`
- `binding_id`
- `provider_ref`

### Reglas normativas

1. `execution_id` debe enlazar snapshot, timeline, recovery candidates y command real del kernel.
2. `contract_id` debe poder cruzarse con la ejecución inspeccionada sin contradicción material.
3. `policy_decision_id` y `approval_request_id` deben quedar visibles cuando expliquen un bloqueo, pausa o release.
4. `result_id` debe quedar visible cuando exista evidence terminal o resultado parcial relevante.
5. `event_ids[]` debe permitir reconstruir el corredor aunque fallen exports OTel.
6. `capability_id`, `bundle_digest`, `binding_id` y `provider_ref` deben exponer la cadena exacta de resolution usada por la ejecución.
7. Si algún enlace no es resoluble, la surface debe declararlo como gap explícito; no puede esconderlo detrás de texto resumido.

## Tratamiento explícito de `blocked`, `failed`, `unknown_outcome`, `compensation_pending`, `compensated`

### `blocked`

- Debe mostrar `reason_code` obligatorio.
- Debe indicar si el bloqueo proviene de policy, approval, precondición, dependency o restricción operativa.
- Debe exponer qué recovery permitido podría destrabarlo y qué constraint sigue activo.
- `blocked` NO equivale a `failed`.

### `failed`

- Debe exponer `reason_code`, evidence refs y tramo del lifecycle afectado.
- Debe distinguir falla técnica retryable de falla terminal no recuperable.
- Debe mostrar si existe compensación requerida, opcional o no disponible.

### `unknown_outcome`

- Debe exponer explícitamente por qué no hay certeza material del outcome.
- Debe mostrar evidencia disponible, gaps de timeline y fuentes consultadas.
- Debe soportar `acknowledge_unknown_outcome` o escalación, nunca “éxito por intuición”.

### `compensation_pending`

- Debe indicar por qué la compensación es requerida.
- Debe mostrar si existe acción soportada por runtime o si sólo cabe compensación manual.
- Debe impedir declarar cierre limpio mientras siga pendiente.

### `compensated`

- Debe exponer evidencia de compensación completada o cierre manual equivalente.
- Debe mantener visible la falla/origen causal; compensado no borra la historia.
- Debe mostrar si el cierre es total, parcial o con review humana adicional.

## Recovery operacional permitido

Recovery permitido v1:

1. `retry_technical_step`
   - sólo para fallas o bloqueos técnicamente retryables;
   - nunca para policy deny vigente ni para estados terminales no recoverables.

2. `resume_after_approval`
   - sólo cuando existe `approval_request_id` correlado y evidencia de release/aprobación suficiente;
   - reanuda el corredor duradero, no crea un bypass.

3. `request_manual_compensation`
   - sólo cuando el runtime declara compensación requerida o razonablemente necesaria y no existe compensación automática suficiente;
   - abre corredor manual gobernado, no inventa rollback físico.

4. `acknowledge_unknown_outcome`
   - sólo para reconocer formalmente un caso donde el sistema no puede afirmar outcome material;
   - deja evidencia explícita y puede habilitar follow-up humano.

5. `escalate_for_human_review`
   - para casos con evidencia mínima suficiente pero no resolubles automáticamente;
   - no muta outcome por sí misma.

## Recovery explícitamente prohibido

El operador NO puede:

- forzar `success`;
- sobrescribir policy;
- saltear approvals;
- editar event log histórico;
- mutar estado canónico directamente;
- cambiar `result_id` histórico para maquillar outcome;
- cambiar `binding_id`, `provider_ref` o `bundle_digest` de la ejecución ya ocurrida para “redefinir” qué pasó;
- declarar `compensated` sin evidencia emitida por el kernel o corredor manual normativo;
- borrar `unknown_outcome` por falta de paciencia operativa.

## Gates previos a ejecutar una acción de recovery

Antes de pedir ejecución real al kernel, deben cumplirse como mínimo:

1. existir `execution_inspection_view` resoluble para el `execution_id`;
2. existir `recovery_action_candidate` persistido o request equivalente normativamente completa;
3. correlación mínima completa entre `tenant_id`, `trace_id`, `contract_id` y `execution_id`;
4. `current_runtime_state` y `current_outcome_state` deben ser compatibles con la acción pedida;
5. `preconditions_refs[]` resolubles;
6. `reason_codes[]` presentes;
7. no existir `blocking_constraints[]` activos incompatibles con la acción;
8. no existir deny de policy vigente que bloquee la acción;
9. si aplica, `approval_request_id` y evidencia de release/aprobación deben estar presentes;
10. evidence refs mínimas deben estar completas;
11. la acción debe pasar deduplicación/idempotencia;
12. el actor solicitante debe quedar auditado;
13. si `requires_additional_approval = true`, la ejecución real debe quedar bloqueada hasta cerrar ese gate.

## Evidencia mínima de inspección y recovery

La surface debe poder persistir/consultar como mínimo:

- `inspection_view_id`
- `recovery_action_candidate_id`
- `tenant_id`
- `trace_id`
- `contract_id`
- `execution_id`
- `policy_decision_id`
- `approval_request_id`
- `result_id`
- `event_ids[]`
- `capability_id`
- `bundle_digest`
- `binding_id`
- `provider_ref`
- `current_runtime_state`
- `current_outcome_state`
- `reason_codes[]`
- `preconditions_refs[]`
- `blocking_constraints[]`
- `requested_by_subject_id`
- `requested_at`
- `kernel_action_request_id` o referencia equivalente cuando se ejecute recovery real
- `kernel_action_outcome`
- `kernel_action_recorded_at`

### Evidence trail obligatorio de recovery

Debe quedar trazable:

- quién pidió la acción (`requested_by_subject_id`);
- qué pidió (`requested_action`);
- por qué (`reason_codes[]` + summary);
- bajo qué estado (`current_runtime_state`, `current_outcome_state`);
- con qué precondiciones/evidencia;
- si el kernel la aceptó, rechazó, deduplicó o dejó pendiente;
- con qué outcome final quedó.

## Reason codes mínimos de recovery

La surface debe contemplar al menos:

- `recovery.retry_technical_step`
- `recovery.resume_after_approval`
- `recovery.request_manual_compensation`
- `recovery.acknowledge_unknown_outcome`
- `recovery.escalate_for_human_review`
- `recovery.blocked.policy_still_denies`
- `recovery.blocked.approval_missing`
- `recovery.blocked.compensation_unavailable`
- `recovery.blocked.state_not_recoverable`
- `recovery.not_allowed.operator_boundary`

### Reason codes adicionales recomendados

- `recovery.blocked.missing_execution_id`
- `recovery.blocked.evidence_missing`
- `recovery.blocked.correlation_broken`
- `recovery.blocked.duplicate_request`
- `recovery.blocked.timeline_incomplete`
- `recovery.blocked.provider_resolution_missing`
- `recovery.blocked.execution_already_closed`

## Integración con runtime/Temporal

- D.4 consume queries/signals/estado del runtime definidos por C.2; no los reemplaza.
- La inspección debe usar `execution_record` como snapshot canónico y puede complementar con queries de workflow para estado vivo.
- La acción real de recovery debe traducirse a comandos/signals compatibles con `execution_workflow` o seams del runtime, nunca a writes directos sobre PostgreSQL.
- `retry_technical_step` sólo puede invocar corredor runtime compatible con retry técnico ya modelado por C.2.
- `resume_after_approval` sólo puede operar sobre pausas durables compatibles con `release_execution` o `release_application`.
- `request_manual_compensation` debe respetar `compensation_coordinator` y estados de compensación de C.2.
- `acknowledge_unknown_outcome` y `escalate_for_human_review` deben registrarse como acciones gobernadas sin maquillar el outcome anterior.

## Integración con event log y observabilidad

- Toda inspección debe privilegiar `event_record` canónico append-only definido en C.4.
- `execution_timeline_view` se reconstruye desde `event_record`; observabilidad derivada sólo enriquece cuando existe.
- Si falla export OTel/LGTM pero el event log está íntegro, la inspección sigue siendo válida.
- Si la observabilidad derivada muestra algo que no está soportado por event log canónico, la surface debe marcar inconsistencia, no promover esa lectura a verdad.
- Recovery real ejecutado por el kernel debe emitir evidencia canónica propia y sólo después proyectarse a observabilidad.
- El operador no puede editar ni “corregir” el event log histórico; sólo puede agregar acciones/evidencia nuevas por caminos gobernados.

## Integración con policy/approvals cuando apliquen

- D.4 puede leer `policy_decision_record` y `approval_request` correlados; no los reescribe.
- Si una acción de recovery requiere permiso adicional, la surface debe marcar `requires_additional_approval = true`.
- `resume_after_approval` exige `approval_request_id` correlado y evidencia material de release/aprobación.
- Si policy sigue en `deny`, `retry_technical_step` y `resume_after_approval` deben quedar bloqueados con `recovery.blocked.policy_still_denies`.
- Escalaciones y compensaciones manuales pueden requerir approval/policy según sensibilidad; eso debe verse antes de ejecutar, no después.

## Idempotencia y deduplicación de recovery actions

- Debe existir deduplicación por acción materialmente equivalente sobre la misma ejecución y mismo estado base.
- La deduplicación mínima debe considerar `tenant_id`, `execution_id`, `requested_action`, `current_runtime_state`, `current_outcome_state`, `reason_codes[]` relevantes y fingerprint de precondiciones.
- Doble click, retry de red o reenvío operativo no pueden disparar recoveries duplicados.
- Si el kernel detecta una solicitud equivalente ya en curso o ya aceptada, debe devolver outcome de deduplicación en vez de ejecutar de nuevo.
- Un cambio material del estado base invalida la equivalencia y puede requerir nuevo candidate.

## Tests borde mínimos

1. vista sin `execution_id`;
2. correlación rota entre `execution_id` y `contract_id`;
3. `blocked` sin `reason_code`;
4. `failed` sin evidence refs;
5. `unknown_outcome` sin timeline suficiente;
6. `compensation_pending` sin action disponible;
7. operator intenta forzar success;
8. operator intenta editar evidence histórica;
9. retry técnico pedido en estado no recuperable;
10. `resume_after_approval` pedido sin `approval_request_id`;
11. policy sigue deny pero operator insiste en retry;
12. recovery action candidate sin preconditions refs;
13. recovery duplicado por doble click/retry;
14. timeline incompleta por export observability fallida pero event log canónico presente;
15. timeline presente en observability pero faltante en event log canónico;
16. `resolved_provider_ref` ausente en inspection view;
17. binding cambiado después del fallo original;
18. execution ya cerrada intenta recovery activo;
19. human review escalada sin evidence mínima;
20. compensation requerida pero no disponible;
21. `blocked` por approval liberada pero sin evento canónico `approval.released`;
22. `unknown_outcome` con evidencia parcial en Temporal pero sin `event_record` canónico suficiente;
23. `compensated` visible sin evento/evidencia de compensación correlada;
24. recovery candidate listo sobre inspection view stale respecto del runtime actual;
25. intento de `acknowledge_unknown_outcome` sobre ejecución ya marcada `failed` con outcome cierto;
26. trace correlado pero `tenant_id` inconsistente entre execution y events;
27. `request_manual_compensation` sin capability de compensación ni pathway manual documentado;
28. operator intenta reescribir `provider_ref` histórico para justificar retry;
29. event log íntegro pero observabilidad contradice timestamps del timeline;
30. recovery real aceptado por kernel pero sin evidence trail de quién lo solicitó.

## Criterios de aceptación de D.4

- Queda explícito que la surface de inspección lee y correlaciona; NO redefine estado canónico.
- Queda explícito el boundary entre `execution_inspection_view`, `recovery_action_candidate` y acción real ejecutada por el kernel.
- Queda explícito que la surface puede preparar y solicitar recovery, pero el kernel sigue siendo quien ejecuta y persiste.
- Quedan definidos los shapes mínimos de `execution_inspection_view` y `recovery_action_candidate`.
- Quedan definidas las vistas mínimas `execution_summary_view`, `execution_timeline_view`, `policy_and_approval_view`, `evidence_and_event_view` y `recovery_decision_view`.
- Queda definida la correlación visible mínima entre IDs y artifacts materiales del corredor.
- Queda tratamiento explícito para `blocked`, `failed`, `unknown_outcome`, `compensation_pending` y `compensated`.
- Queda fijado el catálogo mínimo de recovery permitido y la lista de acciones prohibidas al operador.
- Quedan definidos gates previos a ejecutar recovery, evidencia mínima y reason codes mínimos.
- Queda explícita la integración con runtime/Temporal, event log, observabilidad y policy/approvals sin contradecir C.2, C.4, C.6 ni D.3.
- Queda exigida idempotencia y deduplicación de recovery actions.
- Queda cubierto un set mínimo de tests borde suficiente para implementación y verificación de la surface v1.
