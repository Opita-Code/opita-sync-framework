# B.6 — Conversational configuration comparison

## Objetivo de la comparativa

Definir qué modelo conversacional debe adoptarse como baseline provisional para solicitar, preparar, revisar y gobernar cambios de configuración en Opyta Sync sin romper los invariantes ya cerrados en Fase A ni las decisiones provisionales de B.1-B.5.

La pregunta de B.6 NO es “qué chat se siente más cómodo”. La pregunta real es: **qué abstracción conversacional permite capturar intención humana/IA, transformarla en propuesta declarativa gobernada, simularla, aprobarla y aplicarla con evidencia durable, sin convertir la conversación en source of truth ni en apply directo a producción**.

---

## Qué exige Opyta Sync de la configuración conversacional

Derivado de Fase A y de las decisiones ya cerradas en B.1-B.5, la configuración conversacional debe sostener como mínimo estas exigencias:

- permitir que una persona o una IA expresen **solicitudes e intención** en lenguaje natural, pero que el resultado operativo sea siempre un **artefacto declarativo gobernado**;
- separar con claridad **conversación para solicitar/intencionar**, **propuesta de cambio revisable** y **activación tenant**;
- mantener authoring humano legible en **YAML** y una forma compilada/materializada en **JSON canónico determinístico**;
- producir **diff humano sobre YAML** y **diff material sobre JSON canónico** sobre el MISMO cambio propuesto;
- correr **preview, simulación, policy, approval y classification** sobre la **propuesta de cambio** y no sobre la charla informal;
- preservar la decisión de B.5: el **OCI bundle** sigue siendo el artifact baseline cuando el cambio impacta packaging/publicación, pero la conversación no sustituye artifact ni catálogo;
- sostener governance fuerte, approvals, reason codes, clasificación y evidencia multi-tenant;
- dejar audit trail verificable de **solicitud**, **propuesta**, **simulación**, **aprobación** y **apply**;
- impedir que el chat aplique directo a producción sin **artefacto + diff + preview + simulación + aprobación**;
- preparar una superficie IA-friendly para Fase D, donde diff, preview, simulación y explicación sean first-class y no accesorios tardíos.

En otras palabras: Opyta Sync no necesita “editar YAML con un chatbot”. Necesita una **tubería gobernada de cambio** donde la conversación inicia o colabora, pero la verdad ejecutable vive en objetos declarativos versionables y auditables.

---

## Invariantes no negociables antes de B.6

- El **source of truth** sigue siendo declarativo y gobernado; **la conversación no es la verdad del sistema**.
- El authoring sigue siendo **YAML** y la forma normalizada/compilada sigue siendo **JSON canónico determinístico**.
- El **OCI bundle** sigue siendo el artifact baseline cuando corresponde empaquetar/publicar contenido gobernado.
- **Packaging/publicación** sigue separado de **activación tenant**.
- La activación tenant sigue siendo un **objeto/acto gobernado aparte** y no un side effect del chat.
- Deben existir seams explícitos para **policy, approvals, classification, simulation, evidence y apply**.
- Fase D va a exigir **diff/preview, simulación y operación IA-friendly**, por lo que B.6 no puede elegir un modelo que esconda la semántica real detrás de una caja negra conversacional.
- La conversación puede ayudar a redactar o refinar cambios, pero **nunca reemplaza** el mecanismo formal de propuesta, revisión y aplicación.

---

## Candidatos evaluados y por qué entran al shortlist

### 1. Intent → Change Proposal → Governed Patchset

Entra al shortlist porque separa explícitamente tres capas con semántica distinta:

1. **intent request**: solicitud inicial en lenguaje natural o structured intent liviano;
2. **change proposal**: materialización revisable del cambio propuesto con diff, preview y simulación;
3. **governed patchset**: patch declarativo aprobado y aplicable sobre artifacts/config gobernada.

Es fuerte porque equilibra muy bien UX conversacional con governance dura: la conversación sirve para pedir y refinar, pero el sistema aterriza en un objeto de cambio verificable antes de tocar configuración efectiva.

### 2. Conversational Workspace / PR-style Draft

Entra al shortlist porque maximiza seguridad operativa y familiaridad humana: el cambio vive en un workspace o draft estilo PR, con comentario, review, diff, approvals y merge/apply explícitos.

Es especialmente atractivo si la prioridad dominante es evitar toda ambigüedad entre “hablar de un cambio” y “tener un cambio listo para gobernanza”, aun a costa de mayor fricción y más pasos visibles.

### 3. Schema-guided Structured Editor behind Chat

