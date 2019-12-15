package instana_test

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

var testCustomEventSpecificationWithSystemRuleProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceCustomEventSpecificationWithSystemRuleDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_custom_event_spec_system_rule" "example" {
  name = "name"
  entity_type = "entity_type"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
	rule_severity = "warning"
	rule_system_rule_id = "system-rule-id"
	downstream_integration_ids = [ "integration-id-1", "integration-id-2" ]
	downstream_broadcast_to_all_alerting_configs = true
}
`

const (
	customSystemEventApiPath                             = restapi.CustomEventSpecificationResourcePath + "/{id}"
	testCustomEventSpecificationWithSystemRuleDefinition = "instana_custom_event_spec_system_rule.example"

	customSystemEventID                       = "custom-system-event-id"
	customSystemEventName                     = "name"
	customSystemEventEntityType               = "entity_type"
	customSystemEventQuery                    = "query"
	customSystemEventExpirationTime           = 60000
	customSystemEventDescription              = "description"
	customSystemEventRuleSystemRuleId         = "system-rule-id"
	customSystemEventDownStringIntegrationId1 = "integration-id-1"
	customSystemEventDownStringIntegrationId2 = "integration-id-2"

	customSystemEventMessageNotAValidSeverity           = "not a valid severity"
	customSystemEventTestMessageExpectedInvalidSeverity = "Expected to get error that the provided severity is not valid"

	constSystemEventContentType = "Content-Type"
)

var customSystemEventRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCreateResourceCustomEventSpecificationWithThresholdRuleResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, customSystemEventApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, customSystemEventApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, customSystemEventApiPath, func(w http.ResponseWriter, r *http.Request) {
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
			"rules" : [ { "ruleType" : "system", "severity" : 5, "systemRuleId" : "system-rule-id" } ],
			"downstream" : {
				"integrationIds" : ["integration-id-1", "integration-id-2"],
				"broadcastToAllAlertingConfigs" : true
			}
		}
		`, "{{id}}", vars["id"])
		w.Header().Set(constSystemEventContentType, r.Header.Get(constSystemEventContentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceCustomEventSpecificationWithSystemRuleDefinition := strings.ReplaceAll(resourceCustomEventSpecificationWithSystemRuleDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testCustomEventSpecificationWithSystemRuleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceCustomEventSpecificationWithSystemRuleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithSystemRuleDefinition, "id"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldName, customSystemEventName),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldEntityType, customSystemEventEntityType),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldQuery, customSystemEventQuery),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldTriggering, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldDescription, customSystemEventDescription),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customSystemEventExpirationTime)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldEnabled, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationDownstreamIntegrationIds+".0", customSystemEventDownStringIntegrationId1),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationDownstreamIntegrationIds+".1", customSystemEventDownStringIntegrationId2),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationRuleSeverity, customSystemEventRuleSeverity),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, SystemRuleSpecificationSystemRuleID, customSystemEventRuleSystemRuleId),
				),
			},
		},
	})
}

func TestResourceCustomEventSpecificationWithSystemRuleDefinition(t *testing.T) {
	resource := CreateResourceCustomEventSpecificationWithSystemRule()

	validateCustomEventSpecificationWithSystemRuleResourceSchema(resource.Schema, t)

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

func validateCustomEventSpecificationWithSystemRuleResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
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
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SystemRuleSpecificationSystemRuleID)
}

func TestShouldSuccessfullyReadCustomEventSpecificationWithSystemRuleFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := createBaseTestCustomEventSpecificationWithSystemRuleModel()
	testShouldSuccessfullyReadCustomEventSpecificationWithSystemRuleFromInstanaAPI(expectedModel, t)
}

func TestShouldSuccessfullyReadCustomEventSpecificationWithSystemRuleFromInstanaAPIWhenFullDataIsReturned(t *testing.T) {
	expectedModel := createTestCustomEventSpecificationWithSystemRuleModelWithFullDataSet()
	testShouldSuccessfullyReadCustomEventSpecificationWithSystemRuleFromInstanaAPI(expectedModel, t)
}

func testShouldSuccessfullyReadCustomEventSpecificationWithSystemRuleFromInstanaAPI(expectedModel restapi.CustomEventSpecification, t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithSystemRuleResourceData()
		resourceData.SetId(customSystemEventID)

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customSystemEventID)).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Read(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyCustomEventSpecificationWithSystemRuleModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldFailToReadCustomEventSpecificationWithSystemRuleFromInstanaAPIWhenSeverityFromAPICannotBeMappedToSeverityOfTerraformState(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		expectedModel := createTestCustomEventSpecificationWithSystemRuleModelWithFullDataSet()
		expectedModel.Rules[0].Severity = 999
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithSystemRuleResourceData()
		resourceData.SetId(customSystemEventID)

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customSystemEventID)).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Read(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
			t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithSystemRuleFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithSystemRuleResourceData()

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Read(resourceData, providerMeta)

		if err == nil || !strings.HasPrefix(err.Error(), "ID of custom event specification") {
			t.Fatal("Expected error to occur because of missing id")
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithSystemRuleFromInstanaAPIAndDeleteResourceWhenCustomEventDoesNotExist(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithSystemRuleResourceData()
		resourceData.SetId(customSystemEventID)

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customSystemEventID)).Return(restapi.CustomEventSpecification{}, restapi.ErrEntityNotFound).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Read(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithSystemRuleFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithSystemRuleResourceData()
		resourceData.SetId(customSystemEventID)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customSystemEventID)).Return(restapi.CustomEventSpecification{}, expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Read(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error should be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID should still be set")
		}
	})
}

