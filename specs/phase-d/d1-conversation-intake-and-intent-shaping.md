# D.1 — Conversation intake and intent shaping

## Objetivo

Definir la surface v1 que recibe conversación libre y la transforma, cuando corresponde, en artifacts gobernados de entrada (`intent_input` o `proposal_draft`) sin permitir ejecución directa desde chat, sin reinterpretar seams ya cerrados del kernel y sin contradecir D.0, C.1 ni C.6.

El objetivo NO es implementar UX final ni un runtime conversacional autónomo. El objetivo es fijar el boundary exacto, el shape mínimo, las decisiones de flow, la evidencia y la trazabilidad para que la conversación deje de ser texto suelto y pase a ser intake gobernado auditable.

## Principios de implementación del intake conversacional

- El intake layer es una **surface de interpretación gobernada**, no un seam nuevo del kernel.
- Chat libre NO es input ejecutable del kernel.
- La surface puede interpretar, resumir, estructurar y pedir aclaraciones, pero NO puede aplicar cambios ni disparar ejecución directa.
- Toda salida relevante del intake debe cerrar con una decisión explícita del intake: `continue_free_chat`, `ask_clarification`, `emit_intent_input`, `emit_proposal_draft` o `stop_out_of_scope`.
- El intake debe preservar por separado: afirmaciones del usuario, inferencias de la surface, dudas abiertas y decisión tomada.
- El intake debe fallar en modo **fail-closed** ante ambigüedad crítica, intento de saltear governance o pérdida de correlación mínima.
- El intake puede producir `proposal_draft` para cambios/configuración, pero nunca salteando diff, preview o simulación posterior.
- El intake debe apoyarse en IDs y evidencia compatibles con C.1 y C.6 para que el handoff no requiera reinterpretación libre posterior.
- El intake no redefine artifacts, policy, approvals, clasificación, fingerprinting ni runtime states: los consume como boundaries heredados.
- El intake no reintroduce distribution layer, activation tenant-scoped ni rollout operacional como parte de su scope.

## Boundary exacto entre chat libre, intent gobernado y proposal-trigger

### 1. `conversation_turn`

Es la unidad mínima de entrada conversacional. Representa un mensaje recibido o emitido dentro de una sesión. Puede contener texto libre, adjuntos y metadata contextual. Por sí solo NO autoriza compile, proposal ni apply.

### 2. `intent_candidate`

Es una interpretación intermedia, mutable y todavía incompleta de la intención. Sirve para shaping interno del intake. Puede contener supuestos, conflictos y preguntas abiertas. NO entra al kernel y NO puede usarse como `intent_input` mientras requiera reinterpretación libre.

### 3. `intent_input`

Es el objeto gobernado que ya satisface el boundary de C.1. El handoff a `intent_input` sólo puede ocurrir cuando el shape requerido por C.1 es construible sin reinterpretación libre posterior, con `objetivo` claro, alcance suficiente, evidencia de origen y correlación mínima.

### 4. `proposal_draft`

Es el artifact gobernado de surface para cambios/configuración que requiere proposal workspace. Puede originarse desde intake cuando la intención ya es suficientemente clara para abrir un workspace de propuesta, aunque todavía falten diff, preview o simulación. `proposal_draft` NO es apply candidate automático.

### Regla de boundary

- `conversation_turn` captura conversación.
- `intent_candidate` interpreta conversación.
- `intent_input` entrega entrada gobernada al compilador.
- `proposal_draft` entrega entrada gobernada al proposal flow.

Nunca deben colapsarse en el mismo objeto ni exponerse como equivalentes.

## Responsabilidades explícitas del intake layer

- Recibir `conversation_turn` con correlación mínima por tenant, sesión y sujeto.
- Clasificar cada turn como chat libre, shaping potencial, aclaración, out-of-scope o proposal-trigger potencial.
- Construir y actualizar `intake_session` como estado de shaping conversacional.
- Extraer `intent_candidate` desde uno o más turns relevantes.
- Separar hechos declarados, constraints declaradas, artifacts mencionados, supuestos inferidos y preguntas abiertas.
- Detectar ambigüedad crítica vs tolerable.
- Pedir aclaraciones cuando falta contexto material para avanzar.
- Cortar el flow cuando el pedido es incompatible, fuera de scope o intenta saltear governance.
- Emitir `intent_input` cuando ya existe shape suficiente y compatible con C.1.
- Emitir `proposal_draft` cuando se trata de cambio/configuración y ya puede abrirse un workspace gobernado sin saltarse preview/simulación.
- Persistir evidencia mínima de intake y reason codes auditables.
- Publicar correlación observable hacia event log y trazas derivadas sin reemplazar la verdad canónica.
- Aplicar idempotencia y deduplicación de intake a nivel de sesión/turn/decisión material.

