package instana_test

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

/*
const resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate = `
resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name %d"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rule_severity = "warning"
  rule_metric_name = "metric_name"
  rule_rollup = "40000"
  rule_condition_operator = "="
  rule_condition_value = "1.2"
}
`

const resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate = `
resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name %d"
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
  rule_condition_operator = "{{CONDITION_OPERATOR}}"
  rule_condition_value = 1.2
}
`

const resourceCustomEventSpecificationWithThresholdRuleAndMetricPatternDefinitionTemplate = `
resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name %d"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
  rule_severity = "warning"
  rule_window = 60000
  rule_aggregation = "sum"
  rule_condition_operator = "="
  rule_condition_value = 1.2
  rule_metric_pattern_prefix = "prefix"
  rule_metric_pattern_postfix = "postfix"
  rule_metric_pattern_placeholder = "placeholder"
  rule_metric_pattern_operator = "startsWith"
}
`

const (
	testCustomEventSpecificationWithThresholdRuleDefinition = "instana_custom_event_spec_threshold_rule.example"

	customEventSpecificationWithThresholdRuleID             = "custom-system-event-id"
	customEventSpecificationWithThresholdRuleEntityType     = "entity_type"
	customEventSpecificationWithThresholdRuleQuery          = "query"
	customEventSpecificationWithThresholdRuleExpirationTime = 60000
	customEventSpecificationWithThresholdRuleDescription    = "description"
	customEventSpecificationWithThresholdRuleMetricName     = "metric_name"
	customEventSpecificationWithThresholdRuleRollup         = 40000
	customEventSpecificationWithThresholdRuleWindow         = 60000
	customEventSpecificationWithThresholdRuleAggregation    = restapi.AggregationSum
	customEventSpecificationWithThresholdRuleConditionValue = float64(1.2)
)

var CustomEventSpecificationWithThresholdRuleRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithRollupResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName" : "metric_name", "rollup" : 40000, "conditionOperator" : "=", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate,
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldRollup, strconv.FormatInt(customEventSpecificationWithThresholdRuleRollup, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(restapi.ConditionOperatorEquals.InstanaAPIValue())),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithWindowResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName": "metric_name", "window" : 60000, "aggregation": "sum", "conditionOperator" : "=", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		strings.ReplaceAll(resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate, "{{CONDITION_OPERATOR}}", "="),
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(restapi.ConditionOperatorEquals.InstanaAPIValue())),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithWindowAndAlternativeConditionOperatorRepresentationResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName": "metric_name", "window" : 60000, "aggregation": "sum", "conditionOperator" : "=", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		strings.ReplaceAll(resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate, "{{CONDITION_OPERATOR}}", "=="),
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(restapi.ConditionOperatorEquals.InstanaAPIValue())),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithMetricPatternResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "window" : 60000, "aggregation": "sum", "conditionOperator" : "=", "conditionValue" : 1.2, "metricPattern" : { "prefix" : "prefix", "postfix" : "postfix", "placeholder" : "placeholder", "operator" : "startsWith" } }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		resourceCustomEventSpecificationWithThresholdRuleAndMetricPatternDefinitionTemplate,
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(restapi.ConditionOperatorEquals.InstanaAPIValue())),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricPatternPrefix, "prefix"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricPatternPostfix, "postfix"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricPatternPlaceholder, "placeholder"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricPatternOperator, string(restapi.MetricPatternOperatorTypeStartsWith)),
	)
}

const httpServerResponseTemplate = `
{
	"id" : "%s",
	"name" : "prefix name %d suffix",
	"entityType" : "entity_type",
	"query" : "query",
	"enabled" : true,
	"triggering" : true,
	"description" : "description",
	"expirationTime" : 60000,
	"rules" : [ %s ]
}
`

func testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(t *testing.T, terraformDefinition, ruleAsJson string, ruleTestCheckFunctions ...resource.TestCheckFunc) {
	httpServer := createMockHttpServerForResource(restapi.CustomEventSpecificationResourcePath, httpServerResponseTemplate, ruleAsJson)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: appendProviderConfig(fmt.Sprintf(terraformDefinition, 0), httpServer.GetPort()),
				Check:  resource.ComposeTestCheckFunc(createTestCheckFunctions(ruleTestCheckFunctions, 0)...),
			},
			testStepImport(testCustomEventSpecificationWithThresholdRuleDefinition),
			{
				Config: appendProviderConfig(fmt.Sprintf(terraformDefinition, 1), httpServer.GetPort()),
				Check:  resource.ComposeTestCheckFunc(createTestCheckFunctions(ruleTestCheckFunctions, 1)...),
			},
			testStepImport(testCustomEventSpecificationWithThresholdRuleDefinition),
		},
	})
}

func createTestCheckFunctions(ruleTestCheckFunctions []resource.TestCheckFunc, iteration int) []resource.TestCheckFunc {
	defaultCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithThresholdRuleDefinition, "id"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldName, formatResourceName(iteration)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldFullName, formatResourceFullName(iteration)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldTriggering, trueAsString),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEventSpecificationWithThresholdRuleExpirationTime)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldEnabled, trueAsString),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationRuleSeverity, CustomEventSpecificationWithThresholdRuleRuleSeverity),
	}
	allFunctions := append(defaultCheckFunctions, ruleTestCheckFunctions...)
	return allFunctions
}
*/

