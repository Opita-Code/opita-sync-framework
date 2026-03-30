# D.3 — Simulation, diff and preview surface

## Objetivo

Definir la surface v1 que toma un `proposal_draft` ya gobernado y un `governed_patchset_candidate` consistente para producir una **lectura previa, auditable y promocionable** del cambio antes de crear un `apply_candidate`, sin confundir preview con apply real y sin contradecir D.2, A.3, A.4, C.3, C.4 ni C.6.

El objetivo de D.3 NO es ejecutar el cambio, NO es emitir approvals reales, NO es sustituir enforcement de policy del kernel, NO es compilar contratos ni correr runtime. El objetivo es fijar el contrato exacto de la preview surface donde el cambio candidateado se vuelve legible, simulable y gateable antes de la promoción explícita hacia `apply_candidate`.

## Principios de implementación de preview/simulación

- Preview y simulación corren sobre `proposal_draft` + `governed_patchset_candidate`, nunca sobre charla informal.
- La conversación puede originar el proceso, pero la verdad operativa de preview vive en artifacts persistidos, versionados y correlables.
- `preview_candidate` NO es `apply_candidate`.
- `simulation_result` NO sustituye enforcement real del kernel.
- El `material_diff` sigue siendo la base normativa; el `human_diff` existe para revisión humana, no para reemplazar materialidad.
- Toda simulación debe declarar qué inputs usó, con qué fidelidad y con qué límites.
- Si la simulación no puede hacerse con suficiente fidelidad, el resultado debe ser `preview_incomplete`, no un falso `preview_ok`.
- La preview surface puede recomendar no avanzar o dejar un candidate fuera de promoción, pero no puede bloquear runtime real fuera del pipeline gobernado.
- El fail mode por defecto es fail-closed respecto de promoción a `apply_candidate`.
- Preview, diff y simulación deben ser tenant-scoped, trazables, idempotentes y auditables.
- Resultados compuestos no pueden maquillar sub-resultados bloqueados, stale o incompletos.

## Boundary exacto entre `proposal_draft`, `preview_candidate`, `simulation_result` y `apply_candidate`

### 1. `proposal_draft`

Es la propuesta gobernada revisable definida en D.2. Explica qué se quiere cambiar, por qué, en qué scope y con qué evidencias, pero todavía no representa una sesión formal de preview.

### 2. `preview_candidate`

Es el artifact de D.3 que declara que un `proposal_draft` y un `governed_patchset_candidate` determinados entran a la surface de preview/simulación. Reúne el diff humano, el diff material, el scope y las solicitudes de simulación necesarias para revisar el cambio candidateado.

`preview_candidate` NO aplica cambios, NO crea `execution_record`, NO consume approvals reales y NO habilita por sí mismo el paso a apply.

### 3. `simulation_result`

Es el output normalizado de una simulación específica sobre un `preview_candidate`. Cada `simulation_result` pertenece a una familia cerrada (`policy`, `approval`, `classification`, `risk`) y describe estado, evidencia, límites de fidelidad, findings y confidence.

`simulation_result` es evidencia previa, no decisión final del kernel.

### 4. `apply_candidate`

Es el artifact posterior al preview que declara que el cambio candidateado quedó apto para entrar al flujo de apply real. Sólo puede existir después de que el preview compuesto quede consistente y de que cierren los gates previos definidos por D.2 y D.3.

### Regla normativa de boundary

- `proposal_draft` = propuesta gobernada revisable.
- `preview_candidate` = paquete formal de preview/simulación sobre esa propuesta y su patchset material.
- `simulation_result` = salida auditable de cada simulación de preview.
- `apply_candidate` = promoción posterior hacia apply real.

Nunca deben colapsarse en el mismo objeto ni tratarse como equivalentes.

## Responsabilidades explícitas de la preview surface

