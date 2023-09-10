package instana_test

import (
	"fmt"
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAlertingChannelDataSource(t *testing.T) {
	unitTest := &dataSourceAlertingChannelUnitTest{}
	t.Run("integration test read of office 365 alerting channel", alertingChannelOffice365DataSourceIntegrationTest().testRead)
	t.Run("schema should be valid", unitTest.schemaShouldBeValid)

}

const dataSourceAlertingChannelDefinitionPath = "data.instana_alerting_channel.example"

func alertingChannelOffice365DataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {

	return newDataSourceAlertingChannelIntegrationTest(
		"my-office-356-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelOffice365, AlertingChannelWebhookBasedFieldWebhookURL), "webhook-url-office-365"),
		},
	)
}

func newDataSourceAlertingChannelIntegrationTest(channelName string, additionalChecks []resource.TestCheckFunc) *dataSourceAlertingChannelIntegrationTest {
	return &dataSourceAlertingChannelIntegrationTest{
		channelName:      channelName,
		additionalChecks: additionalChecks,
	}
}

type dataSourceAlertingChannelIntegrationTest struct {
	channelName      string
	additionalChecks []resource.TestCheckFunc
}

func (r *dataSourceAlertingChannelIntegrationTest) testRead(t *testing.T) {
	serverResponse := `
[{
	"id"     	 : "12345",
	"name"   	 : "other1",
	"kind"   	 : "OFFICE_365",
	"webhookUrl" : "webhook url"
},{
	"id"     	 : "23456",
	"name"   	 : "my-office-356-channel",
	"kind"   	 : "OFFICE_365",
	"webhookUrl" : "webhook-url-office-365"
},{
	"id"     	 : "34567",
	"name"   	 : "other2",
	"kind"   	 : "OFFICE_365",
	"webhookUrl" : "webhook url"
}]
`
	httpServer := createMockHttpServerForDataSource(restapi.AlertingChannelsResourcePath, newStringContentResponseProvider(serverResponse))
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			r.createTestStep(httpServer.GetPort()),
		},
	})
}

func (r *dataSourceAlertingChannelIntegrationTest) createTestStep(httpPort int64) resource.TestStep {
	dataSourceAlertingChannelDefinitionTemplate := `
data "instana_alerting_channel" "example" {
  name = "%s"
}
`
	config := appendProviderConfig(fmt.Sprintf(dataSourceAlertingChannelDefinitionTemplate, r.channelName), httpPort)
	checks := append([]resource.TestCheckFunc{resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, "id", "23456")}, r.additionalChecks...)
	return resource.TestStep{
		Config: config,
		Check:  resource.ComposeTestCheckFunc(checks...),
	}
}

type dataSourceAlertingChannelUnitTest struct{}

func (r *dataSourceAlertingChannelUnitTest) schemaShouldBeValid(t *testing.T) {
	schemaData := NewAlertingChannelDataSource().CreateResource().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaData, t)
	require.Len(t, schemaData, 10)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)

	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelEmail)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelOpsGenie)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelPageDuty)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelSlack)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelSplunk)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelVictorOps)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelWebhook)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelOffice365)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfResource(AlertingChannelFieldChannelGoogleChat)

	r.validateEmailChannelSchema(t, schemaData[AlertingChannelFieldChannelEmail].Elem.(*schema.Resource).Schema)
	r.validateOpsGenieChannelSchema(t, schemaData[AlertingChannelFieldChannelOpsGenie].Elem.(*schema.Resource).Schema)
	r.validatePagerDutyChannelSchema(t, schemaData[AlertingChannelFieldChannelPageDuty].Elem.(*schema.Resource).Schema)
	r.validateSlackChannelSchema(t, schemaData[AlertingChannelFieldChannelSlack].Elem.(*schema.Resource).Schema)
	r.validateSplunkChannelSchema(t, schemaData[AlertingChannelFieldChannelSplunk].Elem.(*schema.Resource).Schema)
	r.validateVictorOpsChannelSchema(t, schemaData[AlertingChannelFieldChannelVictorOps].Elem.(*schema.Resource).Schema)
	r.validateWebhookChannelSchema(t, schemaData[AlertingChannelFieldChannelWebhook].Elem.(*schema.Resource).Schema)
	r.validateWebhookBasedChannelSchema(t, schemaData[AlertingChannelFieldChannelOffice365].Elem.(*schema.Resource).Schema)
	r.validateWebhookBasedChannelSchema(t, schemaData[AlertingChannelFieldChannelGoogleChat].Elem.(*schema.Resource).Schema)
}

func (r *dataSourceAlertingChannelUnitTest) validateEmailChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 1)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeSetOfStrings(AlertingChannelEmailFieldEmails)
}

func (r *dataSourceAlertingChannelUnitTest) validateOpsGenieChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 3)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelOpsGenieFieldAPIKey)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelOpsGenieFieldRegion)
	schemaAssert.AssertSchemaIsComputedAndOfTypeListOfStrings(AlertingChannelOpsGenieFieldTags)
}

func (r *dataSourceAlertingChannelUnitTest) validatePagerDutyChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 1)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelPagerDutyFieldServiceIntegrationKey)
}

func (r *dataSourceAlertingChannelUnitTest) validateSlackChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 3)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelSlackFieldWebhookURL)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelSlackFieldIconURL)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelSlackFieldChannel)
}

func (r *dataSourceAlertingChannelUnitTest) validateSplunkChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 2)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelSplunkFieldURL)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelSplunkFieldToken)
}

func (r *dataSourceAlertingChannelUnitTest) validateVictorOpsChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 2)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelVictorOpsFieldAPIKey)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelVictorOpsFieldRoutingKey)
}

func (r *dataSourceAlertingChannelUnitTest) validateWebhookChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 2)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeSetOfStrings(AlertingChannelWebhookFieldWebhookURLs)
	schemaAssert.AssertSchemaIsComputedAndOfTypeMapOfStrings(AlertingChannelWebhookFieldHTTPHeaders)
}

func (r *dataSourceAlertingChannelUnitTest) validateWebhookBasedChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 1)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelWebhookBasedFieldWebhookURL)
}

func (r *dataSourceAlertingChannelUnitTest) shouldHaveSchemaVersion0(t *testing.T) {
	require.Equal(t, 0, NewAlertingChannelResourceHandle().MetaData().SchemaVersion)
}
