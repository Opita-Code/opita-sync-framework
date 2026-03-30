# Runbook técnico mínimo

## Objetivo

Explicar cómo levantar y recorrer la vertical slice actual de **Opita Sync Framework** sin depender de memoria tribal.

## Modos de ejecución

### Modo memoria

Se activa si **no** existe:

- `OSF_DATABASE_URL`

En este modo:

- los artifacts viven en stores en memoria
- sirve para bootstrap y validación rápida
- no debe confundirse con baseline durable final

### Modo PostgreSQL

Se activa si existe:

- `OSF_DATABASE_URL`

En este modo:

- core y varios artifacts de surface persisten en PostgreSQL
- se aplica `internal/platform/postgres/schema.sql`

## Policy

### Stub en memoria

Si no existe:

- `OSF_CERBOS_URL`

se usa el policy engine en memoria.

### Cerbos real

Si existe:

- `OSF_CERBOS_URL`

el servicio usa el cliente HTTP real de Cerbos.

## Arranque esperado

Servicio actual:

- `cmd/intent-service`

El servicio:

1. carga registry declarativo desde `definitions/capabilities/`
2. selecciona policy engine (memory o Cerbos)
3. selecciona stores (memory o postgres)
4. corre `Warmup()` del orquestador
5. expone endpoints HTTP del slice

## Corredor recomendado para prueba manual

1. crear intake turn
2. crear proposal draft
3. crear patchset
4. crear preview
5. compilar intención
6. inspeccionar execution/foundation run
7. si aplica, consultar approval y liberar
8. inspeccionar eventos y debug view

### Demo de referencia

Existe un demo reproducible en:

- `demo/reference/README.md`
- `demo/reference/demo.http`

## Endpoints clave por rol

### Engine

- compile intent
- get contract
- get execution
- get events
- get registry resolution
- get/release approval

### Surface

- intake turns / sessions / candidates
- proposals / patchsets
- previews / simulations
- inspection / recovery
- debug / maintenance

## Qué revisar si algo falla

### Registry

- `definitions/capabilities/*.yaml`
- compatibilidad de `contract_version`
- `result_type`
- `environment`

### Policy

- si `OSF_CERBOS_URL` está seteada, verificar Cerbos accesible
- si no, confirmar que se usa stub en memoria

### Persistencia

- si `OSF_DATABASE_URL` está seteada, validar conexión y esquema
- si no, confirmar que el modo memoria es el esperado

### Correlación

Revisar al menos:

- `contract_id`
- `contract_fingerprint`
- `execution_id`
- `trace_id`
- `approval_request_id`
- `policy_decision_id`

## Límite explícito del runbook

Este runbook cubre la **vertical slice actual**.

No cubre:

- distribution layer
- activation/rollout tenant-scoped
- UI final
- operación enterprise final
