# A.2 — Contrato de intención/inspección v1

## Principios base

- El contrato de intención/inspección es la unidad operativa central del sistema. Todo lo que el motor hace parte de él.
- El contrato no es un mensaje libre ni un prompt: es una estructura compilada, validada y gobernable.
- El contrato tiene dos momentos diferenciados: la intención inicial del usuario (parcial, en lenguaje natural) y el contrato compilado (completo, validado, listo para planear o ejecutar).
- La edición del contrato es controlada: el usuario puede editar sus campos, el sistema controla y enriquece los suyos. Ningún campo del sistema puede ser sobreescrito por el usuario directamente.
- El contrato es auditable en todas sus versiones. Cada cambio material genera un nuevo snapshot y un nuevo fingerprint.
- Un contrato no ejecutable no puede avanzar a ejecución. El motor puede planear, inspeccionar y proponer desde un contrato incompleto, pero nunca ejecutar.

---

## Schema completo del contrato

### Grupo A — Campos del usuario (editables)

Campos que el usuario puede proveer, editar o corregir directamente.

- `objetivo` — intención en lenguaje natural; lo que el usuario quiere lograr
- `alcance` — sobre qué entidades, sistemas o períodos aplica
- `tipo_de_resultado_esperado` — uno de los 8 tipos canónicos (ver A.3); puede ser sugerido por el sistema si no se declara
- `autonomia_solicitada` — nivel de autonomía pedido: `manual` | `assisted` | `autonomous`
- `aprobacion_requerida` — si el usuario pide aprobación explícita; el sistema puede elevar pero nunca bajar
- `criterios_de_exito` — condiciones que el usuario considera éxito; requerido para `report`, `execution`, `system_update`
- `restricciones_declaradas` — restricciones explícitas que el usuario declara (sistemas que no tocar, datos que no usar, etc.)
- `notas_del_usuario` — contexto adicional libre que el usuario quiere que el motor considere

### Grupo B — Campos compilados por el sistema (no editables por el usuario)

Campos que el motor resuelve, enriquece o calcula a partir del contexto, permisos, policies y clasificación.

- `sistemas_confirmados` — sistemas reales accesibles y autorizados para esta ejecución (compilado desde `sistemas_posibles` + permisos)
- `datos_permitidos` — datos a los que el actor tiene acceso según permisos, clasificación y delegación
- `datos_restringidos` — datos identificados pero no accesibles para este actor en este contexto
- `herramientas_permitidas` — tools habilitadas para este contrato según capability, permisos y policies
- `herramientas_bloqueadas` — tools identificadas pero no usables en este contrato
- `nivel_de_riesgo` — `risk_level` compilado: `low` | `medium` | `high` | `critical` (ver A.4 risk model)
- `business_risk_score` — score de riesgo de negocio (ver A.4)
- `security_risk_score` — score de riesgo de seguridad (ver A.4)
- `approval_mode_efectivo` — modo de aprobación derivado: `auto` | `pre_execution` | `pre_application` | `double`
- `classification_level` — clasificación compilada del contrato: `public` | `internal` | `confidential` | `restricted`
- `restricciones_del_sistema` — restricciones adicionales impuestas por policies, clasificación o governance
- `contexto_relevante` — contexto recuperado de memoria que el motor considera relevante
- `capability_id` — capability que va a ejecutar este contrato
- `notas_de_contexto` — notas del sistema sobre asunciones, ambigüedades resueltas o decisiones de compilación

### Grupo C — Campos técnicos (generados por el sistema)

Campos de identidad, trazabilidad y ciclo de vida. No editables.

- `contract_id` — identificador único del contrato
- `contract_version` — versión del contrato (incrementa con cada cambio material)
- `parent_contract_id` — si este contrato deriva de otro (re-compilación)
- `fingerprint` — hash determinístico de los campos materiales del contrato (ver sección de material change)
- `contract_state` — estado actual del contrato (ver a2-contract-states-compilation-v1)
- `tenant_id`
- `environment` — `dev` | `staging` | `prod`
- `workspace_id`
- `user_id` — actor que originó la intención
- `acting_for_subject_id` — si el actor actúa en nombre de otro (delegación)
- `delegation_id` — referencia a la delegación activa si aplica
- `session_id`
- `capability_ref` — referencia completa a la capability incluyendo versión
- `created_at`
- `updated_at`
- `compiled_at` — timestamp de la última compilación completa
- `expires_at` — TTL del contrato si aplica

### Grupo D — Snapshots (inmutables una vez capturados)

