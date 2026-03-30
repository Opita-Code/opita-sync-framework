# E.1 — Full regression and hardening baseline

## Objetivo

E.1 existe para fijar el baseline **implementable y auditable** de regresión integral y hardening del **engine + surface** ya cerrados en C.6 y D.6. Su función es dejar un piso duro reusable para futuras iteraciones sin reabrir seams del motor, sin exigir apply real y sin introducir distribution layer.

Este bloque no busca cobertura infinita del producto. Busca una definición suficiente, rigurosa y operable para afirmar si el baseline reusable sigue sano o si ya entró en drift estructural.

## Principios del baseline de regresión y hardening

1. **Cobertura integral del baseline actual, no del producto infinito.** “Regresión completa” en Fase E significa regresión integral del **engine + surface** ya cerrados, no exploración exhaustiva de todas las variantes futuras.
2. **Verdad canónica primero.** PostgreSQL, artifacts persistidos y records canónicos siguen siendo la base normativa; observabilidad derivada es apoyo, nunca reemplazo.
3. **Hardening orientado a invariantes.** El baseline debe detectar rápido roturas de contrato, correlación, governance, evidence trail y resolución de capabilities.
4. **Fail-fast selectivo.** Debe existir fail-fast para corrupción estructural del baseline reusable. No debe existir fail-fast para degradación no crítica si la evidencia sigue siendo suficiente.
5. **Acumulación de evidencia obligatoria.** Si una corrida pierde el evidence trail mínimo, el baseline no puede considerarse sano aunque varios checks individuales hayan pasado.
6. **Distinción normativa de severidad.** Toda corrida full debe distinguir al menos entre `hard_fail`, `soft_fail` y `warning`.
7. **Sin dependencia de distribution layer.** Distribution layer no forma parte de la regresión full actual y no puede aparecer como prerequisito para considerar sano el baseline.
8. **Sin apply real obligatorio.** La salud del baseline reusable debe poder demostrarse sin convertir apply real en requisito normativo de E.1.

## Qué significa “regresión completa” en el scope actual

En Fase E, “regresión completa” significa que una corrida full puede recorrer y validar con evidencia suficiente los dominios cerrados de engine + surface que hoy definen el baseline reusable.

Eso implica cubrir como mínimo:

- contrato/compilador;
- policy;
- approvals/governance;
- runtime durable;
- result/outcome;
- event log y evidencia;
- registry/resolution;
- intake;
- proposal;
- preview/simulation;
- inspection/recovery;
- AI-friendly maintenance surface.

No implica:

- coverage infinita de todas las combinaciones de capabilities;
- pruebas de distribution layer, rollout o activation tenant-scoped;
- stress/performance profundo como criterio de pase de E.1;
- apply real como condición obligatoria del smoke path base.

La regresión full de E.1 debe probar que el baseline reusable sigue siendo **coherente, correlacionable, gobernado y reconstruible** de punta a punta dentro del boundary actual.

## Qué significa “hardening” en el scope actual

En Fase E, “hardening” significa endurecer los puntos donde una evolución futura podría romper el baseline reusable sin que el sistema lo detecte a tiempo.

Hardening en este scope implica:

- fijar invariantes críticas y sus gates de pase/falla;
- endurecer correlación mínima entre engine y surface;
- endurecer persistencia canónica obligatoria de policy, runtime, eventos y artifacts de surface;
- endurecer fail-safe y degradación segura;
- endurecer tratamiento de no determinismo y flakes;
- endurecer la política de severidad (`hard_fail`, `soft_fail`, `warning`) para que la corrida no colapse todo lo degradado en una única categoría;
- endurecer qué gaps pueden tolerarse sin declarar sano algo estructuralmente roto.

Hardening no significa reescribir arquitectura ni abrir nuevos seams. Significa **cerrar mejor lo que ya existe**.

## Suites mínimas obligatorias de regresión full

La corrida full de E.1 debe incluir como mínimo estas suites:

| suite_name | objetivo | dominios principales | blocking_mode por defecto |
|---|---|---|---|
| `contract_compiler_full_suite` | validar determinismo, versionado material/no material, compilación y correlación base | contract/compiler, result/outcome, runtime | `hard_gate` |
| `policy_governance_full_suite` | validar policy input, persistencia de decisión, approvals/governance y fail-closed | policy, approvals/governance, result/outcome | `hard_gate` |
| `runtime_full_suite` | validar lifecycle durable, estados, idempotencia, compensación y cierre canónico | runtime, result/outcome, governance | `hard_gate` |
| `evidence_and_eventlog_full_suite` | validar event log append-only, evidence trail y autonomía de la verdad canónica | event log/evidence, observability | `hard_gate` |
| `registry_resolution_full_suite` | validar cadena completa de resolución de capability y compatibilidad material | registry/resolution, runtime | `hard_gate` |
| `surface_intake_proposal_preview_full_suite` | validar intake, shaping, proposal, patchset, preview y simulación | intake, proposal, preview/simulation | `hard_gate` |
| `surface_inspection_recovery_full_suite` | validar inspection/recovery sin contradicción con verdad canónica | inspection/recovery, runtime, evidence | `hard_gate` |
| `surface_ai_maintenance_full_suite` | validar maintenance/debugging asistido sin bypass de governance | AI-friendly maintenance, governance, evidence | `hard_gate` |
| `end_to_end_smoke_full_suite` | validar el corredor mínimo integral engine + surface con evidencia reconstruible | engine + surface end-to-end | `hard_gate` |

### Regla de completitud

- Una corrida full de E.1 NO puede declararse completa si falta cualquiera de estas suites.
- Una suite puede degradar por `soft_fail` o `warning` sólo si el evidence trail completo sigue disponible y la degradación no rompe invariantes estructurales.
- Si una suite detecta corrupción estructural, debe disparar `hard_fail` y la corrida deja de ser apta para declarar baseline sano.

## Matriz de cobertura mínima engine + surface

| dominio | qué debe quedar cubierto | suites mínimas que lo cubren | evidencia mínima requerida |
|---|---|---|---|
| contract/compiler | compilación determinística, fingerprint, versionado material/no material, correlación con runtime | `contract_compiler_full_suite`, `end_to_end_smoke_full_suite` | `compiled_contract`, `compilation_report`, `contract_fingerprint`, refs correladas |
| policy | input canonizado, decisión persistida, mapping normativo a runtime/classification/approval | `policy_governance_full_suite`, `end_to_end_smoke_full_suite` | request canonizado o equivalente, `policy_decision_record`, reason codes |
| approvals/governance | gates previos, SoD, autoridad, invalidación por cambio material, bloqueos correctos | `policy_governance_full_suite`, `runtime_full_suite` | approvals refs, governance decision refs, evidence trail |
| runtime | lifecycle durable, separation execution/application, idempotencia, compensation, unknown outcome | `runtime_full_suite`, `end_to_end_smoke_full_suite` | `execution_record`, estados, transiciones, causalidad |
| result/outcome | consistencia de resultado, severity y evidence refs mínimas | `contract_compiler_full_suite`, `policy_governance_full_suite`, `runtime_full_suite` | `result_outcome`, reason codes, evidence refs |
| event log/evidence | append-only, correlación, reconstrucción del caso, independencia de observabilidad derivada | `evidence_and_eventlog_full_suite`, `end_to_end_smoke_full_suite` | `event_record`, IDs correlados, run summary |
| registry/resolution | cadena `capability -> bundle -> binding -> provider_ref`, compatibilidad material | `registry_resolution_full_suite`, `runtime_full_suite` | resolution refs, `bundle_digest`, `binding_id`, `provider_ref` |
| intake | boundary con conversación libre, evidencia mínima, reason codes, correlación | `surface_intake_proposal_preview_full_suite`, `end_to_end_smoke_full_suite` | `conversation_turn`, `intake_session`, refs de evidencia |
| proposal | `proposal_draft`, `source_intent_refs[]`, evidence refs, patchset material | `surface_intake_proposal_preview_full_suite` | `proposal_draft`, source refs, material diff refs |
| preview/simulation | `preview_candidate`, patchset vigente, simulación reproducible y no flaky | `surface_intake_proposal_preview_full_suite`, `end_to_end_smoke_full_suite` | `preview_candidate`, `simulation_result`, `inputs_refs[]` |
| inspection/recovery | vistas coherentes con records canónicos y recovery permitido | `surface_inspection_recovery_full_suite`, `end_to_end_smoke_full_suite` | `execution_inspection_view`, recovery refs, correlación |
| AI-friendly maintenance | maintenance/debugging asistido con governance explícita y sin bypass | `surface_ai_maintenance_full_suite`, `end_to_end_smoke_full_suite` | `semantic_debug_view`, `maintenance_action_candidate`, governance refs |

### Regla de interpretación de cobertura

- Un dominio no queda cubierto sólo porque exista una vista o un trace derivado.
- Un dominio sólo cuenta como cubierto si puede reconstruirse desde artifacts y records canónicos suficientes.
- Si falta evidencia mínima de un dominio, ese dominio queda descubierto aunque existan señales parciales en observabilidad.

