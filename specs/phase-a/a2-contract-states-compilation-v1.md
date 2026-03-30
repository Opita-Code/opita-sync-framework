# A.2 — Estados del contrato y proceso de compilación v1

## Estados del contrato

- `draft` — el contrato fue creado con la intención inicial del usuario; puede ser incompleto o ambiguo
- `inspecting` — el motor está inspeccionando contexto, memoria, permisos y sistemas para enriquecer el contrato
- `incomplete` — la inspección reveló que faltan campos obligatorios que el usuario debe proveer; el motor espera input
- `compiling` — el motor está compilando el contrato: calculando riesgo, clasificación, approval mode, tools y snapshots
- `compiled` — el contrato está compilado; todos los campos del sistema están resueltos; puede avanzar a planificación
- `planning` — el motor está generando un plan a partir del contrato compilado
- `plan_proposed` — el plan fue generado y está pendiente de revisión/aprobación por el usuario o governance
- `plan_approved` — el plan fue aprobado; el `plan_snapshot` está fijado; el contrato puede avanzar a ejecución
- `executable` — el contrato tiene todo lo necesario para ejecutar: plan aprobado, approval válida (si aplica), criterios definidos
- `executing` — el motor está ejecutando el contrato activamente
- `executed` — la ejecución completó; el resultado fue producido
- `superseded` — el contrato fue invalidado por un cambio material después de haber sido aprobado; no puede ejecutarse
- `cancelled` — el contrato fue cancelado explícitamente por el usuario o por governance antes de ejecutar
- `closed` — el contrato cerró de forma terminal (exitoso, fallido o cancelado)

## Reglas operativas clave

- `draft` es el único estado inicial válido. Ningún contrato empieza en otro estado.
- El motor puede avanzar de `draft` a `inspecting` automáticamente si la autonomía configurada lo permite.
- `incomplete` es un estado de espera activo: el motor identificó exactamente qué falta y se lo comunica al usuario.
- El motor nunca ejecuta desde `compiled` directamente: debe pasar por `planning` → `plan_approved` → `executable`.
- `superseded` puede ocurrir desde `plan_proposed`, `plan_approved` o `executable` si un campo material cambia.
- Todo contrato `superseded` genera un nuevo contrato en `draft` si el motor reinicia el proceso.
- `executing` es un estado exclusivo: un contrato no puede estar en `executing` si no está en `executable` primero.
- Todo estado terminal (`executed`, `cancelled`, `closed`) es irreversible.

## Transiciones válidas

```
draft -> inspecting
draft -> cancelled

inspecting -> incomplete (faltan campos del usuario)
inspecting -> compiling (contexto suficiente para compilar)
inspecting -> cancelled

incomplete -> inspecting (el usuario proveyó los campos faltantes)
incomplete -> cancelled

compiling -> compiled
compiling -> incomplete (la compilación reveló campos que necesitan input del usuario)

compiled -> planning
compiled -> cancelled

planning -> plan_proposed
planning -> incomplete (el plan no se pudo generar por ambigüedad residual)

plan_proposed -> plan_approved (el usuario o governance aprobó el plan)
plan_proposed -> planning (el usuario rechazó y pidió replanning)
plan_proposed -> superseded (un campo material cambió)
plan_proposed -> cancelled

plan_approved -> executable (todas las validaciones de ejecución pasaron)
plan_approved -> superseded (un campo material cambió después de aprobar)
plan_approved -> cancelled

executable -> executing
executable -> superseded (un campo material cambió antes de ejecutar)
executable -> cancelled

executing -> executed
executing -> cancelled (cancelación de emergencia durante ejecución — solo si reversible)

executed -> closed
superseded -> draft (el motor puede crear un nuevo contrato derivado)
superseded -> closed
cancelled -> closed
```

## Transiciones prohibidas

- `draft -> executable` (debe compilar y planear primero)
- `draft -> executing`
- `compiled -> executable` (debe planear y aprobar el plan)
- `compiled -> executing`
- `plan_approved -> executing` (debe pasar por `executable`)
- `executed -> executing` (no se puede re-ejecutar el mismo contrato; crear uno nuevo)
- `closed -> cualquier estado` (estado terminal irreversible)
- `cancelled -> executing`
- `superseded -> executing`
- `superseded -> executable`

