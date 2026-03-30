# C.2 — Execution runtime skeleton

## Objetivo

Definir el skeleton implementable del runtime durable del kernel sobre **Temporal** para que un `compiled_contract` persistido en C.1 pueda materializarse como un `execution_record` durable, inspeccionable, idempotente y auditable, sin colapsar ejecución técnica con aplicación de negocio ni reabrir decisiones ya cerradas en A.5, B.1 y C.1.

Este bloque NO implementa todavía todo el runtime del producto. Fija el corredor mínimo operativo que debe existir para que Fase C pueda seguir con policy, event log y activación sin ambigüedad estructural.

## Principios de implementación del runtime skeleton

- **Temporal es el único runtime durable de verdad** del kernel en Fase C. No se introduce una segunda semántica de workflow ni una capa declarativa que compita con Temporal.
- `execution_record` **no es reemplazado por Temporal**. Sigue siendo el objeto first-class del dominio; Temporal sólo materializa su semántica durable y su progresión real.
- Debe existir **un workflow durable principal por `execution_id`**. Ese workflow es la única autoridad de progreso durable de la ejecución viva.
- Debe existir **separación estricta entre workflow durable y activities técnicas**. El workflow decide lifecycle y transiciones; las activities hacen trabajo acotado y retornan evidencia/resultados.
- Debe existir **separación explícita entre `execution_completed` y `application_completed`**. En mutation no pueden colapsar en una sola mutación de estado.
- El runtime skeleton debe ser **idempotente** ante reintentos de arranque con el mismo `execution_id`.
- Debe existir **estado `blocked` separado de `failed`**.
- Debe existir soporte explícito para **`awaiting_approval`** como pausa durable.
- Debe existir camino explícito para **`unknown_outcome`** cuando un timeout o corte externo rompe la certeza material del efecto.
- La evidencia durable debe correlacionarse, como mínimo, con `trace_id`, `contract_id`, `tenant_id`, `result_id`, `approval_request_id` y `execution_id`.
- El runtime debe nacer listo para integrar Cerbos, event log operativo y observabilidad derivada **sin mover el boundary del workflow principal**.

## Boundary exacto del runtime skeleton

El runtime skeleton empieza cuando existe un `compiled_contract` persistido y elegible para arranque, junto con request técnica para iniciar una ejecución. Termina cuando la ejecución queda en un estado terminal o en un estado no terminal duradero que exige intervención/reanudación externa controlada.

Queda **dentro** del boundary:

- creación idempotente de `execution_record`;
- arranque y recuperación del workflow durable principal;
- validación de precondiciones mínimas de ejecución;
- evaluación del gate de policy en el punto correcto del lifecycle;
- espera durable por approvals o liberaciones externas;
- ejecución técnica y eventual fase de aplicación;
- timers, deadlines, retries y pausas durables;
- detección de `unknown_outcome`;
- coordinación de compensación o escalación manual;
- proyección de estado canónico hacia `execution_record`;
- emisión de eventos operativos derivados del lifecycle.

Queda **fuera** del boundary:

- compilación o reinterpretación semántica del contrato;
- decisión final de policy dentro del workflow como lógica ad hoc;
- modelado profundo de user tasks o bandejas humanas;
- resolución de catálogo/capability registry;
- distribution layer y resolution de capabilities fuera del runtime skeleton;
- analytics, dashboards o telemetría como source of truth;
- implementación final de providers/workers de negocio.

## Responsabilidades explícitas

- Crear o reutilizar idempotentemente el `execution_record` para un `execution_id` dado.
- Iniciar el workflow durable principal `execution_workflow` con correlación estable.
- Validar que el `compiled_contract` existe, está persistido y se encuentra en estado apto para ejecución.
- Materializar transiciones de runtime compatibles con A.5.
- Separar fase de **execution** de fase de **application**.
- Esperar approvals/liberaciones externas mediante señales durables.
- Aplicar timers y retries diferenciando retry técnico, retry de ejecución y replay de auditoría.
- Detectar y marcar `blocked`, `awaiting_approval`, `failed`, `unknown_outcome`, `compensation_pending` y cierres terminales relevantes.
- Persistir/proyectar el estado canónico y la evidencia terminal mínima.
- Exponer queries y signals mínimas para inspección y control operativo.
- Preparar el seam para Cerbos sin incrustar policy final en el runtime.
- Preparar el seam para event log/OTel sin convertirlos en la verdad durable.

