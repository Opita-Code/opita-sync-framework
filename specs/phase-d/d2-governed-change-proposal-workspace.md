# D.2 — Governed change proposal workspace

## Objetivo

Definir la surface v1 donde un `intent_input` ya emitido por D.1, o una salida equivalente de intake compatible con el pipeline conversacional elegido, se transforma en un **`proposal_draft` gobernado, auditable, refinable y promocionable** hasta producir un **`governed_patchset_candidate` todavía no aplicado**, sin colapsar proposal, preview ni apply en el mismo acto y sin contradecir D.1, D.3, B.6, C.1 ni C.6.

El objetivo de D.2 NO es ejecutar cambios, NO es correr preview final, NO es aprobar, NO es compilar contratos ni disparar runtime. El objetivo es fijar el contrato exacto del workspace donde el cambio propuesto deja de ser texto conversacional y pasa a ser un artifact gobernado de surface listo para revisión, preview y promoción explícita.

## Principios de implementación del proposal workspace

- El proposal workspace es una **surface gobernada de shaping y promoción**, no un apply shell.
- `proposal_draft` es un **objeto/artefacto gobernado de surface**, no texto libre ni resumen informal.
- `proposal_draft` NO es `apply_candidate`.
- Nada se promociona por “quedó bastante claro”; toda promoción exige estado, gates y evidencia explícita.
- El workspace debe sostener simultáneamente dos lecturas del cambio: `human_diff` y `material_diff`.
- El `material_diff` es la base normativa para preview, policy, approval, classification y apply posterior.
- El workspace puede volver hacia intake o hacia revisión, pero NO puede saltear directamente a apply.
- El workspace puede emitir `governed_patchset_candidate` como artifact todavía no aplicado.
- Apply real sigue fuera de D.2 y de la verdad del workspace.
- El workspace debe operar fail-closed ante pérdida de scope material, evidencia insuficiente, contradicción entre diffs o intento de bypass de governance.
- La conversación puede originar o refinar la propuesta, pero la verdad operativa del workspace vive en artifacts persistidos, versionables y correlables.

## Boundary exacto entre `intent_input`, `proposal_draft`, `preview_candidate`, `apply_candidate`

### 1. `intent_input`

Es la entrada gobernada que ya representa la intención estructurada emitida por D.1 y compatible con C.1. Describe **qué se quiere lograr** con suficiente shape de origen, alcance y evidencia, pero todavía NO describe necesariamente el cambio material completo como proposal workspace auditable.

### 2. `proposal_draft`

Es el artifact gobernado de surface donde la intención se convierte en **propuesta revisable**. Debe explicitar scope, cambios propuestos, constraints, supuestos, preguntas abiertas, evidencia de origen y dos referencias de diff (`human_diff_ref`, `material_diff_ref`). Puede estar incompleto o en revisión, pero ya no es texto libre.

### 3. `preview_candidate`

Es la representación lista para D.3 que toma un `proposal_draft` suficientemente cerrado y un `governed_patchset_candidate` consistente para exponer simulación, diff operativo, lectura de policy, approval y classification. `preview_candidate` NO se produce desde chat libre ni desde `intent_input` sin pasar por proposal workspace.

### 4. `governed_patchset_candidate`

Es el artifact material derivado del `proposal_draft` que enumera operaciones gobernadas todavía no aplicadas. Representa el cambio material candidateado que podrá sostener preview, policy, approval, classification y eventual apply. No implica ejecución ni mutación efectiva.

### 5. `apply_candidate`

Es el artifact posterior al preview y a los gates correspondientes que declara que un cambio material candidateado ya quedó apto para pasar al flujo de apply real. `apply_candidate` sigue fuera de D.2 como acto final; D.2 sólo puede dejar al proposal en `ready_for_apply_candidate` y dejar el `governed_patchset_candidate` preparado para esa promoción posterior.

### Regla normativa de boundary

- `intent_input` = intención gobernada.
- `proposal_draft` = propuesta gobernada revisable.
- `preview_candidate` = paquete listo para simulación/preview.
- `governed_patchset_candidate` = cambio material candidateado todavía no aplicado.
- `apply_candidate` = entrada posterior para apply real.

Nunca deben colapsarse en un solo objeto ni tratarse como equivalentes.

