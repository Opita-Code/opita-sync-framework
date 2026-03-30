# B.0 — Marco común de decisión para comparativas irreversibles

## Principios base

- Fase B existe para cerrar decisiones con alto costo de reversión antes de construir el kernel en Fase C.
- Una comparativa irrelevante o superficial agrega ruido; una comparativa irreversible debe reducir riesgo sistémico real.
- Ninguna decisión puede optimizar una capa aislada rompiendo invariantes ya cerrados en Fase A.
- El criterio rector no es la conveniencia local de implementación, sino la coherencia operativa del sistema completo: multi-tenant, gobernado, auditable y operable por humanos e IA.
- Toda recomendación debe explicitar tradeoffs, condiciones de validez y costo de salida futura.
- Ningún score compensa una violación de hard gates. Primero se filtra por viabilidad estructural; recién después se compara por puntaje.

---

## Qué cuenta como decisión irreversible en Opyta Sync

En Opyta Sync, una decisión se considera irreversible cuando cambiarla más adelante obliga a reescribir o revalidar partes estructurales del kernel, sus artifacts o su modelo operativo. No se trata solo de costo técnico; también importa el costo de reauditoría, recertificación conceptual y retrabajo sobre seams ya consumidos por otras capas.

Una decisión entra en esta categoría si cumple una o más de estas condiciones:

| Criterio | Qué implica |
|---|---|
| Reescribe el modelo de ejecución | Cambia cómo se representan estados, retries, compensaciones, pausas o reanudaciones |
| Afecta evidencia y auditoría | Obliga a redefinir trazabilidad, event log, snapshots o reconstrucción de historia |
| Condiciona contratos entre subsistemas | Fuerza cambios en policy, extensibilidad, memoria, telemetría o configuración conversacional |
| Eleva lock-in operativo | Hace costoso migrar tooling, runtime o mecanismos de observabilidad ya integrados |
| Cambia la superficie de operación | Exige reaprender cómo inspeccionar, depurar, gobernar o operar el sistema |

Para esta fase, las decisiones irreversibles prioritarias son: durable runtime, policy engine, memoria operativa y telemetría, extensibilidad, packaging de capabilities y configuración conversacional.

---

## Orden recomendado de comparativas con razón técnica

El orden de comparativas no es arbitrario. Debe seguir la dirección de dependencia técnica del kernel, empezando por aquello que define el modelo real de ejecución y terminando por aquello que gobierna artefactos ya estabilizados.

| Orden | Comparativa | Razón técnica |
|---|---|---|
| 1 | durable runtime | Va primero porque tiene el mayor blast radius técnico. Condiciona event log, retries, compensación, inspección de estado, correlación de evidencia y seams de extensibilidad. |
| 2 | policy engine | Va segundo porque cruza permisos, approvals y clasificación, pero debe encajar sobre el modelo real de ejecución y sobre los artifacts que el runtime expone. |
| 3 | memoria operativa y telemetría | Va tercero porque depende del runtime y del event model para observar, reconstruir, explicar y operar ejecuciones sin inventar semántica paralela. |
| 4 | extensibilidad | Va después porque los seams de extensión deben montarse sobre un núcleo operativo ya definido y no al revés. |
| 5 | packaging de capabilities | Va después de extensibilidad porque empaquetar, distribuir e instalar capabilities depende del modelo de extensión, ownership y lifecycle ya decidido. |
| 6 | configuración conversacional | Va al final porque gobierna y modifica artefactos ya definidos por todas las decisiones previas; no debe imponer estructura antes de que exista el núcleo estable. |

---

## Hard gates obligatorios previos a cualquier score

Antes de asignar cualquier puntaje, toda alternativa debe pasar un filtro binario de viabilidad. Si falla uno o más hard gates, la alternativa queda marcada como `hard_gate_fail` y no participa del scoring comparativo final.