## No-responsabilidades explícitas

- No recompilar intención ni regenerar `compiled_contract`.
- No reinterpretar el significado del contrato compilado.
- No decidir policy por cuenta propia cuando el diseño exige PDP externo.
- No reemplazar approvals ni resolver autoridad de aprobadores.
- No persistir observabilidad derivada como fuente primaria de estado.
- No asumir que toda falla implica compensación automática.
- No asumir que toda compensación es rollback físico.
- No reabrir la comparativa de runtime de B.1.
- No permitir que una activity técnica cambie estado canónico por fuera del workflow durable.

## Mapping exacto entre `execution_record` y Temporal

La relación exacta es **1 `execution_record` canónico : 1 `execution_workflow` durable principal por `execution_id`**.

### Identidad y correlación

| Dominio / `execution_record` | Temporal / runtime | Regla v1 |
|---|---|---|
| `execution_id` | `workflow_id` de `execution_workflow` | Igual valor. Es la clave de unicidad durable. |
| `tenant_id` | namespace lógico / search attribute | Debe quedar visible para aislamiento y búsqueda operativa. |
| `contract_id` | workflow memo + search attribute | Referencia durable al contrato compilado consumido. |
| `contract_fingerprint` | memo + search attribute | Permite diagnosticar drift e idempotencia. |
| `trace_id` | header / memo / search attribute | Correlación transversal con observabilidad. |
| `result_id` | campo mutable del estado workflow + proyección al record | Puede nacer nulo y asignarse después. |
| `approval_request_id` | campo mutable del estado workflow + search attribute cuando exista | Necesario para pausas de approval. |
| `idempotency_key` | parte del estado workflow y del record | No reemplaza `execution_id`; explica deduplicación de intención. |
| `status` canónico | workflow state proyectado | Temporal no es el contrato de estado público; lo materializa. |
| evidence / reason codes | history + payloads resumidos + proyección a PostgreSQL | History complementa evidencia; no sustituye el record canónico. |

### Regla de verdad operativa

- **Temporal** conserva la continuidad durable del proceso vivo, timers, signals, retries y recovery.
- **PostgreSQL / `execution_record`** conserva la representación canónica de estado operativo consultable por el kernel.
- El **event history** de Temporal es evidencia de soporte, no reemplazo del objeto canónico.
- El **workflow** es la autoridad de transición; el **projector** es la autoridad de escritura canónica derivada.

### Campos mínimos que el workflow debe proyectar a `execution_record`

- `execution_id`
- `tenant_id`
- `contract_id`
- `contract_fingerprint`
- `trace_id`
- `status`
- `execution_phase_status`
- `application_phase_status`
- `approval_status`
- `failure_status`
- `compensation_status`
- `result_id` opcional
- `approval_request_id` opcional
- `started_at`
- `updated_at`
- `terminal_at` opcional
- `reason_codes[]`
- `unknown_outcome` bool
- `blocked_reason_code` opcional
- `retry_counters`

## Workflow types mínimos a materializar primero

### 1. `execution_workflow`

Workflow durable principal y obligatorio por `execution_id`. Coordina lifecycle completo, espera señales, programa timers, decide transiciones y ordena activities.

### 2. `application_phase_workflow`

En v1 **no se separa como workflow independiente**. La fase de application se modela como subfase explícita dentro de `execution_workflow` por tres motivos:

1. la prioridad de C.2 es fijar una sola autoridad durable por `execution_id`;
2. separar demasiado temprano duplicaría correlación, queries y recovery sin valor suficiente en el primer corte;
3. la separación **semántica** execution/application sigue siendo obligatoria aunque la separación **física** de workflow se postergue.

Queda reservado como seam futuro si la application phase exige paralelismo, approvals o compensaciones con lifecycle independiente.

### 3. `compensation_workflow`

En v1 **no se materializa como workflow aparte**. La compensación se modela dentro de `execution_workflow` mediante subestado y activities de compensación, porque:

1. C.2 necesita primero fijar el corredor mínimo de compensación sin crear orquestación secundaria;
2. la correlación de causa/falla/compensación queda más simple dentro del mismo `execution_id`;
3. la futura extracción a `compensation_workflow` queda abierta si aparecen compensaciones largas, humanas o multipaso.

## Activity types mínimas a materializar primero

### 1. `validate_execution_preconditions`

Valida que el contrato existe, está compilado, es ejecutable para el runtime y que no faltan prerequisitos estructurales obvios.

### 2. `evaluate_execution_policy_gate`

Prepara/invoca de forma controlada el gate de policy en el punto exacto del lifecycle. En C.2 puede operar como adapter mínimo o stub gobernado, pero el workflow debe asumir que su salida puede producir `allow`, `blocked`, `awaiting_approval` o `failed_closed`.

### 3. `run_execution_step`

Ejecuta el tramo técnico previo a la aplicación efectiva: validaciones operativas, preparación, lectura, simulación, calls técnicas o pasos reversibles definidos por capability.

### 4. `run_application_step`

Ejecuta el side effect o tramo efectivo de aplicación de negocio cuando la operación no es read-only.

### 5. `record_execution_evidence`

Registra evidencia técnica y reason codes del tramo execution.

### 6. `record_application_evidence`

Registra evidencia del tramo application, incluyendo confirmación, parcialidad o incertidumbre del outcome.

### 7. `run_compensation_step`

Ejecuta compensación física o lógica cuando el runtime la habilita.

## Componentes mínimos del runtime skeleton

1. `execution_runtime_orchestrator`
   - recibe request de arranque;
   - asegura idempotencia de start;
   - crea/reutiliza `execution_record`;
   - inicia o adjunta a `execution_workflow`.

2. `execution_workflow`
   - autoridad durable del lifecycle por `execution_id`;
   - maneja estado, timers, signals, queries y recovery.

3. `execution_activity_executor`
   - invoca `validate_execution_preconditions`, `evaluate_execution_policy_gate`, `run_execution_step`, `record_execution_evidence`.

4. `application_activity_executor`
   - invoca `run_application_step`, `record_application_evidence`;
   - no puede ejecutarse si la operación es read-only.

5. `approval_wait_gateway`
   - encapsula la espera durable, expiración y liberación de approvals;
   - expone el seam de señales `release_execution` y `release_application`.

6. `policy_gate_adapter`
   - encapsula la preparación/invocación controlada al gate de policy;
   - prepara el seam con Cerbos sin mezclar PEP/PDP ad hoc.

7. `compensation_coordinator`
   - decide si la compensación es requerida, posible o manual;
   - ordena `run_compensation_step` o transición a escalación.

8. `execution_state_projector`
   - proyecta el estado workflow a `execution_record`;
   - garantiza consistencia visible para consultas canónicas.

9. `runtime_event_emitter`
   - emite eventos operativos derivados y señalización de observabilidad;
   - si falla, no puede romper la persistencia durable del workflow.

## Señales/queries mínimas del workflow durable

### Signals

#### `release_execution`

Libera una pausa previa a la fase de execution. Debe validar que el workflow está en estado compatible con liberación.

#### `release_application`

Libera una pausa previa a la fase de application. Si llega antes de `execution_completed`, no puede adelantar estado; debe quedar rechazada o registrada como no aplicable.

#### `block_execution`

Fuerza transición a `blocked` con `reason_code` auditable. No equivale a `failed`.

#### `request_compensation`

Solicita abrir corredor de compensación. Puede originarse por fallo, verificación externa o decisión manual.

#### `close_execution`

Solicita cierre explícito cuando existe evidencia terminal suficiente. No puede saltarse invariantes de evidencia/correlación.

### Queries

#### `get_execution_state`

Devuelve snapshot del estado actual, subfases y timestamps relevantes.

#### `get_release_status`

Devuelve si execution/application están liberadas, pendientes o expiradas.

#### `get_failure_status`

Devuelve `failed`, `blocked`, `unknown_outcome`, reason codes y retryability.

#### `get_compensation_status`

