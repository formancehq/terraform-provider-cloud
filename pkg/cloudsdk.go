package pkg

import (
	"context"
	"net/http"

	"github.com/formancehq/go-libs/v3/pointer"
	formancesdkcloudgo "github.com/formancehq/terraform-provider-cloud/pkg/membership_client"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/operations"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
)

//go:generate mockgen -typed -destination=cloud_generated.go -package=pkg . CloudSDK
type CloudSDK interface {
	CreateStack(ctx context.Context, organizationID string, body *shared.CreateStackRequest) (*operations.CreateStackResponse, error)
	ReadStack(ctx context.Context, organizationID, stackID string) (*operations.GetStackResponse, error)
	UpdateStack(ctx context.Context, organizationID, stackID string, body *shared.StackData) (*operations.UpdateStackResponse, error)
	DeleteStack(ctx context.Context, organizationID, stackID string, force bool) (*operations.DeleteStackResponse, error)

	UpgradeStack(ctx context.Context, organizationID, stackID, version string) (*operations.UpgradeStackResponse, error)
	ListStacks(ctx context.Context, organizationID string) (*operations.ListStacksResponse, error)

	ReadStackUserAccess(ctx context.Context, organizationID, stackID, userId string) (*operations.ReadStackUserAccessResponse, error)
	UpsertStackUserAccess(ctx context.Context, organizationID, stackID string, userId string, body *shared.UpdateStackUserRequest) (*operations.UpsertStackUserAccessResponse, error)
	DeleteStackUserAccess(ctx context.Context, organizationID, stackID string, userId string) (*operations.DeleteStackUserAccessResponse, error)

	EnableModule(ctx context.Context, organizationID, stackID, moduleName string) (*operations.EnableModuleResponse, error)
	DisableModule(ctx context.Context, organizationID, stackID, moduleName string) (*operations.DisableModuleResponse, error)
	ListModules(ctx context.Context, organizationID, stackID string) (*operations.ListModulesResponse, error)

	ListRegions(ctx context.Context, organizationID string) (*operations.ListRegionsResponse, error)
	GetRegion(ctx context.Context, organizationID, regionID string) (*operations.GetRegionResponse, error)
	GetRegionVersions(ctx context.Context, organizationID, regionID string) (*operations.GetRegionVersionsResponse, error)

	ReadOrganization(ctx context.Context, organizationID string) (*operations.ReadOrganizationResponse, error)

	CreateInvitation(ctx context.Context, organizationID, email string) (*operations.CreateInvitationResponse, error)
	DeleteInvitation(ctx context.Context, organizationID, invitationID string) (*operations.DeleteInvitationResponse, error)
	ListOrganizationInvitations(ctx context.Context, organizationID string) (*operations.ListInvitationsResponse, error)

	ReadUserOfOrganization(ctx context.Context, organizationID, userID string) (*operations.ReadUserOfOrganizationResponse, error)
	DeleteUserOfOrganization(ctx context.Context, organizationID, userID string) (*operations.DeleteUserFromOrganizationResponse, error)
	UpsertUserOfOrganization(ctx context.Context, organizationID string, userID string, body *shared.UpdateOrganizationUserRequest) (*operations.UpsertOrganizationUserResponse, error)
}

var _ CloudSDK = &sdkImpl{}

type sdkImpl struct {
	sdk *formancesdkcloudgo.FormanceCloud
}

func (s *sdkImpl) ReadStack(ctx context.Context, organizationID string, stackID string) (*operations.GetStackResponse, error) {
	return s.sdk.GetStack(ctx, organizationID, stackID)
}

func (s *sdkImpl) CreateStack(ctx context.Context, organizationID string, body *shared.CreateStackRequest) (*operations.CreateStackResponse, error) {
	return s.sdk.CreateStack(ctx, organizationID, body)
}

func (s *sdkImpl) UpdateStack(ctx context.Context, organizationID string, stackID string, body *shared.StackData) (*operations.UpdateStackResponse, error) {
	return s.sdk.UpdateStack(ctx, organizationID, stackID, body)
}

func (s *sdkImpl) DeleteStack(ctx context.Context, organizationID, stackID string, force bool) (*operations.DeleteStackResponse, error) {
	var forcePtr *bool
	if force {
		forcePtr = pointer.For(true)
	}
	return s.sdk.DeleteStack(ctx, organizationID, stackID, forcePtr)
}

func (s *sdkImpl) ListStacks(ctx context.Context, organizationID string) (*operations.ListStacksResponse, error) {
	return s.sdk.ListStacks(ctx, organizationID, nil, nil)
}

func (s *sdkImpl) UpgradeStack(ctx context.Context, organizationID, stackID string, version string) (*operations.UpgradeStackResponse, error) {
	return s.sdk.UpgradeStack(ctx, organizationID, stackID, &shared.StackVersion{
		Version: pointer.For(version),
	})
}

