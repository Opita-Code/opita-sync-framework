# Phase 4 — Artifact plane, cache and retrieval

## Objetivo

Definir el plano convergente de artifacts, caché y retrieval, fijando el boundary exacto entre PostgreSQL, object storage, Valkey y OpenSearch, junto con el baseline de ingestión documental segura.

## Principios del artifact/cache/retrieval plane

- Cada plano existe para una función distinta y no debe invadir al otro.
- La verdad operativa durable sigue en PostgreSQL.
- El artifact plane persiste blobs, adjuntos y evidencia pesada fuera del truth plane relacional.
- El cache plane acelera; no gobierna.
- El retrieval plane recupera corpus y contexto; no decide verdad ni policy.

## Boundary exacto entre PostgreSQL, object storage, Valkey y OpenSearch

- **PostgreSQL**: truth plane operativo durable, metadatos relacionales, correlación, event log operativo, referencias y estado de control.
- **S3-compatible object storage**: blobs, attachments, evidence pesada, snapshots, exports y derivados binarios/documentales.
- **Valkey**: caché efímera, locks cortos, buckets de rate limiting, hints de idempotencia y material temporal no autoritativo.
- **OpenSearch**: retrieval/corpus plane e índices de búsqueda híbrida para documentos, artefactos proyectados y contexto amplio.

## Qué va a PostgreSQL

- contratos y referencias a contratos compilados
- approvals, execution records y runtime/event metadata canónica
- event log operativo append-only ya cerrado
- catálogo, bindings y referencias a providers
- metadata estructural de artifacts y evidence
- estados de ingesta y correlación con objetos externos
- punteros a object storage y a índices de retrieval

## Qué va a S3-compatible object storage

- attachments de OCI bundles
- evidencia pesada y exportable
- snapshots de input/output redactados cuando correspondan
- archivos de demo, reportes, exports y anexos
- documentos originales o derivados que no deban vivir inline en PostgreSQL

## Qué va a Valkey

- caches calientes de resolución y lookups repetidos
- idempotency hints y dedup temporal acelerada
- locks cortos de coordinación
- session hints y buckets de rate limiting
- material efímero regenerable usado para performance operacional

## Qué va a OpenSearch

- índices de corpus documental
- índices de retrieval híbrido y semántico complementario
- representaciones indexables de artefactos y evidence seleccionada
- metadatos documentales útiles para búsqueda amplia
- fragmentos enriquecidos para recuperación controlada

## Qué nunca debe vivir en cada plano

### Nunca en PostgreSQL
- blobs pesados que rompan el perfil operativo del truth plane
- corpus completo de documentos para retrieval masivo

### Nunca en S3-compatible object storage
- estado transaccional autoritativo
- decisiones de policy o approvals como única copia viva

### Nunca en Valkey
- verdad durable
- evidencia única irrecuperable
- estado canónico de runtime o governance

### Nunca en OpenSearch
- source of truth del dominio
- observabilidad general del sistema
- approvals, decisiones de policy o runtime state como autoridad primaria

## Retrieval/corpus plane con OpenSearch

OpenSearch queda fijado como retrieval/corpus plane complementario. Su función es indexar y recuperar corpus documental, artefactos seleccionados y contexto amplio para inspección, búsqueda y ayuda operativa. No reemplaza PostgreSQL, no reemplaza OTel/LGTM y no absorbe la semántica del event log operativo.

## Document ingestion pipeline baseline

La cadena baseline de ingestión documental queda así:

1. **ClamAV** para escaneo de malware y bloqueo inicial.
2. **Apache Tika** para extracción de texto y metadata útil.
3. **Presidio** para detección de PII/sensibilidad y apoyo a redacción.
4. **Embeddings** vía plano IA derivado cuando esté justificado.
5. **Metadata estructural** persistida en PostgreSQL.
6. **Contenido indexable derivado** proyectado a OpenSearch.
7. **Documento original/derivado pesado** persistido en object storage cuando corresponda.

## Rules de redacción, clasificación y evidence en ingestión

- Ningún documento entra al retrieval plane sin pasar por clasificación mínima y scanning inicial.
- Si el documento contiene material sensible, la versión indexable debe ser redactada o reducida según policy.
- PostgreSQL guarda metadata, estado de ingesta, referencias y correlación; no el corpus masivo.
- Object storage guarda el original o derivado pesado bajo referencias gobernadas.
- OpenSearch solo indexa material permitido para retrieval.
- Toda etapa de ingesta debe emitir evidence refs y correlación con `artifact_ref`, `trace_ref` y `classification_level`.
- El pipeline debe distinguir claramente original, derivado indexable y derivado redactado.

## Tests borde mínimos

1. documento limpio entra a pipeline y genera metadata + refs correctas
2. documento con malware queda bloqueado antes de extracción
3. documento vacío produce resultado controlado sin indexación inválida
4. documento con PII exige derivado redactado para retrieval
5. documento sensible puede persistir original en object storage pero no texto completo en OpenSearch
6. metadata de ingesta queda correlacionada en PostgreSQL
7. artifact_ref roto bloquea publicación al retrieval plane
8. reingesta del mismo documento respeta idempotencia esperada
9. embeddings fallan y pipeline degrada sin perder metadata ni evidence
10. Tika extrae texto parcial y eso queda marcado explícitamente
11. OpenSearch recibe solo el derivado permitido y no el original restringido
12. Valkey puede cachear lookup de ingesta sin convertirse en autoridad del estado
13. object storage indisponible bloquea cierre exitoso si el original era obligatorio
14. PostgreSQL jamás se usa como blob store masivo en la corrida normal
15. search hit de OpenSearch debe poder resolverse a refs autorizadas y no a datos crudos no permitidos
16. export de evidencia conserva clasificación y trazabilidad
17. borrado lógico de artifact invalida retrieval asociado de forma consistente
18. redacción fallida bloquea indexación cuando el material no es publicable completo
19. documento no soportado genera error clasificable y evidence mínima
20. OpenSearch jamás se usa para reconstruir la verdad operativa del documento

## Criterios de aceptación

1. Queda fijado el boundary exacto entre los cuatro planos.
2. Queda explícito qué vive y qué nunca vive en cada plano.
3. Queda fijado OpenSearch como retrieval/corpus plane y nada más.
4. Queda fijado S3-compatible object storage como artifact/evidence plane persistente.
5. Queda fijado Valkey como plano efímero y no autoritativo.
6. Queda descrita la cadena baseline de ingestión documental.
7. Quedan fijadas reglas de redacción, clasificación y evidence.
8. Existen tests borde suficientes para validar el plane convergente.
