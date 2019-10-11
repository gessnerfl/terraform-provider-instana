package instana_test

import (
	"errors"
	"fmt"
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

resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name {{ITERATION}}"
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

resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "name {{ITERATION}}"
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
	customEventSpecificationWithThresholdRuleApiPath        = restapi.CustomEventSpecificationResourcePath + "/{id}"
	testCustomEventSpecificationWithThresholdRuleDefinition = "instana_custom_event_spec_threshold_rule.example"

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
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName" : "metric_name", "rollup" : 40000, "conditionOperator" : "==", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		resourceCustomEventSpecificationWithThresholdRuleAndRollupDefinitionTemplate,
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldRollup, strconv.FormatInt(customEventSpecificationWithThresholdRuleRollup, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(customEventSpecificationWithThresholdRuleConditionOperator)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

func TestCRUDOfCustomEventSpecificationWithThresholdRuleWithWindowResourceWithMockServer(t *testing.T) {
	ruleAsJson := `{ "ruleType" : "threshold", "severity" : 5, "metricName": "metric_name", "window" : 60000, "aggregation": "sum", "conditionOperator" : "==", "conditionValue" : 1.2 }`
	testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(
		t,
		resourceCustomEventSpecificationWithThresholdRuleAndWindowDefinitionTemplate,
		ruleAsJson,
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldMetricName, customEventSpecificationWithThresholdRuleMetricName),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldWindow, strconv.FormatInt(customEventSpecificationWithThresholdRuleWindow, 10)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldAggregation, string(customEventSpecificationWithThresholdRuleAggregation)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionOperator, string(customEventSpecificationWithThresholdRuleConditionOperator)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, ThresholdRuleFieldConditionValue, "1.2"),
	)
}

const httpServerResponseTemplate = `
{
	"id" : "{{id}}",
	"name" : "name (TF managed)",
	"entityType" : "entity_type",
	"query" : "query",
	"enabled" : true,
	"triggering" : true,
	"description" : "description",
	"expirationTime" : 60000,
	"rules" : [ {{rule}} ],
	"downstream" : {
		"integrationIds" : ["integration-id-1", "integration-id-2"],
		"broadcastToAllAlertingConfigs" : true
	}
}
`

func testCRUDOfResourceCustomEventSpecificationThresholdRuleResourceWithMockServer(t *testing.T, terraformDefinition, ruleAsJson string, ruleTestCheckFunctions ...resource.TestCheckFunc) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, customEventSpecificationWithThresholdRuleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, customEventSpecificationWithThresholdRuleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, customEventSpecificationWithThresholdRuleApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(strings.ReplaceAll(httpServerResponseTemplate, "{{id}}", vars["id"]), "{{rule}}", ruleAsJson)
		w.Header().Set(constSystemEventContentType, r.Header.Get(constSystemEventContentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	completeTerraformDefinitionWithoutName := strings.ReplaceAll(terraformDefinition, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	completeTerraformDefinitionWithName1 := strings.ReplaceAll(completeTerraformDefinitionWithoutName, "{{ITERATION}}", "0")
	completeTerraformDefinitionWithName2 := strings.ReplaceAll(completeTerraformDefinitionWithoutName, "{{ITERATION}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testCustomEventSpecificationWithThresholdRuleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: completeTerraformDefinitionWithName1,
				Check:  resource.ComposeTestCheckFunc(createTestCheckFunctions(ruleTestCheckFunctions, 0)...),
			},
			resource.TestStep{
				Config: completeTerraformDefinitionWithName2,
				Check:  resource.ComposeTestCheckFunc(createTestCheckFunctions(ruleTestCheckFunctions, 1)...),
			},
		},
	})
}