func TestCustomEventSpecificationResource(t *testing.T) {
	unitTest := &customerEventSpecificationUnitTest{}
	t.Run("schema should be valid", unitTest.schemaShouldBeValid)
	t.Run("should have schema version 0", unitTest.shouldHaveSchemaVersion0)
	t.Run("should have no state upgrader", unitTest.shouldHaveNoStateUpgraders)
	t.Run("should have correct resource name", unitTest.shouldHaveCorrectResourceName)
	t.Run("should map entity verification rule to state", unitTest.shouldMapEntityVerificationRuleToState)
	t.Run("should map system rule to state", unitTest.shouldMapSystemRuleToState)
	t.Run("should map threshold rule and metric name to state", unitTest.shouldMapThresholdRuleAndMetricNameToState)
	t.Run("should map threshold rule and metric pattern to state", unitTest.shouldMapThresholdRuleAndMetricPatternToState)
	t.Run("should fail to map rule when severity is not valid", unitTest.shouldFailToMapRuleWhenSeverityIsNotValid)
	t.Run("should fail to map rule when rule type is not valid", unitTest.shouldFailToMapRuleWhenRuleTypeIsNotValid)
	t.Run("should map state of entity verification rule to data model", unitTest.shouldMapStateOfEntityVerificationRuleToDataModel)
	t.Run("should fail to map state of entity verification rule when severity is not valid", unitTest.shouldFailToMapStateOfEntityVerificationRuleToDataModelWhenSeverityIsNotValid)
	t.Run("should map state of system rule to data model", unitTest.shouldMapStateOfSystemRuleToDataModel)
	t.Run("should fail to map state of system rule when severity is not valid", unitTest.shouldFailToMapStateOfSystemRuleToDataModelWhenSeverityIsNotValid)
	t.Run("should map state of threshold rule with metric name to data model", unitTest.shouldMapStateOfThresholdRuleWithMetricNameToDataModel)
	t.Run("should map state of threshold rule with metric pattern to data model", unitTest.shouldMapStateOfThresholdRuleWithMetricPatternToDataModel)
	t.Run("should fail to map state of threshold rule when severity is not valid", unitTest.shouldFailToMapStateOfThresholdRuleToDataModelWhenSeverityIsNotValid)
	t.Run("should fail to map state when no rule is provided", unitTest.shouldFailToMapStateWhenNoRuleIsProvided)
}

type customerEventSpecificationUnitTest struct{}

func (r *customerEventSpecificationUnitTest) schemaShouldBeValid(t *testing.T) {
	schemaData := NewCustomEventSpecificationResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaData, t)
	require.Len(t, schemaData, 8)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfResource(CustomEventSpecificationFieldRule)

	r.validateRuleSchema(t, schemaData[CustomEventSpecificationFieldRule].Elem.(*schema.Resource).Schema)
}