## No-responsabilidades explícitas del intake layer

- No compilar contratos.
- No evaluar policy final en Cerbos.
- No emitir approvals ni marcar approvals como satisfechas.
- No generar `execution_record` ni arrancar workflows.
- No ejecutar providers, tools o mutaciones.
- No convertir chat libre en apply directo.
- No materializar diff definitivo ni preview final; eso pertenece a D.2/D.3.
- No escribir runtime state como verdad operativa.
- No inferir artifacts materiales sin marcarlos como sospechados o supuestos.
- No ocultar contradicciones o ambigüedades bajo resúmenes optimistas.
- No tratar distribution layer como scope válido de surface de Fase D.

## Shape mínimo de `conversation_turn`

`conversation_turn` debe cubrir como mínimo:

- `conversation_turn_id`
- `session_id`
- `tenant_id`
- `subject_id`
- `message_role`
- `raw_text`
- `attachments_refs`
- `timestamp`
- `turn_classification`

### Campos normativos adicionales recomendados para operar bien

- `trace_id` si ya existe; si no existe, debe poder asociarse en cuanto aparezca el primer hecho canónico.
- `reply_to_turn_id` opcional.
- `language` opcional.
- `ingestion_channel` opcional.
- `parse_status` para adjuntos/texto (`parsed`, `partial`, `failed`).
- `suspected_artifacts[]` como lectura no material todavía.
- `declared_constraints[]` extraídas literalmente.
- `surface_notes` sólo para metadata técnica, nunca para mezclar inferencias con texto del usuario.

### Semántica mínima

- `raw_text` conserva el texto original sin normalización destructiva.
- `turn_classification` puede ser, como mínimo: `free_chat`, `intent_signal`, `clarification_answer`, `governance_sensitive`, `out_of_scope`, `empty_or_greeting`.
- Un `conversation_turn` puede ser útil para shaping aunque no sea suficiente para emitir artifacts gobernados.

## Shape mínimo de `intake_session`

`intake_session` debe cubrir como mínimo:

- `intake_session_id`
- `session_id`
- `tenant_id`
- `subject_id`
- `current_state`
- `open_questions[]`
- `resolved_facts[]`
- `ambiguity_level`
- `last_decision`

### Campos normativos adicionales recomendados

- `trace_id`
- `started_at`
- `updated_at`
- `source_turn_ids[]`
- `active_intent_candidate_id` opcional
- `reason_codes[]`
- `critical_ambiguities[]`
- `tolerable_ambiguities[]`
- `evidence_refs[]`
- `deduplication_key`

### States mínimos esperados del intake session

- `free_chat`
- `shaping`
- `awaiting_clarification`
- `intent_ready`
- `proposal_ready`
- `out_of_scope`
- `closed`

### Reglas de transición mínimas

- `free_chat -> shaping` cuando aparece una señal material de objetivo o cambio.
- `shaping -> awaiting_clarification` cuando hay ambigüedad crítica.
- `shaping -> intent_ready` cuando el `intent_input` ya es construible sin reinterpretación libre.
- `shaping -> proposal_ready` cuando el `proposal_draft` ya es construible sin freeform adicional material.
- `* -> out_of_scope` cuando el pedido viola scope o governance.
- `intent_ready | proposal_ready | out_of_scope -> closed` cuando la decisión final del intake ya quedó emitida y evidenciada.

## Shape mínimo de `intent_candidate`

`intent_candidate` debe cubrir como mínimo:

- `intent_candidate_id`
- `source_turn_ids[]`
- `objetivo_candidate`
- `alcance_candidate`
- `artifacts_candidate[]`
- `constraints_candidate[]`
- `assumptions[]`
- `open_questions[]`
- `confidence_level`
- `ready_for_intent_input` bool
- `ready_for_proposal_draft` bool

### Semántica obligatoria

- `objetivo_candidate` nunca puede quedar implícito si `ready_for_intent_input = true`.
- `artifacts_candidate[]` distingue entre `mentioned`, `suspected` y `confirmed` a nivel semántico aunque el encoding exacto quede abierto.
- `assumptions[]` debe contener TODA inferencia no textual hecha por la surface.
- `open_questions[]` no puede omitirse cuando exista ambigüedad crítica o respuesta parcial.
- `confidence_level` no habilita por sí solo el handoff; sólo acompaña la decisión.
- `ready_for_intent_input` y `ready_for_proposal_draft` no pueden ser verdaderos por default: deben derivarse de reglas explícitas de shaping.

