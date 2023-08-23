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

const resourceAlertingChannelPagerDutyDefinitionTemplate = `
resource "instana_alerting_channel_pager_duty" "example" {
  name = "name %d"
  service_integration_key = "service integration key"
}
`

const alertingChannelPagerDutyServerResponseTemplate = `
{
	"id"     : "%s",
	"name"   : "name %d",
	"kind"   : "PAGER_DUTY",
	"serviceIntegrationKey" : "service integration key"
}
`

const testAlertingChannelPagerDutyDefinition = "instana_alerting_channel_pager_duty.example"

func TestCRUDOfAlertingChannelPagerDutyResourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.AlertingChannelsResourcePath, alertingChannelPagerDutyServerResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAlertingChannelPagerDutyResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testAlertingChannelPagerDutyDefinition),
			createAlertingChannelPagerDutyResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testAlertingChannelPagerDutyDefinition),
		},
	})
}

func createAlertingChannelPagerDutyResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceAlertingChannelPagerDutyDefinitionTemplate, iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testAlertingChannelPagerDutyDefinition, "id"),
			resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelPagerDutyFieldServiceIntegrationKey, "service integration key"),
		),
	}
}

func TestResourceAlertingChannelPagerDutyDefinition(t *testing.T) {
	resourceHandle := NewAlertingChannelPagerDutyResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelPagerDutyFieldServiceIntegrationKey)
}

func TestShouldUpdateResourceStateForAlertingChannePagerDuty(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelPagerDutyResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	integrationKey := "integration key"
	data := restapi.AlertingChannel{
		ID:                    "id",
		Name:                  resourceName,
		ServiceIntegrationKey: &integrationKey,
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, integrationKey, resourceData.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey))
}

func TestShouldConvertStateOfAlertingChannelPagerDutyToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelPagerDutyResourceHandle()
	integrationKey := "integration key"
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, "name")
	setValueOnResourceData(t, resourceData, AlertingChannelPagerDutyFieldServiceIntegrationKey, integrationKey)

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceName, model.Name, "name should be equal to full name")
	require.Equal(t, integrationKey, *model.ServiceIntegrationKey, "service integration key should be equal")
}

func TestAlertingChannelPagerDutyShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewAlertingChannelPagerDutyResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelPagerDutyShouldHaveOneStateUpgrader(t *testing.T) {
	require.Equal(t, 1, len(NewAlertingChannelPagerDutyResourceHandle().StateUpgraders()))
	require.Equal(t, 0, NewAlertingChannelPagerDutyResourceHandle().StateUpgraders()[0].Version)
}

func TestAlertingChannelPagerDutyShouldMigrateFullNameToNameWhenExecutingFirstStateUpgraderAndFullNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewAlertingChannelPagerDutyResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, AlertingChannelFieldFullName)
	require.Contains(t, result, AlertingChannelFieldName)
	require.Equal(t, "test", result[AlertingChannelFieldName])
}

func TestAlertingChannelPagerDutyShouldDoNothingWhenExecutingFirstStateUpgraderAndFullNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewAlertingChannelPagerDutyResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelPagerDuty(t *testing.T) {
	name := NewAlertingChannelPagerDutyResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel_pager_duty")
}
