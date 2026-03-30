# A.1 — Executable truth format and naming v1

## Decisión del formato fuente de verdad ejecutable

Se adopta una estrategia de doble formato con responsabilidades separadas:

- **Authoring source of truth**: YAML declarativo legible por humanos.
- **Compiled canonical runtime format**: JSON canónico determinístico.

Esta decisión aplica a todos los objetos declarativos del core en Fase A, especialmente catálogo, policies, tenancy bootstrap y configuraciones gobernadas por tenant.

### Decisión normativa

1. Los humanos escriben y revisan YAML.
2. El sistema valida, normaliza y compila a JSON canónico.
3. El runtime consume JSON canónico, no YAML arbitrario.
4. Toda comparación, fingerprint, promoción y auditoría material se hace sobre el JSON compilado determinístico.

---

## Tradeoff: authoring format vs compiled canonical format

### Por qué YAML para authoring

- es legible por humanos
- facilita revisión, diff semántico y mantenimiento
- reduce fricción para authors de policies, catálogo y governance

### Riesgos de YAML

- admite ambigüedades sintácticas y representacionales
- puede introducir alias, anchors, orden no confiable y coerciones implícitas
- no es buen formato final para fingerprinting estable

### Por qué JSON canónico para runtime

- es más estricto y fácil de validar deterministicamente
- soporta hashing y comparación material estables
- reduce ambigüedad en compilación, promoción y replay

### Tradeoff final decidido

YAML se usa para authoring porque optimiza trabajo humano. JSON canónico se usa para runtime porque optimiza consistencia mecánica. Confundir ambos roles sería un error de arquitectura.

---

## Envelope canónico exacto para objetos declarativos

Todo objeto declarativo debe compilar al siguiente envelope exacto:

```json
{
  "api_version": "v1",
  "kind": "approval_profile",
  "metadata": {},
  "spec": {},
  "status": {}
}
```

### Reglas del envelope

- `api_version` — versión del schema publicado del objeto.
- `kind` — kind canónico en snake_case singular.
- `metadata` — identidad, naming, labels de gobierno, ownership y datos de promoción; nunca lógica operativa principal.
- `spec` — definición declarativa material del objeto.
- `status` — estado observado o publicado por sistema/control plane; no reemplaza `spec`.

### Reglas adicionales

- Ningún campo top-level adicional es válido fuera de este envelope salvo decisión posterior explícita.
- El orden de serialización del JSON compilado debe ser determinístico.
- `status` puede estar vacío en authoring, pero el key debe existir en el canon compilado.

---

## Layout de carpetas recomendado

Se recomienda organizar el source declarativo por dominio:

- `tenancy/`
- `governance/`
- `catalog/`
- `runtime/`
- `memory/`
- `observability/`

### Layout de ejemplo

```text
specs/
  canonical/
    tenancy/
      tenant--acme-prod.yaml
      subject--finance-admin.yaml
      delegation_grant--billing-ops.yaml
    governance/
      policy_artifact--default-approval-floor.yaml
      approval_profile--prod-change-double-approval.yaml
    catalog/
      connector--salesforce-core.yaml
      capability--sync-billing-workflow.yaml
      result_type--execution.yaml
    runtime/
      intent_contract--ctr-2026-0001.yaml
      approval_request--apr-2026-0001.yaml
    memory/
      memory_record--customer-onboarding-playbook.yaml
    observability/
      telemetry_event--contract-created-example.yaml
```

### Regla de layout

- El directorio expresa dominio, no ambiente.
- El ambiente vive en `metadata.environment` o en metadata de promoción, no en el `kind`.
- Los artefactos compilados JSON deben poder ubicarse en un árbol paralelo sin perder identidad lógica.

Ejemplo:

```text
compiled/
  catalog/
    capability--sync-billing-workflow.json
  governance/
    approval_profile--prod-change-double-approval.json
```

---

## Naming conventions globales

### Regla general