Entra al shortlist porque promete mucha seguridad sintáctica: el chat opera detrás de un editor guiado por schema, con campos válidos, constraints explícitos y menor riesgo de producir configuraciones mal formadas.

Es un candidato serio cuando el problema principal es la validez estructural de formularios complejos. Sin embargo, entra con sospecha: puede resolver MUY bien edición asistida, pero no necesariamente resuelve mejor el lifecycle completo de propuesta, simulación, evidence y apply gobernado.

---

## Qué enfoques quedan como mecanismos internos/complementarios y no como modelo principal

### Canonical JSON AST / JSON Patch as Control Plane

Queda como **mecanismo interno/complementario** y no como modelo principal porque sirve muy bien como plano de control material para diffing, merge, apply determinístico, simulación y trazabilidad técnica. Pero NO debería ser la experiencia primaria para humanos o para la capa conversacional.

Su lugar correcto es detrás del modelo principal:

- como representación canónica de cambios sobre JSON normalizado;
- como soporte para diff material, patch validation, drift detection y replay;
- como formato técnico consumido por simulación, policy y apply engine.

### Plan-and-Simulate conversational orchestrator

Queda como **mecanismo interno/complementario** y no como modelo principal porque aporta muchísimo valor para exploración, reasoning y simulación previa, pero no debe convertirse en la unidad formal de cambio.

Su lugar correcto es:

- ayudar a transformar intención en propuesta;
- sugerir alternativas, impactos y pasos;
- ejecutar simulaciones conversacionales sobre una propuesta formal;
- explicar resultados y riesgos.

Lo que NO debe hacer es reemplazar proposal/diff/approval/apply con un “agente que decide y ejecuta”. Opyta Sync necesita gobernanza operable, no magia opaca.

---

## Hard gates aplicados

Los hard gates vienen de B.0 y se aplican antes del scoring.

| Hard gate | Intent → Change Proposal → Governed Patchset | Conversational Workspace / PR-style Draft | Schema-guided Structured Editor behind Chat | Lectura |
|---|---|---|---|---|
| soporte real para multi-tenant | **pass** | **pass** | **pass parcial** | Los tres pueden operar tenant-scoped, pero el schema-guided depende más de cómo se modele el lifecycle externo que del editor en sí. |
| trazabilidad auditable | **pass** | **pass** | **pass con tensión** | Intent→Proposal→Patchset y PR-style dejan evidencia natural por etapas; schema-guided necesita capas adicionales para dejar historia rica más allá de eventos de edición. |
| hooks o seams para approvals / classification / policy | **pass** | **pass** | **pass parcial** | Los dos primeros hacen natural colgar simulación y governance sobre la propuesta/draft; schema-guided tiende a concentrarse en validez estructural más que en workflow gobernado. |
| estado durable inspeccionable | **pass** | **pass** | **pass con tensión** | Proposal y draft pueden persistirse como objetos durables; schema-guided suele exponer mejor estado del formulario que estado del cambio gobernado completo. |
| ausencia de cajas negras imposibles de operar | **pass** | **pass** | **pass con tensión** | Los dos primeros conservan artefactos explícitos y diffables; schema-guided corre riesgo de esconder demasiada semántica en lógica UI/schema orchestration. |

### Lectura de gates

- **Intent → Change Proposal → Governed Patchset** pasa el filtro estructural con el mejor balance entre conversación usable y cambio formal gobernado.
- **Conversational Workspace / PR-style Draft** también pasa con solidez: es probablemente el más conservador y el más fácil de explicar operacionalmente.
- **Schema-guided Structured Editor behind Chat** no queda descartado del todo, pero entra tensionado porque su fortaleza principal está en la edición validada, NO en el modelo de evidencia, proposal lifecycle y apply gobernado end-to-end.

---

## Evaluación comparativa usando criterios/pesos de B.0

Se usa la base común de B.0 con escala 1-5. El score ponderado se normaliza a 0-100.

### Tensiones específicas de esta B.6

Aunque B.0 fija el marco principal, en esta comparativa el peso real se concentra en seis tensiones:

1. si el modelo separa con nitidez **intención conversacional**, **propuesta de cambio efectiva** y **activación tenant**;
2. si preview, simulación, policy, approval y classification corren sobre la **propuesta** y no sobre la charla informal;
3. si el cambio puede producir **diff humano sobre YAML** y **diff material sobre JSON canónico** de forma coherente;
4. si la evidencia de **request → proposal → simulation → approval → apply** queda persistida sin inventar reconstrucciones frágiles;
5. si la UX sigue siendo suficientemente buena para humanos e IA sin sacrificar governance;
6. si el baseline elegido deja espacio a Fase D para operación IA-friendly sin mover el source of truth fuera de artifacts declarativos gobernados.

