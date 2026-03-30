# B.1 — Durable runtime comparison

## Objetivo de la comparativa

Definir qué runtime durable debe usarse como baseline del kernel en Fase C para ejecutar workflows con estado durable, reintentos, pausas, compensación, trazabilidad auditable y operación multi-tenant sin romper los invariantes ya cerrados en Fase A.

---

## Qué exige Opyta Sync de un runtime durable

Derivado de Fase A, el runtime elegido debe poder sostener, como mínimo, estas exigencias estructurales:

- representar `execution_record` como verdad operacional durable y auditable;
- separar explícitamente `execution_completed` de `application_completed`;
- soportar pausas, espera de approvals, bloqueos por governance, fallas técnicas y compensación;
- dejar historia reconstruible de eventos, evidencia y correlación estable por `tenant_id`, `environment`, `contract_id`, `execution_id` y snapshots;
- permitir deduplicación, retries controlados, resume seguro y manejo de `unknown outcome`;
- operar con aislamiento multi-tenant real, porque el tenant es frontera canónica de datos, governance, approvals, memoria y operación;
- dejar seams claros para policy, approvals, clasificación, memoria operativa y telemetría;
- ser operable por humanos e IA, con inspección de estado e historial sin cajas negras imposibles de explicar.

En otras palabras: Opyta Sync no necesita solo “orquestar pasos”. Necesita un runtime que pueda convertirse en la base durable del core gobernado y auditable.

---

## Candidatos evaluados y por qué entran al shortlist

### 1. Temporal

Entra al shortlist porque la evidencia verificada muestra durable execution, retries, event history explícito y patrones documentados de multi-tenant con namespaces, visibility/search attributes y task queues por tenant. Eso lo vuelve un candidato natural para el baseline técnico del core.

### 2. Camunda 8

Entra al shortlist porque la evidencia verificada muestra fortaleza real en process orchestration, user tasks, assignment, candidate groups, forms y APIs/tooling operacionales para procesos e incidentes. Es relevante porque Opyta Sync tiene approvals, bloqueos y operación humana gobernada.

### 3. Conductor OSS

Entra al shortlist porque la evidencia verificada muestra durable execution, wait/timer patterns, human task, replay/restart/pause/resume, full execution history, compensation flows tipo saga, workflows JSON-native determinísticos y workers polyglot. Además se presenta como self-hosted/open source sin cloud lock-in.

---

## Candidatos descartados por ahora y por qué no entran a la comparación final

### Restate

Queda fuera por ahora porque, aunque era candidato considerado, en esta iteración B.1 no forma parte del shortlist confirmado. Mantener la comparación acotada reduce ruido y fuerza una decisión sobre alternativas ya priorizadas.

### Dapr Workflow

Queda fuera por ahora por el mismo motivo: fue considerado, pero no está en la comparación final acordada para B.1.

### Azure Durable Functions / Durable Task

Queda fuera por ahora porque tampoco integra el shortlist final confirmado. Además, en esta fase conviene no mezclar la decisión del runtime baseline con opciones que no quedaron priorizadas en la investigación previa.

---

## Hard gates aplicados a los 3 candidatos

Los hard gates vienen de B.0 y se aplican antes del scoring.

| Hard gate | Temporal | Camunda 8 | Conductor OSS | Lectura |
|---|---|---|---|---|
| soporte real para multi-tenant | **pass** | **riesgo / no demostrado en evidencia actual** | **riesgo / no demostrado en evidencia actual** | En la evidencia verificada, solo Temporal trae patrones documentados explícitos de multi-tenant. |
| trazabilidad auditable | **pass** | **pass parcial** | **pass** | Temporal y Conductor tienen historia explícita verificada; Camunda muestra tooling/API operacional, pero el foco de la evidencia está más en operación BPMN/user-task que en baseline durable multi-tenant del core. |
| hooks o seams para approvals / classification / policy | **pass provisional** | **pass** | **pass** | Camunda y Conductor muestran user/human tasks; Temporal permite modelar la espera durable y las decisiones externas, pero el seam exacto con policy queda para B.2-B.6. |
| estado durable inspeccionable | **pass** | **pass** | **pass** | Los tres muestran alguna forma verificada de inspección/historial/operación. |
| ausencia de cajas negras imposibles de operar | **pass** | **pass** | **pass** | Los tres exponen tooling o historia operable; no se asume opacidad cerrada en la evidencia disponible. |

### Lectura de gates

- **Temporal** es el único candidato que, con evidencia verificada actual, pasa de forma clara el gate más sensible para Opyta Sync: multi-tenant real y operable.
- **Camunda 8** y **Conductor OSS** no quedan descartados por incapacidad demostrada, pero sí quedan tensionados por una falta de evidencia verificada en el punto que B.0 trata como gate duro: multi-tenant/governance tenant-scoped.
- Por lo tanto, el scoring comparativo de abajo debe leerse como **diagnóstico de fit relativo**, no como permiso para ignorar ese riesgo.

