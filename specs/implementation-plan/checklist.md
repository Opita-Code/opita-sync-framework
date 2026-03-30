# Checklist operativo del implementation plan

## Phase 1 — Engine foundation slice

- [x] materializar compiler contract path
- [x] materializar runtime skeleton mínimo
- [x] materializar policy integration mínima
- [x] materializar event log canónico mínimo
- [x] materializar registry/resolution mínimo

## Phase 2 — Governed execution slice

- [x] approvals + blocking/release path
- [x] preview/simulation kernel hooks
- [x] compensation/unknown outcome path
- [x] evidence trail completo
- [x] smoke path engine-only

## Phase 3 — Surface operational slice

- [x] intake/shaping mínimo
- [x] proposal workspace mínimo
- [x] preview surface mínima
- [x] inspection/recovery surface mínima
- [x] AI-friendly maintenance surface mínima

## Phase 4 — Hardening and reusable baseline slice

- [x] full regression baseline ejecutable
- [x] connector sdk baseline inicial
- [x] object storage / Valkey / retrieval baseline
- [x] docs/playbooks mínimos operativos
- [x] reference demo ejecutable

## Phase 5 — Go/No-Go checkpoint

- [x] verificar consistencia total
- [x] verificar readiness del baseline
- [x] registrar gaps reales
- [x] decidir freeze / seguir / abrir roadmap nuevo

## Snapshot de progreso actual

- [x] existe vertical slice mínima engine + surface
- [x] existe modo memory y modo postgres opcional para el core
- [x] existe preview/simulation mínima
- [x] existe approval release path básico
- [x] existe surface mínima de intake/proposal/inspection/maintenance
- [x] existe adapter Cerbos real opcional con tests básicos
- [x] existe persistencia postgres opcional de core y surface mínima
- [x] existe primera corrida real de tests en verde
- [x] existe compensation/unknown outcome materializados de forma mínima
- [x] existe test end-to-end del corredor completo en verde
- [ ] falta hardening del slice antes de freeze

## Interpretación operativa actual

- [x] Phase 1 cerrada en código
- [x] Phase 2 cerrada en modo mínimo
- [x] Phase 3 cerrada en modo mínimo
- [x] Phase 4 cerrada en un baseline técnico suficientemente materializado
- [x] Phase 5 evaluada con checkpoint inicial explícito
