# D.6 — Phase D integration checkpoint

## Objetivo

D.6 cierra el checkpoint de integración de **Fase D** para demostrar que la **surface conversacional y operativa v1** ya puede recorrer un corredor mínimo end-to-end sobre el kernel cerrado en C.6, sin reabrir seams del motor, sin colapsar boundaries definidos en D.1-D.5 y sin extender alcance a distribution layer.

Este checkpoint NO agrega breadth funcional nuevo. Valida que intake, shaping, proposal, preview, inspection/debugging y evidencia ya nacen compatibles entre sí a nivel de **surface v1** sobre un kernel ya cerrado.

## Qué valida exactamente este checkpoint

Este checkpoint valida la **surface**, no el kernel ni el distribution layer.

Debe poder probarse, como mínimo, que la surface:

1. recibe y persiste un `conversation_turn` dentro de un `intake_session` trazable;
2. produce `intent_candidate` y/o `intent_input` sin colapsar conversación libre con artifacts gobernados;
3. produce un `proposal_draft` gobernado y correlado con el intake;
4. produce un `preview_candidate` consistente con el `proposal_draft` y su patchset material;
5. produce al menos un `simulation_result` con `inputs_refs[]` y evidencia suficiente;
6. expone una `execution_inspection_view` o una `semantic_debug_view` coherente con artifacts y records reales del kernel;
7. demuestra correlación completa entre artifacts de surface y evidencia canónica del kernel;
8. puede cerrar el smoke path mínimo sin `apply` real, siempre que el corredor mínimo quede íntegro, evidenciado y gateado.

D.6 NO valida:

- recompilación del kernel como seam nuevo;
- distribution layer, rollout, activación operacional de tenants o consumo downstream;
- apply real como requisito para aprobar el checkpoint;
- UX final, dashboards maduros o automatización avanzada de Fase E;
- bypass de governance, approvals o boundaries ya cerrados.

## Principios del checkpoint de integración de surface

- **Surface-only scope.** D.6 cubre intake, proposal, preview, inspection y debugging asistido; no reabre seams internos del kernel.
- **Kernel cerrado, surface integrada.** C.6 ya cerró el corredor engine-only; D.6 valida cómo la surface se apoya sobre ese kernel sin reinterpretarlo.
- **Artifacts primero.** La verdad operativa de la surface vive en artifacts persistidos y correlados, no en resúmenes conversacionales ni UI state efímero.
- **Boundaries explícitos.** `conversation_turn`, `intake_session`, `intent_candidate`, `intent_input`, `proposal_draft`, `preview_candidate`, `simulation_result`, `execution_inspection_view` y `semantic_debug_view` siguen siendo objetos distintos.
- **Fail-closed por defecto.** Si la surface salta governance, saltea boundaries o explica sin evidence trail suficiente, el checkpoint no puede pasar.
- **Correlación íntegra.** La surface sólo puede declararse integrada si puede enlazar sus artifacts con records canónicos del kernel.
- **Observabilidad derivada desacoplada.** Si falla observabilidad derivada, el checkpoint puede pasar sólo si artifacts y evidencia canónica siguen íntegros.
- **Integración mínima, no completitud del producto.** Vale más cerrar el corredor mínimo end-to-end con evidencia robusta que abrir breadth adicional.

## Supuestos heredados de D.1-D.5 y C.6

D.6 hereda sin reabrir estas decisiones:

- D.1 fija el boundary entre chat libre, shaping, `intent_input` y `proposal_draft` sin permitir ejecución directa desde conversación.
- D.1 fija `conversation_turn`, `intake_session` e `intent_candidate` con evidencia mínima, correlación mínima y fail-closed ante ambigüedad crítica.
- D.2 fija `proposal_draft` como artifact gobernado y `governed_patchset_candidate` como cambio material candidateado todavía no aplicado.
- D.2 fija que `proposal_draft` NO es `apply_candidate` y que no existe salto directo desde intake/chat a apply.
- D.3 fija `preview_candidate` y `simulation_result` como artifacts separados, con simulaciones auditables y gates previos a `apply_candidate`.
- D.3 fija que `preview_ok` no equivale a apply ni sustituye enforcement real del kernel.
- D.4 fija `execution_inspection_view` y `recovery_action_candidate` como lectura/preparación sobre evidencia canónica, no como segunda verdad de runtime.
- D.5 fija `semantic_debug_view` y `maintenance_action_candidate` como ayuda asistida gobernada, sin convertir IA en autoridad de mutación.
- D.5 fija que explicación asistida sin `evidence_refs[]` suficientes no es aceptable como lectura normativa.
- C.6 ya validó el corredor mínimo del kernel entre `intent_input`, `compiled_contract`, runtime durable, policy, event log y capability resolution.
- C.6 dejó explícito que distribution layer sigue fuera del roadmap actual y no forma parte del cierre de Fase D.

