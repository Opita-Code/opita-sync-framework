# Phase 0 — OSF alignment matrix

## Objetivo

Determinar qué partes de la propuesta OSF pueden incorporarse al baseline reusable v1 de forma controlada, qué partes deben reinterpretarse para respetar decisiones ya cerradas y qué partes quedan diferidas o rechazadas por ahora.

## Principios de convergencia

- La convergencia existe para **enriquecer** el baseline reusable v1, no para reemplazarlo.
- Toda adopción debe respetar el orden de verdad ya fijado: A-E primero, convergencia después.
- Solo se incorpora lo que mejora implementabilidad, operabilidad o claridad contractual sin romper seams ya cerrados.
- Cuando una propuesta OSF choca con una decisión dura del baseline, gana el baseline.
- La convergencia debe producir decisiones implementables, no una nueva ronda abierta de exploración teórica.

## Qué significa converger OSF sin romper el baseline

Converger OSF significa tomar componentes, perfiles técnicos y mejoras de naming, contratos internos, SDKs, retrieval, artifact plane y topología de servicios, pero **proyectándolos sobre el baseline ya cerrado**. Eso implica reinterpretar OSF como una capa de alineación y profesionalización del motor base, no como una nueva fuente normativa. En la práctica: OSF puede refinar cómo se implementa, expone o empaqueta el sistema; no puede redefinir qué runtime manda, qué PDP gobierna, cuál es la verdad operativa ni cómo se separan source of truth y contracts de transporte.

## Reglas duras de convergencia

- A-E sigue siendo la fuente de verdad normativa.
- OSF convergence no reabre Fase B salvo evidencia extraordinaria.
- No se reemplaza Cerbos como PDP principal.
- No se reemplaza el split-plane PostgreSQL + OTel + LGTM.
- No se reintroduce distribution layer.
- No se convierte Protobuf en source of truth.
- El contrato declarativo sigue viviendo en YAML con compilación determinística a JSON; los wire contracts son derivados.
- Temporal sigue siendo el runtime durable baseline y único runtime de verdad para ejecución gobernada.

## Tabla principal de alineación

