# Phase 5 — Service topology and platform profile

## Objetivo

Mapear la topología de servicios propuesta por OSF sobre los seams ya cerrados del baseline reusable v1 y fijar un platform profile recomendado sin convertirlo en requisito normativo del core.

## Principios del platform profile

- Primero existe el seam lógico; después, si hace falta, el servicio desplegable.
- No todo seam requiere su propio proceso o microservicio.
- La topología debe clarificar ownership y despliegue, no fragmentar por moda.
- La plataforma recomendada orienta implementación realista; no redefine la arquitectura normativa.
- Componentes opcionales de plataforma deben declararse como opcionales de verdad.

## Diferencia entre seam lógico y servicio desplegable

- **Seam lógico**: boundary conceptual o funcional del baseline, definido por responsabilidad, invariantes y contratos.
- **Servicio desplegable**: unidad operacional concreta que implementa uno o más seams según necesidades de escala, aislamiento, seguridad o autonomía operativa.

Regla dura: los seams del baseline no dependen de una topología 1:1 de microservicios para seguir siendo válidos.

## Mapeo de servicios OSF a seams actuales

| servicio_osf | seam_actual_principal | rol en el baseline convergente | decisión |
|---|---|---|---|
| edge-gateway | surface conversacional / borde de interacción | termina protocolos de entrada como MCP o APIs compatibles, sin redefinir contrato normativo | adaptar |
| intent-service | compilación intención -> contrato | normaliza intención, prepara compilación y preserva boundary con conversación libre | adaptar |
| inspection-service | inspección/context gathering | reúne contexto, retrieval y restricciones previas a propuesta/ejecución | adaptar |
| planner-service | proposal/planning surface | genera o asiste planes/proposals sin saltarse governance | adaptar |
| workflow-orchestrator | runtime durable | materializa Temporal como runtime operativo del kernel | adoptar |
| approval-service | approvals/governance surface | gestiona tareas humanas, SLAs y tracking de decisiones | adaptar |
| policy-decision-service | boundary PEP/PDP | encapsula integración con Cerbos y decisiones auditables de policy | adaptar |
| catalog-service | registry/catalog seam | gestiona catálogo, manifests y visibilidad de capabilities | adaptar |
| connector-gateway | provider resolution / execution boundary | resuelve y encamina llamadas a providers remotos vía SDK estándar | adaptar |
| memory-service | truth plane / memory metadata | opera sobre PostgreSQL y metadata de memoria operativa, no sobre corpus masivo | adaptar |
| retrieval-service | retrieval/corpus plane | expone búsqueda amplia sobre OpenSearch y material indexable | adaptar |
| artifact-service | artifact/evidence plane | controla object storage, refs, snapshots y adjuntos | adaptar |
| evidence-service | evidence assembly/reporting | consolida refs, snapshots y reportes reproducibles | adaptar |
| model-gateway | plano IA | abstrae routing de modelos y proveedores LLM sin tocar governance core | adaptar |
| observability-eval-service | observability/evals complementarios | integra evals y trazas IA complementarias, con Langfuse opcional y OTel base | adaptar |

## Topología lógica de referencia

La topología lógica de referencia queda organizada en seis grupos:

1. **Borde y surfaces**: edge-gateway, intent-service, inspection-service, planner-service.
2. **Kernel operativo**: workflow-orchestrator, policy-decision-service, approval-service.
3. **Catálogo y extensibilidad**: catalog-service, connector-gateway.
4. **Datos, artifacts y retrieval**: memory-service, artifact-service, retrieval-service, evidence-service.
5. **Plano IA**: model-gateway, observability-eval-service.
6. **Plataforma**: ingress/gateway, autoscaling, observabilidad base, storage y componentes opcionales de identidad.

Esta topología es una referencia lógica. Puede desplegarse como menos servicios al inicio si se preservan boundaries y trazabilidad.

## Profile de plataforma recomendado

- **Kubernetes** como profile recomendado de despliegue y aislamiento operativo.
- **Gateway API** como profile recomendado de exposición L4/L7.
- **HPA/KEDA** como profile recomendado de autoscaling por métricas/eventos.
- **Object storage compatible con S3** como baseline persistente de artifacts/evidence.
- **Valkey** como componente recomendado de aceleración efímera.
- **OpenSearch** como componente recomendado cuando exista retrieval/corpus plane real.

## Qué entra como plataforma baseline vs qué queda opcional

### Entra como baseline convergente

- profile recomendado Kubernetes/Gateway API/HPA/KEDA
- object storage compatible con S3
- Valkey
- soporte para OpenSearch cuando se active retrieval/corpus plane

### Queda opcional

- Keycloak como IdP/SSO enterprise
- Langfuse como complemento del plano IA
- vLLM para hosting propio de inferencia/embeddings
- despliegues más fragmentados o dedicados por necesidad futura

## Evaluación de Keycloak como IdP/SSO opcional

Keycloak tiene encaje como opción de plataforma para autenticación humana y SSO enterprise. Su adopción no modifica el modelo principal de governance ni desplaza Cerbos. Queda recomendado como opción para entornos con necesidad real de OIDC/SAML y federación empresarial, pero no como requisito de validez del baseline convergente.

## Evaluación de Langfuse como complemento del plano IA

Langfuse queda aceptado como complemento especializado para trazas LLM, prompts y evals del plano IA. No reemplaza OTel/LGTM como observabilidad base del sistema ni se convierte en evidence plane general.

## Kubernetes/Gateway API/HPA/KEDA como platform profile y no como requisito normativo

Estos componentes quedan fijados como profile recomendado para implementación realista y operable. No forman parte de la fuente normativa A-E y no son condición necesaria para que la arquitectura siga siendo válida. El baseline converge sobre ellos como guía de plataforma, no como axioma del dominio.

## Tests borde mínimos

1. topología mínima colapsada en pocos servicios preserva seams lógicos
2. despliegue más fragmentado no altera autoridad de runtime/policy/truth planes
3. edge-gateway no puede saltarse intent/proposal/governance boundaries
4. workflow-orchestrator sigue siendo Temporal aunque cambie el layout de servicios
5. policy-decision-service no reemplaza ni duplica al PDP principal Cerbos
6. retrieval-service usa OpenSearch sin asumir truth plane
7. artifact-service persiste blobs fuera de PostgreSQL
8. memory-service no absorbe responsabilidades de retrieval masivo
9. model-gateway no adquiere autoridad sobre decisiones de negocio o policy
10. observability-eval-service complementa OTel sin convertirse en plano observability único
11. Keycloak opcional ausente no invalida el baseline convergente
12. Langfuse opcional ausente no invalida trazabilidad base del sistema
13. Kubernetes ausente en un entorno pequeño no invalida la arquitectura lógica
14. Gateway API/HPA/KEDA presentes no introducen nuevos seams normativos
15. connector-gateway sigue resolviendo providers según manifests/bindings aprobados
16. approval-service mantiene separación entre decisión humana y ejecución runtime

## Criterios de aceptación

1. Existe tabla de mapeo entre servicios OSF y seams actuales.
2. Queda fijada la diferencia entre seam lógico y servicio desplegable.
3. Queda definida una topología lógica de referencia coherente con A-E.
4. Queda fijado un platform profile recomendado.
5. Queda claro qué entra como baseline y qué queda opcional.
6. Queda evaluado Keycloak como IdP/SSO opcional.
7. Queda evaluado Langfuse como complemento del plano IA.
8. Queda explícito que Kubernetes/Gateway API/HPA/KEDA son profile y no requisito normativo.
9. Existen tests borde suficientes para validar la fase.
