# B.5 — Capability packaging comparison

## Objetivo de la comparativa

Definir qué modelo de packaging debe adoptarse como baseline provisional para capabilities en Opyta Sync, de modo que los artifacts puedan versionarse, firmarse, promoverse entre ambientes, instalarse y activarse por tenant sin romper los invariantes ya cerrados en Fase A ni las decisiones provisionales de B.1-B.4.

La pregunta de B.5 NO es “cómo publicamos cosas al catálogo” ni “cómo activamos algo en un tenant”. La pregunta real es: **qué unidad estructural empaqueta la capability declarativa y sus adjuntos gobernados de forma inmutable, auditable y compatible con el modelo de extensibilidad elegido en B.4**.

---

## Qué exige Opyta Sync del packaging de capabilities

Derivado de Fase A y de las decisiones ya cerradas en B.1-B.4, el packaging elegido debe poder sostener como mínimo estas exigencias:

- separar con claridad **artifact de packaging** de **distribución/publicación** y de **instalación/activación tenant**;
- permitir que la **promoción entre ambientes** mueva idealmente el MISMO artifact por digest y no recompilaciones distintas por ambiente;
- representar una capability como artifact **inmutable, firmable, verificable y auditable**;
- adjuntar schemas, manifests, contracts, policy attachments, metadata, documentación operativa y referencias versionadas sin volver el package una bolsa amorfa;
- convivir con el modelo de B.4: **declarative manifest + remote provider/worker model**, donde packaging no sustituye ni confunde capability, binding y activation;
- permitir relacionar explícitamente **contract version**, **binding version** y **provider runtime version**;
- sostener catálogo central gobernado, aprobación central y límites claros para overlays tenant;
- habilitar install/upgrade/rollback governado sin permitir que overlays tenant reescriban invariantes del package base;
- dejar authoring en **YAML**, pero producir **JSON canónico determinístico** como representación compilada/verificable del artifact.

En otras palabras: Opyta Sync no necesita solo “un formato para subir archivos”. Necesita una **identidad durable de artifact** que luego pueda circular por catálogo, planes de instalación y activación tenant sin perder trazabilidad ni semántica.

---

## Invariantes no negociables antes de B.5

- El authoring sigue siendo **YAML** y la forma normalizada debe ser **JSON canónico determinístico**.
- El catálogo de capabilities sigue siendo **central y gobernado**, no federado libremente por tenants.
- El control de catálogo, conectores y publication workflow sigue estando del lado de **plataforma/superadmin**.
- Los tenants pueden aplicar overlays o activaciones **solo dentro de límites explícitos** y nunca reescribir invariantes del package base.
- B.4 ya fijó como baseline provisional **declarative manifest + remote provider/worker model**.
- Packaging artifact ≠ publicación ≠ activación tenant.
- La activación tenant debe seguir siendo un **objeto/acto gobernado aparte** del artifact base.
- El package base debe poder vincularse con **contract version**, **binding version** y **provider runtime version** sin colapsar todo en una sola versión opaca.

---

## Candidatos evaluados y por qué entran al shortlist

### 1. OCI bundle inmutable + firma + attachments

Entra al shortlist porque ofrece una historia fuerte de artifact inmutable por digest, firma verificable, distribución estándar y attachments/adicionales versionados sin inventar un registry propietario desde cero. Es especialmente fuerte si Opyta Sync quiere una baseline clara para packaging artifact y deja publication/install/activation como capas superiores.

### 2. Catálogo + subscription/install plan + approval gates

Entra al shortlist porque modela muy bien la gobernanza de publicación, aprobación e instalación. Es un candidato serio cuando la prioridad dominante es el lifecycle operacional visible para plataforma y tenants.

Sin embargo, su fortaleza principal está más en **distribution/install governance** que en la pureza del **artifact de packaging**. Por eso entra al shortlist, pero también bajo sospecha: corre el riesgo de mezclar artifact, publicación y activación en una sola abstracción demasiado gorda.

### 3. Split package graph: definition package + binding package + tenant activation package

Entra al shortlist porque maximiza la pureza conceptual entre tres capas distintas:

- **definition package**: capability declarativa, contracts, schemas, policy attachments y metadata base;
- **binding package**: resolución hacia provider/binding/environment compatibility;
- **tenant activation package**: acto gobernado de instalación/activación tenant-scoped.

Es especialmente relevante si Opyta Sync quiere preservar de manera FORMAL la diferencia entre artifact base, publication workflow y activation object, aun al costo de introducir más complejidad de modelado.

---

## Hard gates aplicados

Los hard gates vienen de B.0 y se aplican antes del scoring.

| Hard gate | OCI bundle inmutable + firma + attachments | Catálogo + subscription/install plan + approval gates | Split package graph | Lectura |
|---|---|---|---|---|
| soporte real para multi-tenant | **pass** | **pass** | **pass** | Los tres pueden operar con boundaries tenant-scoped, pero OCI y split graph separan mejor artifact base de activation tenant. |
| trazabilidad auditable | **pass** | **pass parcial** | **pass** | OCI gana por digest + firma + provenance del artifact; split graph gana por explicitud semántica; catálogo/install plan depende demasiado del workflow y menos del artifact inmutable. |
| hooks o seams para approvals / classification / policy | **pass** | **pass** | **pass** | Los tres permiten insertar governance, pero el candidato catálogo lo hace como fortaleza principal, mientras OCI y split graph lo ubican en capas superiores sin contaminar el artifact base. |
| estado durable inspeccionable | **pass** | **pass** | **pass** | El estado durable de install/approval/activation debe vivir fuera del artifact en los tres casos; split graph lo explicita mejor, OCI lo exige disciplinadamente, catálogo tiende a mezclar estado operativo con packaging. |
| ausencia de cajas negras imposibles de operar | **pass** | **pass con tensión** | **pass con costo** | OCI usa primitives conocidas y verificables; catálogo/install plan puede volverse framework propietario gordo; split graph es explícito pero suma complejidad de coordinación y version matrices. |

### Lectura de gates

- **OCI bundle inmutable + firma + attachments** pasa el filtro estructural con el mejor balance porque da una unidad de artifact fuerte sin obligar a mezclar packaging con instalación tenant.
- **Catálogo + subscription/install plan + approval gates** no queda descartado, pero queda tensionado porque su centro de gravedad está en governance workflow y NO en la definición del artifact base.
- **Split package graph** también pasa y es arquitectónicamente serio, pero paga un costo real de complejidad que solo se justifica si Opyta Sync decide que la pureza formal entre definition, binding y activation merece esa inversión desde ya.

---

## Evaluación comparativa usando criterios/pesos de B.0

Se usa la base común de B.0 con escala 1-5. El score ponderado se normaliza a 0-100.

### Tensiones específicas de esta B.5

Aunque B.0 fija el marco principal, en esta comparativa el peso real se concentra en cinco tensiones:

1. si el modelo preserva la diferencia entre **artifact**, **publication/distribution** y **tenant activation**;
2. si la promoción entre ambientes puede mover el **mismo artifact por digest** sin recompilar por ambiente;
3. si overlays tenant quedan limitados a capas permitidas y NO reescriben invariantes del package base;
4. si el package puede relacionarse limpiamente con **contract version**, **binding version** y **provider runtime version**;
5. si el modelo habilita gobierno central y evidencia auditable sin inventar una plataforma de packaging innecesariamente compleja para Fase C.

---

## Tabla de scoring comparativo

