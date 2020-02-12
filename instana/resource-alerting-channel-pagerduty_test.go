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

var testAlertingChannelPagerDutyProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceAlertingChannelPagerDutyDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_pager_duty" "example" {
  name = "name {{ITERATOR}}"
  service_integration_key = "service integration key"
}
`

const alertingChannelPagerDutyServerResponseTemplate = `
{
	"id"     : "{{id}}",
	"name"   : "prefix name suffix",
	"kind"   : "PAGER_DUTY",
	"serviceIntegrationKey" : "service integration key"
}
`

const alertingChannelPagerDutyApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelPagerDutyDefinition = "instana_alerting_channel_pager_duty.example"

func TestCRUDOfAlertingChannelPagerDutyResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelPagerDutyApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelPagerDutyApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelPagerDutyApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelPagerDutyServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelPagerDutyDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingChannelPagerDutyProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelPagerDutyDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelPagerDutyFieldServiceIntegrationKey, "service integration key"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelPagerDutyDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelPagerDutyFieldServiceIntegrationKey, "service integration key"),
				),
			},
		},
	})
}

func TestResourceAlertingChannelPagerDutyDefinition(t *testing.T) {
	resource := NewAlertingChannelPagerDutyResourceHandle()

	schemaMap := resource.Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelPagerDutyFieldServiceIntegrationKey)
}

func TestShouldUpdateResourceStateForAlertingChannePagerDuty(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelPagerDutyResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	integrationKey := "integration key"
	data := restapi.AlertingChannel{
		ID:                    "id",
		Name:                  "name",
		ServiceIntegrationKey: &integrationKey,
	}

	err := resourceHandle.UpdateState(resourceData, data)

	assert.Nil(t, err)
	assert.Equal(t, "id", resourceData.Id(), "id should be equal")
	assert.Equal(t, "name", resourceData.Get(AlertingChannelFieldFullName), "name should be equal to full name")
	assert.Equal(t, integrationKey, resourceData.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey), "service integration key should be equal")
}

func TestShouldConvertStateOfAlertingChannelPagerDutyToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelPagerDutyResourceHandle()
	integrationKey := "integration key"
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelPagerDutyFieldServiceIntegrationKey, integrationKey)

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	assert.Equal(t, "id", model.GetID())
	assert.Equal(t, "prefix name suffix", model.(restapi.AlertingChannel).Name, "name should be equal to full name")
	assert.Equal(t, integrationKey, *model.(restapi.AlertingChannel).ServiceIntegrationKey, "service integration key should be equal")
}

func TestAlertingChannelPagerDutyShouldHaveSchemaVersionZero(t *testing.T) {
	assert.Equal(t, 0, NewAlertingChannelPagerDutyResourceHandle().SchemaVersion)
}

func TestAlertingChannelPagerDutyShouldHaveNoStateUpgrader(t *testing.T) {
	assert.Equal(t, 0, len(NewAlertingChannelPagerDutyResourceHandle().StateUpgraders))
}

func TestShouldReturnCorrectResourceNameForAlertingChannelPagerDuty(t *testing.T) {
	name := NewAlertingChannelPagerDutyResourceHandle().ResourceName

	assert.Equal(t, name, "instana_alerting_channel_pager_duty")
}
