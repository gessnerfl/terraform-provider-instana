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
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	testutils "github.com/gessnerfl/terraform-provider-instana/test-utils"
)

var testRuleProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceRuleDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_rule" "example" {
  name = "name"
  entity_type = "entity_type"
  metric_name = "metric_name"
  rollup = 100
  window = 20000
  aggregation = "sum"
  condition_operator = ">"
  condition_value = 1.1
}
`

const ruleApiPath = restapi.RulesResourcePath + "/{id}"
const testRuleDefinition = "instana_rule.example"

func TestCRUDOfRuleResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, ruleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, ruleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, ruleApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"name" : "name",
			"entityType" : "entity_type",
			"metricName" : "metric_name",
			"rollup" : 100,
			"window" : 20000,
			"aggregation" : "sum",
			"conditionOperator" : ">",
			"conditionValue" : 1.1
		}
		`, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceRuleDefinition := strings.ReplaceAll(resourceRuleDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testRuleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceRuleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testRuleDefinition, "id"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldName, "name"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldEntityType, "entity_type"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldMetricName, "metric_name"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldRollup, "100"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldWindow, "20000"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldAggregation, "sum"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldConditionOperator, ">"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldConditionValue, "1.1"),
				),
			},
		},
	})
}

func TestResourceRuleDefinition(t *testing.T) {
	resource := CreateResourceRule()

	validateRuleResourceSchema(resource.Schema, t)

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

func validateRuleResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	validateRequiredSchemaOfTypeString(RuleFieldName, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleFieldEntityType, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleFieldMetricName, schemaMap, t)
	validateOptionalSchemaOfTypeInt(RuleFieldRollup, schemaMap, t)
	validateRequiredSchemaOfTypeInt(RuleFieldWindow, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleFieldAggregation, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleFieldConditionOperator, schemaMap, t)
	validateRequiredSchemaOfTypeFloat(RuleFieldConditionValue, schemaMap, t)
}

func TestShouldSuccessfullyReadRuleFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := createBaseTestRuleModel()
	testShouldSuccessfullyReadRuleFromInstanaAPI(expectedModel, t)
}

func TestShouldSuccessfullyReadRuleFromInstanaAPIWhenBaseDataWithRollupIsReturned(t *testing.T) {
	expectedModel := createTestRuleModelWithRollup()
	testShouldSuccessfullyReadRuleFromInstanaAPI(expectedModel, t)
}

func testShouldSuccessfullyReadRuleFromInstanaAPI(expectedModel restapi.Rule, t *testing.T) {
	resourceData := createEmptyRuleResourceData(t)
	ruleID := "rule-id"
	resourceData.SetId(ruleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().GetOne(gomock.Eq(ruleID)).Return(expectedModel, nil).Times(1)

	err := ReadRule(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	verifyRuleModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldFailToReadRuleFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	resourceData := createEmptyRuleResourceData(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	err := ReadRule(resourceData, mockInstanaAPI)

	if err == nil || !strings.HasPrefix(err.Error(), "ID of rule") {
		t.Fatal("Expected error to occur because of missing id")
	}
}

func TestShouldFailToReadRuleFromInstanaAPIAndDeleteResourceWhenRoleDoesNotExist(t *testing.T) {
	resourceData := createEmptyRuleResourceData(t)
	ruleID := "rule-id"
	resourceData.SetId(ruleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().GetOne(gomock.Eq(ruleID)).Return(restapi.Rule{}, restapi.ErrEntityNotFound).Times(1)

	err := ReadRule(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldFailToReadRuleFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	resourceData := createEmptyRuleResourceData(t)
	ruleID := "rule-id"
	resourceData.SetId(ruleID)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().GetOne(gomock.Eq(ruleID)).Return(restapi.Rule{}, expectedError).Times(1)

	err := ReadRule(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error should be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID should still be set")
	}
}

func TestShouldCreateRuleThroughInstanaAPI(t *testing.T) {
	data := createFullTestRuleData()
	resourceData := createRuleResourceData(t, data)
	expectedModel := createTestRuleModelWithRollup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.Rule{})).Return(expectedModel, nil).Times(1)

	err := CreateRule(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	verifyRuleModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldReturnErrorWhenCreateRuleFailsThroughInstanaAPI(t *testing.T) {
	data := createFullTestRuleData()
	resourceData := createRuleResourceData(t, data)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.Rule{})).Return(restapi.Rule{}, expectedError).Times(1)

	err := CreateRule(resourceData, mockInstanaAPI)

	if err == nil || expectedError != err {
		t.Fatal("Expected definned error to be returned")
	}
}

func TestShouldDeleteRuleThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestRuleData()
	resourceData := createRuleResourceData(t, data)
	resourceData.SetId(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(nil).Times(1)

	err := DeleteRule(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldReturnErrorWhenDeleteRuleFailsThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestRuleData()
	resourceData := createRuleResourceData(t, data)
	resourceData.SetId(id)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(expectedError).Times(1)

	err := DeleteRule(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error to be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
	}
}

func verifyRuleModelAppliedToResource(model restapi.Rule, resourceData *schema.ResourceData, t *testing.T) {
	if model.ID != resourceData.Id() {
		t.Fatal("Expected ID to be identical")
	}
	if model.Name != resourceData.Get(RuleFieldName).(string) {
		t.Fatal("Expected Name to be identical")
	}
	if model.EntityType != resourceData.Get(RuleFieldEntityType).(string) {
		t.Fatal("Expected EntityType to be identical")
	}
	if model.MetricName != resourceData.Get(RuleFieldMetricName).(string) {
		t.Fatal("Expected MetricName to be identical")
	}
	if model.Rollup != resourceData.Get(RuleFieldRollup).(int) {
		t.Fatal("Expected Rollup to be identical")
	}
	if model.Window != resourceData.Get(RuleFieldWindow).(int) {
		t.Fatal("Expected Window to be identical")
	}
	if model.Aggregation != resourceData.Get(RuleFieldAggregation).(string) {
		t.Fatal("Expected Aggregation to be identical")
	}
	if model.ConditionOperator != resourceData.Get(RuleFieldConditionOperator).(string) {
		t.Fatal("Expected ConditionOperator to be identical")
	}
	if model.ConditionValue != resourceData.Get(RuleFieldConditionValue).(float64) {
		t.Fatal("Expected ConditionValue to be identical")
	}
}

func createTestRuleModelWithRollup() restapi.Rule {
	data := createBaseTestRuleModel()
	data.Rollup = 1234
	return data
}

func createBaseTestRuleModel() restapi.Rule {
	return restapi.Rule{
		ID:                "id",
		Name:              "name",
		EntityType:        "entityType",
		MetricName:        "metricName",
		Window:            9876,
		Aggregation:       "sum",
		ConditionOperator: ">",
		ConditionValue:    1.1,
	}
}

func createFullTestRuleData() map[string]interface{} {
	data := make(map[string]interface{})
	data[RuleFieldName] = "name"
	data[RuleFieldEntityType] = "entityType"
	data[RuleFieldMetricName] = "metricName"
	data[RuleFieldRollup] = 1234
	data[RuleFieldWindow] = 9876
	data[RuleFieldAggregation] = "sum"
	data[RuleFieldConditionOperator] = ">"
	data[RuleFieldConditionValue] = 1.1
	return data
}
