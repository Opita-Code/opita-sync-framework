# Fase D — Estado actual

## Objetivo de Fase D

Fase D existe para montar la **surface conversacional y la operación IA-friendly** sobre el kernel ya cerrado en Fase C. El foco no es reabrir decisiones del motor, sino definir cómo conversación, proposals, preview, inspección operativa y mantenimiento gobernado se apoyan sobre artifacts, configuración y evidencia ya compatibles con el engine.

## Estado general

- Estado: **cerrada a nivel de surface v1**
- D.1 Conversation intake and intent shaping: **cerrado a nivel de surface v1**
- D.2 Governed change proposal workspace: **cerrado a nivel de surface v1**
- D.3 Simulation, diff and preview surface: **cerrado a nivel de surface v1**
- D.4 Execution inspection and operator recovery surface: **cerrado a nivel de surface v1**
- D.5 AI-friendly development and maintenance surface: **cerrado a nivel de surface v1**
- D.6 Phase D integration checkpoint: **cerrado a nivel de surface v1**
- Próximo bloque recomendado: **Fase E — Hardening, cierre y base reusable**

## Bloques de trabajo ordenados

1. **conversation intake and intent shaping**
2. **governed change proposal workspace**
3. **simulation, diff and preview surface**
4. **execution inspection and operator recovery surface**
5. **AI-friendly development and maintenance surface**
6. **phase D integration checkpoint**

## Baseline heredado de Fase C

- El compilador **intent → compiled contract** ya fija el boundary entre intención libre y contrato ejecutable.
- El runtime durable sobre **Temporal** ya resuelve ejecución, bloqueo, falla, compensación y trazabilidad base.
- La integración con **Cerbos** ya fija autorización contextual, governance y fail-closed como comportamiento del motor.
- El **event log operativo append-only** y la observabilidad derivada con **OTel/LGTM** ya proveen evidencia y correlación base.
- El **capability registry and resolution** ya fija cómo el motor descubre artifacts, bindings y providers sin reabrir el seam de extensibilidad.
- El checkpoint end-to-end de Fase C ya validó el corredor mínimo entre contrato compilado, runtime, policy, event log y capability resolution.
- La surface de Fase D debe operar sobre **artifacts/configuración gobernada**, no sobre distribución, instalación o activación operacional de tenants.
- D.6 cerró el checkpoint end-to-end de la surface validando el corredor mínimo conversacional/operativo sobre artifacts gobernados y evidencia correlada con el kernel.
- El smoke path de D.6 puede cerrar sin apply real, siempre que `conversation_turn`, intake, proposal, preview, simulación y lectura operativa/semántica queden íntegros y correlados.
- La surface conversacional/operativa quedó integrada a nivel v1 sobre el kernel cerrado, sin habilitar bypass de governance ni mutación implícita.

distribution layer queda fuera del roadmap actual
