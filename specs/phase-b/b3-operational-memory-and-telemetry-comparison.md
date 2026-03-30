# B.3 — Operational memory and telemetry comparison

## Objetivo de la comparativa

Definir qué estrategia split-plane debe adoptarse como baseline provisional en Fase C para resolver **memoria operativa estable** y **telemetría/observabilidad de ejecución** sin romper los invariantes ya cerrados en Fase A, ni las decisiones provisionales de B.1 y B.2.

La pregunta correcta de B.3 NO es “qué base única usamos para todo”. La pregunta correcta es: **qué combinación separa mejor el sistema de verdad operacional del sistema de evidencia analítica**, manteniendo correlación, gobernanza tenant-scoped y operabilidad real.

---

## Invariantes no negociables para memoria y telemetría en Opyta Sync

- la **memoria operativa** NO debe vivir en el backend de observabilidad;
- la **telemetría analítica** NO debe convertirse en sistema de verdad de negocio;
- `tenant_id` sigue siendo frontera canónica de aislamiento, gobernanza, approvals, clasificación y operación;
- el event log/runtime evidence puede alimentar observabilidad, pero **NO debe competir** con la memoria operativa como fuente de verdad estable;
- la reconstrucción de contexto operativo debe poder distinguir entre:
  - memoria durable del producto,
  - evidencia de ejecución del runtime,
  - datos analíticos para debugging/observabilidad;
- clasificación, minimización y redacción tenant-scoped deben aplicarse **antes de persistir telemetry sensible**;
- cualquier backend elegido debe preservar correlación auditable entre `tenant_id`, `environment`, `contract_id`, `execution_id`, policy decisions y artifacts relevantes.

---

## Qué exige el proyecto de la capa de memoria operativa

Derivado de Fase A y del source-of-truth ya fijado, la capa de memoria operativa debe poder sostener como mínimo estas exigencias:

- almacenar memoria operativa estable y gobernable, no solo evidencia efímera de runtime;
- soportar aislamiento multi-tenant real con controles consistentes y auditables;
- permitir políticas de acceso y evolución controlada del modelo sin hacks laterales;
- ofrecer mutación transaccional, consistencia suficiente y semántica clara de ownership;
- convivir con el runtime durable y con el policy engine sin duplicar semántica de estado;
- ser operable como sistema de verdad del dominio para contexto, snapshots útiles, facts y memoria reutilizable;
- permitir clasificación y retención diferenciada de datos sensibles por tenant, ambiente y tipo de artifact.

En esta comparativa, **PostgreSQL** queda fijado como base del plano operativo porque ya satisface la necesidad de memoria estable y gobernable mejor que cualquier backend de observabilidad. Además, los claims verificados permitidos confirman que PostgreSQL tiene **Row-Level Security** y que habilitar RLS sin policies efectivas deja un comportamiento de **default deny**, lo que conversa muy bien con el invariant tenant-scoped del proyecto.

---

## Qué exige el proyecto de la capa de telemetría/observabilidad

La capa de observabilidad debe poder sostener, como mínimo, estas exigencias:

- ingestión y consulta de **logs, traces, metrics y eventos de ejecución**;
- correlación fuerte con runtime, policy, approvals y artifacts relevantes;
- aislamiento tenant-scoped razonable para operación, soporte e investigación de incidentes;
- retención, búsqueda y análisis orientados a operación y debugging, no a verdad transaccional del negocio;
- compatibilidad natural con OpenTelemetry como capa de instrumentación y transporte;
- capacidad de redacción/clasificación previa al almacenamiento de señal sensible;
- costo operacional razonable para Fase C, sin montar una plataforma innecesariamente pesada antes de tiempo.

En otras palabras: Opyta Sync no necesita “tirar logs”. Necesita una capa de evidencia operacional que ayude a explicar qué pasó en ejecución SIN colonizar la memoria operativa del producto.

---

## Candidatos evaluados y por qué entran al shortlist

### 1. PostgreSQL + OpenTelemetry + Grafana LGTM

Entra al shortlist porque respeta de forma limpia la separación split-plane: PostgreSQL conserva la memoria operativa y el stack LGTM toma señales analíticas de ejecución. Los claims verificados permitidos además confirman que **Loki** y **Tempo** son multi-tenant y usan `X-Scope-OrgID` para aislar requests/datos por tenant, lo que vuelve a esta opción particularmente fuerte para operación tenant-scoped.

### 2. PostgreSQL + OpenTelemetry + ClickHouse-based observability

Entra al shortlist porque también mantiene la separación correcta entre memoria operativa y telemetría, y representa una familia de backends orientados a ingestión/consulta analítica de alto volumen con buena historia potencial de costo por volumen y explotación analítica. Es un candidato serio cuando la prioridad sube desde observabilidad estándar hacia análisis intensivo de eventos de ejecución.