Devuelve si la compensación es `not_required`, `pending`, `running`, `completed`, `partially_compensated`, `manual_required` o `failed`.

#### `get_linked_ids`

Devuelve `execution_id`, `trace_id`, `contract_id`, `tenant_id`, `result_id`, `approval_request_id` y otros links correlacionables presentes.

## Lifecycle de creación de una ejecución

1. Un caller entrega `execution_start_request` con `execution_id`, `tenant_id`, `contract_id`, `trace_id` y correlación mínima.
2. `execution_runtime_orchestrator` verifica idempotencia de arranque por `execution_id`.
3. Si el `execution_id` ya existe y está vivo o terminal, el arranque NO crea una segunda ejecución; retorna referencia a la existente o conflicto según el caso.
4. Se crea o reutiliza el `execution_record` en estado inicial durable.
5. Se arranca `execution_workflow` con `workflow_id = execution_id`.
6. El workflow proyecta `created` y ejecuta `validate_execution_preconditions`.
7. Si el contrato no es válido/ejecutable, la ejecución cierra como `failed` o `blocked` según el motivo, con evidencia explícita.
8. Si existe approval requerida previa a ejecución, el workflow entra en `awaiting_approval`.
9. Si corresponde gate de policy, invoca `evaluate_execution_policy_gate`.
10. Si policy permite seguir, el workflow entra en fase `executing`.
11. Se ejecuta `run_execution_step` y se registra evidencia técnica.
12. Si la operación es read-only y la evidencia terminal es suficiente, puede cerrar sin fase de application.
13. Si la operación es mutation y execution termina bien, se marca `execution_completed` y recién entonces puede habilitarse la fase `applying`.
14. Se ejecuta `run_application_step` y luego `record_application_evidence`.
15. Según outcome, la ejecución cierra como `application_completed`, `failed`, `blocked`, `compensated`, `partially_compensated`, `manual_closure_pending` o `unknown_outcome` contenido.

## Estados materializados en v1 y transición entre estados

### Estados materializados mínimos

- `created`
- `preconditions_validating`
- `awaiting_approval`
- `policy_evaluating`
- `blocked`
- `executing`
- `execution_completed`
- `applying`
- `application_completed`
- `failed`
- `unknown_outcome`
- `compensation_pending`
- `compensating`
- `compensated`
- `partially_compensated`
- `manual_closure_pending`
- `closed`

### Reglas de transición v1

- `created -> preconditions_validating`
- `preconditions_validating -> awaiting_approval` si falta approval válida
- `preconditions_validating -> policy_evaluating` si puede continuar al gate
- `preconditions_validating -> failed` si el contrato o prerequisitos son inviables
- `awaiting_approval -> policy_evaluating` con `release_execution`
- `awaiting_approval -> blocked` por expiración o rechazo gobernado
- `policy_evaluating -> blocked` si policy bloquea
- `policy_evaluating -> awaiting_approval` si policy exige governance adicional
- `policy_evaluating -> executing` si policy permite seguir
- `executing -> failed` por error no retryable antes de completar execution
- `executing -> unknown_outcome` si se pierde certeza material del outcome técnico relevante
- `executing -> execution_completed` si el tramo execution cierra con evidencia suficiente
- `execution_completed -> closed` sólo para operaciones read-only
- `execution_completed -> applying` sólo para operaciones mutation y nunca antes
- `applying -> application_completed` si el efecto de negocio queda confirmado
- `applying -> unknown_outcome` si hay timeout/corte con resultado externo incierto
- `applying -> compensation_pending` si se confirma necesidad de compensar
- `failed -> compensation_pending` si la estrategia exige compensación posterior
- `unknown_outcome -> compensation_pending` si la verificación externa muestra parcialidad o contención requerida
- `compensation_pending -> compensating`
- `compensating -> compensated`
- `compensating -> partially_compensated`
- `compensating -> manual_closure_pending` si no alcanza para cerrar seguro
- `blocked -> manual_closure_pending` si la resolución exige cierre manual con evidencia
- `application_completed -> closed`
- `compensated -> closed`
- `partially_compensated -> manual_closure_pending`
- `manual_closure_pending -> closed` sólo con evidencia terminal explícita

