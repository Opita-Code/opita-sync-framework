# Checklist operativo de Fase C

## C.0 Kernel construction sequence

- [ ] definir secuencia y dependencias
- [ ] definir seams del kernel
- [ ] definir criterios de integración incremental

## C.1 Intent → Contract compiler

- [x] definir boundary de compilación
- [x] definir interfaces de input/output
- [x] definir normalización inicial
- [x] definir enrich pipeline mínimo
- [x] definir persistencia/versión del contrato compilado

## Profesionalización v1 de C.1
- [x] definir pipeline interno de compilación
- [x] definir representación intermedia entre intención y contrato compilado
- [x] definir componentes mínimos del compilador
- [x] definir `compiled_contract` + `compilation_report`
- [x] definir persistencia en PostgreSQL
- [x] definir política de fingerprint y deduplicación
- [x] definir diagnósticos y reason codes de compilación
- [x] definir idempotencia del compilador
- [x] definir integración preparatoria con Cerbos y Temporal
- [x] definir tests borde mínimos de compilación

## C.2 Execution runtime skeleton

- [x] definir mapeo de execution_record a runtime elegido
- [x] definir creación/cierre de ejecuciones
- [x] definir timers/retries base
- [x] definir separación execution/application
- [x] definir integración con compensation states

## Profesionalización v1 de C.2
- [x] definir mapping exacto entre `execution_record` y Temporal
- [x] definir workflow durable principal por `execution_id`
- [x] definir separación workflow vs activities
- [x] definir workflow/activity/signal/query mínimos
- [x] definir lifecycle de creación y progreso de ejecución
- [x] definir separación execution/application
- [x] definir timers, deadlines y pausas
- [x] definir retries y failure model del runtime skeleton
- [x] definir camino de `blocked`, `failed`, `unknown_outcome` y compensación
- [x] definir idempotencia y correlación operativa del runtime

## C.3 Cerbos integration

- [x] definir PDP boundary
- [x] definir policy inputs canonizados
- [x] definir policy decision mapping al runtime
- [x] definir audit trail mínimo
- [x] definir fallback/failure mode

## Profesionalización v1 de C.3
- [x] definir boundary exacto PEP/PDP
- [x] definir resource/action/principal/context model canonizado
- [x] definir puntos exactos del lifecycle donde se consulta policy
- [x] definir mapping de decisiones Cerbos al runtime
- [x] definir integración con approvals y clasificación
- [x] definir `policy_decision_record` o evidencia equivalente
- [x] definir correlación y persistencia mínima de decisiones
- [x] definir failure mode fail-closed
- [x] definir cache guard e invalidación por policy_version
- [x] definir tests borde mínimos de integración policy

## C.4 Event log and observability base

- [x] definir event log mínimo
- [x] definir esquema de correlación de IDs
- [x] definir instrumentación OTel mínima
- [x] definir proyección a memoria/telemetry
- [x] definir redacción/clasificación de señales

## Profesionalización v1 de C.4
- [x] definir boundary exacto entre event log operativo y telemetría derivada
- [x] definir `telemetry_event` / `event_record` canónico mínimo
- [x] definir event types mínimos del kernel
- [x] definir correlación exacta y propagación de IDs
- [x] definir qué persiste PostgreSQL vs qué se proyecta a OTel/LGTM
- [x] definir redacción y clasificación previa a exportación
- [x] definir strategy de emisión, proyección y backpressure
- [x] definir failure guard para que observabilidad no rompa el kernel
- [x] definir append-only e idempotencia del event log
- [x] definir tests borde mínimos de event log y observabilidad

## C.5 Capability registry and resolution

- [x] definir manifest schema ejecutable
- [x] definir relation manifest/bundle/binding/provider
- [x] definir registry mínimo
- [x] definir resolution de provider/worker
- [x] definir compatibilidad y versionado

## Profesionalización v1 de C.5
- [x] definir boundary exacto entre registry, packaging, binding y resolution
- [x] definir manifest schema ejecutable mínimo
- [x] definir binding model mínimo
- [x] definir relation manifest ↔ bundle ↔ binding ↔ provider
- [x] definir registry mínimo implementable
- [x] definir resolution flow del runtime
- [x] definir compatibilidad de versiones
- [x] definir verificación de artifact/digest/provenance
- [x] definir idempotencia y deduplicación del registry
- [x] definir tests borde mínimos del registry/resolution

## C.6 Kernel integration checkpoint

- [x] verificar consistencia end-to-end
- [x] correr smoke path del kernel
- [x] registrar gaps para Fase D/E

## Profesionalización v1 de C.6
- [x] definir corredor mínimo end-to-end del kernel
- [x] definir smoke path mínimo obligatorio
- [x] definir gates de consistencia entre seams
- [x] definir gates de correlación y evidencia
- [x] definir gates de fail-safe y degradación segura
- [x] definir artefactos mínimos de cierre del checkpoint
- [x] definir queries operativas mínimas
- [x] definir métricas mínimas del checkpoint
- [x] definir gaps aceptables y no aceptables
- [x] definir tests borde mínimos de integración
