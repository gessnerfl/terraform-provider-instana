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

var testCustomEventSpecificationWithThresholdRuleProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_custom_event_spec_threshold_rule" "rollup" {
  name = "name"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rule_severity = "warning"
  rule_metric_name = "metric_name"
  rule_rollup = "40000"
  rule_condition_operator = "=="
  rule_condition_value = "1.2"
  downstream_integration_ids = [ "integration-id-1", "integration-id-2" ]
  downstream_broadcast_to_all_alerting_configs = true
}
`

const resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_custom_event_spec_threshold_rule" "window" {
  name = "name"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = "60000"
  rule_severity = "warning"
  rule_metric_name = "metric_name"
  rule_window = "60000"
  rule_aggregation = "sum"
  rule_condition_operator = "=="
  rule_condition_value = "1.2"
  downstream_integration_ids = [ "integration-id-1", "integration-id-2" ]
  downstream_broadcast_to_all_alerting_configs = true
}
`

const (
	customEventSpecificationWithThresholdRuleApiPath              = restapi.CustomEventSpecificationResourcePath + "/{id}"
	testCustomEventSpecificationWithThresholdRuleRollupDefinition = "instana_custom_event_spec_threshold_rule.rollup"
	testCustomEventSpecificationWithThresholdRuleWindowDefinition = "instana_custom_event_spec_threshold_rule.window"

	customEventSpecificationWithThresholdRuleID                       = "custom-system-event-id"
	customEventSpecificationWithThresholdRuleName                     = "name"
	customEventSpecificationWithThresholdRuleEntityType               = "entity_type"
	customEventSpecificationWithThresholdRuleQuery                    = "query"
	customEventSpecificationWithThresholdRuleExpirationTime           = 60000
	customEventSpecificationWithThresholdRuleDescription              = "description"
	customEventSpecificationWithThresholdRuleMetricName               = "metric_name"
	customEventSpecificationWithThresholdRuleRollup                   = 40000
	customEventSpecificationWithThresholdRuleWindow                   = 60000
	customEventSpecificationWithThresholdRuleAggregation              = restapi.AggregationSum
	customEventSpecificationWithThresholdRuleConditionOperator        = restapi.ConditionOperatorEquals
	customEventSpecificationWithThresholdRuleConditionValue           = float64(1.2)
	customEventSpecificationWithThresholdRuleDownstreamIntegrationId1 = "integration-id-1"
	customEventSpecificationWithThresholdRuleDownstreamIntegrationId2 = "integration-id-2"
)

var CustomEventSpecificationWithThresholdRuleRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithRollupResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, customEventSpecificationWithThresholdRuleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, customEventSpecificationWithThresholdRuleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, customEventSpecificationWithThresholdRuleApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"name" : "name",
			"entityType" : "entity_type",
			"query" : "query",
			"enabled" : true,
			"triggering" : true,
			"description" : "description",
			"expirationTime" : 60000,
			"rules" : [ { "ruleType" : "threshold", "severity" : 5, "metricName" : "metric_name", "rollup" : 40000, "conditionOperator" : "==", "conditionValue" : 1.2 } ],
			"downstream" : {
				"integrationIds" : ["integration-id-1", "integration-id-2"],
				"broadcastToAllAlertingConfigs" : true
			}
		}
		`, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceCustomEventSpecificationWithThresholdRuleDefinition := strings.ReplaceAll(resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testCustomEventSpecificationWithThresholdRuleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceCustomEventSpecificationWithThresholdRuleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithThresholdRuleRollupDefinition, "id"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationFieldName, customEventSpecificationWithThresholdRuleName),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationFieldTriggering, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEventSpecificationWithThresholdRuleExpirationTime)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationFieldEnabled, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationDownstreamIntegrationIds+".0", customEventSpecificationWithThresholdRuleDownstreamIntegrationId1),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationDownstreamIntegrationIds+".1", customEventSpecificationWithThresholdRuleDownstreamIntegrationId2),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, CustomEventSpecificationRuleSeverity, CustomEventSpecificationWithThresholdRuleRuleSeverity),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, ThresholdRuleFieldRollup, strconv.FormatInt(customEventSpecificationWithThresholdRuleRollup, 10)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, ThresholdRuleFieldConditionOperator, string(customEventSpecificationWithThresholdRuleConditionOperator)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleRollupDefinition, ThresholdRuleFieldConditionValue, "1.2"),
				),
			},
		},
	})
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithWindowResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, customEventSpecificationWithThresholdRuleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, customEventSpecificationWithThresholdRuleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, customEventSpecificationWithThresholdRuleApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"name" : "name",
			"entityType" : "entity_type",
			"query" : "query",
			"enabled" : true,
			"triggering" : true,
			"description" : "description",
			"expirationTime" : 60000,
			"rules" : [ { "ruleType" : "threshold", "severity" : 5, "metricName": "metric_name", "window" : 60000, "aggregation": "sum", "conditionOperator" : "==", "conditionValue" : 1.2 } ],
			"downstream" : {
				"integrationIds" : ["integration-id-1", "integration-id-2"],
				"broadcastToAllAlertingConfigs" : true
			}
		}
		`, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceCustomEventSpecificationWithThresholdRuleDefinition := strings.ReplaceAll(resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testCustomEventSpecificationWithThresholdRuleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceCustomEventSpecificationWithThresholdRuleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithThresholdRuleWindowDefinition, "id"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationFieldName, customEventSpecificationWithThresholdRuleName),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationFieldTriggering, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEventSpecificationWithThresholdRuleExpirationTime)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationFieldEnabled, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationDownstreamIntegrationIds+".0", customEventSpecificationWithThresholdRuleDownstreamIntegrationId1),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationDownstreamIntegrationIds+".1", customEventSpecificationWithThresholdRuleDownstreamIntegrationId2),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, CustomEventSpecificationRuleSeverity, CustomEventSpecificationWithThresholdRuleRuleSeverity),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, ThresholdRuleFieldConditionOperator, string(customEventSpecificationWithThresholdRuleConditionOperator)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleWindowDefinition, ThresholdRuleFieldConditionValue, "1.2"),
				),
			},
		},
	})
}

