# C.1 — Intent → Contract compiler

## Objetivo

Definir el compilador que toma un `intent_input` gobernado y lo transforma en un `compiled_contract` persistido en PostgreSQL junto con un `compilation_report`, dejando una salida determinística, auditable y lista para que C.2 cree `execution_record` sin reinterpretación semántica.

Este bloque profesionaliza el primer seam operativo del kernel: la conversación, la propuesta o el patchset gobernado pueden originar intención, pero el runtime durable sólo puede arrancar desde contrato compilado persistido.

## Principios de implementación del compilador

- El compilador es un **pipeline puro hasta persistencia**; no ejecuta side effects de negocio, no llama providers, no dispara workflows ni mutaciones externas.
- El compilador **sí puede leer** contexto, memoria operativa, snapshots y metadata de policy para enriquecer el contrato.
- El compilador **no puede decidir approvals finales** ni reemplazar a Cerbos ni al runtime.
- El compilador debe tener una **representación intermedia explícita** entre `intent_input` y `compiled_contract`.
- El compilador debe ser **determinístico** para el mismo `intent_input` normalizado y el mismo conjunto de snapshots/versiones relevantes.
- El compilador debe separar **errores fatales**, **faltantes de input**, **warnings** y **anotaciones diagnósticas**.
- El compilador produce dos salidas formales: `compiled_contract` persistido + `compilation_report`.
- La semántica material del contrato debe quedar cerrada antes de persistir; C.2 puede ejecutar, pero no reinterpretar significado.

## Boundary exacto del compilador

El compilador empieza cuando existe un `intent_input` ya gobernado, identificable por `tenant_id`, `subject_id`, `session_id` y evidencia de origen. Termina cuando:

1. se resolvió una representación intermedia estable;
2. se validó el contrato compilable contra invariantes de A.2;
3. se capturaron snapshots y versiones relevantes;
4. se calculó fingerprint material determinístico;
5. se persistió en PostgreSQL el `compiled_contract` y su `compilation_report`.

Queda **dentro** del boundary:

- parseo de intención estructurada;
- normalización semántica;
- enriquecimiento con contexto/memoria/snapshots;
- resolución de restricciones, capability y riesgo;
- preparación de input canonizado a Cerbos;
- validación de compilación;
- snapshotting, fingerprint y deduplicación;
- persistencia transaccional del resultado.

Queda **fuera** del boundary:

- evaluación final de policy en Cerbos;
- emisión de approval decision;
- generación de `execution_record`;
- arranque de workflows Temporal;
- ejecución técnica o aplicación de negocio;
- conversación con el usuario para pedir aclaraciones en tiempo real.

## Responsabilidades explícitas

- Recibir `intent_input` gobernado y rechazar entradas inválidas en recepción.
- Construir `normalized_intent` como forma canónica de entrada.
- Leer snapshots/contexto necesarios para enriquecer sin mutar el mundo.
- Resolver `capability_id`, `capability_ref`, restricciones efectivas, tools permitidas/bloqueadas y clasificación/riesgo preliminar según A.2/A.4.
- Construir `policy_evaluation_input` estable para Cerbos, sin ejecutar la evaluación final.
- Validar que la salida sea **compilable**, aunque todavía no sea **ejecutable**.
- Congelar snapshots mínimos requeridos por A.2.
- Calcular fingerprint sobre campos materiales normalizados.
- Aplicar deduplicación/idempotencia basada en fingerprint + versiones/snapshots relevantes.
- Persistir `compiled_contract` y `compilation_report` en PostgreSQL con trazabilidad.

## No-responsabilidades explícitas

- No decidir `allow/deny` final de policy.
- No crear approvals ni marcar approvals como válidas.
- No construir ni aprobar `plan_snapshot`.
- No resolver scheduling, retries o timers de runtime.
- No ejecutar tools ni providers.
- No promover un contrato a `executable`; esta salida puede quedar en `compiled`.
- No usar telemetría derivada como source of truth.
- No reabrir comparativas de Fase B ni redefinir el modelo de A.2.

## Inputs detallados del compilador

### 1. `intent_input`

Objeto de entrada gobernado con, como mínimo:

- `intent_id`
- `tenant_id`
- `workspace_id`
- `environment`
- `session_id`
- `user_id`
- `acting_for_subject_id` opcional
- `delegation_id` opcional
- `objetivo`
- `alcance`
- `tipo_de_resultado_esperado`
- `autonomia_solicitada`
- `aprobacion_requerida` opcional
- `criterios_de_exito` opcional
- `restricciones_declaradas` opcional
- `notas_del_usuario` opcional
- `source_ref` con referencia al proposal/patchset/origen gobernado

