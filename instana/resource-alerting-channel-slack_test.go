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

const resourceAlertingChannelSlackDefinitionTemplate = `
resource "instana_alerting_channel_slack" "example" {
  	name        = "name %d"
	webhook_url = "webhook url"
	icon_url    = "icon url"
	channel     = "channel"
}
`

const alertingChannelSlackServerResponseTemplate = `
{
	"id"     	   : "%s",
	"name"   	   : "name %d",
	"kind"   	   : "SLACK",
	"webhookUrl" : "webhook url",
	"iconUrl"    : "icon url",
	"channel"    : "channel"
}
`

const testAlertingChannelSlackDefinition = "instana_alerting_channel_slack.example"
const testAlertingChannelSlackWebhookURL = "webhook url"
const testAlertingChannelSlackIconURL = "icon url"
const testAlertingChannelSlackChannel = "channel"

func TestCRUDOfAlertingChannelSlackResourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.AlertingChannelsResourcePath, alertingChannelSlackServerResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAlertingChannelSlackResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testAlertingChannelSlackDefinition),
			createAlertingChannelSlackResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testAlertingChannelSlackDefinition),
		},
	})
}

func createAlertingChannelSlackResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceAlertingChannelSlackDefinitionTemplate, iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testAlertingChannelSlackDefinition, "id"),
			resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldWebhookURL, testAlertingChannelSlackWebhookURL),
			resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldIconURL, testAlertingChannelSlackIconURL),
			resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldChannel, testAlertingChannelSlackChannel),
		),
	}
}

func TestResourceAlertingChannelSlackDefinition(t *testing.T) {
	resourceHandle := NewAlertingChannelSlackResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSlackFieldWebhookURL)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingChannelSlackFieldIconURL)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingChannelSlackFieldChannel)
}

func TestShouldUpdateResourceStateForAlertingChanneSlack(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelSlackResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	webhookURL := testAlertingChannelSlackWebhookURL
	iconURL := testAlertingChannelSlackIconURL
	channel := testAlertingChannelSlackChannel
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       resourceName,
		WebhookURL: &webhookURL,
		IconURL:    &iconURL,
		Channel:    &channel,
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, webhookURL, resourceData.Get(AlertingChannelSlackFieldWebhookURL))
	require.Equal(t, iconURL, resourceData.Get(AlertingChannelSlackFieldIconURL))
	require.Equal(t, channel, resourceData.Get(AlertingChannelSlackFieldChannel))
}

func TestShouldConvertStateOfAlertingChannelSlackToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelSlackResourceHandle()
	webhookURL := testAlertingChannelSlackWebhookURL
	iconURL := testAlertingChannelSlackIconURL
	channel := testAlertingChannelSlackChannel
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, "name")
	setValueOnResourceData(t, resourceData, AlertingChannelSlackFieldWebhookURL, webhookURL)
	setValueOnResourceData(t, resourceData, AlertingChannelSlackFieldIconURL, iconURL)
	setValueOnResourceData(t, resourceData, AlertingChannelSlackFieldChannel, channel)

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceName, model.Name, "name should be equal to full name")
	require.Equal(t, webhookURL, *model.WebhookURL, "webhook url should be equal")
	require.Equal(t, iconURL, *model.IconURL, "icon url should be equal")
	require.Equal(t, channel, *model.Channel, "channel should be equal")
}

func TestAlertingChannelSlackShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewAlertingChannelSlackResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelSlackShouldHaveOneStateUpgrader(t *testing.T) {
	require.Len(t, NewAlertingChannelSlackResourceHandle().StateUpgraders(), 1)
	require.Equal(t, 0, NewAlertingChannelSlackResourceHandle().StateUpgraders()[0].Version)
}

func TestAlertingChannelSlackShouldMigrateFullNameToNameWhenExecutingFirstStateUpgraderAndFullNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewAlertingChannelSlackResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, AlertingChannelFieldFullName)
	require.Contains(t, result, AlertingChannelFieldName)
	require.Equal(t, "test", result[AlertingChannelFieldName])
}

func TestAlertingChannelSlackShouldDoNothingWhenExecutingFirstStateUpgraderAndFullNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewAlertingChannelSlackResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelSlack(t *testing.T) {
	name := NewAlertingChannelSlackResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_slack")
}