- Recibir un `proposal_draft` vigente y un `governed_patchset_candidate` consistente.
- Validar que el candidate tenga el mínimo necesario para entrar en preview.
- Exponer simultáneamente `human_diff` y `material_diff` del mismo cambio candidateado.
- Crear y persistir un `preview_candidate` correlado con tenant, sujeto, trace y proposal version.
- Ejecutar simulaciones v1 de policy, approval, classification/redaction y risk/impact.
- Publicar `simulation_result` por familia, con evidencia e inputs refs explícitos.
- Construir un `composite_preview` que sintetice sub-previews sin borrar sus diferencias.
- Detectar inconsistencias entre diffs, scope, patchset, inputs y resultados simulados.
- Recomendar continuar, revisar o no promocionar hacia apply candidate.
- Emitir reason codes, findings, warnings y estado compuesto auditable.
- Marcar preview stale cuando cambie el patchset material, los artifacts base o los inputs normativos.
- Aplicar idempotencia y deduplicación para previews/simulaciones equivalentes.
- Integrarse con artifacts del kernel y con el event log sin mover los boundaries de C.3/C.4/C.6.

## No-responsabilidades explícitas de la preview surface

- No ejecutar apply real.
- No crear `execution_record` ni arrancar workflows del runtime.
- No reemplazar el enforcement real de Cerbos ni del kernel.
- No emitir `approval_request` real ni consumir decisiones humanas reales.
- No decidir clasificación/redacción final de outputs productivos.
- No mutar artifacts base fuera del `governed_patchset_candidate`.
- No tratar chat libre, resúmenes informales o prompts como source of truth.
- No promocionar automáticamente un cambio a `apply_candidate` por “se ve razonable”.
- No ocultar sub-previews bloqueados dentro de un resumen optimista.
- No bloquear runtime real fuera del pipeline gobernado; sólo puede bloquear la promoción dentro de esta surface.
- No reemplazar evidencia canónica del kernel con telemetría derivada o UI.

## Shape mínimo de `preview_candidate`

`preview_candidate` debe cubrir como mínimo:

- `preview_candidate_id`
- `proposal_draft_id`
- `patchset_candidate_id`
- `tenant_id`
- `subject_id`
- `artifacts_in_scope[]`
- `human_diff_ref`
- `material_diff_ref`
- `preview_scope`
- `current_preview_state`
- `simulation_requests[]`

### Campos normativos adicionales recomendados

- `preview_version`
- `proposal_version`
- `patchset_version`
- `trace_id`
- `session_id`
- `environment`
- `preview_types_requested[]`
- `material_diff_hash`
- `base_artifact_refs[]`
- `base_artifact_fingerprints[]`
- `evidence_refs[]`
- `reason_codes[]`
- `staleness_status`
- `correlation_refs[]`
- `deduplication_key`
- `created_at`
- `updated_at`

### Semántica mínima obligatoria

- `proposal_draft_id` y `patchset_candidate_id` deben pertenecer al mismo cambio candidateado.
- `artifacts_in_scope[]` debe ser compatible con el scope material de D.2.
- `human_diff_ref` y `material_diff_ref` son obligatorios y separados.
- `preview_scope` debe indicar qué parte del patchset se está previsualizando (`full_change`, `bounded_subset`, `approval_focus`, etc.), pero nunca puede exceder el scope material real.
- `current_preview_state` debe pertenecer al catálogo v1: `preview_ready`, `preview_running`, `preview_completed`, `preview_stale`, `preview_superseded`, `preview_cancelled`.
- `simulation_requests[]` debe declarar explícitamente qué simulaciones fueron pedidas y con qué inputs esperados.
- Si falta correlación mínima (`tenant_id`, `subject_id`, `trace_id` o equivalente resoluble), el preview no puede promocionar.

## Shape mínimo de `simulation_result`

`simulation_result` debe cubrir como mínimo:

- `simulation_result_id`
- `preview_candidate_id`
- `simulation_family` (`policy|approval|classification|risk`)
- `status` (`preview_ok|preview_warning|preview_blocked|preview_incomplete`)
- `reason_codes[]`
- `inputs_refs[]`
- `outputs_summary`
- `blocking_findings[]`
- `warnings[]`
- `confidence_level`

### Campos normativos adicionales recomendados

- `simulation_request_id`
- `trace_id`
- `tenant_id`
- `proposal_draft_id`
- `patchset_candidate_id`
- `material_diff_hash`
- `preview_type`
- `evidence_refs[]`
- `policy_reference_version`
- `approval_profile_refs[]`
- `classification_policy_refs[]`
- `risk_snapshot`
- `fidelity_level`
- `staleness_status`
- `generated_at`
- `deduplication_key`