func TestResourceCustomEventSpecificationWithThresholdRuleDefinition(t *testing.T) {
	resource := CreateResourceCustomEventSpecificationWithThresholdRule()

	validateCustomEventSpecificationWithThresholdRuleResourceSchema(resource.Schema, t)

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

func validateCustomEventSpecificationWithThresholdRuleResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSChemaIsRequiredAndOfTypeListOfStrings(CustomEventSpecificationDownstreamIntegrationIds)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, true)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleSeverity)

	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ThresholdRuleFieldMetricName)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(ThresholdRuleFieldWindow)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(ThresholdRuleFieldRollup)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(ThresholdRuleFieldAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(ThresholdRuleFieldConditionOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeFloat(ThresholdRuleFieldConditionValue)
}

func TestShouldSuccessfullyReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := createBaseTestCustomEventSpecificationWithThresholdRuleModel()
	testShouldSuccessfullyReadCustomEventSpecificationWithThresholdRuleFromInstanaAPI(expectedModel, t)
}

func TestShouldSuccessfullyReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIWhenFullDataIsReturned(t *testing.T) {
	expectedModel := createTestCustomEventSpecificationWithThresholdRuleModelWithFullDataSet()
	testShouldSuccessfullyReadCustomEventSpecificationWithThresholdRuleFromInstanaAPI(expectedModel, t)
}

