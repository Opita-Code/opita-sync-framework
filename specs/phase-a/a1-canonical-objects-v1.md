# A.1 — Canonical objects v1

## Principios base

- El core de Opyta Sync se modela sobre un conjunto finito de objetos canónicos explícitos, no sobre documentos ad hoc ni payloads implícitos.
- Un objeto canónico existe porque necesita identidad estable, reglas formales, lifecycle auditable y relaciones consistentes entre dominios.
- La lista canónica debe ser mínima pero suficiente: si un concepto no necesita gobierno propio, versionado propio o referencias top-level, no entra como first-class.
- El source-of-truth distingue entre objetos declarativos gobernados y records operativos generados por runtime.
- Multi-tenant, delegación explícita, catálogo de capacidades, policies globales del tenant y separación entre memoria operativa y telemetría son restricciones fundacionales de esta fase.

---

## Regla de qué califica como objeto canónico del core

Un concepto **califica como objeto canónico del core** si cumple simultáneamente estas condiciones:

1. Tiene identidad estable y referenciable por otros objetos o eventos.
2. Tiene reglas propias de validez, ownership o lifecycle.
3. Necesita existir como unidad auditable separada, no solo como campo embebido.
4. Su ausencia como top-level object produciría ambigüedad operacional, de governance o de auditoría.
5. Aparece en más de un flujo crítico del sistema (planeamiento, ejecución, approvals, catálogo, memoria u observabilidad).

Si un concepto falla una o más de estas condiciones, debe modelarse como:

- campo embebido,
- snapshot embebido,
- enum,
- subtype,
- payload de evento,
- o detalle interno de otro objeto canónico.

---

## Lista final de objetos críticos del core

La lista final de objetos canónicos de A.1 es la siguiente:

- `tenant`
- `subject`
- `delegation_grant`
- `policy_artifact`
- `connector`
- `capability`
- `result_type`
- `intent_contract`
- `execution_record`
- `result_record`
- `approval_profile`
- `approval_request`
- `approval_decision`
- `memory_record`
- `telemetry_event`

---

## Objetos canónicos por dominio

### 1. Tenancy e identidad

#### `tenant`

**Propósito**

Representa la frontera administrativa, de datos, policy y operación sobre la que se ejecuta Opyta Sync.

**Por qué es first-class**

- Define aislamiento multi-tenant.
- Es el anchor de scope para policies, approvals, delegaciones, memoria y ejecución.
- Debe ser referenciable por catálogo, runtime y auditoría.

**Scope**

`global` como registro platform-governed con límites operativos por tenant; sus instancias delimitan operación `tenant`.

**Relaciones clave**

- un `tenant` contiene `subject`, `delegation_grant`, `approval_profile`, `memory_record`
- un `tenant` consume `connector`, `capability`, `result_type` y `policy_artifact`
- todo `intent_contract`, `execution_record`, `result_record`, `approval_request`, `approval_decision`, `telemetry_event` referencia `tenant_id`

#### `subject`

**Propósito**

Representa al actor identificable que opera directa o delegadamente dentro de un tenant.

**Por qué es first-class**

- Es la unidad mínima de identidad y autoridad efectiva.
- Necesita trazabilidad separada para approvals, delegación y auditoría.
- No puede quedar absorbido por credenciales o claims efímeros.

**Scope**

`tenant`

**Relaciones clave**

- pertenece a un `tenant`
- puede originar `intent_contract`
- puede emitir o recibir `delegation_grant`
- puede solicitar `approval_request` y producir `approval_decision`

#### `delegation_grant`

**Propósito**

Modela la delegación explícita y gobernada por la cual un subject puede actuar en nombre de otro dentro de límites definidos.

**Por qué es first-class**

- La delegación es restricción fundacional del sistema.
- Requiere scope, vigencia, revocación y evidencia propios.
- Debe ser referenciable desde contratos, approvals y ejecución.

**Scope**

`tenant`

**Relaciones clave**

