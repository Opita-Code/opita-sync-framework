# 01 — Plan maestro de 5 fases

## Fase A — Cerrar verdad ejecutable
Objetivo: dejar el proyecto sin ambigüedades bloqueantes para construir el core.

Bloques:
- contrato
- approvals
- clasificación/output control
- capability model
- tenant/runtime/eventos
- evals + cierre

## Fase B — Investigación y comparación de decisiones irreversibles
Objetivo: elegir stack y patrones duros del motor.

Comparativas clave:
- policy engine
- durable runtime
- extensibilidad
- memoria operativa y telemetría
- configuración conversacional
- packaging de capabilities

## Fase C — Construcción del kernel
Objetivo: tener el motor base funcional.

Incluye:
- compilador intención → contrato
- policy
- clasificación/redacción
- approvals
- runtime durable
- event log
- capability registry
- instalación global/tenant
- multi-tenant base

## Fase D — Superficie conversacional y operación IA-friendly
Objetivo: gobernar el motor desde chatbot sin perder control.

Incluye:
- edición conversacional de configuración declarativa
- diff y preview antes de aplicar
- instalación conversacional de tools/workflows
- simuladores de policy/approval/classification
- inspección de ejecuciones

## Fase E — Cierre, endurecimiento y base reusable
Objetivo: dejar Opyta Sync Engine listo como base para futuros proyectos.

Incluye:
- regresión completa
- hardening
- tenant template
- starter kit para capabilities
- playbooks de onboarding y rollback
- documentación IA-first final
- proyecto demo de referencia