---

## Evaluación comparativa usando los criterios y pesos de B.0

Se usa la base común de B.0 con escala 1-5. El score ponderado se normaliza a 0-100.

### Criterios específicos de esta B.1

Aunque B.0 ya fija el marco principal, en esta comparativa el peso real se concentra en cuatro tensiones:

1. si el runtime soporta de verdad el modelo durable/auditable de `execution_record`;
2. si el aislamiento multi-tenant y la operación tenant-scoped están resueltos o al menos documentados con claridad;
3. si approvals/policy/classification pueden montarse como seams del core en lugar de workarounds laterales;
4. si la futura configuración conversacional declarativa quedará alineada o en tensión con el runtime elegido.

---

## Tabla de scoring comparativo

| Criterio | Peso | Temporal | Camunda 8 | Conductor OSS | Nota breve |
|---|---:|---:|---:|---:|---|
| alineación con invariantes de Fase A | 20% | 5 | 3 | 4 | Temporal encaja mejor con durabilidad/event history/multi-tenant; Conductor encaja bien con workflows determinísticos; Camunda luce más fuerte en proceso humano que en baseline técnico del core. |
| gobernanza y seguridad multi-tenant | 18% | 5 | 2 | 2 | Solo Temporal tiene evidencia verificada explícita de patrones multi-tenant. |
| durabilidad, auditabilidad y recuperación | 16% | 5 | 3 | 4 | Temporal y Conductor muestran durable execution e historia fuerte; en Camunda la evidencia verificada está más centrada en BPMN/user tasks y operación. |
| encaje con event log / trazabilidad / evidencia | 12% | 5 | 3 | 4 | Event history explícito favorece a Temporal; full execution history favorece a Conductor. |
| extensibilidad futura sin romper el core | 10% | 4 | 2 | 5 | Conductor puntúa alto por JSON-native deterministic workflows y workers polyglot; Temporal alto pero con más modelado técnico; Camunda queda más condicionado por su centro BPMN/process. |
| operabilidad IA-friendly | 9% | 4 | 3 | 5 | Conductor se alinea mejor con workflows JSON-native y workers polyglot; Temporal es operable y visible; Camunda es operable pero menos natural para una futura capa conversacional declarativa. |
| complejidad de implementación y riesgo de entrega en Fase C | 8% | 3 | 2 | 3 | Temporal y Conductor tienen costo de integración real; Camunda agrega peso operativo/modelado que no parece óptimo para el primer baseline del kernel. |
| lock-in / costo de reversión futura | 4% | 3 | 2 | 5 | Conductor puntúa mejor por su posicionamiento self-hosted/open source sin cloud lock-in. |
| madurez de ecosistema y tooling | 3% | 5 | 4 | 3 | Temporal y Camunda muestran tooling/documentación fuerte en la evidencia considerada. |

### Score ponderado

| Candidato | Score ponderado | Veredicto |
|---|---:|---|
| Temporal | **91.4 / 100** | **preferred** |
| Conductor OSS | **75.2 / 100** | **acceptable_with_tradeoffs** |
| Camunda 8 | **56.2 / 100** | **reject** |

---

## Tradeoffs narrativos por candidato

### Temporal

**Fortalezas**

- Es el candidato con mejor encaje técnico para la base durable del core.
- La evidencia verificada de namespaces, visibility/search attributes y task queues por tenant lo alinea con el invariant multi-tenant de Fase A.
- Durable execution, retries y event history explícito lo acercan de forma natural al modelo de `execution_record`, deduplicación, replay auditable y recuperación.

**Tradeoff principal**

- La tensión real está en la futura capa de authoring/configuración conversacional: Temporal favorece muy fuerte la durabilidad y la auditabilidad del runtime, pero no es la opción que más naturalmente “se deja escribir” como workflow declarativo AI-friendly.

### Camunda 8

**Fortalezas**

- Es el candidato más fuerte del grupo para user tasks, assignment, candidate groups, forms y operación enterprise de procesos/incidentes.
- Eso lo hace especialmente atractivo para approvals y orquestación humana gobernada.

**Tradeoff principal**

- B.1 no está eligiendo todavía el mejor subsistema de approvals o human workflow; está eligiendo el baseline durable del kernel. Ahí, la evidencia disponible lo muestra más centrado en process orchestration/BPMN y operación enterprise que en el tipo de runtime base multi-tenant y event-driven que Opyta Sync necesita como centro del core.

### Conductor OSS

**Fortalezas**

