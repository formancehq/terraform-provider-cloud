package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &StackModule{}
	_ resource.ResourceWithConfigure        = &StackModule{}
	_ resource.ResourceWithConfigValidators = &StackModule{}
	_ resource.ResourceWithValidateConfig   = &StackModule{}
)

type StackModule struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

var SchemaStackModule = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required: true,
		},
		"stack_id": schema.StringAttribute{
			Required: true,
		},
		"organization_id": schema.StringAttribute{
			Required: true,
		},
	},
}

type StackModuleModel struct {
	Name           types.String `tfsdk:"name"`
	StackId        types.String `tfsdk:"stack_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
}

func NewStackModule(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &StackModule{
			logger: logger,
		}
	}
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (s *StackModule) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, res *resource.ValidateConfigResponse) {
	var config StackModuleModel
	res.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.Name.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Invalid Name",
			"The name attribute must not be empty.",
		)
	}

	if config.StackId.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("stack_id"),
			"Invalid Stack ID",
			"The stack_id attribute must not be null.",
		)
	}

	if config.OrganizationId.IsNull() {
		res.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Invalid Organization ID",
			"The organization_id attribute must not be null.",
		)
	}
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (s *StackModule) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return nil
}

// Configure implements resource.ResourceWithConfigure.
func (s *StackModule) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (s *StackModule) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan StackModuleModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	resp, err := s.sdk.EnableModule(ctx, plan.OrganizationId.ValueString(), plan.StackId.ValueString()).Name(plan.Name.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *StackModule) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state StackModuleModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	resp, err := s.sdk.DisableModule(ctx, state.OrganizationId.ValueString(), state.StackId.ValueString()).Name(state.Name.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
		return
	}
}

// Metadata implements resource.Resource.
func (s *StackModule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack_module"
}

// Read implements resource.Resource.
func (s *StackModule) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	ctx = logging.ContextWithLogger(ctx, s.logger.WithField("func", "read"))
	var state StackModuleModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	modules, resp, err := s.sdk.ListModules(ctx, state.OrganizationId.ValueString(), state.StackId.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
		return
	}

	obj := collectionutils.First(modules.Data, func(m sdk.Module) bool {
		return m.Name == state.Name.ValueString()
	})
	if obj.Name == "" {
		res.Diagnostics.AddError("Module not found", fmt.Sprintf("Module with name '%s' not found in stack '%s'", state.Name.ValueString(), state.StackId.ValueString()))
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, &state)...)
}

// Schema implements resource.Resource.
func (s *StackModule) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = SchemaStackModule
}

// Update implements resource.Resource.
func (s *StackModule) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	s.logger.WithField("func", "stack_module.update").Debugf("Update operation is not supported for StackModule resource")
	res.Diagnostics.AddError("Update Operation Not Implemented", "The update operation for StackModule is not supported.")
}
