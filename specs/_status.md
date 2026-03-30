# 00 — Estado actual

## Intención del producto
Construir **Opyta Sync Engine** como base reusable para futuros proyectos, con estas propiedades:

- motor reusable y multi-tenant
- tools y workflows instalables por tenant o globales
- configurable desde chatbot
- IA-friendly para desarrollo, operación y mantenimiento
- gobernanza fuerte en permisos, clasificación, approvals y trazabilidad

## Estado general del plan
- Fase A: **cerrada a nivel de diseño v1**
- Fase B: **cerrada a nivel de decisión provisional**
- Fase C: **cerrada a nivel de construcción v1**
- Fase D: **cerrada a nivel de surface v1**
- Fase E: **cerrada como baseline reusable v1**
- A.6 Evals y testing base: **cerrado a nivel de diseño v1**
- A.5 Tenant, runtime y eventos generales: **cerrado a nivel de diseño v1**
- A.4 Approvals: **cerrado a nivel de diseño v1**
- A.3 Tipos de resultado: **cerrado a nivel de diseño v1**
- A.2 Contrato de intención/inspección: **cerrado a nivel de diseño v1**
- A.1 Objetos canónicos: **cerrado a nivel de diseño v1**
- Próximo bloque recomendado: **ninguno — roadmap actual cerrado**

