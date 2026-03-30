# B.2 — Policy engine comparison

## Objetivo de la comparativa

Definir qué policy engine debe adoptarse como baseline provisional en Fase C para resolver autorización contextual, governance tenant-scoped, decisiones auditables y seams claros para approvals, clasificación y operación asistida sin romper los invariantes ya cerrados en Fase A ni el baseline durable recomendado en B.1.

---

## Qué exige Opyta Sync de un policy engine

Derivado de Fase A y de la decisión provisional de B.1, el motor de policy debe poder sostener como mínimo estas exigencias estructurales:

- evaluar autorización contextual sobre artifacts, ejecuciones, approvals y operaciones sensibles;
- operar con aislamiento multi-tenant real, porque el tenant sigue siendo la frontera canónica de governance;
- dejar evidencia auditable de decisiones, reason path y correlación operativa suficiente para reconstrucción;
- montarse como PDP externo gobernable, sin mezclar policy con lógica ad hoc del core;
- ofrecer seams claros para approvals, clasificación, enforcement y controles operativos;
- evitar que el modelo de policy obligue a rediseñar el runtime durable o a duplicar semántica del event log;
- ser operable por humanos e IA, con políticas y decisiones suficientemente explicables;
- minimizar el riesgo de que Fase C derive en una plataforma de policy propia por complejidad accidental.

En otras palabras: Opyta Sync no necesita solo “un permiso checkeable”. Necesita un motor capaz de gobernar decisiones sensibles del kernel con trazabilidad, boundaries por tenant y costo operacional razonable.

---

## Candidatos evaluados y por qué entran al shortlist

### 1. Cerbos

Entra al shortlist porque la evidencia verificada muestra autorización contextual con políticas declarativas en YAML, scoped policies con `scope` y `scopePermissions` para contextos jerárquicos o tenant, audit logs y decision logs configurables, y operación como PDP externo. Eso lo alinea de forma directa con governance tenant-scoped y operabilidad auditable.

### 2. OpenFGA

Entra al shortlist porque la evidencia verificada muestra un authorization system inspirado en Zanzibar, modelado central basado en tuples/relationships, soporte de conditions/context en relaciones y fortaleza real para relaciones, membresías y usersets. Es especialmente relevante para ownership, memberships y sharing entre actores.

### 3. OPA

Entra al shortlist porque la evidencia verificada lo muestra como general-purpose policy engine, policy-as-code con Rego, decision logs, masking, bundles/revisions y trazabilidad vía `decision_id`/`trace_id`. Es especialmente fuerte para decisiones complejas y auditabilidad técnica.

---

## Candidatos descartados por ahora y por qué

### Oso Cloud

Queda fuera de la comparación final porque fue considerado, pero no integra el shortlist confirmado de B.2. En esta iteración conviene forzar una decisión entre alternativas ya priorizadas.

### AWS Verified Permissions

Queda fuera por ahora por el mismo motivo: fue considerado, pero no forma parte de la comparación final acordada.

### Permit.io

También queda fuera por ahora porque no integra el shortlist final confirmado de B.2 y mezclar más opciones en esta iteración agregaría ruido en lugar de reducir riesgo de reversión.

---

## Hard gates aplicados a los 3 candidatos

Los hard gates vienen de B.0 y se aplican antes del scoring.

