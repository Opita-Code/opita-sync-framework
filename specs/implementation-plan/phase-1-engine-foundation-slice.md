# Phase 1 — Engine foundation slice

## Objetivo

Materializar el corredor mínimo real del engine sobre el baseline ya cerrado: **compiler path**, **runtime durable**, **Cerbos PEP/PDP integration**, **event log canónico** y **registry/resolution**. Esta fase existe para dejar operativo el núcleo mínimo del framework, no para ampliar breadth de producto.

## Qué entra en esta fase

- Compiler path mínimo desde intent/proposal gobernada hasta `compiled_contract` persistido y correlable.
- Runtime skeleton mínimo sobre Temporal con creación de `execution_record`, lifecycle básico y separación explícita entre execution y application.
- Integración mínima Cerbos PEP/PDP con input canonizado y mapping inicial a enforcement de ejecución.
- Event log canónico mínimo append-only en PostgreSQL con correlación básica de IDs y eventos principales del corredor.
- Registry/resolution mínimo para manifest, binding y resolución de provider/worker remoto dentro del baseline ya definido.
- Corredor integrado mínimo que conecte esos cinco seams sin introducir una segunda fuente de verdad.

## Qué queda explícitamente fuera

- Distribution layer, rollout, activation o instalación tenant-scoped.
- Surface conversacional rica, workspace de operador o UX final.
- Hardening integral, suites completas de regresión y baseline reusable final.
- Connector SDK completo, object storage operativo completo, Valkey operativo completo y retrieval plane completo.
- Escala enterprise final, optimizaciones profundas de performance, HA avanzada o breadth de observabilidad fuera del mínimo operativo.
- Cualquier relectura de Temporal, Cerbos, PostgreSQL, OTel/LGTM, OCI bundle o del pipeline Intent -> Change Proposal -> Governed Patchset.

## Qué valida del baseline y qué todavía no valida

### Valida en esta fase

- Que el baseline puede materializar un engine mínimo coherente sobre sus seams principales.
- Que el contrato compilado puede convertirse en ejecución durable gobernada.
- Que policy, evidencia y capability resolution pueden conectarse al corredor sin ambigüedad estructural.

### Todavía no valida

- Approvals/release completos, preview profundo y manejo exhaustivo de `unknown_outcome`.
- Surface operacional usable por operator/developer.
- Hardening reusable, demo final y readiness de freeze.

## Orden interno recomendado

1. **Compiler contract path**  
   Primero se fija la unidad compilada. Sin eso, todo lo demás nace sobre input ambiguo.

2. **Runtime durable mínimo**  
   Segundo se materializa el workflow durable que ejecuta ese contrato y crea evidencia primaria de lifecycle.

3. **Cerbos PEP/PDP integration mínima**  
   Tercero se conecta enforcement externo sobre inputs ya estabilizados por compilador y runtime.

4. **Event log canónico mínimo**  
   Cuarto se fija evidencia append-only y correlación antes de sumar más superficie.

5. **Registry/resolution mínimo**  
   Quinto se cierra la resolución controlada de capabilities ya sobre un engine que compila, ejecuta, decide y registra.

6. **Corredor integrado mínimo**  
   Sexto se prueba el slice completo como unidad de construcción real, no como cinco silos aislados.

## Entregables técnicos mínimos

- Definición operativa del compiler path y de sus artifacts mínimos de entrada/salida.
- Persistencia de `compiled_contract` y `compilation_report` con correlación suficiente para runtime y auditoría.
- Workflow durable mínimo con estados base, timers/retries básicos y camino explícito de execution/application.
- Boundary PEP/PDP implementable con payload canonizado hacia Cerbos y decisión retornable al runtime.
- Event log canónico mínimo con taxonomía inicial de eventos, correlación de IDs y redacción previa a exportación.
- Registry/resolution mínimo capaz de resolver manifest, binding aprobado y provider/worker remoto autorizado.
- Definición del corredor mínimo engine-only que une compilación, ejecución, policy, evidencia y resolution.

## Criterios de done

- Existe un recorrido verificable desde intent gobernada hasta `compiled_contract` usable por runtime.
- Existe una ejecución durable mínima que no colapsa execution y application en un solo estado implícito.
- La decisión de policy se consulta externamente a Cerbos y no queda incrustada ad hoc dentro del runtime.
- El event log registra evidencia canónica mínima y no delega la verdad primaria a observabilidad derivada.
- La resolución de capabilities ocurre por registry/binding/provider y no por acoplamiento directo hardcodeado.
- El slice deja probado el baseline del engine en un corredor mínimo real.
- Queda explícito qué gaps pasan a Phase 2 sin disfrazarlos de cierre completo.

## Riesgos

- Intentar cubrir demasiados casos laterales antes de cerrar el corredor mínimo y terminar con seams difusos.
- Mezclar policy con runtime por conveniencia y perder auditabilidad real.
- Tratar observabilidad derivada como sustituto del event log operativo.
- Resolver providers de forma ad hoc y romper el modelo manifest/binding/provider ya decidido.
- Confundir “mínimo implementado” con “baseline reusable listo”, adelantando conclusiones que todavía no corresponden.

## Dependencias

- Baseline normativo A-E ya cerrado.
- Convergencia OSF v1 ya cerrada como addendum técnico.
- Temporal como único runtime durable de verdad.
- Cerbos como PDP principal.
- PostgreSQL como verdad operativa durable; OTel/LGTM como observabilidad derivada.
- Declarative manifest + remote provider/worker model ya decidido.
- OCI bundle inmutable, firma y attachments ya decididos como artifact base.