## Decisiones clave ya cerradas
- La invalidación por cambio material es obligatoria y dura.
- Lo configurable se resuelve por policy/perfil/contexto, no libremente por ejecución.
- SoD y autoridad de aprobadores son configurables, pero gobernados.
- Riesgo de negocio y riesgo de seguridad se calculan por separado y luego se combinan.
- `pre_application` es el modo por defecto para cambios reales tenant-scoped.
- `double` se usa para irreversible, global, cross-tenant, policy y connector sensitive.
- El chat solo modifica configuración declarativa gobernada.
- No existe autoaprobación por timeout en escenarios `high` o `critical`.
- `aprobación/rechazo` se reemplaza por `governance_decision` como tipo nativo de resultado.
- Los 8 tipos de resultado se agrupan en 3 familias: read-only, mutation, governance.
- El éxito no es un bool: es un objeto `result_outcome` con nivel, evidencia y reason code.
- `system_update` irreversible no soporta `partial_success`: o completo o `failed`.
- Todo resultado pasa por `classifying` antes de ser entregado; no existe entrega sin clasificación aplicada.
- El fallo siempre tiene reason code de la taxonomía normalizada.
- El contrato de intención es la unidad operativa central y tiene dos momentos: intención inicial y contrato compilado.
- El contrato separa 4 grupos de campos: usuario, sistema, técnicos y snapshots.
- Todo cambio material del contrato recalcula fingerprint e invalida approvals previas.
- El contrato usa versionado `major.minor`: `major` para cambios materiales, `minor` para no materiales.
- No existe ejecución desde un contrato que no esté en `executable`.
- La lista final de objetos canónicos first-class del core queda cerrada en 15 objetos.
- `result_outcome` y los snapshots no son objetos top-level: viven embebidos.
- La verdad ejecutable se escribe en YAML y se compila a JSON determinístico.
- El envelope canónico es: `api_version`, `kind`, `metadata`, `spec`, `status`.
- El tenant tiene estado `operable` explícito; existir no alcanza para recibir tráfico real.
- `single_user` relaja organigrama y subadmins, pero no relaja seguridad, approvals ni clasificación.
- `execution_record` es un objeto first-class separado de contrato, approval y resultado.
- `execution_completed` y `application_completed` son estados distintos y obligatoriamente separados en mutation.
- Debe existir `idempotency_key` a nivel ejecución con deduplicación por tenant, capability, contract fingerprint, target scope y operation phase.
- Retry técnico, retry de ejecución y replay de auditoría son conceptos distintos.
- Rollback físico y compensación lógica no son equivalentes; cuando no hay rollback físico debe existir compensación o escalación manual explícita.
- Los evals son validaciones conductuales del motor sobre casos controlados, no solo unit tests.
- Deben existir suites mínimas de regresión: contract, result, approval, classification, runtime, resilience y tenant_onboarding.
- Fase A no exige kernel completo implementado, pero sí verdad ejecutable sin ambigüedades bloqueantes para pasar a Fase B/C.
- Fase B cerró **Temporal** como baseline de durable runtime del kernel por durabilidad, auditabilidad y multi-tenant operable.
- Fase B cerró **Cerbos** como baseline de policy engine para autorización contextual y governance tenant-scoped auditables.
- Fase B cerró **PostgreSQL + OpenTelemetry + Grafana LGTM** como split-plane para memoria operativa estable más observabilidad de ejecución.
- Fase B cerró **Declarative manifest + remote provider/worker model** como seam de extensibilidad sin introducir un segundo runtime de verdad.
- Fase B cerró **OCI bundle inmutable + firma + attachments** como artifact base y **Intent → Change Proposal → Governed Patchset** como pipeline de configuración conversacional gobernada.
- Fase C arranca explícitamente sobre ese baseline y no reabre esas comparativas dentro de la construcción inicial del kernel.
- El primer tramo operativo de Fase C prioriza compilación de contrato, runtime durable, policy y evidencia antes de ampliar superficie de producto.
- La construcción del kernel en Fase C se ordena por seams: contrato compilado, ejecución, enforcement, event log y registry/resolution de capabilities.
- C.1 fijó el compilador como pipeline puro persistido en PostgreSQL con `compiled_contract` + `compilation_report`.
- C.2 fijó el runtime skeleton sobre Temporal con `execution_workflow`, separación execution/application y camino explícito de blocking/failure/compensation.
- C.3 fijó el boundary PEP/PDP, el input canonizado a Cerbos y el mapping de policy decisions a runtime/approval/classification.
- C.4 fijó PostgreSQL como event log operativo append-only y OTel/LGTM como telemetría derivada con redacción previa a exportación.
- C.5 fijó registry, binding y provider resolution como capas separadas sobre OCI bundles y remote providers.
- C.6 cerró el checkpoint end-to-end del motor validando corredor mínimo entre contrato compilado, runtime durable, policy, event log y capability resolution.
- C.6 dejó explícito que distribution layer sigue fuera del roadmap actual y no forma parte del cierre del kernel en Fase C.
- Fase D se monta sobre ese kernel ya cerrado y usa conversación, proposals, preview e inspección como surfaces operativas, no como seams nuevos del motor.
- La surface de Fase D opera sobre artifacts y configuración gobernada; no reabre apply directo desde chat libre.
- D.1 fijó el boundary entre chat libre, intent gobernado y proposal-trigger sin permitir ejecución directa desde conversación.
- D.2 fijó `proposal_draft`, `governed_patchset_candidate`, diff humano/material y promotion gates antes de preview/apply candidate.
- D.3 fijó `preview_candidate`, `simulation_result`, simulaciones de policy/approval/classification/risk y gates previos a `apply_candidate`.
- D.4 fijó vistas de inspección correladas y recovery operacional permitido sin permitir mutación directa del estado canónico.
- D.5 fijó debugging semántico, `maintenance_action_candidate` y límites estrictos de automatización asistida sin habilitar bypass de governance.
- D.6 cerró el checkpoint end-to-end de la surface validando el corredor mínimo conversacional/operativo sobre el kernel ya cerrado en C.6.
- D.6 dejó explícito que el smoke path puede cerrar sin apply real si artifacts, simulación, vistas y evidencia quedan íntegros y correlados.
- Fase D queda cerrada a nivel de surface v1 sin reabrir seams del motor ni convertir la IA en autoridad de ejecución.
- Fase D no incluye distribution layer, tenant activation ni rollout de consumo como parte de su alcance.
- Fase E se recorta explícitamente al cierre reusable del engine + surface ya cerrados en C.6 y D.6.
- E.1 ya fijó la regresión integral y el hardening del baseline reusable de engine + surface como piso implementable de cierre.
- E.2 ya fijó qué artifacts y seams integran el baseline reusable y qué incluye el starter kit mínimo de capabilities compatibles.
- E.3 ya fijó el mapa IA-first final y playbooks diferenciados para operator, developer y debugging/mantenimiento.
- E.4 ya fijó el demo de referencia como corredor probatorio acotado del baseline reusable, sin breadth de producto ni distribution layer.
- E.5 ya cerró el readiness final del roadmap actual como baseline reusable v1.
- E.5 ya cerró el archive/handoff dejando el proyecto legible, transferible y auditable dentro del scope actual.
- El cierre total del baseline reusable no requiere distribution layer y mantiene explícito ese fuera-de-scope.
- La meta de Fase E no es expandir producto, sino endurecer regresión, evidencia, documentación y readiness del baseline reusable.
- El viejo lenguaje de tenant template y onboarding debe reinterpretarse como starter kit y playbooks dentro del boundary engine/surface.
- distribution layer sigue fuera del roadmap actual y no debe reingresar por Fase E de forma implícita.
- Existe una convergencia OSF cerrada como **addendum técnico v1** en `specs/osf-convergence/`.
- Esa convergencia no reabrió Temporal, Cerbos, PostgreSQL + OTel + LGTM ni el resto de decisiones duras del baseline reusable.
- La convergencia enriqueció implementation profile, connector SDK, transport interno, artifact plane y service topology sin cambiar el estado del roadmap principal.
- Existe un **implementation plan** derivado del baseline reusable y de la convergencia OSF en `specs/implementation-plan/`, sin cambiar el estado del roadmap principal.

## Qué no está cerrado todavía
- No hay bloques abiertos del roadmap actual.
