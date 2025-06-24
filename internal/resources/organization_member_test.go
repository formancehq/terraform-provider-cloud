package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
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

func TestOrganizationMemberConfigure(t *testing.T) {
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
				providerData: internal.NewStore(pkg.NewMockCloudSDK(gomock.NewController(t)), fmt.Sprintf("organization_%s", uuid.NewString())),
			},
		} {

			og := resources.NewOrganizationMember(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

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

func TestOrganizationMemberMetadata(t *testing.T) {
	test(t, func(ctx context.Context) {
		og := resources.NewOrganizationMember(logging.FromContext(ctx))().(resource.ResourceWithConfigure)

		res := resource.MetadataResponse{}

		og.Metadata(ctx, resource.MetadataRequest{
			ProviderTypeName: "test",
		}, &res)

		require.Contains(t, res.TypeName, "_organization_member")
	})
}

func TestOrganizationMemberValidateConfig(t *testing.T) {
	type testCase struct {
		email *string
	}

	for _, tc := range []testCase{
		{},
		{
			email: pointer.For(uuid.NewString()),
		},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			test(t, func(ctx context.Context) {
				og := resources.NewOrganizationMember(logging.FromContext(ctx))().(resource.ResourceWithValidateConfig)

				res := resource.ValidateConfigResponse{
					Diagnostics: []diag.Diagnostic{},
				}

				og.ValidateConfig(ctx, resource.ValidateConfigRequest{
					Config: tfsdk.Config{
						Raw: tftypes.NewValue(tftypes.Object{
							AttributeTypes: getSchemaTypes(resources.SchemaOrganizationMember),
						}, map[string]tftypes.Value{
							"email":   tftypes.NewValue(tftypes.String, tc.email),
							"role":    tftypes.NewValue(tftypes.String, nil),
							"user_id": tftypes.NewValue(tftypes.String, nil),
							"id":      tftypes.NewValue(tftypes.String, nil),
						}),
						Schema: resources.SchemaOrganizationMember,
					},
				}, &res)

				if tc.email == nil {
					require.Len(t, res.Diagnostics, 1, "Expected one diagnostic")
					require.Contains(t, collectionutils.Map(res.Diagnostics, func(d diag.Diagnostic) string {
						return d.Summary()
					}), "Invalid Email")
				} else {
					require.Empty(t, res.Diagnostics, "Expected no diagnostics")
				}
			})
		})
	}
}
