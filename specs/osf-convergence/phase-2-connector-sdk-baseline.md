# Phase 2 — Connector SDK baseline

## Objetivo

Definir el baseline del SDK estándar de conectores para que providers remotos, workers y capacidades ejecutables puedan integrarse al runtime convergente con contratos mínimos, evidencia, idempotencia, observabilidad y compatibilidad estable con registry/resolution.

## Principios del connector SDK baseline

- El SDK estandariza integración, no la semántica normativa del dominio.
- El SDK debe reducir variabilidad accidental entre conectores.
- El SDK debe forzar evidencia mínima, idempotencia y señales operativas correlables.
- El SDK no puede bypassear policy, classification, approvals ni runtime governance.
- El SDK debe servir tanto para conectores simples como para conectores con compensación y dry-run.

## Boundary exacto entre capability manifest, binding, provider y connector SDK

- **Capability manifest**: declara qué capability existe, qué contrato expone, qué políticas/boundaries aplica y cómo debe resolverse.
- **Binding**: asocia capability declarada con provider concreto, credenciales/referencias y configuración aprobada por entorno/tenant.
- **Provider**: implementación remota que sabe hablar con un sistema externo o ejecutar una acción real.
- **Connector SDK**: baseline técnico que estructura el provider y su interfaz operativa mínima.

Regla dura: el SDK implementa el lado ejecutable del provider, pero no reemplaza ni manifest ni binding ni runtime.

## Responsabilidades del SDK

- Exponer interfaz mínima homogénea para inspección, simulación y ejecución.
- Validar inputs estructurales antes de interactuar con sistemas externos.
- Propagar `idempotency_key`, `execution_id`, `tenant_id` y referencias de correlación.
- Emitir evidence refs mínimas y spans/eventos OTel obligatorios.
- Normalizar resultados al modelo esperado por runtime y evidence plane.
- Hacer explícito cuándo existe compensación y cuándo no.
- Declarar scopes, riesgo y capacidades soportadas por el provider.

## No-responsabilidades del SDK

- No decide autorización final.
- No define policy ni clasificación normativa.
- No reemplaza approvals.
- No persiste verdad operativa canónica por cuenta propia.
- No define manifests ni bindings normativos.
- No inventa workflow semantics por fuera de Temporal.

## Interfaz estándar mínima del conector

- `inspect()`
- `dry_run()`
- `execute()`
- `get_capabilities()`
- `get_risk_profile()`
- `get_required_scopes()`
- `normalize_result()`
- `compensate()` cuando aplique

## Contrato mínimo de input/output por método

### `inspect()`
- **input mínimo**: `tenant_id`, `capability_id`, `target_ref`, `binding_ref`, `requested_scope`, `trace_ref`
- **output mínimo**: disponibilidad del target, metadatos relevantes, restricciones detectadas, evidencia preliminar, clasificación sugerida, riesgo preliminar

### `dry_run()`
- **input mínimo**: `tenant_id`, `capability_id`, `compiled_contract_ref`, `binding_ref`, `idempotency_key`, `trace_ref`
- **output mínimo**: cambios esperados, objetos potencialmente afectados, riesgos detectados, evidence refs de simulación, resultado normalizado de preview

### `execute()`
- **input mínimo**: `tenant_id`, `capability_id`, `compiled_contract_ref`, `binding_ref`, `idempotency_key`, `execution_id`, `attempt`, `trace_ref`
- **output mínimo**: resultado bruto, resultado normalizado, evidence refs, estado técnico, classification hints, retryability, compensability

### `get_capabilities()`
- **input mínimo**: identidad del provider, versión del SDK, contexto del binding
- **output mínimo**: lista de capabilities soportadas, operaciones disponibles, limitaciones, versiones compatibles

### `get_risk_profile()`
- **input mínimo**: capability u operación objetivo, target scope, contexto del binding
- **output mínimo**: riesgo de negocio, riesgo de seguridad, factores agravantes, necesidad de approval sugerida

### `get_required_scopes()`
- **input mínimo**: capability u operación objetivo
- **output mínimo**: scopes/permisos mínimos requeridos, permisos opcionales, restricciones de uso

### `normalize_result()`
- **input mínimo**: resultado bruto, metadata de ejecución, contract/result context
- **output mínimo**: resultado canonizado, reason codes, severity, evidence refs ligadas, flags de partial/full/failure

