package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                   = &StackModule{}
	_ resource.ResourceWithConfigure      = &StackModule{}
	_ resource.ResourceWithValidateConfig = &StackModule{}
)

type StackModule struct {
	store *internal.Store
}

var SchemaStackModule = schema.Schema{
	Description: "Manages modules within a Formance Cloud stack. Modules are individual services that can be enabled or disabled on a stack.",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "The name of the module to enable. Valid module names include: ledger, payments, webhooks, wallets, search, reconciliation, orchestration, auth, stargate.",
			Required:    true,
		},
		"stack_id": schema.StringAttribute{
			Description: "The ID of the stack where the module will be enabled.",
			Required:    true,
		},
	},
}

type StackModuleModel struct {
	Name    types.String `tfsdk:"name"`
	StackId types.String `tfsdk:"stack_id"`
}

func NewStackModule() func() resource.Resource {
	return func() resource.Resource {
		return &StackModule{}
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
}

// Configure implements resource.ResourceWithConfigure.
func (s *StackModule) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (s *StackModule) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan StackModuleModel
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
	_, err = s.store.GetSDK().EnableModule(ctx, organizationId, plan.StackId.ValueString(), plan.Name.ValueString())
	if err != nil {
		pkg.HandleSDKError(ctx, err, &res.Diagnostics)
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
	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}
	_, err = s.store.GetSDK().DisableModule(ctx, organizationId, state.StackId.ValueString(), state.Name.ValueString())
	if err != nil {
		pkg.HandleSDKError(ctx, err, &res.Diagnostics)
		return
	}
}

// Metadata implements resource.Resource.
func (s *StackModule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack_module"
}

// Read implements resource.Resource.
func (s *StackModule) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state StackModuleModel
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
	modules, err := s.store.GetSDK().ListModules(ctx, organizationId, state.StackId.ValueString())
	if err != nil {
		pkg.HandleSDKError(ctx, err, &res.Diagnostics)
		return
	}

	obj := collectionutils.First(modules.ListModulesResponse.Data, func(m shared.Module) bool {
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
	res.Diagnostics.AddError("Update Operation Not Implemented", "The update operation for StackModule is not supported.")
}
