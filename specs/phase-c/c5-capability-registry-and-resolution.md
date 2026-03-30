# C.5 — Capability registry and resolution

## Objetivo

Definir la construcción operativa v1 del seam de **capability registry and resolution** para que una `capability` canónica, gobernada centralmente y empaquetada como **OCI bundle inmutable + firma + attachments**, pueda resolverse de forma determinística hacia un `provider_ref` remoto ejecutable mediante un `binding` válido, sin absorber lifecycle durable, sin modelar tenant activation y sin abrir una distribution layer fuera del roadmap actual.

## Principios de implementación del registry y resolution

1. **`capability` sigue siendo objeto canónico gobernado centralmente.** El registry no reemplaza el catálogo canónico; lo materializa para consulta y resolución.
2. **El registry responde por catálogo + resolución.** NO responde por lifecycle durable, rollout durable, install plan durable ni distribution.
3. **El runtime resuelve por cadena completa, no por existencia global.** Debe existir relación válida `capability_manifest -> bundle_digest -> binding -> provider_ref`.
4. **Packaging y resolution son capas distintas.** Packaging fija el artifact inmutable; resolution decide si existe un binding vigente y compatible para usarlo.
5. **Tenant activation queda fuera de C.5 y fuera del roadmap actual.** En esta fase no se modelan activaciones tenant-scoped ni overlays tenant.
6. **Fail-closed obligatorio.** Si falta manifest, bundle, binding, provider approval, compatibilidad o evidencia mínima de artifact verification, la resolución falla.
7. **Versionado explícito.** Debe existir compatibilidad verificable entre `capability_version`, `manifest_schema_version`, `binding_version`, `provider_runtime_version` y `contract_schema_version`.
8. **Artifact verification forma parte del resolution path.** Debe verificarse digest y presencia de evidencia de firma/provenance, aunque la verificación criptográfica completa pueda quedar stubbeada en v1.
9. **Disponibilidad del binding es first-class.** El registry debe modelar `active`, `disabled`, `deprecated` y `retired`.
10. **Observabilidad derivada no reemplaza verdad canónica.** La resolución canónica debe persistir aunque falle la proyección a event log o telemetría.
11. **Distribution layer queda fuera del roadmap actual.** C.5 sólo asume artifacts ya publicados/referenciables; no define channels, sync, pull/push ni replication.

## Boundary exacto entre registry, packaging, binding y provider resolution

### `capability_registry`

Responsable de:

- catalogar manifests válidos por `capability_id` + `capability_version`;
- indexar la relación del manifest con `bundle_digest`;
- resolver bindings válidos por environment y compatibilidad;
- entregar una respuesta canónica de resolución o un rechazo explícito.

No hace packaging, no hace distribución, no activa por tenant y no ejecuta providers.

### Packaging

Responsable de:

- producir un artifact OCI inmutable direccionable por `bundle_digest`;
- incluir manifest canonical, attachments y metadata verificable;
- asociar evidencia de firma/provenance.

Packaging NO decide binding vigente, NO resuelve `provider_ref` y NO determina si un runtime puede ejecutar hoy esa capability.

### Binding

Responsable de:

- asociar una capability versionada y un `bundle_digest` con un `provider_ref` concreto;
- restringir esa asociación por `environment`, compatibilidad de contratos, tipos de resultado y `provider_runtime_version`;
- expresar estado operativo (`active`, `disabled`, `deprecated`, `retired`).

Binding NO redefine el manifest, NO recompila el bundle y NO activa tenants.

### Provider resolution

Responsable de:

- localizar el `provider_ref` final ya aprobado y operativo;
- verificar que el provider referenciado satisface runtime, environment y constraints declarados;
- devolver el target remoto ejecutable o fallar cerrado.

Provider resolution NO publica providers ni gobierna su lifecycle durable completo.

## Responsabilidades explícitas del capability registry

El `capability_registry` DEBE:

1. registrar `capability_manifest` canónicos por `capability_id` + `capability_version`;
2. mantener la referencia exacta entre manifest y `bundle_digest`;
3. exponer bindings resolubles por `environment`;
4. resolver sólo contra bindings válidos, vigentes y no ambiguos;
5. aplicar guardas mínimas de compatibilidad entre manifest, binding, provider y contrato compilado;
6. validar que el `provider_ref` referenciado exista y esté aprobado para uso;
7. verificar digest y presencia de evidencia mínima de firma/provenance en el camino de resolución;
8. distinguir estados de binding `active`, `disabled`, `deprecated` y `retired`;
9. detectar entradas duplicadas, stale o conflictivas;
10. persistir resolución canónica e índices/proyecciones necesarias para consulta operativa;
11. integrarse con `compiled_contract`, `execution_workflow`, policy, approvals y event log sin mover sus boundaries;
12. soportar idempotencia de registro y de resolución derivada.

## No-responsabilidades explícitas del capability registry

El `capability_registry` NO DEBE:

1. modelar activación tenant-scoped;
2. modelar overlays tenant;
3. absorber distribution layer, publication workflow o promotion channels;
4. reemplazar packaging, signature service o provenance service;
5. decidir policy contextual completa por fuera de las guardas normativas requeridas para resolver;
6. reemplazar `compiled_contract` ni `execution_record` como objetos canónicos;
7. permitir resolución “best effort” si hay ambigüedad o falta de evidencia mínima;
8. elegir dinámicamente “el provider que mejor parezca” sin binding explícito;
9. reinterpretar manifests o bindings con reglas implícitas no versionadas;
10. persistir telemetría derivada como verdad primaria.

## Manifest schema ejecutable mínimo

El manifest mínimo ejecutable DEBE cubrir estos grupos de campos:

### 1. Identidad

- `capability_id`
- `capability_version`
- `manifest_schema_version`
- `display_name`
- `owner_ref`
- `lifecycle_status` (`draft|valid|deprecated|retired`)

### 2. Contracts / result types soportados

- `contract_schema_version`
- `supported_input_contracts[]`
- `supported_output_contracts[]`
- `supported_result_types[]`
- `operation_refs[]`

### 3. Requirements (`policy` / `classification` / `approvals`)

- `required_policy_profiles[]`
- `required_classification_level`
- `approval_mode`
- `approval_profile_ref`
- `enforcement_hints[]`

### 4. Execution binding refs

- `binding_mode` (`provider_bound` en v1)
- `default_binding_selector` opcional
- `allowed_environments[]`
- `required_provider_capabilities[]`

### 5. Packaging refs

- `bundle_digest`
- `bundle_media_type`
- `required_attachments[]`
- `signature_ref` opcional
- `provenance_ref` opcional

### 6. Compatibility metadata

- `manifest_schema_version`
- `compatibility.contract_schema_versions[]`
- `compatibility.binding_versions[]`
- `compatibility.provider_runtime_versions[]`
- `compatibility.runtime_constraints[]`

### Schema YAML mínimo v1

```yaml
api_version: capability.manifest/v1
kind: capability_manifest
metadata:
  capability_id: cap.document.translate
  capability_version: 1.4.0
  manifest_schema_version: v1
  display_name: Document translation
  owner_ref: platform/capabilities
spec:
  contracts:
    contract_schema_version: v3
    supported_input_contracts:
      - contract.translate.input.v3
    supported_output_contracts:
      - contract.translate.output.v3
    supported_result_types:
      - read_only
      - mutation
    operation_refs:
      - translate
  requirements:
    required_policy_profiles:
      - policy.capability.use
    required_classification_level: internal
    approval_mode: optional
    approval_profile_ref: approval.default
    enforcement_hints:
      - policy_pre_execution
  execution:
    binding_mode: provider_bound
    allowed_environments:
      - staging
      - prod
    required_provider_capabilities:
      - remote-worker
  packaging:
    bundle_digest: sha256:1111111111111111111111111111111111111111111111111111111111111111
    bundle_media_type: application/vnd.opita.capability.bundle.v1+json
    required_attachments:
      - worker-image
      - contract-pack
    signature_ref: sigstore://cap.document.translate/1.4.0
    provenance_ref: provenance://cap.document.translate/1.4.0
  compatibility:
    contract_schema_versions:
      - v3
    binding_versions:
      - v1
    provider_runtime_versions:
      - remote-worker/v1
    runtime_constraints:
      - temporal-adapter/v1
status:
  lifecycle_status: valid
```