### 3. PostgreSQL + OpenTelemetry + OpenSearch

Entra al shortlist porque preserva el split-plane y los claims verificados permitidos confirman que OpenSearch documenta **field-level security, field masking, audit logs y multi-tenancy en Dashboards**. Eso le da credenciales reales de seguridad y gobierno. Aun así, necesita probar que ese poder adicional aporta valor neto y no complejidad innecesaria frente a alternativas más alineadas con observabilidad nativa.

---

## Hard gates aplicados

Los hard gates vienen de B.0 y se aplican antes del scoring.

| Hard gate | PostgreSQL + OTel + Grafana LGTM | PostgreSQL + OTel + ClickHouse-based observability | PostgreSQL + OTel + OpenSearch | Lectura |
|---|---|---|---|---|
| soporte real para multi-tenant | **pass** | **pass provisional** | **pass provisional** | PostgreSQL aporta boundaries claros para memoria operativa; LGTM además tiene evidencia verificada de multi-tenancy en Loki/Tempo. En ClickHouse/OpenSearch el split-plane es válido, pero el aislamiento tenant-scoped completo debe cerrarse con diseño operativo explícito. |
| trazabilidad auditable | **pass** | **pass** | **pass** | Las tres estrategias separan correctamente memoria de evidencia y permiten correlación vía OpenTelemetry más metadatos canónicos. |
| hooks o seams para approvals / classification / policy | **pass** | **pass** | **pass** | Las tres permiten clasificar/redactar antes de persistir y correlacionar decisiones de policy con telemetría. |
| estado durable inspeccionable | **pass** | **pass** | **pass** | PostgreSQL cubre memoria operativa durable; el backend analítico cubre evidencia de ejecución sin usurpar esa verdad. |
| ausencia de cajas negras imposibles de operar | **pass** | **pass con costo** | **pass con costo** | Ninguna opción es una caja negra cerrada, pero ClickHouse/OpenSearch pueden subir complejidad operativa o de modelado antes de dar beneficio claro en Fase C. |

### Lectura de gates

- Las tres alternativas pasan el gate arquitectónico MÁS importante: **mantienen split-plane real** entre memoria operativa y telemetría.
- **Grafana LGTM** llega con la validación más directa para observabilidad tenant-scoped dentro del shortlist confirmado.
- **ClickHouse-based observability** y **OpenSearch** no fallan estructuralmente, pero cargan más incertidumbre operativa y requieren más disciplina para que la separación siga siendo clara en el día a día.

---

## Evaluación comparativa usando criterios/pesos de B.0

Se usa la base común de B.0 con escala 1-5. El score ponderado se normaliza a 0-100.

### Tensiones específicas de esta B.3

Aunque B.0 fija el marco principal, en esta comparativa el peso real se concentra en cinco tensiones:

1. si la opción mantiene sin ambigüedad la separación entre memoria operativa y evidencia analítica;
2. si el aislamiento tenant-scoped puede sostenerse de forma operable en ambos planos;
3. si OpenTelemetry puede actuar como capa común de instrumentación sin acoplar el core a un backend analítico concreto;
4. si logs/traces/metrics/eventos ayudan a explicar la ejecución sin convertirse en “pseudo-memoria” del negocio;
5. si el costo y la complejidad de Fase C quedan proporcionados al valor REAL que se obtiene.

---

## Tabla de scoring comparativo

