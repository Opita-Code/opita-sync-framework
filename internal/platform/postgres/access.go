package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/access"
)

var _ access.Store = (*AccessStore)(nil)

type AccessStore struct {
	store *Store
}

func NewAccessStore(store *Store) *AccessStore { return &AccessStore{store: store} }

func (s *AccessStore) SaveGrant(ctx context.Context, grant access.CapabilityGrant) error {
	raw, err := json.Marshal(grant)
	if err != nil {
		return fmt.Errorf("marshal capability grant: %w", err)
	}
	_, err = s.store.DB.ExecContext(ctx, `
		insert into tenant_capability_grants (grant_id, tenant_id, state, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6)
		on conflict (grant_id) do update set state = excluded.state, payload = excluded.payload, updated_at = excluded.updated_at
	`, grant.GrantID, grant.TenantID, grant.State, raw, grant.CreatedAt, grant.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert capability grant: %w", err)
	}
	return nil
}

func (s *AccessStore) GetGrantByID(ctx context.Context, grantID string) (access.CapabilityGrant, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(ctx, `select payload from tenant_capability_grants where grant_id = $1`, grantID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return access.CapabilityGrant{}, false, nil
		}
		return access.CapabilityGrant{}, false, fmt.Errorf("select capability grant: %w", err)
	}
	var grant access.CapabilityGrant
	if err := json.Unmarshal(raw, &grant); err != nil {
		return access.CapabilityGrant{}, false, fmt.Errorf("unmarshal capability grant: %w", err)
	}
	return grant, true, nil
}

func (s *AccessStore) ListGrantsByTenant(ctx context.Context, tenantID string) ([]access.CapabilityGrant, error) {
	rows, err := s.store.DB.QueryContext(ctx, `select payload from tenant_capability_grants where tenant_id = $1 order by created_at asc`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list capability grants: %w", err)
	}
	defer rows.Close()
	out := make([]access.CapabilityGrant, 0)
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("scan capability grant: %w", err)
		}
		var grant access.CapabilityGrant
		if err := json.Unmarshal(raw, &grant); err != nil {
			return nil, fmt.Errorf("unmarshal capability grant: %w", err)
		}
		out = append(out, grant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate capability grants: %w", err)
	}
	return out, nil
}

func (s *AccessStore) SaveDelegation(ctx context.Context, grant access.DelegationGrant) error {
	raw, err := json.Marshal(grant)
	if err != nil {
		return fmt.Errorf("marshal delegation grant: %w", err)
	}
	_, err = s.store.DB.ExecContext(ctx, `
		insert into tenant_delegation_grants (grant_id, tenant_id, state, payload, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6)
		on conflict (grant_id) do update set state = excluded.state, payload = excluded.payload, updated_at = excluded.updated_at
	`, grant.GrantID, grant.TenantID, grant.State, raw, grant.CreatedAt, grant.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert delegation grant: %w", err)
	}
	return nil
}

func (s *AccessStore) GetDelegationByID(ctx context.Context, grantID string) (access.DelegationGrant, bool, error) {
	var raw []byte
	err := s.store.DB.QueryRowContext(ctx, `select payload from tenant_delegation_grants where grant_id = $1`, grantID).Scan(&raw)
	if err != nil {
		if isNoRows(err) {
			return access.DelegationGrant{}, false, nil
		}
		return access.DelegationGrant{}, false, fmt.Errorf("select delegation grant: %w", err)
	}
	var grant access.DelegationGrant
	if err := json.Unmarshal(raw, &grant); err != nil {
		return access.DelegationGrant{}, false, fmt.Errorf("unmarshal delegation grant: %w", err)
	}
	return grant, true, nil
}

func (s *AccessStore) ListDelegationsByTenant(ctx context.Context, tenantID string) ([]access.DelegationGrant, error) {
	rows, err := s.store.DB.QueryContext(ctx, `select payload from tenant_delegation_grants where tenant_id = $1 order by created_at asc`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list delegation grants: %w", err)
	}
	defer rows.Close()
	out := make([]access.DelegationGrant, 0)
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("scan delegation grant: %w", err)
		}
		var grant access.DelegationGrant
		if err := json.Unmarshal(raw, &grant); err != nil {
			return nil, fmt.Errorf("unmarshal delegation grant: %w", err)
		}
		out = append(out, grant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate delegation grants: %w", err)
	}
	return out, nil
}
