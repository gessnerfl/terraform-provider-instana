package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestValidConfigurationOfProvider(t *testing.T) {
	config := Provider()
	if config.Schema == nil {
		t.Error("Expected Schema to be configured")
	}
	validateSchema(config.Schema, t)

	if config.ResourcesMap == nil {
		t.Error("Expected ResourcesMap to be configured")
	}
	validateResourcesMap(config.ResourcesMap, t)

	if config.ConfigureFunc == nil {
		t.Error("Expected ConfigureFunc to be configured")
	}
}

func validateSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	if len(schemaMap) != 2 {
		t.Error("Expected two configuration options for provider")
	}
	validateRequiredSchemaOfTypeString(SchemaFieldAPIToken, schemaMap, t)
	validateRequiredSchemaOfTypeString(SchemaFieldEndpoint, schemaMap, t)
}

func validateRequiredSchemaOfTypeString(schemaField string, schemaMap map[string]*schema.Schema, t *testing.T) {
	s := schemaMap[schemaField]
	if s == nil {
		t.Errorf("Expected configuration for %s", schemaField)
	}
	if s.Type != schema.TypeString {
		t.Errorf("Expected %s to be of type string", schemaField)
	}
	if !s.Required {
		t.Errorf("Expected %s to be required", schemaField)
	}
	if len(s.Description) == 0 {
		t.Errorf("Expected description for schema of %s", schemaField)
	}
}

func validateResourcesMap(resourceMap map[string]*schema.Resource, t *testing.T) {
	if len(resourceMap) != 2 {
		t.Error("Expected two resources to be configured")
	}

	if resourceMap[ResourceInstanaRule] == nil {
		t.Error("Expected a resources to be configured for instana rule")
	}
	if resourceMap[ResourceInstanaRuleBinding] == nil {
		t.Error("Expected a resources to be configured for instana rule binding")
	}
}
