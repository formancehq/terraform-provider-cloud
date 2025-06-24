package pkg

import (
	"context"
	"net/http"

	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/sdk"
)

//go:generate mockgen -destination=cloud_generated.go -package=pkg . CloudSDK
type CloudSDK interface {
	CreateStack(ctx context.Context, organizationID string, body sdk.CreateStackRequest) (*sdk.CreateStackResponse, *http.Response, error)
	ReadStack(ctx context.Context, organizationID, stackID string) (*sdk.CreateStackResponse, *http.Response, error)
	UpdateStack(ctx context.Context, organizationID, stackID string, body sdk.UpdateStackRequest) (*sdk.CreateStackResponse, *http.Response, error)
	DeleteStack(ctx context.Context, organizationID, stackID string, force bool) (*http.Response, error)

	UpgradeStack(ctx context.Context, organizationID, stackID, version string) (*http.Response, error)
	ListStacks(ctx context.Context, organizationID string) (*sdk.ListStacksResponse, *http.Response, error)

	ReadStackUserAccess(ctx context.Context, organizationID, stackID, userId string) (*sdk.ReadStackUserAccess, *http.Response, error)
	UpsertStackUserAccess(ctx context.Context, organizationID, stackID string, userId string, body sdk.UpdateStackUserRequest) (*http.Response, error)
	DeleteStackUserAccess(ctx context.Context, organizationID, stackID string, userId string) (*http.Response, error)

	EnableModule(ctx context.Context, organizationID, stackID, moduleName string) (*http.Response, error)
	DisableModule(ctx context.Context, organizationID, stackID, moduleName string) (*http.Response, error)
	ListModules(ctx context.Context, organizationID, stackID string) (*sdk.ListModulesResponse, *http.Response, error)

	ListRegions(ctx context.Context, organizationID string) (*sdk.ListRegionsResponse, *http.Response, error)
	GetRegion(ctx context.Context, organizationID, regionID string) (*sdk.GetRegionResponse, *http.Response, error)
	GetRegionVersions(ctx context.Context, organizationID, regionID string) (*sdk.GetRegionVersionsResponse, *http.Response, error)

	ReadOrganization(ctx context.Context, organizationID string) (*sdk.ReadOrganizationResponse, *http.Response, error)

	CreateInvitation(ctx context.Context, organizationID, email string, invitationClaim sdk.InvitationClaim) (*sdk.CreateInvitationResponse, *http.Response, error)
	DeleteInvitation(ctx context.Context, organizationID, invitationID string) (*http.Response, error)
	ListOrganizationInvitations(ctx context.Context, organizationID string) (*sdk.ListInvitationsResponse, *http.Response, error)

	ReadUserOfOrganization(ctx context.Context, organizationID, userID string) (*sdk.ReadOrganizationUserResponse, *http.Response, error)
	DeleteUserOfOrganization(ctx context.Context, organizationID, userID string) (*http.Response, error)
	UpsertUserOfOrganization(ctx context.Context, organizationID string, userID string, body sdk.UpdateOrganizationUserRequest) (*http.Response, error)
}

var _ CloudSDK = &sdkImpl{}

type sdkImpl struct {
	sdk sdk.DefaultAPI
}

func (s *sdkImpl) ReadStack(ctx context.Context, organizationID string, stackID string) (*sdk.CreateStackResponse, *http.Response, error) {
	return s.sdk.GetStack(ctx, organizationID, stackID).Execute()
}

func (s *sdkImpl) CreateStack(ctx context.Context, organizationID string, body sdk.CreateStackRequest) (*sdk.CreateStackResponse, *http.Response, error) {
	return s.sdk.CreateStack(ctx, organizationID).CreateStackRequest(body).Execute()
}

func (s *sdkImpl) UpdateStack(ctx context.Context, organizationID string, stackID string, body sdk.UpdateStackRequest) (*sdk.CreateStackResponse, *http.Response, error) {
	return s.sdk.UpdateStack(ctx, organizationID, stackID).UpdateStackRequest(body).Execute()
}

func (s *sdkImpl) DeleteStack(ctx context.Context, organizationID, stackID string, force bool) (*http.Response, error) {
	if force {
		return s.sdk.DeleteStack(ctx, organizationID, stackID).Force(true).Execute()
	}

	return s.sdk.DeleteStack(ctx, organizationID, stackID).Execute()
}

func (s *sdkImpl) ListStacks(ctx context.Context, organizationID string) (*sdk.ListStacksResponse, *http.Response, error) {
	return s.sdk.ListStacks(ctx, organizationID).Execute()
}