## Reglas de shaping de intención

1. El intake debe trabajar sobre uno o más `conversation_turn` correlados, no sobre un turno aislado si hay contexto material previo.
2. Toda extracción debe distinguir:
   - hecho declarado por usuario;
   - inferencia de surface;
   - duda abierta;
   - conflicto entre turns.
3. `objetivo_candidate` debe expresarse en términos operables y no meramente conversacionales.
4. `alcance_candidate` debe dejar claro si se trata de consulta, cambio, ajuste de configuración, revisión o pedido fuera de scope.
5. Si aparecen múltiples intents mezclados, el intake debe separarlos o pedir partición; no puede fusionarlos silenciosamente.
6. Si el usuario cambia de objetivo a mitad del shaping, el intake debe registrar override, conflicto o reinicio parcial del candidate.
7. Los artifacts mencionados de forma indirecta deben quedar como sospechados, no confirmados.
8. Constraints como urgencia, no downtime, ambiente, approvals o limitaciones operativas deben capturarse explícitamente.
9. Si una restricción fue inferida y no declarada, debe quedar en `assumptions[]`.
10. El intake debe detectar intención incompatible cuando el mismo shaping contiene pedidos mutuamente excluyentes.
11. El intake debe detectar intentos de saltear governance aunque el resto del pedido sea claro.
12. El intake no puede emitir `intent_input` si todavía haría falta reinterpretar texto libre para completar campos materiales requeridos por C.1.
13. El intake puede emitir `proposal_draft` antes de `intent_input` si el flow elegido de surface para cambios/configuración requiere proposal workspace primero, siempre que el draft no habilite ejecución directa.
14. El shaping debe preservar referencias a turns fuente para explicar de dónde salió cada conclusión relevante.
15. Una aclaración respondida parcialmente actualiza el candidate, pero no cierra preguntas pendientes por optimismo.

## Gestión de ambigüedad (crítica vs tolerable)

### Ambigüedad crítica

La ambigüedad es crítica cuando impide construir salida gobernada sin riesgo de reinterpretación material. Incluye como mínimo:

- scope material incierto;
- artifacts inciertos;
- intención incompatible;
- riesgo o clasificación potencialmente alterados por la ambigüedad;
- intento de saltear governance;
- ambiente/tenant/sujeto materialmente inciertos si afectan el caso;
- contradicciones entre turns sobre qué cambiar o qué evitar.

Con ambigüedad crítica, el intake sólo puede:

- `ask_clarification`, o
- `stop_out_of_scope`.

### Ambigüedad tolerable

La ambigüedad es tolerable cuando NO altera el tipo material del caso y puede preservarse como supuesto explícito sin bloquear proposal/intake. Ejemplos:

- wording humano del resumen;
- detalle accesorio de contexto no material;
- artifact secundario sospechado pero no necesario para abrir proposal;
- preferencia de formato de respuesta.

Con ambigüedad tolerable, el intake puede avanzar sólo si:

- la deja explícita en evidencia;
- la marca como supuesto o duda abierta;
- no cambia la clasificación del handoff.

## Preguntas de aclaración y cortes del flow

### Cuándo preguntar

El intake debe preguntar aclaraciones cuando falte alguna de estas piezas materiales:

- objetivo claro;
- alcance suficientemente delimitado;
- artifacts/configuración primarios afectados;
- constraint operativa relevante;
- confirmación de intent conflictivo o contradictorio;
- información mínima para cumplir shape de C.1 o abrir proposal workspace sin reinterpretación libre.

### Cómo preguntar

Las preguntas deben ser:

- puntuales;
- orientadas a cerrar una ambigüedad concreta;
- trazables a un reason code;
- acotadas a lo necesario para avanzar.

El intake no debe abrir entrevistas largas si la causa real es out-of-scope o governance-sensitive.

### Cortes explícitos del flow

Debe cortarse el flow con `stop_out_of_scope` cuando:

- el pedido cae en distribution layer;
- el pedido intenta apply directo salteando proposal/preview/simulación;
- el usuario exige bypass de approvals, policy o clasificación;
- el pedido es incompatible con los seams cerrados del kernel;
- falta contexto mínimo estructural y no hay base razonable para preguntar algo útil.

## Handoff exacto hacia `intent_input` o `proposal_draft`

