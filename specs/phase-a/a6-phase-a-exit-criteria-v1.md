# A.6 — Criterios de salida de Fase A v1

## Principios base

- Fase A existe para cerrar la **verdad ejecutable** del core antes de construir el kernel.
- Cerrar Fase A NO significa tener kernel implementado completo; significa llegar a un punto donde la construcción ya no dependa de ambigüedades conceptuales bloqueantes.
- La salida de Fase A exige cierre **a nivel diseño v1** de A.1, A.2, A.3, A.4, A.5 y A.6.
- A.6 completa la capa de evaluabilidad mínima del sistema: evals, baseline de testing y regresión mínima.
- Después de Fase A, el macro-flujo correcto es: **Fase B** para investigación/comparación de decisiones irreversibles, luego **Fase C** para construcción del kernel.
- Todo criterio de salida debe privilegiar consistencia sistémica, auditabilidad y ausencia de contradicciones entre bloques.
- La regla rectora es simple: si una decisión faltante puede frenar o sesgar de manera material Fase B o Fase C, Fase A NO está cerrada.

---

## Qué significa “verdad ejecutable cerrada” en este proyecto

En Opyta Sync, “verdad ejecutable cerrada” significa que el sistema ya tiene definido, de manera verificable y no ambigua, el conjunto mínimo de contratos, objetos, estados, reglas, eventos, approvals, resultados, runtime y evals necesarios para:

1. describir qué debe existir en el core;
2. compilar esa verdad a un formato canónico determinístico;
3. razonar si un flujo sería válido, bloqueado o fallido aun antes de implementar el kernel completo;
4. evaluar luego la implementación contra criterios ya fijados, en lugar de inventarlos durante la construcción.

La verdad ejecutable está cerrada cuando, como mínimo:

- existe formato de authoring y formato compilado decididos;
- existe envelope canónico y naming global cerrado;
- contrato, resultado, approval y runtime tienen semántica formal compatible;
- existe correlación explícita entre IDs, snapshots, eventos y evidencia;
- existe marco formal para evaluar comportamiento esperado;
- existe baseline de regresión mínima que indique qué NO puede romperse.

La verdad ejecutable NO está cerrada si todavía hay cualquiera de estos defectos:

- dos documentos permiten interpretaciones incompatibles para el mismo flujo;
- falta una decisión material sobre formato, lifecycle, approval, runtime o evaluabilidad;
- no puede determinarse qué evidencia mínima probaría una implementación correcta;
- no puede distinguirse entre gap tolerable y ambigüedad bloqueante.

---

## Preconditions obligatorias para considerar cerrada Fase A

Todas las siguientes preconditions son **must_pass**:

1. A.1-A.6 existen como especificaciones v1 publicadas y coherentes entre sí.
2. A.1-A.6 están cerrados a nivel diseño v1, aunque no estén implementados de punta a punta.
3. El formato de verdad ejecutable está decidido explícitamente: YAML para authoring y JSON canónico determinístico para consumo/runtime.
4. El envelope canónico y el naming global están cerrados y no quedan como “a definir después”.
5. La unidad operativa central (`intent_contract`) y su relación con approvals, runtime y resultados está cerrada a nivel semántico.
6. Existe coherencia explícita entre contrato, resultado, approvals, runtime y evals.
7. Existe baseline formal de evals y baseline formal de regresión mínima.
8. Existe evidencia mínima definida para demostrar cierre de diseño, consistencia y evaluabilidad.
9. Existe lista explícita de temas diferidos a Fase B, C, D y E.
10. No quedan ambigüedades bloqueantes para iniciar Fase B ni para arrancar la construcción en Fase C cuando termine Fase B.

Preconditions **should_pass**:

- existe mapa de trazabilidad entre documentos A.1-A.6;
- existe vocabulario normalizado reutilizable para outcomes, reason codes, states y gates;
- existe criterio explícito para distinguir deuda tolerable de deuda bloqueante.

Preconditions **advisory**:

