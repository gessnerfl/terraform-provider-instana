package instana_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceApplicationConfigDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_application_config" "example" {
  label = "label {{ITERATOR}}"
  scope = "INCLUDE_ALL_DOWNSTREAM"
  boundary_scope = "ALL"
  match_specification = "{{MATCH_SPECIFICATION}}"
}
`

const serverResponseTemplate = `
{
	"id" : "{{id}}",
	"label" : "prefix label {{ITERATOR}} suffix",
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
				"key" : "entity.type",
				"entity" : "DESTINATION",
				"operator" : "EQUALS",
				"value" : "mysql"
			}
		},
		"conjunction" : "OR",
		"right" : {
			"type" : "LEAF",
			"key" : "entity.type",
			"entity" : "SOURCE",
			"operator" : "EQUALS",
			"value" : "elasticsearch"
		}
	}
}
`

const applicationConfigApiPath = restapi.ApplicationConfigsResourcePath + "/{id}"
const testApplicationConfigDefinition = "instana_application_config.example"
const defaultMatchSpecification = "entity.name CONTAINS 'foo' AND entity.type EQUALS 'mysql' OR entity.type@src EQUALS 'elasticsearch'"
const defaultMatchSpecificationNormalized = "entity.name@dest CONTAINS 'foo' AND entity.type@dest EQUALS 'mysql' OR entity.type@src EQUALS 'elasticsearch'"

var defaultMatchSpecificationModel = restapi.NewBinaryOperator(
	restapi.NewBinaryOperator(
		restapi.NewComparisionExpression("entity.name", restapi.MatcherExpressionEntityDestination, restapi.ContainsOperator, "foo"),
		restapi.LogicalAnd,
		restapi.NewComparisionExpression("entity.type", restapi.MatcherExpressionEntityDestination, restapi.EqualsOperator, "mysql"),
	),
	restapi.LogicalOr,
	restapi.NewComparisionExpression("entity.type", restapi.MatcherExpressionEntitySource, restapi.EqualsOperator, "elasticsearch"))

const applicationConfigID = "application-config-id"

func TestCRUDOfApplicationConfigResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, applicationConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, applicationConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, applicationConfigApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := restapi.ApplicationConfigsResourcePath + "/" + vars["id"]
		json := strings.ReplaceAll(strings.ReplaceAll(serverResponseTemplate, "{{id}}", vars["id"]), "{{ITERATOR}}", strconv.Itoa(getZeroBasedCallCount(httpServer, http.MethodPut, path)))
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutLabel := strings.ReplaceAll(
		strings.ReplaceAll(resourceApplicationConfigDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort())),
		"{{MATCH_SPECIFICATION}}",
		defaultMatchSpecification,
	)

	resourceDefinitionWithLabel0 := strings.ReplaceAll(resourceDefinitionWithoutLabel, iteratorPlaceholder, "0")
	resourceDefinitionWithLabel1 := strings.ReplaceAll(resourceDefinitionWithoutLabel, iteratorPlaceholder, "1")

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithLabel0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testApplicationConfigDefinition, "id"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldLabel, "label 0"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldFullLabel, "prefix label 0 suffix"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeAllDownstream)),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll)),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldMatchSpecification, defaultMatchSpecification),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldNormalizedMatchSpecification, defaultMatchSpecificationNormalized),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceDefinitionWithLabel1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testApplicationConfigDefinition, "id"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldLabel, "label 1"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldFullLabel, "prefix label 1 suffix"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeAllDownstream)),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll)),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldMatchSpecification, defaultMatchSpecification),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldNormalizedMatchSpecification, defaultMatchSpecificationNormalized),
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

func TestApplicationConfigSchemaDefinitionIsValid(t *testing.T) {
	schema := NewApplicationConfigResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ApplicationConfigFieldLabel)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(ApplicationConfigFieldFullLabel)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeDefault))
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ApplicationConfigFieldMatchSpecification)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(ApplicationConfigFieldNormalizedMatchSpecification)
}

func TestApplicationConfigResourceShouldHaveSchemaVersionTwo(t *testing.T) {
	require.Equal(t, 2, NewApplicationConfigResourceHandle().MetaData().SchemaVersion)
}

func TestApplicationConfigResourceShouldHaveTwoStateUpgraderForVersionZeroToOne(t *testing.T) {
	resourceHandler := NewApplicationConfigResourceHandle()

	require.Equal(t, 2, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
	require.Equal(t, 1, resourceHandler.StateUpgraders()[1].Version)
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

func TestShouldReturnCorrectResourceNameForApplicationConfigResource(t *testing.T) {
	name := NewApplicationConfigResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_application_config")
}

func TestShouldUpdateApplicationConfigTerraformResourceStateFromModel(t *testing.T) {
	label := "label"
	fullLabel := "prefix label suffix"
	applicationConfig := restapi.ApplicationConfig{
		ID:                 applicationConfigID,
		Label:              fullLabel,
		MatchSpecification: defaultMatchSpecificationModel,
		Scope:              restapi.ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope:      restapi.BoundaryScopeAll,
	}

	testHelper := NewTestHelper(t)
	sut := NewApplicationConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &applicationConfig, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, applicationConfigID, resourceData.Id())
	require.Equal(t, label, resourceData.Get(ApplicationConfigFieldLabel))
	require.Equal(t, fullLabel, resourceData.Get(ApplicationConfigFieldFullLabel))
	require.Equal(t, defaultMatchSpecificationNormalized, resourceData.Get(ApplicationConfigFieldNormalizedMatchSpecification))
	require.Equal(t, string(restapi.ApplicationConfigScopeIncludeNoDownstream), resourceData.Get(ApplicationConfigFieldScope))
	require.Equal(t, string(restapi.BoundaryScopeAll), resourceData.Get(ApplicationConfigFieldBoundaryScope))
}

func TestShouldFailToUpdateApplicationConfigTerraformResourceStateFromModelWhenMatchSpecificationIsNotalid(t *testing.T) {
	comparision := restapi.NewComparisionExpression("entity.name", restapi.MatcherExpressionEntityDestination, "INVALID", "foo")
	label := "label"
	applicationConfig := restapi.ApplicationConfig{
		ID:                 applicationConfigID,
		Label:              label,
		MatchSpecification: comparision,
		Scope:              restapi.ApplicationConfigScopeIncludeNoDownstream,
	}

	testHelper := NewTestHelper(t)
	sut := NewApplicationConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &applicationConfig, testHelper.ResourceFormatter())

	require.NotNil(t, err)
}

func TestShouldSuccessfullyConvertApplicationConfigStateToDataModel(t *testing.T) {
	label := "label"
	testHelper := NewTestHelper(t)
	resourceHandle := NewApplicationConfigResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(applicationConfigID)
	resourceData.Set(ApplicationConfigFieldFullLabel, label)
	resourceData.Set(ApplicationConfigFieldMatchSpecification, defaultMatchSpecification)
	resourceData.Set(ApplicationConfigFieldNormalizedMatchSpecification, defaultMatchSpecificationNormalized)
	resourceData.Set(ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	resourceData.Set(ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll))

	result, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.ApplicationConfig{}, result)
	require.Equal(t, applicationConfigID, result.GetIDForResourcePath())
	require.Equal(t, label, result.(*restapi.ApplicationConfig).Label)
	require.Equal(t, defaultMatchSpecificationModel, result.(*restapi.ApplicationConfig).MatchSpecification)
	require.Equal(t, restapi.ApplicationConfigScopeIncludeNoDownstream, result.(*restapi.ApplicationConfig).Scope)
	require.Equal(t, restapi.BoundaryScopeAll, result.(*restapi.ApplicationConfig).BoundaryScope)
}

func TestShouldFailToConvertApplicationConfigStateToDataModelWhenMatchSpecificationIsNotValid(t *testing.T) {
	label := "label"
	testHelper := NewTestHelper(t)
	resourceHandle := NewApplicationConfigResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(applicationConfigID)
	resourceData.Set(ApplicationConfigFieldFullLabel, label)
	resourceData.Set(ApplicationConfigFieldMatchSpecification, "INVALID")
	resourceData.Set(ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	resourceData.Set(ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll))

	_, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.NotNil(t, err)
}
