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

var testAlertingChannelSplunkProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceAlertingChannelSplunkDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_splunk" "example" {
  name  = "name {{ITERATOR}}"
	url   = "url"
	token = "token"
}
`

const alertingChannelSplunkServerResponseTemplate = `
{
	"id"    : "{{id}}",
	"name"  : "prefix name suffix",
	"kind"  : "SPLUNK",
	"url"   : "url",
	"token" : "token"
}
`

const alertingChannelSplunkApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelSplunkDefinition = "instana_alerting_channel_splunk.example"

func TestCRUDOfAlertingChannelSplunkResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelSplunkApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelSplunkApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelSplunkApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelSplunkServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelSplunkDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingChannelSplunkProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelSplunkDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelSplunkFieldURL, "url"),
					resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelSplunkFieldToken, "token"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelSplunkDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelSplunkFieldURL, "url"),
					resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelSplunkFieldToken, "token"),
				),
			},
		},
	})
}

func TestResourceAlertingChannelSplunkDefinition(t *testing.T) {
	resource := NewAlertingChannelSplunkResourceHandle()

	schemaMap := resource.GetSchema()

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSplunkFieldURL)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSplunkFieldToken)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelSplunk(t *testing.T) {
	name := NewAlertingChannelSplunkResourceHandle().GetResourceName()

	if name != "instana_alerting_channel_splunk" {
		t.Fatal("Expected resource name to be instana_alerting_channel_splunk")
	}
}
