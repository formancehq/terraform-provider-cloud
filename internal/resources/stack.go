package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/mod/semver"
)

var (
	_ resource.Resource                   = &Stack{}
	_ resource.ResourceWithConfigure      = &Stack{}
	_ resource.ResourceWithValidateConfig = &Stack{}
	_ resource.ResourceWithImportState    = &Stack{}
)

var SchemaStack = schema.Schema{
	Description: "Manages a Formance Cloud stack. A stack is an isolated environment where you can deploy and run Formance services.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the stack.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the stack. Must be unique within the organization.",
			Optional:    true,
			Computed:    true,
		},
		"region_id": schema.StringAttribute{
			Description: "The region ID where the stack will be deployed.",
			Required:    true,
		},
		"version": schema.StringAttribute{
			Description: "The version of Formance to deploy. If not specified, the latest version will be used.",
			Optional:    true,
			Computed:    true,
		},
		"force_destroy": schema.BoolAttribute{
			Description: "When set to true, the stack will be forcefully deleted even if it contains data. Use with caution.",
			Optional:    true,
		},
		"uri": schema.StringAttribute{
			Description: "The URI of the deployed stack.",
			Computed:    true,
		},
	},
}

type StackModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	RegionID types.String `tfsdk:"region_id"`
	Version  types.String `tfsdk:"version"`
	URI      types.String `tfsdk:"uri"`

	ForceDestroy types.Bool `tfsdk:"force_destroy"`
}

func (m *StackModel) GetID() string {
	return m.ID.ValueString()
}
func (m *StackModel) GetName() string {
	return m.Name.ValueString()
}

func (m *StackModel) GetRegionID() string {
	return m.RegionID.ValueString()
}

type Stack struct {
	logger logging.Logger
	store  *internal.Store
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

	if config.RegionID.IsNull() {
		res.Diagnostics.AddError("Invalid Region ID", "Region ID cannot be null")
	}
}

// Configure implements resource.ResourceWithConfigure.
func (s *Stack) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	store, ok := req.ProviderData.(*internal.Store)
	if !ok {
		resp.Diagnostics.AddError(
			ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *internal.Store, got: %T", req.ProviderData),
		)
		return
	}

	s.store = store
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
			"github.com/formancehq/terraform-provider-cloud/protected": "true",
		}),
		RegionID: plan.GetRegionID(),
		Name:     plan.GetName(),
		Version:  pointer.For(plan.Version.ValueString()),
	}

	obj, res, err := s.store.GetSDK().CreateStack(ctx, s.store.GetOrganizationID(), createStackRequest)
	if err != nil {
		pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.Name = types.StringValue(obj.Data.Name)
	plan.RegionID = types.StringValue(obj.Data.RegionID)
	plan.URI = types.StringValue(obj.Data.Uri)
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

	res, err := s.store.GetSDK().DeleteStack(ctx, s.store.GetOrganizationID(), plan.GetID(), plan.ForceDestroy.ValueBool())
	if err != nil {
		pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
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

	obj, res, err := s.store.GetSDK().ReadStack(ctx, s.store.GetOrganizationID(), plan.GetID())
	if err != nil {
		pkg.HandleSDKError(ctx, err, res, &resp.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.Name = types.StringValue(obj.Data.Name)
	plan.Version = types.StringNull()
	if obj.Data.Version != nil {
		plan.Version = types.StringValue(*obj.Data.Version)
	}
	plan.RegionID = types.StringValue(obj.Data.RegionID)
	plan.URI = types.StringValue(obj.Data.Uri)

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
	plan.RegionID = state.RegionID
	if plan.Name.ValueString() != state.Name.ValueString() {
		updateRequest := sdk.UpdateStackRequest{
			Name: plan.Name.ValueString(),
			Metadata: pointer.For(map[string]string{
				"github.com/formancehq/terraform-provider-cloud/protected": "true",
			}),
		}
		obj, resp, err := s.store.GetSDK().UpdateStack(ctx, s.store.GetOrganizationID(), plan.GetID(), updateRequest)
		if err != nil {
			pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
			return
		}
		plan.Name = types.StringValue(obj.Data.Name)
		plan.URI = types.StringValue(obj.Data.Uri)
	}

	if state.Version.ValueString() != plan.Version.ValueString() {
		if !semver.IsValid(plan.Version.ValueString()) ||
			(semver.IsValid(plan.Version.ValueString()) && semver.Compare(state.Version.ValueString(), plan.Version.ValueString()) >= 0) {
			resp, err := s.store.GetSDK().UpgradeStack(ctx, s.store.GetOrganizationID(), plan.GetID(), plan.Version.ValueString())
			if err != nil {
				pkg.HandleSDKError(ctx, err, resp, &res.Diagnostics)
				return
			}

			plan.Version = types.StringValue(plan.Version.ValueString())
		}
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}
