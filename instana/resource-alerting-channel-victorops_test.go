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
	"name"       : "prefix name {{ITERATOR}} suffix",
	"kind"       : "VICTOR_OPS",
	"apiKey"     : "api key",
	"routingKey" : "routing key"
}
`

const alertingChannelVictorOpsApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelVictorOpsDefinition = "instana_alerting_channel_victor_ops.example"
const testAlertingChannelVictorOpsRoutingKey = "routing key"
const testAlertingChannelVictorOpsApiKey = "api key"

func TestCRUDOfAlertingChannelVictorOpsResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelVictorOpsApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelVictorOpsApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelVictorOpsApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := restapi.AlertingChannelsResourcePath + "/" + vars["id"]
		json := strings.ReplaceAll(strings.ReplaceAll(alertingChannelVictorOpsServerResponseTemplate, "{{id}}", vars["id"]), "{{ITERATOR}}", strconv.Itoa(getZeroBasedCallCount(httpServer, http.MethodPut, path)))
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelVictorOpsDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelVictorOpsDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldAPIKey, testAlertingChannelVictorOpsApiKey),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldRoutingKey, testAlertingChannelVictorOpsRoutingKey),
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
					resource.TestCheckResourceAttrSet(testAlertingChannelVictorOpsDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldAPIKey, testAlertingChannelVictorOpsApiKey),
					resource.TestCheckResourceAttr(testAlertingChannelVictorOpsDefinition, AlertingChannelVictorOpsFieldRoutingKey, testAlertingChannelVictorOpsRoutingKey),
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

func TestResourceAlertingChannelVictorOpsDefinition(t *testing.T) {
	resource := NewAlertingChannelVictorOpsResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldRoutingKey)
}

func TestShouldUpdateResourceStateForAlertingChanneVictorOps(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelVictorOpsResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	apiKey := testAlertingChannelVictorOpsApiKey
	routingKey := testAlertingChannelVictorOpsRoutingKey
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       "prefix name suffix",
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, "prefix name suffix", resourceData.Get(AlertingChannelFieldFullName))
	require.Equal(t, apiKey, resourceData.Get(AlertingChannelVictorOpsFieldAPIKey))
	require.Equal(t, routingKey, resourceData.Get(AlertingChannelVictorOpsFieldRoutingKey))
}

func TestShouldConvertStateOfAlertingChannelVictorOpsToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelVictorOpsResourceHandle()
	apiKey := testAlertingChannelVictorOpsApiKey
	routingKey := testAlertingChannelVictorOpsRoutingKey
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelVictorOpsFieldAPIKey, apiKey)
	resourceData.Set(AlertingChannelVictorOpsFieldRoutingKey, routingKey)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, "prefix name suffix", model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
	require.Equal(t, apiKey, *model.(*restapi.AlertingChannel).APIKey, "api key should be equal")
	require.Equal(t, routingKey, *model.(*restapi.AlertingChannel).RoutingKey, "routing key should be equal")
}

func TestAlertingChannelVictorOpskShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewAlertingChannelVictorOpsResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelVictorOpsShouldHaveNoStateUpgrader(t *testing.T) {
	require.Equal(t, 0, len(NewAlertingChannelVictorOpsResourceHandle().StateUpgraders()))
}

func TestShouldReturnCorrectResourceNameForAlertingChannelVictorOps(t *testing.T) {
	name := NewAlertingChannelVictorOpsResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_victor_ops")
}
