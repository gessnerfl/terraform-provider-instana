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
	customEntityVerificationEventFullName                = resourceFullName
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
		"name" : "prefix name %d suffix",
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

func createCustomEventSpecificationWithEntityVerificationRuleResourceTestStep(httpPort int, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceCustomEventSpecificationWithEntityVerificationRuleDefinitionTemplate, iteration), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithEntityVerificationRuleDefinition, "id"),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldFullName, formatResourceFullName(iteration)),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldQuery, customEntityVerificationEventQuery),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldTriggering, trueAsString),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldDescription, customEntityVerificationEventDescription),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEntityVerificationEventExpirationTime)),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldEnabled, trueAsString),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingEntityLabel, customEntityVerificationEventRuleMatchingEntityLabel),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingEntityType, customEntityVerificationEventRuleMatchingEntityType),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingOperator, string(customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue())),
			resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldOfflineDuration, strconv.Itoa(customEntityVerificationEventRuleOfflineDuration)),
		),
	}
}

func TestCustomEventSpecificationWithEntityVerificationRuleSchemaDefinitionIsValid(t *testing.T) {
	schema := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldFullName)
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

func TestCustomEventSpecificationWithEntityVerificationRuleResourceShouldHaveSchemaVersionThree(t *testing.T) {
	require.Equal(t, 3, NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().MetaData().SchemaVersion)
}

func TestCustomEventSpecificationWithEntityVerificationRuleShouldHaveThreeStateUpgraderForVersionZeroToTwo(t *testing.T) {
	resourceHandler := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	require.Equal(t, 3, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
	require.Equal(t, 1, resourceHandler.StateUpgraders()[1].Version)
	require.Equal(t, 2, resourceHandler.StateUpgraders()[2].Version)
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

func TestShouldReturnErrorWhenCustomEventSpecificationWithEntityVerificationRuleCannotBeMigratedToVersion3BecuaseOfUnsupportedMatchingOperatorInState(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData[EntityVerificationRuleFieldMatchingOperator] = "invalid"
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "not a supported matching operator")
	require.Equal(t, rawData, result)
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
		Name:           customEntityVerificationEventFullName,
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

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, customEntityVerificationEventID, resourceData.Id())
	require.Equal(t, customEntityVerificationEventName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, customEntityVerificationEventFullName, resourceData.Get(CustomEventSpecificationFieldFullName))
	require.Equal(t, EntityVerificationRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customEntityVerificationEventQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.Equal(t, customEntityVerificationEventRuleMatchingEntityLabel, resourceData.Get(EntityVerificationRuleFieldMatchingEntityLabel))
	require.Equal(t, customEntityVerificationEventRuleMatchingEntityType, resourceData.Get(EntityVerificationRuleFieldMatchingEntityType))
	require.Equal(t, string(customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue()), resourceData.Get(EntityVerificationRuleFieldMatchingOperator))
	require.Equal(t, customEntityVerificationEventRuleOfflineDuration, resourceData.Get(EntityVerificationRuleFieldOfflineDuration))
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), resourceData.Get(CustomEventSpecificationRuleSeverity))
}

func TestShouldFailToUpdateTerraformStateForCustomEventSpecificationWithEntityVerificationRuleWhenMatchingOperatorTypeIsNotSupported(t *testing.T) {
	description := customEntityVerificationEventDescription
	expirationTime := customEntityVerificationEventExpirationTime
	query := customEntityVerificationEventQuery
	spec := &restapi.CustomEventSpecification{
		ID:             customEntityVerificationEventID,
		Name:           customEntityVerificationEventFullName,
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

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec, testHelper.ResourceFormatter())

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid is not a supported matching operator")
}

func TestShouldSuccessfullyConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEntityVerificationEventID)
	resourceData.Set(CustomEventSpecificationFieldFullName, customEntityVerificationEventName)
	resourceData.Set(CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType)
	resourceData.Set(CustomEventSpecificationFieldQuery, customEntityVerificationEventQuery)
	resourceData.Set(CustomEventSpecificationFieldTriggering, true)
	resourceData.Set(CustomEventSpecificationFieldDescription, customEntityVerificationEventDescription)
	resourceData.Set(CustomEventSpecificationFieldExpirationTime, customEntityVerificationEventExpirationTime)
	resourceData.Set(CustomEventSpecificationFieldEnabled, true)
	resourceData.Set(CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity)
	resourceData.Set(EntityVerificationRuleFieldMatchingEntityLabel, customEntityVerificationEventRuleMatchingEntityLabel)
	resourceData.Set(EntityVerificationRuleFieldMatchingEntityType, customEntityVerificationEventRuleMatchingEntityType)
	resourceData.Set(EntityVerificationRuleFieldMatchingOperator, string(customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue()))
	resourceData.Set(EntityVerificationRuleFieldOfflineDuration, customEntityVerificationEventRuleOfflineDuration)

	result, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter(prefixString, suffixString))

	require.Nil(t, err)
	require.IsType(t, &restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(*restapi.CustomEventSpecification)
	require.Equal(t, customEntityVerificationEventID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, customEntityVerificationEventName, customEventSpec.Name)
	require.Equal(t, EntityVerificationRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customEntityVerificationEventQuery, *customEventSpec.Query)
	require.Equal(t, customEntityVerificationEventDescription, *customEventSpec.Description)
	require.Equal(t, customEntityVerificationEventExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, customEntityVerificationEventRuleMatchingEntityLabel, *customEventSpec.Rules[0].MatchingEntityLabel)
	require.Equal(t, customEntityVerificationEventRuleMatchingEntityType, *customEventSpec.Rules[0].MatchingEntityType)
	require.Equal(t, customEntityVerificationEventRuleMatchingOperator.InstanaAPIValue(), *customEventSpec.Rules[0].MatchingOperator)
	require.Equal(t, customEntityVerificationEventRuleOfflineDuration, *customEventSpec.Rules[0].OfflineDuration)
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func TestShouldFailToConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.Set(CustomEventSpecificationRuleSeverity, "INVALID")

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter(prefixString, suffixString))

	require.NotNil(t, err)
}

func TestShouldFailToConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModelWhenMatchingOperatorIsMissing(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.Set(CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity)

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter(prefixString, suffixString))

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "is not a supported matching operator")
}

func TestShouldFailToConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModelWhenMatchingOperatorIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.Set(CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity)
	resourceData.Set(EntityVerificationRuleFieldMatchingOperator, "invalid")

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter(prefixString, suffixString))

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid is not a supported matching operator")
}