func (s *sdkImpl) UpgradeStack(ctx context.Context, organizationID, stackID string, version string) (*http.Response, error) {
	return s.sdk.UpgradeStack(ctx, organizationID, stackID).StackVersion(sdk.StackVersion{
		Version: pointer.For(version),
	}).Execute()
}

func (s *sdkImpl) ListModules(ctx context.Context, organizationID string, stackID string) (*sdk.ListModulesResponse, *http.Response, error) {
	return s.sdk.ListModules(ctx, organizationID, stackID).Execute()
}

func (s *sdkImpl) EnableModule(ctx context.Context, organizationID, stackID, moduleName string) (*http.Response, error) {
	return s.sdk.EnableModule(ctx, organizationID, stackID).Name(moduleName).Execute()
}

func (s *sdkImpl) DisableModule(ctx context.Context, organizationID, stackID, moduleName string) (*http.Response, error) {
	return s.sdk.DisableModule(ctx, organizationID, stackID).Name(moduleName).Execute()
}

func (s *sdkImpl) UpsertStackUserAccess(ctx context.Context, organizationID, stackID, userId string, body sdk.UpdateStackUserRequest) (*http.Response, error) {
	return s.sdk.UpsertStackUserAccess(ctx, organizationID, stackID, userId).UpdateStackUserRequest(body).Execute()
}

func (s *sdkImpl) ReadStackUserAccess(ctx context.Context, organizationID, stackID, userId string) (*sdk.ReadStackUserAccess, *http.Response, error) {
	return s.sdk.ReadStackUserAccess(ctx, organizationID, stackID, userId).Execute()
}

func (s *sdkImpl) DeleteStackUserAccess(ctx context.Context, organizationID, stackID, userId string) (*http.Response, error) {
	return s.sdk.DeleteStackUserAccess(ctx, organizationID, stackID, userId).Execute()
}

func (s *sdkImpl) ListRegions(ctx context.Context, organizationID string) (*sdk.ListRegionsResponse, *http.Response, error) {
	return s.sdk.ListRegions(ctx, organizationID).Execute()
}

func (s *sdkImpl) GetRegion(ctx context.Context, organizationID, regionID string) (*sdk.GetRegionResponse, *http.Response, error) {
	return s.sdk.GetRegion(ctx, organizationID, regionID).Execute()
}

func (s *sdkImpl) GetRegionVersions(ctx context.Context, organizationID, regionID string) (*sdk.GetRegionVersionsResponse, *http.Response, error) {
	return s.sdk.GetRegionVersions(ctx, organizationID, regionID).Execute()
}

func (s *sdkImpl) ReadOrganization(ctx context.Context, organizationID string) (*sdk.ReadOrganizationResponse, *http.Response, error) {
	return s.sdk.ReadOrganization(ctx, organizationID).Execute()
}

func (s *sdkImpl) CreateInvitation(ctx context.Context, organizationID, email string, invitationClaim sdk.InvitationClaim) (*sdk.CreateInvitationResponse, *http.Response, error) {
	if invitationClaim.Role != nil {
		return s.sdk.CreateInvitation(ctx, organizationID).Email(email).InvitationClaim(invitationClaim).Execute()
	}

	return s.sdk.CreateInvitation(ctx, organizationID).Email(email).Execute()
}

func (s *sdkImpl) DeleteInvitation(ctx context.Context, organizationID, invitationID string) (*http.Response, error) {
	return s.sdk.DeleteInvitation(ctx, organizationID, invitationID).Execute()
}

func (s *sdkImpl) ListOrganizationInvitations(ctx context.Context, organizationID string) (*sdk.ListInvitationsResponse, *http.Response, error) {
	return s.sdk.ListOrganizationInvitations(ctx, organizationID).Execute()
}

func (s *sdkImpl) ReadUserOfOrganization(ctx context.Context, organizationID, userID string) (*sdk.ReadOrganizationUserResponse, *http.Response, error) {
	return s.sdk.ReadUserOfOrganization(ctx, organizationID, userID).Execute()
}

func (s *sdkImpl) DeleteUserOfOrganization(ctx context.Context, organizationID, userID string) (*http.Response, error) {
	return s.sdk.DeleteUserFromOrganization(ctx, organizationID, userID).Execute()
}

func (s *sdkImpl) UpsertUserOfOrganization(ctx context.Context, organizationID string, userID string, body sdk.UpdateOrganizationUserRequest) (*http.Response, error) {
	return s.sdk.UpsertOrganizationUser(ctx, organizationID, userID).
		UpdateOrganizationUserRequest(body).
		Execute()
}

type CloudFactory func(creds Creds, transport http.RoundTripper) CloudSDK

func NewCloudSDK() CloudFactory {
	return func(creds Creds, transport http.RoundTripper) CloudSDK {
		return &sdkImpl{sdk: NewSDK(creds, transport)}
	}
}
