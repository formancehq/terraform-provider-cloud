package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/mod/semver"
)

var (
	_ resource.Resource                     = &Stack{}
	_ resource.ResourceWithConfigure        = &Stack{}
	_ resource.ResourceWithConfigValidators = &Stack{}
	_ resource.ResourceWithValidateConfig   = &Stack{}
	_ resource.ResourceWithImportState      = &Stack{}
)

var SchemaStack = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"organization_id": schema.StringAttribute{
			Required: true,
		},
		"region_id": schema.StringAttribute{
			Required: true,
		},
		"version": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"force_destroy": schema.BoolAttribute{
			Optional: true,
		},
	},
}

type StackModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	OrganizationID types.String `tfsdk:"organization_id"`
	RegionID       types.String `tfsdk:"region_id"`
	Version        types.String `tfsdk:"version"`

	ForceDestroy types.Bool `tfsdk:"force_destroy"`
}

func (m *StackModel) GetID() string {
	return m.ID.ValueString()
}
func (m *StackModel) GetName() string {
	return m.Name.ValueString()
}
func (m *StackModel) GetOrganizationID() string {
	return m.OrganizationID.ValueString()
}

func (m *StackModel) GetRegionID() string {
	return m.RegionID.ValueString()
}

type Stack struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

func NewStack(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &Stack{
			logger: logger,
		}
	}
}

// ImportState implements resource.ResourceWithImportState.
func (s *Stack) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	var (
		id             string
		organizationId string
	)

	res.Diagnostics.Append(req.Identity.GetAttribute(ctx, path.Root("id"), &id)...)
	res.Diagnostics.Append(req.Identity.GetAttribute(ctx, path.Root("organization_id"), &organizationId)...)
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *Stack) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var config StackModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.OrganizationID.IsNull() {
		res.Diagnostics.AddError("Invalid Organization ID", "Organization ID cannot be null")
	}

	if config.RegionID.IsNull() {
		res.Diagnostics.AddError("Invalid Region ID", "Region ID cannot be null")
	}
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *Stack) ConfigValidators(context.Context) []resource.ConfigValidator {
	return nil
}

// Configure implements resource.ResourceWithConfigure.
func (s *Stack) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(sdk.DefaultAPI)
	if !ok {
		resp.Diagnostics.AddError(
			ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *FormanceCloudProviderModel, got: %T", req.ProviderData),
		)
		return
	}

	s.sdk = sdk
}

// Create implements resource.Resource.
func (s *Stack) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "stack_create"))
	logging.FromContext(ctx).Debugf("Creating stack")

	var plan StackModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createStackRequest := sdk.CreateStackRequest{
		Metadata: pointer.For(map[string]string{
			"github.com/formancehq/terraform-provider/protected": "true",
		}),
		RegionID: plan.GetRegionID(),
		Name:     plan.GetName(),
		Version:  pointer.For(plan.Version.ValueString()),
	}

	obj, res, err := s.sdk.CreateStack(ctx, plan.GetOrganizationID()).CreateStackRequest(createStackRequest).Execute()
	if err != nil {
		internal.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.Name = types.StringValue(obj.Data.Name)
	plan.OrganizationID = types.StringValue(obj.Data.OrganizationId)
	plan.RegionID = types.StringValue(obj.Data.RegionID)
	plan.Version = types.StringNull()
	if obj.Data.Version != nil {
		plan.Version = types.StringValue(*obj.Data.Version)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *Stack) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "stack.delete"))
	logging.FromContext(ctx).Debugf("Deleting stack")
	var plan StackModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := s.sdk.DeleteStack(ctx, plan.GetOrganizationID(), plan.GetID()).Force(plan.ForceDestroy.ValueBool()).Execute()
	if err != nil {
		internal.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}
}

// Metadata implements resource.Resource.
func (s *Stack) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack"
}

// Read implements resource.Resource.
func (s *Stack) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "stack_read"))
	logging.FromContext(ctx).Debugf("Reading stack")

	var plan StackModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, res, err := s.sdk.GetStack(ctx, plan.GetOrganizationID(), plan.GetID()).Execute()
	if err != nil {
		internal.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.Name = types.StringValue(obj.Data.Name)
	plan.OrganizationID = types.StringValue(obj.Data.OrganizationId)
	plan.Version = types.StringNull()
	if obj.Data.Version != nil {
		plan.Version = types.StringValue(*obj.Data.Version)
	}
	plan.RegionID = types.StringValue(obj.Data.RegionID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Schema implements resource.Resource.
func (s *Stack) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SchemaStack
}

// Update implements resource.Resource.
func (s *Stack) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {

	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "stack.update"))
	logging.FromContext(ctx).Debugf("Updating stack")

	var plan StackModel
	var state StackModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.OrganizationID = state.OrganizationID
	plan.RegionID = state.RegionID
	if plan.Name.ValueString() != state.Name.ValueString() {
		updateRequest := sdk.UpdateStackRequest{
			Name: plan.Name.ValueString(),
			Metadata: pointer.For(map[string]string{
				"github.com/formancehq/terraform-provider/protected": "true",
			}),
		}
		obj, resp, err := s.sdk.UpdateStack(ctx, plan.GetOrganizationID(), plan.GetID()).UpdateStackRequest(updateRequest).Execute()
		if err != nil {
			internal.HandleSDKError(ctx, resp, &res.Diagnostics)
			return
		}
		plan.Name = types.StringValue(obj.Data.Name)
	}

	if state.Version.ValueString() != plan.Version.ValueString() {
		if !semver.IsValid(plan.Version.ValueString()) ||
			(semver.IsValid(plan.Version.ValueString()) && semver.Compare(state.Version.ValueString(), plan.Version.ValueString()) >= 0) {
			resp, err := s.sdk.UpgradeStack(ctx, plan.GetOrganizationID(), plan.GetID()).StackVersion(sdk.StackVersion{
				Version: pointer.For(plan.Version.ValueString()),
			}).Execute()
			if err != nil {
				internal.HandleSDKError(ctx, resp, &res.Diagnostics)
				return
			}

			plan.Version = types.StringValue(plan.Version.ValueString())
		}
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}
