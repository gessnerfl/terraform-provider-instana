package instana_test

import (
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

var testAlertingChannelWebhookProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceAlertingChannelWebhookDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_webhook" "example" {
  name = "name {{ITERATOR}}"
  webhook_urls = [ "url1", "url2" ]
  http_headers = {
	  key1 = "value1"
	  key2 = "value2"
  }
}
`

const alertingChannelWebhookServerResponseTemplate = `
{
	"id"     : "{{id}}",
	"name"   : "prefix name suffix",
	"kind"   : "WEB_HOOK",
	"webhookUrls" : [ "url1", "url2" ],
	"headers" : [ "key1: value1", "key2: value2" ]
}
`

const alertingChannelWebhookApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelWebhookDefinition = "instana_alerting_channel_webhook.example"

func TestCRUDOfAlertingChannelWebhookResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelWebhookApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelWebhookApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelWebhookApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelWebhookServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelWebhookDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingChannelWebhookProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelWebhookDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldWebhookURLs+".0", "url1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldWebhookURLs+".1", "url2"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldHTTPHeaders+".key1", "value1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldHTTPHeaders+".key2", "value2"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelWebhookDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldWebhookURLs+".0", "url1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldWebhookURLs+".1", "url2"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldHTTPHeaders+".key1", "value1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldHTTPHeaders+".key2", "value2"),
				),
			},
		},
	})
}

func TestResourceAlertingChannelWebhookDefinition(t *testing.T) {
	resource := NewAlertingChannelWebhookResourceHandle()

	schemaMap := resource.GetSchema()

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfStrings(AlertingChannelWebhookFieldWebhookURLs)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelWebhook(t *testing.T) {
	name := NewAlertingChannelWebhookResourceHandle().GetResourceName()

	if name != "instana_alerting_channel_webhook" {
		t.Fatal("Expected resource name to be instana_alerting_channel_webhook")
	}
}