### 2. `compiler_context`

Metadata no editable por el usuario necesaria para compilar:

- `request_id`
- `compilation_id`
- `received_at`
- `compiler_version`
- `schema_version`
- `expected_contract_state = compiling`

### 3. `reference_snapshots`

Lecturas de contexto y snapshots requeridos para la compilación:

- `tenant_snapshot`
- `subject_snapshot`
- `delegation_snapshot` si aplica
- `policy_snapshot_ref`
- `capability_registry_snapshot`
- `permission_snapshot_source`
- `memory_context_source`
- `classification_inputs`
- `risk_inputs`

### 4. `system_constraints`

Restricciones del sistema y baseline heredado:

- floors del tipo de resultado (A.3)
- reglas de riesgo/clasificación (A.4)
- reglas de materialidad/versionado (A.2)
- availability del tenant/capability/subject
- compatibilidad environment/capability

## Outputs detallados del compilador

### 1. `compiled_contract`

Objeto persistido en PostgreSQL con forma canónica alineada a A.2:

- campos del usuario normalizados;
- campos compilados del sistema resueltos;
- campos técnicos (`contract_id`, `contract_version`, `fingerprint`, `compiled_at`, etc.);
- snapshots inmutables capturados durante compilación;
- `contract_state = compiled` si la compilación cierra bien;
- `contract_state = incomplete` si la compilación detecta input faltante resoluble por usuario;
- nunca `executable` desde C.1.

### 2. `compilation_report`

Objeto de diagnóstico persistido y retornable al caller. Debe incluir como mínimo:

- `compilation_id`
- `contract_id`
- `compiler_version`
- `status` = `compiled` | `incomplete` | `rejected` | `failed`
- `reason_codes[]`
- `errors[]`
- `warnings[]`
- `missing_user_inputs[]`
- `material_fields_used_for_fingerprint[]`
- `normalized_input_hash`
- `snapshot_versions_used`
- `deduplication_result` = `new` | `reused` | `conflict`
- `persistence_result`
- `started_at`
- `finished_at`

### 3. `policy_evaluation_input`

Subobjeto derivado y persistido/referenciado dentro del contrato o reporte, listo para que C.3/C.2 consulten Cerbos sin reinterpretar significado.

## Componentes internos mínimos del compilador

1. `intent_parser`
   - valida shape de entrada;
   - rechaza intención vacía o estructura imposible;
   - materializa `parsed_intent`.

2. `intent_normalizer`
   - resuelve defaults, enums, orden canónico y equivalencias semánticas;
   - produce `normalized_intent`.

3. `context_enricher`
   - lee memoria operativa, snapshots y metadata de tenant/sujeto/capability;
   - produce `enriched_context`.

4. `constraint_resolver`
   - combina restricciones del usuario y del sistema;
   - resuelve capability, tools, permisos y conflictos;
   - produce `resolved_constraints`.

5. `risk_classifier`
   - calcula riesgo y clasificación preliminar/compilada;
   - aplica floors y señales restrictivas;
   - produce `risk_and_classification`.

6. `policy_input_builder`
   - canoniza principal, recurso, acción y contexto;
   - produce `policy_evaluation_input`.

7. `contract_validator`
   - verifica invariantes de recepción/compilación de A.2;
   - decide `compiled` vs `incomplete` vs `rejected`.

8. `snapshot_builder`
   - congela snapshots y versiones exactas utilizadas;
   - produce `compilation_snapshots`.

9. `fingerprint_calculator`
   - calcula fingerprint sólo sobre campos materiales normalizados;
   - produce `contract_fingerprint`.

10. `contract_persister`
    - persiste contrato + reporte en PostgreSQL de forma transaccional;
    - aplica deduplicación por fingerprint.

11. `compilation_report_builder`
    - agrega diagnósticos, reason codes y metadata del pipeline;
    - produce `compilation_report` final.

## Pipeline de compilación paso a paso

### Paso 1 — Recepción controlada

- validar presencia de IDs técnicos mínimos;
- validar tenant activo, sujeto activo y environment permitido;
- rechazar si falta `objetivo` o si `tipo_de_resultado_esperado` es inválido.

Salida: `parsed_intent` o rechazo de recepción.

### Paso 2 — Parseo semántico

- separar campos del usuario, técnicos y referencias externas;
- detectar nulidad semántica, strings vacíos, listas inválidas, enums fuera de catálogo.

Salida: `parsed_intent`.

### Paso 3 — Normalización