## Separación execution vs application a nivel runtime

La separación no es cosmética; define safety del runtime.

- **Execution** = tramo técnico/orquestador que valida, prepara, consulta policy, espera governance y corre pasos previos/controlados.
- **Application** = tramo donde se materializa el efecto externo o mutación de negocio.

Reglas obligatorias:

- `execution_completed` y `application_completed` **NO pueden colapsar** en mutation.
- Una operación **read-only** puede cerrar legítimamente desde `execution_completed` hacia `closed` sin entrar a application.
- Una operación **mutation** que intente cerrar en `execution_completed` debe considerarse inválida.
- `release_application` no habilita nada si `execution_completed` todavía no existe.
- Retry técnico en execution no autoriza re-aplicar automáticamente el tramo application.

## Timers, deadlines y pausas mínimas

El runtime skeleton debe definir, como mínimo, cuatro familias de tiempo:

1. **start timeout**
   - deadline para que la ejecución pase de `created` a corredor válido;
   - si expira antes de precondiciones mínimas, cierra con `failed` o `blocked` según causa.

2. **execution step timeout**
   - controla `run_execution_step`;
   - puede habilitar retry técnico si no existe riesgo de doble efecto.

3. **application step timeout**
   - controla `run_application_step`;
   - si rompe certeza material, el default es `unknown_outcome`, no retry ciego.

4. **approval wait deadline**
   - controla `awaiting_approval`;
   - si expira, pasa a `blocked` o `manual_closure_pending` según policy/gobernanza.

Pausas mínimas durables:

- pausa por approval antes de execution;
- pausa por approval antes de application si el contrato/policy lo exige;
- pausa por bloqueo manual;
- pausa por verificación externa ante `unknown_outcome`.

## Política de retries del runtime skeleton

La política v1 debe respetar A.5 y separar con precisión tres conceptos:

### 1. Retry técnico de activity

- Permitido para `validate_execution_preconditions`, `evaluate_execution_policy_gate`, `run_execution_step`, `record_execution_evidence`, `record_application_evidence` cuando la falla sea transitoria.
- Permitido para `run_application_step` **solo** si el paso es idempotente o el sistema puede verificar que el efecto no fue aplicado.

### 2. Retry de ejecución

- Se limita al mismo `execution_id` y al mismo workflow durable.
- Nunca crea una segunda ejecución para el mismo arranque.
- Debe quedar trazado como retry, no como nueva corrida.

### 3. Retry prohibido

- Prohibido cuando existe `idempotency_conflict`.
- Prohibido cuando policy bloquea.
- Prohibido cuando approval faltante no fue resuelta.
- Prohibido cuando existe `unknown_outcome` sin verificación externa suficiente.
- Prohibido para application irreversible con riesgo de doble aplicación.

## Failure model del runtime skeleton

El runtime skeleton debe distinguir al menos estas familias:

1. **precondition_failure**
   - contrato inexistente, no compilado, no ejecutable o contexto estructural inválido.

2. **policy_blocked**
   - policy deniega o exige contención. Debe mapear a `blocked`, no a `failed` genérico.

3. **approval_missing_or_expired**
   - approval requerida ausente, inválida o vencida. Debe mapear a `awaiting_approval` o `blocked`.

4. **technical_retryable_failure**
   - timeout transitorio, rate limit, fallo de red o disponibilidad donde repetir controladamente es seguro.

5. **technical_non_retryable_failure**
   - schema invalid, capability incompatible, input imposible, error estructural.

6. **application_failure**
   - la aplicación efectiva devolvió fallo confirmado.

7. **unknown_outcome_failure**
   - la plataforma pierde certeza material sobre el efecto externo. No se simplifica a `failed`.

8. **compensation_failure**
   - la compensación intentada no logró cerrar el remanente.

## Compensación y escalación manual

La compensación v1 debe existir como corredor explícito aunque la implementación inicial sea corta.

