package instana_test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceAlertingConfigTerraformTemplateWithRuleIds = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_config" "rule_ids" {
  alert_name = "name {{ITERATOR}}"
  integration_ids = [ "integration_id1", "integration_id2" ]
  event_filter_query = "query"
  event_filter_rule_ids = [ "rule-1", "rule-2" ]
}
`

const alertingConfigServerResponseTemplateWithRuleIds = `
{
	"id" : "{{id}}",
	"alertName" : "prefix name {{ITERATOR}} suffix",
	"integrationIds" : [ "integration_id2", "integration_id1" ],
	"eventFilteringConfiguration" : {
		"query" : "query",
		"ruleIds" : [ "rule-2", "rule-1" ]
	}
}
`

const resourceAlertingConfigTerraformTemplateWithEventTypes = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_config" "event_types" {
  alert_name = "name {{ITERATOR}}"
  integration_ids = [ "integration_id1", "integration_id2" ]
  event_filter_query = "query"
  event_filter_event_types = [ "incident", "critical" ]
}
`

const alertingConfigServerResponseTemplateWithEventTypes = `
{
	"id" : "{{id}}",
	"alertName" : "prefix name {{ITERATOR}} suffix",
	"integrationIds" : [ "integration_id2", "integration_id1" ],
	"eventFilteringConfiguration" : {
		"query" : "query",
		"eventTypes" : [ "critical", "incident" ]
	}
}
`

const iteratorPlaceholder = "{{ITERATOR}}"
const alertingConfigApiPath = restapi.AlertsResourcePath + "/{id}"
const testAlertingConfigDefinitionWithRuleIds = "instana_alerting_config.rule_ids"
const testAlertingConfigDefinitionWithEventTypes = "instana_alerting_config.event_types"

func TestCRUDOfAlertingConfigurationWithRuleIds(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingConfigApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := restapi.AlertsResourcePath + "/" + vars["id"]
		json := strings.ReplaceAll(strings.ReplaceAll(alertingConfigServerResponseTemplateWithRuleIds, "{{id}}", vars["id"]), "{{ITERATOR}}", strconv.Itoa(getZeroBasedCallCount(httpServer, http.MethodPut, path)))
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingConfigTerraformTemplateWithRuleIds, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	rule1 := "rule-1"
	rule2 := "rule-2"
	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithRuleIds, 0),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, 0), rule1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, 1), rule2),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithRuleIds, 1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, 0), rule1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, 1), rule2),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestCRUDOfAlertingConfigurationWithEventTypes(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingConfigApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := restapi.AlertsResourcePath + "/" + vars["id"]
		json := strings.ReplaceAll(strings.ReplaceAll(alertingConfigServerResponseTemplateWithEventTypes, "{{id}}", vars["id"]), "{{ITERATOR}}", strconv.Itoa(getZeroBasedCallCount(httpServer, http.MethodPut, path)))
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingConfigTerraformTemplateWithEventTypes, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithEventTypes, 0),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, 1), string(restapi.IncidentAlertEventType)),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, 0), string(restapi.CriticalAlertEventType)),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithEventTypes, 1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, 1), string(restapi.IncidentAlertEventType)),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, 0), string(restapi.CriticalAlertEventType)),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func CreateTestCheckFunctionForComonResourceAttributes(config string, iteration int) resource.TestCheckFunc {
	integrationId1 := "integration_id1"
	integrationId2 := "integration_id2"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(config, "id"),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldAlertName, fmt.Sprintf("name %d", iteration)),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldFullAlertName, fmt.Sprintf("prefix name %d suffix", iteration)),
		resource.TestCheckResourceAttr(config, fmt.Sprintf("%s.%d", AlertingConfigFieldIntegrationIds, 0), integrationId1),
		resource.TestCheckResourceAttr(config, fmt.Sprintf("%s.%d", AlertingConfigFieldIntegrationIds, 1), integrationId2),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldEventFilterQuery, "query"),
	)
}