| Criterio | Peso | OCI bundle inmutable + firma + attachments | Catálogo + subscription/install plan + approval gates | Split package graph | Nota breve |
|---|---:|---:|---:|---:|---|
| alineación con invariantes de Fase A | 20% | 5 | 3 | 5 | OCI y split graph respetan mejor artifact base gobernado, overlays limitados y separación de activation; catálogo/install plan tiende a mezclar capas. |
| gobernanza y seguridad multi-tenant | 18% | 4 | 5 | 5 | Catálogo/install plan y split graph son muy fuertes para gates tenant-scoped; OCI necesita que esa gobernanza viva arriba del artifact, pero eso es coherente y suficiente. |
| durabilidad, auditabilidad y recuperación | 16% | 5 | 3 | 4 | OCI gana por digest, firma y reproducibilidad del mismo artifact; split graph mantiene buena trazabilidad pero con más piezas; catálogo depende más de estado operacional que del artifact. |
| encaje con event log / trazabilidad / evidencia | 12% | 4 | 4 | 5 | Split graph expone explícitamente definition/binding/activation como objetos trazables; OCI también correlaciona bien si activation/install plan quedan como objetos separados; catálogo es aceptable pero menos limpio semánticamente. |
| extensibilidad futura sin romper el core | 10% | 4 | 3 | 5 | Split graph deja seams muy claros para crecer; OCI es suficientemente extensible con attachments y refs; catálogo/install plan arriesga volverse demasiado acoplado al workflow actual. |
| operabilidad IA-friendly | 9% | 4 | 3 | 4 | OCI y split graph conservan artifacts legibles y determinísticos; catálogo/install plan corre riesgo de esconder demasiada semántica dentro del lifecycle engine. |
| complejidad de implementación y riesgo de entrega en Fase C | 8% | 4 | 2 | 2 | OCI permite avanzar con primitives maduras; catálogo/install plan y split graph exigen más modelado/plataforma antes de capturar valor. |
| lock-in / costo de reversión futura | 4% | 4 | 2 | 4 | OCI y split graph dejan mejor portabilidad conceptual; catálogo/install plan podría encerrar demasiado en un workflow propietario. |
| madurez de ecosistema y tooling | 3% | 5 | 2 | 2 | OCI reaprovecha tooling estándar de registry, firma y distribución; split graph y catálogo requerirían más tooling propio. |

### Score ponderado

| Candidato | Score ponderado | Veredicto |
|---|---:|---|
| OCI bundle inmutable + firma + attachments | **87.6 / 100** | **preferred** |
| Split package graph: definition package + binding package + tenant activation package | **90.2 / 100** | **acceptable_with_tradeoffs** |
| Catálogo + subscription/install plan + approval gates | **65.8 / 100** | **reject** |

### Lectura del resultado

Sí, el **split package graph** obtiene un score numérico apenas superior. Y sin embargo la recomendación final NO tiene por qué seguir ciegamente el número. B.0 ya dejó clarísimo que el veredicto final también depende de la naturaleza del tradeoff, del costo de reversión y del riesgo de entrega.

En esta B.5, la pregunta NO es cuál es el modelo conceptualmente más puro en abstracto. La pregunta es cuál conviene adoptar como **baseline provisional para Fase C** sin adelantar complejidad estructural innecesaria. Bajo ese criterio, OCI bundle queda mejor posicionado como decisión operable hoy, mientras split graph queda como opción arquitectónicamente válida si más adelante se justifica formalizar las tres capas como packages separados.

---

## Tradeoffs narrativos por candidato

### OCI bundle inmutable + firma + attachments

**Fortalezas**

- Define una baseline MUY fuerte para el **artifact base**: inmutable, direccionable por digest, firmable y promovible entre ambientes sin recompilar.
- Permite que authoring siga en YAML y que el artifact publicado se materialice como **JSON canónico determinístico** junto con attachments versionados.
- Encaja muy bien con el catálogo central gobernado: el catálogo publica y aprueba referencias a artifacts OCI, pero NO redefine el formato del artifact.
- Mantiene limpia la arquitectura por capas:
  1. **package artifact** = OCI bundle firmado;
  2. **publication/distribution** = catálogo, channels, approvals, provenance;
  3. **install/tenant activation** = objeto gobernado aparte, tenant-scoped, auditable y reversible.
- Hace natural la promoción entre ambientes moviendo el MISMO artifact por digest, mientras bindings o activaciones por ambiente/tenant quedan fuera del bundle base.
- Permite relacionar el artifact con **contract version**, **binding version** y **provider runtime version** mediante metadata y references explícitas, sin colapsar todo en una única versión accidental.

**Tradeoff principal**

- OCI resuelve MUY bien la identidad del artifact, pero no resuelve por sí solo el lifecycle completo de install plan, approvals o tenant activation. Eso hay que modelarlo arriba. Y está BIEN que así sea: esa separación es precisamente una de las virtudes de la opción.

### Catálogo + subscription/install plan + approval gates

**Fortalezas**

