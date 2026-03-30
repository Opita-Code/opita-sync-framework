# Phase 1 — Implementation profile

## Objetivo

Fijar el perfil técnico de implementación derivado de OSF que puede incorporarse al baseline reusable v1 sin reemplazar la fuente normativa A-E ni reabrir decisiones duras ya cerradas.

## Principios del implementation profile

- El implementation profile traduce decisiones normativas a elecciones técnicas implementables.
- El profile solo puede enriquecer, nunca reemplazar, el baseline A-E.
- La prioridad es reducir ambigüedad de implementación sin crear una segunda arquitectura paralela.
- Todo componente adoptado debe respetar los seams ya cerrados: contrato, runtime, policy, event log, registry, artifact plane y surfaces.
- Los componentes opcionales deben entrar como complementos controlados, no como dependencias normativas invisibles.

## Qué problema resuelve

La matriz de Phase 0 define qué conviene adoptar, adaptar, diferir o rechazar. Esta fase resuelve el problema siguiente: convertir esa alineación en un perfil técnico concreto para implementación, de modo que equipos futuros sepan qué stack entra realmente al baseline convergente, bajo qué límites y con qué guardrails.

## Boundaries: source of truth vs implementation profile

El source of truth normativo define el modelo operativo, los objetos canónicos, las reglas de gobierno y la semántica ejecutable. El implementation profile define con qué tecnologías, profiles de despliegue y contratos internos derivados se implementa ese modelo sin modificarlo.

Regla dura: si un componente del implementation profile entra en tensión con la semántica normativa, prevalece la semántica normativa.

## Source of truth normativo vs wire contracts derivados

- El source of truth normativo sigue siendo **YAML declarativo** compilado a **JSON canónico determinístico**.
- Los wire contracts derivados existen para transporte interno, tipado, compatibilidad y ergonomía de clientes/servicios.
- Protobuf, Buf, ConnectRPC y gRPC no pueden redefinir contratos canónicos, campos normativos ni reglas de versionado material/no material.
- Toda interfaz interna derivada debe poder regenerarse desde la verdad declarativa sin divergencia semántica.

## Perfil técnico adoptado de OSF

El implementation profile convergente queda fijado así:

- **Temporal** como durable runtime baseline y único runtime de verdad.
- **Cerbos** como PDP principal del sistema.
- **PostgreSQL + OpenTelemetry + Grafana LGTM** como split-plane base ya cerrado.
- **S3-compatible object storage** como artifact/evidence plane persistente.
- **Valkey** como plano efímero de caché, hints de idempotencia, locks cortos y aceleración operativa.
- **OpenSearch** como retrieval/corpus plane complementario, nunca como truth plane.
- **Langfuse** como complemento del plano IA y evals, nunca como reemplazo de OTel/LGTM.
- **Protobuf + Buf + ConnectRPC/gRPC** como profile interno derivado de transporte.
- **Keycloak** como IdP/SSO opcional de plataforma.
- **OpenFGA** como capacidad futura diferida si aparece necesidad real de ReBAC profundo.

## Decisiones adoptadas directamente

- **Temporal** queda ratificado sin cambios.
- **S3-compatible object storage** entra directamente al implementation profile.
- **Valkey** entra directamente al implementation profile.
- **connector SDK estándar** queda asumido como baseline convergente para providers remotos.

## Decisiones adoptadas con adaptación

- **Protobuf/Buf/ConnectRPC/gRPC**: se adoptan solo como contratos internos derivados y capa de transporte.
- **OpenSearch**: se adopta solo como retrieval/corpus plane e índice documental complementario.
- **Langfuse**: se adopta como complemento del plano IA, sin invadir observabilidad general ni evidence trail canónico.
- **MCP edge**: se acepta como protocolo de borde compatible, no como redefinición del modelo central.
- **Kubernetes/Gateway API/HPA/KEDA**: se aceptan como platform profile recomendado, no como requisito normativo del baseline reusable.
- **PostgreSQL + pgvector**: PostgreSQL sigue siendo base; pgvector se admite como enriquecimiento controlado de memoria operativa semántica.

## Decisiones diferidas

- **Keycloak** queda diferido como IdP/SSO opcional hasta Phase 5, sin cambio obligatorio de modelo principal.
- **OpenFGA** queda diferido a futuro complemento de relaciones si aparece ReBAC profundo como necesidad dominante.
- **vLLM** queda diferido a escenarios de hosting propio, embeddings internos o presión real de costos/latencia.

## Decisiones rechazadas por ahora

- **OPA** como PDP principal queda rechazado por ahora porque reabriría Fase B y duplicaría autoridad con Cerbos.
- Cualquier uso de **Protobuf** como source of truth queda rechazado por ahora.
- Cualquier uso de **OpenSearch** como truth plane o observability plane queda rechazado por ahora.
- Cualquier reintroducción de **distribution layer** dentro de la convergencia queda rechazada por ahora.

## Reglas de consistencia con el baseline actual

- YAML + JSON canónico siguen siendo la fuente de verdad ejecutable.
- Cerbos sigue siendo el PDP principal.
- PostgreSQL conserva la verdad operativa durable.
- OTel + LGTM conservan la observabilidad general derivada.
- OpenSearch solo vive en retrieval/corpus plane.
- Langfuse no reemplaza ni duplica autoridad observability del split-plane base.
- Keycloak, si entra, solo entra como opción de plataforma y SSO.
- OpenFGA no entra al core mientras no exista evidencia extraordinaria de necesidad relacional profunda.
- El implementation profile no cambia la semántica de Intent -> Change Proposal -> Governed Patchset.
- Ningún componente del profile puede crear una segunda verdad de runtime, policy o artifacts.

## Criterios de aceptación de Phase 1

Phase 1 se considera cerrada si se cumplen todas estas condiciones:

1. Queda escrito que YAML + JSON canónico siguen siendo source of truth.
2. Queda escrito que Protobuf/Buf/ConnectRPC son solo derivados internos.
3. Queda confirmado que Cerbos sigue como PDP principal.
4. Queda confirmado que OpenSearch vive solo como retrieval/corpus plane.
5. Queda confirmado que Keycloak es opcional como IdP/SSO.
6. Queda confirmado que OpenFGA queda diferido hasta necesidad real de ReBAC profundo.
7. Queda confirmado que Langfuse complementa al plano IA sin reemplazar OTel/LGTM.
8. Queda confirmado que S3-compatible object storage y Valkey entran al implementation profile.
9. No se reabre ninguna decisión dura de Fase B.
10. El profile resulta implementable y coherente con las fases 2-6.