### Handoff a `intent_input`

Sólo puede ocurrir cuando:

- existe `objetivo` claro;
- existe alcance utilizable;
- la salida es compatible con el shape requerido por C.1;
- la correlación mínima está completa o es generable antes del primer hecho canónico;
- las inferencias no resueltas NO obligan reinterpretación libre posterior;
- la evidencia textual mínima ya quedó persistida.

El handoff debe producir un artifact que ya pueda alimentar el compilador de C.1 sin que el compilador tenga que volver al chat a “entender” significado.

### Handoff a `proposal_draft`

Debe ser posible para cambios/configuración cuando:

- el usuario está pidiendo cambio material o configuración gobernada;
- ya hay suficiente claridad para abrir un workspace de propuesta;
- todavía corresponde diff/preview/simulación antes de cualquier apply;
- los artifacts primarios están confirmados o al menos sospechados explícitamente;
- la evidencia deja claro qué parte es intención del usuario y qué parte es interpretación de surface.

El handoff a `proposal_draft` NO puede saltearse diff, preview, simulación ni governance posterior.

## Evidencia mínima de intake

Todo paso relevante del intake debe preservar evidencia textual mínima con estas categorías:

- afirmaciones del usuario;
- inferencias de la surface;
- dudas abiertas;
- decisión tomada por el intake.

### Evidencia mínima requerida por decisión

#### Para `continue_free_chat`
- turn o conjunto de turns fuente;
- motivo por el cual todavía no hay shaping material.

#### Para `ask_clarification`
- ambigüedad crítica detectada;
- pregunta emitida;
- fields bloqueados por esa ambigüedad.

#### Para `emit_intent_input`
- resumen gobernado del objetivo y alcance;
- turns fuente;
- constraints declaradas e inferidas;
- confirmación de que ya no hace falta reinterpretación libre posterior.

#### Para `emit_proposal_draft`
- resumen gobernado del cambio/configuración pedido;
- artifacts afectados confirmados o sospechados;
- supuestos abiertos tolerables;
- confirmación de que preview/simulación siguen pendientes.

#### Para `stop_out_of_scope`
- razón del corte;
- evidencia textual suficiente del pedido incompatible;
- boundary o constraint violado.

## Reason codes mínimos de intake

Debe existir, como mínimo, esta base de reason codes:

- `intake.continue_free_chat`
- `intake.ask_clarification.scope_ambiguous`
- `intake.ask_clarification.artifact_ambiguous`
- `intake.ask_clarification.intent_conflict`
- `intake.ask_clarification.governance_sensitive`
- `intake.emit_intent_input`
- `intake.emit_proposal_draft`
- `intake.stop_out_of_scope`
- `intake.stop_distribution_layer_out_of_scope`
- `intake.stop_missing_minimum_context`

### Regla normativa

Toda decisión final de intake por paso relevante debe poder mapearse al menos a uno de esos reason codes, pudiendo agregarse taxonomía más fina después sin romper compatibilidad.

## Integración con compilador/proposal flow

- `emit_intent_input` integra con C.1 como upstream gobernado.
- `emit_proposal_draft` integra con D.2 como upstream del proposal workspace.
- El intake no elige por conveniencia de UX: elige según tipo de caso y suficiencia del shape disponible.
- Si el caso es de cambio/configuración, la surface puede priorizar `proposal_draft` como artifact inicial, manteniendo la posibilidad de derivar luego a `intent_input` cuando corresponda al compilador.
- Si el caso es de consulta o acción ya suficientemente formalizable, puede emitirse `intent_input` directo sin abrir proposal workspace, siempre respetando C.1.
- Ningún camino desde intake puede saltar directo a ejecución o apply.

## Integración con event log y observabilidad

El intake debe correlacionarse como mínimo con:

- `tenant_id`
- `session_id`
- `trace_id`
- `subject_id`
- `conversation_turn_id`
- `intake_session_id`

### Reglas mínimas

- Si `trace_id` no existe al inicio, debe generarse antes del primer hecho canónico de intake o registrarse un mecanismo de generación determinable.
- El event log puede recibir hechos de intake como evidencia operativa, pero esos eventos no reemplazan los artifacts gobernados emitidos.
- Observabilidad derivada puede resumir shaping, tiempos y decisiones, pero nunca sustituye evidencia textual mínima ni razón normativa del handoff/corte.
- Debe poder reconstruirse qué turns originaron qué candidate, qué sesión de intake tomó qué decisión y con qué correlación terminó el handoff.

## Idempotencia y deduplicación del intake

