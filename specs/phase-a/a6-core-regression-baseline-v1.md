# A.6 — Regresión mínima del core v1

## Principios base

- Esta especificación define la **línea mínima obligatoria** de regresión del core para Phase A. No intenta agotar todos los casos futuros; fija el piso duro que NO puede faltar.
- La regresión mínima del core valida comportamiento observable del motor sobre las verdades ya cerradas en A.1-A.5 y sobre la familia de evals que A.6 introduce.
- El objetivo NO es maximizar cobertura exhaustiva, sino detectar con alta señal las regresiones que rompen contrato, approvals, runtime, resultados, clasificación, idempotencia, compensación, onboarding de tenant y trazabilidad.
- Toda corrida de regresión debe dejar evidencia auditable, correlacionable y suficiente para reconstruir qué se evaluó, contra qué fixtures/artifacts y con qué outcome.
- Debe distinguirse entre **corrupción estructural** del sistema y **degradación de calidad no crítica**. La primera habilita `fail-fast`; la segunda debe seguir acumulando evidencia.
- La regresión mínima del core es un estándar de seguridad operativa. La cobertura exhaustiva futura podrá ampliar dominios, variantes, datasets, scale cases y escenarios probabilísticos, pero NO reemplaza este baseline.

---

## Qué significa “regresión mínima” en Opyta Sync

En Opyta Sync, “regresión mínima” significa el conjunto más chico de suites, casos, fixtures, métricas y artifacts que todavía permite afirmar con fundamento que el core sigue respetando sus contratos sistémicos.

Implica validar, como mínimo, que:

1. el contrato sigue compilando y correlacionando correctamente con runtime y resultado,
2. approvals siguen gobernando ejecución y aplicación como A.4 exige,
3. runtime sigue respetando estados, releases, deduplicación y compensación de A.5,
4. resultados y clasificación siguen siendo coherentes y auditables,
5. la plataforma deja trazabilidad, métricas y referencias a traces/ids según source-of-truth,
6. un tenant nuevo puede alcanzar baseline operable sin romper aislamiento ni invariantes.

No significa:

- cobertura completa de todas las capabilities,
- validación exhaustiva de todas las combinaciones de policy,
- exploración masiva de escenarios probabilísticos,
- reemplazo de stress/performance/chaos testing dedicado.

### Diferencia entre cobertura estructural mínima y cobertura exhaustiva futura

- **Cobertura estructural mínima**: prueba una vez por invariante crítica del core, con suficientes variantes para detectar ruptura conceptual.
- **Cobertura exhaustiva futura**: multiplica tipos de fixtures, escalas, combinatorias, tenants, ambientes, policies, conectores y escenarios históricos.
- Regla v1: si una invariante crítica no está cubierta por la regresión mínima, el baseline está incompleto; si solo falta una variante adicional no crítica, eso pertenece al roadmap de cobertura exhaustiva.

---

## Suites mínimas obligatorias del core

Las suites mínimas obligatorias v1 son exactamente estas:

| suite_name | objetivo | dominios cubiertos | nivel mínimo donde corre (`smoke` / `standard` / `release_candidate`) | blocking_mode (`hard_gate` / `soft_gate`) | artifacts requeridos |
|---|---|---|---|---|---|
| `contract_suite` | Detectar corrupción del contrato compilado, fingerprints y compatibilidad estructural del input del core | contract, runtime, observability | `smoke` | `hard_gate` | resumen, assertions estructurales, fingerprints observados/esperados, failures, trace refs |
| `result_suite` | Verificar coherencia entre ejecución y `result_record`, outcome formal y evidencia mínima | result, runtime, observability | `smoke` | `hard_gate` | resumen, outputs observados, reason codes, evidence refs, métricas de completitud |
| `approval_suite` | Validar lifecycle, releases, vigencia, SoD y bloqueo correcto de approvals | approval, runtime, observability | `standard` | `hard_gate` | resumen, decisiones, transiciones, failures, approval ids, event refs |
| `classification_suite` | Confirmar clasificación efectiva, redacción, bloqueo y consistencia del output entregable | classification, result, observability | `standard` | `hard_gate` | resumen, snapshots de clasificación, artifacts redactados, mismatches, trace refs |
| `runtime_suite` | Validar máquina de estados, correlación entre records y cierre auditable de la ejecución | runtime, contract, result, approval, observability | `smoke` | `hard_gate` | resumen, timeline de estados, event ids, execution/result ids, métricas de transición |
| `resilience_suite` | Verificar idempotencia, deduplicación, retries, unknown outcome y compensación | resilience, runtime, result, observability | `standard` | `hard_gate` | resumen, intentos, causal links, duplicate refs, compensation refs, métricas de retry |
| `tenant_onboarding_suite` | Probar que el tenant alcanza mínimo operable y conserva aislamiento/gobernance base | tenant, contract, approval, classification, observability | `standard` | `hard_gate` | resumen, tenant snapshot, validation summary, bootstrap refs, event ids, métricas de onboarding |

### Regla de completitud por nivel

- `smoke` DEBE correr un subconjunto crítico de suites mínimas.
- `standard` DEBE cubrir todas las invariantes estructurales mínimas del core y puede dejar casos extensivos para `release_candidate`.
- `release_candidate` DEBE incluir **TODAS** las suites mínimas obligatorias.

---

## Matriz de cobertura mínima por dominio

| dominio | invariante mínima obligatoria | suites que la cubren | evidencia mínima requerida |
|---|---|---|---|
| contract | contrato compilable, resoluble y con fingerprint consistente con la ejecución | `contract_suite`, `runtime_suite`, `tenant_onboarding_suite` | contract id, contract fingerprint, schema/fingerprint assertion, trace id |
| result | `result_record` consistente con tipo, outcome formal y evidence refs mínimas | `result_suite`, `runtime_suite`, `classification_suite` | result id, outcome, reason codes, evidence refs |
| approval | release y bloqueo correctos según lifecycle, vigencia, SoD e invalidación material | `approval_suite`, `runtime_suite`, `tenant_onboarding_suite` | approval request/decision ids, state timeline, invalidation evidence |
| classification | clasificación efectiva coherente con output, redacción y restricciones de entrega | `classification_suite`, `result_suite` | classification snapshot, output redacted/unredacted refs, policy refs |
| runtime | máquina de estados válida, correlación entre artifacts y cierre auditable | `runtime_suite`, `resilience_suite`, `approval_suite` | execution id, event ids, state transitions, terminal evidence |
| resilience | deduplicación, retry seguro, unknown outcome defensivo y compensación trazable | `resilience_suite` | idempotency key, causal links, retry attempts, compensation refs |
| tenant | tenant no operable hasta completar mínimo duro; operable cuando bootstrap resuelve referencias críticas | `tenant_onboarding_suite` | tenant snapshot, validation summary, operable flag, bootstrap event ids |
| observability | toda corrida deja resumen, fallos, métricas y referencias trazables a traces/ids | **todas** | run summary, metrics bundle, failures bundle, trace/correlation refs |

---

## Casos mínimos obligatorios por familia

Cada suite DEBE incluir al menos estos casos mínimos v1.

### `contract_suite`

1. contrato válido compila y produce fingerprint estable;
2. cambio material de contrato cambia fingerprint y rompe compatibilidad previa;
3. schema incompatible o artifact corrupto dispara falla estructural dura;
4. referencias obligatorias faltantes impiden avanzar a ejecución;
5. contrato compilado correlaciona con `execution_record` correcto.

### `result_suite`

1. resultado exitoso mínimo con evidence refs válidas;
2. `result_record` sin evidencia mínima falla;
3. `result_type` incompatible con contrato falla;
4. reason code y `result_outcome` deben ser consistentes;
5. resultado terminal debe quedar correlacionado con `execution_id` y `trace_id`.