### Semántica mínima obligatoria

- `simulation_family` pertenece al catálogo cerrado v1; no se aceptan familias libres.
- `status` expresa resultado de preview, no outcome de runtime ni decisión final del kernel.
- `inputs_refs[]` debe permitir reconstruir exactamente con qué artifacts e inputs se produjo la simulación.
- `outputs_summary` debe ser legible para humano pero anclado a evidencia y refs, no texto suelto.
- `blocking_findings[]` contiene findings que impiden promoción; `warnings[]` contiene findings no bloqueantes.
- `confidence_level` debe explicitar al menos `high|medium|low`.
- Si `outputs_summary` no puede sostenerse con evidencia mínima, el resultado debe ser `preview_incomplete`.

## Tipos de preview mínimos

La surface debe soportar como mínimo estos tipos:

1. `preview_summary`
   - lectura ejecutiva del cambio candidateado;
   - resume estado compuesto, scope, evidencias y findings clave.

2. `policy_preview`
   - anticipa cómo el kernel/PDP podría evaluar autorización/governance contextual sobre el cambio candidateado;
   - nunca sustituye enforcement real.

3. `approval_preview`
   - anticipa approval mode, autoridad esperable, SoD y potenciales bloqueos de aprobación;
   - nunca constituye aprobación real.

4. `classification_preview`
   - anticipa visibilidad, redacción y restricciones de evidencia/diff/output dentro de la surface;
   - nunca sustituye clasificación final productiva.

5. `risk_preview`
   - anticipa riesgo e impacto esperado usando el baseline de A.4 risk model y señales del patchset;
   - no redefine el motor de riesgo.

6. `composite_preview`
   - agrega el estado compuesto de los sub-previews;
   - debe permanecer alineado con los sub-resultados y no puede contradicirlos.

## Simulación de policy

La simulación de policy debe anticipar, sobre el `governed_patchset_candidate`, cómo se construiría el input de evaluación relevante para policy sin mandar texto libre a Cerbos.

### Debe cubrir como mínimo

- `decision_point` simulado relevante (`compile`, `execute`, `release_execution`, `release_application` o `view_restricted`, según el caso);
- principal/subject/tenant/context relevantes;
- `resource_kind` y `action` compatibles con C.3;
- clasificación, risk level, approval mode efectivo y external effect inferibles del candidate;
- posibles salidas materiales: `allow`, `deny_block`, `require_approval`, `require_escalation`, `restricted_view`.

### Reglas normativas

- La simulación de policy debe apoyarse en shapes canonizados compatibles con `policy_input_v1` de C.3.
- Si falta información mínima para construir una proyección seria del input canonizado, el resultado es `preview_incomplete`.
- Si la simulación detecta conflicto de governance evidente o no mapeable, el resultado es `preview_blocked`.
- `allow` simulado no habilita apply por sí solo; sólo elimina un posible bloqueo de esa familia.
- Debe persistirse referencia a policy version o snapshot usada por la simulación cuando exista.

## Simulación de approval

La simulación de approval debe anticipar qué approval mode efectivo requeriría el cambio candidateado y qué bloqueos previsibles aparecerían antes de apply.

### Debe cubrir como mínimo

- `approval_mode` inferido (`auto`, `pre_execution`, `pre_application`, `double`);
- razones de elevación por risk, clasificación, external effect, scope o governance sensitivity;
- autoridad/approver types esperables;
- constraints de SoD relevantes;
- potencial invalidez por cambio material si el preview queda stale;
- si el cambio candidateado es resoluble con approval o si está estructuralmente bloqueado.

### Reglas normativas

- Debe respetar floors y overrides duros de A.4.
- Si no puede inferirse un approval mode serio, el resultado debe ser `preview_incomplete`.
- Si la autoridad requerida es incompatible o no resoluble bajo governance conocida, el resultado debe ser `preview_blocked`.
- `approval_preview` nunca equivale a `approval_request` ni a `approval_decision`.
- Si la policy preview requiere approval pero approval preview no puede resolverse, la surface compuesta debe quedar bloqueada.

## Simulación de clasificación/redacción

La simulación de clasificación/redacción debe anticipar qué partes del preview y de la evidencia quedarían visibles, resumidas o redaccionadas dentro de la surface.