- Es el candidato más fuerte para gobernanza operacional visible: publication workflow, approvals, subscriptions, install plans y gates por tenant/environment.
- Hace natural modelar separaciones de deberes entre plataforma, tenant admin y operadores.
- Puede integrarse muy bien con reason codes, policy y evidencia de instalación/upgrade.

**Tradeoff principal**

- Su problema es conceptual: corre el riesgo de convertir el lifecycle de publicación/instalación en el “package” mismo. Cuando eso pasa, **artifact**, **distribution** y **activation** quedan mezclados y la arquitectura pierde nitidez. Para Opyta Sync eso es una mala señal, porque la activación tenant debe seguir siendo un acto gobernado aparte y la promoción ideal entre ambientes debe mover el mismo artifact por digest.

### Split package graph: definition package + binding package + tenant activation package

**Fortalezas**

- Es el candidato más puro conceptualmente. Hace EXPLÍCITA la diferencia entre:
  - package de definición estable y portable,
  - package de binding/resolución hacia providers y entornos,
  - package u objeto de activación tenant-scoped.
- Refuerza de manera nativa que overlays tenant no pueden reescribir invariantes del package base porque directamente operan en otra capa del grafo.
- Hace sobresaliente la trazabilidad entre **contract version**, **binding version** y **provider runtime version**.
- Si Opyta Sync quisiera máxima claridad formal para supply chain, compatibility matrices y activation semantics, esta opción tiene muchísimo valor.

**Tradeoff principal**

- La pureza conceptual se paga. Hay más packages, más referencias cruzadas, más validaciones, más matrices de compatibilidad y más trabajo de tooling. Es una opción BUENÍSIMA si el proyecto decide que esa pureza vale el costo adicional, pero para un baseline provisional de Fase C puede ser demasiado peso demasiado temprano.

---

## Riesgos y dudas abiertas por candidato

### OCI bundle inmutable + firma + attachments

- Definir el layout exacto del bundle: qué va dentro del artifact, qué queda como attachment y qué queda solo como referencia externa gobernada.
- Definir el esquema de metadata para relacionar **capability version**, **contract version**, **binding compatibility** y **provider runtime compatibility**.
- Confirmar qué mecanismo de firma/provenance se usará como baseline y cómo se valida en catálogo e instalación.
- Definir el objeto exacto de **tenant activation** y cómo referencia el digest del bundle sin copiar ni recompilar el artifact.

### Catálogo + subscription/install plan + approval gates

- Validar cómo evitar que el workflow de catálogo termine reemplazando la semántica del package artifact.
- Definir si subscriptions/install plans apuntan a artifacts inmutables o a “latest approved” abstractions más opacas.
- Confirmar cómo se promueve el MISMO artifact por digest entre ambientes sin reempaquetar ni recalcular bundles.

### Split package graph: definition package + binding package + tenant activation package

- Validar si la complejidad adicional entra razonablemente en Fase C sin frenar entrega de valor.
- Definir reglas de compatibilidad entre definition package, binding package y activation object/package.
- Confirmar si tenant activation debe ser realmente un “package” persistente o más bien un objeto gobernado derivado del grafo publicado.
- Evaluar cuánto tooling extra hace falta para UX de catálogo, diffing, rollback y troubleshooting.

---

## Recomendación final de B.5

**OCI bundle inmutable + firma + attachments** como baseline provisional del packaging de capabilities.

Es la recomendación más sólida porque ofrece la mejor base para lo que B.5 realmente debe cerrar hoy: **la identidad del artifact**. Gana como baseline de artifact y firma sin confundir esa capa con publicación, instalación o activación tenant.

### Qué significa exactamente esta recomendación

Esto tiene que quedar CRISTALINO:

1. **package artifact** = OCI bundle inmutable, firmado y direccionable por digest;
2. **publication/distribution** = catálogo central gobernado, canales de promoción, approvals y metadata de release;
3. **tenant install/activation** = objeto/acto gobernado aparte que referencia artifacts y bindings aprobados;
4. **tenant overlays** = ajustes limitados por schema/policy que jamás reescriben invariantes del package base.

O sea: elegir OCI bundle NO significa que catálogo, install plan o tenant activation desaparecen. Significa que **quedan en la capa correcta**.

### Por qué no se recomienda catálogo/install plan como modelo principal

