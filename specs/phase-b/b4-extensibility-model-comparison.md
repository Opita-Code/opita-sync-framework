# B.4 — Extensibility model comparison

## Objetivo de la comparativa

Definir qué modelo de extensibilidad debe adoptarse como baseline provisional para Fase C, de modo que Opyta Sync pueda incorporar nuevas capabilities, actions y handlers ejecutables sin romper los invariantes ya cerrados en Fase A y Fase B, ni degradar gobierno, auditabilidad, multi-tenant o authoring futuro asistido por IA.

La pregunta de B.4 NO es “cómo enchufamos código de terceros”. La pregunta real es: **qué seam estructural permite declarar capacidades de negocio gobernadas y resolver su ejecución contra workers/handlers aislables, versionables y auditables**.

---

## Qué exige Opyta Sync de su modelo de extensibilidad

Derivado de Fase A y de las decisiones provisionales ya cerradas en B.1-B.3, el modelo elegido debe poder sostener como mínimo estas exigencias:

- separar con claridad la **capability declarativa** del **handler/worker ejecutable**;
- permitir que el usuario interactúe sobre acciones de negocio, mientras el sistema resuelve contra sistemas base y handlers concretos;
- preservar `capability` como catálogo gobernado centralmente, con metadata, clasificación, ownership, versionado y policy visibles;
- montar approvals, clasificación, policy y boundaries multi-tenant como parte del camino principal, no como add-ons laterales;
- convivir con el baseline durable de B.1 sin duplicar semántica de estado, retries, pausas, compensación o evidencia;
- soportar aislamiento razonable entre kernel y extensiones, con una historia operable de sandboxing, blast radius y trust boundaries;
- habilitar upgrades/versionado de capabilities y handlers sin romper compatibilidad de contratos ni dejar ejecuciones a mitad de camino en estado ambiguo;
- dejar authoring futuro IA-friendly: una superficie declarativa legible, gobernable y compilable a artifacts estables.

En otras palabras: Opyta Sync no necesita un “plugin system” genérico. Necesita un **modelo gobernado de capacidades declarativas más ejecución resoluble**.

---

## Invariantes no negociables ya cerrados antes de B.4

- `capability` sigue siendo un **catálogo gobernado centralmente** y no un registro libre de código suelto.
- El usuario interactúa sobre **acciones de negocio**; el sistema decide qué sistemas base, providers o handlers concretos intervienen.
- Gobierno fuerte en **approvals, clasificación, policy y multi-tenant** es obligatorio.
- La futura **configuración conversacional IA-friendly** es un requisito importante, no un nice-to-have.
- La decisión provisional de **B.1** deja a **Temporal** como baseline provisional del runtime durable; por lo tanto, la extensibilidad NO debe introducir un segundo runtime de verdad.
- La decisión provisional de **B.2** deja a **Cerbos** como baseline provisional de policy; por lo tanto, las extensiones no deben bypassear evaluation y enforcement gobernado.
- La decisión provisional de **B.3** deja a **PostgreSQL + OpenTelemetry + Grafana LGTM** como baseline provisional de memoria operativa y telemetría; por lo tanto, las extensiones deben emitir evidencia correlacionable sin colonizar la memoria operativa.

---

## Candidatos evaluados y por qué entran al shortlist

### 1. Declarative manifest + remote provider/worker model

Entra al shortlist porque maximiza la separación entre:

- **manifest/capability declaration** como artifact gobernado, y
- **provider/worker ejecutable** como componente resoluble por contrato.

Es fuerte para catálogo central, authoring declarativo, upgrades graduales y boundaries operativos claros entre kernel y extensiones remotas.

### 2. Workflow-definition-first DSL + external activity workers

Entra al shortlist porque ofrece una superficie MUY atractiva para authoring declarativo y conversacional: primero se define el workflow/flujo en DSL y luego se resuelven las activities en workers externos. Es especialmente relevante si el futuro del producto prioriza composición declarativa rica como seam principal.

### 3. WASM sandboxed modules + central manifests

Entra al shortlist porque maximiza aislamiento técnico a nivel de ejecución embebida y promete una historia fuerte de sandboxing y portabilidad binaria. Es un candidato serio cuando la prioridad dominante es ejecutar lógica extensible dentro de fronteras muy estrictas del propio kernel.

---

## Qué enfoques se tratan como invariante/complemento y no como candidatos principales

### Capability façade + governed central catalog

