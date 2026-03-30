# Fase C — Estado actual

## Objetivo de Fase C

Fase C existe para construir el kernel operativo de Opyta Sync sobre el baseline ya cerrado en Fase B. El foco no es expandir superficie de producto, sino fijar la secuencia de construcción, los seams estables del core y los checkpoints mínimos de integración que permitan pasar de arquitectura decidida a kernel ejecutable sin reabrir decisiones estructurales.

## Estado general

- Estado: **cerrada a nivel de construcción v1**
- C.1 Intent → Contract compiler: **cerrado a nivel de construcción v1**
- C.2 Execution runtime skeleton: **cerrado a nivel de construcción v1**
- C.3 Cerbos integration: **cerrado a nivel de construcción v1**
- C.4 Event log and observability base: **cerrado a nivel de construcción v1**
- C.5 Capability registry and resolution: **cerrado a nivel de construcción v1**
- C.6 Kernel integration checkpoint: **cerrado a nivel de construcción v1**
- Próximo bloque recomendado: **Fase D — Surface conversational and AI-friendly operations**

## Bloques de trabajo ordenados

1. **intent → contract compiler**
2. **execution runtime skeleton**
3. **Cerbos integration**
4. **event log and observability base**
5. **capability registry and resolution**
6. **kernel integration checkpoint**

## Baseline heredado de Fase B

- **Temporal** queda como runtime durable único del kernel.
- **Cerbos** queda como policy engine baseline para decisiones contextuales y auditables.
- **PostgreSQL** conserva la memoria operativa durable y gobernable del core.
- **OpenTelemetry + Grafana LGTM** absorben observabilidad derivada, no la verdad operativa.
- **Declarative manifest + remote provider/worker model** define el seam de extensibilidad.
- **OCI bundle inmutable + firma + attachments** fija el artifact base distribuible y trazable.
- **Intent → Change Proposal → Governed Patchset** queda como baseline de configuración conversacional gobernada.

## Lectura operativa final de Fase C

Fase C cierra sobre un baseline ya elegido y deja integrado el **kernel engine-only** a nivel de construcción v1. Eso significa que compilación de contrato, ejecución durable, enforcement vía Cerbos, evidencia canónica y capability resolution ya tienen un corredor mínimo end-to-end coherente, trazable y compatible entre seams.

El cierre de Fase C no expande distribution layer ni superficie conversacional. Deja listo el motor para que Fase D trabaje sobre surfaces operativas y AI-friendly sin reabrir la estructura del kernel.
