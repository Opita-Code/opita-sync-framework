# Fase B — Estado actual

## Objetivo de Fase B

Fase B existe para investigar, comparar y cerrar las decisiones irreversibles que van a condicionar la construcción del kernel en Fase C. Su propósito no es implementar todavía el motor, sino reducir costo de reversión, fijar seams estables y evitar que decisiones de alto impacto queden abiertas durante la ejecución del core.

## Estado general

- Estado: **cerrada a nivel de decisión provisional**
- B.0 Marco común de decisión: **cerrado**
- B.1 Durable runtime comparison: **cerrado**
- Decisión provisional B.1: **Temporal queda recomendado como baseline provisional del durable runtime por durabilidad, multi-tenant y auditabilidad; la tensión con la futura capa declarativa conversacional queda abierta para B.6.**
- B.2 Policy engine comparison: **cerrado**
- Decisión provisional B.2: **Cerbos queda recomendado como baseline provisional del policy engine por operabilidad, auditabilidad, multi-tenant y scoped policies; el riesgo abierto es el encaje futuro de ReBAC profundo.**
- B.3 Operational memory and telemetry comparison: **cerrado**
- Decisión provisional B.3: **PostgreSQL + OpenTelemetry + Grafana LGTM queda recomendado como baseline provisional equilibrado; PostgreSQL sostiene memoria operativa estable y gobernable, y la observabilidad queda separada para trazas, logs, métricas y eventos de ejecución.**
- B.4 Extensibility model comparison: **cerrado**
- Decisión provisional B.4: **Declarative manifest + remote provider/worker model queda recomendado como baseline provisional de extensibilidad; preserva el catálogo gobernado, separa capability declarativa de handler ejecutable y mantiene abierta una capa conversacional futura sin competir con Temporal como runtime de verdad.**
- B.5 Capability packaging comparison: **cerrado**
- Decisión provisional B.5: **OCI bundle inmutable + firma + attachments queda recomendado como baseline provisional de packaging; fija el artifact base por digest, deja publication/distribution en el catálogo gobernado y mantiene tenant activation como objeto/acto separado.**
- B.6 Conversational configuration comparison: **cerrado**
- Decisión provisional B.6: **Intent → Change Proposal → Governed Patchset queda recomendado como baseline provisional de configuración conversacional; la conversación captura intención y ayuda a refinar, pero la verdad operativa vive en una propuesta declarativa gobernada con diff, simulación, aprobación y apply separado de la activación tenant.**
- Resumen final B.1-B.6: **Temporal + Cerbos + PostgreSQL/OTel/Grafana LGTM + declarative manifest/remote provider-worker + OCI bundle firmado + Intent→Change Proposal→Governed Patchset forman el baseline consolidado para construir el kernel sin reabrir decisiones irreversibles en Fase C.**
- Próximo bloque recomendado: **Fase C — Construcción del kernel**

## Comparativas pendientes

1. ninguna dentro de Fase B

## Decisión ya tomada

La recomendación arquitectónica ya cerrada es **empezar por durable runtime**, porque es la decisión con mayor costo de reversión y mayor efecto sistémico sobre todo el kernel. Condiciona el event log, la estrategia de retries, la compensación, la trazabilidad y los seams de extensibilidad que luego deberán respetar policy, memoria, packaging y superficie conversacional.
