# Phase 3 — Surface operational slice

## Objetivo

Construir la **surface operacional mínima** que permita usar el baseline gobernado de forma real: **intake**, **proposal workspace**, **preview surface**, **inspection/recovery** y **AI-friendly maintenance**. Esta fase no redefine el engine; lo vuelve operable.

## Qué entra

- Intake/shaping mínimo que convierta conversación o input libre en intent gobernada.
- Proposal workspace mínimo para materializar el pipeline Intent -> Change Proposal -> Governed Patchset.
- Preview surface mínima para leer diff, simulación, riesgo, approvals y gates antes de ejecutar.
- Inspection/recovery surface mínima para seguir ejecuciones, inspeccionar evidencia y operar recovery permitido.
- AI-friendly maintenance surface mínima para debugging semántico, mantenimiento asistido y lectura de artifacts.
- Corredor mínimo de surface conectado al engine ya validado en Phase 2.

## Qué queda fuera

- UX final de producto, detalle visual cerrado o breadth funcional de frontend.
- Distribution layer, activation, rollout o instalación tenant-scoped.
- Automatización asistida que bypassée governance o aplique cambios sin gates explícitos.
- Hardening reusable final, demo final y freeze del baseline.
- Tooling enterprise completo de observabilidad, analytics o administración masiva.

## Qué valida del baseline y qué todavía no valida

### Valida en esta fase

- Que el baseline puede operarse con una surface útil sin romper seams del engine.
- Que conversación, proposal, preview e inspección pueden convivir sobre artifacts gobernados.
- Que mantenimiento asistido puede existir sin convertir la IA en autoridad soberana.

### Todavía no valida

- Hardening integral reusable.
- Readiness final de baseline freeze.
- Breadth de producto o escala operacional final.

## Intake

- Debe separar con claridad input libre de intent gobernada.
- Debe producir artifacts operables, no instrucciones sueltas pegadas al runtime.
- Debe rechazar explícitamente todo intento de apply directo desde conversación libre.
- Debe dejar claro qué información mínima necesita el proposal workspace para continuar.

## Proposal workspace

- Debe materializar `change_proposal` y `governed_patchset` candidato dentro del pipeline ya cerrado.
- Debe preservar diffs humanos y materiales sin inventar una segunda semántica de configuración.
- Debe dejar visibles los gates pendientes antes de preview o apply.
- Debe operar sobre artifacts gobernados y no sobre mutación directa del estado canónico.

## Preview surface

- Debe presentar diff, simulación, riesgo, approvals requeridas y restricciones relevantes.
- Debe reutilizar hooks/evidencia del kernel, no recalcular soberanamente su propia verdad.
- Debe mostrar con claridad qué está validado, qué es predicción y qué sigue bloqueado.

## Inspection / recovery

- Debe permitir inspeccionar execution lifecycle, outcome, evidence trail y correlación entre artifacts.
- Debe exponer recovery permitido sin mutar directamente el estado canónico fuera del corredor gobernado.
- Debe hacer visible `blocked`, `failed`, `compensating` y `unknown_outcome` como estados operables reales.

## AI-friendly maintenance

- Debe permitir lectura guiada de artifacts, debugging semántico y mantenimiento asistido.
- Debe respetar policy, approvals, clasificación y límites de automatización.
- Debe dejar claro cuándo la IA propone, cuándo resume y cuándo solo ayuda a inspeccionar.
- No puede transformarse en un canal lateral de ejecución privilegiada.

## Corredor mínimo de surface

- Intake mínimo -> proposal workspace mínimo -> preview surface mínima -> ejecución gobernada del engine -> inspection/recovery mínima.
- Ese corredor debe dejar evidencia correlada y utilizable por developer y operator.
- El objetivo es demostrar uso real del baseline; no cerrar todas las variantes de experiencia.

## Criterios de done

- Existe un intake mínimo que no filtra chat libre directo al motor.
- Existe un proposal workspace mínimo que materializa artifacts gobernados.
- Existe una preview surface mínima conectada al kernel real.
- Existe una inspection/recovery surface mínima capaz de seguir estados y evidencia relevantes.
- Existe una maintenance surface mínima IA-friendly sin bypass de governance.
- El corredor mínimo de surface funciona sobre el engine ya validado y deja claros los gaps que pasan a Phase 4.

## Riesgos

- Reabrir seams del engine por problemas de UX mal encuadrados.
- Convertir la IA en operador soberano por conveniencia de surface.
- Diseñar demasiada interface y poca operabilidad real del corredor mínimo.
- Duplicar lógica del kernel dentro de la surface y romper consistencia.
- Confundir “surface mínima operable” con “producto completo”.
