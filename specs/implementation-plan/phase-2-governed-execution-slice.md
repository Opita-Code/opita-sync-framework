# Phase 2 — Governed execution slice

## Objetivo

Extender el engine foundation slice hacia ejecución gobernada real: **approvals + blocking/release path**, **preview/simulation kernel hooks**, **compensation/unknown outcome path**, **evidence trail completo** y **smoke path engine-only**. Esta fase valida que el baseline no solo ejecuta, sino que gobierna correctamente.

## Qué entra

- Camino explícito de approvals con bloqueo, release y reanudación de ejecución.
- Hooks de preview/simulation dentro del kernel para evaluar impacto antes de apply real.
- Manejo explícito de compensación y de `unknown_outcome` cuando el runtime no puede afirmar cierre limpio.
- Evidence trail completo y correlado entre contrato, decisión, ejecución, approval y resultado.
- Smoke path engine-only que demuestre corredor gobernado mínimo sin depender todavía de surface rica.

## Qué queda fuera

- Surface operacional completa para intake, workspace y recovery visual.
- Hardening reusable final, regression baseline integral y demo final de referencia.
- Distribution layer, rollout, activation tenant-scoped o packaging de consumo.
- Automatización avanzada de operator workflows fuera del corredor mínimo gobernado.
- Escenarios enterprise exhaustivos de volumen, multiregión o compliance ampliado.

## Qué valida del baseline y qué todavía no valida

### Valida en esta fase

- Que approvals y governance gates pueden frenar y liberar ejecución durable sin romper consistencia.
- Que preview/simulation puede vivir como hook del kernel y no como motor paralelo.
- Que el baseline soporta evidence trail suficiente para auditoría operativa real.

### Todavía no valida

- Ergonomía operacional de la surface.
- Hardening reusable completo y readiness de freeze.
- Baseline extendido de SDK, object storage, Valkey y retrieval.

## Approvals / release path

- Debe existir punto claro de bloqueo antes de ejecución efectiva cuando el contrato o la policy lo requieran.
- Debe existir liberación explícita, correlada y auditable; no hay desbloqueo implícito por conveniencia.
- Debe quedar claro qué invalida approvals previas y cómo vuelve a evaluarse el fingerprint.
- Debe preservarse la separación entre decisión de governance y resultado operativo de ejecución.

## Compensation / unknown outcome

- Debe existir camino explícito para compensación cuando la aplicación no pueda cerrarse limpiamente.
- Debe existir estado operativo claro para `unknown_outcome`; no puede degradarse a éxito o fallo por falta de precisión.
- Debe quedar documentado qué evidencia mínima se exige antes de compensar, escalar o cerrar manualmente.
- Debe quedar explícito qué parte del baseline sigue siendo recovery manual y qué parte ya queda automatizada.

## Preview / simulation hooks del kernel

- Los hooks deben apoyarse en artifacts y seams reales del kernel, no en una maqueta lateral.
- Deben poder consultar policy, approvals esperables, clasificación y riesgo con insumos canónicos.
- Deben producir evidencia reutilizable por la surface futura, sin convertir preview en truth plane separada.
- Deben dejar claro qué predicen y qué no garantizan todavía.

## Evidence trail completo

- Debe correlacionar intent/proposal, contrato compilado, approval decision, execution lifecycle, outcome y compensación si existe.
- Debe conservar append-only en el event log operativo y exportar solo señales derivadas a observabilidad.
- Debe registrar decisiones y transiciones suficientes para auditoría, debugging y recovery.
- Debe aplicar redacción/clasificación antes de toda exportación fuera del plano operativo duradero.

## Smoke path engine-only

- Debe existir un corredor mínimo que permita demostrar: compilar, evaluar, bloquear si corresponde, liberar, ejecutar, registrar evidencia y cerrar outcome.
- Ese smoke path debe poder fallar de forma explicable y dejar evidencia suficiente para diagnóstico.
- El objetivo no es cobertura máxima; es consistencia real del slice gobernado.

## Criterios de done

- La ejecución puede bloquearse y liberarse por governance sin romper el workflow durable.
- Preview/simulation existe como hook real del kernel y deja artifacts/evidencia correlables.
- `unknown_outcome` y compensación tienen camino explícito y no quedan escondidos en errores genéricos.
- El evidence trail completo permite reconstruir el corredor gobernado sin huecos estructurales obvios.
- El smoke path engine-only pasa como demostración mínima del slice.
- Queda claro qué pasa a Phase 3 y qué todavía no está validado del baseline reusable final.

## Riesgos

- Convertir approvals en lógica lateral en vez de integrarlas al lifecycle real del runtime.
- Diseñar preview/simulation como segundo motor y no como hook del kernel.
- Subestimar `unknown_outcome` y terminar falseando cierre operativo.
- Producir demasiada telemetría y poca evidencia canónica útil.
- Intentar resolver UX de operator en esta fase y distraer el objetivo del slice gobernado.
