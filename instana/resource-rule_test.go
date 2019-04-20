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
const ruleID = "rule-id"
const metricNameFieldValue = "metric_name"
const entityTypeFieldValue = "entity_type"
const ruleNameFieldValue = "name"
const aggregationFieldValue = "sum"
const conditionOperatorFieldValue = ">"

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
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldName, ruleNameFieldValue),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldEntityType, entityTypeFieldValue),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldMetricName, metricNameFieldValue),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldRollup, "100"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldWindow, "20000"),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldAggregation, aggregationFieldValue),
					resource.TestCheckResourceAttr(testRuleDefinition, RuleFieldConditionOperator, conditionOperatorFieldValue),
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
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(RuleFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(RuleFieldEntityType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(RuleFieldMetricName)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(RuleFieldRollup)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeInt(RuleFieldWindow)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(RuleFieldAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(RuleFieldConditionOperator)
	schemaAssert.AssertSchemaIsRequiredAndTypeFloat(RuleFieldConditionValue)
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
	resourceData := NewTestHelper(t).CreateEmptyRuleResourceData()
	resourceData.SetId(ruleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().GetOne(gomock.Eq(ruleID)).Return(expectedModel, nil).Times(1)

	err := ReadRule(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	verifyRuleModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldFailToReadRuleFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyRuleResourceData()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	err := ReadRule(resourceData, mockInstanaAPI)

	if err == nil || !strings.HasPrefix(err.Error(), "ID of rule") {
		t.Fatal("Expected error to occur because of missing id")
	}
}

func TestShouldFailToReadRuleFromInstanaAPIAndDeleteResourceWhenRoleDoesNotExist(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyRuleResourceData()
	resourceData.SetId(ruleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().GetOne(gomock.Eq(ruleID)).Return(restapi.Rule{}, restapi.ErrEntityNotFound).Times(1)

	err := ReadRule(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldFailToReadRuleFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyRuleResourceData()
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
	resourceData := NewTestHelper(t).CreateRuleResourceData(data)
	expectedModel := createTestRuleModelWithRollup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.Rule{})).Return(expectedModel, nil).Times(1)

	err := CreateRule(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	verifyRuleModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldReturnErrorWhenCreateRuleFailsThroughInstanaAPI(t *testing.T) {
	data := createFullTestRuleData()
	resourceData := NewTestHelper(t).CreateRuleResourceData(data)
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
	resourceData := NewTestHelper(t).CreateRuleResourceData(data)
	resourceData.SetId(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleApi := mocks.NewMockRuleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().Rules().Return(mockRuleApi).Times(1)
	mockRuleApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(nil).Times(1)

	err := DeleteRule(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldReturnErrorWhenDeleteRuleFailsThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestRuleData()
	resourceData := NewTestHelper(t).CreateRuleResourceData(data)
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
		Name:              ruleNameFieldValue,
		EntityType:        entityTypeFieldValue,
		MetricName:        metricNameFieldValue,
		Window:            9876,
		Aggregation:       aggregationFieldValue,
		ConditionOperator: conditionOperatorFieldValue,
		ConditionValue:    1.1,
	}
}

func createFullTestRuleData() map[string]interface{} {
	data := make(map[string]interface{})
	data[RuleFieldName] = ruleNameFieldValue
	data[RuleFieldEntityType] = entityTypeFieldValue
	data[RuleFieldMetricName] = metricNameFieldValue
	data[RuleFieldRollup] = 1234
	data[RuleFieldWindow] = 9876
	data[RuleFieldAggregation] = aggregationFieldValue
	data[RuleFieldConditionOperator] = conditionOperatorFieldValue
	data[RuleFieldConditionValue] = 1.1
	return data
}