## Corredor mínimo end-to-end de la surface

El corredor mínimo obligatorio de D.6 queda definido así:

1. **Entrada conversacional gobernable**
   - existe `conversation_turn` tenant-scoped;
   - existe `intake_session` correlado;
   - el intake conserva turns fuente, evidencia y reason codes.

2. **Shaping de intención**
   - existe `intent_candidate` o `intent_input` emitido de forma compatible con D.1;
   - no hay reinterpretación libre posterior para sostener lo emitido;
   - ambigüedades críticas quedan cerradas o explicitadas como bloqueo.

3. **Proposal workspace**
   - existe `proposal_draft` gobernado, correlado con el intake y con `source_intent_refs[]` resolubles;
   - existe diff humano/material separado y evidencia de origen suficiente.

4. **Preview formal**
   - existe `preview_candidate` consistente con el `proposal_draft` y su patchset candidateado;
   - la surface no promociona preview desde charla informal ni desde artifacts incompletos.

5. **Simulación mínima**
   - existe al menos un `simulation_result` válido;
   - el resultado declara `simulation_family`, `status`, `inputs_refs[]`, findings y evidencia.

6. **Lectura operativa o semántica coherente**
   - existe una `execution_inspection_view` o una `semantic_debug_view` correlada con el caso;
   - la vista no contradice artifacts ni records canónicos;
   - si la vista es parcial, debe degradar explícitamente.

7. **Cierre mínimo del checkpoint**
   - el smoke path puede terminar sin `apply` real;
   - alcanza con que la surface complete su corredor mínimo con evidencia y gates coherentes;
   - debe existir evidencia correlada con el kernel para sostener el cierre.

## Smoke path conversacional mínimo obligatorio

El smoke path mínimo obligatorio de D.6 no exige apply real; exige **integridad estructural del corredor de surface**.

Secuencia mínima:

1. ingresar un `conversation_turn` dentro de una `session_id` válida;
2. crear o actualizar `intake_session` correlado;
3. producir `intent_candidate` y decidir `emit_intent_input` o `emit_proposal_draft` por reglas explícitas;
4. persistir `proposal_draft` con `source_intent_refs[]` y evidencia de origen;
5. producir `preview_candidate` correlado con el draft y el patchset material;
6. producir al menos un `simulation_result` con `inputs_refs[]` resolubles;
7. exponer una `execution_inspection_view` o una `semantic_debug_view` coherente con artifacts de surface y records de kernel;
8. demostrar evidencia correlada con el kernel mediante IDs mínimos obligatorios;
9. cerrar el caso en alguno de estos estados válidos para D.6:
   - `preview_completed` con evidencia suficiente;
   - `preview_blocked` con `reason_codes[]` y trail íntegro;
   - `debug_incomplete` o `inspection_incomplete` con degradación explícita y evidence trail suficiente.

Reglas normativas del smoke path:

- puede terminar antes de `apply_candidate` o apply real;
- no puede declarar `ready for apply` si faltan gates explícitos o evidence trail suficiente;
- debe dejar trazabilidad completa entre surface y kernel aunque falle observabilidad derivada;
- debe ser repetible sin duplicar materialmente artifacts canónicos del caso;
- distribution layer NO puede aparecer como dependencia necesaria del smoke path.

## Gates de consistencia entre surface y kernel

Para que D.6 pase, los siguientes gates deben cumplirse:

1. **Conversation -> intake gate**
   - no existe shaping material sin `conversation_turn` y `intake_session` correlados;
   - si existe turn material y no existe `intake_session`, el corredor falla.

