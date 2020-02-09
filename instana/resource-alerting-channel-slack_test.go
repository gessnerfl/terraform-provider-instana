package instana_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

var testAlertingChannelSlackProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceAlertingChannelSlackDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_slack" "example" {
  	name        = "name {{ITERATOR}}"
	webhook_url = "webhook url"
	icon_url    = "icon url"
	channel     = "channel"
}
`

const alertingChannelSlackServerResponseTemplate = `
{
	"id"     	   : "{{id}}",
	"name"   	   : "prefix name suffix",
	"kind"   	   : "SLACK",
	"webhookUrl" : "webhook url",
	"iconUrl"    : "icon url",
	"channel"    : "channel"
}
`

const alertingChannelSlackApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelSlackDefinition = "instana_alerting_channel_slack.example"
const testAlertingChannelSlackWebhookURL = "webhook url"
const testAlertingChannelSlackIconURL = "icon url"
const testAlertingChannelSlackChannel = "channel"

func TestCRUDOfAlertingChannelSlackResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelSlackApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelSlackApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelSlackApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelSlackServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelSlackDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingChannelSlackProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelSlackDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldWebhookURL, testAlertingChannelSlackWebhookURL),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldIconURL, testAlertingChannelSlackIconURL),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldChannel, testAlertingChannelSlackChannel),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelSlackDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldWebhookURL, testAlertingChannelSlackWebhookURL),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldIconURL, testAlertingChannelSlackIconURL),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldChannel, testAlertingChannelSlackChannel),
				),
			},
		},
	})
}

func TestResourceAlertingChannelSlackDefinition(t *testing.T) {
	resource := NewAlertingChannelSlackResourceHandle()

	schemaMap := resource.GetSchema()

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSlackFieldWebhookURL)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingChannelSlackFieldIconURL)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingChannelSlackFieldChannel)
}

func TestShouldUpdateResourceStateForAlertingChanneSlack(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelSlackResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	webhookURL := testAlertingChannelSlackWebhookURL
	iconURL := testAlertingChannelSlackIconURL
	channel := testAlertingChannelSlackChannel
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       "name",
		WebhookURL: &webhookURL,
		IconURL:    &iconURL,
		Channel:    &channel,
	}

	resourceHandle.UpdateState(resourceData, data)

	assert.Equal(t, "id", resourceData.Id(), "id should be equal")
	assert.Equal(t, "name", resourceData.Get(AlertingChannelFieldFullName), "name should be equal to full name")
	assert.Equal(t, webhookURL, resourceData.Get(AlertingChannelSlackFieldWebhookURL), "webhook url should be equal")
	assert.Equal(t, iconURL, resourceData.Get(AlertingChannelSlackFieldIconURL), "icon url should be equal")
	assert.Equal(t, channel, resourceData.Get(AlertingChannelSlackFieldChannel), "channel should be equal")
}

func TestShouldConvertStateOfAlertingChannelSlackToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelSlackResourceHandle()
	webhookURL := testAlertingChannelSlackWebhookURL
	iconURL := testAlertingChannelSlackIconURL
	channel := testAlertingChannelSlackChannel
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelSlackFieldWebhookURL, webhookURL)
	resourceData.Set(AlertingChannelSlackFieldIconURL, iconURL)
	resourceData.Set(AlertingChannelSlackFieldChannel, channel)

	model := resourceHandle.ConvertStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.IsType(t, restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	assert.Equal(t, "id", model.GetID())
	assert.Equal(t, "prefix name suffix", model.(restapi.AlertingChannel).Name, "name should be equal to full name")
	assert.Equal(t, webhookURL, *model.(restapi.AlertingChannel).WebhookURL, "webhook url should be equal")
	assert.Equal(t, iconURL, *model.(restapi.AlertingChannel).IconURL, "icon url should be equal")
	assert.Equal(t, channel, *model.(restapi.AlertingChannel).Channel, "channel should be equal")
}

func TestShouldReturnCorrectResourceNameForAlertingChannelSlack(t *testing.T) {
	name := NewAlertingChannelSlackResourceHandle().GetResourceName()

	assert.Equal(t, name, "instana_alerting_channel_slack")
}
