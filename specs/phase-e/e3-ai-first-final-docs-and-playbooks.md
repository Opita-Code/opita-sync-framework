# E.3 — AI-first final docs and operator/developer playbooks

## Objetivo

E.3 existe para fijar la capa documental final que vuelve a Opyta Sync **navegable, reconstruible y operable** por humans + IA sobre el baseline reusable ya cerrado en E.1 y E.2, sin convertir la documentación en marketing, sin crear una verdad paralela y sin reintroducir distribution layer bajo ningún nombre alternativo.

El propósito de E.3 no es agregar más superficie funcional, sino dejar un sistema documental final que enseñe el modelo mental correcto, permita reconstruir el corredor engine + surface desde cero y habilite operación, extensión y debugging gobernados sin contradicción con los artifacts normativos.

## Principios de la documentación IA-first final

1. **IA-first no significa marketing-friendly.** “IA-first docs” significa documentación navegable y reconstruible por humans + IA, no texto aspiracional o comercial.
2. **La documentación enseña arquitectura, no sólo nombres de archivos.** Debe transmitir boundaries, seams, invariantes, correlaciones y límites de autoridad.
3. **Truth artifacts primero.** La documentación explica y orienta; los artifacts normativos mandan.
4. **Separación explícita de mapas.** Debe existir separación clara entre mapa del sistema, mapa de artifacts canónicos, mapa de seams/flows, mapa de invariantes/fail-safes y mapa de evidencias/queries operativas.
5. **Playbooks diferenciados por rol.** Operator, developer y debugging/mantenimiento necesitan playbooks distintos porque consumen el sistema con objetivos distintos.
6. **Descubribilidad antes que completitud dispersa.** Debe ser posible encontrar rápido qué leer primero, qué leer después y cómo reconstruir un caso sin barrer todo el repositorio.
7. **Sin contradicción con seams normativos.** La documentación no puede redefinir contracts, seams o artefactos ya cerrados.
8. **Sin distribution layer en el mapa actual.** Distribution layer debe quedar explícitamente fuera del mapa documental final de esta fase.

## Qué significa “IA-first docs” en este proyecto

En este proyecto, “IA-first docs” significa una capa documental que permite que una persona o una IA lean Opyta Sync desde cero y reconstruyan correctamente:

- qué es engine y qué es surface;
- qué seams existen y cuáles son sus boundaries;
- qué artifacts son canónicos;
- qué IDs/correlaciones son obligatorias;
- qué invariantes no se pueden romper;
- cómo recorrer un caso end-to-end;
- qué queda explícitamente fuera del roadmap actual.

### Regla explícita de lectura para una IA desde cero

La documentación final debe incluir una regla de lectura equivalente a esta secuencia conceptual:

1. leer primero el **mapa del sistema** para entender boundary engine/surface;
2. leer luego el **mapa de artifacts canónicos** para saber qué objetos son verdad operativa;
3. leer después el **mapa de seams/flows** para entender cómo viaja un caso;
4. leer luego el **mapa de invariantes/fail-safes** para entender qué no puede romperse;
5. leer el **mapa de evidencias/queries operativas** para saber cómo reconstruir y auditar casos;
6. recién después leer los **playbooks por rol**;
7. si una doc contradice un artifact normativo, priorizar artifact normativo y marcar la doc como inconsistente.

Si la documentación final no permite esta lectura secuencial desde cero, NO cumple el objetivo IA-first de E.3.

## Mapa final de documentación IA-first

El mapa final debe quedar separado en al menos estas capas documentales:

### 1. Mapa del sistema

Debe cubrir:

- qué es engine;
- qué es surface;
- dónde está el boundary entre ambos;
- qué decisiones quedaron cerradas en A-D y profundizadas en E.1-E.2;
- qué partes del sistema son baseline reusable;
- qué queda fuera del roadmap actual, incluyendo distribution layer.

### 2. Mapa de artifacts canónicos

Debe cubrir:

- qué artifacts son canónicos;
- qué artifacts son vistas derivadas y cuáles NO son source of truth;
- qué rol cumple cada artifact en el corredor;
- cuáles son las correlaciones mínimas entre artifacts;
- cómo se distingue artifact normativo vs documentation/reference artifact.

### 3. Mapa de seams/flows

Debe cubrir:

- seam `intent -> compiled_contract`;
- seam `compiled_contract -> execution_record`;
- seam PEP/PDP con input canonizado;
- seam event log canónico / observabilidad derivada;
- seam `capability_manifest -> bundle_digest -> binding -> provider_ref`;
- seam proposal/preview/inspection/maintenance gobernados;
- boundaries fuertes de D.5 entre ayuda asistida, debug semántico, candidate y acción real.

### 4. Mapa de invariantes/fail-safes

Debe cubrir:

- invariantes que no pueden romperse;
- qué constituye corrupción estructural;
- qué degradaciones son aceptables y cuáles no;
- por qué observabilidad derivada no reemplaza verdad canónica;
- por qué apply directo desde chat o bypass de governance está prohibido.

### 5. Mapa de evidencias/queries operativas

Debe cubrir:

- qué evidencia mínima debe existir por corredor;
- qué IDs deben correlacionarse;
- cómo reconstruir un caso end-to-end;
- qué preguntas operativas mínimas debe poder responder operator/developer;
- cómo distinguir evidencia canónica de señal derivada.

## Información mínima que cada capa documental debe cubrir

### Capa: mapa del sistema

Información mínima:

- definición de engine y surface;
- relación entre ambos;
- seams cerrados y no reabribles;
- lectura de scope actual;
- exclusiones explícitas del roadmap.

### Capa: mapa de artifacts canónicos

Información mínima:

- listado mínimo de artifacts centrales: `intent_input`, `compiled_contract`, `compilation_report`, `execution_record`, `policy_decision_record`, `event_record`, `proposal_draft`, `preview_candidate`, `simulation_result`, `execution_inspection_view`, `semantic_debug_view`, `maintenance_action_candidate` donde aplique;
- semántica de cada artifact;
- qué artifacts son obligatorios para reconstrucción;
- qué artifacts son derivados o secundarios.

### Capa: mapa de seams/flows

Información mínima:

- qué entra y qué sale por cada seam;
- qué contracts gobiernan cada paso;
- qué IDs se transportan entre seams;
- qué fallas deben degradar o bloquear el corredor;
- qué no puede inferirse libremente entre seams.

### Capa: mapa de invariantes/fail-safes

Información mínima:

- invariantes críticas del baseline reusable;
- bloqueos estructurales;
- límites de automatización asistida;
- definición operativa de fail-closed;
- relación entre evidence trail y salud del baseline.

### Capa: mapa de evidencias/queries operativas

Información mínima:

- lista de IDs mínimos obligatorios;
- queries operativas mínimas por rol;
- caminos de reconstrucción de casos;
- señales de docs rotas vs engine/surface rotos;
- criterio de suficiencia de evidence trail.

## Playbooks de operator

Los playbooks de operator deben ser distintos del resto y cubrir al menos:

### 1. Cómo leer un caso end-to-end

- partir de `tenant_id`, `trace_id`, `execution_id`, `proposal_draft_id` o artifact equivalente;
- reconstruir intake/proposal/preview cuando aplique;
- correlacionar contrato, policy, runtime y event log;
- decidir si la evidencia es suficiente o si el caso queda degradado.

### 2. Cómo interpretar estados operativos

Debe explicar explícitamente:

- `blocked`;
- `awaiting_approval`;
- `failed`;
- `unknown_outcome`;
- `compensation`.

### 3. Cómo escalar correctamente

Debe cubrir:

- cuándo escalar por falta de evidencia;
- cuándo escalar por contradicción documental vs artifact normativo;
- cuándo escalar por inconsistencias de correlación;
- cuándo escalar por sospecha de bypass de governance.

### 4. Qué recovery está permitido y cuál no

Debe distinguir:

- recovery permitido según D.4/D.5;
- recovery que sólo puede candidatearse;
- recovery prohibido;
- acciones que parecen operativas pero salen del boundary actual.

### Regla fuerte del playbook de operator

El playbook de operator NO puede recomendar:

- apply directo desde chat;
- uso de observability como única verdad;
- mutación fuera del corredor gobernado;
- actiones de distribution layer.

## Playbooks de developer

Los playbooks de developer deben cubrir al menos:

### 1. Cómo extender capabilities sin romper seams

- cómo usar el starter kit definido en E.2;
- cómo respetar manifest/binding/provider compatibility;
- cómo respetar `contract_schema_version` y `supported_result_types`;
- cómo integrar proposal/preview/inspection/maintenance donde aplique.