- canonizar enums, casing, defaults, sort de listas semánticamente no ordenadas;
- normalizar restricciones declaradas y notas sin volverlas materiales;
- derivar un `normalized_input_hash` para trazabilidad del input normalizado.

Salida: `normalized_intent`.

### Paso 4 — Enriquecimiento de contexto

- leer tenant, sujeto, capability registry, memoria operativa, políticas activas y delegación;
- recuperar `contexto_relevante` estructurado;
- registrar vacío explícito si el enrichment no devuelve nada útil.

Salida: `enriched_context`.

### Paso 5 — Resolución de restricciones

- intersectar restricciones declaradas por usuario con restricciones del sistema;
- resolver permisos efectivos;
- filtrar systems/tools incompatibles o bloqueadas;
- determinar si la intención sigue siendo compilable o queda en conflicto.

Salida: `resolved_constraints`.

### Paso 6 — Resolución de capability

- confirmar capability explícita o seleccionar una resoluble según objetivo y tipo;
- fijar `capability_id` y `capability_ref` exactos;
- fallar si no existe capability resoluble.

Salida: capability resuelta o error.

### Paso 7 — Clasificación y riesgo

- calcular `business_risk_score`, `security_risk_score` y `nivel_de_riesgo`;
- calcular `classification_level` desde datos, sistemas y capability;
- si faltan insumos de riesgo, asumir peor caso según A.2;
- elevar clasificación si aparecen datos restringidos.

Salida: `risk_and_classification`.

### Paso 8 — Preparación de inputs a policy

- construir `policy_evaluation_input` canonizado;
- adjuntar `principal`, `resource`, `action`, `tenant_scope`, riesgo, clasificación, contract state y atributos materiales relevantes;
- registrar versión/snapshot de policy usado como referencia.

Salida: `policy_evaluation_input`.

### Paso 9 — Validación del contrato compilable

- verificar sistemas confirmados;
- verificar tools permitidas;
- verificar approval floor mínimo;
- verificar consistencia entre restricciones del usuario y del sistema;
- verificar snapshots obligatorios para compilar.

Salida: `validation_result` con `compiled`, `incomplete` o `rejected`.

### Paso 10 — Construcción de snapshots

- congelar `policy_snapshot`, `classification_snapshot`, `risk_snapshot`, `permission_snapshot`, `context_snapshot`, `delegation_snapshot` si aplica;
- asociar versiones exactas y timestamps de captura.

Salida: `compilation_snapshots`.

### Paso 11 — Cálculo de fingerprint

- seleccionar sólo campos materiales normalizados;
- serializar en orden estable;
- hashear con algoritmo estándar del sistema;
- comparar con contratos existentes equivalentes.

Salida: `contract_fingerprint`.

### Paso 12 — Persistencia transaccional

- persistir `compiled_contract` + `compilation_report` en PostgreSQL;
- si fingerprint coincide con uno ya persistido para mismo tenant/scope/versiones relevantes, reutilizar o referenciar sin duplicar;
- si falla persistencia, no debe quedar contrato “medio escrito”.

Salida: persistencia exitosa o falla fatal.

## Contrato interno entre etapas

La implementación v1 debe usar DTOs internos explícitos entre etapas. Mínimos:

### `parsed_intent`

- refleja exactamente lo recibido, validado estructuralmente;
- todavía no resuelve equivalencias ni defaults.

### `normalized_intent`

- forma canónica del input;
- base para idempotencia del compilador;
- todavía no incluye campos derivados de contexto externo.

### `enriched_context`

- agrega datos leídos de memoria, tenant, sujeto, policy, permissions, capability registry y delegación;
- nunca agrega side effects, sólo lecturas.

### `resolved_constraints`

- expresa intersección final de restricciones, herramientas y accesos permitidos/bloqueados;
- hace visibles contradicciones.

### `risk_and_classification`

- encapsula scores, nivel resultante, clasificación y evidence refs.

### `compilation_candidate`

- objeto previo al contrato final;
- incluye todos los campos del futuro `compiled_contract` excepto fingerprint/persistencia.

### `validated_compilation_candidate`

- `compilation_candidate` + resultado formal de validación;
- define si la salida es `compiled`, `incomplete` o `rejected`.

### `persistable_compiled_contract`

- candidato validado + snapshots + fingerprint + metadata de versionado listo para PostgreSQL.

## Persistencia y versionado del contrato compilado

- La persistencia del `compiled_contract` va a **PostgreSQL** como memoria operativa durable del core.
- `compiled_contract` y `compilation_report` deben persistirse en la misma transacción lógica.
- El almacenamiento debe preservar:
  - `contract_id`
  - `contract_version`
  - `parent_contract_id` si deriva de otro
  - `fingerprint`
  - `contract_state`
  - snapshots completos o referencias durables a snapshots
  - `compiler_version`
  - `source_ref`

