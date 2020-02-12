package instana_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
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
  custom_payload = "custom"
  event_filter_query = "query"
  event_filter_rule_ids = [ "rule-1", "rule-2" ]
}
`

const alertingConfigServerResponseTemplateWithRuleIds = `
{
	"id" : "{{id}}",
	"alertName" : "prefix name suffix",
	"customPayload" : "custom",
	"integrationIds" : [ "integration_id1", "integration_id2" ],
	"eventFilteringConfiguration" : {
		"query" : "query",
		"ruleIds" : [ "rule-1", "rule-2" ]
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
  custom_payload = "custom"
  event_filter_query = "query"
  event_filter_event_types = [ "incident", "critical" ]
}
`

const alertingConfigServerResponseTemplateWithEventTypes = `
{
	"id" : "{{id}}",
	"alertName" : "prefix name suffix",
	"customPayload" : "custom",
	"integrationIds" : [ "integration_id1", "integration_id2" ],
	"eventFilteringConfiguration" : {
		"query" : "query",
		"eventTypes" : [ "incident", "critical" ]
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

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingConfigProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithRuleIds, 0),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, AlertingConfigFieldEventFilterRuleIDs+".0", "rule-1"),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, AlertingConfigFieldEventFilterRuleIDs+".1", "rule-2"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithRuleIds, 1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, AlertingConfigFieldEventFilterRuleIDs+".0", "rule-1"),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithRuleIds, AlertingConfigFieldEventFilterRuleIDs+".1", "rule-2"),
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

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingConfigProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithEventTypes, 0),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, AlertingConfigFieldEventFilterEventTypes+".0", "incident"),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, AlertingConfigFieldEventFilterEventTypes+".1", "critical"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					CreateTestCheckFunctionForComonResourceAttributes(testAlertingConfigDefinitionWithEventTypes, 1),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, AlertingConfigFieldEventFilterEventTypes+".0", "incident"),
					resource.TestCheckResourceAttr(testAlertingConfigDefinitionWithEventTypes, AlertingConfigFieldEventFilterEventTypes+".1", "critical"),
				),
			},
		},
	})
}

func CreateTestCheckFunctionForComonResourceAttributes(config string, iteration int) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(config, "id"),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldAlertName, fmt.Sprintf("name %d", iteration)),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldFullAlertName, fmt.Sprintf("prefix name %d suffix", iteration)),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldIntegrationIds+".0", "integration_id1"),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldIntegrationIds+".1", "integration_id2"),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldCustomPayload, "custom"),
		resource.TestCheckResourceAttr(config, AlertingConfigFieldEventFilterQuery, "query"),
	)
}

func TestResourceAlertingConfigDefinition(t *testing.T) {
	resource := NewAlertingConfigResourceHandle()

	schemaMap := resource.Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingConfigFieldAlertName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingConfigFieldFullAlertName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfStrings(AlertingConfigFieldIntegrationIds)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingConfigFieldCustomPayload)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingConfigFieldEventFilterQuery)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfStrings(AlertingConfigFieldEventFilterEventTypes)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfStrings(AlertingConfigFieldEventFilterRuleIDs)
}

func TestAlertingConfigShouldHaveSchemaVersionZero(t *testing.T) {
	assert.Equal(t, 0, NewAlertingConfigResourceHandle().SchemaVersion)
}

func TestAlertingConfigShouldHaveNoStateUpgrader(t *testing.T) {
	assert.Equal(t, 0, len(NewAlertingConfigResourceHandle().StateUpgraders))
}

func TestShouldReturnCorrectResourceNameForAlertingConfig(t *testing.T) {
	name := NewAlertingConfigResourceHandle().ResourceName

	assert.Equal(t, "instana_alerting_config", name, "Expected resource name to be instana_alerting_config")
}

const (
	alertingConfigID             = "alerting-id"
	alertingConfigName           = "alerting-name"
	alertingConfigIntegrationId1 = "alerting-integration-id1"
	alertingConfigIntegrationId2 = "alerting-integration-id2"
	alertingConfigRuleId1        = "alerting-rule-id1"
	alertingConfigRuleId2        = "alerting-rule-id2"
	alertingConfigQuery          = "alerting-query"
	alertingConfigCustomPayload  = "alerting-custom-payload"
)