Esto NO compite como candidato principal porque ya quedó fijado como invariante base del proyecto. La pregunta de B.4 no es si habrá façade/catálogo gobernado, sino **qué modelo de ejecución/extensión vive debajo de esa façade**.

### Hooks / interceptors

Tampoco compiten como modelo principal de extensibilidad. Se tratan como **mecanismo complementario** para cross-cutting concerns, enforcement, observabilidad o policy injection, pero NO como el seam central para representar capabilities de negocio.

### Distinción estructural clave: capability declarativa vs handler/worker ejecutable

Esta distinción debe quedar CRISTALINA:

- **capability declarativa** = artifact gobernado que describe intención de negocio, inputs/outputs, políticas, clasificación, contracts, versionado, compatibilidad, owner y metadata operativa;
- **handler/worker ejecutable** = implementación concreta que realiza la acción contra un sistema base o lógica operativa;
- el catálogo central registra y gobierna la capability declarativa;
- el runtime resuelve qué handler/worker ejecutable usar según tenant, versión, policy, environment, rollout y disponibilidad.

Si estas dos cosas se mezclan, el sistema deja de tener authoring gobernable y vuelve a un plugin registry disfrazado.

---

## Impacto transversal por dimensión crítica

### Sandbox / aislamiento

- **Declarative manifest + remote provider/worker**: el aislamiento principal ocurre por boundary remoto, identidad, tenancy, policy, network controls y execution envelope. No depende de ejecutar código embebido dentro del kernel para estar gobernado.
- **Workflow-definition-first DSL + external activity workers**: comparte una historia de aislamiento parecida al modelo remoto, pero suma una capa extra de compilación/proyección entre definición y ejecución.
- **WASM sandboxed modules + central manifests**: ofrece el sandbox técnico más fuerte dentro del host, pero ese beneficio puede venir acompañado de peor inspección operativa y más complejidad de plataforma.

### Versionado y compatibilidad de upgrades

- **Declarative manifest + remote provider/worker**: favorece versionar por separado manifest, contract y provider binding; eso permite upgrades graduales, compatibility matrices y rollouts por tenant/environment.
- **Workflow-definition-first DSL + external activity workers**: el versionado se vuelve más delicado porque hay que compatibilizar DSL source, artifacts compilados y ejecución durable in-flight sobre Temporal.
- **WASM sandboxed modules + central manifests**: obliga a sumar compatibilidad binaria/ABI, lifecycle de módulos y observabilidad de upgrades, lo que endurece la operación antes de tiempo.

### Authoring conversacional futuro

- **Declarative manifest + remote provider/worker**: deja una superficie declarativa muy buena para authoring IA-friendly, aunque menos rica para composición de flujos complejos si no se agrega una capa superior más adelante.
- **Workflow-definition-first DSL + external activity workers**: es el mejor para authoring conversacional como lenguaje principal de composición.
- **WASM sandboxed modules + central manifests**: es claramente el peor como interfaz primaria para authoring conversacional; sirve para ejecución, no para expresión de intención.

### Tensión con la decisión provisional de B.1 durable runtime

- **Declarative manifest + remote provider/worker**: la tensión es baja y manejable; el runtime durable sigue siendo Temporal y la extensibilidad solo resuelve handlers/workers.
- **Workflow-definition-first DSL + external activity workers**: la tensión es alta; el DSL corre riesgo de competir con Temporal como representación principal del flujo y exigir una compilación/proyección costosa.
- **WASM sandboxed modules + central manifests**: la tensión no viene por authoring sino por el riesgo de ocultar demasiado comportamiento ejecutable detrás del módulo, debilitando la semántica visible del runtime.

---

## Hard gates aplicados

Los hard gates vienen de B.0 y se aplican antes del scoring.

| Hard gate | Declarative manifest + remote provider/worker | Workflow-definition-first DSL + external activity workers | WASM sandboxed modules + central manifests | Lectura |
|---|---|---|---|---|
| soporte real para multi-tenant | **pass** | **pass provisional** | **pass provisional** | Los tres pueden segmentar por tenant, pero manifest + remote provider conversa mejor con routing tenant-scoped, rollout y ownership explícito por provider. |
| trazabilidad auditable | **pass** | **pass** | **pass parcial** | Manifest + provider y DSL + workers dejan contratos/ejecución más fáciles de correlacionar con runtime y evidencia; WASM tiende a comprimir más lógica dentro de un módulo menos visible semánticamente. |
| hooks o seams para approvals / classification / policy | **pass** | **pass** | **pass provisional** | Los tres pueden incorporar governance, pero el modelo remoto y el DSL exponen mejor seams antes/durante/después de la ejecución. |
| estado durable inspeccionable | **pass** | **pass con tensión** | **pass parcial** | El modelo manifest-resolved delega estado durable al runtime de B.1; el DSL también puede hacerlo, pero introduce tensión de compilación/representación; WASM arriesga empujar lógica y estado implícito al módulo. |
| ausencia de cajas negras imposibles de operar | **pass** | **pass con costo** | **riesgo / parcial** | WASM mejora aislamiento técnico, pero puede empeorar inspección, debugging y explicabilidad operativa si se transforma en contenedor opaco de comportamiento. |

