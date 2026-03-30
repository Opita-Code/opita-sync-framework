# A.1 — Object ownership and versioning v1

## Principios de ownership

- Todo objeto canónico debe tener owner de source-of-truth explícito y owner operativo explícito.
- Ownership no equivale a “quién lo usa”, sino a quién puede definirlo, promoverlo, mutarlo o extinguirlo válidamente.
- En Opyta Sync el ownership está condicionado por plataforma, tenant, runtime y governance híbrida; no se resuelve ad hoc por request.
- Los objetos declarativos se versionan por promoción y publicación; los objetos runtime se versionan por append-only, snapshots o supersession.
- Si un objeto controla seguridad, approvals, catálogo o alcance multi-tenant, la plataforma conserva el control superior aunque exista customización tenant-scoped.

---

## Tipos de ownership

### `platform-owned`

Objeto cuyo source-of-truth lo controla la plataforma. El tenant puede consumirlo y, en algunos casos, parametrizarlo dentro de límites definidos, pero no redefinir el canon base.

### `tenant-owned`

Objeto cuyo source-of-truth efectivo lo administra un tenant dentro de límites impuestos por platform policy.

### `runtime-owned`

Objeto generado y persistido por el sistema durante operación. No se crea manualmente como truth declarativa primaria.

### `hybrid-governed`

Objeto con ownership compartido: una capa base o plantilla la gobierna la plataforma y una capa contextual o binding efectivo la gobierna el tenant o runtime dentro de límites predefinidos.

---

## Tabla canónica por objeto

| object_kind | source_of_truth_owner | runtime_owner | mutable_by | scope | initial_schema_version | initial_object_version |
|---|---|---|---|---|---|---|
| `tenant` | platform-owned | platform control plane | platform operators only | global/tenant | `v1` | `1.0` |
| `subject` | tenant-owned | identity/runtime services | tenant admins; platform under exceptional governance | tenant | `v1` | `1.0` |
| `delegation_grant` | tenant-owned | governance/runtime services | tenant admins within policy limits | tenant | `v1` | `1.0` |
| `policy_artifact` | hybrid-governed | policy engine | platform operators for global/base templates; tenant admins for tenant-scoped overlays if allowed | global/tenant | `v1` | `1.0` |
| `connector` | platform-owned | connector control plane and runtime | platform operators only; tenant can enable/use within allowed bounds | global | `v1` | `1.0` |
| `capability` | platform-owned | catalog/runtime services | platform operators only | global | `v1` | `1.0` |
| `result_type` | platform-owned | runtime resolution services | platform operators only | global | `v1` | `1.0` |
| `intent_contract` | runtime-owned | contract compiler/runtime | system-generated; user may edit only allowed user fields through controlled flows | runtime | `v1` | `1.0` |
| `execution_record` | runtime-owned | execution engine | system-generated append-only updates only | runtime | `v1` | `1.0` |
| `result_record` | runtime-owned | execution/result services | system-generated append-only updates; redaction/supersession via governed records | runtime | `v1` | `1.0` |
| `approval_profile` | tenant-owned | governance services | tenant admins within platform policy limits | tenant | `v1` | `1.0` |
| `approval_request` | runtime-owned | approval runtime | system-generated append-only lifecycle transitions only | runtime | `v1` | `1.0` |
| `approval_decision` | runtime-owned | approval runtime | system-generated from policy/human decision flows; never freeform admin overwrite | runtime | `v1` | `1.0` |
| `memory_record` | runtime-owned | memory services | system-generated; governed retention/redaction flows may add superseding records | tenant/runtime | `v1` | `1.0` |
| `telemetry_event` | runtime-owned | observability pipeline | system-generated only | runtime | `v1` | `1.0` |

---

## Interpretación normativa de ownership por objeto

### Platform-governed con límites tenant

Los siguientes objetos son gobernados por plataforma aunque su uso efectivo quede limitado por tenant:

- `tenant`
- `connector`
- `capability`
- `result_type`
- plantillas base de `policy_artifact`