### `approval_suite`

1. ejecución que requiere approval no corre sin release vigente;
2. approval expirada/revocada/superseded bloquea correctamente;
3. `pre_application` permite ejecución técnica pero no aplicación;
4. `double` respeta SoD y autoridad;
5. cambio material invalida release previa.

### `classification_suite`

1. output permitido se entrega sin sobre-redacción;
2. output restringido se redacta o bloquea según policy;
3. clasificación incoherente entre snapshot y output falla;
4. artifact entregado conserva evidencia de redacción aplicada;
5. reason code de bloqueo/clasificación queda trazado.

### `runtime_suite`

1. lifecycle básico válido desde `created` hasta terminal consistente;
2. transición inválida de estado falla;
3. `execution_completed` no equivale automáticamente a `application_completed`;
4. runtime bloqueado por governance no se clasifica como falla técnica;
5. cierre terminal exige trazabilidad completa de artifacts.

### `resilience_suite`

1. intento duplicado equivalente devuelve outcome existente o se adjunta como duplicado;
2. cambio material rompe grupo de deduplicación;
3. timeout con `unknown_outcome` entra en modo defensivo y no reaplica ciegamente;
4. retry de paso técnico idempotente conserva coherencia;
5. compensación deja causalidad y evidencia explícita.

### `tenant_onboarding_suite`

1. tenant incompleto no entra en `operable`;
2. tenant con mínimo duro completo entra en `operable`;
3. referencias críticas no resolubles bloquean onboarding;
4. tenant `single_user` baseline funciona sin relajar governance;
5. aislamiento por `tenant_id` y `environment` se conserva en artifacts emitidos.

---

## Regresión smoke vs standard vs release-candidate

### `smoke`

- Objetivo: detectar rápido corrupción estructural y roturas obvias del core.
- Suites mínimas: `contract_suite`, `result_suite`, `runtime_suite`.
- Debe terminar rápido y priorizar señal alta.
- Puede usar dataset reducido, pero no puede omitir invariantes estructurales críticas.

### `standard`

- Objetivo: baseline operativo diario/pre-merge con cobertura estructural completa del core mínimo.
- Suites mínimas: todas las de `smoke` + `approval_suite` + `classification_suite` + `resilience_suite` + `tenant_onboarding_suite`.
- Debe acumular evidencia suficiente para diagnosticar degradaciones no críticas.

### `release_candidate`

- Objetivo: confirmar aptitud de release sobre baseline mínimo completo y artifacts auditables.
- Debe incluir **TODAS** las suites mínimas obligatorias.
- Debe usar fixtures/versiones congeladas y policy explícita de flakes.
- Debe dejar artifacts completos, métricas agregadas, referencias trazables y estado final de gating.

---

## Reglas por ambiente (`dev`, `staging`, `prod`)

| ambiente | objetivo de la corrida | reglas mínimas |
|---|---|---|
| `dev` | feedback rápido y validación estructural temprana | permite datos sintéticos, `smoke` obligatorio, `standard` recomendado, puede tolerar suites `soft_gate` adicionales fuera del baseline |
| `staging` | validación pre-release y consistencia de integración | `standard` obligatorio; `release_candidate` recomendado antes de promover; artifacts y métricas completas obligatorias |
| `prod` | confirmación controlada de salud/regresión con riesgo mínimo | solo smoke/regresiones seguras, datos/fixtures no destructivos, sin efectos irreversibles, trazabilidad reforzada y referencias auditables obligatorias |

Reglas transversales:

- En `prod` no se habilitan corridas que puedan mutar efectos reales irreversibles sin control explícito.
- Un mismo `eval_case_id` puede correr en distintos ambientes, pero la corrida debe registrar `environment` como dimensión obligatoria.
- Los thresholds pueden endurecerse por ambiente, pero nunca relajarse para `release_candidate` en `staging` o validaciones equivalentes previas a release.

