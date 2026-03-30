# D.5 — AI-friendly development and maintenance surface

## Objetivo

Definir la surface v1 para developers, reviewers y operators que permita **leer, explicar, depurar semánticamente, inspeccionar artifacts y preparar mantenimiento gobernado** sobre el kernel ya cerrado en C.1-C.6 y las surfaces D.1-D.4, sin crear una verdad paralela, sin convertir a la IA en autoridad de apply y sin reintroducir distribution layer bajo el nombre de “maintenance”.

D.5 trabaja sobre artifacts, contracts, runtime, event log y evidence ya cerrados. No recompone semántica desde texto libre ni desde dashboards aislados: consume la misma base de verdad del kernel y de las surfaces previas.

## Principios de implementación de la surface IA-friendly

- La surface IA-friendly **asiste lectura y preparación**, no ejecuta mutación canónica por sí misma.
- La misma evidencia canónica que sostiene operators y runtime debe sostener la ayuda asistida; no existen resúmenes paralelos como source of truth.
- La IA puede asistir en **lectura, explicación, shaping, troubleshooting y preparación de acciones gobernadas**.
- La IA **NO** puede convertirse en autoridad de apply, recovery o maintenance execution.
- Toda ayuda relevante debe ser **tenant-scoped, auditable, correlable e idempotente**.
- Si la surface no puede sostener una explicación con evidencia suficiente, debe degradar a ambigüedad explícita o `debug_incomplete`; nunca inventar certeza.
- `semantic_debug_view`, `maintenance_action_candidate` y la acción real gobernada del sistema son boundaries distintos y no colapsables.
- La observabilidad derivada puede enriquecer una explicación, pero **PostgreSQL/event log/records canónicos** siguen siendo la base normativa.
- Distribution layer, tenant activation y rollout de consumo siguen fuera del roadmap y no pueden reaparecer como maintenance.
- El fail mode por defecto es **fail-closed** ante evidencia insuficiente, clasificación restrictiva, correlación incompleta o intento de bypass de governance.

## Boundary exacto entre ayuda asistida, debugging semántico y acción gobernada

### 1. Ayuda asistida

Es la capacidad de la surface para:

- leer artifacts y records existentes;
- explicar correlaciones;
- resumir findings anclados a evidencia;
- proponer shaping de un diagnóstico o de una solicitud de mantenimiento;
- preparar candidates gobernados.

La ayuda asistida NO cambia estado canónico ni equivale a una orden ejecutable.

### 2. `semantic_debug_view`

Es la vista calculada que vuelve legible un caso correlando conversación, intent, proposal, preview, contrato, runtime, policy, event log y evidencia. Sirve para entender **qué pasó, por qué parece haber pasado y qué gaps siguen abiertos**.

`semantic_debug_view`:

- puede leer y correlacionar records ya persistidos;
- puede exponer explicación, ambigüedades y evidence refs;
- puede sugerir próximos pasos o candidatear mantenimiento permitido;
- NO muta artifacts ni runtime;
- NO reemplaza `execution_inspection_view`, `compiled_contract`, `proposal_draft` ni `event_record`.

### 3. `maintenance_action_candidate`

Es la formalización gobernada de una acción de mantenimiento potencial preparada desde esta surface. Expresa qué se quiere pedir, sobre qué targets, con qué precondiciones y bajo qué governance.

`maintenance_action_candidate`:

- puede existir sin ejecutar nada;
- puede quedar bloqueado o requerir review humana;
- puede ser deduplicado o marcado stale;
- NO equivale a una acción real del kernel;
- NO puede materializar apply directo.

### 4. Acción real de mantenimiento gobernada por el sistema

Es la solicitud formal aceptada por el sistema sobre el corredor canónico correspondiente. Sólo el sistema/kernel la ejecuta, la audita y la persiste.

La acción real:

- sólo puede originarse desde un `maintenance_action_candidate` válido o request equivalente normativamente completa;
- debe pasar policy, approvals, clasificación, deduplicación e idempotencia;
- debe emitir su propio evidence trail y outcome;
- puede ser aceptada, rechazada, bloqueada, deduplicada o derivada a review humana.

### Regla normativa fuerte

- ayuda asistida != `semantic_debug_view` != `maintenance_action_candidate` != acción real gobernada.
- La IA puede sugerir y preparar; el sistema gobierna y ejecuta.
- Ninguna UI, API, agente o shell puede puentear este boundary.

## Responsabilidades explícitas de la surface D.5

