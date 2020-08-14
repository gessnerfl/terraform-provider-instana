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

var testCustomEventSpecificationWithThresholdRuleProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name {{ITERATION}}"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rule_severity = "warning"
  rule_metric_name = "metric_name"
  rule_rollup = "40000"
  rule_condition_operator = "=="
  rule_condition_value = "1.2"
}
`

const resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name {{ITERATION}}"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
  rule_severity = "warning"
  rule_metric_name = "metric_name"
  rule_window = 60000
  rule_aggregation = "sum"
  rule_condition_operator = "=="
  rule_condition_value = 1.2
}
`

const (
	customEventSpecificationWithThresholdRuleApiPath        = restapi.CustomEventSpecificationResourcePath + "/{id}"
	testCustomEventSpecificationWithThresholdRuleDefinition = "instana_custom_event_spec_threshold_rule.example"

	customEventSpecificationWithThresholdRuleID                = "custom-system-event-id"
	customEventSpecificationWithThresholdRuleName              = "name"
	customEventSpecificationWithThresholdRuleEntityType        = "entity_type"
	customEventSpecificationWithThresholdRuleQuery             = "query"
	customEventSpecificationWithThresholdRuleExpirationTime    = 60000
	customEventSpecificationWithThresholdRuleDescription       = "description"
	customEventSpecificationWithThresholdRuleMetricName        = "metric_name"
	customEventSpecificationWithThresholdRuleRollup            = 40000
	customEventSpecificationWithThresholdRuleWindow            = 60000
	customEventSpecificationWithThresholdRuleAggregation       = restapi.AggregationSum
	customEventSpecificationWithThresholdRuleConditionOperator = restapi.ConditionOperatorEquals
	customEventSpecificationWithThresholdRuleConditionValue    = float64(1.2)
)

var CustomEventSpecificationWithThresholdRuleRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithRollupResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName" : "metric_name", "rollup" : 40000, "conditionOperator" : "==", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate,
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldRollup, strconv.FormatInt(customEventSpecificationWithThresholdRuleRollup, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(customEventSpecificationWithThresholdRuleConditionOperator)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithWindowResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName": "metric_name", "window" : 60000, "aggregation": "sum", "conditionOperator" : "==", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate,
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(customEventSpecificationWithThresholdRuleConditionOperator)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

const httpServerResponseTemplate = `
{
	"id" : "{{id}}",
	"name" : "prefix name suffix",
	"entityType" : "entity_type",
	"query" : "query",
	"enabled" : true,
	"triggering" : true,
	"description" : "description",
	"expirationTime" : 60000,
	"rules" : [ {{rule}} ]
}
`

func testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(t *testing.T, terraformDefinition, ruleAsJson string, ruleTestCheckFunctions ...resource.TestCheckFunc) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, customEventSpecificationWithThresholdRuleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, customEventSpecificationWithThresholdRuleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, customEventSpecificationWithThresholdRuleApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(strings.ReplaceAll(httpServerResponseTemplate, "{{id}}", vars["id"]), "{{rule}}", ruleAsJson)
		w.Header().Set(constSystemEventContentType, r.Header.Get(constSystemEventContentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	completeTerraformDefinitionWithoutName := strings.ReplaceAll(terraformDefinition, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	completeTerraformDefinitionWithName1 := strings.ReplaceAll(completeTerraformDefinitionWithoutName, "{{ITERATION}}", "0")
	completeTerraformDefinitionWithName2 := strings.ReplaceAll(completeTerraformDefinitionWithoutName, "{{ITERATION}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testCustomEventSpecificationWithThresholdRuleProviders,
		Steps: []resource.TestStep{
			{
				Config: completeTerraformDefinitionWithName1,
				Check:  resource.ComposeTestCheckFunc(createTestCheckFunctions(ruleTestCheckFunctions, 0)...),
			},
			{
				Config: completeTerraformDefinitionWithName2,
				Check:  resource.ComposeTestCheckFunc(createTestCheckFunctions(ruleTestCheckFunctions, 1)...),
			},
		},
	})
}

func createTestCheckFunctions(ruleTestCheckFunctions []resource.TestCheckFunc, iteration int) []resource.TestCheckFunc {
	defaultCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithThresholdRuleDefinition, "id"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldName, customEventSpecificationWithThresholdRuleName+fmt.Sprintf(" %d", iteration)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldFullName, fmt.Sprintf("%s %s %d %s", "prefix", customEventSpecificationWithThresholdRuleName, iteration, "suffix")),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldTriggering, "true"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEventSpecificationWithThresholdRuleExpirationTime)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldEnabled, "true"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationRuleSeverity, CustomEventSpecificationWithThresholdRuleRuleSeverity),
	}
	allFunctions := append(defaultCheckFunctions, ruleTestCheckFunctions...)
	return allFunctions
}

