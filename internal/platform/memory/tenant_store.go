package memory

import (
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/tenant"
)

type TenantStore struct {
	mu      sync.RWMutex
	records map[string]tenant.BootstrapRecord
}

func NewTenantStore() *TenantStore {
	return &TenantStore{records: map[string]tenant.BootstrapRecord{}}
}

func (s *TenantStore) Save(record tenant.BootstrapRecord) error {
	if record.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[record.TenantID] = record
	return nil
}

func (s *TenantStore) GetByTenantID(tenantID string) (tenant.BootstrapRecord, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, found := s.records[tenantID]
	return record, found, nil
}
