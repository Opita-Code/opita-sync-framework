# E.2 — Reusable engine baseline and starter kit

## Objetivo

E.2 existe para fijar de manera **implementable, portable y auditable** qué parte de Opyta Sync constituye el **baseline reusable real** para proyectos futuros y qué contenido mínimo debe ofrecer el starter kit técnico para introducir capabilities nuevas sin reabrir el motor ni la surface gobernada.

Este bloque no productiza distribución, rollout ni activación tenant-scoped. Su función es cerrar la base reusable del **engine + surface gobernada + contracts + documentación/playbooks mínimos** dentro del boundary ya estabilizado en C.5, D.5, E.0 y E.1.

## Principios del baseline reusable

1. **Reusable baseline no equivale a producto completo.** En el scope actual, reusable significa reutilizable como base técnica y operativa, no listo para consumo downstream masivo.
2. **Engine + surface gobernada forman una sola base reusable.** No alcanza con el kernel aislado; también deben ser reutilizables los flows gobernados de proposal, preview, inspection y maintenance.
3. **Contracts primero.** Lo reusable no es sólo un conjunto de archivos: es un conjunto de seams, artifacts y contratos normativos que deben mantenerse estables o versionarse explícitamente.
4. **Starter kit separado del baseline normativo.** Debe distinguirse con claridad entre artifacts normativos del baseline, starter kit para crear capabilities nuevas y artifacts/documentación de referencia.
5. **Reutilizable no significa libre.** Una capability nueva sólo es compatible si respeta envelopes, contracts, governance, evidence trail, registry/resolution y boundaries de surface ya cerrados.
6. **Sin distribution layer implícita.** El baseline reusable no puede depender de channels, rollout, activation o templates tenant-scoped productizados.
7. **Observabilidad derivada no reemplaza portabilidad real.** Un proyecto futuro debe poder operar el baseline reusable desde artifacts normativos y evidencia canónica, no desde dashboards o traces solamente.
8. **Documentación mínima obligatoria.** Un baseline reusable sin docs, checklists y playbooks mínimos no está realmente listo para ser reutilizado.

## Qué cuenta como baseline reusable real en el scope actual

En el scope actual, cuenta como baseline reusable real el conjunto coherente de:

- engine cerrado con compilación, runtime, policy, event log y registry/resolution compatibles entre sí;
- surface gobernada cerrada con intake, proposal, preview, inspection y maintenance asistido dentro de boundaries explícitos;
- contracts normativos que fijan envelopes, schemas, fingerprints, correlación y evidence trail mínimo;
- documentación y playbooks mínimos que permiten a un proyecto futuro adoptar ese baseline sin reinterpretar su arquitectura.

Por lo tanto, **reusable baseline = engine + surface gobernada + contracts + docs/playbooks mínimos**, no distribution layer.

### Distinción obligatoria dentro de E.2

E.2 debe dejar separadas estas tres capas:

1. **Artifacts normativos del baseline**
   - son la base reusable obligatoria;
   - definen la semántica operativa del sistema.

2. **Starter kit para crear capabilities nuevas**
   - es un kit técnico de arranque;
   - no es baseline por sí mismo, pero debe ser compatible con él.

3. **Artifacts/documentación de referencia**
   - explican cómo usar, extender y validar el baseline;
   - no reemplazan los contracts normativos.

## Qué artifacts integran el baseline reusable

Los artifacts que integran el baseline reusable deben quedar agrupados al menos así:

### 1. Artifacts normativos del dominio y del motor

- envelopes canónicos (`api_version`, `kind`, `metadata`, `spec`, `status`);
- versionado de contratos, fingerprints y reason codes;
- `intent_input`, `compiled_contract`, `compilation_report`;
- `execution_record`, estados canónicos y records de outcome;
- `policy_input_v1`, `policy_decision_record` o equivalente normativa persistida;
- `event_record` canónico y evidence trail mínimo;
- artifacts de resolution (`capability_manifest`, `bundle_digest`, `binding`, `provider_ref`).

### 2. Artifacts normativos de surface gobernada

- `conversation_turn` e `intake_session` donde aplique;
- `intent_candidate`, `proposal_draft`, `governed_patchset_candidate`;
- `preview_candidate`, `simulation_result`;
- `execution_inspection_view` y `semantic_debug_view` como vistas derivadas normativamente acotadas;
- `maintenance_action_candidate` dentro de los límites de D.5.

### 3. Artifacts de baseline operativo reusable