Campos que capturan el estado del mundo en un momento específico. Críticos para auditoría y fingerprint.

- `policy_snapshot` — versión de las policies activas al momento de compilar
- `classification_snapshot` — clasificación de cada dato y sistema involucrado al compilar
- `risk_snapshot` — scores de riesgo calculados con sus variables al compilar
- `delegation_snapshot` — estado de la delegación activa al compilar (si aplica)
- `permission_snapshot` — permisos efectivos del actor al compilar
- `context_snapshot` — contexto de memoria recuperado al compilar
- `plan_snapshot` — plan generado y aprobado (se popula cuando existe un plan aprobado)
- `destination_snapshot` — sistema destino snapshotteado (para `system_update`)

---

## Separación obligatorio / opcional / derivado

### Campos obligatorios para que el contrato sea válido (no vacío)

Mínimo para que el motor acepte el contrato y pueda trabajar con él:

- `objetivo`
- `tipo_de_resultado_esperado`
- `alcance`
- `autonomia_solicitada`
- `tenant_id`
- `user_id`
- `session_id`
- `contract_id`

### Campos obligatorios para que el contrato sea compilable (listo para planear)

El motor debe poder derivar estos antes de generar un plan:

- Todo lo del grupo anterior, más:
- `sistemas_confirmados` — al menos uno
- `datos_permitidos` — resuelto (aunque sea vacío con razón)
- `herramientas_permitidas` — al menos una
- `nivel_de_riesgo` — calculado
- `approval_mode_efectivo` — derivado
- `classification_level` — compilado
- `capability_id` — asignada
- `policy_snapshot`
- `classification_snapshot`
- `risk_snapshot`

### Campos obligatorios para que el contrato sea ejecutable (listo para ejecutar)

El motor solo puede ejecutar si además tiene:

- Todo lo compilable, más:
- `criterios_de_exito` — definidos (para `execution`, `system_update`, `report`)
- `plan_snapshot` — plan aprobado presente y con fingerprint válido
- `aprobacion_requerida` — resuelta y consistente con `approval_mode_efectivo`
- Approval decision válida referenciada (si `approval_mode_efectivo` ≠ `auto`)
- `restricciones_declaradas` + `restricciones_del_sistema` — ambas resueltas sin contradicción

### Campos opcionales

Presentes solo cuando aplican:

- `notas_del_usuario` — contexto adicional libre
- `notas_de_contexto` — notas del sistema
- `parent_contract_id` — solo en re-compilaciones
- `acting_for_subject_id` — solo con delegación activa
- `delegation_id` — solo con delegación activa
- `delegation_snapshot` — solo con delegación activa
- `destination_snapshot` — solo para `system_update`
- `expires_at` — solo si el contrato tiene TTL configurado

### Campos derivados (calculados, nunca provistos por el usuario)

- `nivel_de_riesgo` — derivado de `business_risk_score` + `security_risk_score` + floors
- `business_risk_score` — calculado con las 10 variables del modelo A.4
- `security_risk_score` — calculado con las 10 variables del modelo A.4
- `approval_mode_efectivo` — derivado de `nivel_de_riesgo` + `tipo_de_resultado_esperado` + `aprobacion_requerida` + overrides de policy
- `classification_level` — derivado del nivel más restrictivo entre datos, sistemas y capability
- `fingerprint` — hash determinístico de los campos materiales
- `sistemas_confirmados` — intersección de `sistemas_posibles` mencionados + permisos efectivos
- `herramientas_permitidas` / `herramientas_bloqueadas` — derivadas de capability + permisos + policies

---

## Validaciones duras

Reglas que el motor aplica antes de avanzar de estado. Una violación bloquea el contrato.

### Validaciones en recepción (antes de aceptar el contrato)

- `objetivo` no puede estar vacío o ser un string de solo espacios.
- `tipo_de_resultado_esperado` debe ser uno de los 8 tipos canónicos definidos en A.3.
- `tenant_id` debe corresponder a un tenant activo.
- `user_id` debe corresponder a un sujeto activo en ese tenant.
- `environment` debe ser `dev`, `staging` o `prod`.

### Validaciones en compilación (antes de marcar como compilable)

- `sistemas_confirmados` debe tener al menos un sistema con permiso efectivo para el actor.
- `herramientas_permitidas` debe tener al menos una tool compatible con `tipo_de_resultado_esperado`.
- `nivel_de_riesgo` no puede quedar sin resolver; si las variables son insuficientes, el motor asume el peor caso (`critical`).
- `approval_mode_efectivo` no puede ser inferior al `min_approval_mode` del `tipo_de_resultado_esperado`.
- `classification_level` no puede ser inferior al nivel más restrictivo de los datos o sistemas involucrados.
- Si `delegation_id` está presente, la delegación debe estar activa, vigente y dentro de scope.
- `policy_snapshot` debe capturar la versión activa al momento de compilar; no puede quedar sin versión.

