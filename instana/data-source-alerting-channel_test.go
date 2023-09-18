package instana_test

import (
	"errors"
	"fmt"
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"testing"
)

func TestAlertingChannelDataSource(t *testing.T) {
	unitTest := &dataSourceAlertingChannelUnitTest{}
	t.Run("integration test read of email alerting channel", alertingChannelEmailDataSourceIntegrationTest().testRead)
	t.Run("integration test read of ops genie alerting channel", alertingChannelOpsGenieDataSourceIntegrationTest().testRead)
	t.Run("integration test read of pager duty alerting channel", alertingChannelPagerDutyDataSourceIntegrationTest().testRead)
	t.Run("integration test read of slack 365 alerting channel", alertingChannelSlackDataSourceIntegrationTest().testRead)
	t.Run("integration test read of splunk alerting channel", alertingChannelSplunkDataSourceIntegrationTest().testRead)
	t.Run("integration test read of victor ops alerting channel", alertingChannelVictorOpsDataSourceIntegrationTest().testRead)
	t.Run("integration test read of webhook alerting channel", alertingChannelWebhookDataSourceIntegrationTest().testRead)
	t.Run("integration test read of office 365 alerting channel", alertingChannelOffice365DataSourceIntegrationTest().testRead)
	t.Run("integration test read of google chat alerting channel", alertingChannelGoogleChatDataSourceIntegrationTest().testRead)
	t.Run("schema should be valid", unitTest.schemaShouldBeValid)
	t.Run("schema version should be 0", unitTest.shouldHaveSchemaVersion0)
	t.Run("should successfully read channel", unitTest.shouldSuccessfullyReadChannel)
	t.Run("should fail to read channel when api call fails", unitTest.shouldFailToReadChannelWhenApiCallFails)
	t.Run("should fail to read channel when no channel is found for the given name", unitTest.shouldFailToReadChannelWhenNoChannelIsFoundForTheGivenName)
	t.Run("should fail to read channel when no channel type is not supported", unitTest.shouldFailToReadChannelWhenChannelTypeIsNotSupported)

}

const dataSourceAlertingChannelDefinitionPath = "data.instana_alerting_channel.example"

func alertingChannelEmailDataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666661",
		"my-email-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelEmail, AlertingChannelEmailFieldEmails), 0), "EMAIL1"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelEmail, AlertingChannelEmailFieldEmails), 1), "EMAIL2"),
		},
	)
}

func alertingChannelOpsGenieDataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666662",
		"my-ops-genie-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelOpsGenie, AlertingChannelOpsGenieFieldAPIKey), "api-key"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelOpsGenie, AlertingChannelOpsGenieFieldRegion), "EU"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelOpsGenie, AlertingChannelOpsGenieFieldTags), 0), "tag1"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelOpsGenie, AlertingChannelOpsGenieFieldTags), 1), "tag2"),
		},
	)
}

func alertingChannelPagerDutyDataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666663",
		"my-pager-duty-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelPageDuty, AlertingChannelPagerDutyFieldServiceIntegrationKey), "service-integration-key"),
		},
	)
}

func alertingChannelSlackDataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666664",
		"my-slack-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelSlack, AlertingChannelSlackFieldWebhookURL), "webhook-url"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelSlack, AlertingChannelSlackFieldIconURL), "icon-url"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelSlack, AlertingChannelSlackFieldChannel), "channel"),
		},
	)
}

func alertingChannelSplunkDataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666665",
		"my-splunk-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelSplunk, AlertingChannelSplunkFieldURL), "url"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelSplunk, AlertingChannelSplunkFieldToken), "token"),
		},
	)
}

func alertingChannelVictorOpsDataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666666",
		"my-victor-ops-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelVictorOps, AlertingChannelVictorOpsFieldAPIKey), "api-key"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelVictorOps, AlertingChannelVictorOpsFieldRoutingKey), "routing-key"),
		},
	)
}

func alertingChannelWebhookDataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666667",
		"my-webhook-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldWebhookURLs), 0), "url1"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%d", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldWebhookURLs), 1), "url2"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%s", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldHTTPHeaders), "key1"), "value1"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%s", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldHTTPHeaders), "key2"), "value2"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf("%s.%s", fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelWebhook, AlertingChannelWebhookFieldHTTPHeaders), "key3"), ""),
		},
	)
}

func alertingChannelOffice365DataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666668",
		"my-office-356-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelOffice365, AlertingChannelWebhookBasedFieldWebhookURL), "webhook-url-office-365"),
		},
	)
}

func alertingChannelGoogleChatDataSourceIntegrationTest() *dataSourceAlertingChannelIntegrationTest {
	return newDataSourceAlertingChannelIntegrationTest(
		"666669",
		"my-google-chat-channel",
		[]resource.TestCheckFunc{
			resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, fmt.Sprintf(alertingChannelChannelFieldPattern, AlertingChannelFieldChannelGoogleChat, AlertingChannelWebhookBasedFieldWebhookURL), "webhook-url-google-chat"),
		},
	)
}

func newDataSourceAlertingChannelIntegrationTest(id, channelName string, additionalChecks []resource.TestCheckFunc) *dataSourceAlertingChannelIntegrationTest {
	return &dataSourceAlertingChannelIntegrationTest{
		id:               id,
		channelName:      channelName,
		additionalChecks: additionalChecks,
	}
}

type dataSourceAlertingChannelIntegrationTest struct {
	id               string
	channelName      string
	additionalChecks []resource.TestCheckFunc
}

