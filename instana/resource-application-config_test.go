package instana_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceApplicationConfigWithMatchSpecificationDefinitionTemplate = `
resource "instana_application_config" "example" {
  label = "name %d"
  scope = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "ALL"
  match_specification = "%s"
}
`

// Important if a match specification is provided the corresponding tag filter is also available.
const serverResponseWithMatchSpecificationTemplate = `
{
	"id" : "%s",
	"label" : "name %d",
	"scope" : "INCLUDE_ALL_DOWNSTREAM",
	"boundaryScope" : "ALL",
	"matchSpecification" : {
		"type" : "BINARY_OP",
		"left" : {
			"type" : "BINARY_OP",
			"left" : {
				"type" : "LEAF",
				"key" : "entity.name",
				"entity" : "DESTINATION",
				"operator" : "CONTAINS",
				"value" : "foo"
			},
			"conjunction" : "AND",
			"right" : {
				"type" : "LEAF",
				"key" : "agent.tag.environment",
				"entity" : "DESTINATION",
				"operator" : "EQUALS",
				"value" : "dev-speedboot-local-gessnerfl"
			}
		},
		"conjunction" : "OR",
		"right" : {
			"type" : "LEAF",
			"key" : "call.http.status",
			"entity" : "NOT_APPLICABLE",
			"operator" : "EQUALS",
			"value" : "404"
		}
	},
	"tagFilterExpression" : {
		"type" : "EXPRESSION",
		"logicalOperator": "OR",
		"elements" : [
			{
				"type" : "EXPRESSION",
				"logicalOperator": "AND",
				"elements" : [
					{
						"type" : "TAG_FILTER",
						"name" : "entity.name",
						"entity" : "DESTINATION",
						"operator" : "CONTAINS",
						"stringValue" : "foo",
						"value" : "foo"
					},
					{
						"type" : "TAG_FILTER",
						"name" : "agent.tag",
						"entity" : "DESTINATION",
						"operator" : "EQUALS",
						"stringValue" : "environment=dev-speedboot-local-gessnerfl",
						"key": "environment",
						"value": "dev-speedboot-local-gessnerfl"
					}
				]
			},
			{
				"type" : "TAG_FILTER",
				"name" : "call.http.status",
				"entity" : "NOT_APPLICABLE",
				"operator" : "EQUALS",
				"numberValue" : 404,
				"value" : 404
			}
		]
	}
}
`

const resourceApplicationConfigWithTagFilterDefinitionTemplate = `
resource "instana_application_config" "example" {
  label = "name %d"
  scope = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "ALL"
  tag_filter = "%s"
}
`

// Important if a match specification is not provided only the tag filter will be available.
const serverResponseWithTagFilterTemplate = `
{
	"id" : "%s",
	"label" : "name %d",
	"scope" : "INCLUDE_ALL_DOWNSTREAM",
	"boundaryScope" : "ALL",
	"tagFilterExpression" : {
		"type" : "EXPRESSION",
		"logicalOperator": "OR",
		"elements" : [
			{
				"type" : "EXPRESSION",
				"logicalOperator": "AND",
				"elements" : [
					{
						"type" : "TAG_FILTER",
						"name" : "entity.name",
						"entity" : "DESTINATION",
						"operator" : "CONTAINS",
						"stringValue" : "foo",
						"value" : "foo"
					},
					{
						"type" : "TAG_FILTER",
						"name" : "agent.tag",
						"entity" : "DESTINATION",
						"operator" : "EQUALS",
						"stringValue" : "environment=dev-speedboot-local-gessnerfl",
						"key": "environment",
						"value": "dev-speedboot-local-gessnerfl"
					}
				]
			},
			{
				"type" : "TAG_FILTER",
				"name" : "call.http.status",
				"entity" : "NOT_APPLICABLE",
				"operator" : "EQUALS",
				"numberValue" : 404,
				"value" : 404
			}
		]
	}
}
`

