# Phase 3 — Internal contracts and transport

## Objetivo

Definir el profile de contratos internos y transporte derivado para servicios del baseline convergente, sin mover la fuente de verdad normativa fuera de YAML + JSON canónico.

## Principios del transport/internal contracts profile

- El contrato normativo manda; el contrato de transporte deriva.
- El transporte interno debe ser fuerte, versionable y regenerable.
- La ergonomía técnica no justifica una segunda fuente de verdad.
- La compatibilidad se gestiona explícitamente, no por accidente.
- El profile debe servir para servicios internos, admin APIs y clients controlados sin contaminar la semántica de dominio.

## Qué problema resuelve

El baseline reusable v1 ya define la verdad ejecutable, pero no fija todavía un profile interno uniforme para intercambio tipado entre servicios. Esta fase resuelve esa falta de disciplina técnica: define cómo derivar Protobuf, qué papel juega Buf, cuándo usar ConnectRPC/gRPC y cómo evitar triple source of truth.

## Boundary exacto entre contratos normativos y contratos de transporte

- **Contratos normativos**: objetos canónicos, envelopes, semántica de campos, invariantes, versionado material/no material y estados del dominio.
- **Contratos de transporte**: mensajes, RPCs, services, envelopes internos, metadata técnica y convenciones de serialización para mover datos entre componentes.

Regla dura: ningún cambio aceptable puede nacer en el contrato de transporte y recién después “subirse” al normativo. La dirección válida es siempre normativo -> derivado.

## Artefactos que generan Protobuf

Protobuf puede generar artefactos derivados para:

- mensajes internos entre servicios
- contratos RPC internos
- stubs de cliente/servidor
- validaciones estructurales de transporte
- documentación de APIs internas
- checks de compatibilidad sobre schemas derivados

No genera la semántica normativa del dominio; solo la representa para transporte.

## Reglas de generación derivada YAML/JSON -> Protobuf

- Los objetos canónicos se modelan primero en YAML normativo.
- El JSON canónico compilado es la proyección ejecutable estable.
- Los mensajes Protobuf se generan o mantienen como representación derivada de ese modelo.
- Los nombres y campos derivados deben mantener correspondencia trazable con el modelo normativo.
- Si el modelo normativo cambia materialmente, primero cambia YAML/JSON; luego se regenera o ajusta Protobuf.
- Si aparece un campo técnico solo de transporte, debe estar claramente marcado como derivado y no normativo.

## Rol de Buf

- Mantener disciplina de versionado y compatibilidad de schemas derivados.
- Detectar cambios incompatibles backward/forward en contratos internos.
- Estandarizar generación de stubs y documentación técnica derivada.
- Hacer visible el costo de cambio en wire contracts antes de romper consumers internos.

## Rol de ConnectRPC/gRPC

- **ConnectRPC**: profile HTTP-friendly para clientes internos, UI/admin surfaces y operabilidad pragmática.
- **gRPC**: profile binario eficiente para comunicación interna donde su costo esté justificado.
- Ambos son mecanismos de transporte interno; ninguno es la autoridad normativa del sistema.

## Versionado y compatibilidad backward/forward

- El versionado normativo del dominio sigue reglas `major.minor` ya cerradas.
- El versionado de Protobuf debe alinearse con la evolución normativa y no inventar semántica independiente.
- Los consumers internos deben tolerar adición backward-compatible de campos derivados cuando no altera semántica normativa.
- Los cambios incompatibles en transporte requieren coordinación explícita y no pueden ocultar una ruptura semántica del contrato normativo.
- Los aliases, deprecations y ventanas de migración deben ser explícitos.

## Política anti triple-source-of-truth

- YAML normativo es la primera verdad.
- JSON canónico compilado es la segunda forma de la misma verdad, no una verdad independiente.
- Protobuf es representación derivada de transporte.
- Queda prohibido mantener divergencias manuales persistentes entre YAML, JSON y Protobuf.
- Queda prohibido introducir campos “solo existen en Protobuf” si alteran semántica de dominio.
- Toda regeneración o ajuste manual debe dejar trazabilidad de qué objeto normativo originó ese cambio.

## Riesgos y guardrails

### Riesgos

- que Protobuf capture la evolución real del dominio por comodidad técnica
- que ConnectRPC/gRPC se conviertan en pseudo-API pública normativa
- que Buf valide compatibilidad técnica pero se escape una incompatibilidad semántica
- que aparezcan envelopes internos que contradigan el envelope canónico del dominio

### Guardrails

- revisión obligatoria contra objetos normativos antes de publicar cambios de wire contracts
- naming derivado alineado con el canon del dominio
- documentación explícita de campos técnicos vs normativos
- no publicar RPCs internas que oculten estados o invariantes centrales del dominio
- no aceptar compatibilidad técnica si rompe semántica de negocio o governance

## Tests borde mínimos

1. cambio no material en YAML propaga actualización backward-compatible a Protobuf
2. cambio material en YAML obliga revisión explícita de contratos derivados
3. mensaje Protobuf con campo extra técnico no altera el modelo normativo
4. intento de agregar semántica solo en Protobuf debe bloquearse
5. Buf detecta ruptura backward en un message derivado
6. ConnectRPC expone endpoint compatible con schema derivado sin alterar canon
7. gRPC transporta payload interno sin perder envelope técnico requerido
8. consumer viejo ignora campo nuevo opcional sin romperse
9. consumer nuevo tolera ausencia de campo legado deprecado durante ventana válida
10. serialización JSON derivada mantiene correspondencia con envelope canónico
11. un rename normativo exige mapping explícito y no alias silencioso permanente
12. contract generation falla si no existe traza al objeto normativo origen
13. cambio de wire contract no aprobado no puede publicarse como si fuera normativo
14. schema derivado no puede introducir estados de runtime no existentes en canon
15. documentación generada distingue correctamente campo normativo de campo técnico
16. compatibilidad técnica aprobada pero incompatibilidad semántica detectada debe bloquear release derivada

## Criterios de aceptación

1. Queda fijado el boundary entre contratos normativos y de transporte.
2. Queda explícito que Protobuf es derivado y no source of truth.
3. Queda escrito el rol de Buf y ConnectRPC/gRPC.
4. Queda fijada la política anti triple-source-of-truth.
5. Queda fijado el esquema de versionado y compatibilidad derivada.
6. Quedan documentados riesgos y guardrails concretos.
7. Existen tests borde suficientes para validar el profile.
8. La fase no contradice la regla normativa YAML + JSON canónico.
