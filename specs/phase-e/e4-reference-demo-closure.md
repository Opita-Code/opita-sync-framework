# E.4 — Reference demo closure

## Objetivo

E.4 existe para fijar un **demo de referencia probatorio, acotado y auditable** que demuestre el baseline reusable de Opyta Sync ya cerrado en engine + surface, sin inflar scope, sin convertir el demo en marketing de producto y sin introducir breadth funcional ajena al boundary actual.

El objetivo no es mostrar todo lo que el sistema podría llegar a hacer, sino dejar un corredor demostrable que permita probar con artifacts y evidencia que el baseline reusable realmente puede recorrerse, explicarse y reconstruirse de punta a punta.

## Principios del demo de referencia

1. **Probatorio y acotado.** El demo de referencia es probatorio y acotado, no breadth demo ni marketing demo.
2. **Artifacts primero.** Si el demo depende de explicación verbal sin artifacts/evidence suficientes, no sirve como cierre normativo.
3. **Baseline reusable, no producto completo.** Debe demostrar engine + surface gobernada, no distribution layer ni rollout downstream.
4. **Corredor mínimo explícito.** Debe existir un corredor mínimo demostrable y reconstruible, no una secuencia ambigua de pantallas o relato verbal.
5. **Puede cerrar sin apply real.** El demo puede cerrar sin apply real si deja un evidence trail suficiente y consistente con D.6 y E.1.
6. **Debe declarar límites.** El demo debe dejar claro qué valida y qué NO valida.
7. **No puede mostrar flows fuera de scope.** Si muestra distribution layer, rollout o activación tenant-scoped, el demo queda mal definido.
8. **Representatividad por integridad, no por espectacularidad.** Vale más un corredor sobrio y completo que un demo vistoso pero no auditable.

## Qué significa “demo de referencia” en este proyecto

En este proyecto, “demo de referencia” significa una demostración controlada que permite a un tercero verificar que:

- el baseline reusable de engine + surface es real y no meramente conceptual;
- los seams cerrados pueden recorrerse de forma coherente;
- el caso queda auditado por artifacts y evidence trail suficientes;
- la surface puede cerrar un caso demostrable sin saltar governance ni exigir apply real;
- el sistema puede explicarse como baseline reusable sin prometer breadth de producto fuera del roadmap.

No significa:

- una demo comercial del producto completo;
- una demo que tape huecos con narración oral;
- una demo de distribution layer;
- una demo de provisioning/activation/rollout tenant-scoped;
- una demo donde el valor proviene de UX o storytelling pero no de evidencia.

## Alcance exacto del demo

El demo de referencia debe quedar acotado al boundary actual de engine + surface gobernada y cubrir, como mínimo:

- intake gobernado;
- intent/proposal;
- preview/simulación;
- paso por policy/runtime/resolution/event log en el corredor correspondiente;
- inspección o reconstrucción final del caso.

### Regla de scope

El demo NO debe:

- depender de distribution layer;
- exigir rollout o activation tenant-scoped;
- sugerir que existe apply directo desde chat;
- mostrar maintenance como ejecución directa;
- inflar el alcance con flows no cerrados normativamente.

## Corredor mínimo que debe demostrar

El corredor mínimo obligatorio del demo debe poder demostrarse así:

1. **Entrada gobernada**
   - existe `conversation_turn`, `intent_input` o artifact equivalente según el corredor elegido;
   - la entrada está correlacionada y no se sostiene sólo en explicación verbal.

2. **Intent/proposal**
   - existe `proposal_draft` gobernado;
   - el draft tiene referencias suficientes al origen y al cambio candidateado.

3. **Preview/simulación**
   - existe `preview_candidate` correlado con el draft;
   - existe al menos un `simulation_result` con `inputs_refs[]` y evidencia mínima.

4. **Paso por engine**
   - existe evidencia de `compiled_contract` cuando aplique;
   - existe evidencia de `policy_decision_record` o equivalente normativa;
   - existe evidencia de `execution_record`, `event_record` y/o resolution chain cuando aplique al corredor.