| Criterio | Peso | PostgreSQL + OTel + Grafana LGTM | PostgreSQL + OTel + ClickHouse-based observability | PostgreSQL + OTel + OpenSearch | Nota breve |
|---|---:|---:|---:|---:|---|
| alineación con invariantes de Fase A | 20% | 5 | 4 | 3 | LGTM preserva mejor el modelo split-plane con semántica clara de observabilidad; ClickHouse también, aunque más orientado a analytics; OpenSearch mete más ambigüedad de propósito. |
| gobernanza y seguridad multi-tenant | 18% | 5 | 4 | 4 | PostgreSQL aporta RLS para memoria operativa; LGTM además trae claims verificados de multi-tenancy en Loki/Tempo; OpenSearch aporta controles potentes, pero con más superficie operacional. |
| durabilidad, auditabilidad y recuperación | 16% | 4 | 4 | 3 | En las tres, la durabilidad fuerte de memoria queda en PostgreSQL; LGTM y ClickHouse encajan mejor como evidencia operacional que OpenSearch como backend más generalista. |
| encaje con event log / trazabilidad / evidencia | 12% | 5 | 4 | 3 | LGTM conversa de forma más natural con telemetría de ejecución; ClickHouse puede ser muy fuerte para análisis, pero requiere más diseño; OpenSearch no es el fit más limpio para esta capa. |
| extensibilidad futura sin romper el core | 10% | 4 | 4 | 3 | LGTM y ClickHouse preservan mejor el desacople vía OTel; OpenSearch corre más riesgo de atraer usos laterales ajenos al core de observabilidad. |
| operabilidad IA-friendly | 9% | 4 | 4 | 3 | LGTM ofrece una lectura más directa del plano de observabilidad; ClickHouse puede ser muy útil para analytics avanzados; OpenSearch agrega más complejidad conceptual. |
| complejidad de implementación y riesgo de entrega en Fase C | 8% | 4 | 3 | 3 | LGTM es el baseline más equilibrado; ClickHouse y OpenSearch exigen más decisiones operativas para que el valor compense. |
| lock-in / costo de reversión futura | 4% | 4 | 4 | 3 | PostgreSQL + OTel reduce acoplamiento; OpenSearch tiende a expandir su rol si no se disciplina fuerte. |
| madurez de ecosistema y tooling | 3% | 5 | 4 | 4 | LGTM llega con una historia de observabilidad muy alineada al caso; ClickHouse/OpenSearch también son serios, pero no mejoran de forma decisiva el baseline. |

### Score ponderado

| Candidato | Score ponderado | Veredicto |
|---|---:|---|
| PostgreSQL + OTel + Grafana LGTM | **90.6 / 100** | **preferred** |
| PostgreSQL + OTel + ClickHouse-based observability | **79.6 / 100** | **acceptable_with_tradeoffs** |
| PostgreSQL + OTel + OpenSearch | **64.2 / 100** | **reject** |

---

## Tradeoffs narrativos por candidato

### PostgreSQL + OTel + Grafana LGTM

**Fortalezas**

- Es la opción que deja MÁS clara la arquitectura correcta: **Postgres para memoria operativa estable y gobernable; LGTM para logs, traces, metrics y evidencia de ejecución**.
- OpenTelemetry funciona como seam limpio entre el core y la capa analítica.
- Los claims verificados permitidos sobre multi-tenancy de Loki y Tempo refuerzan la operación tenant-scoped sin pedir inventos raros.

**Tradeoff principal**

- No es necesariamente la opción más “analytics-heavy” del shortlist. Si el proyecto priorizara análisis masivo y exploración analítica avanzada por encima del baseline equilibrado de observabilidad, ClickHouse-based observability podría volverse más atractivo.

### PostgreSQL + OTel + ClickHouse-based observability

**Fortalezas**

- Mantiene correctamente la separación split-plane.
- Tiene buen perfil como baseline orientado a costo/analytics cuando la telemetría de ejecución crezca mucho y el análisis sobre eventos agregados gane peso operativo.
- Preserva el desacople del core si OpenTelemetry sigue siendo la capa común de emisión.

**Tradeoff principal**

- Para Fase C puede introducir más complejidad de modelado y operación que LGTM sin entregar necesariamente mejor valor inmediato para observabilidad estándar del kernel.

### PostgreSQL + OTel + OpenSearch

**Fortalezas**

- Tiene señales reales de gobierno y seguridad: field-level security, field masking, audit logs y multi-tenancy documentados en Dashboards.
- Puede resultar atractivo si hubiera una necesidad MUY dominante de búsqueda avanzada y controles de acceso finos sobre grandes volúmenes de datos analíticos.

**Tradeoff principal**

- Acá hay que ser brutalmente honestos: **B.3 no está eligiendo una plataforma generalista para todo tipo de análisis y search**, está eligiendo el baseline split-plane más claro para memoria operativa + observabilidad. En ese frame, OpenSearch agrega poder, sí, pero no demuestra suficiente ventaja para justificar su mayor superficie conceptual y el riesgo de terminar usándolo como pseudo-sistema de verdad analítico-operacional.

---

## Riesgos y dudas abiertas por candidato

### PostgreSQL + OTel + Grafana LGTM

- Validar el patrón exacto de clasificación/redacción tenant-scoped antes de enviar telemetry sensible a Loki, Tempo o métricas derivadas.
- Definir hasta qué punto los eventos de ejecución se envían completos, resumidos o derivados para no competir con la memoria operativa.
- Confirmar convenciones canónicas de correlación entre runtime, policy y observabilidad.

### PostgreSQL + OTel + ClickHouse-based observability

- Validar si el beneficio analítico/costo realmente aparece en el horizonte inmediato de Fase C o si sería optimización prematura.
- Definir guardrails para que el event log runtime no derive en almacenamiento analítico con semántica de negocio paralela.
- Confirmar cómo quedará resuelto el aislamiento tenant-scoped end-to-end del plano analítico.