**Regla**: el tenant no puede cambiar semántica canónica, schema, identidad global ni invariantes duras de estos objetos. Solo puede consumirlos, habilitarlos o boundearlos donde una policy lo permita.

### Tenant-owned dentro de límites de platform policy

- `approval_profile` es tenant-owned dentro de límites de platform policy.
- `delegation_grant` es tenant-owned y administrado por admins dentro de límites de policy.
- `subject` queda tenant-owned a nivel de alta/baja/rol operativo, salvo restricciones de identidad central.

**Regla**: el tenant administra el objeto, pero no puede relajar restricciones globales, floors de approval, constraints de clasificación ni límites de separación de duties impuestos por plataforma.

### Runtime-owned / system-generated

Los siguientes objetos son runtime-owned y system-generated:

- `intent_contract`
- `execution_record`
- `result_record`
- `approval_request`
- `approval_decision`
- `memory_record`
- `telemetry_event`

**Regla**: el runtime puede enriquecer, cerrar, superseder o appendear estados, pero no permite overwrite destructivo ni edición manual libre del truth histórico.

---

## Reglas por familia de objetos

### 1. Catálogo y control plane

Aplica a `tenant`, `connector`, `capability`, `result_type`, base `policy_artifact`.

- Son objetos declarativos y promovibles entre ambientes.
- Su mutación requiere pipeline de revisión/publicación, no edición directa en runtime.
- La versión publicada en `prod` debe ser inmutable salvo supersession por nueva versión material.

### 2. Governance tenant-scoped

Aplica a `approval_profile`, `delegation_grant`, overlays tenant de `policy_artifact`, aspectos operativos de `subject`.

- El tenant es owner del truth efectivo dentro de su scope.
- La plataforma puede imponer validaciones duras y rechazar configuraciones fuera de policy.
- Cambios con impacto en SoD, autoridad o clasificación requieren nueva versión material.

### 3. Runtime append-only

Aplica a `intent_contract`, `execution_record`, `result_record`, `approval_request`, `approval_decision`, `memory_record`, `telemetry_event`.

- Se crean por el sistema, no por authoring manual.
- Para objetos append-only, la “mutación” se modela como nuevos eventos, nuevos snapshots, nuevos estados válidos o nuevos records superseding.
- Nunca se permite borrar o reescribir destructivamente el historial persistido.

---

## Reglas de versionado inicial

### Convención base

- **Initial schema version**: `v1`
- **Initial object version**: `1.0`

### Distinción obligatoria

- `schema_version` identifica la versión del shape y de las reglas estructurales del objeto.
- `object_version` identifica la evolución material de una instancia u objeto declarativo publicado.

### Semántica de `schema_version`

- Cambia solo cuando cambia el contrato estructural o la semántica validable del tipo de objeto.
- `v1` cubre la primera versión formal aprobada en Fase A.

### Semántica de `object_version`

- Arranca en `1.0` para toda instancia u objeto publicado por primera vez.
- Puede evolucionar como `major.minor`.
- `major` cambia cuando se rompe compatibilidad material o cambia semántica central.
- `minor` cambia cuando hay ampliaciones compatibles o cambios no disruptivos pero materialmente relevantes.

---

## Reglas de promoción entre ambientes (`dev`, `staging`, `prod`)

### Objetos declarativos promovibles

Promocionan entre ambientes:

- `tenant` cuando aplique bootstrap/control plane
- `policy_artifact`
- `connector`
- `capability`
- `result_type`
- `approval_profile`

### Regla de promoción

1. Se authora y valida en `dev`.
2. Se promueve a `staging` sin alterar identidad canónica.
3. Se promueve a `prod` solo si el contenido compilado determinístico coincide con la versión aprobada.
4. La promoción nunca reusa una versión material con contenido distinto.

### Restricciones

- No se permite “editar en prod” un objeto declarativo y mantener la misma `object_version`.
- Si `dev` y `staging` divergen materialmente, deben emitir nuevas versiones, no compartir una misma etiqueta.
- Las referencias runtime deben capturar la versión efectiva promovida al ambiente de ejecución.