func (r *customerEventSpecificationUnitTest) validateRuleSchema(t *testing.T, ruleSchema map[string]*schema.Schema) {
	schemaAssert := testutils.NewTerraformSchemaAssert(ruleSchema, t)
	require.Len(t, ruleSchema, 3)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(CustomEventSpecificationFieldEntityVerificationRule)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(CustomEventSpecificationFieldSystemRule)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(CustomEventSpecificationFieldThresholdRule)

	r.validateEntityVerificationRuleSchema(t, ruleSchema[CustomEventSpecificationFieldEntityVerificationRule].Elem.(*schema.Resource).Schema)
	r.validateSystemRuleSchema(t, ruleSchema[CustomEventSpecificationFieldSystemRule].Elem.(*schema.Resource).Schema)
	r.validateThresholdRuleSchema(t, ruleSchema[CustomEventSpecificationFieldThresholdRule].Elem.(*schema.Resource).Schema)
}

func (r *customerEventSpecificationUnitTest) validateEntityVerificationRuleSchema(t *testing.T, ruleSchema map[string]*schema.Schema) {
	require.Len(t, ruleSchema, 5)
	schemaAssert := testutils.NewTerraformSchemaAssert(ruleSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleFieldSeverity)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeInt(CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration)
}

func (r *customerEventSpecificationUnitTest) validateSystemRuleSchema(t *testing.T, ruleSchema map[string]*schema.Schema) {
	require.Len(t, ruleSchema, 2)
	schemaAssert := testutils.NewTerraformSchemaAssert(ruleSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleFieldSeverity)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationSystemRuleFieldSystemRuleId)
}

func (r *customerEventSpecificationUnitTest) validateThresholdRuleSchema(t *testing.T, ruleSchema map[string]*schema.Schema) {
	require.Len(t, ruleSchema, 8)
	schemaAssert := testutils.NewTerraformSchemaAssert(ruleSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleFieldSeverity)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricName)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(CustomEventSpecificationThresholdRuleFieldMetricPattern)
	r.validateMetricPatternSchema(t, ruleSchema[CustomEventSpecificationThresholdRuleFieldMetricPattern].Elem.(*schema.Resource).Schema)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationThresholdRuleFieldRollup)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeInt(CustomEventSpecificationThresholdRuleFieldWindow)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldConditionOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeFloat(CustomEventSpecificationThresholdRuleFieldConditionValue)
}

func (r *customerEventSpecificationUnitTest) validateMetricPatternSchema(t *testing.T, metricPatternSchema map[string]*schema.Schema) {
	require.Len(t, metricPatternSchema, 4)
	schemaAssert := testutils.NewTerraformSchemaAssert(metricPatternSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationThresholdRuleFieldMetricPatternOperator)
}

func (r *customerEventSpecificationUnitTest) shouldHaveSchemaVersion0(t *testing.T) {
	require.Equal(t, 0, NewCustomEventSpecificationResourceHandle().MetaData().SchemaVersion)
}

func (r *customerEventSpecificationUnitTest) shouldHaveNoStateUpgraders(t *testing.T) {
	resourceHandler := NewCustomEventSpecificationResourceHandle()

	require.Equal(t, 0, len(resourceHandler.StateUpgraders()))
}

func (r *customerEventSpecificationUnitTest) shouldHaveCorrectResourceName(t *testing.T) {
	name := NewCustomEventSpecificationResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_custom_event_specification")
}

func (r *customerEventSpecificationUnitTest) shouldMapEntityVerificationRuleToState(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery
	matchingEntityLabel := "matching-entity-label"
	matchingEntityType := "matching-entity-type"
	matchingOperator := "is"
	offlineDuration := 1234

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceName,
		EntityType:     EntityVerificationRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:               restapi.EntityVerificationRuleType,
				Severity:            restapi.SeverityWarning.GetAPIRepresentation(),
				MatchingEntityLabel: &matchingEntityLabel,
				MatchingEntityType:  &matchingEntityType,
				MatchingOperator:    &matchingOperator,
				OfflineDuration:     &offlineDuration,
			},
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, EntityVerificationRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.IsType(t, []interface{}{}, resourceData.Get(CustomEventSpecificationFieldRule))
	require.Len(t, resourceData.Get(CustomEventSpecificationFieldRule).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(CustomEventSpecificationFieldRule).([]interface{})[0])

	rules := resourceData.Get(CustomEventSpecificationFieldRule).([]interface{})[0].(map[string]interface{})
	require.Len(t, rules, 3)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule])
	require.Len(t, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{})[0])
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldSystemRule])
	require.Len(t, rules[CustomEventSpecificationFieldSystemRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldThresholdRule])
	require.Len(t, rules[CustomEventSpecificationFieldThresholdRule].([]interface{}), 0)

	rule := rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{})[0].(map[string]interface{})
	require.Len(t, rule, 5)
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), rule[CustomEventSpecificationRuleFieldSeverity])
	require.Equal(t, matchingEntityType, rule[CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType])
	require.Equal(t, matchingEntityLabel, rule[CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel])
	require.Equal(t, matchingOperator, rule[CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator])
	require.Equal(t, offlineDuration, rule[CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration])
}

