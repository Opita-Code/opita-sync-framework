# Opita Sync Framework

> Repo scope: **framework/kernel reusable**.  
> Product scope lives separately in the sibling repo: **Opita Sync**.

## Qué es esto

**Opita Sync Framework (OSF)** es el framework/kernel reusable sobre el que se puede construir **Opita Sync**.

La separación es obligatoria:

- **OSF** = framework / kernel gobernado y canónico
- **Opita Sync** = producto / control plane IA-First que consume ese kernel

Este repositorio hoy contiene:

- el **baseline normativo** documentado en `specs/`
- la **convergencia OSF** como addendum técnico en `specs/osf-convergence/`
- un **implementation plan** en `specs/implementation-plan/`
- una **vertical slice implementada** en Go del kernel + surface mínima

La línea de **productization gap closure** del producto vive en el repo hermano **Opita Sync**.

`Opita Sync` no es este framework.  
`Opita Sync` es el producto que consume OSF para ofrecer una experiencia **IA-First** de operación gobernada, y vive en un repo separado.

La frase rectora del proyecto es:

> **Opita Sync es IA-First en experiencia y OSF-First en autoridad.**

## Estado actual

La vertical slice actual queda congelada como **baseline reusable v1** dentro del scope actual.

Su madurez operativa sigue siendo de **alfa técnica**.

Ya materializa, de forma mínima:

- compiler path
- policy integration (stub + Cerbos opcional)
- runtime básico
- approvals/release path
- event log canónico
- capability registry/resolution
- intake
- proposal
- preview/simulation
- inspection/recovery
- semantic debug / maintenance candidates

## Qué es Opita Sync como producto

**Opita Sync** es un **control plane gobernado, IA-First**, construido sobre OSF.

Su objetivo inicial es dejar un tenant `operable` y permitir este corredor:

`intake -> proposal -> preview -> governance -> execution -> inspection/recovery`

La primera vertical de producto es:

- **gobernanza de cambios operativos del tenant**

Roles iniciales del producto:

- `admin tenant`
- `operator`
- `approver`

Opita Sync genera valor en dos etapas:

1. **implementación del tenant** hasta dejarlo operable
2. **uso posterior del sistema** sobre el corredor gobernado

## Cómo correrlo

### Requisitos mínimos

- Go 1.24+
- Opcional: PostgreSQL accesible si querés modo durable real
- Opcional: Cerbos si querés usar PDP real en lugar del stub

### Variables de entorno

- `OSF_DATABASE_URL`
  - si está presente, habilita persistencia en PostgreSQL
  - si no está, el sistema corre en modo memoria

- `OSF_CERBOS_URL`
  - si está presente, usa Cerbos real como PDP
  - si no está, usa un policy engine en memoria para bootstrap

### Arranque

Servicio principal actual:

- `cmd/intent-service`

Runbook detallado:

- `docs/RUNBOOK.md`

Límites y alcance de la alfa técnica:

- `docs/ALPHA_SCOPE.md`

## Endpoints actuales

### Engine / core

- `GET /healthz`
- `POST /v1/intents/compile`
- `GET /v1/contracts/{contract_id}`
- `GET /v1/executions/{execution_id}`
- `GET /v1/events?execution_id=...`
- `GET /v1/foundation/runs/{execution_id}`
- `GET /v1/registry/resolve?...`
- `GET /v1/approvals/{approval_request_id}`
- `POST /v1/approvals/{approval_request_id}/release`

### Surface mínima

- `POST /v1/intake/turns`
- `GET /v1/intake/sessions/{id}`
- `GET /v1/intake/candidates/{id}`
- `POST /v1/proposals`
- `GET /v1/proposals/{id}`
- `POST /v1/patchsets`
- `GET /v1/patchsets/{id}`
- `POST /v1/previews`
- `GET /v1/previews/{id}`
- `GET /v1/simulations?preview_id=...`
- `GET /v1/inspection/executions/{execution_id}`
- `POST /v1/recovery-actions`
- `GET /v1/recovery-actions/{id}`
- `POST /v1/recovery-actions/{id}/execute`
- `GET /v1/debug/semantic?execution_id=...`
- `POST /v1/maintenance-actions`
- `GET /v1/maintenance-actions/{id}`

## Corredor mínimo actual

`conversation_turn -> intake_session -> intent_candidate -> proposal_draft -> patchset_candidate -> preview_candidate -> simulation_result -> compiled_contract -> policy_decision -> registry_resolution -> execution_record -> approval_request/release -> unknown_outcome/compensation_pending -> event_log -> inspection/debug/maintenance artifacts`

## Boundary recomendado

Si una responsabilidad toca:

- runtime
- policy
- approvals core
- contracts
- evidence canónica
- event log

entonces pertenece a **OSF**.

Si una responsabilidad toca:

- UX
- onboarding de tenant
- catálogo visible
- operator/admin/approver surfaces
- realtime derivado
- activación y uso del producto

entonces pertenece a **Opita Sync**.

## Validación actual

La vertical slice ya pasó una primera corrida real de:

- `go test ./...`

## Qué NO es esto todavía

Esto **no** es todavía:

- producto completo
- distribution layer
- rollout tenant-scoped
- producto comercial final
- performance tuning enterprise
- UI final

## Dónde mirar según necesidad

- fuente de verdad del framework: `specs/`
- enriquecimiento técnico OSF: `specs/osf-convergence/`
- ejecución real por fases: `specs/implementation-plan/`
- gap closure de producto: `specs/productization-gap-closure/`
- definición central del producto: `source-of-truth/`
- slice implementado: `cmd/`, `internal/`, `definitions/`