## Responsabilidades explícitas del proposal workspace

- Recibir `intent_input` o handoff de intake compatible con apertura formal de propuesta.
- Crear, persistir y versionar `proposal_draft` con correlación por tenant, sesión, sujeto y fuente.
- Mantener `artifacts_in_scope[]` como scope material revisable.
- Construir y actualizar `proposed_changes[]` como lista explícita de cambios propuestos.
- Preservar `constraints[]`, `assumptions[]` y `open_questions[]` por separado.
- Mantener evidencia mínima de de dónde salió la propuesta: turns, inputs fuente, artifacts base, revisiones y supuestos.
- Emitir y correlacionar `human_diff_ref` y `material_diff_ref`.
- Detectar desalineación entre diff humano y diff material.
- Generar `governed_patchset_candidate` cuando el material diff ya es suficientemente estable.
- Evaluar gates previos a `ready_for_preview` y previos a `ready_for_apply_candidate`.
- Permitir promoción, devolución a revisión, rechazo, cancelación o cierre con evidence trail.
- Marcar stale o bloquear drafts cuando cambian artifacts base o hashes materiales relevantes.
- Aplicar idempotencia y deduplicación para evitar proposals equivalentes duplicadas por retries o turns repetidos.
- Publicar reason codes auditables y correlación hacia event log/observabilidad.

## No-responsabilidades explícitas del proposal workspace

- No ejecutar apply real.
- No crear `execution_record` ni arrancar workflows.
- No reemplazar preview/simulación detallada de D.3.
- No decidir enforcement final de policy.
- No emitir approvals reales ni consumir approvals.
- No realizar clasificación/redacción final.
- No compilar `compiled_contract` ni persistir outputs de C.1.
- No convertir texto libre directamente en `apply_candidate`.
- No tratar el diff humano como fuente material de verdad.
- No ocultar preguntas abiertas materiales bajo un resumen optimista.
- No reabrir distribution layer, tenant activation ni rollout operativo.

## Shape mínimo de `proposal_draft`

`proposal_draft` debe cubrir como mínimo:

- `proposal_draft_id`
- `tenant_id`
- `session_id`
- `subject_id`
- `source_intent_refs[]`
- `title`
- `summary`
- `artifacts_in_scope[]`
- `proposed_changes[]`
- `constraints[]`
- `assumptions[]`
- `open_questions[]`
- `current_state`
- `confidence_level`
- `human_diff_ref`
- `material_diff_ref`

### Campos normativos adicionales recomendados

- `proposal_version`
- `trace_id`
- `workspace_id`
- `created_at`
- `updated_at`
- `created_from_intake_session_id`
- `base_artifact_refs[]`
- `base_artifact_fingerprints[]`
- `review_iteration`
- `reason_codes[]`
- `evidence_refs[]`
- `staleness_status`
- `last_material_diff_hash`
- `preview_readiness_snapshot`
- `apply_candidate_readiness_snapshot`
- `deduplication_key`

### Semántica mínima obligatoria

- `source_intent_refs[]` no puede estar vacío.
- `title` y `summary` son lectura humana gobernada, no texto libre sin estructura de origen.
- `artifacts_in_scope[]` debe representar artifacts afectados o sospechados con enough clarity para revisión; si el scope material sigue incierto, el draft no puede promocionar a preview.
- `proposed_changes[]` debe poder correlacionarse con el `material_diff_ref`.
- `constraints[]` distingue restricciones declaradas, heredadas o inferidas.
- `assumptions[]` debe contener TODA inferencia no confirmada por el usuario o por artifacts base.
- `open_questions[]` no puede omitirse si existe ambigüedad material.
- `human_diff_ref` y `material_diff_ref` son obligatorios como referencias separadas; nunca un solo diff sirve para ambos fines.
- `confidence_level` informa la madurez percibida, pero NO habilita promoción por sí solo.

## Shape mínimo de `governed_patchset_candidate`

`governed_patchset_candidate` debe cubrir como mínimo:

- `patchset_candidate_id`
- `proposal_draft_id`
- `target_artifacts[]`
- `material_operations[]`
- `material_diff_hash`
- `policy_preview_inputs_ref`
- `approval_preview_inputs_ref`
- `classification_preview_inputs_ref`
- `ready_for_preview` bool
- `ready_for_apply_candidate` bool

