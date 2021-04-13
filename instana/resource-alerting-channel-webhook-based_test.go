package instana_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceAlertingChannelWebhookBasedDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_{{CHANNEL_TYPE}}" "example" {
  name = "name {{ITERATOR}}"
  webhook_url = "webhook url"
}
`

const alertingChannelWebhookBasedServerResponseTemplate = `
{
	"id"     	 : "{{id}}",
	"name"   	 : "prefix name {{ITERATOR}} suffix",
	"kind"   	 : "{{type}}",
	"webhookUrl" : "webhook url"
}
`

const alertingChannelWebhookBasedApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelWebhookBasedDefinition = "instana_alerting_channel_%s.example"
const alertingChannelWebhookBasedWebhookUrl = "webhook url"

var supportedAlertingChannelWebhookTypes = []restapi.AlertingChannelType{restapi.GoogleChatChannelType, restapi.Office365ChannelType}

func TestCRUDOfAlertingChannelWebhookBasedResourceWithMockServer(t *testing.T) {
	for _, channelType := range supportedAlertingChannelWebhookTypes {
		t.Run(fmt.Sprintf("TestResourceAlertingChannelWebhookBasedDefinition%s", channelType), func(t *testing.T) {
			testutils.DeactivateTLSServerCertificateVerification()
			httpServer := testutils.NewTestHTTPServer()
			httpServer.AddRoute(http.MethodPut, alertingChannelWebhookBasedApiPath, testutils.EchoHandlerFunc)
			httpServer.AddRoute(http.MethodDelete, alertingChannelWebhookBasedApiPath, testutils.EchoHandlerFunc)
			httpServer.AddRoute(http.MethodGet, alertingChannelWebhookBasedApiPath, func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				path := restapi.AlertingChannelsResourcePath + "/" + vars["id"]
				json := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(alertingChannelWebhookBasedServerResponseTemplate, "{{id}}", vars["id"]), "{{type}}", string(channelType)), "{{ITERATOR}}", strconv.Itoa(getZeroBasedCallCount(httpServer, http.MethodPut, path)))
				w.Header().Set(contentType, r.Header.Get(contentType))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(json))
			})
			httpServer.Start()
			defer httpServer.Close()

			channelTypeString := strings.ToLower(string(channelType))
			resourceDefinitionWithoutName := strings.ReplaceAll(strings.ReplaceAll(resourceAlertingChannelWebhookBasedDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort())), "{{CHANNEL_TYPE}}", channelTypeString)
			resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
			resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")
			resourceName := fmt.Sprintf(testAlertingChannelWebhookBasedDefinition, channelTypeString)

			resource.ParallelTest(t, resource.TestCase{
				Providers: testProviders,
				Steps: []resource.TestStep{
					{
						Config: resourceDefinitionWithoutName0,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(resourceName, "id"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelFieldName, "name 0"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelFieldFullName, "prefix name 0 suffix"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelWebhookBasedFieldWebhookURL, alertingChannelWebhookBasedWebhookUrl),
						),
					},
					{
						ResourceName:      testApplicationConfigDefinition,
						ImportState:       true,
						ImportStateVerify: true,
					},
					{
						Config: resourceDefinitionWithoutName1,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(resourceName, "id"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelFieldName, "name 1"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelFieldFullName, "prefix name 1 suffix"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelWebhookBasedFieldWebhookURL, alertingChannelWebhookBasedWebhookUrl),
						),
					},
					{
						ResourceName:      testApplicationConfigDefinition,
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			})
		})
	}
}

func TestResourceAlertingChannelGoogleChatDefinition(t *testing.T) {
	testResourceAlertingChannelWebhookBasedDefinition(t, NewAlertingChannelGoogleChatResourceHandle())
}

func TestResourceAlertingChannelOffice365Definition(t *testing.T) {
	testResourceAlertingChannelWebhookBasedDefinition(t, NewAlertingChannelOffice356ResourceHandle())
}

func testResourceAlertingChannelWebhookBasedDefinition(t *testing.T, resourceHandle ResourceHandle) {
	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelWebhookBasedFieldWebhookURL)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelGoogleChat(t *testing.T) {
	name := NewAlertingChannelGoogleChatResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_alerting_channel_google_chat", name)
}

func TestShouldUpdateResourceStateForAlertingChanneWebhookBased(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelGoogleChatResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	webhookURL := alertingChannelWebhookBasedWebhookUrl
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       "prefix name suffix",
		WebhookURL: &webhookURL,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, "prefix name suffix", resourceData.Get(AlertingChannelFieldFullName))
	require.Equal(t, webhookURL, resourceData.Get(AlertingChannelWebhookBasedFieldWebhookURL))
}

func TestShouldConvertStateOfAlertingChannelWebhookBasedToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelGoogleChatResourceHandle()
	webhookURL := alertingChannelWebhookBasedWebhookUrl
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelWebhookBasedFieldWebhookURL, webhookURL)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, "prefix name suffix", model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
	require.Equal(t, webhookURL, *model.(*restapi.AlertingChannel).WebhookURL, "webhook url should be equal")
}

func TestAlertingChannelWebhookBasedShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewAlertingChannelOffice356ResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelWebhookBasedShouldHaveNoStateUpgrader(t *testing.T) {
	require.Equal(t, 0, len(NewAlertingChannelOffice356ResourceHandle().StateUpgraders()))
}

func TestShouldReturnCorrectResourceNameForAlertingChannelOffice365(t *testing.T) {
	name := NewAlertingChannelOffice356ResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_office_365")
}
