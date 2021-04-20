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
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:%d"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_victor_ops" "example" {
  name  = "name %d"
	api_key   = "api key"
	routing_key = "routing key"
}
`

const alertingChannelVictorOpsServerResponseTemplate = `
{
	"id"         : "%s",
	"name"       : "prefix name %d suffix",
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
		Providers: testProviders,
		Steps: []resource.TestStep{
			createAlertingChannelVictorOpsResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testAlertingChannelVictorOpsDefinition),
			createAlertingChannelVictorOpsResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testAlertingChannelVictorOpsDefinition),
		},
	})
}

func createAlertingChannelVictorOpsResourceTestStep(httpPort int, iteration int) resource.TestStep {
	config := fmt.Sprintf(resourceAlertingChannelVictorOpsDefinitionTemplate, httpPort, iteration)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testAlertingChannelVictorOpsDefinition, "id"),
			resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldFullName, formatResourceFullName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldAPIKey, testAlertingChannelVictorOpsApiKey),
			resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldRoutingKey, testAlertingChannelVictorOpsRoutingKey),
		),
	}
}

func TestResourceAlertingChannelVictorOpsDefinition(t *testing.T) {
	resource := NewAlertingChannelVictorOpsResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldRoutingKey)
}

func TestShouldUpdateResourceStateForAlertingChanneVictorOps(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelVictorOpsResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	apiKey := testAlertingChannelVictorOpsApiKey
	routingKey := testAlertingChannelVictorOpsRoutingKey
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       resourceFullName,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, resourceFullName, resourceData.Get(AlertingChannelFieldFullName))
	require.Equal(t, apiKey, resourceData.Get(AlertingChannelVictorOpsFieldAPIKey))
	require.Equal(t, routingKey, resourceData.Get(AlertingChannelVictorOpsFieldRoutingKey))
}

func TestShouldConvertStateOfAlertingChannelVictorOpsToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelVictorOpsResourceHandle()
	apiKey := testAlertingChannelVictorOpsApiKey
	routingKey := testAlertingChannelVictorOpsRoutingKey
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, resourceFullName)
	resourceData.Set(AlertingChannelVictorOpsFieldAPIKey, apiKey)
	resourceData.Set(AlertingChannelVictorOpsFieldRoutingKey, routingKey)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceFullName, model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
	require.Equal(t, apiKey, *model.(*restapi.AlertingChannel).APIKey, "api key should be equal")
	require.Equal(t, routingKey, *model.(*restapi.AlertingChannel).RoutingKey, "routing key should be equal")
}

func TestAlertingChannelVictorOpskShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewAlertingChannelVictorOpsResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelVictorOpsShouldHaveNoStateUpgrader(t *testing.T) {
	require.Equal(t, 0, len(NewAlertingChannelVictorOpsResourceHandle().StateUpgraders()))
}

func TestShouldReturnCorrectResourceNameForAlertingChannelVictorOps(t *testing.T) {
	name := NewAlertingChannelVictorOpsResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_victor_ops")
}