### 2. Cómo leer contracts/runtime/policy/event log/registry/surface

- qué artifacts mirar primero;
- cómo recorrer el corredor de engine;
- cómo correlacionar surface y truth artifacts;
- cómo distinguir seam reutilizable de implementación incidental.

### 3. Cómo introducir cambios sin violar invariantes

- qué cambios exigen versionado explícito;
- qué cambios rompen el baseline reusable;
- cómo validar que no se contradiga E.1/E.2;
- cómo distinguir mejora local vs ruptura de seam reusable.

### 4. Cómo usar el starter kit de capabilities

- cómo partir de templates de manifest y binding;
- cómo declarar `provider_ref` y compatibilidades;
- cómo declarar evidence refs mínimas;
- cómo validar la integración de governance y surface.

### Regla fuerte del playbook de developer

El playbook de developer no puede asumir:

- distribution layer;
- rollout tenant-scoped;
- bypass de Cerbos;
- apply directo desde surface;
- observabilidad derivada como reemplazo de truth artifacts.

## Playbooks de debugging y mantenimiento

Los playbooks de debugging y mantenimiento deben cubrir al menos:

### 1. Cómo reconstruir un corredor por correlación

- partir de IDs mínimos;
- seguir source refs, contract refs, execution refs, policy refs y event refs;
- distinguir gaps documentales de gaps materiales;
- detectar cuándo un caso no puede reconstruirse con evidencia suficiente.

### 2. Cómo detectar drift entre evidencia canónica y observabilidad derivada

- identificar cuándo PostgreSQL/event log/artifacts dicen una cosa y traces/dashboards otra;
- priorizar truth artifacts;
- escalar inconsistencia documental si la doc sugiere lo contrario.

### 3. Cómo candidatear mantenimiento gobernado

- distinguir `maintenance_action_candidate` de acción real;
- exigir `governance_requirements[]`, `target_refs[]`, `preconditions_refs[]` y `evidence_refs[]`;
- bloquear maintenance si falta evidencia o se sale del boundary.

### 4. Cómo distinguir un problema de docs vs un problema de engine/surface

- docs rotas: artifacts normativos sanos pero explicación incoherente o inconsistente;
- engine/surface rotos: evidence trail, correlación o invariantes materiales fallan;
- si una doc contradice artifacts normativos, el baseline documental está roto aunque el engine no lo esté.

### Regla fuerte del playbook de debugging/mantenimiento

El playbook de maintenance NO puede:

- sugerir bypass de governance;
- tratar candidate como acción real;
- usar observability derivada como única fuente;
- reintroducir distribution layer como maintenance operativo.

## Reglas de navegación y descubribilidad documental

La capa documental final debe fijar reglas explícitas de navegación:

1. debe indicar qué leer primero, segundo y tercero;
2. debe exponer enlaces o referencias claras entre mapas y playbooks;
3. debe usar términos consistentes para cada artifact;
4. debe evitar nombres alternativos para el mismo objeto salvo alias explícitos y controlados;
5. debe poder recorrerse por rol y también por corredor end-to-end;
6. debe incluir una ruta explícita “leer el sistema desde cero” para humanos e IA;
7. debe declarar de forma visible qué queda fuera del roadmap actual.

### Regla de descubribilidad mínima

Si una IA o una persona no pueden inferir qué leer primero para reconstruir el sistema, la documentación NO es IA-first de cierre.

## Reglas de consistencia entre docs y truth artifacts

Las docs finales deben obedecer estas reglas:

1. los artifacts normativos mandan sobre cualquier explicación documental;
2. la documentación no puede contradecir contracts, seams, invariantes ni boundaries normativos;
3. si una explicación documental contradice los artifacts normativos, el baseline documental está roto;
4. un playbook no puede introducir acciones o flows que no existan en el baseline reusable;
5. starter kit y docs finales deben ser coherentes con E.2;
6. D.5 sigue mandando sobre límites de maintenance asistido y action candidates;
7. distribution layer debe figurar como fuera de scope donde corresponda y nunca como dependencia implícita.

## Artefactos mínimos/documentos mínimos que deben existir al cerrar E.3

Al cerrar E.3 deben existir, como mínimo, estos documentos o artifacts documentales:

1. mapa del sistema final;
2. mapa final de artifacts canónicos;
3. mapa final de seams/flows;
4. mapa final de invariantes/fail-safes;
5. mapa final de evidencias/queries operativas;
6. playbook final de operator;
7. playbook final de developer;
8. playbook final de debugging/mantenimiento;
9. guía explícita de lectura desde cero para humans + IA;
10. referencia explícita de qué queda fuera del roadmap actual.

## Gaps aceptables al cerrar E.3

Son aceptables al cerrar E.3 sólo gaps que no rompan navegabilidad ni consistencia normativa, por ejemplo:

- mejoras futuras de redacción o pedagogía;
- ejemplos adicionales de casos por rol;
- índices auxiliares o atajos de navegación extra;
- mayor profundidad de casos no críticos ya cubiertos conceptualmente;
- materiales de onboarding complementario fuera del baseline mínimo.

## Gaps NO aceptables al cerrar E.3

No son aceptables al cerrar E.3:

- no distinguir los cinco mapas documentales principales;
- no tener playbooks diferenciados por rol;
- docs que contradigan seams o contracts normativos;
- docs que no indiquen qué leer primero;
- docs que no permitan reconstruir un caso end-to-end;
- docs que no mencionen correlación mínima de IDs;
- docs que omitan qué queda fuera del roadmap;
- docs que contradigan el starter kit o baseline reusable de E.2;
- docs que describan apply directo desde chat o bypass de governance;
- docs que usen observability derivada como única verdad.

## Tests borde mínimos (al menos 18)

La definición mínima de E.3 debe incluir, como piso, estos tests borde:

1. **doc final contradice seam normativo** => baseline documental roto.
2. **mapa de artifacts omite `execution_record`** => mapa canónico incompleto.
3. **mapa de seams no explica PEP/PDP** => falta crítica de arquitectura documental.
4. **playbook de operator recomienda acción prohibida** => playbook inválido.
5. **playbook de developer asume distribution layer** => contradicción de scope.
6. **playbook de maintenance sugiere bypass de governance** => contradicción con D.5.
7. **doc IA-first no indica qué leer primero** => falla de descubribilidad mínima.
8. **doc usa términos diferentes para el mismo artifact** => inconsistencia semántica documental.
9. **doc describe apply directo desde chat** => contradicción dura del baseline.
10. **doc de debugging usa observability derivada como única verdad** => documentación inválida.
11. **no hay forma de reconstruir un caso end-to-end desde docs** => E.3 no puede cerrarse.
12. **docs no mencionan correlación mínima de IDs** => capa de evidencias incompleta.
13. **docs no aclaran qué está fuera del roadmap** => mapa del sistema incompleto.
14. **starter kit explicado distinto a E.2** => inconsistencia entre baseline reusable y docs.
15. **operator playbook no cubre `unknown_outcome`** => playbook incompleto.
16. **developer playbook no cubre `provider_ref` / `binding`** => playbook incompleto.
17. **maintenance playbook no distingue candidate vs acción real** => contradicción con D.5.
18. **la IA no podría reconstruir boundaries leyendo solo docs finales** => documentación no IA-first.
19. **mapa de evidencias no aclara que observabilidad derivada no reemplaza truth artifacts** => docs rotas.
20. **playbook final describe un flujo de rollout tenant-scoped** => contradicción fuera de scope.

## Criterios de aceptación de E.3

E.3 puede considerarse cerrado cuando:

1. queda explicitado que “IA-first docs” significa documentación navegable y reconstruible por humans + IA, no texto marketing-friendly;
2. queda separada la documentación final en mapa del sistema, mapa de artifacts canónicos, mapa de seams/flows, mapa de invariantes/fail-safes y mapa de evidencias/queries operativas;
3. queda definido qué información mínima cubre cada capa documental;
4. quedan definidos playbooks diferenciados para operator, developer y debugging/mantenimiento;
5. la documentación final enseña el modelo mental correcto y no sólo lista archivos;
6. queda explícito que la documentación no puede contradecir contracts/seams normativos;
7. queda explícito que si una explicación documental contradice los artifacts normativos, el baseline documental está roto;
8. queda explícitamente fuera del mapa documental actual todo lo relativo a distribution layer;
9. existe una regla explícita para cómo una IA debería leer el sistema desde cero usando los docs finales;
10. quedan definidos artifacts mínimos, reglas de navegación, reglas de consistencia y tests borde suficientes como para tratar E.3 como cierre implementable de la capa documental final.