| Hard gate | Exigencia mínima |
|---|---|
| soporte real para multi-tenant | Debe soportar aislamiento, gobernanza y operación tenant-scoped sin soluciones cosméticas ni dependencia de particiones manuales frágiles. |
| trazabilidad auditable | Debe permitir reconstruir qué pasó, cuándo pasó, por qué pasó y qué evidencia quedó, sin depender de observabilidad informal. |
| hooks o seams para approvals/classification/policy | Debe ofrecer puntos de integración claros para decisiones de governance y clasificación sin hacks laterales ni bifurcaciones del flujo principal. |
| estado durable inspeccionable | Debe permitir inspeccionar estado, historial, pausas, retries, compensaciones o bloqueos con semántica operable. |
| ausencia de cajas negras imposibles de operar | Debe evitar dependencias opacas que no puedan explicarse, depurarse ni gobernarse con criterio de producción. |

### Regla de aplicación

- Los hard gates se evalúan antes del score ponderado.
- Un hard gate fallido no se “compensa” con fortalezas en otros criterios.
- Si una alternativa requiere supuestos no probados para pasar un hard gate, se considera fallido hasta demostrar lo contrario.

---

## Criterios comunes de evaluación con pesos exactos

Las comparativas B.1-B.6 deben compartir una misma base de evaluación para que el resultado sea comparable y acumulable a nivel de fase. Cada comparativa puede agregar criterios específicos del dominio evaluado, pero no alterar esta base común ni sus pesos.

| Criterio común | Peso | Qué se evalúa |
|---|---:|---|
| alineación con invariantes de Fase A | 20% | Compatibilidad con contratos, governance, runtime semantics, auditabilidad y verdad ejecutable ya cerrados en Fase A |
| gobernanza y seguridad multi-tenant | 18% | Aislamiento, control de acceso, separación de deberes, boundaries por tenant y capacidad de gobernar operación sensible |
| durabilidad, auditabilidad y recuperación | 16% | Persistencia confiable, recuperación ante fallos, replay, reconstrucción de estado y evidencia durable |
| encaje con event log / trazabilidad / evidencia | 12% | Cómo se integra con correlación de eventos, event log, snapshots, reason codes y evidencia operacional |
| extensibilidad futura sin romper el core | 10% | Capacidad de agregar nuevas capacidades, hooks y módulos sin reescribir el núcleo ni degradar sus invariantes |
| operabilidad IA-friendly | 9% | Facilidad para inspección, explicación, simulación, observabilidad y operación asistida por IA sin opacidad excesiva |
| complejidad de implementación y riesgo de entrega en Fase C | 8% | Riesgo real de entrega, curva de implementación, complejidad accidental y probabilidad de bloquear el roadmap inmediato |
| lock-in / costo de reversión futura | 4% | Costo técnico y operativo de migrar o reemplazar la opción una vez integrada |
| madurez de ecosistema y tooling | 3% | Calidad del tooling, documentación, comunidad y señales de estabilidad de uso |

**Total:** 100%

---

## Escala de scoring recomendada (1-5) y cómo se calcula score ponderado

La escala recomendada para cada criterio es de 1 a 5, donde el número expresa aptitud relativa para Opyta Sync y no fama de mercado genérica.

| Score | Significado |
|---|---|
| 1 | Deficiente: genera conflicto serio con el criterio o exige workarounds de alto riesgo |
| 2 | Débil: podría usarse, pero con tradeoffs fuertes, deuda significativa o fragilidad operativa |
| 3 | Aceptable: cumple lo mínimo razonable, aunque sin ventajas claras y con restricciones explícitas |
| 4 | Fuerte: encaja bien con el criterio y deja pocos compromisos relevantes |
| 5 | Excelente: encaje sobresaliente, con ventajas claras y bajo costo operacional relativo |

### Cálculo del score ponderado

1. Se asigna un score de 1 a 5 a cada criterio común.
2. Cada score se multiplica por su peso porcentual.
3. La suma de todos los resultados se divide por 5 para normalizar el total a una escala final de 0 a 100.

### Fórmula de referencia

`score_ponderado = ((score_criterio × peso) + ...) / 5`

### Ejemplo conceptual de lectura

- Una alternativa con puntajes altos en criterios críticos y sin hard gates fallidos tenderá a ubicarse en la zona `preferred`.
- Una alternativa con score razonable pero tradeoffs estructurales explícitos puede quedar en `acceptable_with_tradeoffs`.
- Una alternativa con score alto pero un hard gate fallido sigue siendo `hard_gate_fail`.

---

## Qué significa `hard_gate_fail`, `preferred`, `acceptable_with_tradeoffs`, `reject`

