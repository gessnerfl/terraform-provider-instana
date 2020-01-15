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

var testAlertingChannelOpsGenieProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceAlertingChannelOpsGenieDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_ops_genie" "example" {
  name = "name {{ITERATOR}}"
  api_key = "api-key"
  tags = [ "tag1", "tag2" ]
  region = "EU"
}
`

const alertingChannelOpsGenieServerResponseTemplate = `
{
	"id"     : "{{id}}",
	"name"   : "prefix name suffix",
	"kind"   : "OPS_GENIE",
	"apiKey" : "api-key",
	"region" : "EU",
	"tags"   : "tag1, tag2"
}
`

const alertingChannelOpsGenieApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelOpsGenieDefinition = "instana_alerting_channel_ops_genie.example"
const alertingChannelOpsGenieID = "alerting-channel-ops_genie-id"

func TestCRUDOfAlertingChannelOpsGenieResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelOpsGenieApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelOpsGenieApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelOpsGenieApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelOpsGenieServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelOpsGenieDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingChannelOpsGenieProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelOpsGenieDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".0", "tag1"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".1", "tag2"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelOpsGenieDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".0", "tag1"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".1", "tag2"),
				),
			},
		},
	})
}

func TestResourceAlertingChannelOpsGenieDefinition(t *testing.T) {
	resource := NewAlertingChannelOpsGenieResourceHandle()

	schemaMap := resource.GetSchema()

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelOpsGenieFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelOpsGenieFieldRegion)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfStrings(AlertingChannelOpsGenieFieldTags)
}