### Campos normativos adicionales recomendados

- `tenant_id`
- `session_id`
- `subject_id`
- `trace_id`
- `patchset_version`
- `base_artifact_refs[]`
- `base_artifact_fingerprints[]`
- `human_diff_ref`
- `material_diff_ref`
- `evidence_refs[]`
- `reason_codes[]`
- `staleness_status`
- `derived_from_proposal_version`
- `generated_at`

### Semántica mínima obligatoria

- `target_artifacts[]` debe ser compatible con `artifacts_in_scope[]` del draft vigente.
- `material_operations[]` es la enumeración material que luego sostendrá preview/policy/approval/classification/apply.
- `material_diff_hash` debe cambiar cuando cambie materialmente el patchset candidateado.
- `policy_preview_inputs_ref`, `approval_preview_inputs_ref` y `classification_preview_inputs_ref` deben apuntar a insumos consistentes con el mismo `material_diff_hash`.
- `ready_for_preview = true` sólo si los gates previos a preview cerraron.
- `ready_for_apply_candidate = true` sólo si además cerraron los gates previos a apply candidate.

## Lifecycle exacto del proposal draft

1. **Apertura**
   - D.1 emite handoff válido y se crea `proposal_draft` en `draft_open`.
   - Se registran `source_intent_refs[]`, evidencia inicial, artifacts base conocidos y primer resumen gobernado.

2. **Refinamiento**
   - El draft entra en `draft_refining` mientras se clarifican scope, artifacts, constraints, assumptions y proposed changes.
   - Pueden agregarse turns, revisiones humanas o hallazgos sobre artifacts base.

3. **Solicitud de revisión**
   - Si aparece ambigüedad material, contradicción, evidencia insuficiente o attachments no parseables, el draft pasa a `awaiting_revision`.
   - Desde ahí sólo puede volver a refinamiento, ser rechazado o cancelado.

4. **Preparación para preview**
   - Cuando existe `human_diff_ref`, `material_diff_ref`, artifacts en scope estables y evidence trail mínimo, el workspace produce `governed_patchset_candidate` y puede promocionar a `ready_for_preview`.

5. **Bloqueo de preview**
   - Si el patchset candidateado detecta sensibilidad de policy, classification, falta de evidencia o inconsistencia material, el draft puede pasar a `preview_blocked`.
   - Este estado no aplica cambios; sólo frena la promoción y exige revisión/gates adicionales.

6. **Preparación para apply candidate**
   - Después de la instancia de preview/simulación correspondiente en D.3 y de cerrar las observaciones requeridas, el draft puede quedar en `ready_for_apply_candidate`.
   - Ese estado declara preparación de surface, no apply real.

7. **Cierre terminal**
   - `rejected` si el cambio entra en conflicto con governance o resulta inviable.
   - `cancelled` si se aborta por pedido del usuario o retiro explícito.
   - `closed` si el workspace ya dejó evidencia íntegra y el draft terminó su función dentro de la cadena de proposal.

## Estados mínimos y transiciones válidas/prohibidas

### Estados mínimos

- `draft_open`
- `draft_refining`
- `awaiting_revision`
- `ready_for_preview`
- `preview_blocked`
- `ready_for_apply_candidate`
- `rejected`
- `cancelled`
- `closed`

### Transiciones válidas

- `draft_open -> draft_refining`
- `draft_open -> awaiting_revision`
- `draft_open -> cancelled`
- `draft_refining -> awaiting_revision`
- `draft_refining -> ready_for_preview`
- `draft_refining -> rejected`
- `draft_refining -> cancelled`
- `awaiting_revision -> draft_refining`
- `awaiting_revision -> rejected`
- `awaiting_revision -> cancelled`
- `ready_for_preview -> preview_blocked`
- `ready_for_preview -> draft_refining` cuando el preview detecta inconsistencias o el base cambia
- `ready_for_preview -> ready_for_apply_candidate` sólo con evidencia de preview suficiente y patchset consistente
- `preview_blocked -> draft_refining`
- `preview_blocked -> rejected`
- `preview_blocked -> cancelled`
- `ready_for_apply_candidate -> closed`
- `rejected -> closed`
- `cancelled -> closed`