func TestResourceAlertingConfigDefinition(t *testing.T) {
	resource := NewAlertingConfigResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingConfigFieldAlertName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingConfigFieldFullAlertName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingConfigFieldIntegrationIds)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingConfigFieldEventFilterQuery)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeSetOfStrings(AlertingConfigFieldEventFilterEventTypes)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeSetOfStrings(AlertingConfigFieldEventFilterRuleIDs)
}

func TestShouldReturnCorrectResourceNameForAlertingConfig(t *testing.T) {
	name := NewAlertingConfigResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_alerting_config", name, "Expected resource name to be instana_alerting_config")
}

func TestAlertingConfigShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewAlertingConfigResourceHandle().MetaData().SchemaVersion)
}

func TestAlertingConfigShouldHaveOneStateUpgraderForVersionZero(t *testing.T) {
	resourceHandler := NewAlertingConfigResourceHandle()

	require.Equal(t, 1, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
}

func TestShouldReturnStateOfAlertingConfigWithRuleIdsUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[AlertingConfigFieldAlertName] = "name"
	rawData[AlertingConfigFieldFullAlertName] = "fullname"
	rawData[AlertingConfigFieldIntegrationIds] = []interface{}{"integration-id1", "integration-id2"}
	rawData[AlertingConfigFieldEventFilterQuery] = "filter"
	rawData[AlertingConfigFieldEventFilterRuleIDs] = []interface{}{"rule-id1", "rule-id2"}
	meta := "dummy"
	ctx := context.Background()

	result, err := NewAlertingConfigResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, rawData, result)
}

func TestShouldReturnStateOfAlertingConfigWithEventTypesUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[AlertingConfigFieldAlertName] = "name"
	rawData[AlertingConfigFieldFullAlertName] = "fullname"
	rawData[AlertingConfigFieldIntegrationIds] = []interface{}{"integration-id1", "integration-id2"}
	rawData[AlertingConfigFieldEventFilterQuery] = "filter"
	rawData[AlertingConfigFieldEventFilterEventTypes] = []interface{}{"incident", "critical"}
	meta := "dummy"
	ctx := context.Background()

	result, err := NewAlertingConfigResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, rawData, result)
}

const (
	alertingConfigID             = "alerting-id"
	alertingConfigName           = "alerting-name"
	alertingConfigFullName       = "prefix alerting-name suffix"
	alertingConfigIntegrationId1 = "alerting-integration-id1"
	alertingConfigIntegrationId2 = "alerting-integration-id2"
	alertingConfigRuleId1        = "alerting-rule-id1"
	alertingConfigRuleId2        = "alerting-rule-id2"
	alertingConfigQuery          = "alerting-query"
)

func TestShouldUpdateResourceStateForAlertingConfigWithRuleIds(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	query := alertingConfigQuery

	data := restapi.AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigFullName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:   &query,
			RuleIDs: []string{alertingConfigRuleId1, alertingConfigRuleId2},
		},
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, alertingConfigID, resourceData.Id())
	require.Equal(t, alertingConfigName, resourceData.Get(AlertingConfigFieldAlertName))
	require.Equal(t, alertingConfigFullName, resourceData.Get(AlertingConfigFieldFullAlertName))
	require.Equal(t, alertingConfigQuery, resourceData.Get(AlertingConfigFieldEventFilterQuery))
	requireIntegrationIdOFAlertingConfigResourceDataUpdated(t, resourceData)

	ruleIDs := resourceData.Get(AlertingConfigFieldEventFilterRuleIDs).(*schema.Set)
	requireSetMatchesToValues(t, ruleIDs, alertingConfigRuleId1, alertingConfigRuleId2)
}

func TestShouldUpdateResourceStateForAlertingConfigWithEventTypes(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	query := alertingConfigQuery

	data := restapi.AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigFullName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:      &query,
			EventTypes: []restapi.AlertEventType{restapi.IncidentAlertEventType, restapi.CriticalAlertEventType},
		},
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, alertingConfigID, resourceData.Id())
	require.Equal(t, alertingConfigName, resourceData.Get(AlertingConfigFieldAlertName))
	require.Equal(t, alertingConfigFullName, resourceData.Get(AlertingConfigFieldFullAlertName))
	require.Equal(t, alertingConfigQuery, resourceData.Get(AlertingConfigFieldEventFilterQuery))
	requireIntegrationIdOFAlertingConfigResourceDataUpdated(t, resourceData)

	eventTypes := resourceData.Get(AlertingConfigFieldEventFilterEventTypes).(*schema.Set)
	requireSetMatchesToValues(t, eventTypes, string(restapi.CriticalAlertEventType), string(restapi.IncidentAlertEventType))
}