### `compensate()`
- **input mínimo**: referencia a ejecución previa, target afectado, idempotency_key de compensación, trace_ref
- **output mínimo**: estado de compensación, evidencia asociada, efecto logrado o pendiente, escalación manual requerida si no aplica rollback completo

## Evidence refs mínimas

Todo método mutante o de inspección relevante debe devolver al menos:

- `trace_ref`
- `provider_call_ref`
- `input_snapshot_ref` o equivalente redactado
- `output_snapshot_ref` o equivalente redactado
- `artifact_ref` cuando haya adjuntos o archivos
- `error_ref` cuando exista fallo técnico o de dominio

## Idempotency key obligatoria

- `dry_run()` y `execute()` requieren `idempotency_key` obligatoria.
- `compensate()` requiere una key propia de compensación o una derivación explícita de la ejecución original.
- La deduplicación debe poder correlacionarse por tenant, capability, target scope, contract fingerprint y phase.
- El SDK debe tratar reintento técnico, reejecución controlada y replay de evidencia como cosas distintas.

## Spans/eventos OTel mínimos

Todo conector debe emitir al menos:

- span de `connector.inspect`
- span de `connector.dry_run`
- span de `connector.execute`
- span de `connector.compensate` cuando aplique
- evento `connector.request_sent`
- evento `connector.response_received`
- evento `connector.normalization_completed`
- evento `connector.evidence_emitted`
- evento `connector.retry_scheduled` cuando aplique
- evento `connector.failure_classified` cuando aplique

## Clasificación, riesgo y scopes mínimos que debe devolver o respetar

- Debe devolver o respetar `classification_level` esperado para outputs y evidence.
- Debe respetar `output_policy_ref` y `visibility_policy_ref` cuando lleguen resueltos desde capas superiores.
- Debe exponer riesgo de negocio y riesgo de seguridad al menos a nivel de hints o profile.
- Debe declarar scopes mínimos requeridos y no asumir privilegios implícitos.
- Debe marcar si un resultado incluye datos sensibles, parciales, redactados o no exportables.

## Compatibilidad con runtime y registry/resolution

- El SDK debe ser resoluble desde el registry a través de manifests y bindings aprobados.
- El runtime sigue orquestando secuencia, retries, pausas, compensación y approvals; el conector no toma ese control.
- El provider debe poder ser tratado como endpoint remoto gobernado dentro del model de remote provider/worker.
- Las versiones del SDK deben compatibilizarse con capability version, binding version y contract version sin romper el baseline normativo.

## Tests borde mínimos

1. inspección exitosa con target existente y evidencia mínima completa
2. inspección con target inexistente y failure normalizado
3. dry-run con cambios previstos sin efectos laterales reales
4. dry-run repetido con misma idempotency key sin divergencia material
5. execute exitoso con output bruto y normalizado correlados
6. execute repetido con misma idempotency key sin duplicación efectiva
7. execute con timeout remoto y clasificación correcta de retry técnico
8. execute con fallo funcional no reintentable
9. execute con evidence refs incompletas debe fallar validación del SDK
10. execute con datos sensibles exige output marcado/redactado
11. normalize_result sobre respuesta parcial válida
12. normalize_result sobre error remoto sin schema completo
13. get_capabilities devuelve lista consistente con manifest/binding
14. get_risk_profile distingue riesgo de negocio y riesgo de seguridad
15. get_required_scopes detecta scopes mínimos y opcionales
16. compensate exitosa sobre efecto parcialmente reversible
17. compensate no disponible obliga escalación manual explícita
18. provider emite spans/eventos OTel mínimos en camino feliz
19. provider emite eventos de fallo y retry en camino degradado
20. binding inválido bloquea ejecución antes de llamar al sistema externo

## Criterios de aceptación

1. Existe boundary claro entre manifest, binding, provider y SDK.
2. La interfaz mínima obligatoria queda fijada y completa.
3. Quedan fijados inputs/outputs mínimos por método.
4. La evidence mínima obligatoria queda escrita.
5. La `idempotency_key` queda exigida para operaciones que la requieren.
6. Los spans/eventos OTel mínimos quedan fijados.
7. Queda explícito que el SDK no bypass ea runtime, policy ni approvals.
8. Queda explícita la compatibilidad con registry/resolution.
9. Existen tests borde suficientes para validar el baseline.
10. El perfil es coherente con el seam declarative manifest + remote provider/worker ya cerrado.