### Reglas normativas del manifest

1. `capability_id + capability_version` identifican una capability versionada única.
2. `bundle_digest` es obligatorio y debe apuntar al artifact inmutable exacto.
3. `required_attachments[]` forma parte del contrato de resolubilidad; si falta un attachment requerido, no hay resolución válida.
4. `supported_result_types[]` debe ser consistente con A.3 y con el contrato compilado que intente usar la capability.
5. `lifecycle_status = retired` bloquea resolución nueva.
6. `manifest_schema_version` incompatible bloquea carga del manifest en el registry.

## Binding model mínimo

El binding mínimo implementable DEBE incluir:

- `binding_id`
- `capability_id`
- `bundle_digest`
- `provider_ref`
- `environment`
- `supported_contract_versions`
- `supported_result_types`
- `provider_runtime_version`
- `status`

### Schema YAML mínimo v1

```yaml
api_version: capability.binding/v1
kind: binding
metadata:
  binding_id: bnd-cap-document-translate-prod-001
  binding_version: v1
spec:
  capability_id: cap.document.translate
  capability_version: 1.4.0
  bundle_digest: sha256:1111111111111111111111111111111111111111111111111111111111111111
  provider_ref: provider://translation-worker/prod/v2026-03-29
  environment: prod
  supported_contract_versions:
    - v3
  supported_result_types:
    - read_only
    - mutation
  provider_runtime_version: remote-worker/v1
  effective_from: 2026-03-29T00:00:00Z
  effective_until: null
status:
  status: active
  approval_state: approved
```

### Reglas normativas del binding

1. Un binding referencia exactamente una `capability_id` + `capability_version` + `bundle_digest`.
2. `provider_ref` es obligatorio y debe existir en el catálogo/aprobación de providers permitido por plataforma.
3. `environment` es parte del material resoluble; no puede ignorarse.
4. `supported_contract_versions[]` y `supported_result_types[]` deben ser subconjunto compatible del manifest.
5. `provider_runtime_version` debe satisfacer simultáneamente lo declarado por manifest y binding.
6. `status` puede ser sólo `active`, `disabled`, `deprecated` o `retired`.
7. Sólo `active` es resoluble por defecto.
8. `deprecated` sigue siendo resoluble si continúa compatible y no existe policy que lo prohíba.
9. `disabled` y `retired` no son resolubles.
10. Si `effective_until` expiró, el binding se trata como no resoluble aunque el status diga `active`.

## Relation exacta `capability_manifest -> bundle_digest -> binding -> provider_ref`

La cadena canónica obligatoria queda fijada así:

1. **`capability_manifest`** define identidad, contratos, tipos de resultado, requirements y compatibilidad.
2. **`bundle_digest`** fija el artifact inmutable exacto que contiene ese manifest y sus attachments requeridos.
3. **`binding`** asocia ese manifest versionado y ese digest con un `provider_ref` para un `environment` concreto.
4. **`provider_ref`** identifica el worker/provider remoto concreto que el runtime puede invocar.

### Regla dura de resolución

Si cualquiera de estos eslabones falta, no coincide, está incompatible, está stale, no está aprobado o carece de evidencia mínima de artifact verification, la resolución debe fallar cerrada.

## Registry mínimo implementable

### Componentes mínimos esperados

1. `capability_registry`
2. `manifest_loader`
3. `bundle_resolver`
4. `binding_resolver`
5. `provider_locator`
6. `compatibility_guard`
7. `artifact_verifier`
8. `registry_projection_writer`

### Capacidades mínimas del registry v1

- alta idempotente de manifests válidos;
- alta idempotente de bindings válidos;
- consulta de capabilities por `capability_id` y versión;
- consulta de bindings por environment;
- resolución determinística para runtime;
- rechazo explícito por conflicto o ambigüedad;
- persistencia de `resolution_record`/equivalente canónico e índices derivados;
- soporte a reproyección si falla event log/telemetría.

### Persistencia mínima recomendada

- `capability_manifest_record`
- `bundle_reference_record`
- `binding_record`
- `resolution_record`
- `artifact_verification_record`
- `provider_reference_record` o vista aprobada equivalente