- Exponer debugging semántico sobre artifacts ya existentes.
- Permitir tracing semántico desde conversación -> intent -> proposal -> preview -> runtime.
- Permitir inspección gobernada de config, artifacts, contracts, workflows, runtime y evidence.
- Mostrar correlación entre IDs, refs, hashes, fingerprints y eventos relevantes.
- Explicar gaps, ambigüedades, staleness y contradicciones sin maquillarlas.
- Construir `semantic_debug_view` con evidence trail suficiente.
- Construir `maintenance_action_candidate` sólo para acciones permitidas del catálogo v1.
- Aplicar gates previos antes de declarar una acción “lista para ejecución”.
- Registrar por qué una sugerencia queda bloqueada, degradada o derivada a review humana.
- Integrarse con intake/proposal/preview/inspection sin reabrir sus boundaries.
- Integrarse con artifacts del kernel y event log sin reemplazar source of truth.

## No-responsabilidades explícitas de la surface D.5

- No ejecutar apply real.
- No ejecutar recovery real por fuera de D.4.
- No mutar `compiled_contract`, `execution_record`, `policy_decision_record`, `approval_request`, `result_record` ni `event_record` directamente.
- No reemplazar policy enforcement, approval authority ni clasificación final.
- No inventar resumen operativo donde falta evidencia.
- No usar observabilidad derivada como única fuente para debugging o maintenance.
- No recompilar, reproyectar, reconstruir o refrescar nada por mera sugerencia textual sin gates formales.
- No operar distribution layer, tenant activation, rollout de consumo, instalación o bootstrap operacional fuera del roadmap.
- No tratar conversación libre o respuesta de IA como evidencia suficiente de verdad material.
- No convertir `ready_for_execution = true` en ejecución implícita.

## Shape mínimo de `semantic_debug_view`

Debe cubrir como mínimo:

- `debug_view_id`
- `tenant_id`
- `trace_id`
- `source_refs[]`
- `artifact_refs[]`
- `contract_refs[]`
- `execution_refs[]`
- `policy_refs[]`
- `event_refs[]`
- `semantic_summary`
- `explanation_chain[]`
- `ambiguities[]`
- `evidence_refs[]`

### Campos normativos adicionales recomendados

- `debug_scope`
- `debug_status`
- `generated_at`
- `requested_by_subject_id`
- `conversation_turn_refs[]`
- `intent_refs[]`
- `proposal_refs[]`
- `preview_refs[]`
- `runtime_state_snapshot_refs[]`
- `classification_visibility`
- `reason_codes[]`
- `evidence_gaps[]`
- `staleness_status`
- `maintenance_candidate_refs[]`
- `deduplication_key`

### Semántica mínima obligatoria

- `debug_view_id` identifica la vista calculada, no el caso fuente.
- `tenant_id` es obligatorio; no existe debug cross-tenant implícito.
- `trace_id` es obligatorio salvo casos heredados explícitamente degradados; si falta, la vista no puede declararse completa.
- `source_refs[]` debe apuntar a turns, intents, proposals, previews u otros artifacts de surface que expliquen el caso.
- `artifact_refs[]`, `contract_refs[]`, `execution_refs[]`, `policy_refs[]` y `event_refs[]` deben ser resolubles o declararse como gaps explícitos.
- `semantic_summary` debe ser legible para humano pero anclado a evidencia; no puede afirmar hechos que `explanation_chain[]` o `evidence_refs[]` no sostienen.
- `explanation_chain[]` debe distinguir pasos confirmados, inferidos, bloqueados y ambiguos.
- `ambiguities[]` debe listar TODA incertidumbre material abierta; una ambigüedad no resuelta no puede presentarse como hecho cerrado.
- `evidence_refs[]` es obligatorio; no existe debug view normativo basado sólo en resumen.

## Shape mínimo de `maintenance_action_candidate`

Debe cubrir como mínimo:

- `maintenance_action_candidate_id`
- `tenant_id`
- `requested_by_subject_id`
- `action_type`
- `target_refs[]`
- `preconditions_refs[]`
- `governance_requirements[]`
- `reason_codes[]`
- `ready_for_execution` bool
- `requires_human_review` bool

### Campos normativos adicionales recomendados

- `trace_id`
- `debug_view_id`
- `candidate_status`
- `requested_at`
- `target_scope`
- `blocking_constraints[]`
- `evidence_refs[]`
- `classification_visibility`
- `policy_refs[]`
- `approval_refs[]`
- `expected_outcome`
- `candidate_deduplication_key`
- `supersedes_candidate_id`
- `staleness_status`