---

## Proceso de compilación: intención → contrato compilable

La compilación es el proceso por el cual el motor transforma la intención inicial del usuario en un contrato estructurado, validado y gobernable. Es el paso más crítico del ciclo.

### Paso 1 — Recepción y normalización

El motor recibe la intención del usuario (texto libre o estructura parcial) y:
1. Crea el contrato en estado `draft` con los campos del Grupo A que el usuario proveyó.
2. Asigna `contract_id`, `tenant_id`, `user_id`, `session_id`, `created_at`.
3. Si el usuario no declaró `tipo_de_resultado_esperado`, el motor infiere uno candidato basado en el `objetivo` y lo registra como sugerencia en `notas_de_contexto`.
4. Si el usuario no declaró `autonomia_solicitada`, aplica el default del tenant.

### Paso 2 — Inspección de contexto

El motor pasa a `inspecting` y ejecuta en paralelo:
- Recuperar memoria operativa relevante para el `objetivo` y `alcance` → popula `contexto_relevante`
- Identificar sistemas candidatos según el `objetivo` → popula candidatos de `sistemas_confirmados`
- Recuperar policies activas del tenant → inicia `policy_snapshot`
- Recuperar delegaciones activas del actor si aplica → inicia `delegation_snapshot`

Si la inspección revela campos obligatorios del Grupo A que el usuario no proveyó y que el motor no puede inferir, el contrato pasa a `incomplete` y el motor lista exactamente qué falta.

### Paso 3 — Resolución de permisos y acceso

El motor resuelve:
- Permisos efectivos del actor sobre los sistemas candidatos → descarta los no autorizados
- Datos accesibles para el actor en ese contexto → popula `datos_permitidos` y `datos_restringidos`
- Tools habilitadas para la capability candidata → popula `herramientas_permitidas` y `herramientas_bloqueadas`
- Delegación activa → valida scope y vigencia; si inválida, elimina y registra en `notas_de_contexto`
- Popula `permission_snapshot`

### Paso 4 — Cálculo de riesgo

El motor calcula con las variables del modelo A.4:
- `business_risk_score` — 10 variables de riesgo de negocio
- `security_risk_score` — 10 variables de riesgo de seguridad
- `nivel_de_riesgo` — combinación con floors duros
- Registra todas las variables y sus valores en `risk_snapshot`

Si las variables son insuficientes para calcular con confianza, el motor asume el nivel más alto posible y lo registra en `notas_de_contexto`.

### Paso 5 — Compilación de clasificación

El motor determina `classification_level`:
- Evalúa la clasificación de cada dato en `datos_permitidos`
- Evalúa la clasificación de cada sistema en `sistemas_confirmados`
- Evalúa la clasificación de la capability
- Toma el nivel más restrictivo
- Registra la clasificación de cada elemento en `classification_snapshot`

### Paso 6 — Derivación del approval mode

El motor deriva `approval_mode_efectivo`:
1. Toma el `min_approval_mode` del `tipo_de_resultado_esperado` (floor de A.3)
2. Aplica los overrides de A.4 según `nivel_de_riesgo`, `classification_level`, `external_effect` y `scope_size`
3. Respeta el piso declarado por el usuario en `aprobacion_requerida` (el usuario puede pedir más, nunca menos)
4. Aplica overrides duros de policy del tenant
5. Resultado: `approval_mode_efectivo`

### Paso 7 — Asignación de capability

El motor confirma la capability:
- Si el usuario especificó una capability, la valida contra permisos y policies
- Si no, el motor selecciona la más adecuada según `objetivo` y `tipo_de_resultado_esperado`
- Popula `capability_id` y `capability_ref` con la versión exacta

### Paso 8 — Generación de fingerprint

El motor calcula el `fingerprint` del contrato:
- Hash determinístico sobre todos los campos materiales (listados en a2-intent-contract-v1)
- El fingerprint se recalcula cada vez que un campo material cambia
- El fingerprint es la base del `material_change_fingerprint` del `approval_request` de A.4