| Estado | Significado operativo |
|---|---|
| `hard_gate_fail` | La opción queda descartada antes del scoring final porque no cumple una condición estructural obligatoria. |
| `preferred` | Es la alternativa recomendada porque pasa hard gates y ofrece el mejor balance global para Opyta Sync dentro del contexto real del proyecto. |
| `acceptable_with_tradeoffs` | Puede adoptarse si existe razón contextual suficiente, pero exige aceptar costos, límites o mitigaciones explícitas. |
| `reject` | No se recomienda adoptar la opción aunque haya sido evaluable, porque su balance final es inferior o introduce riesgo innecesario frente a otras alternativas. |

### Regla práctica

La etiqueta final no depende solo del score numérico. También depende de la naturaleza de los tradeoffs, del costo de reversión y del nivel de tensión con los invariantes del sistema.

---

## Formato común que deben seguir las comparativas de B.1-B.6

Toda comparativa de Fase B debe seguir una estructura homogénea para que la revisión sea auditable, comparable y reusable.

| Sección obligatoria | Contenido esperado |
|---|---|
| Contexto y decisión a resolver | Qué se está decidiendo, por qué es irreversible y qué subsistemas toca |
| Invariantes y constraints heredados | Qué decisiones de Fase A no pueden romperse |
| Candidatos comparados | Lista explícita de alternativas evaluadas y alcance real de cada una |
| Hard gates | Evaluación binaria de gates obligatorios con evidencia o razonamiento técnico |
| Criterios específicos del bloque | Criterios adicionales propios del dominio comparado |
| Scoring común | Tabla con criterios comunes, pesos, score por candidato y lectura resumida |
| Tradeoffs narrativos | Costos, ventajas, riesgos y condiciones de adopción |
| Recomendación | Opción recomendada y justificación técnica |
| Riesgos asumidos y mitigaciones | Qué se acepta conscientemente y cómo se piensa reducir el riesgo |
| Impacto sobre Fase C | Qué habilita, qué condiciona y qué artefactos deberán respetar la decisión |

---

## Artefactos mínimos que debe dejar cada comparativa

Cada comparativa no solo debe concluir una recomendación; también debe dejar evidencia suficiente para que Fase C no tenga que reinterpretar la decisión.

| Artefacto mínimo | Propósito |
|---|---|
| documento comparativo publicado | Dejar la comparación completa, trazable y revisable |
| tabla de scoring común | Hacer visible la evaluación homogénea entre candidatos |
| veredicto final por candidato | Marcar `hard_gate_fail`, `preferred`, `acceptable_with_tradeoffs` o `reject` |
| decisión recomendada explícita | Evitar ambigüedad sobre qué opción guía Fase C |
| lista de tradeoffs aceptados | Documentar el costo consciente de la decisión |
| riesgos y mitigaciones | Preparar implementación y operación con criterios realistas |
| impacto sobre artifacts del kernel | Indicar qué contratos, eventos, seams o interfaces quedan condicionados |
| preguntas abiertas no bloqueantes | Separar deuda tolerable de deuda que reabriría la decisión |

---

## Criterios de aceptación de B.0

B.0 puede considerarse cerrado cuando se cumplen, como mínimo, todas las condiciones siguientes:

1. existe un marco común explícito para comparar decisiones irreversibles de Fase B;
2. el orden recomendado de comparativas está definido y justificado técnicamente;
3. los hard gates obligatorios están publicados y se aplican antes de cualquier score;
4. los criterios comunes de evaluación tienen pesos exactos y suman 100%;
5. la escala de scoring y el método de cálculo del score ponderado están definidos sin ambigüedad;
6. los estados `hard_gate_fail`, `preferred`, `acceptable_with_tradeoffs` y `reject` tienen semántica operativa clara;
7. el formato común de comparativas B.1-B.6 está definido de forma reusable;
8. los artefactos mínimos exigidos por comparativa están explicitados;
9. el documento es consistente con el cierre de Fase A y con el objetivo de reducir riesgo de reversión antes de Fase C.

### Resultado esperado de B.0

Si B.0 está correctamente cerrado, cada comparativa posterior podrá discutir alternativas distintas sin rediscutir el método base. Ese es el objetivo REAL: separar la discusión de contenido de la discusión de marco, para que Fase B produzca decisiones acumulables y no opiniones sueltas.
