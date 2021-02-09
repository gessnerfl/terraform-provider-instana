package instana_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const resourceAlertingChannelWebhookDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_webhook" "example" {
  name = "name {{ITERATOR}}"
  webhook_urls = [ "url1", "url2" ]
  http_headers = {
	  key1 = "value1"
	  key2 = "value2"
  }
}
`

const alertingChannelWebhookServerResponseTemplate = `
{
	"id"     : "{{id}}",
	"name"   : "prefix name suffix",
	"kind"   : "WEB_HOOK",
	"webhookUrls" : [ "url1", "url2" ],
	"headers" : [ "key1: value1", "key2: value2" ]
}
`

const alertingChannelWebhookApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelWebhookDefinition = "instana_alerting_channel_webhook.example"

func TestCRUDOfAlertingChannelWebhookResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelWebhookApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelWebhookApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelWebhookApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelWebhookServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelWebhookDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	hashFunctionUrls := schema.HashSchema(AlertingChannelWebhookWebhookURLsSchemaField.Elem.(*schema.Schema))
	url1 := "url1"
	url2 := "url2"
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelWebhookDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, fmt.Sprintf("%s.%d", AlertingChannelWebhookFieldWebhookURLs, hashFunctionUrls(url1)), url1),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, fmt.Sprintf("%s.%d", AlertingChannelWebhookFieldWebhookURLs, hashFunctionUrls(url2)), url2),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldHTTPHeaders+".key1", "value1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldHTTPHeaders+".key2", "value2"),
				),
			},
			{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelWebhookDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, fmt.Sprintf("%s.%d", AlertingChannelWebhookFieldWebhookURLs, hashFunctionUrls(url1)), "url1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, fmt.Sprintf("%s.%d", AlertingChannelWebhookFieldWebhookURLs, hashFunctionUrls(url2)), "url2"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldHTTPHeaders+".key1", "value1"),
					resource.TestCheckResourceAttr(testAlertingChannelWebhookDefinition, AlertingChannelWebhookFieldHTTPHeaders+".key2", "value2"),
				),
			},
		},
	})
}

func TestResourceAlertingChannelWebhookDefinition(t *testing.T) {
	resource := NewAlertingChannelWebhookResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingChannelWebhookFieldWebhookURLs)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelWebhook(t *testing.T) {
	name := NewAlertingChannelWebhookResourceHandle().MetaData().ResourceName

	assert.Equal(t, name, "instana_alerting_channel_webhook")
}

func TestAlertingChannelWebhookShouldHaveSchemaVersionOne(t *testing.T) {
	assert.Equal(t, 1, NewAlertingChannelWebhookResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelWebhookShouldHaveOneStateUpgraderForVersionZero(t *testing.T) {
	resourceHandler := NewAlertingChannelWebhookResourceHandle()

	assert.Equal(t, 1, len(resourceHandler.StateUpgraders()))
	assert.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
}

func TestShouldReturnStateOfAlertingChannelWebhookUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[AlertingChannelFieldName] = "name"
	rawData[AlertingChannelFieldFullName] = "fullname"
	rawData[AlertingChannelWebhookFieldWebhookURLs] = []interface{}{"url1", "url2"}
	rawData[AlertingChannelWebhookFieldHTTPHeaders] = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	meta := "dummy"

	result, err := NewAlertingChannelWebhookResourceHandle().StateUpgraders()[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Equal(t, rawData, result)
}

func TestShouldUpdateResourceStateForAlertingChanneWebhookWhenNoHeaderIsProvided(t *testing.T) {
	testShouldUpdateResourceStateForAlertingChanneWebhook(t, []string{}, make(map[string]interface{}))
}

func TestShouldUpdateResourceStateForAlertingChanneWebhookWhenHeadersAreProvided(t *testing.T) {
	headers := []string{"key1: value1", "key2: value2"}
	expectedHeaderMap := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	testShouldUpdateResourceStateForAlertingChanneWebhook(t, headers, expectedHeaderMap)
}

func TestShouldUpdateResourceStateForAlertingChanneWebhookWhenHeaderValueIsNotDefined(t *testing.T) {
	headers := []string{"key1", "key2:"}
	expectedHeaderMap := map[string]interface{}{
		"key1": "",
		"key2": "",
	}
	testShouldUpdateResourceStateForAlertingChanneWebhook(t, headers, expectedHeaderMap)
}

func testShouldUpdateResourceStateForAlertingChanneWebhook(t *testing.T, headersFromApi []string, headersMapped map[string]interface{}) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelWebhookResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	webhookURLs := []string{"url1", "url2"}
	data := restapi.AlertingChannel{
		ID:          "id",
		Name:        "name",
		WebhookURLs: webhookURLs,
		Headers:     headersFromApi,
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	assert.Nil(t, err)
	assert.Equal(t, "id", resourceData.Id(), "id should be equal")
	assert.Equal(t, "name", resourceData.Get(AlertingChannelFieldFullName), "name should be equal to full name")
	assert.Equal(t, headersMapped, resourceData.Get(AlertingChannelWebhookFieldHTTPHeaders))
	urls := resourceData.Get(AlertingChannelWebhookFieldWebhookURLs).(*schema.Set)
	assert.Equal(t, 2, urls.Len())
	assert.Contains(t, urls.List(), "url1")
	assert.Contains(t, urls.List(), "url2")
}

func TestShouldConvertStateOfAlertingChannelWebhookToDataModelWhenNoHeaderIsAvailable(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelWebhookResourceHandle()
	webhookURLs := []string{"url1", "url2"}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelWebhookFieldWebhookURLs, webhookURLs)

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	assert.Equal(t, "id", model.GetID())
	assert.Equal(t, "prefix name suffix", model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
	assert.Len(t, model.(*restapi.AlertingChannel).WebhookURLs, 2)
	assert.Contains(t, model.(*restapi.AlertingChannel).WebhookURLs, "url1")
	assert.Contains(t, model.(*restapi.AlertingChannel).WebhookURLs, "url2")
	assert.Equal(t, []string{}, model.(*restapi.AlertingChannel).Headers, "There should be no headers")
}