## Resolution flow del runtime contra el registry

### Input mínimo del runtime

El runtime debe resolver con, como mínimo:

- `capability_id`
- `environment`
- `contract_schema_version`
- `required_result_type`
- `compiled_contract_id`
- `trace_id`

### Flow canónico v1

1. `execution_workflow` o su adapter de preparación solicita resolución al `capability_registry`.
2. `capability_registry` busca el `capability_manifest` vigente para `capability_id` y versión objetivo.
3. `manifest_loader` valida `manifest_schema_version`, `lifecycle_status`, grupos mínimos obligatorios y consistencia declarativa.
4. `bundle_resolver` obtiene el `bundle_digest` asociado al manifest.
5. `artifact_verifier` confirma que el digest existe y que hay evidencia mínima de firma/provenance registrada.
6. `bundle_resolver` verifica presencia de `required_attachments[]`.
7. `binding_resolver` busca bindings del mismo `capability_id` + `capability_version` + `bundle_digest` para el `environment` pedido.
8. `binding_resolver` descarta bindings `disabled`, `retired`, expirados o stale.
9. `compatibility_guard` verifica contract version, result type, binding version y provider runtime version.
10. `provider_locator` confirma que `provider_ref` exista y esté aprobado.
11. Si queda exactamente un binding resoluble, el registry produce resolución canónica.
12. `registry_projection_writer` persiste la resolución y sus correlaciones canónicas.
13. Se proyecta a event log/observabilidad de forma derivada.
14. Si la proyección derivada falla, la resolución canónica ya persistida sigue siendo válida.

### Output mínimo recomendado

```yaml
resolution_id: reso-...
capability_id: cap.document.translate
capability_version: 1.4.0
bundle_digest: sha256:1111111111111111111111111111111111111111111111111111111111111111
binding_id: bnd-cap-document-translate-prod-001
provider_ref: provider://translation-worker/prod/v2026-03-29
environment: prod
contract_schema_version: v3
result_type: mutation
verification_state:
  digest_verified: true
  signature_evidence_present: true
  provenance_evidence_present: true
```

## Compatibilidad de versiones

La compatibilidad explícita obligatoria en C.5 queda entre:

- `capability_version`
- `manifest_schema_version`
- `binding_version`
- `provider_runtime_version`
- `contract_schema_version`

### Reglas v1

1. `capability_version` versiona la semántica funcional de la capability.
2. `manifest_schema_version` versiona la forma del manifest.
3. `binding_version` versiona la semántica del binding.
4. `provider_runtime_version` versiona el runtime esperado del provider remoto.
5. `contract_schema_version` debe ser soportado por manifest y binding al mismo tiempo.

### Matriz mínima de compatibilidad

| Relación | Regla v1 |
|---|---|
| `capability_version` ↔ `manifest_schema_version` | una capability puede tener nuevas versiones sobre el mismo schema; si el schema cambia de forma incompatible, el loader debe conocer explícitamente la nueva versión |
| `capability_version` ↔ `binding_version` | un binding debe apuntar a una capability versionada explícita; no existe binding `latest` |
| `manifest_schema_version` ↔ `binding_version` | ambos deben ser entendibles por la versión del registry activa |
| `binding_version` ↔ `contract_schema_version` | el binding sólo resuelve contratos incluidos en `supported_contract_versions[]` |
| `provider_runtime_version` ↔ `binding_version` | el binding sólo es resoluble si declara una versión de runtime soportada por el provider locator |
| `provider_runtime_version` ↔ `contract_schema_version` | el provider runtime no puede ejecutar contracts fuera de la matriz soportada |

### Regla de cambios incompatibles

Si un cambio requiere reinterpretar manifest, binding o runtime sin compatibilidad declarada, debe emitirse nueva versión explícita. NO se permite compatibilidad implícita.

## Estrategia de validación y verificación de artifacts

La validación de artifact en C.5 se divide en cuatro pasos:

1. **Schema validation** de manifest y binding.
2. **Reference validation** de `bundle_digest`, `provider_ref` y `required_attachments[]`.
3. **Digest verification** del artifact registrado.
4. **Evidence validation** de firma y provenance.