---

## Tabla de scoring comparativo

| Criterio | Peso | Intent → Change Proposal → Governed Patchset | Conversational Workspace / PR-style Draft | Schema-guided Structured Editor behind Chat | Nota breve |
|---|---:|---:|---:|---:|---|
| alineación con invariantes de Fase A | 20% | 5 | 4 | 3 | Intent→Proposal→Patchset protege mejor source of truth declarativo y separación entre charla, proposal y activation. |
| gobernanza y seguridad multi-tenant | 18% | 5 | 5 | 3 | Los dos primeros hacen natural approvals, classification y boundaries tenant-scoped; schema-guided no resuelve eso como fortaleza central. |
| durabilidad, auditabilidad y recuperación | 16% | 5 | 5 | 3 | Proposal/patchset y draft/PR dejan historia durable y recuperable; schema-guided necesita más capas para auditar intención y revisión. |
| encaje con event log / trazabilidad / evidencia | 12% | 5 | 4 | 3 | Intent→Proposal→Patchset hace muy explícita la cadena de evidencia; PR-style también funciona bien, pero con más ruido conversacional alrededor del draft. |
| extensibilidad futura sin romper el core | 10% | 4 | 4 | 3 | Los dos primeros dejan buen seam para simulación, AI assist y distintos UIs sobre el mismo modelo gobernado. |
| operabilidad IA-friendly | 9% | 5 | 4 | 3 | Intent→Proposal→Patchset permite pedir en lenguaje natural pero operar sobre objetos explícitos; PR-style es fuerte pero más friccional; schema-guided limita flexibilidad conversacional. |
| complejidad de implementación y riesgo de entrega en Fase C | 8% | 4 | 3 | 3 | Intent→Proposal→Patchset tiene complejidad razonable y captura valor pronto; PR-style exige más superficie/workspace/review UX; schema-guided exige más UI/schema orchestration. |
| lock-in / costo de reversión futura | 4% | 4 | 4 | 3 | Los dos primeros son portables conceptualmente; schema-guided puede acoplar demasiado la semántica a un editor particular. |
| madurez de ecosistema y tooling | 3% | 4 | 5 | 4 | PR-style reaprovecha patrones muy maduros de draft/review/diff; intent→proposal también se apoya en primitives conocidas; schema-guided tiene tooling, pero menos alineado al problema completo. |

### Score ponderado

| Candidato | Score ponderado | Veredicto |
|---|---:|---|
| Intent → Change Proposal → Governed Patchset | **95.0 / 100** | **preferred** |
| Conversational Workspace / PR-style Draft | **85.6 / 100** | **acceptable_with_tradeoffs** |
| Schema-guided Structured Editor behind Chat | **61.0 / 100** | **reject** |

### Lectura del resultado

La diferencia importante NO está en que PR-style sea malo. De hecho, es un candidato MUY serio y operacionalmente seguro. Pero B.6 no busca solo seguridad por fricción; busca el mejor balance entre UX conversacional, proposal formal, diff/simulación y evidencia gobernada.

Bajo ese criterio, **Intent → Change Proposal → Governed Patchset** queda mejor posicionado porque preserva una experiencia conversacional fuerte SIN convertir el chat en editor efectivo ni en apply engine. El modelo PR-style sigue siendo valioso y puede sobrevivir como una vista o modalidad de review del mismo pipeline.

---

## Tradeoffs narrativos por candidato

### Intent → Change Proposal → Governed Patchset

**Fortalezas**

- Separa con claridad tres actos diferentes:
  1. **solicitar/intencionar** un cambio;
  2. **proponer** un cambio formalmente revisable;
  3. **aplicar** un patchset gobernado ya aprobado.
- Hace natural que el chat ayude a capturar intención, recopilar contexto y refinar objetivos sin tocar todavía la configuración efectiva.
- Permite que la propuesta produzca simultáneamente:
  - **diff humano sobre YAML** para revisión comprensible;
  - **diff material sobre JSON canónico** para simulación, policy, apply y evidencia técnica.
- Encaja perfecto con una cadena de evidence bien gobernada:
  - request,
  - proposal,
  - simulation,
  - approval,
  - apply.
- Hace muy explícito que **policy/approval/classification/simulation** corren sobre la **change proposal** o sobre el **governed patchset**, nunca sobre la charla.
- Es IA-friendly sin perder disciplina: la IA puede ayudar a construir y explicar, pero el sistema opera sobre artefactos declarativos versionables.

**Tradeoff principal**

