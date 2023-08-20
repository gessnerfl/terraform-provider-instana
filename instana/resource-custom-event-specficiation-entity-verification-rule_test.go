package instana_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const resourceCustomEventSpecificationWithEntityVerificationRuleDefinitionTemplate = `
resource "instana_custom_event_spec_entity_verification_rule" "example" {
  name = "name %d"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
  rule_severity = "warning"
  rule_matching_entity_type = "matching-entity-type"
  rule_matching_operator = "starts_with"
  rule_matching_entity_label = "matching-entity-label"
  rule_offline_duration = 60000
}
`

const (
	testCustomEventSpecificationWithEntityVerificationRuleDefinition = ResourceInstanaCustomEventSpecificationEntityVerificationRule + ".example"

	customEntityVerificationEventID                      = "custom-entity-verification-event-id"
	customEntityVerificationEventName                    = resourceName
	customEntityVerificationEventQuery                   = "query"
	customEntityVerificationEventExpirationTime          = 60000
	customEntityVerificationEventDescription             = "description"
	customEntityVerificationEventRuleMatchingEntityLabel = "matching-entity-label"
	customEntityVerificationEventRuleMatchingEntityType  = "matching-entity-type"
	customEntityVerificationEventRuleOfflineDuration     = 60000

	suffixString = " suffix"
	prefixString = "prefix "
)

var customEntityVerificationEventRuleMatchingOperator = restapi.MatchingOperatorStartsWith
var customEntityVerificationEventRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCreateResourceCustomEventSpecificationWithEntityVerificationRuleResourceWithMockServer(t *testing.T) {
	responseTemplate := `
	{
		"id" : "%s",
		"name" : "name %d",
		"query" : "query",
		"entityType" : "host",
		"enabled" : true,
		"triggering" : true,
		"description" : "description",
		"expirationTime" : 60000,
		"rules" : [ { "ruleType" : "entity_verification", "severity" : 5, "matchingEntityLabel" : "matching-entity-label", "matchingEntityType" : "matching-entity-type", "matchingOperator" : "startsWith", "offlineDuration" : 60000 } ]
	}
	`
	httpServer := createMockHttpServerForResource(restapi.CustomEventSpecificationResourcePath, responseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createCustomEventSpecificationWithEntityVerificationRuleResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testCustomEventSpecificationWithEntityVerificationRuleDefinition),
			createCustomEventSpecificationWithEntityVerificationRuleResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testCustomEventSpecificationWithEntityVerificationRuleDefinition),
		},
	})
}

func createCustomEventSpecificationWithEntityVerificationRuleResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceCustomEventSpecificationWithEntityVerificationRuleDefinitionTemplate, iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithEntityVerificationRuleDefinition, "id"),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldQuery, customEntityVerificationEventQuery),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldTriggering, trueAsString),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldDescription, customEntityVerificationEventDescription),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEntityVerificationEventExpirationTime)),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldEnabled, trueAsString),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingEntityLabel, customEntityVerificationEventRuleMatchingEntityLabel),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingEntityType, customEntityVerificationEventRuleMatchingEntityType),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingOperator, customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue()),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldOfflineDuration, strconv.Itoa(customEntityVerificationEventRuleOfflineDuration)),
		),
	}
}

func TestCustomEventSpecificationWithEntityVerificationRuleSchemaDefinitionIsValid(t *testing.T) {
	schema := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleSeverity)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingEntityLabel)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingEntityType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeInt(EntityVerificationRuleFieldOfflineDuration)
}

func TestCustomEventSpecificationWithEntityVerificationRuleResourceShouldHaveSchemaVersionFour(t *testing.T) {
	require.Equal(t, 4, NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().MetaData().SchemaVersion)
}