- baseline de regresión y hardening de E.1;
- checklist operativo de Fase E para continuidad de cierre;
- definiciones de seams reutilizables y guardrails de extensibilidad;
- documentación mínima para reconstrucción del corredor engine + surface.

### 4. Artifacts de referencia y apoyo

- starter kit de capabilities;
- templates normativos mínimos de manifest y binding;
- guías de packaging, evidence refs y provider compatibility;
- documentación/playbooks mínimos de adopción y extensión.

## Qué seams/contracts se consideran reutilizables

Los siguientes seams/contracts quedan explícitamente declarados como reutilizables en el baseline de E.2:

1. **seam `intent -> compiled_contract`**
   - reusable para proyectos que adopten la misma disciplina de intención gobernada y compilación determinística.

2. **seam `compiled_contract -> execution_record`**
   - reusable para proyectos que necesiten mantener correlación fuerte entre contrato compilado y runtime durable.

3. **seam PEP/PDP con input canonizado**
   - reusable para autorización contextual y governance sin lógica ad hoc fuera del policy boundary.

4. **seam event log canónico / observabilidad derivada**
   - reusable para proyectos que necesiten evidencia operativa primaria desacoplada de traces y dashboards.

5. **seam `capability_manifest -> bundle_digest -> binding -> provider_ref`**
   - reusable como modelo estable de packaging + resolution sin introducir distribution layer.

6. **seam proposal/preview/inspection/maintenance gobernados**
   - reusable como pattern de surface operativa gobernada sobre artifacts y evidencia canónica.

### Regla de estabilidad de seams

- Un seam reutilizable no puede cambiar semánticamente sin versionado explícito.
- Si un proyecto futuro necesita romper un seam reutilizable, eso ya no cuenta como adopción del baseline, sino como derivación arquitectónica.
- La documentación de referencia no puede redefinir un seam por fuera de su contrato normativo.

## Qué NO forma parte del baseline reusable

Queda explícitamente fuera del baseline reusable:

- distribution layer;
- activation/rollout tenant-scoped;
- templates de tenant productizados;
- canales de entrega/consumo downstream;
- publication workflows, sync, pull/push o replication como parte de la base reusable;
- cualquier capability que exija bypass de Cerbos, bypass de governance o apply directo desde surface;
- cualquier forma de maintenance que reintroduzca ejecución implícita o autoridad no gobernada.

### Regla dura de scope

Si un artifact, template o playbook implica distribuir, activar, desplegar o consumir fuera del boundary engine/surface ya cerrado, NO entra en E.2.

## Starter kit mínimo para capabilities

El starter kit mínimo para capabilities nuevas debe ser un **kit técnico de arranque** y no un paquete productizado de rollout.

Debe incluir como mínimo:

- template de `capability_manifest`;
- template de `binding`;
- guía de `bundle_digest` / packaging refs;
- guía de `provider_ref` / provider compatibility;
- checklist de governance/policy/classification expectations;
- checklist de observabilidad y evidence refs;
- checklist de integración con proposal/preview/inspection;
- checklist de regresión mínima para capability nueva.

### Naturaleza del starter kit

- Sirve para crear capabilities nuevas compatibles con el baseline.
- No instala tenants.
- No resuelve distribution.
- No activa rollout.
- No reemplaza los contracts normativos del motor.

## Estructura mínima esperada del starter kit

El starter kit debe poder organizarse de forma equivalente, pero como mínimo debe cubrir estas piezas lógicas:

1. **Plantillas normativas**
   - template de `capability_manifest` compatible con envelope canónico;
   - template de `binding` compatible con runtime baseline.

2. **Guías de compatibilidad**
   - guía de `contract_schema_version` soportadas;
   - guía de `supported_result_types`;
   - guía de `provider_ref` y compatibilidad de runtime/provider.

3. **Guías de packaging y resolución**
   - guía de `bundle_digest` y evidencia mínima de packaging;
   - guía de attachments, signature/provenance refs cuando apliquen;
   - guía de cómo conectar manifest, binding y provider_ref sin saltos implícitos.

4. **Checklists de governance y evidence**
   - checklist de policy/approval/classification;
   - checklist de evidence refs mínimas;
   - checklist de correlación mínima con runtime y event log.

5. **Checklists de surface integration**
   - hooks mínimos para proposal;
   - hooks mínimos para preview/simulation;
   - hooks mínimos para inspection/debug/maintenance donde aplique.

6. **Checklists de regresión de capability nueva**
   - casos mínimos de compatibilidad de manifest;
   - casos mínimos de binding/provider;
   - casos mínimos de evidence trail;
   - casos mínimos de governance y surface integration.