- Requiere diseñar bien la transición entre intent y proposal. Si esa frontera queda difusa, el sistema corre el riesgo de generar “propuestas implícitas” demasiado temprano o de esconder heurísticas de traducción. La buena noticia es que ese riesgo se controla con un contract claro de proposal object.

### Conversational Workspace / PR-style Draft

**Fortalezas**

- Es el candidato más fácil de explicar a nivel operativo: un draft/workspace estilo PR concentra diff, comentarios, review, approvals y apply final.
- Refuerza MUCHO la seguridad operacional porque hace muy visible que hablar, proponer y aplicar son pasos distintos.
- Reaprovecha convenciones humanas maduras: review threads, requested changes, approval states, merge/apply semantics, snapshots del draft.
- Es excelente para entornos donde la fricción adicional es aceptable a cambio de máxima inspección humana previa.

**Tradeoff principal**

- La seguridad extra viene con fricción extra. Si todo cambio conversacional obliga demasiado temprano a abrir y gestionar un workspace pesado, la UX puede volverse rígida para iteración rápida y para experiencias IA-friendly más fluidas. El peligro es sobrediseñar el review shell antes de optimizar el modelo de proposal subyacente.

### Schema-guided Structured Editor behind Chat

**Fortalezas**

- Reduce errores de forma y mejora la edición de estructuras complejas con validaciones tempranas.
- Puede ser excelente como UI complementaria para campos de alto riesgo, overlays limitados o configuraciones muy estructuradas.
- Hace más difícil producir payloads mal tipados o inconsistentes con schema.

**Tradeoff principal**

- Confunde fácilmente el problema de **editar bien un documento** con el problema real de Opyta Sync, que es **gobernar el lifecycle completo del cambio**. Validez estructural NO alcanza. Sin un proposal model fuerte, evidencia, simulación y approvals, el editor guiado queda corto como baseline principal.

---

## Riesgos y dudas abiertas por candidato

### Intent → Change Proposal → Governed Patchset

- Definir con precisión el **schema del proposal object**: qué campos mínimos registra de intent, alcance, artifacts afectados, tenant scope, clasificación, simulación y apply plan.
- Definir cómo se representa el **governed patchset**: si como patch declarativo independiente, como manifest de cambios, o como combinación de YAML diff + canonical JSON patch derivado.
- Definir reglas de transición entre estados: borrador, simulated, approved, rejected, expired, applied, superseded.
- Confirmar cómo se correlaciona evidencia entre proposal, approvals, runtime events y apply outcome en Temporal/event log.

### Conversational Workspace / PR-style Draft

- Definir si el workspace es el objeto principal o solo una vista/review shell sobre una proposal subyacente.
- Evitar que el draft absorba demasiada semántica y termine reemplazando el modelo formal de patchset.
- Confirmar cómo se modelan cambios pequeños o iteraciones rápidas sin crear una UX demasiado pesada.
- Definir cómo se separa con precisión comentario conversacional informal de evidencia formal del cambio aprobado.

### Schema-guided Structured Editor behind Chat

- Confirmar si puede producir evidencia rica de proposal/review/simulation o si solo mejora la edición puntual.
- Evitar acoplar la semántica del sistema a un UI/editor específico y dificultar otros clientes o agentes.
- Definir cómo soportaría diffs narrativos y materiales sin depender de una sesión de edición viva.
- Validar si realmente agrega valor suficiente como baseline o si conviene reservarlo para configuraciones de alto constraint como mecanismo complementario.

---

## Recomendación final de B.6

**Intent → Change Proposal → Governed Patchset** como baseline provisional de configuración conversacional.

Es la recomendación más sólida porque logra el mejor equilibrio entre **UX conversacional** y **governance fuerte**. Permite pedir cambios en lenguaje natural, refinar intención con ayuda humana o de IA, y aun así aterrizar SIEMPRE en una propuesta formal con diff, preview, simulación, approvals y apply gobernado.

### Qué significa exactamente esta recomendación

Esto tiene que quedar CRISTALINO:

1. **conversación para solicitar/intencionar** ≠ **conversación para editar config efectiva** ≠ **activación tenant**;
2. el chat puede originar o refinar un pedido, pero **nunca aplica directo a producción**;
3. todo cambio real pasa por **artefacto/propuesta + diff + preview + simulación + aprobación + apply**;
4. debe coexistir **diff humano sobre YAML** con **diff material sobre JSON canónico**;
5. **policy / approval / classification / simulation** corren sobre la **propuesta de cambio** o su patchset derivado, no sobre la charla informal;
6. debe guardarse evidencia durable de **solicitud, propuesta, simulación, aprobación y apply**;
7. la **activación tenant** sigue siendo una capa/acto separado, gobernado y auditado.