- La misma conversación reingresada no debe duplicar `intake_session` si mantiene misma identidad material (`tenant_id`, `session_id`, `subject_id` y set de turns equivalente).
- Un mismo `conversation_turn` no debe recontarse como turn nuevo si reingresa con el mismo identificador material o fingerprint equivalente.
- Reprocesar el mismo set de turns debe reutilizar o reconciliar el `intent_candidate` previo, salvo cambio material detectado.
- Un cambio material de objetivo, scope o artifacts principales debe abrir nueva versión lógica del candidate o nueva decisión de intake; NO debe sobreescribirse silenciosamente.
- La deduplicación del intake es local al boundary conversacional; no reemplaza la deduplicación de C.1 ni la idempotencia de runtime.

## Tests borde mínimos

1. mensaje vacío o solo saludo → clasifica `continue_free_chat` sin emitir candidate material.
2. múltiples intents mezclados en un mismo turn → pide partición o aclaración, no fusiona silenciosamente.
3. usuario cambia objetivo a mitad del shaping → registra conflicto/override y reevalúa readiness.
4. artifacts mencionados de forma ambigua → `ask_clarification.artifact_ambiguous`.
5. intento explícito de saltear approvals/policy → `ask_clarification.governance_sensitive` o `stop_out_of_scope`, nunca apply.
6. pedido que cae en distribution layer → `stop_distribution_layer_out_of_scope`.
7. scope material incierto → `ask_clarification.scope_ambiguous`.
8. usuario pide apply directo sin proposal/preview → corte o redirección gobernada, nunca ejecución directa.
9. attachment presente pero no parseable → evidencia de parse failure y no se lo toma como hecho confirmado.
10. intake emite `intent_input` sin `objetivo` claro → debe fallar validación del intake.
11. intake emite `proposal_draft` sin artifacts sospechados → debe fallar validación del draft inicial.
12. misma conversación reingresada no duplica intake session.
13. `trace_id` ausente al inicio → se genera o se registra estrategia antes del primer hecho canónico.
14. cambio de tenant/session en mitad del intake → corta o reinicia boundary; no mezcla contextos.
15. surface infiere constraint no dicho y no lo marca como suposición → debe fallar evidencia mínima.
16. pregunta de aclaración respondida parcialmente → mantiene preguntas abiertas y no emite salida final optimista.
17. cierre en out_of_scope sin evidencia textual → debe fallar validación de evidencia.
18. conversación larga con contradicciones entre turns → preserva conflicto explícito y evita handoff prematuro.
19. turn con lenguaje coloquial pero objetivo claro → puede avanzar si el shape estructural queda claro y auditable.
20. usuario confirma artifacts después de una sospecha inicial → actualiza candidate de sospechado a confirmado sin perder histórico.
21. misma intención con wording distinto en turns consecutivos → no crea candidates duplicados si no cambia materialmente.
22. turn del asistente resume mal y el usuario corrige → prevalece corrección del usuario y se registra la inferencia previa como descartada.
23. caso read-only suficientemente claro → permite `emit_intent_input` sin proposal draft.
24. cambio/configuración suficientemente claro pero todavía sin preview → permite `emit_proposal_draft` y bloquea cualquier apply directo.

## Criterios de aceptación de D.1

- Queda explícito que chat libre NO es input ejecutable del kernel.
- Queda fijado el boundary entre `conversation_turn`, `intent_candidate`, `intent_input` y `proposal_draft`.
- Queda explícito que el intake interpreta, resume y pregunta, pero NO aplica cambios ni dispara ejecución directa.
- Toda decisión relevante del intake termina en una de estas salidas: `continue_free_chat`, `ask_clarification`, `emit_intent_input`, `emit_proposal_draft` o `stop_out_of_scope`.
- Queda definida la distinción entre ambigüedad crítica y tolerable.
- La ambigüedad crítica incluye scope incierto, artifacts inciertos, intención incompatible, riesgo/clasificación potencialmente alterados e intento de saltear governance.
- El handoff a `intent_input` sólo ocurre cuando C.1 puede consumirse sin reinterpretación libre posterior.
- El handoff a `proposal_draft` queda habilitado para cambios/configuración sin saltearse diff, preview ni simulación.
- Queda definida evidencia textual mínima con afirmaciones del usuario, inferencias de surface, dudas abiertas y decisión tomada.
- Quedan definidos reason codes mínimos y correlación mínima del intake.
- Quedan definidos estados mínimos, reglas de idempotencia y tests borde suficientes para validar surface v1.
