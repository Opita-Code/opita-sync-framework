# 02 — Checklist operativo de Fase A

## A.1 Objetos canónicos
- [x] definir lista final de objetos críticos del core
- [x] definir ownership de cada objeto
- [x] definir versión inicial de cada objeto
- [x] decidir formato fuente de verdad ejecutable
- [x] definir naming y convenciones globales

## Profesionalización v1 de A.1
- [x] definir criterio formal de qué califica como objeto canónico first-class
- [x] cerrar lista final de 15 objetos críticos del core
- [x] separar objetos top-level de snapshots y objetos embebidos
- [x] definir dominios de objetos: tenancy, governance, catalog, runtime, memory, observability
- [x] definir ownership por objeto: platform-owned, tenant-owned, runtime-owned, hybrid-governed
- [x] definir runtime_owner y mutable_by por objeto
- [x] definir versión inicial de schema y objeto
- [x] definir reglas de mutación permitida y prohibida por familia de objetos
- [x] definir reglas de promoción entre dev, staging y prod
- [x] decidir authoring source of truth en YAML
- [x] decidir formato compilado canónico en JSON determinístico
- [x] definir envelope canónico: api_version, kind, metadata, spec, status
- [x] definir layout de carpetas por dominio
- [x] definir naming conventions globales para kinds, ids, fields, enums, events y archivos
- [x] definir reglas de validación y compilación de la verdad ejecutable
- [x] definir tests borde mínimos del formato y naming
- [x] definir criterios de aceptación de A.1

## A.2 Contrato de intención/inspección
- [x] definir schema exacto del contrato
- [x] separar campos obligatorios, opcionales y derivados
- [x] definir validaciones duras
- [x] definir estados del contrato
- [x] definir transición intención → contrato compilable
- [x] definir diffs/versionado del contrato

## Profesionalización v1 de A.2
- [x] separar formalmente intención inicial vs contrato compilado
- [x] definir 4 grupos de campos: usuario, sistema, técnicos, snapshots
- [x] definir requisitos para contrato válido, compilable y ejecutable
- [x] definir material change fields y non-material fields
- [x] definir fingerprint determinístico del contrato
- [x] definir relación contrato vs approvals (fingerprint, snapshots, reapproval)
- [x] definir relación contrato vs clasificación
- [x] definir relación contrato vs memoria/contexto
- [x] definir 14 estados del contrato
- [x] definir transiciones válidas y prohibidas
- [x] definir proceso de compilación en 10 pasos
- [x] definir proceso de planificación hasta executable
- [x] definir política de versionado major/minor
- [x] definir semantic_diff y presentation_diff
- [x] definir reason codes de diff
- [x] definir eventos canónicos del contrato
- [x] definir payload mínimo por evento
- [x] definir 20 tests borde mínimos
- [x] definir criterios de aceptación de A.2

## A.3 Tipos de resultado
- [x] listar tipos de resultado soportados por el motor
- [x] definir input/output contract por tipo
- [x] definir criterios formales de éxito por tipo
- [x] definir evidencia mínima por tipo
- [x] definir causas de fallo por tipo

## Profesionalización v1 de A.3
- [x] proponer ajuste taxonómico: governance_decision reemplaza aprobación/rechazo
- [x] definir familias de tipos: read-only, mutation, governance
- [x] definir objeto result_type canónico con campos completos
- [x] definir relación resultado vs clasificación
- [x] definir relación resultado vs approvals (floors y overrides)
- [x] definir relación resultado vs auditoría (capas de telemetría por tipo)
- [x] definir objeto result_outcome con niveles formales
- [x] definir failure taxonomy con reason codes normalizados (4 familias)
- [x] definir reason codes de éxito parcial y degradado
- [x] definir métricas por tipo de resultado
- [x] definir matriz de outcome permitido por tipo
- [x] definir estados del ciclo de resultado (14 estados)
- [x] definir transiciones válidas y prohibidas
- [x] definir eventos canónicos de resultado (ciclo principal + por causa)
- [x] definir payload mínimo por evento
- [x] definir reglas de resultados parciales y redactados
- [x] definir qué se puede redactar por tipo
- [x] definir 20 tests borde mínimos
- [x] definir criterios de aceptación de A.3