## Contratos mínimos que debe respetar una capability nueva

Toda capability nueva que pretenda ser compatible con el baseline reusable debe respetar, como mínimo, estos contratos:

1. **Envelope canónico**
   - el manifest y artifacts normativos asociados deben respetar `api_version`, `kind`, `metadata`, `spec`, `status`.

2. **Compatibilidad de schema**
   - debe declarar `contract_schema_version` compatible;
   - no puede romper versiones soportadas sin versionado explícito.

3. **Declaración de tipos de resultado**
   - debe declarar `supported_result_types` compatibles con A.3;
   - no puede asumir outcomes implícitos no modelados.

4. **Compatibilidad con registry/resolution**
   - debe poder resolverse por la cadena `manifest -> bundle_digest -> binding -> provider_ref`;
   - no puede depender de resolución best-effort o heurística.

5. **Compatibilidad con governance**
   - debe declarar expectations de policy/approval/classification;
   - no puede requerir bypass de Cerbos ni saltar approvals normativos.

6. **Compatibilidad con evidence trail**
   - debe emitir evidence refs mínimas suficientes;
   - no puede depender sólo de observabilidad derivada para sostener su operación.

7. **Compatibilidad con surface gobernada**
   - debe integrarse con proposal/preview/inspection/maintenance cuando corresponda;
   - no puede asumir apply directo desde surface.

### Regla de compatibilidad mínima

Una capability nueva no es “compatible con el baseline” sólo porque compile o tenga provider. Es compatible sólo si respeta contracts, governance, evidence trail y surface integration dentro del boundary reusable.

## Guardrails de extensibilidad reutilizable

La extensibilidad reutilizable de E.2 debe quedar protegida por estos guardrails:

1. no reinterpretar manifests ni bindings con reglas implícitas no versionadas;
2. no aceptar capabilities que rompan `contract_schema_version` sin estrategia explícita de compatibilidad;
3. no aceptar capabilities que omitan `supported_result_types`;
4. no aceptar capabilities sin evidence refs mínimas;
5. no aceptar provider compatibility por parecido nominal: debe ser material y verificable;
6. no aceptar integración de surface que salte proposal/preview/inspection donde sea requerido;
7. no aceptar starter kits que incluyan pasos de distribution layer, rollout o tenant activation;
8. no aceptar documentación de referencia que contradiga contracts normativos;
9. no aceptar cambios en seams reutilizables sin versionado explícito;
10. no aceptar baselines que sólo operen si observabilidad derivada está disponible.

## Artefactos/documentación obligatoria que debe incluir el baseline reusable

El baseline reusable de E.2 debe incluir, como mínimo:

1. inventario explícito de artifacts normativos del baseline;
2. inventario explícito de seams/contracts reutilizables;
3. definición del starter kit y su estructura mínima;
4. templates normativos mínimos de manifest y binding;
5. guía de packaging refs y `bundle_digest` evidence;
6. guía de `provider_ref` y compatibility rules;
7. checklists de governance, classification y evidence refs;
8. checklists de integration con proposal/preview/inspection/maintenance;
9. criterios de compatibilidad hacia proyectos futuros;
10. lista de elementos fuera de scope para evitar desborde hacia distribution layer.

## Criterios de compatibilidad hacia proyectos futuros

Un proyecto futuro puede considerarse compatible con el baseline reusable sólo si:

- puede adoptar los artifacts normativos sin redefinir sus envelopes y contracts centrales;
- puede recorrer el seam `intent -> compiled_contract -> execution_record` sin reinterpretación ad hoc;
- puede sostener policy/governance vía input canonizado y persistencia normativa suficiente;
- puede operar con event log canónico sin depender de observabilidad derivada para demostrar verdad;
- puede introducir capabilities nuevas a través del starter kit respetando manifest/binding/provider compatibility;
- puede integrar surface gobernada sin convertir proposal/preview/inspection/maintenance en flows libres;
- puede reconstruir correlación mínima entre artifacts, runtime y evidence trail.

### Incompatibilidades duras hacia proyectos futuros

Un proyecto futuro NO es compatible si:

- requiere distribution layer para arrancar el baseline;
- requiere apply directo desde surface;
- requiere bypass de Cerbos o governance equivalente;
- no puede operar sin traces/dashboards como fuente primaria;
- no puede reconstruir correlación mínima con artifacts del starter kit y del baseline normativo.

## Gaps aceptables al cerrar E.2

Son aceptables al cerrar E.2 sólo gaps que no rompan la reutilización real del baseline, por ejemplo:

- examples adicionales de capabilities no críticas;
- refinamientos de estilo en templates o naming auxiliar;
- documentación complementaria que profundice adopción pero no cambie contratos;
- variantes futuras de starter kit para casos especializados fuera del baseline mínimo;
- mejoras de ergonomía en playbooks de referencia que no alteren semántica normativa.

## Gaps NO aceptables al cerrar E.2

No son aceptables al cerrar E.2:

- no distinguir entre baseline normativo, starter kit y documentación de referencia;
- no explicitar qué seams/contracts son reutilizables;
- no explicitar qué queda fuera del baseline reusable;
- starter kit sin template de manifest o binding;
- starter kit que incluya pasos de distribution layer o rollout;
- capabilities nuevas compatibles “de palabra” pero sin contracts mínimos definidos;
- dependencia de observabilidad derivada para operar el baseline reusable;
- falta de criterios de compatibilidad hacia proyectos futuros;
- documentación de referencia contradiciendo artifacts normativos;
- seams reutilizables que puedan cambiar sin versionado explícito.

## Tests borde mínimos (al menos 18)

La definición mínima de E.2 debe incluir, como piso, estos tests borde:

1. **proyecto futuro intenta reutilizar baseline sin event log canónico** => debe fallar por incompatibilidad estructural.
2. **capability nueva no declara `supported_result_types`** => no puede considerarse compatible.
3. **binding template omitido en starter kit** => starter kit incompleto, E.2 no cierra sano.
4. **provider_ref fuera de compatibilidad** => binding/resolution deben rechazar la capability.
5. **capability nueva rompe `contract_schema_version`** => requiere versionado explícito o queda fuera del baseline.
6. **capability nueva no emite evidence refs mínimas** => incompatibilidad dura con baseline reusable.
7. **proposal/preview hooks ausentes en capability nueva** => integración de surface insuficiente donde aplique.
8. **starter kit incluye paso de distribution layer** => contradicción de scope, debe rechazarse.
9. **documentación de reference contradice contracts normativos** => la documentación no puede prevalecer sobre el contrato.
10. **seam reutilizable cambia sin versionado explícito** => ruptura del baseline reusable.
11. **baseline reusable depende de observability derivada para operar** => incompatibilidad dura.
12. **capability nueva requiere bypass de Cerbos** => debe rechazarse.
13. **capability nueva asume apply directo desde surface** => debe rechazarse por contradicción con D.5.
14. **missing bundle digest evidence** => packaging/resolution no puede considerarse sano.
15. **template de manifest no respeta envelope canónico** => starter kit inválido.
16. **binding template incompatible con runtime baseline** => incompatibilidad estructural del starter kit.
17. **playbook final describe flujo fuera de scope actual** => debe marcarse como documentación inválida para el baseline.
18. **proyecto futuro no puede reconstruir correlación mínima con artifacts del starter kit** => baseline no reutilizable de forma real.
19. **starter kit omite checklist de governance/policy/classification** => baseline extensible de forma insegura.
20. **capability nueva sólo ofrece evidencia en traces y no en artifacts canónicos** => incompatibilidad con baseline reusable.

## Criterios de aceptación de E.2

E.2 puede considerarse cerrado cuando:

1. queda definido qué cuenta como baseline reusable real dentro del scope actual;
2. queda explícito que reusable baseline = **engine + surface gobernada + contracts + docs/playbooks mínimos**, no distribution layer;
3. quedan diferenciados artifacts normativos del baseline, starter kit para capabilities nuevas y artifacts/documentación de referencia;
4. quedan enumerados los artifacts que integran el baseline reusable;
5. quedan fijados los seams/contracts reutilizables mínimos, incluyendo `intent -> compiled_contract`, `compiled_contract -> execution_record`, PEP/PDP canonizado, event log/observabilidad derivada, `capability_manifest -> bundle_digest -> binding -> provider_ref` y proposal/preview/inspection/maintenance gobernados;
6. queda explícito qué NO forma parte del baseline reusable: distribution layer, activation/rollout tenant-scoped, templates de tenant productizados y canales de entrega/consumo;
7. queda definido el starter kit mínimo para capabilities como kit técnico de arranque y no como paquete productizado de rollout;
8. queda definida la estructura mínima del starter kit y los contratos mínimos que una capability nueva debe respetar;
9. quedan definidos guardrails de extensibilidad reutilizable y criterios de compatibilidad hacia proyectos futuros;
10. quedan definidos tests borde suficientes como para tratar E.2 como baseline implementable de cierre reusable.