| componente_osf | estado_actual_en_osf_baseline | decision | rationale | impacto | riesgo | fase_recomendada |
|---|---|---|---|---|---|---|
| MCP edge | presente en OSF como borde agente/IA hacia el motor; no forma parte del baseline duro A-E | adaptar | Sirve como edge protocol de integración con agentes y hosts IA, pero debe quedar como surface/edge compatible y no como redefinición del núcleo operativo | mejora interoperabilidad de borde y estandariza entrada agente-friendly | confundir edge protocol con contrato normativo del core | phase 5 |
| Temporal | ya alineado y cerrado en Fase B como durable runtime baseline | adoptar | Ya es decisión cerrada y compatible con el modelo OSF; no requiere reevaluación | cero fricción conceptual y máxima continuidad con C/D/E | reabrir comparativa innecesariamente | phase 0 |
| Protobuf + Buf + ConnectRPC/gRPC | propuesto en OSF para contratos internos; no forma parte hoy de la verdad normativa | adaptar | Puede ordenar contratos internos y compatibilidad entre servicios, siempre como proyección derivada del modelo declarativo y no como source of truth | mejora disciplina de interfaces internas, versionado y ergonomía de clientes internos | crear triple source of truth si se usa como contrato canónico primario | phase 3 |
| PostgreSQL + pgvector | PostgreSQL ya está cerrado como memoria operativa durable; pgvector aparece como enriquecimiento posible | adaptar | PostgreSQL ya es base. pgvector puede ampliar memoria operativa y recuperación semántica local sin desplazar el plano transaccional | refuerza memory/retrieval operativo sin romper el split-plane base | mezclar retrieval semántico con truth transaccional o sobredimensionar embeddings tempranamente | phase 4 |
| OpenSearch | no es baseline normativo actual; aparece como componente OSF para corpus y retrieval | adaptar | Tiene encaje como retrieval/corpus plane y búsqueda híbrida, pero no como truth plane ni observability plane | habilita corpus amplio, búsquedas híbridas e ingesta documental escalable | usarlo como pseudo-source-of-truth o duplicar memoria operativa | phase 4 |
| Keycloak | propuesto por OSF para autenticación humana/SSO empresarial; no es decisión cerrada en A-E | diferir | Puede ser valioso como IdP/SSO opcional de plataforma, pero no debe forzar un cambio de modelo principal en esta convergencia inicial | abre camino a integración enterprise sin tocar el núcleo actual | acoplar demasiado pronto identidad/plataforma a una implementación específica | phase 5 |
| OpenFGA | propuesto por OSF para relaciones y delegación; no es baseline vigente | diferir | Tiene sentido solo si aparece necesidad real de ReBAC profundo; hoy Cerbos cubre el baseline principal y el riesgo quedó explicitado como futuro | preserva opción futura de modelo relacional avanzado | introducir complejidad de grafo sin evidencia de necesidad dominante | phase 5 |
| OPA | OSF lo propone en gobierno; choca con Cerbos ya cerrado como PDP principal | rechazar_por_ahora | Reemplazar Cerbos reabriría una decisión dura de Fase B. Como mucho podría evaluarse a futuro para compliance/complementos no centrales | evita duplicación de PDP y preserva coherencia del baseline | crear conflicto de autoridad entre engines de policy | phase 6 |
| LiteLLM | propuesto en OSF como gateway de modelos | adaptar | Encaja como plano de model routing y abstracción de proveedores, siempre subordinado al contrato, governance y evidence trail ya definidos | mejora portabilidad de modelos, fallbacks y control operativo del plano IA | inflar prematuramente el plano IA y mezclar routing con decisión de negocio | phase 5 |
| vLLM | propuesto en OSF como inferencia propia opcional | diferir | Tiene valor sólo para escenarios de hosting propio, embeddings o costos/latencia específicos; no es requisito del baseline reusable v1 | deja preparada una ruta de escala futura | meter complejidad infra/ML sin necesidad inmediata | phase 5 |
| Langfuse | no es baseline actual, pero complementa OTel para trazas/evals LLM | adaptar | Puede enriquecer observabilidad del plano IA sin competir con OTel/LGTM ni con la evidencia operativa canónica | mejora trazabilidad de prompts, evals y linking de ejecución IA | duplicar observabilidad o tratar Langfuse como observability plane general | phase 5 |
| Kubernetes + Gateway API + HPA/KEDA | propuesto por OSF como perfil de plataforma de referencia | adaptar | Tiene sentido como platform profile de despliegue y topología de referencia, no como requisito normativo del core reusable | ordena despliegue, escalado y exposición de servicios | sobrediseñar infraestructura antes de fijar implementation profile real | phase 5 |
| S3-compatible object storage | consistente con attachments, evidencia y artifacts derivados del baseline | adoptar | Encaja directamente con OCI attachments, evidencia y storage de artefactos sin tensionar decisiones cerradas | fortalece artifact plane y evidencia durable | dispersar políticas de retención o clasificación si no se gobierna bien | phase 4 |
| Valkey | no estaba cerrado nominalmente en A-E, pero encaja con caché efímera, locks e idempotencia | adoptar | Complementa PostgreSQL sin competir con él y mejora runtime/connector ergonomics para hot paths | mejora caché, idempotency hints, locks cortos y performance operacional | abusarlo como estado durable o authority plane | phase 4 |
| connector SDK estándar | alineado con el seam de manifest + remote provider/worker model ya cerrado | adoptar | Formaliza la interfaz mínima del provider y endurece compatibilidad operativa del ecosistema de conectores | sube consistencia de providers, evidence, OTel e idempotencia | fijar un SDK demasiado rígido sin validar diversidad de conectores | phase 2 |
| document ingestion pipeline (ClamAV/Tika/Presidio) | no está definido como baseline actual, pero encaja con retrieval y output control | adaptar | Aporta una cadena clara para ingesta segura y utilizable de documentos, siempre subordinada a clasificación, redacción y evidence trail del baseline | mejora corpus ingestion, seguridad y preparación de contexto | mezclar ingestion con policy/output control o absorber scope excesivo muy temprano | phase 4 |
| service topology de referencia | OSF propone mapa explícito de servicios; el baseline actual está más descrito por seams | adaptar | Conviene mapear la topología OSF sobre seams ya cerrados para profesionalizar implementación y operación sin cambiar el modelo conceptual | clarifica boundaries, ownership y despliegue de componentes | convertir una topología sugerida en obligación prematura o introducir fragmentación excesiva | phase 5 |

