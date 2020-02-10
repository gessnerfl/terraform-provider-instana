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

	schemaMap := resource.Schema()

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSplunkFieldURL)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSplunkFieldToken)
}

func TestShouldUpdateResourceStateForAlertingChanneSplunk(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelSplunkResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	url := "url"
	token := "token"
	data := restapi.AlertingChannel{
		ID:    "id",
		Name:  "name",
		URL:   &url,
		Token: &token,
	}

	resourceHandle.UpdateState(resourceData, data)

	assert.Equal(t, "id", resourceData.Id(), "id should be equal")
	assert.Equal(t, "name", resourceData.Get(AlertingChannelFieldFullName), "name should be equal to full name")
	assert.Equal(t, url, resourceData.Get(AlertingChannelSplunkFieldURL), "url should be equal")
	assert.Equal(t, token, resourceData.Get(AlertingChannelSplunkFieldToken), "token should be equal")
}

func TestShouldConvertStateOfAlertingChannelSplunkToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelSplunkResourceHandle()
	url := "url"
	token := "token"
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelSplunkFieldURL, url)
	resourceData.Set(AlertingChannelSplunkFieldToken, token)

	model := resourceHandle.ConvertStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.IsType(t, restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	assert.Equal(t, "id", model.GetID())
	assert.Equal(t, "prefix name suffix", model.(restapi.AlertingChannel).Name, "name should be equal to full name")
	assert.Equal(t, url, *model.(restapi.AlertingChannel).URL, "url should be equal")
	assert.Equal(t, token, *model.(restapi.AlertingChannel).Token, "token should be equal")
}

func TestAlertingChannelSplunkkShouldHaveSchemaVersionZero(t *testing.T) {
	assert.Equal(t, 0, NewAlertingChannelSplunkResourceHandle().SchemaVersion())
}

func TestAlertingChannelSplunkShouldHaveNoStateUpgrader(t *testing.T) {
	assert.Equal(t, 0, len(NewAlertingChannelSplunkResourceHandle().StateUpgraders()))
}

func TestShouldReturnCorrectResourceNameForAlertingChannelSplunk(t *testing.T) {
	name := NewAlertingChannelSplunkResourceHandle().ResourceName()

	assert.Equal(t, name, "instana_alerting_channel_splunk")
}
