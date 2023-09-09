package instana_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

func TestAlertingChannelResource(t *testing.T) {
	unitTest := &alertingChannelUnitTest{}
	t.Run("CRUD integration test of with Email Channel", alertingChannelEmailIntegrationTest().testCrud)
	t.Run("CRUD integration test of with OpsGenie Channel", alertingChannelOpsGenieIntegrationTest().testCrud)
	t.Run("CRUD integration test of with PagerDuty Channel", alertingChannelPagerDutyIntegrationTest().testCrud)
	t.Run("CRUD integration test of with Slack Channel", alertingChannelSlackIntegrationTest().testCrud)
	t.Run("CRUD integration test of with Splunk Channel", alertingChannelSplunkIntegrationTest().testCrud)
	t.Run("CRUD integration test of with VictorOps Channel", alertingChannelVictorOpsIntegrationTest().testCrud)
	t.Run("CRUD integration test of with Webhook Channel", alertingChannelWebhookIntegrationTest().testCrud)
	t.Run("CRUD integration test of with Office 365 Channel", alertingChannelOffice365IntegrationTest().testCrud)
	t.Run("CRUD integration test of with Google Chat Channel", alertingChannelGoogleChatIntegrationTest().testCrud)
	t.Run("schema should be valid", unitTest.schemaShouldBeValid)
	t.Run("should have schema version 0", unitTest.shouldHaveSchemaVersion0)
	t.Run("should have no state upgrader", unitTest.shouldHaveNoStateUpgraders)
	t.Run("should have correct resource name", unitTest.shouldHaveCorrectResourceName)
	t.Run("should map email channel to state", unitTest.shouldMapEmailChannelToState)
	t.Run("should map OpsGenie channel to state", unitTest.shouldMapOpsGenieChannelToState)
	t.Run("should map PagerDuty channel to state", unitTest.shouldMapPagerDutyChannelToState)
	t.Run("should map Slack channel to state", unitTest.shouldMapSlackChannelToState)
	t.Run("should map Splunk channel to state", unitTest.shouldMapSplunkChannelToState)
	t.Run("should map VictorOps channel to state", unitTest.shouldMapVictorOpsChannelToState)
	t.Run("should map Webhook channel to state", unitTest.shouldMapWebhookChannelToState)
	t.Run("should map Office 365 channel to state", unitTest.shouldMapOffice365ChannelToState)
	t.Run("should map Google Chat channel to state", unitTest.shouldMapGoogleChatChannelToState)
	t.Run("should fail to map when channel type is not valid", unitTest.shouldFailToMapChannelWhenTypeIsNotValid)
	t.Run("should map state of Email channel to data model", unitTest.shouldMapStateOfEmailChannelToDataModel)
	t.Run("should map state of OpsGenie channel to data model", unitTest.shouldMapStateOfOpsGenieChannelToDataModel)
	t.Run("should map state of PagerDuty channel to data model", unitTest.shouldMapStateOfPagerDutyChannelToDataModel)
	t.Run("should map state of Slack channel to data model", unitTest.shouldMapStateOfSlackChannelToDataModel)
	t.Run("should map state of Splunk channel to data model", unitTest.shouldMapStateOfSplunkChannelToDataModel)
	t.Run("should map state of VictorOps channel to data model", unitTest.shouldMapStateOfVictorOpsChannelToDataModel)
	t.Run("should map state of Webhook channel to data model", unitTest.shouldMapStateOfWebhookChannelToDataModel)
	t.Run("should map state of Webhook channel with headers to data model", unitTest.shouldMapStateOfWebhookChannelWithHeadersToDataModel)
	t.Run("should map state of Office 365 channel to data model", unitTest.shouldMapStateOfOffice365ChannelToDataModel)
	t.Run("should map state of Google Chat channel to data model", unitTest.shouldMapStateOfGoogleChatChannelToDataModel)
	t.Run("should fail to map state when no channel is provided", unitTest.shouldFailToMapStateWhenNoChannelIsProvided)
}

const (
	alertingChannelTestResourceName    = "instana_alerting_channel.example"
	alertingChannelChannelFieldPattern = "%s.0.%s.0.%s"
)

func alertingChannelEmailIntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    email {
      emails = [ "EMAIL1", "EMAIL2" ]
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "EMAIL",
	"emails": [ "EMAIL1", "EMAIL2" ]
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelEmail, AlertingChannelEmailFieldEmails), 0), "EMAIL1"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelEmail, AlertingChannelEmailFieldEmails), 1), "EMAIL2"),
		},
	)
}

func alertingChannelOpsGenieIntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    ops_genie {
      api_key = "api-key"
	  tags 	  = [ "tag1", "tag2" ]
	  region  = "EU"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "OPS_GENIE",
	"apiKey": "api-key",
	"region": "EU",
	"tags": "tag1, tag2"
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelOpsGenie, AlertingChannelOpsGenieFieldAPIKey), "api-key"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelOpsGenie, AlertingChannelOpsGenieFieldRegion), "EU"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelOpsGenie, AlertingChannelOpsGenieFieldTags), 0), "tag1"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelOpsGenie, AlertingChannelOpsGenieFieldTags), 1), "tag2"),
		},
	)
}

func alertingChannelPagerDutyIntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    pager_duty {
  		service_integration_key = "service-integration-key"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "PAGER_DUTY",
	"serviceIntegrationKey": "service-integration-key"
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelPageDuty, AlertingChannelPagerDutyFieldServiceIntegrationKey), "service-integration-key")},
	)
}

func alertingChannelSlackIntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    slack {
		webhook_url = "webhook-url"
		icon_url    = "icon-url"
		channel     = "channel"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "SLACK",
	"webhookUrl": "webhook-url",
	"iconUrl": "icon-url",
	"channel": "channel"
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelSlack, AlertingChannelSlackFieldWebhookURL), "webhook-url"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelSlack, AlertingChannelSlackFieldIconURL), "icon-url"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelSlack, AlertingChannelSlackFieldChannel), "channel"),
		},
	)
}

func alertingChannelSplunkIntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    splunk {
		url   = "url"
		token = "token"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "SPLUNK",
	"url": "url",
	"token": "token"
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelSplunk, AlertingChannelSplunkFieldURL), "url"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelSplunk, AlertingChannelSplunkFieldToken), "token"),
		},
	)
}

func alertingChannelVictorOpsIntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    victor_ops {
		api_key   = "api-key"
		routing_key = "routing-key"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "VICTOR_OPS",
	"apiKey": "api-key",
	"routingKey": "routing-key"
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelVictorOps, AlertingChannelVictorOpsFieldAPIKey), "api-key"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelVictorOps, AlertingChannelVictorOpsFieldRoutingKey), "routing-key"),
		},
	)
}

func alertingChannelWebhookIntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    webhook {
      webhook_urls = [ "url1", "url2" ]
	  http_headers = {
		  key1 = "value1"
		  key2 = "value2"
	  }
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "WEB_HOOK",
	"webhookUrls": [ "url1", "url2" ],
	"headers": [ "key1: value1", "key2: value2" ]
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldWebhookURLs), 0), "url1"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldWebhookURLs), 1), "url2"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf("%s.%s", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldHTTPHeaders), "key1"), "value1"),
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf("%s.%s", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldHTTPHeaders), "key2"), "value2"),
		},
	)
}

func alertingChannelOffice365IntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    office_365 {
      webhook_url = "webhook-url"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "OFFICE_365",
	"webhookUrl": "webhook-url"
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelOffice365, AlertingChannelWebhookBasedFieldWebhookURL), "webhook-url"),
		},
	)
}

func alertingChannelGoogleChatIntegrationTest() *alertingChannelIntegrationTest {
	resourceTemplate := `
resource "instana_alerting_channel" "example" {
  name = "name %d"
  channel {
    google_chat {
      webhook_url = "webhook-url"
    }
  }
}`

	httpServerResponseTemplate := `
{
	"id": "%s",
	"name": "name %d",
	"kind": "GOOGLE_CHAT",
	"webhookUrl": "webhook-url"
}`

	return newAlertingChannelIntegrationTest(
		resourceTemplate,
		alertingChannelTestResourceName,
		httpServerResponseTemplate,
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(alertingChannelTestResourceName, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannel, AlertingChannelFieldChannelGoogleChat, AlertingChannelWebhookBasedFieldWebhookURL), "webhook-url"),
		},
	)
}

