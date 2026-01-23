package internal

import (
	"context"
	"sync"

	"github.com/formancehq/terraform-provider-cloud/pkg"
)

// Store provides a shared storage for provider-wide data
type Store struct {
	sync.Mutex
	organizationID string

	tp  pkg.TokenProviderImpl
	sdk pkg.CloudSDK
}

// NewStore creates a new Store instance
func NewStore(sdkClient pkg.CloudSDK, tp pkg.TokenProviderImpl) *Store {
	return &Store{
		sdk: sdkClient,
		tp:  tp,
	}
}

// GetSDK returns the SDK client
func (s *Store) GetSDK() pkg.CloudSDK {
	return s.sdk
}

// GetOrganizationID returns the current organization ID
func (s *Store) GetOrganizationID(ctx context.Context) (string, error) {
	s.Lock()
	defer s.Unlock()
	if s.organizationID == "" {
		orgId, err := s.tp.OrganizationId(ctx)
		if err != nil {
			return "", err
		}
		s.organizationID = orgId
	}
	return s.organizationID, nil
}
