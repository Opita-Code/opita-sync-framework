package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/intent"
)

type ContractRepository struct {
	store *Store
}

func NewContractRepository(store *Store) *ContractRepository {
	return &ContractRepository{store: store}
}

func (r *ContractRepository) GetByFingerprint(ctx context.Context, fingerprint string) (intent.CompiledContract, bool, error) {
	var raw []byte
	err := r.store.DB.QueryRowContext(ctx, `select payload from compiled_contracts where fingerprint = $1`, fingerprint).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return intent.CompiledContract{}, false, nil
		}
		return intent.CompiledContract{}, false, fmt.Errorf("select contract by fingerprint: %w", err)
	}
	var contract intent.CompiledContract
	if err := json.Unmarshal(raw, &contract); err != nil {
		return intent.CompiledContract{}, false, fmt.Errorf("unmarshal contract payload: %w", err)
	}
	return contract, true, nil
}

func (r *ContractRepository) Save(ctx context.Context, contract intent.CompiledContract) error {
	raw, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("marshal contract payload: %w", err)
	}
	_, err = r.store.DB.ExecContext(ctx, `
		insert into compiled_contracts (contract_id, fingerprint, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
		on conflict (contract_id) do update set
		  fingerprint = excluded.fingerprint,
		  payload = excluded.payload,
		  updated_at = excluded.updated_at
	`, contract.ContractID, contract.Fingerprint, raw, contract.CreatedAt, contract.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert compiled contract: %w", err)
	}
	return nil
}

func (r *ContractRepository) GetByID(ctx context.Context, contractID string) (intent.CompiledContract, bool, error) {
	var raw []byte
	err := r.store.DB.QueryRowContext(ctx, `select payload from compiled_contracts where contract_id = $1`, contractID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return intent.CompiledContract{}, false, nil
		}
		return intent.CompiledContract{}, false, fmt.Errorf("select contract by id: %w", err)
	}
	var contract intent.CompiledContract
	if err := json.Unmarshal(raw, &contract); err != nil {
		return intent.CompiledContract{}, false, fmt.Errorf("unmarshal contract payload: %w", err)
	}
	return contract, true, nil
}
