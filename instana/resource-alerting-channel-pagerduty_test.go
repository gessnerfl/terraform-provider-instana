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
	"name"   : "prefix name %d suffix",
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
			resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelFieldFullName, formatResourceFullName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelPagerDutyDefinition, AlertingChannelPagerDutyFieldServiceIntegrationKey, "service integration key"),
		),
	}
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
		Name:                  resourceFullName,
		ServiceIntegrationKey: &integrationKey,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, "name", resourceData.Get(AlertingChannelFieldName))
	require.Equal(t, resourceFullName, resourceData.Get(AlertingChannelFieldFullName))
	require.Equal(t, integrationKey, resourceData.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey))
}

func TestShouldConvertStateOfAlertingChannelPagerDutyToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelPagerDutyResourceHandle()
	integrationKey := "integration key"
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, resourceFullName)
	resourceData.Set(AlertingChannelPagerDutyFieldServiceIntegrationKey, integrationKey)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceFullName, model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
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