func (r *customerEventSpecificationUnitTest) shouldMapSystemRuleToState(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery
	systemRuleId := "system-rule-id"

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceName,
		EntityType:     SystemRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:        restapi.SystemRuleType,
				Severity:     restapi.SeverityWarning.GetAPIRepresentation(),
				SystemRuleID: &systemRuleId,
			},
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, SystemRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.IsType(t, []interface{}{}, resourceData.Get(CustomEventSpecificationFieldRule))
	require.Len(t, resourceData.Get(CustomEventSpecificationFieldRule).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(CustomEventSpecificationFieldRule).([]interface{})[0])

	rules := resourceData.Get(CustomEventSpecificationFieldRule).([]interface{})[0].(map[string]interface{})
	require.Len(t, rules, 3)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule])
	require.Len(t, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldSystemRule])
	require.Len(t, rules[CustomEventSpecificationFieldSystemRule].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rules[CustomEventSpecificationFieldSystemRule].([]interface{})[0])
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldThresholdRule])
	require.Len(t, rules[CustomEventSpecificationFieldThresholdRule].([]interface{}), 0)

	rule := rules[CustomEventSpecificationFieldSystemRule].([]interface{})[0].(map[string]interface{})
	require.Len(t, rule, 2)
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), rule[CustomEventSpecificationRuleFieldSeverity])
	require.Equal(t, systemRuleId, rule[CustomEventSpecificationSystemRuleFieldSystemRuleId])
}

func (r *customerEventSpecificationUnitTest) shouldMapThresholdRuleAndMetricNameToState(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery

	window := customEventSpecificationWithThresholdRuleWindow
	rollup := customEventSpecificationWithThresholdRuleRollup
	aggregation := customEventSpecificationWithThresholdRuleAggregation
	conditionValue := customEventSpecificationWithThresholdRuleConditionValue
	metricName := customEventSpecificationWithThresholdRuleMetricName
	conditionOperator := restapi.ConditionOperatorEquals.InstanaAPIValue()

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceName,
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
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.IsType(t, []interface{}{}, resourceData.Get(CustomEventSpecificationFieldRule))
	require.Len(t, resourceData.Get(CustomEventSpecificationFieldRule).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(CustomEventSpecificationFieldRule).([]interface{})[0])

	rules := resourceData.Get(CustomEventSpecificationFieldRule).([]interface{})[0].(map[string]interface{})
	require.Len(t, rules, 3)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule])
	require.Len(t, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldSystemRule])
	require.Len(t, rules[CustomEventSpecificationFieldSystemRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldThresholdRule])
	require.Len(t, rules[CustomEventSpecificationFieldThresholdRule].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rules[CustomEventSpecificationFieldThresholdRule].([]interface{})[0])

	rule := rules[CustomEventSpecificationFieldThresholdRule].([]interface{})[0].(map[string]interface{})
	require.Len(t, rule, 8)
	require.Equal(t, metricName, rule[CustomEventSpecificationThresholdRuleFieldMetricName])
	require.Equal(t, []interface{}{}, rule[CustomEventSpecificationThresholdRuleFieldMetricPattern])
	require.Equal(t, window, rule[CustomEventSpecificationThresholdRuleFieldWindow])
	require.Equal(t, rollup, rule[CustomEventSpecificationThresholdRuleFieldRollup])
	require.Equal(t, string(aggregation), rule[CustomEventSpecificationThresholdRuleFieldAggregation])
	require.Equal(t, conditionOperator, rule[CustomEventSpecificationThresholdRuleFieldConditionOperator])
	require.Equal(t, conditionValue, rule[CustomEventSpecificationThresholdRuleFieldConditionValue])
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), rule[CustomEventSpecificationRuleFieldSeverity])
}