### Lectura de gates

- **Declarative manifest + remote provider/worker** es el candidato que mejor pasa el filtro estructural completo porque preserva catálogo gobernado, boundaries operativos y execution resolution sin inventar un segundo modelo de verdad.
- **Workflow-definition-first DSL + external activity workers** también pasa, pero con una tensión estructural REAL respecto de B.1: si el DSL se vuelve la representación dominante, hay que compilarlo/proyectarlo limpiamente sobre Temporal sin crear semántica paralela.
- **WASM sandboxed modules + central manifests** no queda descartado como tecnología útil, pero sí queda tensionado en explicabilidad operativa, authoring conversacional y costo de entrega para Fase C.

---

## Evaluación comparativa usando criterios/pesos de B.0

Se usa la base común de B.0 con escala 1-5. El score ponderado se normaliza a 0-100.

### Tensiones específicas de esta B.4

Aunque B.0 fija el marco principal, en esta comparativa el peso real se concentra en cinco tensiones:

1. si la capability puede seguir siendo declarativa y gobernada sin confundirse con el código ejecutable;
2. si el modelo elegido se monta SOBRE el runtime durable de B.1 y no compite contra él;
3. si sandbox/aislamiento ayudan a gobierno real y no a esconder complejidad operativa;
4. si versionado y upgrades pueden sostener compatibilidad de contracts, providers y ejecuciones in-flight;
5. si la superficie futura de authoring conversacional sigue siendo legible, compilable y gobernable.

---

## Tabla de scoring comparativo

| Criterio | Peso | Declarative manifest + remote provider/worker | Workflow-definition-first DSL + external activity workers | WASM sandboxed modules + central manifests | Nota breve |
|---|---:|---:|---:|---:|---|
| alineación con invariantes de Fase A | 20% | 5 | 4 | 3 | Manifest + provider preserva mejor catálogo gobernado y separación capability/handler; DSL tensiona B.1; WASM inclina el centro hacia la implementación. |
| gobernanza y seguridad multi-tenant | 18% | 4 | 4 | 5 | WASM gana en aislamiento técnico local; manifest y DSL resuelven mejor governance operativa y tenant-scoped rollout sin requerir binarios embebidos como seam principal. |
| durabilidad, auditabilidad y recuperación | 16% | 4 | 4 | 3 | Manifest y DSL conversan mejor con el runtime durable y con evidencia explícita; WASM arriesga encapsular más lógica fuera de la semántica visible del runtime. |
| encaje con event log / trazabilidad / evidencia | 12% | 4 | 4 | 3 | Manifest/provider deja correlación clara entre declaration, resolution y execution; DSL también, pero necesita una proyección cuidadosa; WASM queda menos semántico para evidencia. |
| extensibilidad futura sin romper el core | 10% | 5 | 4 | 4 | Manifest/provider es el seam más modular para agregar providers y capabilities nuevas; DSL también extiende bien pero con más costo de plataforma; WASM extiende ejecución, no necesariamente authoring gobernado. |
| operabilidad IA-friendly | 9% | 5 | 5 | 2 | Manifest y DSL son declarativos y más aptos para authoring conversacional; WASM es pésimo como interfaz primaria para humanos/IA. |
| complejidad de implementación y riesgo de entrega en Fase C | 8% | 4 | 2 | 1 | Manifest/provider permite avanzar sin construir compilador/plataforma pesada; DSL y WASM suben mucho el costo inicial. |
| lock-in / costo de reversión futura | 4% | 4 | 4 | 3 | Manifest/provider y DSL dejan artifacts más portables a futuro; WASM puede fijar demasiado temprano el mecanismo de ejecución. |
| madurez de ecosistema y tooling | 3% | 4 | 3 | 2 | Manifest/provider reaprovecha tooling estándar de contracts/workers; DSL exige tooling propio; WASM suma complejidad especializada. |