### Transiciones prohibidas

- `draft_open -> ready_for_apply_candidate`
- `draft_open -> closed` sin evidence trail mínimo
- `draft_refining -> ready_for_apply_candidate` sin pasar por `ready_for_preview`
- `awaiting_revision -> ready_for_apply_candidate`
- `preview_blocked -> ready_for_apply_candidate` sin volver a refinamiento o nuevo cierre de gates
- `rejected -> ready_for_preview`
- `cancelled -> ready_for_preview`
- `closed -> draft_refining` sin nueva versión explícita
- cualquier estado -> `apply_candidate` directo dentro de D.2

## Diff humano vs diff material

### `human_diff`

Es la lectura explicativa del cambio para operador, reviewer o developer. Debe responder:

- qué se quiere cambiar;
- por qué se propone;
- qué artifacts o zonas quedarían afectadas;
- qué impacto esperado se comunica en términos humanos;
- qué dudas, riesgos o supuestos siguen abiertos.

Puede expresarse como narrativa estructurada, resumen comparativo, bullets o vista PR-style, pero siempre gobernada y referenciada.

### `material_diff`

Es la lectura precisa, estructurada y determinística del cambio material. Debe responder:

- qué artifacts concretos cambian;
- qué operaciones materiales se proponen;
- qué campos, bindings, reglas o nodos se agregan, editan, eliminan o reordenan materialmente;
- cuál es el hash o fingerprint material del candidate.

### Relación normativa

- Ambos diffs deben hablar del MISMO cambio candidateado.
- El `human_diff` puede resumir; el `material_diff` no puede resumir materialidad crítica.
- El `material_diff` es la base normativa para preview, policy, approval, classification y apply posterior.
- Si `human_diff` y `material_diff` divergen materialmente, el draft debe bloquear promoción.
- El workspace debe poder demostrar trazabilidad entre `summary`/`proposed_changes[]` y `material_operations[]`.

## Gates previos a preview

Para promover un draft a `ready_for_preview`, deben cumplirse como mínimo estos gates:

1. `proposal_draft` persistido y con versión vigente.
2. `source_intent_refs[]` presentes y resolubles.
3. `artifacts_in_scope[]` no vacío.
4. `proposed_changes[]` no vacío.
5. `human_diff_ref` presente y legible.
6. `material_diff_ref` presente y resoluble.
7. `governed_patchset_candidate` derivado y consistente con el draft.
8. `material_diff_hash` calculado.
9. `assumptions[]` y `open_questions[]` no contienen ambigüedad crítica abierta.
10. evidence trail mínimo completo.
11. artifacts base no stale respecto del patchset candidateado.
12. reason code explícito de promoción.

Si cualquiera falla, el draft queda en `draft_refining`, `awaiting_revision` o `preview_blocked` según corresponda.

## Gates previos a apply candidate

Para promover un draft a `ready_for_apply_candidate`, además de los gates previos a preview, deben cumplirse como mínimo:

1. existir `governed_patchset_candidate.ready_for_preview = true`.
2. existir evidencia de preview/simulación suficiente proveniente de D.3.
3. `policy_preview_inputs_ref`, `approval_preview_inputs_ref` y `classification_preview_inputs_ref` presentes para el mismo `material_diff_hash`.
4. no existir `assumptions[]` materiales sin resolver que alteren el patchset.
5. no existir `open_questions[]` críticas.
6. no existir inconsistencia entre summary humano y `material_operations[]`.
7. no existir señal de `staleness_status = stale`.
8. si el cambio es policy-sensitive, existir evidencia reforzada suficiente.
9. si el cambio es classification-sensitive, existir preview específico correspondiente.
10. existir promoción explícita con reason code de preparación para apply candidate.

Si cualquiera falla, el draft NO puede quedar `ready_for_apply_candidate`.

## Relación con `Intent -> Change Proposal -> Governed Patchset`

D.2 hace operable la etapa central del pipeline ya cerrado en B.6:

- **Intent**: llega desde D.1 como `intent_input` o handoff equivalente gobernado.
- **Change Proposal**: vive en D.2 como `proposal_draft` versionado, revisable y gateable.
- **Governed Patchset**: se deriva en D.2 como `governed_patchset_candidate`, todavía no aplicado.

