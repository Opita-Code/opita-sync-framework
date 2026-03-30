package memory

import (
	"context"
	"sync"

	"opita-sync-framework/internal/engine/intent"
)

type ContractRepository struct {
	mu            sync.RWMutex
	byFingerprint map[string]intent.CompiledContract
	byID          map[string]intent.CompiledContract
}

func NewContractRepository() *ContractRepository {
	return &ContractRepository{
		byFingerprint: map[string]intent.CompiledContract{},
		byID:          map[string]intent.CompiledContract{},
	}
}

func (r *ContractRepository) GetByFingerprint(_ context.Context, fingerprint string) (intent.CompiledContract, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	contract, found := r.byFingerprint[fingerprint]
	return contract, found, nil
}

func (r *ContractRepository) Save(_ context.Context, contract intent.CompiledContract) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byFingerprint[contract.Fingerprint] = contract
	r.byID[contract.ContractID] = contract
	return nil
}

func (r *ContractRepository) GetByID(_ context.Context, contractID string) (intent.CompiledContract, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	contract, found := r.byID[contractID]
	return contract, found, nil
}