func testShouldSuccessfullyReadCustomEventSpecificationWithThresholdRuleFromInstanaAPI(expectedModel restapi.CustomEventSpecification, t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()
	resourceData.SetId(customEventSpecificationWithThresholdRuleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(expectedModel, nil).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Read(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	verifyCustomEventSpecificationWithThresholdRuleModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldFailToReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Read(resourceData, mockInstanaAPI)

	if err == nil || !strings.HasPrefix(err.Error(), "ID of custom event specification") {
		t.Fatal("Expected error to occur because of missing id")
	}
}

func TestShouldFailToReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIAndDeleteResourceWhenCustomEventDoesNotExist(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()
	resourceData.SetId(customEventSpecificationWithThresholdRuleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(restapi.CustomEventSpecification{}, restapi.ErrEntityNotFound).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Read(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldFailToReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()
	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(restapi.CustomEventSpecification{}, expectedError).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Read(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error should be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID should still be set")
	}
}

func TestShouldFailToReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIWhenSeverityFromAPICannotBeMappedToSeverityOfTerraformState(t *testing.T) {
	expectedModel := createTestCustomEventSpecificationWithThresholdRuleModelWithFullDataSet()
	expectedModel.Rules[0].Severity = 999
	resourceData := NewTestHelper(t).CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()
	resourceData.SetId(customEventSpecificationWithThresholdRuleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(expectedModel, nil).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Read(resourceData, mockInstanaAPI)

	if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
		t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
	}
}

func TestShouldCreateCustomEventSpecificationWithThresholdRuleThroughInstanaAPI(t *testing.T) {
	data := createFullTestCustomEventSpecificationWithThresholdRuleData()
	resourceData := NewTestHelper(t).CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
	expectedModel := createTestCustomEventSpecificationWithThresholdRuleModelWithFullDataSet()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(expectedModel, nil).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Create(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	verifyCustomEventSpecificationWithThresholdRuleModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithThresholdRuleFailsThroughInstanaAPI(t *testing.T) {
	data := createFullTestCustomEventSpecificationWithThresholdRuleData()
	resourceData := NewTestHelper(t).CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(restapi.CustomEventSpecification{}, expectedError).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Create(resourceData, mockInstanaAPI)

	if err == nil || expectedError != err {
		t.Fatal("Expected definned error to be returned")
	}
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithThresholdRuleFailsBecauseOfInvalidSeverityConfiguredInTerraform(t *testing.T) {
	data := createFullTestCustomEventSpecificationWithThresholdRuleData()
	data[CustomEventSpecificationRuleSeverity] = "invalid"
	resourceData := NewTestHelper(t).CreateCustomEventSpecificationWithThresholdRuleResourceData(data)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Create(resourceData, mockInstanaAPI)

	if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
		t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
	}
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithThresholdRuleFailsBecauseOfInvalidSeverityReturnedFromInstanaAPI(t *testing.T) {
	data := createFullTestCustomEventSpecificationWithThresholdRuleData()
	resourceData := NewTestHelper(t).CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
	expectedModel := createTestCustomEventSpecificationWithThresholdRuleModelWithFullDataSet()
	expectedModel.Rules[0].Severity = 999

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(expectedModel, nil).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Create(resourceData, mockInstanaAPI)

	if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
		t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
	}
}

func TestShouldDeleteCustomEventSpecificationWithThresholdRuleThroughInstanaAPI(t *testing.T) {
	data := createFullTestCustomEventSpecificationWithThresholdRuleData()
	resourceData := NewTestHelper(t).CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
	resourceData.SetId(customEventSpecificationWithThresholdRuleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().DeleteByID(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(nil).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Delete(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldReturnErrorWhenDeleteCustomEventSpecificationWithThresholdRuleFailsThroughInstanaAPI(t *testing.T) {
	data := createFullTestCustomEventSpecificationWithThresholdRuleData()
	resourceData := NewTestHelper(t).CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
	resourceData.SetId(customEventSpecificationWithThresholdRuleID)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
	mockCustomEventAPI.EXPECT().DeleteByID(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(expectedError).Times(1)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Delete(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error to be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
	}
}

func TestShouldFailToDeleteCustomEventSpecificationWithThresholdRuleWhenInvalidSeverityIsConfiguredInTerraform(t *testing.T) {
	data := createFullTestCustomEventSpecificationWithThresholdRuleData()
	data[CustomEventSpecificationRuleSeverity] = "invalid"
	resourceData := NewTestHelper(t).CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
	resourceData.SetId(customEventSpecificationWithThresholdRuleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	resource := CreateResourceCustomEventSpecificationWithThresholdRule()
	err := resource.Delete(resourceData, mockInstanaAPI)

	if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
		t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
	}
}

func verifyCustomEventSpecificationWithThresholdRuleModelAppliedToResource(model restapi.CustomEventSpecification, resourceData *schema.ResourceData, t *testing.T) {
	verifyCustomEventSpecificationModelAppliedToResource(model, resourceData, t)
	verifyCustomEventSpecificationDownstreamModelAppliedToResource(model, resourceData, t)

	ruleSpec := model.Rules[0]
	convertedSeverity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if convertedSeverity != resourceData.Get(CustomEventSpecificationRuleSeverity).(string) {
		t.Fatal("Expected Severity to be identical")
	}

	if ruleSpec.MetricName != resourceData.Get(ThresholdRuleFieldMetricName).(string) {
		t.Fatal("Expected metric name to be identical")
	}
	if *ruleSpec.Window != resourceData.Get(ThresholdRuleFieldWindow).(int) {
		t.Fatal("Expected window to be identical")
	}
	if *ruleSpec.Aggregation != restapi.AggregationType(resourceData.Get(ThresholdRuleFieldAggregation).(string)) {
		t.Fatal("Expected aggregation to be identical")
	}
	if ruleSpec.ConditionOperator != restapi.ConditionOperatorType(resourceData.Get(ThresholdRuleFieldConditionOperator).(string)) {
		t.Fatal("Expected condition operator to be identical")
	}
	if *ruleSpec.ConditionValue != resourceData.Get(ThresholdRuleFieldConditionValue).(float64) {
		t.Fatal("Expected System Rule ID to be identical")
	}
}

func createTestCustomEventSpecificationWithThresholdRuleModelWithFullDataSet() restapi.CustomEventSpecification {
	description := customEventSpecificationWithThresholdRuleDescription
	expirationTime := customEventSpecificationWithThresholdRuleExpirationTime
	query := customEventSpecificationWithThresholdRuleQuery

	data := createBaseTestCustomEventSpecificationWithThresholdRuleModel()
	data.Query = &query
	data.Description = &description
	data.ExpirationTime = &expirationTime
	data.Downstream = &restapi.EventSpecificationDownstream{
		IntegrationIds:                []string{customEventSpecificationWithThresholdRuleDownstreamIntegrationId1, customEventSpecificationWithThresholdRuleDownstreamIntegrationId2},
		BroadcastToAllAlertingConfigs: true,
	}
	return data
}

func createBaseTestCustomEventSpecificationWithThresholdRuleModel() restapi.CustomEventSpecification {
	window := customEventSpecificationWithThresholdRuleWindow
	aggregation := customEventSpecificationWithThresholdRuleAggregation
	conditionValue := customEventSpecificationWithThresholdRuleConditionValue

	return restapi.CustomEventSpecification{
		ID:         customEventSpecificationWithThresholdRuleID,
		Name:       customEventSpecificationWithThresholdRuleName,
		EntityType: customEventSpecificationWithThresholdRuleEntityType,
		Triggering: false,
		Enabled:    true,
		Rules: []restapi.RuleSpecification{
			restapi.RuleSpecification{
				DType:             restapi.ThresholdRuleType,
				Severity:          restapi.SeverityWarning.GetAPIRepresentation(),
				MetricName:        customEventSpecificationWithThresholdRuleMetricName,
				Window:            &window,
				Aggregation:       &aggregation,
				ConditionOperator: customEventSpecificationWithThresholdRuleConditionOperator,
				ConditionValue:    &conditionValue,
			},
		},
	}
}

func createFullTestCustomEventSpecificationWithThresholdRuleData() map[string]interface{} {
	data := make(map[string]interface{})
	data[CustomEventSpecificationFieldName] = customEventSpecificationWithThresholdRuleName
	data[CustomEventSpecificationFieldEntityType] = customEventSpecificationWithThresholdRuleEntityType
	data[CustomEventSpecificationFieldQuery] = customEventSpecificationWithThresholdRuleQuery
	data[CustomEventSpecificationFieldTriggering] = "true"
	data[CustomEventSpecificationFieldDescription] = customEventSpecificationWithThresholdRuleDescription
	data[CustomEventSpecificationFieldExpirationTime] = customEventSpecificationWithThresholdRuleExpirationTime
	data[CustomEventSpecificationFieldEnabled] = "true"
	data[CustomEventSpecificationDownstreamIntegrationIds] = []string{customEventSpecificationWithThresholdRuleDownstreamIntegrationId1, customEventSpecificationWithThresholdRuleDownstreamIntegrationId2}
	data[CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs] = "true"
	data[CustomEventSpecificationRuleSeverity] = CustomEventSpecificationWithThresholdRuleRuleSeverity
	data[ThresholdRuleFieldMetricName] = customEventSpecificationWithThresholdRuleMetricName
	data[ThresholdRuleFieldWindow] = customEventSpecificationWithThresholdRuleWindow
	data[ThresholdRuleFieldAggregation] = customEventSpecificationWithThresholdRuleAggregation
	data[ThresholdRuleFieldConditionOperator] = customEventSpecificationWithThresholdRuleConditionOperator
	data[ThresholdRuleFieldConditionValue] = customEventSpecificationWithThresholdRuleConditionValue
	return data
}
