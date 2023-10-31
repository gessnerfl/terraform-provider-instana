package instana_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

func TestAlertingConfig(t *testing.T) {
	unitTest := &alertingConfigResourceUnitTest{}
	t.Run("CRUD integration test with rule IDs", newAlertingConfigResourceIntegrationTestWithRuleIds())
	t.Run("CRUD integration test with event types", newAlertingConfigResourceIntegrationTestWithEventTypes())
	t.Run("schema should be valid", unitTest.resourceSchemaShouldBeValid)
	t.Run("should return correct schema name", unitTest.shouldReturnCorrectResourceNameForAlertingConfig)
	t.Run("should have schema version two", unitTest.shouldHaveSchemaVersionTwo)
	t.Run("should have two state upgrader for version zero and one", unitTest.shouldHaveTwoStateUpgraderForVersionZeroAndOne)
	t.Run("should migrate full alert name to alert name when executing second state upgrader and full alert name is available", unitTest.shouldMigrateFullAlertNameToAlertNameWhenExecutingSecondStateUpgraderAndFullAlertNameIsAvailable)
	t.Run("should do nothing when executing second state upgrader and full alert name is not available", unitTest.shouldDoNothingWhenExecutingSecondStateUpgraderAndFullAlertNameIsNotAvailable)
	t.Run("should return state with rule ids unchanged when migrating from version0 to version1", unitTest.shouldReturnStateWithRuleIdsUnchangedWhenMigratingFromVersion0ToVersion1)
	t.Run("should return state with event types unchanged when migrating from version0 to version1", unitTest.shouldReturnStateWithEventTypesUnchangedWhenMigratingFromVersion0ToVersion1)
	t.Run("should update resource state with rule ids", unitTest.shouldUpdateResourceStateWithRuleIds)
	t.Run("should update resource state with event types", unitTest.shouldUpdateResourceStateWithEventTypes)
	t.Run("should convert state to data model with rule ids", unitTest.shouldConvertStateToDataModelWithRuleIds)
	t.Run("should convert state to data model with event types", unitTest.shouldConvertStateToDataModelWithEventTypes)
}

const alertingConfigResourceDefinition = "instana_alerting_config.example"

func newAlertingConfigResourceIntegrationTestWithRuleIds() func(*testing.T) {
	resourceTemplate := `
resource "instana_alerting_config" "example" {
	alert_name = "name %d"
	integration_ids = [ "integration_id1", "integration_id2" ]
	event_filter_query = "query"
	event_filter_rule_ids = [ "rule-1", "rule-2" ]

	custom_payload_field {
		key    = "static-key"
		value  = "static-value"
	}
}
`
	serverResponseTemplate := `
{
	"id" : "%s",
	"alertName" : "name %d",
	"integrationIds" : [ "integration_id2", "integration_id1" ],
	"eventFilteringConfiguration" : {
		"query" : "query",
		"ruleIds" : [ "rule-2", "rule-1" ]
	},
    "customPayloadFields": [
		{
			"type": "staticString",
			"key": "static-key",
			"value": "static-value"
      	}
	]
}
`
	useCaseSpecificChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, 0), "rule-1"),
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterRuleIDs, 1), "rule-2"),
	}

	instance := &alertingConfigResourceIntegrationTest{
		resourceTemplate:       resourceTemplate,
		serverResponseTemplate: serverResponseTemplate,
		useCaseSpecificChecks:  useCaseSpecificChecks,
	}
	return instance.testCRUD()
}

func newAlertingConfigResourceIntegrationTestWithEventTypes() func(*testing.T) {
	resourceTemplate := `
resource "instana_alerting_config" "example" {
	alert_name = "name %d"
	integration_ids = [ "integration_id1", "integration_id2" ]
	event_filter_query = "query"
	event_filter_event_types = [ "incident", "critical" ]

	custom_payload_field {
		key    = "static-key"
		value  = "static-value"
	}
}
`
	serverResponseTemplate := `
{
	"id" : "%s",
	"alertName" : "name %d",
	"integrationIds" : [ "integration_id2", "integration_id1" ],
	"eventFilteringConfiguration" : {
		"query" : "query",
		"eventTypes" : [ "critical", "incident" ]
	},
    "customPayloadFields": [
		{
			"type": "staticString",
			"key": "static-key",
			"value": "static-value"
      	}
	]
}
`
	useCaseSpecificChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, 1), string(restapi.IncidentAlertEventType)),
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, fmt.Sprintf("%s.%d", AlertingConfigFieldEventFilterEventTypes, 0), string(restapi.CriticalAlertEventType)),
	}

	instance := &alertingConfigResourceIntegrationTest{
		resourceTemplate:       resourceTemplate,
		serverResponseTemplate: serverResponseTemplate,
		useCaseSpecificChecks:  useCaseSpecificChecks,
	}
	return instance.testCRUD()
}

