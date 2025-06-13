package pkg

import (
	"context"
	"errors"
	"sync"
	
	"github.com/formancehq/terraform-provider-cloud/sdk"
)

var (
	// ErrNoOrganization is returned when no organization is found for the authenticated user
	ErrNoOrganization = errors.New("no organization found")
)

// Store provides a shared storage for provider-wide data
type Store struct {
	sdk sdk.DefaultAPI
	mu  sync.RWMutex
	
	// organizationID caches the current organization ID
	organizationID string
}

// NewStore creates a new Store instance
func NewStore(sdkClient sdk.DefaultAPI) *Store {
	return &Store{
		sdk: sdkClient,
	}
}

// GetSDK returns the SDK client
func (s *Store) GetSDK() sdk.DefaultAPI {
	return s.sdk
}

// SetOrganizationID sets the current organization ID
func (s *Store) SetOrganizationID(orgID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.organizationID = orgID
}

// GetOrganizationID returns the current organization ID
func (s *Store) GetOrganizationID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.organizationID
}

// FetchAndSetCurrentOrganization fetches the current organization and sets it in the store
func (s *Store) FetchAndSetCurrentOrganization(ctx context.Context) (string, error) {
	// List all organizations for the authenticated user
	orgsResp, _, err := s.sdk.ListOrganizationsExpanded(ctx).Execute()
	if err != nil {
		return "", err
	}
	
	if len(orgsResp.Data) == 0 {
		return "", ErrNoOrganization
	}
	
	// Use the first organization as the "current" organization
	orgID := orgsResp.Data[0].Id
	s.SetOrganizationID(orgID)
	
	return orgID, nil
}