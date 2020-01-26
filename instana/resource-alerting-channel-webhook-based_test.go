package instana_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

var testAlertingChannelWebhookBasedProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

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
	"name"   	 : "prefix name suffix",
	"kind"   	 : "{{type}}",
	"webhookUrl" : "webhook url"
}
`

const alertingChannelWebhookBasedApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelWebhookBasedDefinition = "instana_alerting_channel_%s.example"

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
				json := strings.ReplaceAll(strings.ReplaceAll(alertingChannelWebhookBasedServerResponseTemplate, "{{id}}", vars["id"]), "{{type}}", string(channelType))
				w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(json))
			})
			httpServer.Start()
			defer httpServer.Close()

			channelTypeString := strings.ToLower(string(channelType))
			resourceDefinitionWithoutName := strings.ReplaceAll(strings.ReplaceAll(resourceAlertingChannelWebhookBasedDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort())), "{{CHANNEL_TYPE}}", channelTypeString)
			resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
			resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")
			resourceName := fmt.Sprintf(testAlertingChannelWebhookBasedDefinition, channelTypeString)

			resource.UnitTest(t, resource.TestCase{
				Providers: testAlertingChannelWebhookBasedProviders,
				Steps: []resource.TestStep{
					resource.TestStep{
						Config: resourceDefinitionWithoutName0,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(resourceName, "id"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelFieldName, "name 0"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelFieldFullName, "prefix name 0 suffix"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelWebhookBasedFieldWebhookURL, "webhook url"),
						),
					},
					resource.TestStep{
						Config: resourceDefinitionWithoutName1,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(resourceName, "id"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelFieldName, "name 1"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelFieldFullName, "prefix name 1 suffix"),
							resource.TestCheckResourceAttr(resourceName, AlertingChannelWebhookBasedFieldWebhookURL, "webhook url"),
						),
					},
				},
			})
		})
	}
}

func TestResourceAlertingChannelWebhookBasedDefinition(t *testing.T) {
	for _, channelType := range supportedAlertingChannelWebhookTypes {
		t.Run(fmt.Sprintf("TestResourceAlertingChannelWebhookBasedDefinition%s", channelType), func(t *testing.T) {
			resource := NewAlertingChannelWebhookBasedResourceHandle(channelType)

			schemaMap := resource.GetSchema()

			schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
			schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
			schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
			schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelWebhookBasedFieldWebhookURL)
		})
	}
}

func TestShouldReturnCorrectResourceNameForAlertingChannelWebhookBased(t *testing.T) {
	for _, channelType := range supportedAlertingChannelWebhookTypes {
		t.Run(fmt.Sprintf("TestShouldReturnCorrectResourceNameForAlertingChannelWebhookBased%s", channelType), func(t *testing.T) {
			name := NewAlertingChannelWebhookBasedResourceHandle(channelType).GetResourceName()

			if name != fmt.Sprintf("instana_alerting_channel_%s", strings.ToLower(string(channelType))) {
				t.Fatalf("Expected resource name to be instana_alerting_channel_%s", channelType)
			}
		})
	}
}