### Política de versionado

- `major.minor` sigue A.2.
- C.1 no decide superseding de estados posteriores, pero sí fija la versión compilada actual.
- cambio material detectado durante recompilación incrementa `major`;
- cambio no material incrementa `minor`;
- si el fingerprint y snapshots/versiones relevantes coinciden, no nace una nueva versión material.

## Política de fingerprint y deduplicación

### Campos base del fingerprint

El fingerprint se calcula sólo sobre campos **materiales normalizados**. Debe incluir, al menos, cuando apliquen:

- `objetivo`
- `alcance`
- `tipo_de_resultado_esperado`
- `sistemas_confirmados`
- `datos_permitidos`
- `herramientas_permitidas`
- `capability_id`
- `classification_level`
- `nivel_de_riesgo`
- `approval_mode_efectivo`
- `policy_snapshot.version`
- `destination_snapshot` si aplica
- `criterios_de_exito` cuando sea material según A.2

### Campos excluidos del fingerprint

- `notas_del_usuario`
- `notas_de_contexto`
- metadata cosmética o de presentación
- ordenamiento no semántico de listas
- IDs técnicos no materiales para semántica operativa

### Reglas de deduplicación

- mismo `tenant_id` + mismo `normalized_intent` + mismo set de snapshots/versiones relevantes + mismos campos materiales => mismo fingerprint;
- mismo fingerprint => el compilador debe reutilizar contrato persistido o registrar intento como deduplicado, no duplicar artefacto material;
- cambio no material NO cambia fingerprint;
- cambio material SÍ cambia fingerprint;
- si las versiones de snapshots relevantes cambian, puede cambiar fingerprint aunque el input del usuario no cambie;
- la deduplicación debe ser tenant-scoped.

## Diagnósticos y errores de compilación

El compilador debe emitir reason codes normalizados. Mínimos v1:

### Rechazo de recepción

- `compile.input.empty_intent`
- `compile.input.invalid_result_type`
- `compile.tenant.inactive`
- `compile.subject.inactive`
- `compile.environment.invalid`

### Incompleto

- `compile.input.missing_scope`
- `compile.input.missing_success_criteria`
- `compile.context.insufficient_user_input`
- `compile.tools.none_permitted`

### Error de compilación

- `compile.capability.unresolvable`
- `compile.snapshot.missing_required`
- `compile.policy.snapshot_stale`
- `compile.delegation.expired`
- `compile.constraints.conflict_user_system`
- `compile.persistence.failed`

### Warnings diagnósticos

- `compile.context.empty_enrichment`
- `compile.risk.assumed_worst_case`
- `compile.output.compiled_not_executable`
- `compile.dedup.reused_existing_contract`

Reglas:

- `compilation_report` puede contener múltiples reason codes;
- debe separar `fatal`, `recoverable`, `warning`;
- si el contrato es compilable pero no ejecutable, eso es **warning o nota diagnóstica**, no error fatal de C.1.

## Integración con memoria operativa

- PostgreSQL conserva la verdad operativa del `compiled_contract` y del `compilation_report`.
- El compilador puede leer memoria operativa relevante para poblar `contexto_relevante` y `context_snapshot`.
- La lectura de memoria debe ser estructurada y versionable; no texto libre sin shape.
- Si el enrichment devuelve vacío, el compilador debe registrar `compile.context.empty_enrichment`, pero no fallar automáticamente.
- La memoria se usa para enriquecer y explicar, no para decidir side effects ni approvals finales.

## Integración con Cerbos (solo preparación de inputs, no evaluación final)

- El compilador prepara `policy_evaluation_input`, pero no invoca la evaluación final que determine `allow/deny` de ejecución.
- Debe dejar canonizados:
  - `principal`
  - `resource`
  - `action`
  - `tenant_scope`
  - atributos de clasificación, riesgo, capability y contract state
  - correlación (`tenant_id`, `contract_id`, `session_id`, `capability_ref`)
- Debe persistir la referencia de policy snapshot/version usada para compilar.
- Si el snapshot de policy requerido está desactualizado o inconsistente, la compilación falla con reason code explícito.

## Integración con runtime/Temporal

- La salida de C.1 debe quedar lista para que C.2 cree `execution_record` sin reinterpretación semántica.
- Temporal no participa en la compilación; participa después, consumiendo contrato compilado persistido.
- El compilador debe producir metadata suficiente para runtime:
  - `contract_id`
  - `contract_version`
  - `fingerprint`
  - `capability_ref`
  - `approval_mode_efectivo`
  - `classification_level`
  - `nivel_de_riesgo`
  - snapshots relevantes referenciables
