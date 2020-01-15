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

var testAlertingChannelVictorOpsProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceAlertingChannelVictorOpsDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_victor_ops" "example" {
  name  = "name {{ITERATOR}}"
	api_key   = "api key"
	routing_key = "routing key"
}
`

const alertingChannelVictorOpsServerResponseTemplate = `
{
	"id"         : "{{id}}",
	"name"       : "prefix name suffix",
	"kind"       : "VICTOR_OPS",
	"apiKey"     : "api key",
	"routingKey" : "routing key"
}
`

const alertingChannelVictorOpsApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelVictorOpsDefinition = "instana_alerting_channel_victor_ops.example"

func TestCRUDOfAlertingChannelVictorOpsResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelVictorOpsApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelVictorOpsApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelVictorOpsApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelVictorOpsServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelVictorOpsDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingChannelVictorOpsProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelVictorOpsDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldAPIKey, "api key"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldRoutingKey, "routing key"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelVictorOpsDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldAPIKey, "api key"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldRoutingKey, "routing key"),
				),
			},
		},
	})
}

func TestResourceAlertingChannelVictorOpsDefinition(t *testing.T) {
	resource := NewAlertingChannelVictorOpsResourceHandle()

	schemaMap := resource.GetSchema()

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldRoutingKey)
}