### Debe cubrir como mínimo

- clasificación máxima inferida del cambio candidateado;
- visibilidad de `human_diff`, `material_diff`, findings y evidence refs;
- campos o secciones que quedarían redaccionados;
- evidencia mínima todavía visible para operar con seguridad;
- incompatibilidades entre lo que el preview dice y lo que realmente podría mostrarse.

### Reglas normativas

- Debe alinearse con el principio de C.4: clasificación y redacción antes de salir.
- Si el sistema no tiene metadata suficiente de clasificación, el resultado debe ser `preview_incomplete`.
- Si la clasificación bloquea evidencia crítica para entender el cambio o sostener gates, el resultado debe ser `preview_blocked`.
- Si se predice redacción parcial, el `human_diff` no puede presentarse como completo; debe reflejar esa degradación.
- La simulación puede producir un preview resumido o parcial, pero nunca inventar visibilidad que el sistema no tendría.

## Simulación de riesgo/impacto

La simulación de riesgo/impacto debe anticipar el impacto probable del cambio candidateado usando el baseline de A.4 risk model, sin crear un motor paralelo y contradictorio.

### Debe cubrir como mínimo

- `business_risk_score` y/o sus señales suficientes cuando existan;
- `security_risk_score` y/o sus señales suficientes cuando existan;
- `effective_risk_level` esperado (`low|medium|high|critical`);
- `external_effect`, blast radius, classification sensitivity, policy sensitivity y rollback feasibility inferibles;
- impacto sobre scope, sistemas, approvals y operator recovery esperable.

### Reglas normativas

- Debe usar los buckets, floors duros y semántica base de A.4 risk model.
- Si existen signals parciales pero no suficientes para un cálculo serio, el resultado debe ser `preview_incomplete`.
- Si el riesgo esperado eleva approvals, classification o governance constraints, eso debe reflejarse en el compuesto.
- `risk_preview` alto no puede convivir con un `preview_summary` que declare low risk sin explicación explícita.

## Diff humano vs diff material dentro del preview

### `human_diff` dentro del preview

Es la lectura explicativa para operador, reviewer o developer. Debe responder:

- qué se propone cambiar;
- por qué;
- qué impacto se comunica en términos humanos;
- qué dudas, warnings o degradaciones siguen abiertas;
- qué partes están resumidas o redaccionadas.

### `material_diff` dentro del preview

Es la lectura estructurada y determinística del cambio candidateado. Debe responder:

- qué artifacts concretos cambian;
- qué operaciones materiales se proponen;
- qué nodos/campos/reglas/bindings se agregan, editan o eliminan;
- qué hash/fingerprint del patchset sostiene la preview.

### Relación normativa

- Ambos diffs deben seguir alineados sobre el mismo `material_diff_hash`.
- El `human_diff` puede resumir; el `material_diff` no puede omitir materialidad crítica.
- Si el `human_diff` y el `material_diff` divergen materialmente, el preview compuesto no puede quedar `preview_ok`.
- Si classification preview anticipa redacción parcial, esa degradación debe reflejarse también en el `human_diff` y en el summary.
- El `material_diff` es la base de simulación; el `human_diff` es la base de lectura humana.

## Gates de promoción desde preview hacia apply candidate

Para promover desde preview hacia `apply_candidate`, deben cumplirse como mínimo estos gates:

1. existir `preview_candidate` persistido, vigente y no stale;
2. el `preview_candidate` debe apuntar a un `proposal_draft` no cancelado ni rechazado;
3. el `patchset_candidate_id` debe corresponder al mismo `material_diff_hash` usado por todos los sub-previews relevantes;
4. existir `human_diff_ref` y `material_diff_ref` resolubles;
5. existir `policy_preview`, `approval_preview`, `classification_preview` y `risk_preview` para el scope requerido;
6. ningún sub-preview obligatorio puede quedar en `preview_blocked`;
7. ningún sub-preview obligatorio puede quedar en `preview_incomplete`;
8. el `composite_preview` debe quedar explícitamente en `preview_ok`;
9. reason codes y evidence refs deben estar completos;
10. no debe existir contradicción material entre sub-previews;
11. no debe existir desalineación entre diff humano, diff material y outputs resumidos;
12. no debe existir `staleness_status = stale` en preview ni en resultados usados para componerlo;
13. si el preview compuesto es `preview_warning`, la promoción a `apply_candidate` NO está permitida en v1;
14. debe existir promoción explícita con actor, timestamp y reason code.