Decisión normativa: preview, policy, approval, classification y eventual apply deben apoyarse en el **patchset material candidateado**, no en la conversación ni en el resumen textual del draft.

## Evidencia mínima del proposal workspace

El workspace debe poder persistir y consultar como mínimo esta evidencia:

- `proposal_draft_id` y versión.
- `tenant_id`, `session_id`, `subject_id`, `trace_id`.
- `source_intent_refs[]`.
- turns fuente relevantes o referencias a ellos.
- artifacts base usados para construir el draft.
- fingerprints o hashes base relevantes.
- `human_diff_ref`.
- `material_diff_ref`.
- `material_diff_hash` del patchset candidateado.
- `assumptions[]` y `open_questions[]` vigentes por versión.
- revisiones, comentarios estructurados o decisiones de shaping.
- promotions de estado con timestamp, actor y reason codes.
- evidencia de bloqueo, rechazo o cancelación cuando ocurra.

Sin este trail, el workspace no puede considerarse gobernado.

## Reason codes mínimos del workspace

El workspace debe contemplar al menos estos reason codes:

- `proposal.continue_refining`
- `proposal.ask_revision.diff_incomplete`
- `proposal.ask_revision.artifact_scope_unclear`
- `proposal.ask_revision.material_diff_missing`
- `proposal.ready_for_preview`
- `proposal.preview_blocked.policy_sensitive`
- `proposal.preview_blocked.classification_sensitive`
- `proposal.ready_for_apply_candidate`
- `proposal.rejected.governance_conflict`
- `proposal.cancelled.user_request`

### Reason codes adicionales recomendados

- `proposal.ask_revision.attachments_unparseable`
- `proposal.ask_revision.assumptions_unresolved`
- `proposal.ask_revision.base_artifact_stale`
- `proposal.preview_blocked.evidence_insufficient`
- `proposal.preview_blocked.diff_inconsistent`
- `proposal.rejected.intent_conflict`
- `proposal.rejected.scope_not_governable`
- `proposal.closed.surface_complete`

## Integración con D.1 intake

- D.1 sigue siendo el único boundary autorizado entre chat libre e artifacts gobernados de entrada.
- D.2 recibe `source_intent_refs[]` emitidos por D.1; no reinterpreta conversación libre como si fuera source of truth.
- Si el workspace detecta scope incierto, source refs vacíos, evidencia insuficiente o conflicto entre turns, debe devolver el caso a intake/revisión y no seguir promocionando.
- Un `proposal_draft` puede abrirse desde D.1 aun con supuestos tolerables, pero esos supuestos deben quedar explícitos en `assumptions[]`.
- D.2 no corrige silenciosamente errores de shaping de D.1; los evidencía y devuelve reason codes cuando haga falta.

## Integración con D.3 preview/simulación

- D.2 prepara la entrada a D.3, pero NO reemplaza su preview surface.
- `ready_for_preview` sólo significa que el draft ya puede entrar en preview/simulación; no significa que el preview exista o esté aprobado.
- D.3 debe consumir el `governed_patchset_candidate`, no el chat libre ni sólo el summary humano.
- El output de D.3 puede devolver observaciones, bloqueos o señales de sensibilidad que hagan volver el draft a `draft_refining` o `preview_blocked`.
- `ready_for_apply_candidate` sólo puede declararse después de considerar la evidencia relevante de D.3.

## Integración con C.1 compiler (cuando corresponda)

- D.2 NO invoca el compilador como responsabilidad propia.
- Cuando el flujo posterior necesite entrada compatible con C.1, el workspace debe proveer `source_ref` gobernado y evidencia suficiente para derivar un `intent_input` o handoff equivalente SIN reinterpretación libre posterior.
- El `material_diff` y el `governed_patchset_candidate` pueden informar el handoff al compilador, pero no sustituyen `compiled_contract` ni `compilation_report`.
- Ningún estado de D.2 implica `compiled`, `executable` ni arranque de runtime; esos seams siguen cerrados en C.1-C.6.

## Integración con event log y observabilidad