- Es el candidato que mejor conversa con la visión futura de workflows declarativos/AI-friendly: JSON-native deterministic workflows, human task, compensation flows, replay/restart/pause/resume y workers polyglot.
- Además mejora la historia de reversibilidad estratégica porque se presenta como self-hosted/open source sin cloud lock-in.

**Tradeoff principal**

- La gran tensión no es técnica de workflow sino de governance operable: en la evidencia verificada actual no quedó demostrado con la misma fuerza que Temporal cómo resolver multi-tenant real, boundaries de operación y gobernanza tenant-scoped del core.

---

## Riesgos y dudas abiertas por candidato

### Temporal

- Validar cuánto esfuerzo de modelado adicional exigirá para representar authoring declarativo conversacional sin deformar la UX futura.
- Validar patrón exacto de integración con policy engine, clasificación y approvals como artifacts externos gobernados.

### Camunda 8

- Validar si su fortaleza en BPMN/user-task no termina sobredimensionando el baseline del kernel.
- Validar multi-tenant/governance tenant-scoped con el mismo nivel de explicitud que exige B.0.
- Validar si la futura capa conversacional declarativa compilaría de forma limpia o quedaría demasiado subordinada al modelado BPMN.

### Conductor OSS

- Validar multi-tenant real, aislamiento y gobernanza tenant-scoped como gate duro, no como mejora futura.
- Validar cómo se mapean ownership, evidence, snapshots y correlación canónica de Fase A sin crear semántica paralela.
- Validar si la operabilidad cotidiana para humanos/IAs conserva el mismo nivel de claridad cuando el sistema crezca en tenants, approvals y artifacts gobernados.

---

## Recomendación final de B.1

**Opción A — Temporal recomendado como baseline** por durabilidad, multi-tenant y auditabilidad.

Es la recomendación más sólida porque es la única alternativa que, con la evidencia verificada actual, alinea de forma clara el corazón durable del runtime con el hard gate multi-tenant y con la necesidad de historia explícita/reconstruible para auditoría y recovery.

---

## Decisión provisional y su justificación

**Decisión provisional:** adoptar **Temporal** como baseline provisional del durable runtime para Fase C, dejando explícita una tensión abierta con la futura configuración declarativa conversacional.

### Justificación

1. **B.0 obliga a priorizar gates duros antes que simpatía de modelado.**
2. **Temporal** es el único candidato con evidencia verificada explícita de multi-tenant patterns y runtime durable con event history claro.
3. **Conductor OSS** es el mejor segundo candidato y probablemente el más atractivo si la prioridad principal fuera authoring declarativo/AI-friendly, pero hoy carga un riesgo no validado justo en el punto donde Opyta Sync no puede improvisar: governance multi-tenant operable.
4. **Camunda 8** aporta mucho para approvals/human tasks, pero no aparece como mejor baseline del durable runtime del core.

La honestidad arquitectónica acá importa: **sí existe una tensión real** entre elegir el runtime más fuerte para durabilidad/auditoría/multi-tenant y elegir el runtime más natural para workflows declarativos conversacionales futuros. En B.1 corresponde resolver a favor del baseline durable del core, y no al revés.

---

## Qué supuestos quedan pendientes para B.2-B.6

- **B.2 Policy engine comparison:** definir cómo policy se monta sobre el runtime elegido sin duplicar semántica de estado, approval o evidencia.
- **B.3 Memoria operativa y telemetría:** definir cómo event log, visibility, snapshots y evidencia alimentan observabilidad y memoria operativa sin crear un segundo sistema de verdad.
- **B.4 Extensibilidad:** definir seams de hooks/modules/workers para capacidades nuevas sin romper invariantes del core.
- **B.5 Packaging de capabilities:** definir cómo se distribuyen e instalan capabilities y workflows tenant/global sobre el runtime elegido.
- **B.6 Configuración conversacional:** validar cómo la capa conversacional compila a artifacts gobernados sobre Temporal y qué abstracción declarativa hará falta para que esa tensión no contamine la UX.

---

## Criterios de aceptación de B.1

B.1 puede considerarse cerrado si se cumplen todas estas condiciones:

1. existe shortlist explícito de candidatos comparados;
2. existen candidatos descartados por ahora con razón explícita;
3. se aplican hard gates de B.0 antes del scoring;
4. se usa la tabla de criterios/pesos comunes de B.0;
5. hay scoring comparativo visible y veredicto por candidato;
6. hay tradeoffs narrativos honestos por candidato;
7. hay riesgos y dudas abiertas por candidato;
8. existe una recomendación final explícita para el durable runtime;
9. la decisión provisional deja claro qué habilita para Fase C y qué tensiones traslada a B.2-B.6;
10. el documento no oculta la tensión entre durabilidad técnica y declaratividad conversacional futura.
