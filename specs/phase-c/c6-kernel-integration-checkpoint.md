 # C.6 — Kernel integration checkpoint

 ## Objetivo

 C.6 cierra el checkpoint de integración de **Fase C** para demostrar que el **kernel engine-only** ya puede recorrer un corredor mínimo end-to-end sin reabrir decisiones de C.1-C.5 y sin extender alcance a distribution layer, surface conversacional o UX operativa. Este checkpoint no agrega breadth de producto: valida que compilación, runtime durable, enforcement de policy, evidencia canónica y capability resolution ya nacen compatibles entre sí a nivel de construcción v1.

 ## Qué valida exactamente este checkpoint

 Este checkpoint valida la integración del **motor**, no del distribution layer ni de la UX conversacional.

 Debe poder probarse, como mínimo, que el kernel:

 1. recibe un `intent_input` ya gobernado;
 2. produce y persiste un `compiled_contract` con `contract_fingerprint`;
 3. crea un `execution_record` durable y proyectado de forma canónica;
 4. consulta policy vía Cerbos con input canonizado y persistencia de decisión/evidencia;
 5. resuelve una capability por cadena completa `registry -> binding -> provider_ref`;
 6. persiste al menos un `event_record` canónico correlacionado;
 7. cierra el smoke path en `blocked`, `awaiting_approval` o `execution_completed`, según corresponda al caso, pero siempre dejando evidencia íntegra.

 C.6 NO valida:

 - distribution layer, channels, sync, pull/push o replication;
 - shells conversacionales avanzadas;
 - UX humana de approvals;
 - analytics o dashboards como criterio de verdad;
 - hardening final, performance tuning profundo o escalado fino de Fase E.

 ## Principios del checkpoint de integración

 - **Engine-only scope.** El checkpoint cubre el corredor del motor y sus seams internos ya definidos en C.1-C.5.
 - **Verdad canónica primero.** PostgreSQL y los records canónicos mandan; traces, logs y Temporal history son evidencia complementaria.
 - **No reinterpretación entre seams.** Cada seam consume la salida persistida y versionada del seam anterior.
 - **Fail-closed donde corresponde.** Si falta evidencia estructural, correlación mínima o resolución válida, el checkpoint no puede considerarse sano.
 - **Integración mínima, no completitud del producto.** Alcanzar el corredor end-to-end mínimo vale más que abrir casos laterales.
 - **Observabilidad derivada desacoplada.** Si falla exportación OTel, el checkpoint puede pasar sólo si la verdad operativa y el event log canónico quedan íntegros.
 - **Persistencia canónica obligatoria.** Si falla persistencia canónica, el checkpoint NO puede pasar aunque existan traces, logs o history de Temporal.

 ## Supuestos heredados de C.1-C.5

 C.6 hereda sin reabrir estas decisiones:

 - C.1 fija `intent_input -> compiled_contract` como compilación determinística y persistida en PostgreSQL.
 - C.1 fija `compiled_contract` + `compilation_report` y `contract_fingerprint` como salida formal del compilador.
 - C.2 fija `execution_record` como objeto first-class distinto de Temporal history.
 - C.2 fija un `execution_workflow` durable por `execution_id` y separación semántica execution/application.
 - C.3 fija al kernel como PEP y a Cerbos como PDP con `policy_input_v1` canonizado.
 - C.3 fija `policy_decision_record` o evidencia equivalente como persistencia obligatoria de policy.
 - C.4 fija `event_record` append-only en PostgreSQL como evidencia operativa canónica.
 - C.4 fija OTel/LGTM como observabilidad derivada, nunca como verdad operativa.
 - C.5 fija la resolución por cadena completa `capability_manifest -> bundle_digest -> binding -> provider_ref`.
 - C.5 deja explícito que distribution layer y activación tenant-scoped quedan fuera del roadmap actual.

 ## Corredor mínimo end-to-end del kernel

 El corredor mínimo obligatorio de C.6 queda definido así:

 1. **Entrada gobernada**
    - existe `intent_input` válido, tenant-scoped y trazable;
    - incluye `tenant_id`, `trace_id` o capacidad de generarlo, origen gobernado y metadata suficiente para compilación.

 2. **Compilación a contrato**
    - el compilador produce `compiled_contract` persistido;
    - existe `contract_id`, `contract_fingerprint`, `capability_id`, snapshots/versiones relevantes y `policy_evaluation_input` preparado.

 3. **Arranque durable de ejecución**
    - el runtime crea o reutiliza idempotentemente `execution_record`;
    - la proyección canónica deja `execution_id`, `contract_id`, `contract_fingerprint`, `tenant_id`, `trace_id` y estado inicial consistente.

 4. **Gate de policy**
    - el kernel consulta Cerbos mediante `policy_input_v1` canonizado;
    - la decisión queda persistida como `policy_decision_record` o evidencia equivalente normativa;
    - la salida mapea a `allow`, `blocked` o `awaiting_approval` sin lógica ad hoc.

 5. **Resolución de capability**
    - el runtime o su seam de ejecución resuelve `capability_id` mediante registry;
    - existe `bundle_digest` verificable;
    - existe `binding_id` resoluble y `provider_ref` compatible.

 6. **Evidencia canónica**
    - se persiste al menos un `event_record` canónico del corredor;
    - la evidencia correlaciona IDs materiales del trayecto.

 7. **Cierre mínimo del smoke path**
    - el caso termina en `blocked`, `awaiting_approval` o `execution_completed`;
    - en todos los casos queda evidence trail íntegro y consultable.

 ## Smoke path mínimo obligatorio

 El smoke path mínimo obligatorio de C.6 no exige éxito funcional externo; exige **integridad estructural del corredor**.

 Secuencia mínima:

 1. ingresar `intent_input` gobernado;
 2. compilar y persistir `compiled_contract`;
 3. crear/proyectar `execution_record`;
 4. canonizar `policy_input_v1` e invocar Cerbos;
 5. persistir `policy_decision_record` y/o evidencia normativa equivalente;
 6. resolver `capability_id` por registry + binding + `provider_ref`;
 7. persistir al menos un `event_record` canónico del corredor;
 8. cerrar el caso en uno de estos estados permitidos:
    - `blocked` con `reason_code`;
    - `awaiting_approval` con `approval_request_id`;
    - `execution_completed` con evidencia suficiente de cierre de fase de execution.

 Reglas normativas del smoke path:

 - puede terminar antes de application phase; C.6 valida motor mínimo, no breadth de efectos externos;
 - NO puede terminar en `unknown_outcome` para pasar el smoke path base, salvo caso borde explícito y con evidence trail completo;
 - debe dejar evidencia íntegra aunque la exportación derivada falle;
 - debe ser repetible sin duplicar canónicamente `execution_record` ni `event_record` materiales.

 ## Gates de consistencia entre seams

 Para que C.6 pase, los siguientes gates deben cumplirse:

 1. **Compiler -> runtime gate**
    - no existe arranque de ejecución sin `compiled_contract` persistido;
    - `execution_record.contract_id` y `execution_record.contract_fingerprint` deben coincidir con el contrato compilado.

 2. **Compiler -> policy gate**
    - el input a Cerbos se construye desde campos canonizados del contrato y del contexto, no desde texto libre;
    - el `capability_id` consultado por policy debe coincidir con el contrato compilado.

 3. **Policy -> runtime gate**
    - la decisión de Cerbos debe mapear a transición conocida del runtime;
    - respuestas no mapeables o inputs incompletos cierran en fail-closed.

 4. **Runtime -> registry gate**
    - el runtime sólo puede intentar ejecución si la capability es resoluble por cadena completa;
    - no alcanza con que exista `capability_id`; debe existir binding vigente y `provider_ref` compatible.

 5. **Runtime -> event log gate**
    - los hitos materiales del smoke path deben poder persistirse como `event_record` canónico;
    - si falla esta persistencia, el checkpoint falla.

 6. **Event log -> observability gate**
    - la exportación derivada puede fallar sin invalidar el checkpoint;
    - nunca puede existir dependencia inversa donde OTel complete semántica faltante del event log canónico.

 ## Gates de correlación y evidencia

 Debe poder demostrarse correlación completa entre:

 - `tenant_id`
 - `trace_id`
 - `contract_id`
 - `contract_fingerprint`
 - `execution_id`
 - `policy_decision_id`
 - `event_id`
 - `capability_id`
 - `bundle_digest`
 - `binding_id`
 - `provider_ref`

 Gates obligatorios:

 1. `tenant_id` debe sobrevivir sin ambigüedad desde intención hasta event log.
 2. `trace_id` debe correlacionar compilación, policy, runtime y evidencia; si no viene, debe generarse antes del primer hecho canónico.
 3. `contract_id` y `contract_fingerprint` deben existir en `execution_record`, `policy_decision_record` y eventos materiales del corredor.
 4. `execution_id` debe enlazar workflow durable, estado canónico y eventos runtime-bound.
 5. `policy_decision_id` debe enlazar la respuesta PDP con el evento y el estado que causó.
 6. `capability_id`, `bundle_digest`, `binding_id` y `provider_ref` deben permitir reconstruir la cadena exacta de resolución usada por el smoke path.
 7. `event_id` debe ser único, append-only y suficiente para ubicar el hecho canónico aunque fallen spans o logs derivados.

 ## Gates de fail-safe / degradación segura

 C.6 debe validar degradación segura y no optimista:

 - Si Cerbos no puede evaluarse correctamente, la salida debe ser fail-closed para casos sensibles o mutaciones.
 - Si falta binding resoluble o `provider_ref` compatible, no puede maquillarse como éxito parcial del corredor ejecutable.
 - Si falla observabilidad derivada, el checkpoint puede seguir pasando sólo si `compiled_contract`, `execution_record`, `policy_decision_record` y `event_record` quedan íntegros.
 - Si falla persistencia canónica de cualquier record estructural requerido, el checkpoint no puede pasar.
 - Si el smoke path cierra en `blocked`, debe existir `reason_code` auditable.
 - Si el smoke path cierra en `awaiting_approval`, debe existir `approval_request_id` auditable.
 - Si aparece `unknown_outcome`, debe existir evidence trail suficiente para demostrar por qué no hay certeza material.

 ## Artefactos mínimos que deben existir al cerrar C.6

 Al cerrar C.6 deben existir, como mínimo, estos artefactos normativos o persistidos en el corredor:

 1. `intent_input` gobernado con referencia de origen.
 2. `compiled_contract` persistido.
 3. `compilation_report` persistido.
 4. `execution_record` persistido y proyectado.
 5. `policy_input_v1` canonizado o evidencia persistida equivalente del request enviado.
 6. `policy_decision_record` persistido o evidencia normativa equivalente que conserve `policy_decision_id`.
 7. evidencia de resolution con `capability_id`, `bundle_digest`, `binding_id` y `provider_ref`.
 8. al menos un `event_record` canónico persistido del corredor.
 9. evidence trail del estado final del smoke path (`blocked`, `awaiting_approval` o `execution_completed`).

 ## Queries operativas mínimas que deben poder responderse

 Al cerrar C.6, el kernel ya debe poder responder preguntas operativas mínimas como estas:

 1. ¿Qué `compiled_contract` dio origen a este `execution_id`?
 2. ¿Qué `contract_fingerprint` estaba vigente cuando arrancó esta ejecución?
 3. ¿Qué decisión de Cerbos habilitó, bloqueó o derivó a approval este caso?
 4. ¿Cuál fue el `policy_decision_id` y con qué `trace_id` quedó correlacionado?
 5. ¿Qué `capability_id` intentó usar el contrato?
 6. ¿Con qué `binding_id`, `bundle_digest` y `provider_ref` se resolvió esa capability?
 7. ¿Cuál es el estado canónico actual de la ejecución y cuál fue su último `reason_code`?
 8. ¿Qué `event_record` canónicos existen para reconstruir el smoke path?
 9. ¿La evidencia canónica está completa aunque OTel/LGTM haya fallado?
 10. ¿El caso terminó `blocked`, `awaiting_approval` o `execution_completed`, y qué evidencia soporta ese cierre?

 ## Métricas mínimas del checkpoint

 C.6 no exige observabilidad madura, pero sí métricas mínimas para afirmar que el checkpoint es operable:

 - porcentaje de smoke paths que persisten `compiled_contract` correctamente;
 - porcentaje de smoke paths que crean/proyectan `execution_record` sin drift de correlación;
 - porcentaje de decisiones de policy con `policy_decision_record` correlacionado completo;
 - porcentaje de resolution paths con `binding_id` + `provider_ref` válidos;
 - porcentaje de smoke paths con al menos un `event_record` canónico persistido;
 - conteo de fallas donde OTel export falla pero el event log canónico permanece íntegro;
 - conteo de intentos rechazados por fail-closed estructural;
 - tasa de duplicados canónicos detectados en `execution_record` o `event_record`.

 El umbral normativo de C.6 no es performance; es **integridad estructural del corredor**.

 ## Gaps aceptables al cerrar C.6

 Son aceptables al cierre de C.6:

 - exportación OTel parcial, diferida o no confiable, si la verdad operativa sigue íntegra;
 - providers reales todavía stubbeados, siempre que la cadena de resolution canónica sea válida y trazable;
 - application phase no profundizada para todos los tipos de resultado;
 - UX humana, bandejas y surface conversacional todavía fuera del checkpoint;
 - dashboards, analytics y búsquedas avanzadas aún inmaduros;
 - verificaciones criptográficas profundas de firma/provenance todavía acotadas a evidencia mínima v1.

 ## Gaps NO aceptables al cerrar C.6

 NO son aceptables al cierre de C.6:

 - `compiled_contract` no persistido o persistido sin `contract_fingerprint` material;
 - `execution_record` ausente, duplicado materialmente o sin correlación mínima;
 - consultas a Cerbos sin `policy_input_v1` canonizado o sin evidencia persistida de decisión;
 - capability “resuelta” sin `binding_id` o sin `provider_ref` verificable;
 - depender de Temporal history, traces o logs para suplir ausencia de estado canónico en PostgreSQL;
 - smoke path sin `event_record` canónico persistido;
 - estado final `blocked` sin `reason_code`;
 - estado final `awaiting_approval` sin `approval_request_id`;
 - correlación rota entre `tenant_id`, `trace_id`, `contract_id`, `execution_id` y `policy_decision_id`;
 - intento de declarar checkpoint aprobado cuando sólo existe evidencia derivada y no verdad operativa durable.

 ## Tests borde mínimos

 Como mínimo, C.6 debe contemplar y documentar estos tests borde de integración:

 1. el compilador produce contrato pero falla proyección a runtime;
 2. el runtime crea ejecución pero falla policy input canonization;
 3. Cerbos responde `allow` pero falta binding resoluble;
 4. event log canónico persiste pero OTel export falla;
 5. OTel export existe pero falta event log canónico;
 6. binding válido pero provider runtime incompatible;
 7. `policy_decision_record` persiste con correlación incompleta;
 8. `execution_record` existe pero sin `contract_fingerprint`;
 9. `execution.completed` emitido sin `execution_record` previo;
 10. capability resoluble pero contract schema incompatible;
 11. mismatch entre `capability_id` del contrato y el registry;
 12. retry del smoke path duplica `execution_record`;
 13. `event_record` duplicado por replay técnico;
 14. checkpoint con `unknown_outcome` sin evidence trail;
 15. checkpoint con `blocked` sin `reason_code`;
 16. checkpoint con `awaiting_approval` sin `approval_request_id`;
 17. correlación rota entre `trace_id` y `execution_id`;
 18. checkpoint intenta pasar con history de Temporal pero sin estado canónico en PostgreSQL;
 19. existe `policy_decision_id` pero no puede enlazarse con `event_id` material;
 20. existe `binding_id` pero el `bundle_digest` persistido no coincide con el manifest resuelto.

 Criterio de lectura de estos bordes:

 - los casos 4 y equivalentes prueban que observabilidad derivada puede degradarse sin romper el checkpoint;
 - los casos 5, 8, 9, 18 y equivalentes prueban que sin verdad canónica el checkpoint falla aunque existan señales auxiliares;
 - los casos 3, 6, 10, 11 y 20 prueban que la resolution chain no puede maquillarse con best effort.

 ## Criterios de aceptación de C.6

 C.6 se considera cerrado a nivel de construcción v1 sólo si se cumplen todos estos criterios:

 1. existe un documento normativo del checkpoint de integración del kernel engine-only;
 2. queda explícito que distribution layer y UX conversacional no forman parte del checkpoint;
 3. queda definido el corredor mínimo end-to-end del kernel sin contradecir C.1-C.5;
 4. queda definido el smoke path mínimo obligatorio y sus cierres válidos (`blocked`, `awaiting_approval`, `execution_completed`);
 5. quedan definidos gates de consistencia entre compilación, runtime, policy, event log y resolution;
 6. queda definida la correlación mínima completa entre `tenant_id`, `trace_id`, `contract_id`, `contract_fingerprint`, `execution_id`, `policy_decision_id`, `event_id`, `capability_id`, `bundle_digest`, `binding_id` y `provider_ref`;
 7. queda explícito que OTel/LGTM puede fallar sin invalidar el checkpoint si la verdad operativa queda íntegra;
 8. queda explícito que si falla persistencia canónica, el checkpoint NO pasa aunque existan traces/logs/history;
 9. quedan listados artefactos mínimos, queries operativas mínimas, métricas mínimas y tests borde mínimos;
 10. queda explícito qué se difiere a Fase D y Fase E.

 ## Qué se difiere explícitamente a Fase D/E

 Se difiere a **Fase D**:

 - surface conversational and AI-friendly operations;
 - workflows de operador, UX humana y affordances de inspección más amigables;
 - APIs/surfaces pensadas para operación conversacional y consumo por agentes.

 Se difiere a **Fase E**:

 - hardening operacional, resiliencia ampliada y cierre de producción;
 - performance tuning, escalado fino, suites extensas de resiliencia y endurecimiento supply-chain;
 - profesionalización final de observabilidad, alerting y operación reusable a nivel producto.
