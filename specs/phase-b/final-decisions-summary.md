# Fase B — Decisiones finales consolidadas

## Objetivo de Fase B y problema que resolvió

Fase B existió para cerrar las decisiones estructurales de mayor costo de reversión antes de construir el kernel en Fase C. El problema a resolver era claro: evitar entrar a la implementación del core con ambigüedad en runtime durable, policy, memoria operativa, observabilidad, extensibilidad, packaging y configuración conversacional.

El resultado de la fase no es código ni despliegue, sino un baseline arquitectónico coherente para construir el kernel sin reabrir comparativas centrales durante la ejecución.

## Resumen ejecutivo de decisiones finales B.1-B.6

- **B.1 Durable runtime:** se consolida el baseline durable ya cerrado en B.1, con **Temporal** como runtime de referencia para ejecución durable, auditabilidad, retries, pausas, compensación y operación multi-tenant.
- **B.2 Policy engine:** se consolida **Cerbos** como PDP baseline por su mejor equilibrio entre policy contextual, auditabilidad, multi-tenant governance y operabilidad en Fase C.
- **B.3 Operational memory + telemetry:** se consolida el split-plane **PostgreSQL + OpenTelemetry + Grafana LGTM**, donde PostgreSQL sostiene memoria operativa durable y LGTM absorbe señales analíticas de ejecución.
- **B.4 Extensibility model:** se consolida **Declarative manifest + remote provider/worker model** como seam principal para separar capability gobernada de handler ejecutable.
- **B.5 Capability packaging:** se consolida **OCI bundle inmutable + firma + attachments** como artifact base, manteniendo publicación/distribución y activación tenant como capas separadas.
- **B.6 Conversational configuration:** se consolida **Intent → Change Proposal → Governed Patchset** como pipeline conversacional gobernado, donde la conversación inicia o refina cambios, pero nunca reemplaza la propuesta formal ni el apply controlado.

## Tabla consolidada

| bloque | decisión | por qué ganó | tensión/tradeoff abierto | impacto directo en Fase C |
|---|---|---|---|---|
| B.1 | **Temporal** como baseline de durable runtime | Fue la opción con mejor encaje para durabilidad, auditabilidad, event history y multi-tenant operable | La capa conversacional/declarativa futura no debe competir con el runtime ni crear una segunda verdad de flujo | El kernel debe modelar ejecución, pausas, approvals, retries, compensación y evidencia sobre Temporal |
| B.2 | **Cerbos** como baseline de policy engine | Equilibra mejor PDP externo, policy contextual, auditabilidad y tenant-scoping sin empujar una mini-plataforma propia | Riesgo abierto si el dominio evoluciona hacia ReBAC profundo como necesidad central | Fase C debe integrar enforcement, approvals y decisiones auditables sin policy ad hoc dispersa |
| B.3 | **PostgreSQL + OTel + Grafana LGTM** como baseline split-plane | Separa correctamente memoria operativa estable de observabilidad analítica y conserva correlación auditable | Hay que fijar redacción/clasificación y qué evidencia se proyecta al plano analítico sin competir con memoria | Fase C debe construir memoria operativa en PostgreSQL e instrumentación/telemetría sobre OTel y LGTM |
| B.4 | **Declarative manifest + remote provider/worker model** | Preserva catálogo gobernado, separación capability/handler y compatibilidad con Temporal como único runtime de verdad | Si más adelante se quiere composición declarativa más rica, deberá vivir en una capa superior sin romper este seam | Fase C debe definir manifests, bindings, resolución de providers/workers y boundaries de ejecución |
| B.5 | **OCI bundle inmutable + firma + attachments** | Fija identidad fuerte del artifact por digest, firma y promoción estable entre ambientes | Publication, install plan y tenant activation deben modelarse arriba del artifact y no mezclarse con él | Fase C debe definir bundle schema, metadata, compatibilidad y objeto de tenant activation separado |
| B.6 | **Intent → Change Proposal → Governed Patchset** | Es el mejor equilibrio entre UX conversacional y governance fuerte con diff, simulación, aprobación y apply | La transición intent → proposal debe quedar explícita para no generar cambios implícitos ni chat-to-prod | Fase C debe definir proposal object, governed patchset, preview, simulación y correlación de evidencia |

## Compatibilidad entre decisiones

Estas decisiones NO compiten entre sí; encajan por capas y refuerzan el mismo modelo operativo:

1. **Durable runtime**: Temporal sostiene la ejecución durable real del kernel. Es la capa donde viven pausas, retries, compensación, approvals pendientes y evidencia de ejecución.
2. **Policy engine**: Cerbos se monta como PDP externo sobre ese runtime. No reemplaza estado ni flujo; decide autorización, enforcement y governance contextual sobre artifacts, operaciones y activaciones.
3. **Operational memory + telemetry**: PostgreSQL guarda memoria operativa estable y gobernable; OpenTelemetry + Grafana LGTM reciben señales derivadas de runtime, policy y operación para observabilidad. La telemetría explica; no gobierna ni reemplaza memoria.
4. **Extensibility model**: los manifests declarativos describen capabilities gobernadas; el runtime resuelve providers/workers remotos para ejecutar. Esto evita mezclar catálogo de negocio con código suelto y conserva a Temporal como único runtime durable.
5. **Capability packaging**: el OCI bundle transporta manifests, contracts, schemas, attachments y referencias versionadas como artifact base inmutable. No reemplaza bindings ni activaciones tenant; los estabiliza.
6. **Configuración conversacional**: la conversación genera intención; esa intención aterriza en una change proposal y luego en un governed patchset. Ese patchset modifica artifacts/configuración gobernada sin puentear policy, approvals, packaging ni runtime.