### Score ponderado

| Candidato | Score ponderado | Veredicto |
|---|---:|---|
| Declarative manifest + remote provider/worker model | **87.8 / 100** | **preferred** |
| Workflow-definition-first DSL + external activity workers | **78.0 / 100** | **acceptable_with_tradeoffs** |
| WASM sandboxed modules + central manifests | **63.6 / 100** | **reject** |

---

## Tradeoffs narrativos por candidato

### Declarative manifest + remote provider/worker model

**Fortalezas**

- Es el candidato que mejor preserva la separación FUNDAMENTAL entre capability declarativa y handler/worker ejecutable.
- El catálogo central puede gobernar manifests, contracts, clasificación, policies, owners, compatibilidad y rollout sin acoplar esa capa al código concreto del proveedor.
- Conversa bien con el runtime elegido en B.1 porque el runtime sigue siendo quien modela estado durable, retries, pausas, compensaciones y evidencia; el modelo de extensibilidad solo resuelve qué worker/provider ejecuta cada acción.
- Mantiene authoring futuro IA-friendly porque la superficie principal sigue siendo declarativa: manifests, contracts, bindings, policies y mappings resolubles.

**Tradeoff principal**

- La tensión honesta es que este modelo NO convierte al workflow declarativo en el artefacto principal. Si mañana el producto quisiera authoring conversacional MUY orientado a composición de flujos complejos, habrá que construir una capa superior que genere manifests/bindings/workflow intents sin perder gobernanza.

### Workflow-definition-first DSL + external activity workers

**Fortalezas**

- Es el candidato más fuerte para una visión de authoring conversacional donde humanos o IA describen flujos de alto nivel y el sistema los compila a ejecución real.
- Separa bastante bien la definición declarativa del workflow respecto de las activities ejecutables.
- Puede ofrecer una experiencia poderosa para composición, validación y simulación de flujos de negocio complejos.

**Tradeoff principal**

- Acá está la tensión BRAVA con B.1 y hay que decirla sin maquillaje: **Temporal ya quedó elegido provisionalmente como durable runtime**. Si ahora el DSL pasa a ser el seam principal, Opyta Sync tendría que diseñar una traducción estable entre:
  - DSL declarativo,
  - artifacts gobernados del catálogo,
  - y el modelo real de ejecución de Temporal.

Eso implica costo de compilación/proyección, compatibilidad semántica, debugging de “fuente declarativa vs ejecución real”, y riesgo de terminar con **dos verdades**: la del DSL y la del runtime.

### WASM sandboxed modules + central manifests

**Fortalezas**

- Su mejor historia está en sandbox/aislamiento técnico del código ejecutable.
- Puede ser atractivo cuando el problema dominante es ejecutar lógica no confiable con límites estrictos dentro de un host controlado.

**Tradeoff principal**

- Opyta Sync NO está eligiendo una tecnología de aislamiento por sí misma; está eligiendo el modelo principal de extensibilidad del kernel. En ese frame, WASM corre el riesgo de empujar la conversación hacia “cómo embebemos módulos” en vez de “cómo gobernamos capabilities declarativas y su resolución operativa”.

---

## Riesgos y dudas abiertas por candidato

### Declarative manifest + remote provider/worker model

- Validar el contrato exacto entre manifest, capability version, provider binding y worker runtime para que la resolución sea estable y auditable.
- Definir política explícita de compatibilidad entre manifest versions y handler/provider versions.
- Confirmar cómo se resuelve aislamiento por trust level: provider remoto confiable, semi-confiable o third-party.

### Workflow-definition-first DSL + external activity workers

- Validar cuánto cuesta construir y gobernar el compilador/proyector de DSL sobre Temporal sin inventar semántica paralela.
- Definir debugging end-to-end entre fuente DSL, artifacts derivados y execution history real.
- Confirmar cómo se gobiernan upgrades del DSL cuando existen ejecuciones in-flight iniciadas con una versión anterior.

### WASM sandboxed modules + central manifests

- Validar si el aislamiento técnico compensa el deterioro potencial en explicabilidad, debugging y authoring.
- Confirmar qué superficie de APIs host necesitaría el kernel y si eso no recrea un framework propietario prematuro.
- Evaluar si el costo de packaging, firma, compatibilidad binaria y observabilidad de módulos entra razonablemente en Fase C.

---

## Recomendación final de B.4

