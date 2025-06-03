package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &OrganizationMember{}
	_ resource.ResourceWithConfigure        = &OrganizationMember{}
	_ resource.ResourceWithConfigValidators = &OrganizationMember{}
	_ resource.ResourceWithValidateConfig   = &OrganizationMember{}
)

type OrganizationMember struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

type OrganizationMemberModel struct {
	ID types.String `tfsdk:"id"`

	Role types.String `tfsdk:"role"`

	Email          types.String `tfsdk:"email"`
	OrganizationId types.String `tfsdk:"organization_id"`
	UserId         types.String `tfsdk:"user_id"`
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
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"organization_id": schema.StringAttribute{
			Required: true,
		},
		"email": schema.StringAttribute{
			Required: true,
		},
		"user_id": schema.StringAttribute{
			Computed: true,
		},
		"role": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	},
}

// Schema implements resource.Resource.
func (s *OrganizationMember) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaOrganizationMember
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *OrganizationMember) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var config OrganizationMemberModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.OrganizationId.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Invalid Organization ID",
			"The organization_id attribute must not be null.",
		)
	}

	if config.Email.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("email"),
			"Invalid Email",
			"The email attribute must not be null.",
		)
	}
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *OrganizationMember) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return nil
}

// Configure implements resource.ResourceWithConfigure.
func (s *OrganizationMember) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(sdk.DefaultAPI)
	if !ok {
		res.Diagnostics.AddError(
			ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *FormanceCloudProviderModel, got: %T", req.ProviderData),
		)
		return
	}

	s.sdk = sdk
}

// Create implements resource.Resource.
func (s *OrganizationMember) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan OrganizationMemberModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	// Call the SDK method to create the resource here
	sdkReq := s.sdk.CreateInvitation(ctx, plan.OrganizationId.ValueString()).
		Email(plan.Email.ValueString())

	if plan.Role.ValueString() != "" {
		sdkReq = sdkReq.InvitationClaim(sdk.InvitationClaim{
			Role: pointer.For(sdk.Role(plan.Role.ValueString())),
		})
	}

	obj, resp, err := sdkReq.Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.OrganizationId = types.StringValue(obj.Data.OrganizationId)
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
	var state OrganizationMemberModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	objs, resp, err := s.sdk.ListOrganizationInvitations(ctx, state.OrganizationId.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
		return
	}

	obj := collectionutils.First(objs.Data, func(inv sdk.Invitation) bool {
		return inv.Id == state.ID.ValueString()
	})

	switch obj.Status {
	case "PENDING":
		resp, err := s.sdk.DeleteInvitation(ctx, state.OrganizationId.ValueString(), state.ID.ValueString()).Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
			return
		}
	case "ACCEPTED":
		resp, err := s.sdk.DeleteUserFromOrganization(ctx, state.OrganizationId.ValueString(), state.UserId.ValueString()).Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
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

	objs, resp, err := s.sdk.ListOrganizationInvitations(ctx, state.OrganizationId.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
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
		user, resp, err := s.sdk.ReadUserOfOrganization(ctx, state.OrganizationId.ValueString(), state.UserId.ValueString()).Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
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
	var state OrganizationMemberModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	objs, resp, err := s.sdk.ListOrganizationInvitations(ctx, state.OrganizationId.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
		return
	}

	obj := collectionutils.First(objs.Data, func(inv sdk.Invitation) bool {
		return inv.Id == state.ID.ValueString()
	})

	switch obj.Status {
	case "PENDING":
		resp, err := s.sdk.DeleteInvitation(ctx, state.OrganizationId.ValueString(), state.ID.ValueString()).Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
			return
		}

		sdkReq := s.sdk.CreateInvitation(ctx, state.OrganizationId.ValueString()).
			Email(state.Email.ValueString())
		if state.Role.ValueString() != "" {
			sdkReq = sdkReq.InvitationClaim(sdk.InvitationClaim{
				Role: pointer.For(sdk.Role(state.Role.ValueString())),
			})
		}
		obj, respCreate, err := sdkReq.Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, respCreate, &res.Diagnostics)
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
		resp, err := s.sdk.UpsertOrganizationUser(ctx, state.OrganizationId.ValueString(), state.UserId.ValueString()).UpdateOrganizationUserRequest(body).Execute()
		if err != nil {
			pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
			return
		}
		state.Role = types.StringValue(state.Role.ValueString())

	}

}