func createTestCheckFunctions(ruleTestCheckFunctions []resource.TestCheckFunc, iteration int) []resource.TestCheckFunc {
	defaultCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithThresholdRuleDefinition, "id"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldName, customEventSpecificationWithThresholdRuleName+fmt.Sprintf(" %d", iteration)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldFullName, customEventSpecificationWithThresholdRuleName+fmt.Sprintf(" %d%s", iteration, TerraformManagedResourceNameSuffix)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldEntityType, customEventSpecificationWithThresholdRuleEntityType),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldQuery, customEventSpecificationWithThresholdRuleQuery),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldTriggering, "true"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldDescription, customEventSpecificationWithThresholdRuleDescription),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEventSpecificationWithThresholdRuleExpirationTime)),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationFieldEnabled, "true"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationDownstreamIntegrationIds+".0", customEventSpecificationWithThresholdRuleDownstreamIntegrationId1),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationDownstreamIntegrationIds+".1", customEventSpecificationWithThresholdRuleDownstreamIntegrationId2),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, "true"),
		resource.TestCheckResourceAttr(testCustomEventSpecificationWithThresholdRuleDefinition, CustomEventSpecificationRuleSeverity, CustomEventSpecificationWithThresholdRuleRuleSeverity),
	}
	allFunctions := append(defaultCheckFunctions, ruleTestCheckFunctions...)
	return allFunctions
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
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSChemaIsOptionalAndOfTypeListOfStrings(CustomEventSpecificationDownstreamIntegrationIds)
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
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(expectedModel, nil).Times(1)
		resource := CreateResourceCustomEventSpecificationWithThresholdRule()

		err := resource.Read(resourceData, providerMeta)
		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyCustomEventSpecificationWithThresholdRuleModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldFailToReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIWhenSeverityFromAPICannotBeMappedToSeverityOfTerraformState(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		expectedModel := createTestCustomEventSpecificationWithThresholdRuleModelWithFullDataSet()
		expectedModel.Rules[0].Severity = 999
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(expectedModel, nil).Times(1)
		resource := CreateResourceCustomEventSpecificationWithThresholdRule()

		err := resource.Read(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
			t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Read(resourceData, providerMeta)

		if err == nil || !strings.HasPrefix(err.Error(), "ID of custom event specification") {
			t.Fatal("Expected error to occur because of missing id")
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIAndDeleteResourceWhenCustomEventDoesNotExist(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(restapi.CustomEventSpecification{}, restapi.ErrEntityNotFound).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Read(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithThresholdRuleFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithThresholdRuleResourceData()
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(restapi.CustomEventSpecification{}, expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Read(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error should be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID should still be set")
		}
	})
}

func TestShouldCreateCustomEventSpecificationWithThresholdRuleThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithThresholdRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
		expectedModel := createTestCustomEventSpecificationWithThresholdRuleModelWithFullDataSet()

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Create(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyCustomEventSpecificationWithThresholdRuleModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithThresholdRuleFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithThresholdRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(restapi.CustomEventSpecification{}, expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || expectedError != err {
			t.Fatal("Expected definned error to be returned")
		}
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithThresholdRuleFailsBecauseOfInvalidSeverityConfiguredInTerraform(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithThresholdRuleData()
		data[CustomEventSpecificationRuleSeverity] = "invalid"
		resourceData := testHelper.CreateCustomEventSpecificationWithThresholdRuleResourceData(data)

		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
			t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithThresholdRuleFailsBecauseOfInvalidSeverityReturnedFromInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithThresholdRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
		expectedModel := createTestCustomEventSpecificationWithThresholdRuleModelWithFullDataSet()
		expectedModel.Rules[0].Severity = 999

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
			t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldDeleteCustomEventSpecificationWithThresholdRuleThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithThresholdRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().DeleteByID(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Delete(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldReturnErrorWhenDeleteCustomEventSpecificationWithThresholdRuleFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithThresholdRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().DeleteByID(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Delete(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error to be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
		}
	})
}

func TestShouldFailToDeleteCustomEventSpecificationWithThresholdRuleWhenInvalidSeverityIsConfiguredInTerraform(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithThresholdRuleData()
		data[CustomEventSpecificationRuleSeverity] = "invalid"
		resourceData := testHelper.CreateCustomEventSpecificationWithThresholdRuleResourceData(data)
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)

		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)

		resource := CreateResourceCustomEventSpecificationWithThresholdRule()
		err := resource.Delete(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
			t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
		}
	})
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