### Regla normativa fuerte

`preview_candidate` NO se promociona a `apply_candidate` salvo que el compuesto esté en `preview_ok`. `preview_warning`, `preview_blocked` y `preview_incomplete` no alcanzan.

## Evidencia mínima de preview/simulación

La preview surface debe poder persistir y consultar como mínimo esta evidencia:

- `preview_candidate_id`, versión y estado.
- `proposal_draft_id`, `patchset_candidate_id`, `tenant_id`, `subject_id`, `trace_id`.
- `proposal_version`, `patchset_version` y `material_diff_hash`.
- `human_diff_ref` y `material_diff_ref`.
- `artifacts_in_scope[]` y `preview_scope`.
- `simulation_requests[]` emitidas.
- `simulation_result_id` por familia.
- `inputs_refs[]` usados por cada simulación.
- `outputs_summary` con references/evidence refs suficientes.
- `blocking_findings[]`, `warnings[]`, `reason_codes[]` y `confidence_level`.
- snapshots/versiones de policy, approvals, clasificación y riesgo usados por la lectura simulada cuando existan.
- promotion events, stale marks, supersessions y reason codes correspondientes.

Sin este trail, el preview no puede considerarse gobernado ni promocionable.

## Reason codes mínimos de preview/simulación

La preview surface debe contemplar al menos estos reason codes:

- `preview.ready`
- `preview.warning.policy_sensitive`
- `preview.warning.classification_partial`
- `preview.warning.risk_high`
- `preview.blocked.missing_material_diff`
- `preview.blocked.governance_conflict`
- `preview.blocked.approval_unresolvable`
- `preview.blocked.classification_blocking`
- `preview.incomplete.insufficient_inputs`
- `preview.incomplete.stale_patchset`

### Reason codes adicionales recomendados

- `preview.blocked.proposal_terminal_state`
- `preview.blocked.subpreview_conflict`
- `preview.blocked.diff_misaligned`
- `preview.blocked.scope_out_of_bounds`
- `preview.warning.redaction_hides_noncritical_detail`
- `preview.warning.approval_elevation_expected`
- `preview.warning.risk_score_provisional`
- `preview.incomplete.missing_trace`
- `preview.incomplete.policy_input_unresolved`
- `preview.incomplete.outputs_without_evidence`

## Integración con D.2 proposal workspace

- D.3 consume `proposal_draft` y `governed_patchset_candidate` emitidos por D.2; no reemplaza el workspace.
- `preview_candidate` sólo puede abrirse desde artifacts ya gobernados, nunca desde chat libre o `intent_input` sin pasar por D.2.
- D.2 sigue siendo dueño de la preparación de `human_diff`, `material_diff`, scope material y patchset candidateado.
- D.3 puede devolver findings que hagan volver el draft a `draft_refining`, `preview_blocked` o stale.
- Si cambia el `material_diff_hash`, el preview queda stale y debe recomponerse; no puede reciclarse por inercia.
- D.2 mantiene el gate de `ready_for_preview`; D.3 agrega el gate de `preview_ok` compuesto antes de `apply_candidate`.

## Integración con C.3/C.4 y artifacts del kernel

- `policy_preview` debe ser compatible con el boundary PEP/PDP de C.3 y con `policy_input_v1`; no manda texto libre al PDP.
- `approval_preview` debe respetar el subsystem de approvals de A.4 sin fingir requests o decisions reales.
- `classification_preview` debe respetar redacción previa a salida según C.4.
- `risk_preview` debe reutilizar la semántica de A.4 risk model y no crear un scoring paralelo incompatible.
- La preview surface debe emitir evidencia canónica o equivalente normativo correlable con `tenant_id`, `trace_id`, `proposal_draft_id`, `patchset_candidate_id` y `preview_candidate_id`.
- Observabilidad derivada puede resumir o referenciar preview/simulación, pero nunca reemplaza la evidencia canónica ni sus reason codes.
- D.3 no crea `compiled_contract`, `execution_record`, `policy_decision_record` ni `approval_request` reales; sólo se integra con sus artifacts, snapshots o inputs compatibles.

