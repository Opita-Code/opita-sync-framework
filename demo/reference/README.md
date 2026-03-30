# Reference Demo

## Objetivo

Este demo de referencia muestra el corredor mínimo actual de **Opita Sync Framework** sobre la vertical slice implementada.

No es un demo comercial ni breadth demo. Es un demo **probatorio, acotado y reproducible**.

## Qué demuestra

1. intake gobernado
2. proposal draft
3. patchset candidate
4. preview + simulation
5. compile hacia engine
6. inspection / debug final

## Qué no demuestra

- distribution layer
- rollout tenant-scoped
- apply real de negocio
- UI final
- scale enterprise

## Archivos del demo

- `requests/01-intake-turn.json`
- `requests/02-proposal.json`
- `requests/03-patchset.json`
- `requests/04-preview.json`
- `requests/05-compile-intent.json`
- `demo.http`

## Corredor recomendado

1. `POST /v1/intake/turns`
2. `POST /v1/proposals`
3. `POST /v1/patchsets`
4. `POST /v1/previews`
5. `POST /v1/intents/compile`
6. `GET /v1/foundation/runs/{execution_id}`
7. `GET /v1/debug/semantic?execution_id=...`

## Evidence trail esperado

El demo debe dejar visibles, como mínimo:

- `conversation_turn_id`
- `intake_session_id`
- `intent_candidate_id`
- `proposal_draft_id`
- `patchset_candidate_id`
- `preview_candidate_id`
- `simulation_result_id`
- `contract_id`
- `contract_fingerprint`
- `execution_id`
- `policy_decision_id`

## Criterio de éxito

El demo se considera exitoso si el corredor se puede reconstruir de punta a punta con artifacts reales y evidence trail suficiente, sin explicación verbal como única fuente.
