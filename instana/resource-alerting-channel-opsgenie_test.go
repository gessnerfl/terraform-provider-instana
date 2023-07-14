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
	"name"   : "prefix name %d suffix",
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
			resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelFieldFullName, formatResourceFullName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".0", "tag1"),
			resource.TestCheckResourceAttr(testAlertingChannelOpsGenieDefinition, AlertingChannelOpsGenieFieldTags+".1", "tag2"),
		),
	}
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
	require.Equal(t, resourceFullName, resourceData.Get(AlertingChannelFieldFullName))
	require.Equal(t, "apiKey", resourceData.Get(AlertingChannelOpsGenieFieldAPIKey))
	require.Equal(t, "EU", resourceData.Get(AlertingChannelOpsGenieFieldRegion))
}

func createAlertingChannelEmailModelForResourceUpdateWithoutTags() *restapi.AlertingChannel {
	apiKey := "apiKey"
	region := restapi.EuOpsGenieRegion
	return &restapi.AlertingChannel{
		ID:     "id",
		Name:   resourceFullName,
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
	resourceData.Set(AlertingChannelFieldFullName, resourceFullName)
	resourceData.Set(AlertingChannelOpsGenieFieldAPIKey, "api key")
	resourceData.Set(AlertingChannelOpsGenieFieldRegion, "EU")
	resourceData.Set(AlertingChannelOpsGenieFieldTags, tags)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceFullName, model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
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
