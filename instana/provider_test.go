package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
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
	if len(schemaMap) != 2 {
		t.Fatal("Expected three configuration options for provider")
	}
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SchemaFieldAPIToken)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SchemaFieldEndpoint)
}

func validateResourcesMap(resourceMap map[string]*schema.Resource, t *testing.T) {
	if len(resourceMap) != 4 {
		t.Fatal("Expected 4 resources to be configured")
	}

	if resourceMap[ResourceInstanaUserRole] == nil {
		t.Fatal("Expected a resources to be configured for instana user roles")
	}
	if resourceMap[ResourceInstanaApplicationConfig] == nil {
		t.Fatal("Expected a resources to be configured for instana application config")
	}
	if resourceMap[ResourceInstanaCustomEventSpecificationSystemRule] == nil {
		t.Fatal("Expected a resources to be configured for instana custom event specification system rule")
	}
	if resourceMap[ResourceInstanaCustomEventSpecificationThresholdRule] == nil {
		t.Fatal("Expected a resources to be configured for instana custom event specification threshold rule")
	}
}

func validateConfigureFunc(schemaMap map[string]*schema.Schema, configureFunc func(*schema.ResourceData) (interface{}, error), t *testing.T) {
	data := make(map[string]interface{})
	data[SchemaFieldAPIToken] = "api-token"
	data[SchemaFieldEndpoint] = "instana.io"
	resourceData := schema.TestResourceDataRaw(t, schemaMap, data)

	result, err := configureFunc(resourceData)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if _, ok := result.(restapi.InstanaAPI); !ok {
		t.Fatal("expected to get instance of InstanaAPI")
	}
}