5. **Cierre de lectura final**
   - existe `execution_inspection_view` o `semantic_debug_view` final;
   - el caso puede reconstruirse de punta a punta por correlación mínima.

### Regla de cierre

El demo puede cerrar sin apply real, siempre que:

- el corredor quede íntegro;
- la evidencia sea suficiente;
- los artifacts sean persistibles o normativamente reproducibles;
- la reconstrucción del caso no dependa de memoria humana o relato oral.

## Qué prueba el demo

El demo de referencia debe probar, como mínimo:

1. que el baseline reusable de engine + surface puede recorrerse end-to-end dentro del boundary actual;
2. que intake, proposal, preview, policy, runtime, resolution, event log e inspection/debug pueden sostener un caso coherente;
3. que existe evidence trail suficiente para reconstruir el caso;
4. que el sistema puede explicarse como baseline reusable sin breadth extra de producto;
5. que la demostración no depende de observabilidad derivada como única verdad;
6. que el corredor puede entenderse sin exigir apply real.

## Qué NO prueba el demo

El demo NO prueba:

- distribution layer;
- rollout, activación o bootstrap tenant-scoped;
- tenant templates productizados;
- breadth de todas las capabilities futuras;
- performance a escala;
- madurez comercial total del producto;
- automatización operativa fuera del boundary engine/surface.

### Regla de honestidad del demo

El demo debe declarar de forma explícita qué no valida. Si oculta sus límites y parece prometer producto completo, queda mal definido.

## Artefactos mínimos que debe producir o exhibir el demo

El demo debe exhibir, como mínimo:

1. artifact de entrada gobernada (`conversation_turn` / `intent_input` o equivalente según el corredor);
2. `proposal_draft`;
3. `preview_candidate`;
4. al menos un `simulation_result`;
5. evidencia de `compiled_contract` cuando aplique;
6. evidencia de `policy_decision_record` o equivalente;
7. evidencia de `execution_record` / `event_record` / resolution chain cuando aplique;
8. una `execution_inspection_view` o `semantic_debug_view` final;
9. una explicación explícita de límites del demo.

### Regla documental del demo

Estos artifacts deben ser visibles, referenciables o reproducibles normativamente. No alcanza con decir que “existen detrás”.

## Evidence trail mínimo del demo

El demo debe dejar un evidence trail mínimo que permita reconstrucción independiente del relato del presentador.

Debe incluir, como mínimo:

- IDs correlados del corredor (`tenant_id`, `trace_id`, `proposal_draft_id`, `preview_candidate_id`, `simulation_result_id`, `contract_id`, `execution_id`, `event_id` cuando apliquen);
- evidence refs que conecten surface con engine;
- evidencia de policy/governance cuando el demo diga que está mostrando enforcement;
- evidencia de runtime/event log si el demo afirma que está mostrando paso por engine;
- evidencia de inspection/debug final alineada con los truth artifacts;
- explicación explícita de cuál es el límite del corredor y por qué el demo sigue siendo representativo.

### Regla dura de evidence trail

Si el demo no deja evidencia suficiente para reconstrucción end-to-end, NO sirve como cierre probatorio de E.4.

## Criterios de representatividad del demo

El demo será representativo sólo si:

1. recorre un corredor real del baseline reusable y no un mock narrativo;
2. usa artifacts normativos o evidence refs válidas;
3. muestra tanto surface como engine en los puntos necesarios del corredor;
4. deja clara la relación entre proposal, preview, simulation, policy, runtime y evidence;
5. puede ser entendido por operator/developer sin inferencias heroicas;
6. no necesita breadth funcional adicional para demostrar el valor del baseline reusable;
7. puede explicarse como referencia reusable, no como producto completo.

## Criterios de exclusión (qué haría al demo inválido o engañoso)

El demo queda inválido o engañoso si ocurre cualquiera de estas condiciones:

- muestra chat pero no artifacts gobernados;
- muestra proposal pero no preview;
- muestra preview sin `simulation_result`;
- usa observability derivada como única evidencia;
- no puede reconstruir correlación mínima;
- dice probar governance pero no muestra policy evidence;
- dice probar runtime pero no muestra `execution_record`/`event_record` o evidencia equivalente;
- depende de apply real para que el caso sea comprensible;
- usa un caso fuera de scope actual;
- oculta límites y parece prometer producto completo;
- no distingue qué valida y qué no valida;
- no muestra inspection/debug final;
- contradice el baseline reusable fijado en E.2;
- requiere distribution layer implícito.

## Gaps aceptables al cerrar E.4

Son aceptables al cerrar E.4 sólo gaps que no rompan el carácter probatorio del demo, por ejemplo:

- mejoras futuras de presentación visual;
- variantes adicionales del mismo corredor para audiencias distintas;
- documentación complementaria de apoyo a la demo;
- material pedagógico adicional que no cambie artifacts ni evidencia;
- refinamientos narrativos mientras el demo ya sea auditable y no engañoso.

## Gaps NO aceptables al cerrar E.4

No son aceptables al cerrar E.4:

- demo sin corredor mínimo explícito;
- demo sin artifacts mínimos obligatorios;
- demo sin evidence trail reconstruible;
- demo que dependa de apply real como requisito normativo;
- demo que use flows fuera de scope;
- demo que prometa breadth de producto o distribution layer;
- demo que contradiga E.1, D.6, E.2 o `specs/_status.md`;
- demo que no explicite qué valida y qué no valida;
- demo marketing-friendly pero no auditable.

## Tests borde mínimos (al menos 18)

La definición mínima de E.4 debe incluir, como piso, estos tests borde:

1. **demo muestra chat pero no artifacts gobernados** => demo inválido.
2. **demo muestra proposal pero no preview** => corredor incompleto.
3. **demo muestra preview sin simulation_result** => demo no probatorio.
4. **demo usa observability derivada como única evidencia** => invalida el cierre.
5. **demo no puede reconstruir correlación mínima** => falla estructural del demo.
6. **demo dice probar governance pero no muestra policy evidence** => demo engañoso.
7. **demo dice probar runtime pero no muestra `execution_record`/`event_record`** => demo incompleto.
8. **demo depende de apply real para ser comprensible** => contradicción con el scope de E.4.
9. **demo usa caso fuera de scope actual** => demo mal definido.
10. **demo oculta límites y parece prometer producto completo** => demo engañoso.
11. **demo no distingue qué valida vs qué no valida** => demo no aceptable.
12. **demo no muestra inspection/debug final** => corredor no cerrable.
13. **demo tiene simulation_result pero sin inputs refs** => evidencia insuficiente.
14. **demo muestra maintenance como ejecución directa** => contradicción con D.5.
15. **demo contradice baseline reusable de E.2** => demo inválido.
16. **demo requiere distribution layer implícito** => fuera de scope.
17. **demo no deja evidence trail persistible** => no cumple carácter probatorio.
18. **demo marketing-friendly pero no auditable** => no sirve como reference demo.
19. **demo muestra policy decision pero sin mapping a runtime/estado** => evidencia incompleta.
20. **demo muestra semantic debug final que contradice artifacts canónicos** => demo roto aunque la narración sea buena.

## Criterios de aceptación de E.4

E.4 puede considerarse cerrado cuando:

1. queda explícito que el demo de referencia es **probatorio y acotado**, no breadth demo ni marketing demo;
2. queda explícito que el demo demuestra el baseline reusable de engine + surface, no distribution layer;
3. queda fijado un corredor mínimo demostrable que incluye intake gobernado, intent/proposal, preview/simulación, paso por policy/runtime/resolution/event log en el corredor correspondiente e inspección/reconstrucción final del caso;
4. queda explícito que el demo puede cerrar sin apply real;
5. quedan definidos los artifacts mínimos obligatorios del demo;
6. queda definido el evidence trail mínimo necesario para reconstrucción end-to-end;
7. queda claro qué prueba y qué NO prueba el demo;
8. quedan definidos criterios de representatividad y criterios de exclusión del demo;
9. queda explícito que si el demo depende de explicación verbal sin artifacts/evidence suficientes, no sirve;
10. quedan definidos tests borde suficientes como para tratar E.4 como cierre implementable del demo de referencia.