type alertingConfigResourceIntegrationTest struct {
	resourceTemplate       string
	serverResponseTemplate string
	useCaseSpecificChecks  []resource.TestCheckFunc
}

func (it *alertingConfigResourceIntegrationTest) testCRUD() func(t *testing.T) {
	return func(t *testing.T) {
		httpServer := createMockHttpServerForResource(restapi.AlertsResourcePath, it.serverResponseTemplate)
		httpServer.Start()
		defer httpServer.Close()

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: testProviderFactory,
			Steps: []resource.TestStep{
				it.createTestCheckFunction(httpServer.GetPort(), 0),
				testStepImport(alertingConfigResourceDefinition),
				it.createTestCheckFunction(httpServer.GetPort(), 1),
				testStepImport(alertingConfigResourceDefinition),
			},
		})
	}
}

func (it *alertingConfigResourceIntegrationTest) createTestCheckFunction(httpPort int, iteration int) resource.TestStep {
	integrationId1 := "integration_id1"
	integrationId2 := "integration_id2"
	customPayloadFieldStaticKey := fmt.Sprintf("%s.0.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldKey)
	customPayloadFieldStaticValue := fmt.Sprintf("%s.0.%s", DefaultCustomPayloadFieldsName, CustomPayloadFieldsFieldStaticStringValue)
	defaultChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, AlertingConfigFieldAlertName, fmt.Sprintf("name %d", iteration)),
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, fmt.Sprintf("%s.%d", AlertingConfigFieldIntegrationIds, 0), integrationId1),
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, fmt.Sprintf("%s.%d", AlertingConfigFieldIntegrationIds, 1), integrationId2),
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, AlertingConfigFieldEventFilterQuery, "query"),
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, customPayloadFieldStaticKey, "static-key"),
		resource.TestCheckResourceAttr(alertingConfigResourceDefinition, customPayloadFieldStaticValue, "static-value"),
	}
	checks := append(defaultChecks, it.useCaseSpecificChecks...)
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(it.resourceTemplate, iteration), httpPort),
		Check:  resource.ComposeTestCheckFunc(checks...),
	}
}

type alertingConfigResourceUnitTest struct{}

func (ut *alertingConfigResourceUnitTest) resourceSchemaShouldBeValid(t *testing.T) {
	resourceHandle := NewAlertingConfigResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingConfigFieldAlertName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingConfigFieldIntegrationIds)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(AlertingConfigFieldEventFilterQuery)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeSetOfStrings(AlertingConfigFieldEventFilterEventTypes)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeSetOfStrings(AlertingConfigFieldEventFilterRuleIDs)
}

func (ut *alertingConfigResourceUnitTest) shouldReturnCorrectResourceNameForAlertingConfig(t *testing.T) {
	name := NewAlertingConfigResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_alerting_config", name, "Expected resource name to be instana_alerting_config")
}

func (ut *alertingConfigResourceUnitTest) shouldHaveSchemaVersionTwo(t *testing.T) {
	require.Equal(t, 2, NewAlertingConfigResourceHandle().MetaData().SchemaVersion)
}

func (ut *alertingConfigResourceUnitTest) shouldHaveTwoStateUpgraderForVersionZeroAndOne(t *testing.T) {
	resourceHandler := NewAlertingConfigResourceHandle()

	require.Equal(t, 2, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
	require.Equal(t, 1, resourceHandler.StateUpgraders()[1].Version)
}

func (ut *alertingConfigResourceUnitTest) shouldMigrateFullAlertNameToAlertNameWhenExecutingSecondStateUpgraderAndFullAlertNameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_alert_name": "test",
	}
	result, err := NewAlertingConfigResourceHandle().StateUpgraders()[1].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, AlertingConfigFieldFullAlertName)
	require.Contains(t, result, AlertingConfigFieldAlertName)
	require.Equal(t, "test", result[AlertingConfigFieldAlertName])
}

func (ut *alertingConfigResourceUnitTest) shouldDoNothingWhenExecutingSecondStateUpgraderAndFullAlertNameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"alert_name": "test",
	}
	result, err := NewAlertingConfigResourceHandle().StateUpgraders()[1].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func (ut *alertingConfigResourceUnitTest) shouldReturnStateWithRuleIdsUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[AlertingConfigFieldAlertName] = resourceName
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

func (ut *alertingConfigResourceUnitTest) shouldReturnStateWithEventTypesUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[AlertingConfigFieldAlertName] = resourceName
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
	alertingConfigIntegrationId1 = "alerting-integration-id1"
	alertingConfigIntegrationId2 = "alerting-integration-id2"
	alertingConfigRuleId1        = "alerting-rule-id1"
	alertingConfigRuleId2        = "alerting-rule-id2"
	alertingConfigQuery          = "alerting-query"
)

