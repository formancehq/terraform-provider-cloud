package resources

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-cloud-go/pkg/models/shared"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/pkg"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &OrganizationMember{}
	_ resource.ResourceWithConfigure = &OrganizationMember{}
)

var SchemaOrganizationMember = schema.Schema{
	Description: "Manages organization members and invitations in Formance Cloud. This resource can be used to invite users to an organization and manage their access levels.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the invitation or membership.",
			Computed:    true,
		},
		"email": schema.StringAttribute{
			Description: "The email address of the user to invite or add to the organization.",
			Required:    true,
		},
		"user_id": schema.StringAttribute{
			Description: "The user ID once the invitation has been accepted.",
			Computed:    true,
		},
	},
}

type OrganizationMemberModel struct {
	ID     types.String `tfsdk:"id"`
	Email  types.String `tfsdk:"email"`
	UserId types.String `tfsdk:"user_id"`
}

func (m *OrganizationMemberModel) GetID() string {
	return m.ID.ValueString()
}

func (m *OrganizationMemberModel) GetEmail() string {
	return m.Email.ValueString()
}

type OrganizationMember struct {
	store *internal.Store
}

func NewOrganizationMember() func() resource.Resource {
	return func() resource.Resource {
		return &OrganizationMember{}
	}
}

// Schema implements resource.Resource.
func (s *OrganizationMember) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaOrganizationMember
}

// Configure implements resource.ResourceWithConfigure.
func (s *OrganizationMember) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	store, ok := req.ProviderData.(*internal.Store)
	if !ok {
		res.Diagnostics.AddError(
			ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *internal.Store, got: %T", req.ProviderData),
		)
		return
	}

	s.store = store
}

// Create implements resource.Resource.
func (s *OrganizationMember) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan OrganizationMemberModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}

	operation, err := s.store.GetSDK().CreateInvitation(ctx, organizationId, plan.GetEmail())
	if err != nil {
		pkg.HandleSDKError(ctx, err, &res.Diagnostics)
		return
	}

	if operation.CreateInvitationResponse == nil || operation.CreateInvitationResponse.Data == nil {
		res.Diagnostics.AddError(
			"Invalid response",
			"CreateInvitation returned an invalid response",
		)
		return
	}

	invitation := operation.CreateInvitationResponse.Data
	plan.ID = types.StringValue(invitation.ID)
	plan.Email = types.StringValue(invitation.UserEmail)
	plan.UserId = types.StringNull()
	if invitation.UserID != nil {
		plan.UserId = types.StringValue(*invitation.UserID)
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *OrganizationMember) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state OrganizationMemberModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}

	operation, err := s.store.GetSDK().ListOrganizationInvitations(ctx, organizationId)
	if err != nil {
		pkg.HandleSDKError(ctx, err, &res.Diagnostics)
		return
	}

	if operation.ListInvitationsResponse == nil {
		res.Diagnostics.AddError(
			"Invalid response",
			"ListOrganizationInvitations returned an invalid response",
		)
		return
	}

	var invitation *shared.Invitation
	for i := range operation.ListInvitationsResponse.Data {
		if operation.ListInvitationsResponse.Data[i].ID == state.GetID() {
			invitation = &operation.ListInvitationsResponse.Data[i]
			break
		}
	}

	if invitation == nil {
		// Invitation not found, might have been deleted already
		return
	}

	switch invitation.Status {
	case shared.InvitationStatusPending:
		operation, err := s.store.GetSDK().DeleteInvitation(ctx, organizationId, state.GetID())
		if err != nil {
			if operation.StatusCode == http.StatusNotFound {
				res.Diagnostics.AddWarning(
					"Invitation not found",
					"The invitation was not found. It may have already been deleted outside of Terraform.",
				)
				return
			}
			pkg.HandleSDKError(ctx, err, &res.Diagnostics)
			return
		}
	case shared.InvitationStatusAccepted:
		operation, err := s.store.GetSDK().DeleteUserOfOrganization(ctx, organizationId, state.UserId.ValueString())
		if err != nil {
			if operation.StatusCode == http.StatusNotFound {
				res.Diagnostics.AddWarning(
					"User not found",
					"The user was not found. They may have already been removed outside of Terraform.",
				)
				return
			}
			pkg.HandleSDKError(ctx, err, &res.Diagnostics)
			return
		}
	}
}

// Metadata implements resource.Resource.
func (s *OrganizationMember) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_member"
}