| Hard gate | Cerbos | OpenFGA | OPA | Lectura |
|---|---|---|---|---|
| soporte real para multi-tenant | **pass** | **pass parcial** | **pass provisional** | Cerbos aporta scopes jerárquicos y tenant-scoped policies; OpenFGA resuelve muy bien boundaries relacionales, pero no toda la necesidad contextual; OPA puede aislar por input/policy, pero exige más disciplina de plataforma propia para que el tenant boundary no sea accidental. |
| trazabilidad auditable | **pass** | **riesgo / parcial** | **pass** | Cerbos y OPA traen evidencia verificada explícita de decision/audit logs; en OpenFGA la fortaleza verificada está más en el modelo relacional que en la evidencia auditable rica del decision path. |
| hooks o seams para approvals / classification / policy | **pass** | **pass parcial** | **pass** | Los tres pueden participar como motor de decisión, pero OpenFGA encaja más naturalmente en relaciones/membresías que en clasificación y reglas contextuales amplias. |
| estado durable inspeccionable | **pass provisional** | **pass parcial** | **pass provisional** | En policy engine este gate se interpreta como inspección de artifacts y de decisiones; Cerbos y OPA encajan mejor por logs y artifacts de policy, mientras OpenFGA depende más de inspección del grafo relacional que de explicación contextual rica. |
| ausencia de cajas negras imposibles de operar | **pass** | **pass** | **pass con costo** | Ninguno es caja negra cerrada; la tensión real está en cuánto trabajo operacional propio exige cada uno para seguir siendo explicable. |

### Lectura de gates

- **Cerbos** pasa mejor el filtro estructural porque combina PDP externo, políticas declarativas y señales explícitas de auditabilidad sin obligar a construir demasiada plataforma alrededor.
- **OPA** también pasa, pero con una advertencia seria: su flexibilidad extrema no falla el gate, aunque sí traslada más responsabilidad al equipo para diseñar una plataforma de policy operable y consistente.
- **OpenFGA** no queda descartado como tecnología útil, pero sí aparece tensionado en los gates más sensibles para B.2 cuando la necesidad excede relaciones/membresías y entra en autorización contextual auditable del kernel.

---

## Evaluación comparativa usando criterios/pesos de B.0

Se usa la base común de B.0 con escala 1-5. El score ponderado se normaliza a 0-100.

### Tensiones específicas de esta B.2

Aunque B.0 fija el marco principal, en esta comparativa el peso real se concentra en cuatro tensiones:

1. si el engine resuelve autorización contextual y governance tenant-scoped sin deformar el core;
2. si la auditoría y la evidencia de decisión son suficientemente fuertes para operación real;
3. si approvals, clasificación y enforcement pueden montarse como seams del sistema en vez de policy ad hoc dispersa;
4. si la potencia del motor reduce riesgo sistémico o, al contrario, empuja a construir demasiada plataforma propia en Fase C.

---

## Tabla de scoring comparativo

| Criterio | Peso | Cerbos | OpenFGA | OPA | Nota breve |
|---|---:|---:|---:|---:|---|
| alineación con invariantes de Fase A | 20% | 5 | 3 | 4 | Cerbos alinea mejor con governance contextual auditable; OpenFGA encaja excelente en relaciones pero más angosto para la necesidad total; OPA encaja fuerte, aunque con más diseño propio. |
| gobernanza y seguridad multi-tenant | 18% | 5 | 4 | 4 | Cerbos sobresale por scoped policies jerárquicas; OpenFGA es muy sólido para boundaries relacionales; OPA puede hacerlo, pero no lo “trae” tan encuadrado como producto de policy multi-tenant. |
| durabilidad, auditabilidad y recuperación | 16% | 4 | 3 | 5 | OPA destaca por decision logs, masking, bundles/revisions y trazabilidad técnica; Cerbos también es fuerte por audit/decision logs; OpenFGA queda más corto en la evidencia verificada para auditabilidad rica. |
| encaje con event log / trazabilidad / evidencia | 12% | 4 | 3 | 5 | OPA conversa mejor con trazabilidad técnica profunda; Cerbos encaja bien para correlación operativa; OpenFGA no parece el mejor eje para evidencia contextual del kernel. |
| extensibilidad futura sin romper el core | 10% | 4 | 3 | 5 | OPA es el más flexible; Cerbos extiende bien dentro de un marco más gobernado; OpenFGA brilla si la extensión sigue siendo relacional, pero no necesariamente fuera de ese dominio. |
| operabilidad IA-friendly | 9% | 5 | 3 | 3 | YAML declarativo y scopes hacen a Cerbos más legible y gobernable para operación asistida; OPA es potente pero Rego sube complejidad cognitiva; OpenFGA es claro para relaciones, menos para policy contextual amplia. |
| complejidad de implementación y riesgo de entrega en Fase C | 8% | 4 | 3 | 2 | Cerbos reduce tiempo a valor como PDP externo declarativo; OpenFGA exigiría combinarlo con otras piezas para cubrir todo; OPA arriesga sobrecargar Fase C con diseño de plataforma propia. |
| lock-in / costo de reversión futura | 4% | 4 | 4 | 4 | Los tres son evaluables sin quedar atrapados en una caja negra comercial dentro del shortlist actual. |
| madurez de ecosistema y tooling | 3% | 4 | 4 | 5 | OPA destaca por ecosistema general y tooling de policy-as-code; Cerbos y OpenFGA muestran buena señal, aunque la decisión acá no se define por fama de mercado. |

