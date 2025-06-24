package internal

import (
	"strings"

	"github.com/formancehq/terraform-provider-cloud/sdk"
)

// Store provides a shared storage for provider-wide data
type Store struct {
	clientID string
	sdk      sdk.DefaultAPI
}

// NewStore creates a new Store instance
func NewStore(sdkClient sdk.DefaultAPI, clientID string) *Store {
	return &Store{
		sdk:      sdkClient,
		clientID: clientID,
	}
}

// GetSDK returns the SDK client
func (s *Store) GetSDK() sdk.DefaultAPI {
	return s.sdk
}

// GetOrganizationID returns the current organization ID
func (s *Store) GetOrganizationID() string {
	return strings.Split(s.clientID, "organization_")[1]
}
