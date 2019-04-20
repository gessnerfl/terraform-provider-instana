package instana_test

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

var testApplicationConfigProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceApplicationConfigDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_application_config" "example" {
  label = "label"
  scope = "INCLUDE_ALL_DOWNSTREAM"
  match_specification = "{{MATCH_SPECIFICATION}}"
}
`

const applicationConfigApiPath = restapi.ApplicationConfigsResourcePath + "/{id}"
const testApplicationConfigDefinition = "instana_application_config.example"
const defaultMatchSpecification = "entity.name CONTAINS 'foo' AND entity.type EQUALS 'mysql' OR entity.type EQUALS 'elasticsearch'"
const applicationConfigID = "application-config-id"

func TestCRUDOfApplicationConfigResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, applicationConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, applicationConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, applicationConfigApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"label" : "label",
			"scope" : "INCLUDE_ALL_DOWNSTREAM",
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
		`, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinition := strings.ReplaceAll(
		strings.ReplaceAll(resourceApplicationConfigDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort())),
		"{{MATCH_SPECIFICATION}}",
		defaultMatchSpecification,
	)

	resource.UnitTest(t, resource.TestCase{
		Providers: testApplicationConfigProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testApplicationConfigDefinition, "id"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldLabel, "label"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldScope, "INCLUDE_ALL_DOWNSTREAM"),
					resource.TestCheckResourceAttr(testApplicationConfigDefinition, ApplicationConfigFieldMatchSpecification, defaultMatchSpecification),
				),
			},
		},
	})
}

func TestResourceApplicationConfigDefinition(t *testing.T) {
	resource := CreateResourceApplicationConfig()

	validateApplicationConfigResourceSchema(resource.Schema, t)

	if resource.Create == nil {
		t.Fatal("Create function expected")
	}
	if resource.Update == nil {
		t.Fatal("Update function expected")
	}
	if resource.Read == nil {
		t.Fatal("Read function expected")
	}
	if resource.Delete == nil {
		t.Fatal("Delete function expected")
	}
}

func validateApplicationConfigResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ApplicationConfigFieldLabel)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeStringWithDefault(ApplicationConfigFieldScope, ApplicationConfigScopeIncludeNoDownstream)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ApplicationConfigFieldMatchSpecification)
}

func TestShouldSuccessfullyReadApplicationConfigFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := createBaseTestApplicationConfigModel()
	testShouldSuccessfullyReadApplicationConfigFromInstanaAPI(expectedModel, t)
}

func TestShouldSuccessfullyReadApplicationConfigFromInstanaAPIWhenBaseDataWithScopeIsReturned(t *testing.T) {
	expectedModel := createTestApplicationConfigModelWithRollup()
	testShouldSuccessfullyReadApplicationConfigFromInstanaAPI(expectedModel, t)
}