Porque sería empezar la casa por el techo. El catálogo y los install plans son importantes, sí, pero pertenecen a la capa de publication/activation governance. Si se los toma como “el package”, Opyta Sync corre el riesgo de perder trazabilidad limpia del artifact y de mezclar lifecycle operacional con identidad de supply chain.

### Cuándo podría valer la pena evolucionar hacia split package graph

Si en B.6 o en Fase C aparece la necesidad de formalizar fuertemente la separación entre definition, binding y activation para compatibilidad compleja, governance profunda o supply chain multi-etapa, entonces el camino natural es evolucionar desde OCI bundle baseline hacia un **split package graph** apoyado en artifacts OCI y objetos gobernados derivados. Esa evolución queda ABIERTA y no contradice la recomendación actual.

---

## Decisión provisional y justificación

**Decisión provisional:** adoptar **OCI bundle inmutable + firma + attachments** como baseline provisional del packaging de capabilities para Fase C.

### Justificación

1. **B.0 obliga a priorizar gates duros, auditabilidad y costo de reversión antes que pureza conceptual máxima.**
2. OCI bundle define con claridad la unidad que debe permanecer estable entre build/publication/promotion: el **artifact base por digest**.
3. Preserva una separación estructural correcta entre **packaging artifact**, **publication workflow** y **tenant activation**.
4. Encaja bien con B.4 porque el package puede transportar manifest declarativo, contracts, schemas, attachments y references a providers/bindings sin confundir capability con instalación tenant.
5. Permite que overlays tenant sigan existiendo, pero limitados por schema/policy y sin reescribir invariantes del package base.
6. Deja abierta una evolución futura hacia **split package graph** si la complejidad del dominio demuestra que esa pureza adicional compensa el costo extra.

La honestidad arquitectónica importa MUCHO acá: el split package graph es elegantísimo, sí. Pero la pregunta correcta es si esa elegancia ya paga en Fase C. Hoy, la respuesta más responsable es: **todavía no como baseline obligatorio**. Primero fijá bien el artifact. Después, si hace falta, refinás el grafo.

---

## Qué supuestos quedan pendientes para B.6 y Fase C

- **B.6 Conversational configuration comparison** debe definir qué abstracción conversacional produce overlays, bindings o activation intents sin permitir que la capa conversacional reescriba invariantes del package base.
- Fase C debe definir el **schema exacto del OCI bundle**: manifest canonical JSON, contracts, policy attachments, schemas, docs operativas, provenance y references permitidas.
- Fase C debe definir el **objeto de tenant activation**: su lifecycle, approvals, rollback, drift detection y relación exacta con bundle digest y binding version.
- Fase C debe definir la política de compatibilidad entre **capability version**, **contract version**, **binding version** y **provider runtime version**.
- Sigue abierto si algunos bindings deben viajar como attachments del artifact base o si deben publicarse como artifacts separados relacionados por digest/reference.
- Sigue abierto si el camino a largo plazo evoluciona a un **split package graph formal** o si alcanza con OCI bundle + objetos separados de binding/activation.

---

## Criterios de aceptación de B.5

B.5 puede considerarse cerrado si se cumplen todas estas condiciones:

1. existe shortlist explícito de candidatos comparados;
2. queda explícita la diferencia entre **packaging artifact**, **publicación/distribución** y **activación tenant**;
3. queda explícito que la activación tenant sigue siendo un **objeto/acto gobernado aparte**;
4. se documenta que la promoción entre ambientes debe mover idealmente el **mismo artifact por digest** y no recompilar por ambiente;
5. se documenta que overlays tenant no pueden reescribir invariantes del package base;
6. se documenta cómo el package se relaciona con **contract version**, **binding version** y **provider runtime version**;
7. se aplican hard gates de B.0 antes del scoring;
8. se usa la tabla de criterios/pesos comunes de B.0;
9. hay scoring comparativo visible y veredicto por candidato;
10. existe una recomendación final explícita y una decisión provisional clara;
11. queda documentado por qué **catalog + subscription/install plan** no se adopta como modelo principal;
12. quedan identificados los supuestos que deben cerrarse en B.6 y Fase C sin reabrir innecesariamente B.5.
