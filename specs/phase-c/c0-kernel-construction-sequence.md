# C.0 — Secuencia de construcción del kernel

## Principios base

Fase C construye el kernel, no reabre la arquitectura base. Eso implica cinco principios: un solo runtime durable de verdad, una sola memoria operativa gobernable, policy externa y auditable, seams explícitos antes que acoplamientos implícitos, e integración incremental con checkpoints verificables. El objetivo es reducir ambigüedad de implementación sin colapsar diseño, runtime, policy y observabilidad en una sola capa difusa.

## Qué significa “kernel” en Opyta Sync

En Opyta Sync, “kernel” es el conjunto mínimo de capacidades que convierte intención gobernada en ejecución durable y evidencia correlacionable. Incluye: compilación de intención a contrato ejecutable, creación y seguimiento de `execution_record`, enforcement de policy, emisión de event log operativo y resolución controlada de capabilities sobre artifacts inmutables. No incluye todavía shells conversacionales avanzadas, UX pesada, distribution layer ni endurecimiento final de operación.

## Dependencias heredadas de Fase A/B

Desde Fase A llegan los invariantes del dominio: objetos canónicos, contratos versionados, estados y eventos, approvals, clasificación, result types, idempotencia, retries, compensación y separación explícita entre ejecución y aplicación. Desde Fase B llega el baseline de implementación: Temporal como runtime durable, Cerbos como PDP, PostgreSQL como memoria operativa, OTel/LGTM como plano observability, manifest declarativo con provider/worker remoto como seam de extensibilidad, y OCI bundle firmado con pipeline Intent → Change Proposal → Governed Patchset como marco de artifacts gobernados. La distribution layer queda fuera del roadmap actual de Fase C.

## Orden de construcción recomendado y por qué

1. **C.1 Intent → Contract compiler**. Primero hay que fijar la unidad compilada que el resto del kernel va a ejecutar y gobernar. Sin contrato compilado estable, runtime, policy y observabilidad nacen sobre entradas ambiguas.
2. **C.2 Execution runtime skeleton**. Una vez fijado el artefacto compilado, se construye el esqueleto durable que materializa `execution_record`, estados, timers, retries y compensación sobre Temporal.
3. **C.3 Cerbos integration**. Con contrato y runtime mínimos definidos, se incorpora enforcement externo. Hacerlo antes arriesga policy inputs inestables; hacerlo después evita policy ad hoc incrustada en el runtime.
4. **C.4 Event log and observability base**. Cuando compilación, ejecución y policy ya tienen boundaries más estables, se fija la evidencia mínima y la proyección a observabilidad sin duplicar semántica.
5. **C.5 Capability registry and resolution**. Recién entonces conviene cerrar cómo se resuelven manifests, bundles, bindings y providers, porque ya existe un kernel básico capaz de compilar, ejecutar, decidir y evidenciar.
6. **C.6 Kernel integration checkpoint**. El cierre de Fase C no es feature breadth; es demostrar consistencia del camino mínimo end-to-end.

Este orden minimiza retrabajo porque estabiliza primero la unidad operativa (`compiled_contract`), luego el contenedor durable (`execution_record` sobre Temporal), después governance y evidencia, y recién al final catálogo/resolution.

## Seams del kernel que deben existir desde el día 1

- **Seam de compilación** entre intención/propuesta y contrato compilado versionado.
- **Seam de runtime** entre `execution_record` canónico y workflow durable real.
- **Seam de policy** entre PEP del kernel y PDP Cerbos.
- **Seam de evidencia** entre event log operativo y telemetría derivada.
- **Seam de capability resolution** entre manifest declarativo, binding aprobado y provider/worker ejecutable.
- **Seam de clasificación/redacción** antes de persistir o exportar señales fuera de memoria operativa.

Estos seams deben aparecer desde el primer tramo porque después son carísimos de extraer. Si no nacen explícitos, el kernel termina mezclando compilación, policy, runtime y catálogo en una sola capa imposible de gobernar.

## Estrategia de integración incremental

La integración debe avanzar por un corredor mínimo y estable: intención gobernada entra al compilador, el compilador emite un contrato compilado versionado y con fingerprint, el runtime crea una ejecución durable mínima sobre ese contrato, el runtime consulta policy con inputs canonizados, se registran eventos operativos correlacionados y el kernel resuelve una capability autorizada mediante registry/binding. Cada bloque agrega una capacidad nueva, pero debe conectarse al corredor anterior con evidencia verificable y sin introducir una segunda fuente de verdad.

La regla práctica es simple: cada bloque nuevo debe integrarse con el corredor mínimo antes de sumar casos laterales. Primero consistencia, después cobertura.

## Qué NO construir todavía en Fase C

- UX conversacional avanzada o workspace rico de operador.
- Capa declarativa superior que compita con Temporal como runtime real.
- Optimización profunda de analytics, reporting o búsquedas complejas.
- ReBAC profundo o modelos de sharing avanzados más allá del baseline de Cerbos.
- Packaging/distribution workflow completo; la distribution layer queda fuera del roadmap actual más allá del vínculo mínimo manifest-bundle-binding-provider.
- Overlays tenant capaces de reescribir invariantes del bundle o del contrato base.
- Automatizaciones de hardening, supply chain o escalado fino propias de Fase E.

## Criterios de aceptación de C.0

1. Existe una secuencia de construcción explícita, ordenada y justificada.
2. Queda claro qué artifact o boundary estabiliza cada bloque antes del siguiente.
3. Los seams del kernel están nombrados y tienen función operativa concreta.
4. La estrategia de integración incremental define un corredor mínimo end-to-end verificable.
5. Queda documentado qué trabajo se difiere deliberadamente a Fase D/E.
6. La secuencia propuesta respeta sin ambigüedad el baseline elegido en Fase B.