func newAlertingChannelIntegrationTest(resourceTemplate string, resourceName string, serverResponseTemplate string, useCaseSpecificChecks []resource.TestCheckFunc) *alertingChannelIntegrationTest {
	return &alertingChannelIntegrationTest{
		resourceTemplate:       resourceTemplate,
		resourceName:           resourceName,
		serverResponseTemplate: serverResponseTemplate,
		useCaseSpecificChecks:  useCaseSpecificChecks,
	}
}

type alertingChannelIntegrationTest struct {
	resourceTemplate       string
	resourceName           string
	serverResponseTemplate string
	useCaseSpecificChecks  []resource.TestCheckFunc
}

func (r *alertingChannelIntegrationTest) testCrud(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.AlertingChannelsResourcePath, r.serverResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: appendProviderConfig(fmt.Sprintf(r.resourceTemplate, 0), httpServer.GetPort()),
				Check:  r.createTestCheckFunctions(0),
			},
			testStepImport(r.resourceName),
			{
				Config: appendProviderConfig(fmt.Sprintf(r.resourceTemplate, 1), httpServer.GetPort()),
				Check:  r.createTestCheckFunctions(1),
			},
			testStepImport(r.resourceName),
		},
	})
}

func (r *alertingChannelIntegrationTest) createTestCheckFunctions(iteration int) resource.TestCheckFunc {
	defaultCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(r.resourceName, "id"),
		resource.TestCheckResourceAttr(r.resourceName, AlertingChannelFieldName, formatResourceName(iteration)),
	}
	allFunctions := append(defaultCheckFunctions, r.useCaseSpecificChecks...)
	return resource.ComposeTestCheckFunc(allFunctions...)
}

type alertingChannelUnitTest struct{}

func (r *alertingChannelUnitTest) schemaShouldBeValid(t *testing.T) {
	schemaData := NewAlertingChannelResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaData, t)
	require.Len(t, schemaData, 2)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfResource(AlertingChannelFieldChannel)

	r.validateChannelSchema(t, schemaData[AlertingChannelFieldChannel].Elem.(*schema.Resource).Schema)
}

func (r *alertingChannelUnitTest) validateChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	require.Len(t, channelSchema, 9)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelEmail)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelOpsGenie)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelPageDuty)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelSlack)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelSplunk)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelVictorOps)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelWebhook)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelOffice365)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(AlertingChannelFieldChannelGoogleChat)

	r.validateEmailChannelSchema(t, channelSchema[AlertingChannelFieldChannelEmail].Elem.(*schema.Resource).Schema)
	r.validateOpsGenieChannelSchema(t, channelSchema[AlertingChannelFieldChannelOpsGenie].Elem.(*schema.Resource).Schema)
	r.validatePagerDutyChannelSchema(t, channelSchema[AlertingChannelFieldChannelPageDuty].Elem.(*schema.Resource).Schema)
	r.validateSlackChannelSchema(t, channelSchema[AlertingChannelFieldChannelSlack].Elem.(*schema.Resource).Schema)
	r.validateSplunkChannelSchema(t, channelSchema[AlertingChannelFieldChannelSplunk].Elem.(*schema.Resource).Schema)
	r.validateVictorOpsChannelSchema(t, channelSchema[AlertingChannelFieldChannelVictorOps].Elem.(*schema.Resource).Schema)
	r.validateWebhookChannelSchema(t, channelSchema[AlertingChannelFieldChannelWebhook].Elem.(*schema.Resource).Schema)
	r.validateWebhookBasedChannelSchema(t, channelSchema[AlertingChannelFieldChannelOffice365].Elem.(*schema.Resource).Schema)
	r.validateWebhookBasedChannelSchema(t, channelSchema[AlertingChannelFieldChannelGoogleChat].Elem.(*schema.Resource).Schema)
}

func (r *alertingChannelUnitTest) validateEmailChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 1)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingChannelEmailFieldEmails)
}

