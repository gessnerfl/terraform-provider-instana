package instana_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceAlertingChannelVictorOpsDefinitionTemplate = `
resource "instana_alerting_channel_victor_ops" "example" {
  name  = "name %d"
	api_key   = "api key"
	routing_key = "routing key"
}
`

const alertingChannelVictorOpsServerResponseTemplate = `
{
	"id"         : "%s",
	"name"       : "name %d",
	"kind"       : "VICTOR_OPS",
	"apiKey"     : "api key",
	"routingKey" : "routing key"
}
`

const testAlertingChannelVictorOpsDefinition = "instana_alerting_channel_victor_ops.example"
const testAlertingChannelVictorOpsRoutingKey = "routing key"
const testAlertingChannelVictorOpsApiKey = "api key"

func TestCRUDOfAlertingChannelVictorOpsResourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.AlertingChannelsResourcePath, alertingChannelVictorOpsServerResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAlertingChannelVictorOpsResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testAlertingChannelVictorOpsDefinition),
			createAlertingChannelVictorOpsResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testAlertingChannelVictorOpsDefinition),
		},
	})
}

func createAlertingChannelVictorOpsResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceAlertingChannelVictorOpsDefinitionTemplate, iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testAlertingChannelVictorOpsDefinition, "id"),
			resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldAPIKey, testAlertingChannelVictorOpsApiKey),
			resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldRoutingKey, testAlertingChannelVictorOpsRoutingKey),
		),
	}
}

func TestResourceAlertingChannelVictorOpsDefinition(t *testing.T) {
	resourceHandle := NewAlertingChannelVictorOpsResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldRoutingKey)
}

func TestShouldUpdateResourceStateForAlertingChanneVictorOps(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelVictorOpsResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	apiKey := testAlertingChannelVictorOpsApiKey
	routingKey := testAlertingChannelVictorOpsRoutingKey
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       resourceName,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, apiKey, resourceData.Get(AlertingChannelVictorOpsFieldAPIKey))
	require.Equal(t, routingKey, resourceData.Get(AlertingChannelVictorOpsFieldRoutingKey))
}

func TestShouldConvertStateOfAlertingChannelVictorOpsToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelVictorOpsResourceHandle()
	apiKey := testAlertingChannelVictorOpsApiKey
	routingKey := testAlertingChannelVictorOpsRoutingKey
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, "name")
	setValueOnResourceData(t, resourceData, AlertingChannelVictorOpsFieldAPIKey, apiKey)
	setValueOnResourceData(t, resourceData, AlertingChannelVictorOpsFieldRoutingKey, routingKey)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceName, model.Name, "name should be equal to full name")
	require.Equal(t, apiKey, *model.APIKey, "api key should be equal")
	require.Equal(t, routingKey, *model.RoutingKey, "routing key should be equal")
}

func TestAlertingChannelVictorOpsShouldHaveSchemaVersionOne(t *testing.T) {
	require.Len(t, NewAlertingChannelVictorOpsResourceHandle().StateUpgraders(), 1)
	require.Equal(t, 0, NewAlertingChannelVictorOpsResourceHandle().StateUpgraders()[0].Version)
}

func TestAlertingChannelVictorOpsShouldMigrateFullNameToNameWhenExecutingFirstStateUpgraderAndFullNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewAlertingChannelVictorOpsResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, AlertingChannelFieldFullName)
	require.Contains(t, result, AlertingChannelFieldName)
	require.Equal(t, "test", result[AlertingChannelFieldName])
}

func TestAlertingChannelVictorOpsShouldDoNothingWhenExecutingFirstStateUpgraderAndFullNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewAlertingChannelVictorOpsResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelVictorOps(t *testing.T) {
	name := NewAlertingChannelVictorOpsResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_victor_ops")
}
