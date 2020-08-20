package instana_test

import (
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

var testApplicationConfigProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

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
  match_specification = "{{MATCH_SPECIFICATION}}"
	boundary_scope = "DEFAULT"
}
`

const serverResponseTemplate = `
{
	"id" : "{{id}}",
	"label" : "prefix label suffix",
	"scope" : "INCLUDE_ALL_DOWNSTREAM",
	"boundaryScope" : "DEFAULT",
	"matchSpecification" : {
		"type" : "BINARY_OP",
		"left" : {
			"type" : "BINARY_OP",
			"left" : {
				"type" : "LEAF",
				"key" : "entity.name",
				"operator" : "CONTAINS",
				"value" : "foo"
			},
			"conjunction" : "AND",
			"right" : {
				"type" : "LEAF",
				"key" : "entity.type",
				"operator" : "EQUALS",
				"value" : "mysql"
			}
		},
		"conjunction" : "OR",
		"right" : {
			"type" : "LEAF",
			"key" : "entity.type",
			"operator" : "EQUALS",
			"value" : "elasticsearch"
		}
	}
}
`

const applicationConfigApiPath = restapi.ApplicationConfigsResourcePath + "/{id}"
const testApplicationConfigDefinition = "instana_application_config.example"
const defaultMatchSpecification = "entity.name CONTAINS 'foo' AND entity.type EQUALS 'mysql' OR entity.type EQUALS 'elasticsearch'"

var defaultMatchSpecificationModel = restapi.NewBinaryOperator(
	restapi.NewBinaryOperator(
		restapi.NewComparisionExpression("entity.name", restapi.ContainsOperator, "foo"),
		restapi.LogicalAnd,
		restapi.NewComparisionExpression("entity.type", restapi.EqualsOperator, "mysql"),
	),
	restapi.LogicalOr,
	restapi.NewComparisionExpression("entity.type", restapi.EqualsOperator, "elasticsearch"))

const applicationConfigID = "application-config-id"

func TestCRUDOfApplicationConfigResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, applicationConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, applicationConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, applicationConfigApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(serverResponseTemplate, "{{id}}", vars["id"])
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

	resource.UnitTest(t, resource.TestCase{
		Providers: testApplicationConfigProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithLabel0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testApplicationConfigDefinition, "id"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldLabel, "label 0"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldFullLabel, "prefix label 0 suffix"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldScope, "INCLUDE_ALL_DOWNSTREAM"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldMatchSpecification, defaultMatchSpecification),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldBoundaryScope, "DEFAULT"),
				),
			},
			{
				Config: resourceDefinitionWithLabel1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testApplicationConfigDefinition, "id"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldLabel, "label 1"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldFullLabel, "prefix label 1 suffix"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldScope, "INCLUDE_ALL_DOWNSTREAM"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldMatchSpecification, defaultMatchSpecification),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldBoundaryScope, "DEFAULT"),
				),
			},
		},
	})
}

func TestApplicationConfigSchemaDefinitionIsValid(t *testing.T) {
	schema := NewApplicationConfigResourceHandle().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ApplicationConfigFieldLabel)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(ApplicationConfigFieldFullLabel)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(ApplicationConfigFieldScope, ApplicationConfigScopeIncludeNoDownstream)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ApplicationConfigFieldMatchSpecification)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(ApplicationConfigFieldBoundaryScope, ApplicationConfigBoundaryScopeDefault)
}

func TestUserRoleResourceShouldHaveSchemaVersionOne(t *testing.T) {
	assert.Equal(t, 1, NewApplicationConfigResourceHandle().SchemaVersion)
}

func TestUserRoleResourceShouldHaveOneStateUpgraderForVersionZero(t *testing.T) {
	resourceHandler := NewApplicationConfigResourceHandle()

	assert.Equal(t, 1, len(resourceHandler.StateUpgraders))
	assert.Equal(t, 0, resourceHandler.StateUpgraders[0].Version)
}

func TestShouldMigrateApplicationConfigStateAndAddFullLabelWithSameValueAsLabelWhenMigratingFromVersion0To1(t *testing.T) {
	label := "Test Label"
	rawData := make(map[string]interface{})
	rawData[ApplicationConfigFieldLabel] = label
	meta := "dummy"

	result, err := NewApplicationConfigResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Equal(t, label, result[ApplicationConfigFieldFullLabel])
}

func TestShouldMigrateEmptyApplicationConfigStateFromVersion0To1(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"

	result, err := NewApplicationConfigResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Nil(t, result[ApplicationConfigFieldFullLabel])
}

func TestShouldReturnCorrectResourceNameForApplicationConfigResource(t *testing.T) {
	name := NewApplicationConfigResourceHandle().ResourceName

	assert.Equal(t, name, "instana_application_config")
}

func TestShouldUpdateApplicationConfigTerraformResourceStateFromModel(t *testing.T) {
	label := "label"
	applicationConfig := restapi.ApplicationConfig{
		ID:                 applicationConfigID,
		Label:              label,
		MatchSpecification: defaultMatchSpecificationModel,
		Scope:              ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope:      ApplicationConfigBoundaryScopeDefault,
	}

	testHelper := NewTestHelper(t)
	sut := NewApplicationConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, applicationConfig)

	assert.Nil(t, err)
	assert.Equal(t, applicationConfigID, resourceData.Id())
	assert.Equal(t, label, resourceData.Get(ApplicationConfigFieldFullLabel))
	assert.Equal(t, defaultMatchSpecification, resourceData.Get(ApplicationConfigFieldMatchSpecification))
	assert.Equal(t, ApplicationConfigScopeIncludeNoDownstream, resourceData.Get(ApplicationConfigFieldScope))
	assert.Equal(t, ApplicationConfigBoundaryScopeDefault, resourceData.Get(ApplicationConfigFieldBoundaryScope))
}

func TestShouldFailToUpdateApplicationConfigTerraformResourceStateFromModelWhenMatchSpecificationIsNotalid(t *testing.T) {
	comparision := restapi.NewComparisionExpression("entity.name", "INVALID", "foo", "INVALID")
	label := "label"
	applicationConfig := restapi.ApplicationConfig{
		ID:                 applicationConfigID,
		Label:              label,
		MatchSpecification: comparision,
		Scope:              ApplicationConfigScopeIncludeNoDownstream,
		BoundaryScope:      ApplicationConfigBoundaryScopeDefault,
	}

	testHelper := NewTestHelper(t)
	sut := NewApplicationConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, applicationConfig)

	assert.NotNil(t, err)
}

func TestShouldSuccessfullyConvertApplicationConfigStateToDataModel(t *testing.T) {
	label := "label"
	testHelper := NewTestHelper(t)
	resourceHandle := NewApplicationConfigResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(applicationConfigID)
	resourceData.Set(ApplicationConfigFieldFullLabel, label)
	resourceData.Set(ApplicationConfigFieldMatchSpecification, defaultMatchSpecification)
	resourceData.Set(ApplicationConfigFieldScope, ApplicationConfigScopeIncludeNoDownstream)
	resourceData.Set(ApplicationConfigFieldBoundaryScope, ApplicationConfigBoundaryScopeDefault)

	result, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.ApplicationConfig{}, result)
	assert.Equal(t, applicationConfigID, result.GetID())
	assert.Equal(t, label, result.(restapi.ApplicationConfig).Label)
	assert.Equal(t, defaultMatchSpecificationModel, result.(restapi.ApplicationConfig).MatchSpecification)
	assert.Equal(t, ApplicationConfigScopeIncludeNoDownstream, result.(restapi.ApplicationConfig).Scope)
	assert.Equal(t, ApplicationConfigBoundaryScopeDefault, result.(restapi.ApplicationConfig).BoundaryScope)
}

func TestShouldFailToConvertApplicationConfigStateToDataModelWhenMatchSpecificationIsNotValid(t *testing.T) {
	label := "label"
	testHelper := NewTestHelper(t)
	resourceHandle := NewApplicationConfigResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(applicationConfigID)
	resourceData.Set(ApplicationConfigFieldFullLabel, label)
	resourceData.Set(ApplicationConfigFieldMatchSpecification, "INVALID")
	resourceData.Set(ApplicationConfigFieldScope, ApplicationConfigScopeIncludeNoDownstream)
	resourceData.Set(ApplicationConfigFieldBoundaryScope, ApplicationConfigBoundaryScopeDefault)

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.NotNil(t, err)
}