func (r *dataSourceAlertingChannelIntegrationTest) testRead(t *testing.T) {
	serverResponse := `
[{
	"id": "666661",
	"name": "my-email-channel",
	"kind": "EMAIL",
	"emails": [ "EMAIL1", "EMAIL2" ]
},{
	"id": "666662",
	"name": "my-ops-genie-channel",
	"kind": "OPS_GENIE",
	"apiKey": "api-key",
	"region": "EU",
	"tags": "tag1, tag2"
},{
	"id": "666663",
	"name": "my-pager-duty-channel",
	"kind": "PAGER_DUTY",
	"serviceIntegrationKey": "service-integration-key"
},{
	"id": "666664",
	"name": "my-slack-channel",
	"kind": "SLACK",
	"webhookUrl": "webhook-url",
	"iconUrl": "icon-url",
	"channel": "channel"
},{
	"id": "666665",
	"name": "my-splunk-channel",
	"kind": "SPLUNK",
	"url": "url",
	"token": "token"
},{
	"id": "666666",
	"name": "my-victor-ops-channel",
	"kind": "VICTOR_OPS",
	"apiKey": "api-key",
	"routingKey": "routing-key"
},{
	"id": "666667",
	"name": "my-webhook-channel",
	"kind": "WEB_HOOK",
	"webhookUrls": [ "url1", "url2" ],
	"headers": [ "key1: value1", "key2: value2", "key3" ]
},{
	"id"     	 : "666668",
	"name"   	 : "my-office-356-channel",
	"kind"   	 : "OFFICE_365",
	"webhookUrl" : "webhook-url-office-365"
},{
	"id"     	 : "666669",
	"name"   	 : "my-google-chat-channel",
	"kind"   	 : "GOOGLE_CHAT",
	"webhookUrl" : "webhook-url-google-chat"
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

func (r *dataSourceAlertingChannelIntegrationTest) createTestStep(httpPort int) resource.TestStep {
	dataSourceAlertingChannelDefinitionTemplate := `
data "instana_alerting_channel" "example" {
  name = "%s"
}
`
	config := appendProviderConfig(fmt.Sprintf(dataSourceAlertingChannelDefinitionTemplate, r.channelName), httpPort)
	checks := append([]resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, "id", r.id),
		resource.TestCheckResourceAttr(dataSourceAlertingChannelDefinitionPath, AlertingChannelFieldName, r.channelName),
	}, r.additionalChecks...)
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

func (r *dataSourceAlertingChannelUnitTest) shouldSuccessfullyReadChannel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		data := restapi.AlertingChannel{
			ID:     "id",
			Name:   resourceName,
			Kind:   restapi.EmailChannelType,
			Emails: []string{"email1", "email2"},
		}

		AlertingChannelAPI := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)
		AlertingChannelAPI.EXPECT().GetAll().Times(1).Return(&[]*restapi.AlertingChannel{&data}, nil)
		mockInstanaApi.EXPECT().AlertingChannels().Return(AlertingChannelAPI).Times(1)

		sut := NewAlertingChannelDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			AlertingChannelFieldName: resourceName,
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.Nil(t, diag)
		require.Equal(t, data.ID, resourceData.Id())
		require.Equal(t, data.Name, resourceData.Get(AlertingChannelFieldName))

		require.IsType(t, []interface{}{}, resourceData.Get(AlertingChannelFieldChannelEmail))
		require.Len(t, resourceData.Get(AlertingChannelFieldChannelEmail).([]interface{}), 1)
		require.IsType(t, map[string]interface{}{}, resourceData.Get(AlertingChannelFieldChannelEmail).([]interface{})[0])

		channel := resourceData.Get(AlertingChannelFieldChannelEmail).([]interface{})[0].(map[string]interface{})
		require.Len(t, channel, 1)
		emails := channel[AlertingChannelEmailFieldEmails].(*schema.Set)
		require.Equal(t, 2, emails.Len())
		require.Contains(t, emails.List(), "email1")
		require.Contains(t, emails.List(), "email2")
	})
}

func (r *dataSourceAlertingChannelUnitTest) shouldFailToReadChannelWhenApiCallFails(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		expectedError := errors.New("test")

		AlertingChannelAPI := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)
		AlertingChannelAPI.EXPECT().GetAll().Times(1).Return(nil, expectedError)
		mockInstanaApi.EXPECT().AlertingChannels().Return(AlertingChannelAPI).Times(1)

		sut := NewAlertingChannelDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			AlertingChannelFieldName: resourceName,
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.NotNil(t, diag)
		require.True(t, diag.HasError())
		require.Contains(t, diag[0].Summary, expectedError.Error())
	})
}

func (r *dataSourceAlertingChannelUnitTest) shouldFailToReadChannelWhenNoChannelIsFoundForTheGivenName(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		data := restapi.AlertingChannel{
			ID:     "id",
			Name:   "other name",
			Kind:   restapi.EmailChannelType,
			Emails: []string{"email1", "email2"},
		}

		AlertingChannelAPI := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)
		AlertingChannelAPI.EXPECT().GetAll().Times(1).Return(&[]*restapi.AlertingChannel{&data}, nil)
		mockInstanaApi.EXPECT().AlertingChannels().Return(AlertingChannelAPI).Times(1)

		sut := NewAlertingChannelDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			AlertingChannelFieldName: resourceName,
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.NotNil(t, diag)
		require.True(t, diag.HasError())
		require.Contains(t, diag[0].Summary, "no alerting channel found")
	})
}

func (r *dataSourceAlertingChannelUnitTest) shouldFailToReadChannelWhenChannelTypeIsNotSupported(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanaApi *mocks.MockInstanaAPI) {
		data := restapi.AlertingChannel{
			ID:   "id",
			Name: resourceName,
			Kind: restapi.AlertingChannelType("invalid"),
		}

		AlertingChannelAPI := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)
		AlertingChannelAPI.EXPECT().GetAll().Times(1).Return(&[]*restapi.AlertingChannel{&data}, nil)
		mockInstanaApi.EXPECT().AlertingChannels().Return(AlertingChannelAPI).Times(1)

		sut := NewAlertingChannelDataSource().CreateResource()
		resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{
			AlertingChannelFieldName: resourceName,
		})

		diag := sut.ReadContext(nil, resourceData, meta)

		require.NotNil(t, diag)
		require.True(t, diag.HasError())
		require.Contains(t, diag[0].Summary, "received unsupported alerting channel of type invalid")
	})
}