const (
	testApplicationConfigDefinition     = "instana_application_config.example"
	defaultMatchSpecification           = "entity.name CONTAINS 'foo' AND agent.tag.environment EQUALS 'dev-speedboot-local-gessnerfl' OR call.http.status@na EQUALS '404'"
	defaultNormalizedMatchSpecification = "entity.name@dest CONTAINS 'foo' AND agent.tag.environment@dest EQUALS 'dev-speedboot-local-gessnerfl' OR call.http.status@na EQUALS '404'"
	validMatchSpecification             = "entity.type EQUALS 'foo'"
	invalidMatchSpecification           = "entity.type bla bla bla"
	defaultTagFilter                    = "entity.name CONTAINS 'foo' AND agent.tag:environment EQUALS 'dev-speedboot-local-gessnerfl' OR call.http.status@na EQUALS 404"
	defaultNormalizedTagFilter          = "((entity.name@dest CONTAINS 'foo' AND agent.tag:environment@dest EQUALS 'dev-speedboot-local-gessnerfl') OR call.http.status@na EQUALS 404)"
	validTagFilter                      = "entity.type EQUALS 'foo'"
	invalidTagFilter                    = "entity.type bla bla bla"
	defaultLabel                        = "label"
	entityName                          = "entity.name"
	expressionEntityTypeDestEqValue     = "entity.type@dest EQUALS 'foo'"
	expressionEntityTypeSrcEqValue      = "entity.type@src EQUALS 'foo'"
)

var defaultMatchSpecificationModel = restapi.NewBinaryOperator(
	restapi.NewBinaryOperator(
		restapi.NewComparisonExpression(entityName, restapi.MatcherExpressionEntityDestination, restapi.ContainsOperator, "foo"),
		restapi.LogicalAnd,
		restapi.NewComparisonExpression("agent.tag.environment", restapi.MatcherExpressionEntityDestination, restapi.EqualsOperator, "dev-speedboot-local-gessnerfl"),
	),
	restapi.LogicalOr,
	restapi.NewComparisonExpression("call.http.status", restapi.MatcherExpressionEntityNotApplicable, restapi.EqualsOperator, "404"))

var defaultTagFilterModel = restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{
	restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{
		restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, entityName, restapi.ContainsOperator, "foo"),
		restapi.NewTagTagFilter(restapi.TagFilterEntityDestination, "agent.tag", restapi.EqualsOperator, "environment", "dev-speedboot-local-gessnerfl"),
	}),
	restapi.NewNumberTagFilter(restapi.TagFilterEntityNotApplicable, "call.http.status", restapi.EqualsOperator, 404),
})

const applicationConfigID = "application-config-id"

func TestCRUDOfApplicationConfigWithMatchSpecificationResourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.ApplicationConfigsResourcePath, serverResponseWithMatchSpecificationTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createApplicationConfigWithMatchSpecificationResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testApplicationConfigDefinition),
			createApplicationConfigWithMatchSpecificationResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testApplicationConfigDefinition),
		},
	})
}

func createApplicationConfigWithMatchSpecificationResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceApplicationConfigWithMatchSpecificationDefinitionTemplate, iteration, defaultMatchSpecification), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testApplicationConfigDefinition, "id"),
			resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldLabel, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeAllDownstream)),
			resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll)),
			resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldMatchSpecification, defaultNormalizedMatchSpecification),
			resource.TestCheckNoResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldTagFilter),
		),
	}
}

func TestCRUDOfApplicationConfigWithTagFilterResourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.ApplicationConfigsResourcePath, serverResponseWithTagFilterTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createApplicationConfigWithTagFilterResourceTestStep(httpServer.GetPort(), 0),
			testStepImport(testApplicationConfigDefinition),
			createApplicationConfigWithTagFilterResourceTestStep(httpServer.GetPort(), 1),
			testStepImport(testApplicationConfigDefinition),
		},
	})
}

func createApplicationConfigWithTagFilterResourceTestStep(httpPort int64, iteration int) resource.TestStep {
	config := appendProviderConfig(fmt.Sprintf(resourceApplicationConfigWithTagFilterDefinitionTemplate, iteration, defaultTagFilter), httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(testApplicationConfigDefinition, "id"),
			resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldLabel, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeAllDownstream)),
			resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll)),
			resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldTagFilter, defaultNormalizedTagFilter),
			resource.TestCheckNoResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldMatchSpecification),
		),
	}
}

func TestApplicationConfigSchemaDefinitionIsValid(t *testing.T) {
	schema := NewApplicationConfigResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ApplicationConfigFieldLabel)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeDefault))
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ApplicationConfigFieldMatchSpecification)
	require.Equal(t, []string{ApplicationConfigFieldMatchSpecification, ApplicationConfigFieldTagFilter}, schema[ApplicationConfigFieldMatchSpecification].ExactlyOneOf)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ApplicationConfigFieldTagFilter)
	require.Equal(t, []string{ApplicationConfigFieldMatchSpecification, ApplicationConfigFieldTagFilter}, schema[ApplicationConfigFieldTagFilter].ExactlyOneOf)
}