func (r *customerEventSpecificationUnitTest) shouldMapThresholdRuleAndMetricPatternToState(t *testing.T) {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery

	window := customEventSpecificationWithThresholdRuleWindow
	rollup := customEventSpecificationWithThresholdRuleRollup
	aggregation := customEventSpecificationWithThresholdRuleAggregation
	conditionValue := customEventSpecificationWithThresholdRuleConditionValue
	conditionOperator := restapi.ConditionOperatorEquals.InstanaAPIValue()
	prefix := "prefix"
	postfix := "postfix"
	placeholder := "placeholder"
	operator := restapi.MetricPatternOperatorTypeStartsWith

	spec := &restapi.CustomEventSpecification{
		ID:             customEventSpecificationWithThresholdRuleID,
		Name:           resourceName,
		EntityType:     customEventSpecificationWithThresholdRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			{
				DType:    restapi.ThresholdRuleType,
				Severity: restapi.SeverityWarning.GetAPIRepresentation(),
				MetricPattern: &restapi.MetricPattern{
					Prefix:      prefix,
					Postfix:     &postfix,
					Placeholder: &placeholder,
					Operator:    operator,
				},
				Window:            &window,
				Rollup:            &rollup,
				Aggregation:       &aggregation,
				ConditionOperator: &conditionOperator,
				ConditionValue:    &conditionValue,
			},
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, resourceData.Id())
	require.Equal(t, resourceName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.IsType(t, []interface{}{}, resourceData.Get(CustomEventSpecificationFieldRule))
	require.Len(t, resourceData.Get(CustomEventSpecificationFieldRule).([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, resourceData.Get(CustomEventSpecificationFieldRule).([]interface{})[0])

	rules := resourceData.Get(CustomEventSpecificationFieldRule).([]interface{})[0].(map[string]interface{})
	require.Len(t, rules, 3)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldEntityVerificationRule])
	require.Len(t, rules[CustomEventSpecificationFieldEntityVerificationRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldSystemRule])
	require.Len(t, rules[CustomEventSpecificationFieldSystemRule].([]interface{}), 0)
	require.IsType(t, []interface{}{}, rules[CustomEventSpecificationFieldThresholdRule])
	require.Len(t, rules[CustomEventSpecificationFieldThresholdRule].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rules[CustomEventSpecificationFieldThresholdRule].([]interface{})[0])

	rule := rules[CustomEventSpecificationFieldThresholdRule].([]interface{})[0].(map[string]interface{})
	require.Len(t, rule, 8)
	require.Equal(t, "", rule[CustomEventSpecificationThresholdRuleFieldMetricName])
	require.Equal(t, window, rule[CustomEventSpecificationThresholdRuleFieldWindow])
	require.Equal(t, rollup, rule[CustomEventSpecificationThresholdRuleFieldRollup])
	require.Equal(t, string(aggregation), rule[CustomEventSpecificationThresholdRuleFieldAggregation])
	require.Equal(t, conditionOperator, rule[CustomEventSpecificationThresholdRuleFieldConditionOperator])
	require.Equal(t, conditionValue, rule[CustomEventSpecificationThresholdRuleFieldConditionValue])
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), rule[CustomEventSpecificationRuleFieldSeverity])

	require.IsType(t, []interface{}{}, rule[CustomEventSpecificationThresholdRuleFieldMetricPattern])
	require.Len(t, rule[CustomEventSpecificationThresholdRuleFieldMetricPattern].([]interface{}), 1)
	require.IsType(t, map[string]interface{}{}, rule[CustomEventSpecificationThresholdRuleFieldMetricPattern].([]interface{})[0])

	metricPatternData := rule[CustomEventSpecificationThresholdRuleFieldMetricPattern].([]interface{})[0].(map[string]interface{})

	require.Len(t, metricPatternData, 4)
	require.Equal(t, prefix, metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix])
	require.Equal(t, postfix, metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix])
	require.Equal(t, placeholder, metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder])
	require.Equal(t, string(operator), metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternOperator])
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapRuleWhenSeverityIsNotValid(t *testing.T) {
	spec := &restapi.CustomEventSpecification{
		Rules: []restapi.RuleSpecification{
			{
				DType:    restapi.ThresholdRuleType,
				Severity: 123,
			},
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "is not a valid severity")
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapRuleWhenRuleTypeIsNotValid(t *testing.T) {
	spec := &restapi.CustomEventSpecification{
		Rules: []restapi.RuleSpecification{
			{
				DType:    "invalid",
				Severity: restapi.SeverityWarning.GetAPIRepresentation(),
			},
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unsupported rule type invalid")
}

func (r *customerEventSpecificationUnitTest) shouldMapStateOfEntityVerificationRuleToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	matchingEntityLabel := "matching-entity-label"
	matchingEntityType := "matching-entity-type"
	matchingOperator := "is"
	offlineDuration := 1234

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, resourceName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRule, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity:                              restapi.SeverityWarning.GetTerraformRepresentation(),
					CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel: matchingEntityLabel,
					CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType:  matchingEntityType,
					CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator:    matchingOperator,
					CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration:     offlineDuration,
				}},
			CustomEventSpecificationFieldSystemRule:    []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(*restapi.CustomEventSpecification)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, restapi.EntityVerificationRuleType, customEventSpec.Rules[0].DType)
	require.Equal(t, matchingEntityLabel, *customEventSpec.Rules[0].MatchingEntityLabel)
	require.Equal(t, matchingEntityType, *customEventSpec.Rules[0].MatchingEntityType)
	require.Equal(t, matchingOperator, *customEventSpec.Rules[0].MatchingOperator)
	require.Equal(t, offlineDuration, *customEventSpec.Rules[0].OfflineDuration)
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapStateOfEntityVerificationRuleToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRule, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity: "invalid",
				}},
			CustomEventSpecificationFieldSystemRule:    []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Error(t, err)
	require.ErrorContains(t, err, "invalid is not a valid severity")
}