## A.4 Approvals
- [x] definir acciones que nunca requieren aprobación
- [x] definir acciones que siempre requieren aprobación
- [x] definir acciones condicionales según riesgo, clasificación, impacto y ámbito
- [x] definir los 4 modos exactos: auto, pre-ejecución, pre-aplicación, doble aprobación
- [x] definir quién puede aprobar por tenant, rol, delegación y contexto
- [x] definir si una aprobación puede ser humana, política o híbrida
- [x] definir expiración de aprobaciones
- [x] definir revocación antes de ejecutar y antes de aplicar
- [x] definir qué invalida una aprobación ya emitida
- [x] definir si cambios de plan, datos, conectores o destino obligan re-aprobación
- [x] definir evidencia mínima para solicitar aprobación
- [x] definir evidencia mínima para auditar aprobación posterior
- [x] definir resultado si el aprobador no responde
- [x] definir resultado si el aprobador rechaza
- [x] definir resultado si el contexto cambia entre aprobación y ejecución
- [x] definir compatibilidad entre aprobación y delegación
- [x] definir compatibilidad entre aprobación y clasificación de datos
- [x] definir compatibilidad entre aprobación y publicación de workflows/tools
- [x] definir cómo se modela approval request
- [x] definir cómo se modela approval decision
- [x] definir estados del ciclo de aprobación
- [x] definir eventos canónicos de aprobación
- [x] definir reglas de idempotencia y reintento
- [x] definir matriz v1 riesgo × acción × contexto → modo de aprobación
- [x] definir casos borde que deben entrar a tests

## Profesionalización v1 de A.4
- [x] separar formalmente authorization, approval y execution
- [x] definir separation of duties por tipo de acción
- [x] definir SoD configurable por policy, tipo de acción, tenant y contexto
- [x] definir fingerprint/context snapshot de aprobación
- [x] definir invalidación automática por cambio material
- [x] definir reason codes normalizados para approve/reject/revoke/expire
- [x] definir versionado de policy dentro de la aprobación
- [x] definir snapshot de clasificación dentro de la aprobación
- [x] definir snapshot de risk score dentro de la aprobación
- [x] definir redacción de evidencia para approvals y auditoría
- [x] definir SLA/timeout policy por tipo de aprobación
- [x] definir fallback si nadie responde
- [x] definir simulación previa de approval
- [x] definir replay/audit mode
- [x] definir tests de SoD
- [x] definir tests de material change

## A.5 Tipo de tenant, runtime y eventos generales
- [x] definir tenant schema final
- [x] definir mínimo duro de creación
- [x] definir onboarding exacto
- [x] definir estados de ejecución
- [x] definir eventos canónicos generales
- [x] definir idempotencia y reintentos
- [x] definir compensación/rollback lógico

## Profesionalización v1 de A.5
- [x] definir `tenant` como objeto operable con estado explícito
- [x] cerrar schema final de tenant por grupos de campos
- [x] definir tenant operable vs tenant creado
- [x] definir `single_user` sin relajar seguridad, approvals ni clasificación
- [x] formalizar mínimo duro obligatorio de onboarding
- [x] definir proceso exacto de onboarding hasta `operable`
- [x] definir lifecycle del tenant y transiciones válidas/prohibidas
- [x] definir `execution_record` como objeto first-class del runtime
- [x] separar formalmente contract state vs execution state vs approval state vs result state
- [x] definir estados generales del runtime y release model
- [x] definir eventos canónicos generales del runtime
- [x] definir reglas de correlación entre `execution_id`, `trace_id`, `contract_id`, `approval_request_id` y `result_id`
- [x] definir observabilidad mínima obligatoria por ejecución
- [x] definir `idempotency_key` y reglas de deduplicación
- [x] definir retry técnico vs retry de ejecución vs replay de auditoría
- [x] definir política de reintentos por tipo de operación
- [x] definir tratamiento de `unknown outcome` en side effects externos
- [x] separar rollback físico de compensación lógica
- [x] definir estados/flags de compensación en ejecución
- [x] definir tests borde mínimos del tenant, runtime e idempotencia
- [x] definir criterios de aceptación de A.5

## A.6 Evals y testing base
- [x] definir evals por tipo de resultado
- [x] definir evals de policy
- [x] definir evals de approval
- [x] definir evals de clasificación/redacción
- [x] definir regresión mínima del core
- [x] definir criterios de salida de Fase A

## Profesionalización v1 de A.6
- [x] definir framework general de evals con taxonomía por familia
- [x] definir `eval_case` y output mínimo de eval
- [x] definir severidad, scoring y gating modes
- [x] definir evals por tipo de resultado
- [x] definir evals de policy
- [x] definir evals de approval
- [x] definir evals de clasificación y redacción
- [x] definir evals de runtime general
- [x] definir evals de resiliencia operacional
- [x] definir dataset strategy y fixtures mínimas
- [x] definir eventos mínimos de ejecución de evals
- [x] definir regresión mínima del core por suites
- [x] definir smoke vs standard vs release_candidate
- [x] definir cobertura mínima por dominio
- [x] definir artifacts y métricas obligatorias de regresión
- [x] definir política de flakes y quarantined tests
- [x] definir criterios formales de salida de Fase A
- [x] definir gates must_pass / should_pass / advisory
- [x] definir evidencia mínima y sign-off recomendado de cierre de Fase A