**Opción A — Declarative manifest + remote provider/worker model** como baseline provisional del modelo de extensibilidad.

Es la recomendación más sólida porque ofrece el mejor equilibrio entre:

- catálogo central gobernado,
- separación capability/handler,
- compatibilidad con Temporal como runtime durable ya elegido en B.1,
- versionado y upgrades operables,
- y una superficie declarativa suficientemente apta para el authoring conversacional futuro.

### Por qué NO se pierde el objetivo conversacional futuro

Esto tiene que quedar CLARÍSIMO: elegir manifest + provider/worker **NO significa renunciar** a authoring conversacional futuro.

Significa, más bien, ordenar bien las capas:

1. **capa conversacional futura**: humanos/IA expresan intención, composición, restricciones y outcomes esperados;
2. **capa declarativa gobernada**: esa intención se normaliza a manifests, contracts, bindings, policies y artifacts versionables;
3. **capa ejecutable**: el runtime durable resuelve workers/providers concretos y ejecuta con evidencia auditable.

O sea: la conversación futura sigue siendo posible, pero NO se le pide a B.4 que convierta el DSL conversacional en seam principal ANTES de estabilizar packaging, versionado y execution contracts.

---

## Decisión provisional y justificación

**Decisión provisional:** adoptar **Declarative manifest + remote provider/worker model** como baseline provisional del modelo de extensibilidad para Fase C.

### Justificación

1. **B.0 obliga a priorizar gates duros e invariantes antes que elegancia abstracta de authoring.**
2. **Manifest + provider/worker** preserva mejor la separación central entre capability gobernada y ejecución resoluble.
3. Encaja mejor con **Temporal** porque deja el estado durable en el runtime y usa la extensibilidad solo para resolver handlers/workers, en lugar de introducir otro lenguaje de verdad operativa.
4. Deja una ruta limpia para authoring conversacional futuro mediante artifacts declarativos intermedios, sin exigir desde ya un compilador DSL-first sobre el runtime.
5. **Workflow-definition-first DSL** queda como mejor segundo candidato si en B.6 el proyecto decide invertir fuerte en una capa declarativa superior más expresiva.

La honestidad arquitectónica importa MUCHO acá: la tensión no es entre “uno moderno” y “uno viejo”. La tensión es entre **priorizar hoy un seam gobernable y compatible con el runtime ya elegido** versus **adelantar demasiado pronto una plataforma DSL-first** que puede ser correcta más adelante, pero que hoy sube demasiado el riesgo estructural.

---

## Qué supuestos quedan pendientes para B.5-B.6

- **B.5 Capability packaging comparison:** definir cómo se empaquetan manifests, contracts, policy attachments, schemas, bindings y references a providers/workers; ahí se cierra de verdad la historia de instalación/promoción/versionado.
- **B.5 también debe definir** firma, compatibilidad, rollback, promotion channels y ownership de paquetes tenant/global.
- **B.6 Configuración conversacional:** definir qué abstracción conversacional produce artifacts gobernados sobre este modelo sin exponer detalles accidentales de workers/providers al autor humano o IA.
- **Supuesto abierto importante:** puede existir una capa futura de workflow composition o DSL superior, pero debe compilar a manifests/bindings/contracts sin competir con Temporal como runtime de verdad.
- **Supuesto de aislamiento:** no todo sandbox debe resolverse en B.4 como WASM. Parte del aislamiento real puede venir de boundaries remotos, policy, network controls, identity, tenancy y execution envelopes gobernados.

---

## Criterios de aceptación de B.4

B.4 puede considerarse cerrado si se cumplen todas estas condiciones:

1. existe shortlist explícito de candidatos comparados;
2. queda explícito qué enfoques se tratan como invariantes/complementos y no como candidatos principales;
3. la diferencia entre capability declarativa y handler/worker ejecutable queda explicada sin ambigüedad;
4. se aplican hard gates de B.0 antes del scoring;
5. se usa la tabla de criterios/pesos comunes de B.0;
6. hay scoring comparativo visible y veredicto por candidato;
7. se explica el impacto de cada modelo en sandbox/aislamiento;
8. se explica el impacto de cada modelo en versionado y compatibilidad de upgrades;
9. se explica el impacto de cada modelo en authoring conversacional futuro;
10. queda documentada la tensión con la decisión provisional de B.1 durable runtime;
11. existe una recomendación final explícita y una decisión provisional clara;
12. quedan identificados los supuestos que deben cerrarse en B.5-B.6 sin reabrir innecesariamente B.4.