func (r *alertingChannelUnitTest) validateOpsGenieChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 3)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelOpsGenieFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelOpsGenieFieldRegion)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfStrings(AlertingChannelOpsGenieFieldTags)
}

func (r *alertingChannelUnitTest) validatePagerDutyChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 1)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelPagerDutyFieldServiceIntegrationKey)
}

func (r *alertingChannelUnitTest) validateSlackChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 3)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSlackFieldWebhookURL)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingChannelSlackFieldIconURL)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingChannelSlackFieldChannel)
}

func (r *alertingChannelUnitTest) validateSplunkChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 2)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSplunkFieldURL)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelSplunkFieldToken)
}

func (r *alertingChannelUnitTest) validateVictorOpsChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 2)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldAPIKey)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelVictorOpsFieldRoutingKey)
}

func (r *alertingChannelUnitTest) validateWebhookChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 2)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingChannelWebhookFieldWebhookURLs)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeMapOfStrings(AlertingChannelWebhookFieldHTTPHeaders)
}

func (r *alertingChannelUnitTest) validateWebhookBasedChannelSchema(t *testing.T, channelSchema map[string]*schema.Schema) {
	require.Len(t, channelSchema, 1)
	schemaAssert := testutils.NewTerraformSchemaAssert(channelSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelWebhookBasedFieldWebhookURL)
}

func (r *alertingChannelUnitTest) shouldHaveSchemaVersion0(t *testing.T) {
	require.Equal(t, 0, NewAlertingChannelResourceHandle().MetaData().SchemaVersion)
}

func (r *alertingChannelUnitTest) shouldHaveNoStateUpgraders(t *testing.T) {
	resourceHandler := NewAlertingChannelResourceHandle()

	require.Equal(t, 0, len(resourceHandler.StateUpgraders()))
}

func (r *alertingChannelUnitTest) shouldHaveCorrectResourceName(t *testing.T) {
	name := NewAlertingChannelResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_alerting_channel")
}

func (r *alertingChannelUnitTest) shouldMapEmailChannelToState(t *testing.T) {
	data := restapi.AlertingChannel{
		ID:     "id",
		Name:   resourceName,
		Kind:   restapi.EmailChannelType,
		Emails: []string{"email1", "email2"},
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelEmail)

	channel := channels[AlertingChannelFieldChannelEmail].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 1)
	emails := channel[AlertingChannelEmailFieldEmails].(*schema.Set)
	require.Equal(t, 2, emails.Len())
	require.Contains(t, emails.List(), "email1")
	require.Contains(t, emails.List(), "email2")
}

func (r *alertingChannelUnitTest) shouldMapOpsGenieChannelToState(t *testing.T) {
	apiKey := "apiKey"
	region := restapi.EuOpsGenieRegion
	tags := "tag1, tag2"
	data := restapi.AlertingChannel{
		ID:     "id",
		Name:   resourceName,
		Kind:   restapi.OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
		Tags:   &tags,
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelOpsGenie)

	channel := channels[AlertingChannelFieldChannelOpsGenie].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 3)
	require.Equal(t, apiKey, channel[AlertingChannelOpsGenieFieldAPIKey])
	require.Equal(t, string(region), channel[AlertingChannelOpsGenieFieldRegion])
	require.Equal(t, []interface{}{"tag1", "tag2"}, channel[AlertingChannelOpsGenieFieldTags])
}

func (r *alertingChannelUnitTest) shouldMapPagerDutyChannelToState(t *testing.T) {
	integrationKey := "integration key"
	data := restapi.AlertingChannel{
		ID:                    "id",
		Name:                  resourceName,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &integrationKey,
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelPageDuty)

	channel := channels[AlertingChannelFieldChannelPageDuty].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 1)
	require.Equal(t, integrationKey, channel[AlertingChannelPagerDutyFieldServiceIntegrationKey])
}

