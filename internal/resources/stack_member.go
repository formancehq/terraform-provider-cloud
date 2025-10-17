package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &StackMember{}
	_ resource.ResourceWithConfigure = &StackMember{}
)

var SchemaStackMember = schema.Schema{
	Description: "Manages stack members and their access levels in Formance Cloud. This resource allows you to grant users access to specific stacks within an organization.",
	Attributes: map[string]schema.Attribute{
		"user_id": schema.StringAttribute{
			Required:    true,
			Description: "The ID of the user to grant access to the stack. The user must already be a member of the organization.",
		},
		"stack_id": schema.StringAttribute{
			Required:    true,
			Description: "The ID of the stack where the user will be granted access.",
		},
		"role": schema.StringAttribute{
			Required:    true,
			Description: "The role to assign to the user for this stack. Valid values are: GUEST, ADMIN.",
		},
	},
}

type StackMember struct {
	store *internal.Store
}

type StackMemberModel struct {
	Role    types.String `tfsdk:"role"`
	UserId  types.String `tfsdk:"user_id"`
	StackId types.String `tfsdk:"stack_id"`
}

func NewStackMember() func() resource.Resource {
	return func() resource.Resource {
		return &StackMember{}
	}
}

// Configure implements resource.ResourceWithConfigure.
func (s *StackMember) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (s *StackMember) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan StackMemberModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	body := sdk.UpdateStackUserRequest{
		Role: sdk.Role(plan.Role.ValueString()),
	}
	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}
	resp, err := s.store.GetSDK().UpsertStackUserAccess(ctx, organizationId, plan.StackId.ValueString(), plan.UserId.ValueString(), body)
	if err != nil {
		pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *StackMember) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state StackMemberModel
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
	resp, err := s.store.GetSDK().DeleteStackUserAccess(ctx, organizationId, state.StackId.ValueString(), state.UserId.ValueString())
	if err != nil {
		pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
		return
	}
}

// Update implements resource.Resource.
func (s *StackMember) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan StackMemberModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	body := sdk.UpdateStackUserRequest{}
	if r := plan.Role.ValueString(); r != "" {
		body.Role = sdk.Role(r)
	}

	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}

	resp, err := s.store.GetSDK().UpsertStackUserAccess(ctx, organizationId, plan.StackId.ValueString(), plan.UserId.ValueString(), body)
	if err != nil {
		pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Metadata implements resource.Resource.
func (s *StackMember) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack_member"
}

// Read implements resource.Resource.
func (s *StackMember) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state StackMemberModel
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
	userAccess, resp, err := s.store.GetSDK().ReadStackUserAccess(ctx, organizationId, state.StackId.ValueString(), state.UserId.ValueString())
	if err != nil {
		pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
		return
	}

	state.Role = types.StringValue(string(userAccess.Data.Role))

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Schema implements resource.Resource.
func (s *StackMember) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaStackMember
}