### Score ponderado

| Candidato | Score ponderado | Veredicto |
|---|---:|---|
| Cerbos | **89.0 / 100** | **preferred** |
| OPA | **82.8 / 100** | **acceptable_with_tradeoffs** |
| OpenFGA | **65.6 / 100** | **reject** |

---

## Tradeoffs narrativos por candidato

### Cerbos

**Fortalezas**

- Es el candidato con mejor balance entre policy contextual, auditabilidad operativa y gobierno multi-tenant.
- Las scoped policies con `scope` y `scopePermissions` conversan muy bien con contextos jerárquicos por tenant, environment o dominio.
- El modelo declarativo en YAML y el rol de PDP externo ayudan a mantener la policy fuera del core sin volverla opaca.

**Tradeoff principal**

- La tensión honesta es el **ReBAC profundo**: si en Fase C-Fase D Opyta Sync necesitara un modelo dominante de relaciones complejas, usersets ricos y graph authorization como centro de todo, Cerbos no parece el ganador natural frente a OpenFGA.

### OpenFGA

**Fortalezas**

- Es el mejor candidato del grupo para ownership, memberships, sharing y usersets.
- Su raíz Zanzibar-style le da una ventaja clara cuando el problema central es relacional.
- Conditions/context en relaciones lo hacen más expresivo que un simple ACL store.

**Tradeoff principal**

- B.2 no está eligiendo solamente el mejor motor para graph authorization. Está eligiendo el baseline de policy del kernel. Ahí, OpenFGA queda demasiado centrado en relaciones como modelo principal y menos natural para policy contextual amplia, clasificación y evidencia auditable del decision path.

### OPA

**Fortalezas**

- Es el candidato más poderoso para decisiones complejas y auditabilidad técnica detallada.
- Decision logs, masking, bundles/revisions y trazabilidad por `decision_id`/`trace_id` lo hacen muy fuerte para compliance y debugging fino.
- Su flexibilidad permitiría modelar casi cualquier regla futura sin cambiar de motor.

**Tradeoff principal**

- La tensión real es brutal y hay que decirla sin maquillaje: **OPA gana en flexibilidad, pero te empuja a construir una plataforma de policy propia**. Si el equipo no disciplina authoring, testing, packaging, tenancy y gobernanza desde el día uno, esa flexibilidad se convierte en complejidad accidental.

---

## Riesgos y dudas abiertas por candidato

### Cerbos

- Validar hasta dónde alcanza el modelo elegido si ownership, delegación y sharing evolucionan hacia ReBAC más profundo.
- Validar el patrón exacto para combinar policy contextual con futuras necesidades de memberships o relationship checks sin duplicar semántica.
- Confirmar cómo se correlacionarán decision logs con `execution_id`, `tenant_id` y artifacts del runtime elegido en B.1.

### OpenFGA

- Validar cuánto policy contextual adicional quedaría fuera del modelo relacional y terminaría repartido en otras capas.
- Validar qué costo tendría complementar OpenFGA con otro mecanismo para auditoría rica, clasificación o reglas no relacionales.
- Validar si usarlo como motor principal generaría fragmentación conceptual entre graph authorization y policy contextual del core.

### OPA