- `compensation_required` se activa cuando hubo aplicación parcial, unknown outcome contenido o fallo posterior con remanente operativo.
- `run_compensation_step` puede ser físico o lógico.
- Si la operación es irreversible, el runtime no inventa rollback; debe pasar a `manual_closure_pending` o `blocked` con evidencia explícita.
- La escalación manual es obligatoria cuando:
  - el efecto es irreversible;
  - el outcome sigue desconocido tras verificación razonable;
  - la compensación es parcial;
  - faltan precondiciones para compensar;
  - el bloqueo requiere decisión humana de cierre.

## Persistencia y correlación operativa

La persistencia mínima v1 se reparte así:

### Temporal

- continuidad durable del workflow;
- history de signals/timers/retries;
- estado vivo del `execution_workflow`.

### PostgreSQL / memoria operativa canónica

- `execution_record`;
- proyección de estado canónico;
- evidencia terminal resumida/referenciada;
- correlación con `contract_id`, `tenant_id`, `trace_id`, `result_id`, `approval_request_id`.

### Correlación mínima obligatoria

Toda creación, transición material o cierre debe enlazar como mínimo:

- `execution_id`
- `trace_id`
- `contract_id`
- `tenant_id`
- `result_id` cuando exista
- `approval_request_id` cuando exista

## Integración con C.1 (`compiled_contract`)

El runtime consume `compiled_contract` ya persistido y NO recompila nada.

Reglas de integración:

- C.2 sólo puede arrancar desde `compiled_contract` persistido.
- Si C.1 devolvió `compiled` pero no `executable`, C.2 no debe reinterpretar: debe bloquear o fallar de forma explícita según la razón del contrato.
- `contract_id`, `contract_version`, `fingerprint`, `policy_evaluation_input` y snapshots relevantes deben llegar al runtime como referencias consumibles, no como texto a reinterpretar.
- El `execution_record` debe referenciar exactamente el `compiled_contract` ejecutado.
- El workflow debe preservar la relación estable entre `execution_id` y `contract_id` durante toda la vida de la ejecución.

## Integración preparatoria con Cerbos

Aunque C.3 profundiza Cerbos, C.2 ya debe fijar el seam correcto.

- `policy_gate_adapter` es el único punto del runtime skeleton que prepara/invoca policy.
- El workflow no codifica reglas de autorización ad hoc.
- La salida mínima del gate debe mapear a: `allow`, `blocked`, `awaiting_approval`, `failed_closed`.
- Si el gate no responde o responde inválidamente en mutation/sensible, el comportamiento por defecto es **fail closed** con `blocked` o contención equivalente.
- El runtime guarda correlación suficiente para que C.3 agregue audit trail sin rediseñar el lifecycle.

## Integración preparatoria con event log / observabilidad

- `runtime_event_emitter` debe emitir eventos operativos derivados de los hitos mínimos: created, awaiting_approval, policy_evaluated, executing, execution_completed, applying, application_completed, blocked, failed, unknown_outcome, compensation_pending, compensated, closed.
- `execution_state_projector` escribe verdad operativa; `runtime_event_emitter` publica derivados.
- Si el emitter falla, el workflow no pierde estado durable ni revierte la transición principal.
- `trace_id` debe propagarse para spans/logs, pero observabilidad no reemplaza el `execution_record`.

## Idempotencia del runtime skeleton

La idempotencia v1 se define así:

- mismo `execution_id` => mismo `execution_workflow`;
- un start duplicado con mismo `execution_id` no crea nueva instancia ni duplica estados;
- si el arranque repetido refiere otro contrato/fingerprint incompatible para el mismo `execution_id`, debe devolver `idempotency_conflict`;
- retries del workflow ocurren dentro del mismo `execution_id`;
- la proyección de estado debe ser monotónica o conflict-safe, nunca duplicadora;
- evidencia y eventos derivados deben llevar claves de deduplicación correlacionables.

## Tests borde mínimos

1. **start duplicado con mismo `execution_id`**
   - esperado: no se crea segundo workflow ni segundo `execution_record`.

2. **`compiled_contract` no compilado**
   - esperado: falla/precondición inválida sin arranque efectivo.

3. **approval requerida no presente**
   - esperado: transición a `awaiting_approval`.

4. **execution liberada pero policy gate bloquea**
   - esperado: transición a `blocked`, no `failed`.

