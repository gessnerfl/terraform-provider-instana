package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestProviderShouldValidateInternally(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("Expected that provider configuration is valid but got error: %s", err)
	}
}

func TestValidConfigurationOfProvider(t *testing.T) {
	config := Provider()
	if config.Schema == nil {
		t.Fatal("Expected Schema to be configured")
	}
	validateSchema(config.Schema, t)

	if config.ResourcesMap == nil {
		t.Fatal("Expected ResourcesMap to be configured")
	}
	validateResourcesMap(config.ResourcesMap, t)

	if config.ConfigureFunc == nil {
		t.Fatal("Expected ConfigureFunc to be configured")
	}
}

func validateSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	if len(schemaMap) != 3 {
		t.Fatal("Expected three configuration options for provider")
	}
	validateRequiredSchemaOfTypeString(SchemaFieldAPIToken, schemaMap, t)
	validateRequiredSchemaOfTypeString(SchemaFieldEndpoint, schemaMap, t)
	validateSchemaOfTypeBoolWithDefault(SchemaFieldVerifyServerCertificate, true, schemaMap, t)
}

func validateRequiredSchemaOfTypeString(schemaField string, schemaMap map[string]*schema.Schema, t *testing.T) {
	s := schemaMap[schemaField]
	if s == nil {
		t.Fatalf("Expected configuration for %s", schemaField)
	}
	if s.Type != schema.TypeString {
		t.Fatalf("Expected %s to be of type string", schemaField)
	}
	if !s.Required {
		t.Fatalf("Expected %s to be required", schemaField)
	}
	if len(s.Description) == 0 {
		t.Fatalf("Expected description for schema of %s", schemaField)
	}
}

func validateSchemaOfTypeBoolWithDefault(schemaField string, defaultValue bool, schemaMap map[string]*schema.Schema, t *testing.T) {
	s := schemaMap[schemaField]
	if s == nil {
		t.Fatalf("Expected configuration for %s", schemaField)
	}
	if s.Type != schema.TypeBool {
		t.Fatalf("Expected %s to be of type bool", schemaField)
	}
	if s.Required {
		t.Fatalf("Expected %s to be optional", schemaField)
	}
	if len(s.Description) == 0 {
		t.Fatalf("Expected description for schema of %s", schemaField)
	}
	if s.Default != defaultValue {
		t.Fatalf("Expected default value %t", defaultValue)
	}
}

func validateResourcesMap(resourceMap map[string]*schema.Resource, t *testing.T) {
	if len(resourceMap) != 2 {
		t.Fatal("Expected two resources to be configured")
	}

	if resourceMap[ResourceInstanaRule] == nil {
		t.Fatal("Expected a resources to be configured for instana rule")
	}
	if resourceMap[ResourceInstanaRuleBinding] == nil {
		t.Fatal("Expected a resources to be configured for instana rule binding")
	}
}

func validateConfigureFunc(schemaMap map[string]*schema.Schema, configureFunc func(d *schema.ResourceData) (interface{}, error), t *testing.T) {
	data := make(map[string]interface{})
	data[SchemaFieldAPIToken] = "api-token"
	data[SchemaFieldEndpoint] = "instana.io"
	resourceData := schema.TestResourceDataRaw(t, schemaMap, data)

	result, err := configureFunc(resourceData)

	if err != nil {
		t.Fatalf("expected no error but got %s", err)
	}
	if _, ok := result.(restapi.InstanaAPI); ok == false {
		t.Fatal("expected to get instance of InstanaAPI")
	}
}
