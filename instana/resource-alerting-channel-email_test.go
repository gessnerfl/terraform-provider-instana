package instana_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"

	"github.com/stretchr/testify/require"
)

const resourceAlertingChannelEmailDefinitionTemplate = `
resource "instana_alerting_channel_email" "example" {
  name = "name %d"
  emails = [ "EMAIL1", "EMAIL2" ]
}
`

const alertingChannelEmailServerResponseTemplate = `
{
	"id"     : "%s",
	"name"   : "name %d",
	"kind"   : "EMAIL",
	"emails" : [ "EMAIL1", "EMAIL2" ]
}
`

const testAlertingChannelEmailDefinition = "instana_alerting_channel_email.example"

func TestCRUDOfAlertingChannelEmailResourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.AlertingChannelsResourcePath, alertingChannelEmailServerResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAlertingChannelEmailResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testAlertingChannelEmailDefinition),
			createAlertingChannelEmailResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testAlertingChannelEmailDefinition),
		},
	})
}

func createAlertingChannelEmailResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceAlertingChannelEmailDefinitionTemplate, iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
			resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, 0), "EMAIL1"),
			resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, 1), "EMAIL2"),
		),
	}
}

func TestResourceAlertingChannelEmailDefinition(t *testing.T) {
	resourceHandle := NewAlertingChannelEmailResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingChannelEmailFieldEmails)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelEmail(t *testing.T) {
	name := NewAlertingChannelEmailResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_alerting_channel_email", name, "Expected resource name to be instana_alerting_channel_email")
}

func TestAlertingChannelEmailResourceShouldHaveSchemaVersionTwo(t *testing.T) {
	require.Equal(t, 2, NewAlertingChannelEmailResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingChannelEmailShouldHaveTwoStateUpgraderForVersionZeroAndOne(t *testing.T) {
	resourceHandler := NewAlertingChannelEmailResourceHandle()

	require.Equal(t, 2, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
	require.Equal(t, 1, resourceHandler.StateUpgraders()[1].Version)
}

func TestAlertingChannelEmailShouldMigrateFullNameToNameWhenExecutingSecondStateUpgraderAndFullNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewAlertingChannelEmailResourceHandle().StateUpgraders()[1].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, AlertingChannelFieldFullName)
	require.Contains(t, result, AlertingChannelFieldName)
	require.Equal(t, "test", result[AlertingChannelFieldName])
}

func TestAlertingChannelEmailShouldDoNothingWhenExecutingSecondStateUpgraderAndFullNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewAlertingChannelEmailResourceHandle().StateUpgraders()[1].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnStateOfAlertingChannelEmailUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	emails := []interface{}{"email1", "email2"}
	name := resourceName
	fullname := "fullname"
	rawData := make(map[string]interface{})
	rawData[AlertingChannelFieldName] = name
	rawData[AlertingChannelFieldFullName] = fullname
	rawData[AlertingChannelEmailFieldEmails] = emails
	meta := "dummy"
	ctx := context.Background()

	result, err := NewAlertingChannelEmailResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, rawData, result)
}

func TestShouldUpdateResourceStateForAlertingChannelEmail(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelEmailResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.AlertingChannel{
		ID:     "id",
		Name:   resourceName,
		Emails: []string{"email1", "email2"},
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	emails := resourceData.Get(AlertingChannelEmailFieldEmails).(*schema.Set)
	require.Equal(t, 2, emails.Len())
	require.Contains(t, emails.List(), "email1")
	require.Contains(t, emails.List(), "email2")
}

func TestShouldConvertStateOfAlertingChannelEmailToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	resourceHandle := NewAlertingChannelEmailResourceHandle()
	emails := []string{"email1", "email2"}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, "name")
	setValueOnResourceData(t, resourceData, AlertingChannelEmailFieldEmails, emails)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, resourceName, model.Name, "name should be equal to full name")
	require.Len(t, model.Emails, 2)
	require.Contains(t, model.Emails, "email1")
	require.Contains(t, model.Emails, "email2")
}