### Por qué no se recomienda Conversational Workspace / PR-style Draft como baseline principal

Porque, aunque es muy seguro y muy serio, introduce una fricción operativa mayor como punto de partida. Y B.6 necesita una baseline que no sacrifique innecesariamente la fluidez conversacional que Fase D va a querer explotar.

Dicho brutalmente y con cariño: si arrancás con PR-style como modelo principal, corrés el riesgo de diseñar primero el cascarón de review y recién después el modelo real de propuesta. Opyta Sync necesita al revés: **primero proposal object y patchset gobernado; después, todas las vistas y workspaces que quieras arriba de eso**.

### En qué sentido PR-style Draft sigue siendo importante

Que no sea el baseline principal NO significa descartarlo. Al contrario: puede ser una **modalidad de revisión y colaboración** excelente sobre el mismo pipeline recomendado. O sea, la recomendación de B.6 no mata PR-style; lo reubica en la capa correcta.

---

## Decisión provisional y justificación

**Decisión provisional:** adoptar **Intent → Change Proposal → Governed Patchset** como baseline provisional para configuración conversacional en Fase C.

### Justificación

1. **B.0 obliga a priorizar hard gates, auditabilidad y costo de reversión antes que gusto de UX local.**
2. Este modelo preserva el principio central del proyecto: **source of truth declarativo gobernado, no la conversación**.
3. Separa correctamente **intent**, **proposal**, **patch/apply** y **tenant activation**, evitando colapsar actos distintos en una sola interacción difusa.
4. Hace natural que **simulación, policy, approvals y classification** corran sobre objetos formales de cambio.
5. Permite coexistencia limpia entre **YAML authoring** y **JSON canónico** como control plane material.
6. Deja una UX muy fuerte para humanos e IA sin habilitar atajos peligrosos de “chat -> producción”.
7. Permite que un **workspace/PR-style review** exista como capa superior opcional, sin volverla el centro semántico del sistema.

La razón de fondo es simple: **equilibra mejor UX y governance**. No pide la fricción máxima desde el minuto cero, pero tampoco cae en la fantasía irresponsable de que conversar ya equivale a cambiar sistema. Esa disciplina importa MUCHÍSIMO.

---

## Qué supuestos quedan pendientes para Fase C/D

- Fase C debe definir el **proposal object model** exacto: identidad, estado, tenant scope, target artifacts, clasificación, approvals requeridos, simulation snapshot y apply outcome.
- Fase C debe definir el formato exacto del **governed patchset** y su relación con YAML diff y canonical JSON diff/patch.
- Fase C debe definir cómo se materializa el **preview**: render humano, impacto operativo, riesgos, compatibilidad y efectos sobre activaciones tenant.
- Fase C debe definir cómo se correlaciona la evidencia en runtime/event log para la cadena **request → proposal → simulation → approval → apply**.
- Fase C debe definir si el **PR-style workspace** se implementa desde el inicio como vista oficial sobre la proposal o si queda como evolución cercana.
- Fase D debe validar patrones IA-friendly para explicación, diff asistido, simulación guiada y troubleshooting sin mover el source of truth fuera del proposal/patchset.
- Sigue abierto dónde conviene usar **schema-guided editing** como mecanismo complementario: overlays limitados, formularios críticos, configuraciones de alto constraint o assistants especializados.

---

## Criterios de aceptación de B.6

B.6 puede considerarse cerrado si se cumplen todas estas condiciones:

1. existe shortlist explícito de candidatos comparados;
2. queda explícita la diferencia entre **solicitud/intención conversacional**, **propuesta de cambio** y **activación tenant**;
3. queda explícito que el **chat nunca aplica directo a producción**;
4. queda explícito que todo cambio real requiere **artefacto/propuesta + diff + preview + simulación + aprobación**;
5. se documenta la coexistencia obligatoria de **diff humano sobre YAML** y **diff material sobre JSON canónico**;
6. se documenta que **policy, approval, classification y simulation** corren sobre la **propuesta de cambio** y no sobre la conversación informal;
7. se documenta que debe persistirse evidencia de **solicitud, propuesta, simulación, aprobación y apply**;
8. se aplican hard gates de B.0 antes del scoring;
9. se usa la tabla de criterios/pesos comunes de B.0;
10. hay scoring comparativo visible y veredicto por candidato;
11. existe una recomendación final explícita y una decisión provisional clara;
12. queda documentado por qué **schema-guided structured editor** no se adopta como modelo principal;
13. quedan identificados los supuestos que deben cerrarse en Fase C/D sin reabrir innecesariamente la decisión base de B.6.