- referencia `tenant`, subject delegante y subject delegado
- condiciona `intent_contract`, `approval_request` y `approval_decision`
- se snapshottea en `intent_contract` y `approval_request`

---

### 2. Governance y catálogo central

#### `policy_artifact`

**Propósito**

Representa una policy ejecutable o compilable que impone reglas de autorización, approval, clasificación, riesgo, publicación o uso.

**Por qué es first-class**

- Las policies son source-of-truth gobernado.
- Necesitan versionado material y promoción entre ambientes.
- Sus snapshots gobiernan contratos, approvals y ejecución.

**Scope**

`global` o `tenant`, según tipo y binding efectivo.

**Relaciones clave**

- gobierna `capability`, `connector`, `approval_profile`, `intent_contract`
- sus versiones quedan en snapshots de `intent_contract`, `approval_request`, `approval_decision`, `execution_record`
- puede imponer límites sobre `delegation_grant` y `memory_record`

#### `connector`

**Propósito**

Representa un conector gobernado hacia un sistema externo o interno bajo control central.

**Por qué es first-class**

- Los conectores son activos sensibles y centrales.
- Tienen lifecycle, posture, clasificación y restricciones propias.
- Son referencia obligatoria para ejecución y approval.

**Scope**

`global` con habilitación efectiva por `tenant`.

**Relaciones clave**

- es consumido por `capability`
- puede ser referenciado en `intent_contract`, `approval_request`, `execution_record`, `telemetry_event`
- su estado y versión pueden afectar `policy_artifact` y `approval_profile`

#### `capability`

**Propósito**

Define una capacidad canónica del catálogo que el sistema puede exponer y gobernar.

**Por qué es first-class**

- Es la unidad de publicación, autorización, routing y contrato operativo.
- Debe mapearse con tipos de resultado, conectores y approvals.
- Necesita identidad y versionado de catálogo.

**Scope**

`global` con disponibilidad efectiva por `tenant`.

**Relaciones clave**

- cubre tanto `action` como `workflow`; ambos NO se modelan como objetos top-level separados en A.1
- referencia `connector`, `result_type` y `policy_artifact`
- es seleccionada por `intent_contract`
- produce `execution_record` y `result_record`

#### `result_type`

**Propósito**

Define el contrato canónico de comportamiento de un tipo de resultado del motor.

**Por qué es first-class**

- Determina input contract, output contract, evidence, floors de approval y telemetría.
- Es catálogo compartido entre compilación, runtime y auditoría.
- Necesita referencia estable por capabilities y resultados.

**Scope**

`global`

**Relaciones clave**

- es referenciado por `intent_contract`, `capability` y `result_record`
- controla semantics de outcomes y evidencia
- incluye `governance_decision` como tipo canónico

---

### 3. Contratación y ejecución runtime

#### `intent_contract`

**Propósito**

Es la unidad operativa compilada que traduce intención en contrato validable, planeable, aprobable y eventualmente ejecutable.

**Por qué es first-class**

- Es el pivote central entre intención, approvals, plan y ejecución.
- Tiene fingerprint, snapshots, versionado material y lifecycle propio.
- Necesita referencias estables para runtime y auditoría.

**Scope**

`runtime`

**Relaciones clave**

- referencia `tenant`, `subject`, `delegation_grant`, `capability`, `result_type`
- embebe snapshots como `risk_snapshot`, `classification_snapshot`, `plan_snapshot`, `destination_snapshot`
- origina `approval_request`, `execution_record` y `result_record`

#### `execution_record`

**Propósito**

Registra la ejecución concreta de una capability bajo un contrato y contexto determinados.

**Por qué es first-class**

- La ejecución tiene lifecycle, idempotencia, estados y evidencia propios.
- Debe persistir aunque no produzca éxito.
- Es necesaria para trazabilidad entre contrato, approvals, resultados y telemetría.

**Scope**

`runtime`

**Relaciones clave**

- referencia `intent_contract`, `capability`, `tenant`
- puede depender de `approval_decision`
- produce `result_record`
- emite `telemetry_event`

#### `result_record`