## Hardening de invariantes críticos

Los siguientes invariantes quedan endurecidos como parte obligatoria de E.1:

1. no existe ejecución material sin `compiled_contract` persistido;
2. no existe `compiled_contract` válido sin correlación suficiente hacia runtime y evidence trail;
3. no existe decisión de policy válida si no queda persistida como `policy_decision_record` o evidencia normativa equivalente;
4. no existe cierre sano del runtime si falta al menos un `event_record` canónico material;
5. no existe éxito normativo si `result_outcome` contradice el estado real del corredor;
6. no existe resolución aceptable si `binding_id` o `provider_ref` son incompatibles materialmente;
7. no existe `proposal_draft` válido sin evidence refs mínimas y `source_intent_refs[]` suficientes;
8. no existe `preview_candidate` sano sobre patchset stale o inputs no resolubles;
9. no existe `execution_inspection_view` o `semantic_debug_view` aceptable si contradice records canónicos;
10. no existe `maintenance_action_candidate` aceptable sin `governance_requirements[]` explícitos cuando corresponda;
11. no existe `blocked` aceptable sin `reason_code` auditable;
12. no existe `unknown_outcome` aceptable sin evidence trail suficiente;
13. no existe dependencia normativa de observabilidad derivada para demostrar verdad;
14. no existe baseline sano si la corrida full pierde correlación mínima reconstruible.

### Política de severidad de invariantes

- Violación estructural del baseline reusable => `hard_fail`
- Degradación relevante pero con evidence trail suficiente y sin ruptura del baseline => `soft_fail`
- Señal menor, riesgo de calidad o deuda de robustez sin invalidación estructural => `warning`

## Hardening de policy/runtime/event log/registry/surface

### Policy

Debe endurecerse que:

- el input a policy provenga de campos canonizados y no de texto libre reinterpretado;
- la respuesta de Cerbos no alcance por sí sola: debe quedar persistida y correlada;
- decisiones no mapeables o incompletas degraden a fail-closed;
- approvals/governance no queden colapsados dentro de un outcome genérico sin trazabilidad.

### Runtime

Debe endurecerse que:

- `execution_completed` no equivalga automáticamente a `application_completed`;
- retries técnicos, retries de ejecución y replay de auditoría se mantengan separados;
- la deduplicación no rompa correlación ni duplique hechos materiales;
- blocked, failed, compensation y unknown outcome se distingan semánticamente;
- el runtime no cierre sano si pierde persistencia canónica mínima.

### Event log

Debe endurecerse que:

- el event log siga siendo append-only y suficiente para reconstrucción;
- la ausencia de traces derivados NO invalide una corrida si el baseline canónico está íntegro;
- la presencia de traces derivados NO sane una corrida sin PostgreSQL/event records suficientes;
- replays o retries no dupliquen eventos materiales sin causalidad explícita.

### Registry

Debe endurecerse que:

- la cadena `capability_manifest -> bundle_digest -> binding -> provider_ref` sea reconstruible;
- compatibilidad material y vigencia de binding sean parte de la evaluación, no sólo existencia nominal;
- resolución incompleta o incompatible bloquee ejecución;
- el registry no introduzca atajos implícitos fuera del seam ya cerrado.

### Surface

Debe endurecerse que:

- intake no salte directo a apply ni a artifacts gobernados sin evidencia suficiente;
- proposal y preview mantengan separación material y correlación explícita;
- inspection/recovery lean verdad canónica y no la reemplacen;
- debugging y maintenance asistidos no se conviertan en autoridad implícita de ejecución;
- apply real no aparezca como requisito normativo del baseline full.

## Hardening de fail-safe y degradación segura

La corrida full de E.1 debe fijar una política explícita de fail-safe y degradación segura:

### Casos que deben disparar `hard_fail`

- corrupción estructural del contrato o pérdida de determinismo con el mismo input;
- respuesta de policy sin persistencia normativa suficiente;
- cierre de runtime sin `event_record` canónico material;
- evidencia sólo en observabilidad y no en PostgreSQL/artifacts canónicos;
- imposibilidad de reconstruir el smoke path mínimo por IDs correlados;
- apply real apareciendo como requisito obligatorio del baseline;
- pérdida del evidence trail mínimo de la corrida full.

### Casos que pueden degradar a `soft_fail`

