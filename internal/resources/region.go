package resources

import (
	"context"
	"fmt"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &Region{}
	_ resource.ResourceWithConfigure        = &Region{}
	_ resource.ResourceWithConfigValidators = &Region{}
	_ resource.ResourceWithValidateConfig   = &Region{}
	_ resource.ResourceWithImportState      = &Region{}
)

type Region struct {
	logger logging.Logger
	sdk    sdk.DefaultAPI
}

var SchemaRegion = schema.Schema{
	Description: "Manages a private region in Formance Cloud. A private region provides dedicated infrastructure for your Formance stacks.",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "The name of the region. Must be unique within the organization.",
			Optional:    true,
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The unique identifier of the region.",
			Computed:    true,
		},
		"organization_id": schema.StringAttribute{
			Description: "The organization ID where the region will be created.",
			Required:    true,
		},
		"base_url": schema.StringAttribute{
			Description: "The base URL of the region API endpoint.",
			Computed:    true,
		},
		"secret": schema.StringAttribute{
			Description: "The secret key for authenticating with the region. This value is only available during creation.",
			Computed:    true,
			Sensitive:   true,
		},
	},
}

type RegionModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	OrganizationID types.String `tfsdk:"organization_id"`
	BaseURL        types.String `tfsdk:"base_url"`
	Secret         types.String `tfsdk:"secret"`
}

func NewRegion(logger logging.Logger) func() resource.Resource {
	return func() resource.Resource {
		return &Region{
			logger: logger,
		}
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *Region) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	panic("ImportState is not implemented for Region resource. Please use the Read operation to fetch the region details.")
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (r *Region) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config RegionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if config.OrganizationID.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Organization ID",
			"While validating a private region, the organization_id was not found in the configuration.",
		)
	}
}

// ConfigValidators implements resource.ResourceWithConfigValidators.
func (r *Region) ConfigValidators(context.Context) []resource.ConfigValidator {
	return nil
}

// Configure implements resource.ResourceWithConfigure.
func (r *Region) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.sdk = sdk
}

// Create implements resource.ResourceWithConfigure.
func (r *Region) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger.WithField("func", "region_create"))
	logging.FromContext(ctx).Debugf("Creating region")

	var plan RegionModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationID.ValueString()
	body := sdk.CreatePrivateRegionRequest{}
	if v := plan.Name.ValueString(); v != "" {
		body.Name = v
	}
	obj, resp, err := r.sdk.CreatePrivateRegion(ctx, organizationId).CreatePrivateRegionRequest(body).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
		return
	}

	plan.ID = types.StringValue(obj.Data.Id)
	plan.BaseURL = types.StringValue(obj.Data.BaseUrl)
	plan.Name = types.StringValue(obj.Data.Name)
	plan.OrganizationID = types.StringValue(obj.Data.OrganizationID)
	plan.Secret = types.StringValue(*obj.Data.Secret.Clear)

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}

// Read implements resource.ResourceWithConfigure.
func (r *Region) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RegionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	objs, res, err := r.sdk.ListRegions(ctx, data.OrganizationID.ValueString()).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, res, &resp.Diagnostics)
		return
	}

	obj := collectionutils.First(objs.Data, func(o sdk.AnyRegion) bool {
		return o.Name == data.Name.ValueString()
	})
	if obj.Id == "" {
		resp.Diagnostics.AddError(
			"Region not found",
			fmt.Sprintf("No region found with name '%s' in organization '%s'", data.Name.ValueString(), data.OrganizationID.ValueString()),
		)
		return
	}

	data.ID = types.StringValue(obj.Id)
	data.Name = types.StringValue(obj.Name)
	data.OrganizationID = types.StringNull()
	if obj.OrganizationID != nil {
		data.OrganizationID = types.StringValue(*obj.OrganizationID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.ResourceWithConfigure.
func (r *Region) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddWarning(
		"Update not supported",
		"Updating a region is not supported. Please delete and recreate the region.",
	)
}

// Delete implements resource.ResourceWithConfigure.
func (r *Region) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = logging.ContextWithLogger(ctx, r.logger.WithField("func", "delete_region"))
	logging.FromContext(ctx).Debugf("Deleting region")

	var plan RegionModel
	diags := req.State.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	if plan.OrganizationID.IsNull() {
		return
	}

	orgId := plan.OrganizationID.ValueString()
	regionId := plan.ID.ValueString()
	resp, err := r.sdk.DeleteRegion(ctx, orgId, regionId).Execute()
	if err != nil {
		pkg.HandleSDKError(ctx, resp, &res.Diagnostics)
		return
	}
}

// Metadata implements resource.ResourceWithConfigure.
func (r *Region) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_region"
}

// Schema implements resource.ResourceWithConfigure.
func (r *Region) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SchemaRegion
}