- El workspace debe emitir eventos canónicos o equivalentes normativos para: apertura de draft, revisiones materiales, generación de patchset candidateado, promociones, bloqueos, rechazo, cancelación y cierre.
- Debe preservarse correlación por `tenant_id`, `session_id`, `subject_id`, `trace_id`, `proposal_draft_id` y `patchset_candidate_id`.
- Observabilidad derivada puede exponer métricas o spans, pero nunca reemplaza la evidencia canónica del workspace.
- Un draft visible como `ready_for_preview` o `ready_for_apply_candidate` debe poder trazarse a su promotion event y a su reason code correspondiente.

## Idempotencia y deduplicación del workspace

- El workspace debe derivar una `deduplication_key` a partir de tenant, sesión, source intents, scope material y hash material vigente.
- Retries técnicos no deben duplicar `proposal_draft` si el mismo cambio material ya está abierto.
- Turns duplicados o replay de intake no deben crear dos drafts equivalentes sin señal explícita de nueva versión.
- Si cambia el `material_diff_hash`, debe abrirse nueva versión o revisión explícita; no puede mantenerse el mismo readiness state por inercia.
- Si cambia un artifact base relevante, el draft debe marcarse stale o volver a refinamiento.
- Reabrir un draft `cancelled` o `closed` exige nueva versión o nuevo draft correlacionado; nunca mutación silenciosa del historial.

## Tests borde mínimos (al menos 20)

1. proposal draft sin artifacts en scope.
2. proposal draft sin diff material.
3. diff humano y diff material inconsistentes.
4. múltiples artifacts con cambios contradictorios.
5. proposal generado desde intent ambiguo.
6. `source_intent_refs` vacíos.
7. `material_diff_hash` cambia sin actualizar proposal state.
8. proposal intenta saltar a `apply_candidate` sin preview.
9. `preview_candidate` generado sin `governed_patchset_candidate`.
10. proposal con assumptions no resueltas pero marcado `ready_for_apply_candidate`.
11. policy-sensitive change sin evidence suficiente.
12. classification-sensitive change sin preview específico.
13. proposal cancelado reabierto sin nueva versión.
14. retry del workspace duplica `proposal_draft`.
15. turns duplicados generan dos proposals equivalentes.
16. artifact base cambió y proposal queda stale.
17. proposal con attachments no parseables.
18. patchset candidate inconsistente con human summary.
19. proposal rechazado pero sigue visible como `ready_for_preview`.
20. proposal cerrada sin evidence trail.
21. `governed_patchset_candidate.target_artifacts[]` no coincide con `artifacts_in_scope[]`.
22. `material_operations[]` vacío con `ready_for_preview = true`.
23. promoción a `ready_for_preview` sin reason code.
24. promoción a `ready_for_apply_candidate` con `preview_blocked` vigente.
25. draft con `open_questions[]` críticas omitidas por error de shaping.
26. base fingerprints cambian entre draft y patchset candidateado.
27. `human_diff_ref` existe pero apunta a versión previa del draft.
28. `policy_preview_inputs_ref` corresponde a otro `material_diff_hash`.
29. rechazo por governance conflict sin persistir evidencia de causal.
30. draft `closed` sin correlación a promotion o terminal event.

## Criterios de aceptación de D.2

- Existe boundary explícito y no ambiguo entre `intent_input`, `proposal_draft`, `preview_candidate`, `governed_patchset_candidate` y `apply_candidate`.
- `proposal_draft` queda definido como artifact gobernado de surface y explícitamente distinto de `apply_candidate`.
- `governed_patchset_candidate` queda definido como artifact material todavía no aplicado.
- El workspace sostiene simultáneamente `human_diff` y `material_diff`, con prioridad normativa del segundo para preview/policy/approval/classification/apply.
- El lifecycle del draft queda cerrado con estados mínimos, transiciones válidas y transiciones prohibidas.
- Quedan definidos los gates previos a preview y previos a apply candidate.
- Queda explícita la evidencia mínima del workspace y los reason codes mínimos.
- Queda explícito que el workspace puede volver a intake o revisión, pero no puede saltear directo a apply.
- Queda explícita la integración con D.1, D.3, C.1, event log y observabilidad sin reabrir seams del kernel.
- Queda documentada idempotencia/deduplicación y el set mínimo de tests borde para validar la surface.
