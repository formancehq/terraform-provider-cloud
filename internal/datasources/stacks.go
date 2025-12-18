package datasources

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/formancehq/formance-sdk-cloud-go/pkg/models/shared"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                     = &Stack{}
	_ datasource.DataSourceWithConfigure        = &Stack{}
	_ datasource.DataSourceWithConfigValidators = &Stack{}
)

type Stack struct {
	store *internal.Store
}

// ConfigValidators implements datasource.DataSourceWithConfigValidators.
func (s *Stack) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

var SchemaStack = schema.Schema{
	Description: "Retrieves information about a Formance Cloud stack. If id is specified, returns a specific stack by ID. Otherwise, returns the first available stack sorted alphabetically by name for predictable behavior.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the stack. If not specified, returns the first available stack sorted alphabetically by name.",
			Optional:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the stack.",
			Optional:    true,
		},
		"region_id": schema.StringAttribute{
			Description: "The region ID where the stack is installed.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The current status of the stack.",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "The current state of the stack.",
			Computed:    true,
		},
	},
}

// Configure implements datasource.DataSourceWithConfigure.
func (s *Stack) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	store, ok := req.ProviderData.(*internal.Store)
	if !ok {
		res.Diagnostics.AddError(
			resources.ErrProviderDataNotSet.Error(),
			fmt.Sprintf("Expected *internal.Store, got: %T", req.ProviderData),
		)
		return
	}

	s.store = store
}

type StackModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	RegionID types.String `tfsdk:"region_id"`
	Status   types.String `tfsdk:"status"`
	State    types.String `tfsdk:"state"`
}

func NewStacks() func() datasource.DataSource {
	return func() datasource.DataSource {
		return &Stack{}
	}
}

func (s *Stack) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stacks"
}

func (s *Stack) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SchemaStack
}

func (s *Stack) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StackModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
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
	var stack shared.Stack
	switch {
	case data.ID.ValueString() != "":
		operation, err := s.store.GetSDK().ReadStack(ctx, organizationId, data.ID.ValueString())
		if err != nil {
			pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
			return
		}

		stack = *operation.CreateStackResponse.Data
	case data.Name.ValueString() != "":
		operation, err := s.store.GetSDK().ListStacks(ctx, organizationId)
		if err != nil {
			pkg.HandleSDKError(ctx, err, &resp.Diagnostics)
			return
		}

		if len(operation.ListStacksResponse.Data) == 0 {
			resp.Diagnostics.AddError(
				"No stacks found",
				fmt.Sprintf("No stacks found in organization '%s'", organizationId),
			)
			return
		}

		sort.Slice(operation.ListStacksResponse.Data, func(i, j int) bool {
			return strings.ToLower(operation.ListStacksResponse.Data[i].Name) < strings.ToLower(operation.ListStacksResponse.Data[j].Name)
		})

		stack = operation.ListStacksResponse.Data[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Stack Identifier",
			"Either 'id' or 'name' must be specified to retrieve a stack.",
		)
		return
	}

	data.ID = types.StringValue(stack.ID)
	data.Name = types.StringValue(stack.Name)
	data.RegionID = types.StringValue(stack.RegionID)
	data.Status = types.StringValue(string(stack.Status))
	data.State = types.StringValue(string(stack.State))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
