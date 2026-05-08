package memory

import (
	"context"
	"errors"
	"sync"

	"opita-sync-framework/internal/engine/tenant"
)

var _ tenant.Store = (*TenantStore)(nil)

type TenantStore struct {
	mu      sync.RWMutex
	records map[string]tenant.BootstrapRecord
}

func NewTenantStore() *TenantStore {
	return &TenantStore{records: map[string]tenant.BootstrapRecord{}}
}

func (s *TenantStore) Save(ctx context.Context, record tenant.BootstrapRecord) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if record.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[record.TenantID] = record
	return nil
}

func (s *TenantStore) GetByTenantID(ctx context.Context, tenantID string) (tenant.BootstrapRecord, bool, error) {
	select {
	case <-ctx.Done():
		return tenant.BootstrapRecord{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, found := s.records[tenantID]
	return record, found, nil
}
