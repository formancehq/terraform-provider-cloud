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

var (
	claimOrganizationID = "organization_id"
)

// GetOrganizationID returns the current organization ID
func (s *Store) GetOrganizationID(ctx context.Context) string {
	s.Lock()
	defer s.Unlock()
	if s.organizationID == "" {
		introspection, err := s.tp.IntrospectToken(ctx)
		if err != nil {
			panic(err)
		}
		s.organizationID = introspection.Claims[claimOrganizationID].(string)
	}
	return s.organizationID
}