func (ut *alertingConfigResourceUnitTest) shouldUpdateResourceStateWithRuleIds(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingConfiguration](t)
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
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{
			{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   "static-key-1",
				Value: restapi.StaticStringCustomPayloadFieldValue("static-value-1"),
			},
			{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   "static-key-2",
				Value: restapi.StaticStringCustomPayloadFieldValue("static-value-2"),
			},
		},
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)
	require.Equal(t, alertingConfigID, resourceData.Id())
	require.Equal(t, alertingConfigName, resourceData.Get(AlertingConfigFieldAlertName))
	require.Equal(t, alertingConfigQuery, resourceData.Get(AlertingConfigFieldEventFilterQuery))
	ut.requireIntegrationIdsSetInResource(t, resourceData)

	ruleIDs := resourceData.Get(AlertingConfigFieldEventFilterRuleIDs).(*schema.Set)
	ut.requireSetMatchesToValues(t, ruleIDs, alertingConfigRuleId1, alertingConfigRuleId2)
	ut.requireCustomPayloadFieldsSetInResource(t, resourceData, data.CustomerPayloadFields)
}

func (ut *alertingConfigResourceUnitTest) shouldUpdateResourceStateWithEventTypes(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingConfiguration](t)
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
		CustomerPayloadFields: []restapi.CustomPayloadField[any]{
			{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   "static-key-1",
				Value: restapi.StaticStringCustomPayloadFieldValue("static-value-1"),
			},
			{
				Type:  restapi.StaticStringCustomPayloadType,
				Key:   "static-key-2",
				Value: restapi.StaticStringCustomPayloadFieldValue("static-value-2"),
			},
		},
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)
	require.Equal(t, alertingConfigID, resourceData.Id())
	require.Equal(t, alertingConfigName, resourceData.Get(AlertingConfigFieldAlertName))
	require.Equal(t, alertingConfigQuery, resourceData.Get(AlertingConfigFieldEventFilterQuery))
	ut.requireIntegrationIdsSetInResource(t, resourceData)

	eventTypes := resourceData.Get(AlertingConfigFieldEventFilterEventTypes).(*schema.Set)
	ut.requireSetMatchesToValues(t, eventTypes, string(restapi.CriticalAlertEventType), string(restapi.IncidentAlertEventType))
	ut.requireCustomPayloadFieldsSetInResource(t, resourceData, data.CustomerPayloadFields)
}

func (ut *alertingConfigResourceUnitTest) requireIntegrationIdsSetInResource(t *testing.T, resourceData *schema.ResourceData) {
	integrationIDs := resourceData.Get(AlertingConfigFieldIntegrationIds).(*schema.Set)
	ut.requireSetMatchesToValues(t, integrationIDs, alertingConfigIntegrationId1, alertingConfigIntegrationId2)
}

func (ut *alertingConfigResourceUnitTest) requireSetMatchesToValues(t *testing.T, set *schema.Set, values ...string) {
	require.Equal(t, len(values), set.Len())
	for _, v := range values {
		require.Contains(t, set.List(), v)
	}
}

func (ut *alertingConfigResourceUnitTest) requireCustomPayloadFieldsSetInResource(t *testing.T, resourceData *schema.ResourceData, models []restapi.CustomPayloadField[any]) {
	fields := resourceData.Get(DefaultCustomPayloadFieldsName)

	require.NotNil(t, fields)
	require.IsType(t, &schema.Set{}, fields)
	fieldList := fields.(*schema.Set).List()
	require.Len(t, fieldList, len(models))

	for _, val := range fieldList {
		require.IsType(t, map[string]interface{}{}, val)
		field := val.(map[string]interface{})

		key := field[CustomPayloadFieldsFieldDynamicKey].(string)
		model := ut.getCustomPayloadFieldByKey(key, models)
		require.NotNil(t, model)

		require.NotNil(t, field[CustomPayloadFieldsFieldStaticStringValue])
		require.Equal(t, string(model.Value.(restapi.StaticStringCustomPayloadFieldValue)), field[CustomPayloadFieldsFieldStaticStringValue])
		require.Nil(t, field[CustomPayloadFieldsFieldDynamicValue])
	}
}

func (ut *alertingConfigResourceUnitTest) getCustomPayloadFieldByKey(key string, models []restapi.CustomPayloadField[any]) *restapi.CustomPayloadField[any] {
	for _, f := range models {
		if f.Key == key {
			return &f
		}
	}
	return nil
}

