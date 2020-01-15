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
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldWebhookURL, "webhook url"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldIconURL, "icon url"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldChannel, "channel"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelSlackDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldWebhookURL, "webhook url"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldIconURL, "icon url"),
					resource.TestCheckResourceAttr(testAlertingChannelSlackDefinition, AlertingChannelSlackFieldChannel, "channel"),
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
