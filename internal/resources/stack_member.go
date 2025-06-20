package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &StackMember{}
	_ resource.ResourceWithConfigure        = &StackMember{}
	_ resource.ResourceWithConfigValidators = &StackMember{}
	_ resource.ResourceWithValidateConfig   = &StackMember{}
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
		"organization_id": schema.StringAttribute{
			Required:    true,
			Description: "The organization ID that owns the stack.",
		},
		"role": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The role to assign to the user for this stack. Valid values are: NONE, READ, WRITE.",
		},
	},
}

type StackMember struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

type StackMemberModel struct {
	Role           types.String `tfsdk:"role"`
	UserId         types.String `tfsdk:"user_id"`
	StackId        types.String `tfsdk:"stack_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
}

func NewStackMember(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &StackMember{
			logger: logger,
		}
	}
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *StackMember) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var config StackMemberModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.OrganizationId.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Invalid organization ID",
			"The organization_id attribute must not be null.",
		)
	}

	if config.StackId.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("stack_id"),
			"Invalid stack ID",
			"The stack_id attribute must not be null.",
		)
	}
	if config.UserId.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("user_id"),
			"Invalid user ID",
			"The user_id attribute must not be null.",
		)
	}
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *StackMember) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return nil
}

// Configure implements resource.ResourceWithConfigure.
func (s *StackMember) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (s *StackMember) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan StackMemberModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	body := sdk.UpdateStackUserRequest{}
	if r := plan.Role.ValueString(); r != "" {
		body.Role = sdk.Role(r)
	}
	resp, err := s.sdk.UpsertStackUserAccess(ctx, plan.OrganizationId.ValueString(), plan.StackId.ValueString(), plan.UserId.ValueString()).UpdateStackUserRequest(body).Execute()
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

	resp, err := s.sdk.DeleteStackUserAccess(ctx, state.OrganizationId.ValueString(), state.StackId.ValueString(), state.UserId.ValueString()).Execute()
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
	resp, err := s.sdk.UpsertStackUserAccess(ctx, plan.OrganizationId.ValueString(), plan.StackId.ValueString(), plan.UserId.ValueString()).UpdateStackUserRequest(body).Execute()
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

	userAccess, resp, err := s.sdk.ReadStackUserAccess(ctx, state.OrganizationId.ValueString(), state.StackId.ValueString(), state.UserId.ValueString()).Execute()
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
