package resources_test

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func test(t *testing.T, fn func(ctx context.Context)) {
	t.Parallel()

	ctx := logging.TestingContext()

	fn(ctx)
}

func getSchemaTypes(schema schema.Schema) map[string]tftypes.Type {
	attributeTypes := make(map[string]tftypes.Type)
	schemaAt := schema.Attributes
	for name, attr := range schemaAt {
		t := attr.GetType()
		attributeTypes[name] = t.TerraformType(logging.TestingContext())
	}

	return attributeTypes
}
