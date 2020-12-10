package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestProviderShouldValidateInternally(t *testing.T) {
	err := Provider().InternalValidate()

	assert.Nil(t, err)
}

func TestValidConfigurationOfProvider(t *testing.T) {
	config := Provider()

	assert.NotNil(t, config.Schema)
	validateSchema(config.Schema, t)

	assert.NotNil(t, config.ResourcesMap)
	validateResourcesMap(config.ResourcesMap, t)

	assert.NotNil(t, config.ConfigureFunc)
}

func validateSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	assert.Equal(t, 4, len(schemaMap))

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SchemaFieldAPIToken)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SchemaFieldEndpoint)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(SchemaFieldDefaultNamePrefix, "")
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(SchemaFieldDefaultNameSuffix, "(TF managed)")
}

func validateResourcesMap(resourceMap map[string]*schema.Resource, t *testing.T) {
	assert.Equal(t, 16, len(resourceMap))

	assert.NotNil(t, resourceMap[ResourceInstanaUserRole])
	assert.NotNil(t, resourceMap[ResourceInstanaApplicationConfig])

	validateResourcesMapForCustomEvents(resourceMap, t)
	validateResourcesMapForAlerting(resourceMap, t)
}

func validateResourcesMapForCustomEvents(resourceMap map[string]*schema.Resource, t *testing.T) {
	assert.NotNil(t, resourceMap[ResourceInstanaCustomEventSpecificationSystemRule])
	assert.NotNil(t, resourceMap[ResourceInstanaCustomEventSpecificationThresholdRule])
	assert.NotNil(t, resourceMap[ResourceInstanaCustomEventSpecificationEntityVerificationRule])
}

func validateResourcesMapForAlerting(resourceMap map[string]*schema.Resource, t *testing.T) {
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelEmail])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelGoogleChat])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelSlack])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelOffice365])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelOpsGenie])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelPagerDuty])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelSplunk])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelVictorOps])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingChannelWebhook])
	assert.NotNil(t, resourceMap[ResourceInstanaAlertingConfig])
}

func validateConfigureFunc(schemaMap map[string]*schema.Schema, configureFunc func(*schema.ResourceData) (interface{}, error), t *testing.T) {
	data := make(map[string]interface{})
	data[SchemaFieldAPIToken] = "api-token"
	data[SchemaFieldEndpoint] = "instana.io"
	resourceData := schema.TestResourceDataRaw(t, schemaMap, data)

	result, err := configureFunc(resourceData)

	assert.Nil(t, err)
	_, ok := result.(restapi.InstanaAPI)
	assert.True(t, ok)
}