func TestShouldReturnTrueWhenCheckingForSchemaDiffSuppressForMatchSpecificationOfApplicationConfigAndValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := expressionEntityTypeDestEqValue
	newValue := "entity.type  EQUALS    'foo'"

	require.True(t, schema[ApplicationConfigFieldMatchSpecification].DiffSuppressFunc(ApplicationConfigFieldMatchSpecification, oldValue, newValue, nil))
}

func TestShouldReturnFalseWhenCheckingForSchemaDiffSuppressForMatchSpecificationOfApplicationConfigAndValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := expressionEntityTypeSrcEqValue
	newValue := validMatchSpecification

	require.False(t, schema[ApplicationConfigFieldMatchSpecification].DiffSuppressFunc(ApplicationConfigFieldMatchSpecification, oldValue, newValue, nil))
}

func TestShouldReturnTrueWhenCheckingForSchemaDiffSuppressForMatchSpecificationOfApplicationConfigAndValueCannotBeNormalizedAndOldAndNewValueAreEqual(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	invalidValue := invalidMatchSpecification

	require.True(t, schema[ApplicationConfigFieldMatchSpecification].DiffSuppressFunc(ApplicationConfigFieldMatchSpecification, invalidValue, invalidValue, nil))
}

func TestShouldReturnFalseWhenCheckingForSchemaDiffSuppressForMatchSpecificationOfApplicationConfigAndValueCannotBeNormalizedAndOldAndNewValueAreNotEqual(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := invalidMatchSpecification
	newValue := "entity.type foo foo foo"

	require.False(t, schema[ApplicationConfigFieldMatchSpecification].DiffSuppressFunc(ApplicationConfigFieldMatchSpecification, oldValue, newValue, nil))
}

func TestShouldReturnNormalizedValueForMatchSpecificationOfApplicationConfigWhenStateFuncIsCalledAndValueCanBeNormalized(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	expectedValue := expressionEntityTypeDestEqValue
	newValue := validMatchSpecification

	require.Equal(t, expectedValue, schema[ApplicationConfigFieldMatchSpecification].StateFunc(newValue))
}

func TestShouldReturnProvidedValueForMatchSpecificationOfApplicationConfigWhenStateFuncIsCalledAndValueCannotBeNormalized(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := invalidMatchSpecification

	require.Equal(t, value, schema[ApplicationConfigFieldMatchSpecification].StateFunc(value))
}

func TestShouldReturnNoErrorsAndWarningsWhenValidationOfMatchSpecificationOfApplicationConfigIsCalledAndValueCanBeParsed(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := validMatchSpecification

	warns, errs := schema[ApplicationConfigFieldMatchSpecification].ValidateFunc(value, ApplicationConfigFieldMatchSpecification)
	require.Empty(t, warns)
	require.Empty(t, errs)
}

func TestShouldReturnOneErrorAndNoWarningsWhenValidationOfMatchSpecificationOfApplicationConfigIsCalledAndValueCannotBeParsed(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := invalidMatchSpecification

	warns, errs := schema[ApplicationConfigFieldMatchSpecification].ValidateFunc(value, ApplicationConfigFieldMatchSpecification)
	require.Empty(t, warns)
	require.Len(t, errs, 1)
}

func TestShouldReturnTrueWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationConfigAndValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := expressionEntityTypeDestEqValue
	newValue := "entity.type  EQUALS    'foo'"

	require.True(t, schema[ApplicationConfigFieldTagFilter].DiffSuppressFunc(ApplicationConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnFalseWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationConfigAndValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := expressionEntityTypeSrcEqValue
	newValue := validTagFilter

	require.False(t, schema[ApplicationConfigFieldTagFilter].DiffSuppressFunc(ApplicationConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnTrueWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationConfigAndValueCannotBeNormalizedAndOldAndNewValueAreEqual(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	invalidValue := invalidTagFilter

	require.True(t, schema[ApplicationConfigFieldTagFilter].DiffSuppressFunc(ApplicationConfigFieldTagFilter, invalidValue, invalidValue, nil))
}

func TestShouldReturnFalseWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationConfigAndValueCannotBeNormalizedAndOldAndNewValueAreNotEqual(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := invalidTagFilter
	newValue := "entity.type foo foo foo"

	require.False(t, schema[ApplicationConfigFieldTagFilter].DiffSuppressFunc(ApplicationConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnNormalizedValueForTagFilterOfApplicationConfigWhenStateFuncIsCalledAndValueCanBeNormalized(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	expectedValue := expressionEntityTypeDestEqValue
	newValue := validTagFilter

	require.Equal(t, expectedValue, schema[ApplicationConfigFieldTagFilter].StateFunc(newValue))
}

func TestShouldReturnProvidedValueForTagFilterOfApplicationConfigWhenStateFuncIsCalledAndValueCannotBeNormalized(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := invalidTagFilter

	require.Equal(t, value, schema[ApplicationConfigFieldTagFilter].StateFunc(value))
}

func TestShouldReturnNoErrorsAndWarningsWhenValidationOfTagFilterOfApplicationConfigIsCalledAndValueCanBeParsed(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := validTagFilter

	warns, errs := schema[ApplicationConfigFieldTagFilter].ValidateFunc(value, ApplicationConfigFieldTagFilter)
	require.Empty(t, warns)
	require.Empty(t, errs)
}

func TestShouldReturnOneErrorAndNoWarningsWhenValidationOfTagFilterOfApplicationConfigIsCalledAndValueCannotBeParsed(t *testing.T) {
	resourceHandle := NewApplicationConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := invalidTagFilter

	warns, errs := schema[ApplicationConfigFieldTagFilter].ValidateFunc(value, ApplicationConfigFieldTagFilter)
	require.Empty(t, warns)
	require.Len(t, errs, 1)
}

func TestApplicationConfigResourceShouldHaveSchemaVersionFour(t *testing.T) {
	require.Equal(t, 4, NewApplicationConfigResourceHandle().MetaData().SchemaVersion)
}

func TestApplicationConfigResourceShouldHaveFourStateUpgraderForVersionZeroToThree(t *testing.T) {
	resourceHandler := NewApplicationConfigResourceHandle()

	require.Equal(t, 4, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
	require.Equal(t, 1, resourceHandler.StateUpgraders()[1].Version)
	require.Equal(t, 2, resourceHandler.StateUpgraders()[2].Version)
	require.Equal(t, 3, resourceHandler.StateUpgraders()[3].Version)
}

func TestShouldMigrateApplicationConfigStateAndAddFullLabelWithSameValueAsLabelWhenMigratingFromVersion0To1(t *testing.T) {
	label := "Test Label"
	rawData := make(map[string]interface{})
	rawData[ApplicationConfigFieldLabel] = label
	meta := "dummy"
	ctx := context.Background()

	result, err := NewApplicationConfigResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, label, result[ApplicationConfigFieldFullLabel])
}

func TestShouldMigrateEmptyApplicationConfigStateFromVersion0To1(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewApplicationConfigResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result[ApplicationConfigFieldFullLabel])
}

func TestShouldHarmonizeMatchSpecificationWhenMigratingStateFromVersion1To2(t *testing.T) {
	input := "entity.name EQUALS 'foo'"
	expectedResult := "entity.name@dest EQUALS 'foo'"
	rawData := make(map[string]interface{})
	rawData[ApplicationConfigFieldMatchSpecification] = input
	meta := "dummy"
	ctx := context.Background()

	result, err := NewApplicationConfigResourceHandle().StateUpgraders()[1].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, input, result[ApplicationConfigFieldMatchSpecification])
	require.Equal(t, expectedResult, result[ApplicationConfigFieldNormalizedMatchSpecification])
}

func TestShouldFailToHarmonizeMatchSpecificationWhenMigratingStateFromVersion1To2(t *testing.T) {
	input := "invalid match spec"
	rawData := make(map[string]interface{})
	rawData[ApplicationConfigFieldMatchSpecification] = input
	meta := "dummy"
	ctx := context.Background()

	_, err := NewApplicationConfigResourceHandle().StateUpgraders()[1].Upgrade(ctx, rawData, meta)

	require.Error(t, err)
}

func TestShouldMigrateEmptyApplicationConfigStateFromVersion1To2(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewApplicationConfigResourceHandle().StateUpgraders()[1].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result[ApplicationConfigFieldMatchSpecification])
	require.Nil(t, result[ApplicationConfigFieldNormalizedMatchSpecification])
}

func TestShouldRemoveHarmonizedMatchSpecificationWhenMigratingApplicationConfigStateFromVersion2To3AndHarmonizedMatchSpecificationIsSet(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	rawData["id"] = applicationConfigID
	rawData[ApplicationConfigFieldFullLabel] = defaultLabel
	rawData[ApplicationConfigFieldMatchSpecification] = defaultMatchSpecification
	rawData[ApplicationConfigFieldScope] = string(restapi.ApplicationConfigScopeIncludeNoDownstream)
	rawData[ApplicationConfigFieldBoundaryScope] = string(restapi.BoundaryScopeAll)
	expectedResult := copyMap(rawData)
	rawData[ApplicationConfigFieldNormalizedMatchSpecification] = defaultNormalizedMatchSpecification

	result, err := NewApplicationConfigResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestShouldRemoveHarmonizedMatchSpecificationWhenMigratingApplicationConfigStateFromVersion2To3AndNoHarmonizedMatchSpecificationIsSet(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	rawData["id"] = applicationConfigID
	rawData[ApplicationConfigFieldFullLabel] = defaultLabel
	rawData[ApplicationConfigFieldMatchSpecification] = defaultMatchSpecification
	rawData[ApplicationConfigFieldScope] = string(restapi.ApplicationConfigScopeIncludeNoDownstream)
	rawData[ApplicationConfigFieldBoundaryScope] = string(restapi.BoundaryScopeAll)
	expectedResult := copyMap(rawData)

	result, err := NewApplicationConfigResourceHandle().StateUpgraders()[2].Upgrade(ctx, rawData, meta)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestApplicationConfigShouldMigrateFullLabelToLabelWhenExecutingThirdStateUpgraderAndFullLabelIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_label": "test",
	}
	result, err := NewApplicationConfigResourceHandle().StateUpgraders()[3].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, ApplicationConfigFieldFullLabel)
	require.Contains(t, result, ApplicationConfigFieldLabel)
	require.Equal(t, "test", result[ApplicationConfigFieldLabel])
}

func TestApplicationConfigShouldDoNothingWhenExecutingThirdStateUpgraderAndFullLabelIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"label": "test",
	}
	result, err := NewApplicationConfigResourceHandle().StateUpgraders()[3].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForApplicationConfigResource(t *testing.T) {
	name := NewApplicationConfigResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_application_config")
}

func TestShouldUpdateApplicationConfigTerraformResourceStateFromModelWhenMatchSpecificationIsProvided(t *testing.T) {
	applicationConfig := restapi.ApplicationConfig{
		ID:                 applicationConfigID,
		Label:              defaultLabel,
		MatchSpecification: defaultMatchSpecificationModel,
		Scope:              restapi.ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope:      restapi.BoundaryScopeAll,
	}

	testHelper := NewTestHelper[*restapi.ApplicationConfig](t)
	sut := NewApplicationConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &applicationConfig)

	require.NoError(t, err)
	require.Equal(t, applicationConfigID, resourceData.Id())
	require.Equal(t, defaultLabel, resourceData.Get(ApplicationConfigFieldLabel))
	require.Equal(t, defaultNormalizedMatchSpecification, resourceData.Get(ApplicationConfigFieldMatchSpecification))
	_, tagFilterSet := resourceData.GetOk(ApplicationConfigFieldTagFilter)
	require.False(t, tagFilterSet)
	require.Equal(t, string(restapi.ApplicationConfigScopeIncludeNoDownstream), resourceData.Get(ApplicationConfigFieldScope))
	require.Equal(t, string(restapi.BoundaryScopeAll), resourceData.Get(ApplicationConfigFieldBoundaryScope))
}

func TestShouldFailToUpdateApplicationConfigTerraformResourceStateFromModelWhenMatchSpecificationIsNotValid(t *testing.T) {
	comparison := restapi.NewComparisonExpression(entityName, restapi.MatcherExpressionEntityDestination, "INVALID", "foo")
	applicationConfig := restapi.ApplicationConfig{
		ID:                 applicationConfigID,
		Label:              defaultLabel,
		MatchSpecification: comparison,
		Scope:              restapi.ApplicationConfigScopeIncludeNoDownstream,
	}

	testHelper := NewTestHelper[*restapi.ApplicationConfig](t)
	sut := NewApplicationConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &applicationConfig)

	require.Error(t, err)
}

