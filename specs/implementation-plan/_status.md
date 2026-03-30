# Implementation Plan — Estado actual

## Objetivo del plan

Traducir el baseline reusable ya cerrado y la convergencia técnica de OSF en una secuencia de construcción real, priorizada y ejecutable para el **Opita Sync Framework (OSF)**, sin abrir un roadmap nuevo de producto.

## Estado general

- Estado: **cerrado como reusable implementation baseline v1**
- Próximo bloque recomendado: **ninguno — implementation plan actual cerrado**

## Fases del implementation plan

1. engine foundation slice
2. governed execution slice
3. surface operational slice
4. hardening and reusable baseline slice
5. go/no-go checkpoint

## Nota de alcance

Este implementation plan implementa el baseline ya cerrado; no lo redefine ni reabre decisiones duras del roadmap principal.

## Snapshot de progreso real

- Phase 1 está **materializada** en un primer slice ejecutable.
- Phase 2 está **cerrada en modo mínimo**: approvals/release, preview/simulation, smoke path y paths básicos de `unknown_outcome` / `compensation_pending` existen.
- Phase 3 está **materializada en modo mínimo** con intake, proposal, preview, inspection/recovery y maintenance candidate.
- Phase 4 está **cerrada**: existe regresión ejecutable del corredor, connector SDK baseline, artifact/cache/retrieval baseline, docs/playbooks mínimos y demo técnico reproducible.
- Existe una primera pasada real de `go test ./...` en verde sobre el slice actual.
- Existe adapter Cerbos real opcional con tests básicos y PostgreSQL opcional también para varios artifacts de surface.
- Existe documentación operativa mínima (`README`, runbook técnico y alcance de alfa técnica`) para handoff inicial del slice.
- Existe un connector SDK baseline inicial materializado en código con mock provider y tests básicos.
- Existe un demo de referencia reproducible y un checkpoint explícito de go/no-go.
- El implementation plan actual queda cerrado: el baseline reusable v1 se congela dentro del scope actual y cualquier trabajo siguiente debe abrirse como hardening/productización explícita o nuevo roadmap.

## Lectura de checkpoint

- Estado actual: **baseline reusable v1 congelado**.
- Madurez operativa: **alfa técnica**.
- Decisión: el implementation plan se considera cerrado; cualquier trabajo posterior debe tratarse como hardening/productización explícita y no como continuación implícita de estas fases.