**Propósito**

Materializa el resultado producido por una ejecución conforme al `result_type` canónico.

**Por qué es first-class**

- Es el artefacto de salida auditable del motor.
- Debe persistir clasificación, evidencia, outcome y redacción.
- Es la unidad consumible por usuario, auditoría y flujos posteriores.

**Scope**

`runtime`

**Relaciones clave**

- referencia `execution_record`, `intent_contract`, `result_type`, `capability`
- puede referenciar `approval_decision`
- se relaciona con `telemetry_event` y `memory_record`
- **`result_outcome` no es objeto first-class**: vive embebido en `result_record`

---

### 4. Governance de approvals

#### `approval_profile`

**Propósito**

Define el comportamiento reusable de approvals para una familia de decisiones o contextos.

**Por qué es first-class**

- Centraliza reglas de modo, autoridad, SoD, expiración e invalidación.
- Debe ser gobernable por tenant dentro de límites de platform policy.
- Es referencia reusable de requests concretos.

**Scope**

`tenant`

**Relaciones clave**

- referencia `policy_artifact`
- gobierna `approval_request`
- condiciona `approval_decision`, `intent_contract`, `execution_record`

#### `approval_request`

**Propósito**

Captura una solicitud concreta de aprobación con snapshots inmutables del contexto relevante.

**Por qué es first-class**

- Tiene lifecycle, expiración, invalidación e idempotencia propios.
- Debe conservar evidencia y fingerprint material.
- Es frontera formal entre análisis/plan y autorización efectiva.

**Scope**

`runtime`

**Relaciones clave**

- referencia `intent_contract`, `approval_profile`, `tenant`, `subject`
- embebe snapshots de riesgo, clasificación, plan, conectores y destino
- recibe una o más `approval_decision`

#### `approval_decision`

**Propósito**

Registra cada decisión política o humana asociada a un approval request.

**Por qué es first-class**

- La decisión es evidencia material separada del request.
- Requiere trazabilidad de quién decidió, bajo qué policy y con qué snapshot.
- Debe soportar supersession, revocación y expiración.

**Scope**

`runtime`

**Relaciones clave**

- pertenece a `approval_request`
- referencia `subject`, `delegation_grant`, `policy_artifact`
- puede habilitar o bloquear `execution_record`

---

### 5. Memoria y observabilidad

#### `memory_record`

**Propósito**

Representa memoria operativa estructurada, recuperable y gobernada para asistir compilación, planificación, ejecución o explicación.

**Por qué es first-class**

- La memoria operativa es source separado de telemetría.
- Necesita reglas de retención, clasificación, recuperación y redacción.
- Debe poder referenciarse desde contratos y resultados.

**Scope**

`tenant` o `runtime`, según origen y binding operativo.

**Relaciones clave**

- puede alimentar `intent_contract` mediante `context_snapshot`
- puede ser producida o actualizada a partir de `result_record`
- está gobernada por `policy_artifact`

#### `telemetry_event`

**Propósito**

Registra observabilidad técnica y auditable del sistema en forma de eventos append-only.

**Por qué es first-class**

- La telemetría requiere ingestión, correlación y retención separadas.
- Debe existir aunque no haya resultado usable.
- Es esencial para auditoría, SRE, debugging y compliance.

**Scope**

`runtime`

**Relaciones clave**

- referencia `tenant`, `intent_contract`, `execution_record`, `result_record`, `approval_request`
- NO reemplaza `memory_record`
- alimenta observabilidad y auditoría, no contexto operativo primario

---

## Tabla consolidada de objetos canónicos