func TestCustomEventSpecificationWithThresholdRuleSchemaDefinitionIsValid(t *testing.T) {
	schema := NewCustomEventSpecificationWithThresholdRuleResourceHandle().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleSeverity)

	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ThresholdRuleFieldMetricName)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(ThresholdRuleFieldWindow)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(ThresholdRuleFieldRollup)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ThresholdRuleFieldAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ThresholdRuleFieldConditionOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeFloat(ThresholdRuleFieldConditionValue)
}

func TestCustomEventSpecificationWithThresholdRuleResourceShouldHaveSchemaVersionTwo(t *testing.T) {
	assert.Equal(t, 2, NewCustomEventSpecificationWithThresholdRuleResourceHandle().SchemaVersion)
}

func TestCustomEventSpecificationWithThresholdRuleShouldHaveTwoStateUpgraderForVersionZeroAndOne(t *testing.T) {
	resourceHandler := NewCustomEventSpecificationWithThresholdRuleResourceHandle()

	assert.Equal(t, 2, len(resourceHandler.StateUpgraders))
	assert.Equal(t, 0, resourceHandler.StateUpgraders[0].Version)
	assert.Equal(t, 1, resourceHandler.StateUpgraders[1].Version)
}

func TestShouldMigrateCustomEventSpecificationWithThresholdRuleStateAndAddFullNameWithSameValueAsNameWhenMigratingFromVersion0To1(t *testing.T) {
	name := "Test Name"
	rawData := make(map[string]interface{})
	rawData[CustomEventSpecificationFieldName] = name
	meta := "dummy"

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Equal(t, name, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateEmptyCustomEventSpecificationWithThresholdRuleStateFromVersion0To1(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Nil(t, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateCustomEventSpecificationWithThresholdRuleStateToVersion2WhenDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData["downstream_integration_ids"] = []interface{}{"id1", "id2"}
	rawData["downstream_broadcast_to_all_alerting_configs"] = true
	meta := "dummy"

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders[1].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Nil(t, result["downstream_integration_ids"])
	assert.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldMigrateCustomEventSpecificationWithThresholdRuleStateToVersion2WhenNoDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"

	result, err := NewCustomEventSpecificationWithThresholdRuleResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Nil(t, result["downstream_integration_ids"])
	assert.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldReturnCorrectResourceNameForCustomEventSpecificationWithThresholdRuleResource(t *testing.T) {
	name := NewCustomEventSpecificationWithThresholdRuleResourceHandle().ResourceName

	assert.Equal(t, name, "instana_custom_event_spec_threshold_rule")
}

func TestShouldUpdateCustomEventSpecificationWithThresholdRuleTerraformStateFromApiObject(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery

	window := customEventSpecificationWithThresholdRuleWindow
	rollup := customEventSpecificationWithThresholdRuleRollup
	aggregation := customEventSpecificationWithThresholdRuleAggregation
	conditionValue := customEventSpecificationWithThresholdRuleConditionValue
	metricName := customEventSpecificationWithThresholdRuleMetricName
	conditionOperator := customEventSpecificationWithThresholdRuleConditionOperator

	spec := restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           customEventSpecificationWithThresholdRuleName,
		EntityType:     customEventSpecificationWithThresholdRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:             restapi.ThresholdRuleType,
				Severity:          restapi.SeverityWarning.GetAPIRepresentation(),
				MetricName:        &metricName,
				Window:            &window,
				Rollup:            &rollup,
				Aggregation:       &aggregation,
				ConditionOperator: &conditionOperator,
				ConditionValue:    &conditionValue,
			},
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationWithThresholdRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec)

	assert.Nil(t, err)
	assert.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	assert.Equal(t, customEventSpecificationWithThresholdRuleName, resourceData.Get(CustomEventSpecificationFieldFullName))
	assert.Equal(t, customEventSpecificationWithThresholdRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	assert.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	assert.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	assert.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	assert.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	assert.Equal(t, metricName, resourceData.Get(ThresholdRuleFieldMetricName))
	assert.Equal(t, window, resourceData.Get(ThresholdRuleFieldWindow))
	assert.Equal(t, rollup, resourceData.Get(ThresholdRuleFieldRollup))
	assert.Equal(t, string(aggregation), resourceData.Get(ThresholdRuleFieldAggregation))
	assert.Equal(t, string(conditionOperator), resourceData.Get(ThresholdRuleFieldConditionOperator))
	assert.Equal(t, conditionValue, resourceData.Get(ThresholdRuleFieldConditionValue))
	assert.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), resourceData.Get(CustomEventSpecificationRuleSeverity))
}

func TestShouldSuccessfullyConvertCustomEventSpecificationWithThresholdRuleStateToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithThresholdRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	resourceData.Set(CustomEventSpecificationFieldFullName, customEventSpecificationWithThresholdRuleName)
	resourceData.Set(CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	resourceData.Set(CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	resourceData.Set(CustomEventSpecificationFieldTriggering, true)
	resourceData.Set(CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	resourceData.Set(CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	resourceData.Set(CustomEventSpecificationFieldEnabled, true)
	resourceData.Set(CustomEventSpecificationRuleSeverity, restapi.SeverityWarning.GetTerraformRepresentation())
	resourceData.Set(ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName)
	resourceData.Set(ThresholdRuleFieldWindow, customEventSpecificationWithThresholdRuleWindow)
	resourceData.Set(ThresholdRuleFieldRollup, customEventSpecificationWithThresholdRuleRollup)
	resourceData.Set(ThresholdRuleFieldAggregation, customEventSpecificationWithThresholdRuleAggregation)
	resourceData.Set(ThresholdRuleFieldConditionOperator, customEventSpecificationWithThresholdRuleConditionOperator)
	resourceData.Set(ThresholdRuleFieldConditionValue, customEventSpecificationWithThresholdRuleConditionValue)

	result, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(restapi.CustomEventSpecification)
	assert.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetID())
	assert.Equal(t, customEventSpecificationWithThresholdRuleName, customEventSpec.Name)
	assert.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	assert.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	assert.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	assert.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	assert.True(t, customEventSpec.Triggering)
	assert.True(t, customEventSpec.Enabled)

	assert.Equal(t, 1, len(customEventSpec.Rules))
	assert.Equal(t, customEventSpecificationWithThresholdRuleMetricName, *customEventSpec.Rules[0].MetricName)
	assert.Equal(t, customEventSpecificationWithThresholdRuleWindow, *customEventSpec.Rules[0].Window)
	assert.Equal(t, customEventSpecificationWithThresholdRuleRollup, *customEventSpec.Rules[0].Rollup)
	assert.Equal(t, customEventSpecificationWithThresholdRuleAggregation, *customEventSpec.Rules[0].Aggregation)
	assert.Equal(t, customEventSpecificationWithThresholdRuleConditionOperator, *customEventSpec.Rules[0].ConditionOperator)
	assert.Equal(t, customEventSpecificationWithThresholdRuleConditionValue, *customEventSpec.Rules[0].ConditionValue)
	assert.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func TestShouldFailToConvertCustomEventSpecificationWithThresholdRuleStateToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithThresholdRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.Set(CustomEventSpecificationRuleSeverity, "INVALID")

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.NotNil(t, err)
}
