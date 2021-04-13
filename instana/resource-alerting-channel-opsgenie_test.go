package instana_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

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
	"name"   : "prefix name {{ITERATOR}} suffix",
	"kind"   : "OPS_GENIE",
	"apiKey" : "api-key",
	"region" : "EU",
	"tags"   : "tag1, tag2"
}
`

const alertingChannelOpsGenieApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelOpsGenieDefinition = "instana_alerting_channel_ops_genie.example"

func TestCRUDOfAlertingChannelOpsGenieResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelOpsGenieApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelOpsGenieApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelOpsGenieApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := restapi.AlertingChannelsResourcePath + "/" + vars["id"]
		json := strings.ReplaceAll(strings.ReplaceAll(alertingChannelOpsGenieServerResponseTemplate, "{{id}}", vars["id"]), "{{ITERATOR}}", strconv.Itoa(getZeroBasedCallCount(httpServer, http.MethodPut, path)))
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelOpsGenieDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelOpsGenieDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".0", "tag1"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".1", "tag2"),
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
					resource.TestCheckResourceAttrSet(testAlertingChannelOpsGenieDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".0", "tag1"),
					resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".1", "tag2"),
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

func TestResourceAlertingChannelOpsGenieDefinition(t *testing.T) {
	resource := NewAlertingChannelOpsGenieResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelOpsGenieFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelOpsGenieFieldRegion)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfStrings(AlertingChannelOpsGenieFieldTags)
}

func TestShouldUpdateResourceStateForAlertingChannelOpsGenieWhenSingleTagIsProvided(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelOpsGenieResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	tags := "tag1"
	data := createAlertingChannelEmailModelForResourceUpdateWithoutTags()
	data.Tags = &tags

	err := resourceHandle.UpdateState(resourceData, data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	requireBasicAlertingChannelEmailsFieldsSet(t, resourceData)
	require.Equal(t, []interface{}{"tag1"}, resourceData.Get(AlertingChannelOpsGenieFieldTags), "list of tags should be equal")
}

func TestShouldUpdateResourceStateForAlertingChannelOpsGenieWhenMultipleTagsAreProvided(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelOpsGenieResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	tags := "tag1, tag2"
	data := createAlertingChannelEmailModelForResourceUpdateWithoutTags()
	data.Tags = &tags

	err := resourceHandle.UpdateState(resourceData, data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	requireBasicAlertingChannelEmailsFieldsSet(t, resourceData)
	require.Equal(t, []interface{}{"tag1", "tag2"}, resourceData.Get(AlertingChannelOpsGenieFieldTags), "list of tags should be equal")
}

func requireBasicAlertingChannelEmailsFieldsSet(t *testing.T, resourceData *schema.ResourceData) {
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, "prefix name suffix", resourceData.Get(AlertingChannelFieldFullName))
	require.Equal(t, "apiKey", resourceData.Get(AlertingChannelOpsGenieFieldAPIKey))
	require.Equal(t, "EU", resourceData.Get(AlertingChannelOpsGenieFieldRegion))
}

func createAlertingChannelEmailModelForResourceUpdateWithoutTags() *restapi.AlertingChannel {
	apiKey := "apiKey"
	region := restapi.EuOpsGenieRegion
	return &restapi.AlertingChannel{
		ID:     "id",
		Name:   "prefix name suffix",
		APIKey: &apiKey,
		Region: &region,
	}
}

func TestShouldConvertStateOfAlertingChannelOpsGenieToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelOpsGenieResourceHandle()
	tags := []string{"tag1", "tag2"}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelOpsGenieFieldAPIKey, "api key")
	resourceData.Set(AlertingChannelOpsGenieFieldRegion, "EU")
	resourceData.Set(AlertingChannelOpsGenieFieldTags, tags)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, "prefix name suffix", model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
	require.Equal(t, "api key", *model.(*restapi.AlertingChannel).APIKey, "api key should be equal")
	require.Equal(t, restapi.EuOpsGenieRegion, *model.(*restapi.AlertingChannel).Region, "region should be EU")
	require.Equal(t, "tag1,tag2", *model.(*restapi.AlertingChannel).Tags, "tags should be equal")
}

func TestAlertingChannelOpsGenieShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewAlertingChannelOpsGenieResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelOpsGenieShouldHaveNoStateUpgrader(t *testing.T) {
	require.Equal(t, 0, len(NewAlertingChannelOpsGenieResourceHandle().StateUpgraders()))
}

func TestShouldReturnCorrectResourceNameForAlertingChannelOpsGenie(t *testing.T) {
	name := NewAlertingChannelOpsGenieResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_ops_genie")
}