---

## Orden recomendado de ejecución de suites

1. `contract_suite`
2. `runtime_suite`
3. `result_suite`
4. `approval_suite`
5. `classification_suite`
6. `resilience_suite`
7. `tenant_onboarding_suite`

### Justificación del orden

- `contract_suite` primero porque detecta corrupción de schema/fingerprint antes de contaminar el resto.
- `runtime_suite` temprano porque valida la columna vertebral de correlación y estados.
- `result_suite` luego porque necesita contrato y runtime ya confiables.
- `approval_suite` y `classification_suite` después para validar enforcement especializado sobre base estructural sana.
- `resilience_suite` una vez confirmado el flujo base, porque evalúa retry/dedup/compensación contra una línea estable.
- `tenant_onboarding_suite` al final porque valida baseline operativo integral del tenant y depende de varias verdades previas.

---

## Datos/fixtures mínimas para correr regresión

La regresión mínima DEBE contar, como mínimo, con estas fixtures/versiones controladas:

1. un contrato read-only baseline válido;
2. un contrato mutation reversible válido;
3. un contrato mutation con riesgo alto o irreversible controlado para paths defensivos;
4. snapshot de policy baseline;
5. snapshot de clasificación baseline;
6. approval request/decision válidos para `pre_execution`, `pre_application` y `double`;
7. execution records con lifecycle válido e inválido controlado;
8. result records exitosos, fallidos y con evidencia incompleta;
9. fixture de duplicate attempt con misma `idempotency_key`/fingerprint material;
10. fixture de `unknown_outcome` post-timeout;
11. fixture de compensación parcial o total;
12. tenant baseline `single_user` operable;
13. tenant incompleto no operable;
14. conjunto mínimo de `trace_id`, `event_id`, `execution_id`, `result_id`, `approval_id` correlacionables.

Reglas mínimas de fixtures:

- Deben ser determinísticas, versionadas y reusables entre niveles de corrida.
- Deben permitir reproducir tanto pass cases como fail cases esperados.
- Deben evitar secretos y datos sensibles innecesarios.
- Deben incluir fingerprints y snapshots suficientes para detectar drift estructural.

---

## Criterios de pase por suite

### Regla general

Una suite pasa solo si:

1. ejecutó todos sus casos obligatorios no quarantined,
2. no registró `hard_fail` en invariantes de esa suite,
3. produjo artifacts obligatorios completos,
4. dejó métricas mínimas de corrida,
5. conservó trazabilidad hacia IDs, events y evidence refs relevantes.

### Criterios mínimos por suite

- `contract_suite`: 0 incompatibilidades de schema/fingerprint/refs obligatorias.
- `result_suite`: 0 inconsistencias duras entre contrato, ejecución y `result_record`; 100% de evidence refs requeridas presentes en casos pass.
- `approval_suite`: 0 bypasses de governance; 0 transiciones ilegales aceptadas.
- `classification_suite`: 0 entregas incompatibles con clasificación efectiva; 0 redacciones faltantes en casos bloqueantes.
- `runtime_suite`: 0 transiciones estructuralmente inválidas aceptadas; 100% de corridas con correlación mínima completa.
- `resilience_suite`: 0 re-aplicaciones inseguras sobre duplicate/unknown outcome; 100% de compensaciones trazadas cuando el caso las exige.
- `tenant_onboarding_suite`: 0 tenants incompletos marcados operables; 0 tenants completos correctamente rechazados sin razón válida.

---

## Criterios de fail-fast vs continue-on-failure

### Debe existir `fail-fast` para corrupción estructural

La corrida DEBE detenerse inmediatamente cuando detecta cualquiera de estas condiciones:

