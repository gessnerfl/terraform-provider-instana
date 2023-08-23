package instana_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceAlertingChannelOpsGenieDefinitionTemplate = `
resource "instana_alerting_channel_ops_genie" "example" {
  name = "name %d"
  api_key = "api-key"
  tags = [ "tag1", "tag2" ]
  region = "EU"
}
`

const alertingChannelOpsGenieServerResponseTemplate = `
{
	"id"     : "%s",
	"name"   : "name %d",
	"kind"   : "OPS_GENIE",
	"apiKey" : "api-key",
	"region" : "EU",
	"tags"   : "tag1, tag2"
}
`

const testAlertingChannelOpsGenieDefinition = "instana_alerting_channel_ops_genie.example"

func TestCRUDOfAlertingChannelOpsGenieResourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.AlertingChannelsResourcePath, alertingChannelOpsGenieServerResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAlertingChannelOpsGenielResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testAlertingChannelOpsGenieDefinition),
			createAlertingChannelOpsGenielResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testAlertingChannelOpsGenieDefinition),
		},
	})
}

func createAlertingChannelOpsGenielResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceAlertingChannelOpsGenieDefinitionTemplate, iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testAlertingChannelOpsGenieDefinition, "id"),
			resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".0", "tag1"),
			resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".1", "tag2"),
		),
	}
}

func TestResourceAlertingChannelOpsGenieDefinition(t *testing.T) {
	resourceHandle := NewAlertingChannelOpsGenieResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelOpsGenieFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelOpsGenieFieldRegion)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfStrings(AlertingChannelOpsGenieFieldTags)
}

func TestShouldUpdateResourceStateForAlertingChannelOpsGenieWhenSingleTagIsProvided(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelOpsGenieResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	tags := "tag1"
	data := createAlertingChannelEmailModelForResourceUpdateWithoutTags()
	data.Tags = &tags

	err := resourceHandle.UpdateState(resourceData, data)

	require.Nil(t, err)
	requireBasicAlertingChannelEmailsFieldsSet(t, resourceData)
	require.Equal(t, []interface{}{"tag1"}, resourceData.Get(AlertingChannelOpsGenieFieldTags), "list of tags should be equal")
}

func TestShouldUpdateResourceStateForAlertingChannelOpsGenieWhenMultipleTagsAreProvided(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelOpsGenieResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	tags := "tag1, tag2"
	data := createAlertingChannelEmailModelForResourceUpdateWithoutTags()
	data.Tags = &tags

	err := resourceHandle.UpdateState(resourceData, data)

	require.Nil(t, err)
	requireBasicAlertingChannelEmailsFieldsSet(t, resourceData)
	require.Equal(t, []interface{}{"tag1", "tag2"}, resourceData.Get(AlertingChannelOpsGenieFieldTags), "list of tags should be equal")
}

func requireBasicAlertingChannelEmailsFieldsSet(t *testing.T, resourceData *schema.ResourceData) {
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, "apiKey", resourceData.Get(AlertingChannelOpsGenieFieldAPIKey))
	require.Equal(t, "EU", resourceData.Get(AlertingChannelOpsGenieFieldRegion))
}

func createAlertingChannelEmailModelForResourceUpdateWithoutTags() *restapi.AlertingChannel {
	apiKey := "apiKey"
	region := restapi.EuOpsGenieRegion
	return &restapi.AlertingChannel{
		ID:     "id",
		Name:   resourceName,
		APIKey: &apiKey,
		Region: &region,
	}
}

func TestShouldConvertStateOfAlertingChannelOpsGenieToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelOpsGenieResourceHandle()
	tags := []string{"tag1", "tag2"}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, "name")
	setValueOnResourceData(t, resourceData, AlertingChannelOpsGenieFieldAPIKey, "api key")
	setValueOnResourceData(t, resourceData, AlertingChannelOpsGenieFieldRegion, "EU")
	setValueOnResourceData(t, resourceData, AlertingChannelOpsGenieFieldTags, tags)

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceName, model.Name, "name should be equal to full name")
	require.Equal(t, "api key", *model.APIKey, "api key should be equal")
	require.Equal(t, restapi.EuOpsGenieRegion, *model.Region, "region should be EU")
	require.Equal(t, "tag1,tag2", *model.Tags, "tags should be equal")
}

func TestAlertingChannelOpsGenieShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewAlertingChannelOpsGenieResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelOpsGenieShouldHaveOneStateUpgrader(t *testing.T) {
	require.Equal(t, 1, len(NewAlertingChannelOpsGenieResourceHandle().StateUpgraders()))
	require.Equal(t, 0, NewAlertingChannelOpsGenieResourceHandle().StateUpgraders()[0].Version)
}

func TestAlertingChannelOpsGenieShouldMigrateFullNameToNameWhenExecutingFirstStateUpgraderAndFullNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewAlertingChannelOpsGenieResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, AlertingChannelFieldFullName)
	require.Contains(t, result, AlertingChannelFieldName)
	require.Equal(t, "test", result[AlertingChannelFieldName])
}

func TestAlertingChannelOpsGenieShouldDoNothingWhenExecutingFirstStateUpgraderAndFullNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewAlertingChannelOpsGenieResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelOpsGenie(t *testing.T) {
	name := NewAlertingChannelOpsGenieResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_ops_genie")
}
