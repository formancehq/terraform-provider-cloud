package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider/internal/membership"
	"github.com/formancehq/terraform-provider/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &Organization{}
	_ resource.ResourceWithConfigure        = &Organization{}
	_ resource.ResourceWithConfigValidators = &Organization{}
	_ resource.ResourceWithValidateConfig   = &Organization{}
	_ resource.ResourceWithImportState      = &Organization{}
)

var SchemaOrganization = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"domain": schema.StringAttribute{
			Optional: true,
		},
		"default_organization_access": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"default_stack_access": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	},
}

type Organization struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

// ImportState implements resource.ResourceWithImportState.
func (o *Organization) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx = logging.ContextWithLogger(ctx, o.logger.WithField("func", "organization_import"))
	logging.FromContext(ctx).Debugf("Importing organization")

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

type OrganizationModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Domain                    types.String `tfsdk:"domain"`
	DefaultOrganizationAccess types.String `tfsdk:"default_organization_access"`
	DefaultStackAccess        types.String `tfsdk:"default_stack_access"`
}

func NewOrganization(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &Organization{
			logger: logger,
		}
	}
}

// Metadata implements resource.Resource.
func (o *Organization) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

// Schema implements resource.Resource.
func (o *Organization) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SchemaOrganization
}

// Create implements resource.Resource.
func (o *Organization) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx = logging.ContextWithLogger(ctx, o.logger.WithField("func", "organization_create"))
	logging.FromContext(ctx).Debugf("Creating organization")

	var plan OrganizationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.OrganizationData{}
	if v := plan.Name.ValueString(); v != "" {
		body.Name = v
	}
	if v := plan.Domain.ValueString(); v != "" {
		body.Domain = pointer.For(v)
	}
	if v := plan.DefaultOrganizationAccess.ValueString(); v != "" {
		body.DefaultOrganizationAccess = pointer.For(sdk.Role(v))
	}
	if v := plan.DefaultStackAccess.ValueString(); v != "" {
		body.DefaultStackAccess = pointer.For(sdk.Role(v))
	}

	obj, res, err := o.sdk.CreateOrganization(ctx).Body(body).Execute()
	if err != nil {
		membership.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.Domain = types.StringNull()
	if obj.Data.Domain != nil {
		plan.Domain = types.StringValue(*obj.Data.Domain)
	}
	plan.Name = types.StringValue(obj.Data.Name)
	plan.DefaultOrganizationAccess = types.StringNull()
	if obj.Data.DefaultOrganizationAccess != nil {
		plan.DefaultOrganizationAccess = types.StringValue(string(*obj.Data.DefaultOrganizationAccess))
	}
	plan.DefaultStackAccess = types.StringNull()
	if obj.Data.DefaultStackAccess != nil {
		plan.DefaultStackAccess = types.StringValue(string(*obj.Data.DefaultStackAccess))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read implements resource.Resource.
func (o *Organization) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = logging.ContextWithLogger(ctx, o.logger.WithField("func", "organization_read"))
	logging.FromContext(ctx).Debugf("Reading organization")
	var state OrganizationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, res, err := o.sdk.ReadOrganization(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		membership.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

	state.Name = types.StringValue(obj.Data.Name)
	if obj.Data.Domain != nil {
		state.Domain = types.StringValue(*obj.Data.Domain)
	}
	state.DefaultOrganizationAccess = types.StringNull()
	if obj.Data.DefaultOrganizationAccess != nil {
		state.DefaultOrganizationAccess = types.StringValue(string(*obj.Data.DefaultOrganizationAccess))
	}
	state.DefaultStackAccess = types.StringNull()
	if obj.Data.DefaultStackAccess != nil {
		state.DefaultStackAccess = types.StringValue(string(*obj.Data.DefaultStackAccess))
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (o *Organization) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx = logging.ContextWithLogger(ctx, o.logger.WithField("func", "organization_update"))
	logging.FromContext(ctx).Debugf("Updating organization")
	var plan, state OrganizationModel
	req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	data := sdk.OrganizationData{}
	if plan.Name.ValueString() != "" {
		data.Name = plan.Name.ValueString()
	}
	if plan.Domain.ValueString() != "" {
		data.Domain = pointer.For(plan.Domain.ValueString())
	}
	if plan.DefaultOrganizationAccess.ValueString() != "" {
		data.DefaultOrganizationAccess = pointer.For(sdk.Role(plan.DefaultOrganizationAccess.ValueString()))
	}
	if plan.DefaultStackAccess.ValueString() != "" {
		data.DefaultStackAccess = pointer.For(sdk.Role(plan.DefaultStackAccess.ValueString()))
	}

	logging.FromContext(ctx).Debugf("Updating organization with data: %v", data)
	obj, res, err := o.sdk.UpdateOrganization(ctx, state.ID.ValueString()).OrganizationData(data).Execute()
	if err != nil {
		membership.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.Name = types.StringValue(obj.Data.Name)
	plan.Domain = types.StringNull()
	if obj.Data.Domain != nil {
		plan.Domain = types.StringValue(*obj.Data.Domain)
	}
	plan.DefaultOrganizationAccess = types.StringNull()
	if obj.Data.DefaultOrganizationAccess != nil {
		plan.DefaultOrganizationAccess = types.StringValue(string(*obj.Data.DefaultOrganizationAccess))
	}
	plan.DefaultStackAccess = types.StringNull()
	if obj.Data.DefaultStackAccess != nil {
		plan.DefaultStackAccess = types.StringValue(string(*obj.Data.DefaultStackAccess))
	}

	resp.State.Set(ctx, &plan)
}

// Delete implements resource.Resource.
func (o *Organization) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = logging.ContextWithLogger(ctx, o.logger.WithField("func", "organization_delete"))
	logging.FromContext(ctx).Debugf("Deleting organization")
	var state OrganizationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := o.sdk.DeleteOrganization(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		membership.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

}

// Configure implements resource.ResourceWithConfigure.
func (o *Organization) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	o.sdk = sdk
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (o *Organization) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	// Add config validators here if needed
	return nil
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (o *Organization) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config OrganizationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Name.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Invalid Organization Name",
			"Organization name cannot be null or unknown.",
		)
	}
}
