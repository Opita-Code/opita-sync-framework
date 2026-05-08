package memory

import (
	"context"
	"sync"

	"opita-sync-framework/internal/engine/access"
)

var _ access.Store = (*AccessStore)(nil)

type AccessStore struct {
	mu          sync.RWMutex
	grants      map[string]access.CapabilityGrant
	delegations map[string]access.DelegationGrant
}

func NewAccessStore() *AccessStore {
	return &AccessStore{grants: map[string]access.CapabilityGrant{}, delegations: map[string]access.DelegationGrant{}}
}

func (s *AccessStore) SaveGrant(ctx context.Context, grant access.CapabilityGrant) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.grants[grant.GrantID] = grant
	return nil
}

func (s *AccessStore) GetGrantByID(ctx context.Context, grantID string) (access.CapabilityGrant, bool, error) {
	select {
	case <-ctx.Done():
		return access.CapabilityGrant{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	grant, found := s.grants[grantID]
	return grant, found, nil
}

func (s *AccessStore) ListGrantsByTenant(ctx context.Context, tenantID string) ([]access.CapabilityGrant, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]access.CapabilityGrant, 0)
	for _, grant := range s.grants {
		if grant.TenantID == tenantID {
			out = append(out, grant)
		}
	}
	return out, nil
}

func (s *AccessStore) SaveDelegation(ctx context.Context, grant access.DelegationGrant) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.delegations[grant.GrantID] = grant
	return nil
}

func (s *AccessStore) GetDelegationByID(ctx context.Context, grantID string) (access.DelegationGrant, bool, error) {
	select {
	case <-ctx.Done():
		return access.DelegationGrant{}, false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	grant, found := s.delegations[grantID]
	return grant, found, nil
}

func (s *AccessStore) ListDelegationsByTenant(ctx context.Context, tenantID string) ([]access.DelegationGrant, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]access.DelegationGrant, 0)
	for _, grant := range s.delegations {
		if grant.TenantID == tenantID {
			out = append(out, grant)
		}
	}
	return out, nil
}
