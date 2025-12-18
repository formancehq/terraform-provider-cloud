package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/pkg/membership_client/pkg/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
		"metadata": schema.MapAttribute{
			Description: "A map of metadata key-value pairs to associate with the stack.",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.Map{
				mapplanmodifier.UseStateForUnknown(),
			},
		},
	},
}

type StackModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	RegionID types.String `tfsdk:"region_id"`
	Version  types.String `tfsdk:"version"`
	URI      types.String `tfsdk:"uri"`

	Metadata types.Map `tfsdk:"metadata"`

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
	store *internal.Store
}

func NewStack() func() resource.Resource {
	return func() resource.Resource {
		return &Stack{}
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
	// Retrieve the plan
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
	var plan StackModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}

	createStackRequest := &shared.CreateStackRequest{
		Metadata: map[string]string{},
		RegionID: plan.GetRegionID(),
		Name:     plan.GetName(),
		Version:  pointer.For(plan.Version.ValueString()),
	}
	if !plan.Metadata.IsNull() {
		plan.Metadata.ElementsAs(ctx, &createStackRequest.Metadata, false)
	}
	createStackRequest.Metadata["github.com/formancehq/terraform-provider-cloud/protected"] = "true"

	operation, err := s.store.GetSDK().CreateStack(ctx, organizationId, createStackRequest)
	if err != nil {
		pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
		return
	}

	plan.ID = types.StringValue(operation.CreateStackResponse.Data.ID)
	plan.Name = types.StringValue(operation.CreateStackResponse.Data.Name)
	plan.RegionID = types.StringValue(operation.CreateStackResponse.Data.RegionID)
	plan.URI = types.StringValue(operation.CreateStackResponse.Data.URI)
	plan.Version = types.StringNull()
	if operation.CreateStackResponse.Data.Version != nil {
		plan.Version = types.StringValue(*operation.CreateStackResponse.Data.Version)
	}
	plan.Metadata = types.MapNull(types.StringType)
	if len(operation.CreateStackResponse.Data.Metadata) > 0 {
		md := make(map[string]attr.Value, len(operation.CreateStackResponse.Data.Metadata))
		for k, v := range operation.CreateStackResponse.Data.Metadata {
			md[k] = types.StringValue(v)
		}
		//FixMe(hack): Remove the protected metadata to match the config
		delete(md, "github.com/formancehq/terraform-provider-cloud/protected")
		plan.Metadata = types.MapValueMust(types.StringType, md)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (s *Stack) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan StackModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}
	operation, err := s.store.GetSDK().DeleteStack(ctx, organizationId, plan.GetID(), plan.ForceDestroy.ValueBool())
	if err != nil {
		if operation.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning(
				"Stack not found",
				"The stack was not found. It may have already been deleted outside of Terraform.",
			)
			return
		}

		pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
		return
	}
}

// Metadata implements resource.Resource.
func (s *Stack) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stack"
}

// Read implements resource.Resource.
func (s *Stack) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan StackModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}
	op, err := s.store.GetSDK().ReadStack(ctx, organizationId, plan.GetID())
	if err != nil {
		pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
		return
	}

	res := op.CreateStackResponse
	plan.ID = types.StringValue(res.Data.ID)
	plan.Name = types.StringValue(res.Data.Name)
	plan.Version = types.StringNull()
	if res.Data.Version != nil {
		plan.Version = types.StringValue(*res.Data.Version)
	}
	plan.RegionID = types.StringValue(res.Data.RegionID)
	plan.URI = types.StringValue(res.Data.URI)
	plan.Metadata = types.MapNull(types.StringType)
	if len(res.Data.Metadata) > 0 {
		md := make(map[string]attr.Value, len(res.Data.Metadata))
		for k, v := range res.Data.Metadata {
			md[k] = types.StringValue(v)
		}
		//FixMe(hack): Remove the protected metadata to match the config
		delete(md, "github.com/formancehq/terraform-provider-cloud/protected")
		plan.Metadata = types.MapValueMust(types.StringType, md)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Schema implements resource.Resource.
func (s *Stack) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SchemaStack
}

// Update implements resource.Resource.
func (s *Stack) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan StackModel
	var state StackModel
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.RegionID = state.RegionID
	organizationId, err := s.store.GetOrganizationID(ctx)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to get organization ID",
			fmt.Sprintf("Error retrieving organization ID: %s", err),
		)
		return
	}
	if plan.Name.ValueString() != state.Name.ValueString() {
		updateRequest := &shared.StackData{
			Name:     plan.Name.ValueString(),
			Metadata: map[string]string{},
		}

		if !plan.Metadata.IsNull() {
			plan.Metadata.ElementsAs(ctx, &updateRequest.Metadata, false)
		}
		updateRequest.Metadata["github.com/formancehq/terraform-provider-cloud/protected"] = "true"

		operation, err := s.store.GetSDK().UpdateStack(ctx, organizationId, plan.GetID(), updateRequest)
		if err != nil {
			pkg.HandleSDKError(ctx, err, &res.Diagnostics)
			return
		}
		plan.Name = types.StringValue(operation.CreateStackResponse.Data.Name)
		plan.URI = types.StringValue(operation.CreateStackResponse.Data.URI)
		plan.Metadata = types.MapNull(types.StringType)
		if len(operation.CreateStackResponse.Data.Metadata) > 0 {
			md := make(map[string]attr.Value, len(operation.CreateStackResponse.Data.Metadata))
			for k, v := range operation.CreateStackResponse.Data.Metadata {
				md[k] = types.StringValue(v)
			}
			//FixMe(hack): Remove the protected metadata to match the config
			delete(md, "github.com/formancehq/terraform-provider-cloud/protected")
			plan.Metadata = types.MapValueMust(types.StringType, md)
		}
	}

	if state.Version.ValueString() != plan.Version.ValueString() {
		if !semver.IsValid(plan.Version.ValueString()) ||
			(semver.IsValid(plan.Version.ValueString()) && semver.Compare(state.Version.ValueString(), plan.Version.ValueString()) >= 0) {
			_, err := s.store.GetSDK().UpgradeStack(ctx, organizationId, plan.GetID(), plan.Version.ValueString())
			if err != nil {
				pkg.HandleSDKError(ctx, err, &res.Diagnostics)
				return
			}

			plan.Version = types.StringValue(plan.Version.ValueString())
		}
	}

	res.Diagnostics.Append(res.State.Set(ctx, &plan)...)
}
