# Checklist operativo de Fase D

## D.0 Surface construction sequence

- [ ] definir secuencia y dependencias de la surface
- [ ] definir boundaries con el kernel
- [ ] definir criterios de integración progresiva

## D.1 Conversation intake and intent shaping

- [x] definir boundary entre chat libre e intent gobernado
- [x] definir shape de conversation turn útil para el motor
- [x] definir handoff a compilador/proposal flow
- [x] definir gestión de ambigüedad crítica
- [x] definir evidencia mínima de intake

## Profesionalización v1 de D.1
- [x] definir boundary exacto entre chat libre, intent gobernado y proposal-trigger
- [x] definir `conversation_turn`, `intake_session` e `intent_candidate`
- [x] definir reglas de shaping de intención
- [x] definir ambigüedad crítica vs tolerable
- [x] definir preguntas de aclaración y cortes del flow
- [x] definir handoff exacto a `intent_input` o `proposal_draft`
- [x] definir evidencia mínima de intake
- [x] definir reason codes de intake
- [x] definir correlación mínima del intake
- [x] definir tests borde mínimos de intake

## D.2 Governed change proposal workspace

- [x] definir workspace de propuesta gobernada
- [x] definir lifecycle de proposal draft
- [x] definir diff humano vs diff material
- [x] definir gating previo a apply
- [x] definir evidencia de propuesta y aprobación

## Profesionalización v1 de D.2
- [x] definir boundary exacto entre intent, proposal, preview y apply candidate
- [x] definir `proposal_draft` y `governed_patchset_candidate`
- [x] definir lifecycle completo del proposal draft
- [x] definir diff humano vs diff material
- [x] definir gates previos a preview
- [x] definir gates previos a apply candidate
- [x] definir evidencia mínima del proposal workspace
- [x] definir reason codes del workspace
- [x] definir correlación con intake y preview
- [x] definir tests borde mínimos del proposal workspace

## D.3 Simulation, diff and preview surface

- [x] definir preview mínimo
- [x] definir simulación de policy
- [x] definir simulación de approval
- [x] definir simulación de clasificación/redacción
- [x] definir lectura de riesgo y efectos esperados

## Profesionalización v1 de D.3
- [x] definir boundary exacto entre proposal, preview, simulation result y apply candidate
- [x] definir `preview_candidate` y `simulation_result`
- [x] definir tipos mínimos de preview
- [x] definir simulación de policy
- [x] definir simulación de approval
- [x] definir simulación de clasificación/redacción
- [x] definir simulación de riesgo/impacto
- [x] definir diff humano vs diff material dentro del preview
- [x] definir gates de promoción desde preview
- [x] definir tests borde mínimos de preview/simulación

## D.4 Execution inspection and operator recovery surface

- [x] definir vistas mínimas de ejecución
- [x] definir correlación visible entre contrato/policy/runtime/eventos
- [x] definir herramientas mínimas de recovery operacional
- [x] definir surface de blocked/failed/unknown_outcome
- [x] definir auditoría consultable por operador

## Profesionalización v1 de D.4
- [x] definir boundary exacto entre inspección, recovery permitido y acción prohibida
- [x] definir `execution_inspection_view` y `recovery_action_candidate`
- [x] definir vistas mínimas de inspección
- [x] definir correlación visible entre IDs y artifacts
- [x] definir tratamiento explícito de blocked/failed/unknown_outcome/compensation
- [x] definir recovery permitido y recovery prohibido
- [x] definir gates previos a recovery
- [x] definir evidencia mínima y reason codes de recovery
- [x] definir integración con runtime/event log/policy/approvals
- [x] definir tests borde mínimos de operator surface

## D.5 AI-friendly development and maintenance surface

- [x] definir vistas/flows para desarrollo IA-friendly
- [x] definir surface de debugging semántico
- [x] definir surface de inspección de artifacts/configuraciones
- [x] definir surface de mantenimiento gobernado
- [x] definir límites de automatización asistida

## Profesionalización v1 de D.5
- [x] definir boundary exacto entre ayuda asistida, debugging semántico y acción gobernada
- [x] definir `semantic_debug_view` y `maintenance_action_candidate`
- [x] definir flows mínimos para developers y operators
- [x] definir debugging semántico
- [x] definir inspección gobernada de artifacts/configs/contracts/workflows
- [x] definir surface de mantenimiento gobernado
- [x] definir límites de automatización asistida
- [x] definir gates previos a maintenance actions
- [x] definir evidence trail y reason codes de D.5
- [x] definir tests borde mínimos de surface IA-friendly

## D.6 Phase D integration checkpoint

- [x] verificar consistencia end-to-end de la surface
- [x] correr smoke path conversacional mínimo
- [x] registrar gaps para Fase E

## Profesionalización v1 de D.6
- [x] definir corredor mínimo end-to-end de la surface
- [x] definir smoke path conversacional mínimo obligatorio
- [x] definir gates de consistencia entre surface y kernel
- [x] definir gates de correlación y evidencia
- [x] definir gates de fail-safe y límites de automatización
- [x] definir artefactos mínimos de cierre del checkpoint
- [x] definir queries operativas mínimas
- [x] definir métricas mínimas del checkpoint
- [x] definir gaps aceptables y no aceptables
- [x] definir tests borde mínimos de integración de surface
