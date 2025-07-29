package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &OrganizationMember{}
	_ resource.ResourceWithConfigure = &OrganizationMember{}
)

type OrganizationMember struct {
	logger logging.Logger
	store  *internal.Store
}

type OrganizationMemberModel struct {
	ID types.String `tfsdk:"id"`

	Role types.String `tfsdk:"role"`

	Email  types.String `tfsdk:"email"`
	UserId types.String `tfsdk:"user_id"`
}

type Roles struct {
	Organization types.String `tfsdk:"organization"`
	Stack        types.String `tfsdk:"stack"`
}

func NewOrganizationMember(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &OrganizationMember{
			logger: logger,
		}
	}
}

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
		"role": schema.StringAttribute{
			Description: "The role to assign to the user in the organization. Valid values are: GUEST, ADMIN.",
			Optional:    true,
			Computed:    true,
		},
	},
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
	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "organization_member_create"))
	s.logger.Debug("Creating organization member")

	var plan OrganizationMemberModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	claim := sdk.InvitationClaim{}
	if plan.Role.ValueString() != "" {
		claim.Role = pointer.For(sdk.Role(plan.Role.ValueString()))
	}
	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}
	obj, resp, err := s.store.GetSDK().CreateInvitation(ctx, organizationId, plan.Email.ValueString(), claim)
	if err != nil {
		pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.Role = types.StringValue(string(obj.Data.Role))
	plan.Email = types.StringValue(obj.Data.UserEmail)
	plan.UserId = types.StringNull()
	if obj.Data.UserId != nil {
		plan.UserId = types.StringValue(*obj.Data.UserId)
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...) // Save the plan as state
}

// Delete implements resource.Resource.
func (s *OrganizationMember) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "organization_member_delete"))
	s.logger.Debug("Deleting organization member")

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
	objs, resp, err := s.store.GetSDK().ListOrganizationInvitations(ctx, organizationId)
	if err != nil {
		pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
		return
	}

	obj := collectionutils.First(objs.Data, func(inv sdk.Invitation) bool {
		return inv.Id == state.ID.ValueString()
	})

	switch obj.Status {
	case "PENDING":
		resp, err := s.store.GetSDK().DeleteInvitation(ctx, organizationId, state.ID.ValueString())
		if err != nil {
			pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
			return
		}
	case "ACCEPTED":
		resp, err := s.store.GetSDK().DeleteUserOfOrganization(ctx, organizationId, state.UserId.ValueString())
		if err != nil {
			pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
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
	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "organization_member_read"))
	s.logger.Debug("Reading organization member")

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
	objs, resp, err := s.store.GetSDK().ListOrganizationInvitations(ctx, organizationId)
	if err != nil {
		pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
		return
	}

	obj := collectionutils.First(objs.Data, func(inv sdk.Invitation) bool {
		return inv.Id == state.ID.ValueString()
	})

	switch obj.Status {
	default:
		state.Role = types.StringValue(string(obj.Role))
		state.Email = types.StringValue(obj.UserEmail)
		state.UserId = types.StringNull()
		if obj.UserId != nil {
			state.UserId = types.StringValue(*obj.UserId)
		}
		state.ID = types.StringValue(obj.Id)
	case "ACCEPTED":
		user, resp, err := s.store.GetSDK().ReadUserOfOrganization(ctx, organizationId, state.UserId.ValueString())
		if err != nil {
			pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
			return
		}
		state.Role = types.StringValue(string(user.Data.Role))
		state.Email = types.StringValue(user.Data.Email)
		state.UserId = types.StringValue(user.Data.Id)
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (s *OrganizationMember) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "organization_member_update"))
	s.logger.Debug("Updating organization member")

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
	objs, resp, err := s.store.GetSDK().ListOrganizationInvitations(ctx, organizationId)
	if err != nil {
		pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
		return
	}

	obj := collectionutils.First(objs.Data, func(inv sdk.Invitation) bool {
		return inv.Id == state.ID.ValueString()
	})

	switch obj.Status {
	case "PENDING":
		resp, err := s.store.GetSDK().DeleteInvitation(ctx, organizationId, state.ID.ValueString())
		if err != nil {
			pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
			return
		}

		claim := sdk.InvitationClaim{}
		if state.Role.ValueString() != "" {
			claim.Role = pointer.For(sdk.Role(state.Role.ValueString()))
		}

		obj, respCreate, err := s.store.GetSDK().CreateInvitation(ctx, organizationId, state.Email.ValueString(), claim)
		if err != nil {
			pkg.HandleSDKError(ctx, err, respCreate, &res.Diagnostics)
			return
		}
		state.ID = types.StringValue(obj.Data.Id)
		state.Role = types.StringValue(string(obj.Data.Role))
		state.Email = types.StringValue(obj.Data.UserEmail)
		state.UserId = types.StringNull()
		if obj.Data.UserId != nil {
			state.UserId = types.StringValue(*obj.Data.UserId)
		}
		res.Diagnostics.Append(res.State.Set(ctx, &state)...)
	case "ACCEPTED":
		if state.Role.ValueString() == "" {
			return
		}
		body := sdk.UpdateOrganizationUserRequest{
			Role: sdk.Role(state.Role.ValueString()),
		}
		resp, err := s.store.GetSDK().UpsertUserOfOrganization(ctx, organizationId, state.UserId.ValueString(), body)
		if err != nil {
			pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
			return
		}
		state.Role = types.StringValue(state.Role.ValueString())

	}

}