func testShouldSuccessfullyReadApplicationConfigFromInstanaAPI(expectedModel restapi.ApplicationConfig, t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyApplicationConfigResourceData()
	resourceData.SetId(applicationConfigID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApplicationConfigApi := mocks.NewMockApplicationConfigResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().ApplicationConfigs().Return(mockApplicationConfigApi).Times(1)
	mockApplicationConfigApi.EXPECT().GetOne(gomock.Eq(applicationConfigID)).Return(expectedModel, nil).Times(1)

	err := ReadApplicationConfig(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	verifyApplicationConfigModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldFailToReadApplicationConfigFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyApplicationConfigResourceData()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	err := ReadApplicationConfig(resourceData, mockInstanaAPI)

	if err == nil || !strings.HasPrefix(err.Error(), "ID of application config") {
		t.Fatal("Expected error to occur because of missing id")
	}
}

func TestShouldFailToReadApplicationConfigFromInstanaAPIAndDeleteResourceWhenRoleDoesNotExist(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyApplicationConfigResourceData()
	resourceData.SetId(applicationConfigID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApplicationConfigApi := mocks.NewMockApplicationConfigResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().ApplicationConfigs().Return(mockApplicationConfigApi).Times(1)
	mockApplicationConfigApi.EXPECT().GetOne(gomock.Eq(applicationConfigID)).Return(restapi.ApplicationConfig{}, restapi.ErrEntityNotFound).Times(1)

	err := ReadApplicationConfig(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldFailToReadApplicationConfigFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyApplicationConfigResourceData()
	resourceData.SetId(applicationConfigID)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApplicationConfigApi := mocks.NewMockApplicationConfigResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().ApplicationConfigs().Return(mockApplicationConfigApi).Times(1)
	mockApplicationConfigApi.EXPECT().GetOne(gomock.Eq(applicationConfigID)).Return(restapi.ApplicationConfig{}, expectedError).Times(1)

	err := ReadApplicationConfig(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error should be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID should still be set")
	}
}

func TestShouldCreateApplicationConfigThroughInstanaAPI(t *testing.T) {
	data := createFullTestApplicationConfigData()
	resourceData := NewTestHelper(t).CreateApplicationConfigResourceData(data)
	expectedModel := createTestApplicationConfigModelWithRollup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApplicationConfigApi := mocks.NewMockApplicationConfigResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().ApplicationConfigs().Return(mockApplicationConfigApi).Times(1)
	mockApplicationConfigApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.ApplicationConfig{})).Return(expectedModel, nil).Times(1)

	err := CreateApplicationConfig(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	verifyApplicationConfigModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldReturnErrorWhenCreateApplicationConfigFailsThroughInstanaAPI(t *testing.T) {
	data := createFullTestApplicationConfigData()
	resourceData := NewTestHelper(t).CreateApplicationConfigResourceData(data)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApplicationConfigApi := mocks.NewMockApplicationConfigResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().ApplicationConfigs().Return(mockApplicationConfigApi).Times(1)
	mockApplicationConfigApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.ApplicationConfig{})).Return(restapi.ApplicationConfig{}, expectedError).Times(1)

	err := CreateApplicationConfig(resourceData, mockInstanaAPI)

	if err == nil || expectedError != err {
		t.Fatal("Expected definned error to be returned")
	}
}

func TestShouldDeleteApplicationConfigThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestApplicationConfigData()
	resourceData := NewTestHelper(t).CreateApplicationConfigResourceData(data)
	resourceData.SetId(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApplicationConfigApi := mocks.NewMockApplicationConfigResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().ApplicationConfigs().Return(mockApplicationConfigApi).Times(1)
	mockApplicationConfigApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(nil).Times(1)

	err := DeleteApplicationConfig(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldReturnErrorWhenDeleteApplicationConfigFailsThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestApplicationConfigData()
	resourceData := NewTestHelper(t).CreateApplicationConfigResourceData(data)
	resourceData.SetId(id)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApplicationConfigApi := mocks.NewMockApplicationConfigResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().ApplicationConfigs().Return(mockApplicationConfigApi).Times(1)
	mockApplicationConfigApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(expectedError).Times(1)

	err := DeleteApplicationConfig(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error to be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
	}
}

func verifyApplicationConfigModelAppliedToResource(model restapi.ApplicationConfig, resourceData *schema.ResourceData, t *testing.T) {
	if model.ID != resourceData.Id() {
		t.Fatal("Expected ID to be identical")
	}
	if model.Label != resourceData.Get(ApplicationConfigFieldLabel).(string) {
		t.Fatal("Expected Label to be identical")
	}
	if model.Scope != resourceData.Get(ApplicationConfigFieldScope).(string) {
		t.Fatal("Expected Scope to be identical")
	}
	if resourceData.Get(ApplicationConfigFieldMatchSpecification).(string) != defaultMatchSpecification {
		t.Fatal("Expected MatchSpecification to be identical")
	}
}

func createTestApplicationConfigModelWithRollup() restapi.ApplicationConfig {
	cfg := createBaseTestApplicationConfigModel()
	cfg.Scope = ApplicationConfigScopeIncludeNoDownstream
	return cfg
}

func createBaseTestApplicationConfigModel() restapi.ApplicationConfig {
	comparision1 := restapi.NewComparisionExpression("entity.name", restapi.ContainsOperator, "foo")
	comparision2 := restapi.NewComparisionExpression("entity.type", restapi.EqualsOperator, "mysql")
	comparision3 := restapi.NewComparisionExpression("entity.type", restapi.EqualsOperator, "elasticsearch")
	logicalAnd := restapi.NewBinaryOperator(comparision1, restapi.LogicalAnd, comparision2)
	logicalOr := restapi.NewBinaryOperator(logicalAnd, restapi.LogicalOr, comparision3)
	return restapi.ApplicationConfig{
		ID:                 "id",
		Label:              "label",
		MatchSpecification: logicalOr,
	}
}

func createFullTestApplicationConfigData() map[string]interface{} {
	data := make(map[string]interface{})
	data[ApplicationConfigFieldLabel] = "label"
	data[ApplicationConfigFieldScope] = ApplicationConfigScopeIncludeNoDownstream
	data[ApplicationConfigFieldMatchSpecification] = defaultMatchSpecification
	return data
}