func (r *alertingChannelUnitTest) shouldMapSlackChannelToState(t *testing.T) {
	webhookURL := testAlertingChannelSlackWebhookURL
	iconURL := testAlertingChannelSlackIconURL
	slackChannel := testAlertingChannelSlackChannel
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       resourceName,
		Kind:       restapi.SlackChannelType,
		WebhookURL: &webhookURL,
		IconURL:    &iconURL,
		Channel:    &slackChannel,
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelSlack)

	channel := channels[AlertingChannelFieldChannelSlack].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 3)
	require.Equal(t, webhookURL, channel[AlertingChannelSlackFieldWebhookURL])
	require.Equal(t, iconURL, channel[AlertingChannelSlackFieldIconURL])
	require.Equal(t, slackChannel, channel[AlertingChannelSlackFieldChannel])
}

func (r *alertingChannelUnitTest) shouldMapSplunkChannelToState(t *testing.T) {
	url := "url"
	token := "token"
	data := restapi.AlertingChannel{
		ID:    "id",
		Name:  resourceName,
		Kind:  restapi.SplunkChannelType,
		URL:   &url,
		Token: &token,
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelSplunk)

	channel := channels[AlertingChannelFieldChannelSplunk].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 2)
	require.Equal(t, url, channel[AlertingChannelSplunkFieldURL])
	require.Equal(t, token, channel[AlertingChannelSplunkFieldToken])
}

func (r *alertingChannelUnitTest) shouldMapVictorOpsChannelToState(t *testing.T) {
	apiKey := testAlertingChannelVictorOpsApiKey
	routingKey := testAlertingChannelVictorOpsRoutingKey
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       resourceName,
		Kind:       restapi.VictorOpsChannelType,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelVictorOps)

	channel := channels[AlertingChannelFieldChannelVictorOps].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 2)
	require.Equal(t, apiKey, channel[AlertingChannelVictorOpsFieldAPIKey])
	require.Equal(t, routingKey, channel[AlertingChannelVictorOpsFieldRoutingKey])
}

func (r *alertingChannelUnitTest) shouldMapWebhookChannelToState(t *testing.T) {
	webhookURLs := []string{"url1", "url2"}
	headers := []string{"key1", "key2:"}
	data := restapi.AlertingChannel{
		ID:          "id",
		Name:        resourceName,
		Kind:        restapi.WebhookChannelType,
		WebhookURLs: webhookURLs,
		Headers:     headers,
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelWebhook)

	channel := channels[AlertingChannelFieldChannelWebhook].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 2)
	require.Equal(t, []interface{}{"url1", "url2"}, channel[AlertingChannelWebhookFieldWebhookURLs].(*schema.Set).List())
	require.Equal(t, map[string]interface{}{
		"key1": "",
		"key2": "",
	}, channel[AlertingChannelWebhookFieldHTTPHeaders])
}

func (r *alertingChannelUnitTest) shouldMapOffice365ChannelToState(t *testing.T) {
	webhookURL := "webhookUrl"
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       resourceName,
		Kind:       restapi.Office365ChannelType,
		WebhookURL: &webhookURL,
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelOffice365)

	channel := channels[AlertingChannelFieldChannelOffice365].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 1)
	require.Equal(t, webhookURL, channel[AlertingChannelWebhookBasedFieldWebhookURL])
}

func (r *alertingChannelUnitTest) shouldMapGoogleChatChannelToState(t *testing.T) {
	webhookURL := "webhookUrl"
	data := restapi.AlertingChannel{
		ID:         "id",
		Name:       resourceName,
		Kind:       restapi.GoogleChatChannelType,
		WebhookURL: &webhookURL,
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(AlertingChannelFieldName))

	require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannel))
	require.Len(t, resourceData.Get(AlertingChannelFieldChannel).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0])

	channels := resourceData.Get(AlertingChannelFieldChannel).([]interface{})[0].(map[string]interface{})
	r.verifyChannelIsMappedToResource(t, channels, AlertingChannelFieldChannelGoogleChat)

	channel := channels[AlertingChannelFieldChannelGoogleChat].([]interface{})[0].(map[string]interface{})
	require.Len(t, channel, 1)
	require.Equal(t, webhookURL, channel[AlertingChannelWebhookBasedFieldWebhookURL])
}