## Conflictos estructurales identificados

1. **Cerbos vs OPA como autoridad principal de policy**: conflicto directo. Se mantiene Cerbos como PDP principal y OPA queda fuera del baseline core por ahora.
2. **YAML/JSON normativo vs Protobuf como contrato interno**: hay tensión si Protobuf intenta convertirse en verdad primaria. Debe quedar estrictamente derivado.
3. **PostgreSQL truth plane vs OpenSearch retrieval plane**: si OpenSearch absorbe funciones de memoria operativa o evidencia, rompe el split-plane cerrado.
4. **OSF service catalog vs seams ya cerrados en C/D/E**: la topología de servicios debe mapearse sobre seams existentes, no reemplazarlos ni redefinir responsabilidades nucleares.
5. **Gobierno compuesto OSF vs baseline ya consolidado**: Keycloak, OpenFGA y OPA no pueden entrar juntos como paquete de reemplazo del baseline de policy/gobernanza ya cerrado.
6. **Plataforma de referencia vs baseline reusable**: Kubernetes y componentes de plataforma sirven como profile de despliegue, no como requisito para declarar válida la arquitectura base.

## Oportunidades de mejora de alto ROI

- Formalizar un **connector SDK estándar** con contratos mínimos de provider, evidence, OTel e idempotencia.
- Incorporar **S3-compatible object storage** y **Valkey** como artifact/cache plane explícito del baseline convergente.
- Derivar **Protobuf + Buf + ConnectRPC/gRPC** para interfaces internas sin tocar la verdad normativa declarativa.
- Agregar **OpenSearch** y una **document ingestion pipeline** como retrieval/corpus plane controlado.
- Definir una **service topology de referencia** que haga más implementable el baseline reusable en entornos reales.
- Sumar **Langfuse** como complemento específico del plano IA y evals, sin competir con OTel/LGTM.
- Evaluar **pgvector** como mejora incremental de memoria operativa semántica sobre PostgreSQL.

## Lo que definitivamente NO se toca

- Temporal como runtime durable baseline.
- Cerbos como PDP principal.
- PostgreSQL + OpenTelemetry + Grafana LGTM como split-plane base.
- Declarative manifest + remote provider/worker model.
- OCI bundle inmutable + firma + attachments.
- Intent -> Change Proposal -> Governed Patchset.
- distribution layer fuera del roadmap actual.
- La regla de que YAML es la verdad declarativa y JSON determinístico la proyección ejecutable.

## Criterios de aceptación de Phase 0

Phase 0 se considera aceptada si se cumplen todas estas condiciones:

1. Existe una matriz explícita que mapea los principales componentes OSF contra el baseline reusable v1.
2. Cada componente relevante queda clasificado como adoptar, adaptar, diferir o rechazar_por_ahora.
3. Quedan documentados los conflictos estructurales sin reabrir decisiones ya cerradas en Fase B.
4. Se explicita qué mejoras tienen ROI alto y en qué fase conviene tratarlas.
5. Queda escrito qué elementos no se tocan bajo ninguna convergencia normal.
6. El próximo bloque recomendado queda fijado como **Phase 1 — Implementation profile**.