### Validaciones en ejecución (antes de ejecutar)

- `plan_snapshot` debe estar presente y su fingerprint debe coincidir con el `fingerprint` del contrato al momento de la aprobación.
- Si `approval_mode_efectivo` ≠ `auto`, debe existir una `approval_decision` válida con estado `approved` o `execution_released`.
- `criterios_de_exito` debe estar definido para tipos `execution`, `system_update` y `report`; no puede ser vacío.
- Para `system_update`: `external_effect` debe estar declarado; si es `irreversible`, `rollback_plan` debe estar presente o explícitamente marcado como no disponible con razón.
- `restricciones_declaradas` y `restricciones_del_sistema` no pueden ser contradictorias entre sí.
- El `contract_state` debe ser `executable` antes de ejecutar; ningún otro estado permite ejecución.

### Validaciones de material change (detección de cambio que invalida aprobación)

Los siguientes campos son **materiales**: un cambio en cualquiera de ellos invalida cualquier aprobación previa y obliga a recompilar el fingerprint.

- `objetivo`
- `alcance`
- `tipo_de_resultado_esperado`
- `sistemas_confirmados`
- `datos_permitidos`
- `herramientas_permitidas`
- `plan_snapshot`
- `destination_snapshot`
- `capability_id`
- `classification_level`
- `nivel_de_riesgo`
- `policy_snapshot` (versión)
- `approval_mode_efectivo`

Campos **no materiales** (su cambio no invalida aprobación):
- `notas_del_usuario`
- `notas_de_contexto`
- `contexto_relevante`
- `criterios_de_exito` — EXCEPCIÓN: sí es material si cambia para tipos `execution` o `system_update` después de que el plan fue aprobado

---

## Relación contrato ↔ tipo de resultado

El `tipo_de_resultado_esperado` del contrato determina:

1. Qué campos son obligatorios para compilar (ver A.3 input contract por tipo)
2. Qué `min_approval_mode` aplica como floor
3. Qué campos del Grupo D son obligatorios (`plan_snapshot` para `execution`; `destination_snapshot` para `system_update`)
4. Qué capas de telemetría se activan
5. Qué criterios de éxito son requeridos

El sistema puede sugerir el `tipo_de_resultado_esperado` si el usuario no lo declara, basándose en el `objetivo` y la capability seleccionada. Si la sugerencia es aceptada, el usuario la confirma y queda registrada como campo del Grupo A.

---

## Relación contrato ↔ approvals

- El `approval_mode_efectivo` del contrato es el resultado de compilar: `min_approval_mode` del tipo + riesgo calculado + overrides de policy + `aprobacion_requerida` del usuario. Siempre se toma el más restrictivo.
- El contrato lleva `policy_snapshot`, `classification_snapshot` y `risk_snapshot` que son los mismos que quedan capturados en el `approval_request` de A.4.
- El `fingerprint` del contrato es la base del `material_change_fingerprint` del `approval_request`. Si el contrato cambia después de una aprobación, el fingerprint ya no coincide y la aprobación queda `superseded`.
- El contrato no aprueba ni rechaza — eso es responsabilidad del subsistema de approvals de A.4. El contrato provee los datos necesarios para que ese subsistema opere.

---

## Relación contrato ↔ clasificación

- La clasificación del contrato (`classification_level`) se compila a partir del nivel más restrictivo entre: datos accedidos, sistemas involucrados, capability seleccionada.
- Un contrato `restricted` nunca puede ser ejecutado por un actor sin clearance `restricted`.
- La clasificación del contrato hereda a todos los snapshots derivados de él.
- Si la clasificación sube durante la ejecución (porque se accede a datos más sensibles de lo previsto), el contrato debe recompilarse y la aprobación previa queda `superseded`.

---

## Relación contrato ↔ memoria y contexto

- El `contexto_relevante` del contrato es el resultado de recuperar memoria operativa relevante para el objetivo. No es contenido libre; es memoria estructurada recuperada.
- El `context_snapshot` del Grupo D captura el estado de ese contexto al momento de compilar. Si la memoria cambia de forma material antes de ejecutar, el motor puede recompilar.
- Las asunciones tomadas por el motor durante la compilación quedan en `notas_de_contexto`, no en campos materiales.
