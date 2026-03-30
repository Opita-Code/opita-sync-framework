# Phase 6 — Final convergence checkpoint

## Objetivo

Cerrar formalmente la convergencia OSF como addendum técnico v1 del baseline reusable, validando que las mejoras adoptadas no reabren decisiones duras ni alteran la fuente normativa A-E.

## Qué valida el checkpoint final

- que la convergencia enriqueció implementation profile, SDK, transport, artifact plane y topología
- que A-E sigue siendo la fuente normativa dominante
- que no se reabrieron Temporal, Cerbos, PostgreSQL + OTel + LGTM, manifest/provider model, OCI bundle ni pipeline conversacional gobernado
- que los componentes adoptados/adaptados quedan ubicados en el plano correcto
- que los diferidos y rechazados quedaron explícitos y no escondidos como dependencias implícitas

## Gates de convergencia

1. **Gate de autoridad normativa**: A-E sigue mandando sobre toda decisión convergente.
2. **Gate de no-reapertura**: ninguna decisión dura de Fase B fue reabierta.
3. **Gate de source of truth**: YAML + JSON canónico siguen siendo la verdad ejecutable.
4. **Gate de PDP**: Cerbos sigue siendo el PDP principal.
5. **Gate de split-plane**: PostgreSQL + OTel + LGTM siguen siendo la base de truth/observability.
6. **Gate de retrieval**: OpenSearch queda acotado al retrieval/corpus plane.
7. **Gate de extensibilidad**: el seam declarative manifest + remote provider/worker se mantiene intacto.
8. **Gate de artifacts**: object storage y evidence plane quedan alineados con OCI bundle + attachments.
9. **Gate de platform profile**: la plataforma recomendada no se convierte en requisito normativo.

## Qué mejoras quedaron adoptadas

- S3-compatible object storage como artifact/evidence plane persistente.
- Valkey como plano efímero de aceleración, locks cortos e idempotency hints.
- connector SDK estándar como baseline convergente para providers.
- ratificación de Temporal como runtime durable baseline dentro del addendum.

## Qué mejoras quedaron adaptadas

- Protobuf + Buf + ConnectRPC/gRPC como profile interno derivado.
- OpenSearch como retrieval/corpus plane, no truth plane ni observability plane.
- Langfuse como complemento del plano IA.
- MCP edge como protocolo de borde compatible.
- PostgreSQL + pgvector como base relacional con enriquecimiento semántico opcional y controlado.
- service topology de referencia mapeada a seams actuales.
- Kubernetes + Gateway API + HPA/KEDA como platform profile recomendado.
- document ingestion pipeline con ClamAV/Tika/Presidio y derivados controlados.

## Qué quedó diferido

- Keycloak como IdP/SSO opcional de plataforma.
- OpenFGA como complemento futuro si aparece ReBAC profundo.
- vLLM para inferencia/embeddings propios cuando exista justificación real.

## Qué quedó rechazado por ahora

- OPA como PDP principal.
- Protobuf como source of truth.
- OpenSearch como truth plane.
- OpenSearch como observability plane.
- reintroducción de distribution layer dentro de esta convergencia.

## Criterios de no-regresión respecto del baseline A-E

- no cambia la semántica del contrato canónico ni su versionado material/no material
- no cambia el runtime durable baseline
- no cambia el PDP principal
- no cambia el split-plane truth/observability ya cerrado
- no cambia el seam principal de extensibilidad
- no cambia el artifact base OCI bundle + firma + attachments
- no cambia la secuencia Intent -> Change Proposal -> Governed Patchset
- no mueve distribution layer dentro del scope actual

## Criterios para considerar cerrada la convergencia OSF

La convergencia OSF se considera cerrada si:

1. Las phases 0-6 existen y son coherentes entre sí.
2. Checklist y status quedan actualizados como cerrados.
3. Los componentes adoptados, adaptados, diferidos y rechazados están explícitos.
4. El addendum no contradice ninguna decisión dura del baseline reusable v1.
5. El resultado sirve como guía implementable para evolución futura sin necesidad de reinterpretación teórica adicional.

## Artifacts mínimos de cierre

- `_status.md` de convergencia cerrado
- `checklist.md` de convergencia completo
- `phase-0-alignment-matrix.md`
- `phase-1-implementation-profile.md`
- `phase-2-connector-sdk-baseline.md`
- `phase-3-internal-contracts-and-transport.md`
- `phase-4-artifact-plane-cache-and-retrieval.md`
- `phase-5-service-topology-and-platform-profile.md`
- `phase-6-final-convergence-checkpoint.md`

## Tests borde mínimos

1. intento de usar Protobuf como contrato normativo debe ser rechazado
2. intento de reintroducir OPA como PDP principal debe ser rechazado
3. intento de mover truth plane a OpenSearch debe ser rechazado
4. intento de usar Langfuse como observabilidad general debe ser rechazado
5. ausencia de Keycloak no invalida la convergencia cerrada
6. ausencia de OpenFGA no invalida la convergencia cerrada
7. ausencia de vLLM no invalida la convergencia cerrada
8. object storage y Valkey pueden entrar sin alterar A-E
9. connector SDK puede endurecer providers sin redefinir manifests/bindings
10. platform profile recomendado no se vuelve requisito normativo por error documental
11. document ingestion pipeline no invade truth plane ni policy plane
12. service topology de referencia no obliga 1:1 microservicios
13. pgvector opcional no desplaza a PostgreSQL como base transaccional
14. OpenSearch retrieval complementario no reemplaza evidence ni event log
15. el addendum completo mantiene distribution layer fuera de scope
16. status final indica cierre sin reabrir roadmap principal

## Criterios de aceptación

1. Los gates de convergencia quedan explícitos.
2. Las mejoras adoptadas, adaptadas, diferidas y rechazadas quedan consolidadas.
3. Los criterios de no-regresión quedan escritos y verificables.
4. Los artifacts mínimos de cierre quedan listados.
5. Existen tests borde suficientes para validar el checkpoint.
6. La convergencia queda cerrada como addendum técnico v1 y no como nuevo roadmap.
