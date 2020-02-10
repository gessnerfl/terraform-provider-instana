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
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
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

	schemaMap := resource.Schema()

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

	err := resourceHandle.UpdateState(resourceData, data)

	assert.Nil(t, err)
	assertBasicAlertingChannelEmailsFieldsSet(t, resourceData)
	assert.Equal(t, []interface{}{"tag1"}, resourceData.Get(AlertingChannelOpsGenieFieldTags), "list of tags should be equal")
}

func TestShouldUpdateResourceStateForAlertingChannelOpsGenieWhenMultipleTagsAreProvided(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelOpsGenieResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	tags := "tag1, tag2"
	data := createAlertingChannelEmailModelForResourceUpdateWithoutTags()
	data.Tags = &tags

	err := resourceHandle.UpdateState(resourceData, data)

	assert.Nil(t, err)
	assertBasicAlertingChannelEmailsFieldsSet(t, resourceData)
	assert.Equal(t, []interface{}{"tag1", "tag2"}, resourceData.Get(AlertingChannelOpsGenieFieldTags), "list of tags should be equal")
}

func assertBasicAlertingChannelEmailsFieldsSet(t *testing.T, resourceData *schema.ResourceData) {
	assert.Equal(t, "id", resourceData.Id(), "id should be equal")
	assert.Equal(t, "name", resourceData.Get(AlertingChannelFieldFullName), "name should be equal to full name")
	assert.Equal(t, "apiKey", resourceData.Get(AlertingChannelOpsGenieFieldAPIKey), "api key should be equal")
	assert.Equal(t, "EU", resourceData.Get(AlertingChannelOpsGenieFieldRegion), "region should be EU")
}

func createAlertingChannelEmailModelForResourceUpdateWithoutTags() restapi.AlertingChannel {
	apiKey := "apiKey"
	region := restapi.EuOpsGenieRegion
	return restapi.AlertingChannel{
		ID:     "id",
		Name:   "name",
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

	model, err := resourceHandle.ConvertStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	assert.Equal(t, "id", model.GetID())
	assert.Equal(t, "prefix name suffix", model.(restapi.AlertingChannel).Name, "name should be equal to full name")
	assert.Equal(t, "api key", *model.(restapi.AlertingChannel).APIKey, "api key should be equal")
	assert.Equal(t, restapi.EuOpsGenieRegion, *model.(restapi.AlertingChannel).Region, "region should be EU")
	assert.Equal(t, "tag1,tag2", *model.(restapi.AlertingChannel).Tags, "tags should be equal")
}

func TestAlertingChannelOpsGenieShouldHaveSchemaVersionZero(t *testing.T) {
	assert.Equal(t, 0, NewAlertingChannelOpsGenieResourceHandle().SchemaVersion())
}

func TestAlertingChannelOpsGenieShouldHaveNoStateUpgrader(t *testing.T) {
	assert.Equal(t, 0, len(NewAlertingChannelOpsGenieResourceHandle().StateUpgraders()))
}

func TestShouldReturnCorrectResourceNameForAlertingChannelOpsGenie(t *testing.T) {
	name := NewAlertingChannelOpsGenieResourceHandle().ResourceName()

	assert.Equal(t, name, "instana_alerting_channel_ops_genie")
}