## Idempotencia y deduplicación de preview/simulación

- El sistema debe derivar una `deduplication_key` a partir de `tenant_id`, `proposal_draft_id`, `patchset_candidate_id`, `material_diff_hash`, `preview_scope` y familia/tipo de simulación.
- Retries técnicos no deben duplicar `preview_candidate` ni `simulation_result` equivalentes.
- Si los inputs materiales no cambian, repetir preview debe reutilizar o superseder limpiamente resultados equivalentes, no crear ruido canónico duplicado.
- Si cambia el `material_diff_hash`, debe abrirse nueva versión o nueva corrida de preview; nunca reusar `simulation_result` viejo como si fuera vigente.
- Si cambian base artifact fingerprints, policy version, classification metadata o approval constraints materiales, el preview debe marcarse stale.
- El `composite_preview` sólo puede componerse con sub-resultados vigentes para el mismo `preview_candidate_id` y el mismo `material_diff_hash`.

## Tests borde mínimos (al menos 20)

1. preview sobre proposal sin patchset candidate.
2. preview con diff humano presente pero diff material ausente.
3. patchset stale respecto a artifacts base.
4. policy preview con inputs incompletos.
5. approval preview sin approval mode inferible.
6. classification preview sin metadata de clasificación suficiente.
7. risk preview con artifacts fuera de scope.
8. preview dice ok pero un sub-preview está blocked.
9. sub-previews contradictorios entre sí.
10. preview generado sobre proposal cancelada.
11. preview sobre proposal rechazada.
12. preview repetido no duplica `simulation_result` equivalente.
13. preview candidate sin trace/correlation suficiente.
14. classification preview predice redacción parcial pero human diff no la refleja.
15. policy preview allow pero approval preview bloquea.
16. risk preview alto pero proposal sigue marcado low risk en resumen.
17. preview candidate promovido a apply candidate sin `preview_ok` compuesto.
18. composite preview usa resultados viejos con patchset nuevo.
19. preview incompleto tratado como warning.
20. `outputs_summary` sin evidence refs.
21. preview candidate con `artifacts_in_scope[]` incompatibles con el patchset.
22. policy preview construido desde summary humano en vez de material diff.
23. approval preview ignora floor duro por `restricted` o `irreversible`.
24. classification preview permite visibilidad completa pese a evidence restricted.
25. risk preview omite `external_effect` irreversible y devuelve medium.
26. composite preview queda `preview_ok` pese a `preview_warning` de policy-sensitive change usado como bloqueo en governance local.
27. preview candidate queda vigente aunque cambió `policy_version` material.
28. sub-preview de classification refiere a otro `preview_candidate_id`.
29. policy preview marca `require_approval` pero approval preview informa `auto` sin explicación.
30. `preview_summary` declara “listo para apply” sin reason code ni promotion event.

## Criterios de aceptación de D.3

- Existe boundary explícito y no ambiguo entre `proposal_draft`, `preview_candidate`, `simulation_result` y `apply_candidate`.
- Queda explícito que preview y simulación corren sobre `proposal_draft`/`governed_patchset_candidate`, nunca sobre charla informal.
- `preview_candidate` queda definido como artifact de surface distinto de `apply_candidate`.
- `simulation_result` queda definido como evidencia previa y explícitamente distinto del enforcement real del kernel.
- Quedan definidas las simulaciones mínimas de `policy`, `approval`, `classification` y `risk`.
- Quedan definidos los tipos mínimos de preview y el resultado explícito `preview_ok | preview_warning | preview_blocked | preview_incomplete`.
- Queda definida la coexistencia obligatoria de diff humano y diff material dentro del preview, con alineación material obligatoria.
- Quedan definidos los gates de promoción desde preview hacia `apply_candidate`, incluyendo la exigencia de `preview_ok` compuesto.
- Queda explícita la evidencia mínima requerida para cada simulación, incluyendo `inputs_refs[]` y soporte de `outputs_summary`.
- Quedan definidos reason codes mínimos, integración con D.2/C.3/C.4/C.6, reglas de idempotencia/deduplicación y tests borde mínimos de la surface.