- inconsistencia no estructural en una vista derivada mientras la verdad canónica sigue íntegra y auditable;
- ausencia de observabilidad derivada cuando artifacts y event log siguen completos;
- warning repetible de calidad del preview/simulación sin drift estructural del patchset ni pérdida de evidence trail;
- debt de robustez en una surface asistida que sigue bloqueada correctamente por governance.

### Casos que pueden quedar como `warning`

- señales tempranas de complejidad o mantenibilidad;
- lentitud o fricción operativa que no afecte verdad, correlación ni governance;
- evidencia adicional deseable, pero no mínima;
- ruido controlado en salida derivada sin impacto estructural.

### Regla clave

- Debe existir fail-fast para corrupción estructural del baseline reusable.
- No debe existir fail-fast para degradación no crítica si la evidencia sigue siendo suficiente.

## Artefactos mínimos de una corrida full

Toda corrida full de E.1 debe dejar, como mínimo, estos artifacts auditables:

1. `full_run_manifest` con suites ejecutadas, versión de fixtures y ambiente;
2. `full_run_summary` con estado agregado y severidades (`hard_fail`, `soft_fail`, `warning`);
3. bundle de métricas agregadas por suite y por dominio;
4. lista de failures/warnings con reason codes y referencias;
5. evidencia de `compiled_contract` + `compilation_report` para casos mínimos usados;
6. evidencia de `policy_decision_record` o equivalente normativa persistida;
7. evidencia de `execution_record` y estados terminales/relevantes;
8. evidencia de al menos un `event_record` material por corredor validado;
9. evidencia de resolution con `capability_id`, `bundle_digest`, `binding_id`, `provider_ref`;
10. artifacts de surface (`conversation_turn`, `intake_session`, `proposal_draft`, `preview_candidate`, `simulation_result`) para los casos aplicables;
11. artifacts de inspection/debug/maintenance para los casos aplicables;
12. mapa mínimo de correlación por IDs para reconstruir el smoke path;
13. registro explícito de gaps aceptados y gaps rechazados de la corrida.

Sin estos artifacts, la corrida puede haber producido señales útiles, pero no puede considerarse baseline full sano.

## Métricas mínimas del baseline

La corrida full de E.1 debe poder informar, como mínimo:

- porcentaje de suites mínimas ejecutadas vs requeridas;
- porcentaje de dominios cubiertos con evidencia suficiente;
- conteo de `hard_fail`, `soft_fail` y `warning`;
- porcentaje de casos con `policy_decision_record` persistido correctamente;
- porcentaje de casos con `event_record` material persistido correctamente;
- porcentaje de smoke paths reconstruibles de punta a punta por IDs;
- porcentaje de simulaciones reproducibles sin drift con igual input;
- tasa de duplicados materiales detectados en replay/retry;
- porcentaje de vistas de inspection/debug consistentes con records canónicos;
- porcentaje de maintenance candidates con governance requirements completos;
- conteo de casos donde observabilidad derivada falta pero baseline canónico sigue íntegro;
- conteo de casos inválidos por depender sólo de observability y no de PostgreSQL.

## Política de flakes/no determinismo en E.1

E.1 debe tratar el no determinismo como problema de baseline, no como detalle cosmético.

### Reglas

1. un caso repetido con mismos inputs y mismo baseline no puede cambiar outcome material sin explicación normativa;
2. una simulación que cambia sin cambio de inputs debe considerarse sospecha de flake estructural;
3. un flake que afecte determinismo de compilación, correlación, policy persistida o event log material debe escalar a `hard_fail`;
4. un flake que afecte sólo salida derivada o presentación puede degradar a `soft_fail` o `warning` según evidencia disponible;
5. reintentar un caso flaky no puede borrar ni ocultar el primer outcome observado;
6. la corrida full debe registrar explícitamente qué casos fueron considerados flakes, con causa observada o hipótesis documentada.

### Política normativa

- “flaky simulation que cambia sin cambio de inputs” es caso mínimo obligatorio y no puede normalizarse como ruido aceptable por defecto.
- “compilación determinística falla con mismo input” es `hard_fail` automático.
- si el no determinismo impide reconstruir una decisión o un corredor, el baseline no puede pasar.

## Gaps aceptables al cerrar E.1

Pueden aceptarse al cerrar E.1 sólo gaps que NO invaliden el baseline reusable, por ejemplo:

- warnings de ergonomía, claridad o fricción operativa;
- cobertura adicional deseable de variantes no críticas de capabilities;
- mejoras futuras de visualización derivada o dashboards;
- mayor profundidad de performance/scale fuera del objetivo de E.1;
- documentación operativa complementaria que pertenezca más a E.3 que a E.1;
- refinamientos adicionales de grouping o reporting mientras la señal estructural ya sea suficiente.