- incompatibilidad de schema que vuelve inválidos los artifacts del core,
- mismatch estructural de fingerprint que invalida correlación material,
- incompatibilidad de state machine que hace ilegibles o imposibles las transiciones normativas,
- imposibilidad de resolver referencias mínimas necesarias para interpretar el caso,
- corrupción del manifest de fixtures que impide saber qué se ejecutó.

Estas condiciones son `hard_fail` estructurales y bloquean el nivel en curso.

### NO debe existir `fail-fast` para issues de calidad no críticas

Los siguientes casos NO deben cortar la corrida; deben seguir acumulando evidencia:

- score degradado sin ruptura estructural,
- warning de completitud menor,
- calidad subóptima de output no bloqueante,
- timing peor al esperado sin violar umbral duro,
- issue aislado de clasificación advisory fuera del baseline hard gate,
- inconsistencia menor de metadata no material.

Regla: si el sistema todavía puede producir evidencia válida y comparable, la corrida debe continuar.

---

## Artefactos obligatorios que debe dejar una corrida de regresión

Toda corrida de regresión mínima DEBE dejar, como mínimo:

1. **resumen de corrida** con nivel, suites ejecutadas, estado final y timestamps;
2. **bundle de evidencias** por suite/caso (`evidence_refs`, snapshots, outputs relevantes);
3. **bundle de fallos** con `hard_fail`, `soft_fail`, warnings y reason codes;
4. **bundle de métricas** de la corrida;
5. **referencias a traces/ids**: `trace_id`, `event_id`, `execution_id`, `result_id`, `approval_id`, `tenant_id` cuando aplique;
6. **manifest de fixtures/artifacts evaluados** con versiones/fingerprints;
7. **estado de gating** por suite y global;
8. **registro de reruns/quarantine** cuando existan flakes.

Regla dura: no existe corrida “válida” si solo deja pass/fail agregado sin artifacts trazables.

---

## Métricas mínimas de la corrida

La corrida DEBE capturar, como mínimo, estas métricas:

- cantidad total de suites ejecutadas,
- cantidad total de casos ejecutados,
- cantidad de `pass`, `hard_fail`, `soft_fail`, `warning`, `quarantined`,
- tiempo total de corrida,
- tiempo por suite,
- tasa de cobertura de artifacts obligatorios,
- tasa de correlación válida entre records (`contract`/`execution`/`result`/`approval`),
- cantidad de reruns por flake,
- cantidad de casos con `unknown_outcome`,
- cantidad de casos con compensación requerida vs completada,
- cantidad de casos con evidencia incompleta,
- cantidad de referencias a traces/ids faltantes,
- distribución por ambiente,
- distribución por dominio,
- estado de gating final.

Estas métricas satisfacen el mínimo de trazabilidad y observabilidad exigido por `source-of-truth/memory-analytics-and-observability.md` al cubrir negocio, orquestación, ejecución, calidad y rendimiento.

---

## Gestión de flakes / no determinismo

La regresión mínima DEBE asumir que puede existir no determinismo controlado, pero NO puede normalizarlo sin disciplina.

### Política mínima v1

- Un caso puede rerunearse hasta **2 veces** después del intento original; máximo total: **3 intentos**.
- Solo se permite rerun automático para fallos potencialmente no deterministas y nunca para corrupción estructural.
- Si un caso pasa en rerun pero falló antes, debe marcarse `flaky_detected`.
- Si un mismo caso flakea en **2 corridas estándar consecutivas** o **2 de las últimas 5 corridas elegibles**, debe proponerse quarantine.
- Si un caso flakea en `release_candidate`, la suite queda al menos en estado degradado y requiere decisión explícita según severidad.

### Regla de bloqueo de release

- Un flake en caso `hard_gate` crítico **bloquea release** hasta confirmar que no es defecto real o hasta quarantinarlo formalmente con reemplazo/cobertura equivalente aprobada.
- Un flake en caso `soft_gate` no bloquea automáticamente, pero debe quedar visiblemente reportado.