func TestShouldCreateCustomEventSpecificationWithSystemRuleThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithSystemRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithSystemRuleResourceData(data)
		expectedModel := createTestCustomEventSpecificationWithSystemRuleModelWithFullDataSet()

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Create(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyCustomEventSpecificationWithSystemRuleModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithSystemRuleFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithSystemRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithSystemRuleResourceData(data)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(restapi.CustomEventSpecification{}, expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || expectedError != err {
			t.Fatal("Expected definned error to be returned")
		}
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithSystemRuleFailsBecauseOfInvalidSeverityConfiguredInTerraform(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithSystemRuleData()
		data[CustomEventSpecificationRuleSeverity] = "invalid"
		resourceData := testHelper.CreateCustomEventSpecificationWithSystemRuleResourceData(data)

		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
			t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithSystemRuleFailsBecauseOfInvalidSeverityReturnedFromInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithSystemRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithSystemRuleResourceData(data)
		expectedModel := createTestCustomEventSpecificationWithSystemRuleModelWithFullDataSet()
		expectedModel.Rules[0].Severity = 999

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
			t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldDeleteCustomEventSpecificationWithSystemRuleThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithSystemRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithSystemRuleResourceData(data)
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().DeleteByID(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Delete(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldReturnErrorWhenDeleteCustomEventSpecificationWithSystemRuleFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithSystemRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithSystemRuleResourceData(data)
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockCustomEventSpecificationResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().DeleteByID(gomock.Eq(customEventSpecificationWithThresholdRuleID)).Return(expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Delete(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error to be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
		}
	})
}

func TestShouldFailToDeleteCustomEventSpecificationWithSystemRuleWhenInvalidSeverityIsConfiguredInTerraform(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithSystemRuleData()
		data[CustomEventSpecificationRuleSeverity] = "invalid"
		resourceData := testHelper.CreateCustomEventSpecificationWithSystemRuleResourceData(data)
		resourceData.SetId(customEventSpecificationWithThresholdRuleID)

		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)

		resource := CreateResourceCustomEventSpecificationWithSystemRule()
		err := resource.Delete(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customSystemEventMessageNotAValidSeverity) {
			t.Fatal(customSystemEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func verifyCustomEventSpecificationWithSystemRuleModelAppliedToResource(model restapi.CustomEventSpecification, resourceData *schema.ResourceData, t *testing.T) {
	verifyCustomEventSpecificationModelAppliedToResource(model, resourceData, t)
	verifyCustomEventSpecificationDownstreamModelAppliedToResource(model, resourceData, t)

	convertedSeverity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(model.Rules[0].Severity)
	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if convertedSeverity != resourceData.Get(CustomEventSpecificationRuleSeverity).(string) {
		t.Fatal("Expected Severity to be identical")
	}

	if model.Rules[0].SystemRuleID != resourceData.Get(SystemRuleSpecificationSystemRuleID).(string) {
		t.Fatal("Expected System Rule ID to be identical")
	}
}

func verifyCustomEventSpecificationModelAppliedToResource(model restapi.CustomEventSpecification, resourceData *schema.ResourceData, t *testing.T) {
	if model.ID != resourceData.Id() {
		t.Fatal("Expected ID to be identical")
	}
	if model.Name != resourceData.Get(CustomEventSpecificationFieldFullName).(string) {
		t.Fatal("Expected Full Name to be identical")
	}
	if model.EntityType != resourceData.Get(CustomEventSpecificationFieldEntityType).(string) {
		t.Fatal("Expected EntityType to be identical")
	}
	verifyCustomEventSpecificationQueryAppliedToResource(model, resourceData, t)
	if model.Triggering != resourceData.Get(CustomEventSpecificationFieldTriggering).(bool) {
		t.Fatal("Expected Triggering to be identical")
	}
	verifyCustomEventSpecificationDescriptionAppliedToResource(model, resourceData, t)
	verifyCustomEventSpecificationExpirationTimeAppliedToResource(model, resourceData, t)
	if model.Enabled != resourceData.Get(CustomEventSpecificationFieldEnabled).(bool) {
		t.Fatal("Expected Enabled to be identical")
	}
}

func verifyCustomEventSpecificationQueryAppliedToResource(model restapi.CustomEventSpecification, resourceData *schema.ResourceData, t *testing.T) {
	if model.Query != nil {
		if *model.Query != resourceData.Get(CustomEventSpecificationFieldQuery).(string) {
			t.Fatal("Expected Query to be identical")
		}
	} else {
		if _, ok := resourceData.GetOk(CustomEventSpecificationFieldQuery); ok {
			t.Fatal("Expected Query not to be defined")
		}
	}
}

func verifyCustomEventSpecificationDescriptionAppliedToResource(model restapi.CustomEventSpecification, resourceData *schema.ResourceData, t *testing.T) {
	if model.Description != nil {
		if *model.Description != resourceData.Get(CustomEventSpecificationFieldDescription).(string) {
			t.Fatal("Expected Description to be identical")
		}
	} else {
		if _, ok := resourceData.GetOk(CustomEventSpecificationFieldDescription); ok {
			t.Fatal("Expected Description not to be defined")
		}
	}
}

func verifyCustomEventSpecificationExpirationTimeAppliedToResource(model restapi.CustomEventSpecification, resourceData *schema.ResourceData, t *testing.T) {
	if model.ExpirationTime != nil {
		if *model.ExpirationTime != resourceData.Get(CustomEventSpecificationFieldExpirationTime).(int) {
			t.Fatal("Expected Expiration Time to be identical")
		}
	} else {
		if _, ok := resourceData.GetOk(CustomEventSpecificationFieldExpirationTime); ok {
			t.Fatal("Expected Expiration Time not to be defined")
		}
	}
}

func verifyCustomEventSpecificationDownstreamModelAppliedToResource(model restapi.CustomEventSpecification, resourceData *schema.ResourceData, t *testing.T) {
	if model.Downstream != nil {
		if !cmp.Equal(model.Downstream.IntegrationIds, ReadStringArrayParameterFromResource(resourceData, CustomEventSpecificationDownstreamIntegrationIds)) {
			t.Fatal("Expected Integration IDs to be identical")
		}
		if model.Downstream.BroadcastToAllAlertingConfigs != resourceData.Get(CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs).(bool) {
			t.Fatal("Expected Broadcast to All Alert Configs to be identical")
		}
	} else {
		if _, ok := resourceData.GetOk(CustomEventSpecificationDownstreamIntegrationIds); ok {
			t.Fatal("Expected Integration IDs not to be defined")
		}
		if !resourceData.Get(CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs).(bool) {
			t.Fatalf("Expected Broadcast to All Alert Configs to have the default value set")
		}
	}
}

func createTestCustomEventSpecificationWithSystemRuleModelWithFullDataSet() restapi.CustomEventSpecification {
	description := customSystemEventDescription
	expirationTime := customSystemEventExpirationTime
	query := customSystemEventQuery

	data := createBaseTestCustomEventSpecificationWithSystemRuleModel()
	data.Query = &query
	data.Description = &description
	data.ExpirationTime = &expirationTime
	data.Downstream = &restapi.EventSpecificationDownstream{
		IntegrationIds:                []string{customSystemEventDownStringIntegrationId1, customSystemEventDownStringIntegrationId2},
		BroadcastToAllAlertingConfigs: true,
	}
	return data
}

func createBaseTestCustomEventSpecificationWithSystemRuleModel() restapi.CustomEventSpecification {
	return restapi.CustomEventSpecification{
		ID:         customSystemEventID,
		Name:       customSystemEventName,
		EntityType: customSystemEventEntityType,
		Triggering: false,
		Enabled:    true,
		Rules: []restapi.RuleSpecification{
			restapi.NewSystemRuleSpecification(customSystemEventRuleSystemRuleId, restapi.SeverityWarning.GetAPIRepresentation()),
		},
	}
}

func createFullTestCustomEventSpecificationWithSystemRuleData() map[string]interface{} {
	data := make(map[string]interface{})
	data[CustomEventSpecificationFieldName] = customSystemEventName
	data[CustomEventSpecificationFieldEntityType] = customSystemEventEntityType
	data[CustomEventSpecificationFieldQuery] = customSystemEventQuery
	data[CustomEventSpecificationFieldTriggering] = "true"
	data[CustomEventSpecificationFieldDescription] = customSystemEventDescription
	data[CustomEventSpecificationFieldExpirationTime] = customSystemEventExpirationTime
	data[CustomEventSpecificationFieldEnabled] = "true"
	integrationIds := make([]interface{}, 2)
	integrationIds[0] = customSystemEventDownStringIntegrationId1
	integrationIds[1] = customSystemEventDownStringIntegrationId2
	data[CustomEventSpecificationDownstreamIntegrationIds] = integrationIds
	data[CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs] = "true"
	data[CustomEventSpecificationRuleSeverity] = customSystemEventRuleSeverity
	data[SystemRuleSpecificationSystemRuleID] = customSystemEventRuleSystemRuleId
	return data
}