func (r *customerEventSpecificationUnitTest) shouldMapStateOfSystemRuleToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	systemRuleId := "system-rule-id"

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, resourceName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRule, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity:           restapi.SeverityWarning.GetTerraformRepresentation(),
					CustomEventSpecificationSystemRuleFieldSystemRuleId: systemRuleId,
				}},
			CustomEventSpecificationFieldThresholdRule: []interface{}{},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(*restapi.CustomEventSpecification)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, restapi.SystemRuleType, customEventSpec.Rules[0].DType)
	require.Equal(t, systemRuleId, *customEventSpec.Rules[0].SystemRuleID)
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapStateOfSystemRuleToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRule, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity: "invalid",
				}},
			CustomEventSpecificationFieldThresholdRule: []interface{}{},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Error(t, err)
	require.ErrorContains(t, err, "invalid is not a valid severity")
}

func (r *customerEventSpecificationUnitTest) shouldMapStateOfThresholdRuleWithMetricNameToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, resourceName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRule, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule:             []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity:                   restapi.SeverityWarning.GetTerraformRepresentation(),
					CustomEventSpecificationThresholdRuleFieldMetricName:        customEventSpecificationWithThresholdRuleMetricName,
					CustomEventSpecificationThresholdRuleFieldMetricPattern:     []interface{}{},
					CustomEventSpecificationThresholdRuleFieldRollup:            customEventSpecificationWithThresholdRuleRollup,
					CustomEventSpecificationThresholdRuleFieldWindow:            customEventSpecificationWithThresholdRuleWindow,
					CustomEventSpecificationThresholdRuleFieldAggregation:       customEventSpecificationWithThresholdRuleAggregation,
					CustomEventSpecificationThresholdRuleFieldConditionOperator: restapi.ConditionOperatorEquals.InstanaAPIValue(),
					CustomEventSpecificationThresholdRuleFieldConditionValue:    customEventSpecificationWithThresholdRuleConditionValue,
				},
			},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(*restapi.CustomEventSpecification)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, restapi.ThresholdRuleType, customEventSpec.Rules[0].DType)
	require.Equal(t, customEventSpecificationWithThresholdRuleMetricName, *customEventSpec.Rules[0].MetricName)
	require.Nil(t, customEventSpec.Rules[0].MetricPattern)
	require.Equal(t, customEventSpecificationWithThresholdRuleWindow, *customEventSpec.Rules[0].Window)
	require.Equal(t, customEventSpecificationWithThresholdRuleRollup, *customEventSpec.Rules[0].Rollup)
	require.Equal(t, customEventSpecificationWithThresholdRuleAggregation, *customEventSpec.Rules[0].Aggregation)
	require.Equal(t, restapi.ConditionOperatorEquals.InstanaAPIValue(), *customEventSpec.Rules[0].ConditionOperator)
	require.Equal(t, customEventSpecificationWithThresholdRuleConditionValue, *customEventSpec.Rules[0].ConditionValue)
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func (r *customerEventSpecificationUnitTest) shouldMapStateOfThresholdRuleWithMetricPatternToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	prefix := "prefix"
	postfix := "postfix"
	placeholder := "placeholder"
	operator := restapi.MetricPatternOperatorTypeStartsWith

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, resourceName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEventSpecificationWithThresholdRuleExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRule, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule:             []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity: restapi.SeverityWarning.GetTerraformRepresentation(),
					CustomEventSpecificationThresholdRuleFieldMetricPattern: []interface{}{
						map[string]interface{}{
							CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix:      prefix,
							CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix:     postfix,
							CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder: placeholder,
							CustomEventSpecificationThresholdRuleFieldMetricPatternOperator:    operator,
						},
					},
					CustomEventSpecificationThresholdRuleFieldRollup:            customEventSpecificationWithThresholdRuleRollup,
					CustomEventSpecificationThresholdRuleFieldWindow:            customEventSpecificationWithThresholdRuleWindow,
					CustomEventSpecificationThresholdRuleFieldAggregation:       customEventSpecificationWithThresholdRuleAggregation,
					CustomEventSpecificationThresholdRuleFieldConditionOperator: restapi.ConditionOperatorEquals.InstanaAPIValue(),
					CustomEventSpecificationThresholdRuleFieldConditionValue:    customEventSpecificationWithThresholdRuleConditionValue,
				},
			},
		},
	})

	result, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(*restapi.CustomEventSpecification)
	require.Equal(t, customEventSpecificationWithThresholdRuleID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, resourceName, customEventSpec.Name)
	require.Equal(t, customEventSpecificationWithThresholdRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEventSpecificationWithThresholdRuleQuery, *customEventSpec.Query)
	require.Equal(t, customEventSpecificationWithThresholdRuleDescription, *customEventSpec.Description)
	require.Equal(t, customEventSpecificationWithThresholdRuleExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, restapi.ThresholdRuleType, customEventSpec.Rules[0].DType)
	require.Nil(t, customEventSpec.Rules[0].MetricName)
	require.NotNil(t, customEventSpec.Rules[0].MetricPattern)
	require.Equal(t, prefix, customEventSpec.Rules[0].MetricPattern.Prefix)
	require.Equal(t, postfix, *customEventSpec.Rules[0].MetricPattern.Postfix)
	require.Equal(t, placeholder, *customEventSpec.Rules[0].MetricPattern.Placeholder)
	require.Equal(t, operator, customEventSpec.Rules[0].MetricPattern.Operator)
	require.Equal(t, customEventSpecificationWithThresholdRuleWindow, *customEventSpec.Rules[0].Window)
	require.Equal(t, customEventSpecificationWithThresholdRuleRollup, *customEventSpec.Rules[0].Rollup)
	require.Equal(t, customEventSpecificationWithThresholdRuleAggregation, *customEventSpec.Rules[0].Aggregation)
	require.Equal(t, restapi.ConditionOperatorEquals.InstanaAPIValue(), *customEventSpec.Rules[0].ConditionOperator)
	require.Equal(t, customEventSpecificationWithThresholdRuleConditionValue, *customEventSpec.Rules[0].ConditionValue)
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapStateOfThresholdRuleToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRule, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule:             []interface{}{},
			CustomEventSpecificationFieldThresholdRule: []interface{}{
				map[string]interface{}{
					CustomEventSpecificationRuleFieldSeverity: "invalid",
				},
			},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Error(t, err)
	require.ErrorContains(t, err, "invalid is not a valid severity")
}

func (r *customerEventSpecificationUnitTest) shouldFailToMapStateWhenNoRuleIsProvided(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldRule, []interface{}{
		map[string]interface{}{
			CustomEventSpecificationFieldEntityVerificationRule: []interface{}{},
			CustomEventSpecificationFieldSystemRule:             []interface{}{},
			CustomEventSpecificationFieldThresholdRule:          []interface{}{},
		},
	})

	_, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Error(t, err)
	require.ErrorContains(t, err, "no supported rule defined")
}