- existe propuesta inicial de orden recomendado de validación al comenzar Fase C;
- existe recomendación de ownership del sign-off por bloque.

---

## Checklist de salida por bloque (A.1-A.6)

| Bloque | Condición de salida | Evidencia mínima requerida | Gate |
|---|---|---|---|
| A.1 | Objetos canónicos, ownership, versionado, formato y naming cerrados a nivel v1 | Documentos A.1 publicados; lista final de objetos first-class; decisión YAML→JSON; envelope canónico; convenciones de naming y validación | must_pass |
| A.2 | Contrato de intención/inspección definido con schema, estados, compilación, fingerprint y versionado | Documentos A.2 publicados; criterios de aceptación; eventos canónicos; tests borde; regla formal de `executable` | must_pass |
| A.3 | Tipos de resultado, outcomes, evidencia mínima y lifecycle formalmente definidos | Documentos A.3 publicados; matriz por tipo; failure taxonomy; reglas de clasificación/redacción vinculadas; tests borde | must_pass |
| A.4 | Subsystem de approvals definido con modos, SoD, invalidación material, SLA y evidencia | Documentos A.4 publicados; matriz de approval mode; lifecycle; eventos; reason codes; tests y criterios de aceptación | must_pass |
| A.5 | Tenant, runtime, correlación de IDs, idempotencia, retries y compensación cerrados a nivel v1 | Documentos A.5 publicados; schema de tenant; runtime states; eventos; deduplicación; compensación; criterios de aceptación | must_pass |
| A.6 | Framework de evals, baseline de regresión mínima y criterio formal de salida de Fase A definidos | `a6-evals-framework-v1.md`, `a6-core-regression-baseline-v1.md`, este documento; suites mínimas; gates de evaluabilidad; evidencia de consistencia | must_pass |

### Regla de bloque completo

Fase A solo puede declararse cerrada si **todos** los bloques A.1-A.6 pasan su condición de salida en nivel `must_pass`.

Un bloque puede tener mejoras futuras pendientes y seguir cerrado solo si esas pendientes:

- no cambian la semántica v1 ya decidida;
- no reabren contratos entre bloques;
- no impiden evaluar una implementación futura contra la verdad ya definida.

---

## Gates de calidad obligatorios

### `must_pass`

- No existen contradicciones normativas entre documentos A.1-A.6.
- Todo término crítico (`contract`, `approval`, `execution`, `application`, `result`, `classification`, `eval`) tiene semántica diferenciada.
- Todo objeto o artifact crítico tiene identidad, estado o referencia suficiente para auditoría.
- Todo flujo mutativo distingue correctamente `execution_completed` de `application_completed`.
- Todo cambio material relevante invalida o reevalúa lo que corresponda según bloque.
- La regresión mínima está definida con suites obligatorias, artifacts requeridos y política de pase/fallo.
- Los evals están definidos con `eval_case_id`, oracle, expected outcome, severity, gating mode y evidencia.

### `should_pass`

- Existe terminología consistente entre tablas, ejemplos y criterios de aceptación.
- Cada bloque incluye tests borde suficientes para detectar regresión conceptual.
- La evidencia mínima por bloque está expresada de forma reusable para tooling futuro.
- Los criterios de pase evitan lenguaje aspiracional del tipo “idealmente”, “más adelante se verá”.

### `advisory`

- Existe una matriz resumida de dependencies entre bloques.
- Existe recomendación de priorización para convertir esta verdad en fixtures/casos ejecutables en Fase C.
- Existe sugerencia de owners para revisión de cierre.

---

## Gates de consistencia entre bloques

Los siguientes gates son obligatorios y todos son `must_pass`:

1. **A.1 ↔ A.2**: el contrato usa envelope, naming, versionado y reglas de compilación compatibles con la verdad ejecutable definida en A.1.
2. **A.2 ↔ A.3**: cada contrato ejecutable puede resolver un `result_type` permitido y un outcome evaluable según A.3.
3. **A.2 ↔ A.4**: approvals se calculan sobre contrato compilado, fingerprint y snapshots relevantes; no sobre intención ambigua.
4. **A.2 ↔ A.5**: runtime no ejecuta contratos no `executable` y conserva correlación estable con `contract_id` y fingerprint.
5. **A.3 ↔ A.4**: floors, overrides, governance decisions y restricciones de resultado son compatibles con policy y approval modes.
6. **A.3 ↔ A.5**: resultado, evidencia, clasificación y cierre terminal se correlacionan con `execution_record` sin ambigüedad.
7. **A.4 ↔ A.5**: release de approvals y release/runtime/application mantienen separación formal y trazabilidad completa.
8. **A.3/A.4/A.5 ↔ A.6**: evals y regresión mínima pueden expresar y verificar outcomes, bloqueos, retries, compensación y trazabilidad usando las verdades previas sin inventar semántica nueva.

### Condición de bloqueo

Si cualquiera de estos cruces requiere “interpretación humana libre” para decidir qué documento manda, Fase A NO está cerrada.

---

## Gates de evaluabilidad (que el sistema ya sea evaluable aunque todavía no esté implementado)

Fase A puede cerrarse sin kernel completo solo si el sistema ya es **evaluable por diseño**.

### `must_pass`

- Existe definición formal de `eval_case` y `eval_run`.
- Existen familias mínimas de evals: `result_eval`, `policy_eval`, `approval_eval`, `classification_eval`, `runtime_eval`, `resilience_eval`.
- Existe baseline de regresión mínima con siete suites obligatorias.
- Existen fixtures/artifacts mínimos conceptuales para reconstruir casos.
- Existe criterio explícito de `pass`, `hard_fail`, `soft_fail` y `warning`.
- Existe trazabilidad mínima requerida a `contract_id`, `execution_id`, `result_id`, `approval_request_id`/`approval_decision_id` cuando aplica.
- Existe evidencia mínima suficiente para distinguir defecto de producto, bloqueo correcto por governance y degradación tolerable.

### `should_pass`

- Existe mapeo preliminar entre suites de regresión y familias de evals.
- Existe orden recomendado de ejecución de suites.
- Existe política de fail-fast para corrupción estructural y continue-on-failure para degradaciones no críticas.

### `advisory`

- Existe propuesta de priorización de casos golden vs edge cases.
- Existe convención sugerida de naming para manifests o fixtures de eval.

### Regla central

Si hoy no se puede escribir una suite futura sin tomar decisiones nuevas sobre oracle, evidencia o artifacts, entonces la evaluabilidad todavía NO está cerrada.

---

## Qué queda explícitamente fuera de Fase A

Queda fuera de Fase A, y por lo tanto NO bloquea su salida, todo lo siguiente siempre que esté explícitamente diferido:

### Diferido a Fase B

- comparación e investigación de decisiones irreversibles de arquitectura, almacenamiento, orquestación o integración;
- tradeoffs profundos entre alternativas de implementación del kernel;
- selección final de componentes cuando la decisión dependa de costo/rendimiento/operación real.

### Diferido a Fase C

- construcción efectiva del kernel;
- implementación de compiladores, validadores, state machines y ejecutores runtime;
- automatización inicial de suites, eval runners y pipelines asociados.

### Diferido a Fase D

- ampliación de cobertura, hardening operativo, observabilidad expandida y automatización avanzada;
- soporte exhaustivo de variantes no críticas, fixtures masivas y expansión de catálogo.

### Diferido a Fase E

- optimización de performance, escalado, tuning fino, cobertura extensiva multi-tenant y maduración de operación productiva.

### Regla de exclusión

Nada de lo anterior puede reintroducir dudas sobre la semántica base ya cerrada en Fase A. Si una decisión diferida cambia la verdad ejecutable central, entonces estaba mal diferida.

---

## Riesgos aceptables para salir de Fase A

Son aceptables solo si están documentados, acotados y no bloquean Fase B/C:

- falta de implementación completa del kernel;
- falta de automatización completa de regresión/evals;
- ausencia de tooling definitivo para compilar, correr o reportar evals;
- falta de cobertura exhaustiva de combinatorias, scale cases y escenarios probabilísticos;
- placeholders de datasets o fixtures concretas siempre que la estructura mínima ya esté definida;
- decisiones de performance, storage físico o deployment aún no resueltas cuando no alteran la verdad lógica v1.

### Cuándo un gap todavía es tolerable

Un gap es tolerable si cumple simultáneamente estas cuatro condiciones:

1. no cambia contratos semánticos ya fijados;
2. no impide diseñar o evaluar el kernel contra la verdad actual;
3. tiene fase destino explícita (B/C/D/E);
4. tiene impacto acotado a implementación, cobertura o optimización, no a definición conceptual.

---

## Riesgos NO aceptables para salir de Fase A

Son bloqueantes y por lo tanto impiden declarar Fase A cerrada:

- formato de verdad ejecutable todavía discutible o naming no cerrado;
- falta de cierre de cualquier bloque A.1-A.6 a nivel diseño v1;
- contradicción entre contrato, resultado, approvals, runtime y evals;
- ausencia de baseline mínimo de regresión;
- ausencia de definición formal de evaluabilidad;
- imposibilidad de determinar evidencia mínima de corrección;
- ambigüedad sobre cuándo una approval habilita ejecución vs aplicación;
- ambigüedad sobre correlación entre IDs/artifacts clave;
- temas diferidos sin lista explícita o diferidos de manera incorrecta;
- falta de sign-off o evidencia mínima para sostener el cierre.

### Cuándo un gap pasa a ser bloqueante

Un gap es bloqueante si ocurre cualquiera de estas condiciones:

1. obliga a redefinir una semántica core durante Fase B o C;
2. permite más de una interpretación válida para el mismo flujo crítico;
3. impide construir tests/evals sin inventar reglas nuevas;
4. rompe coherencia entre artifacts o lifecycles;
5. compromete seguridad, governance, auditabilidad o trazabilidad mínima.

---

## Evidencia mínima requerida para declarar Fase A cerrada

Para declarar cierre de Fase A debe existir, como mínimo, esta evidencia:

1. documento de estado o checklist mostrando A.1-A.6 cerrados a nivel diseño v1;
2. documentos A.1-A.6 publicados y referenciables;
3. decisión explícita del formato de verdad ejecutable y naming global;
4. criterios de aceptación publicados por bloque o subbloque relevante;
5. baseline de evals publicado;
6. baseline de regresión mínima publicado;
7. este documento de criterios de salida publicado;
8. tabla o registro de temas diferidos a Fase B/C/D/E;
9. registro de inconsistencias conocidas, si existen, con clasificación: tolerable o bloqueante;
10. sign-off recomendado con responsables y fecha.

### Evidencia complementaria `should_pass`

- matriz de trazabilidad A.1-A.6;
- lista consolidada de invariantes del core;
- mapping preliminar de suites ↔ fixtures ↔ artifacts.

### Evidencia `advisory`

- propuesta inicial de orden de implementación en Fase C;
- propuesta de secuencia de investigación para Fase B.

---

## Decision record / sign-off recomendado

Se recomienda registrar el cierre de Fase A en un decision record o acta de sign-off con estos campos concretos:

| Campo | Obligatorio | Descripción |
|---|---|---|
| `phase_id` | sí | `phase_a` |
| `phase_version` | sí | versión del cierre, por ejemplo `v1` |
| `decision_type` | sí | `phase_exit_signoff` |
| `status` | sí | `approved` / `approved_with_tolerated_gaps` / `rejected` |
| `decision_date` | sí | fecha del sign-off |
| `approved_by` | sí | responsables que firman el cierre |
| `reviewed_blocks` | sí | lista A.1-A.6 revisada |
| `must_pass_summary` | sí | resumen de gates duros cumplidos |
| `should_pass_summary` | sí | resumen de deuda no bloqueante |
| `tolerated_gaps` | sí | gaps permitidos con fase destino |
| `blocking_gaps` | sí | debe ser vacío para aprobar |
| `linked_artifacts` | sí | referencias a specs, checklist y status |
| `next_phase` | sí | `phase_b` |
| `notes` | no | observaciones adicionales |