### Alcance explícito v1

- C.5 exige **presencia y correlación verificable** de evidencia de firma/provenance.
- C.5 NO exige todavía implementar toda la verificación criptográfica end-to-end.
- C.5 SÍ exige que el resolution path no ignore ausencia de digest/signature/provenance evidence.

## Persistencia y correlación operativa

### Records mínimos recomendados

- `capability_manifest_record`
- `binding_record`
- `artifact_verification_record`
- `resolution_record`

### Correlación mínima obligatoria por `resolution_record`

- `resolution_id`
- `capability_id`
- `capability_version`
- `bundle_digest`
- `binding_id`
- `provider_ref`
- `environment`
- `contract_schema_version`
- `compiled_contract_id` cuando aplique
- `execution_id` cuando aplique
- `trace_id`
- `resolution_status`
- `reason_codes[]`
- `resolved_at`

### Reglas v1

1. El `resolution_record` representa el hecho canónico de que el registry resolvió o rechazó una capability bajo un input material dado.
2. Debe ser idempotente por input material equivalente.
3. Debe poder correlacionarse con `compiled_contract`, `execution_workflow` y `event_record`.
4. Las proyecciones secundarias no pueden duplicar ni alterar la resolución canónica persistida.

## Integración con `compiled_contract`

1. El registry NO recompila ni modifica el `compiled_contract`.
2. El `compiled_contract` entrega al registry `capability_id`, `contract_schema_version`, `required_result_type`, `environment` y correlación.
3. La resolución sólo es válida si el `contract_schema_version` pedido es compatible con manifest y binding.
4. Si el `compiled_contract` refiere una capability retirada, incompatible o sin binding resoluble, el runtime debe fallar antes de iniciar trabajo de provider.
5. El `resolution_record` debe referenciar `compiled_contract_id` o `contract_id` cuando la resolución nace de ejecución real.

## Integración con `execution_workflow`

1. `execution_workflow` sigue siendo la autoridad durable del lifecycle de ejecución según C.2.
2. La resolución de capability ocurre **antes** de invocar el provider remoto efectivo.
3. Un fallo de resolución debe reflejarse como rechazo/bloqueo/falla temprana compatible con C.2.
4. `execution_workflow` consume la salida del registry, pero no reemplaza ni duplica la lógica de resolución.
5. Reintentos técnicos del workflow deben reutilizar la resolución persistida si el input material no cambió.
6. Si cambia el input material, debe generarse nueva resolución.

## Integración con policy y approvals

1. C.5 no reemplaza C.3 ni A.4.
2. El registry puede exigir guardas normativas mínimas como “provider aprobado” o “binding activo”, pero la policy contextual sigue viviendo en el seam de policy.
3. Si una capability requiere approvals o clasificación según manifest, esa exigencia debe reflejarse en metadata resoluble y ser consumible por policy/runtime.
4. `provider_ref` no aprobado debe bloquear resolución.
5. Un binding `deprecated` puede seguir resolviendo si policy/approval no lo prohíben y la compatibilidad sigue siendo válida.
6. Un binding `retired` no puede resolverse aunque exista provider aprobado.

## Integración con event log y observabilidad

1. La resolución debe poder emitir eventos derivados como `capability.resolution_requested`, `capability.resolution_succeeded` y `capability.resolution_rejected`.
2. Esos eventos son derivados y no reemplazan `resolution_record`.
3. La correlación mínima debe incluir `trace_id`, `capability_id`, `bundle_digest`, `binding_id`, `provider_ref`, `contract_id/compiled_contract_id` y `execution_id` cuando aplique.
4. Si el event log falla pero `resolution_record` persiste, la resolución canónica sigue siendo válida.
5. Retry de proyección no debe duplicar records derivados ni cambiar el resultado canónico ya persistido.

## Idempotencia y deduplicación del registry

### Deduplicación de catálogo

- Key mínima recomendada: `capability_id + capability_version + bundle_digest`

### Deduplicación de bindings

- Key mínima recomendada: `capability_id + capability_version + bundle_digest + environment + provider_ref + binding_version`

### Deduplicación de resolución

