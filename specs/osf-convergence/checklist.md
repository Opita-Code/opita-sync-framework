# Checklist operativo de convergencia OSF

## Phase 0 — Alignment matrix

- [x] mapear componentes de OSF al baseline actual
- [x] identificar conflictos con decisiones cerradas
- [x] clasificar adoptar / adaptar / diferir / rechazar
- [x] definir equivalencias de naming y conceptos
- [x] definir límites de convergencia

## Professionalization v1 de Phase 0

- [x] fijar reglas duras de convergencia
- [x] consolidar matriz de alineación principal
- [x] identificar conflictos estructurales
- [x] fijar oportunidades de alto ROI
- [x] documentar lo que definitivamente no se toca
- [x] dejar criterios de aceptación explícitos

## Phase 1 — Implementation profile

- [x] fijar profile técnico de implementación
- [x] separar source of truth vs wire contracts
- [x] fijar límites de servicios/plataforma
- [x] documentar decisiones derivadas

## Professionalization v1 de Phase 1

- [x] fijar boundary entre source of truth y implementation profile
- [x] ratificar YAML + JSON canónico como verdad ejecutable
- [x] fijar Protobuf/Buf/ConnectRPC como capa derivada
- [x] incorporar S3-compatible object storage al profile
- [x] incorporar Valkey al profile
- [x] ubicar OpenSearch solo en retrieval/corpus plane
- [x] fijar Keycloak como opcional de plataforma
- [x] fijar Langfuse como complemento del plano IA

## Phase 2 — Connector SDK baseline

- [x] definir SDK estándar de conectores
- [x] definir contratos mínimos del provider
- [x] definir requirements de evidence/OTel/idempotencia
- [x] definir compatibilidad con registry/runtime

## Professionalization v1 de Phase 2

- [x] fijar boundary entre manifest, binding, provider y SDK
- [x] definir interfaz estándar mínima del conector
- [x] fijar input/output mínimo por método
- [x] exigir evidence refs mínimas
- [x] exigir idempotency key obligatoria
- [x] fijar spans y eventos OTel mínimos
- [x] fijar clasificación/riesgo/scopes mínimos
- [x] fijar compatibilidad con runtime y registry

## Phase 3 — Internal contracts and transport

- [x] definir uso derivado de Protobuf/Buf/ConnectRPC
- [x] fijar versionado entre YAML/JSON y wire contracts
- [x] definir compatibilidad backward/forward
- [x] evitar triple source of truth

## Professionalization v1 de Phase 3

- [x] fijar boundary entre contratos normativos y de transporte
- [x] definir artefactos generados desde Protobuf
- [x] fijar reglas de generación derivada
- [x] definir rol de Buf
- [x] definir rol de ConnectRPC/gRPC
- [x] fijar política anti triple-source-of-truth
- [x] documentar riesgos y guardrails
- [x] fijar tests borde del transport profile

## Phase 4 — Artifact plane, cache and retrieval

- [x] definir object storage baseline
- [x] definir Valkey baseline
- [x] definir retrieval/document ingestion baseline
- [x] fijar límites de OpenSearch

## Professionalization v1 de Phase 4

- [x] fijar boundary exacto entre PostgreSQL, object storage, Valkey y OpenSearch
- [x] definir qué vive en PostgreSQL
- [x] definir qué vive en object storage
- [x] definir qué vive en Valkey
- [x] definir qué vive en OpenSearch
- [x] fijar qué nunca debe vivir en cada plano
- [x] definir ingestion pipeline baseline
- [x] fijar reglas de redacción, clasificación y evidence

## Phase 5 — Service topology and platform profile

- [x] mapear servicios OSF a seams actuales
- [x] fijar topología lógica de referencia
- [x] evaluar IdP/SSO y componentes complementarios
- [x] fijar perfil de despliegue/plataforma

## Professionalization v1 de Phase 5

- [x] fijar diferencia entre seam lógico y servicio desplegable
- [x] mapear servicios OSF a seams A-E
- [x] definir topología lógica de referencia
- [x] fijar platform profile recomendado
- [x] separar baseline de plataforma vs opcionales
- [x] evaluar Keycloak como opcional
- [x] evaluar Langfuse como complemento IA
- [x] fijar Kubernetes/Gateway API/HPA/KEDA como profile no normativo

## Phase 6 — Final convergence checkpoint

- [x] verificar que no se reabren decisiones cerradas
- [x] consolidar mejoras adoptadas
- [x] registrar elementos diferidos/rechazados
- [x] cerrar convergencia

## Professionalization v1 de Phase 6

- [x] fijar gates de convergencia
- [x] consolidar mejoras adoptadas
- [x] consolidar mejoras adaptadas
- [x] consolidar diferidos
- [x] consolidar rechazados por ahora
- [x] fijar criterios de no-regresión respecto de A-E
- [x] fijar artifacts mínimos de cierre
- [x] declarar convergencia cerrada como addendum técnico v1