- Validar cuánto esfuerzo real implicaría definir convenciones sólidas de Rego, tenancy, bundles, testing y authoring para no crear caos operativo.
- Validar si Fase C puede absorber ese costo sin desviar foco del kernel hacia una mini-plataforma de policy.
- Validar cómo se mantendrá explicabilidad suficiente para operación humana e IA cuando las reglas crezcan en volumen y complejidad.

---

## Recomendación final de B.2

**Opción A — Cerbos recomendado como baseline provisional del policy engine.**

Es la recomendación más sólida porque ofrece el mejor equilibrio entre operabilidad, auditabilidad, multi-tenant governance y autorización contextual declarativa SIN empujar a Opyta Sync a construir demasiada plataforma propia en Fase C.

La tensión principal debe quedar EXPLÍCITA: **Cerbos gana por operabilidad/auditoría/multi-tenant/scoped policies, pero deja abierto un riesgo si el dominio evoluciona hacia ReBAC profundo como necesidad central.** Ese riesgo existe. No hay que barrerlo abajo de la alfombra.

---

## Decisión provisional y justificación

**Decisión provisional:** adoptar **Cerbos** como baseline provisional del policy engine para Fase C, dejando explícito como riesgo abierto el eventual encaje de ReBAC profundo.

### Justificación

1. **B.0 obliga a priorizar gates duros antes que potencia abstracta.**
2. **Cerbos** es el candidato que mejor equilibra tenant-scoping, policy contextual declarativa, PDP externo y auditabilidad operable.
3. **OPA** es el mejor segundo candidato y probablemente el más potente si la prioridad absoluta fuera flexibilidad extrema y policy-as-code sin restricciones, pero hoy carga un riesgo serio de complejidad/plataforma propia para Fase C.
4. **OpenFGA** es excelente como solución relacional, pero no aparece como mejor motor principal para TODA la superficie de policy que Opyta Sync necesita en esta fase.

La honestidad arquitectónica acá importa MUCHO: la tensión no es entre “uno bueno” y “uno malo”. La tensión es entre **un motor más gobernable y operable hoy** versus **un motor más flexible pero más peligroso de operacionalizar mal**. Para B.2 conviene resolver a favor del baseline más gobernable del kernel.

---

## Qué supuestos quedan pendientes para B.3-B.6

- **B.3 Operational memory and telemetry comparison:** definir cómo las decisiones de policy y sus logs se correlacionan con memoria operativa, event log y observabilidad sin duplicar sistemas de verdad.
- **B.4 Extensibilidad:** definir si habrá seams específicos para plugins/capabilities que agreguen reglas, actions o evaluadores complementarios al baseline de Cerbos.
- **B.5 Packaging de capabilities:** definir cómo se versionan, distribuyen y promueven policies por tenant/global sin romper gobernanza ni trazabilidad.
- **B.6 Configuración conversacional:** validar qué abstracción declarativa conversacional podrá compilar o proyectar policies sin exponer complejidad accidental a autores humanos o IA.
- **Supuesto transversal abierto:** si el dominio real exige relationship-based authorization profundo como primer problema y no como necesidad secundaria, habrá que reevaluar el rol de OpenFGA como complemento o incluso reabrir parte de esta decisión.

---

## Criterios de aceptación de B.2

B.2 puede considerarse cerrado si se cumplen todas estas condiciones:

1. existe shortlist explícito de candidatos comparados;
2. existen candidatos descartados por ahora con razón explícita;
3. se aplican hard gates de B.0 antes del scoring;
4. se usa la tabla de criterios/pesos comunes de B.0;
5. hay scoring comparativo visible y veredicto por candidato;
6. hay tradeoffs narrativos honestos por candidato;
7. hay riesgos y dudas abiertas por candidato;
8. existe una recomendación final explícita para el policy engine;
9. la decisión provisional deja claro qué habilita para Fase C y qué tensiones traslada a B.3-B.6;
10. el documento deja explícita la tensión entre baseline gobernable/contextual y posible necesidad futura de ReBAC profundo.