// Read implements resource.Resource.
func (s *OrganizationMember) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state OrganizationMemberModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}

	operation, err := s.store.GetSDK().ListOrganizationInvitations(ctx, organizationId)
	if err != nil {
		pkg.HandleSDKError(ctx, err, &res.Diagnostics)
		return
	}

	if operation.ListInvitationsResponse == nil {
		res.Diagnostics.AddError(
			"Invalid response",
			"ListOrganizationInvitations returned an invalid response",
		)
		return
	}

	var invitation *shared.Invitation
	for i := range operation.ListInvitationsResponse.Data {
		if operation.ListInvitationsResponse.Data[i].ID == state.GetID() {
			invitation = &operation.ListInvitationsResponse.Data[i]
			break
		}
	}

	if invitation == nil {
		// Invitation not found, mark as removed
		res.State.RemoveResource(ctx)
		return
	}

	state.UserId = types.StringNull()
	if invitation.UserID != nil {
		state.UserId = types.StringValue(*invitation.UserID)
	}
	switch invitation.Status {
	default:
		state.ID = types.StringValue(invitation.ID)
		state.Email = types.StringValue(invitation.UserEmail)
	case shared.InvitationStatusAccepted:
		operation, err := s.store.GetSDK().ReadUserOfOrganization(ctx, organizationId, state.UserId.ValueString())
		if err != nil {
			if operation.StatusCode == http.StatusNotFound {
				res.State.RemoveResource(ctx)
				return
			}
			pkg.HandleSDKError(ctx, err, &res.Diagnostics)
			return
		}

		if operation.ReadOrganizationUserResponse == nil || operation.ReadOrganizationUserResponse.Data == nil {
			res.Diagnostics.AddError(
				"Invalid response",
				"ReadUserOfOrganization returned an invalid response",
			)
			return
		}

		user := operation.ReadOrganizationUserResponse.Data
		state.Email = types.StringValue(user.Email)
		state.UserId = types.StringValue(user.ID)
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *OrganizationMember) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan OrganizationMemberModel
	var state OrganizationMemberModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}
	defer res.Diagnostics.Append(res.State.Set(ctx, &plan)...)

	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}

	operation, err := s.store.GetSDK().ListOrganizationInvitations(ctx, organizationId)
	if err != nil {
		pkg.HandleSDKError(ctx, err, &res.Diagnostics)
		return
	}

	if operation.ListInvitationsResponse == nil {
		res.Diagnostics.AddError(
			"Invalid response",
			"ListOrganizationInvitations returned an invalid response",
		)
		return
	}

	var invitation *shared.Invitation
	for i := range operation.ListInvitationsResponse.Data {
		if operation.ListInvitationsResponse.Data[i].ID == state.GetID() {
			invitation = &operation.ListInvitationsResponse.Data[i]
			break
		}
	}

	if invitation == nil {
		res.Diagnostics.AddError(
			"Invitation not found",
			"The invitation was not found",
		)
		return
	}

	plan.ID = state.ID
	plan.UserId = state.UserId

	switch invitation.Status {
	case shared.InvitationStatusPending:
		// Delete and recreate invitation if email changed
		if invitation.ExpiresAt == nil {
			return
		}
		if time.Now().Before(*invitation.ExpiresAt) {
			return
		}
		deleteoperation, err := s.store.GetSDK().DeleteInvitation(ctx, organizationId, state.GetID())
		if err != nil && deleteoperation.StatusCode != http.StatusNotFound {
			pkg.HandleSDKError(ctx, err, &res.Diagnostics)
			return

		}

		operation, err := s.store.GetSDK().CreateInvitation(ctx, organizationId, plan.GetEmail())
		if err != nil {
			pkg.HandleSDKError(ctx, err, &res.Diagnostics)
			return
		}

		newInvitation := operation.CreateInvitationResponse.Data
		plan.ID = types.StringValue(newInvitation.ID)
		plan.Email = types.StringValue(newInvitation.UserEmail)
		plan.UserId = types.StringNull()
		if newInvitation.UserID != nil {
			plan.UserId = types.StringValue(*newInvitation.UserID)
		}

	case shared.InvitationStatusAccepted:
		// For accepted invitations, email cannot be changed
		// Keep the existing state as email is tied to the user account
		if plan.GetEmail() != state.GetEmail() {
			res.Diagnostics.AddWarning(
				"Email cannot be changed",
				"Email cannot be changed for accepted invitations. The existing email will be preserved.",
			)
		}
		// No changes needed, keep existing state
		plan = state
	}

}
