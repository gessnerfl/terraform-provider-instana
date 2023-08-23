package instana_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceAlertingChannelSplunkDefinitionTemplate = `
resource "instana_alerting_channel_splunk" "example" {
	name  = "name %d"
	url   = "url"
	token = "token"
}
`

const alertingChannelSplunkServerResponseTemplate = `
{
	"id"    : "%s",
	"name"  : "name %d",
	"kind"  : "SPLUNK",
	"url"   : "url",
	"token" : "token"
}
`

const testAlertingChannelSplunkDefinition = "instana_alerting_channel_splunk.example"

func TestCRUDOfAlertingChannelSplunkResourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.AlertingChannelsResourcePath, alertingChannelSplunkServerResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAlertingChannelSplunkResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testAlertingChannelSplunkDefinition),
			createAlertingChannelSplunkResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testAlertingChannelSplunkDefinition),
		},
	})
}

func createAlertingChannelSplunkResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceAlertingChannelSplunkDefinitionTemplate, iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testAlertingChannelSplunkDefinition, "id"),
			resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelSplunkFieldURL, "url"),
			resource.TestCheckResourceAttr(testAlertingChannelSplunkDefinition, AlertingChannelSplunkFieldToken, "token"),
		),
	}
}

func TestResourceAlertingChannelSplunkDefinition(t *testing.T) {
	resourceHandle := NewAlertingChannelSplunkResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSplunkFieldURL)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSplunkFieldToken)
}

func TestShouldUpdateResourceStateForAlertingChanneSplunk(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelSplunkResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	url := "url"
	token := "token"
	data := restapi.AlertingChannel{
		ID:    "id",
		Name:  resourceName,
		URL:   &url,
		Token: &token,
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, url, resourceData.Get(AlertingChannelSplunkFieldURL))
	require.Equal(t, token, resourceData.Get(AlertingChannelSplunkFieldToken))
}

func TestShouldConvertStateOfAlertingChannelSplunkToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelSplunkResourceHandle()
	url := "url"
	token := "token"
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, "name")
	setValueOnResourceData(t, resourceData, AlertingChannelSplunkFieldURL, url)
	setValueOnResourceData(t, resourceData, AlertingChannelSplunkFieldToken, token)

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceName, model.Name, "name should be equal to full name")
	require.Equal(t, url, *model.URL, "url should be equal")
	require.Equal(t, token, *model.Token, "token should be equal")
}

func TestAlertingChannelSplunkkShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewAlertingChannelSplunkResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelSplunkShouldHaveOneStateUpgrader(t *testing.T) {
	require.Len(t, NewAlertingChannelSplunkResourceHandle().StateUpgraders(), 1)
	require.Equal(t, 0, NewAlertingChannelSplunkResourceHandle().StateUpgraders()[0].Version)
}

func TestAlertingChannelSplunkShouldMigrateFullNameToNameWhenExecutingFirstStateUpgraderAndFullNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewAlertingChannelSplunkResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, AlertingChannelFieldFullName)
	require.Contains(t, result, AlertingChannelFieldName)
	require.Equal(t, "test", result[AlertingChannelFieldName])
}

func TestAlertingChannelSplunkShouldDoNothingWhenExecutingFirstStateUpgraderAndFullNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewAlertingChannelSplunkResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelSplunk(t *testing.T) {
	name := NewAlertingChannelSplunkResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_splunk")
}
