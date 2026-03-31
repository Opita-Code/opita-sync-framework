package memory

import (
	"sync"

	"opita-sync-framework/internal/engine/access"
)

type AccessStore struct {
	mu          sync.RWMutex
	grants      map[string]access.CapabilityGrant
	delegations map[string]access.DelegationGrant
}

func NewAccessStore() *AccessStore {
	return &AccessStore{grants: map[string]access.CapabilityGrant{}, delegations: map[string]access.DelegationGrant{}}
}

func (s *AccessStore) SaveGrant(grant access.CapabilityGrant) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.grants[grant.GrantID] = grant
	return nil
}

func (s *AccessStore) GetGrantByID(grantID string) (access.CapabilityGrant, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	grant, found := s.grants[grantID]
	return grant, found, nil
}

func (s *AccessStore) ListGrantsByTenant(tenantID string) ([]access.CapabilityGrant, error) {
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

func (s *AccessStore) SaveDelegation(grant access.DelegationGrant) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.delegations[grant.GrantID] = grant
	return nil
}

func (s *AccessStore) GetDelegationByID(grantID string) (access.DelegationGrant, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	grant, found := s.delegations[grantID]
	return grant, found, nil
}

func (s *AccessStore) ListDelegationsByTenant(tenantID string) ([]access.DelegationGrant, error) {
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