func requireIntegrationIdOFAlertingConfigResourceDataUpdated(t *testing.T, resourceData *schema.ResourceData) {
	integrationIDs := resourceData.Get(AlertingConfigFieldIntegrationIds).(*schema.Set)
	requireSetMatchesToValues(t, integrationIDs, alertingConfigIntegrationId1, alertingConfigIntegrationId2)
}

func requireSetMatchesToValues(t *testing.T, set *schema.Set, values ...string) {
	require.Equal(t, len(values), set.Len())
	for _, v := range values {
		require.Contains(t, set.List(), v)
	}
}

func TestShouldConvertStateOfAlertingConfigToDataModelWithRuleIds(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingConfigResourceHandle()
	integrationIds := []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2}
	ruleIds := []string{alertingConfigRuleId1, alertingConfigRuleId2}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(alertingConfigID)
	resourceData.Set(AlertingConfigFieldAlertName, alertingConfigName)
	resourceData.Set(AlertingConfigFieldFullAlertName, alertingConfigName)
	resourceData.Set(AlertingConfigFieldIntegrationIds, integrationIds)
	resourceData.Set(AlertingConfigFieldEventFilterQuery, alertingConfigQuery)
	resourceData.Set(AlertingConfigFieldEventFilterRuleIDs, ruleIds)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingConfiguration{}, model)
	require.Equal(t, alertingConfigID, model.GetIDForResourcePath())
	require.Equal(t, alertingConfigName, model.(*restapi.AlertingConfiguration).AlertName)

	requireIntegrationIdOFAlertingConfigModel(t, model.(*restapi.AlertingConfiguration))
	require.Equal(t, alertingConfigQuery, *model.(*restapi.AlertingConfiguration).EventFilteringConfiguration.Query)
	requireSliceValuesMatchesToValues(t, model.(*restapi.AlertingConfiguration).EventFilteringConfiguration.RuleIDs, alertingConfigRuleId1, alertingConfigRuleId2)
}

func TestShouldConvertStateOfAlertingConfigToDataModelWithEventTypes(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingConfigResourceHandle()
	integrationIds := []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(alertingConfigID)
	resourceData.Set(AlertingConfigFieldAlertName, alertingConfigName)
	resourceData.Set(AlertingConfigFieldFullAlertName, alertingConfigName)
	resourceData.Set(AlertingConfigFieldIntegrationIds, integrationIds)
	resourceData.Set(AlertingConfigFieldEventFilterQuery, alertingConfigQuery)
	resourceData.Set(AlertingConfigFieldEventFilterEventTypes, []string{"incident", "critical"})

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingConfiguration{}, model)
	require.Equal(t, alertingConfigID, model.GetIDForResourcePath())
	require.Equal(t, alertingConfigName, model.(*restapi.AlertingConfiguration).AlertName)

	requireIntegrationIdOFAlertingConfigModel(t, model.(*restapi.AlertingConfiguration))
	require.Equal(t, alertingConfigQuery, *model.(*restapi.AlertingConfiguration).EventFilteringConfiguration.Query)

	eventTypes := model.(*restapi.AlertingConfiguration).EventFilteringConfiguration.EventTypes
	require.Len(t, eventTypes, 2)
	require.Contains(t, eventTypes, restapi.CriticalAlertEventType)
	require.Contains(t, eventTypes, restapi.IncidentAlertEventType)
}

func requireIntegrationIdOFAlertingConfigModel(t *testing.T, model *restapi.AlertingConfiguration) {
	requireSliceValuesMatchesToValues(t, model.IntegrationIDs, alertingConfigIntegrationId1, alertingConfigIntegrationId2)
}

func requireSliceValuesMatchesToValues(t *testing.T, data []string, values ...string) {
	require.Equal(t, len(values), len(data))
	for _, v := range values {
		require.Contains(t, data, v)
	}
}
