# D.0 — Secuencia de construcción de la surface

## Principios base

- La surface existe para **operar y gobernar** el kernel; no para redefinirlo.
- Toda interacción conversacional debe terminar en artifacts, evidencia y estados compatibles con el motor ya cerrado.
- La surface debe separar con claridad **intención libre**, **propuesta gobernada**, **preview/simulación** y **apply**.
- Ninguna convenience de UX puede saltearse policy, approvals, clasificación, fingerprints o trazabilidad.
- La progresión de Fase D debe maximizar validación incremental sin abrir seams nuevos dentro del kernel.

## Qué significa “surface” en Opyta Sync

En Opyta Sync, “surface” significa la capa operacional y conversacional que permite a developers y operators:

- capturar intención útil desde lenguaje natural;
- transformarla en proposal gobernada;
- inspeccionar diffs, simulaciones y riesgo antes de aplicar;
- observar ejecuciones y responder a estados operativos relevantes;
- mantener artifacts, configuraciones y evidencias de manera IA-friendly.

No implica definir una UI concreta, un canal único ni un runtime paralelo. Implica definir contracts de interacción, evidencia mínima, boundaries y flujos operativos sobre el motor existente.

## Dependencias heredadas de Fase C

- El compilador ya produce `compiled_contract` y `compilation_report` persistidos y correlables.
- El runtime durable ya ejecuta sobre `execution_id` con separación entre execution y application.
- Policy enforcement, clasificación y governance ya tienen puntos de evaluación definidos sobre el lifecycle.
- El event log operativo ya conserva evidencia append-only y correlación de IDs.
- El registry/resolution de capabilities ya determina cómo se resuelven artifacts y providers.
- El checkpoint end-to-end ya cerró el corredor mínimo del kernel.

Estas dependencias implican que Fase D debe apoyarse en seams cerrados, no reinterpretarlos.

## Orden de construcción recomendado y por qué

1. **D.1 Conversation intake and intent shaping**  
   Primero, porque sin boundary claro entre chat libre e intent gobernado la surface entera queda ambigua y termina filtrando lenguaje natural directo al motor.

2. **D.2 Governed change proposal workspace**  
   Segundo, porque la intención útil necesita un espacio intermedio gobernado antes de hablar de apply. Ahí se materializa el pipeline `Intent → Change Proposal → Governed Patchset` ya heredado.

3. **D.3 Simulation, diff and preview surface**  
   Tercero, porque el proposal workspace necesita una lectura previa de impacto, policy, approvals y redacción antes de cualquier ejecución real.

4. **D.4 Execution inspection and operator recovery surface**  
   Cuarto, porque una vez que existe apply gobernado, la operator surface debe poder inspeccionar correlación y responder a estados como `blocked`, `failed` o `unknown_outcome`.

5. **D.5 AI-friendly development and maintenance surface**  
   Quinto, porque las vistas de desarrollo y mantenimiento deben construirse sobre flows ya definidos de intake, proposal, preview e inspección, no en paralelo desalineado.

6. **D.6 Phase D integration checkpoint**  
   Último, para verificar consistencia end-to-end de la surface y registrar gaps reales hacia Fase E.

## Boundaries que la surface NO puede romper

- No puede ejecutar cambios saltándose el compilador, el proposal flow o los gates de governance.
- No puede escribir directamente en runtime state como si la conversación fuera la verdad operativa.
- No puede desactivar policy, approvals, clasificación o redacción por conveniencia conversacional.
- No puede introducir un segundo modelo de artifacts distinto del ya cerrado en el kernel.
- No puede reintroducir **distribution layer**, rollout de consumo, instalación tenant-scoped o activación operacional de tenants como scope de Fase D.
- No puede convertir observabilidad derivada en fuente primaria de verdad por encima del event log operativo.

## Estrategia de integración progresiva

- Arrancar con el mínimo corredor: conversación acotada → intent gobernado → proposal draft.
- Agregar después preview/simulación como capa previa obligatoria a apply.
- Conectar luego inspección operacional sobre IDs y artifacts ya correlados por el kernel.
- Recién después expandir vistas IA-friendly de desarrollo y mantenimiento reutilizando las mismas primitivas.
- Cerrar con un checkpoint que valide coherencia entre lenguaje conversacional, proposal workspace, simulaciones, apply y recovery.

La clave es que cada bloque agregue **surface** sin introducir nueva lógica soberana fuera del engine.

## Qué NO construir todavía en Fase D

- Distribution layer o estrategias de distribución de capabilities.
- Tenant activation, onboarding operativo o rollout de consumo.
- Un sistema de instalación conversacional de tools/workflows como concern separado del artifact governance.
- UX final de producto, detalle visual de pantallas o decisiones concretas de frontend.
- Automatización asistida con apply implícito o auto-ejecución fuera de gates explícitos.
- Nuevos seams del kernel para compensar problemas de surface mal definidos.

## Criterios de aceptación de D.0

- La secuencia de construcción de Fase D queda explícita y justificada.
- Los boundaries entre surface y kernel quedan definidos sin ambigüedad.
- La integración progresiva preserva el pipeline gobernado ya cerrado en Fase B/C.
- Queda explícito qué items NO pertenecen a Fase D.
- Queda explícito que la surface opera sobre motor, artifacts y configuración gobernada, no sobre distribución.