func (r *alertingChannelUnitTest) verifyChannelIsMappedToResource(t *testing.T, channels map[string]interface{}, expectedChannel string) {
	require.Len(t, channels, 9)
	for k := range channels {
		require.IsType(t, []interface{}{}, channels[k])
		if k == expectedChannel {
			require.Len(t, channels[k].([]interface{}), 1)
		} else {
			require.Len(t, channels[k].([]interface{}), 0)
		}
	}
}

func (r *alertingChannelUnitTest) shouldFailToMapChannelWhenTypeIsNotValid(t *testing.T) {
	data := restapi.AlertingChannel{
		ID:     "id",
		Name:   resourceName,
		Kind:   restapi.AlertingChannelType("invalid"),
		Emails: []string{"email1", "email2"},
	}

	testHelper := NewTestHelper(t)
	sut := NewAlertingChannelResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &data, nil)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "received unsupported alerting channel of type invalid")
}

func (r *alertingChannelUnitTest) shouldMapStateOfEmailChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	emails := []string{"email1", "email2"}
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail: []interface{}{
				map[string]interface{}{
					AlertingChannelEmailFieldEmails: emails,
				},
			},
			AlertingChannelFieldChannelOpsGenie:   []interface{}{},
			AlertingChannelFieldChannelPageDuty:   []interface{}{},
			AlertingChannelFieldChannelSlack:      []interface{}{},
			AlertingChannelFieldChannelSplunk:     []interface{}{},
			AlertingChannelFieldChannelVictorOps:  []interface{}{},
			AlertingChannelFieldChannelWebhook:    []interface{}{},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.EmailChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Len(t, result.(*restapi.AlertingChannel).Emails, 2)
	require.Contains(t, result.(*restapi.AlertingChannel).Emails, "email1")
	require.Contains(t, result.(*restapi.AlertingChannel).Emails, "email2")
}

func (r *alertingChannelUnitTest) shouldMapStateOfOpsGenieChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	tags := []string{"tag1", "tag2"}
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail: []interface{}{},
			AlertingChannelFieldChannelOpsGenie: []interface{}{
				map[string]interface{}{
					AlertingChannelOpsGenieFieldAPIKey: "api-key",
					AlertingChannelOpsGenieFieldRegion: "EU",
					AlertingChannelOpsGenieFieldTags:   tags,
				},
			},
			AlertingChannelFieldChannelPageDuty:   []interface{}{},
			AlertingChannelFieldChannelSlack:      []interface{}{},
			AlertingChannelFieldChannelSplunk:     []interface{}{},
			AlertingChannelFieldChannelVictorOps:  []interface{}{},
			AlertingChannelFieldChannelWebhook:    []interface{}{},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.OpsGenieChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Equal(t, "api-key", *result.(*restapi.AlertingChannel).APIKey)
	require.Equal(t, restapi.EuOpsGenieRegion, *result.(*restapi.AlertingChannel).Region)
	require.Equal(t, "tag1,tag2", *result.(*restapi.AlertingChannel).Tags)
}

func (r *alertingChannelUnitTest) shouldMapStateOfPagerDutyChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	integrationKey := "integration key"
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:    []interface{}{},
			AlertingChannelFieldChannelOpsGenie: []interface{}{},
			AlertingChannelFieldChannelPageDuty: []interface{}{
				map[string]interface{}{
					AlertingChannelPagerDutyFieldServiceIntegrationKey: integrationKey,
				},
			},
			AlertingChannelFieldChannelSlack:      []interface{}{},
			AlertingChannelFieldChannelSplunk:     []interface{}{},
			AlertingChannelFieldChannelVictorOps:  []interface{}{},
			AlertingChannelFieldChannelWebhook:    []interface{}{},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.PagerDutyChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Equal(t, integrationKey, *result.(*restapi.AlertingChannel).ServiceIntegrationKey)
}