| object_kind | dominio | propósito resumido | first-class porque | scope principal | relaciones clave |
|---|---|---|---|---|---|
| `tenant` | tenancy | frontera administrativa y de datos | aislamiento multi-tenant y anchor de governance | global/tenant | `subject`, `policy_artifact`, `intent_contract` |
| `subject` | identity | actor identificable | identidad y autoridad auditable | tenant | `tenant`, `delegation_grant`, `approval_decision` |
| `delegation_grant` | identity/governance | delegación explícita | scope, vigencia y revocación propios | tenant | `subject`, `intent_contract`, `approval_request` |
| `policy_artifact` | governance | policy ejecutable/compilable | gobierno, versionado y promoción | global/tenant | `capability`, `approval_profile`, snapshots |
| `connector` | catalog | conexión gobernada a sistemas | activo sensible con lifecycle propio | global | `capability`, `execution_record`, `approval_request` |
| `capability` | catalog | capacidad canónica del motor | unidad de publicación, routing y autorización | global | `connector`, `result_type`, `intent_contract` |
| `result_type` | catalog | contrato de comportamiento del resultado | gobierna input/output/evidence/approval floors | global | `capability`, `intent_contract`, `result_record` |
| `intent_contract` | runtime | contrato compilado desde intención | pivote entre intención, approval y ejecución | runtime | `capability`, `approval_request`, `execution_record` |
| `execution_record` | runtime | registro de ejecución concreta | lifecycle e idempotencia propios | runtime | `intent_contract`, `result_record`, `telemetry_event` |
| `result_record` | runtime | resultado materializado | salida auditable y consumible | runtime | `execution_record`, `result_type`, `memory_record` |
| `approval_profile` | governance | perfil reusable de approval | SoD, modos y constraints configurables | tenant | `policy_artifact`, `approval_request` |
| `approval_request` | governance/runtime | solicitud concreta de approval | snapshots, expiración e invalidación propias | runtime | `intent_contract`, `approval_profile`, `approval_decision` |
| `approval_decision` | governance/runtime | decisión política o humana | evidencia material separada | runtime | `approval_request`, `subject`, `execution_record` |
| `memory_record` | memory | memoria operativa estructurada | recuperación y gobierno propios | tenant/runtime | `intent_contract`, `result_record` |
| `telemetry_event` | observability | evento de observabilidad | correlación y retención append-only | runtime | `execution_record`, `result_record`, `approval_request` |

---

## Qué NO entra como objeto first-class y por qué

### `result_outcome`

No entra como objeto top-level.

**Razón**: no necesita identidad ni lifecycle independiente; es parte del contrato semántico del `result_record`. Su valor, reason codes y niveles viven embebidos en el resultado producido.

### Snapshots (`risk_snapshot`, `classification_snapshot`, `policy_snapshot`, `plan_snapshot`, `destination_snapshot`, etc.)

No entran como objetos top-level.

**Razón**: son capturas inmutables contextuales de otro objeto, no entidades gobernadas autónomamente. Viven embebidos en `intent_contract`, `approval_request`, `approval_decision` o `execution_record`.

### `governance_decision`

No entra como objeto top-level separado.

**Razón**: en A.3 se definió como **tipo canónico de resultado**, no como entidad raíz. Por lo tanto vive como valor de `result_type` y se materializa en `result_record` cuando corresponde.

### `action` y `workflow`

No entran como objetos top-level separados en A.1.

**Razón**: quedan absorbidos por `capability`, que cubre ambos modelos bajo una única unidad de catálogo, autorización y versionado.

### Credenciales, secrets y sesiones efímeras

No entran como objetos canónicos del core de A.1.

**Razón**: son concerns operativos o de infraestructura. Pueden existir en subsistemas específicos, pero no cumplen la regla de objeto canónico transversal definida en esta sección.

### Métricas agregadas, vistas derivadas y caches

No entran como objetos canónicos.

**Razón**: son proyecciones derivadas de `telemetry_event`, `execution_record` o `result_record`, y no constituyen source-of-truth primario.

---

## Cierre normativo de A.1 sobre objetos

- La lista de 15 objetos anteriores es la lista final de objetos críticos del core para v1.
- Ningún nuevo objeto top-level debe agregarse en Fase A sin demostrar que cumple la regla de objeto canónico definida arriba.
- Todo concepto adicional de esta fase debe modelarse como campo, snapshot, enum, subtype o vista derivada salvo decisión arquitectónica posterior explícita.
