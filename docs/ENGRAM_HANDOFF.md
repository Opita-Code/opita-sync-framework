# Engram handoff — Opita Sync Framework

## Repo

- GitHub: `https://github.com/Opita-Code/opita-sync-framework`
- Scope: **framework / kernel reusable**
- Repo hermano: `https://github.com/Opita-Code/opita-sync`

## Boundary

- **OSF** conserva autoridad sobre runtime, policy, approvals, evidence y execution semantics
- **Opita Sync** conserva autoridad sobre UX, tenant onboarding, catálogo visible, roadmap y surfaces de producto

## Estado técnico actual

La suite completa pasa con:

- `go test ./...`

## Bloques ya implementados en este repo

- `S2` PostgreSQL hardening inicial
  - índices para correlación
  - tests `sqlmock` de roundtrip/persistencia

- `S3` approvals and release hardening
  - `release/reject/escalate`
  - actor decisor, comentario, reason codes
  - `source_contract_fingerprint`
  - invalidación por fingerprint mismatch

- `S4` evidence trail hardening
  - refs canónicos en `events.Record`
  - correlación surface -> compile -> runtime -> inspection

- `S5` recovery and compensation minimum
  - subset soportado y auditado
  - candidates `blocked` cuando no cumplen precondiciones
  - `unknown_outcome` preservado

- `S6` tenant bootstrap baseline
  - dominio `tenant`
  - stores memory/postgres
  - `POST /v1/tenants/bootstrap`
  - `GET /v1/tenants/{tenant_id}`

- `M3.2`
  - baseline profiles reales para policy / approval / classification

- `M3.3`
  - catálogo visible por tenant
  - conectores habilitados por tenant
  - `GET /v1/tenants-catalog/{tenant_id}`
  - `GET /v1/tenants-connectors/{tenant_id}`

- `M4.1`
  - `GET /v1/workspaces/intake-proposal`

- `M4.2`
  - `GET /v1/readable-previews/{preview_id}`

## Próximo paso recomendado

Seguir con:

1. `M4.3` inspection/recovery usable
2. `M5.1` modelado de primera vertical
3. `M5.2` conectores del dominio
4. `M5.3` piloto

## Archivos clave para retomar

- `internal/app/operatorsurface/http.go`
- `internal/app/previewservice/http.go`
- `internal/app/surfaceservice/http.go`
- `internal/app/tenantservice/http.go`
- `internal/engine/foundation/orchestrator.go`
- `internal/platform/postgres/schema.sql`
- `internal/e2e/slice_test.go`

## Orden mental correcto

1. leer `docs/REPO_BOUNDARY.md`
2. correr `go test ./...`
3. leer el repo de producto (`Opita Sync`) para entender el source of truth
4. implementar en este repo siguiendo esa dirección

## Nota sobre Engram

La memoria local de Engram no viaja sola con git.

Este archivo resume el estado útil actual para continuar en otra PC después de un `fetch`.