### Action types mínimos esperados

- `request_recompile_contract`
- `request_rebuild_preview`
- `request_refresh_registry_resolution`
- `request_reproject_event_view`
- `request_human_review`

### Semántica mínima obligatoria

- `action_type` pertenece al catálogo cerrado v1; no se aceptan acciones libres.
- `requested_by_subject_id` es obligatorio para auditoría.
- `target_refs[]` es obligatorio y debe ser compatible con `action_type`.
- `preconditions_refs[]` debe listar precondiciones resolubles o declarar explícitamente que faltan.
- `governance_requirements[]` debe enumerar policy/approval/classification/gates aplicables.
- `reason_codes[]` debe explicar tanto la motivación como el bloqueo si existe.
- `ready_for_execution = true` no ejecuta nada; sólo indica que el candidate superó los gates de surface.
- `requires_human_review = true` bloquea la ejecución automática aunque el resto parezca completo.

## Vistas/flows mínimos para developers y operators

### Vistas mínimas

1. `semantic_trace_view`
   - explica el corredor conversación -> intent -> proposal -> preview -> runtime.

2. `artifact_inspection_view`
   - inspecciona artifacts, config, manifests, bundles, contracts y evidence refs visibles.

3. `contract_runtime_debug_view`
   - correlaciona contrato, runtime, policy, approvals y event log.

4. `maintenance_preparation_view`
   - muestra candidates, blockers, governance requirements y evidencia.

5. `assisted_explanation_view`
   - presenta explicación legible, ambigüedades, fuentes y límites de confianza.

### Flows mínimos esperados

#### `semantic_trace_flow`

1. recibir pregunta o pedido de explicación;
2. resolver `tenant_id` + `trace_id` o degradar explícitamente;
3. correlacionar `conversation_turn`/`intent_input`/`proposal_draft`/`preview_candidate`/runtime/eventos;
4. construir `semantic_debug_view`;
5. exponer summary, chain, ambiguities y evidence trail.

#### `artifact_inspection_flow`

1. recibir artifact/config/contract/workflow ref;
2. validar acceso por tenant, clasificación y visibility;
3. resolver artifact base desde source of truth canónica;
4. exponer shape, version, fingerprint, evidence y vínculos relacionados;
5. registrar si la lectura quedó parcial, redaccionada o bloqueada.

#### `contract_and_runtime_debug_flow`

1. partir de `contract_ref`, `execution_ref` o `trace_id`;
2. correlacionar compiled contract, policy decisions, approvals, execution y events;
3. detectar contradicciones, gaps o staleness;
4. construir explicación semántica;
5. si corresponde, candidatear mantenimiento permitido.

#### `maintenance_candidate_flow`

1. partir de un hallazgo validado en `semantic_debug_view` o `execution_inspection_view`;
2. mapear el hallazgo a un `action_type` permitido;
3. exigir targets, precondiciones, governance requirements y evidence refs;
4. correr gates previos;
5. persistir `maintenance_action_candidate` listo, bloqueado o derivado a review humana.

#### `assisted_explanation_flow`

1. recibir una pregunta de developer u operator;
2. resolver el scope permitido de lectura;
3. recuperar evidencia canónica relevante;
4. producir explicación con `semantic_summary`, `explanation_chain[]` y `ambiguities[]`;
5. dejar trazabilidad de que hubo asistencia y sobre qué evidencia se apoyó.

## Debugging semántico

El debugging semántico de D.5 debe responder como mínimo:

- qué pidió originalmente la conversación y cómo fue interpretado;
- qué intent se emitió o por qué no se emitió;
- qué proposal se abrió, con qué diff humano/material;
- qué previews/simulaciones influyeron en el estado posterior;
- qué contrato, policy, approval, runtime o event log explican el resultado observado;
- qué parte está confirmada y qué parte sigue ambigua.

### Reglas normativas

1. Toda explicación debe apoyarse en evidence refs resolubles.
2. La cadena explicativa debe separar hechos confirmados de inferencias.
3. Una contradicción entre artifacts debe exponerse como contradicción, no resolverse por optimismo.
4. Si la evidencia está redaccionada o incompleta, la explicación debe declarar sus límites.
5. La surface puede explicar por qué algo quedó bloqueado, falló o quedó incompleto; no puede reinterpretar policy o runtime para “hacer que cierre”.
6. La explicación semántica puede atravesar surfaces D.1-D.4, pero siempre leyendo artifacts ya persistidos.
7. El resultado mínimo del debugging no es “tener razón”, sino producir una lectura auditable y accionable.

