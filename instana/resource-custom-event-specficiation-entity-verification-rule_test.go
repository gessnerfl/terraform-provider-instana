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

var testCustomEventSpecificationWithEntityVerificationRuleProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceCustomEventSpecificationWithEntityVerificationRuleDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_custom_event_spec_entity_verification_rule" "example" {
  name = "name"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
  rule_severity = "warning"
  rule_matching_entity_type = "matching-entity-type"
  rule_matching_operator = "is"
  rule_matching_entity_label = "matching-entity-label"
  rule_offline_duration = 60000
  downstream_integration_ids = [ "integration-id-1", "integration-id-2" ]
  downstream_broadcast_to_all_alerting_configs = true
}
`

const (
	customEntityVerificationEventApiPath                             = restapi.CustomEventSpecificationResourcePath + "/{id}"
	testCustomEventSpecificationWithEntityVerificationRuleDefinition = ResourceInstanaCustomEventSpecificationEntityVerificationRule + ".example"

	customEntityVerificationEventID                       = "custom-entity-verification-event-id"
	customEntityVerificationEventName                     = "name"
	customEntityVerificationEventQuery                    = "query"
	customEntityVerificationEventExpirationTime           = 60000
	customEntityVerificationEventDescription              = "description"
	customEntityVerificationEventRuleMatchingEntityLabel  = "matching-entity-label"
	customEntityVerificationEventRuleMatchingEntityType   = "matching-entity-type"
	customEntityVerificationEventRuleMatchingOperator     = restapi.MatchingOperatorIs
	customEntityVerificationEventRuleOfflineDuration      = 60000
	customEntityVerificationEventDownStringIntegrationId1 = "integration-id-1"
	customEntityVerificationEventDownStringIntegrationId2 = "integration-id-2"

	customEntityVerificationEventMessageNotAValidSeverity           = "not a valid severity"
	customEntityVerificationEventTestMessageExpectedInvalidSeverity = "Expected to get error that the provided severity is not valid"

	constEntityVerificationEventContentType = "Content-Type"
)

var customEntityVerificationEventRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCreateResourceCustomEventSpecificationWithEntityVerificationRuleResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, customEntityVerificationEventApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, customEntityVerificationEventApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, customEntityVerificationEventApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"name" : "name",
			"query" : "query",
			"entityType" : "host",
			"enabled" : true,
			"triggering" : true,
			"description" : "description",
			"expirationTime" : 60000,
			"rules" : [ { "ruleType" : "entity_verification", "severity" : 5, "matchingEntityLabel" : "matching-entity-label", "matchingEntityType" : "matching-entity-type", "matchingOperator" : "is", "offlineDuration" : 60000 } ],
			"downstream" : {
				"integrationIds" : ["integration-id-1", "integration-id-2"],
				"broadcastToAllAlertingConfigs" : true
			}
		}
		`, "{{id}}", vars["id"])
		w.Header().Set(constEntityVerificationEventContentType, r.Header.Get(constEntityVerificationEventContentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceCustomEventSpecificationWithEntityVerificationRuleDefinition := strings.ReplaceAll(resourceCustomEventSpecificationWithEntityVerificationRuleDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testCustomEventSpecificationWithEntityVerificationRuleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceCustomEventSpecificationWithEntityVerificationRuleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithEntityVerificationRuleDefinition, "id"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldName, customEntityVerificationEventName),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldQuery, customEntityVerificationEventQuery),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldTriggering, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldDescription, customEntityVerificationEventDescription),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEntityVerificationEventExpirationTime)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldEnabled, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationDownstreamIntegrationIds+".0", customEntityVerificationEventDownStringIntegrationId1),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationDownstreamIntegrationIds+".1", customEntityVerificationEventDownStringIntegrationId2),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingEntityLabel, customEntityVerificationEventRuleMatchingEntityLabel),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingEntityType, customEntityVerificationEventRuleMatchingEntityType),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingOperator, string(customEntityVerificationEventRuleMatchingOperator)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldOfflineDuration, strconv.Itoa(customEntityVerificationEventRuleOfflineDuration)),
				),
			},
		},
	})
}

func TestResourceCustomEventSpecificationWithEntityVerificationRuleDefinition(t *testing.T) {
	resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()

	validateCustomEventSpecificationWithEntityVerificationRuleResourceSchema(resource.Schema, t)

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

func validateCustomEventSpecificationWithEntityVerificationRuleResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldFullName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfStrings(CustomEventSpecificationDownstreamIntegrationIds)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, true)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleSeverity)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingEntityLabel)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingEntityType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeInt(EntityVerificationRuleFieldOfflineDuration)
}

func TestShouldSuccessfullyReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := createBaseTestCustomEventSpecificationWithEntityVerificationRuleModel()
	testShouldSuccessfullyReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPI(expectedModel, t)
}

func TestShouldSuccessfullyReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPIWhenFullDataIsReturned(t *testing.T) {
	expectedModel := createTestCustomEventSpecificationWithEntityVerificationRuleModelWithFullDataSet()
	testShouldSuccessfullyReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPI(expectedModel, t)
}

func testShouldSuccessfullyReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPI(expectedModel restapi.CustomEventSpecification, t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithEntityVerificationRuleResourceData()
		resourceData.SetId(customEntityVerificationEventID)

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEntityVerificationEventID)).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Read(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyCustomEventSpecificationWithEntityVerificationRuleModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldFailToReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPIWhenSeverityFromAPICannotBeMappedToSeverityOfTerraformState(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		expectedModel := createTestCustomEventSpecificationWithEntityVerificationRuleModelWithFullDataSet()
		expectedModel.Rules[0].Severity = 999
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithEntityVerificationRuleResourceData()
		resourceData.SetId(customEntityVerificationEventID)

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEntityVerificationEventID)).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Read(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customEntityVerificationEventMessageNotAValidSeverity) {
			t.Fatal(customEntityVerificationEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithEntityVerificationRuleResourceData()

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Read(resourceData, providerMeta)

		if err == nil || !strings.HasPrefix(err.Error(), "ID of custom event specification") {
			t.Fatal("Expected error to occur because of missing id")
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPIAndDeleteResourceWhenCustomEventDoesNotExist(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithEntityVerificationRuleResourceData()
		resourceData.SetId(customEntityVerificationEventID)

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEntityVerificationEventID)).Return(restapi.CustomEventSpecification{}, restapi.ErrEntityNotFound).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Read(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldFailToReadCustomEventSpecificationWithEntityVerificationRuleFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyCustomEventSpecificationWithEntityVerificationRuleResourceData()
		resourceData.SetId(customEntityVerificationEventID)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockCustomEventAPI.EXPECT().GetOne(gomock.Eq(customEntityVerificationEventID)).Return(restapi.CustomEventSpecification{}, expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Read(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error should be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID should still be set")
		}
	})
}

func TestShouldCreateCustomEventSpecificationWithEntityVerificationRuleThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithEntityVerificationRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithEntityVerificationRuleResourceData(data)
		expectedModel := createTestCustomEventSpecificationWithEntityVerificationRuleModelWithFullDataSet()

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Create(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyCustomEventSpecificationWithEntityVerificationRuleModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithEntityVerificationRuleFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithEntityVerificationRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithEntityVerificationRuleResourceData(data)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(restapi.CustomEventSpecification{}, expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || expectedError != err {
			t.Fatal("Expected definned error to be returned")
		}
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithEntityVerificationRuleFailsBecauseOfInvalidSeverityConfiguredInTerraform(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithEntityVerificationRuleData()
		data[CustomEventSpecificationRuleSeverity] = "invalid"
		resourceData := testHelper.CreateCustomEventSpecificationWithEntityVerificationRuleResourceData(data)

		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customEntityVerificationEventMessageNotAValidSeverity) {
			t.Fatal(customEntityVerificationEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldReturnErrorWhenCreateCustomEventSpecificationWithEntityVerificationRuleFailsBecauseOfInvalidSeverityReturnedFromInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithEntityVerificationRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithEntityVerificationRuleResourceData(data)
		expectedModel := createTestCustomEventSpecificationWithEntityVerificationRuleModelWithFullDataSet()
		expectedModel.Rules[0].Severity = 999

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.CustomEventSpecification{})).Return(expectedModel, nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Create(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customEntityVerificationEventMessageNotAValidSeverity) {
			t.Fatal(customEntityVerificationEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func TestShouldDeleteCustomEventSpecificationWithEntityVerificationRuleThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithEntityVerificationRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithEntityVerificationRuleResourceData(data)
		resourceData.SetId(customEntityVerificationEventID)

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().DeleteByID(gomock.Eq(customEntityVerificationEventID)).Return(nil).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Delete(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldReturnErrorWhenDeleteCustomEventSpecificationWithEntityVerificationRuleFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithEntityVerificationRuleData()
		resourceData := testHelper.CreateCustomEventSpecificationWithEntityVerificationRuleResourceData(data)
		resourceData.SetId(customEntityVerificationEventID)
		expectedError := errors.New("test")

		mockCustomEventAPI := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().CustomEventSpecifications().Return(mockCustomEventAPI).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)
		mockCustomEventAPI.EXPECT().DeleteByID(gomock.Eq(customEntityVerificationEventID)).Return(expectedError).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Delete(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error to be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
		}
	})
}

func TestShouldFailToDeleteCustomEventSpecificationWithEntityVerificationRuleWhenInvalidSeverityIsConfiguredInTerraform(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createFullTestCustomEventSpecificationWithEntityVerificationRuleData()
		data[CustomEventSpecificationRuleSeverity] = "invalid"
		resourceData := testHelper.CreateCustomEventSpecificationWithEntityVerificationRuleResourceData(data)
		resourceData.SetId(customEntityVerificationEventID)

		mockResourceNameFormatter.EXPECT().Format(data[CustomEventSpecificationFieldName]).Return(data[CustomEventSpecificationFieldName]).Times(1)

		resource := CreateResourceCustomEventSpecificationWithEntityVerificationRule()
		err := resource.Delete(resourceData, providerMeta)

		if err == nil || !strings.Contains(err.Error(), customEntityVerificationEventMessageNotAValidSeverity) {
			t.Fatal(customEntityVerificationEventTestMessageExpectedInvalidSeverity)
		}
	})
}

func verifyCustomEventSpecificationWithEntityVerificationRuleModelAppliedToResource(model restapi.CustomEventSpecification, resourceData *schema.ResourceData, t *testing.T) {
	verifyCustomEventSpecificationModelAppliedToResource(model, resourceData, t)
	verifyCustomEventSpecificationDownstreamModelAppliedToResource(model, resourceData, t)

	convertedSeverity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(model.Rules[0].Severity)
	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if convertedSeverity != resourceData.Get(CustomEventSpecificationRuleSeverity).(string) {
		t.Fatal("Expected Severity to be identical")
	}

	if *model.Rules[0].MatchingEntityLabel != resourceData.Get(EntityVerificationRuleFieldMatchingEntityLabel).(string) {
		t.Fatal("Expected EntityVerification MatchingEntityLabel to be identical")
	}
	if *model.Rules[0].MatchingEntityType != resourceData.Get(EntityVerificationRuleFieldMatchingEntityType).(string) {
		t.Fatal("Expected EntityVerification MatchingEntityType to be identical")
	}
	if string(*model.Rules[0].MatchingOperator) != resourceData.Get(EntityVerificationRuleFieldMatchingOperator).(string) {
		t.Fatal("Expected EntityVerification MatchingOperator to be identical")
	}
	if *model.Rules[0].OfflineDuration != resourceData.Get(EntityVerificationRuleFieldOfflineDuration).(int) {
		t.Fatal("Expected EntityVerification OfflineDuration to be identical")
	}
}

func createTestCustomEventSpecificationWithEntityVerificationRuleModelWithFullDataSet() restapi.CustomEventSpecification {
	description := customEntityVerificationEventDescription
	expirationTime := customEntityVerificationEventExpirationTime
	query := customEntityVerificationEventQuery

	data := createBaseTestCustomEventSpecificationWithEntityVerificationRuleModel()
	data.Query = &query
	data.Description = &description
	data.ExpirationTime = &expirationTime
	data.Downstream = &restapi.EventSpecificationDownstream{
		IntegrationIds:                []string{customEntityVerificationEventDownStringIntegrationId1, customEntityVerificationEventDownStringIntegrationId2},
		BroadcastToAllAlertingConfigs: true,
	}
	return data
}

func createBaseTestCustomEventSpecificationWithEntityVerificationRuleModel() restapi.CustomEventSpecification {
	return restapi.CustomEventSpecification{
		ID:         customEntityVerificationEventID,
		Name:       customEntityVerificationEventName,
		EntityType: EntityVerificationRuleEntityType,
		Triggering: false,
		Enabled:    true,
		Rules: []restapi.RuleSpecification{
			restapi.NewEntityVerificationRuleSpecification(customEntityVerificationEventRuleMatchingEntityLabel,
				customEntityVerificationEventRuleMatchingEntityType,
				customEntityVerificationEventRuleMatchingOperator,
				customEntityVerificationEventRuleOfflineDuration,
				restapi.SeverityWarning.GetAPIRepresentation()),
		},
	}
}

func createFullTestCustomEventSpecificationWithEntityVerificationRuleData() map[string]interface{} {
	data := make(map[string]interface{})
	data[CustomEventSpecificationFieldName] = customEntityVerificationEventName
	data[CustomEventSpecificationFieldEntityType] = EntityVerificationRuleEntityType
	data[CustomEventSpecificationFieldQuery] = customEntityVerificationEventQuery
	data[CustomEventSpecificationFieldTriggering] = "true"
	data[CustomEventSpecificationFieldDescription] = customEntityVerificationEventDescription
	data[CustomEventSpecificationFieldExpirationTime] = customEntityVerificationEventExpirationTime
	data[CustomEventSpecificationFieldEnabled] = "true"
	integrationIds := make([]interface{}, 2)
	integrationIds[0] = customEntityVerificationEventDownStringIntegrationId1
	integrationIds[1] = customEntityVerificationEventDownStringIntegrationId2
	data[CustomEventSpecificationDownstreamIntegrationIds] = integrationIds
	data[CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs] = "true"
	data[CustomEventSpecificationRuleSeverity] = customEntityVerificationEventRuleSeverity
	data[EntityVerificationRuleFieldMatchingEntityLabel] = customEntityVerificationEventRuleMatchingEntityLabel
	data[EntityVerificationRuleFieldMatchingEntityType] = customEntityVerificationEventRuleMatchingEntityType
	data[EntityVerificationRuleFieldMatchingOperator] = string(customEntityVerificationEventRuleMatchingOperator)
	return data
}