- Si el contrato está `compiled` pero todavía no `executable`, runtime no puede arrancar ejecución real; C.2 deberá crear el corredor correcto posterior.

## Idempotencia del compilador

- El compilador es idempotente para el mismo `intent_input` normalizado + mismo set de snapshots/versiones relevantes.
- Idempotencia significa:
  - mismo fingerprint;
  - mismo `compiled_contract` material;
  - no duplicación de persistencia material;
  - mismo `compilation_report.status`, salvo metadata temporal no material.
- Si cambia cualquier snapshot/version relevante, la compilación debe tratarse como nueva evaluación material.
- El retry técnico posterior a cálculo de fingerprint no puede crear contratos duplicados si el fingerprint coincide.

## Tests borde mínimos (al menos 18)

1. **intención vacía**
   - esperado: rechazo con `compile.input.empty_intent`.

2. **`tipo_de_resultado_esperado` inválido**
   - esperado: rechazo con `compile.input.invalid_result_type`.

3. **tenant inactivo**
   - esperado: rechazo de recepción con `compile.tenant.inactive`.

4. **sujeto inactivo**
   - esperado: rechazo de recepción con `compile.subject.inactive`.

5. **capability no resoluble**
   - esperado: falla con `compile.capability.unresolvable`.

6. **mismos inputs producen mismo fingerprint**
   - esperado: fingerprint idéntico e intento deduplicado.

7. **cambio no material NO cambia fingerprint**
   - esperado: mismo fingerprint, `minor` o sólo actualización diagnóstica.

8. **cambio material sí cambia fingerprint**
   - esperado: fingerprint distinto y nueva versión material.

9. **falta snapshot requerido**
   - esperado: falla con `compile.snapshot.missing_required`.

10. **context enrichment devuelve vacío**
    - esperado: compilación sigue si lo demás alcanza; warning `compile.context.empty_enrichment`.

11. **policy snapshot desactualizado**
    - esperado: falla con `compile.policy.snapshot_stale`.

12. **delegation expirada**
    - esperado: falla o incompleto según dependencia de scope; reason code `compile.delegation.expired`.

13. **clasificación sube por dato restringido**
    - esperado: `classification_level` se eleva al nivel más restrictivo.

14. **tools permitidas vacías**
    - esperado: `incomplete` o `failed`, nunca `compiled` exitoso para ejecutar.

15. **restricción del usuario contradice restricción del sistema**
    - esperado: `compile.constraints.conflict_user_system`.

16. **output compilado no es ejecutable pero sí compilable**
    - esperado: `compiled_contract` en `compiled` + warning `compile.output.compiled_not_executable`.

17. **retry de compilación no duplica contrato si fingerprint coincide**
    - esperado: deduplicación tenant-scoped; no duplica fila material.

18. **persistencia falla después de calcular fingerprint**
    - esperado: `compile.persistence.failed`; sin contrato parcialmente persistido.

19. **riesgo incompleto obliga peor caso**
    - esperado: `nivel_de_riesgo = critical` + warning `compile.risk.assumed_worst_case`.

20. **reordenamiento cosmético de lista no cambia fingerprint**
    - esperado: normalización lo estabiliza.

21. **mismo intent normalizado pero distinta versión de policy cambia identidad material**
    - esperado: fingerprint distinto si `policy_snapshot.version` es materialmente distinta.

22. **capability resoluble pero bloqueada por constraints efectivas**
    - esperado: falla o `incomplete`; nunca contrato compilado engañosamente válido.

## Criterios de aceptación de C.1

1. Queda explícito que el compilador es un pipeline puro hasta persistencia y no ejecuta side effects de negocio.
2. Existe una representación intermedia clara entre `intent_input` y `compiled_contract`.
3. Quedan definidos los once componentes internos mínimos del compilador y su responsabilidad.
4. El compilador produce formalmente `compiled_contract` persistido + `compilation_report`.
5. La persistencia del resultado queda fijada en PostgreSQL y no contradice el split-plane de Fase B.
6. El fingerprint queda definido sobre campos materiales normalizados, excluyendo notas cosméticas.
7. La política de deduplicación e idempotencia queda explicitada para retries y recompilaciones.
8. Quedan definidos diagnósticos y reason codes suficientes para rechazo, incompleto, falla y warning.
9. La integración con memoria, Cerbos y Temporal queda delimitada sin invadir sus responsabilidades.
10. La salida queda lista para que C.2 cree `execution_record` sin reinterpretación semántica.