func TestShouldUpdateResourceStateForAlertingConfigWithRuleIds(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	query := alertingConfigQuery
	customPayload := alertingConfigCustomPayload

	data := restapi.AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		CustomPayload:  &customPayload,
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:   &query,
			RuleIDs: []string{alertingConfigRuleId1, alertingConfigRuleId2},
		},
	}

	err := resourceHandle.UpdateState(resourceData, data)

	assert.Nil(t, err)
	assert.Equal(t, alertingConfigID, resourceData.Id())
	assert.Equal(t, alertingConfigName, resourceData.Get(AlertingConfigFieldFullAlertName))
	assert.Equal(t, []interface{}{alertingConfigIntegrationId1, alertingConfigIntegrationId2}, resourceData.Get(AlertingConfigFieldIntegrationIds))
	assert.Equal(t, alertingConfigCustomPayload, resourceData.Get(AlertingConfigFieldCustomPayload))
	assert.Equal(t, alertingConfigQuery, resourceData.Get(AlertingConfigFieldEventFilterQuery))
	assert.Equal(t, []interface{}{alertingConfigRuleId1, alertingConfigRuleId2}, resourceData.Get(AlertingConfigFieldEventFilterRuleIDs))
}

func TestShouldUpdateResourceStateForAlertingConfigWithEventTypes(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	query := alertingConfigQuery
	customPayload := alertingConfigCustomPayload

	data := restapi.AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		CustomPayload:  &customPayload,
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:      &query,
			EventTypes: []restapi.AlertEventType{restapi.IncidentAlertEventType, restapi.CriticalAlertEventType},
		},
	}

	err := resourceHandle.UpdateState(resourceData, data)

	assert.Nil(t, err)
	assert.Equal(t, alertingConfigID, resourceData.Id())
	assert.Equal(t, alertingConfigName, resourceData.Get(AlertingConfigFieldFullAlertName))
	assert.Equal(t, []interface{}{alertingConfigIntegrationId1, alertingConfigIntegrationId2}, resourceData.Get(AlertingConfigFieldIntegrationIds))
	assert.Equal(t, alertingConfigCustomPayload, resourceData.Get(AlertingConfigFieldCustomPayload))
	assert.Equal(t, alertingConfigQuery, resourceData.Get(AlertingConfigFieldEventFilterQuery))
	assert.Equal(t, []interface{}{string(restapi.IncidentAlertEventType), string(restapi.CriticalAlertEventType)}, resourceData.Get(AlertingConfigFieldEventFilterEventTypes))
}

func TestShouldConvertStateOfAlertingConfigToDataModelWithRuleIds(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingConfigResourceHandle()
	integrationIds := []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2}
	ruleIds := []string{alertingConfigRuleId1, alertingConfigRuleId1}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(alertingConfigID)
	resourceData.Set(AlertingConfigFieldAlertName, alertingConfigName)
	resourceData.Set(AlertingConfigFieldFullAlertName, alertingConfigName)
	resourceData.Set(AlertingConfigFieldIntegrationIds, integrationIds)
	resourceData.Set(AlertingConfigFieldCustomPayload, alertingConfigCustomPayload)
	resourceData.Set(AlertingConfigFieldEventFilterQuery, alertingConfigQuery)
	resourceData.Set(AlertingConfigFieldEventFilterRuleIDs, ruleIds)

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.AlertingConfiguration{}, model)
	assert.Equal(t, alertingConfigID, model.GetID())
	assert.Equal(t, alertingConfigName, model.(restapi.AlertingConfiguration).AlertName)
	assert.Equal(t, integrationIds, model.(restapi.AlertingConfiguration).IntegrationIDs)
	assert.Equal(t, alertingConfigCustomPayload, *model.(restapi.AlertingConfiguration).CustomPayload)
	assert.Equal(t, alertingConfigQuery, *model.(restapi.AlertingConfiguration).EventFilteringConfiguration.Query)
	assert.Equal(t, ruleIds, model.(restapi.AlertingConfiguration).EventFilteringConfiguration.RuleIDs)
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
	resourceData.Set(AlertingConfigFieldCustomPayload, alertingConfigCustomPayload)
	resourceData.Set(AlertingConfigFieldEventFilterQuery, alertingConfigQuery)
	resourceData.Set(AlertingConfigFieldEventFilterEventTypes, []string{"incident", "critical"})

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.AlertingConfiguration{}, model)
	assert.Equal(t, alertingConfigID, model.GetID())
	assert.Equal(t, alertingConfigName, model.(restapi.AlertingConfiguration).AlertName)
	assert.Equal(t, integrationIds, model.(restapi.AlertingConfiguration).IntegrationIDs)
	assert.Equal(t, alertingConfigCustomPayload, *model.(restapi.AlertingConfiguration).CustomPayload)
	assert.Equal(t, alertingConfigQuery, *model.(restapi.AlertingConfiguration).EventFilteringConfiguration.Query)
	assert.Equal(t, []restapi.AlertEventType{restapi.IncidentAlertEventType, restapi.CriticalAlertEventType}, model.(restapi.AlertingConfiguration).EventFilteringConfiguration.EventTypes)
}