## Inspección gobernada de config/artifacts/contracts/workflows

La surface debe permitir inspeccionar de forma gobernada, como mínimo:

- configuraciones declarativas relevantes del cambio o del caso;
- artifacts de surface (`intent_input`, `proposal_draft`, `governed_patchset_candidate`, `preview_candidate`, `simulation_result`);
- `compiled_contract` y evidencia de compilación;
- `execution_record`, `result_record`, `policy_decision_record`, `approval_request`;
- `event_record` y projections consultables derivadas;
- refs de capability, bundle, binding y provider cerrados en C.5;
- workflows o snapshots de runtime siempre como lectura subordinada a la verdad canónica.

### Reglas normativas

- La inspección respeta clasificación, redacción y permisos existentes.
- Un artifact no accesible por clasificación debe responder con bloqueo o vista redaccionada, nunca con bypass implícito.
- El sistema debe distinguir explícitamente qué lectura viene de records canónicos y cuál es enriquecimiento derivado.
- Una vista de workflow no puede reemplazar event log ni `execution_record`.
- La inspección sobre registry/resolution sólo puede leer y preparar refresh permitido; no puede mutar resolution histórica ya usada por una ejecución pasada.

## Surface de mantenimiento gobernado

D.5 sólo permite preparar mantenimiento dentro de un catálogo cerrado de acciones permitidas:

1. `request_recompile_contract`
   - para casos donde el contrato o evidencia de compilación deban ser regenerados formalmente;
   - requiere `contract_ref` resoluble y evidencia de por qué la recompilación es pertinente.

2. `request_rebuild_preview`
   - para regenerar preview/simulación sobre proposal/patchset vigentes cuando haya staleness, inputs actualizados o evidencia inconsistente;
   - nunca para reabrir proposal terminal cerrada sin corredor normativo.

3. `request_refresh_registry_resolution`
   - para pedir refresco gobernado de capability/binding/provider resolution sobre artifacts vigentes;
   - no redefine retrospectivamente la cadena usada por ejecuciones históricas.

4. `request_reproject_event_view`
   - para regenerar proyecciones o vistas derivadas basadas en event log canónico;
   - requiere event log suficiente; no inventa eventos faltantes.

5. `request_human_review`
   - para derivar el caso cuando la evidencia existe pero la surface no puede cerrarlo sin review humana.

### Acción de mantenimiento explícitamente prohibida

- apply directo desde IA o desde candidate.
- recovery directo fuera del corredor D.4.
- mutación histórica de event log, contracts o results.
- operación sobre distribution layer.
- rollout, activación tenant, instalación o despliegue operativo bajo el alias de maintenance.

## Límites de automatización asistida

- La IA puede sugerir, resumir, explicar, troubleshooting y preparar candidates.
- La IA no puede ejecutar por sí sola ninguna acción con efecto material.
- La IA no puede marcar una aprobación como satisfecha ni reinterpretar un deny como allow.
- La IA no puede promover un candidate a acción real por “alta confianza”.
- La IA no puede usar datos redaccionados insuficientes para afirmar causalidad cerrada.
- La IA no puede salir de la base de verdad del kernel para inventar correlaciones desde memoria conversacional o observabilidad aislada.
- La IA no puede crear nuevos `action_type` fuera del catálogo v1.
- La IA no puede operar distribution layer, aunque el pedido se formule como debugging o maintenance.

## Gates previos a maintenance actions

Antes de declarar un `maintenance_action_candidate` como listo para ejecución, deben cumplirse como mínimo:

1. existir `semantic_debug_view` o vista de inspección correlada suficiente para justificar la acción;
2. `tenant_id` resoluble y consistente;
3. `trace_id` resoluble o degradación explícita aprobada por policy del caso;
4. `action_type` perteneciente al catálogo permitido;
5. `target_refs[]` presentes y resolubles;
6. `preconditions_refs[]` presentes y resolubles;
7. `governance_requirements[]` completos;
8. `reason_codes[]` presentes y compatibles con la acción;
9. evidencia suficiente anclada a records canónicos;
10. no existir `blocking_constraints[]` activos incompatibles;
11. clasificación/redacción compatible con el nivel de detalle mostrado y con la acción pedida;
12. deduplicación/idempotencia superada;
13. si aplica, approval/policy adicional resuelta o pendiente explícitamente;
14. si el caso mantiene ambigüedad material, `requires_human_review = true`.