### Regla de aprobación recomendada

- `approved`: todos los `must_pass` cumplidos y sin gaps bloqueantes.
- `approved_with_tolerated_gaps`: todos los `must_pass` cumplidos, sin gaps bloqueantes y con deuda explícita diferida.
- `rejected`: falta al menos un `must_pass` o existe cualquier gap bloqueante.

---

## Tests borde de criterios de salida (mínimo 12)

1. Si A.1-A.5 están cerrados pero A.6 no define regresión mínima, Fase A NO puede cerrarse.
2. Si existe framework de evals pero no hay `eval_case_id` estable, Fase A NO puede cerrarse.
3. Si el naming global sigue abierto a variantes incompatibles, Fase A NO puede cerrarse.
4. Si contrato y runtime discrepan sobre cuándo una ejecución es liberable, Fase A NO puede cerrarse.
5. Si approvals permiten interpretar igual `execution_released` y `application_released`, Fase A NO puede cerrarse.
6. Si la automatización todavía no existe pero los casos, artifacts y gates ya están definidos, ese gap SÍ es tolerable.
7. Si falta elegir storage físico del kernel pero no afecta semántica v1, ese gap SÍ es tolerable y se difiere a Fase B/C.
8. Si un tema supuestamente diferido a Fase C en realidad redefine fingerprint, lifecycle o approval semantics, ese gap es bloqueante y NO puede diferirse.
9. Si existen suites mínimas definidas pero sin artifacts obligatorios ni evidencia requerida, Fase A NO puede cerrarse.
10. Si el documento de sign-off marca `approved_with_tolerated_gaps` pero `blocking_gaps` no está vacío, el cierre es inválido.
11. Si hay contradicción entre A.3 y A.6 sobre outcome formal, Fase A NO puede cerrarse hasta resolverla.
12. Si ya puede formularse una batería de evals y regresión sin decisiones nuevas sobre verdad base, el sistema SÍ es evaluable por diseño.
13. Si no existe lista explícita de temas diferidos a Fase B/C/D/E, Fase A NO puede cerrarse.
14. Si un bloque tiene “cerrado” pero sus criterios de aceptación permiten interpretaciones libres, ese cierre es inválido.
15. Si todos los `must_pass` se cumplen y solo quedan deudas de tooling/automatización/cobertura extensiva, Fase A SÍ puede cerrarse.

---

## Criterios de aceptación del documento

Este documento se considera aceptado solo si se cumplen simultáneamente estas condiciones:

- define de forma explícita qué significa “verdad ejecutable cerrada” para Opyta Sync;
- establece preconditions obligatorias para cerrar Fase A;
- incluye checklist de salida por bloque A.1-A.6 con condición y evidencia;
- define gates `must_pass`, `should_pass` y `advisory`;
- fija gates de calidad, consistencia y evaluabilidad;
- explicita qué queda fuera de Fase A y a qué fase se difiere;
- distingue con precisión riesgos aceptables vs no aceptables;
- define exactamente cuándo un gap es tolerable y cuándo es bloqueante;
- define evidencia mínima requerida para declarar cierre;
- propone un sign-off concreto con campos normativos;
- incluye al menos 12 tests borde propios;
- no contradice el contexto rector de Fase A, Fase B y Fase C.

---

## Conclusión normativa

Fase A v1 puede declararse cerrada únicamente cuando la verdad ejecutable del core ya quedó formalmente especificada, evaluable y consistente entre A.1-A.6, aun sin kernel implementado completo. Si falta implementación pero no falta verdad, Fase A puede cerrar. Si falta verdad, aunque sobren documentos sueltos, NO puede cerrar.