- Key mínima recomendada: `capability_id + capability_version + bundle_digest + binding_id + provider_ref + environment + contract_schema_version + result_type + compiled_contract_id`

### Reglas v1

1. Repetir el mismo manifest material no crea dos entradas canónicas.
2. Repetir el mismo binding material no crea dos bindings resolubles.
3. Retry técnico con mismo input material reutiliza el mismo `resolution_record` o una reproyección del mismo.
4. Retry de event log/telemetría no duplica la resolución canónica.
5. Una colisión material incompatible debe rechazarse o marcar conflicto explícito.

## Tests borde mínimos (al menos 20)

1. capability existe pero no tiene binding resoluble.
2. bundle digest inexistente.
3. provider_ref no aprobado.
4. binding apunta a environment incompatible.
5. manifest schema version incompatible.
6. binding version incompatible con contract version.
7. capability retirada pero todavía referenciada.
8. bundle firmado pero manifest inválido.
9. manifest válido pero attachment requerido faltante.
10. runtime intenta resolver capability global sin binding vigente.
11. provider runtime version fuera de compatibilidad.
12. bundle digest correcto pero binding stale.
13. duplicate registry entry mismo capability+version+binding.
14. resolution ambigua con dos bindings activos.
15. binding disabled durante resolución.
16. manifest/result type incompatibles.
17. event log falla pero resolution canónica persiste.
18. retry de resolution no duplica records derivados.
19. capability marcada deprecated pero todavía compatible.
20. bundle digest verificado pero signature provenance no disponible.
21. binding `retired` todavía presente en storage pero no resoluble.
22. binding `deprecated` único y compatible sigue resolviendo con evidencia de deprecación.
23. `supported_contract_versions[]` del binding no incluye la versión del `compiled_contract`.
24. `supported_result_types[]` del binding excluye el result type pedido por runtime.
25. manifest declara `allowed_environments` que no incluye el environment del binding.
26. digest coincide pero attachment requerido pertenece a otra versión del bundle.
27. provider existe pero `provider_runtime_version` declarada por binding no coincide con provider registrado.
28. dos manifests comparten `capability_id + capability_version` pero difieren en `bundle_digest` sin versión nueva explícita.
29. `resolution_record` canónico persiste y falla la reproyección asíncrona a observabilidad.
30. retry concurrente del mismo registro de binding no crea dos rows resolubles.

## Criterios de aceptación de C.5

1. `capability` permanece explícitamente como objeto canónico gobernado centralmente.
2. Queda explícito que el registry responde por **catálogo + resolución**, NO por lifecycle durable ni distribution.
3. Queda explícito que tenant activation y overlays tenant quedan fuera de C.5 y fuera del roadmap actual.
4. La cadena `capability_manifest -> bundle_digest -> binding -> provider_ref` queda fijada sin ambigüedad.
5. El runtime queda obligado a resolver capability usando registry + binding válido, no por mera existencia global.
6. El documento separa claramente registry, packaging, binding y provider resolution.
7. Existe manifest schema ejecutable mínimo con identidad, contracts/result types, requirements, execution binding refs, packaging refs y compatibility metadata.
8. Existe binding model mínimo con `binding_id`, `capability_id`, `bundle_digest`, `provider_ref`, `environment`, `supported_contract_versions`, `supported_result_types`, `provider_runtime_version` y `status`.
9. Queda definida compatibilidad explícita entre `capability_version`, `manifest_schema_version`, `binding_version`, `provider_runtime_version` y `contract_schema_version`.
10. Queda exigida verificación de digest y evidencia de firma/provenance como parte del resolution path.
11. Quedan definidos los componentes mínimos: `capability_registry`, `manifest_loader`, `bundle_resolver`, `binding_resolver`, `provider_locator`, `compatibility_guard`, `artifact_verifier` y `registry_projection_writer`.
12. Queda definida persistencia y correlación operativa mínima del seam.
13. Queda definida integración con `compiled_contract`, `execution_workflow`, policy, approvals, event log y observabilidad.
14. Quedan definidas reglas de idempotencia y deduplicación para manifest, binding y resolution.
15. Existe batería mínima de tests borde suficiente para implementar el seam v1 sin reabrir packaging ni distribution layer.