func TestCustomEventSpecificationWithEntityVerificationRuleShouldHaveFourStateUpgraderForVersionZeroToThree(t *testing.T) {
	resourceHandler := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	require.Equal(t, 4, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
	require.Equal(t, 1, resourceHandler.StateUpgraders()[1].Version)
	require.Equal(t, 2, resourceHandler.StateUpgraders()[2].Version)
	require.Equal(t, 3, resourceHandler.StateUpgraders()[3].Version)
}

func TestShouldMigrateCustomEventSpecificationWithEntityVerificationRuleStateAndAddFullNameWithSameValueAsNameWhenMigratingFromVersion0To1(t *testing.T) {
	name := "Test Name"
	rawData := make(map[string]interface{})
	rawData[CustomEventSpecificationFieldName] = name
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, name, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateEmptyCustomEventSpecificationWithEntityVerificationRuleStateFromVersion0To1(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[1].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateCustomEventSpecificationWithEntityVerificationRuleStateToVersion2WhenDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData["downstream_integration_ids"] = []interface{}{"id1", "id2"}
	rawData["downstream_broadcast_to_all_alerting_configs"] = true
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[1].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result["downstream_integration_ids"])
	require.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldMigrateCustomEventSpecificationWithEntityVerificationRuleStateToVersion2WhenNoDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result["downstream_integration_ids"])
	require.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldMigrateCustomEventSpecificationWithEntityVerificationRuleStateToVersion3WhenMatchingOperatorIsDefinedAndValid(t *testing.T) {
	for _, mo := range restapi.SupportedMatchingOperators {
		for _, v := range mo.TerraformSupportedValues() {
			t.Run(fmt.Sprintf("TestShouldMigrateCustomEventSpecificationWithEntityVerificationRuleStateToVersion3WhenMatchingOperatorIsDefinedAndValid%s", v), createTestCaseForSuccessfulMigrationOfCustomEventSpecificationWithEntityVerificationRuleToVersion3(mo, v))
		}
	}
}

func createTestCaseForSuccessfulMigrationOfCustomEventSpecificationWithEntityVerificationRuleToVersion3(mo restapi.MatchingOperator, value string) func(*testing.T) {
	return func(t *testing.T) {
		rawData := make(map[string]interface{})
		rawData[EntityVerificationRuleFieldMatchingOperator] = value
		meta := "dummy"
		ctx := context.Background()

		result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

		require.Nil(t, err)
		require.Equal(t, mo.InstanaAPIValue(), result[EntityVerificationRuleFieldMatchingOperator])
	}
}

func TestShouldDoNothingWhenMigratingCustomEventSpecificationWithEntityVerificationRuleToVersion3AndNoMatchingOperatorIsDefined(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result[EntityVerificationRuleFieldMatchingOperator])
}

func TestShouldReturnErrorWhenCustomEventSpecificationWithEntityVerificationRuleCannotBeMigratedToVersion3BecauseOfUnsupportedMatchingOperatorInState(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[EntityVerificationRuleFieldMatchingOperator] = "invalid"
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "not a supported matching operator")
	require.Equal(t, rawData, result)
}

func TestCustomEventSpecificationWithEntityVerificationRuleShouldMigrateFullnameToNameWhenExecutingForthStateUpgraderAndFullnameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[3].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, CustomEventSpecificationFieldFullName)
	require.Contains(t, result, CustomEventSpecificationFieldName)
	require.Equal(t, "test", result[CustomEventSpecificationFieldName])
}

func TestCustomEventSpecificationWithEntityVerificationRuleShouldDoNothingWhenExecutingForthStateUpgraderAndFullnameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[3].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForCustomEventSpecificationWithEntityVerificationRuleResource(t *testing.T) {
	name := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_custom_event_spec_entity_verification_rule")
}

func TestShouldUpdateCustomEventSpecificationWithEntityVerificationRuleTerraformStateFromApiObject(t *testing.T) {
	description := customEntityVerificationEventDescription
	expirationTime := customEntityVerificationEventExpirationTime
	query := customEntityVerificationEventQuery
	spec := &restapi.CustomEventSpecification{
		ID:             customEntityVerificationEventID,
		Name:           customEntityVerificationEventName,
		EntityType:     EntityVerificationRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			restapi.NewEntityVerificationRuleSpecification(customEntityVerificationEventRuleMatchingEntityLabel,
				customEntityVerificationEventRuleMatchingEntityType,
				customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue(),
				customEntityVerificationEventRuleOfflineDuration,
				restapi.SeverityWarning.GetAPIRepresentation()),
		},
	}

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	sut := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, customEntityVerificationEventID, resourceData.Id())
	require.Equal(t, customEntityVerificationEventName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, EntityVerificationRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEntityVerificationEventQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.Equal(t, customEntityVerificationEventRuleMatchingEntityLabel, resourceData.Get(EntityVerificationRuleFieldMatchingEntityLabel))
	require.Equal(t, customEntityVerificationEventRuleMatchingEntityType, resourceData.Get(EntityVerificationRuleFieldMatchingEntityType))
	require.Equal(t, customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue(), resourceData.Get(EntityVerificationRuleFieldMatchingOperator))
	require.Equal(t, customEntityVerificationEventRuleOfflineDuration, resourceData.Get(EntityVerificationRuleFieldOfflineDuration))
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), resourceData.Get(CustomEventSpecificationRuleSeverity))
}

