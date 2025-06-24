package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/internal"
	"github.com/formancehq/terraform-provider-cloud/internal/resources"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStackMemberConfigure(t *testing.T) {
	test(t, func(ctx context.Context) {

		type testCase struct {
			providerData any
			expectedErr  error
		}

		for _, tc := range []testCase{
			{},
			{
				providerData: "something",
				expectedErr:  resources.ErrProviderDataNotSet,
			},
			{
				providerData: internal.NewStore(pkg.NewMockDefaultAPI(gomock.NewController(t)), fmt.Sprintf("organization_%s", uuid.NewString())),
			},
		} {

			og := resources.NewStackMember(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

			res := resource.ConfigureResponse{
				Diagnostics: []diag.Diagnostic{},
			}

			og.Configure(ctx, resource.ConfigureRequest{
				ProviderData: tc.providerData,
			}, &res)

			if tc.expectedErr != nil {
				require.Len(t, res.Diagnostics, 1, "Expected one diagnostic")
				require.Equal(t, res.Diagnostics[0].Summary(), tc.expectedErr.Error())
			} else {
				require.Empty(t, res.Diagnostics, "Expected no diagnostics")
			}

		}
	})
}

func TestStackMemberMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := resources.NewStackMember(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

		res := resource.MetadataResponse{}

		og.Metadata(ctx, resource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_stack_member")
	})
}

func TestStackMemberValidateConfig(t *testing.T) {
	type testCase struct {
		stackID *string
		userID  *string
	}
	for _, tc := range []testCase{} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				og := resources.NewStackMember(logging.FromContext(ctx))().(resource.ResourceWithValidateConfig)

				res := resource.ValidateConfigResponse{
					Diagnostics: []diag.Diagnostic{},
				}

				og.ValidateConfig(ctx, resource.ValidateConfigRequest{
					Config: tfsdk.Config{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: getSchemaTypes(resources.SchemaStackMember),
						}, map[string]tftypes.Value{
							"stack_id": tftypes.NewValue(tftypes.String, tc.stackID),
							"user_id":  tftypes.NewValue(tftypes.String, tc.userID),
							"role":     tftypes.NewValue(tftypes.String, nil),
						}),
						Schema: resources.SchemaStackMember,
					},
				}, &res)

				require.Empty(t, res.Diagnostics, "Expected no diagnostics")
			})
		})
	}
}