## Evidencia mínima de debugging/mantenimiento

Debe quedar persistible/consultable, como mínimo:

- `debug_view_id`
- `maintenance_action_candidate_id`
- `tenant_id`
- `trace_id`
- `requested_by_subject_id`
- `source_refs[]`
- `artifact_refs[]`
- `contract_refs[]`
- `execution_refs[]`
- `policy_refs[]`
- `event_refs[]`
- `semantic_summary`
- `explanation_chain[]`
- `ambiguities[]`
- `evidence_refs[]`
- `action_type`
- `target_refs[]`
- `preconditions_refs[]`
- `governance_requirements[]`
- `reason_codes[]`
- `ready_for_execution`
- `requires_human_review`
- `candidate_deduplication_key` o equivalente
- `recorded_at`

### Evidence trail obligatorio

Debe quedar trazable:

- quién pidió la explicación o candidateó la acción;
- qué evidencia fue consultada;
- qué summary se produjo;
- qué ambigüedades quedaron abiertas;
- qué acción se sugirió y por qué;
- qué gates pasaron y cuáles bloquearon;
- si hubo derivación a review humana;
- si luego la acción real fue aceptada, rechazada o deduplicada por el sistema.

## Reason codes mínimos de D.5

La surface debe soportar como mínimo:

- `debug.explain_trace`
- `debug.inspect_artifact`
- `debug.inspect_contract`
- `debug.inspect_runtime`
- `maintenance.request_recompile_contract`
- `maintenance.request_rebuild_preview`
- `maintenance.request_refresh_registry_resolution`
- `maintenance.request_reproject_event_view`
- `maintenance.request_human_review`
- `maintenance.blocked.operator_boundary`

### Reason codes adicionales recomendados

- `debug.blocked.missing_trace`
- `debug.blocked.insufficient_evidence`
- `debug.blocked.classification_restricted`
- `maintenance.blocked.missing_target`
- `maintenance.blocked.missing_preconditions`
- `maintenance.blocked.governance_requirements_missing`
- `maintenance.blocked.out_of_scope`
- `maintenance.blocked.distribution_layer`
- `maintenance.blocked.duplicate_candidate`
- `maintenance.blocked.requires_human_review`

## Integración con intake/proposal/preview/inspection

- Con D.1: puede leer `conversation_turn`, `intake_session`, `intent_candidate`, `intent_input` para explicar cómo se shapeó la intención, pero no reabre decisions de intake.
- Con D.2: puede leer `proposal_draft` y `governed_patchset_candidate` para explicar cambios propuestos, diffs y gates, pero no promociona por fuera del workspace.
- Con D.3: puede leer `preview_candidate`, `simulation_result` y `composite_preview` para explicar bloqueos o staleness, pero no sustituye simulación ni crea `apply_candidate`.
- Con D.4: puede leer `execution_inspection_view` y findings operativos para profundizar debugging o preparar maintenance permitido, pero no reemplaza recovery surface.
- Si un caso cruza múltiples surfaces, D.5 debe mantener correlación visible entre sus IDs y no resumirlas como un único blob textual.

## Integración con artifacts del kernel y event log

- `semantic_debug_view` y `maintenance_action_candidate` deben referenciar records canónicos del kernel, no snapshots paralelos de UI.
- El event log append-only de C.4 es la base para timeline y causalidad; OTel/LGTM sólo enriquecen.
- `compiled_contract`, `execution_record`, `policy_decision_record`, `approval_request`, `result_record` y `event_record` siguen siendo los artifacts normativos.
- Las refs a capability/binding/provider deben ser compatibles con C.5 y mostrar la cadena resoluble.
- Cuando una proyección derivada contradiga al event log o al record canónico, prevalece la verdad canónica y la contradicción debe quedar explicitada.
- La surface IA-friendly trabaja sobre la misma base de verdad del kernel, no sobre resúmenes paralelos.

## Idempotencia y deduplicación de maintenance actions

- Debe existir deduplicación por tenant, `action_type`, targets materiales, trace/caso y precondiciones equivalentes.
- Un retry técnico o una repetición conversacional no debe producir candidates materialmente duplicados.
- `candidate_deduplication_key` o equivalente debe ser estable para solicitudes equivalentes.
- Si cambia materialmente el target, las precondiciones o la evidencia base, debe emitirse un nuevo candidate o marcar el anterior stale/superseded.
- La deduplicación nunca puede ocultar evidencia: debe quedar visible que hubo retry, reuse o supersession.
- La acción real del sistema debe volver a validar idempotencia, aunque el candidate de surface ya haya sido deduplicado.