2. **Intake -> governed intent/proposal gate**
   - `intent_candidate` o `intent_input` deben sostenerse con evidencia de intake suficiente;
   - la surface no puede emitir handoff gobernado si faltan facts materiales o si persiste ambigüedad crítica.

3. **Proposal -> preview gate**
   - no existe `preview_candidate` válido sin `proposal_draft` correlado, `source_intent_refs[]` y patchset material consistente;
   - `preview_ok` no puede coexistir con `open_questions` críticas sin bloqueo o degradación explícita.

4. **Preview -> simulation gate**
   - todo `simulation_result` debe derivar de un `preview_candidate` real y de `inputs_refs[]` resolubles;
   - la surface no puede sintetizar simulaciones sin base material verificable.

5. **Surface -> kernel evidence gate**
   - la surface debe poder enlazar artifacts de surface con `contract_id`, `execution_id` y `event_id` o declarar explícitamente por qué no aplica;
   - explicación o inspección sin evidence trail suficiente no puede pasar el checkpoint.

6. **Inspection/debug -> canonical truth gate**
   - `execution_inspection_view` y `semantic_debug_view` deben alinearse con records canónicos del kernel;
   - si la UI/IA contradice artifacts reales, el checkpoint falla.

7. **Governance boundary gate**
   - si la surface salta governance, saltea boundaries o convierte asistencia en autoridad implícita de ejecución, el checkpoint falla;
   - ningún candidate puede transformarse en apply o maintenance real sin gates explícitos ya definidos en D.2-D.5.

## Gates de correlación y evidencia

Debe poder demostrarse correlación completa entre surface y kernel usando como mínimo:

- `tenant_id`
- `session_id`
- `trace_id`
- `conversation_turn_id`
- `intake_session_id`
- `proposal_draft_id`
- `preview_candidate_id`
- `simulation_result_id`
- `contract_id`
- `execution_id`
- `event_id`

Gates obligatorios:

1. `tenant_id` debe sobrevivir sin ambigüedad desde conversación hasta evidencia canónica.
2. `session_id` debe enlazar el corredor conversacional de surface.
3. `trace_id` debe correlacionar intake, proposal, preview y lectura operativa/semántica; si no existía, debe aparecer antes del primer hecho canónico que lo requiera.
4. `conversation_turn_id` debe permitir reconstruir el origen del shaping.
5. `intake_session_id` debe enlazar turns, decisiones de intake y handoff gobernado.
6. `proposal_draft_id`, `preview_candidate_id` y `simulation_result_id` deben poder recorrerse como cadena de surface sin saltos implícitos.
7. `contract_id`, `execution_id` y `event_id` deben permitir probar que la surface no está explicando un caso inventado o desacoplado del kernel.
8. La evidencia canónica debe poder sostener la lectura aunque fallen traces o dashboards derivados.
9. Si existe evidencia sólo en observability y no en records canónicos, el checkpoint no puede pasar.

## Gates de fail-safe y límites de automatización

D.6 debe validar degradación segura y límites estrictos de automatización:

- si la surface salta governance o saltea boundaries, el checkpoint no puede pasar;
- si la IA produce explicación sin evidence trail suficiente, el checkpoint no puede pasar;
- si aparece un `apply_candidate` implícito sin gates explícitos, el checkpoint no puede pasar;
- si un `maintenance_action_candidate` carece de `governance_requirements[]`, el checkpoint no puede pasar;
- si una vista asistida no puede sostener certeza material, debe degradar a ambigüedad explícita, `debug_incomplete` o `inspection_incomplete`;
- si falla observabilidad derivada, el checkpoint puede pasar sólo si artifacts y evidencia canónica siguen íntegros;
- si falla persistencia o correlación canónica material, el checkpoint no puede pasar aunque existan resúmenes, traces o dashboards;
- la surface puede sugerir recovery o maintenance sólo para acciones permitidas y gobernadas; nunca para mutación implícita ni bypass de approvals/policy.

## Artefactos mínimos que deben existir al cerrar D.6

Al cerrar D.6 deben existir, como mínimo, estos artefactos normativos o persistidos en el corredor:

1. `conversation_turn` persistido y correlado.
2. `intake_session` persistido con evidencia mínima.
3. `intent_candidate` o `intent_input` persistido/emitido con reason codes.
4. `proposal_draft` persistido con `source_intent_refs[]`.
5. `governed_patchset_candidate` o evidencia normativa equivalente necesaria para sostener preview.
6. `preview_candidate` persistido.
7. al menos un `simulation_result` persistido con `inputs_refs[]`.
8. una `execution_inspection_view` o una `semantic_debug_view` persistida o reproducible normativamente.
9. evidence trail que enlace el corredor de surface con `contract_id`, `execution_id` y `event_id` cuando aplique.
10. reason codes y/o findings del cierre del smoke path.

## Queries operativas mínimas que debe poder responder la surface

Al cerrar D.6, la surface ya debe poder responder preguntas operativas mínimas como estas:

1. ¿Qué `conversation_turn_id` y qué `intake_session_id` dieron origen a este caso?
2. ¿Qué `intent_candidate` o `intent_input` se emitió y con qué evidencia mínima?
3. ¿Qué `proposal_draft_id` se abrió para este pedido y cuáles son sus `source_intent_refs[]`?
4. ¿Qué `preview_candidate_id` corresponde a ese draft y con qué `material_diff` quedó correlacionado?
5. ¿Qué `simulation_result_id` existe, de qué familia es y con qué `inputs_refs[]` se generó?
6. ¿Qué preguntas abiertas siguen críticas y por qué bloquean o degradan la surface?
7. ¿La surface está mostrando una lectura coherente con `contract_id`, `execution_id` y `event_id` reales?
8. ¿Qué `trace_id` conecta intake, proposal, preview y evidencia del kernel?
9. ¿La explicación asistida tiene `evidence_refs[]` suficientes o quedó degradada?
10. ¿El caso quedó en `preview_completed`, `preview_blocked`, `debug_incomplete` o `inspection_incomplete`, y qué evidencia soporta ese cierre?

## Métricas mínimas del checkpoint

D.6 no exige observabilidad madura, pero sí métricas mínimas para afirmar que la surface integrada es operable:

- porcentaje de smoke paths con `conversation_turn` + `intake_session` correlados completos;
- porcentaje de handoffs de intake con evidencia mínima suficiente;
- porcentaje de `proposal_draft` con `source_intent_refs[]` y diffs separados válidos;
- porcentaje de `preview_candidate` con correlación íntegra hacia draft y patchset material;
- porcentaje de `simulation_result` con `inputs_refs[]` resolubles y evidence trail suficiente;
- porcentaje de vistas `execution_inspection_view`/`semantic_debug_view` sin contradicción contra evidencia canónica;
- conteo de casos donde falla observabilidad derivada pero artifacts y records canónicos siguen íntegros;
- conteo de bloqueos por governance bypass, correlación rota o explicación sin evidencia suficiente;
- tasa de pérdida de `trace_id` entre intake y preview;
- tasa de candidates mal promovidos a apply/maintenance sin gates explícitos.

El umbral normativo de D.6 no es UX ni throughput; es **integridad estructural del corredor de surface**.

## Gaps aceptables al cerrar D.6

Son aceptables al cierre de D.6:

- apply real no ejercitado dentro del smoke path mínimo;
- algunas vistas operativas todavía austeras, siempre que sean coherentes y auditables;
- observabilidad derivada parcial, diferida o inmadura, si la evidencia canónica y los artifacts siguen íntegros;
- simulaciones mínimas v1 todavía acotadas en fidelidad, siempre que lo declaren explícitamente;
- UX humana, dashboards y affordances avanzadas todavía diferidas a Fase E;
- automatización asistida todavía conservadora y bloqueante por defecto.

## Gaps NO aceptables al cerrar D.6

NO son aceptables al cierre de D.6:

- `conversation_turn` material sin `intake_session` correlado;
- handoff de intake sin evidencia mínima suficiente;
- `proposal_draft` sin `source_intent_refs[]` o sin evidencia de origen;
- `preview_candidate` sin patchset material consistente;
- `simulation_result` sin `inputs_refs[]` resolubles;
- lectura de surface que contradiga artifacts o records canónicos del kernel;
- pérdida de `trace_id` o correlación rota entre los IDs materiales mínimos;
- explicación asistida sin `evidence_refs[]` suficientes;
- promoción implícita a apply o maintenance sin gates explícitos;
- dependencia de observabilidad derivada para suplir ausencia de evidencia canónica;
- distribución, rollout o distribution layer introducidos como dependencia necesaria del checkpoint.