El naming del core privilegia consistencia mecánica antes que creatividad humana. Si un nombre no puede validarse automáticamente, no sirve como canon.

### Convenciones obligatorias

- `kind` en **snake_case singular**: `intent_contract`, `approval_profile`
- fields en **snake_case**
- enums en **snake_case**
- IDs humanos en **kebab-case**
- event types en **dot.notation**: `contract.created`
- archivos declarativos por convención **`<kind>--<id>.yaml`**
- directorios por dominio: `tenancy/`, `governance/`, `catalog/`, `runtime/`, `memory/`, `observability/`

### Ejemplos válidos

- `kind: approval_profile`
- `metadata.id: prod-change-double-approval`
- `spec.approval_mode: pre_application`
- `event_type: contract.created`
- archivo: `capability--sync-billing-workflow.yaml`

### Ejemplos inválidos

- `kind: ApprovalProfile`
- `metadata.id: Prod_Change_DoubleApproval`
- `spec.approvalMode: preApplication`
- `event_type: contract_created`
- archivo: `approval_profile_prod.yaml`

---

## Convenciones de IDs

### IDs humanos declarativos

- deben estar en `kebab-case`
- deben ser estables, legibles y no depender de orden incidental
- deben evitar espacios, mayúsculas, underscores y caracteres ambiguos

### Regla recomendada

```text
<contexto>-<propósito>-<calificador>
```

Ejemplos:

- `acme-prod`
- `salesforce-core`
- `sync-billing-workflow`
- `prod-change-double-approval`

### IDs runtime

Los objetos runtime pueden usar IDs opacos o semiestructurados, pero deben seguir una convención determinable por `kind` y dominio.

Ejemplos válidos:

- `ctr-2026-0001`
- `exe-2026-0042`
- `res-2026-0042`
- `apr-2026-0010`

### Regla de unicidad

- `metadata.id` debe ser único dentro de su `kind` y scope efectivo.
- No se permite reutilizar un mismo ID humano para distinto contenido material sin nueva versión explícita.

---

## Convenciones de fields, enums, events y archivos

### Fields

- todos los field names van en snake_case
- booleans deben empezar con prefijo semántico claro cuando aplique: `is_`, `has_`, `can_`, `requires_`
- referencias a otros objetos usan sufijo `_id` o `_ref` según corresponda
- timestamps usan sufijo `_at`
- versiones usan sufijo `_version`

### Enums

- todos los enum values van en snake_case
- no se permiten sinónimos equivalentes para el mismo valor
- los catálogos cerrados deben definirse explícitamente por objeto o dominio

Ejemplos válidos:

- `pre_execution`
- `double`
- `restricted`
- `system_update`

### Events

- los event types usan `dot.notation`
- el primer segmento debe mapear al dominio canónico: `contract`, `approval`, `execution`, `result`, `memory`, `telemetry`
- el segundo segmento expresa el hecho en pasado lógico o transición establecida: `created`, `compiled`, `requested`, `approved`, `started`, `completed`

Ejemplos:

- `contract.created`
- `contract.compiled`
- `approval.requested`
- `approval.decision_recorded`
- `execution.completed`
- `result.redacted`

### Archivos

- authoring files: `<kind>--<id>.yaml`
- compiled files: `<kind>--<id>.json`
- no se permiten espacios, mayúsculas ni sufijos libres en el filename base

---

## Reglas de validación y compilación

### Validación previa de authoring

Antes de compilar, todo YAML debe validar:

1. parseo YAML sin aliases prohibidos ni coerciones ambiguas aceptadas por policy
2. presencia exacta del envelope top-level
3. `api_version` soportado
4. `kind` válido dentro del catálogo canónico
5. `metadata.id` con formato válido
6. fields en snake_case
7. enums dentro de catálogo permitido
8. referencias a objetos existentes o resolubles en el scope correcto
9. ausencia de fields desconocidos cuando el schema sea cerrado

### Compilación canónica

La compilación debe:

1. normalizar tipos primitivos
2. expandir defaults explícitos permitidos por schema
3. remover ambigüedades de representación equivalentes
4. ordenar keys determinísticamente
5. serializar JSON con reglas estables de whitespace y escape
6. producir el mismo output para el mismo input semántico

### Regla de determinismo

Dos YAML semánticamente equivalentes deben compilar al mismo JSON canónico. Si no ocurre, la compilación es inválida.

### Regla de promotion safety

Ningún artefacto puede promocionarse entre `dev`, `staging` y `prod` si el JSON compilado no es determinístico, no valida referencias o no preserva identidad y versión.

---

## Tests borde mínimos

Como mínimo, A.1 debe cubrir estos tests borde:

1. **kind inválido** → rechazar objeto cuyo `kind` no exista en el catálogo canónico.
2. **id con formato inválido** → rechazar `metadata.id` con mayúsculas, espacios o underscore.
3. **field en camelCase** → rechazar cualquier field como `approvalMode`.
4. **enum fuera de catálogo** → rechazar valor no permitido como `preApproval`.
5. **YAML válido pero JSON compilado no determinístico** → bloquear publicación/promoción.
6. **referencia a objeto inexistente** → rechazar `connector_ref`, `policy_ref` o `capability_ref` no resolubles.
7. **scope global mal definido en objeto tenant-owned** → rechazar `approval_profile` con scope global no permitido.
8. **mutación prohibida en objeto append-only** → rechazar overwrite destructivo de `approval_request` o `telemetry_event`.
9. **filename no coincide con `kind` + `metadata.id`** → rechazar artefacto inconsistente.
10. **`api_version` ausente o no soportado** → rechazar compilación.
11. **envelope incompleto** → rechazar si falta alguno de `api_version`, `kind`, `metadata`, `spec`, `status`.
12. **fields top-level extra** → rechazar keys fuera del envelope canónico.
13. **JSON compilado con orden de keys variable** → marcar compilador como inválido.
14. **ID duplicado con contenido material distinto y misma versión** → rechazar publicación.
15. **objeto platform-owned mutado por tenant** → rechazar cambio fuera de autoridad.
16. **objeto tenant-owned fuera de límites de platform policy** → rechazar aunque el YAML sea sintácticamente válido.
17. **referencia cross-tenant no autorizada** → rechazar resolución aunque el objeto exista.
18. **enum correcto pero casing incorrecto** → rechazar `Pre_Execution`.
19. **uso de plural en `kind`** → rechazar `approval_profiles`.
20. **status authoring intenta sobreescribir estado runtime prohibido** → rechazar compilación o publicación según objeto.

---

## Criterios de aceptación de A.1 formato/naming

Se considera cerrado A.1 en formato y naming si y solo si:

1. Queda decidido que el authoring es YAML declarativo.
2. Queda decidido que el runtime consume JSON canónico determinístico.
3. El envelope canónico exacto queda fijado como `api_version`, `kind`, `metadata`, `spec`, `status`.
4. `kind`, fields y enums quedan normalizados en snake_case.
5. Los IDs humanos quedan normalizados en kebab-case.
6. Los eventos quedan normalizados en dot.notation.
7. Los archivos declarativos quedan normalizados como `<kind>--<id>.yaml`.
8. Los directorios quedan organizados por dominio canónico.
9. Las reglas de validación previa y compilación determinística quedan explícitas.
10. Existe una suite mínima de tests borde que cubre naming, envelope, referencias, authority y determinismo.

---

## Decisión final de A.1 para formato ejecutable y naming

- YAML es el formato de authoring humano.
- JSON canónico determinístico es el formato runtime y de comparación material.
- El envelope canónico queda fijado y no es opcional.
- Naming e IDs quedan gobernados por reglas mecánicas uniformes en todo el core.
- Toda promoción futura debe apoyarse en validación estricta y compilación determinística; si no, no hay source-of-truth ejecutable confiable.