func TestShouldUpdateApplicationConfigTerraformResourceStateFromModelWhenTagFilterIsProvided(t *testing.T) {
	applicationConfig := restapi.ApplicationConfig{
		ID:                  applicationConfigID,
		Label:               defaultLabel,
		TagFilterExpression: defaultTagFilterModel,
		Scope:               restapi.ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope:       restapi.BoundaryScopeAll,
	}

	testHelper := NewTestHelper[*restapi.ApplicationConfig](t)
	sut := NewApplicationConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &applicationConfig)

	require.NoError(t, err)
	require.Equal(t, applicationConfigID, resourceData.Id())
	require.Equal(t, defaultLabel, resourceData.Get(ApplicationConfigFieldLabel))
	_, matchSpecificationSet := resourceData.GetOk(ApplicationConfigFieldMatchSpecification)
	require.False(t, matchSpecificationSet)
	require.Equal(t, defaultNormalizedTagFilter, resourceData.Get(ApplicationConfigFieldTagFilter))
	require.Equal(t, string(restapi.ApplicationConfigScopeIncludeNoDownstream), resourceData.Get(ApplicationConfigFieldScope))
	require.Equal(t, string(restapi.BoundaryScopeAll), resourceData.Get(ApplicationConfigFieldBoundaryScope))
}