### Paso 9 — Validación de compilación

El motor ejecuta todas las validaciones duras de compilación (listadas en a2-intent-contract-v1).
- Si pasan: contrato pasa a `compiled`
- Si fallan por falta de input del usuario: contrato pasa a `incomplete`
- Si fallan por policy o governance: contrato pasa a `cancelled` con reason code

### Paso 10 — Cierre de compilación

El motor:
- Popula `policy_snapshot` con las versiones exactas de policies evaluadas
- Marca `compiled_at`
- Actualiza `contract_version`
- El contrato pasa a `compiled`

---

## Proceso de planificación: compilado → plan aprobado

### Plan generation

Desde `compiled`, el motor pasa a `planning` y genera el plan:
- Selecciona tools de `herramientas_permitidas` para cada paso
- Genera los pasos respetando las restricciones del contrato
- Calcula las asunciones necesarias y las registra
- El plan no modifica campos materiales del contrato; vive en `plan_steps` del resultado `plan`

### Plan snapshot

Una vez generado:
- El motor crea el `plan_snapshot` con el plan completo y el fingerprint del contrato en ese momento
- El contrato pasa a `plan_proposed`
- El `plan_snapshot` es inmutable: si el plan necesita modificarse, se regenera y el fingerprint se recalcula

### Plan approval

Cuando el plan es aprobado (por el usuario, por policy auto-approval o por governance):
- El contrato pasa a `plan_approved`
- El `plan_snapshot` queda fijado
- El `approval_request` de A.4 captura el `fingerprint` del contrato en este momento como `material_change_fingerprint`

### Transición a executable

El contrato pasa a `executable` cuando:
1. `plan_snapshot` está presente y su fingerprint es consistente con el contrato actual
2. Si `approval_mode_efectivo` ≠ `auto`: existe `approval_decision` válida referenciada
3. `criterios_de_exito` está definido para los tipos que lo requieren
4. No existe ningún campo material que haya cambiado desde la aprobación (fingerprint match)
5. Todas las validaciones de ejecución pasan

---

## Detección de material change y superseded

### Cuándo se detecta material change

El motor evalúa material change:
- Cada vez que el usuario edita un campo del Grupo A
- Cada vez que el sistema recompila un campo del Grupo B
- Antes de cada transición de estado crítica (plan_approved → executable, executable → executing)

### Proceso de evaluación

1. El motor recalcula el `fingerprint` con los valores actuales de los campos materiales
2. Compara con el `fingerprint` registrado en el `approval_request` vigente
3. Si difieren: material change detectado

### Consecuencias del material change

- Si el contrato está en `plan_proposed` o posterior: pasa a `superseded`
- Se emite evento `contract.material_change_detected` con qué campo cambió
- El `approval_request` vigente pasa a `superseded` (hereda la invalidación de A.4)
- El motor puede crear un nuevo contrato en `draft` derivado del `superseded` si el actor quiere continuar
- El nuevo contrato empieza el ciclo desde `draft` con los campos actualizados

### Campos que siempre recalculan fingerprint

Ver lista completa en a2-intent-contract-v1 sección "Validaciones de material change".

---

## Autonomía y nivel de asunción durante compilación

El nivel de autonomía del contrato (`autonomia_solicitada`) afecta cómo el motor maneja la compilación:

### `manual`
- El motor no avanza de estado sin confirmación explícita del usuario en cada paso crítico
- El motor presenta cada asunción para revisión antes de compilar
- El contrato nunca pasa de `inspecting` a `compiling` sin confirmación
- Más lento, más control

### `assisted`
- El motor avanza automáticamente hasta `plan_proposed`
- El usuario confirma el plan antes de que pase a `plan_approved`
- Las asunciones se documentan pero no requieren confirmación individual
- Balance entre velocidad y control

### `autonomous`
- El motor puede avanzar hasta `executable` sin intervención si el `approval_mode_efectivo` es `auto`
- Si el approval mode requiere intervención humana, el motor espera en el estado correspondiente
- La autonomía no puede saltarse floors de approval ni validaciones duras
- El nivel de asunción configurable del tenant aplica aquí: cuánto puede el motor asumir sin preguntar
