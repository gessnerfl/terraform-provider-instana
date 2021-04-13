package instana_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

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
	"name"   : "prefix name {{ITERATOR}} suffix",
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
		path := restapi.AlertingChannelsResourcePath + "/" + vars["id"]
		json := strings.ReplaceAll(strings.ReplaceAll(alertingChannelPagerDutyServerResponseTemplate, "{{id}}", vars["id"]), "{{ITERATOR}}", strconv.Itoa(getZeroBasedCallCount(httpServer, http.MethodPut, path)))
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelPagerDutyDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelPagerDutyDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelPagerDutyFieldServiceIntegrationKey, "service integration key"),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelPagerDutyDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelPagerDutyFieldServiceIntegrationKey, "service integration key"),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceAlertingChannelPagerDutyDefinition(t *testing.T) {
	resource := NewAlertingChannelPagerDutyResourceHandle()

	schemaMap := resource.MetaData().Schema

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
		Name:                  "prefix name suffix",
		ServiceIntegrationKey: &integrationKey,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, "prefix name suffix", resourceData.Get(AlertingChannelFieldFullName))
	require.Equal(t, integrationKey, resourceData.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey))
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

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, "prefix name suffix", model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
	require.Equal(t, integrationKey, *model.(*restapi.AlertingChannel).ServiceIntegrationKey, "service integration key should be equal")
}

func TestAlertingChannelPagerDutyShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewAlertingChannelPagerDutyResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelPagerDutyShouldHaveNoStateUpgrader(t *testing.T) {
	require.Equal(t, 0, len(NewAlertingChannelPagerDutyResourceHandle().StateUpgraders()))
}

func TestShouldReturnCorrectResourceNameForAlertingChannelPagerDuty(t *testing.T) {
	name := NewAlertingChannelPagerDutyResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_pager_duty")
}