## Tests borde mínimos

Como mínimo, D.6 debe contemplar y documentar estos tests borde de integración de surface:

1. `conversation_turn` existe pero no genera `intake_session`.
2. `intake_session` emite `intent_input` sin evidence suficiente.
3. `proposal_draft` existe pero no tiene `source_intent_refs`.
4. `preview_candidate` existe sin `governed_patchset_candidate`.
5. `simulation_result` existe pero no tiene `inputs_refs`.
6. preview dice `preview_ok` y proposal sigue con `open_questions` críticas.
7. `semantic_debug_view` contradice artifacts reales.
8. `execution_inspection_view` no correlaciona con `execution_id`.
9. `maintenance_action_candidate` aparece sin `governance_requirements`.
10. `apply_candidate` implícito sin gates explícitos.
11. `trace_id` se pierde entre intake y preview.
12. evidence trail existe en observability pero no en records canónicos.
13. event log canónico existe pero la surface muestra resumen contradictorio.
14. IA-assisted explanation sin evidence refs.
15. `out_of_scope` mal tratado como proposal válido.
16. preview incompleto tratado como listo para apply.
17. recovery sugerido en caso no recuperable.
18. distribution layer aparece como dependencia del smoke path.
19. `conversation_turn` duplicado provoca múltiples `proposal_draft` materiales para el mismo caso sin deduplicación explícita.
20. `simulation_result` correlaciona con `preview_candidate_id` correcto pero con `proposal_draft_id` stale.
21. `execution_inspection_view` se arma sólo con observabilidad derivada y sin `event_id` canónico resoluble.
22. `semantic_debug_view` afirma un `contract_id` que no aparece en artifacts ni records reales.

Criterio de lectura de estos bordes:

- los casos 12, 13 y 21 prueban que observabilidad derivada y resúmenes no pueden reemplazar verdad canónica;
- los casos 6, 10, 16 y 17 prueban que la surface no puede promover ni recomendar de manera optimista;
- los casos 7, 8, 14 y 22 prueban que explicación e inspección deben anclarse a evidencia real;
- los casos 15 y 18 prueban que scope y roadmap siguen cerrados y no admiten reinterpretación oportunista.

## Criterios de aceptación de D.6

D.6 se considera cerrado a nivel de surface v1 sólo si se cumplen todos estos criterios:

1. existe un documento normativo del checkpoint de integración de la surface;
2. queda explícito que D.6 valida la **surface**, no el kernel ni el distribution layer;
3. queda definido el corredor mínimo end-to-end de la surface sin contradecir D.1-D.5 ni C.6;
4. queda definido el smoke path conversacional mínimo obligatorio y queda explícito que puede cerrar sin apply real;
5. queda explícito que el corredor mínimo debe incluir `conversation_turn`/`intake_session`, `intent_candidate` o `intent_input`, `proposal_draft`, `preview_candidate`, al menos un `simulation_result`, una vista coherente (`execution_inspection_view` o `semantic_debug_view`) y evidencia correlada con el kernel;
6. quedan definidos gates de consistencia entre surface y kernel;
7. queda definida la correlación mínima completa usando `tenant_id`, `session_id`, `trace_id`, `conversation_turn_id`, `intake_session_id`, `proposal_draft_id`, `preview_candidate_id`, `simulation_result_id`, `contract_id`, `execution_id` y `event_id`;
8. queda explícito que si falla observabilidad derivada, el checkpoint puede pasar sólo si artifacts y evidencia canónica siguen íntegros;
9. queda explícito que si la surface salta governance, saltea boundaries o explica sin evidence trail suficiente, el checkpoint NO pasa;
10. quedan listados artefactos mínimos, queries operativas mínimas, métricas mínimas, gaps aceptables/no aceptables y tests borde mínimos;
11. queda explícito que distribution layer sigue fuera del scope y no forma parte de ningún caso necesario para cerrar D.6;
12. queda explícito que Fase E absorbe el hardening, cierre final y base reusable posterior al checkpoint de surface.
