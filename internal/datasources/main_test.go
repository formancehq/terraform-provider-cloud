package datasources_test

import (
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func getSchemaTypes(schema schema.Schema) map[string]tftypes.Type {
	attributeTypes := make(map[string]tftypes.Type)
	schemaAt := schema.Attributes
	for name, attr := range schemaAt {
		t := attr.GetType()
		attributeTypes[name] = t.TerraformType(logging.TestingContext())
	}

	return attributeTypes
}