func TestShouldFailToUpdateApplicationConfigTerraformResourceStateFromModelWhenTagFilterIsNotValid(t *testing.T) {
	comparison := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, entityName, "INVALID", "foo")
	applicationConfig := restapi.ApplicationConfig{
		ID:                  applicationConfigID,
		Label:               defaultLabel,
		TagFilterExpression: comparison,
		Scope:               restapi.ApplicationConfigScopeIncludeNoDownstream,
	}

	testHelper := NewTestHelper[*restapi.ApplicationConfig](t)
	sut := NewApplicationConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &applicationConfig)

	require.Error(t, err)
}

func TestShouldSuccessfullyConvertApplicationConfigStateToDataModelWhenMatchSpecificationIsAvailable(t *testing.T) {
	testHelper := NewTestHelper[*restapi.ApplicationConfig](t)
	resourceHandle := NewApplicationConfigResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(applicationConfigID)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldLabel, defaultLabel)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldMatchSpecification, defaultMatchSpecification)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll))

	result, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.ApplicationConfig{}, result)
	require.Equal(t, applicationConfigID, result.GetIDForResourcePath())
	require.Equal(t, defaultLabel, result.Label)
	require.Equal(t, defaultMatchSpecificationModel, result.MatchSpecification)
	require.Nil(t, result.TagFilterExpression)
	require.Equal(t, restapi.ApplicationConfigScopeIncludeNoDownstream, result.Scope)
	require.Equal(t, restapi.BoundaryScopeAll, result.BoundaryScope)
}

func TestShouldFailToConvertApplicationConfigStateToDataModelWhenMatchSpecificationIsNotValid(t *testing.T) {
	testHelper := NewTestHelper[*restapi.ApplicationConfig](t)
	resourceHandle := NewApplicationConfigResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(applicationConfigID)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldLabel, defaultLabel)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldMatchSpecification, "INVALID")
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll))

	_, err := resourceHandle.MapStateToDataObject(resourceData)

	require.NotNil(t, err)
}

func TestShouldSuccessfullyConvertApplicationConfigStateToDataModelWhenTagFilterIsAvailable(t *testing.T) {
	testHelper := NewTestHelper[*restapi.ApplicationConfig](t)
	resourceHandle := NewApplicationConfigResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(applicationConfigID)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldLabel, defaultLabel)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldTagFilter, defaultTagFilter)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll))

	result, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.ApplicationConfig{}, result)
	require.Equal(t, applicationConfigID, result.GetIDForResourcePath())
	require.Equal(t, defaultLabel, result.Label)
	require.Nil(t, result.MatchSpecification)
	require.Equal(t, defaultTagFilterModel, result.TagFilterExpression)
	require.Equal(t, restapi.ApplicationConfigScopeIncludeNoDownstream, result.Scope)
	require.Equal(t, restapi.BoundaryScopeAll, result.BoundaryScope)
}

func TestShouldFailToConvertApplicationConfigStateToDataModelWhenTagFilterIsNotValid(t *testing.T) {
	testHelper := NewTestHelper[*restapi.ApplicationConfig](t)
	resourceHandle := NewApplicationConfigResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(applicationConfigID)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldLabel, defaultLabel)
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldTagFilter, "INVALID")
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	setValueOnResourceData(t, resourceData, ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll))

	_, err := resourceHandle.MapStateToDataObject(resourceData)

	require.NotNil(t, err)
}
