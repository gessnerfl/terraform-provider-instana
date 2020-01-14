package instana_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

var testAlertingChannelEmailProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceAlertingChannelEmailDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_email" "example" {
  name = "name {{ITERATOR}}"
  emails = [ "EMAIL1", "EMAIL2" ]
}
`

const alertingChannelEmailServerResponseTemplate = `
{
	"id"     : "{{id}}",
	"name"   : "prefix name suffix",
	"kind"   : "EMAIL",
	"emails" : [ "EMAIL1", "EMAIL2" ]
}
`

const alertingChannelEmailApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelEmailDefinition = "instana_alerting_channel_email.example"
const alertingChannelEmailID = "alerting-channel-email-id"

func TestCRUDOfAlertingChannelEmailResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelEmailApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelEmailApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelEmailApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelEmailServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelEmailDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingChannelEmailProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelEmailFieldEmails+".0", "EMAIL1"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelEmailFieldEmails+".1", "EMAIL2"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelEmailFieldEmails+".0", "EMAIL1"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelEmailFieldEmails+".1", "EMAIL2"),
				),
			},
		},
	})
}

func TestResourceAlertingChannelEmailDefinition(t *testing.T) {
	resource := NewAlertingChannelEmailResourceHandle()

	validateAlertingChannelEmailResourceSchema(resource.GetSchema(), t)
}

func validateAlertingChannelEmailResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfStrings(AlertingChannelEmailFieldEmails)
}