func TestShouldFailToUpdateTerraformStateForCustomEventSpecificationWithEntityVerificationRuleWhenMatchingOperatorTypeIsNotSupported(t *testing.T) {
	description := customEntityVerificationEventDescription
	expirationTime := customEntityVerificationEventExpirationTime
	query := customEntityVerificationEventQuery
	spec := &restapi.CustomEventSpecification{
		ID:             customEntityVerificationEventID,
		Name:           customEntityVerificationEventName,
		EntityType:     EntityVerificationRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			restapi.NewEntityVerificationRuleSpecification(customEntityVerificationEventRuleMatchingEntityLabel,
				customEntityVerificationEventRuleMatchingEntityType,
				"invalid",
				customEntityVerificationEventRuleOfflineDuration,
				restapi.SeverityWarning.GetAPIRepresentation()),
		},
	}

	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	sut := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid is not a supported matching operator")
}

func TestShouldSuccessfullyConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEntityVerificationEventID)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldName, customEntityVerificationEventName)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldQuery, customEntityVerificationEventQuery)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldTriggering, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldDescription, customEntityVerificationEventDescription)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldExpirationTime, customEntityVerificationEventExpirationTime)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationFieldEnabled, true)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity)
	setValueOnResourceData(t, resourceData, EntityVerificationRuleFieldMatchingEntityLabel, customEntityVerificationEventRuleMatchingEntityLabel)
	setValueOnResourceData(t, resourceData, EntityVerificationRuleFieldMatchingEntityType, customEntityVerificationEventRuleMatchingEntityType)
	setValueOnResourceData(t, resourceData, EntityVerificationRuleFieldMatchingOperator, customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue())
	setValueOnResourceData(t, resourceData, EntityVerificationRuleFieldOfflineDuration, customEntityVerificationEventRuleOfflineDuration)

	result, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("foo", "bar"))

	require.Nil(t, err)
	require.IsType(t, &restapi.CustomEventSpecification{}, result)
	require.Equal(t, customEntityVerificationEventID, result.GetIDForResourcePath())
	require.Equal(t, customEntityVerificationEventName, result.Name)
	require.Equal(t, EntityVerificationRuleEntityType, result.EntityType)
	require.Equal(t, customEntityVerificationEventQuery, *result.Query)
	require.Equal(t, customEntityVerificationEventDescription, *result.Description)
	require.Equal(t, customEntityVerificationEventExpirationTime, *result.ExpirationTime)
	require.True(t, result.Triggering)
	require.True(t, result.Enabled)

	require.Equal(t, 1, len(result.Rules))
	require.Equal(t, customEntityVerificationEventRuleMatchingEntityLabel, *result.Rules[0].MatchingEntityLabel)
	require.Equal(t, customEntityVerificationEventRuleMatchingEntityType, *result.Rules[0].MatchingEntityType)
	require.Equal(t, customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue(), *result.Rules[0].MatchingOperator)
	require.Equal(t, customEntityVerificationEventRuleOfflineDuration, *result.Rules[0].OfflineDuration)
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), result.Rules[0].Severity)
}

func TestShouldFailToConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationRuleSeverity, "INVALID")

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter(prefixString, suffixString))

	require.NotNil(t, err)
}

func TestShouldFailToConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModelWhenMatchingOperatorIsMissing(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity)

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter(prefixString, suffixString))

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "is not a supported matching operator")
}

func TestShouldFailToConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModelWhenMatchingOperatorIsNotValid(t *testing.T) {
	testHelper := NewTestHelper[*restapi.CustomEventSpecification](t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	setValueOnResourceData(t, resourceData, CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity)
	setValueOnResourceData(t, resourceData, EntityVerificationRuleFieldMatchingOperator, "invalid")

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter(prefixString, suffixString))

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid is not a supported matching operator")
}
