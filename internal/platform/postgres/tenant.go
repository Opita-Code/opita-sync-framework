package postgres

import (
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/tenant"
)

type TenantStore struct {
	store *Store
}

func NewTenantStore(store *Store) *TenantStore {
	return &TenantStore{store: store}
}

func (s *TenantStore) Save(record tenant.BootstrapRecord) error {
	raw, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal tenant bootstrap: %w", err)
	}
	_, err = s.store.DB.ExecContext(contextBackground(), `
		insert into tenant_bootstrap_records (tenant_id, state, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
		on conflict (tenant_id) do update set state = excluded.state, payload = excluded.payload, updated_at = excluded.updated_at
	`, record.TenantID, record.State, raw, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert tenant bootstrap: %w", err)
	}
	return nil
}

func (s *TenantStore) GetByTenantID(tenantID string) (tenant.BootstrapRecord, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(contextBackground(), `select payload from tenant_bootstrap_records where tenant_id = $1`, tenantID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return tenant.BootstrapRecord{}, false, nil
		}
		return tenant.BootstrapRecord{}, false, fmt.Errorf("select tenant bootstrap: %w", err)
	}
	var record tenant.BootstrapRecord
	if err := json.Unmarshal(raw, &record); err != nil {
		return tenant.BootstrapRecord{}, false, fmt.Errorf("unmarshal tenant bootstrap: %w", err)
	}
	return record, true, nil
}