func (r *alertingChannelUnitTest) shouldMapStateOfSlackChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	webhookURL := testAlertingChannelSlackWebhookURL
	iconURL := testAlertingChannelSlackIconURL
	channel := testAlertingChannelSlackChannel
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:    []interface{}{},
			AlertingChannelFieldChannelOpsGenie: []interface{}{},
			AlertingChannelFieldChannelPageDuty: []interface{}{},
			AlertingChannelFieldChannelSlack: []interface{}{
				map[string]interface{}{
					AlertingChannelSlackFieldWebhookURL: webhookURL,
					AlertingChannelSlackFieldIconURL:    iconURL,
					AlertingChannelSlackFieldChannel:    channel,
				},
			},
			AlertingChannelFieldChannelSplunk:     []interface{}{},
			AlertingChannelFieldChannelVictorOps:  []interface{}{},
			AlertingChannelFieldChannelWebhook:    []interface{}{},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.SlackChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Equal(t, webhookURL, *result.(*restapi.AlertingChannel).WebhookURL)
	require.Equal(t, iconURL, *result.(*restapi.AlertingChannel).IconURL)
	require.Equal(t, channel, *result.(*restapi.AlertingChannel).Channel)
}

func (r *alertingChannelUnitTest) shouldMapStateOfSplunkChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	url := "url"
	token := "token"
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:    []interface{}{},
			AlertingChannelFieldChannelOpsGenie: []interface{}{},
			AlertingChannelFieldChannelPageDuty: []interface{}{},
			AlertingChannelFieldChannelSlack:    []interface{}{},
			AlertingChannelFieldChannelSplunk: []interface{}{
				map[string]interface{}{
					AlertingChannelSplunkFieldURL:   url,
					AlertingChannelSplunkFieldToken: token,
				},
			},
			AlertingChannelFieldChannelVictorOps:  []interface{}{},
			AlertingChannelFieldChannelWebhook:    []interface{}{},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.SplunkChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Equal(t, url, *result.(*restapi.AlertingChannel).URL)
	require.Equal(t, token, *result.(*restapi.AlertingChannel).Token)
}

func (r *alertingChannelUnitTest) shouldMapStateOfVictorOpsChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	apiKey := testAlertingChannelVictorOpsApiKey
	routingKey := testAlertingChannelVictorOpsRoutingKey
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:    []interface{}{},
			AlertingChannelFieldChannelOpsGenie: []interface{}{},
			AlertingChannelFieldChannelPageDuty: []interface{}{},
			AlertingChannelFieldChannelSlack:    []interface{}{},
			AlertingChannelFieldChannelSplunk:   []interface{}{},
			AlertingChannelFieldChannelVictorOps: []interface{}{
				map[string]interface{}{
					AlertingChannelVictorOpsFieldAPIKey:     apiKey,
					AlertingChannelVictorOpsFieldRoutingKey: routingKey,
				},
			},
			AlertingChannelFieldChannelWebhook:    []interface{}{},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.VictorOpsChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Equal(t, apiKey, *result.(*restapi.AlertingChannel).APIKey)
	require.Equal(t, routingKey, *result.(*restapi.AlertingChannel).RoutingKey)
}

func (r *alertingChannelUnitTest) shouldMapStateOfWebhookChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	webhookURLs := []string{"url1", "url2"}
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:     []interface{}{},
			AlertingChannelFieldChannelOpsGenie:  []interface{}{},
			AlertingChannelFieldChannelPageDuty:  []interface{}{},
			AlertingChannelFieldChannelSlack:     []interface{}{},
			AlertingChannelFieldChannelSplunk:    []interface{}{},
			AlertingChannelFieldChannelVictorOps: []interface{}{},
			AlertingChannelFieldChannelWebhook: []interface{}{
				map[string]interface{}{
					AlertingChannelWebhookFieldWebhookURLs: webhookURLs,
				},
			},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.WebhookChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Len(t, result.(*restapi.AlertingChannel).WebhookURLs, 2)
	require.Contains(t, result.(*restapi.AlertingChannel).WebhookURLs, "url1")
	require.Contains(t, result.(*restapi.AlertingChannel).WebhookURLs, "url2")
	require.Empty(t, result.(*restapi.AlertingChannel).Headers)
}