Todo gap aceptable debe quedar explícito y no puede ocultar una rotura de correlación, governance o evidencia.

## Gaps NO aceptables al cerrar E.1

No son aceptables al cerrar E.1:

- falta de alguna suite mínima obligatoria;
- falta de cobertura de cualquiera de los dominios mínimos;
- ausencia de distinción entre `hard_fail`, `soft_fail` y `warning`;
- dependencia de observabilidad derivada para demostrar verdad del baseline;
- imposibilidad de reconstruir el smoke path mínimo por IDs y artifacts;
- ausencia de evidence trail mínimo;
- apply real exigido como requisito para declarar sano el baseline;
- pérdida de persistencia normativa de policy, runtime o event log;
- flakes estructurales no clasificados ni tratados;
- registry/resolution aceptando compatibilidades materiales inválidas;
- inspection/debug/maintenance contradiciendo verdad canónica sin fail explícito.

## Tests borde mínimos (al menos 20)

La definición mínima de E.1 debe incluir, como piso, estos tests borde:

1. **fingerprint cambia sin cambio material** => debe fallar por deriva indebida del compilador.
2. **compilación determinística falla con mismo input** => `hard_fail` automático.
3. **Cerbos responde pero falta persistencia de `policy_decision_record`** => la suite de policy no puede pasar.
4. **policy input canonizado incompleto** => debe degradar a fail-closed, no a allow implícito.
5. **runtime cierra sin `event_record` canónico** => baseline no sano.
6. **event log correcto pero trace derivado ausente** => puede pasar si la evidencia canónica es suficiente.
7. **evidencia sólo en observability y no en PostgreSQL** => `hard_fail`.
8. **registry resuelve binding incompatible** => debe bloquearse como incompatibilidad material.
9. **registry encuentra capability_id pero no `provider_ref` compatible** => no puede promocionarse a ejecutable.
10. **surface genera `proposal_draft` sin evidence refs** => debe fallar la suite de surface.
11. **surface genera `proposal_draft` sin `source_intent_refs[]` suficientes** => falla estructural del corredor de proposal.
12. **preview simula con patchset stale** => la simulación no puede considerarse válida.
13. **flaky simulation que cambia sin cambio de inputs** => debe registrarse como flake; si compromete decisión, no pasa.
14. **inspection view contradice `execution_record`** => falla la suite de inspection/recovery.
15. **maintenance action candidate sin governance requirements** => falla la suite AI-friendly maintenance.
16. **replay técnico duplica event log material** => falla por duplicación indebida.
17. **deduplicación rompe correlación** => falla aunque evite duplicado superficial.
18. **unknown outcome sin evidence trail** => baseline no sano.
19. **blocked sin reason_code** => falla estructural.
20. **soft fail agregado como hard fail** => debe detectarse como mala clasificación de severidad.
21. **warning agregado como soft fail sin justificación** => debe marcarse inconsistencia de severidad.
22. **observability caída pero baseline canónico íntegro** => no debe forzar hard fail por sí sola.
23. **end-to-end smoke no puede reconstruirse por IDs** => `hard_fail`.
24. **result_outcome exitoso contradice bloqueo de governance** => falla por incoherencia semántica.
25. **apply real aparece como requisito del baseline** => debe fallar la definición de E.1 por contradicción de scope.

## Criterios de aceptación de E.1

E.1 puede considerarse cerrado cuando:

1. queda definida “regresión completa” como regresión integral del **engine + surface** cerrados, no como cobertura infinita del producto;
2. quedan definidas y obligatorias todas las suites mínimas del baseline full;
3. la matriz de cobertura mínima cubre contract/compiler, policy, approvals/governance, runtime, result/outcome, event log/evidence, registry/resolution, intake, proposal, preview/simulation, inspection/recovery y AI-friendly maintenance;
4. queda explícita la política de severidad con `hard_fail`, `soft_fail` y `warning`;
5. queda explícito el fail-fast para corrupción estructural del baseline reusable;
6. queda explícito que no debe existir fail-fast para degradación no crítica si la evidencia sigue siendo suficiente;
7. queda explícito que si la corrida pierde evidence trail mínimo, el baseline no puede considerarse sano;
8. queda explícito que si la corrida depende de observabilidad derivada para demostrar verdad, no puede pasar;
9. queda explícito que distribution layer no forma parte de la regresión full actual;
10. quedan definidos artifacts mínimos, métricas mínimas, política de flakes y tests borde suficientes como para tratar E.1 como baseline implementable de cierre.