## Tests borde mínimos

1. `semantic_debug_view` sin `trace_id` -> la vista queda `debug_incomplete` o equivalente y no se presenta como completa.
2. `source_refs[]` sin correlación suficiente -> se bloquea el trace semántico con gap explícito.
3. `explanation_chain[]` contradice `evidence_refs[]` -> la vista es inválida y no puede emitirse como cerrada.
4. inspección de artifact no accesible por clasificación -> respuesta redaccionada o bloqueada, nunca bypass.
5. `maintenance_action_candidate` sin `target_refs[]` -> candidate inválido.
6. `maintenance_action_candidate` intenta apply directo -> bloqueo normativo inmediato.
7. IA sugiere acción fuera de scope del sistema -> reason code `maintenance.blocked.out_of_scope`.
8. IA intenta operar distribution layer -> bloqueo explícito con referencia a scope fuera de roadmap.
9. `request_recompile_contract` pedido sin `contract_ref` -> candidate inválido.
10. `request_rebuild_preview` pedido sobre proposal cerrada terminalmente -> bloqueo o derivación a review humana.
11. `request_refresh_registry_resolution` pedido sin refs de capability/binding -> candidate inválido.
12. `request_reproject_event_view` con event log faltante -> bloqueo por evidencia insuficiente.
13. debug view muestra resumen sin `evidence_refs[]` -> inválido.
14. ambigüedad no resuelta pero `semantic_summary` la presenta como hecho -> inválido.
15. suggestion automatizada sin `governance_requirements[]` -> candidate inválido.
16. `ready_for_execution = true` en acción no permitida -> rechazo normativo.
17. candidate duplicado por retry -> debe deduplicarse o marcar superseded.
18. explicación generada sobre datos demasiado redaccionados -> debe degradar confianza y marcar límites.
19. human review requerida pero `requires_human_review = false` -> candidate inválido.
20. la surface usa observabilidad derivada como única fuente -> rechazo por contradicción con C.4.
21. `tenant_id` ausente o inconsistente entre refs -> bloqueo del debug y del maintenance candidate.
22. `trace_id` resuelve a múltiples corredores incompatibles -> se expone contradicción, no se fusiona.
23. artifact inspection intenta leer snapshot de UI no canonizado -> rechazo por fuente inválida.
24. refresh de registry intenta reinterpretar retrospectivamente una ejecución histórica -> bloqueo.
25. reproject de event view intenta inventar eventos no emitidos -> bloqueo.
26. contract/runtime debug con policy deny vigente pero IA sugiere bypass -> bloqueo `maintenance.blocked.operator_boundary`.
27. `request_human_review` sin evidencia mínima del caso -> candidate incompleto.
28. candidate preparado desde debug stale después de cambio material del artifact -> debe marcarse stale o recalcularse.

## Criterios de aceptación de D.5

- Existe boundary explícito entre ayuda asistida, `semantic_debug_view`, `maintenance_action_candidate` y acción real gobernada.
- La IA puede asistir en lectura, explicación, shaping, troubleshooting y preparación de acciones gobernadas, pero NO puede convertirse en autoridad de apply ni bypass de governance.
- La surface define y usa `semantic_debug_view` y `maintenance_action_candidate` con shapes mínimos completos.
- Existen al menos los flows mínimos: `semantic_trace_flow`, `artifact_inspection_flow`, `contract_and_runtime_debug_flow`, `maintenance_candidate_flow` y `assisted_explanation_flow`.
- La surface permite tracing semántico desde conversación -> intent -> proposal -> preview -> runtime.
- La surface permite inspección gobernada de artifacts/configs/contracts/workflows sin romper clasificación ni boundaries previos.
- La surface sólo prepara maintenance actions del catálogo permitido y nunca operación sobre distribution layer.
- Los límites entre sugerir y ejecutar quedan explícitos y fail-closed.
- Toda explicación y toda acción asistida relevante deja evidence trail.
- La surface opera sobre la misma base de verdad del kernel y del event log, no sobre resúmenes paralelos.
- Se definen gates previos, reason codes mínimos, deduplicación e idempotencia.
- Se cubren tests borde suficientes para detectar bypass de governance, evidencia insuficiente, duplicación y confusión entre debug y ejecución.