### PostgreSQL + OTel + OpenSearch

- Validar si la necesidad real del proyecto justifica una plataforma más search-centric que observability-centric.
- Evitar que field-level features y dashboards multipropósito empujen a mezclar datos operativos estables con evidencia analítica en un mismo hábito operacional.
- Confirmar si el costo operacional adicional mejora de verdad el baseline o solo lo vuelve más pesado.

---

## Recomendación final de B.3

**Opción A — PostgreSQL + OpenTelemetry + Grafana LGTM como baseline equilibrado.**

Es la recomendación más sólida porque preserva con mayor claridad el invariant arquitectónico central de B.3: **PostgreSQL se usa para memoria operativa estable y gobernable, mientras el backend de observabilidad se usa para trazas, logs, métricas y eventos de ejecución**.

Además, deja explícito algo que NO se puede negociar en Opyta Sync:

- el **event log/runtime evidence NO debe competir con la memoria operativa**;
- la observabilidad sirve para explicar, operar y depurar la ejecución;
- la memoria operativa sirve para sostener contexto durable, facts y estado gobernable del producto;
- clasificación/redacción tenant-scoped debe aplicarse **antes** de persistir telemetría sensible.

---

## Decisión provisional y justificación

**Decisión provisional:** adoptar **PostgreSQL + OpenTelemetry + Grafana LGTM** como baseline provisional de B.3 para Fase C.

### Justificación

1. **B.0 obliga a priorizar gates duros e invariantes antes que potencia genérica de plataforma.**
2. **PostgreSQL** fija mejor que los backends analíticos el rol de memoria operativa estable, gobernable y tenant-scoped.
3. **Grafana LGTM** ofrece el baseline más equilibrado para observabilidad de ejecución sin invadir la semántica del sistema de verdad.
4. **ClickHouse-based observability** queda como mejor segundo candidato cuando la prioridad pase a costo/analytics y no al baseline equilibrado inicial.
5. **OpenSearch** no queda recomendado salvo que aparezca una necesidad MUY fuerte y dominante de search/security-centric analytics que hoy no justifica desplazar a LGTM o ClickHouse.

La honestidad arquitectónica importa mucho acá: el error clásico sería intentar resolver memoria operativa y telemetría con “una sola cosa”. **Eso sería conceptualmente incorrecto para Opyta Sync.** El split-plane no es capricho; es una defensa contra mezclar verdad de negocio con evidencia de ejecución.

---

## Qué supuestos quedan pendientes para B.4-B.6

- **B.4 Extensibility model comparison:** definir cómo plugins/capabilities pueden emitir telemetría y leer/escribir memoria operativa sin saltarse boundaries tenant-scoped ni clasificación previa.
- **B.5 Packaging de capabilities:** definir qué artifacts de packaging declaran esquemas de memoria operativa, contratos de telemetry y reglas de redacción/retención.
- **B.6 Configuración conversacional:** definir cómo la capa conversacional podrá configurar instrumentación, eventos y memoria sin volver opaca la frontera entre facts del producto y evidence del runtime.
- **Supuesto transversal abierto:** la correlación canónica entre runtime, policy, approvals, memoria y observabilidad debe quedar fijada con naming/metadata estable antes de expandir capabilities.
- **Supuesto transversal abierto:** hay que decidir qué porción del event log se proyecta al plano observability como señal derivada y cuál queda exclusivamente en el plano operativo/auditable.

---

## Criterios de aceptación de B.3

B.3 puede considerarse cerrado si se cumplen todas estas condiciones:

1. existe shortlist explícito de estrategias split-plane comparadas;
2. la comparativa deja claro que B.3 NO evalúa una “base única para todo”;
3. se explicitan invariantes no negociables entre memoria operativa y telemetría;
4. se aplican hard gates de B.0 antes del scoring;
5. se usa la tabla de criterios/pesos comunes de B.0;
6. hay scoring comparativo visible y veredicto por candidato;
7. queda explícito que PostgreSQL sostiene la memoria operativa estable y gobernable;
8. queda explícito que el backend de observabilidad sostiene logs, traces, metrics y eventos de ejecución;
9. queda explícito que event log/runtime evidence NO compite con memoria operativa como sistema de verdad;
10. queda explícito que clasificación/redacción tenant-scoped debe aplicarse antes de persistir telemetry sensible;
11. existe una recomendación final clara entre las dos opciones permitidas como baseline;
12. la decisión provisional deja claro qué habilita para Fase C y qué tensiones traslada a B.4-B.6.