func (s *sdkImpl) ListModules(ctx context.Context, organizationID string, stackID string) (*operations.ListModulesResponse, error) {
	return s.sdk.ListModules(ctx, organizationID, stackID)
}

func (s *sdkImpl) EnableModule(ctx context.Context, organizationID, stackID, moduleName string) (*operations.EnableModuleResponse, error) {
	return s.sdk.EnableModule(ctx, organizationID, stackID, moduleName)
}

func (s *sdkImpl) DisableModule(ctx context.Context, organizationID, stackID, moduleName string) (*operations.DisableModuleResponse, error) {
	return s.sdk.DisableModule(ctx, organizationID, stackID, moduleName)
}

func (s *sdkImpl) UpsertStackUserAccess(ctx context.Context, organizationID, stackID, userId string, body *shared.UpdateStackUserRequest) (*operations.UpsertStackUserAccessResponse, error) {
	return s.sdk.UpsertStackUserAccess(ctx, organizationID, stackID, userId, body)
}

func (s *sdkImpl) ReadStackUserAccess(ctx context.Context, organizationID, stackID, userId string) (*operations.ReadStackUserAccessResponse, error) {
	return s.sdk.ReadStackUserAccess(ctx, organizationID, stackID, userId)
}

func (s *sdkImpl) DeleteStackUserAccess(ctx context.Context, organizationID, stackID, userId string) (*operations.DeleteStackUserAccessResponse, error) {
	return s.sdk.DeleteStackUserAccess(ctx, organizationID, stackID, userId)
}

func (s *sdkImpl) ListRegions(ctx context.Context, organizationID string) (*operations.ListRegionsResponse, error) {
	return s.sdk.ListRegions(ctx, organizationID)
}

func (s *sdkImpl) GetRegion(ctx context.Context, organizationID, regionID string) (*operations.GetRegionResponse, error) {
	return s.sdk.GetRegion(ctx, organizationID, regionID)
}

func (s *sdkImpl) GetRegionVersions(ctx context.Context, organizationID, regionID string) (*operations.GetRegionVersionsResponse, error) {
	return s.sdk.GetRegionVersions(ctx, organizationID, regionID)
}

func (s *sdkImpl) ReadOrganization(ctx context.Context, organizationID string) (*operations.ReadOrganizationResponse, error) {
	return s.sdk.ReadOrganization(ctx, organizationID, nil)
}

func (s *sdkImpl) CreateInvitation(ctx context.Context, organizationID, email string) (*operations.CreateInvitationResponse, error) {
	return s.sdk.CreateInvitation(ctx, organizationID, email)
}

func (s *sdkImpl) DeleteInvitation(ctx context.Context, organizationID, invitationID string) (*operations.DeleteInvitationResponse, error) {
	return s.sdk.DeleteInvitation(ctx, organizationID, invitationID)
}

func (s *sdkImpl) ListOrganizationInvitations(ctx context.Context, organizationID string) (*operations.ListInvitationsResponse, error) {
	orgPtr := pointer.For(organizationID)
	return s.sdk.ListInvitations(ctx, nil, orgPtr)
}

func (s *sdkImpl) ReadUserOfOrganization(ctx context.Context, organizationID, userID string) (*operations.ReadUserOfOrganizationResponse, error) {
	return s.sdk.ReadUserOfOrganization(ctx, organizationID, userID)
}

func (s *sdkImpl) DeleteUserOfOrganization(ctx context.Context, organizationID, userID string) (*operations.DeleteUserFromOrganizationResponse, error) {
	return s.sdk.DeleteUserFromOrganization(ctx, organizationID, userID)
}

func (s *sdkImpl) UpsertUserOfOrganization(ctx context.Context, organizationID string, userID string, body *shared.UpdateOrganizationUserRequest) (*operations.UpsertOrganizationUserResponse, error) {
	return s.sdk.UpsertOrganizationUser(ctx, organizationID, userID, body)
}

type CloudFactory func(creds Creds, transport http.RoundTripper) CloudSDK

func NewCloudSDK(creds Creds, transport http.RoundTripper) CloudSDK {
	tp := NewTokenProvider(transport, creds)
	sdk := NewSDK(creds.Endpoint(), transport, tp)
	return &sdkImpl{
		sdk: sdk,
	}
}

func NewSDK(endpoint string, transport http.RoundTripper, tp TokenProviderImpl) *formancesdkcloudgo.FormanceCloud {
	client := &http.Client{
		Transport: transport,
	}

	return formancesdkcloudgo.New(
		formancesdkcloudgo.WithServerURL(endpoint),
		formancesdkcloudgo.WithClient(client),
		formancesdkcloudgo.WithSecuritySource(func(ctx context.Context) (shared.Security, error) {
			token, err := tp.RefreshToken(ctx)
			if err != nil {
				return shared.Security{}, err
			}
			return shared.Security{
				Oauth2: token.AccessToken,
			}, nil
		}),
	)
}