### Objetos runtime

No promocionan entre ambientes como truth primaria.

**Razón**: su validez depende del contexto operacional real del ambiente. Solo se pueden replicar como evidencia, replay o dataset de auditoría, nunca como configuración viva.

---

## Reglas de mutación permitida / prohibida

### Mutación permitida

#### Declarativos platform-owned o tenant-owned

Se permite:

- crear nueva instancia
- publicar nueva versión material
- deshabilitar o retirar (`retired`) mediante transición gobernada
- agregar campos compatibles si el schema lo permite

No se permite:

- overwrite no auditado de una versión publicada
- cambiar owner efectivo por edición libre
- cambiar scope canónico de una instancia publicada sin nueva versión material

#### Runtime-owned append-only

Se permite:

- agregar nuevos estados válidos de lifecycle
- agregar snapshots posteriores
- marcar superseded, revoked, expired, redacted o closed
- emitir nuevos records derivados o correctivos

No se permite:

- borrar historial
- editar timestamps, fingerprints o evidence de manera destructiva
- reemplazar payload histórico sin dejar rastro explícito

### Regla especial para `intent_contract`

El usuario puede editar campos del grupo editable por el flujo formal de contratación, pero el objeto persistido resultante debe registrarse como nueva versión material o recompilación, nunca como overwrite opaco del contrato histórico.

---

## Qué cambios exigen nueva versión material

### Siempre exigen nueva versión material

- cambio de semántica normativa
- cambio de scope efectivo
- cambio de constraints de seguridad, clasificación o approval
- cambio de bindings a `connector`, `capability`, `result_type` o `policy_artifact`
- cambio de reglas de SoD o autoridad
- cambio de campos materiales usados en fingerprint o compilación
- cambio de behavior de redacción, retención o evidencia mínima

### En runtime append-only

No siempre se crea “nueva versión” por overwrite; se modela como:

- nuevo snapshot
- nuevo evento
- nuevo estado
- nuevo record superseding

**Regla**: si cambia materialmente la base de decisión o de ejecución, debe quedar una nueva materialización auditable aunque el identificador raíz se mantenga correlacionado.

### Cambios que normalmente NO exigen nueva versión material

- correcciones tipográficas en metadata no normativa
- notas explicativas no materiales
- campos derivados no persistidos como canon de decisión

Si un cambio aparentemente menor altera compilación, approval floor, clasificación o auditoría, entonces SÍ es material.

---

## Criterios de aceptación de A.1 ownership/versioning

Se considera cerrada esta parte de A.1 si y solo si se cumplen todos estos puntos:

1. Cada objeto canónico tiene owner de source-of-truth explícito.
2. Cada objeto canónico tiene runtime owner explícito.
3. Cada objeto canónico declara quién puede mutarlo válidamente.
4. Cada objeto canónico tiene scope normativo definido.
5. Todos los objetos arrancan con `schema_version = v1` y `object_version = 1.0`.
6. Queda explícito qué objetos promocionan entre `dev`, `staging` y `prod`.
7. Queda explícito que los objetos runtime append-only no admiten overwrite destructivo.
8. Queda explícito qué cambios fuerzan nueva versión material.
9. Queda consistente con A.2, A.3 y A.4 respecto a contratos, result types y approvals.
10. No quedan zonas grises entre ownership de plataforma, tenant y runtime.

---

## Decisión final de A.1 para ownership y versioning

- El modelo v1 de Opyta Sync usa cuatro tipos de ownership: `platform-owned`, `tenant-owned`, `runtime-owned`, `hybrid-governed`.
- La plataforma gobierna el catálogo base y los invariantes globales.
- El tenant gobierna perfiles y delegaciones dentro de límites explícitos.
- El runtime gobierna records operativos como historial append-only auditable.
- El baseline formal de versionado para todos los objetos canónicos en Fase A es `schema v1` y `object 1.0`.