5. **read-only intenta entrar a application phase**
   - esperado: rechazo de transición; cierre correcto desde `execution_completed`.

6. **mutation intenta cerrar en `execution_completed`**
   - esperado: inválido; debe pasar por `applying` o corredor de fallo/compensación.

7. **timeout durante ejecución técnica**
   - esperado: retry técnico o `failed` según retryability, sin duplicar estado.

8. **timeout con outcome externo desconocido**
   - esperado: `unknown_outcome`, no retry ciego.

9. **retry técnico exitoso sin duplicar estado**
   - esperado: mismo `execution_id`, misma progresión canónica, counters actualizados.

10. **retry completo prohibido por idempotency conflict**
    - esperado: rechazo explícito.

11. **signal `release_application` recibido antes de `execution_completed`**
    - esperado: no adelanta fase; queda rechazado o registrado como no aplicable.

12. **signal de bloqueo después de `application_completed`**
    - esperado: no reabre ni degrada ejecución terminal.

13. **compensación en operación irreversible**
    - esperado: `manual_closure_pending` o escalación, no rollback ficticio.

14. **compensación parcial con cierre pendiente**
    - esperado: `partially_compensated -> manual_closure_pending`.

15. **cierre sin evidencia terminal**
    - esperado: `close_execution` rechazada.

16. **query de estado en ejecución recién creada**
    - esperado: `created` + linked ids mínimos presentes.

17. **result linked ids múltiples**
    - esperado: `get_linked_ids` devuelve `result_id` y demás correlaciones de forma estable.

18. **approval expira mientras está `awaiting_approval`**
    - esperado: transición a `blocked` o cierre manual según policy.

19. **workflow recovery después de crash**
    - esperado: reanuda desde history durable sin crear nueva ejecución.

20. **event emission falla pero workflow no debe perder estado durable**
    - esperado: estado canónico persiste; emitter puede reintentarse aparte.

21. **policy gate devuelve respuesta inválida**
    - esperado: fail closed para mutation/sensible.

22. **`block_execution` durante `awaiting_approval`**
    - esperado: transición inmediata a `blocked` con reason code.

23. **`request_compensation` sobre ejecución sin aplicación previa**
    - esperado: rechazo o `not_required`, nunca compensación ficticia.

24. **`close_execution` sobre `unknown_outcome` sin verificación**
    - esperado: rechazo por evidencia insuficiente.

25. **application confirmada luego de `unknown_outcome`**
    - esperado: transición controlada a `application_completed` sin re-aplicar.

26. **application confirma parcialidad luego de `unknown_outcome`**
    - esperado: `compensation_pending`.

27. **retry de `run_application_step` sobre operación no idempotente**
    - esperado: prohibido por seguridad.

28. **`get_failure_status` sobre ejecución bloqueada por policy**
    - esperado: muestra `blocked` y reason code de policy, no `failed` genérico.

## Criterios de aceptación de C.2

1. Temporal queda fijado explícitamente como único runtime durable de verdad.
2. El documento deja explícito que `execution_record` no es reemplazado por Temporal.
3. Existe un workflow durable principal por `execution_id`.
4. Queda definida la separación workflow durable vs activities técnicas.
5. Quedan definidos workflow types, activity types, signals y queries mínimas.
6. El lifecycle de creación y progreso de una ejecución queda descrito paso a paso.
7. `execution_completed` y `application_completed` quedan separados de forma obligatoria.
8. El corredor read-only puede cerrar sin application y el mutation no puede hacerlo.
9. `blocked`, `failed`, `awaiting_approval` y `unknown_outcome` quedan modelados explícitamente.
10. Timers, deadlines, retries, failure model y compensación quedan fijados al nivel suficiente para implementar v1.
11. La correlación mínima con `trace_id`, `contract_id`, `tenant_id`, `result_id` y `approval_request_id` queda documentada.
12. La integración con C.1, Cerbos y event log/observabilidad queda delimitada sin contradicción estructural.
13. La idempotencia de arranque con mismo `execution_id` queda definida.
14. Existe set mínimo de tests borde que cubre arranque, pausas, blocking, retries, unknown outcome, compensación y recovery.