---

## Reglas para quarantined tests / evals inestables

- Un caso solo puede entrar en quarantine si existe evidencia de inestabilidad repetida y diagnóstico preliminar registrado.
- Quarantine NO elimina el caso del catálogo; lo mueve fuera del baseline blocking y exige ticket/razón/versionado.
- Todo quarantined case debe conservar:
  - `eval_case_id`,
  - motivo de quarantine,
  - fecha de entrada,
  - owner,
  - severidad,
  - criterio de salida.
- Un caso `hard_gate` quarantined requiere cobertura compensatoria equivalente o criterio explícito de aceptación de riesgo.
- `release_candidate` no debe aprobar silenciosamente con quarantined críticos sin decisión explícita y artifact de riesgo.
- Un caso quarantined debe revisarse periódicamente; si no tiene plan de salida, el baseline se considera degradado.

---

## Tests borde de la propia regresión (mínimo 15)

Estos casos prueban la **regresión como sistema de validación**, no solo el producto subyacente.

1. corrida sin manifest de fixtures debe fallar duro;
2. corrida con suite requerida ausente debe fallar;
3. corrida `release_candidate` sin todas las suites mínimas debe fallar;
4. corrida con artifacts incompletos debe fallar aunque los assertions lógicos pasen;
5. corrida sin `trace_id`/refs correlacionables debe fallar;
6. corrida con mismatch entre nivel declarado y suites realmente ejecutadas debe fallar;
7. corrida que marca pass pese a `hard_fail` estructural debe fallar;
8. corrida que corta por un `soft_fail` no crítico debe fallar contra la política;
9. corrida con reruns por encima del máximo permitido debe fallar;
10. corrida que quarantinea automáticamente sin evidencia histórica debe fallar;
11. corrida que reusa `eval_case_id` para un oracle materialmente distinto debe fallar;
12. corrida con métricas mínimas incompletas debe fallar;
13. corrida que no conserva references a `execution_id`/`result_id`/`approval_id` cuando aplican debe fallar;
14. corrida que mezcla artifacts de tenants distintos sin declararlo debe fallar;
15. corrida que ejecuta casos mutativos inseguros en `prod` debe fallar;
16. corrida que no diferencia `execution_completed` de `application_completed` debe fallar;
17. corrida que trata duplicate/unknown outcome como retry ciego debe fallar;
18. corrida que deja quarantined crítico invisible en `release_candidate` debe fallar.

---

## Criterios de aceptación de la regresión mínima del core

La regresión mínima del core v1 se considera aceptada solo si se cumplen simultáneamente todas estas condiciones:

1. existe este baseline documentado con suites, niveles, dominios, thresholds y artifacts obligatorios;
2. existen las siete suites mínimas requeridas:
   - `contract_suite`
   - `result_suite`
   - `approval_suite`
   - `classification_suite`
   - `runtime_suite`
   - `resilience_suite`
   - `tenant_onboarding_suite`
3. `release_candidate` incluye TODAS las suites mínimas;
4. la política de `fail-fast` corta solo ante corrupción estructural (`schema`, `fingerprint`, `state machine`, referencias mínimas incompatibles);
5. la política de continue-on-failure sigue acumulando evidencia para issues no críticas;
6. toda corrida deja resumen, evidencias, fallos, métricas y referencias a traces/ids;
7. existe política explícita de flakes con máximo de reruns, criterio de quarantine y criterio de bloqueo de release;
8. queda explícita la diferencia entre baseline estructural mínimo y cobertura exhaustiva futura;
9. la matriz de cobertura mínima demuestra cobertura de `contract`, `result`, `approval`, `classification`, `runtime`, `resilience`, `tenant` y `observability`;
10. la especificación define criterios verificables de pase, fail y artifacts requeridos por suite.

Si cualquiera de estas condiciones falta, la regresión mínima del core NO está suficientemente especificada para A.6.