func (r *alertingChannelUnitTest) shouldMapStateOfWebhookChannelWithHeadersToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	webhookURLs := []string{"url1", "url2"}
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:     []interface{}{},
			AlertingChannelFieldChannelOpsGenie:  []interface{}{},
			AlertingChannelFieldChannelPageDuty:  []interface{}{},
			AlertingChannelFieldChannelSlack:     []interface{}{},
			AlertingChannelFieldChannelSplunk:    []interface{}{},
			AlertingChannelFieldChannelVictorOps: []interface{}{},
			AlertingChannelFieldChannelWebhook: []interface{}{
				map[string]interface{}{
					AlertingChannelWebhookFieldWebhookURLs: webhookURLs,
					AlertingChannelWebhookFieldHTTPHeaders: map[string]interface{}{"key1": "value1", "key2": ""},
				},
			},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.WebhookChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Len(t, result.(*restapi.AlertingChannel).WebhookURLs, 2)
	require.Contains(t, result.(*restapi.AlertingChannel).WebhookURLs, "url1")
	require.Contains(t, result.(*restapi.AlertingChannel).WebhookURLs, "url2")
	require.Len(t, result.(*restapi.AlertingChannel).Headers, 2)
	require.Contains(t, result.(*restapi.AlertingChannel).Headers, "key1: value1")
	require.Contains(t, result.(*restapi.AlertingChannel).Headers, "key2: ")
}

func (r *alertingChannelUnitTest) shouldMapStateOfOffice365ChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	webhookURL := alertingChannelWebhookBasedWebhookUrl
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:     []interface{}{},
			AlertingChannelFieldChannelOpsGenie:  []interface{}{},
			AlertingChannelFieldChannelPageDuty:  []interface{}{},
			AlertingChannelFieldChannelSlack:     []interface{}{},
			AlertingChannelFieldChannelSplunk:    []interface{}{},
			AlertingChannelFieldChannelVictorOps: []interface{}{},
			AlertingChannelFieldChannelWebhook:   []interface{}{},
			AlertingChannelFieldChannelOffice365: []interface{}{
				map[string]interface{}{
					AlertingChannelWebhookBasedFieldWebhookURL: webhookURL,
				},
			},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.Office365ChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Equal(t, webhookURL, *result.(*restapi.AlertingChannel).WebhookURL)
}

func (r *alertingChannelUnitTest) shouldMapStateOfGoogleChatChannelToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	webhookURL := alertingChannelWebhookBasedWebhookUrl
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldName, resourceName)
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:     []interface{}{},
			AlertingChannelFieldChannelOpsGenie:  []interface{}{},
			AlertingChannelFieldChannelPageDuty:  []interface{}{},
			AlertingChannelFieldChannelSlack:     []interface{}{},
			AlertingChannelFieldChannelSplunk:    []interface{}{},
			AlertingChannelFieldChannelVictorOps: []interface{}{},
			AlertingChannelFieldChannelWebhook:   []interface{}{},
			AlertingChannelFieldChannelOffice365: []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{
				map[string]interface{}{
					AlertingChannelWebhookBasedFieldWebhookURL: webhookURL,
				},
			},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingChannel{}, result)
	require.Equal(t, "id", result.GetIDForResourcePath())
	require.Equal(t, resourceName, result.(*restapi.AlertingChannel).Name)
	require.Equal(t, restapi.GoogleChatChannelType, result.(*restapi.AlertingChannel).Kind)
	require.Equal(t, webhookURL, *result.(*restapi.AlertingChannel).WebhookURL)
}

func (r *alertingChannelUnitTest) shouldFailToMapStateWhenNoChannelIsProvided(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, AlertingChannelFieldChannel, []interface{}{
		map[string]interface{}{
			AlertingChannelFieldChannelEmail:      []interface{}{},
			AlertingChannelFieldChannelOpsGenie:   []interface{}{},
			AlertingChannelFieldChannelPageDuty:   []interface{}{},
			AlertingChannelFieldChannelSlack:      []interface{}{},
			AlertingChannelFieldChannelSplunk:     []interface{}{},
			AlertingChannelFieldChannelVictorOps:  []interface{}{},
			AlertingChannelFieldChannelWebhook:    []interface{}{},
			AlertingChannelFieldChannelOffice365:  []interface{}{},
			AlertingChannelFieldChannelGoogleChat: []interface{}{},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData, nil)

	require.Error(t, err)
	require.ErrorContains(t, err, "no supported alerting channel defined")
}