En síntesis: **la conversación produce cambios gobernados; el packaging los vuelve artifacts trazables; la extensibilidad define qué capability existe y cómo se resuelve; Cerbos gobierna decisiones; PostgreSQL conserva verdad operativa; LGTM explica la ejecución; Temporal ejecuta el kernel durable**.

## Riesgos abiertos aceptados antes de Fase C

- La capa declarativa/conversacional futura debe compilar o proyectarse sobre Temporal sin crear una segunda semántica de workflow.
- Cerbos deja abierto el riesgo de encaje futuro de **ReBAC profundo** si ownership, delegation y sharing se vuelven el problema dominante.
- Debe cerrarse el patrón exacto de correlación entre `tenant_id`, `execution_id`, `contract_id`, decisiones de policy y señales de observabilidad.
- Debe definirse el límite entre memoria operativa, event/runtime evidence y telemetría derivada para evitar duplicación semántica.
- Debe fijarse la política de compatibilidad entre `capability version`, `contract version`, `binding version` y `provider runtime version`.
- Debe definirse el proposal object y el governed patchset con suficiente precisión para que intent no colapse en apply implícito.

## Riesgos NO aceptables para entrar en Fase C

- Entrar a Fase C sin respetar que **Temporal es el único runtime durable de verdad**.
- Mezclar memoria operativa con observabilidad y usar telemetry como pseudo-source-of-truth.
- Permitir que extensiones o providers bypasseen Cerbos, approvals, clasificación o boundaries tenant-scoped.
- Tratar publication/install workflow como si fuera el package artifact base.
- Permitir overlays tenant que reescriban invariantes del bundle base.
- Habilitar chat-to-prod o apply directo sin proposal, diff, simulación, aprobación y apply gobernado.
- Reabrir comparativas B.1-B.6 dentro de la implementación del kernel salvo evidencia nueva extraordinaria y explícita.

## Recomendaciones concretas para arrancar Fase C

1. Definir el **kernel contract map** que conecte Temporal, Cerbos, PostgreSQL, OTel/LGTM y el catálogo de capabilities con IDs canónicos y correlación estable.
2. Diseñar el **proposal object model** y el **governed patchset model** como objetos first-class antes de construir interfaces conversacionales más ricas.
3. Definir el **manifest schema** de capability y el contrato exacto entre manifest, binding, provider/worker y runtime execution.
4. Definir el **OCI bundle layout**: qué va dentro del artifact, qué queda como attachment y qué queda como referencia externa gobernada.
5. Diseñar el **tenant activation object** como acto separado, auditable, aprobable y reversible que referencia bundle digest y bindings aprobados.
6. Fijar reglas de **instrumentación, clasificación y redacción** para que runtime, policy y extensiones emitan evidencia correlacionable sin filtrar datos sensibles.
7. Implementar primero seams del kernel y control plane; dejar para después shells UX avanzadas o workspaces conversacionales pesados.

## Qué queda deliberadamente diferido a Fase D/E

### Diferido a Fase D

- UX conversacional avanzada, asistentes IA-friendly, explicación guiada, troubleshooting asistido y vistas tipo PR/workspace sobre proposal/patchset.
- Posible capa declarativa superior o composición más rica de workflows, siempre que compile al baseline ya cerrado.
- Uso complementario de schema-guided editing para configuraciones de alto constraint.

### Diferido a Fase E

- Hardening operativo, optimización de supply chain, refinamientos de compatibilidad avanzada y evaluación de evolución hacia split package graph formal si el dominio lo exige.
- Revisión de necesidades de analytics más pesadas o de telemetría avanzada que pudieran justificar una evolución posterior del plano observability.
- Reevaluaciones mayores por crecimiento real del dominio, no por preferencia teórica anticipada.

## Criterios de aceptación del cierre de Fase B

Fase B se considera cerrada a nivel de decisión provisional si se cumplen todas estas condiciones:

1. B.0-B.6 están cerrados documentalmente y sin comparativas pendientes dentro de la fase.
2. Existe un resumen consolidado que no contradice las decisiones ya cerradas por bloque.
3. Queda explícito que el baseline arquitectónico para Fase C es coherente entre runtime, policy, memoria, observabilidad, extensibilidad, packaging y conversación.
4. Quedan identificados riesgos abiertos aceptados y riesgos no aceptables de entrada a Fase C.
5. Se declara como próximo bloque recomendado **Fase C — Construcción del kernel**.
6. Se mantiene la naturaleza provisional de las decisiones a nivel de implementación futura, pero se cierra la fase a nivel de baseline arquitectónico.
