package instana_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

var testAlertingConfigProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

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
	"alertName" : "prefix name suffix",
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
	"alertName" : "prefix name suffix",
	"integrationIds" : [ "integration_id2", "integration_id1" ],
	"eventFilteringConfiguration" : {
		"query" : "query",
		"eventTypes" : [ "critical", "incident" ]
	}
}
`

const iteratorPlaceholder = "{{ITERATOR}}"
const contentType = "Content-Type"
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
		json := strings.ReplaceAll(alertingConfigServerResponseTemplateWithRuleIds, "{{id}}", vars["id"])
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
	hashFunctionRules := schema.HashSchema(AlertingConfigSchemaEventFilterRuleIDs.Elem.(*schema.Schema))
	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingConfigProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithRuleIds, 0),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, hashFunctionRules(rule1)), rule1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, hashFunctionRules(rule2)), rule2),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithRuleIds, 1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, hashFunctionRules(rule1)), rule1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, hashFunctionRules(rule2)), rule2),
				),
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
		json := strings.ReplaceAll(alertingConfigServerResponseTemplateWithEventTypes, "{{id}}", vars["id"])
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingConfigTerraformTemplateWithEventTypes, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	hashFunctionEventTypes := schema.HashSchema(AlertingConfigSchemaEventFilterEventTypes.Elem.(*schema.Schema))
	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingConfigProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithEventTypes, 0),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, hashFunctionEventTypes(string(restapi.IncidentAlertEventType))), string(restapi.IncidentAlertEventType)),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, hashFunctionEventTypes(string(restapi.CriticalAlertEventType))), string(restapi.CriticalAlertEventType)),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithEventTypes, 1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, hashFunctionEventTypes(string(restapi.IncidentAlertEventType))), string(restapi.IncidentAlertEventType)),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, hashFunctionEventTypes(string(restapi.CriticalAlertEventType))), string(restapi.CriticalAlertEventType)),
				),
			},
		},
	})
}

func CreateTestCheckFunctionForComonResourceAttributes(config string, iteration int) resource.TestCheckFunc {
	integrationId1 := "integration_id1"
	integrationId2 := "integration_id2"
	hashFunctionIntegrationIds := schema.HashSchema(AlertingConfigSchemaIntegrationIds.Elem.(*schema.Schema))
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(config, "id"),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldAlertName, fmt.Sprintf("name %d", iteration)),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldFullAlertName, fmt.Sprintf("prefix name %d suffix", iteration)),
		resource.TestCheckResourceAttr(config, fmt.Sprintf("%s.%d", AlertingConfigFieldIntegrationIds, hashFunctionIntegrationIds(integrationId1)), integrationId1),
		resource.TestCheckResourceAttr(config, fmt.Sprintf("%s.%d", AlertingConfigFieldIntegrationIds, hashFunctionIntegrationIds(integrationId2)), integrationId2),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldEventFilterQuery, "query"),
	)
}

func TestResourceAlertingConfigDefinition(t *testing.T) {
	resource := NewAlertingConfigResourceHandle()

	schemaMap := resource.Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingConfigFieldAlertName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingConfigFieldFullAlertName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingConfigFieldIntegrationIds)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingConfigFieldEventFilterQuery)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeSetOfStrings(AlertingConfigFieldEventFilterEventTypes)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeSetOfStrings(AlertingConfigFieldEventFilterRuleIDs)
}

func TestShouldReturnCorrectResourceNameForAlertingConfig(t *testing.T) {
	name := NewAlertingConfigResourceHandle().ResourceName

	assert.Equal(t, "instana_alerting_config", name, "Expected resource name to be instana_alerting_config")
}

func TestAlertingConfigShouldHaveSchemaVersionOne(t *testing.T) {
	assert.Equal(t, 1, NewAlertingConfigResourceHandle().SchemaVersion)
}

func TestAlertingConfigShouldHaveOneStateUpgraderForVersionZero(t *testing.T) {
	resourceHandler := NewAlertingConfigResourceHandle()

	assert.Equal(t, 1, len(resourceHandler.StateUpgraders))
	assert.Equal(t, 0, resourceHandler.StateUpgraders[0].Version)
}

func TestShouldReturnStateOfAlertingConfigWithRuleIdsUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[AlertingConfigFieldAlertName] = "name"
	rawData[AlertingConfigFieldFullAlertName] = "fullname"
	rawData[AlertingConfigFieldIntegrationIds] = []interface{}{"integration-id1", "integration-id2"}
	rawData[AlertingConfigFieldEventFilterQuery] = "filter"
	rawData[AlertingConfigFieldEventFilterRuleIDs] = []interface{}{"rule-id1", "rule-id2"}
	meta := "dummy"

	result, err := NewAlertingConfigResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Equal(t, rawData, result)
}

func TestShouldReturnStateOfAlertingConfigWithEventTypesUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[AlertingConfigFieldAlertName] = "name"
	rawData[AlertingConfigFieldFullAlertName] = "fullname"
	rawData[AlertingConfigFieldIntegrationIds] = []interface{}{"integration-id1", "integration-id2"}
	rawData[AlertingConfigFieldEventFilterQuery] = "filter"
	rawData[AlertingConfigFieldEventFilterEventTypes] = []interface{}{"incident", "critical"}
	meta := "dummy"

	result, err := NewAlertingConfigResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Equal(t, rawData, result)
}

const (
	alertingConfigID             = "alerting-id"
	alertingConfigName           = "alerting-name"
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
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:   &query,
			RuleIDs: []string{alertingConfigRuleId1, alertingConfigRuleId2},
		},
	}

	err := resourceHandle.UpdateState(resourceData, data)

	assert.Nil(t, err)
	assert.Equal(t, alertingConfigID, resourceData.Id())
	assert.Equal(t, alertingConfigName, resourceData.Get(AlertingConfigFieldFullAlertName))

	integrationIDs := resourceData.Get(AlertingConfigFieldIntegrationIds).(*schema.Set)
	assert.Equal(t, 2, integrationIDs.Len())
	assert.Contains(t, integrationIDs.List(), alertingConfigIntegrationId1)
	assert.Contains(t, integrationIDs.List(), alertingConfigIntegrationId2)

	assert.Equal(t, alertingConfigQuery, resourceData.Get(AlertingConfigFieldEventFilterQuery))

	ruleIDs := resourceData.Get(AlertingConfigFieldEventFilterRuleIDs).(*schema.Set)
	assert.Equal(t, 2, ruleIDs.Len())
	assert.Contains(t, ruleIDs.List(), alertingConfigRuleId1)
	assert.Contains(t, ruleIDs.List(), alertingConfigRuleId2)
}

func TestShouldUpdateResourceStateForAlertingConfigWithEventTypes(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	query := alertingConfigQuery

	data := restapi.AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:      &query,
			EventTypes: []restapi.AlertEventType{restapi.IncidentAlertEventType, restapi.CriticalAlertEventType},
		},
	}

	err := resourceHandle.UpdateState(resourceData, data)

	assert.Nil(t, err)
	assert.Equal(t, alertingConfigID, resourceData.Id())
	assert.Equal(t, alertingConfigName, resourceData.Get(AlertingConfigFieldFullAlertName))

	integrationIDs := resourceData.Get(AlertingConfigFieldIntegrationIds).(*schema.Set)
	assert.Equal(t, 2, integrationIDs.Len())
	assert.Contains(t, integrationIDs.List(), alertingConfigIntegrationId1)
	assert.Contains(t, integrationIDs.List(), alertingConfigIntegrationId2)

	assert.Equal(t, alertingConfigQuery, resourceData.Get(AlertingConfigFieldEventFilterQuery))

	eventTypes := resourceData.Get(AlertingConfigFieldEventFilterEventTypes).(*schema.Set)
	assert.Equal(t, 2, eventTypes.Len())
	assert.Contains(t, eventTypes.List(), string(restapi.CriticalAlertEventType))
	assert.Contains(t, eventTypes.List(), string(restapi.IncidentAlertEventType))
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

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.AlertingConfiguration{}, model)
	assert.Equal(t, alertingConfigID, model.GetID())
	assert.Equal(t, alertingConfigName, model.(restapi.AlertingConfiguration).AlertName)

	mappedIntegrationIds := model.(restapi.AlertingConfiguration).IntegrationIDs
	assert.Len(t, mappedIntegrationIds, 2)
	assert.Contains(t, mappedIntegrationIds, alertingConfigIntegrationId1)
	assert.Contains(t, mappedIntegrationIds, alertingConfigIntegrationId2)

	assert.Equal(t, alertingConfigQuery, *model.(restapi.AlertingConfiguration).EventFilteringConfiguration.Query)

	mappedRuleIds := model.(restapi.AlertingConfiguration).EventFilteringConfiguration.RuleIDs
	assert.Len(t, mappedRuleIds, 2)
	assert.Contains(t, mappedRuleIds, alertingConfigRuleId1)
	assert.Contains(t, mappedRuleIds, alertingConfigRuleId2)
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

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.AlertingConfiguration{}, model)
	assert.Equal(t, alertingConfigID, model.GetID())
	assert.Equal(t, alertingConfigName, model.(restapi.AlertingConfiguration).AlertName)

	mappedIntegrationIds := model.(restapi.AlertingConfiguration).IntegrationIDs
	assert.Len(t, mappedIntegrationIds, 2)
	assert.Contains(t, mappedIntegrationIds, alertingConfigIntegrationId1)
	assert.Contains(t, mappedIntegrationIds, alertingConfigIntegrationId2)

	assert.Equal(t, alertingConfigQuery, *model.(restapi.AlertingConfiguration).EventFilteringConfiguration.Query)

	eventTypes := model.(restapi.AlertingConfiguration).EventFilteringConfiguration.EventTypes
	assert.Len(t, eventTypes, 2)
	assert.Contains(t, eventTypes, restapi.CriticalAlertEventType)
	assert.Contains(t, eventTypes, restapi.IncidentAlertEventType)
}