func (ut *alertingConfigResourceUnitTest) shouldConvertStateToDataModelWithRuleIds(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingConfiguration](t)
	resourceHandle := NewAlertingConfigResourceHandle()
	integrationIds := []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2}
	ruleIds := []string{alertingConfigRuleId1, alertingConfigRuleId2}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(alertingConfigID)
	setValueOnResourceData(t, resourceData, AlertingConfigFieldAlertName, alertingConfigName)
	setValueOnResourceData(t, resourceData, AlertingConfigFieldIntegrationIds, integrationIds)
	setValueOnResourceData(t, resourceData, AlertingConfigFieldEventFilterQuery, alertingConfigQuery)
	setValueOnResourceData(t, resourceData, AlertingConfigFieldEventFilterRuleIDs, ruleIds)
	setValueOnResourceData(t, resourceData, DefaultCustomPayloadFieldsName, []interface{}{
		map[string]interface{}{
			CustomPayloadFieldsFieldKey:               "static-key-1",
			CustomPayloadFieldsFieldStaticStringValue: "static-value-1",
		},
		map[string]interface{}{
			CustomPayloadFieldsFieldKey:               "static-key-2",
			CustomPayloadFieldsFieldStaticStringValue: "static-value-2",
		},
	})

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingConfiguration{}, model)
	require.Equal(t, alertingConfigID, model.GetIDForResourcePath())
	require.Equal(t, alertingConfigName, model.AlertName)

	ut.requireIntegrationIdsSetInModel(t, model)
	require.Equal(t, alertingConfigQuery, *model.EventFilteringConfiguration.Query)
	ut.requireSliceValuesMatchesToValues(t, model.EventFilteringConfiguration.RuleIDs, alertingConfigRuleId1, alertingConfigRuleId2)

	require.Equal(t, []restapi.CustomPayloadField[any]{
		{
			Type:  restapi.StaticStringCustomPayloadType,
			Key:   "static-key-2",
			Value: restapi.StaticStringCustomPayloadFieldValue("static-value-2"),
		},
		{
			Type:  restapi.StaticStringCustomPayloadType,
			Key:   "static-key-1",
			Value: restapi.StaticStringCustomPayloadFieldValue("static-value-1"),
		},
	}, model.CustomerPayloadFields)
}

func (ut *alertingConfigResourceUnitTest) shouldConvertStateToDataModelWithEventTypes(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingConfiguration](t)
	resourceHandle := NewAlertingConfigResourceHandle()
	integrationIds := []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(alertingConfigID)
	setValueOnResourceData(t, resourceData, AlertingConfigFieldAlertName, alertingConfigName)
	setValueOnResourceData(t, resourceData, AlertingConfigFieldIntegrationIds, integrationIds)
	setValueOnResourceData(t, resourceData, AlertingConfigFieldEventFilterQuery, alertingConfigQuery)
	setValueOnResourceData(t, resourceData, AlertingConfigFieldEventFilterEventTypes, []string{"incident", "critical"})
	setValueOnResourceData(t, resourceData, DefaultCustomPayloadFieldsName, []interface{}{
		map[string]interface{}{
			CustomPayloadFieldsFieldKey:               "static-key-1",
			CustomPayloadFieldsFieldStaticStringValue: "static-value-1",
		},
		map[string]interface{}{
			CustomPayloadFieldsFieldKey:               "static-key-2",
			CustomPayloadFieldsFieldStaticStringValue: "static-value-2",
		},
	})

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.AlertingConfiguration{}, model)
	require.Equal(t, alertingConfigID, model.GetIDForResourcePath())
	require.Equal(t, alertingConfigName, model.AlertName)

	ut.requireIntegrationIdsSetInModel(t, model)
	require.Equal(t, alertingConfigQuery, *model.EventFilteringConfiguration.Query)

	eventTypes := model.EventFilteringConfiguration.EventTypes
	require.Len(t, eventTypes, 2)
	require.Contains(t, eventTypes, restapi.CriticalAlertEventType)
	require.Contains(t, eventTypes, restapi.IncidentAlertEventType)

	require.Equal(t, []restapi.CustomPayloadField[any]{
		{
			Type:  restapi.StaticStringCustomPayloadType,
			Key:   "static-key-2",
			Value: restapi.StaticStringCustomPayloadFieldValue("static-value-2"),
		},
		{
			Type:  restapi.StaticStringCustomPayloadType,
			Key:   "static-key-1",
			Value: restapi.StaticStringCustomPayloadFieldValue("static-value-1"),
		},
	}, model.CustomerPayloadFields)
}

func (ut *alertingConfigResourceUnitTest) requireIntegrationIdsSetInModel(t *testing.T, model *restapi.AlertingConfiguration) {
	ut.requireSliceValuesMatchesToValues(t, model.IntegrationIDs, alertingConfigIntegrationId1, alertingConfigIntegrationId2)
}

func (ut *alertingConfigResourceUnitTest) requireSliceValuesMatchesToValues(t *testing.T, data []string, values ...string) {
	require.Equal(t, len(values), len(data))
	for _, v := range values {
		require.Contains(t, data, v)
	}
}
