# Phase 4 — Hardening and reusable baseline slice

## Objetivo

Endurecer lo construido en Phases 1-3 hasta dejar un **baseline reusable** real del OSF: **regresión integral**, **connector SDK baseline inicial**, **object storage / Valkey / retrieval baseline**, **docs/playbooks mínimos** y **demo ejecutable de referencia**. Esta fase cierra calidad operativa, no expande el producto.

## Qué entra

- Full regression baseline ejecutable sobre engine + surface ya materializados.
- Connector SDK baseline inicial compatible con el modelo manifest/binding/provider ya cerrado.
- Object storage / Valkey / retrieval baseline integrados como complementos controlados del implementation profile.
- Documentación y playbooks mínimos operativos para developer, operator y debugging.
- Demo ejecutable de referencia que pruebe el corredor reusable sin inflar scope.

## Qué queda fuera

- Distribution layer, rollout, activation o industrialización tenant-scoped.
- Escala enterprise final, despliegues multiregión completos o performance tuning exhaustivo.
- Nuevos seams de arquitectura o rediscusión del baseline duro.
- Breadth adicional de producto no necesario para probar el baseline reusable.

## Qué valida del baseline y qué todavía no valida

### Valida en esta fase

- Que el baseline puede sostener regresión y operación mínima de forma reusable.
- Que los complementos del implementation profile entran sin crear segundas verdades.
- Que existe evidencia suficiente para handoff, reuse y freeze técnico.

### Todavía no valida

- Roadmap nuevo de producto.
- Expansiones fuera del boundary engine/surface.
- Distribution y consumo tenant-scoped.

## Regresión integral

- Debe cubrir el corredor mínimo completo de engine + surface ya construido.
- Debe incluir invariantes de contrato, runtime, policy, approvals, evidence trail, preview y recovery mínimos.
- Debe detectar drift estructural del baseline antes de declararlo reusable.
- La prioridad es consistencia de seams y corredores, no vanity coverage.

## Connector SDK baseline inicial

- Debe fijar el baseline mínimo para providers remotos compatibles con manifest/binding/provider.
- Debe dejar claro qué contratos, lifecycle y evidencia espera el engine de un connector.
- Debe ser suficiente para interoperabilidad inicial, no para cubrir todas las familias de connectors imaginables.

## Object storage / Valkey / retrieval baseline

- Object storage entra como plane persistente de artifacts/evidence complementario al truth plane operativo.
- Valkey entra como plano efímero de aceleración, hints y locks cortos; no como verdad durable.
- Retrieval entra como plano complementario de búsqueda/corpus; nunca reemplaza PostgreSQL ni el event log canónico.
- Todo esto debe integrarse sin mover autoridad sobre runtime, policy o evidencia primaria.

## Docs / playbooks mínimos

- Deben explicar cómo operar, diagnosticar y extender el baseline sin reinterpretar arquitectura base.
- Deben cubrir al menos developer flow, operator flow, recovery básico, debugging semántico y límites del sistema.
- Deben escribirse en lenguaje implementable, sin marketing ni promesas de scale final.

## Demo ejecutable de referencia

- Debe demostrar el corredor reusable principal con artifacts, governance, preview, ejecución e inspección mínimas.
- Debe ser acotado y probatorio.
- No debe usarse como excusa para sumar breadth de producto o distribution encubierta.

## Criterios de done

- Existe una regresión integral ejecutable que protege el corredor reusable mínimo.
- Existe un connector SDK baseline inicial coherente con el modelo del engine.
- Object storage, Valkey y retrieval quedan integrados como planos complementarios controlados.
- Existen docs y playbooks mínimos suficientes para operar y reutilizar el baseline.
- Existe un demo ejecutable de referencia que demuestra el baseline sin reabrir scope.
- Queda explícito qué evidencia habilita pasar al checkpoint de go/no-go.

## Riesgos

- Confundir hardening con expansión funcional.
- Empaquetar como reusable algo que todavía depende de conocimiento tácito no documentado.
- Dar autoridad excesiva a object storage, Valkey o retrieval y romper el split-plane ya decidido.
- Diseñar un SDK demasiado amplio demasiado pronto.
- Usar el demo como vitrina de producto en vez de como prueba del baseline.
