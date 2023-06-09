package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestProviderShouldValidateInternally(t *testing.T) {
	err := Provider().InternalValidate()

	assert.Nil(t, err)
}

func TestProviderShouldContainValidSchemaDefinition(t *testing.T) {
	config := Provider()

	assert.NotNil(t, config.Schema)
	assert.Equal(t, 5, len(config.Schema))

	schemaAssert := testutils.NewTerraformSchemaAssert(config.Schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SchemaFieldAPIToken)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SchemaFieldEndpoint)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(SchemaFieldDefaultNamePrefix, "")
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(SchemaFieldDefaultNameSuffix, "(TF managed)")
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SchemaFieldTlsSkipVerify, false)
}

func TestProviderShouldContainValidResourceDefinitions(t *testing.T) {
	config := Provider()

	assert.Equal(t, 23, len(config.ResourcesMap))

	assert.NotNil(t, config.ResourcesMap[ResourceInstanaAPIToken])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaApplicationConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaApplicationAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaGlobalApplicationAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaSliConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaWebsiteMonitoringConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaWebsiteAlertConfig])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaGroup])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaCustomDashboard])
	assert.NotNil(t, config.ResourcesMap[ResourceInstanaSyntheticMonitor])

	validateResourcesMapForCustomEvents(config.ResourcesMap, t)
	validateResourcesMapForAlerting(config.ResourcesMap, t)
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

func TestProviderShouldContainValidDataSourceDefinitions(t *testing.T) {
	config := Provider()

	assert.Equal(t, 2, len(config.DataSourcesMap))

	assert.NotNil(t, config.DataSourcesMap[DataSourceBuiltinEvent])
	assert.NotNil(t, config.DataSourcesMap[DataSourceSyntheticLocation])

}
