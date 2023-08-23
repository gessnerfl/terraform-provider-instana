package instana_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceAlertingChannelWebhookBasedDefinitionTemplate = `
resource "instana_alerting_channel_%s" "example" {
  name = "name %d"
  webhook_url = "webhook url"
}
`

const alertingChannelWebhookBasedServerResponseTemplate = `
{
	"id"     	 : "%s",
	"name"   	 : "name %d",
	"kind"   	 : "%s",
	"webhookUrl" : "webhook url"
}
`

const testAlertingChannelWebhookBasedDefinition = "instana_alerting_channel_%s.example"
const alertingChannelWebhookBasedWebhookUrl = "webhook url"

var supportedAlertingChannelWebhookTypes = []restapi.AlertingChannelType{restapi.GoogleChatChannelType, restapi.Office365ChannelType}

func TestCRUDOfAlertingChannelWebhookBasedResourceWithMockServer(t *testing.T) {
	for _, channelType := range supportedAlertingChannelWebhookTypes {
		t.Run(fmt.Sprintf("TestResourceAlertingChannelWebhookBasedDefinition%s", channelType), func(t *testing.T) {
			httpServer := createMockHttpServerForResource(restapi.AlertingChannelsResourcePath, alertingChannelWebhookBasedServerResponseTemplate, string(channelType))
			httpServer.Start()
			defer httpServer.Close()

			definition := fmt.Sprintf(testAlertingChannelWebhookBasedDefinition, strings.ToLower(string(channelType)))
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testProviderFactory,
				Steps: []resource.TestStep{
					createAlertingChannelWebhookBasedResourceTestStep(definition, httpServer.GetPort(), 0, channelType),
					testStepImport(definition),
					createAlertingChannelWebhookBasedResourceTestStep(definition, httpServer.GetPort(), 1, channelType),
					testStepImport(definition),
				},
			})
		})
	}
}

func createAlertingChannelWebhookBasedResourceTestStep(resourceDefinition string, httpPort int64, iteration int, channelType restapi.AlertingChannelType) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceAlertingChannelWebhookBasedDefinitionTemplate, strings.ToLower(string(channelType)), iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resourceDefinition, "id"),
			resource.TestCheckResourceAttr(resourceDefinition, AlertingChannelFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(resourceDefinition, AlertingChannelWebhookBasedFieldWebhookURL, alertingChannelWebhookBasedWebhookUrl),
		),
	}
}

func TestResourceAlertingChannelGoogleChatDefinition(t *testing.T) {
	testResourceAlertingChannelWebhookBasedDefinition(t, NewAlertingChannelGoogleChatResourceHandle())
}

func TestResourceAlertingChannelOffice365Definition(t *testing.T) {
	testResourceAlertingChannelWebhookBasedDefinition(t, NewAlertingChannelOffice365ResourceHandle())
}

func testResourceAlertingChannelWebhookBasedDefinition[T restapi.InstanaDataObject](t *testing.T, resourceHandle ResourceHandle[T]) {
	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelWebhookBasedFieldWebhookURL)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelGoogleChat(t *testing.T) {
	name := NewAlertingChannelGoogleChatResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_alerting_channel_google_chat", name)
}

func TestShouldUpdateResourceStateForAlertingChanneWebhookBased(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelGoogleChatResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	webhookURL := alertingChannelWebhookBasedWebhookUrl
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       resourceName,
		WebhookURL: &webhookURL,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, webhookURL, resourceData.Get(AlertingChannelWebhookBasedFieldWebhookURL))
}

func TestShouldConvertStateOfAlertingChannelWebhookBasedToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelGoogleChatResourceHandle()
	webhookURL := alertingChannelWebhookBasedWebhookUrl
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, "name")
	setValueOnResourceData(t, resourceData, AlertingChannelWebhookBasedFieldWebhookURL, webhookURL)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceName, model.Name, "name should be equal to full name")
	require.Equal(t, webhookURL, *model.WebhookURL, "webhook url should be equal")
}

func TestAlertingChannelWebhookBasedShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewAlertingChannelOffice365ResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelWebhookBasedShouldHaveOneStateUpgrader(t *testing.T) {
	require.Len(t, NewAlertingChannelOffice365ResourceHandle().StateUpgraders(), 1)
	require.Equal(t, 0, NewAlertingChannelOffice365ResourceHandle().StateUpgraders()[0].Version)
}

func TestAlertingChannelWebhookBasedShouldMigrateFullNameToNameWhenExecutingFirstStateUpgraderAndFullNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewAlertingChannelOffice365ResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, AlertingChannelFieldFullName)
	require.Contains(t, result, AlertingChannelFieldName)
	require.Equal(t, "test", result[AlertingChannelFieldName])
}

func TestAlertingChannelWebhookBasedShouldDoNothingWhenExecutingFirstStateUpgraderAndFullNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewAlertingChannelOffice365ResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelOffice365(t *testing.T) {
	name := NewAlertingChannelOffice365ResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_office_365")
}
